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

type CMDBConfig struct {
	Address  *string
	User     *string
	Password *string
	Database *string
}

var srcConfig CMDBConfig
var dstConfig CMDBConfig

func init() {
	srcConfig.Address = flag.String("source", "127.0.0.1:27017", "")
	srcConfig.User = flag.String("source_user", "cc", "")
	srcConfig.Password = flag.String("source_passwd", "cc", "")
	srcConfig.Database = flag.String("source_db", "cmdb_bak", "")

	dstConfig.Address = flag.String("target", "127.0.0.1:27017", "")
	dstConfig.User = flag.String("target_user", "cc", "")
	dstConfig.Password = flag.String("target_passwd", "cc", "")
	dstConfig.Database = flag.String("target_db", "cmdb", "")
}

func main() {
	flag.Parse()
	srcCli, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{*srcConfig.Address},
		Direct:    false,
		Timeout:   time.Second * 5,
		Database:  *srcConfig.Database,
		Source:    "",
		Username:  *srcConfig.User,
		Password:  *srcConfig.Password,
		PoolLimit: 4096,
		Mechanism: "",
	})
	if err != nil {
		panic(err)
	}
	dstCli, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{*dstConfig.Address},
		Direct:    false,
		Timeout:   time.Second * 5,
		Database:  *dstConfig.Database,
		Source:    "",
		Username:  *dstConfig.User,
		Password:  *dstConfig.Password,
		PoolLimit: 4096,
		Mechanism: "",
	})
	if err != nil {
		panic(err)
	}

	if err := processCompare(srcCli.DB(*srcConfig.Database), dstCli.DB(*dstConfig.Database)); err != nil {
		panic(err)
	}

	log.Println("congratulation, compare success")
}

var empty = map[string]interface{}{}

func processCompare(srcCli, tarCli *mgo.Database) error {
	tableNames, err := srcCli.CollectionNames()
	assertNotErr(err)
	for _, tableName := range tableNames {
		if ignoreTable[tableName] {
			continue
		}
		srcDatas := make([]map[string]interface{}, 0)
		err = srcCli.C(tableName).Find(empty).All(&srcDatas)
		assertNotErr(err)
		for _, srcData := range srcDatas {
			tableKey := tableKeys(tableName)
			condition := util.CopyMap(srcData, tableKey.keys, []string{"_id"})
			tarData := map[string]interface{}{}
			err = tarCli.C(tableName).Find(condition).One(&tarData)
			if err != nil {
				log.Fatalf("tablename %s, condition %v, error: %s", tableName, condition, err.Error())
			}
			compare(tableName, srcData, tarData, tableKey.ignores)
		}
	}
	return nil
}

var ignoreTable = map[string]bool{
	"cc_OperationLog": true,
	"cc_System":       true,
}

func compare(tableName string, srcData, tarData map[string]interface{}, ignores []string) {
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
			log.Fatalf("not equal!! tablename: %s, key %s , expect %#v, actual %#v", tableName, key, tarData[key], srcData[key])
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
	"cc_ApplicationBase":   {keys: []string{"bk_biz_name"}, ignores: []string{"bk_biz_id"}},
	"cc_ModuleBase":        {keys: []string{"bk_module_name"}, ignores: []string{"bk_module_id", "bk_biz_id", "bk_set_id", "bk_parent_id"}},
	"cc_ObjAttDes":         {keys: []string{"bk_obj_id", "bk_property_id"}, ignores: []string{"id"}},
	"cc_ObjClassification": {keys: []string{"bk_classification_id"}, ignores: []string{"id"}},
	"cc_ObjDes":            {keys: []string{"bk_obj_id"}, ignores: []string{"id"}},
	"cc_PlatBase":          {keys: []string{"bk_cloud_name"}, ignores: []string{}},
	"cc_Proc2Module":       {keys: []string{"bk_module_name", "bk_process_id", "bk_biz_id"}, ignores: []string{}},
	"cc_Process":           {keys: []string{"bk_process_name"}, ignores: []string{}},
	"cc_PropertyGroup":     {keys: []string{"bk_obj_id", "bk_group_id"}, ignores: []string{}},
	"cc_SetBase":           {keys: []string{"bk_set_name", "bk_biz_id"}, ignores: []string{"bk_set_id"}},
	"cc_OperationLog":      {keys: []string{"op_type", "inst_id"}, ignores: []string{"op_time"}},
	"cc_ObjAsst":           {keys: []string{"bk_obj_id", "bk_object_att_id", "bk_asst_obj_id"}, ignores: []string{"id"}},
}

func tableKeys(tableName string) *tableKey {
	if _, ok := tableKeysCache[tableName]; ok {
		tableKeysCache[tableName].ignores = append(tableKeysCache[tableName].ignores, "create_time", "last_time", "_id")
		return tableKeysCache[tableName]
	}
	return &tableKey{keys: []string{}, ignores: []string{"create_time", "last_time"}}
}

func assertNotErr(err error) {
	if err != nil {
		panic(err)
	}
}
