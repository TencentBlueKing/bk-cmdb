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
	"sync"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/logics/model"
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

func (t *template) getTemplateAttrByID(kit *rest.Kit, id int64, fields []string) ([]metadata.FieldTemplateAttr,
	errors.CCErrorCoder) {

	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, id),
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: fields,
	}

	// list field template attributes
	res, err := t.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list template attributes failed, template id: %d, err: %v, rid: %s", id, err, kit.Rid)
		return nil, err
	}

	if len(res.Info) == 0 {
		blog.Errorf("no template attributes founded, template id: %d, rid: %s", id, kit.Rid)
		return []metadata.FieldTemplateAttr{}, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKTemplateID)
	}

	return res.Info, nil
}

func (t *template) getAttrSyncStatus(kit *rest.Kit, templateID, objectID int64) (
	*metadata.ListFieldTmpltSyncStatusResult, error) {

	res, ccErr := t.getTemplateAttrByID(kit, templateID, []string{})
	if ccErr != nil {
		return nil, ccErr
	}

	// verify the validity of field attributes and obtain
	// the contents of template fields that need to be synchronized
	opt := &metadata.CompareFieldTmplAttrOption{
		TemplateID: templateID,
		ObjectID:   objectID,
		Attrs:      res,
		IsPartial:  true,
	}

	_, attrRes, err := t.CompareFieldTemplateAttr(kit, opt, false)
	if err != nil {
		blog.Errorf("compare field template failed, cond: %+v, err: %v, rid: %s", opt, err, kit.Rid)
		return nil, err
	}
	return attrRes, nil
}

func (t *template) getUniqueSyncStatus(kit *rest.Kit, templateID, objectID int64) (
	*metadata.ListFieldTmpltSyncStatusResult, error) {

	uniqueFilter := filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, templateID)
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: uniqueFilter},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	// list field template uniques
	res, ccErr := t.clientSet.CoreService().FieldTemplate().ListFieldTemplateUnique(kit.Ctx, kit.Header, listOpt)
	if ccErr != nil {
		blog.Errorf("list field template uniques failed, err: %v, id: %+v, rid: %s", ccErr, templateID, kit.Rid)
		return nil, ccErr
	}

	uniqueOp := &metadata.CompareFieldTmplUniqueOption{
		TemplateID: templateID,
		ObjectID:   objectID,
		IsPartial:  true,
		Uniques:    make([]metadata.FieldTmplUniqueForUpdate, len(res.Info)),
	}

	attrs, ccErr := t.getTemplateAttrByID(kit, templateID, []string{common.BKFieldID, common.BKPropertyIDField})
	if ccErr != nil {
		return nil, ccErr
	}

	attrIDProMap := make(map[int64]string)
	for _, attr := range attrs {
		attrIDProMap[attr.ID] = attr.PropertyID
	}

	propertyIDs := make([]string, 0)
	for index := range res.Info {
		for _, key := range res.Info[index].Keys {
			propertyID, ok := attrIDProMap[key]
			if !ok {
				blog.Errorf("property id not found, attr id: %s, object id: %d, rid: %s", key, objectID, kit.Rid)
				return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound, common.BKPropertyIDField)
			}
			propertyIDs = append(propertyIDs, propertyID)
		}
		uniqueOp.Uniques[index].Keys = propertyIDs
		uniqueOp.Uniques[index].ID = res.Info[index].ID
		propertyIDs = []string{}
	}

	_, uniqueRes, err := t.CompareFieldTemplateUnique(kit, uniqueOp, false)
	if err != nil {
		blog.Errorf("get field template unique failed, cond: %+v, err: %v, rid: %s", uniqueOp, err, kit.Rid)
		return nil, err
	}
	return uniqueRes, nil
}

// ListFieldTemplateSyncStatus get the diff status of templates and models
func (t *template) ListFieldTemplateSyncStatus(kit *rest.Kit, option *metadata.ListFieldTmpltSyncStatusOption) (
	[]metadata.ListFieldTmpltSyncStatusResult, error) {

	// check whether the corresponding relationship between object and templateID is legal
	objIDMap := make(map[int64]struct{})
	for _, id := range option.ObjectIDs {
		objIDMap[id] = struct{}{}
	}

	objIDs := make([]int64, 0)
	for id := range objIDMap {
		objIDs = append(objIDs, id)
	}

	cond := filtertools.GenAtomFilter(common.ObjectIDField, filter.In, objIDs)
	tmplFilter, err := filtertools.And(cond, filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, option.ID))
	if err != nil {
		blog.Errorf("gen filter failed, objIDs: %v, template id: %d, err: %v, rid: %s", objIDs, option.ID, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "template_filter")
	}

	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: tmplFilter},
		Page:               metadata.BasePage{EnableCount: true},
	}

	res, err := t.clientSet.CoreService().FieldTemplate().ListObjFieldTmplRel(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, option, kit.Rid)
		return nil, err
	}
	if len(objIDs) != int(res.Count) {
		blog.Errorf("the number of associations obtained is incorrect, opt: %+v, rid: %s", option, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.ObjectIDField)
	}

	var wg sync.WaitGroup
	var firstErr error
	pipeline := make(chan bool, 5)

	result := make([]metadata.ListFieldTmpltSyncStatusResult, 0)
	var lock sync.Mutex

	// here, the concurrency is performed according to the objectID,
	// and the concurrency internally compares the attributes first
	// in order, and then compares the unique check
	for _, objectID := range option.ObjectIDs {

		pipeline <- true
		wg.Add(1)
		go func(id, objectID int64) {
			defer func() {
				wg.Done()
				<-pipeline
				lock.Unlock()
			}()

			attrStatus, err := t.getAttrSyncStatus(kit, id, objectID)
			if err != nil {
				firstErr = err
			}
			lock.Lock()
			// if a difference in attributes has already been identified,
			// there is no need to continue to calculate whether there is
			// a difference in the unique check
			if attrStatus.NeedSync {
				result = append(result, *attrStatus)
				return
			}

			uniqueStatus, err := t.getUniqueSyncStatus(kit, id, objectID)
			if err != nil {
				firstErr = err
			}
			result = append(result, *uniqueStatus)

		}(option.ID, objectID)
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return result, nil
}
