// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"reflect"

	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"
)

type unmarshalingTestCase struct {
	name  string
	reg   *bsoncodec.Registry
	sType reflect.Type
	want  interface{}
	data  []byte
}

var unmarshalingTestCases = []unmarshalingTestCase{
	{
		"small struct",
		nil,
		reflect.TypeOf(struct {
			Foo bool
		}{}),
		&struct {
			Foo bool
		}{Foo: true},
		docToBytes(D{{"foo", true}}),
	},
	{
		"nested document",
		nil,
		reflect.TypeOf(struct {
			Foo struct {
				Bar bool
			}
		}{}),
		&struct {
			Foo struct {
				Bar bool
			}
		}{
			Foo: struct {
				Bar bool
			}{Bar: true},
		},
		docToBytes(D{{"foo", D{{"bar", true}}}}),
	},
	{
		"simple array",
		nil,
		reflect.TypeOf(struct {
			Foo []bool
		}{}),
		&struct {
			Foo []bool
		}{
			Foo: []bool{true},
		},
		docToBytes(D{{"foo", A{true}}}),
	},
}
