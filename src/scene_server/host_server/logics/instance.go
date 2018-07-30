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
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

type InstNameAsst struct {
	ID         string                 `json:"id"`
	ObjID      string                 `json:"bk_obj_id"`
	ObjIcon    string                 `json:"bk_obj_icon"`
	ObjectID   int                    `json:"bk_inst_id"`
	ObjectName string                 `json:"bk_obj_name"`
	Name       string                 `json:"bk_inst_name"`
	InstInfo   map[string]interface{} `json:"inst_info,omitempty"`
}

func (lgc *Logics) GetObjectAsst(ownerID string, pheader http.Header) (map[string]string, error) {
	// get host attribute info
	opt := hutil.NewOperation().WithOwnerID(ownerID).WithObjID(common.BKInnerObjIDHost).Data()
	attResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, opt)
	if err != nil || (err == nil && !attResult.Result) {
		return nil, fmt.Errorf("get host association failed, err, %v, %v", err, attResult.ErrMsg)
	}

	// key 是关联字段，val 是字段关联的模型ID
	attributes := make(map[string]string)
	for _, item := range attResult.Data {
		if item.PropertyType == common.FieldTypeSingleAsst || item.PropertyType == common.FieldTypeMultiAsst {
			opt := hutil.NewOperation().WithObjID(item.ObjectID).WithOwnerID(item.OwnerID).WithPropertyID(item.PropertyID).Data()
			res, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAssociations(context.Background(), pheader, opt)
			if err != nil || (err == nil && !res.Result) {
				return nil, fmt.Errorf("get host pre data failed, err, %v, %v", err, res.ErrMsg)
			}
			if len(res.Data) > 0 {
				attributes[item.PropertyID] = res.Data[0].AsstObjID
			}
		}
	}

	return attributes, nil
}

func (lgc *Logics) getInstAsst(owerID, objID string, IDs []string, pheader http.Header, query *meta.QueryInput) ([]InstNameAsst, int, error) {
	return lgc.getRawInstAsst(owerID, objID, IDs, pheader, query, false)

}

func (lgc *Logics) getInstAsstDetail(owerID, objID string, IDs []string, pheader http.Header, query *meta.QueryInput) ([]InstNameAsst, int, error) {
	return lgc.getRawInstAsst(owerID, objID, IDs, pheader, query, true)
}

func (lgc *Logics) getRawInstAsst(ownerID, objID string, IDs []string, pheader http.Header, query *meta.QueryInput, isDetail bool) ([]InstNameAsst, int, error) {
	var infos []map[string]interface{}
	var count int
	var instName, instID string
	tmpIDs := []int{}
	for _, ID := range IDs {
		tmpID, err := strconv.Atoi(ID)
		if nil != err {
			return nil, 0, fmt.Errorf("assocate id not integer, ids:%v", strings.Join(IDs, ","))
		}
		tmpIDs = append(tmpIDs, tmpID)
	}
	condition := make(map[string]interface{})
	if nil != query.Condition {
		var ok bool
		condition, ok = query.Condition.(map[string]interface{})
		if false == ok {
			return nil, 0, fmt.Errorf("assocate id not integer, ids:%v", strings.Join(IDs, ","))
		}
	}
	switch objID {
	case common.BKInnerObjIDHost:
		instName = common.BKHostInnerIPField
		instID = common.BKHostIDField
		if 0 != len(tmpIDs) {
			condition[common.BKHostIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.HostController().Host().GetHosts(context.Background(), pheader, query)
		if err != nil || (err == nil && !rtn.Result) {
			return nil, 0, fmt.Errorf("get hosts failed, err, %v, %v", err, rtn.ErrMsg)
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
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, pheader, query)
		if err != nil || (err == nil && !rtn.Result) {
			return nil, 0, fmt.Errorf("get hosts failed, err, %v, %v", err, rtn.ErrMsg)
		}
		for _, tmp := range rtn.Data.Info {
			infos = append(infos, map[string]interface{}(tmp))
		}
		count = rtn.Data.Count

	case common.BKInnerObjIDSet:
		instID = common.BKSetIDField
		instName = common.BKSetNameField
		query.Sort = common.BKSetIDField
		if 0 != len(tmpIDs) {
			condition[common.BKSetIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, pheader, query)
		if err != nil || (err == nil && !rtn.Result) {
			return nil, 0, fmt.Errorf("get hosts failed, err, %v, %v", err, rtn.ErrMsg)
		}
		for _, tmp := range rtn.Data.Info {
			infos = append(infos, map[string]interface{}(tmp))
		}
		count = rtn.Data.Count

	case common.BKInnerObjIDModule:
		instID = common.BKModuleIDField
		instName = common.BKModuleNameField
		query.Sort = common.BKModuleIDField
		if 0 != len(tmpIDs) {
			condition[common.BKModuleIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKObjIDField, pheader, query)
		if err != nil || (err == nil && !rtn.Result) {
			return nil, 0, fmt.Errorf("get hosts failed, err, %v, %v", err, rtn.ErrMsg)
		}
		for _, tmp := range rtn.Data.Info {
			infos = append(infos, map[string]interface{}(tmp))
		}
		count = rtn.Data.Count

	case common.BKInnerObjIDPlat:
		instID = common.BKCloudIDField
		instName = common.BKCloudNameField
		query.Sort = common.BKCloudIDField
		if 0 != len(tmpIDs) {
			condition[common.BKCloudIDField] = map[string]interface{}{"$in": tmpIDs}
			query.Condition = condition
		}
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDPlat, pheader, query)
		if err != nil || (err == nil && !rtn.Result) {
			return nil, 0, fmt.Errorf("get plat failed, err, %v, %v", err, rtn.ErrMsg)
		}

		for _, tmp := range rtn.Data.Info {
			infos = append(infos, map[string]interface{}(tmp))
		}
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
		rtn, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKINnerObjIDObject, pheader, query)
		if err != nil || (err == nil && !rtn.Result) {
			return nil, 0, fmt.Errorf("get hosts failed, err, %v, %v", err, rtn.ErrMsg)
		}
		for _, tmp := range rtn.Data.Info {
			infos = append(infos, map[string]interface{}(tmp))
		}
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
							inst.ObjectID, _ = strconv.Atoi(IDs[idx])
							IDs = delarry(IDs, idx)
							allInst = append(allInst, inst)
							goto next
						}
					}
				} else {
					inst.ID = strconv.FormatInt(itemInstID, 10)
					inst.ObjectID = int(itemInstID)
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
func (lgc *Logics) GetInstDetailsSub(pheader http.Header, objID, ownerID string, input []map[string]interface{}, page meta.BasePage) ([]map[string]interface{}, error) {
	return lgc.getInstDetailsSub(pheader, objID, ownerID, input, page, false)
}

// get inst detail sub with association object detail
func (lgc *Logics) GetInstAsstDetailsSub(pheader http.Header, objID, ownerID string, input []map[string]interface{}, page meta.BasePage) ([]map[string]interface{}, error) {
	return lgc.getInstDetailsSub(pheader, objID, ownerID, input, page, true)
}

func (lgc *Logics) getInstDetailsSub(pheader http.Header, objID, ownerID string, input []map[string]interface{}, page meta.BasePage, isDetail bool) ([]map[string]interface{}, error) {
	asso, err := lgc.GetObjectAsst(ownerID, pheader)
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
				if true == isDetail {

					retData, _, err = lgc.getInstAsstDetail(ownerID, objID, strings.Split(keyItemStr, ","), pheader, query)
				} else {
					retData, _, err = lgc.getInstAsst(ownerID, objID, strings.Split(keyItemStr, ","), pheader, query)
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
