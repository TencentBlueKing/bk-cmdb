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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
	"encoding/json"
	"errors"
	"fmt"

	httpcli "configcenter/src/common/http/httpclient"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

func DefaultAppMigrate(req *restful.Request, cc *api.APIResource, ownerID string) error {
	isExist, err := defaultAppIsExist(req, cc, ownerID)
	if nil != err {
		return err
	}
	if !isExist {
		return addDefaultApp(req, cc, ownerID)
	}
	return nil

}

func addDefaultApp(req *restful.Request, cc *api.APIResource, ownerID string) error {
	params, err := getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDApp)
	if err != nil {
		blog.Errorf("get app fields %s", err.Error())
		return err
	}
	params[common.BKAppNameField] = common.DefaultAppName
	params[common.BKMaintainersField] = "admin"
	params[common.BKProductPMField] = "admin"
	params[common.BKTimeZoneField] = "Asia/Shanghai"
	params[common.BKLanguageField] = "1" //中文
	params[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal

	byteParams, _ := json.Marshal(params)
	url := cc.TopoAPI() + "/topo/v1/app/default/" + ownerID
	blog.Info("migrate add default app url :%s", url)
	blog.Info("migrate add default app content :%s", string(byteParams))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
	blog.Info("migrate add default app return :%s", string(reply))
	if err != nil {
		return err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()

	code, err := util.GetIntByInterface(output[common.HTTPBKAPIErrorCode])
	if err != nil {
		return errors.New(reply)
	}
	if 0 != code {
		return errors.New(fmt.Sprint(output[common.HTTPBKAPIErrorMessage]))
	}

	return nil
}

func defaultAppIsExist(req *restful.Request, cc *api.APIResource, ownerID string) (bool, error) {

	params := make(map[string]interface{})

	params["condition"] = make(map[string]interface{})
	params["fields"] = []string{common.BKAppIDField}
	params["start"] = 0
	params["limit"] = 20

	byteParams, _ := json.Marshal(params)
	url := cc.TopoAPI() + "/topo/v1/app/default/" + ownerID + "/search"
	blog.Info("migrate get default app url :%s", url)
	blog.Info("migrate get default app content :%s", string(byteParams))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, byteParams)
	blog.Info("migrate get default app return :%s", string(reply))
	if err != nil {
		return false, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()

	code, err := util.GetIntByInterface(output["bk_error_code"])
	if err != nil {
		return false, errors.New(reply)
	}
	if 0 != code {
		return false, errors.New(output["message"].(string))
	}
	cnt, err := js.Get("data").Get("count").Int()
	if err != nil {
		return false, errors.New(reply)
	}
	if 0 == cnt {
		return false, nil
	}
	return true, nil
}

func getObjectFields(url string, req *restful.Request, objID string) (common.KvMap, error) {
	url = url + "/topo/v1/objectattr/search"
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"skip": 0, "limit": common.BKNoLimit}}
	byteParams, _ := json.Marshal(conds)
	blog.Info("migrate get object fields url :%s", url)
	blog.Info("migrate get object fields content :%s", string(byteParams))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
	blog.Info("migrate get object fileds return :%s", string(reply))
	if err != nil {
		return nil, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	hostFields, _ := js.Map()
	fields, _ := hostFields["data"].([]interface{})
	ret := common.KvMap{}
	type intOptionType struct {
		Min int
		Max int
	}
	type EnumOptionType struct {
		Name string
		Type string
	}

	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})
		fieldName, _ := mapField["bk_property_id"].(string)
		fieldType, _ := mapField["bk_property_type"].(string)
		option, _ := mapField["option"]
		switch fieldType {
		case common.FieldTypeSingleChar:
			ret[fieldName] = ""
		case common.FieldTypeLongChar:
			ret[fieldName] = ""
		case common.FieldTypeInt:
			ret[fieldName] = nil
		case common.FieldTypeEnum:
			enumOptions := validator.ParseEnumOption(option)
			v := ""
			if len(enumOptions) > 0 {
				var defaultOption *validator.EnumVal
				for _, k := range enumOptions {
					if k.IsDefault {
						defaultOption = &k
						break
					}
				}
				if nil != defaultOption {
					v = defaultOption.ID
				}
			}
			ret[fieldName] = v
		case common.FieldTypeDate:
			ret[fieldName] = ""
		case common.FieldTypeTime:
			ret[fieldName] = ""
		case common.FieldTypeUser:
			ret[fieldName] = ""
		case common.FieldTypeMultiAsst:
			ret[fieldName] = nil
		case common.FieldTypeTimeZone:
			ret[fieldName] = nil
		case common.FieldTypeBool:
			ret[fieldName] = false
		default:
			ret[fieldName] = nil
			continue
		}

	}
	return ret, nil
}
