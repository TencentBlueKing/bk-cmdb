package logics

import (
	"encoding/json"
	"time"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"

	"configcenter/src/common/blog"
	bkc "configcenter/src/common"
)

const (
	// 缓存时间5min
	cacheTime = time.Minute * 5
)

const (
	defaultRelateAttr = "host"
	defaultModelIcon  = "icon-cc-middleware"
)

// Model 模型元数据结构
type Model struct {
	BkClassificationID string `json:"bk_classification_id"`
	BkObjID            string `json:"bk_obj_id"`
	BkObjName          string `json:"bk_obj_name"`
	Keys               string `json:"bk_obj_keys"`
}

// Attr 属性元数据结构
type Attr struct {
	BkPropertyName string `json:"bk_property_name"`
	BkPropertyType string `json:"bk_property_type"`
	BkAsstObjID    string `json:"bk_asst_obj_id"`
	Editable       bool   `json:"editable"`
}

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

// Result 接口返回信息结构
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

// parseListResult 接口返回数据解析
func parseListResult(res []byte) (ListResult, error) {

	var lR ListResult

	if err := json.Unmarshal(res, &lR); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s", err)
		return lR, err
	}

	return lR, nil
}

// parseDetailResult 接口返回数据解析
func parseDetailResult(res []byte) (DetailResult, error) {

	var dR DetailResult

	if err := json.Unmarshal(res, &dR); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s", err)
		return dR, err
	}

	return dR, nil
}

// parseResult 接口返回数据解析
func parseResult(res []byte) (Result, error) {

	var r Result

	if err := json.Unmarshal(res, &r); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s", err)
		return r, err
	}

	return r, nil
}

// parseModel 解析模型元数据
func (d *Discover) parseModel(msg string) (model *Model, err error) {

	model = &Model{}
	modelStr := gjson.Get(msg, "data.meta.model").String()

	if err = json.Unmarshal([]byte(modelStr), &model); err != nil {
		blog.Errorf("parse model error: %s", err)
		return
	}

	return
}

// parseData 解析模型数据
func (d *Discover) parseData(msg string) (data M, err error) {

	dataStr := gjson.Get(msg, "data.data").String()
	if err = json.Unmarshal([]byte(dataStr), &data); err != nil {
		blog.Errorf("parse data error: %s", err)
		return
	}
	return
}

// parseHost 解析主机身份数据
func (d *Discover) parseHost(msg string) (data M, err error) {

	dataStr := gjson.Get(msg, "data.host").String()
	if err = json.Unmarshal([]byte(dataStr), &data); err != nil {
		blog.Errorf("parse host error: %s", err)
		return
	}
	return
}

// parseAttrs 解析属性列表
func (d *Discover) parseAttrs(msg string) (fields map[string]Attr, err error) {

	fieldsStr := gjson.Get(msg, "data.meta.fields").String()
	//blog.Debug("create model attr fieldsStr: %s", fieldsStr)
	if err = json.Unmarshal([]byte(fieldsStr), &fields); err != nil {
		blog.Errorf("create model attr unmarshal error: %s", err)
		return
	}
	return
}

// parseObjID 解析模型ID
func (d *Discover) parseObjID(msg string) string {
	return gjson.Get(msg, "data.meta.model.bk_obj_id").String()
}

// parseObjID 解析开发商id
func (d *Discover) parseOwnerId(msg string) string {
	ownerId := gjson.Get(msg, "data.host.bk_supplier_account").String()

	// 使用默认开发商ID
	if ownerId == "" {
		ownerId = bkc.BKDefaultOwnerID
	}
	return ownerId
}

// GetAttrs 查询模型属性
func (d *Discover) GetAttrs(ownerID, objID, modelAttrKey string, attrs map[string]Attr) (ListResult, error) {

	var nilR = ListResult{}

	filterAttrs := make([]string, 0)
	for attr := range attrs {
		filterAttrs = append(filterAttrs, attr)
	}

	cachedAttrs, err := d.GetModelAttrsFromRedis(modelAttrKey)
	// 差异比较
	if err == nil && len(cachedAttrs) == len(filterAttrs) {
		blog.Infof("attr exist in redis: %s", modelAttrKey)

		var attrMap = make([]MapData, len(cachedAttrs))
		var tmpHash = make(map[string]bool, len(cachedAttrs))

		for i, attr := range cachedAttrs {
			tmpHash[attr] = true
			attrMap[i] = MapData{bkc.BKPropertyIDField: attr}
		}

		totalEqual := true
		for _, filterAttr := range filterAttrs {
			if _, ok := tmpHash[filterAttr]; !ok {
				totalEqual = false
				break
			}
		}

		if totalEqual {
			blog.Infof("attr exist in redis, and equal: %s", modelAttrKey)
			return ListResult{
				ResultBase{
					Result:  true,
					Code:    0,
					Message: "success",
				},
				attrMap,
			}, nil
		}

		blog.Infof("attr exist in redis, but not equal: %s", modelAttrKey)

	}

	// construct the condition
	cond := M{
		bkc.BKObjIDField:   objID,
		bkc.BKOwnerIDField: ownerID,
	}

	// marshal the condition
	condJs, err := cond.toJson()
	if err != nil {
		return nilR, fmt.Errorf("marshal condition error: %s", err)
	}

	// search attr by condition
	url := fmt.Sprintf("%s/topo/v1/objectattr/search", d.cc.TopoAPI())
	blog.Infof("get model attr url=%s, body=%s", url, condJs)

	res, err := d.requests.POST(url, nil, condJs)
	if nil != err {
		blog.Errorf("search model err: %s", err)
		return nilR, err
	}

	//blog.Infof("search attr result: %s", res)

	// parse inst data
	dR, err := parseListResult(res)
	if err != nil {
		blog.Errorf("parse result error: %s", err)
		return nilR, err
	}

	return dR, nil

}

// UpdateOrAppendAttrs 创建或新增模型属性
func (d *Discover) UpdateOrAppendAttrs(msg string) error {

	// parse owner id
	ownerID := d.parseOwnerId(msg)

	// parse object_id
	objID := d.parseObjID(msg)

	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err)
	}

	// create model attr
	attrs, err := d.parseAttrs(msg)
	if err != nil {
		blog.Errorf("create model attr unmarshal error: %s", err)
		return err
	}

	modelAttrKey := d.CreateModelAttrKey(*model, ownerID)

	// get exist attr
	dR, err := d.GetAttrs(ownerID, objID, modelAttrKey, attrs)
	if nil != err {
		return fmt.Errorf("get attr error: %s", err)
	}

	existAttrHash := make(map[string]int, len(dR.Data))
	if dR.Result && len(dR.Data) > 0 {
		for i, propertyMap := range dR.Data {
			if idStr, ok := propertyMap[bkc.BKPropertyIDField].(string); ok {
				existAttrHash[idStr] = i
			}
		}
	}

	// collect final attrs of model
	finalAttrs := make([]string, 0)

	// batch create model attrs
	hasDiff := false
	for propertyId, property := range attrs {

		finalAttrs = append(finalAttrs, propertyId)

		// skip exist attr
		if _, ok := existAttrHash[propertyId]; ok {
			continue
		}

		// skip default field
		if propertyId == bkc.BKInstNameField {
			blog.Infof("skip default field: %s", propertyId)
			continue
		}

		blog.Infof("attr: %s -> %v", propertyId, property)

		attrBody := M{
			bkc.BKObjIDField:         objID,
			bkc.BKPropertyGroupField: bkc.BKDefaultField,
			bkc.BKPropertyIDField:    propertyId,
			bkc.BKAsstObjIDField:     property.BkAsstObjID,
			bkc.BKPropertyNameField:  property.BkPropertyName,
			bkc.BKPropertyTypeField:  property.BkPropertyType,
			bkc.BKOwnerIDField:       ownerID,
			bkc.CreatorField:         bkc.CCSystemCollectorUserName,
			"editable":               property.Editable,
		}

		attrBodyJs, err := attrBody.toJson()
		if err != nil {
			return fmt.Errorf("marshal condition error: %s", err)
		}

		url := fmt.Sprintf("%s/topo/v1/objectattr", d.cc.TopoAPI())

		blog.Infof("create model attr url=%s, body=%s", url, attrBody)
		res, err := d.requests.POST(url, nil, []byte(attrBodyJs))
		if nil != err {
			return fmt.Errorf("create model attr error: %s", err.Error())
		}

		blog.Infof("create model attr result: %s", res)

		resMap, err := parseResult(res)
		if !resMap.Result {
			return fmt.Errorf("create model attr failed: %s", resMap.Message)
		}

		// updated
		hasDiff = true

	}

	// flush to redis, ignore fail
	if dR.Result && len(dR.Data) > 0 && hasDiff {
		attrJs, err := json.Marshal(finalAttrs)
		if err != nil {
			blog.Warnf("%s: flush to redis marshal failed: %s", modelAttrKey, err)
			return nil
		}
		d.TrySetRedis(modelAttrKey, attrJs, cacheTime)
	}

	return nil
}

// GetModelFromRedis 从redis中获取模型元数据
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

// GetModelFromRedis 从redis中获取模型元数据
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

// GetInstFromRedis 从redis中获取实例数据
func (d *Discover) GetInstFromRedis(instKey string) (DetailResult, error) {

	var nilR = DetailResult{}

	val, err := d.redisCli.Get(instKey).Result()
	if err != nil {
		return nilR, fmt.Errorf("%s: get inst cache error: %s", instKey, err)
	}

	var cacheData = DetailResult{}
	err = json.Unmarshal([]byte(val), &cacheData)
	if err != nil {
		return nilR, fmt.Errorf("marshal condition error: %s", err)
	}

	return cacheData, nil

}

// CreateModelKey 根据model生成key
func (d *Discover) CreateModelKey(model Model, ownerID string) string {
	return fmt.Sprintf("cc:v3:model[%s:%s:%s]",
		bkc.CCSystemCollectorUserName,
		ownerID,
		model.BkObjID,
	)
}

// CreateModelAttrKey 根据model生成mode-attr的key
func (d *Discover) CreateModelAttrKey(model Model, ownerID string) string {
	return fmt.Sprintf("cc:v3:attr[%s:%s:%s]",
		bkc.CCSystemCollectorUserName,
		ownerID,
		model.BkObjID,
	)
}

// TrySetRedis 尝试写入redis，忽略失败情况
func (d *Discover) TrySetRedis(key string, value []byte, duration time.Duration) {
	_, err := d.redisCli.Set(key, value, duration).Result()
	if err != nil {
		blog.Warnf("%s: flush to redis failed: %s", key, err)
	} else {

		blog.Infof("%s: flush to redis success", key)
	}
}

// TryUnSetRedis 尝试移除redis缓存，忽略失败情况
func (d *Discover) TryUnsetRedis(key string) {
	_, err := d.redisCli.Del(key).Result()
	if err != nil {
		blog.Warnf("%s: remove from redis failed: %s", key, err)
	} else {

		blog.Infof("%s: remove from redis success", key)
	}
}

// GetModel 查询模型元数据
func (d *Discover) GetModel(model Model, modelKey, ownerID string) (ListResult, error) {

	var nilR = ListResult{}

	// construct the condition
	cond := M{
		bkc.BKObjIDField:   model.BkObjID,
		bkc.BKOwnerIDField: ownerID,
		//bkc.CreatorField:            bkc.CCSystemCollectorUserName,
	}

	// try fetch redis cache
	modelData, err := d.GetModelFromRedis(modelKey)
	if err == nil {
		blog.Infof("model exist in redis: %s", modelKey)
		return ListResult{
			ResultBase{
				Result:  true,
				Code:    0,
				Message: "success",
			},
			[]MapData{modelData},
		}, nil
	}

	blog.Infof("%s: get model from redis error: %s", modelKey, err)

	// marshal the condition
	condJs, err := cond.toJson()
	if err != nil {
		return nilR, fmt.Errorf("marshal condition error: %s", err)
	}

	// search object by condition
	url := fmt.Sprintf("%s/topo/v1/objects", d.cc.TopoAPI())
	blog.Infof("get model url=%s, condition=%s", url, condJs)

	res, err := d.requests.POST(url, nil, condJs)
	if nil != err {
		blog.Errorf("search model err: %s", err)
		return nilR, err
	}

	blog.Infof("search model result: %s", res)

	// parse inst data
	dR, err := parseListResult(res)
	if err != nil {
		blog.Errorf("parse result error: %s", err)
		return nilR, err
	}

	// try flush to redis, maybe fail
	if dR.Result && len(dR.Data) > 0 {

		val, err := M(dR.Data[0]).toJson()
		if err != nil {
			blog.Errorf("%s: flush to redis marshal failed: %s", modelKey, err)
		}
		d.TrySetRedis(modelKey, val, cacheTime)
	}

	return dR, nil

}

// TryCreateModel 创建模型元数据
func (d *Discover) TryCreateModel(msg string) error {
	// parse ownerID
	ownerID := d.parseOwnerId(msg)

	// parse model
	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err)
	}

	modelKey := d.CreateModelKey(*model, ownerID)
	dR, err := d.GetModel(*model, modelKey, ownerID)
	if nil != err {
		return fmt.Errorf("get inst error: %s", err)
	}

	// model exist
	if dR.Result && len(dR.Data) > 0 {
		blog.Infof("model exist, give up create operation")
		return nil
	}

	//create model
	body := M{
		bkc.BKClassificationIDField: model.BkClassificationID,
		bkc.BKObjIDField:            model.BkObjID,
		bkc.BKObjNameField:          model.BkObjName,
		bkc.BKOwnerIDField:          ownerID,
		bkc.BKObjIconField:          defaultModelIcon,
		bkc.CreatorField:            bkc.CCSystemCollectorUserName,
	}

	bodyJs, _ := body.toJson()
	url := fmt.Sprintf("%s/topo/v1/object", d.cc.TopoAPI())
	blog.Infof("create model url=%s, body=%s", bodyJs)

	res, err := d.requests.POST(url, nil, bodyJs)
	if nil != err {
		return fmt.Errorf("create model error: %s", err.Error())
	}

	blog.Debug("create model result: %s", res)

	resMap, err := parseResult(res)
	if !resMap.Result {
		return fmt.Errorf("create model failed: %s", resMap.Message)
	}

	return nil
}

// GetInst 获取模型实例信息
func (d *Discover) GetInst(ownerID, objID string, keys []string, instKey string) (DetailResult, error) {

	var nilR = DetailResult{}

	// try fetch inst cache from redis
	instData, err := d.GetInstFromRedis(instKey)
	if err == nil {
		blog.Infof("inst exist in redis: %s", instKey)
		return instData, nil
	} else {
		blog.Errorf("get inst from redis error: %s", err)
	}

	// construct the condition
	condition := M{
		"fields": []string{},
		"page": M{
			"start": 0,
			"limit": 1,
			"sort":  bkc.BKInstNameField,
		},
		"condition": M{
			bkc.BKCollectorKeyField: instKey,
			bkc.BKObjIDField:        objID,
		},
	}

	// marshal the condition
	condJs, err := condition.toJson()
	if err != nil {
		return nilR, fmt.Errorf("marshal condition error: %s", err)
	}

	// search inst by condition
	url := fmt.Sprintf("%s/topo/v1/inst/search/%s/%s", d.cc.TopoAPI(), ownerID, objID)
	blog.Infof("get inst url=%s, condition=%s", url, condJs)

	res, err := d.requests.POST(url, nil, condJs)
	if nil != err {
		blog.Errorf("search inst err: %s", err)
		return nilR, err
	}

	blog.Debug("search inst result: %s", res)

	// parse inst data
	dR, err := parseDetailResult(res)
	if err != nil {
		blog.Errorf("parse result error: %s", err)
		return nilR, err
	}

	// try flush to redis, maybe fail
	if dR.Result && dR.Data.Count > 0 {

		val, err := json.Marshal(dR)
		if err != nil {
			blog.Errorf("%s: flush to redis marshal failed: %s", instKey, err)
		}
		d.TrySetRedis(instKey, val, cacheTime)
	}

	return dR, nil

}

// UpdateOrCreateInst 更新或新增模型实例信息
func (d *Discover) UpdateOrCreateInst(msg string) error {

	// parse ownerID
	ownerID := d.parseOwnerId(msg)

	// parse object_id
	objID := d.parseObjID(msg)

	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err)
	}

	bodyData, err := d.parseData(msg)
	if err != nil {
		return fmt.Errorf("parse data error: %s", err)
	}

	// try fetch inst cache from redis
	instKey := bodyData[bkc.BKCollectorKeyField]
	instKeyStr, ok := instKey.(string)
	if !ok || instKeyStr == "" {
		return fmt.Errorf("skip inst because of empty collect_key: %s", instKeyStr)
	}

	// fetch key's values
	keys := strings.Split(model.Keys, ",")
	dR, err := d.GetInst(ownerID, objID, keys, instKeyStr)
	if nil != err {
		return fmt.Errorf("get inst error: %s", err)
	}

	blog.Infof("get inst result: count=%d", dR.Data.Count)

	// create inst
	if dR.Data.Count == 0 {

		createJs := gjson.Get(msg, "data.data").String()

		url := fmt.Sprintf("%s/topo/v1/inst/%s/%s", d.cc.TopoAPI(), ownerID, objID)
		blog.Infof("create inst url=%s, body=%s", url, createJs)

		res, err := d.requests.POST(url, nil, []byte(createJs))
		if nil != err {
			return fmt.Errorf("create inst error: %s", err)
		}

		blog.Infof("create inst result: %s", res)

		resMap, err := parseResult(res)
		if !resMap.Result {
			return fmt.Errorf("create inst failed: %s", resMap.Message)
		}

		return nil
	}

	// update exist inst
	info := dR.Data.Info[0]
	instID, ok := info[bkc.BKInstIDField].(float64)
	if !ok {
		return fmt.Errorf("get bk_inst_id failed: %s", info[bkc.BKInstIDField])
	}

	// update info by sample data
	hasDiff := false
	for attrId, attrValue := range bodyData {
		// skip single relation attr: host
		if info[attrId] != attrValue && attrId != defaultRelateAttr {
			blog.Debug("[changed]%s: %v ---> %v", attrId, attrValue, info[attrId])
			hasDiff = true
		}
		info[attrId] = attrValue
	}

	if !hasDiff {
		blog.Infof("no need to update inst")
		return nil
	}

	// remove some keys from query result
	delete(info, bkc.BKObjIDField)
	delete(info, bkc.BKOwnerIDField)
	delete(info, bkc.BKDefaultField)
	delete(info, bkc.BKInstIDField)
	delete(info, bkc.LastTimeField)
	delete(info, bkc.CreateTimeField)

	updateJs, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal inst data error: %s", err)
	}

	url := fmt.Sprintf("%s/topo/v1/inst/%s/%s/%d", d.cc.TopoAPI(), ownerID, objID, int(instID))
	blog.Infof("update inst url=%s, body=%s", url, updateJs)

	res, err := d.requests.PUT(url, nil, updateJs)
	if nil != err {
		return fmt.Errorf("update inst error: %s", err)
	}

	blog.Infof("update inst result: %s", res)

	// remove bad cache
	d.TryUnsetRedis(instKeyStr)

	return nil
}
