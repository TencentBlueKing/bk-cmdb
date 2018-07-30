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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/app/options"
	"github.com/samuel/go-zookeeper/zk"
	"gopkg.in/redis.v5"
)

func (lgc *Logics) GetAgentStatus(appID int64, gseConfg *options.Gse, header http.Header) (*meta.GetAgentStatusResult, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	configCon := map[string][]int64{
		common.BKAppIDField: []int64{appID},
	}
	configData, err := lgc.GetConfigByCond(header, configCon)
	if nil != err {
		blog.Errorf("GetAgentStatus  GetConfigByCond error : %s, input:%v", err.Error(), appID)
		return nil, defErr.Errorf(common.CCErrHostModuleConfigFaild)
	}
	hostIDArr := make([]int64, 0)
	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	hosts, err := lgc.GetHostInfoByConds(header, nil)
	blog.V(3).Infof("GetAgentStatus GetHostInfoByConds get agent status hosts:%v, input:%v", hosts, appID)
	if nil != err {
		blog.Error("GetAgentStatus error :%v", err)
		return nil, defErr.Errorf(common.CCErrHostGetFail)
	}

	hostDataArr := make([]interface{}, 0)
	for _, host := range hosts {
		companyID := 0

		subArea, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if nil != err {
			blog.Errorf("GetAgentStatus get agent status hosts not find cloud id, rawHost:%v, input:%v", host, appID)
			return nil, defErr.Errorf(common.CCErrHostGetFail)
		}
		platID := int(subArea)

		ip, ok := host[common.BKHostInnerIPField].(string)
		if !ok {
			blog.Errorf("GetAgentStatus get agent status hosts not find innerip, rawHost:%v, input:%v", host, appID)
			return nil, defErr.Errorf(common.CCErrHostGetFail)

		}

		intIp := ip2long(ip)
		comID := platID<<22 + companyID
		agentFlag := fmt.Sprintf("agentalive_cloudid_%d", comID)
		cellData := map[string]interface{}{"agentFlag": agentFlag, "offset": intIp}

		hostDataArr = append(hostDataArr, cellData)
	}

	phpapi := lgc.NewPHPAPI(header)
	blog.Infof("get gse hostDataArr:%v", hostDataArr)
	agentStatus, err := phpapi.getGseAgentStatus(hostDataArr, gseConfg)
	if nil != err {
		blog.Error("getGseAgentStatus error :%s, input:%d", err.Error(), appID)
		return nil, defErr.Errorf(common.CCErrHostAgentStatusFail, err.Error())
	}
	agentNorCnt := 0
	agentAbnorCnt := 0

	agentNorList := make([]map[string]interface{}, 0)
	agentAbnorList := make([]map[string]interface{}, 0)
	idx := 0

	blog.V(3).Infof("agentStatus:%v", agentStatus)
	agentStatuLen := len(agentStatus)
	for _, host := range hosts {
		platIdInt, err := util.GetIntByInterface(host[common.BKCloudIDField])
		if nil != err {
			blog.Errorf("GetAgentStatus get agent status hosts not find cloud id, rawHost:%v, input:%v", host, appID)
			return nil, defErr.Errorf(common.CCErrHostGetFail)
		}
		platIdStr := strconv.Itoa(platIdInt)
		hostMapTemp := map[string]interface{}{
			"Ip":        host[common.BKHostInnerIPField],
			"CompanyID": 0,
			"PlatID":    platIdStr,
			"CompanyId": 0,
			"PlatId":    platIdStr,
		}
		status := int64(0)
		if idx < agentStatuLen {
			status = agentStatus[idx]
		}
		if status == 1 {
			agentNorCnt++
			agentNorList = append(agentNorList, hostMapTemp)
		} else {
			agentAbnorCnt++
			agentAbnorList = append(agentAbnorList, hostMapTemp)
		}

		idx++
	}
	ret := &meta.GetAgentStatusResult{
		AgentNorCnt:    agentNorCnt,
		AgentAbnorCnt:  agentAbnorCnt,
		AgentNorList:   agentNorList,
		AgentAbnorList: agentAbnorList,
	}

	return ret, nil

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

var (
	redisIp  string        = ""
	client   *redis.Client = nil
	zkClient *zk.Conn      = nil
)

//getRedisSession
func (phpapi *PHPAPI) getRedisSession(gseConfg *options.Gse) (*redis.Client, error) {
	newIp, port, auth, err := phpapi.getRedisIP(gseConfg)
	blog.V(3).Infof("newIp:%v port:%v err:%v", newIp, port, err)

	if "" == redisIp && "" == newIp || "" == port {
		blog.Errorf("get gse redis ip error:ip:%s, port:%s", newIp, port)
		return nil, err
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

//getRedisIP
func (phpapi *PHPAPI) getRedisIP(gseConfg *options.Gse) (string, string, string, error) {

	zkAddrArr := strings.Split(gseConfg.ZkAddress, ",")
	// TODO: 放到配置文件
	var err error
	if nil == zkClient {
		zkClient, _, err = zk.Connect(zkAddrArr, time.Second*60)
		if err != nil {
			return "", "", "", err
		}
	}
	if zkClient.State() == zk.StateExpired {
		zkClient, _, err = zk.Connect(zkAddrArr, time.Second*60)
		if err != nil {
			return "", "", "", err
		}
	}

	// AddAuth
	auth := gseConfg.ZkUser + ":" + gseConfg.ZkPassword
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

	return string(resIp), gseConfg.RedisPort, gseConfg.RedisPassword, nil
}

//getGseAgentStatus
func (phpapi *PHPAPI) getGseAgentStatus(hostDataArr []interface{}, gseConfg *options.Gse) ([]int64, error) {
	blog.Infof("getGseAgentStatus hostDataArr1:%v", hostDataArr)
	if len(hostDataArr) == 0 {
		return []int64{}, nil
	}

	rdClient, err := phpapi.getRedisSession(gseConfg)
	if nil == rdClient {
		blog.Errorf("getGseAgentStatus getRedisSession error:%v, hostDataArr:%v", err, hostDataArr)
		return nil, err
	}
	length := len(hostDataArr)
	data := make([]int64, length)
	pipe := rdClient.Pipeline()
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
