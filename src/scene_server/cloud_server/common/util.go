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

package common

import (
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
)

// 将不同云厂商的实例状态转为统一的实例状态
func CovertInstState(instState string) string {
	switch strings.ToLower(instState) {
	case "starting", "pending", "rebooting":
		return common.BKCloudHostStatusStarting
	case "running":
		return common.BKCloudHostStatusRunning
	case "stopping", "shutting-down", "terminating":
		return common.BKCloudHostStatusStopping
	case "stopped", "shutdown", "terminated":
		return common.BKCloudHostStatusStopped
	default:
		blog.Infof("convert to unknow state, the origin instState:%s", instState)
		return common.BKCloudHostStatusUnknown
	}
	return instState
}

// NewHeader 创建云资源同步需要的header
func NewHeader() http.Header {
	header := make(http.Header)
	header.Add(common.BKHTTPOwnerID, common.BKSuperOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.BKCloudSyncUser)
	header.Add(common.BKHTTPLanguage, "cn")
	header.Add(common.BKHTTPCCRequestID, util.GenerateRID())
	header.Add("Content-Type", "application/json")
	return header
}

// NewKit 创建新的Kit
func NewKit() *rest.Kit {
	header := NewHeader()
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

// NewReadwKit 创建专用于读操作的kit
// SupplierAccount为superadmin
func NewReadwKit() *rest.Kit {
	return NewKit()
}

// NewWriteKit 创建专用于写操作的kit
// SupplierAccount为与当前环境匹配的开发商
func NewWriteKit(supplierAccount string) *rest.Kit {
	kit := NewKit()
	kit.Header.Set(common.BKHTTPOwnerID, supplierAccount)
	kit.SupplierAccount = supplierAccount
	return kit
}

// CopyHeaderTxnInfo copy transaction info from src to dst
func CopyHeaderTxnInfo(src http.Header, dst http.Header) {
	dst.Set(common.TransactionIdHeader, src.Get(common.TransactionIdHeader))
	dst.Set(common.TransactionTimeoutHeader, src.Get(common.TransactionTimeoutHeader))
}

// DelHeaderTxnInfo delete transaction info from header
func DelHeaderTxnInfo(header http.Header) {
	header.Del(common.TransactionIdHeader)
	header.Del(common.TransactionTimeoutHeader)
}
