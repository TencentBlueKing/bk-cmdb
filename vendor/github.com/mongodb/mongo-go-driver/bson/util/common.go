package util

import (
	"reflect"

	"github.com/mongodb/mongo-go-driver/bson/bsontype"
)

// MapType use map decoder validation type
var MapType = reflect.TypeOf(make(map[string]interface{}))

// IsInterface is interface
func IsInterface(t reflect.Type) bool {
	if t.Kind() == reflect.Interface {
		return true
	}
	return false
}

// documentDecodeUseMapStrInterface Whether to enable bson EmbeddedDocument with bson Unmarshal value of interface to use map[string]interface object parsing
const documentDecodeUseMapStrInterface = true

// IsEmbeddedDocument  is bsontype EmbeddedDocument
func IsEmbeddedDocument(bt bsontype.Type) bool {
	if bt == bsontype.EmbeddedDocument {
		return true
	}
	return false
}

// ExtendEmbeddedDocumentDecoder use map[string]interface decode bson  EmbeddedDocument to map[string]interface
func ExtendEmbeddedDocumentDecoder(vrType bsontype.Type, val reflect.Value) bool {
	/*
		eg: test case
		instanceMap := map[string]interface{}{"value": map[string]interface{}{"a": []int{12, 14}, "b": map[string]interface{}{"cc": 11}}}
		instanceByte, err := bson.Marshal(instanceMap)
		require.NoError(t, err)

		var i map[string]interface{}
		err = bson.Unmarshal(instanceByte, &i)
		require.NoError(t, err)

		result :
		{"value":{"a":[12,14],"b":{"cc":11}}}
		not  [{"Key":"value","Value":[{"Key":"a","Value":12},{"Key":"b","Value":122}]}]


	*/

	if documentDecodeUseMapStrInterface &&
		(vrType == bsontype.Type(0) || IsEmbeddedDocument(vrType)) &&
		IsInterface(val.Type()) {
		return true
	}
	return false

}

// ExtendInterfaceDecode return documentDecodeUseMapStrInterface value
func ExtendInterfaceDecode(val reflect.Value) bool {
	/*
		eg: test case
		instanceMap := map[string]interface{}{"value": map[string]interface{}{"a": []int{12, 14}, "b": map[string]interface{}{"cc": 11}}}
		instanceByte, err := bson.Marshal(instanceMap)
		require.NoError(t, err)

		var i map[string]interface{}
		err = bson.Unmarshal(instanceByte, &i)
		require.NoError(t, err)

		result :
		{"value":{"a":[12,14],"b":{"cc":11}}}
		not  [{"Key":"value","Value":[{"Key":"a","Value":12},{"Key":"b","Value":122}]}]


	*/
	if documentDecodeUseMapStrInterface && IsInterface(val.Type()) {
		return true
	}
	return false
}
