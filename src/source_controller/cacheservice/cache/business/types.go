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

package business

import (
	"time"
)

type cacheCollection struct {
	business    *business
	set         *moduleSet
	module      *moduleSet
	customLevel *customLevel
}

type forUpsertCache struct {
	instID   int64
	parentID int64
	name     string
	doc      []byte

	// keys to be used
	listKey         string
	listExpireKey   string
	detailKey       string
	detailExpireKey string

	// generate list value and parse list value
	parseListKeyValue func(key string) (instIDint64, parentID int64, instName string, err error)
	genListKeyValue   func(instID int64, parentID int64, instName string) string

	// get the instance name with instance id from mongodb
	getInstName func(instID int64) (name string, err error)
}

type refreshInstance struct {
	// the key to store data.
	mainKey        string
	lockKey        string
	expireKey      string
	expireDuration time.Duration
	// detail is to get the data to be refresh.
	getDetail func(instID int64) (string, error)
}

type refreshList struct {
	// the key to store the list .
	mainKey        string
	lockKey        string
	expireKey      string
	expireDuration time.Duration
	// detail is to get all the list keys to be refresh.
	// for business list, the bizID should be ignored.
	getList func(bizID int64) ([]string, error)
}

type BizBaseInfo struct {
	BusinessID   int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	BusinessName string `json:"bk_biz_name" bson:"bk_biz_name"`
}

type ModuleBaseInfo struct {
	ModuleID   int64  `json:"bk_module_id" bson:"bk_module_id"`
	ModuleName string `json:"bk_module_name" bson:"bk_module_name"`
	SetID      int64  `json:"bk_set_id" bson:"bk_set_id"`
}

type SetBaseInfo struct {
	SetID    int64  `json:"bk_set_id" bson:"bk_set_id"`
	SetName  string `json:"bk_set_name" bson:"bk_set_name"`
	ParentID int64  `json:"bk_parent_id" bson:"bk_parent_id"`
}

type MainlineTopoAssociation struct {
	AssociateTo string `json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
	ObjectID    string `json:"bk_obj_id" bson:"bk_obj_id"`
}

type CustomInstanceBase struct {
	ObjectID     string `json:"bk_obj_id" bson:"bk_obj_id"`
	InstanceID   int64  `json:"bk_inst_id" bson:"bk_inst_id"`
	InstanceName string `json:"bk_inst_name" bson:"bk_inst_name"`
	ParentID     int64  `json:"bk_parent_id" bson:"bk_parent_id"`
}

const mainlineTopologyListDoneKey = "<mainlineTopologyListDoneKey>"
