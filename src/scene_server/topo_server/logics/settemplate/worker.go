package settemplate

import (
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/logics/inst"
)

// BackendWorker worker to do backend syncing set template task
type BackendWorker struct {
	ClientSet       apimachinery.ClientSetInterface
	ModuleOperation inst.ModuleOperationInterface
	// InstOperation instance operations, used when updating module name
	InstOperation inst.InstOperationInterface
}

// DoModuleSyncTask do syncing module under set template by service template task
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
			blog.Errorf("delete module %d failed, err: %v, biz: %d, set: %d, rid: %s", moduleID, err, bizID, setID, rid)
			return err
		}
	case metadata.ModuleDiffAdd:
		serviceTemplate, ccErr := bw.ClientSet.CoreService().Process().GetServiceTemplate(ctx, header,
			moduleDiff.ServiceTemplateID)
		if ccErr != nil {
			blog.Errorf("get service temp failed, err: %v, id: %d, rid: %s", ccErr, moduleDiff.ServiceTemplateID, rid)
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
			blog.Errorf("create module(%#v) failed, err: %v, biz: %d, set: %d, rid: %s", data, err, bizID, setID, rid)
			return err
		}
	case metadata.ModuleDiffChanged:
		cond := mapstr.MapStr{
			common.BKAppIDField:    bizID,
			common.BKSetIDField:    setID,
			common.BKModuleIDField: moduleID,
		}
		data := mapstr.MapStr(map[string]interface{}{
			common.BKModuleNameField: moduleDiff.ServiceTemplateName,
		})

		err := bw.InstOperation.UpdateInst(kit, cond, data, common.BKInnerObjIDModule)
		if err != nil {
			blog.Errorf("update module failed, cond: %#v, data: %#v, err: %v, rid: %s", cond, data, err, rid)
			return err
		}
	case metadata.ModuleDiffUnchanged:
		return nil
	default:
		blog.ErrorJSON("module sync task diff type(%s) is invalid, rid: %s", moduleDiff.DiffType, rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "diff_type")
	}
	return nil
}
