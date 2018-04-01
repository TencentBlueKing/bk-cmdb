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
 
package host

import (
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"fmt"
	"io"

	"github.com/emicklei/go-restful"
)

func AddHostFavourite(req *restful.Request, resp *restful.Response) {
	cc := api.NewAPIResource()
	url := cc.HostAPI() + "/host/v1/hosts/favorites"

	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

func EditHostFavourite(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")

	cc := api.NewAPIResource()
	url := cc.HostAPI() + "/host/v1/hosts/favorites/"
	url = fmt.Sprintf("%s%s", url, id)

	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
	return
}

func DeleteHostFavourite(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")

	cc := api.NewAPIResource()
	url := cc.HostAPI() + "/host/v1/hosts/favorites/"
	url = fmt.Sprintf("%s%s", url, id)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)
	io.WriteString(resp, rsp)
	return
}

func GetHostFavourites(req *restful.Request, resp *restful.Response) {

	cc := api.NewAPIResource()
	url := cc.HostAPI() + "/host/v1/hosts/favorites/search"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
	return
}

func IncrHostFavouritesCount(req *restful.Request, resp *restful.Response) {
	cc := api.NewAPIResource()

	id := req.PathParameter("id")
	url := cc.HostAPI() + "/host/v1/hosts/favorites/" + id + "/incr"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
	return
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/favorites/search", Params: nil, Handler: GetHostFavourites, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/favorites", Params: nil, Handler: AddHostFavourite, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/hosts/favorites/{id}", Params: nil, Handler: EditHostFavourite, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/hosts/favorites/{id}", Params: nil, Handler: DeleteHostFavourite, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/hosts/favorites/{id}/incr", Params: nil, Handler: IncrHostFavouritesCount, FilterHandler: nil, Version: v3.APIVersion})

}
