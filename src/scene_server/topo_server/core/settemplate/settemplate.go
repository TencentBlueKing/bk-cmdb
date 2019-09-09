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

package settemplate

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/mitchellh/mapstructure"
)

type SetTemplate interface {
	DiffSetTplWithInst(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.DiffSetTplWithInstOption) ([]metadata.SetDiff, errors.CCErrorCoder)
	SyncSetTplToInst(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.SyncSetTplToInstOption) errors.CCErrorCoder
}

func NewSetTemplate(client apimachinery.ClientSetInterface) SetTemplate {
	return &setTemplate{
		client: client,
	}
}

type setTemplate struct {
	client apimachinery.ClientSetInterface
}

func (st *setTemplate) DiffSetTplWithInst(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.DiffSetTplWithInstOption) ([]metadata.SetDiff, errors.CCErrorCoder) {
	ccError := util.GetDefaultCCError(header)
	if ccError == nil {
		return nil, errors.GlobalCCErrorNotInitialized
	}

	rid := util.GetHTTPCCRequestID(header)

	serviceTemplates, err := st.client.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(ctx, header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("DiffSetTemplateWithInstances failed, ListSetTplRelatedSvcTpl failed, bizID: %d, setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), rid)
		return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
	}
	serviceTemplateMap := make(map[int64]metadata.ServiceTemplate)
	for _, svcTpl := range serviceTemplates {
		serviceTemplateMap[svcTpl.ID] = svcTpl
	}

	setIDs := util.IntArrayUnique(option.SetIDs)
	setFilter := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKSetIDField: map[string]interface{}{
				common.BKDBIN: setIDs,
			},
		}),
	}
	setInstResult, e := st.client.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDSet, setFilter)
	if e != nil {
		blog.Errorf("DiffSetTemplateWithInstances failed, list sets failed, bizID: %d, setTemplateID: %d, setIDs: %+v, err: %s, rid: %s", bizID, setTemplateID, option.SetIDs, err.Error(), rid)
		return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(setInstResult.Data.Info) != len(setIDs) {
		blog.Errorf("DiffSetTemplateWithInstances failed, some setID invalid, input IDs: %+v, valid IDs: %+v, rid: %s", setIDs, setInstResult.Data.Info, rid)
		return nil, ccError.CCErrorf(common.CCErrCommParamsInvalid, "bk_set_ids")
	}
	setMap := make(map[int64]metadata.SetInst)
	for _, setInstance := range setInstResult.Data.Info {
		set := metadata.SetInst{}
		if err := mapstructure.Decode(setInstance, &set); err != nil {
			blog.Errorf("DiffSetTemplateWithInstances failed, decode set instance failed, set: %+v, err: %s, rid: %s", setInstance, err.Error(), rid)
			return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
		}
		setMap[set.SetID] = set
	}

	moduleFilter := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKParentIDField: map[string]interface{}{
				common.BKDBIN: option.SetIDs,
			},
		}),
	}
	modulesInstResult, e := st.client.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDModule, moduleFilter)
	if e != nil {
		blog.Errorf("DiffSetTemplateWithInstances failed, list modules failed, bizID: %d, setTemplateID: %d, setIDs: %+v, err: %s, rid: %s", bizID, setTemplateID, option.SetIDs, err.Error(), rid)
		return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
	}

	setModules := make(map[int64][]metadata.ModuleInst)
	// init before modules loop so that set with no modules could be initial correctly
	for _, setID := range option.SetIDs {
		setModules[setID] = make([]metadata.ModuleInst, 0)
	}
	for _, moduleInstance := range modulesInstResult.Data.Info {
		module := metadata.ModuleInst{}
		if err := mapstructure.Decode(moduleInstance, &module); err != nil {
			blog.Errorf("DiffSetTemplateWithInstances failed, decode module instance failed, module: %+v, err: %s, rid: %s", moduleInstance, err.Error(), rid)
			return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
		}
		if _, exist := setModules[module.ParentID]; exist == false {
			setModules[module.ParentID] = make([]metadata.ModuleInst, 0)
		}
		setModules[module.ParentID] = append(setModules[module.ParentID], module)
	}

	// diff
	setDiffs := make([]metadata.SetDiff, 0)
	for setID, modules := range setModules {
		moduleDiff := DiffServiceTemplateWithModules(serviceTemplates, modules)
		setDiff := metadata.SetDiff{
			ModuleDiffs: moduleDiff,
			SetID:       setID,
		}
		if set, ok := setMap[setID]; ok == true {
			setDiff.SetDetail = set
		}
		setDiffs = append(setDiffs, setDiff)
	}
	return setDiffs, nil
}

func (st *setTemplate) SyncSetTplToInst(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.SyncSetTplToInstOption) errors.CCErrorCoder {
	rid := util.GetHTTPCCRequestID(header)
	diffOption := metadata.DiffSetTplWithInstOption{
		SetIDs: option.SetIDs,
	}
	setDiffs, err := st.DiffSetTplWithInst(ctx, header, bizID, setTemplateID, diffOption)
	if err != nil {
		return err
	}

	// TODO use queue to dispatch task
	// run sync
	for _, setDiff := range setDiffs {
		blog.V(3).Infof("begin to run sync task on set [%s](%d)", setDiff.SetDetail.SetName, setDiff.SetID)
		// TODO: deal with result
		backendWorker := BackendWorker{
			ClientSet: st.client,
		}
		for _, moduleDiff := range setDiff.ModuleDiffs {
			err := backendWorker.AsyncRunModuleSyncTask(header, setDiff.SetDetail, moduleDiff)
			if err != nil {
				blog.Errorf("AsyncRunSetSyncTask failed, err: %+v, rid: %s", err, rid)
				continue
			}
		}
	}
	return nil
}

// DiffServiceTemplateWithModules diff modules with template in one set
func DiffServiceTemplateWithModules(serviceTemplates []metadata.ServiceTemplate, modules []metadata.ModuleInst) []metadata.SetModuleDiff {
	svcTplMap := make(map[int64]metadata.ServiceTemplate)
	svcTplHitMap := make(map[int64]bool)
	for _, svcTpl := range serviceTemplates {
		svcTplMap[svcTpl.ID] = svcTpl
		svcTplHitMap[svcTpl.ID] = false
	}

	moduleMap := make(map[int64]metadata.ModuleInst)
	for _, module := range modules {
		moduleMap[module.ModuleID] = module
	}

	moduleDiffs := make([]metadata.SetModuleDiff, 0)
	for _, module := range modules {
		template, ok := svcTplMap[module.ServiceTemplateID]
		if ok == false {
			moduleDiffs = append(moduleDiffs, metadata.SetModuleDiff{
				ModuleID:            module.ModuleID,
				ModuleName:          module.ModuleName,
				ServiceTemplateID:   module.ServiceTemplateID,
				ServiceTemplateName: "",
				DiffType:            metadata.ModuleDiffRemove,
			})
			continue
		}

		if module.ModuleName != template.Name {
			moduleDiffs = append(moduleDiffs, metadata.SetModuleDiff{
				ModuleID:            module.ModuleID,
				ModuleName:          module.ModuleName,
				ServiceTemplateID:   module.ServiceTemplateID,
				ServiceTemplateName: template.Name,
				DiffType:            metadata.ModuleDiffChanged,
			})
		}
	}

	for templateID, hit := range svcTplHitMap {
		if hit == true {
			continue
		}
		template := svcTplMap[templateID]
		moduleDiffs = append(moduleDiffs, metadata.SetModuleDiff{
			ModuleID:            0,
			ModuleName:          "",
			ServiceTemplateID:   templateID,
			ServiceTemplateName: template.Name,
			DiffType:            metadata.ModuleDiffAdd,
		})
	}
	return moduleDiffs
}
