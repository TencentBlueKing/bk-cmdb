/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package local

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var (
	dvd = bsoncodec.DefaultValueDecoders{}
	// MapType use map decoder validation type
	mapType   = reflect.TypeOf(make(map[string]interface{}))
	ptBool    = reflect.TypeOf((*bool)(nil))
	ptInt8    = reflect.TypeOf((*int8)(nil))
	ptInt16   = reflect.TypeOf((*int16)(nil))
	ptInt32   = reflect.TypeOf((*int32)(nil))
	ptInt64   = reflect.TypeOf((*int64)(nil))
	ptInt     = reflect.TypeOf((*int)(nil))
	ptUint8   = reflect.TypeOf((*uint8)(nil))
	ptUint16  = reflect.TypeOf((*uint16)(nil))
	ptUint32  = reflect.TypeOf((*uint32)(nil))
	ptUint64  = reflect.TypeOf((*uint64)(nil))
	ptUint    = reflect.TypeOf((*uint)(nil))
	ptFloat32 = reflect.TypeOf((*float32)(nil))
	ptFloat64 = reflect.TypeOf((*float64)(nil))
	ptString  = reflect.TypeOf((*string)(nil))

	tBool    = reflect.TypeOf(false)
	tFloat32 = reflect.TypeOf(float32(0))
	tFloat64 = reflect.TypeOf(float64(0))
	tInt     = reflect.TypeOf(int(0))
	tInt8    = reflect.TypeOf(int8(0))
	tInt16   = reflect.TypeOf(int16(0))
	tInt32   = reflect.TypeOf(int32(0))
	tInt64   = reflect.TypeOf(int64(0))
	tString  = reflect.TypeOf("")
	tTime    = reflect.TypeOf(time.Time{})
	tUint    = reflect.TypeOf(uint(0))
	tUint8   = reflect.TypeOf(uint8(0))
	tUint16  = reflect.TypeOf(uint16(0))
	tUint32  = reflect.TypeOf(uint32(0))
	tUint64  = reflect.TypeOf(uint64(0))

	tEmpty      = reflect.TypeOf((*interface{})(nil)).Elem()
	tByteSlice  = reflect.TypeOf([]byte(nil))
	tByte       = reflect.TypeOf(byte(0x00))
	tURL        = reflect.TypeOf(url.URL{})
	tJSONNumber = reflect.TypeOf(json.Number(""))

	tValueMarshaler   = reflect.TypeOf((*bsoncodec.ValueMarshaler)(nil)).Elem()
	tValueUnmarshaler = reflect.TypeOf((*bsoncodec.ValueUnmarshaler)(nil)).Elem()
	tMarshaler        = reflect.TypeOf((*bsoncodec.Marshaler)(nil)).Elem()
	tUnmarshaler      = reflect.TypeOf((*bsoncodec.Unmarshaler)(nil)).Elem()
	tProxy            = reflect.TypeOf((*bsoncodec.Proxy)(nil)).Elem()

	tBinary        = reflect.TypeOf(primitive.Binary{})
	tUndefined     = reflect.TypeOf(primitive.Undefined{})
	tOID           = reflect.TypeOf(primitive.ObjectID{})
	tDateTime      = reflect.TypeOf(primitive.DateTime(0))
	tNull          = reflect.TypeOf(primitive.Null{})
	tRegex         = reflect.TypeOf(primitive.Regex{})
	tCodeWithScope = reflect.TypeOf(primitive.CodeWithScope{})
	tDBPointer     = reflect.TypeOf(primitive.DBPointer{})
	tJavaScript    = reflect.TypeOf(primitive.JavaScript(""))
	tSymbol        = reflect.TypeOf(primitive.Symbol(""))
	tTimestamp     = reflect.TypeOf(primitive.Timestamp{})
	tDecimal       = reflect.TypeOf(primitive.Decimal128{})
	tMinKey        = reflect.TypeOf(primitive.MinKey{})
	tMaxKey        = reflect.TypeOf(primitive.MaxKey{})
	tD             = reflect.TypeOf(primitive.D{})
	tA             = reflect.TypeOf(primitive.A{})
	tE             = reflect.TypeOf(primitive.E{})

	tCoreDocument = reflect.TypeOf(bsoncore.Document{})
)

func init() {

	bsonRegister := bson.NewRegistryBuilder()
	bsonRegister.RegisterDefaultDecoder(reflect.Map, bsoncodec.ValueDecoderFunc(MapDecodeValue))
	bsonRegister.RegisterDefaultDecoder(reflect.Array, bsoncodec.ValueDecoderFunc(ArrayDecodeValue))
	bsonRegister.RegisterDefaultDecoder(reflect.Slice, bsoncodec.ValueDecoderFunc(SliceDecodeValue))
	bsonRegister.RegisterDefaultDecoder(reflect.Struct, defaultStructCodec)
	bsonRegister.RegisterDecoder(tEmpty, bsoncodec.ValueDecoderFunc(EmptyInterfaceDecodeValue))

	bson.DefaultRegistry = bsonRegister.Build()
}

// MapDecodeValue is the ValueDecoderFunc for map[string]* types.
func MapDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {

	if !val.CanSet() {
		return bsoncodec.ValueDecoderError{Name: "MapDecodeValue", Kinds: []reflect.Kind{reflect.Map}, Received: val}
	}

	oldVal := val
	if isInterface(val.Type()) {
		val = reflect.New(mapType).Elem()
	}

	if val.Kind() != reflect.Map || val.Type().Key().Kind() != reflect.String {
		return bsoncodec.ValueDecoderError{Name: "MapDecodeValue", Kinds: []reflect.Kind{reflect.Map}, Received: val}
	}

	switch vr.Type() {
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	default:
		return fmt.Errorf("cannot decode %v into a %s", vr.Type(), val.Type())
	}

	dr, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	if val.IsNil() {
		val.Set(reflect.MakeMap(val.Type()))
	}

	eType := val.Type().Elem()
	decoder, err := dc.LookupDecoder(eType)
	if err != nil {
		return err
	}

	if eType == tEmpty {
		dc.Ancestor = val.Type()
	}

	keyType := val.Type().Key()
	for {
		key, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD {
			break
		}
		if err != nil {
			return err
		}

		elem := reflect.New(eType).Elem()

		err = decoder.DecodeValue(dc, vr, elem)
		if err != nil {
			return err
		}

		val.SetMapIndex(reflect.ValueOf(key).Convert(keyType), elem)
	}
	oldVal.Set(val)
	return nil
}

// ArrayDecodeValue is the ValueDecoderFunc for array types.
func ArrayDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {

	if !val.IsValid() || val.Kind() != reflect.Array {
		return bsoncodec.ValueDecoderError{Name: "ArrayDecodeValue", Kinds: []reflect.Kind{reflect.Array}, Received: val}
	}

	switch vr.Type() {
	case bsontype.Array:
	case bsontype.Type(0), bsontype.EmbeddedDocument:
		if extendEmbeddedDocumentDecoder(vr.Type(), val) {
			return MapDecodeValue(dc, vr, val)
		} else if val.Type().Elem() != tE {
			return fmt.Errorf("cannot decode document into %s", val.Type())
		}
	default:
		return fmt.Errorf("cannot decode %v into an array", vr.Type())
	}

	var elemsFunc func(bsoncodec.DecodeContext, bsonrw.ValueReader, reflect.Value) ([]reflect.Value, error)
	switch val.Type().Elem() {
	case tE:
		elemsFunc = decodeD
	default:
		elemsFunc = decodeDefault
	}

	elems, err := elemsFunc(dc, vr, val)
	if err != nil {
		return err
	}

	if len(elems) > val.Len() {
		return fmt.Errorf("more elements returned in array than can fit inside %s", val.Type())
	}

	for idx, elem := range elems {
		val.Index(idx).Set(elem)
	}

	return nil
}

// SliceDecodeValue is the ValueDecoderFunc for slice types.
func SliceDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {

	if !val.CanSet() || val.Kind() != reflect.Slice {
		return bsoncodec.ValueDecoderError{Name: "SliceDecodeValue", Kinds: []reflect.Kind{reflect.Slice}, Received: val}
	}

	switch vr.Type() {
	case bsontype.Array:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	case bsontype.Type(0), bsontype.EmbeddedDocument:
		if extendEmbeddedDocumentDecoder(vr.Type(), val) {
			return MapDecodeValue(dc, vr, val)
		} else if val.Type().Elem() != tE {
			return fmt.Errorf("cannot decode document into %s", val.Type())
		}
	default:
		return fmt.Errorf("cannot decode %v into a slice", vr.Type())
	}

	var elemsFunc func(bsoncodec.DecodeContext, bsonrw.ValueReader, reflect.Value) ([]reflect.Value, error)
	switch val.Type().Elem() {
	case tE:
		dc.Ancestor = val.Type()
		elemsFunc = decodeD
	default:
		elemsFunc = decodeDefault
	}

	elems, err := elemsFunc(dc, vr, val)
	if err != nil {
		return err
	}

	if val.IsNil() {
		val.Set(reflect.MakeSlice(val.Type(), 0, len(elems)))
	}

	val.SetLen(0)
	val.Set(reflect.Append(val, elems...))

	return nil
}

// EmptyInterfaceDecodeValue is the ValueDecoderFunc for interface{}.
func EmptyInterfaceDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tEmpty {
		return bsoncodec.ValueDecoderError{Name: "EmptyInterfaceDecodeValue", Types: []reflect.Type{tEmpty}, Received: val}
	}

	if extendEmbeddedDocumentDecoder(vr.Type(), val) {
		return MapDecodeValue(dc, vr, val)
	}

	rtype, err := dc.LookupTypeMapEntry(vr.Type())
	if err != nil {
		switch vr.Type() {
		case bsontype.EmbeddedDocument:
			if dc.Ancestor != nil {
				rtype = dc.Ancestor
				break
			}
			rtype = tD
		case bsontype.Null:
			val.Set(reflect.Zero(val.Type()))
			return vr.ReadNull()
		default:
			return err
		}
	}

	decoder, err := dc.LookupDecoder(rtype)
	if err != nil {
		return err
	}

	elem := reflect.New(rtype).Elem()
	err = decoder.DecodeValue(dc, vr, elem)
	if err != nil {
		return err
	}

	val.Set(elem)
	return nil
}

// IsInterface is interface
func isInterface(t reflect.Type) bool {
	if t.Kind() == reflect.Interface {
		return true
	}
	return false
}

// documentDecodeUseMapStrInterface Whether to enable bson EmbeddedDocument with bson Unmarshal value of interface to use map[string]interface object parsing
const documentDecodeUseMapStrInterface = true

// IsEmbeddedDocument  is bsontype EmbeddedDocument
func isEmbeddedDocument(bt bsontype.Type) bool {
	if bt == bsontype.EmbeddedDocument {
		return true
	}
	return false
}

// ExtendEmbeddedDocumentDecoder use map[string]interface decode bson  EmbeddedDocument to map[string]interface
func extendEmbeddedDocumentDecoder(vrType bsontype.Type, val reflect.Value) bool {
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
		(vrType == bsontype.Type(0) || isEmbeddedDocument(vrType)) &&
		isInterface(val.Type()) {
		return true
	}
	return false

}

// extendInterfaceDecode return documentDecodeUseMapStrInterface value
func extendInterfaceDecode(val reflect.Value) bool {
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
	if documentDecodeUseMapStrInterface && isInterface(val.Type()) {
		return true
	}
	return false
}

func decodeDefault(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) ([]reflect.Value, error) {
	elems := make([]reflect.Value, 0)

	ar, err := vr.ReadArray()
	if err != nil {
		return nil, err
	}

	eType := val.Type().Elem()

	decoder, err := dc.LookupDecoder(eType)
	if err != nil {
		return nil, err
	}

	for {
		vr, err := ar.ReadValue()
		if err == bsonrw.ErrEOA {
			break
		}
		if err != nil {
			return nil, err
		}

		elem := reflect.New(eType).Elem()

		err = decoder.DecodeValue(dc, vr, elem)
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}

	return elems, nil
}

func decodeD(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, _ reflect.Value) ([]reflect.Value, error) {
	switch vr.Type() {
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	default:
		return nil, fmt.Errorf("cannot decode %v into a D", vr.Type())
	}

	dr, err := vr.ReadDocument()
	if err != nil {
		return nil, err
	}

	return decodeElemsFromDocumentReader(dc, dr)
}

func decodeElemsFromDocumentReader(dc bsoncodec.DecodeContext, dr bsonrw.DocumentReader) ([]reflect.Value, error) {
	decoder, err := dc.LookupDecoder(tEmpty)
	if err != nil {
		return nil, err
	}

	elems := make([]reflect.Value, 0)
	for {
		key, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD {
			break
		}
		if err != nil {
			return nil, err
		}

		val := reflect.New(tEmpty).Elem()
		err = decoder.DecodeValue(dc, vr, val)
		if err != nil {
			return nil, err
		}

		elems = append(elems, reflect.ValueOf(primitive.E{Key: key, Value: val.Interface()}))
	}

	return elems, nil
}
