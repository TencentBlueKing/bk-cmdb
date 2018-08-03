package mongobyc

// #include <stdlib.h>
// #include "mongo.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// InitMongoc init the mongc lib
func InitMongoc() {
	C.mongoc_init()
}

// CleanupMongoc cleanup the mongc lib
func CleanupMongoc() {
	C.mongoc_cleanup()
}

// Client client for mongo
type Client interface {
	Open() error
	Close() error
	Ping() error
	Collection(collName string) CollectionInterface
}

// NewClient create a mongoc client instance
func NewClient(uri, database string) Client {
	return &client{
		uri:    uri,
		dbName: database,
	}
}

type client struct {
	uri         string
	dbName      string
	db          *C.mongoc_database_t
	innerClient *C.mongoc_client_t
}

func (c *client) Open() error {

	// create client
	var err C.bson_error_t
	uri := C.mongoc_uri_new_with_error(C.CString(c.uri), &err)
	if nil == uri {
		return TransformError(err)
	}

	c.innerClient = C.mongoc_client_new_from_uri(uri)
	if nil == c.innerClient {
		return fmt.Errorf("can not create a client instance")
	}

	C.mongoc_uri_destroy(uri)

	// set app name
	cName := C.CString(c.dbName)
	C.mongoc_client_set_appname(c.innerClient, cName)

	// get database by name
	c.db = C.mongoc_client_get_database(c.innerClient, cName)
	C.free(unsafe.Pointer(cName))

	return nil
}

func (c *client) Close() error {
	if nil != c.innerClient {
		C.mongoc_client_destroy(c.innerClient)
	}
	return nil
}

func (c *client) Collection(collName string) CollectionInterface {
	return newCollection(c.innerClient, c.dbName, collName)
}

func (c *client) Ping() error {

	pingCStr := C.CString("ping")
	pingCmd := C.create_bcon_new_int32(pingCStr, 1)

	var reply C.bson_t
	var err C.bson_error_t
	adminCStr := C.CString("admin")
	ok := C.mongoc_client_command_simple(c.innerClient, adminCStr, pingCmd, nil, &reply, &err)

	C.free(unsafe.Pointer(pingCStr))
	C.free(unsafe.Pointer(adminCStr))
	if !ok {
		return TransformError(err)
	}

	return nil
}
