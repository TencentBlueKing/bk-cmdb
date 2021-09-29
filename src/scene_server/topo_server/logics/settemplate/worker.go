package settemplate

import (
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/logics/inst"
	"configcenter/src/scene_server/topo_server/logics/model"
)

type BackendWorker struct {
	ClientSet       apimachinery.ClientSetInterface
	ObjectOperation model.ObjectOperationInterface
	ModuleOperation inst.ModuleOperationInterface
	Engine          *backbone.Engine
}

// DoModuleSyncTask do module sync task
func (bw BackendWorker) DoModuleSyncTask(header http.Header, set metadata.SetInst,
	moduleDiff metadata.SetModuleDiff) error {

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

	bizID := set.BizID
	setID := set.SetID
	moduleID := moduleDiff.ModuleID
	switch moduleDiff.DiffType {
	case metadata.ModuleDiffRemove:
		err := bw.ModuleOperation.DeleteModule(kit, set.BizID, []int64{setID}, []int64{moduleID})
		if err != nil {
			blog.ErrorJSON("delete module failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err, rid)
			return err
		}
	case metadata.ModuleDiffAdd:
		serviceTemplate, ccErr := bw.ClientSet.CoreService().Process().GetServiceTemplate(ctx, header,
			moduleDiff.ServiceTemplateID)
		if ccErr != nil {
			blog.Errorf("delete module failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, ccErr, rid)
			return ccErr
		}
		data := map[string]interface{}{
			common.BKModuleNameField:        moduleDiff.ServiceTemplateName,
			common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
			common.BKServiceTemplateIDField: moduleDiff.ServiceTemplateID,
			common.BKParentIDField:          set.SetID,
			common.BKSetTemplateIDField:     set.SetTemplateID,
		}

		_, err := bw.ModuleOperation.CreateModule(kit, bizID, setID, data)
		if err != nil {
			blog.ErrorJSON("create module failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err, rid)
			return err
		}
	case metadata.ModuleDiffChanged:
		data := mapstr.MapStr(map[string]interface{}{
			common.BKModuleNameField: moduleDiff.ModuleName,
		})
		err := bw.ModuleOperation.UpdateModule(kit, data, bizID, setID, moduleID)
		if err != nil {
			blog.ErrorJSON("update module failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err, rid)
			return err
		}
	// unchanged module should not be in asynchronous tasks
	default:
		blog.Errorf("do module sync task but diff type is invalid/unchanged, moduleDiff: %#v, rid: %s", moduleDiff, rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "diff_type")
	}
	return nil
}
