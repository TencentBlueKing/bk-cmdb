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

package fieldtmpl

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/logics/model"
	"sync"
)

type comparator struct {
	clientSet apimachinery.ClientSetInterface
	asst      model.AssociationOperationInterface
}

// getObjIDAndValidate validate object, do not allow field template to bind mainline object(except for host)
func (c *comparator) getObjIDAndValidate(kit *rest.Kit, objectID int64) (string, error) {
	objCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKFieldID: objectID},
		Fields:    []string{common.BKObjIDField},
	}

	objRes, err := c.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, objCond)
	if err != nil {
		blog.Errorf("get object by id %d failed, err: %v, rid: %s", objectID, err, kit.Rid)
		return "", err
	}

	if len(objRes.Info) != 1 {
		blog.Errorf("object with id %d count is invalid, res: %+v, rid: %s", objectID, objRes.Info, kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.ObjectIDField)
	}

	objID := objRes.Info[0].ObjectID

	if objID == common.BKInnerObjIDHost {
		return objID, nil
	}

	isMainline, err := c.asst.IsMainlineObject(kit, objID)
	if err != nil {
		blog.Errorf("check if object %s is mainline object failed, err: %v, rid: %s", objID, err, kit.Rid)
		return "", err
	}

	if isMainline {
		blog.Errorf("object %s is mainline object, can not bind field template, rid: %s", objID, kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	return objID, nil
}

// ListFieldTemplateSyncStatus get the diff status of templates and models
func (t *template) ListFieldTemplateSyncStatus(kit *rest.Kit, option *metadata.ListFieldTmpltSyncStatusOption) (
	[]metadata.ListFieldTmpltSyncStatusResult, error) {

	result := make([]metadata.ListFieldTmpltSyncStatusResult, 0)

	var wg sync.WaitGroup
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 5)
	// 这里按照objectID进行并发，其中并发内部按照顺序首先对比属性，然后再对比唯一校验

	for _, objectID := range option.ObjectIDs {

		pipeline <- true
		wg.Add(1)
		go func(id, objectID int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			_, res, err := t.comparator.compareAttrForBackend(kit, compParams, id, objectID, objAttrRes.Info, true)
			if err != nil {
				return
			}
			if res.NeedSync {
				result = append(result, *res)
				return
			}

		}(option.ID, objectID)
	}

	return result, nil
}
