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
	"configcenter/src/common/util"

	"configcenter/src/scene_server/admin_server/migrate_service/logics"

	"github.com/emicklei/go-restful"
)

var clear *clearAction = &clearAction{}

type clearAction struct {
	base.BaseAction
}

func init() {
	clear.CreateAction()

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/clear", Params: nil, Handler: clear.clear})
	// create CC object
}

func (cli *clearAction) clear(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	err := logics.Clear(cli.CC.InstCli)
	if nil != err {
		blog.Errorf("clear error: %v", err)
		cli.ResponseFailed(common.CCErrCommMigrateFailed, defErr.Error(common.CCErrCommMigrateFailed), resp)
		return
	}

	cli.ResponseSuccess("migrate success", resp)
	return

}
