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
	"sort"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CompareFieldTemplateUnique compare field template uniques with object uniques
// @param forUI: defines if all comparison detail is needed for ui display, or only returns basic info for backend sync
func (t *template) CompareFieldTemplateUnique(kit *rest.Kit, opt *metadata.CompareFieldTmplUniqueOption, forUI bool) (
	*metadata.CompareFieldTmplUniquesRes, error) {

	// check compare options that is not related to object unique
	objID, err := t.comparator.getObjIDAndValidate(kit, opt.ObjectID)
	if err != nil {
		return nil, err
	}

	compParams, err := t.comparator.preCheckUnique(kit, opt.Uniques)
	if err != nil {
		return nil, err
	}

	// get object uniques
	objUniqueOpt := metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKObjIDField: objID},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Fields:    []string{common.BKFieldID, common.BKTemplateID, common.BKObjectUniqueKeys},
	}
	objUniqueRes, err := t.clientSet.CoreService().Model().ReadModelAttrUnique(kit.Ctx, kit.Header, objUniqueOpt)
	if err != nil {
		blog.Errorf("get object uniques failed, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		return nil, err
	}

	// if object has no uniques, add all field template uniques
	if len(objUniqueRes.Info) == 0 {
		createRes := make([]metadata.CompareOneFieldTmplUniqueRes, len(opt.Uniques))
		for idx := range opt.Uniques {
			createRes[idx] = metadata.CompareOneFieldTmplUniqueRes{
				Index: idx,
			}
		}
		return &metadata.CompareFieldTmplUniquesRes{Create: createRes}, nil
	}

	attrIDs := make([]uint64, len(objUniqueRes.Info))
	for _, unique := range objUniqueRes.Info {
		for _, key := range unique.Keys {
			attrIDs = append(attrIDs, key.ID)
		}
	}

	// get object uniques related attribute info
	objAttrOpt := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: attrIDs}},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Fields:    []string{common.BKFieldID, common.BKPropertyIDField},
	}
	util.AddModelBizIDCondition(objAttrOpt.Condition, 0)

	objAttrRes, err := t.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, objAttrOpt)
	if err != nil {
		blog.Errorf("get object uniques related attributes failed, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		return nil, err
	}

	if len(objAttrRes.Info) != len(attrIDs) {
		blog.Errorf("object uniques related attributes length is invalid, ids: %+v, rid: %s", attrIDs, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAttributeIDField)
	}

	for _, attribute := range objAttrRes.Info {
		compParams.attrMap[attribute.ID] = attribute.PropertyID
	}

	// cross-compare object uniques with template uniques
	if forUI {
		return t.comparator.compareUniqueForUI(kit, compParams, objUniqueRes.Info)
	}
	return t.comparator.compareUniqueForBackend(kit, compParams, objUniqueRes.Info)
}

type compUniqueParams struct {
	tmplIDMap map[int64]metadata.FieldTmplUniqueForUpdate
	tmpKeyMap map[string]metadata.FieldTmplUniqueForUpdate
	// tmplPropIDMap field template key property id to unique map, used for conflict check
	tmplPropIDMap map[string][]metadata.FieldTmplUniqueForUpdate
	// tmplIndexMap field template attribute property id to original index map, used to trace back to original unique
	tmplIndexMap map[string]int
	// createTmplMap to be created(has no matching/conflict unique) field template unique key map
	createTmplMap map[string]struct{}
	// attrMap unique keys' attribute id to property id map
	attrMap map[int64]string
}

func (c *comparator) preCheckUnique(kit *rest.Kit, uniques []metadata.FieldTmplUniqueForUpdate) (*compUniqueParams,
	error) {

	params := &compUniqueParams{
		tmplIDMap:     make(map[int64]metadata.FieldTmplUniqueForUpdate),
		tmpKeyMap:     make(map[string]metadata.FieldTmplUniqueForUpdate),
		tmplPropIDMap: make(map[string][]metadata.FieldTmplUniqueForUpdate),
		tmplIndexMap:  make(map[string]int),
		createTmplMap: make(map[string]struct{}),
		attrMap:       make(map[int64]string),
	}

	for _, unique := range uniques {
		if len(unique.Keys) == 0 {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjectUniqueKeys)
		}
	}

	// check if field template uniques conflicts with each other, and categorize field template uniques
	for index, unique := range uniques {
		compKey := c.genUniqueKey(unique.Keys)

		if _, exists := params.tmpKeyMap[compKey]; exists {
			blog.Errorf("template unique(key: %s) duplicate, rid: %s", compKey, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "uniques")
		}

		for _, key := range unique.Keys {
			// check if unique keys conflicts with another unique that has the same attribute
			for _, sameAttrUnique := range params.tmplPropIDMap[key] {
				isConflict := c.checkUniqueConflict(unique.Keys, sameAttrUnique.Keys)
				if isConflict {
					blog.Errorf("template unique( %+v and %+v) duplicate, rid: %s", unique, sameAttrUnique, kit.Rid)
					return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "uniques")
				}
			}

			params.tmplPropIDMap[key] = append(params.tmplPropIDMap[key], unique)
		}

		if unique.ID != 0 {
			params.tmplIDMap[unique.ID] = unique
		}

		params.tmpKeyMap[compKey] = unique
		params.tmplIndexMap[compKey] = index
		params.createTmplMap[compKey] = struct{}{}
	}

	return params, nil
}

// genUniqueKey generate key for comparison by unique keys
func (c *comparator) genUniqueKey(keys []string) string {
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

// checkUniqueConflict check if keys are subsets/supersets of each other.
func (c *comparator) checkUniqueConflict(keys1, keys2 []string) bool {
	keyMap := make(map[string]struct{})
	for _, key := range keys1 {
		keyMap[key] = struct{}{}
	}

	isConflict := true
	for _, key := range keys2 {
		if _, exists := keyMap[key]; !exists {
			isConflict = false
		}

		delete(keyMap, key)
	}

	// all keys in keys2 matches keys1
	if isConflict {
		return true
	}

	// all keys in keys1 matches keys2
	if len(keyMap) == 0 {
		return true
	}

	return false
}

type objUniqueForComp struct {
	unique  metadata.ObjectUnique
	compKey string
	keys    []string
}

func (c *comparator) compareUniqueForUI(kit *rest.Kit, params *compUniqueParams, uniques []metadata.ObjectUnique) (
	*metadata.CompareFieldTmplUniquesRes, error) {

	res := new(metadata.CompareFieldTmplUniquesRes)

	// compare object unique with its template
	noTmplUnique := make([]objUniqueForComp, 0)
	for idx := range uniques {
		unique := uniques[idx]

		// convert db unique keys to unique keys for template comparison
		uniqueKeys := make([]string, len(unique.Keys))
		for i, key := range unique.Keys {
			uniqueKeys[i] = params.attrMap[int64(key.ID)]
		}
		compKey := c.genUniqueKey(uniqueKeys)
		compUnique := objUniqueForComp{
			unique:  unique,
			compKey: compKey,
			keys:    uniqueKeys,
		}

		// compare unique without template later, because template id has maximum priority in comparison
		if unique.TemplateID == 0 {
			noTmplUnique = append(noTmplUnique, compUnique)
			continue
		}

		tmplUnique, exists := params.tmplIDMap[unique.TemplateID]
		if !exists {
			// unique template is not exist, check if the unique conflicts with other templates
			if conflictKey, isConflict := c.checkObjUniqueConflict(params, &compUnique); isConflict {
				delete(params.createTmplMap, conflictKey)

				res.Conflict = append(res.Conflict, metadata.CompareOneFieldTmplUniqueRes{
					Index: params.tmplIndexMap[conflictKey],
					Message: kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateUniqueConflict,
						compUnique.unique.ID, conflictKey).Error(),
					Data: &compUnique.unique,
				})
				continue
			}

			res.Unchanged = append(res.Unchanged, unique)
			continue
		}

		// compare the unique with its template, check if their keys are the same
		if isChanged := c.compareOneUniqueInfo(params, &compUnique, &tmplUnique, res); !isChanged {
			res.Unchanged = append(res.Unchanged, unique)
		}
	}

	// compare object uniques without template
	for idx := range noTmplUnique {
		compUnique := noTmplUnique[idx]

		tmplUnique, exists := params.tmpKeyMap[compUnique.compKey]
		if !exists {
			// unique is not related to template, check if its keys conflict with all templates
			if conflictKey, isConflict := c.checkObjUniqueConflict(params, &compUnique); isConflict {
				delete(params.createTmplMap, conflictKey)

				res.Conflict = append(res.Conflict, metadata.CompareOneFieldTmplUniqueRes{
					Index: params.tmplIndexMap[conflictKey],
					Message: kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateUniqueConflict,
						compUnique.unique.ID, conflictKey).Error(),
					Data: &compUnique.unique,
				})
				continue
			}

			res.Unchanged = append(res.Unchanged, compUnique.unique)
			continue
		}

		// compare the unique with its template, check if their keys are the same
		if isChanged := c.compareOneUniqueInfo(params, &compUnique, &tmplUnique, res); !isChanged {
			res.Unchanged = append(res.Unchanged, compUnique.unique)
		}
	}

	// field template unique with no matching object uniques should be created
	for compKey := range params.createTmplMap {
		res.Create = append(res.Create, metadata.CompareOneFieldTmplUniqueRes{
			Index: params.tmplIndexMap[compKey],
		})
	}

	return res, nil
}

func (c *comparator) compareUniqueForBackend(kit *rest.Kit, params *compUniqueParams, uniques []metadata.ObjectUnique) (
	*metadata.CompareFieldTmplUniquesRes, error) {

	res := new(metadata.CompareFieldTmplUniquesRes)

	// compare object unique with its template
	noTmplUnique := make([]objUniqueForComp, 0)
	for idx := range uniques {
		unique := uniques[idx]

		// convert db unique keys to unique keys for template comparison
		uniqueKeys := make([]string, len(unique.Keys))
		for i, key := range unique.Keys {
			uniqueKeys[i] = params.attrMap[int64(key.ID)]
		}
		compKey := c.genUniqueKey(uniqueKeys)
		compUnique := objUniqueForComp{
			unique:  unique,
			compKey: compKey,
			keys:    uniqueKeys,
		}

		// compare unique without template later, because template id has maximum priority in comparison
		if unique.TemplateID == 0 {
			noTmplUnique = append(noTmplUnique, compUnique)
			continue
		}

		tmplUnique, exists := params.tmplIDMap[unique.TemplateID]
		if !exists {
			// unique template is not exist, check if the unique conflicts with other templates
			if conflictKey, isConflict := c.checkObjUniqueConflict(params, &compUnique); isConflict {
				return nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateUniqueConflict, compUnique.unique.ID,
					conflictKey)
			}

			// unique template is deleted, so we update unique template id to zero
			res.Update = append(res.Update, metadata.CompareOneFieldTmplUniqueRes{
				Data: &unique,
			})
			continue
		}

		// compare the unique with its template, check if their keys are the same
		c.compareOneUniqueInfo(params, &compUnique, &tmplUnique, res)
	}

	// compare object uniques without template
	for idx := range noTmplUnique {
		compUnique := noTmplUnique[idx]

		tmplUnique, exists := params.tmpKeyMap[compUnique.compKey]
		if !exists {
			// unique is not related to template, check if its keys conflict with all templates
			if conflictKey, isConflict := c.checkObjUniqueConflict(params, &compUnique); isConflict {
				return nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateUniqueConflict, compUnique.unique.ID,
					conflictKey)
			}
			continue
		}

		// compare the unique with its template, check if their keys are the same
		c.compareOneUniqueInfo(params, &compUnique, &tmplUnique, res)
	}

	// field template unique with no matching object uniques should be created
	for compKey := range params.createTmplMap {
		res.Create = append(res.Create, metadata.CompareOneFieldTmplUniqueRes{
			Index: params.tmplIndexMap[compKey],
		})
	}

	return res, nil
}

// checkObjUniqueConflict check if object unique conflicts with field template uniques
func (c *comparator) checkObjUniqueConflict(params *compUniqueParams, objCompUnique *objUniqueForComp) (string, bool) {
	if _, exists := params.tmpKeyMap[objCompUnique.compKey]; exists {
		return objCompUnique.compKey, true
	}

	for _, key := range objCompUnique.keys {
		// check if object unique keys conflicts with field template uniques that has the same attribute
		for _, tmplUnique := range params.tmplPropIDMap[key] {
			isConflict := c.checkUniqueConflict(objCompUnique.keys, tmplUnique.Keys)
			if !isConflict {
				continue
			}

			conflictKey := c.genUniqueKey(tmplUnique.Keys)
			return conflictKey, true
		}
	}

	return "", false
}

// compareOneUniqueInfo compare one field template and object unique, add update unique info into compare result
func (c *comparator) compareOneUniqueInfo(params *compUniqueParams, objCompUnique *objUniqueForComp,
	tmplUnique *metadata.FieldTmplUniqueForUpdate, res *metadata.CompareFieldTmplUniquesRes) bool {

	tmplCompKey := c.genUniqueKey(tmplUnique.Keys)

	// remove template unique that has matching obj unique, because it can't be related to another one
	delete(params.tmplIDMap, tmplUnique.ID)
	delete(params.createTmplMap, tmplCompKey)

	for _, key := range tmplUnique.Keys {
		delete(params.tmplPropIDMap, key)
	}

	// check if unique keys are the same
	if objCompUnique.compKey == tmplCompKey {
		return false
	}

	res.Update = append(res.Update, metadata.CompareOneFieldTmplUniqueRes{
		Index: params.tmplIndexMap[tmplCompKey],
		Data:  &objCompUnique.unique,
	})
	return true
}
