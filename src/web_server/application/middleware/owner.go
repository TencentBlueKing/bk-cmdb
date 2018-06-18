package middleware

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/http/httpclient"
	webCommon "configcenter/src/web_server/common"
	"encoding/json"
	"fmt"
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
	params := map[string]interface{}{}
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

type CreateAppResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"bk_error_code"`
	Message interface{} `json:"bk_error_msg"`
}
