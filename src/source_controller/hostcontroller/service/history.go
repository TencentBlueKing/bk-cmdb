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
	"encoding/json"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
	"github.com/rs/xid"
	"strconv"
	"time"
)

const HistoryCollection = "cc_History"

func (s *Service) AddHistory(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := req.PathParameter("user")
	ownerID := util.GetOwnerID(pheader)

	bodyData := new(meta.HistoryContent)
	if err := json.NewDecoder(req.Request.Body).Decode(bodyData); err != nil {
		blog.Errorf("add history, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}

	if bodyData.Content == "" {
		blog.Errorf("add history, but history content is empty.")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	id := xid.New().String()
	history := meta.HistoryMeta{
		ID:         id,
		User:       user,
		Content:    bodyData.Content,
		OwnerID:    ownerID,
		CreateTime: time.Now().UTC(),
	}

	_, err := s.Instance.Insert("cc_History", history)
	if nil != err {
		blog.Error("add history failed, err: %v, params: %+v", err, history)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}

	resp.WriteEntity(meta.IDResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     meta.ID{ID: id},
	})

}

func (s *Service) GetHistorys(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := req.PathParameter("user")
	ownerID := util.GetOwnerID(pheader)

	start, err := strconv.Atoi(req.PathParameter("start"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}
	limit, err := strconv.Atoi(req.PathParameter("limit"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	conds := common.KvMap{"user": user}
	conds = util.SetModOwner(conds, ownerID)
	fields := []string{"id", "content", common.CreateTimeField, "user"}
	var result []meta.HistoryMeta
	sort := "-" + common.LastTimeField
	err = s.Instance.GetMutilByCondition(HistoryCollection, fields, conds, &result, sort, start, limit)
	if nil != err {
		blog.Error("query  history failed, err: %v, params: %v", err, conds)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	nums, err := s.Instance.GetCntByCondition(HistoryCollection, conds)
	if nil != err {
		blog.Error("query  history failed, err: %v, params:%v", err, conds)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}

	resp.WriteEntity(meta.GetHistoryResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     meta.HistoryResult{Count: nums, Info: result},
	})
}
