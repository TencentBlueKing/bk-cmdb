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

package service

import (
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

// SelectObjectAttWithParams select object's attribute with some params
func (cli *Service) SelectObjectAttWithParams(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetLanguage(req.Request.Header)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// decode json object
	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Errorf("read request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	page := meta.BasePage{Limit: common.BKNoLimit}
	if pageJS, ok := js.CheckGet("page"); ok {
		tmpMap, _ := pageJS.Map()
		page = meta.ParsePage(tmpMap)
		js.Del("page")
	}

	results := make([]meta.Attribute, 0)
	// select from storage
	selector, err := js.Map()
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	selector = util.SetQueryOwner(selector, ownerID)

	if selErr := db.Table(common.BKTableNameObjAttDes).Find(selector).Start(uint64(page.Start)).Limit(uint64(page.Limit)).Sort(page.Sort).All(ctx, &results); nil != selErr && !db.IsNotFoundError(selErr) {
		blog.Errorf("find object by selector failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selErr.Error())})
		return
	}

	// translate language
	for index := range results {
		results[index].PropertyName = cli.TranslatePropertyName(defLang, &results[index])
		results[index].Placeholder = cli.TranslatePlaceholder(defLang, &results[index])
		if results[index].PropertyType == common.FieldTypeEnum {
			results[index].Option = cli.TranslateEnumName(defLang, &results[index], results[index].Option)
		}
	}
	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})
}
