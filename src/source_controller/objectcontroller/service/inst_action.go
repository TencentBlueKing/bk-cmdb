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
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/metadata"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

// UpdateInstObject UpdateInstObject
func (cli *Service) UpdateInstObject(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetLanguage(req.Request.Header)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParams := req.PathParameters()
	objType := pathParams["obj_type"]

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}
	js, err := simplejson.NewJson([]byte(value))
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	input, err := js.Map()
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	data, ok := input["data"].(map[string]interface{})
	if !ok {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsIsInvalid, "lost data field")})
		return
	}

	data[common.LastTimeField] = time.Now()
	condition := input["condition"]
	condition = util.SetModOwner(condition, ownerID)

	// retrieve original datas
	originDatas := make([]map[string]interface{}, 0)
	getErr := cli.GetObjectByCondition(ctx, db, defLang, objType, nil, condition, &originDatas, "", 0, 0)
	if getErr != nil {
		blog.Errorf("retrieve original datas error:%v", getErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, getErr.Error())})
		return
	}

	blog.Infof("update object type:%s,data:%v,condition:%v", objType, data, condition)
	err = cli.UpdateObjByCondition(ctx, db, objType, data, condition)
	if err != nil {
		blog.Errorf("update object type:%s,data:%v,condition:%v,error:%v", objType, data, condition, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	// record event
	if len(originDatas) > 0 {
		idname := common.GetInstIDField(objType)
		for _, originData := range originDatas {
			newData := map[string]interface{}{}
			id, err := strconv.Atoi(fmt.Sprintf("%v", originData[idname]))
			if err != nil {
				blog.Errorf("create event error:%v", err)
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrEventPushEventFailed)})
				return
			}
			realObjType := objType
			if objType == common.BKInnerObjIDObject {
				var ok bool
				realObjType, ok = originData[common.BKObjIDField].(string)
				if !ok {
					blog.Errorf("create event error: there is no bk_obj_type exist, err: %v, originData: %#v", err, originData)
					resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrEventPushEventFailed)})
					return
				}
			}
			if err := cli.GetObjectByID(ctx, db, realObjType, nil, id, &newData, ""); err != nil {
				blog.Errorf("create event error:%v", err)
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrEventPushEventFailed)})
			} else {
				srcevent := eventclient.NewEventWithHeader(req.Request.Header)
				srcevent.EventType = metadata.EventTypeInstData
				srcevent.ObjType = objType
				srcevent.Action = metadata.EventActionUpdate
				srcevent.Data = []metadata.EventData{
					{
						CurData: newData,
						PreData: originData,
					},
				}
				err = cli.EventC.Push(ctx, srcevent)
				if err != nil {
					blog.Errorf("create event error:%v", err)
					resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrEventPushEventFailed)})
					return
				}
			}
		}

	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})

}
