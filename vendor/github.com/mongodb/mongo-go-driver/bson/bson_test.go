// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mongodb/mongo-go-driver/x/bsonx/bsoncore"
	"github.com/stretchr/testify/require"
)

func noerr(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Errorf("Unexpected error: (%T)%v", err, err)
		t.FailNow()
	}
}

func requireErrEqual(t *testing.T, err1 error, err2 error) { require.True(t, compareErrors(err1, err2)) }

func TestTimeRoundTrip(t *testing.T) {
	val := struct {
		Value time.Time
		ID    string
	}{
		ID: "time-rt-test",
	}

	if !val.Value.IsZero() {
		t.Errorf("Did not get zero time as expected.")
	}

	bsonOut, err := Marshal(val)
	noerr(t, err)
	rtval := struct {
		Value time.Time
		ID    string
	}{}

	err = Unmarshal(bsonOut, &rtval)
	noerr(t, err)
	if !cmp.Equal(val, rtval) {
		t.Errorf("Did not round trip properly. got %v; want %v", val, rtval)
	}
	if !rtval.Value.IsZero() {
		t.Errorf("Did not get zero time as expected.")
	}

	beforeTimeVal := time.Now()
	beforeTimeMap := map[string]interface{}{"ts": beforeTimeVal}
	before, err := Marshal(beforeTimeMap)
	noerr(t, err)

	afterTimeMap := map[string]interface{}{}
	err = Unmarshal(before, &afterTimeMap)
	noerr(t, err)
	afterTime, ok := afterTimeMap["ts"].(time.Time)
	if !ok {
		t.Errorf("after time format error. after time info:%#v", afterTimeMap)
		return
	}
	if afterTime.Unix() != beforeTimeVal.Unix() {
		t.Errorf("after time not equal before time. befere:%#v, after:%#v", beforeTimeVal, afterTime)
		return
	}

}

func TestNonNullTimeRoundTrip(t *testing.T) {
	now := time.Now()
	now = time.Unix(now.Unix(), 0)
	val := struct {
		Value time.Time
		ID    string
	}{
		ID:    "time-rt-test",
		Value: now,
	}

	bsonOut, err := Marshal(val)
	noerr(t, err)
	rtval := struct {
		Value time.Time
		ID    string
	}{}

	err = Unmarshal(bsonOut, &rtval)
	noerr(t, err)
	if !cmp.Equal(val, rtval) {
		t.Errorf("Did not round trip properly. got %v; want %v", val, rtval)
	}
}

func TestD(t *testing.T) {
	t.Run("can marshal", func(t *testing.T) {
		d := D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
		idx, want := bsoncore.AppendDocumentStart(nil)
		want = bsoncore.AppendStringElement(want, "foo", "bar")
		want = bsoncore.AppendStringElement(want, "hello", "world")
		want = bsoncore.AppendDoubleElement(want, "pi", 3.14159)
		want, err := bsoncore.AppendDocumentEnd(want, idx)
		noerr(t, err)
		got, err := Marshal(d)
		noerr(t, err)
		if !bytes.Equal(got, want) {
			t.Errorf("Marshaled documents do not match. got %v; want %v", Raw(got), Raw(want))
		}
	})
	t.Run("can unmarshal", func(t *testing.T) {
		want := D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
		idx, doc := bsoncore.AppendDocumentStart(nil)
		doc = bsoncore.AppendStringElement(doc, "foo", "bar")
		doc = bsoncore.AppendStringElement(doc, "hello", "world")
		doc = bsoncore.AppendDoubleElement(doc, "pi", 3.14159)
		doc, err := bsoncore.AppendDocumentEnd(doc, idx)
		noerr(t, err)
		var got D
		err = Unmarshal(doc, &got)
		noerr(t, err)
		if !cmp.Equal(got, want) {
			t.Errorf("Unmarshaled documents do not match. got %v; want %v", got, want)
		}
	})
}
