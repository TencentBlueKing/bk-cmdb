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
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

var flag *flagAction = &flagAction{}

type flagAction struct {
	base.BaseAction
}

func init() {
	flag.CreateAction()

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/migrate/system/hostcrossbiz/{ownerID}", Params: nil, Handler: flag.Set})
	// create CC object
}

func (cli *flagAction) Set(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	var ownerID string
	pathParameters := req.PathParameters()
	if err := cli.GetParams(cli.CC, &pathParameters, "ownerID", &ownerID, resp); nil != err {
		blog.Error("failed to get params, error info is %s", err.Error())
		cli.ResponseFailed(common.CCErrCommHTTPInputInvalid, defErr.Error(common.CCErrCommHTTPInputInvalid), resp)
		return
	}
	if ownerID == "" {
		ownerID = common.BKDefaultOwnerID
	}

	a := api.GetAPIResource()

	blog.Errorf("modify data for  %s table ", "cc_System")
	cond := map[string]interface{}{
		common.HostCrossBizField: common.HostCrossBizValue}
	data := map[string]interface{}{
		common.HostCrossBizField: common.HostCrossBizValue + ownerID}

	err := a.InstCli.UpdateByCondition("cc_System", data, cond)
	if nil != err {
		blog.Errorf("modify data for  %s table error  %s", "cc_System", err)
		cli.ResponseFailed(common.CCErrCommMigrateFailed, defErr.Error(common.CCErrCommMigrateFailed), resp)
		return
	}

	cli.ResponseSuccess("modify system config success", resp)
	return

}
