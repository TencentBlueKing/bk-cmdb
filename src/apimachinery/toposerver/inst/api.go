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

package inst

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
)

// InstanceInterface instance operation interface
type InstanceInterface interface {
	CreateApp(ctx context.Context, ownerID string, h http.Header, dat map[string]interface{}) (resp *metadata.CreateInstResult, err error)
	DeleteApp(ctx context.Context, ownerID string, appID string, h http.Header) (resp *metadata.Response, err error)
	UpdateApp(ctx context.Context, ownerID string, appID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	UpdateAppDataStatus(ctx context.Context, ownerID string, flag common.DataStatusFlag, appID string, h http.Header) (resp *metadata.Response, err error)
	SearchApp(ctx context.Context, ownerID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error)
	GetAppBasicInfo(ctx context.Context, h http.Header, bizID int64) (resp *metadata.AppBasicInfoResult, err error)
	GetDefaultApp(ctx context.Context, ownerID string, h http.Header) (resp *metadata.SearchInstResult, err error)
	CreateDefaultApp(ctx context.Context, ownerID string, h http.Header, data map[string]interface{}) (resp *metadata.CreateInstResult, err error)
	SearchAuditDict(ctx context.Context, h http.Header) (resp *metadata.Response, err error)
	SearchAuditList(ctx context.Context, h http.Header, input *metadata.AuditQueryInput) (*metadata.Response, error)
	SearchAuditDetail(ctx context.Context, h http.Header, input *metadata.AuditDetailQueryInput) (*metadata.Response, error)
	GetInternalModule(ctx context.Context, ownerID, appID string, h http.Header) (resp *metadata.SearchInnterAppTopoResult, err error)
	SearchBriefBizTopo(ctx context.Context, h http.Header, bizID int64, input map[string]interface{}) (resp *metadata.SearchBriefBizTopoResult, err error)
	CreateInst(ctx context.Context, objID string, h http.Header, dat interface{}) (resp *metadata.CreateInstResult, err error)
	CreateManyCommInst(ctx context.Context, objID string, header http.Header, data metadata.CreateManyCommInst) (resp *metadata.CreateManyCommInstResult, err error)
	DeleteInst(ctx context.Context, objID string, instID int64, h http.Header) (resp *metadata.Response, err error)
	UpdateInst(ctx context.Context, objID string, instID int64, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	SelectInsts(ctx context.Context, ownerID string, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	SelectInstsAndAsstDetail(ctx context.Context, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	InstSearch(ctx context.Context, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	SelectInstsByAssociation(ctx context.Context, objID string, h http.Header, p *metadata.AssociationParams) (resp *metadata.SearchInstResult, err error)
	SelectInst(ctx context.Context, objID string, instID int64, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	SelectTopo(ctx context.Context, objID string, instID int64, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchTopoResult, err error)
	SelectAssociationTopo(ctx context.Context, objID string, instID int64, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchAssociationTopoResult, err error)
	CreateModule(ctx context.Context, appID, setID int64, h http.Header, dat map[string]interface{}) (mapstr.MapStr,
		errors.CCErrorCoder)
	DeleteModule(ctx context.Context, appID, setID, moduleID int64, h http.Header) errors.CCErrorCoder
	UpdateModule(ctx context.Context, appID, setID, moduleID int64, h http.Header,
		dat map[string]interface{}) errors.CCErrorCoder
	SearchModule(ctx context.Context, ownerID string, appID, setID int64, h http.Header, s *params.SearchParams) (
		*metadata.InstResult, errors.CCErrorCoder)
	SearchModuleByCondition(ctx context.Context, appID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error)
	SearchModuleBatch(ctx context.Context, appID string, h http.Header, s *metadata.SearchInstBatchOption) (resp *metadata.MapArrayResponse, err error)
	SearchModuleWithRelation(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.ResponseInstData, err error)
	CreateSet(ctx context.Context, appID int64, h http.Header, dat mapstr.MapStr) (mapstr.MapStr, errors.CCErrorCoder)
	DeleteSet(ctx context.Context, appID, setID int64, h http.Header) errors.CCErrorCoder
	UpdateSet(ctx context.Context, appID, setID int64, h http.Header, dat map[string]interface{}) errors.CCErrorCoder
	SearchSet(ctx context.Context, ownerID string, appID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error)
	SearchSetBatch(ctx context.Context, appID string, h http.Header, s *metadata.SearchInstBatchOption) (resp *metadata.MapArrayResponse, err error)
	SearchInstsNames(ctx context.Context, h http.Header, s *metadata.SearchInstsNamesOption) (resp *metadata.ArrayResponse, err error)
	GetTopoNodeHostAndServiceInstCount(ctx context.Context, h http.Header, objID int64,
		s *metadata.HostAndSerInstCountOption) (resp *metadata.GetHostAndSerInstCountResult, err error)

	// SearchObjectInstances searches object instances.
	SearchObjectInstances(ctx context.Context, header http.Header,
		objID string, input *metadata.CommonSearchFilter) (*metadata.Response, error)

	// CountObjectInstances counts object instances num.
	CountObjectInstances(ctx context.Context, header http.Header,
		objID string, input *metadata.CommonCountFilter) (*metadata.Response, error)

	// CreateBizSet create biz set
	CreateBizSet(ctx context.Context, h http.Header, opt metadata.CreateBizSetRequest) (int64, errors.CCErrorCoder)

	// UpdateBizSet update biz set
	UpdateBizSet(ctx context.Context, h http.Header, opt metadata.UpdateBizSetOption) errors.CCErrorCoder

	// DeleteBizSet delete biz set
	DeleteBizSet(ctx context.Context, h http.Header, opt metadata.DeleteBizSetOption) errors.CCErrorCoder

	// FindBizInBizSet find biz list in biz set
	FindBizInBizSet(ctx context.Context, h http.Header, opt *metadata.FindBizInBizSetOption) (*metadata.InstResult,
		errors.CCErrorCoder)

	// FindBizSetTopo find topo info by parent in biz set
	FindBizSetTopo(ctx context.Context, h http.Header, opt *metadata.FindBizSetTopoOption) ([]mapstr.MapStr,
		errors.CCErrorCoder)

	// SearchBusinessSet search business set
	SearchBusinessSet(ctx context.Context, h http.Header, opt *metadata.QueryBusinessSetRequest) (
		*metadata.InstResult, errors.CCErrorCoder)
}

type instanceClient struct {
	client rest.ClientInterface
}

// NewInstanceClient TODO
func NewInstanceClient(client rest.ClientInterface) InstanceInterface {
	return &instanceClient{
		client: client,
	}
}
