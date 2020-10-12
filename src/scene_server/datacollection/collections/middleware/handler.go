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

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
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

func (d *Discover) parseData(msg *string) (data map[string]interface{}, err error) {
	dataStr := gjson.Get(*msg, "data.data").String()
	if err = json.Unmarshal([]byte(dataStr), &data); err != nil {
		blog.Errorf("parse data error: %s", err)
		return
	}
	return
}

func (d *Discover) parseObjID(msg *string) string {
	return gjson.Get(*msg, "data.meta.model.bk_obj_id").String()
}

func (d *Discover) parseOwnerId(msg *string) string {
	ownerId := gjson.Get(*msg, "data.meta.model.bk_supplier_account").String()

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

	val, err := d.redisCli.Get(d.ctx, instKey).Result()
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
	_, err := d.redisCli.Set(d.ctx, key, value, duration).Result()
	if err != nil {
		blog.Warnf("%s: flush to redis failed: %s", key, err)
	} else {

		blog.Infof("%s: flush to redis success", key)
	}
}

func (d *Discover) TryUnsetRedis(key string) {
	_, err := d.redisCli.Del(d.ctx, key).Result()
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

func (d *Discover) UpdateOrCreateInst(msg *string) error {
	if msg == nil {
		return fmt.Errorf("message nil")
	}

	rid := util.GetHTTPCCRequestID(d.httpHeader)

	ownerID := d.parseOwnerId(msg)

	objID := d.parseObjID(msg)

	// get must check unique to judge if the instance exists
	cond := map[string]interface{}{
		common.BKObjIDField: objID,
		"must_check":        true,
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

	instIDField := common.GetInstIDField(objID)

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

		// add audit log.
		if err := func() error {
			// ready audit interface of instance.
			audit := auditlog.NewInstanceAudit(d.CoreAPI.CoreService())
			kit := &rest.Kit{
				Rid:             rid,
				Header:          d.httpHeader,
				Ctx:             d.ctx,
				CCError:         d.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(d.httpHeader)),
				User:            common.CCSystemCollectorUserName,
				SupplierAccount: common.BKDefaultOwnerID,
			}

			// generate audit log for create instance.
			data := []mapstr.MapStr{mapstr.NewFromMap(bodyData)}
			generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate).WithOperateFrom(metadata.FromDataCollection)
			auditLog, err := audit.GenerateAuditLog(generateAuditParameter, objID, data)
			if err != nil {
				blog.Errorf("generate instance audit log failed after create instance, objID: %s, err: %v, rid: %s",
					objID, err, rid)
				return err
			}

			// save audit log.
			if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
				blog.Errorf("save instance audit log failed after create instance, objID: %s, err: %v, rid: %s",
					objID, err, rid)
				return err
			}

			return nil
		}(); err != nil && blog.V(3) {
			blog.Errorf("save inst create audit log failed, err: %+v, rid: %s", err.Error(), rid)
		}
		return nil
	}

	instID, err := util.GetInt64ByInterface(inst[instIDField])
	if nil != err {
		return fmt.Errorf("get bk_inst_id failed: %s %s", inst[instIDField], err.Error())
	}

	dataChange := map[string]interface{}{}
	hasDiff := false
	for attrId, attrValue := range bodyData {

		if attrId == defaultRelateAttr {
			if relateList, ok := inst[defaultRelateAttr].([]interface{}); ok && len(relateList) == 1 {

				relateObj, ok := relateList[0].(map[string]interface{})

				if ok && (relateObj["id"] != "" && relateObj["id"] != "0" && relateObj["id"] != nil) {
					blog.Infof("skip updating single relation attr: [%s]=%v, since it is existed:%v.", defaultRelateAttr, attrValue, relateObj["id"])
				} else {
					if val, ok := attrValue.(string); ok && val != "" {
						dataChange[defaultRelateAttr] = val
						blog.Debug("[relation changed]  %s: %v ---> %v", defaultRelateAttr, "nil", dataChange[attrId])
						hasDiff = true
					}
				}

				continue
			}

			blog.Errorf("parse relation data failed, skip update: \n%v\n", inst[defaultRelateAttr])
			continue

		}
		if inst[attrId] != attrValue {
			dataChange[attrId] = attrValue
			blog.Debug("[changed]  %s: %v ---> %v", attrId, inst[attrId], dataChange[attrId])
			hasDiff = true
		}
	}

	if !hasDiff {
		blog.Infof("no need to update inst")
		return nil
	}

	// remove unchangeable fields.
	delete(inst, common.BKObjIDField)
	delete(inst, common.BKOwnerIDField)
	delete(inst, common.BKDefaultField)
	delete(inst, instIDField)
	delete(inst, common.LastTimeField)
	delete(inst, common.CreateTimeField)

	// ready audit interface of instance.
	audit := auditlog.NewInstanceAudit(d.CoreAPI.CoreService())
	kit := &rest.Kit{
		Rid:             rid,
		Header:          d.httpHeader,
		Ctx:             d.ctx,
		CCError:         d.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(d.httpHeader)),
		User:            common.CCSystemCollectorUserName,
		SupplierAccount: common.BKDefaultOwnerID,
	}

	// generate audit log before update instance.
	auditCond := map[string]interface{}{instIDField: instID}
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).
		WithOperateFrom(metadata.FromDataCollection).WithUpdateFields(inst)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, auditCond)
	if err != nil {
		blog.Errorf("generate instance audit log failed after create instance, objID: %s, err: %v, rid: %s",
			objID, err, rid)
		return err
	}

	// to update.
	input := metadata.UpdateOption{
		Data: dataChange,
		Condition: map[string]interface{}{
			instIDField: instID,
		},
		CanEditAll: true,
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

	// save audit log.
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("save instance audit log failed after create instance, objID: %s, err: %v, rid: %s",
			objID, err, rid)
		return err
	}

	return nil
}
