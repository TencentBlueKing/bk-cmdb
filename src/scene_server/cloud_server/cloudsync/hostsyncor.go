package cloudsync

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	ccom "configcenter/src/scene_server/cloud_server/common"
	"configcenter/src/scene_server/cloud_server/logics"
)

// 云主机同步器
type HostSyncor struct {
	id        int
	logics    *logics.Logics
	enableTxn bool
	kit       *rest.Kit
	// used in cases which don't need transaction
	kitNoTxn *rest.Kit
}

// 创建云主机同步器
func NewHostSyncor(id int, logics *logics.Logics) *HostSyncor {
	return &HostSyncor{
		id:        id,
		logics:    logics,
		enableTxn: true,
	}
}

// 同步云主机
func (h *HostSyncor) Sync(task *metadata.CloudSyncTask) error {
	defer func() {
		if err := recover(); err != nil {
			blog.Errorf("sync panic err:%#v, rid:%s, debug strace:%s", err, h.kit.Rid, debug.Stack())
		}
	}()

	blog.V(4).Infof("hostSyncor%d start sync", h.id)
	startTime := time.Now()
	// 每次同步生成新的kit
	h.kit = ccom.NewKit()
	h.kitNoTxn = ccom.NewKit()
	// 两个kit的requestID保持一致，方便追踪日志
	h.kitNoTxn.Header.Set(common.BKHTTPCCRequestID, h.kit.Header.Get(common.BKHTTPCCRequestID))
	// 根据账号id获取账号详情
	accountConf, err := h.logics.GetCloudAccountConf(h.kitNoTxn, task.AccountID)
	if err != nil {
		blog.Errorf("hostSyncor%d GetCloudAccountConf fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
		return err
	}

	// 根据任务详情和账号信息获取要同步的云主机资源
	hostResource, err := h.getCloudHostResource(task, accountConf)
	if err != nil {
		blog.Errorf("hostSyncor%d getCloudHostResource fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
		return err
	}
	if len(hostResource.HostResource) == 0 && len(hostResource.DestroyedVpcs) == 0 {
		blog.V(4).Infof("hostSyncor%d hostResource is empty, taskid:%d, rid:%s", h.id, task.TaskID, h.kit.Rid)
		return nil
	}

	blog.V(4).Infof("hostSyncor%d, taskid:%d, destroyed vpc count:%d, other vpc count:%d, rid:%s",
		h.id, task.TaskID, len(hostResource.DestroyedVpcs), len(hostResource.HostResource), h.kit.Rid)

	syncResult := new(metadata.SyncResult)
	syncResult.FailInfo.IPError = make(map[string]string)

	txnErr := h.logics.CoreAPI.CoreService().Txn().AutoRunTxn(h.kit.Ctx, h.enableTxn, h.kit.Header, func() error {
		if len(hostResource.DestroyedVpcs) > 0 {
			// 同步被销毁的VPC相关资源
			err = h.syncDestroyedVpcs(hostResource, syncResult)
			if err != nil {
				blog.Errorf("hostSyncor%d syncDestroyedVpcs fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
				return err
			}
		}

		// 查询vpc对应的云区域并更新云主机资源信息里的云区域id
		err = h.addCLoudId(accountConf, hostResource)
		if err != nil {
			blog.Errorf("hostSyncor%d addCLoudId fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
			return err
		}

		// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
		diffHosts, err := h.getDiffHosts(hostResource)
		if err != nil {
			blog.Errorf("hostSyncor%d getDiffHosts fail, taskid:%d, err:%s, rid:%s",
				h.id, task.TaskID, err.Error(), h.kit.Rid)
			return err
		}

		// 没差异则结束
		if len(diffHosts) == 0 {
			blog.V(4).Infof("hostSyncor%d no diff hosts for taskid:%d, rid:%s", h.id, task.TaskID, h.kit.Rid)
			if syncResult.SuccessInfo.Count == 0 {
				blog.V(4).Infof("hostSyncor%d no any hosts need sync for taskid:%d, rid:%s", h.id, task.TaskID, h.kit.Rid)
				return nil
			}
		}

		// 有差异的更新任务同步状态为同步中
		err = h.updateTaskState(h.kit, task.TaskID, metadata.CloudSyncInProgress, nil)
		if err != nil {
			blog.Errorf("hostSyncor%d updateTaskState fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
			return err
		}

		// 同步有差异的主机数据
		err = h.syncDiffHosts(diffHosts, syncResult)
		if err != nil {
			blog.Errorf("hostSyncor%d syncDiffHosts fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
			return err
		}

		// 设置SyncResult的状态信息
		err = h.SetSyncResultStatus(syncResult, startTime)
		if err != nil {
			blog.Errorf("hostSyncor%d SetSyncResultStatus fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
			return err
		}

		// 增加任务同步历史记录
		_, err = h.addSyncHistory(syncResult, task.TaskID)
		if err != nil {
			blog.Errorf("hostSyncor%d addSyncHistory fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
			return err
		}

		// 完成后更新任务同步状态
		err = h.updateTaskState(h.kit, task.TaskID, syncResult.SyncStatus, &syncResult.StatusDescription)
		if err != nil {
			blog.Errorf("hostSyncor%d updateTaskState fail, taskid:%d, err:%s, rid:%s", h.id, task.TaskID, err.Error(), h.kit.Rid)
			return err
		}

		blog.Infof("hostSyncor%d finish SyncCloudHost, costTime:%ds, syncResult.Detail:%#v, syncResult.FailInfo:%#v, rid:%s",
			h.id, time.Since(startTime)/time.Second, syncResult.Detail, syncResult.FailInfo, h.kit.Rid)

		return nil
	})
	if txnErr != nil {
		blog.Errorf("hostSyncor%d sync fail, taskid:%d, txnErr:%v, rid:%s", h.id, task.TaskID, txnErr, h.kit.Rid)
		err := h.updateTaskState(h.kitNoTxn, task.TaskID, metadata.CloudSyncFail, &metadata.SyncStatusDesc{ErrorInfo: txnErr.Error()})
		if err != nil {
			blog.Errorf("hostSyncor%d updateTaskState fail, taskid:%d, err:%v, rid:%s", h.id, task.TaskID, err, h.kit.Rid)
		}
	} else {
		// 检查同步状态，如果为异常，则更新为正常, 以在没有主机差异要同步时也可以恢复某些问题导致的状态异常；已经是正常的情况下不需要更新
		opt := &metadata.SearchCloudOption{
			Condition: mapstr.MapStr{common.BKCloudSyncTaskID: task.TaskID},
		}

		ret, err := h.logics.CoreAPI.CoreService().Cloud().SearchSyncTask(h.kitNoTxn.Ctx, h.kitNoTxn.Header, opt)
		if err != nil {
			blog.Errorf("hostSyncor%d SearchSyncTask failed, taskid: %v, opt:%#v, err: %s, rid:%s", h.id, task.TaskID, opt, err.Error(), h.kit.Rid)
			return err
		}
		if len(ret.Info) == 0 {
			blog.Errorf("hostSyncor%d SearchSyncTask failed, taskid %d is not found, opt:%#v, rid:%s", h.id, task.TaskID, opt, h.kit.Rid)
			return fmt.Errorf("taskID %d is not found", task.TaskID)
		}
		if ret.Info[0].SyncStatus == metadata.CloudSyncFail {
			costTime, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", float64(time.Since(startTime)/time.Millisecond)/1000.0), 64)
			err := h.updateTaskState(h.kitNoTxn, task.TaskID, metadata.CloudSyncSuccess, &metadata.SyncStatusDesc{CostTime: costTime})
			if err != nil {
				blog.Errorf("hostSyncor%d updateTaskState fail, taskid:%d, err:%v, rid:%s", h.id, task.TaskID, err, h.kit.Rid)
			}
			blog.Infof("hostSyncor%d update taskid:%d status from fail to success, rid:%s", h.id, task.TaskID, h.kit.Rid)
		}
	}

	blog.V(4).Infof("hostSyncor%d sync loop is over, costTime:%ds, txnErr:%v, rid:%s",
		h.id, time.Since(startTime)/time.Second, txnErr, h.kit.Rid)

	return nil
}

// 根据任务详情和账号信息获取要同步的云主机资源
func (h *HostSyncor) getCloudHostResource(task *metadata.CloudSyncTask, accountConf *metadata.CloudAccountConf) (*metadata.CloudHostResource, error) {
	hostResource, err := h.logics.GetCloudHostResource(h.kitNoTxn, *accountConf, task.SyncVpcs)
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
	blog.V(4).Infof("Destroyed cloudIDs: %#v, rid:%s", cloudIDs, h.kit.Rid)

	// 更新属于被销毁vpc下的主机信息，将内外网ip置空，状态置为已销毁
	condition := mapstr.MapStr{common.BKCloudIDField: map[string]interface{}{
		common.BKDBIN: cloudIDs,
	}}

	query := &metadata.QueryCondition{
		Condition: condition,
	}
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(h.kit.Ctx, h.kit.Header, common.BKInnerObjIDHost, query)
	if nil != err {
		blog.Errorf("syncDestroyedVpcs ReadInstance failed, error: %v query:%#v, rid:%s", err, query, h.kit.Rid)
		return err
	}
	if false == res.Result {
		blog.Errorf("syncDestroyedVpcs failed, query:%#v, err code:%d, err msg:%s, rid:%s", query, res.Code, res.ErrMsg, h.kit.Rid)
		return fmt.Errorf("%s", res.ErrMsg)
	}
	hostIDs := make([]int64, 0)
	for _, host := range res.Data.Info {
		hostID, _ := host.Int64(common.BKHostIDField)
		hostIDs = append(hostIDs, hostID)
	}
	sResult, err := h.deleteDestroyedHosts(hostIDs)
	if err != nil {
		blog.Errorf("syncDestroyedVpcs deleteDestroyedHosts fail, cloudIDs:%#v, err:%s, rid:%s", cloudIDs, err.Error(), h.kit.Rid)
		return err
	}
	syncResult.Detail.Update.Count += sResult.SuccessInfo.Count
	syncResult.Detail.Update.IPs = append(syncResult.Detail.Update.IPs, sResult.SuccessInfo.IPs...)
	syncResult.SuccessInfo.Count += sResult.SuccessInfo.Count
	syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, sResult.SuccessInfo.IPs...)

	// 更新被销毁vpc对应的云区域状态为异常
	if err := h.updateDestroyedCloudArea(cloudIDs); err != nil {
		blog.Errorf("syncDestroyedVpcs updateDestroyedCloudArea fail, cloudIDs:%#v, err:%s, rid:%s", cloudIDs, err.Error(), h.kit.Rid)
		return err
	}

	vpcs := make(map[string]bool)
	for _, vpcInfo := range hostResource.DestroyedVpcs {
		vpcs[vpcInfo.VpcID] = true
	}
	// 更新同步任务里的vpc状态为被销毁
	if err := h.updateDestroyedTaskVpc(hostResource.TaskID, vpcs); err != nil {
		blog.Errorf("syncDestroyedVpcs updateDestroyedTaskVpc fail, cloudIDs:%#v, err:%s, rid:%s", cloudIDs, err.Error(), h.kit.Rid)
		return err
	}

	return nil
}

// 查询vpc对应的云区域
func (h *HostSyncor) addCLoudId(accountConf *metadata.CloudAccountConf, hostResource *metadata.CloudHostResource) error {
	for _, hostRes := range hostResource.HostResource {
		cloudID, err := h.getCloudId(hostRes.Vpc.VpcID)
		if err != nil {
			blog.Errorf("addCLoudId getCloudId err:%s, vpcID:%s, rid:%s", err.Error(), hostRes.Vpc.VpcID, h.kit.Rid)
			return err
		}
		if cloudID == 0 {
			blog.Errorf("addCLoudId getCloudId err:%s, vpcID:%s, rid:%s", "the correspond cloudID for the vpc can't be found", hostRes.Vpc.VpcID, h.kit.Rid)
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
	blog.V(4).Infof("taskid:%d, host cloudIDs:%#v, rid:%s", hostResource.TaskID, cloudIDs, h.kit.Rid)

	// 本地已有的云主机
	localHosts, err := h.getLocalHosts(cloudIDs)
	if err != nil {
		return nil, err
	}
	blog.V(4).Infof("taskid:%d, len(localHosts):%d, rid:%s", hostResource.TaskID, len(localHosts), h.kit.Rid)
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
				blog.Errorf("syncDiffHosts fail, err:%s, rid:%s", err.Error(), h.kit.Rid)
				return err
			}
			syncResult.Detail.NewAdd.Count += result.SuccessInfo.Count
			syncResult.Detail.NewAdd.IPs = append(syncResult.Detail.NewAdd.IPs, result.SuccessInfo.IPs...)
		case "update":
			result, err = h.updateHosts(hosts)
			if err != nil {
				blog.Errorf("syncDiffHosts fail, err:%s, rid:%s", err.Error(), h.kit.Rid)
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
				blog.Errorf("syncDiffHosts fail, err:%s, rid:%s", err.Error(), h.kit.Rid)
				return err
			}
			syncResult.Detail.Update.Count += result.SuccessInfo.Count
			syncResult.Detail.Update.IPs = append(syncResult.Detail.Update.IPs, result.SuccessInfo.IPs...)
		default:
			blog.Errorf("syncDiffHosts fail, op:%s is invalid, rid:%s", op, h.kit.Rid)
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
	result, err := h.logics.CreateSyncHistory(h.kit, &syncHistory)
	if err != nil {
		blog.Errorf("addSyncHistory err:%v, rid:%s", err.Error(), h.kit.Rid)
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
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(h.kit.Ctx, h.kit.Header, common.BKInnerObjIDPlat, query)
	if nil != err {
		blog.Errorf("getCloudId failed, error: %v query:%#v, rid:%s", err, query, h.kit.Rid)
		return 0, err
	}
	if false == res.Result {
		blog.Errorf("getCloudId failed, query:%#v, err code:%d, err msg:%s, rid:%s", query, res.Code, res.ErrMsg, h.kit.Rid)
		return 0, fmt.Errorf("%s", res.ErrMsg)
	}
	if len(res.Data.Info) == 0 {
		return 0, nil
	}
	cloudID, err := res.Data.Info[0].Int64(common.BKCloudIDField)
	if err != nil {
		blog.Errorf("getCloudId failed, err:%v, query:%#v, rid:%s", err, query, h.kit.Rid)
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

	createRes, err := h.logics.CoreAPI.CoreService().Instance().CreateInstance(h.kit.Ctx, h.kit.Header, common.BKInnerObjIDPlat, instInfo)
	if nil != err {
		blog.Errorf("createCloudArea failed, error: %s, input:%#v, rid:%s", err.Error(), cloudArea, h.kit.Rid)
		return 0, err
	}

	if false == createRes.Result {
		blog.Errorf("createCloudArea failed, error code:%d,err msg:%s,input:%#v, rid:%s", createRes.Code, createRes.ErrMsg, cloudArea, h.kit.Rid)
		return 0, fmt.Errorf("%s", createRes.ErrMsg)
	}

	cloudID := int64(createRes.Data.Created.ID)

	auditLog := h.logics.NewCloudAreaLog(h.kit)
	// create auditLog
	if err := auditLog.WithCurrent(cloudID); err != nil {
		blog.Errorf("updateDestroyedCloudArea WithCurrent err:%s, cloudID:%#v, rid:%s", err.Error(), cloudID, h.kit.Rid)
	}
	if err := auditLog.SaveAuditLog(metadata.AuditUpdate); err != nil {
		blog.Errorf("updateDestroyedCloudArea SaveAuditLog err:%s, cloudID:%#v, rid:%s", err.Error(), cloudID, h.kit.Rid)
	}

	return cloudID, nil
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
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(h.kit.Ctx, h.kit.Header, common.BKInnerObjIDHost, query)
	if nil != err {
		blog.Errorf("getLocalHosts failed, error: %v query:%#v, rid:%s", err, query, h.kit.Rid)
		return nil, err
	}
	if false == res.Result {
		blog.Errorf("getLocalHosts failed, query:%#v, err code:%d, err msg:%s, rid:%s", query, res.Code, res.ErrMsg, h.kit.Rid)
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
			blog.Errorf("addHosts err:%s, rid:%s", err.Error(), h.kit.Rid)
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
	auditLogs := make([]metadata.AuditLog, 0)

	curData, err := h.getHostDetailByInstIDs(h.kit, instIDs)
	if err != nil {
		blog.Errorf("addHosts getHostDetailByInstIDs err:%s, instIDs:%#v, rid:%s", err.Error(), instIDs, h.kit.Rid)
	}
	

	for _, cur := range curData {
		auditLog, err := h.logics.GetAddHostLog(h.kit, cur)
		if err != nil {
			blog.Errorf("addHosts GetAddHostLog err:%s, cur:%#v, rid:%s", err.Error(), cur, h.kit.Rid)
			return nil, err
		}
		auditLogs = append(auditLogs, *auditLog)
	}

	if len(auditLogs) > 0 {
		_, err := h.logics.CoreAPI.CoreService().Audit().SaveAuditLog(h.kit.Ctx, h.kit.Header, auditLogs...)
		if err != nil {
			blog.Errorf("addHosts SaveAuditLog err:%s, rid:%s", err.Error(), h.kit.Rid)
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
	var err error
	result, err := h.logics.CoreAPI.CoreService().Instance().CreateInstance(h.kit.Ctx, h.kit.Header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, input:%+v, rid:%s", err.Error(), host, h.kit.Rid)
		return "", err
	}
	if !result.Result {
		blog.Errorf("addHost fail,err:%s, input:%+v, rid:%s", result.ErrMsg, host, h.kit.Rid)
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
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(h.kit.Ctx, h.kit.Header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, cond:%+v, rid:%s", err.Error(), *cond, h.kit.Rid)
		return "", err
	}
	if !res.Result {
		blog.Errorf("addHost fail,err:%s, cond:%+v, rid:%s", res.ErrMsg, *cond, h.kit.Rid)
		return "", fmt.Errorf("%s", res.ErrMsg)
	}

	if len(res.Data.Info) == 0 {
		blog.Errorf("addHost fail,err:%s, cond:%+v, rid:%s", "no default biz is found", *cond, h.kit.Rid)
		return "", fmt.Errorf("%s", "no default biz is found")
	}

	appID, err := res.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, cond:%+v, rid:%s", err.Error(), *cond, h.kit.Rid)
		return "", err
	}

	// 添加主机和同步目录模块的关系
	opt := &metadata.TransferHostToInnerModule{
		ApplicationID: appID,
		ModuleID:      cHost.SyncDir,
		HostID:        []int64{hostID},
	}
	hResult, err := h.logics.CoreAPI.CoreService().Host().TransferToInnerModule(h.kit.Ctx, h.kit.Header, opt)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, opt:%+v, rid:%s", err.Error(), *opt, h.kit.Rid)
		return "", err
	}
	if !hResult.Result {
		blog.Errorf("addHost fail,err:%s, opt:%+v, rid:%s", hResult.ErrMsg, *opt, h.kit.Rid)
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
	auditLogs := make([]metadata.AuditLog, 0)
	preDataMap := make(map[int64]map[string]interface{})
	CurDataMap := make(map[int64]map[string]interface{})

	instIDs := make([]string, len(hosts))
	for i, host := range hosts {
		instIDs[i] = host.InstanceId
	}

	preData, err := h.getHostDetailByInstIDs(h.kit, instIDs)
	if err != nil {
		blog.Errorf("updateHosts getHostDetailByInstIDs err:%s, instIDs:%#v, rid:%s", err.Error(), instIDs, h.kit.Rid)
	}
	for _, data := range preData {
		hostID, err := data.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("updateHosts Int64 err:%s, data:%#v, rid:%s", err.Error(), data, h.kit.Rid)
			return nil, err
		}
		preDataMap[hostID] = data
	}

	for _, host := range hosts {
		updateInfo := mapstr.MapStr{
			common.BKCloudIDField:         host.CloudID,
			common.BKHostInnerIPField:     host.PrivateIp,
			common.BKHostOuterIPField:     host.PublicIp,
			common.BKCloudHostStatusField: host.InstanceState,
		}
		if err := h.updateHost(host.InstanceId, updateInfo); err != nil {
			blog.Errorf("updateHosts err:%v, rid:%s", err.Error(), h.kit.Rid)
			syncResult.FailInfo.Count++
			syncResult.FailInfo.IPError[host.PrivateIp] = err.Error()
			return nil, err
		} else {
			syncResult.SuccessInfo.Count++
			syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, host.PrivateIp)
		}
	}

	curData, err := h.getHostDetailByInstIDs(h.kit, instIDs)
	if err != nil {
		blog.Errorf("deleteDestroyedHosts getHostDetailByHostIDs err:%s, instIDs:%#v, rid:%s", err.Error(), instIDs, h.kit.Rid)
	}
	for _, data := range curData {
		hostID, err := data.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("deleteDestroyedHosts Int64 err:%s, data:%#v, rid:%s", err.Error(), data, h.kit.Rid)
			return nil, err
		}
		CurDataMap[hostID] = data
	}

	for hostID, cur := range CurDataMap {
		auditLog, err := h.logics.GetUpdateHostLog(h.kit, preDataMap[hostID], cur)
		if err != nil {
			blog.Errorf("updateHosts GetUpdateHostLog err:%s, rid:%s", err.Error(), h.kit.Rid)
			return nil, err
		}
		auditLogs = append(auditLogs, *auditLog)
	}

	if len(auditLogs) > 0 {
		_, err := h.logics.CoreAPI.CoreService().Audit().SaveAuditLog(h.kit.Ctx, h.kit.Header, auditLogs...)
		if err != nil {
			blog.Errorf("updateHosts SaveAuditLog err:%s, rid:%s", err.Error(), h.kit.Rid)
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
	uResult, err := h.logics.CoreAPI.CoreService().Instance().UpdateInstance(h.kit.Ctx, h.kit.Header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("updateHost fail,err:%s, input:%+v, rid:%s", err.Error(), *input, h.kit.Rid)
		return err
	}
	if !uResult.Result {
		blog.Errorf("updateHost fail,err:%s, input:%+v, rid:%s", uResult.ErrMsg, *input, h.kit.Rid)
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

	auditLogs := make([]metadata.AuditLog, 0)
	preDataMap := make(map[int64]map[string]interface{})
	CurDataMap := make(map[int64]map[string]interface{})

	preData, err := h.getHostDetailByHostIDs(h.kit, hostIDs)
	if err != nil {
		blog.Errorf("deleteDestroyedHosts getHostDetailByHostIDs err:%s, hostIDs:%#v, rid:%s", err.Error(), hostIDs, h.kit.Rid)
		return nil, err
	}
	innerIPs := make([]string, 0)
	for _, data := range preData {
		hostID, err := data.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("deleteDestroyedHosts Int64 err:%s, data:%#v, rid:%s", err.Error(), data, h.kit.Rid)
			return nil, err
		}
		preDataMap[hostID] = data

		innerIP, err := data.String(common.BKHostInnerIPField)
		if err != nil {
			blog.Errorf("deleteDestroyedHosts Int64 err:%s, data:%#v, rid:%s", err.Error(), data, h.kit.Rid)
			return nil, err
		}
		innerIPs = append(innerIPs, innerIP)
	}

	result.SuccessInfo.Count = int64(len(preData))
	result.SuccessInfo.IPs = innerIPs

	err = h.logics.CoreAPI.CoreService().Cloud().DeleteDestroyedHostRelated(h.kit.Ctx, h.kit.Header, &metadata.DeleteDestroyedHostRelatedOption{HostIDs: hostIDs})
	if err != nil {
		blog.Errorf("deleteDestroyedHosts failed, err:%s, hostIDs:%#v, rid:%s", err.Error(), hostIDs, h.kit.Rid)
		return nil, err
	}

	curData, err := h.getHostDetailByHostIDs(h.kit, hostIDs)
	if err != nil {
		blog.Errorf("deleteDestroyedHosts getHostDetailByHostIDs err:%s, hostIDs:%#v, rid:%s", err.Error(), hostIDs, h.kit.Rid)
	}
	for _, data := range curData {
		hostID, err := data.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("deleteDestroyedHosts Int64 err:%s, data:%#v, rid:%s", err.Error(), data, h.kit.Rid)
			return nil, err
		}
		CurDataMap[hostID] = data
	}

	for hostID, cur := range CurDataMap {
		auditLog, err := h.logics.GetUpdateHostLog(h.kit, preDataMap[hostID], cur)
		if err != nil {
			blog.Errorf("updateHosts GetUpdateHostLog err:%s, rid:%s", err.Error(), h.kit.Rid)
			return nil, err
		}
		auditLogs = append(auditLogs, *auditLog)
	}

	if len(auditLogs) > 0 {
		_, err := h.logics.CoreAPI.CoreService().Audit().SaveAuditLog(h.kit.Ctx, h.kit.Header, auditLogs...)
		if err != nil {
			blog.Errorf("updateHosts SaveAuditLog err:%s, rid:%s", err.Error(), h.kit.Rid)
			return nil, err
		}
	}

	return result, nil
}

// 更新被销毁vpc对应的云区域状态为异常
func (h *HostSyncor) updateDestroyedCloudArea(cloudIDs []int64) error {
	auditLog := h.logics.NewCloudAreaLog(h.kit)
	if err := auditLog.WithPrevious(cloudIDs...); err != nil {
		blog.Errorf("updateDestroyedCloudArea WithPrevious err:%s, cloudIDs:%#v, rid:%s", err.Error(), cloudIDs, h.kit.Rid)
	}

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

	uResult, err := h.logics.CoreAPI.CoreService().Instance().UpdateInstance(h.kit.Ctx, h.kit.Header, common.BKInnerObjIDPlat, input)
	if err != nil {
		blog.Errorf("updateDestroyedCloudArea fail,err:%s, input:%+v, rid:%s", err.Error(), *input, h.kit.Rid)
		return err
	}
	if !uResult.Result {
		blog.Errorf("updateDestroyedCloudArea fail,err:%s, input:%+v, rid:%s", uResult.ErrMsg, *input, h.kit.Rid)
		return uResult.CCError()
	}

	// update auditLog
	if err := auditLog.WithCurrent(cloudIDs...); err != nil {
		blog.Errorf("updateDestroyedCloudArea WithCurrent err:%s, cloudIDs:%#v, rid:%s", err.Error(), cloudIDs, h.kit.Rid)
	}
	if err := auditLog.SaveAuditLog(metadata.AuditUpdate); err != nil {
		blog.Errorf("updateDestroyedCloudArea SaveAuditLog err:%s, cloudIDs:%#v, rid:%s", err.Error(), cloudIDs, h.kit.Rid)
	}

	return nil
}

// 更新同步任务里的vpc状态为被销毁
func (h *HostSyncor) updateDestroyedTaskVpc(taskID int64, vpcs map[string]bool) error {
	opt := &metadata.SearchCloudOption{
		Condition: mapstr.MapStr{common.BKCloudSyncTaskID: taskID},
	}

	ret, err := h.logics.CoreAPI.CoreService().Cloud().SearchSyncTask(h.kit.Ctx, h.kit.Header, opt)
	if err != nil {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, opt: err: %s, rid:%s", taskID, opt, err.Error(), h.kit.Rid)
		return err
	}
	if len(ret.Info) == 0 {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, opt: err: %s, rid:%s", taskID, opt, "no task is found", h.kit.Rid)
		return fmt.Errorf("no task is found")
	}

	syncInfo := ret.Info[0].SyncVpcs
	for i, info := range syncInfo {
		if vpcs[info.VpcID] {
			syncInfo[i].Destroyed = true
		}
	}

	option := map[string]interface{}{common.BKCloudSyncVpcs: syncInfo}
	if err := h.logics.UpdateSyncTask(h.kit, taskID, option); err != nil {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, err: %s, rid:%s", taskID, err.Error(), h.kit.Rid)
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

	if err := h.logics.UpdateSyncTask(kit, taskid, option); err != nil {
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
