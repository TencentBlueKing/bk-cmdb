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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"

	com "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (lgc *Logics) AddCloudTask(taskList *meta.CloudTaskList, pheader http.Header) (string, error) {
	// TaskName Uniqueness check
	resp, err := lgc.CoreAPI.HostController().Cloud().TaskNameCheck(context.Background(), pheader, taskList)
	if err != nil {
		return "", err
	}

	if resp.Data != 0.0 {
		blog.Errorf("task name %s already exits.", taskList.TaskName)
		errString := "task name " + taskList.TaskName + " already exits."
		return errString, nil
	}

	// Encode secretKey
	taskList.SecretKey = base64.StdEncoding.EncodeToString([]byte(taskList.SecretKey))

	if _, err := lgc.CoreAPI.HostController().Cloud().AddCloudTask(context.Background(), pheader, taskList); err != nil {
		return "", err
	}

	return "", nil
}

func (lgc *Logics) CloudTaskSync(taskList mapstr.MapStr, pheader http.Header) error {
	tickerStart := make(chan bool)
	ticker := time.NewTicker(5 * time.Minute)
	var nextTrigger int64

	PeriodType, errType := taskList.String("bk_period_type")
	if errType != nil {
		blog.Errorf("mapstr interface convert to string failed.")
		return errType
	}
	Period, errP := taskList.String("bk_period")
	if errP != nil {
		blog.Errorf("mapstr interface convert to string failed.")
		return errP
	}

	if PeriodType != "minute" {
		nextTrigger = lgc.UnixSubtract(PeriodType, Period)
	}

	status, errStatus := taskList.Bool("bk_status")
	if errStatus != nil {
		blog.Errorf("mapstr interface convert to bool failed.")
		return errP
	}

	blog.Debug("taskList.Status: %v", status)
	blog.Debug("nextTrigger: %v", nextTrigger)
	blog.Debug("PeriodType: %v", PeriodType)

	if status {
		switch PeriodType {
		case "day":
			timer := time.NewTimer(time.Duration(nextTrigger) * time.Second)
			go func() {
				for {
					select {
					case <-timer.C:
						tickerStart <- true
						blog.Info("case day")
						lgc.ExecSync(taskList, pheader)
					}
				}
			}()
			ticker = time.NewTicker(24 * time.Hour)
			if <-tickerStart {
				go func() {
					for {
						select {
						case <-ticker.C:
							lgc.ExecSync(taskList, pheader)
							blog.Info("case day")
						}

					}
				}()
			}
		case "hour":
			timer := time.NewTimer(time.Duration(nextTrigger) * time.Second)
			go func() {
				for {
					select {
					case <-timer.C:
						tickerStart <- true
						blog.Info("case hour")
						lgc.ExecSync(taskList, pheader)
					}
				}
			}()
			ticker = time.NewTicker(1 * time.Hour)
			if <-tickerStart {
				go func() {
					for {
						select {
						case <-ticker.C:
							lgc.ExecSync(taskList, pheader)
							blog.Info("case hour, Ticker")
						}

					}
				}()
			}
		case "minute":
			go func() {
				for {
					select {
					case <-ticker.C:
						lgc.ExecSync(taskList, pheader)
						blog.Info("case minute")
					}
				}
			}()
		}
	} else {
		ticker.Stop()
		blog.Info("bk_status: false, stop cloud sync")
	}
	return nil
}

func (lgc *Logics) ExecSync(taskList mapstr.MapStr, pheader http.Header) error {
	cloudHistory := new(meta.CloudHistory)

	blog.V(3).Info("start sync")
	taskObjID, errObj := taskList.String("bk_obj_id")
	if errObj != nil {
		blog.Errorf("mapstr key-value convert to string failed.")
		return errObj
	}

	taskID, errTaskID := taskList.Int64("bk_task_id")
	if errTaskID != nil {
		blog.Errorf("mapstr key-value convert to int64 failed.")
		return errTaskID
	}

	cloudHistory.ObjID = taskObjID
	cloudHistory.TaskID = taskID
	startTime := time.Now().Unix()

	defer lgc.CloudHistory(taskID, startTime, cloudHistory, pheader)

	// obtain the hosts from cc_HostBase
	body := new(meta.HostCommonSearch)
	host, err := lgc.SearchHost(pheader, body, false)
	if err != nil {
		blog.Errorf("search host failed, err: %v", err)
		cloudHistory.Status = "失败"
		return err
	}

	existHostList := make([]string, 0)
	for i := 0; i < host.Count; i++ {
		hostInfo, err := mapstr.NewFromInterface(host.Info[i]["host"])
		if err != nil {
			blog.Errorf("get hostInfo failed with err: %v", err)
			cloudHistory.Status = "失败"
			return err
		}

		ip, errH := hostInfo.String(common.BKHostInnerIPField)
		if errH != nil {
			blog.Errorf("get hostIp failed with err: %v")
			cloudHistory.Status = "失败"
			return errH
		}

		existHostList = append(existHostList, ip)
	}

	// obtain hosts from TencentCloud needs secretID and secretKey
	secretID, errS := taskList.String("bk_secret_id")
	if errS != nil {
		blog.Errorf("mapstr convert to string failed.")
		return errS
	}
	secretKeyEncrypted, errKey := taskList.String("bk_secret_key")
	if errKey != nil {
		blog.Errorf("mapstr convert to string failed.")
		return errKey
	}

	decodeBytes, errDecode := base64.StdEncoding.DecodeString(secretKeyEncrypted)
	if errDecode != nil {
		blog.Errorf("Base64 decode secretKey failed.")
		return errDecode
	}
	secretKey := string(decodeBytes)

	// ObtainCloudHosts obtain cloud hosts
	cloudHostInfo, err := lgc.ObtainCloudHosts(secretID, secretKey)
	if err != nil {
		blog.Errorf("obtain cloud hosts failed with err: %v", err)
		cloudHistory.Status = "失败"
		return err
	}

	// pick out the new add cloud hosts
	newAddHost := make([]string, 0)
	newCloudHost := make([]mapstr.MapStr, 0)
	for _, hostInfo := range cloudHostInfo {
		newHostInnerip, ok := hostInfo[common.BKHostInnerIPField].(string)
		if !ok {
			blog.Errorf("interface convert to string failed, err: %v", err)
			cloudHistory.Status = "失败"
		}
		if !util.InStrArr(existHostList, newHostInnerip) {
			newAddHost = append(newAddHost, newHostInnerip)
			newCloudHost = append(newCloudHost, hostInfo)
		}
	}

	// pick out the hosts that has changed attributes
	cloudHostAttr := make([]mapstr.MapStr, 0)
	for _, hostInfo := range cloudHostInfo {
		newHostInnerip, ok := hostInfo[common.BKHostInnerIPField].(string)
		if !ok {
			blog.Errorf("interface convert to string failed, err: %v", err)
			cloudHistory.Status = "失败"
			break
		}
		newHostOuterip, oK := hostInfo[common.BKHostOuterIPField].(string)
		if !oK {
			blog.Errorf("interface convert to string failed, err: %v", err)
			cloudHistory.Status = "失败"
			break
		}
		newHostOsname, _ := hostInfo[common.BKOSNameField].(string)

		for i := 0; i < host.Count; i++ {
			existHostInfo, err := mapstr.NewFromInterface(host.Info[i]["host"])
			if err != nil {
				blog.Errorf("get hostInfo failed with err: %v", err)
				cloudHistory.Status = "失败"
				return err
			}

			existHostIp, ok := existHostInfo.String(common.BKHostInnerIPField)
			if ok != nil {
				blog.Errorf("get hostIp failed with err: %v", ok)
				cloudHistory.Status = "失败"
				break
			}
			existHostOsname, osOk := existHostInfo.String(common.BKOSNameField)
			if osOk != nil {
				blog.Errorf("get os name failed with err: %v", ok)
				cloudHistory.Status = "失败"
				break
			}

			existHostOuterip, ipOk := existHostInfo.String(common.BKHostOuterIPField)
			if ipOk != nil {
				blog.Errorf("get outerip failed with err: %v", ok)
				cloudHistory.Status = "失败"
				break
			}

			existHostID, idOk := existHostInfo.String(common.BKHostIDField)
			if idOk != nil {
				blog.Errorf("get hostID failed with err: %v", ok)
				cloudHistory.Status = "失败"
				break
			}

			if existHostIp == newHostInnerip {
				if existHostOsname != newHostOsname || existHostOuterip != newHostOuterip {
					hostInfo[common.BKHostIDField] = existHostID
					cloudHostAttr = append(cloudHostAttr, hostInfo)
				}
			}
		}
	}

	cloudHistory.NewAdd = len(newAddHost)
	cloudHistory.AttrChanged = len(cloudHostAttr)

	attrConfirm, errAttr := taskList.Bool("bk_attr_confirm")
	if errAttr != nil {
		blog.Errorf("mapstr convert to bool failed.")
		return errAttr
	}

	resourceConfirm, errR := taskList.Bool("bk_confirm")
	if errR != nil {
		blog.Errorf("mapstr convert to string failed.")
		return errR
	}

	if !resourceConfirm && !attrConfirm {
		if len(newAddHost) > 0 {
			err := lgc.AddCloudHosts(pheader, newCloudHost)
			if err != nil {
				blog.Errorf("add cloud hosts failed, err: %v", err)
				cloudHistory.Status = "失败"
				return err
			}
		}
		if len(cloudHostAttr) > 0 {
			err := lgc.UpdateCloudHosts(pheader, cloudHostAttr)
			if err != nil {
				blog.Errorf("update cloud hosts failed, err: %v", err)
				cloudHistory.Status = "失败"
				return err
			}
		}
	}

	if resourceConfirm {
		err := lgc.NewAddConfirm(taskList, pheader, newAddHost, newCloudHost)
		if err != nil {
			blog.Errorf("newly add cloud resource confirm failed, err: %v", err)
			cloudHistory.Status = "失败"
			return err
		}
		cloudHistory.Status = "队列中"
	}

	if attrConfirm && len(cloudHostAttr) > 0 {
		blog.Debug("attr chang")
		for _, host := range cloudHostAttr {
			resourceConfirm := mapstr.MapStr{}
			resourceConfirm["bk_obj_id"] = taskList["bk_obj_id"]
			innerIp, errIp := host.String(common.BKHostInnerIPField)
			if errIp != nil {
				blog.Debug("mapstr.Map convert to string failed.")
				cloudHistory.Status = "失败"
				return errIp
			}

			resourceConfirm[common.BKHostInnerIPField] = innerIp
			resourceConfirm["bk_resource"] = cloudHostAttr
			resourceConfirm["bk_source_type"] = "云同步"
			resourceConfirm["bk_task_id"] = taskList["bk_task_id"]
			resourceConfirm["bk_attr_confirm"] = attrConfirm
			resourceConfirm["bk_confirm"] = false
			resourceConfirm["bk_task_name"] = taskList["bk_task_name"]
			resourceConfirm["bk_account_type"] = taskList["bk_account_type"]
			resourceConfirm["bk_account_admin"] = taskList["bk_account_admin"]

			_, err := lgc.CoreAPI.HostController().Cloud().ResourceConfirm(context.Background(), pheader, resourceConfirm)
			if err != nil {
				blog.Errorf("add resource confirm failed with err: %v", err)
				cloudHistory.Status = "失败"
				return err
			}
		}
		cloudHistory.Status = "队列中"
		return nil
	}

	cloudHistory.Status = "成功"
	blog.V(3).Info("finish sync")
	return nil
}

func (lgc *Logics) AddCloudHosts(pheader http.Header, newCloudHost []mapstr.MapStr) error {
	hostList := new(meta.HostList)
	hostInfoMap := make(map[int64]map[string]interface{}, 0)
	appID := hostList.ApplicationID

	if appID == 0 {
		// get default app id
		var err error
		appID, err = lgc.GetDefaultAppIDWithSupplier(pheader)
		if err != nil {
			blog.Errorf("add host, but get default appid failed, err: %v", err)
			return err
		}
	}

	cond := hutil.NewOperation().WithModuleName(common.DefaultResModuleName).WithAppID(appID).Data()
	cond[common.BKDefaultField] = common.DefaultResModuleFlag
	moduleID, err := lgc.GetResoulePoolModuleID(pheader, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %s", err.Error())
		return err
	}

	blog.V(3).Info("resource confirm add new hosts")
	for index, hostInfo := range newCloudHost {
		if _, ok := hostInfoMap[int64(index)]; !ok {
			hostInfoMap[int64(index)] = make(map[string]interface{}, 0)
		}

		//resource, ok := hostInfo["bk_resource"].([]mapstr.MapStr)
		//if !ok {
		//	blog.Errorf("interface convert to []mapstr.MapStr failed")
		//	break
		//}

		hostInfoMap[int64(index)][common.BKHostInnerIPField] = hostInfo[common.BKHostInnerIPField]
		//hostInfoMap[int64(index)][common.BKHostOuterIPField] = resource[0][common.BKHostOuterIPField]
		//hostInfoMap[int64(index)][common.BKOSNameField] = resource[0][common.BKOSNameField]
		hostInfoMap[int64(index)]["import_from"] = "3"
		hostInfoMap[int64(index)]["bk_cloud_id"] = 1
	}

	succ, updateErrRow, errRow, ok := lgc.AddHost(appID, []int64{moduleID}, util.GetOwnerID(pheader), pheader, hostInfoMap, hostList.InputType)
	if ok != nil {
		blog.Errorf("add host failed, succ: %v, update: %v, err: %v, %v", succ, updateErrRow, ok, errRow)
		return ok
	}

	return nil
}

func (lgc *Logics) UpdateCloudHosts(pheader http.Header, cloudHostAttr []mapstr.MapStr) error {
	for _, hostInfo := range cloudHostAttr {
		hostID, err := hostInfo.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("hostID convert to string failed")
			return err
		}

		delete(hostInfo, common.BKHostIDField)
		delete(hostInfo, "bk_confirm")
		delete(hostInfo, "bk_attr_confirm")
		opt := mapstr.MapStr{"condition": mapstr.MapStr{common.BKHostIDField: hostID}, "data": hostInfo}

		blog.V(3).Info("opt: %v", opt)
		result, err := lgc.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDHost, pheader, opt)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("update host batch failed, ids[%v], err: %v, %v", hostID, err, result.ErrMsg)
			return err
		}
	}
	return nil
}

func (lgc *Logics) NewAddConfirm(taskList mapstr.MapStr, pheader http.Header, newAddHost []string, newCloudHost []mapstr.MapStr) error {
	// Check whether the host is already exist in resource confirm.
	opt := make(map[string]interface{})
	confirmHosts, errS := lgc.CoreAPI.HostController().Cloud().SearchConfirm(context.Background(), pheader, opt)
	if errS != nil {
		blog.Errorf("get confirm info failed with err: %v", errS)
		return errS
	}

	confirmIpList := make([]string, 0)
	if confirmHosts.Count > 0 {
		for _, confirmInfo := range confirmHosts.Info {
			for _, ip := range confirmInfo[common.BKHostInnerIPField].([]string) {
				confirmIpList = append(confirmIpList, ip)
			}
		}
	}

	newHostIp := make([]string, 0)
	for _, ip := range newAddHost {
		if !util.InStrArr(confirmIpList, ip) {
			newHostIp = append(newHostIp, ip)
		}
	}

	// newly added cloud hosts confirm
	if len(newHostIp) > 0 {
		for _, innerIp := range newHostIp {
			resourceConfirm := mapstr.MapStr{}
			resourceConfirm["bk_obj_id"] = taskList["bk_obj_id"]
			resourceConfirm[common.BKHostInnerIPField] = innerIp
			resourceConfirm["bk_source_type"] = "云同步"
			resourceConfirm["bk_task_id"] = taskList["bk_task_id"]
			resourceConfirm["bk_resource"] = newCloudHost
			resourceConfirm["bk_confirm"] = true
			resourceConfirm["bk_attr_confirm"] = false
			resourceConfirm["bk_task_name"] = taskList["bk_task_name"]
			resourceConfirm["bk_account_type"] = taskList["bk_account_type"]
			resourceConfirm["bk_account_admin"] = taskList["bk_account_admin"]

			_, err := lgc.CoreAPI.HostController().Cloud().ResourceConfirm(context.Background(), pheader, resourceConfirm)
			if err != nil {
				blog.Errorf("add resource confirm failed with err: %v", err)
				return err
			}
		}
	}
	return nil
}

func (lgc *Logics) UnixSubtract(periodType string, period string) int64 {
	timeLayout := "2006-01-02 15:04:05" // transfer model
	toBeCharge := period
	var unixSubtract int64
	nowStr := time.Unix(time.Now().Unix(), 0).Format(timeLayout)

	blog.Debug("periodType: %v", periodType)
	blog.Debug("period: %v", period)
	if periodType == "day" {
		intHour, _ := strconv.Atoi(toBeCharge[:2])
		intMinute, _ := strconv.Atoi(toBeCharge[3:])
		if intHour > time.Now().Hour() {
			toBeCharge = fmt.Sprintf("%s%s%s", nowStr[:11], toBeCharge, ":00")
		}
		if intHour < time.Now().Hour() {
			toBeCharge = fmt.Sprintf("%s%d %s%s", nowStr[:8], time.Now().Day()+1, toBeCharge, ":00")
		}
		if intHour == time.Now().Hour() && intMinute > time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%s%s", nowStr[:11], toBeCharge, ":00")
		}
		if intHour == time.Now().Hour() && intMinute <= time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d %s%s", nowStr[:8], time.Now().Day()+1, toBeCharge, ":00")
		}

		loc, _ := time.LoadLocation("Local")
		theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
		sr := theTime.Unix()
		unixSubtract = sr - time.Now().Unix()
	}

	if periodType == "hour" {
		intToBeCharge, err := strconv.Atoi(toBeCharge)
		if err != nil {
			blog.Errorf("period transfer to int failed with err: %v", err)
			return 0
		}

		if intToBeCharge >= 10 && intToBeCharge > time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d:%s:%s", nowStr[:11], time.Now().Hour(), toBeCharge, "00")
		}
		if intToBeCharge >= 10 && intToBeCharge < time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d:%s:%s", nowStr[:11], time.Now().Hour()+1, toBeCharge, "00")
		}
		if intToBeCharge < 10 && intToBeCharge > time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d:0%s:%s", nowStr[:11], time.Now().Hour(), toBeCharge, "00")
		}
		if intToBeCharge < 10 && intToBeCharge < time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d:0%s:%s", nowStr[:11], time.Now().Hour()+1, toBeCharge, "00")
		}

		loc, _ := time.LoadLocation("Local")
		theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
		sr := theTime.Unix()
		unixSubtract = sr - time.Now().Unix()
	}

	return unixSubtract
}

func (lgc *Logics) CloudHistory(taskID int64, startTime int64, cloudHistory *meta.CloudHistory, pheader http.Header) error {
	finishTime := time.Now().Unix()
	timeConsumed := finishTime - startTime
	if timeConsumed > 60 {
		minute := timeConsumed / 60
		seconds := timeConsumed % 60
		cloudHistory.TimeConsume = fmt.Sprintf("%dmin%ds", minute, seconds)
	} else {
		cloudHistory.TimeConsume = fmt.Sprintf("%ds", timeConsumed)
	}

	timeLayout := "2006-01-02 15:04:05" // transfer model
	startTimeStr := time.Unix(startTime, 0).Format(timeLayout)
	cloudHistory.StartTime = startTimeStr

	blog.V(3).Info(cloudHistory.TimeConsume)

	updateData := mapstr.MapStr{}
	updateTime := time.Now()
	updateData["bk_last_sync_time"] = updateTime
	updateData["bk_task_id"] = taskID
	updateData["bk_sync_status"] = cloudHistory.Status
	updateData["new_add"] = cloudHistory.NewAdd
	updateData["attr_changed"] = cloudHistory.AttrChanged

	if _, err := lgc.CoreAPI.HostController().Cloud().UpdateCloudTask(context.Background(), pheader, updateData); err != nil {
		blog.Errorf("update task failed with decode body err: %v", err)
		return err
	}

	if _, err := lgc.CoreAPI.HostController().Cloud().CloudHistory(context.Background(), pheader, cloudHistory); err != nil {
		blog.Errorf("add cloud history table failed, err: %v", err)
		return err
	}

	return nil
}

func (lgc *Logics) ObtainCloudHosts(secretID string, secretKey string) ([]map[string]interface{}, error) {
	credential := com.NewCredential(
		secretID,
		secretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "GET"
	cpf.HttpProfile.ReqTimeout = 10
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"

	ClientRegion, _ := cvm.NewClient(credential, regions.Guangzhou, cpf)
	regionRequest := cvm.NewDescribeRegionsRequest()
	Response, err := ClientRegion.DescribeRegions(regionRequest)

	if err != nil {
		return nil, err
	}

	data := Response.ToJsonString()
	regionResponse := new(meta.RegionResponse)
	if err := json.Unmarshal([]byte(data), regionResponse); err != nil {
		blog.Errorf("json unmarsha1 error :%v\n", err)
		return nil, err
	}

	cloudHostInfo := make([]map[string]interface{}, 0)
	for _, region := range regionResponse.Response.Data {
		var inneripList string
		var outeripList string
		var osName string
		regionHosts := make(map[string]interface{})

		client, _ := cvm.NewClient(credential, region.Region, cpf)
		instRequest := cvm.NewDescribeInstancesRequest()
		response, err := client.DescribeInstances(instRequest)

		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s", err)
			return nil, err
		}
		if err != nil {
			panic(err)
		}

		data := response.ToJsonString()
		Hosts := meta.HostResponse{}
		if err := json.Unmarshal([]byte(data), &Hosts); err != nil {
			fmt.Printf("json unmarsha1 error :%v\n", err)
		}

		instSet := Hosts.HostResponse.InstanceSet
		for _, obj := range instSet {
			osName = obj.OsName
			if len(obj.PrivateIpAddresses) > 0 {
				inneripList = obj.PrivateIpAddresses[0]
			}
		}

		for _, obj := range instSet {
			if len(obj.PublicIpAddresses) > 0 {
				outeripList = obj.PublicIpAddresses[0]
			}
		}

		if len(instSet) > 0 {
			regionHosts["bk_cloud_region"] = region.Region
			regionHosts["bk_host_innerip"] = inneripList
			regionHosts["bk_host_outerip"] = outeripList
			regionHosts["bk_os_name"] = osName
			cloudHostInfo = append(cloudHostInfo, regionHosts)
		}
	}
	return cloudHostInfo, nil
}
