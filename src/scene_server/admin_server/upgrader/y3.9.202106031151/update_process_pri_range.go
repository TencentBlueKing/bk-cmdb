package y3_9_202106031151

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"context"
)

func updatePriorityProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: "priority",
		common.BKAppIDField:      0,
	}
	doc := map[string]interface{}{
		"option": map[string]interface{}{
			"min": common.MinProcessPrio,
			"max": common.MaxProcessPrio,
		},
		"placeholder": "批量启动进程依据优先级从小到大排序操作，停止进程按反序操作。优先级数值范围支持输入 [ -100 ~ 10000 ]",
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("update failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
		return err
	}

	return nil
}
