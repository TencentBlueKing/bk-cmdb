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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/thirdparty/esbserver/nodeman"

	"gopkg.in/yaml.v2"
)

// Netdevicebeat netdevicebeat collector name
const Netdevicebeat = "netdevicebeat"

func (lgc *Logics) SearchCollector(header http.Header, cond metadata.ParamNetcollectorSearch) (int64, []metadata.Netcollector, error) {
	collectors := make([]metadata.Netcollector, 0)

	// fetch package info
	packageResp, err := lgc.ESB.NodemanSrv().SearchPackage(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] SearchPackage by %s failed: %v", Netdevicebeat, err)
		return 0, nil, err
	}
	if !packageResp.Result {
		blog.Errorf("[NetDevice][SearchCollector] SearchPackage by %s failed: %+v", Netdevicebeat, packageResp)
		return 0, nil, fmt.Errorf("search plugin host from nodeman failed: %s", packageResp.Message)
	}
	var pkg nodeman.PluginPackage
	if len(packageResp.Data) > 0 { // nodeman team ensure that they will sort by id in descending order
		pkg = packageResp.Data[0]
	}

	// fetch hosts
	pluginHostResp, err := lgc.ESB.NodemanSrv().SearchPluginHost(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] SearchPluginHost failed: %v", err)
		return 0, nil, err
	}
	if !pluginHostResp.Result {
		blog.Errorf("[NetDevice][SearchCollector] SearchPluginHost by %s failed: %+v", Netdevicebeat, packageResp)
		return 0, nil, fmt.Errorf("search plugin host from nodeman by %s failed: %s", Netdevicebeat, pluginHostResp.Message)
	}

	// build collectors
	cloudIDs := make([]int64, 0)
	ips := make([]string, 0)
	for _, pHost := range pluginHostResp.Data {
		cloudIDs = append(cloudIDs, pHost.Host.BkCloudID)
		ips = append(ips, pHost.Host.InnerIP)

		collectorStatus := metadata.CollectorStatusAbnormal
		if strings.ToUpper(pHost.Status) == "RUNNING" {
			collectorStatus = metadata.CollectorStatusNormal
		}

		collector := metadata.Netcollector{
			BizID:   pHost.Host.BkBizID,
			CloudID: pHost.Host.BkCloudID,
			InnerIP: pHost.Host.InnerIP,
			Version: pHost.Version,
			Status: metadata.NetcollectorStatus{
				CollectorStatus: collectorStatus,
			},
			LatestVersion: pkg.Version,
		}
		collectors = append(collectors, collector)
	}

	cloudCond := map[string]interface{}{
		common.BKCloudIDField: map[string]interface{}{
			common.BKDBIN: cloudIDs,
		},
	}
	cloudMap, err := lgc.findInstMap(header, common.BKInnerObjIDPlat, &metadata.QueryCondition{Condition: cloudCond})
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] find clouds by %+v failed: %v", cloudCond, err)
		return 0, nil, err
	}

	cloudCond[common.BKHostInnerIPField] = map[string]interface{}{
		common.BKDBIN: ips,
	}
	collectorMap, err := lgc.findCollectorMap(cloudCond)
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] get collector config by %+v failed %v", cloudCond, err)
		return 0, nil, err
	}

	// fill config field from our db
	for index := range collectors {
		collector := &collectors[index]

		if clodInst, ok := cloudMap[collector.CloudID]; ok {
			cloudName, err := clodInst.String(common.BKCloudNameField)
			if err != nil {
				blog.Errorf("[NetDevice][SearchCollector] bk_cloud_name field invalid: %v, inst: %+v", err, clodInst)
			}
			collector.CloudName = cloudName
		}

		key := collectorMapKey(collector.CloudID, collector.InnerIP)
		existsOne, ok := collectorMap[key]
		if !ok {
			blog.Warnf("[NetDevice][SearchCollector] get collector config for %s failed", key)
		}
		collector.Config = existsOne.Config
		collector.ReportTotal = existsOne.ReportTotal
		collector.TaskID = existsOne.TaskID
		collector.DeployTime = existsOne.DeployTime

		if existsOne.Status.ConfigStatus == "" || existsOne.Status.ConfigStatus == metadata.CollectorConfigStatusPending {
			var taskStatus string
			if existsOne.TaskID <= 0 {
				taskStatus = metadata.CollectorConfigStatusAbnormal
			} else {
				taskStatus, err = lgc.queryCollectTask(header, collector.BizID, existsOne.TaskID)
				if err != nil {
					blog.Warnf("[NetDevice][SearchCollector] queryNodemanTask by BizID [%v], TaskID [%v], failed: %v", collector.BizID, existsOne.TaskID, err)
				}
			}
			existsOne.Status.ConfigStatus = taskStatus
			if taskStatus == metadata.CollectorConfigStatusNormal || taskStatus == metadata.CollectorConfigStatusAbnormal {
				if err = lgc.saveCollectTask(collector, existsOne.TaskID, taskStatus); err != nil {
					blog.Warnf("[NetDevice][SearchCollector] saveCollectTask for %+v failed: %v", collector, err)
				}
			}
		}
		if existsOne.Status.ReportStatus == "" {
			existsOne.Status.ReportStatus = metadata.CollectorReportStatusAbnormal
		}
		collector.Status.ConfigStatus = existsOne.Status.ConfigStatus
		collector.Status.ReportStatus = existsOne.Status.ReportStatus
	}
	return int64(len(collectors)), collectors, nil
}

func (lgc *Logics) queryCollectTask(header http.Header, bizID int64, taskID int64) (string, error) {
	resp, err := lgc.ESB.NodemanSrv().SearchTask(context.Background(), header, bizID, taskID)
	if err != nil {
		return metadata.CollectorConfigStatusPending, err
	}
	if !resp.Result {
		blog.Errorf("[NetDevice][queryNodemanTask] failed: %+v", resp.Message)
		return metadata.CollectorConfigStatusPending, fmt.Errorf(resp.Message)
	}

	for _, host := range resp.Data.Hosts {
		switch host.Status {
		case "FAILED":
			return metadata.CollectorConfigStatusAbnormal, nil
		case "SUCCESS":
			return metadata.CollectorConfigStatusNormal, nil
		case "QUEUE", "RUNNING":
			return metadata.CollectorConfigStatusPending, nil
		default:
			return metadata.CollectorConfigStatusPending, nil
		}
	}
	return metadata.CollectorConfigStatusPending, nil
}

func collectorMapKey(cloudID int64, ip string) string {
	return fmt.Sprintf("%d:%s", cloudID, ip)
}

func (lgc *Logics) findCollectorMap(cond interface{}) (map[string]metadata.Netcollector, error) {
	collectors := make([]metadata.Netcollector, 0)
	err := lgc.db.Table(common.BKTableNameNetcollectConfig).Find(cond).All(lgc.ctx, &collectors)
	if err != nil {
		return nil, err
	}
	collectorMap := map[string]metadata.Netcollector{}
	for index := range collectors {
		key := collectorMapKey(collectors[index].CloudID, collectors[index].InnerIP)
		collectorMap[key] = collectors[index]
	}
	return collectorMap, nil
}

func (lgc *Logics) UpdateCollector(header http.Header, config metadata.Netcollector) error {
	filter := map[string]interface{}{
		common.BKCloudIDField:     config.CloudID,
		common.BKHostInnerIPField: config.InnerIP,
	}

	count, err := lgc.db.Table(common.BKTableNameNetcollectConfig).Find(filter).Count(lgc.ctx)
	if err != nil {
		blog.Errorf("[UpdateCollector] count by %+v error: %v", filter, err)
		return err
	}
	if count > 0 {
		err = lgc.db.Table(common.BKTableNameNetcollectConfig).Update(lgc.ctx, filter, config)
		if err != nil {
			blog.Errorf("[UpdateCollector] UpdateByCondition by %+v to %+v error: %v", filter, config, err)
			return err
		}
		return lgc.DiscoverNetDevice(header, []metadata.Netcollector{config})
	}

	err = lgc.db.Table(common.BKTableNameNetcollectConfig).Insert(lgc.ctx, config)
	if err != nil {
		blog.Errorf("[UpdateCollector] Insert %+v error: %v", config, err)
		return err
	}

	return lgc.DiscoverNetDevice(header, []metadata.Netcollector{config})
}

func (lgc *Logics) DiscoverNetDevice(header http.Header, configs []metadata.Netcollector) error {

	// fetch global_params
	pkgResp, err := lgc.ESB.NodemanSrv().SearchPackage(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchPackage by %v failed %v", Netdevicebeat, err)
		return err
	}
	if !pkgResp.Result {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchPackage by %s, failed %s", Netdevicebeat, pkgResp.Message)
		return fmt.Errorf("search plugin host from nodeman failed: %s", pkgResp.Message)
	}
	var pkg nodeman.PluginPackage
	if len(pkgResp.Data) > 0 {
		pkg = pkgResp.Data[0]
	}

	procResp, err := lgc.ESB.NodemanSrv().SearchProcess(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchProcess by %s failed %v", Netdevicebeat, err)
		return err
	}
	if !procResp.Result {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchProcess by %s failed %+v", Netdevicebeat, pkgResp)
		return fmt.Errorf("search plugin host from nodeman failed: %s", procResp.Message)
	}
	procInfoResp, err := lgc.ESB.NodemanSrv().SearchProcessInfo(context.Background(), header, Netdevicebeat)
	if err != nil {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchProcess by %v failed %v", Netdevicebeat, err)
		return err
	}
	if !procInfoResp.Result {
		blog.Errorf("[NetDevice][DiscoverNetDevice] SearchProcess by %v failed %+v", Netdevicebeat, pkgResp)
		return fmt.Errorf("search plugin host from nodeman failed: %s", procInfoResp.Message)
	}

	cloudIDs := make([]int64, 0)
	ips := make([]string, 0)
	for _, config := range configs {
		cloudIDs = append(cloudIDs, config.CloudID)
		ips = append(ips, config.InnerIP)
	}
	cloudFilter := map[string]interface{}{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: ips,
		},
		common.BKCloudIDField: map[string]interface{}{
			common.BKDBIN: cloudIDs,
		},
	}
	collectorMap, err := lgc.findCollectorMap(cloudFilter)
	if err != nil {
		blog.Errorf("[NetDevice][SearchCollector] get collector config by %+v failed, %v", cloudFilter, err)
		return err
	}

	for _, config := range configs {
		key := collectorMapKey(config.CloudID, config.InnerIP)
		collector, ok := collectorMap[key]
		if !ok {
			blog.Errorf("[NetDevice][DiscoverNetDevice] get collector config for %s failed", key)
			continue
		}

		upgradeReq, err := lgc.buildUpgradePluginRequest(&collector, util.GetUser(header), &pkg, &procResp.Data, &procInfoResp.Data)
		if err != nil {
			blog.Errorf("[NetDevice][DiscoverNetDevice] buildUpgradePluginRequest %s failed, %v", key, err)
		}
		blog.InfoJSON("[NetDevice][DiscoverNetDevice] UpgradePlugin request %s", upgradeReq)
		upgradeResp, err := lgc.ESB.NodemanSrv().UpgradePlugin(context.Background(), header, strconv.FormatInt(collector.BizID, 10), upgradeReq)
		if err != nil {
			blog.Errorf("[NetDevice][DiscoverNetDevice] UpgradePlugin %s failed", key)
			continue
		}
		if !upgradeResp.Result {
			blog.Errorf("NetDevice][DiscoverNetDevice] search plugin host from nodeman failed: %+v", upgradeResp)
			continue
		}

		blog.V(3).Infof("[NetDevice][DiscoverNetDevice] UpgradePlugin response %+v ", upgradeResp)
		if err := lgc.saveCollectTask(&collector, upgradeResp.Data.ID, metadata.CollectorConfigStatusPending); err != nil {
			blog.Errorf("[NetDevice][DiscoverNetDevice] saveCollectTask %s failed, %v", key, err)
		}
	}

	return err
}

func (lgc *Logics) saveCollectTask(collector *metadata.Netcollector, taskID int64, status string) error {
	data := map[string]interface{}{
		"task_id":              taskID,
		"status.config_status": status,
	}

	filter := map[string]interface{}{
		common.BKCloudIDField: collector.CloudID,
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: collector.InnerIP,
		},
	}

	return lgc.db.Table(common.BKTableNameNetcollectConfig).Update(lgc.ctx, filter, data)
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
		Hosts: []nodeman.UpgradePluginHostField{
			{
				InnerIPs: collector.InnerIP,
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
			Content:  base64.StdEncoding.EncodeToString(pluginConfig),
		},
	}

	return &upgradeReq, nil
}

func (lgc *Logics) buildNetdevicebeatConfigFile(collector *metadata.Netcollector) ([]byte, error) {
	customs, err := lgc.findCustom()
	if err != nil {
		blog.Errorf("[NetDevice][buildNetdevicebeatConfigFile] findCustom for %+v failed: %v", collector, err)
		return []byte(""), err
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
			Timeout:   10,
			Retries:   3,
			MaxOids:   10,
		},
		Customs: customs,
		Report: Report{
			Debug: true,
		},
	}

	configContent := map[string]interface{}{
		"netdevicebeat": config,
		"output":        map[string]interface{}{"gse": map[string]string{"endpoint": "/var/run/ipc.state.report"}},
		"path": map[string]interface{}{
			"data": "/var/lib/gse",
			"logs": "/var/log/gse",
			"pid":  "/var/run/gse",
		},
		"logging": map[string]interface{}{
			"to_files": true,
		},
	}

	return yaml.Marshal(&configContent)
}

func (lgc *Logics) findCustom() ([]Custom, error) {
	customs := make([]Custom, 0)
	properties := make([]metadata.NetcollectProperty, 0)
	if err := lgc.db.Table(common.BKTableNameNetcollectProperty).Find(nil).All(lgc.ctx, &properties); err != nil {
		blog.Errorf("[NetDevice] failed to query the propertys, error info %v", err)
		return nil, err
	}
	devices := make([]metadata.NetcollectDevice, 0)
	if err := lgc.db.Table(common.BKTableNameNetcollectDevice).Find(nil).All(lgc.ctx, &devices); err != nil {
		blog.Errorf("[NetDevice] failed to query the devices, error info %v", err)
		return nil, err
	}

	deviceMap := map[uint64]metadata.NetcollectDevice{}
	for _, device := range devices {
		deviceMap[device.DeviceID] = device
	}

	for _, property := range properties {
		device := deviceMap[property.DeviceID]

		custom := Custom{}
		custom.BkVendor = device.BkVendor
		custom.DeviceModel = device.DeviceModel
		custom.ObjectID = device.ObjectID
		custom.Method = property.Action
		custom.OID = property.OID
		custom.Period = property.Period
		custom.PropertyID = property.PropertyID
		customs = append(customs, custom)
	}

	return customs, nil
}

type NetDeviceConfig struct {
	DataID      int64      `yaml:"dataid,omitempty"`
	CloudID     int64      `yaml:"bk_cloud_id,omitempty"`
	OwnerID     string     `yaml:"bk_supplier_account"`
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
	DeviceModel string `yaml:"device_model,omitempty"`
	ObjectID    string `yaml:"bk_obj_id,omitempty"`
	BkVendor    string `yaml:"bk_vendor,omitempty"`
	PropertyID  string `json:"bk_property_id" bson:"bk_property_id,omitempty"`

	Method string `yaml:"method,omitempty"`
	Period string `yaml:"period,omitempty"`
	OID    string `yaml:"oid,omitempty"`
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
	Timeout int `yaml:"timeout,omitempty"`
	// Set the number of retries to attempt within timeout.
	Retries int `yaml:"retries,omitempty"`
	// MaxOids is the maximum number of oids allowed in a Get()
	// (default: 10)
	MaxOids int `yaml:"max_oids,omitempty"`
}

type Version string

// SnmpVersion constant
const (
	Version1  Version = "SNMPv1"
	Version2c Version = "SNMPv2c"
	Version3  Version = "SNMPv3"
)
