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

	"configcenter/pkg/filter"
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
	*metadata.CompareFieldTmplUniquesRes, *metadata.ListFieldTmpltSyncStatusResult, error) {

	// check compare options that is not related to object unique
	objID, err := t.comparator.getObjIDAndValidate(kit, opt.ObjectID)
	if err != nil {
		return nil, nil, err
	}

	compParams, err := t.comparator.preCheckUnique(kit, objID, opt.Uniques, forUI)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	// if object has no uniques, add all field template uniques
	if len(objUniqueRes.Info) == 0 {
		createRes := make([]metadata.CompareOneFieldTmplUniqueRes, len(opt.Uniques))
		for idx := range opt.Uniques {
			createRes[idx] = metadata.CompareOneFieldTmplUniqueRes{Index: idx}
		}
		if opt.IsPartial {
			result := &metadata.ListFieldTmpltSyncStatusResult{ObjectID: opt.ObjectID, NeedSync: true}
			return nil, result, nil
		}
		return &metadata.CompareFieldTmplUniquesRes{Create: createRes}, nil, nil
	}

	attrs, err := t.getObjAttrsByID(kit, objID, objUniqueRes.Info)
	if err != nil {
		return nil, nil, err
	}

	for _, attribute := range attrs {
		compParams.attrMap[attribute.ID] = attribute.PropertyID
	}

	// get object uniques related field template ids that does not belong to field template for comparison
	compParams.otherTmplIDMap, err = t.getOtherUniqueTmpl(kit, opt.TemplateID, objUniqueRes.Info)
	if err != nil {
		return nil, nil, err
	}

	// cross-compare object uniques with template uniques
	if forUI {
		result, err := t.comparator.compareUniqueForUI(kit, compParams, objUniqueRes.Info)
		if err != nil {
			return nil, nil, err
		}
		return result, nil, nil
	}

	if opt.IsPartial {
		_, statusResult, err := t.comparator.compareUniqueForBackend(kit, compParams, objUniqueRes.Info,
			opt.ObjectID, true)
		if err != nil {
			return nil, nil, err
		}
		return nil, statusResult, nil
	}

	uniquesRes, _, err := t.comparator.compareUniqueForBackend(kit, compParams, objUniqueRes.Info, opt.ObjectID, false)
	if err != nil {
		return nil, nil, err
	}
	return uniquesRes, nil, nil
}

func (t *template) getObjAttrsByID(kit *rest.Kit, objID string, uniques []metadata.ObjectUnique) (
	[]metadata.Attribute, error) {

	attrMap := make(map[uint64]struct{})
	for _, unique := range uniques {
		for _, key := range unique.Keys {
			attrMap[key.ID] = struct{}{}
		}
	}

	attrIDs := make([]uint64, 0)
	for id := range attrMap {
		if id == 0 {
			continue
		}
		attrIDs = append(attrIDs, id)
	}

	if len(attrIDs) == 0 {
		return []metadata.Attribute{}, nil
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
		blog.Errorf("get object uniques related attrs failed, opt: %+v,err: %v, rid: %s", objAttrOpt, err, kit.Rid)
		return nil, err
	}

	if len(objAttrRes.Info) != len(attrIDs) {
		blog.Errorf("object uniques related attributes length is invalid, ids: %+v, rid: %s", attrIDs, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAttributeIDField)
	}
	return objAttrRes.Info, nil
}

func (t *template) getOtherUniqueTmpl(kit *rest.Kit, templateID int64, uniques []metadata.ObjectUnique) (
	map[int64]struct{}, error) {

	otherTmplIDMap := make(map[int64]struct{})

	if len(uniques) == 0 {
		return otherTmplIDMap, nil
	}

	templateIDs := make([]int64, len(uniques))
	for idx, unique := range uniques {
		templateIDs[idx] = unique.TemplateID
	}

	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: &filter.Expression{
			RuleFactory: &filter.CombinedRule{
				Condition: filter.And,
				Rules: []filter.RuleFactory{
					&filter.AtomRule{
						Field:    common.BKFieldID,
						Operator: filter.In.Factory(),
						Value:    templateIDs,
					}, &filter.AtomRule{
						Field:    common.BKTemplateID,
						Operator: filter.NotEqual.Factory(),
						Value:    templateID,
					},
				},
			},
		}},
		Fields: []string{common.BKFieldID},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	res, err := t.clientSet.CoreService().FieldTemplate().ListFieldTemplateUnique(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field template uniques failed, err: %v, opt: %+v, rid: %s", err, listOpt, kit.Rid)
		return nil, err
	}

	for _, unique := range res.Info {
		otherTmplIDMap[unique.ID] = struct{}{}
	}

	return otherTmplIDMap, nil
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
	// the corresponding relationship between the template propertyID
	// and the self-incrementing ID of the model attribute, which is
	// used to change from single unique to joint unique in the subsequent
	// update unique verification scenario.
	tmplProToIDMap map[string]int64
	// otherTmplIDMap other field template's unique ids map, used to check if unique template is deleted or conflict
	otherTmplIDMap map[int64]struct{}
}

func (c *comparator) preCheckUnique(kit *rest.Kit, objID string, uniques []metadata.FieldTmplUniqueForUpdate,
	forUI bool) (*compUniqueParams, error) {

	params := &compUniqueParams{
		tmplIDMap:      make(map[int64]metadata.FieldTmplUniqueForUpdate),
		tmpKeyMap:      make(map[string]metadata.FieldTmplUniqueForUpdate),
		tmplPropIDMap:  make(map[string][]metadata.FieldTmplUniqueForUpdate),
		tmplIndexMap:   make(map[string]int),
		createTmplMap:  make(map[string]struct{}),
		attrMap:        make(map[int64]string),
		tmplProToIDMap: make(map[string]int64),
	}

	tmplPropertyIDMap := make(map[string]struct{})
	for _, unique := range uniques {
		if len(unique.Keys) == 0 {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjectUniqueKeys)
		}
		for _, key := range unique.Keys {
			tmplPropertyIDMap[key] = struct{}{}
		}
	}

	tmplPropertyIDs := make([]string, 0)
	for propertyID := range tmplPropertyIDMap {
		tmplPropertyIDs = append(tmplPropertyIDs, propertyID)
	}

	// obtain the auto-increment ID of the corresponding attribute of the model through the propertyID on
	// the field combination template. In the newly added scenario for subsequent new unique verification,
	// find the auto-increment ID of the corresponding model attribute according to the propertyID.
	objAttrOpt := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKPropertyIDField: mapstr.MapStr{common.BKDBIN: tmplPropertyIDs}},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Fields:    []string{common.BKFieldID, common.BKPropertyIDField},
	}
	util.AddModelBizIDCondition(objAttrOpt.Condition, 0)

	objAttrRes, err := c.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, objAttrOpt)
	if err != nil {
		blog.Errorf("get object attrs failed, opt: %+v, err: %v, rid: %s", objAttrOpt, err, kit.Rid)
		return nil, err
	}

	// when info is 0, it is a diff scene, and if it is greater than 0, it is a synchronous scene. If it is a
	// synchronous scene, the number of model attributes and the num of template attrs must be equal and greater
	// than 0 however, there may be some attributes that can be found for comparison scenarios.
	if len(objAttrRes.Info) != len(tmplPropertyIDs) && len(objAttrRes.Info) != 0 && !forUI {
		blog.Errorf("object attrs length is invalid, property ids: %+v, rid: %s", tmplPropertyIDs, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAttributeIDField)
	}

	// in the synchronization scenario, the attributes of the template have been synchronized to the model, and info
	// must be greater than 0 at this time. in the previous diff comparison scenario, the info field is 0 because the
	// attribute has not been synchronized to the model.
	if len(objAttrRes.Info) != 0 {
		for _, info := range objAttrRes.Info {
			params.tmplProToIDMap[info.PropertyID] = info.ID
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

	// compare object unique with its template
	res, noTmplUnique := new(metadata.CompareFieldTmplUniquesRes), make([]objUniqueForComp, 0)
	for idx := range uniques {
		unique := uniques[idx]

		// convert db unique keys to unique keys for template comparison
		uniqueKeys := make([]string, len(unique.Keys))
		for i, key := range unique.Keys {
			uniqueKeys[i] = params.attrMap[int64(key.ID)]
		}

		compUnique := objUniqueForComp{unique: unique, compKey: c.genUniqueKey(uniqueKeys), keys: uniqueKeys}

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
		isChanged, err := c.compareOneUniqueInfo(kit, params, &compUnique, &tmplUnique, res, false, true)
		if err != nil {
			return nil, err
		}
		if !isChanged {
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
		isChanged, err := c.compareOneUniqueInfo(kit, params, &compUnique, &tmplUnique, res, true, true)
		if err != nil {
			return nil, err
		}
		if !isChanged {
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

func (c *comparator) compareUniqueForBackend(kit *rest.Kit, params *compUniqueParams, uniques []metadata.ObjectUnique,
	objectID int64, isPartial bool) (*metadata.CompareFieldTmplUniquesRes, *metadata.ListFieldTmpltSyncStatusResult,
	error) {

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
				if isPartial {
					result := &metadata.ListFieldTmpltSyncStatusResult{ObjectID: objectID, NeedSync: true}
					return nil, result, nil
				}
				return nil, nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateUniqueConflict, compUnique.unique.ID,
					conflictKey)
			}

			// unique template belongs to other field template, do not update unique template id
			if _, isOther := params.otherTmplIDMap[unique.TemplateID]; isOther {
				continue
			}

			if isPartial {
				result := &metadata.ListFieldTmpltSyncStatusResult{ObjectID: objectID, NeedSync: true}
				return nil, result, nil
			}

			// unique template is deleted, so we update unique template id to -1
			res.Update = append(res.Update, metadata.CompareOneFieldTmplUniqueRes{
				Index: -1,
				Data:  &unique,
			})
			continue
		}

		// compare the unique with its template, check if their keys are the same
		isChanged, err := c.compareOneUniqueInfo(kit, params, &compUnique, &tmplUnique, res, false, false)
		if err != nil {
			return nil, nil, err
		}

		if isChanged && isPartial {
			result := &metadata.ListFieldTmpltSyncStatusResult{ObjectID: objectID, NeedSync: true}
			return nil, result, nil
		}
	}

	// compare object uniques without template
	isChanged, err := c.dealNoTmplUnique(kit, params, noTmplUnique, isPartial, res)
	if err != nil {
		return nil, nil, err
	}

	if isPartial {
		partialRes, err := c.dealUniquePartialResult(params, isChanged, objectID)
		if err != nil {
			return nil, nil, err
		}
		return nil, partialRes, nil
	}

	// field template unique with no matching object uniques should be created
	for compKey := range params.createTmplMap {
		res.Create = append(res.Create, metadata.CompareOneFieldTmplUniqueRes{
			Index: params.tmplIndexMap[compKey],
		})
	}
	return res, nil, nil
}

func (c *comparator) dealUniquePartialResult(params *compUniqueParams, isChanged bool, objectID int64) (
	*metadata.ListFieldTmpltSyncStatusResult, error) {

	result := &metadata.ListFieldTmpltSyncStatusResult{
		ObjectID: objectID,
	}

	if isChanged || len(params.createTmplMap) > 0 {
		result.NeedSync = true
		return result, nil
	}

	return result, nil
}

func (c *comparator) dealNoTmplUnique(kit *rest.Kit, params *compUniqueParams, noTmplUnique []objUniqueForComp,
	isPartial bool, res *metadata.CompareFieldTmplUniquesRes) (bool, error) {

	// compare object uniques without template
	for idx := range noTmplUnique {
		compUnique := noTmplUnique[idx]
		tmplUnique, exists := params.tmpKeyMap[compUnique.compKey]
		if !exists {
			// unique is not related to template, check if its keys conflict with all templates
			if conflictKey, isConflict := c.checkObjUniqueConflict(params, &compUnique); isConflict {
				if isPartial {
					return true, nil
				}
				return false, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateUniqueConflict,
					compUnique.unique.ID, conflictKey)
			}
			continue
		}

		// compare the unique with its template, check if their keys are the same
		isChanged, err := c.compareOneUniqueInfo(kit, params, &compUnique, &tmplUnique, res, true, false)
		if err != nil {
			return false, err
		}

		if isChanged && isPartial {
			return isChanged, nil
		}
	}
	return false, nil
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
func (c *comparator) compareOneUniqueInfo(kit *rest.Kit, params *compUniqueParams, objCompUnique *objUniqueForComp,
	tmplUnique *metadata.FieldTmplUniqueForUpdate, res *metadata.CompareFieldTmplUniquesRes, isNoTmpl bool,
	forUI bool) (bool, error) {

	tmplCompKey := c.genUniqueKey(tmplUnique.Keys)

	// remove template unique that has matching obj unique, because it can't be related to another one
	delete(params.tmplIDMap, tmplUnique.ID)
	delete(params.createTmplMap, tmplCompKey)

	for _, key := range tmplUnique.Keys {
		delete(params.tmplPropIDMap, key)
	}

	// check if unique keys are the same
	if objCompUnique.compKey == tmplCompKey {
		if !forUI && isNoTmpl {
			// the template attribute in the unique verification refers to the
			// auto-increment ID corresponding to the unique verification of the template
			objCompUnique.unique.TemplateID = tmplUnique.ID
			res.Update = append(res.Update, metadata.CompareOneFieldTmplUniqueRes{
				Index: params.tmplIndexMap[tmplCompKey],
				Data:  &objCompUnique.unique,
			})
			return true, nil
		}
		return false, nil
	}

	// keys need to be processed separately again because it is possible to change from joint
	// unique to single unique for the same unique check. Or change a single unique to a joint unique
	if !forUI {
		objCompUnique.unique.Keys = []metadata.UniqueKey{}
		for id := range tmplUnique.Keys {
			// tmplUnique Here is the propertyID of the template, according to
			// this propertyID, the attribute auto-increment ID of the object is obtained
			objAttrID, ok := params.tmplProToIDMap[tmplUnique.Keys[id]]
			if !ok {
				blog.Errorf("get obj attr id failed, template property id: %v, rid: %s", tmplUnique.Keys[id], kit.Rid)
				return false, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKPropertyIDField)
			}
			objCompUnique.unique.Keys = append(objCompUnique.unique.Keys, metadata.UniqueKey{
				ID:   uint64(objAttrID),
				Kind: metadata.UniqueKeyKindProperty,
			})
		}
	}

	res.Update = append(res.Update, metadata.CompareOneFieldTmplUniqueRes{
		Index: params.tmplIndexMap[tmplCompKey],
		Data:  &objCompUnique.unique,
	})

	return true, nil
}
