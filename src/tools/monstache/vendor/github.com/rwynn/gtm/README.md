gtm
===
gtm (go tail mongo) is a utility written in Go which tails the MongoDB oplog and 
sends create, update, delete events to your code.
It can be used to send emails to new users, [index documents](https://www.github.com/rwynn/monstache), 
[write time series data](https://www.github.com/rwynn/mongofluxd), or something else.

This branch is a port of the original gtm to use the new official golang driver from MongoDB.
The original gtm uses the community mgo driver. To use the community mgo driver use the `legacy` branch.

### Requirements ###
+ [Go](http://golang.org/doc/install)
+ [mongodb go driver](https://github.com/mongodb/mongo-go-driver)
+ [mongodb](http://www.mongodb.org/)

### Installation ###

	go get github.com/rwynn/gtm/v2

### Setup ###

gtm uses the MongoDB [oplog](https://docs.mongodb.com/manual/core/replica-set-oplog/) as an event source. 
You will need to ensure that MongoDB is configured to produce an oplog by 
[deploying a replica set](http://docs.mongodb.org/manual/tutorial/deploy-replica-set/).

If you haven't already done so, follow the 5 step 
[procedure](https://docs.mongodb.com/manual/tutorial/deploy-replica-set/#procedure) to initiate and 
validate your replica set. For local testing your replica set may contain a 
[single member](https://docs.mongodb.com/manual/tutorial/convert-standalone-to-replica-set/).

### Usage ###

```golang
package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/rwynn/gtm/v2"
	"reflect"
	"time"
)

func main() {
	rb := bson.NewRegistryBuilder()
	//rb.RegisterTypeMapEntry(bsontype.Timestamp, reflect.TypeOf(time.Time{}))
	rb.RegisterTypeMapEntry(bsontype.DateTime, reflect.TypeOf(time.Time{}))
	reg := rb.Build()
	clientOptions := options.Client()
	clientOptions.SetRegistry(reg)
	clientOptions.ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}
	ctxm, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctxm)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())
	ctx := gtm.Start(client, &gtm.Options{
		DirectReadNs: []string{"test.test"},
		ChangeStreamNs: []string{"test.test"},
		MaxWaitSecs: 10,
		OpLogDisabled: true,
	})
	for {
		select {
		case err := <-ctx.ErrC:
			fmt.Printf("got err %+v", err)
			break
		case op := <-ctx.OpC:
			fmt.Printf("got op %+v", op)
			break
		}
	}
}
```

### Configuration ###

```golang
func PipeBuilder(namespace string, changeStream bool) ([]interface{}, error) {

	// to build your pipelines for change events you will want to reference
	// the MongoDB reference for change events at 
	// https://docs.mongodb.com/manual/reference/change-events/

	// you will only receive changeStream == true when you configure gtm with
	// ChangeStreamNS (requies MongoDB 3.6+).  You cannot build pipelines for
	// changes using legacy direct oplog tailing

	if namespace == "users.users" {
		// given a set of docs like {username: "joe", email: "joe@email.com", amount: 1}
		if changeStream {
			return []interface{}{
				bson.M{"$match": bson.M{"fullDocument.username": "joe"}},
			}, nil
		} else {
			return []interface{}{
				bson.M{"$match": bson.M{"username": "joe"}},
			}, nil
		}
	} else if namespace == "users.status" && changeStream {
		// return a pipeline that only receives events when a document is 
		// inserted, deleted, or a specific field is changed. In this case
		// only a change to field1 is processed.  Changes to other fields
		// do not match the pipeline query and thus you won't receive the event.
		return []interface{}{
			bson.M{"$match": bson.M{"$or": []interface{} {
				bson.M{"updateDescription": bson.M{"$exists": false}},
				bson.M{"updateDescription.updatedFields.field1": bson.M{"$exists": true}},
			}}},
		}, nil
	}
	return nil, nil
}

func NewUsers(op *gtm.Op) bool {
	return op.Namespace == "users.users" && op.IsInsert()
}

// if you want to listen only for certain events on certain collections
// pass a filter function in options
ctx := gtm.Start(client, &gtm.Options{
	NamespaceFilter: NewUsers, // only receive inserts in the user collection
})
// more options are available for tuning
ctx := gtm.Start(client, &gtm.Options{
	NamespaceFilter      nil,           // op filter function that has access to type/ns ONLY
	Filter               nil,           // op filter function that has access to type/ns/data
	After:               nil,     	    // if nil defaults to gtm.LastOpTimestamp; not yet supported for ChangeStreamNS
	OpLogDisabled:       false,         // true to disable tailing the MongoDB oplog
	OpLogDatabaseName:   nil,     	    // defaults to "local"
	OpLogCollectionName: nil,     	    // defaults to "oplog.rs"
	ChannelSize:         0,       	    // defaults to 20
	BufferSize:          25,            // defaults to 50. used to batch fetch documents on bursts of activity
	BufferDuration:      0,             // defaults to 750 ms. after this timeout the batch is force fetched
	WorkerCount:         8,             // defaults to 1. number of go routines batch fetching concurrently
	Ordering:            gtm.Document,  // defaults to gtm.Oplog. ordering guarantee of events on the output channel as compared to the oplog
	UpdateDataAsDelta:   false,         // set to true to only receive delta information in the Data field on updates (info straight from oplog)
	DirectReadNs:        []string{"db.users"}, // set to a slice of namespaces (collections or views) to read data directly from
	DirectReadSplitMax:  9,             // the max number of times to split a collection for concurrent reads (impacts memory consumption)
	Pipe:                PipeBuilder,   // an optional function to build aggregation pipelines
	PipeAllowDisk:       false,         // true to allow MongoDB to use disk for aggregation pipeline options with large result sets
	Log:                 myLogger,      // pass your own logger
	ChangeStreamNs       []string{"db.col1", "db.col2"}, // MongoDB 3.6+ only; set to a slice to namespaces to read via MongoDB change streams
})
```

### Direct Reads ###

If, in addition to tailing the oplog, you would like to also read entire collections you can set the DirectReadNs field
to a slice of MongoDB namespaces.  Documents from these collections will be read directly and output on the ctx.OpC channel.  

You can wait till all the collections have been fully read by using the DirectReadWg wait group on the ctx.

```golang
go func() {
	ctx.DirectReadWg.Wait()
	fmt.Println("direct reads are done")
}()
```

### Pause, Resume, Since, and Stop ###

You can pause, resume, or seek to a timestamp from the oplog. These methods effect only change events and not direct reads.

```golang
go func() {
	ctx.Pause()
	time.Sleep(time.Duration(2) * time.Minute)
	ctx.Resume()
	ctx.Since(previousTimestamp)
}()
```

You can stop all goroutines created by `Start` or `StartMulti`. You cannot resume a context once it has been stopped. You would need to create a new one.

```golang
go func() {
	ctx.Stop()
	fmt.Println("all go routines are stopped")
}
```

### Custom Unmarshalling ###

If you'd like to unmarshall MongoDB documents into your own struct instead of the document getting
unmarshalled to a generic map[string]interface{} you can use a custom unmarshal function:

```golang
type MyDoc struct {
	Id interface{} "_id"
	Foo string "foo"
}

func custom(namespace string, data []byte) (interface{}, error) {
	// use namespace, e.g. db.col, to map to a custom struct
	if namespace == "test.test" {
		var doc MyDoc
		if err := bson.Unmarshal(data, &doc); err == nil {
			return doc, nil
		} else {
			return nil, err
		}
	}
	return nil, errors.New("unsupported namespace")
}

ctx := gtm.Start(client, &gtm.Options{
	Unmarshal: custom,
}

for {
	select {
	case op:= <-ctx.OpC:
		if op.Namespace == "test.test" {
			doc := op.Doc.(MyDoc)
			fmt.Println(doc.Foo)
		}
	}
}
```

### Workers ###

You may want to distribute event handling between a set of worker processes on different machines.
To do this you can leverage the **github.com/rwynn/gtm/consistent** package.  

Create a TOML document containing a list of all the event handlers.

```toml
Workers = [ "Tom", "Dick", "Harry" ] 
```

Create a consistent filter to distribute the work between Tom, Dick, and Harry. A consistent filter
needs to acces the Data attribute of each op so it needs to be set as a Filter as opposed to a 
NamespaceFilter.

```golang
name := flag.String("name", "", "the name of this worker")
flag.Parse()
filter, filterErr := consistent.ConsistentHashFilterFromFile(*name, "/path/to/toml")
if filterErr != nil {
	panic(filterErr)
}

// there is also a method **consistent.ConsistentHashFilterFromDocument** which allows
// you to pass a Mongo document representing the config if you would like to avoid
// copying the same config file to multiple servers
```

Pass the filter into the options when calling gtm.Tail

```golang
ctx := gtm.Start(client, &gtm.Options{Filter: filter})
```

If you have your multiple filters you can use the gtm utility method ChainOpFilters

```golang
func ChainOpFilters(filters ...OpFilter) OpFilter
```
