package y3_9_202107271940

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

const step uint64 = 10000

// deleteTaskInvalidHistory delete task which settemplate or set has been removed
func deleteTaskInvalidHistory(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := mapstr.New()
	cond.Set(common.BKSetTemplateIDField, map[string]interface{}{common.BKDBNE: 0})

	setIDs, err := db.Table(common.BKTableNameBaseSet).Distinct(ctx, common.BKSetIDField, cond)
	if err != nil {
		blog.Errorf("get set id failed, tablename: %s, err: %v", common.BKTableNameBaseSet, err)
		return err
	}

	flags := []string{"set_template_sync"}
	for _, setID := range setIDs {
		flags = append(flags, fmt.Sprintf("set_template_sync:%d", setID))
	}

	deleteCond := mapstr.New()
	deleteCond.Set("flag", map[string]interface{}{common.BKDBNIN: flags})
	flag := 0
	for {
		flag += 1
		blog.Info("delete cc_APITask's task, entering the %d cycle", flag)

		task := make([]map[string]interface{}, 0)
		err := db.Table(common.BKTableNameAPITask).Find(deleteCond).Fields(common.BKTaskIDField).
			Start(0).Limit(step).All(ctx, &task)

		if err != nil {
			blog.Errorf("get apitask id to delete failed, tablename: %s, err: %v", common.BKTableNameAPITask, err)
			return err
		}

		if len(task) == 0 {
			break
		}

		taskIDs := make([]interface{}, len(task))
		for index, taskID := range task {
			taskIDfield, exist := taskID[common.BKTaskIDField]
			if !exist {
				blog.Errorf("get api task id failed, task: %+v, err: %v", taskID, err)
				return err
			}
			taskIDs[index] = taskIDfield
		}

		filter := map[string]interface{}{
			common.BKTaskIDField: map[string]interface{}{common.BKDBIN: taskIDs},
		}

		err = db.Table(common.BKTableNameAPITask).Delete(ctx, filter)
		if err != nil {
			blog.Errorf("delete invalid task history failed, tablename: %s, err: %v", common.BKTableNameAPITask, err)
			return err
		}
	}

	blog.Info("delete cc_APITask's task success")

	for _, flag := range flags {
		if flag == "set_template_sync" {
			continue
		}

		split := strings.Split(flag, ":")
		if len(split) != 2 {
			continue
		}
		instID, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			blog.Errorf("change string to int64 failed, value: %s, type: %T, err: %s",
				split[1], split[1], err.Error())
			return err
		}
		filter := map[string]interface{}{"flag": flag}
		flagDoc := mapstr.MapStr{
			"flag":       split[0],
			"bk_inst_id": instID,
		}
		if err = db.Table(common.BKTableNameAPITask).Update(ctx, filter, flagDoc); err != nil {
			blog.Errorf("upsert task flag failed, tablename: %s, err: %s", common.BKTableNameAPITask, err.Error())
			return err
		}
	}

	return nil
}
