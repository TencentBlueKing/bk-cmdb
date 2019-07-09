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
	"bytes"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"context"
	"github.com/gin-gonic/gin/json"
	"io"
	"net/http"

	"github.com/emicklei/go-restful"
)

func DecodeJSON(r io.Reader, v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := json.NewDecoder(io.TeeReader(r, buf)).Decode(v)
	return buf.Bytes(), err
}

// CreateObjectUnique create object's unique

// UpdateObjectUnique update object's unique
// DeleteObjectUnique delte object's unique
// SearchObjectUnique delte object's unique
func (cli *Service) SearchObjectUnique(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	objID := req.PathParameter(common.BKObjIDField)

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ownerID)

	uniques, err := cli.searchObjectUnique(ctx, db, ownerID, objID)
	if nil != err {
		blog.Errorf("[SearchObjectUnique] Search error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	resp.WriteEntity(metadata.SearchUniqueResult{BaseResp: metadata.SuccessBaseResp, Data: uniques})
}

func (cli *Service) searchObjectUnique(ctx context.Context, db dal.RDB, ownerID, objID string) ([]metadata.ObjectUnique, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ownerID)

	uniques := []metadata.ObjectUnique{}
	err := db.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).All(ctx, &uniques)
	return uniques, err
}
