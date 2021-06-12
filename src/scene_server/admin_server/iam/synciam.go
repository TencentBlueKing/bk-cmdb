/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package iam

import (
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/service"
	iamcli "configcenter/src/ac/iam"
)

const (
	// 同步周期最小值
	SyncIAMPeriodMinutesMin = 1
	// 同步周期默认值
	SyncIAMPeriodMinutesDefault = 5
)

// 同步周期
var SyncIAMPeriodMinutes int

// newHeader 创建IAM同步需要的header
func newHeader() http.Header {
	header := make(http.Header)
	header.Add(common.BKHTTPOwnerID, common.BKSuperOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.BKIAMSyncUser)
	header.Add(common.BKHTTPLanguage, "cn")
	header.Add(common.BKHTTPCCRequestID, util.GenerateRID())
	header.Add("Content-Type", "application/json")
	return header
}

// newKit 创建新的Kit
func newKit() *rest.Kit {
	header := newHeader()
	if header.Get(common.BKHTTPCCRequestID) == "" {
		header.Set(common.BKHTTPCCRequestID, util.GenerateRID())
	}
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)
	user := util.GetUser(header)
	supplierAccount := util.GetOwnerID(header)
	defaultCCError := util.GetDefaultCCError(header)

	return &rest.Kit{
		Rid:             rid,
		Header:          header,
		Ctx:             ctx,
		CCError:         defaultCCError,
		User:            user,
		SupplierAccount: supplierAccount,
	}
}

// SyncIAM sync the system instances resource between CMDB and IAM
func SyncIAM(srv *service.Service, iamCli *iamcli.IAM) {
	if !auth.EnableAuthorize() {
		return
	}

	// delay some time to sync at beginning, leave some time for upgrade program
	time.Sleep(time.Minute * 10)

	for {
		blog.Infof("start sync iam")

		// only master can run it
		if !srv.ServiceManageInterface.IsMaster() {
			blog.Infof("it is not master, skip sync iam")
			time.Sleep(20 * time.Second)
			continue
		}

		// new kit with a different rid, header
		kit := newKit()

		// get all custom objects (without mainline objects) in cmdb
		objects, err := srv.GetCustomObjects(kit.Header)
		if err != nil {
			blog.Errorf("sync iam failed, get custom objects err: %s ,rid: %s", err, kit.Rid)
			time.Sleep(time.Duration(SyncIAMPeriodMinutes) * time.Minute)
			continue
		}

		if err := iamCli.SyncIAMSysInstances(kit.Ctx, objects); err != nil {
			blog.Errorf("sync iam failed, sync iam system instances err: %s ,rid: %s", err, kit.Rid)
			time.Sleep(time.Duration(SyncIAMPeriodMinutes) * time.Minute)
			continue
		}

		blog.Infof("finish sync iam successfully")
		time.Sleep(time.Duration(SyncIAMPeriodMinutes) * time.Minute)
	}
}
