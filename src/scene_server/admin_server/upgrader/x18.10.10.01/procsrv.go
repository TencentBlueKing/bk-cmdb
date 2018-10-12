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
package x18_10_10_01

import (
	"fmt"
	"strings"

	"gopkg.in/mgo.v2"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func addProcOpTaskTable(db storage.DI, conf *upgrader.Config) error {
	tableName := common.BKTableNameProcOperateTask
	exists, err := db.HasTable(tableName)
	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	indexs := []storage.Index{
		storage.Index{Name: "idx_taskID_gseTaskID", Columns: []string{common.BKTaskIDField, common.BKGseOpTaskIDField}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	for index := range indexs {
		if err = db.Index(tableName, &indexs[index]); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	return nil
}
func addProcInstanceModelTable(db storage.DI, conf *upgrader.Config) error {
	tableName := common.BKTableNameProcInstanceModel
	exists, err := db.HasTable(tableName)
	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	indexs := []storage.Index{
		storage.Index{Name: "idx_bkBizID_bkSetID_bkModuleID_bkHostInstanceID", Columns: []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField, "bk_host_instance_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "idx_bkBizID_bkHostID", Columns: []string{common.BKAppIDField, common.BKHostIDField}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "idx_bkBizID_bkProcessID", Columns: []string{common.BKAppIDField, common.BKProcessIDField}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	for index := range indexs {
		if err = db.Index(tableName, &indexs[index]); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	return nil
}
func addProcInstanceDetailTable(db storage.DI, conf *upgrader.Config) error {
	tableName := common.BKTableNameProcInstaceDetail
	exists, err := db.HasTable(tableName)
	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	indexs := []storage.Index{
		storage.Index{Name: "idx_bkBizID_bkModuleID_bkProcessID", Columns: []string{common.BKAppIDField, common.BKModuleIDField, common.BKProcessIDField}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "idx_bkBizID_status", Columns: []string{common.BKAppIDField, common.BKStatusField}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "idx_bkBizID_bkHostID", Columns: []string{common.BKAppIDField, common.BKHostIDField}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	for index := range indexs {
		if err = db.Index(tableName, &indexs[index]); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	return nil
}
func addProcFreshInstance(db storage.DI, conf *upgrader.Config) error {
	if "" != conf.CCApiSrvAddr {
		tableName := common.BKTableNameSubscription
		sID, err := db.GetIncID(tableName)
		if nil != err {
			return err
		}
		SubscriptionName := "process instance refresh [incorrect deletion]"
		cnt, err := db.GetCntByCondition(tableName, mapstr.MapStr{common.BKSubscriptionNameField: SubscriptionName, common.BKOperatorField: conf.User})
		if nil != err {
			return err
		}
		if 0 < cnt {
			return nil
		}
		subscription := metadata.Subscription{
			SubscriptionID:   sID,
			SubscriptionName: SubscriptionName,
			SystemName:       "cmdb",
			CallbackURL:      fmt.Sprintf("http://%s/api/v3/proc/process/refresh/hostinstnum", strings.Trim(conf.CCApiSrvAddr, "/")),
			ConfirmMode:      metadata.ConfirmmodeHttpstatus,
			ConfirmPattern:   "200",
			TimeOut:          120,
			SubscriptionForm: "hostupdate,moduletransfer,update,processmodule,processupdate",
			OwnerID:          common.BKDefaultOwnerID,
			Operator:         conf.User,
		}
		_, err = db.Insert(tableName, subscription)
		return err
	}
	return nil
}
