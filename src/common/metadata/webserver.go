/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package metadata

import (
	"net/http"
	"time"

	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/querybuilder"
	"configcenter/src/web_server/app/options"

	"github.com/gin-gonic/gin"
)

// LoginUserInfoTenantUinList login user tenant uin list
type LoginUserInfoTenantUinList struct {
	TenantID   string `json:"id"`
	TenantName string `json:"name"`
}

// LoginUserInfo login user info
type LoginUserInfo struct {
	UserName     string                       `json:"username"`
	ChName       string                       `json:"chname"`
	BkToken      string                       `json:"bk_token"`
	BkTicket     string                       `json:"bk_ticket"`
	TenantUin    string                       `json:"tenant_id"`
	TenantUinArr []LoginUserInfoTenantUinList `json:"tenant_list"` // user all tenant uin
	Extra        map[string]interface{}       `json:"extra"`       // custom information
	Language     string                       `json:"-"`
	AvatarUrl    string                       `json:"avatar_url"`
	MultiTenant  bool                         `json:"multi_tenant"`
	TimeZone     string                       `json:"time_zone"`
}

// LoginPluginInfo login user plugin info
type LoginPluginInfo struct {
	Name       string // plugin info
	Version    string // In what version is used
	HandleFunc LoginUserPluginInerface
}

// LoginUserPluginParams login user plugin params
type LoginUserPluginParams struct {
	Url           string
	IsMultiTenant bool
	Cookie        []*http.Cookie // Reserved word, not used now
	Header        http.Header    // Reserved word, not used now
}

// LoginUserPluginInerface login user plugin interface
type LoginUserPluginInerface interface {
	LoginUser(c *gin.Context, config options.Config, isMultiTenant bool) (user *LoginUserInfo, loginSucc bool)
	GetLoginUrl(c *gin.Context, config map[string]string, input *LogoutRequestParams) string
	GetUserList(c *gin.Context, options *GetUserListOptions) ([]*LoginSystemUserInfo, *errors.RawErrorInfo)
}

// GetUserListOptions is the get user list options
type GetUserListOptions struct {
	NeedAll   bool
	Usernames []string
}

// LoginSystemUserInfo login system user info
type LoginSystemUserInfo struct {
	CnName string `json:"chinese_name"`
	EnName string `json:"english_name"`
}

// LonginSystemUserListResult login system user list result
type LonginSystemUserListResult struct {
	BaseResp `json:",inline"`
	Data     []*LoginSystemUserInfo `json:"data"`
}

// LoginUserInfoDetail login user info detail
type LoginUserInfoDetail struct {
	UserName     string                       `json:"username"`
	ChName       string                       `json:"chname"`
	TenantUin    string                       `json:"current_tenant"`
	TenantUinArr []LoginUserInfoTenantUinList `json:"tenant_list"` // user all tenant uin
	AvatarUrl    string                       `json:"avatar_url"`
	MultiTenant  bool                         `json:"multi_tenant"`
}

// LoginUserInfoResult login user info result
type LoginUserInfoResult struct {
	BaseResp `json:",inline"`
	Data     LoginUserInfoDetail `json:"data"`
}

// LoginChangeTenantResult login change tenant result
type LoginChangeTenantResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		TenantID string `json:"tenant_id"`
	} `json:"data"`
}

// LogoutResult logout result
type LogoutResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		LogoutURL string `json:"url"`
	} `json:"data"`
}

// LogoutRequestParams logout request params
type LogoutRequestParams struct {
	HTTPScheme string `json:"http_scheme"`
}

// ExcelAssociationOperate excel association operate
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

// ChangelogDetailConfigOption changelog detail's config
type ChangelogDetailConfigOption struct {
	Version string `json:"version"`
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

// ExcelExportHostInput excel export host input
type ExcelExportHostInput struct {
	// 导出的主机字段
	CustomFields []string `json:"export_custom_fields"`
	// 指定需要导出的主机ID, 设置本参数后， ExportCond限定条件无效
	HostIDArr []int64 `json:"bk_host_ids"`
	// 需要导出主机业务id
	AppID int64 `json:"bk_biz_id"`
	// 导出主机查询参数,就是search host 主机参数
	ExportCond HostCommonSearch `json:"export_condition"`

	// 用来限定导出关联关系，map[bk_obj_id]object_unique_id 2021年05月17日
	AssociationCond map[string]int64 `json:"association_condition"`
	// 用来限定当前操作对象导出数据的时候，需要使用的唯一校验关系，
	// 自关联的时候，规定左边对象使用到的唯一索引
	ObjectUniqueID int64 `json:"object_unique_id"`
}

// ExcelImportAddHostInput excel import add host input
type ExcelImportAddHostInput struct {
	ModuleID int64 `json:"bk_module_id"`
	OpType   int64 `json:"op"`
	// 用来限定导出关联关系，map[bk_obj_id]object_unique_id 2021年05月17日
	AssociationCond map[string]int64 `json:"association_condition"`
	// 用来限定当前操作对象导出数据的时候，需要使用的唯一校验关系，
	// 自关联的时候，规定左边对象使用到的唯一索引
	ObjectUniqueID int64 `json:"object_unique_id"`
}

// ExcelImportUpdateHostInput excel import update host input
type ExcelImportUpdateHostInput struct {
	BizID  int64 `json:"bk_biz_id"`
	OpType int64 `json:"op"`
	// 用来限定导出关联关系，map[bk_obj_id]object_unique_id 2021年05月17日
	AssociationCond map[string]int64 `json:"association_condition"`
	// 用来限定当前操作对象导出数据的时候，需要使用的唯一校验关系，
	// 自关联的时候，规定左边对象使用到的唯一索引
	ObjectUniqueID int64 `json:"object_unique_id"`
}

// ListFieldTmplWithObjOption list field template with object condition option
type ListFieldTmplWithObjOption struct {
	TemplateFilter *filter.Expression `json:"template_filter"`
	ObjectFilter   *filter.Expression `json:"object_filter"`
	Page           BasePage           `json:"page"`
	Fields         []string           `json:"fields"`
}

// ExportBusinessRequest 导出业务请求结构体
type ExportBusinessRequest struct {
	Page          BasePage                  `json:"page"`
	TimeCondition *TimeCondition            `json:"time_condition"`
	Filter        *querybuilder.QueryFilter `json:"filter,omitempty"`
}

// Validate list field template with object condition option
func (l *ListFieldTmplWithObjOption) Validate() errors.RawErrorInfo {
	if err := l.Page.ValidateWithEnableCount(false, common.BKMaxLimitSize); err.ErrCode != 0 {
		return err
	}

	opt := filter.NewDefaultExprOpt(nil)
	opt.IgnoreRuleFields = true

	if l.ObjectFilter != nil {
		if err := l.ObjectFilter.Validate(opt); err != nil {
			return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{err.Error()}}
		}
	}

	if l.TemplateFilter != nil {
		if err := l.TemplateFilter.Validate(opt); err != nil {
			return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{err.Error()}}
		}
	}

	return errors.RawErrorInfo{}
}

// CountFieldTmplResOption list field templates' related resource count option
type CountFieldTmplResOption struct {
	TemplateIDs []int64 `json:"bk_template_ids"`
}

// Validate list field templates' related resource count option
func (c *CountFieldTmplResOption) Validate() errors.RawErrorInfo {
	if len(c.TemplateIDs) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"bk_template_ids"}}
	}

	if len(c.TemplateIDs) > common.BKMaxLimitSize {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"bk_template_ids", common.BKMaxLimitSize},
		}
	}

	return errors.RawErrorInfo{}
}

// FieldTmplResCount field template related resource count info
type FieldTmplResCount struct {
	TemplateID int64 `json:"bk_template_id"`
	Count      int   `json:"count"`
}

// CountFieldTemplateAttrResult count field template attr result
type CountFieldTemplateAttrResult struct {
	BaseResp `json:",inline"`
	Data     []FieldTmplResCount `json:"data"`
}

// CountByIDsOption count by ids option
type CountByIDsOption struct {
	IDs []int64 `json:"ids"`
}

// Validate CountByIDsOption
func (c *CountByIDsOption) Validate() errors.RawErrorInfo {
	if len(c.IDs) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	if len(c.IDs) > common.BKMaxUpdateOrCreatePageSize {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", common.BKMaxUpdateOrCreatePageSize},
		}
	}

	return errors.RawErrorInfo{}
}

// IDCountInfo id to count info
type IDCountInfo struct {
	ID    int64 `json:"id"`
	Count int   `json:"count"`
}
