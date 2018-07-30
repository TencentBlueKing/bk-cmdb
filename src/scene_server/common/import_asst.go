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

package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/language"
	"configcenter/src/common/util"
	sceneUtil "configcenter/src/scene_server/common/util"
	sourceAPI "configcenter/src/source_controller/api/object"
	//"configcenter/src/source_controller/common/commondata"
)

// AsstObjectInst  instances assocate object fields value
type AsstObjectInst struct {
	req            *restful.Request
	ownerID        string
	objAddr        string
	defLang        language.DefaultCCLanguageIf
	fields         map[string]*sourceAPI.ObjAttDes  // import object fields [property]val
	asstPrimaryKey map[string][]sourceAPI.ObjAttDes // import object assocate object primary fields [object id][]object fiedls
	asstInstConds  map[string][]interface{}         // import object assocate fields search instance data condition
	// import object assocate object primary fields and instance data relation map, map[object id]map[primary key val string]instance id
	instPrimayIDMap map[string]map[string]int64
}

// NewAsstObjectInst get Asst object instnace struct,
// NewAsstObjectInst use handle multiple instance  assocate object value
func NewAsstObjectInst(req *restful.Request, ownerID, objAddr string, fields map[string]*sourceAPI.ObjAttDes, defLang language.DefaultCCLanguageIf) *AsstObjectInst {
	return &AsstObjectInst{
		req:     req,
		ownerID: ownerID,
		objAddr: objAddr,
		fields:  fields,
		defLang: defLang,
	}
}

// GetObjAsstObjectPrimaryKey  get instance assocate object primary property fields
func (a *AsstObjectInst) GetObjAsstObjectPrimaryKey() error {
	ret := make(map[string][]sourceAPI.ObjAttDes)

	for _, f := range a.fields {
		if util.IsAssocateProperty(f.PropertyType) {
			fields, err := a.getObjectFields(f.AssociationID, common.BKPropertyIDField)
			if nil != err {
				blog.Errorf("get object  assocate property %s error:%s", f.PropertyID, err.Error())
				return err
			}
			var primaryFields []sourceAPI.ObjAttDes
			for _, f := range fields {
				if true == f.IsOnly {
					primaryFields = append(primaryFields, f)
				}
			}
			ret[f.AssociationID] = primaryFields
		}
	}
	a.asstPrimaryKey = ret
	return nil

}

// SetObjAsstPropertyVal set instance assocate object value to property fields
func (a *AsstObjectInst) SetObjAsstPropertyVal(inst map[string]interface{}) error {

	for key, val := range inst {
		f, ok := a.fields[key]
		if false == ok {
			continue
		}

		if util.IsAssocateProperty(f.PropertyType) {
			strInsts, _ := val.(string)
			if common.ExcelDelAsstObjectRelation == strings.TrimSpace(strInsts) {
				inst[key] = ""
				continue
			}
			// assocate object relation no change
			if "" == strInsts {
				continue
			}

			insts := strings.Split(strInsts, common.ExcelAsstPrimaryKeyRowChar)
			var strIds []string
			for _, inst := range insts {
				inst = strings.TrimSpace(inst)
				if "" == inst {
					continue
				}
				id, err := a.GetIDsByExcelStr(f.AssociationID, inst)
				if nil != err {
					return err
				}
				strIds = append(strIds, fmt.Sprintf("%d", id))
			}

			inst[key] = strings.Join(strIds, common.InstAsstIDSplit)

		}
	}
	return nil
}

// InitInstFromData get assocate object instance data, return map[row]error,  error
func (a *AsstObjectInst) InitInstFromData(infos map[int]map[string]interface{}) (map[int]error, error) {
	rowErr := a.getAsstObjectConds(infos)
	if 0 != len(rowErr) {
		return rowErr, nil
	}

	err := a.getAsstInstByAsstObjectConds()

	if nil != err {
		return nil, err
	}

	return nil, nil
}

// GetIDsByExcelStr Get a string of data based on multiple primary key values (multiple primary key values are sorted by field name, separated by ##)
func (a *AsstObjectInst) GetIDsByExcelStr(objID, key string) (int64, error) {
	if nil == a.instPrimayIDMap {
		return 0, fmt.Errorf("%s primary %s not found", objID, key)
	}
	objIDInst, ok := a.instPrimayIDMap[objID]
	if false == ok {
		return 0, fmt.Errorf("%s primary %s not found", objID, key)
	}
	id, ok := objIDInst[key]
	if false == ok {
		return 0, fmt.Errorf("%s primary %s not found", objID, key)
	}
	return id, nil

}

// getAsstInstByAsstObjectConds  parse import instance assocate property field value to search paramaters
func (a *AsstObjectInst) getAsstInstByAsstObjectConds() error {
	ret := make(map[string]map[string]int64)
	for objID, conds := range a.asstInstConds {
		isExist := false
		for _, f := range a.fields {
			if f.AssociationID == objID {
				isExist = true
			}
		}

		if false == isExist {
			continue
		}
		searchObjID := objID

		if !util.IsInnerObject(objID) {
			searchObjID = common.BKINnerObjIDObject
		}

		url := fmt.Sprintf("%s/object/v1/insts/%s/search", a.objAddr, searchObjID)
		condition := make(map[string]interface{})
		condition["condition"] = common.KvMap{common.BKDBOR: conds}
		isSuccess, message, _, data := a.getInstData(url, common.HTTPSelectPost, condition)
		if false == isSuccess {
			return fmt.Errorf(message)
		}
		if nil == data {
			ret[objID] = make(map[string]int64, 0)
			continue
		}
		primaryKey, _ := a.asstPrimaryKey[objID]
		for _, item := range data {
			mapItem, ok := item.(map[string]interface{})
			if false == ok {
				return fmt.Errorf("not inst data:%v", item)
			}
			keys := []string{}
			for _, f := range primaryKey {
				key, ok := mapItem[f.PropertyID]
				if false == ok {
					errMsg := a.defLang.Languagef("import_str_asst_str_query_data_format_error", objID, f.PropertyID)
					blog.Error(errMsg)
					return errors.New(errMsg)
				}
				keys = append(keys, fmt.Sprintf("%v", key))
			}
			_, ok = ret[objID]
			if false == ok {
				ret[objID] = make(map[string]int64)
			}
			id, err := a.getInstIDFromMapByObjectID(objID, mapItem)
			if nil != err {
				err := fmt.Errorf("get object %s  inst id error, inst info:%v, err:%s ", objID, mapItem, err.Error())
				blog.Errorf("%s ", err.Error())

				return err
			}

			ret[objID][strings.Join(keys, common.ExcelAsstPrimaryKeySplitChar)] = id

		}

	}
	a.instPrimayIDMap = ret
	return nil

}

func (a *AsstObjectInst) getInstIDFromMapByObjectID(objType string, mapInst map[string]interface{}) (int64, error) {
	idField := common.GetInstIDField(objType)
	idInterface, ok := mapInst[idField]
	if false == ok {
		return 0, fmt.Errorf("%s %s not found", objType, idField)
	}
	return util.GetInt64ByInterface(idInterface)
}

// getInstData  get instance data by condition
func (a *AsstObjectInst) getInstData(url, method string, params interface{}) (bool, string, int, []interface{}) {
	var strParams []byte
	switch params.(type) {
	case string:
		strParams = []byte(params.(string))
	default:
		strParams, _ = json.Marshal(params)

	}

	blog.Info("get request url:%s", url)
	blog.Info("get request info  params:%v", string(strParams))
	reply, err := httpcli.ReqHttp(a.req, url, method, []byte(strParams))

	blog.Info("get request result:%v", string(reply))
	if err != nil {
		blog.Error("http do error, params:%s, error:%s", strParams, err.Error())
		return false, err.Error(), 0, nil
	}

	addReply, err := simplejson.NewJson([]byte(reply))
	if err != nil {
		blog.Error("http do error, params:%s, reply:%s, error:%s", strParams, reply, err.Error())
		return false, err.Error(), 0, nil
	}
	isSuccess, err := addReply.Get("result").Bool()
	if nil != err || !isSuccess {
		errMsg, _ := addReply.Get("message").String()
		blog.Error("http do error, url:%s, params:%s, error:%s", url, strParams, errMsg)
		return false, errMsg, 0, nil
	}
	cnt, err := addReply.Get("data").Get("count").Int()
	if err != nil {
		blog.Error("http do error, params:%s, reply:%s, error:%s", strParams, reply, err.Error())
		return false, err.Error(), 0, nil
	}
	data, err := addReply.Get("data").Get("info").Array()
	if err != nil {
		blog.Error("http do error, params:%s, reply:%s, error:%s", strParams, reply, err.Error())
		return false, err.Error(), 0, nil
	}
	return true, "", cnt, data
}

// getAsstObjectConds get all  assocatie field object condition
func (a *AsstObjectInst) getAsstObjectConds(infos map[int]map[string]interface{}) map[int]error {

	errs := make(map[int]error, 0)
	asstMap := make(map[string][]interface{}) //map[AssociationID][]condition

	for rowIndex, info := range infos {
		for key, val := range info {
			f, ok := a.fields[key]
			if false == ok {
				continue
			}
			if util.IsAssocateProperty(f.PropertyType) {

				asstFields, ok := a.asstPrimaryKey[f.AssociationID]
				if false == ok {
					errs[rowIndex] = errors.New(a.defLang.Languagef("import_asst_property_str_not_found", key))
					continue
				}

				strVal, ok := val.(string)
				if false == ok {
					errs[rowIndex] = errors.New(a.defLang.Languagef("import_property_str_format_error", key))
					continue
				}

				if common.ExcelDelAsstObjectRelation == strings.TrimSpace(strVal) {
					continue
				}
				rows := strings.Split(strVal, common.ExcelAsstPrimaryKeyRowChar)

				asstConds := make([]interface{}, 0)
				for _, row := range rows {
					if "" == row {
						continue
					}
					primaryKeys := strings.Split(row, common.ExcelAsstPrimaryKeySplitChar)
					if len(primaryKeys) != len(asstFields) {
						errs[rowIndex] = errors.New(a.defLang.Languagef("import_asst_property_str_primary_count_len", key))
						continue
					}
					conds := common.KvMap{}
					if false == util.IsInnerObject(f.AssociationID) {
						conds[common.BKObjIDField] = f.AssociationID
					}
					for i, val := range primaryKeys {

						asstf := asstFields[i]
						var err error
						conds[asstf.PropertyID], err = sceneUtil.ConvByPropertytype(&asstf, val)
						if nil != err {
							errs[rowIndex] = errors.New(a.defLang.Languagef("import_asst_property_str_primary_count_len", key))
							continue
						}
					}
					asstConds = append(asstConds, conds)

				}

				_, ok = asstMap[f.AssociationID]
				if ok {
					asstMap[f.AssociationID] = append(asstMap[f.AssociationID], asstConds...)
				} else {
					asstMap[f.AssociationID] = asstConds
				}

			}

		}

	}
	a.asstInstConds = asstMap

	return errs
}

//GetObjectFields get object fields
func (a *AsstObjectInst) getObjectFields(objID, sort string) ([]sourceAPI.ObjAttDes, error) {
	data := make(map[string]interface{})
	data[common.BKOwnerIDField] = a.ownerID
	data[common.BKObjIDField] = objID
	data["page"] = common.KvMap{
		"start": 0,
		"limit": common.BKNoLimit,
		"sort":  sort,
	}
	info, _ := json.Marshal(data)
	forward := &sourceAPI.ForwardParam{Header: a.req.Request.Header}
	client := sourceAPI.NewClient(a.objAddr)
	atts, err := client.SearchMetaObjectAtt(forward, []byte(info))
	if nil != err {
		return nil, err
	}

	for idx, a := range atts {
		if !util.IsAssocateProperty(a.PropertyType) {
			continue
		}
		// read property group
		condition := map[string]interface{}{
			"bk_object_att_id":    a.PropertyID, // tmp.PropertyGroup,
			common.BKOwnerIDField: a.OwnerID,
			"bk_obj_id":           a.ObjectID,
		}
		objasstval, jserr := json.Marshal(condition)
		if nil != jserr {
			blog.Error("mashar json failed, error information is %v", jserr)
			return nil, jserr
		}
		asstMsg, err := client.SearchMetaObjectAsst(forward, objasstval)
		if nil != err {
			return nil, err
		}
		if 0 < len(asstMsg) {
			atts[idx].AssociationID = asstMsg[0].AsstObjID // by the rules, only one id
			atts[idx].AsstForward = asstMsg[0].AsstForward // by the rules, only one id
		}
	}
	return atts, nil
}

// SetMapFields set import object property fields
func (a *AsstObjectInst) SetMapFields(objID string) error {
	fields, err := a.getObjectFields(objID, "")
	if nil != err {
		return err
	}
	ret := make(map[string]*sourceAPI.ObjAttDes)
	for index, f := range fields {
		ret[f.PropertyID] = &fields[index]
	}
	a.fields = ret
	return nil
}
