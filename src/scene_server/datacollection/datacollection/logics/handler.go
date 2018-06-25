package logics

import (
	"encoding/json"
	"io/ioutil"
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
	// 配置文件路径
	jsonFile = "/etc/gse/host/hostid"
)

var Host = M{
	"bk_host_id": "",
}

func init() {
	// parse json file
	jf, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		blog.Errorf("discover: read file error: %s")
	} else {
		err = json.Unmarshal(jf, &Host)
		if err != nil {
			blog.Errorf("discover: parse file error: %s")
		}
	}

	blog.Infof("Host Detail: %s\n%v\n", jf, Host)
}

// 模型元数据结构
type Model struct {
	BkClassificationID string `json:"bk_classification_id"`
	BkObjID            string `json:"bk_obj_id"`
	BkObjName          string `json:"bk_obj_name"`
	Keys               string `json:"keys"`
}

// 属性元数据结构
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
	for k := range m {
		keys = append(keys, k)
	}

	return
}

// 接口返回信息结构
type Result struct {
	ResultBase
	Data interface{} `json:"data"`
}

func (r *Result) mapData() (MapData, error) {
	if m, ok := r.Data.(MapData); ok {
		return m, nil
	}
	return nil, fmt.Errorf("parse map data error: %v", r)
}

// 接口返回数据解析
func parseListResult(res []byte) (ListResult, error) {

	var lR ListResult

	if err := json.Unmarshal(res, &lR); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s\n", err)
		return lR, err
	}

	return lR, nil
}

func parseDetailResult(res []byte) (DetailResult, error) {

	var dR DetailResult

	if err := json.Unmarshal(res, &dR); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s\n", err)
		return dR, err
	}

	return dR, nil
}

func parseResult(res []byte) (Result, error) {

	var r Result

	if err := json.Unmarshal(res, &r); nil != err {
		blog.Errorf("failed to unmarshal the result, error info is: %s\n", err)
		return r, err
	}

	return r, nil
}

// parseModel 解析模型元数据
func (d *Discover) parseModel(msg string) (model *Model, err error) {

	model = &Model{}
	modelStr := gjson.Get(msg, "data.meta.model").String()

	if err = json.Unmarshal([]byte(modelStr), &model); err != nil {
		blog.Errorf("unmarshal error: %s", err)
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

// parseAttrs 解析属性列表
func (d *Discover) parseAttrs(msg string) (fields map[string]Attr, err error) {

	fieldsStr := gjson.Get(msg, "data.meta.fields").String()
	//blog.Debug("create model attr fieldsStr: %s\n", fieldsStr)
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

// GetAttrs 查询模型属性
func (d *Discover) GetAttrs(objID string, modelAttrKey string, attrs map[string]Attr) (ListResult, error) {

	var nilR = ListResult{}

	filterAttrs := make([]string, 0)
	for k := range attrs {
		// skip empty string
		if k == "" {
			continue
		}
		filterAttrs = append(filterAttrs, k)
	}

	cachedAttrs, err := d.GetModelAttrsFromRedis(modelAttrKey)

	// 差异比较
	if err == nil && len(cachedAttrs) == len(filterAttrs) {
		blog.Infof("attr exist in redis: %s\n", modelAttrKey)

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
			blog.Infof("attr exist in redis, and equal: %s\n", modelAttrKey)
			return ListResult{
				ResultBase{
					Result:  true,
					Code:    0,
					Message: "success",
				},
				attrMap,
			}, nil
		}

		blog.Infof("attr exist in redis, but not equal: %s\n", modelAttrKey)

	}

	// construct the condition
	cond := M{
		bkc.BKObjIDField:   objID,
		bkc.BKOwnerIDField: bkc.BKDefaultOwnerID,
	}

	// marshal the condition
	condJs, err := cond.toJson()
	if err != nil {
		return nilR, fmt.Errorf("marshal condition error: %s", err)
	}

	// search attr by condition
	url := fmt.Sprintf("%s/topo/v1/objectattr/search", d.cc.TopoAPI())
	blog.Infof("get model attr url=%s, body=%s\n", url, condJs)

	res, err := d.requests.POST(url, nil, condJs)
	if nil != err {
		blog.Errorf("search model err: %s\n", err)
		return nilR, err
	}

	//blog.Infof("search attr result: %s\n", res)

	// parse inst data
	dR, err := parseListResult(res)
	if err != nil {
		blog.Errorf("parse result error: %s\n", err)
		return nilR, err
	}

	return dR, nil

}

// UpdateOrAppendAttrs 创建或新增模型属性
func (d *Discover) UpdateOrAppendAttrs(msg string) error {

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

	modelAttrKey := d.CreateModelAttrKey(*model)

	// get exist attr
	dR, err := d.GetAttrs(objID, modelAttrKey, attrs)
	if nil != err {
		return fmt.Errorf("get attr error: %s", err)
	}

	existAttrHash := make(map[string]int, len(dR.Data))
	//existAttrs := make([]string, len(dR.Data))
	if dR.Result && len(dR.Data) > 0 {
		for i, v := range dR.Data {
			if idStr, ok := v[bkc.BKPropertyIDField].(string); ok {
				existAttrHash[idStr] = i
				//existAttrs = append(existAttrs, idStr)
			}
		}
		//blog.Infof("attr exist: %v\n", existAttrs)
	}

	// add extra relation attr associate to host
	attrs["host"] = Attr{
		BkPropertyName: "主机",
		BkPropertyType: bkc.FieldTypeSingleAsst,
		BkAsstObjID:    "host",
	}

	// add collector key
	attrs[bkc.BKCollectorKeyField] = Attr{
		BkPropertyName: "采集标识",
		BkPropertyType: bkc.FieldTypeSingleChar,
		Editable:       false,
	}


	// collect final attrs of model
	finalAttrs := make([]string, 0)

	// batch create model attrs
	for instId, v := range attrs {

		finalAttrs = append(finalAttrs, instId)

		// skip exist attr
		if _, ok := existAttrHash[instId]; ok {
			//blog.Infof("skip exist field: %s", instId)
			continue
		}

		blog.Infof("attr: %s -> %v\n", instId, v)

		// skip default field
		if instId == bkc.BKInstNameField {
			blog.Infof("skip default field: %s", instId)
			continue
		}

		attrBody := M{
			bkc.BKObjIDField:         objID,
			bkc.BKPropertyGroupField: bkc.BKDefaultField,
			bkc.BKPropertyIDField:    instId,
			bkc.BKAsstObjIDField:     v.BkAsstObjID,
			bkc.BKPropertyNameField:  v.BkPropertyName,
			bkc.BKPropertyTypeField:  v.BkPropertyType,
			bkc.BKOwnerIDField:       bkc.BKDefaultOwnerID,
			bkc.CreatorField:         bkc.CCSystemCollectorUserName,
			"editable":               v.Editable,
		}

		attrBodyJs, _ := attrBody.toJson()
		url := fmt.Sprintf("%s/topo/v1/objectattr", d.cc.TopoAPI())

		blog.Infof("create model attr url=%s, body=%s\n", url, attrBody)
		res, err := d.requests.POST(url, nil, []byte(attrBodyJs))
		if nil != err {
			return fmt.Errorf("create model attr error: %s", err.Error())
		}

		blog.Infof("create model attr result: %s\n", res)

		resMap, err := parseResult(res)
		if !resMap.Result {
			return fmt.Errorf("create model attr failed: %s\n", resMap.Message)
		}

	}

	// flush to redis
	if dR.Result && len(dR.Data) > 0 {
		attrJs, err := json.Marshal(finalAttrs)
		if err != nil {
			blog.Infof("%s: flush to redis marshal failed: %s\n", modelAttrKey, err)
			return err
		}

		_, err = d.redisCli.Set(modelAttrKey, attrJs, cacheTime).Result()
		if err != nil {
			blog.Infof("%s: flush to redis failed: %s\n", modelAttrKey, err)
			return err
		}

		blog.Errorf("%s: flush to redis success: %s\n", modelAttrKey, attrJs)

	}

	return nil
}

// GetModelFromRedis 从redis中获取模型元数据
func (d *Discover) GetModelFromRedis(modelKey string) (MapData, error) {

	var nilR = MapData{}

	val, err := d.redisCli.Get(modelKey).Result()
	if err != nil {
		return nilR, fmt.Errorf("%s: get model cache error: %s\n", modelKey, err)
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
		return cacheData, fmt.Errorf("%s: get attr cache error: %s\n", modelAttrKey, err)
	}

	err = json.Unmarshal([]byte(val), &cacheData)
	if err != nil {
		return cacheData, fmt.Errorf("marshal condition error: %s", err)
	}

	return cacheData, nil

}

// CreateModelKey 根据model生成key
func (d *Discover) CreateModelKey(model Model) string {
	return fmt.Sprintf("cc:v3:model[%s:%s:%s:%s]",
		bkc.CCSystemCollectorUserName,
		bkc.BKDefaultOwnerID,
		model.BkClassificationID,
		model.BkObjID,
	)
}

// 根据model生成mode-attr的key
func (d *Discover) CreateModelAttrKey(model Model) string {
	return fmt.Sprintf("cc:v3:attr[%s:%s:%s:%s]",
		bkc.CCSystemCollectorUserName,
		bkc.BKDefaultOwnerID,
		model.BkClassificationID,
		model.BkObjID,
	)
}

// GetModel 查询模型元数据
func (d *Discover) GetModel(msg string) (ListResult, error) {

	var nilR = ListResult{}

	model, err := d.parseModel(msg)
	if err != nil {
		return nilR, fmt.Errorf("parse model error: %s", err)
	}

	// construct the condition
	cond := M{
		bkc.BKObjIDField:            model.BkObjID,
		bkc.BKClassificationIDField: model.BkClassificationID,
		bkc.BKOwnerIDField:          bkc.BKDefaultOwnerID,
		bkc.CreatorField:            bkc.CCSystemCollectorUserName,
	}

	// try fetch redis cache
	modelKey := d.CreateModelKey(*model)
	modelData, err := d.GetModelFromRedis(modelKey)
	if err == nil {
		blog.Infof("model exist in redis: %s\n", modelKey)
		return ListResult{
			ResultBase{
				Result:  true,
				Code:    0,
				Message: "success",
			},
			[]MapData{modelData},
		}, nil
	}

	blog.Infof("%s: get model from redis error: %s\n", modelKey, err)

	// marshal the condition
	condJs, err := cond.toJson()
	if err != nil {
		return nilR, fmt.Errorf("marshal condition error: %s", err)
	}

	// search object by condition
	url := fmt.Sprintf("%s/topo/v1/objects", d.cc.TopoAPI())
	blog.Infof("get model url=%s, condition=%s\n", url, condJs)

	res, err := d.requests.POST(url, nil, condJs)
	if nil != err {
		blog.Errorf("search model err: %s\n", err)
		return nilR, err
	}

	blog.Infof("search model result: %s\n", res)

	// parse inst data
	dR, err := parseListResult(res)
	if err != nil {
		blog.Errorf("parse result error: %s\n", err)
		return nilR, err
	}

	// flush to redis
	if dR.Result && len(dR.Data) > 0 {

		val, err := M(dR.Data[0]).toJson()
		if err != nil {
			blog.Infof("%s: flush to redis marshal failed: %s\n", modelKey, err)
			return dR, nil
		}

		_, err = d.redisCli.Set(modelKey, val, cacheTime).Result()
		if err != nil {
			blog.Infof("%s: flush to redis failed: %s\n", modelKey, err)
			return dR, nil
		}

		blog.Errorf("%s: flush to redis success\n", modelKey)

	}

	return dR, nil

}

// TryCreateModel 创建模型元数据
func (d *Discover) TryCreateModel(msg string) error {

	dR, err := d.GetModel(msg)
	if nil != err {
		return fmt.Errorf("get inst error: %s", err)
	}

	// model exist
	if dR.Result && len(dR.Data) > 0 {
		blog.Infof("model exist, give up create operation\n")
		return nil
	}

	//create model
	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err.Error())
	}

	body := M{
		bkc.BKClassificationIDField: model.BkClassificationID,
		bkc.BKObjIDField:            model.BkObjID,
		bkc.BKObjNameField:          model.BkObjName,
		bkc.BKOwnerIDField:          bkc.BKDefaultOwnerID,
		bkc.BKObjIconField:          "icon-cc-middleware",
		bkc.CreatorField:            bkc.CCSystemCollectorUserName,
	}

	bodyJs, _ := body.toJson()
	url := fmt.Sprintf("%s/topo/v1/object", d.cc.TopoAPI())
	blog.Infof("create model url=%s, body=%s\n", bodyJs)

	res, err := d.requests.POST(url, nil, bodyJs)
	if nil != err {
		return fmt.Errorf("create model error: %s", err.Error())
	}

	blog.Debug("create model result: %s\n", res)

	resMap, err := parseResult(res)
	if !resMap.Result {
		return fmt.Errorf("create model failed: %s\n", resMap.Message)
	}

	return nil
}

// GetInst 获取模型实例信息
func (d *Discover) GetInst(objID string, keys []string, bodyMap M) (DetailResult, error) {

	var nilR = DetailResult{}

	// build condition
	condition := M{
		//bkc.CreatorField: bkc.CCSystemCollectorUserName,
		bkc.BKObjIDField: objID,
	}

	for _, key := range keys {
		keyStr := string(key)
		condition[keyStr] = bodyMap[keyStr]
	}

	// construct the condition
	cond := M{
		"fields": []string{},
		"page": M{
			"start": 0,
			"limit": 1,
			"sort":  bkc.BKInstNameField,
		},
		"condition": condition,
	}

	// marshal the condition
	condJs, err := cond.toJson()
	if err != nil {
		return nilR, fmt.Errorf("marshal condition error: %s", err)
	}

	// search inst by condition
	url := fmt.Sprintf("%s/topo/v1/inst/search/%s/%s", d.cc.TopoAPI(), bkc.BKDefaultOwnerID, objID)
	blog.Infof("get inst url=%s, condition=%s\n", url, condJs)

	res, err := d.requests.POST(url, nil, condJs)
	if nil != err {
		blog.Errorf("search inst err: %s\n", err)
		return nilR, err
	}

	//blog.Debug("search inst result: %s\n", res)

	// parse inst data
	dR, err := parseDetailResult(res)
	if err != nil {
		blog.Errorf("parse result error: %s\n", err)
		return nilR, err
	}

	return dR, nil

}

// UpdateOrCreateInst 更新或新增模型实例信息
func (d *Discover) UpdateOrCreateInst(msg string) error {

	// parse object_id
	objID := d.parseObjID(msg)

	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err)
	}

	bodyMap, err := d.parseData(msg)
	if err != nil {
		return fmt.Errorf("parse data error: %s", err)
	}

	// fetch key's values
	keys := strings.Split(model.Keys, ",")
	dR, err := d.GetInst(objID, keys, bodyMap)
	if nil != err {
		return fmt.Errorf("get inst error: %s", err)
	}

	blog.Infof("get inst result: count=%d\n", dR.Data.Count)

	bodyData, err := d.parseData(msg)
	if nil != err {
		return fmt.Errorf("parse inst data error: %s", err)
	}

	// create inst
	if dR.Data.Count == 0 {

		values := make([]string, len(keys))
		for i, key := range keys {
			keyStr := string(key)
			valueStr, ok := bodyMap[keyStr].(string)
			if !ok {
				return fmt.Errorf("parse key error: %v(%s)\n", bodyMap, keyStr)
			}
			values[i] = valueStr
		}

		// add unique key to collector_key and reassign it to inst_name
		//createJs := gjson.Get(msg, "data.data").String()
		bodyData[bkc.BKCollectorKeyField] = strings.Join(values, "-")
		bodyData[bkc.BKInstNameField] = bodyData[bkc.BKCollectorKeyField]
		// add related host
		bodyData["host"] = Host[bkc.BKHostIDField]

		// marshal the condition
		createJs, err := bodyData.toJson()
		if err != nil {
			return fmt.Errorf("marshal condition error: %s", err)
		}

		url := fmt.Sprintf("%s/topo/v1/inst/%s/%s", d.cc.TopoAPI(), bkc.BKDefaultOwnerID, objID)
		blog.Infof("create inst url=%s, body=%s\n", url, createJs)

		res, err := d.requests.POST(url, nil, createJs)
		if nil != err {
			return fmt.Errorf("create inst error: %s", err)
		}

		blog.Infof("create inst result: %s\n", res)

		resMap, err := parseResult(res)
		if !resMap.Result {
			return fmt.Errorf("create inst failed: %s\n", resMap.Message)
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
	for k, v := range bodyData {
		if info[k] != v {
			blog.Debug("[changed]%s: %v ---> %v", k, v, info[k])
			hasDiff = true
		}
		info[k] = v
	}

	if !hasDiff {
		blog.Infof("no need to update inst")
		return nil
	}

	// remove some keys
	delete(info, bkc.BKObjIDField)
	delete(info, bkc.BKOwnerIDField)
	delete(info, bkc.BKDefaultField)
	delete(info, bkc.BKInstIDField)
	delete(info, bkc.LastTimeField)
	delete(info, bkc.CreateTimeField)

	//info[bkc.CreatorField] = bkc.CCSystemCollectorUserName

	updateJs, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal inst data error: %s", err)
	}

	url := fmt.Sprintf("%s/topo/v1/inst/%s/%s/%d", d.cc.TopoAPI(), bkc.BKDefaultOwnerID, objID, int(instID))
	blog.Infof("update inst url=%s, body=%s\n", url, updateJs)

	res, err := d.requests.PUT(url, nil, updateJs)
	if nil != err {
		return fmt.Errorf("update inst error: %s", err)
	}

	blog.Infof("update inst result: %s\n", res)

	return nil
}
