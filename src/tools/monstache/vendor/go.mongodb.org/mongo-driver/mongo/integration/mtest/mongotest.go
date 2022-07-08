// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mtest

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/testutil/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

var (
	// Background is a no-op context.
	Background = context.Background()
	// MajorityWc is the majority write concern.
	MajorityWc = writeconcern.New(writeconcern.WMajority())
	// PrimaryRp is the primary read preference.
	PrimaryRp = readpref.Primary()
	// LocalRc is the local read concern
	LocalRc = readconcern.Local()
	// MajorityRc is the majority read concern
	MajorityRc = readconcern.Majority()
)

const (
	namespaceExistsErrCode int32 = 48
)

// FailPoint is a representation of a server fail point.
// See https://github.com/mongodb/specifications/tree/master/source/transactions/tests#server-fail-point
// for more information regarding fail points.
type FailPoint struct {
	ConfigureFailPoint string `bson:"configureFailPoint"`
	// Mode should be a string, FailPointMode, or map[string]interface{}
	Mode interface{}   `bson:"mode"`
	Data FailPointData `bson:"data"`
}

// FailPointMode is a representation of the Failpoint.Mode field.
type FailPointMode struct {
	Times int32 `bson:"times"`
	Skip  int32 `bson:"skip"`
}

// FailPointData is a representation of the FailPoint.Data field.
type FailPointData struct {
	FailCommands                  []string               `bson:"failCommands,omitempty"`
	CloseConnection               bool                   `bson:"closeConnection,omitempty"`
	ErrorCode                     int32                  `bson:"errorCode,omitempty"`
	FailBeforeCommitExceptionCode int32                  `bson:"failBeforeCommitExceptionCode,omitempty"`
	ErrorLabels                   *[]string              `bson:"errorLabels,omitempty"`
	WriteConcernError             *WriteConcernErrorData `bson:"writeConcernError,omitempty"`
	BlockConnection               bool                   `bson:"blockConnection,omitempty"`
	BlockTimeMS                   int32                  `bson:"blockTimeMS,omitempty"`
	AppName                       string                 `bson:"appName,omitempty"`
}

// WriteConcernErrorData is a representation of the FailPoint.Data.WriteConcern field.
type WriteConcernErrorData struct {
	Code        int32     `bson:"code"`
	Name        string    `bson:"codeName"`
	Errmsg      string    `bson:"errmsg"`
	ErrorLabels *[]string `bson:"errorLabels,omitempty"`
	ErrInfo     bson.Raw  `bson:"errInfo,omitempty"`
}

// T is a wrapper around testing.T.
type T struct {
	*testing.T

	// members for only this T instance
	createClient      *bool
	createCollection  *bool
	runOn             []RunOnBlock
	mockDeployment    *mockDeployment // nil if the test is not being run against a mock
	mockResponses     []bson.D
	createdColls      []*Collection // collections created in this test
	proxyDialer       *proxyDialer
	dbName, collName  string
	failPointNames    []string
	minServerVersion  string
	maxServerVersion  string
	validTopologies   []TopologyKind
	auth              *bool
	enterprise        *bool
	dataLake          *bool
	ssl               *bool
	collCreateOpts    bson.D
	connsCheckedOut   int // net number of connections checked out during test execution
	requireAPIVersion *bool

	// options copied to sub-tests
	clientType  ClientType
	clientOpts  *options.ClientOptions
	collOpts    *options.CollectionOptions
	shareClient *bool

	baseOpts *Options // used to create subtests

	// command monitoring channels
	monitorLock sync.Mutex
	started     []*event.CommandStartedEvent
	succeeded   []*event.CommandSucceededEvent
	failed      []*event.CommandFailedEvent

	Client *mongo.Client
	DB     *mongo.Database
	Coll   *mongo.Collection
}

func newT(wrapped *testing.T, opts ...*Options) *T {
	t := &T{
		T: wrapped,
	}
	for _, opt := range opts {
		for _, optFn := range opt.optFuncs {
			optFn(t)
		}
	}

	if err := t.verifyConstraints(); err != nil {
		t.Skipf("skipping due to environmental constraints: %v", err)
	}

	if t.collName == "" {
		t.collName = t.Name()
	}
	if t.dbName == "" {
		t.dbName = TestDb
	}
	t.collName = sanitizeCollectionName(t.dbName, t.collName)

	// create a set of base options for sub-tests
	t.baseOpts = NewOptions().ClientOptions(t.clientOpts).CollectionOptions(t.collOpts).ClientType(t.clientType)
	if t.shareClient != nil {
		t.baseOpts.ShareClient(*t.shareClient)
	}

	return t
}

// New creates a new T instance with the given options. If the current environment does not satisfy constraints
// specified in the options, the test will be skipped automatically.
func New(wrapped *testing.T, opts ...*Options) *T {
	t := newT(wrapped, opts...)

	// only create a client if it needs to be shared in sub-tests
	// otherwise, a new client will be created for each subtest
	if t.shareClient != nil && *t.shareClient {
		t.createTestClient()
	}

	return t
}

// Close cleans up any resources associated with a T. There should be one Close corresponding to every New.
func (t *T) Close() {
	if t.Client == nil {
		return
	}

	// only clear collections and fail points if the test is not running against a mock
	if t.clientType != Mock {
		t.ClearCollections()
		t.ClearFailPoints()
	}

	// always disconnect the client regardless of clientType because Client.Disconnect will work against
	// all deployments
	_ = t.Client.Disconnect(Background)
}

// Run creates a new T instance for a sub-test and runs the given callback. It also creates a new collection using the
// given name which is available to the callback through the T.Coll variable and is dropped after the callback
// returns.
func (t *T) Run(name string, callback func(*T)) {
	t.RunOpts(name, NewOptions(), callback)
}

// RunOpts creates a new T instance for a sub-test with the given options. If the current environment does not satisfy
// constraints specified in the options, the new sub-test will be skipped automatically. If the test is not skipped,
// the callback will be run with the new T instance. RunOpts creates a new collection with the given name which is
// available to the callback through the T.Coll variable and is dropped after the callback returns.
func (t *T) RunOpts(name string, opts *Options, callback func(*T)) {
	t.T.Run(name, func(wrapped *testing.T) {
		sub := newT(wrapped, t.baseOpts, opts)

		// add any mock responses for this test
		if sub.clientType == Mock && len(sub.mockResponses) > 0 {
			sub.AddMockResponses(sub.mockResponses...)
		}

		// for shareClient, inherit the client from the parent
		if sub.shareClient != nil && *sub.shareClient && sub.clientType == t.clientType {
			sub.Client = t.Client
		}
		// only create a client if not already set
		if sub.Client == nil {
			if sub.createClient == nil || *sub.createClient {
				sub.createTestClient()
			}
		}
		// create a collection for this test
		if sub.Client != nil {
			sub.createTestCollection()
		}

		// defer dropping all collections if the test is using a client
		defer func() {
			if sub.Client == nil {
				return
			}

			// store number of sessions and connections checked out here but assert that they're equal to 0 after
			// cleaning up test resources to make sure resources are always cleared
			sessions := sub.Client.NumberSessionsInProgress()
			conns := sub.connsCheckedOut

			if sub.clientType != Mock {
				sub.ClearFailPoints()
				sub.ClearCollections()
			}
			// only disconnect client if it's not being shared
			if sub.shareClient == nil || !*sub.shareClient {
				_ = sub.Client.Disconnect(Background)
			}
			assert.Equal(sub, 0, sessions, "%v sessions checked out", sessions)
			assert.Equal(sub, 0, conns, "%v connections checked out", conns)
		}()

		// clear any events that may have happened during setup and run the test
		sub.ClearEvents()
		callback(sub)
	})
}

// AddMockResponses adds responses to be returned by the mock deployment. This should only be used if T is being run
// against a mock deployment.
func (t *T) AddMockResponses(responses ...bson.D) {
	t.mockDeployment.addResponses(responses...)
}

// ClearMockResponses clears all responses in the mock deployment.
func (t *T) ClearMockResponses() {
	t.mockDeployment.clearResponses()
}

// GetStartedEvent returns the most recent CommandStartedEvent, or nil if one is not present.
// This can only be called once per event.
// TODO(GODRIVER-2075): GetStartedEvent documents that it returns the most recent event, but
// actually returns the first event. Update either the documentation or implementation.
func (t *T) GetStartedEvent() *event.CommandStartedEvent {
	if len(t.started) == 0 {
		return nil
	}
	e := t.started[0]
	t.started = t.started[1:]
	return e
}

// GetSucceededEvent returns the most recent CommandSucceededEvent, or nil if one is not present.
// This can only be called once per event.
// TODO(GODRIVER-2075): GetSucceededEvent documents that it returns the most recent event, but
// actually returns the first event. Update either the documentation or implementation.
func (t *T) GetSucceededEvent() *event.CommandSucceededEvent {
	if len(t.succeeded) == 0 {
		return nil
	}
	e := t.succeeded[0]
	t.succeeded = t.succeeded[1:]
	return e
}

// GetFailedEvent returns the most recent CommandFailedEvent, or nil if one is not present.
// This can only be called once per event.
// TODO(GODRIVER-2075): GetFailedEvent documents that it returns the most recent event, but actually
// returns the first event. Update either the documentation or implementation.
func (t *T) GetFailedEvent() *event.CommandFailedEvent {
	if len(t.failed) == 0 {
		return nil
	}
	e := t.failed[0]
	t.failed = t.failed[1:]
	return e
}

// GetAllStartedEvents returns a slice of all CommandStartedEvent instances for this test. This can be called multiple
// times.
func (t *T) GetAllStartedEvents() []*event.CommandStartedEvent {
	return t.started
}

// GetAllSucceededEvents returns a slice of all CommandSucceededEvent instances for this test. This can be called multiple
// times.
func (t *T) GetAllSucceededEvents() []*event.CommandSucceededEvent {
	return t.succeeded
}

// GetAllFailedEvents returns a slice of all CommandFailedEvent instances for this test. This can be called multiple
// times.
func (t *T) GetAllFailedEvents() []*event.CommandFailedEvent {
	return t.failed
}

// FilterStartedEvents filters the existing CommandStartedEvent instances for this test using the provided filter
// callback. An event will be retained if the filter returns true. The list of filtered events will be used to overwrite
// the list of events for this test and will therefore change the output of t.GetAllStartedEvents().
func (t *T) FilterStartedEvents(filter func(*event.CommandStartedEvent) bool) {
	var newEvents []*event.CommandStartedEvent
	for _, evt := range t.started {
		if filter(evt) {
			newEvents = append(newEvents, evt)
		}
	}
	t.started = newEvents
}

// FilterSucceededEvents filters the existing CommandSucceededEvent instances for this test using the provided filter
// callback. An event will be retained if the filter returns true. The list of filtered events will be used to overwrite
// the list of events for this test and will therefore change the output of t.GetAllSucceededEvents().
func (t *T) FilterSucceededEvents(filter func(*event.CommandSucceededEvent) bool) {
	var newEvents []*event.CommandSucceededEvent
	for _, evt := range t.succeeded {
		if filter(evt) {
			newEvents = append(newEvents, evt)
		}
	}
	t.succeeded = newEvents
}

// FilterFailedEvents filters the existing CommandFailedEVent instances for this test using the provided filter
// callback. An event will be retained if the filter returns true. The list of filtered events will be used to overwrite
// the list of events for this test and will therefore change the output of t.GetAllFailedEvents().
func (t *T) FilterFailedEvents(filter func(*event.CommandFailedEvent) bool) {
	var newEvents []*event.CommandFailedEvent
	for _, evt := range t.failed {
		if filter(evt) {
			newEvents = append(newEvents, evt)
		}
	}
	t.failed = newEvents
}

// GetProxiedMessages returns the messages proxied to the server by the test. If the client type is not Proxy, this
// returns nil.
func (t *T) GetProxiedMessages() []*ProxyMessage {
	if t.proxyDialer == nil {
		return nil
	}
	return t.proxyDialer.Messages()
}

// NumberConnectionsCheckedOut returns the number of connections checked out from the test Client.
func (t *T) NumberConnectionsCheckedOut() int {
	return t.connsCheckedOut
}

// ClearEvents clears the existing command monitoring events.
func (t *T) ClearEvents() {
	t.started = t.started[:0]
	t.succeeded = t.succeeded[:0]
	t.failed = t.failed[:0]
}

// ResetClient resets the existing client with the given options. If opts is nil, the existing options will be used.
// If t.Coll is not-nil, it will be reset to use the new client. Should only be called if the existing client is
// not nil. This will Disconnect the existing client but will not drop existing collections. To do so, ClearCollections
// must be called before calling ResetClient.
func (t *T) ResetClient(opts *options.ClientOptions) {
	if opts != nil {
		t.clientOpts = opts
	}

	_ = t.Client.Disconnect(Background)
	t.createTestClient()
	t.DB = t.Client.Database(t.dbName)
	t.Coll = t.DB.Collection(t.collName, t.collOpts)

	for _, coll := range t.createdColls {
		// If the collection was created using a different Client, it doesn't need to be reset.
		if coll.hasDifferentClient {
			continue
		}

		// If the namespace is the same as t.Coll, we can use t.Coll.
		if coll.created.Name() == t.collName && coll.created.Database().Name() == t.dbName {
			coll.created = t.Coll
			continue
		}

		// Otherwise, reset the collection to use the new Client.
		coll.created = t.Client.Database(coll.DB).Collection(coll.Name, coll.Opts)
	}
}

// Collection is used to configure a new collection created during a test.
type Collection struct {
	Name       string
	DB         string        // defaults to mt.DB.Name() if not specified
	Client     *mongo.Client // defaults to mt.Client if not specified
	Opts       *options.CollectionOptions
	CreateOpts bson.D

	hasDifferentClient bool
	created            *mongo.Collection // the actual collection that was created
}

// CreateCollection creates a new collection with the given configuration. The collection will be dropped after the test
// finishes running. If createOnServer is true, the function ensures that the collection has been created server-side
// by running the create command. The create command will appear in command monitoring channels.
func (t *T) CreateCollection(coll Collection, createOnServer bool) *mongo.Collection {
	if coll.DB == "" {
		coll.DB = t.DB.Name()
	}
	if coll.Client == nil {
		coll.Client = t.Client
	}
	coll.hasDifferentClient = coll.Client != t.Client

	db := coll.Client.Database(coll.DB)

	if createOnServer && t.clientType != Mock {
		cmd := bson.D{{"create", coll.Name}}
		cmd = append(cmd, coll.CreateOpts...)

		if err := db.RunCommand(Background, cmd).Err(); err != nil {
			// ignore NamespaceExists errors for idempotency

			cmdErr, ok := err.(mongo.CommandError)
			if !ok || cmdErr.Code != namespaceExistsErrCode {
				t.Fatalf("error creating collection %v on server: %v", coll.Name, err)
			}
		}
	}

	coll.created = db.Collection(coll.Name, coll.Opts)
	t.createdColls = append(t.createdColls, &coll)
	return coll.created
}

// ClearCollections drops all collections previously created by this test.
func (t *T) ClearCollections() {
	// Collections should not be dropped when testing against Atlas Data Lake because the data is pre-inserted.
	if !testContext.dataLake {
		for _, coll := range t.createdColls {
			_ = coll.created.Drop(Background)
		}
	}
	t.createdColls = t.createdColls[:0]
}

// SetFailPoint sets a fail point for the client associated with T. Commands to create the failpoint will appear
// in command monitoring channels. The fail point will automatically be disabled after this test has run.
func (t *T) SetFailPoint(fp FailPoint) {
	// ensure mode fields are int32
	if modeMap, ok := fp.Mode.(map[string]interface{}); ok {
		var key string
		var err error

		if times, ok := modeMap["times"]; ok {
			key = "times"
			modeMap["times"], err = t.interfaceToInt32(times)
		}
		if skip, ok := modeMap["skip"]; ok {
			key = "skip"
			modeMap["skip"], err = t.interfaceToInt32(skip)
		}

		if err != nil {
			t.Fatalf("error converting %s to int32: %v", key, err)
		}
	}

	if err := SetFailPoint(fp, t.Client); err != nil {
		t.Fatal(err)
	}
	t.failPointNames = append(t.failPointNames, fp.ConfigureFailPoint)
}

// SetFailPointFromDocument sets the fail point represented by the given document for the client associated with T. This
// method assumes that the given document is in the form {configureFailPoint: <failPointName>, ...}. Commands to create
// the failpoint will appear in command monitoring channels. The fail point will be automatically disabled after this
// test has run.
func (t *T) SetFailPointFromDocument(fp bson.Raw) {
	if err := SetRawFailPoint(fp, t.Client); err != nil {
		t.Fatal(err)
	}

	name := fp.Index(0).Value().StringValue()
	t.failPointNames = append(t.failPointNames, name)
}

// TrackFailPoint adds the given fail point to the list of fail points to be disabled when the current test finishes.
// This function does not create a fail point on the server.
func (t *T) TrackFailPoint(fpName string) {
	t.failPointNames = append(t.failPointNames, fpName)
}

// ClearFailPoints disables all previously set failpoints for this test.
func (t *T) ClearFailPoints() {
	db := t.Client.Database("admin")
	for _, fp := range t.failPointNames {
		cmd := bson.D{
			{"configureFailPoint", fp},
			{"mode", "off"},
		}
		err := db.RunCommand(Background, cmd).Err()
		if err != nil {
			t.Fatalf("error clearing fail point %s: %v", fp, err)
		}
	}
	t.failPointNames = t.failPointNames[:0]
}

// CloneDatabase modifies the default database for this test to match the given options.
func (t *T) CloneDatabase(opts *options.DatabaseOptions) {
	t.DB = t.Client.Database(t.dbName, opts)
}

// CloneCollection modifies the default collection for this test to match the given options.
func (t *T) CloneCollection(opts *options.CollectionOptions) {
	var err error
	t.Coll, err = t.Coll.Clone(opts)
	assert.Nil(t, err, "error cloning collection: %v", err)
}

func sanitizeCollectionName(db string, coll string) string {
	// Collections can't have "$" in their names, so we substitute it with "%".
	coll = strings.Replace(coll, "$", "%", -1)

	// Namespaces can only have 120 bytes max.
	if len(db+"."+coll) >= 120 {
		// coll len must be <= remaining
		remaining := 120 - (len(db) + 1) // +1 for "."
		coll = coll[len(coll)-remaining:]
	}
	return coll
}

func (t *T) createTestClient() {
	clientOpts := t.clientOpts
	if clientOpts == nil {
		// default opts
		clientOpts = options.Client().SetWriteConcern(MajorityWc).SetReadPreference(PrimaryRp)
	}
	// set ServerAPIOptions to latest version if required
	if clientOpts.Deployment == nil && t.clientType != Mock && clientOpts.ServerAPIOptions == nil && testContext.requireAPIVersion {
		clientOpts.SetServerAPIOptions(options.ServerAPI(driver.TestServerAPIVersion))
	}
	// command monitor
	clientOpts.SetMonitor(&event.CommandMonitor{
		Started: func(_ context.Context, cse *event.CommandStartedEvent) {
			t.monitorLock.Lock()
			defer t.monitorLock.Unlock()
			t.started = append(t.started, cse)
		},
		Succeeded: func(_ context.Context, cse *event.CommandSucceededEvent) {
			t.monitorLock.Lock()
			defer t.monitorLock.Unlock()
			t.succeeded = append(t.succeeded, cse)
		},
		Failed: func(_ context.Context, cfe *event.CommandFailedEvent) {
			t.monitorLock.Lock()
			defer t.monitorLock.Unlock()
			t.failed = append(t.failed, cfe)
		},
	})
	// only specify connection pool monitor if no deployment is given
	if clientOpts.Deployment == nil {
		previousPoolMonitor := clientOpts.PoolMonitor

		clientOpts.SetPoolMonitor(&event.PoolMonitor{
			Event: func(evt *event.PoolEvent) {
				if previousPoolMonitor != nil {
					previousPoolMonitor.Event(evt)
				}

				switch evt.Type {
				case event.GetSucceeded:
					t.connsCheckedOut++
				case event.ConnectionReturned:
					t.connsCheckedOut--
				}
			},
		})
	}

	var err error
	switch t.clientType {
	case Pinned:
		// pin to first mongos
		pinnedHostList := []string{testContext.connString.Hosts[0]}
		uriOpts := options.Client().ApplyURI(testContext.connString.Original).SetHosts(pinnedHostList)
		t.Client, err = mongo.NewClient(uriOpts, clientOpts)
	case Mock:
		// clear pool monitor to avoid configuration error
		clientOpts.PoolMonitor = nil
		t.mockDeployment = newMockDeployment()
		clientOpts.Deployment = t.mockDeployment
		t.Client, err = mongo.NewClient(clientOpts)
	case Proxy:
		t.proxyDialer = newProxyDialer()
		clientOpts.SetDialer(t.proxyDialer)

		// After setting the Dialer, fall-through to the Default case to apply the correct URI
		fallthrough
	case Default:
		// Use a different set of options to specify the URI because clientOpts may already have a URI or host seedlist
		// specified.
		var uriOpts *options.ClientOptions
		if clientOpts.Deployment == nil {
			// Only specify URI if the deployment is not set to avoid setting topology/server options along with the
			// deployment.
			uriOpts = options.Client().ApplyURI(testContext.connString.Original)
		}

		// Pass in uriOpts first so clientOpts wins if there are any conflicting settings.
		t.Client, err = mongo.NewClient(uriOpts, clientOpts)
	}
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	if err := t.Client.Connect(Background); err != nil {
		t.Fatalf("error connecting client: %v", err)
	}
}

func (t *T) createTestCollection() {
	t.DB = t.Client.Database(t.dbName)
	t.createdColls = t.createdColls[:0]

	// Collections should not be explicitly created when testing against Atlas Data Lake because they already exist in
	// the server with pre-seeded data.
	createOnServer := (t.createCollection == nil || *t.createCollection) && !testContext.dataLake
	t.Coll = t.CreateCollection(Collection{
		Name:       t.collName,
		CreateOpts: t.collCreateOpts,
		Opts:       t.collOpts,
	}, createOnServer)
}

// verifyVersionConstraints returns an error if the cluster's server version is not in the range [min, max]. Server
// versions will only be checked if they are non-empty.
func verifyVersionConstraints(min, max string) error {
	if min != "" && CompareServerVersions(testContext.serverVersion, min) < 0 {
		return fmt.Errorf("server version %q is lower than min required version %q", testContext.serverVersion, min)
	}
	if max != "" && CompareServerVersions(testContext.serverVersion, max) > 0 {
		return fmt.Errorf("server version %q is higher than max version %q", testContext.serverVersion, max)
	}
	return nil
}

// verifyTopologyConstraints returns an error if the cluster's topology kind does not match one of the provided
// kinds. If the topologies slice is empty, nil is returned without any additional checks.
func verifyTopologyConstraints(topologies []TopologyKind) error {
	if len(topologies) == 0 {
		return nil
	}

	for _, topo := range topologies {
		// For ShardedReplicaSet, we won't get an exact match because testContext.topoKind will be Sharded so we do an
		// additional comparison with the testContext.shardedReplicaSet field.
		if topo == testContext.topoKind || (topo == ShardedReplicaSet && testContext.shardedReplicaSet) {
			return nil
		}
	}
	return fmt.Errorf("topology kind %q does not match any of the required kinds %q", testContext.topoKind, topologies)
}

func verifyServerParametersConstraints(serverParameters map[string]bson.RawValue) error {
	for param, expected := range serverParameters {
		actual, err := testContext.serverParameters.LookupErr(param)
		if err != nil {
			return fmt.Errorf("server does not support parameter %q", param)
		}
		if !expected.Equal(actual) {
			return fmt.Errorf("mismatched values for server parameter %q; expected %s, got %s", param, expected, actual)
		}
	}
	return nil
}

func verifyAuthConstraint(expected *bool) error {
	if expected != nil && *expected != testContext.authEnabled {
		return fmt.Errorf("test requires auth value: %v, cluster auth value: %v", *expected, testContext.authEnabled)
	}
	return nil
}

// verifyRunOnBlockConstraint returns an error if the current environment does not match the provided RunOnBlock.
func verifyRunOnBlockConstraint(rob RunOnBlock) error {
	if err := verifyVersionConstraints(rob.MinServerVersion, rob.MaxServerVersion); err != nil {
		return err
	}
	if err := verifyTopologyConstraints(rob.Topology); err != nil {
		return err
	}
	if err := verifyAuthConstraint(rob.Auth); err != nil {
		return err
	}

	return verifyServerParametersConstraints(rob.ServerParameters)
}

// verifyConstraints returns an error if the current environment does not match the constraints specified for the test.
func (t *T) verifyConstraints() error {
	// Check constraints not specified as runOn blocks
	if err := verifyVersionConstraints(t.minServerVersion, t.maxServerVersion); err != nil {
		return err
	}
	if err := verifyTopologyConstraints(t.validTopologies); err != nil {
		return err
	}
	if err := verifyAuthConstraint(t.auth); err != nil {
		return err
	}
	if t.ssl != nil && *t.ssl != testContext.sslEnabled {
		return fmt.Errorf("test requires ssl value: %v, cluster ssl value: %v", *t.ssl, testContext.sslEnabled)
	}
	if t.enterprise != nil && *t.enterprise != testContext.enterpriseServer {
		return fmt.Errorf("test requires enterprise value: %v, cluster enterprise value: %v", *t.enterprise,
			testContext.enterpriseServer)
	}
	if t.dataLake != nil && *t.dataLake != testContext.dataLake {
		return fmt.Errorf("test requires cluster to be data lake: %v, cluster is data lake: %v", *t.dataLake,
			testContext.dataLake)
	}
	if t.requireAPIVersion != nil && *t.requireAPIVersion != testContext.requireAPIVersion {
		return fmt.Errorf("test requires RequireAPIVersion value: %v, local RequireAPIVersion value: %v", *t.requireAPIVersion,
			testContext.requireAPIVersion)
	}

	// Check runOn blocks. The test can be executed if there are no blocks or at least block matches the current test
	// setup.
	if len(t.runOn) == 0 {
		return nil
	}

	// Stop once we find a RunOnBlock that matches the current environment. Record all errors as we go because if we
	// don't find any matching blocks, we want to report the comparison errors for each block.
	var runOnErrors []error
	for _, runOn := range t.runOn {
		err := verifyRunOnBlockConstraint(runOn)
		if err == nil {
			return nil
		}

		runOnErrors = append(runOnErrors, err)
	}
	return fmt.Errorf("no matching RunOnBlock; comparison errors: %v", runOnErrors)
}

func (t *T) interfaceToInt32(i interface{}) (int32, error) {
	switch conv := i.(type) {
	case int:
		return int32(conv), nil
	case int32:
		return conv, nil
	case int64:
		return int32(conv), nil
	case float64:
		return int32(conv), nil
	}

	return 0, fmt.Errorf("type %T cannot be converted to int32", i)
}
