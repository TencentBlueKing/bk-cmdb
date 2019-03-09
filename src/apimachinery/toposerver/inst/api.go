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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

type InstanceInterface interface {
	CreateApp(ctx context.Context, ownerID string, h http.Header, dat map[string]interface{}) (resp *metadata.CreateInstResult, err error)
	DeleteApp(ctx context.Context, ownerID string, appID string, h http.Header) (resp *metadata.Response, err error)
	UpdateApp(ctx context.Context, ownerID string, appID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	UpdateAppDataStatus(ctx context.Context, ownerID string, flag common.DataStatusFlag, appID string, h http.Header) (resp *metadata.Response, err error)
	SearchApp(ctx context.Context, ownerID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error)
	GetDefaultApp(ctx context.Context, ownerID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error)
	CreateDefaultApp(ctx context.Context, ownerID string, h http.Header, data map[string]interface{}) (resp *metadata.CreateInstResult, err error)
	QueryAudit(ctx context.Context, ownerID string, h http.Header, input *metadata.QueryInput) (resp *metadata.Response, err error)
	GetInternalModule(ctx context.Context, ownerID, appID string, h http.Header) (resp *metadata.SearchInnterAppTopoResult, err error)
	CreateInst(ctx context.Context, ownerID string, objID string, h http.Header, dat interface{}) (resp *metadata.CreateInstResult, err error)
	DeleteInst(ctx context.Context, ownerID string, objID string, instID int64, h http.Header) (resp *metadata.Response, err error)
	UpdateInst(ctx context.Context, ownerID string, objID string, instID int64, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	SelectInsts(ctx context.Context, ownerID string, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	SelectInstsAndAsstDetail(ctx context.Context, ownerID string, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	InstSearch(ctx context.Context, ownerID string, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	SelectInstsByAssociation(ctx context.Context, ownerID string, objID string, h http.Header, p *metadata.AssociationParams) (resp *metadata.SearchInstResult, err error)
	SelectInst(ctx context.Context, ownerID string, objID string, instID string, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	SelectTopo(ctx context.Context, ownerID string, objID string, instID string, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	SelectAssociationTopo(ctx context.Context, ownerID string, objID string, instID string, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchInstResult, err error)
	CreateModule(ctx context.Context, appID string, setID string, h http.Header, dat map[string]interface{}) (resp *metadata.CreateInstResult, err error)
	DeleteModule(ctx context.Context, appID string, setID string, moduleID string, h http.Header) (resp *metadata.Response, err error)
	UpdateModule(ctx context.Context, appID string, setID string, moduleID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	SearchModule(ctx context.Context, ownerID string, appID string, setID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error)
	CreateSet(ctx context.Context, appID string, h http.Header, dat mapstr.MapStr) (resp *metadata.CreateInstResult, err error)
	DeleteSet(ctx context.Context, appID string, setID string, h http.Header) (resp *metadata.Response, err error)
	UpdateSet(ctx context.Context, appID string, setID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	SearchSet(ctx context.Context, ownerID string, appID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error)
}

type instanceClient struct {
	client rest.ClientInterface
}

func NewInstanceClient(client rest.ClientInterface) InstanceInterface {
	return &instanceClient{
		client: client,
	}
}
