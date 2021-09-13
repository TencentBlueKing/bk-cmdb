// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package integration

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal/testutil/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var (
	awsAccessKeyID         = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey     = os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsTempAccessKeyID     = os.Getenv("CSFLE_AWS_TEMP_ACCESS_KEY_ID")
	awsTempSecretAccessKey = os.Getenv("CSFLE_AWS_TEMP_SECRET_ACCESS_KEY")
	awsTempSessionToken    = os.Getenv("CSFLE_AWS_TEMP_SESSION_TOKEN")
	azureTenantID          = os.Getenv("AZURE_TENANT_ID")
	azureClientID          = os.Getenv("AZURE_CLIENT_ID")
	azureClientSecret      = os.Getenv("AZURE_CLIENT_SECRET")
	gcpEmail               = os.Getenv("GCP_EMAIL")
	gcpPrivateKey          = os.Getenv("GCP_PRIVATE_KEY")
)

// Helper functions to do read JSON spec test files and convert JSON objects into the appropriate driver types.
// Functions in this file should take testing.TB rather than testing.T/mtest.T for generality because they
// do not do any database communication.

// generate a slice of all JSON file names in a directory
func jsonFilesInDir(t testing.TB, dir string) []string {
	t.Helper()

	files := make([]string, 0)

	entries, err := ioutil.ReadDir(dir)
	assert.Nil(t, err, "unable to read json file: %v", err)

	for _, entry := range entries {
		if entry.IsDir() || path.Ext(entry.Name()) != ".json" {
			continue
		}

		files = append(files, entry.Name())
	}

	return files
}

// create client options from a map
func createClientOptions(t testing.TB, opts bson.Raw) *options.ClientOptions {
	t.Helper()

	clientOpts := options.Client()
	elems, _ := opts.Elements()
	for _, elem := range elems {
		name := elem.Key()
		opt := elem.Value()

		switch name {
		case "retryWrites":
			clientOpts.SetRetryWrites(opt.Boolean())
		case "w":
			switch opt.Type {
			case bson.TypeInt32:
				w := int(opt.Int32())
				clientOpts.SetWriteConcern(writeconcern.New(writeconcern.W(w)))
			case bson.TypeDouble:
				w := int(opt.Double())
				clientOpts.SetWriteConcern(writeconcern.New(writeconcern.W(w)))
			case bson.TypeString:
				clientOpts.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
			default:
				t.Fatalf("unrecognized type for w client option: %v", opt.Type)
			}
		case "readConcernLevel":
			clientOpts.SetReadConcern(readconcern.New(readconcern.Level(opt.StringValue())))
		case "readPreference":
			clientOpts.SetReadPreference(readPrefFromString(opt.StringValue()))
		case "heartbeatFrequencyMS":
			hf := convertValueToMilliseconds(t, opt)
			clientOpts.SetHeartbeatInterval(hf)
		case "retryReads":
			clientOpts.SetRetryReads(opt.Boolean())
		case "autoEncryptOpts":
			clientOpts.SetAutoEncryptionOptions(createAutoEncryptionOptions(t, opt.Document()))
		case "appname":
			clientOpts.SetAppName(opt.StringValue())
		case "connectTimeoutMS":
			ct := convertValueToMilliseconds(t, opt)
			clientOpts.SetConnectTimeout(ct)
		case "serverSelectionTimeoutMS":
			sst := convertValueToMilliseconds(t, opt)
			clientOpts.SetServerSelectionTimeout(sst)
		default:
			t.Fatalf("unrecognized client option: %v", name)
		}
	}

	return clientOpts
}

func createAutoEncryptionOptions(t testing.TB, opts bson.Raw) *options.AutoEncryptionOptions {
	t.Helper()

	aeo := options.AutoEncryption()
	var kvnsFound bool
	elems, _ := opts.Elements()

	for _, elem := range elems {
		name := elem.Key()
		opt := elem.Value()

		switch name {
		case "kmsProviders":
			aeo.SetKmsProviders(createKmsProvidersMap(t, opt.Document()))
		case "schemaMap":
			var schemaMap map[string]interface{}
			err := bson.Unmarshal(opt.Document(), &schemaMap)
			if err != nil {
				t.Fatalf("error creating schema map: %v", err)
			}

			aeo.SetSchemaMap(schemaMap)
		case "keyVaultNamespace":
			kvnsFound = true
			aeo.SetKeyVaultNamespace(opt.StringValue())
		case "bypassAutoEncryption":
			aeo.SetBypassAutoEncryption(opt.Boolean())
		default:
			t.Fatalf("unrecognized auto encryption option: %v", name)
		}
	}
	if !kvnsFound {
		aeo.SetKeyVaultNamespace("keyvault.datakeys")
	}

	return aeo
}

func createKmsProvidersMap(t testing.TB, opts bson.Raw) map[string]map[string]interface{} {
	t.Helper()

	kmsMap := make(map[string]map[string]interface{})
	elems, _ := opts.Elements()

	for _, elem := range elems {
		provider := elem.Key()
		providerOpt := elem.Value()

		switch provider {
		case "aws":
			awsMap := map[string]interface{}{
				"accessKeyId":     awsAccessKeyID,
				"secretAccessKey": awsSecretAccessKey,
			}
			kmsMap["aws"] = awsMap
		case "azure":
			kmsMap["azure"] = map[string]interface{}{
				"tenantId":     azureTenantID,
				"clientId":     azureClientID,
				"clientSecret": azureClientSecret,
			}
		case "gcp":
			kmsMap["gcp"] = map[string]interface{}{
				"email":      gcpEmail,
				"privateKey": gcpPrivateKey,
			}
		case "local":
			_, key := providerOpt.Document().Lookup("key").Binary()
			localMap := map[string]interface{}{
				"key": key,
			}
			kmsMap["local"] = localMap
		case "awsTemporary":
			if awsTempAccessKeyID == "" {
				t.Fatal("AWS temp access key ID not set")
			}
			if awsTempSecretAccessKey == "" {
				t.Fatal("AWS temp secret access key not set")
			}
			if awsTempSessionToken == "" {
				t.Fatal("AWS temp session token not set")
			}
			awsMap := map[string]interface{}{
				"accessKeyId":     awsTempAccessKeyID,
				"secretAccessKey": awsTempSecretAccessKey,
				"sessionToken":    awsTempSessionToken,
			}
			kmsMap["aws"] = awsMap
		case "awsTemporaryNoSessionToken":
			if awsTempAccessKeyID == "" {
				t.Fatal("AWS temp access key ID not set")
			}
			if awsTempSecretAccessKey == "" {
				t.Fatal("AWS temp secret access key not set")
			}
			awsMap := map[string]interface{}{
				"accessKeyId":     awsTempAccessKeyID,
				"secretAccessKey": awsTempSecretAccessKey,
			}
			kmsMap["aws"] = awsMap
		default:
			t.Fatalf("unrecognized KMS provider: %v", provider)
		}
	}

	return kmsMap
}

// create session options from a map
func createSessionOptions(t testing.TB, opts bson.Raw) *options.SessionOptions {
	t.Helper()

	sessOpts := options.Session()
	elems, _ := opts.Elements()
	for _, elem := range elems {
		name := elem.Key()
		opt := elem.Value()

		switch name {
		case "causalConsistency":
			sessOpts = sessOpts.SetCausalConsistency(opt.Boolean())
		case "defaultTransactionOptions":
			txnOpts := createTransactionOptions(t, opt.Document())
			if txnOpts.ReadConcern != nil {
				sessOpts.SetDefaultReadConcern(txnOpts.ReadConcern)
			}
			if txnOpts.ReadPreference != nil {
				sessOpts.SetDefaultReadPreference(txnOpts.ReadPreference)
			}
			if txnOpts.WriteConcern != nil {
				sessOpts.SetDefaultWriteConcern(txnOpts.WriteConcern)
			}
			if txnOpts.MaxCommitTime != nil {
				sessOpts.SetDefaultMaxCommitTime(txnOpts.MaxCommitTime)
			}
		default:
			t.Fatalf("unrecognized session option: %v", name)
		}
	}

	return sessOpts
}

// create database options from a BSON document.
func createDatabaseOptions(t testing.TB, opts bson.Raw) *options.DatabaseOptions {
	t.Helper()

	do := options.Database()
	elems, _ := opts.Elements()
	for _, elem := range elems {
		name := elem.Key()
		opt := elem.Value()

		switch name {
		case "readConcern":
			do.SetReadConcern(createReadConcern(opt))
		case "writeConcern":
			do.SetWriteConcern(createWriteConcern(t, opt))
		default:
			t.Fatalf("unrecognized database option: %v", name)
		}
	}

	return do
}

// create collection options from a map
func createCollectionOptions(t testing.TB, opts bson.Raw) *options.CollectionOptions {
	t.Helper()

	co := options.Collection()
	elems, _ := opts.Elements()
	for _, elem := range elems {
		name := elem.Key()
		opt := elem.Value()

		switch name {
		case "readConcern":
			co.SetReadConcern(createReadConcern(opt))
		case "writeConcern":
			co.SetWriteConcern(createWriteConcern(t, opt))
		case "readPreference":
			co.SetReadPreference(createReadPref(opt))
		default:
			t.Fatalf("unrecognized collection option: %v", name)
		}
	}

	return co
}

// create transaction options from a map
func createTransactionOptions(t testing.TB, opts bson.Raw) *options.TransactionOptions {
	t.Helper()

	txnOpts := options.Transaction()
	elems, _ := opts.Elements()
	for _, elem := range elems {
		name := elem.Key()
		opt := elem.Value()

		switch name {
		case "writeConcern":
			txnOpts.SetWriteConcern(createWriteConcern(t, opt))
		case "readPreference":
			txnOpts.SetReadPreference(createReadPref(opt))
		case "readConcern":
			txnOpts.SetReadConcern(createReadConcern(opt))
		case "maxCommitTimeMS":
			t := time.Duration(opt.Int32()) * time.Millisecond
			txnOpts.SetMaxCommitTime(&t)
		default:
			t.Fatalf("unrecognized transaction option: %v", opt)
		}
	}
	return txnOpts
}

// create a read concern from a map
func createReadConcern(opt bson.RawValue) *readconcern.ReadConcern {
	return readconcern.New(readconcern.Level(opt.Document().Lookup("level").StringValue()))
}

// create a read concern from a map
func createWriteConcern(t testing.TB, opt bson.RawValue) *writeconcern.WriteConcern {
	wcDoc, ok := opt.DocumentOK()
	if !ok {
		return nil
	}

	var opts []writeconcern.Option
	elems, _ := wcDoc.Elements()
	for _, elem := range elems {
		key := elem.Key()
		val := elem.Value()

		switch key {
		case "wtimeout":
			wtimeout := convertValueToMilliseconds(t, val)
			opts = append(opts, writeconcern.WTimeout(wtimeout))
		case "j":
			opts = append(opts, writeconcern.J(val.Boolean()))
		case "w":
			switch val.Type {
			case bson.TypeString:
				if val.StringValue() != "majority" {
					break
				}
				opts = append(opts, writeconcern.WMajority())
			case bson.TypeInt32:
				w := int(val.Int32())
				opts = append(opts, writeconcern.W(w))
			default:
				t.Fatalf("unrecognized type for w: %v", val.Type)
			}
		default:
			t.Fatalf("unrecognized write concern option: %v", key)
		}
	}
	return writeconcern.New(opts...)
}

// create a read preference from a string.
// returns readpref.Primary() if the string doesn't match any known read preference modes.
func readPrefFromString(s string) *readpref.ReadPref {
	switch strings.ToLower(s) {
	case "primary":
		return readpref.Primary()
	case "primarypreferred":
		return readpref.PrimaryPreferred()
	case "secondary":
		return readpref.Secondary()
	case "secondarypreferred":
		return readpref.SecondaryPreferred()
	case "nearest":
		return readpref.Nearest()
	}
	return readpref.Primary()
}

// create a read preference from a map.
func createReadPref(opt bson.RawValue) *readpref.ReadPref {
	mode := opt.Document().Lookup("mode").StringValue()
	return readPrefFromString(mode)
}

// transform a slice of BSON documents to a slice of interface{}.
func rawSliceToInterfaceSlice(docs []bson.Raw) []interface{} {
	out := make([]interface{}, len(docs))

	for i, doc := range docs {
		out[i] = doc
	}

	return out
}

// transform a BSON raw array to a slice of interface{}.
func rawArrayToInterfaceSlice(docs bson.Raw) []interface{} {
	vals, _ := docs.Values()

	out := make([]interface{}, len(vals))
	for i, val := range vals {
		out[i] = val.Document()
	}

	return out
}

// retrieve the error associated with a result.
func errorFromResult(t testing.TB, result interface{}) *operationError {
	t.Helper()

	// embedded doc will be unmarshalled as Raw
	raw, ok := result.(bson.Raw)
	if !ok {
		return nil
	}

	var expected operationError
	err := bson.Unmarshal(raw, &expected)
	if err != nil {
		return nil
	}
	if expected.ErrorCodeName == nil && expected.ErrorContains == nil && len(expected.ErrorLabelsOmit) == 0 &&
		len(expected.ErrorLabelsContain) == 0 {
		return nil
	}

	return &expected
}

// errorDetails is a helper type that holds information that can be returned by driver functions in different error
// types.
type errorDetails struct {
	name   string
	labels []string
}

// extractErrorDetails creates an errorDetails instance based on the provided error. It returns the details and an "ok"
// value which is true if the provided error is of a known type that can be processed.
func extractErrorDetails(err error) (errorDetails, bool) {
	var details errorDetails

	switch converted := err.(type) {
	case mongo.CommandError:
		details.name = converted.Name
		details.labels = converted.Labels
	case mongo.WriteException:
		if converted.WriteConcernError != nil {
			details.name = converted.WriteConcernError.Name
		}
		details.labels = converted.Labels
	case mongo.BulkWriteException:
		if converted.WriteConcernError != nil {
			details.name = converted.WriteConcernError.Name
		}
		details.labels = converted.Labels
	default:
		return errorDetails{}, false
	}

	return details, true
}

// verify that an error returned by an operation matches the expected error.
func verifyError(expected *operationError, actual error) error {
	// The spec test format doesn't treat ErrNoDocuments or ErrUnacknowledgedWrite as errors, so set actual to nil
	// to indicate that no error occurred.
	if actual == mongo.ErrNoDocuments || actual == mongo.ErrUnacknowledgedWrite {
		actual = nil
	}

	if expected == nil && actual != nil {
		return fmt.Errorf("did not expect error but got %v", actual)
	}
	if expected != nil && actual == nil {
		return fmt.Errorf("expected error but got nil")
	}
	if expected == nil {
		return nil
	}

	// check ErrorContains for all error types
	if expected.ErrorContains != nil {
		emsg := strings.ToLower(*expected.ErrorContains)
		amsg := strings.ToLower(actual.Error())
		if !strings.Contains(amsg, emsg) {
			return fmt.Errorf("expected error message %q to contain %q", amsg, emsg)
		}
	}

	// Get an errorDetails instance for the error. If this fails but the test has expectations about the error name or
	// labels, fail because we can't verify them.
	details, ok := extractErrorDetails(actual)
	if !ok {
		if expected.ErrorCodeName != nil || len(expected.ErrorLabelsContain) > 0 || len(expected.ErrorLabelsOmit) > 0 {
			return fmt.Errorf("failed to extract details from error %v of type %T", actual, actual)
		}
		return nil
	}

	if expected.ErrorCodeName != nil {
		if *expected.ErrorCodeName != details.name {
			return fmt.Errorf("expected error name %v, got %v", *expected.ErrorCodeName, details.name)
		}
	}
	for _, label := range expected.ErrorLabelsContain {
		if !stringSliceContains(details.labels, label) {
			return fmt.Errorf("expected error %v to contain label %q", actual, label)
		}
	}
	for _, label := range expected.ErrorLabelsOmit {
		if stringSliceContains(details.labels, label) {
			return fmt.Errorf("expected error %v to not contain label %q", actual, label)
		}
	}
	return nil
}

// get the underlying value of i as an int64. returns nil if i is not an int, int32, or int64 type.
func getIntFromInterface(i interface{}) *int64 {
	var out int64

	switch v := i.(type) {
	case int:
		out = int64(v)
	case int32:
		out = int64(v)
	case int64:
		out = v
	case float32:
		f := float64(v)
		if math.Floor(f) != f || f > float64(math.MaxInt64) {
			break
		}

		out = int64(f)
	case float64:
		if math.Floor(v) != v || v > float64(math.MaxInt64) {
			break
		}

		out = int64(v)
	default:
		return nil
	}

	return &out
}

func createCollation(t testing.TB, m bson.Raw) *options.Collation {
	var collation options.Collation
	elems, _ := m.Elements()

	for _, elem := range elems {
		switch elem.Key() {
		case "locale":
			collation.Locale = elem.Value().StringValue()
		case "caseLevel":
			collation.CaseLevel = elem.Value().Boolean()
		case "caseFirst":
			collation.CaseFirst = elem.Value().StringValue()
		case "strength":
			collation.Strength = int(elem.Value().Int32())
		case "numericOrdering":
			collation.NumericOrdering = elem.Value().Boolean()
		case "alternate":
			collation.Alternate = elem.Value().StringValue()
		case "maxVariable":
			collation.MaxVariable = elem.Value().StringValue()
		case "normalization":
			collation.Normalization = elem.Value().Boolean()
		case "backwards":
			collation.Backwards = elem.Value().Boolean()
		default:
			t.Fatalf("unrecognized collation option: %v", elem.Key())
		}
	}
	return &collation
}

func createChangeStreamOptions(t testing.TB, opts bson.Raw) *options.ChangeStreamOptions {
	t.Helper()

	csOpts := options.ChangeStream()
	elems, _ := opts.Elements()
	for _, elem := range elems {
		key := elem.Key()
		opt := elem.Value()

		switch key {
		case "batchSize":
			csOpts.SetBatchSize(opt.Int32())
		default:
			t.Fatalf("unrecognized change stream option: %v", key)
		}
	}
	return csOpts
}

func convertValueToMilliseconds(t testing.TB, val bson.RawValue) time.Duration {
	t.Helper()

	int32Val, ok := val.Int32OK()
	if !ok {
		t.Fatalf("failed to convert value of type %s to int32", val.Type)
	}
	return time.Duration(int32Val) * time.Millisecond
}

func stringSliceContains(stringSlice []string, target string) bool {
	for _, str := range stringSlice {
		if str == target {
			return true
		}
	}
	return false
}
