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

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/topo_service/actions/inst"
	"configcenter/src/source_controller/common/commondata"
)

type InstanceInterface interface {
	// app operation
	CreateApp(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error)
	DeleteApp(ctx context.Context, appID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdateApp(ctx context.Context, appID string, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	UpdateAppDataStatus(ctx context.Context, flag common.DataStatusFlag, appID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	SearchApp(ctx context.Context, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error)
	GetDefaultApp(ctx context.Context, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error)
	CreateDefaultApp(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)

	// inst operation
	CreateInst(ctx context.Context, objID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
	DeleteInst(ctx context.Context, objID string, instID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdateInst(ctx context.Context, objID string, instID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SelectInsts(ctx context.Context, objID string, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error)
	SelectInstsAndAsstDetail(ctx context.Context, objID string, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error)
	InstSearch(ctx context.Context, objID string, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error)
	SelectInstsByAssociation(ctx context.Context, objID string, h util.Headers, p *inst.AssociationParams) (resp *api.BKAPIRsp, err error)
	SelectInst(ctx context.Context, objID string, instID string, h util.Headers, p *params.SearchParams) (resp *api.BKAPIRsp, err error)
	SelectTopo(ctx context.Context, objID string, instID string, h util.Headers, p *params.SearchParams) (resp *api.BKAPIRsp, err error)
	SelectAssociationTopo(ctx context.Context, objID string, instID string, h util.Headers, p *params.SearchParams) (resp *api.BKAPIRsp, err error)

	// module operation
	CreateModule(ctx context.Context, appID string, setID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	DeleteModule(ctx context.Context, appID string, setID string, moduleID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdateModule(ctx context.Context, appID string, setID string, moduleID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SearchModule(ctx context.Context, appID string, setID string, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error)

	// set operation
	CreateSet(ctx context.Context, appID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	DeleteSet(ctx context.Context, appID string, setID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdateSet(ctx context.Context, appID string, setID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SearchSet(ctx context.Context, ownerID string, appID string, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error)

	// common operation
	QueryAudit(ctx context.Context, h util.Headers, input *commondata.ObjQueryInput) (resp *api.BKAPIRsp, err error)
	GetInternalModule(ctx context.Context, ownerID, appID string, h util.Headers) (resp *api.BKAPIRsp, err error)
}

type instanceClient struct {
	client rest.ClientInterface
}

func NewInstanceClient(client rest.ClientInterface) InstanceInterface {
	return &instanceClient{
		client: client,
	}
}
