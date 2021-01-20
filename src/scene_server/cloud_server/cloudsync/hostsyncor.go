package cloudsync

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	ccom "configcenter/src/scene_server/cloud_server/common"
	"configcenter/src/scene_server/cloud_server/logics"
)

// 云主机同步器
type HostSyncor struct {
	logics *logics.Logics
	// readKit used for read operation
	readKit *rest.Kit
	// writeKit used for write operation
	writeKit *rest.Kit
}

// 创建云主机同步器
func NewHostSyncor(logics *logics.Logics) *HostSyncor {
	return &HostSyncor{
		logics: logics,
	}
}

// 同步云主机
func (h *HostSyncor) Sync(task *metadata.CloudSyncTask) error {
	defer func() {
		if err := recover(); err != nil {
			blog.Errorf("sync panic err:%#v, rid:%s, debug strace:%s", err, h.readKit.Rid, debug.Stack())
		}
	}()

	// 每次同步生成新的kit
	h.readKit = ccom.NewKit()
	// 将云同步任务的开发商ID作为写kit的开发商ID
	h.writeKit = ccom.NewWriteKit(task.OwnerID)
	// 让读写kit的requestID保持一致，以追踪同一个task的日志
	h.writeKit.Header.Set(common.BKHTTPCCRequestID, h.readKit.Header.Get(common.BKHTTPCCRequestID))

	startTime := time.Now()
	blog.Infof("start sync taskid:%d, rid:%s", task.TaskID, h.readKit.Rid)

	// 根据账号id获取账号详情
	accountConf, err := h.logics.GetCloudAccountConf(h.readKit, task.AccountID)
	if err != nil {
		blog.Errorf("GetCloudAccountConf fail, taskid:%d, err:%s, rid:%s", task.TaskID,
			err.Error(), h.readKit.Rid)
		return err
	}

	// 根据任务详情和账号信息获取要同步的云主机资源
	hostResource, err := h.getCloudHostResource(task, accountConf)
	if err != nil {
		blog.Errorf("getCloudHostResource fail, taskid:%d, err:%s, rid:%s", task.TaskID,
			err.Error(), h.readKit.Rid)
		return err
	}
	if len(hostResource.HostResource) == 0 && len(hostResource.DestroyedVpcs) == 0 {
		blog.Infof("hostResource is empty, taskid:%d, rid:%s", task.TaskID, h.readKit.Rid)
		return nil
	}

	blog.Infof(" taskid:%d, destroyed vpc count:%d, other vpc count:%d, rid:%s", task.TaskID, len(hostResource.DestroyedVpcs), len(hostResource.HostResource), h.readKit.Rid)

	syncResult := new(metadata.SyncResult)
	syncResult.FailInfo.IPError = make(map[string]string)

	txnErr := h.logics.CoreAPI.CoreService().Txn().AutoRunTxn(h.readKit.Ctx, h.readKit.Header, func() error {
		// 让writeKit的header含有同样的事务信息，以保证同一个事务里写操作后的数据能够被读到
		ccom.CopyHeaderTxnInfo(h.readKit.Header, h.writeKit.Header)
		if len(hostResource.DestroyedVpcs) > 0 {
			// 同步被销毁的VPC相关资源
			err = h.syncDestroyedVpcs(hostResource, syncResult)
			if err != nil {
				blog.Errorf("syncDestroyedVpcs fail, taskid:%d, err:%s, rid:%s", task.TaskID,
					err.Error(), h.readKit.Rid)
				return err
			}
		}

		// 查询vpc对应的云区域并更新云主机资源信息里的云区域id
		err = h.addCLoudId(accountConf, hostResource)
		if err != nil {
			blog.Errorf("addCLoudId fail, taskid:%d, err:%s, rid:%s", task.TaskID, err.Error(),
				h.readKit.Rid)
			return err
		}

		// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
		diffHosts, err := h.getDiffHosts(hostResource)
		if err != nil {
			blog.Errorf("getDiffHosts fail, taskid:%d, err:%s, rid:%s", task.TaskID, err.Error(), h.readKit.Rid)
			return err
		}

		// 没差异则结束
		if len(diffHosts) == 0 {
			blog.Infof("no diff hosts for taskid:%d, rid:%s", task.TaskID, h.readKit.Rid)
			if syncResult.SuccessInfo.Count == 0 {
				blog.Infof("no any hosts need sync for taskid:%d, rid:%s", task.TaskID,
					h.readKit.Rid)
				return nil
			}
		}

		// 有差异的更新任务同步状态为同步中
		err = h.updateTaskState(h.writeKit, task.TaskID, metadata.CloudSyncInProgress, nil)
		if err != nil {
			blog.Errorf("updateTaskState fail, taskid:%d, err:%s, rid:%s", task.TaskID,
				err.Error(), h.readKit.Rid)
			return err
		}

		// 同步有差异的主机数据
		err = h.syncDiffHosts(diffHosts, syncResult)
		if err != nil {
			blog.Errorf("syncDiffHosts fail, taskid:%d, err:%s, rid:%s", task.TaskID, err.Error(),
				h.readKit.Rid)
			return err
		}

		// 设置SyncResult的状态信息
		err = h.SetSyncResultStatus(syncResult, startTime)
		if err != nil {
			blog.Errorf("SetSyncResultStatus fail, taskid:%d, err:%s, rid:%s", task.TaskID,
				err.Error(), h.readKit.Rid)
			return err
		}

		// 增加任务同步历史记录
		_, err = h.addSyncHistory(syncResult, task.TaskID)
		if err != nil {
			blog.Errorf("addSyncHistory fail, taskid:%d, err:%s, rid:%s", task.TaskID,
				err.Error(), h.readKit.Rid)
			return err
		}

		// 完成后更新任务同步状态
		err = h.updateTaskState(h.writeKit, task.TaskID, syncResult.SyncStatus, &syncResult.StatusDescription)
		if err != nil {
			blog.Errorf("updateTaskState fail, taskid:%d, err:%s, rid:%s", task.TaskID,
				err.Error(), h.readKit.Rid)
			return err
		}

		blog.Infof("sync success, finish sync taskid:%s, costTime:%ds, Detail:%#v, FailInfo:%#v, rid:%s",
			task.TaskID, time.Since(startTime)/time.Second, syncResult.Detail, syncResult.FailInfo, h.readKit.Rid)
		return nil
	})

	// 事务结束，去掉readKit、writeKit中header的事务信息
	ccom.DelHeaderTxnInfo(h.readKit.Header)
	ccom.DelHeaderTxnInfo(h.writeKit.Header)

	if txnErr != nil {
		blog.Errorf("sync fail, taskid:%d, txnErr:%v, rid:%s", task.TaskID, txnErr, h.readKit.Rid)
		err := h.updateTaskState(h.writeKit, task.TaskID, metadata.CloudSyncFail, &metadata.SyncStatusDesc{ErrorInfo: txnErr.Error()})
		if err != nil {
			blog.Errorf("updateTaskState fail, taskid:%d, err:%v, rid:%s", task.TaskID, err,
				h.readKit.Rid)
		}
	} else {
		// 检查同步状态，如果为异常，则更新为正常, 以在没有主机差异要同步时也可以恢复某些问题导致的状态异常；已经是正常的情况下不需要更新
		opt := &metadata.SearchCloudOption{
			Condition: mapstr.MapStr{common.BKCloudSyncTaskID: task.TaskID},
		}

		ret, err := h.logics.CoreAPI.CoreService().Cloud().SearchSyncTask(h.readKit.Ctx, h.readKit.Header, opt)
		if err != nil {
			blog.Errorf("SearchSyncTask failed, taskid: %v, opt:%#v, err: %s, rid:%s",
				task.TaskID, opt, err.Error(), h.readKit.Rid)
			return err
		}
		if len(ret.Info) == 0 {
			blog.Errorf("SearchSyncTask failed, taskid %d is not found, opt:%#v, rid:%s",
				task.TaskID, opt, h.readKit.Rid)
			return fmt.Errorf("taskID %d is not found", task.TaskID)
		}
		if ret.Info[0].SyncStatus == metadata.CloudSyncFail {
			costTime, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", float64(time.Since(startTime)/time.Millisecond)/1000.0), 64)
			err := h.updateTaskState(h.writeKit, task.TaskID, metadata.CloudSyncSuccess, &metadata.SyncStatusDesc{CostTime: costTime})
			if err != nil {
				blog.Errorf("updateTaskState fail, taskid:%d, err:%v, rid:%s", task.TaskID, err,
					h.readKit.Rid)
			}
			blog.Infof("update taskid:%d status from fail to success, rid:%s", task.TaskID,
				h.readKit.Rid)
		}
		blog.Errorf("sync success, taskid:%d, rid:%s", task.TaskID, h.readKit.Rid)
	}

	blog.Infof("sync loop for taskid:%d is over, costTime:%ds, rid:%s", task.TaskID, time.Since(startTime)/time.Second,
		h.readKit.Rid)

	return nil
}

// 根据任务详情和账号信息获取要同步的云主机资源
func (h *HostSyncor) getCloudHostResource(task *metadata.CloudSyncTask, accountConf *metadata.CloudAccountConf) (*metadata.CloudHostResource, error) {
	hostResource, err := h.logics.GetCloudHostResource(h.readKit, *accountConf, task.SyncVpcs)
	if err != nil {
		return nil, err
	}
	hostResource.AccountConf = accountConf
	hostResource.TaskID = task.TaskID
	return hostResource, err
}

// 同步被销毁的VPC相关资源
func (h *HostSyncor) syncDestroyedVpcs(hostResource *metadata.CloudHostResource, syncResult *metadata.SyncResult) error {
	if len(hostResource.DestroyedVpcs) == 0 {
		return nil
	}
	cloudIDs := make([]int64, 0)
	for _, vpcInfo := range hostResource.DestroyedVpcs {
		cloudIDs = append(cloudIDs, vpcInfo.CloudID)
	}
	blog.Infof("Destroyed cloudIDs: %#v, rid:%s", cloudIDs, h.readKit.Rid)

	// 更新属于被销毁vpc下的主机信息，将内外网ip置空，状态置为已销毁
	condition := mapstr.MapStr{common.BKCloudIDField: map[string]interface{}{
		common.BKDBIN: cloudIDs,
	}}

	query := &metadata.QueryCondition{
		Condition: condition,
	}
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(h.readKit.Ctx, h.readKit.Header,
		common.BKInnerObjIDHost, query)
	if nil != err {
		blog.Errorf("syncDestroyedVpcs ReadInstance failed, error: %v query:%#v, rid:%s", err, query, h.readKit.Rid)
		return err
	}
	if false == res.Result {
		blog.Errorf("syncDestroyedVpcs failed, query:%#v, err code:%d, err msg:%s, rid:%s", query, res.Code,
			res.ErrMsg, h.readKit.Rid)
		return fmt.Errorf("%s", res.ErrMsg)
	}

	hostIDs := make([]int64, 0)
	for _, host := range res.Data.Info {
		hostID, _ := host.Int64(common.BKHostIDField)
		hostIDs = append(hostIDs, hostID)
	}
	sResult, err := h.deleteDestroyedHosts(hostIDs)
	if err != nil {
		blog.Errorf("syncDestroyedVpcs deleteDestroyedHosts fail, cloudIDs:%#v, err:%s, rid:%s", cloudIDs,
			err.Error(), h.readKit.Rid)
		return err
	}
	syncResult.Detail.Update.Count += sResult.SuccessInfo.Count
	syncResult.Detail.Update.IPs = append(syncResult.Detail.Update.IPs, sResult.SuccessInfo.IPs...)
	syncResult.SuccessInfo.Count += sResult.SuccessInfo.Count
	syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, sResult.SuccessInfo.IPs...)

	// 更新被销毁vpc对应的云区域状态为异常
	if err := h.updateDestroyedCloudArea(cloudIDs); err != nil {
		blog.Errorf("syncDestroyedVpcs updateDestroyedCloudArea fail, cloudIDs:%#v, err:%s, rid:%s", cloudIDs,
			err.Error(), h.readKit.Rid)
		return err
	}

	vpcs := make(map[string]bool)
	for _, vpcInfo := range hostResource.DestroyedVpcs {
		vpcs[vpcInfo.VpcID] = true
	}
	// 更新同步任务里的vpc状态为被销毁
	if err := h.updateDestroyedTaskVpc(hostResource.TaskID, vpcs); err != nil {
		blog.Errorf("syncDestroyedVpcs updateDestroyedTaskVpc fail, cloudIDs:%#v, err:%s, rid:%s", cloudIDs,
			err.Error(), h.readKit.Rid)
		return err
	}

	return nil
}

// 查询vpc对应的云区域
func (h *HostSyncor) addCLoudId(accountConf *metadata.CloudAccountConf, hostResource *metadata.CloudHostResource) error {
	for _, hostRes := range hostResource.HostResource {
		cloudID, err := h.getCloudId(hostRes.Vpc.VpcID)
		if err != nil {
			blog.Errorf("addCLoudId getCloudId err:%s, vpcID:%s, rid:%s", err.Error(), hostRes.Vpc.VpcID, h.readKit.Rid)
			return err
		}
		if cloudID == 0 {
			blog.Errorf("addCLoudId getCloudId err:%s, vpcID:%s, rid:%s",
				"the correspond cloudID for the vpc can't be found", hostRes.Vpc.VpcID, h.readKit.Rid)
			return fmt.Errorf("the correspond cloudID for the vpc %s can't be found,vpc name: %s", hostRes.Vpc.VpcID, hostRes.Vpc.VpcName)
		}
		hostRes.CloudID = cloudID
	}
	return nil
}

// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
func (h *HostSyncor) getDiffHosts(hostResource *metadata.CloudHostResource) (map[string][]*metadata.CloudHost, error) {
	// 云端的主机
	remoteHostsMap := make(map[string]*metadata.CloudHost)
	for _, hostRes := range hostResource.HostResource {
		for _, host := range hostRes.Instances {
			remoteHostsMap[host.InstanceId] = &metadata.CloudHost{
				Instance:   *host,
				CloudID:    hostRes.CloudID,
				VendorName: hostResource.AccountConf.VendorName,
				SyncDir:    hostRes.Vpc.SyncDir,
			}
		}
	}

	cloudIDs := make([]int64, 0)
	for _, vpcInfo := range hostResource.HostResource {
		cloudIDs = append(cloudIDs, vpcInfo.CloudID)
	}
	blog.V(4).Infof("taskid:%d, host cloudIDs:%#v, rid:%s", hostResource.TaskID, cloudIDs, h.readKit.Rid)

	// 本地已有的云主机
	localHosts, err := h.getLocalHosts(cloudIDs)
	if err != nil {
		return nil, err
	}
	blog.V(4).Infof("taskid:%d, len(localHosts):%d, rid:%s", hostResource.TaskID, len(localHosts), h.readKit.Rid)
	localIdHostsMap := make(map[string]*metadata.CloudHost)
	for _, h := range localHosts {
		localIdHostsMap[h.InstanceId] = h
	}

	// 有差异的主机
	diffHosts := make(map[string][]*metadata.CloudHost)
	// 本地需要同步新增和更新的主机
	for _, h := range remoteHostsMap {
		if _, ok := localIdHostsMap[h.InstanceId]; ok {
			lh := localIdHostsMap[h.InstanceId]
			// 判断云主机和本地主机是否有差异，有则需要更新
			// 如果是已经同步过的已销毁主机，就不用再同步了，保持被销毁主机的内外网ip始终为空
			if lh.InstanceState == common.BKCloudHostStatusDestroyed {
				continue
			}
			if h.InstanceState != lh.InstanceState || h.PublicIp != lh.PublicIp ||
				h.PrivateIp != lh.PrivateIp || h.CloudID != lh.CloudID {
				diffHosts["update"] = append(diffHosts["update"], h)
			}
		} else {
			diffHosts["add"] = append(diffHosts["add"], h)
		}
	}

	// 云端已销毁的主机,也就是本地需要删除的主机
	for id, h := range localIdHostsMap {
		// 已经同步过的已销毁主机，就不用再同步了，保持被销毁主机的内外网ip始终为空
		if h.InstanceState == common.BKCloudHostStatusDestroyed {
			continue
		}
		if _, ok := remoteHostsMap[id]; !ok {
			diffHosts["delete"] = append(diffHosts["delete"], h)
		}
	}

	return diffHosts, nil
}

// 同步有差异的主机数据
func (h *HostSyncor) syncDiffHosts(diffhosts map[string][]*metadata.CloudHost, syncResult *metadata.SyncResult) error {
	result := new(metadata.SyncResult)
	var err error
	for op, hosts := range diffhosts {
		switch op {
		case "add":
			result, err = h.addHosts(hosts)
			if err != nil {
				blog.Errorf("syncDiffHosts fail, err:%s, rid:%s", err.Error(), h.readKit.Rid)
				return err
			}
			syncResult.Detail.NewAdd.Count += result.SuccessInfo.Count
			syncResult.Detail.NewAdd.IPs = append(syncResult.Detail.NewAdd.IPs, result.SuccessInfo.IPs...)
		case "update":
			result, err = h.updateHosts(hosts)
			if err != nil {
				blog.Errorf("syncDiffHosts fail, err:%s, rid:%s", err.Error(), h.readKit.Rid)
				return err
			}
			syncResult.Detail.Update.Count += result.SuccessInfo.Count
			syncResult.Detail.Update.IPs = append(syncResult.Detail.Update.IPs, result.SuccessInfo.IPs...)
		case "delete":
			hostIDs := make([]int64, 0)
			for _, h := range hosts {
				hostIDs = append(hostIDs, h.HostID)
			}
			result, err = h.deleteDestroyedHosts(hostIDs)
			if err != nil {
				blog.Errorf("syncDiffHosts fail, err:%s, rid:%s", err.Error(), h.readKit.Rid)
				return err
			}
			syncResult.Detail.Update.Count += result.SuccessInfo.Count
			syncResult.Detail.Update.IPs = append(syncResult.Detail.Update.IPs, result.SuccessInfo.IPs...)
		default:
			blog.Errorf("syncDiffHosts fail, op:%s is invalid, rid:%s", op, h.readKit.Rid)
			return fmt.Errorf("syncDiffHosts op:%s is invalid", op)
		}
		syncResult.SuccessInfo.Count += result.SuccessInfo.Count
		syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, result.SuccessInfo.IPs...)
		syncResult.FailInfo.Count += result.FailInfo.Count
		for ip, errinfo := range result.FailInfo.IPError {
			syncResult.FailInfo.IPError[ip] = errinfo
		}
	}

	return nil
}

// 增加任务同步历史记录
func (h *HostSyncor) addSyncHistory(syncResult *metadata.SyncResult, taskid int64) (*metadata.SyncHistory, error) {
	syncHistory := metadata.SyncHistory{
		TaskID:            taskid,
		SyncStatus:        syncResult.SyncStatus,
		StatusDescription: syncResult.StatusDescription,
		Detail:            syncResult.Detail,
	}
	result, err := h.logics.CreateSyncHistory(h.writeKit, &syncHistory)
	if err != nil {
		blog.Errorf("addSyncHistory err:%v, rid:%s", err.Error(), h.readKit.Rid)
		return nil, err
	}
	return result, nil
}

// 设置SyncResult的状态信息
func (h *HostSyncor) SetSyncResultStatus(syncResult *metadata.SyncResult, startTime time.Time) error {
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
	syncResult.SyncStatus = syncStatus
	syncResult.StatusDescription = statusDesc
	return nil
}

// 根据账号vpcID获取云区域ID
func (h *HostSyncor) getCloudId(vpcID string) (int64, error) {
	cond := mapstr.MapStr{common.BKVpcID: vpcID}
	query := &metadata.QueryCondition{
		Condition: cond,
	}
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(h.readKit.Ctx, h.readKit.Header,
		common.BKInnerObjIDPlat, query)
	if nil != err {
		blog.Errorf("getCloudId failed, error: %v query:%#v, rid:%s", err, query, h.readKit.Rid)
		return 0, err
	}
	if false == res.Result {
		blog.Errorf("getCloudId failed, query:%#v, err code:%d, err msg:%s, rid:%s", query, res.Code, res.ErrMsg,
			h.readKit.Rid)
		return 0, fmt.Errorf("%s", res.ErrMsg)
	}
	if len(res.Data.Info) == 0 {
		return 0, nil
	}
	cloudID, err := res.Data.Info[0].Int64(common.BKCloudIDField)
	if err != nil {
		blog.Errorf("getCloudId failed, err:%v, query:%#v, rid:%s", err, query, h.readKit.Rid)
		return 0, nil
	}
	return cloudID, nil
}

// 创建vpc对应的云区域
func (h *HostSyncor) createCloudArea(vpc *metadata.VpcSyncInfo, accountConf *metadata.CloudAccountConf) (int64, error) {
	cloudArea := map[string]interface{}{
		common.BKCloudNameField:  fmt.Sprintf("%d_%s", accountConf.AccountID, vpc.VpcID),
		common.BKCloudVendor:     accountConf.VendorName,
		common.BKVpcID:           vpc.VpcID,
		common.BKVpcName:         vpc.VpcName,
		common.BKRegion:          vpc.Region,
		common.BKCloudAccountID:  accountConf.AccountID,
		common.BKCreator:         common.BKCloudSyncUser,
		common.BKLastEditor:      common.BKCloudSyncUser,
		common.BkSupplierAccount: common.BKDefaultOwnerID,
		common.BKStatus:          "1",
	}

	instInfo := &metadata.CreateModelInstance{
		Data: mapstr.NewFromMap(cloudArea),
	}

	createRes, err := h.logics.CoreAPI.CoreService().Instance().CreateInstance(h.writeKit.Ctx, h.writeKit.Header, common.BKInnerObjIDPlat, instInfo)
	if nil != err {
		blog.Errorf("createCloudArea failed, error: %s, input:%#v, rid:%s", err.Error(), cloudArea, h.readKit.Rid)
		return 0, err
	}

	if false == createRes.Result {
		blog.Errorf("createCloudArea failed, error code:%d,err msg:%s,input:%#v, rid:%s", createRes.Code,
			createRes.ErrMsg, cloudArea, h.readKit.Rid)
		return 0, fmt.Errorf("%s", createRes.ErrMsg)
	}

	cloudID := int64(createRes.Data.Created.ID)

	// generate audit log.
	audit := auditlog.NewCloudAreaAuditLog(h.logics.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(h.readKit, metadata.AuditCreate)
	logs, err := audit.GenerateAuditLog(generateAuditParameter, []int64{cloudID})
	if err != nil {
		blog.Errorf("generate audit log failed after create cloud area, err: %v, rid: %s", err, h.readKit.Rid)
		return cloudID, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(h.writeKit, logs...); err != nil {
		blog.Errorf("save audit log failed after create cloud area, err: %v, rid: %s", err, h.readKit.Rid)
		return cloudID, err
	}

	return cloudID, nil
}

// getHostIDAndIP get hostID and innerIP by hostInfo.
func getHostIDAndIP(hostInfo map[string]interface{}) (int64, string, error) {
	var hostID int64
	var innerIP string
	if hostIDI, ok := hostInfo[common.BKHostIDField]; ok {
		if hostIDVal, err := strconv.ParseInt(fmt.Sprintf("%v", hostIDI), 10, 64); err == nil {
			hostID = hostIDVal
		}
	}

	if innerIPI, ok := hostInfo[common.BKHostInnerIPField]; ok {
		innerIP = fmt.Sprintf("%s", innerIPI)
	}

	if hostID == 0 {
		blog.Errorf("getHostIDAndIP fail,hostID is 0, hostInfo:%+v", hostInfo)
		return 0, "", fmt.Errorf("%s", "hostID is 0")
	}

	return hostID, innerIP, nil
}

// 获取本地数据库中的主机信息
func (h *HostSyncor) getLocalHosts(cloudIDs []int64) ([]*metadata.CloudHost, error) {
	result := make([]*metadata.CloudHost, 0)
	cond := mapstr.MapStr{
		common.BKCloudIDField: mapstr.MapStr{common.BKDBIN: cloudIDs},
		// 必须带有实例id，说明是云主机
		common.BKCloudInstIDField: mapstr.MapStr{common.BKDBNIN: []interface{}{nil, ""}},
	}
	query := &metadata.QueryCondition{
		Condition: cond,
	}
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(h.readKit.Ctx, h.readKit.Header,
		common.BKInnerObjIDHost, query)
	if nil != err {
		blog.Errorf("getLocalHosts failed, error: %v query:%#v, rid:%s", err, query, h.readKit.Rid)
		return nil, err
	}
	if false == res.Result {
		blog.Errorf("getLocalHosts failed, query:%#v, err code:%d, err msg:%s, rid:%s", query, res.Code, res.ErrMsg,
			h.readKit.Rid)
		return nil, fmt.Errorf("%s", res.ErrMsg)
	}
	if len(res.Data.Info) == 0 {
		return nil, nil
	}

	for _, host := range res.Data.Info {
		instID, _ := host.String(common.BKCloudInstIDField)
		hostStatus, _ := host.String(common.BKCloudHostStatusField)
		privateIp, _ := host.String(common.BKHostInnerIPField)
		publicIp, _ := host.String(common.BKHostOuterIPField)
		cloudID, _ := host.Int64(common.BKCloudIDField)
		hostID, _ := host.Int64(common.BKHostIDField)
		result = append(result, &metadata.CloudHost{
			Instance: metadata.Instance{
				InstanceId:    instID,
				InstanceState: hostStatus,
				PrivateIp:     privateIp,
				PublicIp:      publicIp,
			},
			CloudID: cloudID,
			HostID:  hostID,
		})
	}

	return result, nil
}

// 添加云主机到本地数据库和主机资源池目录对应关系
func (h *HostSyncor) addHosts(hosts []*metadata.CloudHost) (*metadata.SyncResult, error) {
	syncResult := new(metadata.SyncResult)
	syncResult.FailInfo.IPError = make(map[string]string)

	instIDs := make([]string, 0)
	for _, host := range hosts {
		_, err := h.addHost(host)
		if err != nil {
			blog.Errorf("addHosts err:%s, rid:%s", err.Error(), h.readKit.Rid)
			syncResult.FailInfo.Count++
			syncResult.FailInfo.IPError[host.PrivateIp] = err.Error()
			return nil, err
		} else {
			instIDs = append(instIDs, host.InstanceId)
			syncResult.SuccessInfo.Count++
			syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, host.PrivateIp)
		}
	}

	// 添加审计日志
	audit := auditlog.NewHostAudit(h.logics.CoreAPI.CoreService())
	logContext := make([]metadata.AuditLog, 0)
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(h.readKit,
		metadata.AuditCreate).WithOperateFrom(metadata.FromCloudSync)
	curData, err := h.getHostDetailByInstIDs(h.readKit, instIDs)
	if err != nil {
		blog.Errorf("addHosts getHostDetailByInstIDs err:%s, instIDs:%#v, rid:%s", err.Error(), instIDs, h.readKit.Rid)
	}

	for _, data := range curData {
		hostID, innerIP, err := getHostIDAndIP(data)
		if err != nil {
			blog.Errorf("generate audit log failed after create host, failed to get hostID and hostIP, err: %v, rid: %s",
				err, h.readKit.Rid)
			return nil, err
		}

		// generate audit log.
		tmpAuditLog, err := audit.GenerateAuditLogByHostIDGetBizID(generateAuditParameter, hostID, innerIP, data)
		if err != nil {
			blog.Errorf("generate audit log failed after create host, hostID: %d, innerIP: %s, err: %v, rid: %s",
				hostID, innerIP, err, h.readKit.Rid)
			return nil, err
		}

		logContext = append(logContext, *tmpAuditLog)
	}

	// save audit log.
	if len(logContext) > 0 {
		if err := audit.SaveAuditLog(h.writeKit, logContext...); err != nil {
			blog.Errorf("save audit log failed after create host, err: %v, rid: %s", err, h.readKit.Rid)
			return nil, err
		}
	}

	return syncResult, nil
}

// 添加云主机
func (h *HostSyncor) addHost(cHost *metadata.CloudHost) (string, error) {
	host := mapstr.MapStr{
		common.BKCloudIDField:         cHost.CloudID,
		common.BKCloudInstIDField:     cHost.InstanceId,
		common.BKHostInnerIPField:     cHost.PrivateIp,
		common.BKHostOuterIPField:     cHost.PublicIp,
		common.BKCloudHostStatusField: cHost.InstanceState,
		common.BKCloudVendor:          cHost.VendorName,
	}
	input := &metadata.CreateModelInstance{
		Data: host,
	}

	result, err := h.logics.CoreAPI.CoreService().Instance().CreateInstance(h.writeKit.Ctx, h.writeKit.Header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, input:%+v, rid:%s", err.Error(), host, h.readKit.Rid)
		return "", err
	}
	if !result.Result {
		blog.Errorf("addHost fail,err:%s, input:%+v, rid:%s", result.ErrMsg, host, h.readKit.Rid)
		return "", fmt.Errorf("%s", result.ErrMsg)
	}

	hostID := int64(result.Data.Created.ID)

	// 获取资源池业务id
	condition := mapstr.MapStr{
		common.BKDefaultField: common.DefaultAppFlag,
	}
	cond := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField},
		Condition: condition,
	}
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(h.readKit.Ctx, h.readKit.Header,
		common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, cond:%+v, rid:%s", err.Error(), *cond, h.readKit.Rid)
		return "", err
	}
	if !res.Result {
		blog.Errorf("addHost fail,err:%s, cond:%+v, rid:%s", res.ErrMsg, *cond, h.readKit.Rid)
		return "", fmt.Errorf("%s", res.ErrMsg)
	}

	if len(res.Data.Info) == 0 {
		blog.Errorf("addHost fail,err:%s, cond:%+v, rid:%s", "no default biz is found", *cond, h.readKit.Rid)
		return "", fmt.Errorf("%s", "no default biz is found")
	}

	appID, err := res.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, cond:%+v, rid:%s", err.Error(), *cond, h.readKit.Rid)
		return "", err
	}

	// 添加主机和同步目录模块的关系
	opt := &metadata.TransferHostToInnerModule{
		ApplicationID: appID,
		ModuleID:      cHost.SyncDir,
		HostID:        []int64{hostID},
	}
	hResult, err := h.logics.CoreAPI.CoreService().Host().TransferToInnerModule(h.writeKit.Ctx, h.writeKit.Header, opt)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, opt:%+v, rid:%s", err.Error(), *opt, h.readKit.Rid)
		return "", err
	}
	if !hResult.Result {
		blog.Errorf("addHost fail,err:%s, opt:%+v, rid:%s", hResult.ErrMsg, *opt, h.readKit.Rid)
		if len(hResult.Data) > 0 {
			return "", fmt.Errorf("%s", hResult.Data[0].Message)
		}
		return "", hResult.CCError()
	}

	return cHost.PrivateIp, nil
}

// 更新云主机到本地数据库
func (h *HostSyncor) updateHosts(hosts []*metadata.CloudHost) (*metadata.SyncResult, error) {
	syncResult := new(metadata.SyncResult)
	syncResult.FailInfo.IPError = make(map[string]string)

	// for audit log.
	audit := auditlog.NewHostAudit(h.logics.CoreAPI.CoreService())
	logContext := make([]metadata.AuditLog, 0)

	for _, host := range hosts {
		updateInfo := mapstr.MapStr{
			common.BKCloudIDField:         host.CloudID,
			common.BKHostInnerIPField:     host.PrivateIp,
			common.BKHostOuterIPField:     host.PublicIp,
			common.BKCloudHostStatusField: host.InstanceState,
		}

		// generate audit log.
		preData, err := h.getHostDetailByInstIDs(h.readKit, []string{host.InstanceId})
		if err != nil {
			blog.Errorf("get host detail failed, err: %v, instID: %s, rid:%s", err, host.InstanceId, h.readKit.Rid)
			return nil, err
		}
		if len(preData) <= 0 {
			blog.Errorf("generate audit log failed, not find host data, instID: %s, rid: %s", host.InstanceId,
				h.readKit.Rid)
			return nil, fmt.Errorf("generate audit log failed, not find host data when bk_cloud_inst_id is %s", host.InstanceId)
		}

		hostID, innerIP, err := getHostIDAndIP(preData[0])
		if err != nil {
			blog.Errorf("generate audit log failed before update host, failed to get hostID and hostIP, err: %v, rid: %s",
				err, h.readKit.Rid)
			return nil, err
		}

		// generate audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(h.readKit, metadata.AuditUpdate).
			WithOperateFrom(metadata.FromCloudSync).WithUpdateFields(updateInfo)
		tmpAuditLog, err := audit.GenerateAuditLogByHostIDGetBizID(generateAuditParameter, hostID, innerIP, preData[0])
		if err != nil {
			blog.Errorf("generate audit log failed before update host, hostID: %d, innerIP: %s, err: %v, rid: %s",
				hostID, innerIP, err, h.readKit.Rid)
			return nil, err
		}

		// to update.
		if err := h.updateHost(host.InstanceId, updateInfo); err != nil {
			blog.Errorf("updateHosts err:%v, rid:%s", err.Error(), h.readKit.Rid)
			syncResult.FailInfo.Count++
			syncResult.FailInfo.IPError[host.PrivateIp] = err.Error()
			return nil, err
		} else {
			syncResult.SuccessInfo.Count++
			syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, host.PrivateIp)
		}

		// add audit log.
		logContext = append(logContext, *tmpAuditLog)
	}

	// save audit log.
	if len(logContext) > 0 {
		if err := audit.SaveAuditLog(h.writeKit, logContext...); err != nil {
			blog.Errorf("save audit log failed after update host, err: %v, rid: %s", err, h.readKit.Rid)
			return nil, err
		}
	}

	return syncResult, nil
}

// 更新云主机
func (h *HostSyncor) updateHost(cloudInstID string, updateInfo map[string]interface{}) error {
	input := &metadata.UpdateOption{
		CanEditAll: true,
	}
	input.Condition = map[string]interface{}{common.BKCloudInstIDField: cloudInstID}
	input.Data = updateInfo
	uResult, err := h.logics.CoreAPI.CoreService().Instance().UpdateInstance(h.writeKit.Ctx, h.writeKit.Header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("updateHost fail,err:%s, input:%+v, rid:%s", err.Error(), *input, h.readKit.Rid)
		return err
	}
	if !uResult.Result {
		blog.Errorf("updateHost fail,err:%s, input:%+v, rid:%s", uResult.ErrMsg, *input, h.readKit.Rid)
		return uResult.CCError()
	}
	return nil
}

// 删除被销毁云主机相关联的数据
func (h *HostSyncor) deleteDestroyedHosts(hostIDs []int64) (*metadata.SyncResult, error) {
	result := new(metadata.SyncResult)
	result.FailInfo.IPError = make(map[string]string)

	if len(hostIDs) == 0 {
		return result, nil
	}

	// for audit log.
	audit := auditlog.NewHostAudit(h.logics.CoreAPI.CoreService())
	logContext := make([]metadata.AuditLog, 0)
	innerIPs := make([]string, 0)
	preData, err := h.getHostDetailByHostIDs(h.readKit, hostIDs)
	if err != nil {
		blog.Errorf("deleteDestroyedHosts getHostDetailByHostIDs err:%s, hostIDs:%#v, rid:%s", err.Error(), hostIDs,
			h.readKit.Rid)
		return nil, err
	}

	updateHostData := mapstr.MapStr{
		common.BKHostInnerIPField:     []string{},
		common.BKHostOuterIPField:     []string{},
		common.BKCloudHostStatusField: common.BKCloudHostStatusDestroyed,
	}

	// generate audit log.
	for _, data := range preData {
		hostID, err := data.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("generate audit log failed before update host, failed to get hostID, err:%s, data:%#v, rid:%s",
				err.Error(), data, h.readKit.Rid)
			return nil, err
		}
		innerIP, err := data.String(common.BKHostInnerIPField)
		if err != nil {
			blog.Errorf("generate audit log failed before update host, failed to get InnerIP, err:%s, data:%#v, rid:%s",
				err.Error(), data, h.readKit.Rid)
			return nil, err
		}

		innerIPs = append(innerIPs, innerIP)

		// generate audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(h.readKit, metadata.AuditUpdate).
			WithOperateFrom(metadata.FromCloudSync).WithUpdateFields(updateHostData)
		tmpAuditLog, err := audit.GenerateAuditLogByHostIDGetBizID(generateAuditParameter, hostID, innerIP, data)
		if err != nil {
			blog.Errorf("generate audit log failed before update host, hostID: %d, err: %v, rid: %s",
				hostID, err, h.readKit.Rid)
			return nil, err
		}

		// add audit log.
		logContext = append(logContext, *tmpAuditLog)
	}

	result.SuccessInfo.Count = int64(len(hostIDs))
	result.SuccessInfo.IPs = innerIPs

	// to change state of cloud host.
	err = h.logics.CoreAPI.CoreService().Cloud().DeleteDestroyedHostRelated(h.writeKit.Ctx, h.writeKit.Header, &metadata.DeleteDestroyedHostRelatedOption{HostIDs: hostIDs})
	if err != nil {
		blog.Errorf("deleteDestroyedHosts failed, err:%s, hostIDs:%#v, rid:%s", err.Error(), hostIDs, h.readKit.Rid)
		return nil, err
	}

	// save audit log.
	if len(logContext) > 0 {
		if err := audit.SaveAuditLog(h.writeKit, logContext...); err != nil {
			blog.Errorf("save audit log failed after update host, err: %v, rid: %s", err, h.readKit.Rid)
			return nil, err
		}
	}

	return result, nil
}

// 更新被销毁vpc对应的云区域状态为异常
func (h *HostSyncor) updateDestroyedCloudArea(cloudIDs []int64) error {
	input := &metadata.UpdateOption{
		// must set CanEditAll as true to update the field which can't be editable
		CanEditAll: true,
	}
	input.Condition = mapstr.MapStr{common.BKCloudIDField: map[string]interface{}{
		common.BKDBIN: cloudIDs,
	}}
	input.Data = mapstr.MapStr{
		common.BKStatus:       common.BKCloudAreaStatusAbnormal,
		common.BKStatusDetail: "the correspond vpc is destroyed",
	}

	// generate audit log.
	audit := auditlog.NewCloudAreaAuditLog(h.logics.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(h.readKit,
		metadata.AuditUpdate).WithUpdateFields(input.Data)
	logs, err := audit.GenerateAuditLog(generateAuditParameter, cloudIDs)
	if err != nil {
		blog.Errorf("generate audit log failed before update cloud area, err: %v, rid: %s", err, h.readKit.Rid)
		return err
	}

	// to update.
	uResult, err := h.logics.CoreAPI.CoreService().Instance().UpdateInstance(h.writeKit.Ctx, h.writeKit.Header, common.BKInnerObjIDPlat, input)
	if err != nil {
		blog.Errorf("updateDestroyedCloudArea fail,err:%s, input:%+v, rid:%s", err.Error(), *input, h.readKit.Rid)
		return err
	}
	if !uResult.Result {
		blog.Errorf("updateDestroyedCloudArea fail,err:%s, input:%+v, rid:%s", uResult.ErrMsg, *input, h.readKit.Rid)
		return uResult.CCError()
	}

	// save audit log.
	if err := audit.SaveAuditLog(h.writeKit, logs...); err != nil {
		blog.Errorf("save audit log failed after update cloud area, err: %v, rid: %s", err, h.readKit.Rid)
		return err
	}

	return nil
}

// 更新同步任务里的vpc状态为被销毁
func (h *HostSyncor) updateDestroyedTaskVpc(taskID int64, vpcs map[string]bool) error {
	opt := &metadata.SearchCloudOption{
		Condition: mapstr.MapStr{common.BKCloudSyncTaskID: taskID},
	}

	ret, err := h.logics.CoreAPI.CoreService().Cloud().SearchSyncTask(h.readKit.Ctx, h.readKit.Header, opt)
	if err != nil {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, opt: err: %s, rid:%s", taskID, opt, err.Error(),
			h.readKit.Rid)
		return err
	}
	if len(ret.Info) == 0 {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, opt: err: %s, rid:%s", taskID, opt,
			"no task is found", h.readKit.Rid)
		return fmt.Errorf("no task is found")
	}

	syncInfo := ret.Info[0].SyncVpcs
	for i, info := range syncInfo {
		if vpcs[info.VpcID] {
			syncInfo[i].Destroyed = true
		}
	}

	option := map[string]interface{}{common.BKCloudSyncVpcs: syncInfo}
	if err := h.logics.UpdateSyncTask(h.writeKit, taskID, option); err != nil {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, err: %s, rid:%s", taskID, err.Error(), h.readKit.Rid)
		return err
	}

	return nil
}

// 更新任务同步状态
func (h *HostSyncor) updateTaskState(kit *rest.Kit, taskid int64, status string, syncStatusDesc *metadata.SyncStatusDesc) error {
	option := mapstr.MapStr{common.BKCloudSyncStatus: status}
	if status == metadata.CloudSyncSuccess || status == metadata.CloudSyncFail {
		ts := time.Now()
		option.Set(common.BKCloudLastSyncTime, &ts)
		option.Set(common.BKCloudSyncStatusDescription, syncStatusDesc)
	}

	if err := h.logics.UpdateSyncTask(h.writeKit, taskid, option); err != nil {
		blog.Errorf("UpdateSyncTask failed, taskid: %v, err: %s, rid:%s", taskid, err.Error(), kit.Rid)
		return err
	}

	return nil
}

// 根据主机实例ID获取主机详情
func (h *HostSyncor) getHostDetailByInstIDs(kit *rest.Kit, instIDs []string) ([]mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKCloudInstIDField: mapstr.MapStr{
			common.BKDBIN: instIDs,
		},
	}
	return h.getHostDetail(kit, cond)
}

// 根据主机ID获取主机详情
func (h *HostSyncor) getHostDetailByHostIDs(kit *rest.Kit, hostIDs []int64) ([]mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKHostIDField: mapstr.MapStr{
			common.BKDBIN: hostIDs,
		},
	}
	return h.getHostDetail(kit, cond)
}

// 获取主机详情
func (h *HostSyncor) getHostDetail(kit *rest.Kit, cond mapstr.MapStr) ([]mapstr.MapStr, error) {
	query := &metadata.QueryCondition{
		Condition: cond,
	}
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost, query)
	if nil != err {
		blog.Errorf("getHostDetail failed, error: %v query:%#v, rid:%s", err, query, kit.Rid)
		return nil, err
	}
	if false == res.Result {
		blog.Errorf("getHostDetail failed, query:%#v, err msg:%s, rid:%s", query, res.ErrMsg, kit.Rid)
		return nil, fmt.Errorf("%s", res.ErrMsg)
	}
	if len(res.Data.Info) == 0 {
		blog.Errorf("getHostDetail fail, host is not found, query:%#v, rid:%s", query, kit.Rid)
		return nil, fmt.Errorf("host is not found")
	}

	return res.Data.Info, nil
}
