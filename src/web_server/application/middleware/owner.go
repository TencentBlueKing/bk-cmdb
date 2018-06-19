package middleware

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/scene_server/validator"
	webCommon "configcenter/src/web_server/common"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

type OwnerManager struct {
	httpCli  *httpclient.HttpClient
	OwnerID  string
	UserName string
}

func NewOwnerManager(userName, ownerID, language string) *OwnerManager {
	ownerManager := new(OwnerManager)
	ownerManager.UserName = userName
	ownerManager.OwnerID = ownerID
	ownerManager.httpCli = httpclient.NewHttpClient()
	ownerManager.httpCli.SetHeader(common.BKHTTPHeaderUser, userName)
	ownerManager.httpCli.SetHeader(common.BKHTTPLanguage, language)
	ownerManager.httpCli.SetHeader(common.BKHTTPOwnerID, common.BKSuperOwnerID)
	return ownerManager
}

func (m *OwnerManager) InitOwner() error {
	blog.Infof("init owner %s", m.OwnerID)
	exist, err := m.defaultAppIsExist()
	if err != nil {
		return err
	}
	if !exist {
		err = m.addDefaultApp()
		if nil != err {
			return err
		}
	}
	return nil
}

func (m *OwnerManager) addDefaultApp() error {
	blog.Info("addDefaultApp")
	params, err := m.getObjectFields(common.BKInnerObjIDApp)
	if err != nil {
		return err
	}
	params[common.BKAppNameField] = common.DefaultAppName
	params[common.BKMaintainersField] = "admin"
	params[common.BKProductPMField] = "admin"
	params[common.BKTimeZoneField] = "Asia/Shanghai"
	params[common.BKLanguageField] = "1" //中文
	params[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal

	byteParams, _ := json.Marshal(params)
	url := fmt.Sprintf("%s/api/%s/biz/default/%s", api.GetAPIResource().APIAddr(), webCommon.API_VERSION, m.OwnerID)
	blog.Info("migrate add default app url :%s", url)
	blog.Info("migrate add default app content :%s", string(byteParams))
	m.httpCli.POST(url, nil, byteParams)
	reply, err := m.httpCli.POST(url, nil, byteParams)
	blog.Info("migrate add default app return :%s", string(reply))
	if err != nil {
		return err
	}

	result := CreateAppResult{}
	err = json.Unmarshal(reply, &result)
	if nil != err {
		return err
	}

	if result.Code != common.CCSuccess {
		return fmt.Errorf("create app faild %s", result.Message)
	}
	return nil
}

func (m *OwnerManager) defaultAppIsExist() (bool, error) {
	params := make(map[string]interface{})
	params["condition"] = make(map[string]interface{})
	params["fields"] = []string{common.BKAppIDField}
	params["start"] = 0
	params["limit"] = 20

	byteParams, _ := json.Marshal(params)
	url := fmt.Sprintf("%s/api/%s/biz/default/%s/search", api.GetAPIResource().APIAddr(), webCommon.API_VERSION, m.OwnerID)

	blog.Info("migrate get default app url :%s", url)
	blog.Info("migrate get default app content :%s", string(byteParams))
	reply, err := m.httpCli.POST(url, nil, byteParams)
	blog.Info("migrate get default app return :%s", string(reply))
	if err != nil {
		return false, err
	}

	result := SearchAppResult{}
	err = json.Unmarshal(reply, &result)
	if nil != err {
		return false, err
	}

	if result.Code != common.CCSuccess {
		return false, fmt.Errorf("search default app err: %s", result.Message)
	}

	if 0 >= result.Data.Count {
		return false, nil
	}
	return true, nil
}

func (m *OwnerManager) getObjectFields(objID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/%s/object/attr/search", api.GetAPIResource().APIAddr(), webCommon.API_VERSION)
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"skip": 0, "limit": common.BKNoLimit}}
	byteParams, _ := json.Marshal(conds)
	blog.Info("migrate get object fields url :%s", url)
	blog.Info("migrate get object fields content :%s", string(byteParams))
	reply, err := m.httpCli.POST(url, nil, byteParams)
	blog.Info("migrate get object fileds return :%s", string(reply))
	if err != nil {
		return nil, err
	}

	replyVal := gjson.ParseBytes(reply)
	if !replyVal.Get("result").Bool() {
		return nil, fmt.Errorf("get object fields faile: %s", replyVal.Get(common.HTTPBKAPIErrorMessage))
	}

	fields := []map[string]interface{}{}
	json.Unmarshal([]byte(replyVal.Get("data").String()), &fields)

	ret := map[string]interface{}{}
	type intOptionType struct {
		Min int
		Max int
	}
	type EnumOptionType struct {
		Name string
		Type string
	}

	for _, mapField := range fields {
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
		}
		blog.Infof("set field %s to %+v", fieldName, ret[fieldName])

	}
	return ret, nil
}

type CreateAppResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"bk_error_code"`
	Message interface{} `json:"bk_error_msg"`
}
