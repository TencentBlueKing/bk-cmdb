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
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type ApiServerClientInterface interface {
	AddDefaultApp(ctx context.Context, h http.Header, ownerID string, params mapstr.MapStr) (resp *metadata.Response, err error)
	SearchDefaultApp(ctx context.Context, h http.Header, ownerID string, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error)
	GetRolePrivilege(ctx context.Context, h http.Header, ownerID, objID, role string) (resp *metadata.RolePriResult, err error)
	GetAppRole(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.RoleAppResult, err error)
	GetUserPrivilegeApp(ctx context.Context, h http.Header, ownerID, userName string, params mapstr.MapStr) (resp *metadata.AppQueryResult, err error)
	GetUserPrivilegeConfig(ctx context.Context, h http.Header, ownerID, userName string) (resp *metadata.UserPriviResult, err error)
	GetAllMainLineObject(ctx context.Context, h http.Header, ownerID, userName string) (resp *metadata.MainLineResult, err error)
	GetObjectData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrBatchResult, err error)
	GetInstDetail(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error)
	GetObjectAttr(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrResult, err error)
	GetHostData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error)
	GetObjectGroup(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.ObjectAttrGroupResult, err error)
	AddHost(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error)
	AddInst(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error)
	AddObjectBatch(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.Response, err error)
	SearchAssociationInst(ctx context.Context, h http.Header, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error)
	SearchInsts(ctx context.Context, h http.Header, objID string, cond condition.Condition) (resp *metadata.ResponseInstData, err error)
	ImportAssociation(ctx context.Context, h http.Header, objID string, input *metadata.RequestImportAssociation) (resp *metadata.ResponeImportAssociation, err error)
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
