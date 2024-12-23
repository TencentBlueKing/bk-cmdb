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
	"strconv"
	"strings"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

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

	// update instance api
	ws.PUT("/table/hosts/batch", s.UpdateHostWithTable)
	ws.PUT("/table/updatemany/biz_set", s.UpdateBizSetWithTable)
	ws.PUT("/table/biz/:bk_supplier_account/:id", s.UpdateOneInstWithTable(common.BKInnerObjIDApp))
	ws.PUT("/table/set/:bk_biz_id/:id", s.UpdateOneInstWithTable(common.BKInnerObjIDSet))
	ws.PUT("/table/module/:bk_biz_id/:bk_set_id/:id", s.UpdateOneInstWithTable(common.BKInnerObjIDModule))
	ws.PUT("/table/updatemany/project", s.UpdateProjectWithTable)
	ws.PUT("/table/update/instance/object/:bk_obj_id/inst/:id", s.UpdateCommonInstWithTable)

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

		txnErr := s.ApiCli.Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
			// get table attributes
			attrOpt := mapstr.MapStr{common.BKObjIDField: objID}
			attributes, err := s.ApiCli.ModelQuote().GetObjectAttrWithTable(kit.Ctx, kit.Header, attrOpt)
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
				if len(attrVal) > 0 {
					tableOpt := &metadata.BatchCreateQuotedInstOption{
						ObjID:      objID,
						PropertyID: attr.PropertyID,
						Data:       attrVal,
					}
					ids, err := s.ApiCli.ModelQuote().BatchCreateQuotedInstance(kit.Ctx, kit.Header, tableOpt)
					if err != nil {
						return err
					}

					// set table instance ids to attach instance
					data[attr.PropertyID] = ids
				}
			}

			// proxy request to api server
			resp := new(metadata.Response)
			err = s.ApiCli.Client().Post().
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

// updateInstanceWithTable update instance with table attributes, ** ONLY FOR UI **
func (s *Service) updateInstanceWithTable(objID string, id int64, data mapstr.MapStr) gin.HandlerFunc {
	return func(c *gin.Context) {
		kit := rest.NewKitFromHeader(c.Request.Header, s.CCErr)

		txnErr := s.ApiCli.Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
			// get table attributes
			attrOpt := mapstr.MapStr{common.BKObjIDField: objID}
			attributes, err := s.ApiCli.ModelQuote().GetObjectAttrWithTable(kit.Ctx, kit.Header, attrOpt)
			if err != nil {
				blog.Errorf("get object(%s) attributes failed, err: %v, rid: %s", objID, err, kit.Rid)
				return err
			}

			for _, attr := range attributes {
				if attr.PropertyType != common.FieldTypeInnerTable {
					continue
				}

				if !attr.IsEditable {
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

				// update table attributes
				if err = s.updateTableAttr(kit, objID, attr.PropertyID, id, attrVal); err != nil {
					return err
				}

				delete(data, attr.PropertyID)
			}

			// proxy request to api server
			resp := new(metadata.Response)
			err = s.ApiCli.Client().Put().
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

		c.JSON(http.StatusOK, metadata.NewSuccessResp(nil))
	}
}

func (s *Service) updateTableAttr(kit *rest.Kit, objID, attrID string, instID int64, attrVal []mapstr.MapStr) error {
	// list table instances
	idMap, err := s.getQuotedInstIDMap(kit, objID, attrID, instID)
	if err != nil {
		return err
	}

	// cross compare table instances
	createData := make([]mapstr.MapStr, 0)
	for _, val := range attrVal {
		idVal, exists := val[common.BKFieldID]
		if !exists {
			val[common.BKInstIDField] = instID
			createData = append(createData, val)
			continue
		}

		id, err := util.GetInt64ByInterface(idVal)
		if err != nil {
			blog.Errorf("parse input id failed, err: %v, id: %+v, rid: %s", err, idVal, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID)
		}

		// update table instances
		updateOpt := &metadata.BatchUpdateQuotedInstOption{
			ObjID:      objID,
			PropertyID: attrID,
			IDs:        []uint64{uint64(id)},
			Data:       val,
		}
		err = s.ApiCli.ModelQuote().BatchUpdateQuotedInstance(kit.Ctx, kit.Header, updateOpt)
		if err != nil {
			blog.Errorf("update quoted instance failed, err: %v, opt: %+v, rid: %s", err, updateOpt, kit.Rid)
			return err
		}

		delete(idMap, uint64(id))
	}

	deleteIDs := make([]uint64, 0)
	for id := range idMap {
		deleteIDs = append(deleteIDs, id)
	}

	// delete redundant table instances
	if len(deleteIDs) > 0 {
		deleteOpt := &metadata.BatchDeleteQuotedInstOption{
			ObjID:      objID,
			PropertyID: attrID,
			IDs:        deleteIDs,
		}
		err = s.ApiCli.ModelQuote().BatchDeleteQuotedInstance(kit.Ctx, kit.Header, deleteOpt)
		if err != nil {
			blog.Errorf("delete quoted instance failed, err: %v, opt: %+v, rid: %s", err, deleteOpt, kit.Rid)
			return err
		}
	}

	// create table instances
	if len(createData) > 0 {
		createOpt := &metadata.BatchCreateQuotedInstOption{
			ObjID:      objID,
			PropertyID: attrID,
			Data:       createData,
		}
		_, err = s.ApiCli.ModelQuote().BatchCreateQuotedInstance(kit.Ctx, kit.Header, createOpt)
		if err != nil {
			blog.Errorf("create quoted instance failed, err: %v, opt: %+v, rid: %s", err, createOpt, kit.Rid)
			return err
		}
	}

	return nil
}

func (s *Service) getQuotedInstIDMap(kit *rest.Kit, objID string, attrID string, instID int64) (
	map[uint64]struct{}, error) {

	listOpt := &metadata.ListQuotedInstOption{
		ObjID:      objID,
		PropertyID: attrID,
		CommonQueryOption: metadata.CommonQueryOption{
			CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(
				common.BKInstIDField, filter.Equal, instID)},
			Page:   metadata.BasePage{Limit: common.BKMaxPageSize},
			Fields: []string{common.BKFieldID},
		},
	}
	listRes, err := s.ApiCli.ModelQuote().ListQuotedInstance(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list quoted instance failed, err: %v, opt: %+v, rid: %s", err, listOpt, kit.Rid)
		return nil, err
	}

	idMap := make(map[uint64]struct{})
	for _, info := range listRes.Info {
		id, err := util.GetInt64ByInterface(info[common.BKFieldID])
		if err != nil {
			blog.Errorf("parse db id failed, err: %v, id: %+v, rid: %s", err, info[common.BKFieldID], kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID)
		}

		idMap[uint64(id)] = struct{}{}
	}
	return idMap, nil
}

// UpdateOneInstWithTable update one instance with table attributes, ** ONLY FOR UI **
func (s *Service) UpdateOneInstWithTable(objID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param(common.BKFieldID)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommParamsNeedSet, ErrMsg: err.Error()})
			return
		}

		data := make(mapstr.MapStr)
		if err = json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
			c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
			return
		}

		handler := s.updateInstanceWithTable(objID, id, data)
		handler(c)
	}
}

// UpdateCommonInstWithTable update common instance with table attributes, ** ONLY FOR UI **
func (s *Service) UpdateCommonInstWithTable(c *gin.Context) {
	objID := c.Param(common.BKObjIDField)
	handler := s.UpdateOneInstWithTable(objID)
	handler(c)
}

// UpdateHostWithTable update one host with table attributes using batch api, ** ONLY FOR UI **
func (s *Service) UpdateHostWithTable(c *gin.Context) {
	data := make(mapstr.MapStr)
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
		return
	}

	idStr, ok := data[common.BKHostIDField].(string)
	if !ok {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommParamsInvalid, ErrMsg: "host id is invalid"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommParamsNeedSet, ErrMsg: err.Error()})
		return
	}

	handler := s.updateInstanceWithTable(common.BKInnerObjIDHost, id, data)
	handler(c)
}

// UpdateBizSetWithTable update one biz set with table attributes using batch api, ** ONLY FOR UI **
func (s *Service) UpdateBizSetWithTable(c *gin.Context) {
	data := new(metadata.UpdateBizSetOption)
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
		return
	}

	if len(data.BizSetIDs) != 1 {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommParamsInvalid, ErrMsg: "id length should == 1"})
		return
	}

	updateData, err := mapstr.Struct2Map(data.Data)
	if err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommParamsInvalid, ErrMsg: err.Error()})
		return
	}

	handler := s.updateInstanceWithTable(common.BKInnerObjIDHost, data.BizSetIDs[0], updateData)
	handler(c)
}

// UpdateProjectWithTable update one project with table attributes using batch api, ** ONLY FOR UI **
func (s *Service) UpdateProjectWithTable(c *gin.Context) {
	data := new(metadata.UpdateProjectOption)
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
		return
	}

	if len(data.IDs) != 1 {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommParamsInvalid, ErrMsg: "id length should == 1"})
		return
	}

	handler := s.updateInstanceWithTable(common.BKInnerObjIDHost, data.IDs[0], data.Data)
	handler(c)
}
