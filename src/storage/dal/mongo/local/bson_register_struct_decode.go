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

// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package local

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var defaultStructCodec = &StructCodec{
	cache:  make(map[reflect.Type]*structDescription),
	parser: bsoncodec.DefaultStructTagParser,
}

// Zeroer allows custom struct types to implement a report of zero
// state. All struct types that don't implement Zeroer or where IsZero
// returns false are considered to be not zero.
type Zeroer interface {
	IsZero() bool
}

// StructCodec is the Codec used for struct values.
type StructCodec struct {
	cache  map[reflect.Type]*structDescription
	l      sync.RWMutex
	parser bsoncodec.StructTagParser
}

// NewStructCodec returns a StructCodec that uses p for struct tag parsing.
func NewStructCodec(p bsoncodec.StructTagParser) (*StructCodec, error) {
	if p == nil {
		return nil, errors.New("a StructTagParser must be provided to NewStructCodec")
	}

	return &StructCodec{
		cache:  make(map[reflect.Type]*structDescription),
		parser: p,
	}, nil
}

// DecodeValue implements the Codec interface.
// By default, map types in val will not be cleared. If a map has existing key/value pairs, it will be extended with the new ones from vr.
// For slices, the decoder will set the length of the slice to zero and append all elements. The underlying array will not be cleared.
func (sc *StructCodec) DecodeValue(r bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != reflect.Struct {
		return bsoncodec.ValueDecoderError{Name: "StructCodec.DecodeValue", Kinds: []reflect.Kind{reflect.Struct}, Received: val}
	}

	switch vr.Type() {
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	default:
		return fmt.Errorf("cannot decode %v into a %s", vr.Type(), val.Type())
	}

	sd, err := sc.describeStruct(r.Registry, val.Type())
	if err != nil {
		return err
	}

	var decoder bsoncodec.ValueDecoder
	var inlineMap reflect.Value
	if sd.inlineMap >= 0 {
		inlineMap = val.Field(sd.inlineMap)
		if inlineMap.IsNil() {
			inlineMap.Set(reflect.MakeMap(inlineMap.Type()))
		}
		decoder, err = r.LookupDecoder(inlineMap.Type().Elem())
		if err != nil {
			return err
		}
	}

	dr, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	for {
		name, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD {
			break
		}
		if err != nil {
			return err
		}

		fd, exists := sd.fm[name]
		if !exists {
			if sd.inlineMap < 0 {
				// The encoding/json package requires a flag to return on error for non-existent fields.
				// This functionality seems appropriate for the struct codec.
				err = vr.Skip()
				if err != nil {
					return err
				}
				continue
			}

			elem := reflect.New(inlineMap.Type().Elem()).Elem()
			err = decoder.DecodeValue(r, vr, elem)
			if err != nil {
				return err
			}

			inlineMap.SetMapIndex(reflect.ValueOf(name), elem)
			continue
		}

		var field reflect.Value
		if fd.inline == nil {
			field = val.Field(fd.idx)
		} else {
			field = val.FieldByIndex(fd.inline)
		}
		if extendEmbeddedDocumentDecoder(vr.Type(), field) {
			decoder, err := r.LookupDecoder(mapType)
			if err == nil {
				fd.decoder = decoder
			}
		}

		// compatible with null value
		if vr.Type() == bsontype.Null {
			if err := vr.ReadNull(); err != nil {
				return err
			}
			continue
		}

		if !field.CanSet() { // Being settable is a super set of being addressable.
			return fmt.Errorf("cannot decode element '%s' into field %v; it is not settable", name, field)
		}
		if field.Kind() == reflect.Ptr && field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Addr()

		dctx := bsoncodec.DecodeContext{Registry: r.Registry, Truncate: fd.truncate}
		if fd.decoder == nil {
			return bsoncodec.ErrNoDecoder{Type: field.Elem().Type()}
		}
		if decoder, ok := fd.decoder.(bsoncodec.ValueDecoder); ok {
			err = decoder.DecodeValue(dctx, vr, field.Elem())
			if err != nil {
				return err
			}
			continue
		}
		err = fd.decoder.DecodeValue(dctx, vr, field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sc *StructCodec) isZero(i interface{}) bool {
	v := reflect.ValueOf(i)

	// check the value validity
	if !v.IsValid() {
		return true
	}

	if z, ok := v.Interface().(Zeroer); ok && (v.Kind() != reflect.Ptr || !v.IsNil()) {
		return z.IsZero()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return false
}

type structDescription struct {
	fm        map[string]fieldDescription
	fl        []fieldDescription
	inlineMap int
}

type fieldDescription struct {
	name      string
	idx       int
	omitEmpty bool
	minSize   bool
	truncate  bool
	inline    []int
	encoder   bsoncodec.ValueEncoder
	decoder   bsoncodec.ValueDecoder
}

func (sc *StructCodec) describeStruct(r *bsoncodec.Registry, t reflect.Type) (*structDescription, error) {
	// We need to analyze the struct, including getting the tags, collecting
	// information about inlining, and create a map of the field name to the field.
	sc.l.RLock()
	ds, exists := sc.cache[t]
	sc.l.RUnlock()
	if exists {
		return ds, nil
	}

	numFields := t.NumField()
	sd := &structDescription{
		fm:        make(map[string]fieldDescription, numFields),
		fl:        make([]fieldDescription, 0, numFields),
		inlineMap: -1,
	}

	for i := 0; i < numFields; i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" {
			// unexported, ignore
			continue
		}

		encoder, err := r.LookupEncoder(sf.Type)
		if err != nil {
			encoder = nil
		}
		decoder, err := r.LookupDecoder(sf.Type)
		if err != nil {
			decoder = nil
		}

		description := fieldDescription{idx: i, encoder: encoder, decoder: decoder}

		stags, err := sc.parser.ParseStructTags(sf)
		if err != nil {
			return nil, err
		}
		if stags.Skip {
			continue
		}
		description.name = stags.Name
		description.omitEmpty = stags.OmitEmpty
		description.minSize = stags.MinSize
		description.truncate = stags.Truncate

		if stags.Inline {
			switch sf.Type.Kind() {
			case reflect.Map:
				if sd.inlineMap >= 0 {
					return nil, errors.New("(struct " + t.String() + ") multiple inline maps")
				}
				if sf.Type.Key() != tString {
					return nil, errors.New("(struct " + t.String() + ") inline map must have a string keys")
				}
				sd.inlineMap = description.idx
			case reflect.Struct:
				inlinesf, err := sc.describeStruct(r, sf.Type)
				if err != nil {
					return nil, err
				}
				for _, fd := range inlinesf.fl {
					if _, exists := sd.fm[fd.name]; exists {
						return nil, fmt.Errorf("(struct %s) duplicated key %s", t.String(), fd.name)
					}
					if fd.inline == nil {
						fd.inline = []int{i, fd.idx}
					} else {
						fd.inline = append([]int{i}, fd.inline...)
					}
					sd.fm[fd.name] = fd
					sd.fl = append(sd.fl, fd)
				}
			default:
				return nil, fmt.Errorf("(struct %s) inline fields must be either a struct or a map", t.String())
			}
			continue
		}

		if _, exists := sd.fm[description.name]; exists {
			return nil, fmt.Errorf("struct %s) duplicated key %s", t.String(), description.name)
		}

		sd.fm[description.name] = description
		sd.fl = append(sd.fl, description)
	}

	sc.l.Lock()
	sc.cache[t] = sd
	sc.l.Unlock()

	return sd, nil
}
