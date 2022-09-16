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

package metadata

import (
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/errors"

	"github.com/gin-gonic/gin"
)

// LoginUserInfoOwnerUinList TODO
type LoginUserInfoOwnerUinList struct {
	OwnerID   string `json:"id"`
	OwnerName string `json:"name"`
	Role      int64  `json:"role"`
}

// LoginUserInfo TODO
type LoginUserInfo struct {
	UserName      string                      `json:"username"`
	ChName        string                      `json:"chname"`
	Phone         string                      `json:"phone"`
	Email         string                      `json:"email"`
	Role          string                      `json:"-"`
	BkToken       string                      `json:"bk_token"`
	OnwerUin      string                      `json:"current_supplier"`
	OwnerUinArr   []LoginUserInfoOwnerUinList `json:"supplier_list"` // user all owner uin
	IsOwner       bool                        `json:"-"`             // is master
	Extra         map[string]interface{}      `json:"extra"`         // custom information
	Language      string                      `json:"-"`
	AvatarUrl     string                      `json:"avatar_url"`
	MultiSupplier bool                        `json:"multi_supplier"`
}

// LoginPluginInfo TODO
type LoginPluginInfo struct {
	Name       string // plugin info
	Version    string // In what version is used
	HandleFunc LoginUserPluginInerface
}

// LoginUserPluginParams TODO
type LoginUserPluginParams struct {
	Url          string
	IsMultiOwner bool
	Cookie       []*http.Cookie // Reserved word, not used now
	Header       http.Header    // Reserved word, not used now
}

// LoginUserPluginInerface TODO
type LoginUserPluginInerface interface {
	LoginUser(c *gin.Context, config map[string]string, isMultiOwner bool) (user *LoginUserInfo, loginSucc bool)
	GetLoginUrl(c *gin.Context, config map[string]string, input *LogoutRequestParams) string
	GetUserList(c *gin.Context, config map[string]string) ([]*LoginSystemUserInfo, *errors.RawErrorInfo)
}

// LoginSystemUserInfo TODO
type LoginSystemUserInfo struct {
	CnName string `json:"chinese_name"`
	EnName string `json:"english_name"`
}

// LonginSystemUserListResult TODO
type LonginSystemUserListResult struct {
	BaseResp `json:",inline"`
	Data     []*LoginSystemUserInfo `json:"data"`
}

// DepartmentResult TODO
type DepartmentResult struct {
	BaseResp `json:",inline"`
	Data     *DepartmentData `json:"data"`
}

// DepartmentProfileResult TODO
type DepartmentProfileResult struct {
	BaseResp `json:",inline"`
	Data     *DepartmentProfileData `json:"data"`
}

// LoginUserInfoDetail TODO
type LoginUserInfoDetail struct {
	UserName      string                      `json:"username"`
	ChName        string                      `json:"chname"`
	OnwerUin      string                      `json:"current_supplier"`
	OwnerUinArr   []LoginUserInfoOwnerUinList `json:"supplier_list"` // user all owner uin
	AvatarUrl     string                      `json:"avatar_url"`
	MultiSupplier bool                        `json:"multi_supplier"`
}

// LoginUserInfoResult TODO
type LoginUserInfoResult struct {
	BaseResp `json:",inline"`
	Data     LoginUserInfoDetail `json:"data"`
}

// LoginChangeSupplierResult TODO
type LoginChangeSupplierResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		ID string `json:"bk_supplier_account"`
	} `json:"data"`
}

// LogoutResult TODO
type LogoutResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		LogoutURL string `json:"url"`
	} `json:"data"`
}

// LogoutRequestParams TODO
type LogoutRequestParams struct {
	HTTPScheme string `json:"http_scheme"`
}

// ExcelAssociationOperate TODO
type ExcelAssociationOperate int

const (
	_ ExcelAssociationOperate = iota
	// ExcelAssociationOperateError TODO
	ExcelAssociationOperateError
	// ExcelAssociationOperateAdd TODO
	ExcelAssociationOperateAdd
	// ExcelAssociationOperateDelete TODO
	// ExcelAssociationOperateUpdate
	ExcelAssociationOperateDelete
)

// ExcelAssociation TODO
type ExcelAssociation struct {
	ObjectAsstID string                  `json:"bk_obj_asst_id"`
	Operate      ExcelAssociationOperate `json:"operate"`
	SrcPrimary   string                  `json:"src_primary_key"`
	DstPrimary   string                  `json:"dst_primary_key"`
}

// ObjectAsstIDStatisticsInfo TODO
type ObjectAsstIDStatisticsInfo struct {
	Create int64 `json:"create"`
	Delete int64 `json:"delete"`
	Total  int64 `json:"total"`
}

// BatchExportObject param of bacth export object
type BatchExportObject struct {
	ObjectID       []int64 `json:"object_id"`
	ExcludedAsstID []int64 `json:"excluded_asst_id"`
	Password       string  `json:"password"`
	Expiration     int64   `json:"expiration"`
	FileName       string  `json:"file_name"`
}

// ListObjectTopoResponse list object with it's topo info response
type ListObjectTopoResponse struct {
	BaseResp `json:",inline"`
	Data     *TotalObjectInfo `json:"data"`
}

// ZipFileAnalysis analysis zip file
type ZipFileAnalysis struct {
	Password string `json:"password"`
}

// AnalysisResult result of analysis zip file
type AnalysisResult struct {
	BaseResp `json:",inline"`
	Data     BatchImportObject `json:"data"`
}

// BatchImportObject param of batch import object
type BatchImportObject struct {
	Object []YamlObject      `json:"import_object"`
	Asst   []AssociationKind `json:"import_asst"`
}

// YamlHeader yaml's common header
type YamlHeader struct {
	ExpireTime int64 `json:"expire_time" yaml:"expire_time"`
	CreateTime int64 `json:"create_time" yaml:"create_time"`
}

// Validate validate yaml common field
func (o *YamlHeader) Validate() errors.RawErrorInfo {
	if o.ExpireTime == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrWebVerifyYamlFail,
			Args:    []interface{}{"expire_time not found"},
		}
	}

	if o.CreateTime == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrWebVerifyYamlFail,
			Args:    []interface{}{"create_time not found"},
		}
	}

	if (o.ExpireTime != o.CreateTime && o.ExpireTime < time.Now().Local().UnixNano()) || o.ExpireTime < o.CreateTime {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrWebVerifyYamlFail,
			Args:    []interface{}{"expire time incorrect"},
		}
	}

	return errors.RawErrorInfo{}
}
