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

package core

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// GetCustomTopoBriefMsg get custom topo brief message
func (d *Client) GetCustomTopoBriefMsg(kit *rest.Kit) ([]TopoBriefMsg, error) {
	objIDs, err := d.getCustomTopoObjIDs(kit)
	if err != nil {
		blog.Errorf("get custom topo objID failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(objIDs) == 0 {
		return nil, nil
	}

	input := &metadata.QueryCondition{
		Fields:    []string{common.BKObjIDField, common.BKObjNameField},
		Condition: mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}},
	}

	objResult, err := d.Engine.CoreAPI.ApiServer().ReadModel(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search mainline obj failed, objIDs: %#v, err: %v, rid: %s", objIDs, err, kit.Rid)
		return nil, err
	}
	idNameMap := make(map[string]string, len(objResult.Info))
	for _, object := range objResult.Info {
		idNameMap[object.ObjectID] = object.ObjectName
	}

	result := make([]TopoBriefMsg, len(objIDs))
	for idx, objID := range objIDs {
		result[len(objIDs)-1-idx] = TopoBriefMsg{ObjID: objID, Name: idNameMap[objID]}
	}

	return result, nil
}

func (d *Client) getCustomTopoObjIDs(kit *rest.Kit) ([]string, error) {
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline},
	}
	mainlineAsstRsp, err := d.Engine.CoreAPI.ApiServer().ReadModuleAssociation(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("search mainline association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	mainlineObjChildMap := make(map[string]string, 0)
	for _, asst := range mainlineAsstRsp.Info {
		if asst.ObjectID == common.BKInnerObjIDHost {
			continue
		}
		mainlineObjChildMap[asst.AsstObjID] = asst.ObjectID
	}

	// get all mainline custom object id
	objIDs := make([]string, 0)
	for objectID := common.BKInnerObjIDApp; len(objectID) != 0; objectID = mainlineObjChildMap[objectID] {
		if objectID == common.BKInnerObjIDApp || objectID == common.BKInnerObjIDSet ||
			objectID == common.BKInnerObjIDModule {
			continue
		}

		objIDs = append(objIDs, objectID)
	}

	return util.ReverseArrayString(objIDs), nil
}

// GetDefaultBizID get resource pool biz ID
func (d *Client) GetDefaultBizID(kit *rest.Kit) (int64, error) {
	resp, err := d.Engine.CoreAPI.ApiServer().SearchDefaultApp(kit.Ctx, kit.Header, kit.SupplierAccount)
	if err != nil {
		blog.Errorf("search default bizID failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.Errorf("search default bizID failed, err: %v, rid: %s", err, kit.Rid)
		return 0, resp.CCError()
	}

	if len(resp.Data.Info) == 0 {
		return 0, kit.CCError.CCError(common.CCErrHostNotResourceFail)
	}
	bizID, err := resp.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		return 0, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp,
			common.BKAppIDField, "int64", err.Error())
	}

	return bizID, nil
}

// GetHostBizRelations get host and biz relations
func (d *Client) GetHostBizRelations(kit *rest.Kit, hosts map[int]map[string]interface{}) (map[int64]int64, error) {
	var hostIDs []int64
	for _, host := range hosts {
		hostID, ok := host[common.BKHostIDField]
		if !ok {
			continue
		}
		if hostIDVal, err := util.GetInt64ByInterface(hostID); err == nil {
			hostIDs = append(hostIDs, hostIDVal)
		}
	}

	hostLen := len(hostIDs)
	if hostLen == 0 {
		return make(map[int64]int64), nil
	}

	hostBizMap := make(map[int64]int64)
	// the length of GetHostModuleRelation's param bk_host_id can't bigger than 500
	for idx := 0; idx < hostLen; idx += 500 {
		endIdx := idx + 500
		if endIdx > hostLen {
			endIdx = hostLen
		}
		ids := hostIDs[idx:endIdx]

		params := mapstr.MapStr{
			common.BKHostIDField: ids,
		}
		resp, err := d.Engine.CoreAPI.ApiServer().GetHostModuleRelation(kit.Ctx, kit.Header, params)
		if err != nil {
			blog.Errorf("get host module relation failed, err: %v, params: %#v, rid: %s", err, params, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !resp.Result {
			blog.Errorf("get host module relation failed, err: %v, params: %#v, rid: %s", err, params, kit.Rid)
			return nil, resp.CCError()
		}

		for _, relation := range resp.Data {
			hostBizMap[relation.HostID] = relation.AppID
		}
	}

	return hostBizMap, nil
}
