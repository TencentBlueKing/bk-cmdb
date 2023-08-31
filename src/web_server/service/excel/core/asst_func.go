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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// GetObjAssociation get object association
func (d *Client) GetObjAssociation(kit *rest.Kit, objID string) ([]*metadata.Association, error) {
	cond := &metadata.SearchAssociationObjectRequest{
		Condition: map[string]interface{}{
			condition.BKDBOR: []mapstr.MapStr{
				{
					common.BKObjIDField: objID,
				},
				{
					common.BKAsstObjIDField: objID,
				},
			},
		},
	}

	// 确定关联标识的列表，定义excel选项下拉栏。此处需要查cc_ObjAsst表。
	resp, err := d.ApiClient.SearchObjectAssociation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.ErrorJSON("get object association list failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err := resp.CCError(); err != nil {
		blog.ErrorJSON("get object association list failed, err: %v, rid: %s", resp.ErrMsg, kit.Rid)
		return nil, err
	}

	return resp.Data, nil
}

// GetAsstOpFlag get association operator flag
func GetAsstOpFlag(op AsstOp) metadata.ExcelAssociationOperate {
	opFlag := metadata.ExcelAssociationOperateError
	switch op {
	case AsstOpAdd:
		opFlag = metadata.ExcelAssociationOperateAdd
	case AsstOpDel:
		opFlag = metadata.ExcelAssociationOperateDelete
	}

	return opFlag
}

// FindAsstByAsstID find association by object association id
func (d *Client) FindAsstByAsstID(kit *rest.Kit, objID string, asstIDs []string) ([]metadata.Association, error) {
	input := metadata.FindAssociationByObjectAssociationIDRequest{
		ObjAsstIDArr: asstIDs,
	}
	resp, err := d.ApiClient.FindAssociationByObjectAssociationID(kit.Ctx, kit.Header, objID, input)
	if err != nil {
		blog.ErrorJSON("find model association by bk_obj_asst_id failed, err: %s, objID: %s, input: %s, rid: %s",
			err.Error(), objID, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if ccErr := resp.CCError(); ccErr != nil {
		blog.ErrorJSON("find model association by bk_obj_asst_id failed, err: %v, objID: %s, input: %s, rid: %s", ccErr,
			objID, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	return resp.Data, nil
}

// ImportAssociation import association
func (d *Client) ImportAssociation(kit *rest.Kit, objID string, input *metadata.RequestImportAssociation) (
	*metadata.ResponeImportAssociationData, error) {

	resp, err := d.ApiClient.ImportAssociation(kit.Ctx, kit.Header, objID, input)
	if err != nil {
		blog.Errorf("import association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	if err := resp.CCError(); err != nil {
		blog.Errorf("import association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return &resp.Data, nil
}

// GetInstAsst get instance associations
func (d *Client) GetInstAsst(kit *rest.Kit, objID string, instIDs []int64, asstIDs []string, hasSelfAsst bool) (
	[]*metadata.InstAsst, error) {

	if len(instIDs) == 0 || len(asstIDs) == 0 {
		return nil, nil
	}

	// 实例作为关联关系源模型
	cond := mapstr.MapStr{
		common.BKObjIDField:              objID,
		common.BKInstIDField:             mapstr.MapStr{common.BKDBIN: instIDs},
		common.AssociationObjAsstIDField: mapstr.MapStr{common.BKDBIN: asstIDs},
	}
	if !hasSelfAsst {
		// 不包含自关联
		cond[common.BKAsstObjIDField] = mapstr.MapStr{common.BKDBNE: objID}
	}
	input := &metadata.SearchAssociationInstRequest{ObjID: objID, Condition: cond}

	objRes, err := d.ApiClient.SearchAssociationInst(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("fetch %s association data failed, input: %+v, err: %v, rid: %s", objID, input, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := objRes.CCError(); ccErr != nil {
		blog.Errorf("fetch %s association data failed, input: %+v, err: %v, rid: %s", objID, input, ccErr, kit.Rid)
		return nil, ccErr
	}

	// 实例作为关联关系目标模型
	cond = mapstr.MapStr{
		common.BKAsstObjIDField:          objID,
		common.BKInstIDField:             mapstr.MapStr{common.BKDBIN: instIDs},
		common.AssociationObjAsstIDField: mapstr.MapStr{common.BKDBIN: asstIDs},
		// 自关联数据，上面已经处理，只需要拉取一次数据即可
		common.BKObjIDField: mapstr.MapStr{common.BKDBNE: objID},
	}
	input.Condition = cond

	asstObjRes, err := d.ApiClient.SearchAssociationInst(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("fetch %s association data failed, input: %+v, err: %v, rid: %s", objID, input, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := asstObjRes.CCError(); ccErr != nil {
		blog.Errorf("fetch %s association data failed, input: %+v, err: %v, rid: %s", objID, input, ccErr, kit.Rid)
		return nil, ccErr
	}

	result := append(objRes.Data[:], asstObjRes.Data...)

	return result, nil
}

// GetInstUniqueKeys get instance unique keys
func (d *Client) GetInstUniqueKeys(kit *rest.Kit, objID string, instIDs []int64, uniqueID int64) (
	map[int64]string, error) {

	instIDKey := metadata.GetInstIDFieldByObjID(objID)
	option := mapstr.MapStr{"condition": mapstr.MapStr{instIDKey: mapstr.MapStr{common.BKDBIN: instIDs}}}

	resp, err := d.ApiClient.GetInstUniqueFields(kit.Ctx, kit.Header, objID, uniqueID, option)
	if err != nil {
		blog.Errorf("get instance unique fields failed, err: %v, option: %v, rid: %s", err, option, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if resp.CCError() != nil {
		blog.ErrorJSON("get inst unique fields failed, err: %v, option: %s, rid: %s", resp.CCError(), option, kit.Rid)
		return nil, resp.CCError()
	}

	uniqueAttr := make(map[string]string)
	for id, name := range resp.Data.UniqueAttribute {
		uniqueAttr[id] = name
	}

	result := make(map[int64]string, 0)
	for _, inst := range resp.Data.Info {
		instID, err := inst.Int64(instIDKey)
		if err != nil {
			blog.Warnf("get %s instance %s field failed, err: %v, inst: %+v, rid:%s", objID, instIDKey, err, inst,
				kit.Rid)
			continue
		}

		isSkip := false
		uniqueKeys := make([]string, 0)
		for attrID, attrName := range uniqueAttr {
			val, err := inst.String(attrID)
			if err != nil {
				blog.WarnJSON("get %s instance %s field failed, err: %v, inst: %v, rid: %s", objID, attrID, err, inst,
					kit.Rid)
				isSkip = true
				break
			}

			key := attrName + common.ExcelAsstPrimaryKeyJoinChar + val
			uniqueKeys = append(uniqueKeys, key)
		}

		if isSkip {
			continue
		}

		result[instID] = strings.Join(uniqueKeys, common.ExcelAsstPrimaryKeySplitChar)

	}

	return result, nil
}
