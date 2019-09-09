package worker

import (
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core"
	"configcenter/src/scene_server/topo_server/types"
)

type BackendWorker struct {
	ClientSet   apimachinery.ClientSetInterface
	Core        core.Core
	Engine      *backbone.Engine
	AuthManager *extensions.AuthManager
	Error       errors.CCErrorIf
	Language    language.CCLanguageIf
}

func (bw BackendWorker) AsyncRunModuleSyncTask(header http.Header, set metadata.SetInst, moduleDiff metadata.SetModuleDiff) error {
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)

	language := util.GetLanguage(header)
	defLang := bw.Language.CreateDefaultCCLanguageIf(language)
	errors.SetGlobalCCError(bw.Error)
	defErr := bw.Error.CreateDefaultCCErrorIf(language)
	params := types.ContextParams{
		Context:         ctx,
		Engin:           bw.Engine,
		Header:          header,
		MaxTopoLevel:    0,
		SupplierAccount: util.GetOwnerID(header),
		User:            util.GetUser(header),
		Err:             defErr,
		Lang:            defLang,
		MetaData:        nil,
		ReqID:           rid,
	}
	obj, err := bw.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("create module failed, failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
		return err
	}

	switch moduleDiff.DiffType {
	case metadata.ModuleDiffRemove:
		err := bw.Core.ModuleOperation().DeleteModule(params, obj, set.BizID, []int64{set.SetID}, []int64{moduleDiff.ModuleID})
		if err != nil {
			blog.Errorf("AsyncRunSetSyncTask failed, DeleteModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
	case metadata.ModuleDiffAdd:
		serviceTemplate, err := bw.ClientSet.CoreService().Process().GetServiceTemplate(params.Context, header, moduleDiff.ServiceTemplateID)
		if err != nil {
			blog.Errorf("AsyncRunSetSyncTask failed, DeleteModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
		data := mapstr.MapStr(map[string]interface{}{
			common.BKModuleNameField:        moduleDiff.ModuleName,
			common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
			common.BKServiceTemplateIDField: moduleDiff.ServiceTemplateID,
			common.BKParentIDField:          set.SetID,
			common.BKSetTemplateIDField:     set.SetTemplateID,
		})
		if _, err := bw.Core.ModuleOperation().CreateModule(params, obj, set.BizID, set.SetID, data); err != nil {
			blog.Errorf("AsyncRunSetSyncTask failed, DeleteModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
	case metadata.ModuleDiffChanged:
		data := mapstr.MapStr(map[string]interface{}{
			common.BKModuleNameField: moduleDiff.ModuleName,
		})
		err := bw.Core.ModuleOperation().UpdateModule(params, data, obj, set.BizID, set.SetID, moduleDiff.ModuleID)
		if err != nil {
			blog.ErrorJSON("AsyncRunSetSyncTask failed, UpdateModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
	}
	return nil
}
