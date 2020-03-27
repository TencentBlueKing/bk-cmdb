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

package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/tidwall/gjson"
)

const (
	cacheTime = time.Minute * 5
)

const (
	defaultRelateAttr = "host"
)

func (d *Discover) parseData(msg string) (data map[string]interface{}, err error) {

	dataStr := gjson.Get(msg, "data.data").String()
	if err = json.Unmarshal([]byte(dataStr), &data); err != nil {
		blog.Errorf("parse data error: %s", err)
		return
	}
	return
}

func (d *Discover) parseObjID(msg string) string {
	return gjson.Get(msg, "data.meta.model.bk_obj_id").String()
}

func (d *Discover) parseOwnerId(msg string) string {
	ownerId := gjson.Get(msg, "data.meta.model.bk_supplier_account").String()

	if ownerId == "" {
		ownerId = common.BKDefaultOwnerID
	}
	return ownerId
}

func (d *Discover) CreateInstKey(objID string, ownerID string, val []string) string {
	return fmt.Sprintf("cc:v3:inst[%s:%s:%s:%s]",
		common.CCSystemCollectorUserName,
		ownerID,
		objID,
		strings.Join(val, ":"),
	)
}

func (d *Discover) GetInstFromRedis(instKey string) (map[string]interface{}, error) {

	val, err := d.redisCli.Get(instKey).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: get inst cache error: %s", instKey, err)
	}

	var cacheData = make(map[string]interface{})
	err = json.Unmarshal([]byte(val), &cacheData)
	if err != nil {
		return nil, fmt.Errorf("marshal condition error: %s", err)
	}

	return cacheData, nil

}

func (d *Discover) TrySetRedis(key string, value []byte, duration time.Duration) {
	_, err := d.redisCli.Set(key, value, duration).Result()
	if err != nil {
		blog.Warnf("%s: flush to redis failed: %s", key, err)
	} else {

		blog.Infof("%s: flush to redis success", key)
	}
}

func (d *Discover) TryUnsetRedis(key string) {
	_, err := d.redisCli.Del(key).Result()
	if err != nil {
		blog.Warnf("%s: remove from redis failed: %s", key, err)
	} else {
		blog.Infof("%s: remove from redis success", key)
	}
}

func (d *Discover) GetInst(ownerID, objID string, instKey string, cond map[string]interface{}) (map[string]interface{}, error) {
	rid := util.GetHTTPCCRequestID(d.httpHeader)
	instData, err := d.GetInstFromRedis(instKey)
	if err == nil {
		blog.Infof("inst exist in redis: %s", instKey)
		return instData, nil
	} else {
		blog.Errorf("get inst from redis error: %s", err)
	}

	resp, err := d.CoreAPI.CoreService().Instance().ReadInstance(d.ctx, d.httpHeader, objID, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search inst failed, cond: %s, error: %s, rid: %s", cond, err.Error(), rid)
		return nil, fmt.Errorf("search inst failed: %s", err.Error())
	}
	if !resp.Result {
		blog.Errorf("search inst failed, cond: %s, error message: %s, rid: %s", cond, resp.ErrMsg, rid)
		return nil, fmt.Errorf("search inst failed: %s", resp.ErrMsg)
	}

	if len(resp.Data.Info) > 0 {
		val, err := json.Marshal(resp.Data.Info[0])
		if err != nil {
			blog.Errorf("%s: flush to redis marshal failed: %s", instKey, err)
		}
		d.TrySetRedis(instKey, val, cacheTime)
		return resp.Data.Info[0], nil
	}

	return nil, nil
}

func (d *Discover) UpdateOrCreateInst(msg string) error {
	rid := util.GetHTTPCCRequestID(d.httpHeader)

	ownerID := d.parseOwnerId(msg)

	objID := d.parseObjID(msg)

	// get must check unique to judge if the instance exists
	cond := map[string]interface{}{
		common.BKObjIDField:   objID,
		"must_check":          true,
	}
	uniqueResp, err := d.CoreAPI.CoreService().Model().ReadModelAttrUnique(d.ctx, d.httpHeader, metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search model unique failed, cond: %s, error: %s, rid: %s", cond, err.Error(), rid)
		return fmt.Errorf("search model unique failed: %s", err.Error())
	}
	if !uniqueResp.Result {
		blog.Errorf("search model unique failed, cond: %s, error message: %s, rid: %s", cond, uniqueResp.ErrMsg, rid)
		return fmt.Errorf("search model unique failed: %s", uniqueResp.ErrMsg)
	}
	if uniqueResp.Data.Count != 1 {
		return fmt.Errorf("model %s has wrong must check unique num", objID)
	}
	keyIDs := make([]int64, 0)
	for _, key := range uniqueResp.Data.Info[0].Keys {
		keyIDs = append(keyIDs, int64(key.ID))
	}
	keys := make([]string, 0)
	cond = map[string]interface{}{
		common.BKObjIDField:   objID,
		common.BKOwnerIDField: ownerID,
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: keyIDs,
		},
	}
	attrResp, err := d.CoreAPI.CoreService().Model().ReadModelAttr(d.ctx, d.httpHeader, objID, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search model attribute failed, cond: %s, error: %s, rid: %s", cond, err.Error(), rid)
		return fmt.Errorf("search model attribute failed: %s", err.Error())
	}
	if !attrResp.Result {
		blog.Errorf("search model attribute failed, cond: %s, error message: %s, rid: %s", cond, attrResp.ErrMsg, rid)
		return fmt.Errorf("search model attribute failed: %s", attrResp.ErrMsg)
	}
	if attrResp.Data.Count <= 0 {
		blog.Errorf("unique model attribute count illegal, cond: %s, rid: %s", cond, rid)
		return fmt.Errorf("search model attribute failed: %s", attrResp.ErrMsg)
	}
	for _, attr := range attrResp.Data.Info {
		keys = append(keys, attr.PropertyID)
	}

	bodyData, err := d.parseData(msg)
	if err != nil {
		return fmt.Errorf("parse data error: %s", err)
	}

	cond = map[string]interface{}{
		common.BKObjIDField:   objID,
		common.BKOwnerIDField: ownerID,
	}
	valArr := make([]string, 0)
	for _, key := range keys {
		val := util.GetStrByInterface(bodyData[key])
		if val == "" {
			return fmt.Errorf("skip inst because of empty unique key %s value", key)
		}
		valArr = append(valArr, val)
		cond[key] = bodyData[key]
	}
	instKeyStr := d.CreateInstKey(objID, ownerID, valArr)
	inst, err := d.GetInst(ownerID, objID, instKeyStr, cond)
	if nil != err {
		return fmt.Errorf("get inst error: %s", err)
	}

	blog.Infof("get inst result: %v", inst)

	if len(inst) <= 0 {
		resp, err := d.CoreAPI.CoreService().Instance().CreateInstance(d.ctx, d.httpHeader, objID, &metadata.CreateModelInstance{Data: bodyData})
		if err != nil {
			blog.Errorf("search model failed %s", err.Error())
			return fmt.Errorf("search model failed: %s", err.Error())
		}
		if !resp.Result {
			blog.Errorf("search model failed %s", resp.ErrMsg)
			return fmt.Errorf("search model failed: %s", resp.ErrMsg)
		}
		blog.Infof("create inst result: %v", resp)
		instID := int64(resp.Data.Created.ID)

		if err := func() error {
			bizID, err := extensions.ParseBizID(bodyData)
			if err != nil {
				if blog.V(5) {
					blog.InfoJSON("ParseBizID from input data: %+v failed, err: %+v", bodyData)
				}
				return err
			}
			auditHeader, err := GetAuditLogHeader(d.CoreAPI, d.httpHeader, objID)
			if err != nil {
				blog.Errorf("GetAuditLogHeader failed, objID: %s, err: %s, rid: %s", objID, err.Error(), rid)
				return err
			}
			instIDField := common.GetInstIDField(objID)
			bodyData[instIDField] = instID
			auditLog := metadata.SaveAuditLogParams{
				ID:    instID,
				Model: objID,
				Content: metadata.Content{
					CurData: bodyData,
					PreData: nil,
					Headers: auditHeader,
				},
				OpDesc: "create " + objID,
				OpType: auditoplog.AuditOpTypeAdd,
				ExtKey: "",
				BizID:  bizID,
			}

			result, err := d.CoreAPI.CoreService().Audit().SaveAuditLog(d.ctx, d.httpHeader, auditLog)
			if err != nil {
				blog.Errorf("create inst audit log failed, http failed, err:%s, rid:%s", err.Error(), rid)
				return err
			}
			if !result.Result {
				blog.Errorf("create inst audit log failed, err code:%d, err msg:%s, rid:%s", result.Code, result.ErrMsg, rid)
				return err
			}
			return nil
		}(); err != nil && blog.V(3) {
			blog.Errorf("save inst create audit log failed, err: %+v, rid: %s", err.Error(), rid)
		}

		// update registry to iam
		if err := d.authManager.RegisterInstancesByID(d.ctx, d.httpHeader, objID, instID); err != nil {
			blog.Errorf("UpdateOrCreateInst success, but RegisterInstancesByID failed, objID: %s, instID: %d, err: %s, rid: %s", objID, instID, err, rid)
			return err
		}
		return nil
	}

	preUpdatedData, err := DeepCopyToMap(inst)
	if err != nil {
		blog.ErrorJSON("DeepCopyToMap pre updated inst failed, inst: %s, err: %s, rid: %s", inst, err.Error(), rid)
		// should not return here, it must try to to more jobs at it best
	}

	instID, err := util.GetInt64ByInterface(inst[common.BKInstIDField])
	if nil != err {
		return fmt.Errorf("get bk_inst_id failed: %s %s", inst[common.BKInstIDField], err.Error())
	}

	hasDiff := false
	for attrId, attrValue := range bodyData {

		if attrId == defaultRelateAttr {
			if relateList, ok := inst[defaultRelateAttr].([]interface{}); ok && len(relateList) == 1 {

				relateObj, ok := relateList[0].(map[string]interface{})

				if ok && (relateObj["id"] != "" && relateObj["id"] != "0" && relateObj["id"] != nil) {
					blog.Infof("skip update exist single relation attr: %s->%v", attrId, attrValue)
				} else if attrValue != "" {
					blog.Debug("[relation changed]  %s: %v ---> %v", attrId, attrValue)
					inst[attrId] = attrValue
					hasDiff = true
				}

				continue
			}

			blog.Errorf("parse relation data failed, skip update: \n%v\n", inst[defaultRelateAttr])
			continue

		}

		if inst[attrId] != attrValue {
			inst[attrId] = attrValue
			blog.Debug("[changed]  %s: %v ---> %v", attrId, attrValue, inst[attrId])
			hasDiff = true
		}
	}

	if !hasDiff {
		blog.Infof("no need to update inst")
		return nil
	}

	delete(inst, common.BKObjIDField)
	delete(inst, common.BKOwnerIDField)
	delete(inst, common.BKDefaultField)
	delete(inst, common.BKInstIDField)
	delete(inst, common.LastTimeField)
	delete(inst, common.CreateTimeField)

	input := metadata.UpdateOption{
		Data: inst,
		Condition: map[string]interface{}{
			common.BKInstIDField: instID,
		},
	}
	resp, err := d.CoreAPI.CoreService().Instance().UpdateInstance(d.ctx, d.httpHeader, objID, &input)
	if err != nil {
		blog.Errorf("search model failed %s", err.Error())
		return fmt.Errorf("search model failed: %s", err.Error())
	}
	if !resp.Result {
		blog.Errorf("search model failed %s", resp.ErrMsg)
		return fmt.Errorf("search model failed: %s", resp.ErrMsg)
	}
	blog.Infof("update inst result: %v", resp)
	d.TryUnsetRedis(instKeyStr)

	if err := func() error {
		bizID, err := extensions.ParseBizID(inst)
		if err != nil {
			if blog.V(5) {
				blog.InfoJSON("ParseBizID from input data: %+v failed, err: %+v", inst)
			}
			return err
		}
		auditHeader, err := GetAuditLogHeader(d.CoreAPI, d.httpHeader, objID)
		if err != nil {
			blog.Errorf("GetAuditLogHeader failed, objID: %s, err: %s, rid: %s", objID, err.Error(), rid)
			return err
		}
		qc := metadata.QueryCondition{
			Condition: map[string]interface{}{
				common.BKInstIDField: instID,
			},
		}
		readResult, err := d.CoreAPI.CoreService().Instance().ReadInstance(d.ctx, d.httpHeader, objID, &qc)
		if err != nil {
			blog.Errorf("read updated inst failed, objID: %s, instID: %d, err: %s, rid: %s", objID, instID, err.Error(), rid)
			return err
		}
		if len(readResult.Data.Info) == 0 {
			blog.Errorf("read updated inst failed, not found, objID: %s, instID: %d, rid: %s", objID, instID, rid)
			return fmt.Errorf("get updated instance failed, not found, objID: %s, instID: %d", objID, instID)
		}
		updatedInst := readResult.Data.Info[0]
		auditLog := metadata.SaveAuditLogParams{
			ID:    instID,
			Model: objID,
			Content: metadata.Content{
				CurData: updatedInst,
				PreData: preUpdatedData,
				Headers: auditHeader,
			},
			OpDesc: "update " + objID,
			OpType: auditoplog.AuditOpTypeModify,
			ExtKey: "",
			BizID:  bizID,
		}

		result, err := d.CoreAPI.CoreService().Audit().SaveAuditLog(d.ctx, d.httpHeader, auditLog)
		if err != nil {
			blog.Errorf("save inst update audit log failed, http failed, err:%s, rid:%s", err.Error(), rid)
			return err
		}
		if !result.Result {
			blog.Errorf("save inst update audit log failed, err code:%d, err msg:%s, rid:%s", result.Code, result.ErrMsg, rid)
			return fmt.Errorf("coreservice save audit log result faield, result: %+v", result)
		}
		return nil
	}(); err != nil && blog.V(3) {
		blog.Errorf("save inst update audit log failed, err: %s, rid: %s", err.Error(), rid)
	}

	// update registry to iam
	if err := d.authManager.UpdateRegisteredInstanceByID(d.ctx, d.httpHeader, objID, instID); err != nil {
		blog.Errorf("UpdateOrCreateInst success, but UpdateRegisteredInstanceByID failed, objID: %s, instID: %d, err: %s, rid: %s", objID, instID, err, rid)
		return err
	}

	return nil
}
