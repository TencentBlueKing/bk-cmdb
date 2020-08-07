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

package types

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDocuments_Decode(t *testing.T) {
	type V struct {
		Name string
		Time *time.Time
	}

	now := time.Now()
	expect := []V{{Name: "A", Time: &now}, {Name: "B", Time: &now}}
	docs := Documents{}
	err := docs.Encode(expect)
	require.NoError(t, err)
	t.Logf("docs: %+v", docs)

	actual := []V{}

	err = docs.Decode(&actual)
	require.NoError(t, err)
	t.Logf("actual: %+v", actual)

	for index, v := range actual {
		require.Equal(t, v.Name, expect[index].Name)
		require.Equal(t, v.Time.Format("2006-01-02 15:04:05"), expect[index].Time.Format("2006-01-02 15:04:05"))
	}

}

func TestDecodeArray(t *testing.T) {
	type V struct {
		Name string
		Time time.Time
	}

	now := time.Now()

	expect := []V{{Name: "A", Time: now}, {Name: "B", Time: now}}
	docs := Documents{}
	err := decodeBsonArray(expect, &docs)
	require.NoError(t, err)
	t.Logf("docs: %+v", docs)

	actual := []V{}
	err = decodeBsonArray(docs, &actual)
	require.NoError(t, err)
	t.Logf("actual: %+v", actual)

	for index, v := range actual {
		require.Equal(t, v.Name, expect[index].Name)
		require.Equal(t, v.Time.Format("2006-01-02 15:04:05"), expect[index].Time.Format("2006-01-02 15:04:05"))
	}

}

func TestDecodeSubStruct(t *testing.T) {
	type SubStruct struct {
		AA map[string]interface{}
	}
	type MainStruct struct {
		S    SubStruct
		Name string
		Time time.Time
	}

	ss := SubStruct{
		AA: map[string]interface{}{"1": 1, "2": "bbb"},
	}
	ms := MainStruct{
		S:    ss,
		Name: "test",
		Time: time.Now(),
	}

	msbytes, err := bson.Marshal(ms)
	if err != nil {
		panic(err)
	}

	tmpMap := make(map[string]interface{}, 0)
	err = bson.Unmarshal(msbytes, &tmpMap)
	if err != nil {
		panic(err)
	}

	msbytes, err = bson.Marshal(tmpMap)
	if err != nil {
		panic(err)
	}

	out := MainStruct{}
	err = bson.Unmarshal(msbytes, &out)
	if err != nil {
		panic(err)
	}

	outArr := []MainStruct{}
	elemt := reflect.ValueOf(outArr).Type().Elem()

	elem := reflect.New(elemt)
	err = bson.Unmarshal(msbytes, elem.Interface())
	if err != nil {
		panic(err)
	}

	resultv := reflect.ValueOf(&outArr)
	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt2 := slicev.Type().Elem()

	elem2 := reflect.New(elemt2)
	err = bson.Unmarshal(msbytes, elem2.Interface())
	if err != nil {
		panic(err)
	}

	t.Logf("%#v   %#v  %#v", out, elem, elem2)
}

// Classification the classification metadata definition
type Classification struct {
	Time time.Time `field:"last_time"  json:"last_time" bson:"last_time"`
}

func TestStruct(t *testing.T) {
	docItem := map[string]interface{}{

		"last_time": "2018-11-16T08:18:13.946Z",
	}

	out, err := bson.Marshal(docItem)
	if nil != err {
		t.Errorf("Decode array error when marshal: %v, source is %s", err, docItem)
	}

	newData := &Classification{}
	err = bson.Unmarshal(out, newData)
	if nil != err {
		t.Errorf("Decode array error when unmarshal: %v, source is %v", err, newData)
	}
	t.Log(newData)

}
