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

package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
)

type ApiServerClientInterface interface {
	Client() rest.ClientInterface

	AddDefaultApp(ctx context.Context, h http.Header, ownerID string, params mapstr.MapStr) (resp *metadata.Response, err error)
	SearchDefaultApp(ctx context.Context, h http.Header, ownerID string) (resp *metadata.QueryInstResult, err error)
	GetObjectData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrBatchResult, err error)
	GetInstDetail(ctx context.Context, h http.Header, objID string, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error)
	CreateObjectAtt(ctx context.Context, h http.Header, obj *metadata.ObjAttDes) (resp *metadata.Response, err error)
	UpdateObjectAtt(ctx context.Context, objID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteObjectAtt(ctx context.Context, objID string, h http.Header) (resp *metadata.Response, err error)
	GetObjectAttr(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrResult, err error)
	GetHostData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error)
	ListHostWithoutApp(ctx context.Context, h http.Header, option metadata.ListHostsWithNoBizParameter) (resp *metadata.ListHostWithoutAppResponse, err error)
	GetObjectGroup(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.ObjectAttrGroupResult, err error)
	AddHost(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error)
	AddHostByExcel(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error)
	UpdateHost(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error)
	GetHostModuleRelation(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.HostModuleResp, err error)
	AddInst(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error)
	AddObjectBatch(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.Response, err error)
	SearchAssociationInst(ctx context.Context, h http.Header, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error)
	ImportAssociation(ctx context.Context, h http.Header, objID string, input *metadata.RequestImportAssociation) (resp *metadata.ResponeImportAssociation, err error)

	GetUserAuthorizedBusinessList(ctx context.Context, h http.Header, user string) (resp *metadata.InstDataInfo, err error)

	SearchNetCollectDevice(ctx context.Context, h http.Header, cond condition.Condition) (resp *metadata.ResponseInstData, err error)
	SearchNetDeviceProperty(ctx context.Context, h http.Header, cond condition.Condition) (resp *metadata.ResponseInstData, err error)
	SearchNetCollectDeviceBatch(ctx context.Context, h http.Header, cond mapstr.MapStr) (resp *metadata.ResponseInstData, err error)
	SearchNetDevicePropertyBatch(ctx context.Context, h http.Header, cond mapstr.MapStr) (resp *metadata.ResponseInstData, err error)

	CreateBiz(ctx context.Context, ownerID string, h http.Header, dat map[string]interface{}) (resp *metadata.CreateInstResult, err error)
	UpdateBiz(ctx context.Context, ownerID string, bizID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	UpdateBizDataStatus(ctx context.Context, ownerID string, flag common.DataStatusFlag, bizID string, h http.Header) (resp *metadata.Response, err error)
	SearchBiz(ctx context.Context, ownerID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error)
}

func NewApiServerClientInterface(c *util.Capability, version string) ApiServerClientInterface {
	base := fmt.Sprintf("/api/%s", version)
	return &apiServer{
		client: rest.NewRESTClient(c, base),
	}
}

type apiServer struct {
	client rest.ClientInterface
}
