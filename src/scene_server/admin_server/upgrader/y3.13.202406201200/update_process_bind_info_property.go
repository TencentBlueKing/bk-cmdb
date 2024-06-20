/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_13_202406201200

import (
	"context"
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func updateProcessAttribute(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	attrFilter := mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: common.BKProcBindInfo,
	}

	var attr metadata.Attribute
	err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).One(ctx, &attr)
	if err != nil {
		blog.Errorf("get exist cloud region and zone attribute failed, err: %v, filter: %+v", err, attrFilter)
		return err
	}
	if attr.ID == 0 {
		blog.Errorf("get process bind info attribute failed, err: attribute data does not exist , filter: %v",
			attrFilter)
		return err
	}

	var option []metadata.Attribute
	optionByte, err := json.Marshal(attr.Option)
	if err != nil {
		blog.Errorf("marshal table type option failed, err: %v, option: %v, rid: %s, ", err, attr.Option)
		return err
	}
	if err := json.Unmarshal(optionByte, &option); err != nil {
		blog.Errorf("unmarshal table type option failed, err: %v, option: %v, rid: %s, ", err, attr.Option)
		return err
	}

	for index, opt := range option {
		switch opt.PropertyID {
		case common.BKIP:
			option[index].Placeholder = "进程绑定的IP"
		case common.BKPort:
			option[index].Placeholder = "进程监听的端口</br> (填写示例)单端口: 8080</br> (填写示例)端口段: 8080-9090"
		case string(common.ProtocolName):
			option[index].Placeholder = "进程使用的协议"
		case common.BKEnable:
			option[index].Placeholder = "是否开启端口监控"
		default:
			continue
		}
	}

	attr.Option = option
	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, map[string]int64{common.BKFieldID: attr.ID}, attr)
	if err != nil {
		blog.Errorf("update process bind info attribute failed, err: %v, attr: %v", err, attr)
		return err
	}

	return nil
}
