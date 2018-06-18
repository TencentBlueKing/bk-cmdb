package middleware

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	webCommon "configcenter/src/web_server/common"
	"encoding/json"
	"fmt"
)

type OwnerManager struct {
	httpCli  *httpclient.HttpClient
	OwnerID  string
	UserName string
	APIAddr  string
}

func NewOwnerManager(userName, APIAddr, ownerID, language string) (*OwnerManager, error) {
	ownerManager := new(OwnerManager)
	ownerManager.UserName = userName
	ownerManager.APIAddr = APIAddr
	ownerManager.OwnerID = ownerID
	ownerManager.httpCli = httpclient.NewHttpClient()
	ownerManager.httpCli.SetHeader(common.BKHTTPHeaderUser, userName)
	ownerManager.httpCli.SetHeader(common.BKHTTPLanguage, language)
	ownerManager.httpCli.SetHeader(common.BKHTTPOwnerID, ownerID)
	return ownerManager, nil
}

func (m OwnerManager) InitOwner(ownerID string) {

}

func (m OwnerManager) addDefaultApp(ownerID string) error {
	params := map[string]interface{}{}
	params[common.BKAppNameField] = common.DefaultAppName
	params[common.BKMaintainersField] = "admin"
	params[common.BKProductPMField] = "admin"
	params[common.BKTimeZoneField] = "Asia/Shanghai"
	params[common.BKLanguageField] = "1" //中文
	params[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal

	byteParams, _ := json.Marshal(params)
	url := fmt.Sprintf("%s/api/%s/topo/app/default/%s", m.APIAddr, webCommon.API_VERSION, m.OwnerID)
	blog.Info("migrate add default app url :%s", url)
	blog.Info("migrate add default app content :%s", string(byteParams))
	m.httpCli.POST(url, nil, byteParams)
	reply, err := m.httpCli.POST(url, nil, byteParams)
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
