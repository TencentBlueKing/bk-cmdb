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

	iamcli "configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/logics"
)

const (
	// 同步周期最小值
	syncIAMPeriodMinutesMin = 1
	// 同步周期默认值
	syncIAMPeriodMinutesDefault = 5
)

// syncor used to sync iam
type syncor struct {
	// 同步周期
	SyncIAMPeriodMinutes int
}

func NewSyncor() *syncor {
	return &syncor{}
}

// SetSyncIAMPeriod set the sync period
func (s *syncor) SetSyncIAMPeriod(periodMinutes int) {
	s.SyncIAMPeriodMinutes = periodMinutes
	if s.SyncIAMPeriodMinutes < syncIAMPeriodMinutesMin {
		s.SyncIAMPeriodMinutes = syncIAMPeriodMinutesDefault
	}
	blog.Infof("sync iam period is %d minutes", s.SyncIAMPeriodMinutes)
}

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
func (s *syncor) SyncIAM(iamCli *iamcli.IAM, lgc *logics.Logics) {
	if !auth.EnableAuthorize() {
		return
	}
	time.Sleep(time.Minute)

	for {
		// new kit with a different rid, header
		kit := newKit()

		// only master can run it
		if !lgc.ServiceManageInterface.IsMaster() {
			blog.V(4).Infof("it is not master, skip sync iam, rid: %s", kit.Rid)
			time.Sleep(time.Minute)
			continue
		}

		blog.Infof("start sync iam, rid: %s", kit.Rid)

		// get all custom objects (without mainline objects) in cmdb
		objects, err := lgc.GetCustomObjects(kit.Ctx, kit.Header)
		if err != nil {
			blog.Errorf("sync iam failed, get custom objects err: %s ,rid: %s", err, kit.Rid)
			time.Sleep(time.Duration(s.SyncIAMPeriodMinutes) * time.Minute)
			continue
		}

		if err := iamCli.SyncIAMSysInstances(kit.Ctx, objects); err != nil {
			blog.Errorf("sync iam failed, sync iam system instances err: %s ,rid: %s", err, kit.Rid)
			time.Sleep(time.Duration(s.SyncIAMPeriodMinutes) * time.Minute)
			continue
		}

		blog.Infof("finish sync iam successfully, rid:%s", kit.Rid)
		time.Sleep(time.Duration(s.SyncIAMPeriodMinutes) * time.Minute)
	}
}
