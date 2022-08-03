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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) getAssociationData(ctx context.Context, header http.Header, objID string,
	instAsstArr []*metadata.InstAsst, modelBizID int64, AsstObjectUniqueIDMap map[string]int64, selfObjectUniqueID int64) (
	map[string]map[int64][]PropertyPrimaryVal, map[int64][]PropertyPrimaryVal, error) {

	// 当前操作对象实例id
	instIDArr := make([]int64, 0)
	// map[objID][]instID
	asstObjIDInstIDArr := make(map[string][]int64)
	for _, instAsst := range instAsstArr {
		if instAsst.ObjectID == objID {
			instIDArr = append(instIDArr, instAsst.InstID)
			_, ok := asstObjIDInstIDArr[instAsst.AsstObjectID]
			if !ok {
				asstObjIDInstIDArr[instAsst.AsstObjectID] = make([]int64, 0)
			}
			asstObjIDInstIDArr[instAsst.AsstObjectID] = append(asstObjIDInstIDArr[instAsst.AsstObjectID], instAsst.AsstInstID)
		} else {
			instIDArr = append(instIDArr, instAsst.AsstInstID)
			_, ok := asstObjIDInstIDArr[instAsst.ObjectID]
			if !ok {
				asstObjIDInstIDArr[instAsst.ObjectID] = make([]int64, 0)
			}
			asstObjIDInstIDArr[instAsst.ObjectID] = append(asstObjIDInstIDArr[instAsst.ObjectID], instAsst.InstID)
		}

	}

	// map[objID]map[inst_id][]Property
	retAsstObjIDInstInfoMap := make(map[string]map[int64][]PropertyPrimaryVal)
	for itemObjID, asstInstIDArr := range asstObjIDInstIDArr {
		uniqueID := AsstObjectUniqueIDMap[itemObjID]
		objPrimaryInfo, err := lgc.fetchInstAssociationData(ctx, header, itemObjID, asstInstIDArr, modelBizID, uniqueID)
		if err != nil {
			return nil, nil, err
		}
		retAsstObjIDInstInfoMap[itemObjID] = objPrimaryInfo
	}

	objPrimaryInfo, err := lgc.fetchInstAssociationData(ctx, header, objID, instIDArr, modelBizID, selfObjectUniqueID)
	if err != nil {
		return nil, nil, err
	}

	return retAsstObjIDInstInfoMap, objPrimaryInfo, err
}

// fetchAssociationData TODO
func (lgc *Logics) fetchAssociationData(ctx context.Context, header http.Header, objID string, instIDArr []int64,
	modelBizID int64, asstIDArr []string, hasSelfAssociation bool) ([]*metadata.InstAsst, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	input := &metadata.SearchAssociationInstRequest{ObjID: objID}

	if len(asstIDArr) == 0 {
		blog.Infof("empty  bk_obj_asst_id. obj: %s, rid: %s", objID, rid)
		return nil, nil
	}
	// 实例作为关联关系源模型
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKInstIDField).In(instIDArr)
	cond.Field(common.AssociationObjAsstIDField).In(asstIDArr)
	if !hasSelfAssociation {
		// 不包含自关联
		cond.Field(common.BKAsstObjIDField).NotEq(objID)
	}
	input.Condition = cond.ToMapStr()
	bkObjRst, err := lgc.CoreAPI.ApiServer().SearchAssociationInst(ctx, header, input)
	if err != nil {
		blog.Errorf("fetch %s association data failed, input: %+v, err: %v, rid: %s", objID, input, err, rid)
		return nil, ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := bkObjRst.CCError(); ccErr != nil {
		blog.Errorf("fetch %s association data failed, input: %+v, err: %v, rid: %s", objID, input, ccErr, rid)
		return nil, ccErr
	}

	// 实例作为关联关系目标模型
	cond = condition.CreateCondition()
	cond.Field(common.BKAsstObjIDField).Eq(objID)
	cond.Field(common.BKAsstInstIDField).In(instIDArr)
	cond.Field(common.AssociationObjAsstIDField).In(asstIDArr)
	// 自关联数据，上面已经处理，只需要拉取一次数据即可
	cond.Field(common.BKObjIDField).NotEq(objID)

	input.Condition = cond.ToMapStr()
	bkAsstObjRst, err := lgc.CoreAPI.ApiServer().SearchAssociationInst(ctx, header, input)
	if err != nil {
		blog.Errorf("fetch %s association data failed, input: %+v, err: %v, rid: %s", objID, input, err, rid)
		return nil, ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := bkObjRst.CCError(); ccErr != nil {
		blog.Errorf("fetch %s association data failed, input: %+v, err: %v, rid: %s", objID, input, ccErr, rid)
		return nil, ccErr
	}
	result := append(bkObjRst.Data[:], bkAsstObjRst.Data...)

	return result, nil
}

func (lgc *Logics) fetchInstAssociationData(ctx context.Context, header http.Header, objID string, instIDArr []int64,
	modelBizID int64, uniqueID int64) (map[int64][]PropertyPrimaryVal, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	instIDKey := metadata.GetInstIDFieldByObjID(objID)
	insts := make([]mapstr.MapStr, 0)
	option := mapstr.MapStr{
		"condition": mapstr.MapStr{
			instIDKey: mapstr.MapStr{
				common.BKDBIN: instIDArr,
			},
		},
	}

	resp, err := lgc.CoreAPI.ApiServer().GetInstUniqueFields(ctx, header, objID, uniqueID, option)
	if err != nil {
		blog.ErrorJSON("fetchInstAssociationData failed, GetInstUniqueFields err:%v, option: %s, rid: %s", err, option, rid)
		return nil, ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.ErrorJSON("fetchInstAssociationData failed, GetInstUniqueFields resp:%s, option: %s, rid: %s", resp, option, rid)
		return nil, resp.CCError()
	}
	insts = resp.Data.Info

	retAsstInstInfo := make(map[int64][]PropertyPrimaryVal, 0)
	for _, inst := range insts {
		instID, err := inst.Int64(instIDKey)
		if err != nil {
			blog.Warnf("fetchInstAssociationData get %s instance %s field error, err:%s, inst:%+v, rid:%s", objID, instIDKey, err.Error(), inst, rid)
			continue
		}
		isSkip := false
		var primaryKeysVal []PropertyPrimaryVal
		for attrID, attrName := range resp.Data.UniqueAttribute {
			// use display , use string
			val, err := inst.String(attrID)
			if err != nil {
				blog.WarnJSON("fetchInstAssociationData get %s instance %s field error, err:%s, inst:%s, rid:%s",
					objID, attrID, err.Error(), inst, rid)
				isSkip = true
				break
			}
			primaryKeysVal = append(primaryKeysVal, PropertyPrimaryVal{
				ID:     attrID,
				Name:   attrName,
				StrVal: val,
			})
		}
		if isSkip {
			continue
		}
		retAsstInstInfo[instID] = primaryKeysVal

	}
	return retAsstInstInfo, nil

}
