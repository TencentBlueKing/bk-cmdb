/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package filter

import (
	"reflect"
	"testing"
	"time"

	"configcenter/src/common"
)

func TestEqualMongoCond(t *testing.T) {
	op := Equal.Factory().Operator()

	// test equal int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBEQ: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test equal string type
	cond, err = op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBEQ: "a"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test equal bool type
	cond, err = op.ToMgo("test", false)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBEQ: false}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid equal type
	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotEqualMongoCond(t *testing.T) {
	op := NotEqual.Factory().Operator()

	// test not equal int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBNE: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test not equal string type
	cond, err = op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBNE: "a"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test not equal bool type
	cond, err = op.ToMgo("test", false)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBNE: false}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not equal type
	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestInMongoCond(t *testing.T) {
	op := In.Factory().Operator()

	// test in int array type
	cond, err := op.ToMgo("test", []int64{1, 2})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBIN: []int64{1, 2}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test in string array type
	cond, err = op.ToMgo("test", []string{"a", "b"})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBIN: []string{"a", "b"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test in bool array type
	cond, err = op.ToMgo("test", []interface{}{true, false})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBIN: []interface{}{true, false}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid in type
	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []interface{}{1, "a"})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotInMongoCond(t *testing.T) {
	op := NotIn.Factory().Operator()

	// test not in int array type
	cond, err := op.ToMgo("test", []int64{1, 2})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBNIN: []int64{1, 2}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test not in string array type
	cond, err = op.ToMgo("test", []string{"a", "b"})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBNIN: []string{"a", "b"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test not in bool array type
	cond, err = op.ToMgo("test", []interface{}{true, false})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBNIN: []interface{}{true, false}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not in type
	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []interface{}{1, "a"})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestLessMongoCond(t *testing.T) {
	op := Less.Factory().Operator()

	// test less int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLT: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test less than 0
	cond, err = op.ToMgo("test", uint64(0))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLT: uint64(0)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test less than a negative number
	cond, err = op.ToMgo("test", int32(-1))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLT: int32(-1)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid less type
	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestLessOrEqualMongoCond(t *testing.T) {
	op := LessOrEqual.Factory().Operator()

	// test less or equal int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLTE: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test less or equal than 0
	cond, err = op.ToMgo("test", uint64(0))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLTE: uint64(0)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test less or equal than a negative number
	cond, err = op.ToMgo("test", int32(-1))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLTE: int32(-1)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid less or equal type
	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestGreaterMongoCond(t *testing.T) {
	op := Greater.Factory().Operator()

	// test greater int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGT: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test greater than 0
	cond, err = op.ToMgo("test", uint64(0))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGT: uint64(0)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test greater than a negative number
	cond, err = op.ToMgo("test", int32(-1))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGT: int32(-1)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid greater type
	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestGreaterOrEqualMongoCond(t *testing.T) {
	op := GreaterOrEqual.Factory().Operator()

	// test greater or equal int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGTE: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test greater or equal than 0
	cond, err = op.ToMgo("test", uint64(0))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGTE: uint64(0)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test greater or equal than a negative number
	cond, err = op.ToMgo("test", int32(-1))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGTE: int32(-1)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid greater or equal type
	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestDatetimeLessMongoCond(t *testing.T) {
	op := DatetimeLess.Factory().Operator()

	// test datetime less time type
	now := time.Now()
	cond, err := op.ToMgo("test", now)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLT: now}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime less timestamp type
	nowUnixTime := time.Unix(now.Unix(), 0)
	cond, err = op.ToMgo("test", now.Unix())
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLT: nowUnixTime}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime less time string
	nowStr := now.Format(common.TimeTransferModel)
	nowStrTime, _ := time.ParseInLocation(common.TimeTransferModel, nowStr, time.Local)
	cond, err = op.ToMgo("test", nowStr)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBLT: nowStrTime.UTC()}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid datetime less type
	cond, err = op.ToMgo("test", "2022")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestDatetimeLessOrEqualMongoCond(t *testing.T) {
	op := DatetimeLessOrEqual.Factory().Operator()

	// test datetime less or equal time type
	now := time.Now()
	cond, err := op.ToMgo("test", now)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLTE: now}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime less or equal timestamp type
	nowUnixTime := time.Unix(now.Unix(), 0)
	cond, err = op.ToMgo("test", now.Unix())
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLTE: nowUnixTime}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime less or equal time string
	nowStr := now.Format(common.TimeTransferModel)
	nowStrTime, _ := time.ParseInLocation(common.TimeTransferModel, nowStr, time.Local)
	cond, err = op.ToMgo("test", nowStr)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBLTE: nowStrTime.UTC()}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid datetime less or equal type
	cond, err = op.ToMgo("test", "2022")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestDatetimeGreaterMongoCond(t *testing.T) {
	op := DatetimeGreater.Factory().Operator()

	// test datetime greater time type
	now := time.Now()
	cond, err := op.ToMgo("test", now)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGT: now}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime greater timestamp type
	nowUnixTime := time.Unix(now.Unix(), 0)
	cond, err = op.ToMgo("test", now.Unix())
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGT: nowUnixTime}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime greater time string
	nowStr := now.Format(common.TimeTransferModel)
	nowStrTime, _ := time.ParseInLocation(common.TimeTransferModel, nowStr, time.Local)
	cond, err = op.ToMgo("test", nowStr)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBGT: nowStrTime.UTC()}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid datetime greater type
	cond, err = op.ToMgo("test", "2022")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestDatetimeGreaterOrEqualMongoCond(t *testing.T) {
	op := DatetimeGreaterOrEqual.Factory().Operator()

	// test datetime greater or equal time type
	now := time.Now()
	cond, err := op.ToMgo("test", now)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGTE: now}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime greater or equal timestamp type
	nowUnixTime := time.Unix(now.Unix(), 0)
	cond, err = op.ToMgo("test", now.Unix())
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBGTE: nowUnixTime}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime greater or equal time string
	nowStr := now.Format(common.TimeTransferModel)
	nowStrTime, _ := time.ParseInLocation(common.TimeTransferModel, nowStr, time.Local)
	cond, err = op.ToMgo("test", nowStr)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBGTE: nowStrTime.UTC()}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid datetime greater or equal type
	cond, err = op.ToMgo("test", "2022")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestBeginsWithMongoCond(t *testing.T) {
	op := BeginsWith.Factory().Operator()

	// test begins with string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLIKE: "^a"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid begins with type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestBeginsWithInsensitiveMongoCond(t *testing.T) {
	op := BeginsWithInsensitive.Factory().Operator()

	// test begins with insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLIKE: "^a",
		common.BKDBOPTIONS: "i"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid begins with insensitive type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotBeginsWithMongoCond(t *testing.T) {
	op := NotBeginsWith.Factory().Operator()

	// test not begins with string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBNot: map[string]interface{}{common.BKDBLIKE: "^a"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not begins with type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotBeginsWithInsensitiveMongoCond(t *testing.T) {
	op := NotBeginsWithInsensitive.Factory().Operator()

	// test not begins with insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBNot: map[string]interface{}{common.BKDBLIKE: "^a", common.BKDBOPTIONS: "i"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not begins with insensitive type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestContainsMongoCond(t *testing.T) {
	op := Contains.Factory().Operator()

	// test contains string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLIKE: "a",
		common.BKDBOPTIONS: "i"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid contains type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestContainsSensitiveMongoCond(t *testing.T) {
	op := ContainsSensitive.Factory().Operator()

	// test contains string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLIKE: "a"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid contains type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestContainsInsensitiveMongoCond(t *testing.T) {
	op := ContainsInsensitive.Factory().Operator()

	// test contains insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLIKE: "a",
		common.BKDBOPTIONS: "i"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid contains insensitive type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotContainsMongoCond(t *testing.T) {
	op := NotContains.Factory().Operator()

	// test not contains string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBNot: map[string]interface{}{common.BKDBLIKE: "a"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not contains type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotContainsInsensitiveMongoCond(t *testing.T) {
	op := NotContainsInsensitive.Factory().Operator()

	// test not contains insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBNot: map[string]interface{}{common.BKDBLIKE: "a", common.BKDBOPTIONS: "i"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not contains insensitive type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestEndsWithMongoCond(t *testing.T) {
	op := EndsWith.Factory().Operator()

	// test ends with string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLIKE: "a$"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid ends with type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestEndsWithInsensitiveMongoCond(t *testing.T) {
	op := EndsWithInsensitive.Factory().Operator()

	// test ends with insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBLIKE: "a$",
		common.BKDBOPTIONS: "i"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid ends with insensitive type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotEndsWithMongoCond(t *testing.T) {
	op := NotEndsWith.Factory().Operator()

	// test not ends with string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBNot: map[string]interface{}{common.BKDBLIKE: "a$"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not ends with type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotEndsWithInsensitiveMongoCond(t *testing.T) {
	op := NotEndsWithInsensitive.Factory().Operator()

	// test not ends with insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBNot: map[string]interface{}{common.BKDBLIKE: "a$", common.BKDBOPTIONS: "i"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not ends with insensitive type
	cond, err = op.ToMgo("test", "")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestIsEmptyMongoCond(t *testing.T) {
	op := IsEmpty.Factory().Operator()

	// test is empty cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBSize: 0}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid is empty field
	cond, err = op.ToMgo("", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestIsNotEmptyMongoCond(t *testing.T) {
	op := IsNotEmpty.Factory().Operator()

	// test is not empty cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		common.BKDBSize: map[string]interface{}{common.BKDBGT: 0}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid is not empty field
	cond, err = op.ToMgo("", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestSizeMongoCond(t *testing.T) {
	op := Size.Factory().Operator()

	// test size int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBSize: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test size equal to 0
	cond, err = op.ToMgo("test", uint64(0))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBSize: uint64(0)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid size type
	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", int32(-1))
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestIsNullMongoCond(t *testing.T) {
	op := IsNull.Factory().Operator()

	// test is null cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBEQ: nil}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid is null field
	cond, err = op.ToMgo("", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestIsNotNullMongoCond(t *testing.T) {
	op := IsNotNull.Factory().Operator()

	// test is not null cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBNE: nil}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid is not null field
	cond, err = op.ToMgo("", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestExistMongoCond(t *testing.T) {
	op := Exist.Factory().Operator()

	// test exist cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBExists: true}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid exist field
	cond, err = op.ToMgo("", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestNotExistMongoCond(t *testing.T) {
	op := NotExist.Factory().Operator()

	// test not exist cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{common.BKDBExists: false}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid not exist field
	cond, err = op.ToMgo("", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestFilterObjectMongoCond(t *testing.T) {
	op := FilterObject.Factory().Operator()

	// test filter object normal type
	cond, err := op.ToMgo("obj", exampleRule)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	expectCond := map[string]interface{}{
		common.BKDBAND: []map[string]interface{}{
			{
				"obj.test": map[string]interface{}{common.BKDBEQ: 1},
			}, {
				common.BKDBOR: []map[string]interface{}{
					{
						common.BKDBAND: []map[string]interface{}{{
							"obj.test1.test2": map[string]interface{}{common.BKDBIN: []string{"b", "c"}},
						}},
					}, {
						"obj.test3": map[string]interface{}{common.BKDBLT: time.Unix(1, 0)},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(cond, expectCond) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid filter object type
	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []interface{}{1, "a"})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", CombinedRule{
		Condition: And,
		Rules: []RuleFactory{
			&AtomRule{
				Field:    "test",
				Operator: Equal.Factory(),
				Value:    1,
			},
			&CombinedRule{
				Condition: Or,
				Rules: []RuleFactory{
					&AtomRule{
						Field:    "test1",
						Operator: In.Factory(),
						Value:    "a",
					},
				},
			},
		},
	})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}

func TestFilterArrayMongoCond(t *testing.T) {
	op := FilterArray.Factory().Operator()

	// test filter array normal type
	cond, err := op.ToMgo("arr", &CombinedRule{
		Condition: And,
		Rules: []RuleFactory{
			&AtomRule{
				Field:    FilterArrayElement,
				Operator: Equal.Factory(),
				Value:    1,
			},
			&CombinedRule{
				Condition: Or,
				Rules: []RuleFactory{
					&AtomRule{
						Field:    "1",
						Operator: NotEqual.Factory(),
						Value:    "a",
					},
					&AtomRule{
						Field:    "4",
						Operator: In.Factory(),
						Value:    []string{"b", "c"},
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	expectCond := map[string]interface{}{
		common.BKDBAND: []map[string]interface{}{{
			"arr": map[string]interface{}{common.BKDBEQ: 1},
		}, {
			common.BKDBOR: []map[string]interface{}{{
				"arr.1": map[string]interface{}{common.BKDBNE: "a"},
			}, {
				"arr.4": map[string]interface{}{common.BKDBIN: []string{"b", "c"}},
			}},
		}},
	}

	if !reflect.DeepEqual(cond, expectCond) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test invalid filter array type
	cond, err = op.ToMgo("test", 1)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", "a")
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", false)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", map[string]interface{}{"test1": 1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", struct{}{})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", nil)
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []int64{1})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", []interface{}{1, "a"})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}

	cond, err = op.ToMgo("test", CombinedRule{
		Condition: And,
		Rules: []RuleFactory{
			&AtomRule{
				Field:    FilterArrayElement,
				Operator: Equal.Factory(),
				Value:    1,
			},
			&CombinedRule{
				Condition: Or,
				Rules: []RuleFactory{
					&AtomRule{
						Field:    "-1",
						Operator: In.Factory(),
						Value:    "a",
					},
				},
			},
		},
	})
	if err == nil {
		t.Errorf("to mongo should return error, but get cond: %+v", cond)
		return
	}
}
