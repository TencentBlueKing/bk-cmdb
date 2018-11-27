package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"context"
)

func (lgc *Logics) CreateCloudTask(ctx context.Context, input interface{}) error {
	objID, err := lgc.Instance.NextSequence(ctx, common.BKTableNameCloudTask)
	if err != nil {
		return err
	}

	inputc := input.(map[string]interface{})
	inputc["bk_task_id"] = objID
	if err := lgc.Instance.Table(common.BKTableNameCloudTask).Insert(ctx, inputc); err != nil {
		return err
	}

	return nil
}

func (lgc *Logics) CreateResourceConfirm(ctx context.Context, input interface{}) error {
	objID, err := lgc.Instance.NextSequence(ctx, common.BKTableNameCloudResourceSync)
	if err != nil {
		return err
	}

	inputc := input.(map[string]interface{})
	inputc["bk_resource_id"] = objID
	if err := lgc.Instance.Table(common.BKTableNameCloudResourceSync).Insert(ctx, inputc); err != nil {
		return err
	}

	blog.Info("CreateResourceSync table bk_CloudResourceSync")
	return nil
}

func (lgc *Logics) CreateCloudHistory(ctx context.Context, input interface{}) error {
	objID, err := lgc.Instance.NextSequence(ctx, common.BKTableNameCloudHistory)
	if err != nil {
		return err
	}

	inputc := input.(map[string]interface{})
	inputc["bk_history_id"] = objID
	if err := lgc.Instance.Table(common.BKTableNameCloudHistory).Insert(ctx, inputc); err != nil {
		return err
	}

	return nil
}
