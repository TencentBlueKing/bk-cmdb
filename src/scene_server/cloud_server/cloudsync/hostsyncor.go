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
	"configcenter/src/scene_server/cloud_server/logics"
)

// 云主机同步器
type HostSyncor struct {
	id     int
	logics *logics.Logics
}

// 创建云主机同步器
func NewHostSyncor(id int, logics *logics.Logics) *HostSyncor {
	return &HostSyncor{id, logics}
}

// 同步云主机
func (h *HostSyncor) Sync(task *metadata.CloudSyncTask) error {
	blog.V(4).Infof("hostSyncor%d start sync", h.id)
	startTime := time.Now()
	// 根据账号id获取账号详情
	accountConf, err := h.logics.GetCloudAccountConf(task.AccountID)
	if err != nil {
		blog.Errorf("hostSyncor%d GetCloudAccountConf fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}

	// 根据任务详情和账号信息获取要同步的云主机资源
	hostResource, err := h.getCloudHostResource(task, accountConf)
	if err != nil {
		blog.Errorf("hostSyncor%d getCloudHostResource fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}
	if len(hostResource.HostResource) == 0 && len(hostResource.DestroyedVpcs) == 0 {
		blog.V(4).Infof("hostSyncor%d hostResource is empty, taskid:%d", h.id, task.TaskID)
		return nil
	}

	blog.V(4).Infof("hostSyncor%d, taskid:%d, destroyed vpc count:%d, other vpc count:%d",
		h.id, task.TaskID, len(hostResource.DestroyedVpcs), len(hostResource.HostResource))

	// 同步被销毁的VPC相关资源
	err = h.syncDestroyedVpcs(hostResource)
	if err != nil {
		blog.Errorf("hostSyncor%d syncDestroyedVpcs fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}

	// 查询vpc对应的云区域，没有则创建,并更新云主机资源信息里的云区域id
	h.addCLoudId(accountConf, hostResource)
	if err != nil {
		blog.Errorf("hostSyncor%d addCLoudId fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}

	// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
	diffHosts, err := h.getDiffHosts(hostResource)
	if err != nil {
		blog.Errorf("hostSyncor%d getDiffHosts fail, taskid:%d, err:%s",
			h.id, task.TaskID, err.Error())
		return err
	}

	// 没差异则结束
	if len(diffHosts) == 0 {
		blog.V(4).Infof("hostSyncor%d no diff hosts for taskid:%d", h.id, task.TaskID)
		return nil
	}

	// 有差异的更新任务同步状态为同步中
	err = h.updateTaskState(task.TaskID, metadata.CloudSyncInProgress, nil)
	if err != nil {
		blog.Errorf("hostSyncor%d updateTaskState fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}

	// 如果中途有错误则保证更新同步任务为失败状态
	defer func() {
		if err != nil {
			h.updateTaskState(task.TaskID, metadata.CloudSyncFail, &metadata.SyncStatusDesc{ErrorInfo: err.Error()})
		}
	}()

	// todo 后面几个表操作放在同一个事务里
	// 同步有差异的主机数据
	syncResult, err := h.syncDiffHosts(diffHosts)
	if err != nil {
		blog.Errorf("hostSyncor%d syncDiffHosts fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}

	// 设置SyncResult的状态信息
	err = h.SetSyncResultStatus(syncResult, startTime)
	if err != nil {
		blog.Errorf("hostSyncor%d SetSyncResultStatus fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}

	// 增加任务同步历史记录
	_, err = h.addSyncHistory(syncResult, task.TaskID)
	if err != nil {
		blog.Errorf("hostSyncor%d addSyncHistory fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}

	// 完成后更新任务同步状态
	err = h.updateTaskState(task.TaskID, syncResult.SyncStatus, &syncResult.StatusDescription)
	if err != nil {
		blog.Errorf("hostSyncor%d updateTaskState fail, taskid:%d, err:%s", h.id, task.TaskID, err.Error())
		return err
	}

	costTime := time.Since(startTime) / time.Second
	blog.Infof("hostSyncor%d finish SyncCloudHost, costTime:%ds, syncResult.Detail:%#v, syncResult.FailInfo:%#v",
		h.id, costTime, syncResult.Detail, syncResult.FailInfo)

	return nil
}

// 根据任务详情和账号信息获取要同步的云主机资源
func (h *HostSyncor) getCloudHostResource(task *metadata.CloudSyncTask, accountConf *metadata.CloudAccountConf) (*metadata.CloudHostResource, error) {
	hostResource, err := h.logics.GetCloudHostResource(*accountConf, task.SyncVpcs)
	if err != nil {
		return nil, err
	}
	hostResource.AccountConf = accountConf
	hostResource.TaskID = task.TaskID
	return hostResource, err
}

// 同步被销毁的VPC相关资源
func (h *HostSyncor) syncDestroyedVpcs(hostResource *metadata.CloudHostResource) error {
	if len(hostResource.DestroyedVpcs) == 0 {
		return nil
	}
	cloudIDs := make([]int64, 0)
	for _, vpcInfo := range hostResource.DestroyedVpcs {
		cloudIDs = append(cloudIDs, vpcInfo.CloudID)
	}
	blog.V(4).Infof("Destroyed cloudIDs: %#v", cloudIDs)

	// 更新属于被销毁vpc下的主机信息，将内外网ip置空，状态置为已销毁
	conditon := mapstr.MapStr{common.BKCloudIDField: map[string]interface{}{
		common.BKDBIN: cloudIDs,
	}}


	if err := h.updateDestroyedHosts(conditon); err != nil {
		blog.Errorf("syncDestroyedVpcs fail, cloudIDs:%#v, err:%s", cloudIDs, err.Error())
		return err
	}

	// 更新被销毁vpc对应的云区域状态为异常
	if err := h.updateDestroyedCloudArea(cloudIDs); err != nil {
		blog.Errorf("syncDestroyedVpcs fail, cloudIDs:%#v, err:%s", cloudIDs, err.Error())
		return err
	}

	vpcs := make(map[string]bool)
	for _, vpcInfo := range hostResource.DestroyedVpcs {
		vpcs[vpcInfo.VpcID] = true
	}
	// 更新同步任务里的vpc状态为被销毁
	if err := h.updateDestroyedTaskVpc(hostResource.TaskID, vpcs); err != nil {
		blog.Errorf("syncDestroyedVpcs fail, cloudIDs:%#v, err:%s", cloudIDs, err.Error())
		return err
	}

	return nil
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
			cloudID, err = h.createCloudArea(hostRes.Vpc, accountConf)
			if err != nil {
				continue
			}
		}
		hostRes.CloudID = cloudID
	}
	return nil, nil
}

// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
func (h *HostSyncor) getDiffHosts(hostResource *metadata.CloudHostResource) (map[string][]*metadata.CloudHost, error) {
	// 云端的主机
	remoteHostsMap := make(map[string]*metadata.CloudHost)
	for _, hostRes := range hostResource.HostResource {
		for _, host := range hostRes.Instances {
			host.InstanceState = metadata.CloudHostStatusIDs[host.InstanceState]
			remoteHostsMap[host.InstanceId] = &metadata.CloudHost{
				Instance:   *host,
				CloudID:    hostRes.CloudID,
				VendorName: metadata.VendorNameIDs[hostResource.AccountConf.VendorName],
				SyncDir:    hostRes.Vpc.SyncDir,
			}
		}
	}

	cloudIDs := make([]int64, 0)
	for _, vpcInfo := range hostResource.HostResource {
		cloudIDs = append(cloudIDs, vpcInfo.Vpc.CloudID)
	}
	blog.V(4).Infof("taskid:%d, host cloudIDs:%#v", hostResource.TaskID, cloudIDs)

	// 本地已有的云主机
	localHosts, err := h.getLocalHosts(cloudIDs)
	if err != nil {
		return nil, err
	}
	blog.V(4).Infof("taskid:%d, len(localHosts):%d", hostResource.TaskID, len(localHosts))
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
			if h.InstanceState == lh.InstanceState && lh.InstanceState == metadata.CloudHostStatusIDs["stopped"] {
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
		if _, ok := remoteHostsMap[id]; !ok {
			diffHosts["delete"] = append(diffHosts["delete"], h)
		}
	}

	return diffHosts, nil
}

// 同步有差异的主机数据
func (h *HostSyncor) syncDiffHosts(diffhosts map[string][]*metadata.CloudHost) (*metadata.SyncResult, error) {
	syncResult := new(metadata.SyncResult)
	syncResult.FailInfo.IPError = make(map[string]string)
	var result *metadata.SyncResult
	var err error
	for op, hosts := range diffhosts {
		switch op {
		case "add":
			result, err = h.addHosts(hosts)
			if err != nil {
				blog.Errorf("syncDiffHosts fail, err:%s", err.Error())
				return nil, err
			}
			syncResult.Detail.NewAdd = result.SuccessInfo
		case "update":
			result, err = h.updateHosts(hosts)
			if err != nil {
				blog.Errorf("syncDiffHosts fail, err:%s", err.Error())
				return nil, err
			}
			syncResult.Detail.Update = result.SuccessInfo
		case "delete":
			insts := make([]string, 0)
			for _, h := range hosts {
				insts = append(insts, h.InstanceId)
			}
			conditon := mapstr.MapStr{common.BKCloudInstIDField: map[string]interface{}{
				common.BKDBIN: insts,
			}}
			err = h.updateDestroyedHosts(conditon)
			if err != nil {
				blog.Errorf("syncDiffHosts fail, err:%s", err.Error())
				return nil, err
			}
		default:
			blog.Errorf("syncDiffHosts fail, op:%s is invalid", op)
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
func (h *HostSyncor) addSyncHistory(syncResult *metadata.SyncResult, taskid int64) (*metadata.SyncHistory, error) {
	syncHistory := metadata.SyncHistory{
		TaskID:            taskid,
		SyncStatus:        syncResult.SyncStatus,
		StatusDescription: syncResult.StatusDescription,
		Detail:            syncResult.Detail,
	}
	result, err := h.logics.CreateSyncHistory(kit, &syncHistory)
	if err != nil {
		blog.Errorf("addSyncHistory err:%v", err.Error())
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

// 根据账号vpcID获取云区域ID，没有则创建
func (h *HostSyncor) getCloudId(vpcID string) (int64, error) {
	cond := mapstr.MapStr{common.BKVpcID: vpcID}
	query := &metadata.QueryCondition{
		Condition: cond,
	}
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), header, common.BKInnerObjIDPlat, query)
	if nil != err {
		blog.Errorf("getCloudId failed, error: %v query:%#v", err, query)
		return 0, err
	}
	if false == res.Result {
		blog.Errorf("getCloudId failed, query:%#v, err code:%d, err msg:%s", query, res.Code, res.ErrMsg)
		return 0, fmt.Errorf("%s", res.ErrMsg)
	}
	if len(res.Data.Info) == 0 {
		return 0, nil
	}
	cloudID, err := res.Data.Info[0].Int64(common.BKCloudIDField)
	if err != nil {
		blog.Errorf("getCloudId failed, err:%v, query:%#v", err, query)
		return 0, nil
	}
	return cloudID, nil
}

// 创建vpc对应的云区域
func (h *HostSyncor) createCloudArea(vpc *metadata.VpcSyncInfo, accountConf *metadata.CloudAccountConf) (int64, error) {
	cloudArea := map[string]interface{}{
		common.BKCloudNameField:  fmt.Sprintf("%d_%s", accountConf.AccountID, vpc.VpcID),
		common.BKCloudVendor:     metadata.VendorNameIDs[accountConf.VendorName],
		common.BKVpcID:           vpc.VpcID,
		common.BKVpcName:         vpc.VpcName,
		common.BKReion:           vpc.Region,
		common.BKCloudAccountID:  accountConf.AccountID,
		common.BKCreator:         common.BKCloudSyncUser,
		common.BKLastEditor:      common.BKCloudSyncUser,
		common.BkSupplierAccount: fmt.Sprintf("%d", common.BKDefaultSupplierID),
		common.BKStatus:          "1",
	}

	instInfo := &metadata.CreateModelInstance{
		Data: mapstr.NewFromMap(cloudArea),
	}

	createRes, err := h.logics.CoreAPI.CoreService().Instance().CreateInstance(context.Background(), header, common.BKInnerObjIDPlat, instInfo)
	if nil != err {
		blog.Errorf("createCloudArea failed, error: %s, input:%#v", err.Error(), cloudArea)
		return 0, err
	}

	if false == createRes.Result {
		blog.Errorf("createCloudArea failed, error code:%d,err msg:%s,input:%#v", createRes.Code, createRes.ErrMsg, cloudArea)
		return 0, fmt.Errorf("%s", createRes.ErrMsg)
	}

	return int64(createRes.Data.Created.ID), nil
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
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), header, common.BKInnerObjIDHost, query)
	if nil != err {
		blog.Errorf("getLocalHosts failed, error: %v query:%#v", err, query)
		return nil, err
	}
	if false == res.Result {
		blog.Errorf("getLocalHosts failed, query:%#v, err code:%d, err msg:%s", query, res.Code, res.ErrMsg)
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
	for _, host := range hosts {
		_, err := h.addHost(host)
		if err != nil {
			blog.Errorf("addHosts err:%s", err.Error())
			syncResult.FailInfo.Count++
			syncResult.FailInfo.IPError[host.PrivateIp] = err.Error()
		} else {
			syncResult.SuccessInfo.Count++
			syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, host.PrivateIp)
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
	result, err := h.logics.CoreAPI.CoreService().Instance().CreateInstance(context.Background(), header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, input:%+v", err.Error(), host)
		return "", err
	}
	if !result.Result {
		blog.Errorf("addHost fail,err:%s, input:%+v", result.ErrMsg, host)
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
	res, err := h.logics.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, cond:%+v", err.Error(), *cond)
		return "", err
	}
	if !result.Result {
		blog.Errorf("addHost fail,err:%s, cond:%+v", result.ErrMsg, *cond)
		return "", fmt.Errorf("%s", result.ErrMsg)
	}

	if len(res.Data.Info) == 0 {
		blog.Errorf("addHost fail,err:%s, cond:%+v", "no default biz is found", *cond)
		return "", fmt.Errorf("%s", "no default biz is found")
	}

	appID, err := res.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, cond:%+v", err.Error(), *cond)
		return "", err
	}

	// 添加主机和同步目录模块的关系
	opt := &metadata.TransferHostToInnerModule{
		ApplicationID: appID,
		ModuleID:      cHost.SyncDir,
		HostID:        []int64{hostID},
	}
	hResult, err := h.logics.CoreAPI.CoreService().Host().TransferToInnerModule(context.Background(), header, opt)
	if err != nil {
		blog.Errorf("addHost fail,err:%s, opt:%+v", err.Error(), *opt)
		return "", err
	}
	if !hResult.Result {
		blog.Errorf("addHost fail,err:%s, opt:%+v", hResult.ErrMsg, *opt)
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
	for _, host := range hosts {
		updateInfo := mapstr.MapStr{
			common.BKCloudIDField:         host.CloudID,
			common.BKHostInnerIPField:     host.PrivateIp,
			common.BKHostOuterIPField:     host.PublicIp,
			common.BKCloudHostStatusField: host.InstanceState,
		}
		if err := h.updateHost(host.InstanceId, updateInfo); err != nil {
			blog.Errorf("updateHosts err:%v", err.Error())
			syncResult.FailInfo.Count++
			syncResult.FailInfo.IPError[host.PrivateIp] = err.Error()
			continue
		} else {
			syncResult.SuccessInfo.Count++
			syncResult.SuccessInfo.IPs = append(syncResult.SuccessInfo.IPs, host.PrivateIp)
		}
	}
	return syncResult, nil
}

// 更新云主机
func (h *HostSyncor) updateHost(cloudInstID string, updateInfo map[string]interface{}) error {
	input := &metadata.UpdateOption{}
	input.Condition = map[string]interface{}{common.BKCloudInstIDField: cloudInstID}
	input.Data = updateInfo
	uResult, err := h.logics.CoreAPI.CoreService().Instance().UpdateInstance(context.Background(), header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("updateHost fail,err:%s, input:%+v", err.Error(), *input)
		return err
	}
	if !uResult.Result {
		blog.Errorf("updateHost fail,err:%s, input:%+v", uResult.ErrMsg, *input)
		return uResult.CCError()
	}
	return nil
}

// 更新已销毁vpc下的云主机
func (h *HostSyncor) updateDestroyedHosts(conditon mapstr.MapStr) error {
	input := &metadata.UpdateOption{}
	input.Condition = conditon
	input.Data = mapstr.MapStr{
		common.BKHostInnerIPField:     "",
		common.BKHostOuterIPField:     "",
		common.BKCloudHostStatusField: metadata.CloudHostStatusIDs["stopped"],
		// 必须有该标识，用来跳过内网ip字段为空的校验
		common.IsDestroyedCloudHost: true,
	}

	uResult, err := h.logics.CoreAPI.CoreService().Instance().UpdateInstance(context.Background(), header, common.BKInnerObjIDHost, input)
	// 没有主机数据要更新并不算错误
	if err != nil  && err.Error() != kit.CCError.CCError(common.CCErrCommNotFound).Error() {
		blog.Errorf("updateDestroyedHosts fail,err:%s, input:%+v", err.Error(), *input)
		return err
	}
	if !uResult.Result && uResult.ErrMsg != kit.CCError.CCError(common.CCErrCommNotFound).Error() {
		blog.Errorf("updateDestroyedHosts fail,err:%s, input:%+v", uResult.ErrMsg, *input)
		return uResult.CCError()
	}

	return nil
}

// 更新被销毁vpc对应的云区域状态为异常
func (h *HostSyncor) updateDestroyedCloudArea(cloudIDs []int64) error {
	input := &metadata.UpdateOption{}
	input.Condition = mapstr.MapStr{common.BKCloudIDField: map[string]interface{}{
		common.BKDBIN: cloudIDs,
	}}
	input.Data = mapstr.MapStr{
		common.BKStatus: metadata.CloudAreaStatusIDs["abnormal"],
	}

	uResult, err := h.logics.CoreAPI.CoreService().Instance().UpdateInstance(context.Background(), header, common.BKInnerObjIDPlat, input)
	if err != nil {
		blog.Errorf("updateDestroyedCloudArea fail,err:%s, input:%+v", err.Error(), *input)
		return err
	}
	if !uResult.Result {
		blog.Errorf("updateDestroyedCloudArea fail,err:%s, input:%+v", uResult.ErrMsg, *input)
		return uResult.CCError()
	}

	return nil
}

// 更新同步任务里的vpc状态为被销毁
func (h *HostSyncor) updateDestroyedTaskVpc(taskID int64, vpcs map[string]bool) error {
	opt := &metadata.SearchSyncTaskOption{
		SearchCloudOption: metadata.SearchCloudOption{Condition: mapstr.MapStr{common.BKCloudSyncTaskID: taskID}},
	}

	ret, err := h.logics.SearchSyncTask(kit, opt)
	if err != nil {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, opt: err: %s", taskID, opt, err.Error())
		return err
	}
	if len(ret.Info) == 0 {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, opt: err: %s", taskID, opt, "no task is found")
		return fmt.Errorf("no task is found")
	}

	syncInfo := ret.Info[0].SyncVpcs
	for i, info := range syncInfo {
		if vpcs[info.VpcID] {
			syncInfo[i].Destroyed = true
		}
	}

	option := map[string]interface{}{common.BKCloudSyncVpcs: syncInfo}
	if err := h.logics.UpdateSyncTask(kit, taskID, option); err != nil {
		blog.Errorf("updateDestroyedTaskVpc failed, taskID: %v, err: %s", taskID, err.Error())
		return err
	}

	return nil
}

// 更新任务同步状态
func (h *HostSyncor) updateTaskState(taskid int64, status string, syncStatusDesc *metadata.SyncStatusDesc) error {
	option := mapstr.MapStr{common.BKCloudSyncStatus: status}
	if status == metadata.CloudSyncSuccess || status == metadata.CloudSyncFail {
		ts := time.Now()
		option.Set(common.BKCloudLastSyncTime, &ts)
		option.Set(common.BKCloudSyncStatusDescription, syncStatusDesc)
	}

	if err := h.logics.UpdateSyncTask(kit, taskid, option); err != nil {
		blog.Errorf("UpdateSyncTask failed, taskid: %v, err: %s", taskid, err.Error())
		return err
	}

	return nil
}
