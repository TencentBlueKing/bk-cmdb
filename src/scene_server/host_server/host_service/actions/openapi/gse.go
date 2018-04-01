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
 
package openapi

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	myCommon "configcenter/src/scene_server/host_server/common"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
	"github.com/samuel/go-zookeeper/zk"
	"gopkg.in/redis.v5"
	//"configcenter/src/source_controller/api/client"
)

var gse *gseAction = &gseAction{}

var (
	redisIp string        = ""
	client  *redis.Client = nil
)

type gseAction struct {
	base.BaseAction
}

// GetAgentStatus: 获取指定业务下agent正常和异常的主机列表
func (cli *gseAction) GetAgentStatus(req *restful.Request, resp *restful.Response) {
	// 获取AppID
	pathParams := req.PathParameters()
	appID, err := strconv.Atoi(pathParams["appid"])
	if nil != err {
		blog.Error("GetAgentStatus error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}

	hosts, err := getHostsByAppID(req, appID)
	blog.Infof("get agent status hosts:%v", hosts)
	if nil != err {
		blog.Error("GetAgentStatus error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, err, resp)
		return
	}

	hostDataArr := make([]interface{}, 0)
	for _, host := range hosts {
		hostMap := host.(map[string]interface{})
		companyID := 0

		subArea, _ := util.GetInt64ByInterface(hostMap[common.BKCloudIDField])
		platID := int(subArea)

		/*if platID == 1 {
			platID = common.BKDefaultDirSubArea
		}*/

		ip, ok := hostMap[common.BKHostInnerIPField].(string)
		if !ok {
			blog.Error("InnerIP is not ok")
			return
		}

		intIp := ip2long(ip)
		comID := platID<<22 + companyID
		agentFlag := fmt.Sprintf("agentalive_cloudid_%d", comID)
		cellData := map[string]interface{}{"agentFlag": agentFlag, "offset": intIp}

		hostDataArr = append(hostDataArr, cellData)
	}
	blog.Infof("get gse hostDataArr:%v", hostDataArr)
	agentStatus, err := getGseAgentStatus(hostDataArr)
	if nil != err {
		blog.Error("getGseAgentStatus error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}
	agentNorCnt := 0
	agentAbnorCnt := 0

	agentNorList := make([]map[string]interface{}, 0)
	agentAbnorList := make([]map[string]interface{}, 0)
	i := 0

	blog.Debug("agentStatus:%v", agentStatus)
	agentStatuLen := len(agentStatus)
	for _, host := range hosts {
		hostMap := host.(map[string]interface{})
		platIdInt, err := util.GetIntByInterface(hostMap[common.BKCloudIDField])
		if nil != err {
			blog.Error("SubArea is not ok")
			return
		}
		platIdStr := strconv.Itoa(platIdInt)
		hostMapTemp := map[string]interface{}{
			"Ip":        hostMap[common.BKHostInnerIPField],
			"CompanyID": 0,
			"PlatID":    platIdStr,
		}
		status := int64(0)
		if i < agentStatuLen {
			status = agentStatus[i]
		}
		if status == 1 {
			agentNorCnt++
			agentNorList = append(agentNorList, hostMapTemp)
		} else {
			agentAbnorCnt++
			agentAbnorList = append(agentAbnorList, hostMapTemp)
		}

		i++
	}

	resData := map[string]interface{}{
		"agentNorCnt":    fmt.Sprintf("%d", agentNorCnt),
		"agentAbnorCnt":  fmt.Sprintf("%d", agentAbnorCnt),
		"agentNorList":   agentNorList,
		"agentAbnorList": agentAbnorList,
	}

	cli.ResponseSuccess(resData, resp)
}

//ip2long
func ip2long(ip string) int64 {
	bits := strings.Split(ip, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

//getGseAgentStatus
func getGseAgentStatus(hostDataArr []interface{}) ([]int64, error) {
	blog.Infof("getGseAgentStatus hostDataArr1:%v", hostDataArr)
	if len(hostDataArr) == 0 {
		return []int64{}, nil
	}

	client, err := getRedisSession()
	if nil == client {
		blog.Error("getRedisSession error:%v", err)
		return nil, err
	}
	length := len(hostDataArr)
	data := make([]int64, length)
	pipe := client.Pipeline()
	for _, hostData := range hostDataArr {
		hostDataMap := hostData.(map[string]interface{})
		blog.Infof("get gse hostDataMap:%v", hostDataMap)
		pipe.GetBit(hostDataMap["agentFlag"].(string), hostDataMap["offset"].(int64))
	}

	result, err := pipe.Exec()
	if err == redis.Nil {
		return []int64{}, nil
	} else if err != nil {
		blog.Error("redis get bit error %s, hostData:%v", err.Error(), hostDataArr)
		return []int64{}, err

	} else {
		for i, re := range result {

			reArr := strings.Split(re.String(), ":")
			var valTemp string
			if len(reArr) == 2 {
				valTemp = reArr[1]
			} else {
				valTemp = reArr[0]
			}
			val := strings.Replace(valTemp, " ", "", -1)

			data[i], err = strconv.ParseInt(val, 10, 64)
			if nil != err {
				blog.Error("get bit error %s, re:%v", err, re)
				return nil, err
			}
		}
	}
	return data, nil

}

//getRedisSession
func getRedisSession() (*redis.Client, error) {
	newIp, port, auth, err := getRedisIP()
	blog.Error(fmt.Sprintf("newIp:%v port:%v auth:%v err:%v", newIp, port, auth, err))

	if "" == redisIp && "" == newIp || "" == fmt.Sprintf("%v", port) {
		blog.Errorf("get gse redis ip error:ip:%s, port:%s", newIp, port)
		return nil, err
	}
	if nil != client && newIp == redisIp {
		return client, errors.New("redis config error ")
	}
	if nil != client {
		client.Close()
	}
	redisIp = newIp
	blog.Infof("redisIp:%v", redisIp)
	client = redis.NewClient(&redis.Options{
		Addr:         redisIp + ":" + port,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		Password:     auth,
		DB:           0,
	})
	return client, nil
}

var zkClient *zk.Conn = nil

//getRedisIP
func getRedisIP() (string, string, string, error) {

	// TODO: 放到配置文件
	GSE_ZK_HOST, AUTH_USER, AUTH_PWD, port, redisAuth := myCommon.GetSetConfig() //
	blog.Errorf("%v %s %s %s %s", GSE_ZK_HOST, AUTH_USER, AUTH_PWD, port, redisAuth)
	var err error
	if nil == zkClient {
		zkClient, _, err = zk.Connect(GSE_ZK_HOST, time.Second*60)
		if err != nil {
			return "", "", "", err
		}
	}
	if zkClient.State() == zk.StateExpired {
		zkClient, _, err = zk.Connect(GSE_ZK_HOST, time.Second*60)
		if err != nil {
			return "", "", "", err
		}
	}

	// AddAuth
	auth := AUTH_USER + ":" + AUTH_PWD
	if err := zkClient.AddAuth("digest", []byte(auth)); err != nil {
		zkClient.Close()
		return "", "", "", err
	}

	path := "/gse/config/server/dbproxy/newall"
	result, _, err := zkClient.Children(path)
	if nil != err {
		blog.Errorf("get redis id from gse zk error:%s", err.Error())
		return "", "", "", err
	}
	minData := result[0]
	for _, re := range result {
		if re < minData {
			minData = re
		}
	}

	fullPath := "/gse/config/server/dbproxy/newall/" + minData
	resIp, _, err := zkClient.Get(fullPath)

	return string(resIp), port, redisAuth, nil
}

//getHostsByAppID
func getHostsByAppID(req *restful.Request, appID int) ([]interface{}, error) {

	configData, err := getConfigByCond(req, gse.CC.HostCtrl(), map[string]interface{}{
		common.BKAppIDField: []interface{}{appID},
	})

	if nil != err {
		return nil, err
	}

	hostIDArr := make([]int, 0)

	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	hostMapCondition := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			"$in": hostIDArr,
		},
	}

	url := gse.CC.HostCtrl() + "/host/v1/hosts/search"
	searchParams := map[string]interface{}{
		"fields":    fmt.Sprintf("%s,%s,%s,%s", common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField, common.BKOwnerIDField),
		"condition": hostMapCondition,
	}
	inputJson, _ := json.Marshal(searchParams)

	appInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	if nil != err {
		return nil, err
	}

	js, err := simplejson.NewJson([]byte(appInfo))
	if nil != err {
		return nil, err
	}

	res, err := js.Map()
	if nil != err {
		return nil, err
	}

	if !res["result"].(bool) {
		return nil, errors.New(res["message"].(string))
	}

	resData := res["data"].(map[string]interface{})
	resDataInfo := resData["info"].([]interface{})

	return resDataInfo, nil
}

//getConfigByCond
func getConfigByCond(req *restful.Request, hostURL string, cond interface{}) ([]map[string]int, error) {

	configArr := make([]map[string]int, 0)
	bodyContent, _ := json.Marshal(cond)
	url := hostURL + "/host/v1/meta/hosts/module/config/search"
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	if err != nil {
		return configArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, err := js.Map()
	configData := output["data"]
	configInfo, _ := configData.([]interface{})
	for _, mh := range configInfo {
		celh := mh.(map[string]interface{})

		hostID, _ := celh[common.BKHostIDField].(json.Number).Int64()
		setID, _ := celh[common.BKSetIDField].(json.Number).Int64()
		moduleID, _ := celh[common.BKModuleIDField].(json.Number).Int64()
		appID, _ := celh[common.BKAppIDField].(json.Number).Int64()
		data := make(map[string]int)
		data[common.BKAppIDField] = int(appID)
		data[common.BKSetIDField] = int(setID)
		data[common.BKModuleIDField] = int(moduleID)
		data[common.BKHostIDField] = int(hostID)
		configArr = append(configArr, data)
	}
	return configArr, nil
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "getAgentStatus/{appid}", Params: nil, Handler: gse.GetAgentStatus})
	// create CC object
	gse.CreateAction()
}
