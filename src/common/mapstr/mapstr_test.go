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

package mapstr_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"configcenter/src/common/mapstr"
)

func TestConvertToMapStrFromInterface(t *testing.T) {

	// construct the test data
	testData := map[string]interface{}{
		"nil": nil,
		"map-int": map[string]int{
			"int-key": 1024,
		},
		"map-int-embed": map[string]interface{}{
			"embed-key": map[string]int{
				"embed-key-int": 1024,
			},
		},
		"struct": struct {
			TestStr     string
			testInt     int
			EmbedStruct interface{}
		}{
			TestStr: "test-str",
			testInt: 1024,
			EmbedStruct: struct {
				EmbedTestStr string
				EmbedTestInt int
			}{
				EmbedTestInt: 1024,
				EmbedTestStr: "embed-test-struct-str",
			},
		},
		"[]byte": []byte(`{"byte-array":"byte-array-valu"}`),
	}

	// execute the test
	for caseName, testItem := range testData {
		result, err := mapstr.NewFromInterface(testItem)
		require.NoError(t, err)
		require.NotNil(t, result)
		json, err := result.ToJSON()
		require.NoError(t, err)
		t.Logf("case:%s, the result:%s", caseName, string(json))
		if "struct" == caseName {
			subResult, err := result.MapStr("EmbedStruct")
			require.NoError(t, err)
			require.NotNil(t, subResult)
			subJson, err := subResult.ToJSON()
			require.NoError(t, err)
			t.Logf("case:%s embed struct:%s", caseName, subJson)
		}
	}

}

func TestMapStrInto(t *testing.T) {
	type testData struct {
		Ignor int
		Data  string `field:"data"`
		Test  int    `field:"test"`
	}
	target := mapstr.New()
	target.Set("test", 245)
	target.Set("data", "test_data")

	tmp := &testData{}
	target.MarshalJSONInto(tmp)
	//t.Logf("the test tmp %#v", tmp)

	maps := mapstr.NewArrayFromInterface([]map[string]interface{}{
		{"k": "value"}, {"i": 0},
	})
	target1 := maps[0]

	target2, err := mapstr.NewFromInterface(map[string]interface{}{
		"k": "v", "i": 1, "j": 2, "time": time.Now(), "map": map[string]interface{}{}, "bool": true,
	})
	if err != nil {
		t.Fail()
	}
	target1.Different(target2)
	target1.Merge(target2)

	target1.Set("set_key", "set_value")
	_, ok := target1.Get("set_key")
	if !ok {
		t.Fail()
	}

	if i, _ := target1.Int64("i"); i != 1 {
		t.Fail()
	}

	if b, _ := target1.Bool("bool"); !b {
		t.Fail()
	}

	if i, _ := target1.Float("i"); i != 1 {
		t.Fail()
	}

	if s, _ := target1.String("k"); s != "v" {
		t.Fail()
	}

	if _, err := target1.Time("time"); err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if _, err := target1.MapStr("map"); err != nil {
		t.Fail()
	}

	target1.Set("maps", maps)
	if _, err := target1.MapStrArray("maps"); err != nil {
		t.Fail()
	}

	if target1.ForEach(func(key string, val interface{}) error {
		return nil
	}) != nil {
		t.Fail()
	}

	target1.Remove("maps")
	if target1.Exists("maps") {
		t.Fail()
	}

	if target1.IsEmpty() {
		t.Fail()
	}

	target1.Reset()

}

func TestMapStrToMapstr(t *testing.T) {

	testData := mapstr.MapStr{"aa": "bb"}

	_, err := mapstr.NewFromInterface(testData)
	if err != nil {
		t.Error(err.Error())
		return
	}

	testData2 := mapstr.MapStr{"aa": []mapstr.MapStr{
		mapstr.MapStr{"aa": "bb"},
	}}
	_, err = testData2.MapStrArray("aa")
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestMapToMapStr(t *testing.T) {

	newMapStr := mapstr.NewFromMap(map[string]interface{}{"aa": "bb"})
	val, err := newMapStr.String("aa")
	if err != nil {
		t.Errorf("get aa value error, err:%s", err.Error())
		return
	}
	if val != "bb" {
		t.Errorf("aa value not equal aa, val:%s", val)
		return
	}
}

type TargetInline struct {
	Field1Inline string `field:"field_inline_one"`
	Field2Inline int    `field:"field_inline_two"`
}
type Label mapstr.MapStr

func TestNewFromStruct(t *testing.T) {

	type targetTest struct {
		Field1   string       `field:"field_one"`
		Field2   int          `field:"field_two"`
		Labels   Label        `field:"field_mapstr"`
		TargetIn TargetInline `field:"field_inline"`
	}

	targetMapStr := mapstr.NewFromStruct(&targetTest{
		Field1: "field1",
		Field2: 2,
		Labels: Label{"key": "value"},
		TargetIn: TargetInline{
			Field1Inline: "field_in_line",
			Field2Inline: 2,
		},
	}, "field")
	t.Logf("target mapstr %v", targetMapStr)
}

func TestConvertToMapStrFromStructInnerPointer(t *testing.T) {

	type targetTest struct {
		Field1   string        `field:"field_one"`
		Field2   int           `field:"field_two"`
		Field3   *string       `field:"field_three"`
		Labels   Label         `field:"field_mapstr"`
		TargetIn *TargetInline `field:"field_inline"`
	}
	tmpStr := "field3-str"

	targetMapStr := mapstr.NewFromStruct(&targetTest{
		Field1: "field1",
		Field2: 2,
		Field3: &tmpStr,
		Labels: Label{"key": "value"},
		TargetIn: &TargetInline{
			Field1Inline: "field_in_line",
			Field2Inline: 2,
		},
	}, "field")
	targteJson, err := targetMapStr.ToJSON()
	require.NoError(t, err)
	t.Logf("target mapstr %s", targteJson)

	resultTmp := targetTest{}
	err = targetMapStr.ToStructByTag(&resultTmp, "field")
	require.NoError(t, err)
	t.Logf("result struct:%#v %s", resultTmp, *resultTmp.Field3)
}

func TestConvertToMapStrFromStructInnerEmbedPointer(t *testing.T) {

	type targetTest struct {
		Field1        string  `field:"field_one"`
		Field2        int     `field:"field_two"`
		Field3        float64 `field:"field_float"`
		Labels        Label   `field:"field_mapstr"`
		*TargetInline `field:"field_inline_target"`
	}

	targetMapStr := mapstr.NewFromStruct(&targetTest{
		Field1: "field1",
		Field2: 2,
		Field3: 3.3,
		Labels: Label{"key": "value"},
		TargetInline: &TargetInline{
			Field1Inline: "field_in_line",
			Field2Inline: 2,
		},
	}, "field")

	t.Logf("target mapstr %v", targetMapStr)

	resultTmp := targetTest{}
	err := targetMapStr.ToStructByTag(&resultTmp, "field")
	require.NoError(t, err)
	t.Logf("result struct :%#v", resultTmp)
}

func TestConvertToMapStrFromStructEmbed(t *testing.T) {

	type targetTest struct {
		Field1       string `field:"field_one"`
		Field2       int    `field:"field_two"`
		Labels       Label  `field:"field_mapstr"`
		TargetInline `field:"field_inline_target"`
	}

	targetMapStr := mapstr.NewFromStruct(&targetTest{
		Field1: "field1",
		Field2: 2,
		Labels: Label{"key": "value"},
		TargetInline: TargetInline{
			Field1Inline: "field_in_line",
			Field2Inline: 2,
		},
	}, "field")

	t.Logf("target mapstr %v", targetMapStr)

	resultTmp := targetTest{}
	err := targetMapStr.ToStructByTag(&resultTmp, "field")
	require.NoError(t, err)
	t.Logf("result struct :%v", resultTmp)
}

func TestEmbedMap(t *testing.T) {

	type Label map[string]string

	type Metadata struct {
		Label Label `field:"label" json:"label" bson:"label"`
	}

	type classification struct {
		Metadata           Metadata `field:"metadata" json:"metadata" bson:"metadata"`
		ID                 int64    `field:"id" json:"id" bson:"id"`
		ClassificationID   string   `field:"bk_classification_id"  json:"bk_classification_id" bson:"bk_classification_id"`
		ClassificationName string   `field:"bk_classification_name" json:"bk_classification_name" bson:"bk_classification_name"`
		ClassificationType string   `field:"bk_classification_type" json:"bk_classification_type" bson:"bk_classification_type"`
		ClassificationIcon string   `field:"bk_classification_icon" json:"bk_classification_icon" bson:"bk_classification_icon"`
		OwnerID            string   `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"  `
	}

	testData := `{
		"bk_supplier_account": "0",
		"bk_classification_id": "test",
		"bk_classification_name": "test",
		"metadata": {
			"label": {
				"bk_biz_id": "1"
			}
		}
	}`

	data, err := mapstr.NewFromInterface(testData)
	require.NoError(t, err)

	out := &classification{}
	err = mapstr.SetValueToStructByTags(out, data)
	require.NoError(t, err)
	t.Logf("output:%#v", out)

}

func TestClone(t *testing.T) {

	type targetTest struct {
		Field1        string  `field:"field_one"`
		Field2        int     `field:"field_two"`
		Field3        float64 `field:"field_float"`
		Labels        Label   `field:"field_mapstr"`
		privateField  string
		*TargetInline `field:"field_inline_target"`
	}

	targetMapStr := mapstr.MapStr{
		"struct": &targetTest{
			Field1:       "field1",
			Field2:       2,
			Field3:       3.3,
			privateField: "hello private",
			Labels:       Label{"key": "value"},
			TargetInline: &TargetInline{
				Field1Inline: "field_in_line",
				Field2Inline: 2,
			},
		},
	}

	cloneInst := targetMapStr.Clone()
	(cloneInst["struct"].(*targetTest)).TargetInline.Field1Inline = "inline_hello"
	t.Logf("origin:%#v private:%s", (targetMapStr["struct"].(*targetTest)).TargetInline.Field1Inline, (targetMapStr["struct"].(*targetTest)).privateField)
	t.Logf("clone:%#v private:%s", (cloneInst["struct"].(*targetTest)).TargetInline.Field1Inline, (cloneInst["struct"].(*targetTest)).privateField)
}

func TestLargeNumber(t *testing.T) {
	testData := mapstr.MapStr{
		"int": int(5000000000000000000),
		"uint": uint(5000000000000000000),
		"int64": int64(5000000000000000000),
		"uint64": uint64(5000000000000000000),
		"float64": float64(5000000000000000000),
		"float32": float32(5000000000000000000),
	}
	str, err := testData.String("int")
	require.NoError(t, err)
	require.Equal(t, str, "5000000000000000000")
	str, err = testData.String("uint")
	require.NoError(t, err)
	require.Equal(t, str, "5000000000000000000")
	str, err = testData.String("int64")
	require.NoError(t, err)
	require.Equal(t, str, "5000000000000000000")
	str, err = testData.String("uint64")
	require.NoError(t, err)
	require.Equal(t, str, "5000000000000000000")
	str, err = testData.String("float64")
	require.NoError(t, err)
	require.Equal(t, str, "5000000000000000000")
	str, err = testData.String("float32")
	require.NoError(t, err)
	require.Equal(t, str, "5000000000000000000")
}
