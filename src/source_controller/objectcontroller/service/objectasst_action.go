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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
)

func (cli *Service) createObjectAssociation(req *restful.Request, resp *restful.Response, enableMainlineAssociationType bool) {
	// enableMainlineAssociationType used for distinguish two creation mode
	// when enableMainlineAssociationType enabled, only bk_mainline type could be create
	// when enableMainlineAssociationType disabled, all type except bk_mainline could be create

	// get the language
	language := util.GetLanguage(req.Request.Header)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	obj := &meta.Association{}
	if err = json.Unmarshal([]byte(value), obj); nil != err {
		blog.Errorf("fail to unmarshal json, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	if enableMainlineAssociationType == false {
		// AsstKindID shouldn't be use bk_mainline
		if obj.AsstKindID == common.AssociationKindMainline {
			blog.Errorf("use inner association type: %v is forbidden", common.AssociationKindMainline)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrorTopoAssociationKindMainlineUnavailable, obj.AsstKindID)})
			return
		}
	} else {
		// AsstKindID could only be bk_mainline
		if obj.AsstKindID != common.AssociationKindMainline {
			blog.Errorf("use CreateMainlineObjectAssociation method but bk_asst_id is: %s", obj.AsstKindID)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrorTopoAssociationKindInconsistent, obj.AsstKindID)})
			return
		}
	}

	// save to the storage
	id, err := db.NextSequence(ctx, common.BKTableNameObjAsst)
	if err != nil {
		blog.Errorf("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	obj.ID = int64(id)
	obj.OwnerID = ownerID
	err = db.Table(common.BKTableNameObjAsst).Insert(ctx, obj)
	if nil != err {
		blog.Errorf("create objectasst failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: []*meta.Association{obj}})
}
