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

package logics

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

type InstNameAsst struct {
	ID         string                 `json:"id"`
	ObjID      string                 `json:"bk_obj_id"`
	ObjIcon    string                 `json:"bk_obj_icon"`
	ObjectID   int64                  `json:"bk_inst_id"`
	ObjectName string                 `json:"bk_obj_name"`
	Name       string                 `json:"bk_inst_name"`
	InstInfo   map[string]interface{} `json:"inst_info,omitempty"`
}

func (lgc *Logics) GetObjectAsst(ctx context.Context, ownerID string) (map[string]string, errors.CCError) {
	// get host attribute info
	opt := hutil.NewOperation().WithOwnerID(ownerID).WithObjID(common.BKInnerObjIDHost).Data()
	attResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(ctx, lgc.header, opt)
	if err != nil || (err == nil && !attResult.Result) {
		return nil, fmt.Errorf("get host association failed, err, %v, %v", err, attResult.ErrMsg)
	}

	// key 是关联字段，val 是字段关联的模型ID
	attributes := make(map[string]string)
	for _, item := range attResult.Data {
		if item.PropertyType == common.FieldTypeSingleAsst || item.PropertyType == common.FieldTypeMultiAsst {
			opt := hutil.NewOperation().WithObjID(item.ObjectID).WithOwnerID(item.OwnerID).WithPropertyID(item.PropertyID).Data()
			res, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAssociations(ctx, lgc.header, opt)
			if err != nil {
				blog.Errorf("GetObjectAsst SelectObjectAssociations http do error, err:%s, input:%+v, rid:%s", err.Error(), opt, lgc.rid)
				return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			if err != nil {
				blog.Errorf("GetObjectAsst SelectObjectAssociations http reponse error, err code:%d, err msg:%s, input:%+v, rid:%s", res.Code, res.ErrMsg, opt, lgc.rid)
				return nil, lgc.ccErr.New(res.Code, res.ErrMsg)
			}

			if len(res.Data) > 0 {
				attributes[item.PropertyID] = res.Data[0].AsstObjID
			}
		}
	}

	return attributes, nil
}

func (lgc *Logics) getInstAsst(ctx context.Context, owerID, objID string, IDs []string, query *meta.QueryInput) ([]InstNameAsst, int, errors.CCError) {
	return lgc.getRawInstAsst(ctx, owerID, objID, IDs, query, false)

}

func (lgc *Logics) getInstAsstDetail(ctx context.Context, owerID, objID string, IDs []string, query *meta.QueryInput) ([]InstNameAsst, int, errors.CCError) {
	return lgc.getRawInstAsst(ctx, owerID, objID, IDs, query, true)
}

func (lgc *Logics) getRawInstAsst(ctx context.Context, ownerID, objID string, IDs []string, query *meta.QueryInput, isDetail bool) ([]InstNameAsst, int, errors.CCError) {
	var infos []mapstr.MapStr
	var count int
	var instName, instID string
	tmpIDs := []int{}
	for _, ID := range IDs {
		if "" == strings.TrimSpace(ID) {
			continue
		}
		tmpID, err := strconv.Atoi(ID)
		if nil != err {
			blog.Errorf("getRawInstAsst get objID(%s) inst id not integer, inst id:(%+v), rid:%s", objID, IDs, lgc.rid)
			return nil, 0, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, objID, "association id", "int", err.Error())
		}
		tmpIDs = append(tmpIDs, tmpID)
	}
	if 0 == len(tmpIDs) {
		return make([]InstNameAsst, 0), 0, nil
	}
	condition := mapstr.New()
	if nil != query.Condition {
		newCondtion, err := mapstr.NewFromInterface(query.Condition)
		if err != nil {
			blog.Errorf("getRawInstAsst get objID(%s) inst id not integer, inst id:(%+v), rid:%s", objID, IDs, lgc.rid)
			return nil, 0, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, objID, "query condition", "map[string]interface{}", err.Error())
		}
		condition = newCondtion
	}
	switch objID {
	case common.BKInnerObjIDHost:
		instName = common.BKHostInnerIPField
		instID = common.BKHostIDField
		if 0 != len(tmpIDs) {
			condition[common.BKHostIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.HostController().Host().GetHosts(ctx, lgc.header, query)
		if err != nil {
			blog.Errorf("getRawInstAsst GetHosts http do error, err:%s,input:%+v,rid:%s", err.Error(), query, lgc.rid)
			return nil, 0, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rtn.Result {
			blog.Errorf("getRawInstAsst GetHosts http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", rtn.Code, rtn.ErrMsg, query, lgc.rid)
			return nil, 0, lgc.ccErr.New(rtn.Code, rtn.ErrMsg)
		}
		infos = rtn.Data.Info
		count = rtn.Data.Count

	case common.BKInnerObjIDApp:
		instName = common.BKAppNameField
		instID = common.BKAppIDField
		if 0 != len(tmpIDs) {
			condition[common.BKAppIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDApp, lgc.header, query)
		if err != nil {
			blog.Errorf("getRawInstAsst SearchObjects http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rtn.Result {
			blog.Errorf("getRawInstAsst SearchObjects http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", rtn.Code, rtn.ErrMsg, objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.New(rtn.Code, rtn.ErrMsg)
		}
		infos = rtn.Data.Info
		count = rtn.Data.Count

	case common.BKInnerObjIDSet:
		instID = common.BKSetIDField
		instName = common.BKSetNameField
		query.Sort = common.BKSetIDField
		if 0 != len(tmpIDs) {
			condition[common.BKSetIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDSet, lgc.header, query)
		if err != nil {
			blog.Errorf("getRawInstAsst SearchObjects http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rtn.Result {
			blog.Errorf("getRawInstAsst SearchObjects http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", rtn.Code, rtn.ErrMsg, objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.New(rtn.Code, rtn.ErrMsg)
		}
		infos = rtn.Data.Info
		count = rtn.Data.Count
	case common.BKInnerObjIDModule:
		instID = common.BKModuleIDField
		instName = common.BKModuleNameField
		query.Sort = common.BKModuleIDField
		if 0 != len(tmpIDs) {
			condition[common.BKModuleIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKObjIDField, lgc.header, query)
		if err != nil {
			blog.Errorf("getRawInstAsst SearchObjects http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rtn.Result {
			blog.Errorf("getRawInstAsst SearchObjects http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", rtn.Code, rtn.ErrMsg, objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.New(rtn.Code, rtn.ErrMsg)
		}
		infos = rtn.Data.Info
		count = rtn.Data.Count

	case common.BKInnerObjIDPlat:
		instID = common.BKCloudIDField
		instName = common.BKCloudNameField
		query.Sort = common.BKCloudIDField
		if 0 != len(tmpIDs) {
			condition[common.BKCloudIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDPlat, lgc.header, query)
		if err != nil {
			blog.Errorf("getRawInstAsst SearchObjects http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rtn.Result {
			blog.Errorf("getRawInstAsst SearchObjects http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", rtn.Code, rtn.ErrMsg, objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.New(rtn.Code, rtn.ErrMsg)
		}
		infos = rtn.Data.Info
		count = rtn.Data.Count
	default:
		instName = common.BKInstNameField
		instID = common.BKInstIDField
		query.Sort = common.BKInstIDField
		condition[common.BKObjIDField] = objID
		if 0 != len(tmpIDs) {
			condition[common.BKInstIDField] = map[string]interface{}{"$in": tmpIDs}
		}
		query.Condition = condition
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDObject, lgc.header, query)
		if err != nil {
			blog.Errorf("getRawInstAsst SearchObjects http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rtn.Result {
			blog.Errorf("getRawInstAsst SearchObjects http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", rtn.Code, rtn.ErrMsg, objID, query, lgc.rid)
			return nil, 0, lgc.ccErr.New(rtn.Code, rtn.ErrMsg)
		}
		infos = rtn.Data.Info
		count = rtn.Data.Count
	}

	delarry := func(s []string, i int) []string {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}

	allInst := make([]InstNameAsst, 0)
	for _, info := range infos {
		if val, exist := info[instName]; exist {
			inst := InstNameAsst{}
			if name, can := val.(string); can {
				inst.Name = name
				inst.ObjID = objID
				if isDetail {
					inst.InstInfo = info
				}
			}

			if dataVal, exist := info[instID]; exist {

				itemInstID, err := util.GetInt64ByInterface(dataVal)
				if nil != err {
					blog.Errorf("not found assocte object ID %d from %v", instID, info)
					return nil, 0, fmt.Errorf("not found assocte object ID %d from %v", instID, info)
				}
				if 0 != len(IDs) {
					for idx, key := range IDs {
						if key == strconv.FormatInt(itemInstID, 10) {
							inst.ID = IDs[idx]
							inst.ObjectID, _ = util.GetInt64ByInterface(IDs[idx])
							IDs = delarry(IDs, idx)
							allInst = append(allInst, inst)
							goto next
						}
					}
				} else {
					inst.ID = strconv.FormatInt(itemInstID, 10)
					inst.ObjectID = itemInstID
					allInst = append(allInst, inst)
				}

			next:
			}
		}
	}

	// get the InstName name
	for _, ID := range IDs {
		allInst = append(allInst, InstNameAsst{ID: ID})
	}

	return allInst, count, nil
}

// get inst detail sub without association object detail
func (lgc *Logics) GetInstDetailsSub(ctx context.Context, objID, ownerID string, input []mapstr.MapStr, page meta.BasePage) ([]mapstr.MapStr, errors.CCError) {
	return lgc.getInstDetailsSub(ctx, objID, ownerID, input, page, false)
}

// get inst detail sub with association object detail
func (lgc *Logics) GetInstAsstDetailsSub(ctx context.Context, objID, ownerID string, input []mapstr.MapStr, page meta.BasePage) ([]mapstr.MapStr, errors.CCError) {
	return lgc.getInstDetailsSub(ctx, objID, ownerID, input, page, true)
}

func (lgc *Logics) getInstDetailsSub(ctx context.Context, objID, ownerID string, input []mapstr.MapStr, page meta.BasePage, isDetail bool) ([]mapstr.MapStr, errors.CCError) {
	asso, err := lgc.GetObjectAsst(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	for _, dataItem := range input {
		// key 是关联字段，val 是字段关联的模型ID
		for key, objID := range asso {
			if keyItem, exist := dataItem[key]; exist {
				if nil == keyItem {
					continue
				}
				keyItemStr := util.GetStrByInterface(keyItem)
				var retData []InstNameAsst
				var err error
				query := &meta.QueryInput{
					Start: page.Start,
					Limit: page.Limit,
					Sort:  page.Sort,
				}
				if "" == strings.TrimSpace(keyItemStr) {
					dataItem[key] = make([]InstNameAsst, 0)
					continue
				}
				if true == isDetail {
					retData, _, err = lgc.getInstAsstDetail(ctx, ownerID, objID, strings.Split(keyItemStr, ","), query)
				} else {
					retData, _, err = lgc.getInstAsst(ctx, ownerID, objID, strings.Split(keyItemStr, ","), query)
				}
				if err != nil {
					return nil, err
				}
				dataItem[key] = retData

			}
		}
	}

	return input, nil
}
