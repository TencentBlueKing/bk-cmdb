package settemplate

import (
	"fmt"
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/operation"
)

type BackendWorker struct {
	ClientSet       apimachinery.ClientSetInterface
	ObjectOperation operation.ObjectOperationInterface
	ModuleOperation operation.ModuleOperationInterface
	Engine          *backbone.Engine
}

func (bw BackendWorker) DoModuleSyncTask(header http.Header, set metadata.SetInst, moduleDiff metadata.SetModuleDiff) error {
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)
	user := util.GetUser(header)
	supplierAccount := util.GetOwnerID(header)
	defaultCCError := util.GetDefaultCCError(header)

	kit := &rest.Kit{
		Rid:             rid,
		Header:          header,
		Ctx:             ctx,
		CCError:         defaultCCError,
		User:            user,
		SupplierAccount: supplierAccount,
	}
	moduleObj, err := bw.ObjectOperation.FindSingleObject(kit, common.BKInnerObjIDModule)
	if err != nil {
		blog.Errorf("DoModuleSyncTask failed, FindSingleObject: module failed, err: %s, rid: %s", err.Error(), rid)
		return err
	}

	bizID := set.BizID
	setID := set.SetID
	moduleID := moduleDiff.ModuleID
	switch moduleDiff.DiffType {
	case metadata.ModuleDiffRemove:
		err := bw.ModuleOperation.DeleteModule(kit, set.BizID, []int64{setID}, []int64{moduleID})
		if err != nil {
			blog.ErrorJSON("DoModuleSyncTask failed, DeleteModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
	case metadata.ModuleDiffAdd:
		serviceTemplate, ccErr := bw.ClientSet.CoreService().Process().GetServiceTemplate(ctx, header, moduleDiff.ServiceTemplateID)
		if ccErr != nil {
			blog.Errorf("DoModuleSyncTask failed, DeleteModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, ccErr.Error(), rid)
			return ccErr
		}
		data := map[string]interface{}{
			common.BKModuleNameField:        moduleDiff.ServiceTemplateName,
			common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
			common.BKServiceTemplateIDField: moduleDiff.ServiceTemplateID,
			common.BKParentIDField:          set.SetID,
			common.BKSetTemplateIDField:     set.SetTemplateID,
		}

		_, err := bw.ModuleOperation.CreateModule(kit, moduleObj, bizID, setID, data)
		if err != nil {
			blog.ErrorJSON("DoModuleSyncTask failed, CreateModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
	case metadata.ModuleDiffChanged:
		data := mapstr.MapStr(map[string]interface{}{
			common.BKModuleNameField: moduleDiff.ModuleName,
		})
		err := bw.ModuleOperation.UpdateModule(kit, data, moduleObj, bizID, setID, moduleID)
		if err != nil {
			blog.ErrorJSON("DoModuleSyncTask failed, UpdateModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
	case metadata.ModuleDiffUnchanged:
		return nil
	default:
		blog.ErrorJSON("DoModuleSyncTask failed, UpdateModule failed, set: %s, moduleDiff: %s, rid: %s", set, moduleDiff, rid)
		return fmt.Errorf("unexpected diff type: %s", moduleDiff.DiffType)
	}
	return nil
}
