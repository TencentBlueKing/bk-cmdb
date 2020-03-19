package cloudsync

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	ccom "configcenter/src/scene_server/cloud_server/common"
	"configcenter/src/scene_server/cloud_server/logics"
	"configcenter/src/storage/dal"
)

// 云主机同步器
type HostSyncor struct {
	logics *logics.Logics
	db     dal.DB
}

// 创建云主机同步器
func NewHostSyncor(logics *logics.Logics, db dal.DB) *HostSyncor {
	return &HostSyncor{logics, db}
}

// 同步云主机
func (h *HostSyncor) Sync(task *metadata.CloudSyncTask) error {
	startTime := time.Now()
	// 根据账号id获取账号详情
	accountConf, err := h.logics.GetCloudAccountConf(task.AccountID)
	if err != nil {
		blog.Errorf("GetCloudAccountConf err, taskid:%d, err:%s", task.TaskID, err.Error())
		return err
	}

	// 根据任务详情和账号信息获取要同步的云主机资源
	hostResource, err := h.getCloudHostResource(task, accountConf)
	if err != nil {
		blog.Errorf("getCloudHostResource err, taskid:%d, err:%s", task.TaskID, err.Error())
		return err
	}
	if len(hostResource.HostResource) == 0 {
		blog.Infof("HostResource is empty, taskid:%d", task.TaskID)
		return nil
	}
	hostResource.TaskID = task.TaskID
	blog.V(4).Infof("taskid:%d, vpc count:%d", task.TaskID, len(hostResource.HostResource))

	// 查询vpc对应的云区域，没有则创建,并更新云主机资源信息里的云区域id
	h.addCLoudId(accountConf, hostResource)
	if err != nil {
		blog.Errorf("addCLoudId err, taskid:%d, err:%s", task.TaskID, err.Error())
		return err
	}

	// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
	diffHosts, err := h.getDiffHosts(hostResource)
	if err != nil {
		blog.Errorf("getDiffHosts err, taskid:%d, err:%s", task.TaskID, err.Error())
		return err
	}

	// 没差异则结束
	if len(diffHosts) == 0 {
		blog.V(3).Infof("no diff hosts for taskid:%d", task.TaskID)
		return nil
	}

	// 有差异的更新任务同步状态为同步中
	err = h.updateTaskState(task.TaskID, metadata.CloudSyncInProgress)
	if err != nil {
		blog.Errorf("updateTaskState err, taskid:%d, err:%s", task.TaskID, err.Error())
		return err
	}

	// todo 后面几个表操作放在同一个事务里
	// 同步有差异的主机数据
	syncResult, err := h.syncDiffHosts(diffHosts)
	if err != nil {
		blog.Errorf("syncDiffHosts err, taskid:%d, err:%s", task.TaskID, err.Error())
		return err
	}

	// 增加任务同步历史记录
	_, err = h.addSyncHistory(syncResult, task.TaskID, startTime)
	if err != nil {
		blog.Errorf("addSyncHistory err, taskid:%d, err:%s", task.TaskID, err.Error())
		return err
	}

	// 完成后更新任务同步状态为成功或失败
	syncState := metadata.CloudSyncSuccess
	if syncResult.FailInfo.Count > 0 {
		syncState = metadata.CloudSyncFail
	}
	err = h.updateTaskState(task.TaskID, syncState)
	if err != nil {
		blog.Errorf("updateTaskState err, taskid:%d, err:%s", task.TaskID, err.Error())
		return err
	}

	costTime := time.Since(startTime) / time.Second
	blog.V(3).Infof("finish SyncCloudHost, costTime:%ds, syncResult.Detail:%#v, syncResult.FailInfo:%#v",
		costTime, syncResult.Detail, syncResult.FailInfo)

	return nil
}

// 根据任务详情和账号信息获取要同步的云主机资源
func (h *HostSyncor) getCloudHostResource(task *metadata.CloudSyncTask, accountConf *metadata.CloudAccountConf) (*metadata.CloudHostResource, error) {
	return h.logics.GetCloudHostResource(*accountConf, task.SyncVpcs)
}

// 查询vpc对应的云区域，没有则创建,并更新云主机资源信息里的云区域id
func (h *HostSyncor) addCLoudId(accountConf *metadata.CloudAccountConf, hostResource *metadata.CloudHostResource) (*metadata.CloudHostResource, error) {
	for _, hostRes := range hostResource.HostResource {
		cloudID, err := h.getCloudId(hostRes.Vpc.VpcID)
		if err != nil {
			continue
		}
		// 没有则创建
		if cloudID == 0 {
			cloudArea, err := h.createCloudArea(hostRes.Vpc, accountConf)
			if err != nil {
				continue
			}
			cloudID = cloudArea.CloudID
		}
		hostRes.CloudID = cloudID
	}
	return nil, nil
}

// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
func (h *HostSyncor) getDiffHosts(hostResource *metadata.CloudHostResource) (map[string][]*metadata.CloudHost, error) {
	hosts := make([]*metadata.CloudHost, 0)
	for _, hostRes := range hostResource.HostResource {
		for _, host := range hostRes.Instances {
			hosts = append(hosts, &metadata.CloudHost{
				Instance: *host,
				CloudID:  hostRes.CloudID,
				SyncDir:  hostRes.Vpc.SyncDir,
			})
		}
	}
	instanceIds := make([]string, 0)
	for _, h := range hosts {
		instanceIds = append(instanceIds, h.InstanceId)
	}
	blog.V(4).Infof("taskid:%d, host instanceIds:%#v", hostResource.TaskID, instanceIds)

	localHosts, err := h.getLocalHosts(instanceIds)
	if err != nil {
		return nil, err
	}
	blog.V(4).Infof("taskid:%d, len(localHosts):%d", hostResource.TaskID, len(localHosts))
	localIdHostsMap := make(map[string]*metadata.CloudHost)
	for _, h := range localHosts {
		localIdHostsMap[h.InstanceId] = h
	}

	diffHosts := make(map[string][]*metadata.CloudHost)
	for _, h := range hosts {
		if _, ok := localIdHostsMap[h.InstanceId]; ok {
			lh := localIdHostsMap[h.InstanceId]
			// 判断云主机和本地主机是否有差异，有则需要更新
			if ccom.CovertInstState(h.InstanceState) != lh.InstanceState || h.PublicIp != lh.PublicIp ||
				h.PrivateIp != lh.PrivateIp || h.CloudID != lh.CloudID {
				diffHosts["update"] = append(diffHosts["update"], h)
			}
		} else {
			diffHosts["add"] = append(diffHosts["add"], h)
		}
	}
	return diffHosts, nil
}

// 同步有差异的主机数据
func (h *HostSyncor) syncDiffHosts(diffhosts map[string][]*metadata.CloudHost) (*metadata.SyncResult, error) {
	syncResult := new(metadata.SyncResult)
	var result *metadata.SyncResult
	var err error
	for op, hosts := range diffhosts {
		switch op {
		case "add":
			result, err = h.addHosts(hosts)
			if err != nil {
				return nil, err
			}
			syncResult.Detail.NewAdd = result.SuccessInfo
		case "update":
			result, err = h.updateHosts(hosts)
			if err != nil {
				return nil, err
			}
			syncResult.Detail.Update = result.SuccessInfo
		default:
			blog.Errorf("syncDiffHosts op:%s is invalid", op)
			return nil, fmt.Errorf("syncDiffHosts op:%s is invalid", op)
		}
		syncResult.SuccessInfo.Count += result.SuccessInfo.Count
		syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, result.SuccessInfo.IPs...)
		syncResult.FailInfo.Count += result.FailInfo.Count
		for ip, errinfo := range result.FailInfo.IPError {
			syncResult.FailInfo.IPError[ip] = errinfo
		}
	}

	return syncResult, nil
}

// 增加任务同步历史记录
func (h *HostSyncor) addSyncHistory(syncResult *metadata.SyncResult, taskid int64, startTime time.Time) (*metadata.SyncHistory, error) {
	id, err := h.db.NextSequence(context.Background(), common.BKTableNameCloudSyncHistory)
	if nil != err {
		blog.Errorf("createCloudArea failed, generate id failed, err: %s", err.Error())
		return nil, err
	}
	syncStatus := metadata.CloudSyncSuccess
	costTime, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", float64(time.Since(startTime)/time.Millisecond)/1000.0), 64)
	statusDesc := metadata.SyncStatusDesc{CostTime: costTime}
	if syncResult.FailInfo.Count > 0 {
		syncStatus = metadata.CloudSyncFail
		for _, errinfo := range syncResult.FailInfo.IPError {
			statusDesc.ErrorInfo = errinfo
			break
		}
	}

	syncHistory := metadata.SyncHistory{
		HistoryID:         int64(id),
		TaskID:            taskid,
		SyncStatus:        syncStatus,
		StatusDescription: statusDesc,
		OwnerID:           fmt.Sprintf("%d", common.BKDefaultSupplierID),
		Detail:            syncResult.Detail,
		CreateTime:        metadata.Now(),
	}
	if err := h.db.Table(common.BKTableNameCloudSyncHistory).Insert(context.Background(), syncHistory); err != nil {
		if err != nil {
			blog.Errorf("addSyncHistory insert err:%v", err.Error())
			return nil, err
		}
	}
	return &syncHistory, nil
}

// 根据账号vpcID获取云区域ID，没有则创建
func (h *HostSyncor) getCloudId(vpcID string) (int64, error) {
	cond := mapstr.MapStr{common.BKVpcID: vpcID}
	result := make([]*metadata.CloudArea, 0)
	err := h.db.Table(common.BKTableNameBasePlat).Find(cond).All(context.Background(), &result)
	if err != nil {
		blog.Errorf("getCloudId err:%v", err.Error())
		return int64(0), err
	}
	if len(result) == 0 {
		return int64(0), nil
	}
	return result[0].CloudID, nil
}

// 创建vpc对应的云区域
func (h *HostSyncor) createCloudArea(vpc *metadata.VpcSyncInfo, accountConf *metadata.CloudAccountConf) (*metadata.CloudArea, error) {
	id, err := h.db.NextSequence(context.Background(), common.BKTableNameBasePlat)
	if nil != err {
		blog.Errorf("createCloudArea failed, generate id failed, err: %s", err.Error())
		return nil, err
	}
	ts := metadata.Now()
	cloudArea := metadata.CloudArea{
		CloudID:     int64(id),
		CloudName:   fmt.Sprintf("%d_%s", accountConf.AccountID, vpc.VpcID),
		Status:      1,
		CloudVendor: accountConf.VendorName,
		OwnerID:     fmt.Sprintf("%d", common.BKDefaultSupplierID),
		VpcID:       vpc.VpcID,
		VpcName:     vpc.VpcName,
		Region:      vpc.Region,
		AccountID:   accountConf.AccountID,
		Creator:     "cc_system",
		CreateTime:  ts,
		LastEditor:  "cc_system",
		LastTime:    ts,
	}
	if err := h.db.Table(common.BKTableNameBasePlat).Insert(context.Background(), cloudArea); err != nil {
		if err != nil {
			blog.Errorf("createCloudArea insert err:%v", err.Error())
			return nil, err
		}
	}
	return &cloudArea, nil
}

// 获取本地数据库中的主机信息
func (h *HostSyncor) getLocalHosts(instanceIds []string) ([]*metadata.CloudHost, error) {
	cond := mapstr.MapStr{common.BKCloudInstIDField: mapstr.MapStr{common.BKDBIN: instanceIds}}
	result := make([]*metadata.CloudHost, 0)
	err := h.db.Table(common.BKTableNameBaseHost).Find(cond).All(context.Background(), &result)
	if err != nil {
		blog.Errorf("getLocalHosts err:%v", err.Error())
		return nil, err
	}
	return result, nil
}

// 添加云主机到本地数据库和主机资源池目录对应关系
func (h *HostSyncor) addHosts(hosts []*metadata.CloudHost) (*metadata.SyncResult, error) {
	syncResult := new(metadata.SyncResult)
	for _, host := range hosts {
		id, err := h.db.NextSequence(context.Background(), common.BKTableNameBaseHost)
		if nil != err {
			blog.Errorf("addHosts failed, generate id failed, err: %s", err.Error())
			syncResult.FailInfo.Count++
			syncResult.FailInfo.IPError[host.PrivateIp] = err.Error()
			continue
		}
		host.HostID = int64(id)
		ts := metadata.Now()
		hostSyncInfo := metadata.HostSyncInfo{
			HostID:        host.HostID,
			CloudID:       host.CloudID,
			InstanceId:    host.InstanceId,
			InstanceName:  host.InstanceName,
			PrivateIp:     host.PrivateIp,
			PublicIp:      host.PublicIp,
			InstanceState: host.InstanceState,
			OsName:        host.OsName,
			CreateTime:    ts,
			LastTime:      ts,
		}
		if err := h.db.Table(common.BKTableNameBaseHost).Insert(context.Background(), hostSyncInfo); err != nil {
			blog.Errorf("addHosts insert err:%v", err.Error())
			syncResult.FailInfo.Count++
			syncResult.FailInfo.IPError[host.PrivateIp] = err.Error()
			continue
		} else {
			syncResult.SuccessInfo.Count++
			syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, host.PrivateIp)
		}

		// 获取资源池目录信息
		cond := mapstr.MapStr{common.BKModuleIDField: host.SyncDir}
		result := make([]*metadata.ModuleInst, 0)
		err = h.db.Table(common.BKTableNameBaseModule).Find(cond).All(context.Background(), &result)
		if err != nil {
			blog.Errorf("resource dir find err:%v", err.Error())
			return nil, err
		}
		if len(result) == 0 {
			blog.Errorf("resource dir %d is not exist", host.SyncDir)
			return nil, fmt.Errorf("resource dir %d is not exist", host.SyncDir)
		}

		// 增加主机资源池目录对应关系
		module := result[0]
		modulehost := metadata.ModuleHost{
			AppID:    module.BizID,
			HostID:   host.HostID,
			ModuleID: module.ModuleID,
			SetID:    module.ParentID,
			OwnerID:  fmt.Sprintf("%d", common.BKDefaultSupplierID),
		}
		if err := h.db.Table(common.BKTableNameModuleHostConfig).Insert(context.Background(), modulehost); err != nil {
			blog.Errorf("add module host relationship err:%s", err.Error())
			return nil, fmt.Errorf("add module host relationship err:%s", err.Error())
		}
	}

	return syncResult, nil
}

// 更新云主机到本地数据库
func (h *HostSyncor) updateHosts(hosts []*metadata.CloudHost) (*metadata.SyncResult, error) {
	syncResult := new(metadata.SyncResult)
	for _, host := range hosts {
		cond := mapstr.MapStr{common.BKCloudInstIDField: host.InstanceId}
		updateInfo := mapstr.MapStr{
			common.BKCloudIDField:         host.CloudID,
			common.BKHostInnerIPField:     host.PrivateIp,
			common.BKHostOuterIPField:     host.PublicIp,
			common.BKCloudHostStatusField: host.InstanceState,
			common.BKHostNameField:        host.InstanceName,
			common.LastTimeField:          metadata.Now(),
		}
		if err := h.db.Table(common.BKTableNameBaseHost).Update(context.Background(), cond, updateInfo); err != nil {
			blog.Errorf("updateHosts update err:%v", err.Error())
			syncResult.FailInfo.Count++
			syncResult.FailInfo.IPError[host.PrivateIp] = err.Error()
			continue
		} else {
			syncResult.SuccessInfo.Count++
			syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, host.PrivateIp)
		}
	}
	return syncResult, nil
	return nil, nil
}

// 更新任务同步状态
func (h *HostSyncor) updateTaskState(taskid int64, status string) error {
	option := mapstr.MapStr{common.BKCloudSyncStatus: status}
	if status == metadata.CloudSyncSuccess || status == metadata.CloudSyncFail {
		option.Set(common.BKCloudLastSyncTime, metadata.Now())
	}

	if err := h.logics.CoreAPI.CoreService().Cloud().UpdateSyncTask(context.Background(), header, taskid, option); err != nil {
		blog.Errorf("UpdateSyncTask failed, taskid: %v, err: %s", taskid, err.Error())
		return err
	}

	return nil
}
