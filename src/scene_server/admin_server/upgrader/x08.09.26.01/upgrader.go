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

package x08_09_26_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

type inputDatas struct {
	data mapstr.MapStr
	cond mapstr.MapStr
}

// dataWithConditon the key is the data, the value is conditon
var dataWithCondition = []inputDatas{
	inputDatas{
		data: mapstr.MapStr{
			metadata.AttributeFieldPropertyIndex: -3,
			metadata.AttributeFieldPlaceHolder:   `对外显示的服务名</br> 比如程序的二进制名称为java的服务zookeeper，则填zookeeper`,
		},
		cond: mapstr.MapStr{
			metadata.AttributeFieldPropertyID: common.BKProcessNameField,
			metadata.AttributeFieldObjectID:   common.BKInnerObjIDProc,
		},
	},
	inputDatas{
		data: mapstr.MapStr{
			metadata.AttributeFieldPropertyIndex: -2,
			metadata.AttributeFieldIsRequired:    true,
			metadata.AttributeFieldPropertyGroup: mCommon.BaseInfo,
			metadata.AttributeFieldPlaceHolder:   `程序的二进制名称</br> 比如zookeeper的二进制名称是java，则填java`,
		},
		cond: mapstr.MapStr{
			metadata.AttributeFieldPropertyID: "bk_func_name",
			metadata.AttributeFieldObjectID:   common.BKInnerObjIDProc,
		},
	},

	inputDatas{
		data: mapstr.MapStr{
			metadata.AttributeFieldPropertyIndex: -1,
			metadata.AttributeFieldPlaceHolder:   `程序启动参数</br> 唯一识别一个进程，比如zookeeper的启动参数包含 zookeeper`,
		},
		cond: mapstr.MapStr{
			metadata.AttributeFieldPropertyID: "bk_start_param_regex",
			metadata.AttributeFieldObjectID:   common.BKInnerObjIDProc,
		},
	},
}

func updateProcessTooltips(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {

	for _, input := range dataWithCondition {
		if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, input.cond, input.data); nil != err {
			return err
		}
	}

	return nil
}
