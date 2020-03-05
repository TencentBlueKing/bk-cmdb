package cloudsync

import (
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

// 处理云主机同步
func (t *taskProcessor) SyncCloudHost(task *metadata.CloudSyncTask) error {
	// 根据账号id获取账号详情
	account, err := t.getAccountDetail(task.TaskID)
	if err != nil {
		return err
	}

	// 根据任务详情和账号信息获取要同步的云主机资源
	hostResource, err := t.getCloudHostResource(task, account)
	if err != nil {
		return err
	}

	// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
	diffHosts, err := t.getDiffHosts(hostResource)
	if err != nil {
		return err
	}

	// 没差异则结束
	if len(diffHosts) == 0 {
		blog.V(3).Infof("no diff hosts for taskid:%d", task.TaskID)
		return nil
	}

	// 有差异的更新任务同步状态为同步中
	err = t.updateTaskState(task.TaskID, metadata.InSync)
	if err != nil {
		return err
	}

	// 同步有差异的主机数据
	syncResult, err := t.syncDiffHosts(diffHosts)
	if err != nil {
		return err
	}

	// 增加任务同步历史记录
	err = t.addHostSyncHistory(syncResult)
	if err != nil {
		return err
	}

	// 完成后更新任务同步状态为成功或失败
	err = t.updateTaskState(task.TaskID, metadata.Success)
	if err != nil {
		return err
	}

	return nil
}

// 根据任务详情和账号信息获取要同步的云主机资源
func (t *taskProcessor) getCloudHostResource(task *metadata.CloudSyncTask, account *metadata.CloudAccount) (*metadata.CloudHostResource, error) {
	return nil, nil
}

// 查询vpc对应的云区域，没有则创建,并更新云主机资源信息里的云区域id
func (t *taskProcessor) updateCLoudId(hostResource *metadata.CloudHostResource) (*metadata.CloudHostResource, error) {
	return nil, nil
}

// 根据主机实例id获取mongo中的主机信息,并获取有差异的主机
func (t *taskProcessor) getDiffHosts(hostResource *metadata.CloudHostResource) (map[string][]*metadata.Instance, error) {
	return nil, nil
}

// 更新任务同步状态
func (t *taskProcessor) updateTaskState(taskid int64, state int) error {
	return nil
}

// 同步有差异的主机数据
func (t *taskProcessor) syncDiffHosts(diffhosts map[string][]*metadata.Instance) (*metadata.SyncHostsResult, error) {
	return nil, nil
}

// 增加任务同步历史记录
func (t *taskProcessor) addHostSyncHistory(*metadata.SyncHostsResult) error {
	return nil
}
