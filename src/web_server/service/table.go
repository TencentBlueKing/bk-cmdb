/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"encoding/json"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/gin-gonic/gin"
)

func (s *Service) initModelQuote(ws *gin.Engine) {
	// create instance api
	ws.POST("/table/create/biz_set", s.CreateInstanceWithTable(common.BKInnerObjIDBizSet))
	ws.POST("/table/biz/:bk_supplier_account", s.CreateInstanceWithTable(common.BKInnerObjIDApp))
	ws.POST("/table/set/:bk_biz_id", s.CreateInstanceWithTable(common.BKInnerObjIDSet))
	ws.POST("/table/module/:bk_biz_id/:bk_set_id", s.CreateInstanceWithTable(common.BKInnerObjIDModule))
	ws.POST("/table/createmany/project", s.CreateInstanceWithTable(common.BKInnerObjIDProject))
	ws.POST("/table/create/instance/object/:bk_obj_id", s.CreateCommonInstWithTable)

}

// CreateInstanceWithTable create instance with table attributes, ** ONLY FOR UI **
func (s *Service) CreateInstanceWithTable(objID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		kit := rest.NewKitFromHeader(c.Request.Header, s.CCErr)

		data := make(mapstr.MapStr)
		if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
			c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
			return
		}

		var res interface{}

		// start transaction ** NOTICE: web-server should not call core-service, this is only a temporary use **
		txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
			// get table attributes
			attrOpt := mapstr.MapStr{common.BKObjIDField: objID}
			attributes, err := s.CoreAPI.ApiServer().ModelQuote().GetObjectAttrWithTable(kit.Ctx, kit.Header, attrOpt)
			if err != nil {
				blog.Errorf("get object(%s) attributes failed, err: %v, rid: %s", objID, err, kit.Rid)
				return err
			}

			for _, attr := range attributes {
				if attr.PropertyType != common.FieldTypeInnerTable {
					continue
				}

				attrData, exists := data[attr.PropertyID]
				if !exists {
					continue
				}

				attrVal, err := data.MapStrArray(attr.PropertyID)
				if err != nil {
					blog.Errorf("table attribute type is invalid, err: %v, attr: %+v, rid: %s", err, attrData, kit.Rid)
					return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, attr.PropertyID)
				}

				// create table instances
				tableOpt := &metadata.BatchCreateQuotedInstOption{
					ObjID:      objID,
					PropertyID: attr.PropertyID,
					Data:       attrVal,
				}
				ids, err := s.CoreAPI.ApiServer().ModelQuote().BatchCreateQuotedInstance(kit.Ctx, kit.Header, tableOpt)
				if err != nil {
					return err
				}

				// set table instance ids to attach instance
				data[attr.PropertyID] = ids
			}

			// proxy request to api server
			resp := new(metadata.Response)
			err = s.CoreAPI.ApiServer().Client().Post().
				WithContext(kit.Ctx).
				Body(data).
				SubResourcef(strings.TrimPrefix(c.Request.URL.Path, "/table")).
				WithHeaders(kit.Header).
				Do().
				Into(resp)

			if err != nil {
				return errors.CCHttpError
			}

			if err = resp.CCError(); err != nil {
				return err
			}

			res = resp.Data
			return nil
		})

		if txnErr != nil {
			errCode := common.CCErrorUnknownOrUnrecognizedError
			ccErr, ok := txnErr.(errors.CCErrorCoder)
			if ok {
				errCode = ccErr.GetCode()
			}
			c.JSON(http.StatusOK, metadata.BaseResp{Code: errCode, ErrMsg: txnErr.Error()})
			return
		}

		c.JSON(http.StatusOK, metadata.NewSuccessResp(res))
	}
}

// CreateCommonInstWithTable create common instance with table attributes, ** ONLY FOR UI **
func (s *Service) CreateCommonInstWithTable(c *gin.Context) {
	objID := c.Param(common.BKObjIDField)
	handler := s.CreateInstanceWithTable(objID)
	handler(c)
}
