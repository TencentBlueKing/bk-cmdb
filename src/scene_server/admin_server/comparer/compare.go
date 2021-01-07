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

package main

import (
	"flag"
	"log"
	"time"

	"gopkg.in/mgo.v2"

	"configcenter/src/common/util"
)

var src = flag.String("source", "127.0.0.1:27017", "")
var srcUser = flag.String("source_user", "cc", "")
var srcPWD = flag.String("source_passwd", "cc", "")
var srcDB = flag.String("source_db", "cmdb_bak", "")
var target = flag.String("target", "127.0.0.1:27017", "")
var targetDB = flag.String("target_db", "cmdb", "")
var targetUser = flag.String("target_user", "cc", "")
var targetPWD = flag.String("target_passwd", "cc", "")

func main() {
	flag.Parse()
	srccli, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{*src},
		Direct:    false,
		Timeout:   time.Second * 5,
		Database:  *srcDB,
		Source:    "",
		Username:  *srcUser,
		Password:  *srcPWD,
		PoolLimit: 4096,
		Mechanism: "",
	})
	if err != nil {
		panic(err)
	}
	tarcli, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{*target},
		Direct:    false,
		Timeout:   time.Second * 5,
		Database:  *targetDB,
		Source:    "",
		Username:  *targetUser,
		Password:  *targetPWD,
		PoolLimit: 4096,
		Mechanism: "",
	})
	if err != nil {
		panic(err)
	}

	if err := processCompare(srccli.DB(*srcDB), tarcli.DB(*targetDB)); err != nil {
		panic(err)
	}

	log.Println("congratulation, compare success")
}

var empty = map[string]interface{}{}

func processCompare(srccli, tarcli *mgo.Database) error {
	tablenames, err := srccli.CollectionNames()
	assertNotErr(err)
	for _, tablename := range tablenames {
		if ignoreTable[tablename] {
			continue
		}
		srcDatas := []map[string]interface{}{}
		err = srccli.C(tablename).Find(empty).All(&srcDatas)
		assertNotErr(err)
		for _, srcData := range srcDatas {
			tableKey := tableKeys(tablename)
			condition := util.CopyMap(srcData, tableKey.keys, []string{"_id"})
			tarData := map[string]interface{}{}
			err = tarcli.C(tablename).Find(condition).One(&tarData)
			if err != nil {
				log.Fatalf("tablename %s, condition %v, error: %s", tablename, condition, err.Error())
			}
			compare(tablename, srcData, tarData, tableKey.ignores)
		}
	}
	return nil
}

var ignoreTable = map[string]bool{
	"cc_OperationLog": true,
	"cc_System":       true,
}

func compare(tablename string, srcData, tarData map[string]interface{}, ignores []string) {
	ignore := map[string]bool{}
	for _, key := range ignores {
		ignore[key] = true
	}
	for key := range srcData {
		if ignore[key] {
			continue
		}
		equal := true
		switch val := srcData[key].(type) {
		case int, int64:
			if toInt64(tarData[key]) != toInt64(val) {
				equal = false
			}
		case []interface{}:
		case map[string]interface{}:

		default:
			if tarData[key] != srcData[key] {
				equal = false
			}
		}
		if !equal {
			log.Fatalf("not equal!! tablename: %s, key %s , expect %#v, actual %#v", tablename, key, tarData[key], srcData[key])
		}
	}
}

func toInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int:
		return int64(v)
	case int64:
		return int64(v)
	}
	return 0
}

type tableKey struct {
	keys    []string
	ignores []string
}

var tableKeysCache = map[string]*tableKey{
	"cc_ApplicationBase":   &tableKey{keys: []string{"bk_biz_name"}, ignores: []string{"bk_biz_id"}},
	"cc_ModuleBase":        &tableKey{keys: []string{"bk_module_name"}, ignores: []string{"bk_module_id", "bk_biz_id", "bk_set_id", "bk_parent_id"}},
	"cc_ObjAttDes":         &tableKey{keys: []string{"bk_obj_id", "bk_property_id"}, ignores: []string{"id"}},
	"cc_ObjClassification": &tableKey{keys: []string{"bk_classification_id"}, ignores: []string{"id"}},
	"cc_ObjDes":            &tableKey{keys: []string{"bk_obj_id"}, ignores: []string{"id"}},
	"cc_PlatBase":          &tableKey{keys: []string{"bk_cloud_name"}, ignores: []string{}},
	"cc_Proc2Module":       &tableKey{keys: []string{"bk_module_name", "bk_process_id", "bk_biz_id"}, ignores: []string{}},
	"cc_Process":           &tableKey{keys: []string{"bk_process_name"}, ignores: []string{}},
	"cc_PropertyGroup":     &tableKey{keys: []string{"bk_obj_id", "bk_group_id"}, ignores: []string{}},
	"cc_SetBase":           &tableKey{keys: []string{"bk_set_name", "bk_biz_id"}, ignores: []string{"bk_set_id"}},
	"cc_OperationLog":      &tableKey{keys: []string{"op_type", "inst_id"}, ignores: []string{"op_time"}},
	"cc_ObjAsst":           &tableKey{keys: []string{"bk_obj_id", "bk_object_att_id", "bk_asst_obj_id"}, ignores: []string{"id"}},
}

func tableKeys(tablename string) *tableKey {
	if _, ok := tableKeysCache[tablename]; ok {
		tableKeysCache[tablename].ignores = append(tableKeysCache[tablename].ignores, "create_time", "last_time", "_id")
		return tableKeysCache[tablename]
	}
	return &tableKey{keys: []string{}, ignores: []string{"create_time", "last_time"}}
}

func assertNotErr(err error) {
	if err != nil {
		panic(err)
	}
}
