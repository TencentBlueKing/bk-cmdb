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
 
package property

import (
	"configcenter/src/api_server/ccapi/actions/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/defs"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"fmt"

	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
)

var obj *objAction = &objAction{}

type objAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Property/getList", Params: nil, Handler: obj.GetObjProperty, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	obj.CreateAction()
}

// GetObjProperty: get object property, 1, 2, 3, 4 represent app，set,module，host
func (cli *objAction) GetObjProperty(req *restful.Request, resp *restful.Response) {
	blog.Debug("getObjProperty start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getObjProperty error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	res, msg := utils.ValidateFormData(formData, []string{"type"})
	if !res {
		blog.Error("getObjProperty error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	objType := formData["type"][0]

	obj, ok := defs.ObjMap[objType]
	if !ok {
		blog.Error("getObjProperty error, non match objType: %s", objType)
		converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "type").Error(), resp)
		return
	}

	objID := obj["ObjectID"]
	idName := obj["IDName"]
	idDisplayName := obj["IDDisplayName"]

	reqParam := make(map[string]interface{})
	reqParam[common.BKObjIDField] = objID
	//reqParam[common.BKIsPre] = true
	reqParamJson, _ := json.Marshal(reqParam)

	url := fmt.Sprintf("%s/topo/v1/objectattr/search", cli.CC.TopoAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(reqParamJson))
	if err != nil {
		blog.Error("getObjProperty url:%s, data:%s, error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	blog.Debug("getObjProperty rspV3:%v", rspV3)
	resDataV2, err := converter.ResToV2ForPropertyList(rspV3, idName, idDisplayName)
	if err != nil {
		blog.Error("convert property res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}
