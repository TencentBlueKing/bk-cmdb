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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/thirdpartyclient/esbserver/nodeman"
)

// Netdevicebeat netdevicebeat collector name
const Netdevicebeat = "netdevicebeat"

func (lgc *Logics) SearchCollector(header http.Header, cond metadata.ParamNetcollectorSearch) (int64, []metadata.Netcollector, error) {
	collectors := []metadata.Netcollector{}

	// fetch package info
	packageResp, err := lgc.ESB.NodemanSrv().SearchPackage(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] SearchPackage failed: %v", err)
		return 0, nil, err
	}
	if !packageResp.Result {
		blog.Errorf("[NetDevice][SearchCollector] SearchPackage failed: %+v", packageResp)
		return 0, nil, fmt.Errorf("search plugin host from nodeman failed: %s", packageResp.Message)
	}
	var pkg nodeman.PluginPackage
	if len(packageResp.Data) > 0 {
		pkg = packageResp.Data[0]
	}

	// fetch hosts
	pluginHostResp, err := lgc.ESB.NodemanSrv().SearchPluginHost(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] SearchPluginHost failed: %v", err)
		return 0, nil, err
	}
	if !pluginHostResp.Result {
		blog.Errorf("[NetDevice][SearchCollector] SearchPluginHost failed: %v", pluginHostResp.Message)
		return 0, nil, fmt.Errorf("search plugin host from nodeman failed: %s", pluginHostResp.Message)
	}

	// build collectors
	cloudIDs := []int64{}
	ips := []string{}
	for _, phost := range pluginHostResp.Data {
		cloudIDs = append(cloudIDs, phost.Host.BkCloudID)
		ips = append(ips, phost.Host.InnerIP)

		collectorStatus := metadata.CollectorStatusAbnormal
		if strings.ToUpper(phost.Status) == "RUNNING" {
			collectorStatus = metadata.CollectorStatusNormal
		}

		collector := metadata.Netcollector{
			BizID:   phost.Host.BkBizID,
			CloudID: phost.Host.BkCloudID,
			InnerIP: phost.Host.InnerIP,
			Version: phost.Version,
			Status: metadata.NetcollectorStatus{
				CollectorStatus: collectorStatus,
			},
			LatestVersion: pkg.Version,
		}
		collectors = append(collectors, collector)
	}

	cloudCond := condition.CreateCondition()
	cloudCond.Field(common.BKCloudIDField).In(cloudIDs)
	cloudMap, err := lgc.findInstMap(header, common.BKInnerObjIDPlat, &metadata.QueryInput{Condition: cloudCond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] find clouds failed: %v", err)
		return 0, nil, err
	}

	cloudCond.Field(common.BKHostInnerIPField).In(ips)
	collectorMap, err := lgc.findCollectorMap(cloudCond.ToMapStr())
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] get collector config failed")
		return 0, nil, err
	}

	// fill config field from our db
	for index := range collectors {
		collector := &collectors[index]

		if clodInst, ok := cloudMap[collector.CloudID]; ok {
			cloudname, err := clodInst.String(common.BKCloudNameField)
			if err != nil {
				blog.Errorf("[NetDevice][SearchCollector] bk_cloud_name field invalied: %v", err)
			}
			collector.CloudName = cloudname
		}

		cond := condition.CreateCondition()
		cond.Field(common.BKCloudIDField).Eq(collector.CloudID)

		key := fmt.Sprintf("%d:%s", collector.CloudID, collector.InnerIP)
		existsOne, ok := collectorMap[key]
		if !ok {
			blog.Errorf("[NetDevice][SearchCollector] get collector config for %s failed", key)
		}
		collector.Config = existsOne.Config
		collector.ReportTotal = existsOne.ReportTotal

		if existsOne.Status.ConfigStatus == "" {
			existsOne.Status.ConfigStatus = metadata.CollectorConfigStatusAbnormal
		}
		if existsOne.Status.ReportStatus == "" {
			existsOne.Status.ReportStatus = metadata.CollectorReportStatusAbnormal
		}
		collector.Status.ConfigStatus = existsOne.Status.ConfigStatus
		collector.Status.ReportStatus = existsOne.Status.ReportStatus
	}
	return int64(len(collectors)), collectors, nil
}

func (lgc *Logics) findCollectorMap(cond interface{}) (map[string]metadata.Netcollector, error) {
	collectors := []metadata.Netcollector{}
	err := lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectConfig, nil, cond, &collectors, "", 0, 0)
	if err != nil {
		return nil, err
	}
	collectorMap := map[string]metadata.Netcollector{}
	for index := range collectors {
		key := fmt.Sprintf("%d:%s", collectors[index].CloudID, collectors[index].InnerIP)
		collectorMap[key] = collectors[index]
	}
	return collectorMap, nil
}

func (lgc *Logics) UpdateCollector(header http.Header, config metadata.Netcollector) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKCloudIDField).Eq(config.CloudID)
	cond.Field(common.BKHostInnerIPField).Eq(config.InnerIP)

	count, err := lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectConfig, cond.ToMapStr())
	if err != nil {
		blog.Errorf("[UpdateCollector] count error: %v", err)
		return err
	}
	if count > 0 {
		err = lgc.Instance.UpdateByCondition(common.BKTableNameNetcollectConfig, config, cond)
		if err != nil {
			blog.Errorf("[UpdateCollector] UpdateByCondition error: %v", err)
			return err
		}
		return nil
	}

	_, err = lgc.Instance.Insert(common.BKTableNameNetcollectConfig, config)
	if err != nil {
		blog.Errorf("[UpdateCollector] UpdateByCondition error: %v", err)
		return err
	}

	return lgc.DiscoverNetDevice(header, []metadata.Netcollector{config})
}

func (lgc *Logics) DiscoverNetDevice(header http.Header, configs []metadata.Netcollector) error {

	// fetch global_params
	pkgResp, err := lgc.ESB.NodemanSrv().SearchPackage(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchPackage failed %v", err)
		return err
	}
	if !pkgResp.Result {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchPackage failed %s", pkgResp.Message)
		return fmt.Errorf("search plugin host from nodeman failed: %s", pkgResp.Message)
	}
	var pkg nodeman.PluginPackage
	if len(pkgResp.Data) > 0 {
		pkg = pkgResp.Data[0]
	}

	procResp, err := lgc.ESB.NodemanSrv().SearchProcess(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchProcess failed %v", err)
		return err
	}
	if !procResp.Result {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchProcess failed %v", pkgResp.Message)
		return fmt.Errorf("search plugin host from nodeman failed: %s", procResp.Message)
	}
	procInfoResp, err := lgc.ESB.NodemanSrv().SearchProcessInfo(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchProcess failed %v", err)
		return err
	}
	if !procInfoResp.Result {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchProcess failed %v", pkgResp.Message)
		return fmt.Errorf("search plugin host from nodeman failed: %s", procInfoResp.Message)
	}

	cloudIDs := []int64{}
	ips := []string{}
	for _, config := range configs {
		cloudIDs = append(cloudIDs, config.CloudID)
		ips = append(ips, config.InnerIP)
	}
	cloudCond := condition.CreateCondition()
	cloudCond.Field(common.BKCloudIDField).In(cloudIDs)
	cloudCond.Field(common.BKHostInnerIPField).In(ips)
	collectorMap, err := lgc.findCollectorMap(cloudCond.ToMapStr())
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] get collector config failed, %v", err)
		return err
	}

	for _, config := range configs {
		key := fmt.Sprintf("%d:%s", config.CloudID, config.InnerIP)
		collector, ok := collectorMap[key]
		if !ok {
			blog.Errorf("[NetDevice][DiscoverNetDevice] get collector config for %s failed", key)
			continue
		}

		upgradeReq, err := lgc.buildUpgradePluginRequest(&collector, util.GetUser(header), &pkg, &procResp.Data, &procInfoResp.Data)
		if err != nil {
			blog.Errorf("[NetDevice][DiscoverNetDevice] get collector config for %s failed, %v", key, err)
		}
		blog.InfoJSON("[NetDevice][DiscoverNetDevice] UpgradePlugin request %s", upgradeReq)
		upgradeResp, err := lgc.ESB.NodemanSrv().UpgradePlugin(context.Background(), header, strconv.FormatInt(collector.BizID, 10), upgradeReq)
		if err != nil {
			blog.Errorf("[NetDevice][DiscoverNetDevice] get collector config for %s failed", key)
			continue
		}
		if !upgradeResp.Result {
			blog.Errorf("NetDevice][DiscoverNetDevice] search plugin host from nodeman failed: %s", upgradeResp.Message)
			continue
		}

		blog.V(3).Infof("[NetDevice][DiscoverNetDevice] UpgradePlugin response %+v ", upgradeResp)
		if err := lgc.saveCollectTask(&collector, ""); err != nil {
			blog.Errorf("[NetDevice][DiscoverNetDevice] saveCollectTask %s failed, %v", key, err)
		}
	}

	return err
}

func (lgc *Logics) saveCollectTask(collector *metadata.Netcollector, taskID string) error {
	data := mapstr.MapStr{}
	data.Set("task_id", taskID)

	cond := condition.CreateCondition()
	cond.Field(common.BKCloudIDField).Eq(collector.CloudID)
	cond.Field(common.BKHostInnerIPField).In(collector.InnerIP)

	return lgc.Instance.UpdateByCondition(common.BKTableNameNetcollectConfig, data, cond.ToMapStr())
}

func (lgc *Logics) buildUpgradePluginRequest(collector *metadata.Netcollector, user string, pkg *nodeman.PluginPackage, proc *nodeman.PluginProcess, porcInfo *nodeman.PluginProcessInfo) (*nodeman.UpgradePluginRequest, error) {
	pluginConfig, err := lgc.buildNetdevicebeatConfigFile(collector)
	if err != nil {
		return nil, err
	}

	upgradeReq := nodeman.UpgradePluginRequest{
		Creator:   user,
		BkCloudID: strconv.FormatInt(collector.CloudID, 10),
		NodeType:  "PLUGIN",
		OpType:    "UPDATE",
		Hosts: []nodeman.UpgradePluginConfig{
			{
				InnerIPs: []string{collector.InnerIP},
			},
		},
	}
	upgradeReq.GlobalParams.Package = pkg
	upgradeReq.GlobalParams.Plugin = proc
	upgradeReq.GlobalParams.Control = porcInfo
	upgradeReq.GlobalParams.UpgradeType = "APPEND"
	upgradeReq.GlobalParams.Configs = []nodeman.UpgradePluginConfig{
		{
			InnerIPs: []string{collector.InnerIP},
			Content:  base64.RawStdEncoding.EncodeToString(pluginConfig),
		},
	}

	return &upgradeReq, nil
}

func (lgc *Logics) buildNetdevicebeatConfigFile(collector *metadata.Netcollector) ([]byte, error) {
	customs, err := lgc.findCustom()
	if err != nil {
		blog.Errorf("[NetDevice][buildNetdevicebeatConfigFile] findCustom failed: %v", err)
	}

	config := NetDeviceConfig{
		DataID:      1014,
		CloudID:     collector.CloudID,
		ScanRange:   collector.Config.ScanRange,
		PingTimeout: 5,
		PingRetry:   3,
		Worker:      10,
		Period:      collector.Config.Period,
		Snmp: SnmpConfig{
			Port:      161,
			Community: collector.Config.Community,
			Version:   Version2c,
			Timeout:   time.Duration(10) * time.Second,
			Retries:   3,
			MaxOids:   10,
		},
		Customs: customs,
		Report: Report{
			Debug: true,
		},
	}

	return yaml.Marshal(&config)
}

func (lgc *Logics) findCustom() ([]Custom, error) {
	customs := []Custom{}
	propertys := []metadata.NetcollectProperty{}
	if err := lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectProperty, nil, nil, &propertys, "", 0, 0); err != nil {
		blog.Errorf("[NetDevice] failed to query the propertys, error info %s", err.Error())
		return nil, err
	}
	devices := []metadata.NetcollectDevice{}
	if err := lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectDevice, nil, nil, &devices, "", 0, 0); err != nil {
		blog.Errorf("[NetDevice] failed to query the devices, error info %s", err.Error())
		return nil, err
	}

	deviceMap := map[int64]metadata.NetcollectDevice{}
	for _, device := range devices {
		deviceMap[device.DeviceID] = device
	}

	for _, property := range propertys {
		device := deviceMap[property.DeviceID]

		custom := Custom{}
		custom.BkVendor = device.BkVendor
		custom.DeviceModel = device.DeviceModel
		custom.ObjectID = device.ObjectID
		custom.Method = property.Action
		custom.OID = property.OID
		custom.Period = property.Period
		custom.PropertyID = property.PropertyID
	}

	return customs, nil
}

type NetDeviceConfig struct {
	DataID      int64      `yaml:"dataid,omitempty"`
	CloudID     int64      `yaml:"bk_cloud_id,omitempty"`
	ScanRange   []string   `yaml:"scan_range,omitempty"`
	Snmp        SnmpConfig `yaml:"snmp,omitempty"`
	PingTimeout int        `yaml:"ping_timeout,omitempty"`
	PingRetry   int        `yaml:"ping_retry,omitempty"`
	Worker      int        `yaml:"worker,omitempty"`
	Period      string     `yaml:"period,omitempty"`
	Customs     []Custom   `yaml:"customs,omitempty"`
	Report      Report     `yaml:"report,omitempty"`
}

type Report struct {
	Debug bool `yaml:"debug,omitempty"`
}

type Custom struct {
	DeviceModel string `yaml:"device_model,omitempty,omitempty"`
	ObjectID    string `yaml:"bk_obj_id,omitempty"`
	BkVendor    string `yaml:"bk_vendor,omitempty,omitempty"`
	PropertyID  string `json:"bk_property_id" bson:"bk_property_id,omitempty,omitempty"`

	Method string `yaml:"method,omitempty"`
	Period string `yaml:"period,omitempty"`
	OID    string `yaml:"oid,omitempty"`
}

var snmpDefault = SnmpConfig{
	Port:      161,
	Community: "public",
	Version:   Version2c,
	Timeout:   time.Duration(10) * time.Second,
	Retries:   3,
	MaxOids:   10,
}

// SnmpConfig snmp config
type SnmpConfig struct {
	// Target is an ipv4 address
	Target string `yaml:"target,omitempty"`
	// Port is a udp port
	Port int `yaml:"port,omitempty"`
	// Community is an SNMP Community string
	Community string `yaml:"community,omitempty"`
	// Version is an SNMP Version
	Version Version `yaml:"version,omitempty"`
	// Timeout is the timeout for the SNMP Query
	Timeout time.Duration `yaml:"timeout,omitempty"`
	// Set the number of retries to attempt within timeout.
	Retries int `yaml:"retries,omitempty"`
	// MaxOids is the maximum number of oids allowed in a Get()
	// (default: 10)
	MaxOids int `yaml:"max_oids,omitempty"`
}

type Version string

// SnmpVersion constanst
const (
	Version1  Version = "SNMPv1"
	Version2c Version = "SNMPv2c"
	Version3  Version = "SNMPv3"
)
