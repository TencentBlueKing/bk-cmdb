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
	bkc "configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
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
	defaultModelIcon  = "icon-cc-middleware"
)

type Model struct {
	BkClassificationID string `json:"bk_classification_id"`
	BkObjID            string `json:"bk_obj_id"`
	BkObjName          string `json:"bk_obj_name"`
	Keys               string `json:"bk_obj_keys"`
}

type Attr struct {
	ID            int    `json:"id"`
	OwnerID       string `json:"bk_supplier_account"`
	ObjID         string `json:"bk_obj_id"`
	PropertyGroup string `json:"bk_property_group"`

	PropertyID    string `json:"bk_property_id"`
	PropertyName  string `json:"bk_property_name"`
	PropertyType  string `json:"bk_property_type"`
	AssociationID string `json:"bk_asst_obj_id"`

	Option  interface{} `json:"option"`
	Creator string      `json:"creator"`

	Editable    bool   `json:"editable"`
	IsRequired  bool   `json:"isrequired"`
	IsReadOnly  bool   `json:"isreadonly"`
	IsOnly      bool   `json:"isonly"`
	Description string `json:"description"`
}

type Related struct {
	BkInstId   int    `json:"bk_inst_id"`
	BkInstName string `json:"bk_inst_name"`
	BkObjIcon  string `json:"bk_obj_icon"`
	BkObjId    string `json:"bk_obj_id"`
	BkObjName  string `json:"bk_obj_name"`
	Id         string `json:"id"`
}

type SingleRelated []Related

type M map[string]interface{}

type MapData M

type ResultBase struct {
	Result  bool   `json:"result"`
	Code    int    `json:"bk_error_code"`
	Message string `json:"bk_err_message"`
}

type DetailResult struct {
	ResultBase
	Data struct {
		Count int       `json:"count"`
		Info  []MapData `json:"info"`
	} `json:"data"`
}

type ListResult struct {
	ResultBase
	Data []MapData `json:"data"`
}

type Result struct {
	ResultBase
	Data interface{} `json:"data"`
}

func (m M) toJson() ([]byte, error) {
	return json.Marshal(m)
}

func (m M) debug() {
	if js, err := m.toJson(); err == nil {
		blog.Infof("=====\n%s\n====", js)
	} else {
		blog.Errorf("debug error: %s", err)
	}
}

func (m M) Keys() (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}

	return
}

func (r *Result) mapData() (MapData, error) {
	if m, ok := r.Data.(MapData); ok {
		return m, nil
	}
	return nil, fmt.Errorf("parse map data error: %v", r)
}

func parseListResult(res []byte) (ListResult, error) {
	var lR ListResult

	if err := json.Unmarshal(res, &lR); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s", err)
		return lR, err
	}

	return lR, nil
}

func parseDetailResult(res []byte) (DetailResult, error) {
	var dR DetailResult

	if err := json.Unmarshal(res, &dR); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s", err)
		return dR, err
	}

	return dR, nil
}

func parseResult(res []byte) (Result, error) {

	var r Result

	if err := json.Unmarshal(res, &r); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s", err)
		return r, err
	}

	return r, nil
}

func (d *Discover) parseModel(msg string) (model *Model, err error) {

	model = &Model{}
	modelStr := gjson.Get(msg, "data.meta.model").String()

	if err = json.Unmarshal([]byte(modelStr), &model); err != nil {
		blog.Errorf("parse model error: %s", err)
		return
	}

	return
}

func (d *Discover) parseData(msg string) (data M, err error) {

	dataStr := gjson.Get(msg, "data.data").String()
	if err = json.Unmarshal([]byte(dataStr), &data); err != nil {
		blog.Errorf("parse data error: %s", err)
		return
	}
	return
}

func (d *Discover) parseHost(msg string) (data M, err error) {

	dataStr := gjson.Get(msg, "data.host").String()
	if err = json.Unmarshal([]byte(dataStr), &data); err != nil {
		blog.Errorf("parse host error: %s", err)
		return
	}
	return
}

func (d *Discover) parseAttrs(msg string) (fields map[string]metadata.ObjAttDes, err error) {

	fieldsStr := gjson.Get(msg, "data.meta.fields").String()

	if err = json.Unmarshal([]byte(fieldsStr), &fields); err != nil {
		blog.Errorf("create model attr unmarshal error: %s", err)
		return
	}
	return
}

func (d *Discover) parseObjID(msg string) string {
	return gjson.Get(msg, "data.meta.model.bk_obj_id").String()
}

func (d *Discover) parseOwnerId(msg string) string {
	ownerId := gjson.Get(msg, "data.host.bk_supplier_account").String()

	if ownerId == "" {
		ownerId = bkc.BKDefaultOwnerID
	}
	return ownerId
}

func (d *Discover) GetAttrs(ownerID, objID, modelAttrKey string, attrs map[string]metadata.ObjAttDes) ([]metadata.Attribute, error) {

	cachedAttrs, err := d.GetModelAttrsFromRedis(modelAttrKey)

	if err == nil && len(cachedAttrs) == len(attrs) {
		blog.Infof("attr exist in redis: %s", modelAttrKey)

		var attrMap = make([]metadata.Attribute, len(cachedAttrs))
		totalEqual := true
		for i, cachedAttr := range cachedAttrs {
			attrMap[i] = metadata.Attribute{PropertyID: cachedAttr}
			if _, ok := attrs[cachedAttr]; !ok {
				totalEqual = false
			}
		}

		if totalEqual {
			blog.Infof("attr exist in redis, and equal: %s", modelAttrKey)
			return attrMap, nil
		}
		blog.Infof("attr exist in redis, but not equal: %s", modelAttrKey)
	}

	cond := mapstr.MapStr{
		bkc.BKObjIDField:   objID,
		bkc.BKOwnerIDField: ownerID,
	}

	resp, err := d.CoreAPI.CoreService().Model().ReadModelAttr(d.ctx, d.httpHeader, objID, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("SelectObjectAttWithParams error %s", err.Error())
		return nil, err
	}
	if !resp.Result {
		blog.Errorf("SelectObjectAttWithParams error %s", resp.ErrMsg)
		return nil, err
	}

	return resp.Data.Info, nil
}

func (d *Discover) UpdateOrAppendAttrs(msg string) error {

	ownerID := d.parseOwnerId(msg)

	objID := d.parseObjID(msg)

	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err)
	}

	attrs, err := d.parseAttrs(msg)
	if err != nil {
		blog.Errorf("create model attr unmarshal error: %s", err)
		return err
	}

	modelAttrKey := d.CreateModelAttrKey(*model, ownerID)

	existAttrs, err := d.GetAttrs(ownerID, objID, modelAttrKey, attrs)
	if nil != err {
		return fmt.Errorf("get attr error: %s", err)
	}

	existAttrHash := make(map[string]bool, len(existAttrs))
	for _, existAttr := range existAttrs {
		existAttrHash[existAttr.PropertyID] = true
	}

	finalAttrs := make([]string, 0)

	hasDiff := false
	for propertyId, property := range attrs {

		finalAttrs = append(finalAttrs, propertyId)

		if existAttrHash[propertyId] {
			continue
		}

		if propertyId == bkc.BKInstNameField {
			blog.Infof("skip default field: %s", propertyId)
			continue
		}

		blog.Infof("attr: %s -> %v", propertyId, property)

		property.ObjectID = objID
		property.OwnerID = ownerID
		property.PropertyID = propertyId
		property.PropertyGroup = bkc.BKDefaultField
		property.Creator = bkc.CCSystemCollectorUserName

		resp, err := d.CoreAPI.TopoServer().Object().CreateObjectAtt(d.ctx, d.httpHeader, &property)
		if err != nil {
			blog.Errorf("create model attr failed %s", err.Error())
			return fmt.Errorf("create model attr failed: %s", err.Error())
		}
		if !resp.Result {
			blog.Errorf("create model attr failed %s", resp.ErrMsg)
			return fmt.Errorf("create model attr failed: %s", resp.ErrMsg)
		}

		hasDiff = true

	}

	if hasDiff {
		attrJs, err := json.Marshal(finalAttrs)
		if err != nil {
			blog.Warnf("%s: flush to redis marshal failed: %s", modelAttrKey, err)
			return nil
		}
		d.TrySetRedis(modelAttrKey, attrJs, cacheTime)
	}

	return nil
}

func (d *Discover) GetModelFromRedis(modelKey string) (MapData, error) {

	var nilR = MapData{}

	val, err := d.redisCli.Get(modelKey).Result()
	if err != nil {
		return nilR, fmt.Errorf("%s: get model cache error: %s", modelKey, err)
	}

	var cacheData = MapData{}
	err = json.Unmarshal([]byte(val), &cacheData)
	if err != nil {
		return nilR, fmt.Errorf("marshal condition error: %s", err)
	}

	return cacheData, nil
}

func (d *Discover) GetModelAttrsFromRedis(modelAttrKey string) ([]string, error) {

	var cacheData = make([]string, 0)

	val, err := d.redisCli.Get(modelAttrKey).Result()
	if err != nil {
		return cacheData, fmt.Errorf("%s: get attr cache error: %s", modelAttrKey, err)
	}

	err = json.Unmarshal([]byte(val), &cacheData)
	if err != nil {
		return cacheData, fmt.Errorf("marshal condition error: %s", err)
	}

	return cacheData, nil

}

func (d *Discover) GetInstFromRedis(instKey string) (map[string]interface{}, error) {

	val, err := d.redisCli.Get(instKey).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: get inst cache error: %s", instKey, err)
	}

	var cacheData = MapData{}
	err = json.Unmarshal([]byte(val), &cacheData)
	if err != nil {
		return nil, fmt.Errorf("marshal condition error: %s", err)
	}

	return cacheData, nil

}

func (d *Discover) CreateModelKey(model Model, ownerID string) string {
	return fmt.Sprintf("cc:v3:model[%s:%s:%s]",
		bkc.CCSystemCollectorUserName,
		ownerID,
		model.BkObjID,
	)
}

func (d *Discover) CreateModelAttrKey(model Model, ownerID string) string {
	return fmt.Sprintf("cc:v3:attr[%s:%s:%s]",
		bkc.CCSystemCollectorUserName,
		ownerID,
		model.BkObjID,
	)
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

func (d *Discover) GetModel(model Model, ownerID string) (bool, error) {
	modelKey := d.CreateModelKey(model, ownerID)

	_, err := d.GetModelFromRedis(modelKey)
	if err == nil {
		blog.Infof("model exist in redis: %s", modelKey)
		return true, nil
	}

	blog.Infof("%s: get model from redis error: %s", modelKey, err)

	cond := mapstr.MapStr{
		bkc.BKObjIDField:   model.BkObjID,
		bkc.BKOwnerIDField: ownerID,
	}
	resp, err := d.CoreAPI.CoreService().Model().ReadModel(d.ctx, d.httpHeader, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search model failed %s", err.Error())
		return false, fmt.Errorf("search model failed: %s", err.Error())
	}
	if !resp.Result {
		blog.Errorf("search model failed %s", resp.ErrMsg)
		return false, fmt.Errorf("search model failed: %s", resp.ErrMsg)
	}
	blog.Infof("search model result: %v", resp.Data)

	if len(resp.Data.Info) > 0 {
		val, err := json.Marshal(resp.Data.Info[0])
		if err != nil {
			blog.Errorf("%s: flush to redis marshal failed: %s", modelKey, err)
		}
		d.TrySetRedis(modelKey, val, cacheTime)
		return true, nil
	}

	return false, nil
}

func (d *Discover) TryCreateModel(msg string) error {
	rid := util.GetHTTPCCRequestID(d.httpHeader)
	ownerID := d.parseOwnerId(msg)

	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err)
	}

	exists, err := d.GetModel(*model, ownerID)
	if nil != err {
		return fmt.Errorf("get inst error: %s", err)
	}

	if exists {
		blog.Infof("model exist, give up create operation")
		return nil
	}

	newObj := metadata.Object{}
	newObj.ObjCls = model.BkClassificationID
	newObj.ObjectID = model.BkObjID
	newObj.ObjectName = model.BkObjName
	newObj.OwnerID = ownerID
	newObj.ObjIcon = defaultModelIcon
	newObj.Creator = bkc.CCSystemCollectorUserName

	input := metadata.CreateModel{
		Spec: newObj,
	}
	resp, err := d.CoreAPI.CoreService().Model().CreateModel(d.ctx, d.httpHeader, &input)
	if err != nil {
		blog.Errorf("create model failed %s", err.Error())
		return fmt.Errorf("create model failed: %s", err.Error())
	}
	if !resp.Result {
		blog.Errorf("create model failed %s", resp.ErrMsg)
		return fmt.Errorf("create model failed: %s", resp.ErrMsg)
	}
	newObj.ID = int64(resp.Data.Created.ID)

	// update registry to iam
	if err := d.authManager.RegisterObject(d.ctx, d.httpHeader, newObj); err != nil {
		blog.Errorf("TryCreateModel success, but RegisterObject failed, object: %+v, err: %s, rid: %s", newObj, err, rid)
		return err
	}

	return nil
}

func (d *Discover) GetInst(ownerID, objID string, keys []string, instKey string) (map[string]interface{}, error) {

	instData, err := d.GetInstFromRedis(instKey)
	if err == nil {
		blog.Infof("inst exist in redis: %s", instKey)
		return instData, nil
	} else {
		blog.Errorf("get inst from redis error: %s", err)
	}

	cond := mapstr.MapStr{
		bkc.BKInstKeyField: instKey,
		bkc.BKObjIDField:   objID,
	}

	resp, err := d.CoreAPI.CoreService().Instance().ReadInstance(d.ctx, d.httpHeader, objID, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search model failed %s", err.Error())
		return nil, fmt.Errorf("search model failed: %s", err.Error())
	}
	if !resp.Result {
		blog.Errorf("search model failed %s", resp.ErrMsg)
		return nil, fmt.Errorf("search model failed: %s", resp.ErrMsg)
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

	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err)
	}

	bodyData, err := d.parseData(msg)
	if err != nil {
		return fmt.Errorf("parse data error: %s", err)
	}

	instKey := bodyData[bkc.BKInstKeyField]
	instKeyStr, ok := instKey.(string)
	if !ok || instKeyStr == "" {
		return fmt.Errorf("skip inst because of empty collect_key: %s", instKeyStr)
	}

	keys := strings.Split(model.Keys, ",")
	inst, err := d.GetInst(ownerID, objID, keys, instKeyStr)
	if nil != err {
		return fmt.Errorf("get inst error: %s", err)
	}

	blog.Infof("get inst result: %v", inst)

	if len(inst) <= 0 {
		data, err := mapstr.NewFromInterface(gjson.Get(msg, "data.data").Value())
		resp, err := d.CoreAPI.CoreService().Instance().CreateInstance(d.ctx, d.httpHeader, objID, &metadata.CreateModelInstance{Data: data})
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
			bizID, err := extensions.ParseBizID(data)
			if err != nil {
				if blog.V(5) {
					blog.InfoJSON("ParseBizID from input data: %+v failed, err: %+v", data)
				}
				return err
			}
			auditHeader, err := GetAuditLogHeader(d.CoreAPI, d.httpHeader, objID)
			if err != nil {
				blog.Errorf("GetAuditLogHeader failed, objID: %s, err: %s, rid: %s", objID, err.Error(), rid)
				return err
			}
			instIDField := common.GetInstIDField(objID)
			data[instIDField] = instID
			auditLog := metadata.SaveAuditLogParams{
				ID:    instID,
				Model: objID,
				Content: metadata.Content{
					CurData: data,
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

	instID, err := util.GetInt64ByInterface(inst[bkc.BKInstIDField])
	if nil != err {
		return fmt.Errorf("get bk_inst_id failed: %s %s", inst[bkc.BKInstIDField], err.Error())
	}

	hasDiff := false
	for attrId, attrValue := range bodyData {

		if attrId == defaultRelateAttr {
			if relateList, ok := inst[defaultRelateAttr].([]interface{}); ok && len(relateList) == 1 {

				relateObj, ok := relateList[0].(map[string]interface{})

				if ok && (relateObj["id"] != "" && relateObj["id"] != "0" && relateObj["id"] != nil) {
					blog.Infof("skip update exist single relation attr: %s->%v", attrId, attrValue)
				} else {

					if attrValue != "" {
						blog.Debug("[relation changed]  %s: %v ---> %v", attrId, attrValue)
						inst[attrId] = attrValue
						hasDiff = true
					}
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

	delete(inst, bkc.BKObjIDField)
	delete(inst, bkc.BKOwnerIDField)
	delete(inst, bkc.BKDefaultField)
	delete(inst, bkc.BKInstIDField)
	delete(inst, bkc.LastTimeField)
	delete(inst, bkc.CreateTimeField)

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
