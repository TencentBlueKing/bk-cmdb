package settemplate

import (
	"net/http"
	"strconv"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type BackendWorker struct {
	ClientSet apimachinery.ClientSetInterface
}

func (bw BackendWorker) AsyncRunModuleSyncTask(header http.Header, set metadata.SetInst, moduleDiff metadata.SetModuleDiff) error {
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)

	bizIDStr := strconv.FormatInt(set.BizID, 10)
	setIDStr := strconv.FormatInt(set.SetID, 10)
	moduleIDStr := strconv.FormatInt(moduleDiff.ModuleID, 10)
	switch moduleDiff.DiffType {
	case metadata.ModuleDiffRemove:
		resp, err := bw.ClientSet.TopoServer().Instance().DeleteModule(ctx, bizIDStr, setIDStr, moduleIDStr, header)
		if err != nil {
			blog.Errorf("AsyncRunSetSyncTask failed, DeleteModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
		if resp.Result == false || resp.Code != 0 {
			blog.Errorf("AsyncRunSetSyncTask failed, DeleteModule failed, set: %s, moduleDiff: %s, rid: %s", set, moduleDiff, rid)
			return errors.New(resp.Code, resp.ErrMsg)
		}
	case metadata.ModuleDiffAdd:
		serviceTemplate, ccErr := bw.ClientSet.CoreService().Process().GetServiceTemplate(ctx, header, moduleDiff.ServiceTemplateID)
		if ccErr != nil {
			blog.Errorf("AsyncRunSetSyncTask failed, DeleteModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, ccErr.Error(), rid)
			return ccErr
		}
		data := map[string]interface{}{
			common.BKModuleNameField:        moduleDiff.ModuleName,
			common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
			common.BKServiceTemplateIDField: moduleDiff.ServiceTemplateID,
			common.BKParentIDField:          set.SetID,
			common.BKSetTemplateIDField:     set.SetTemplateID,
		}

		resp, err := bw.ClientSet.TopoServer().Instance().CreateModule(ctx, bizIDStr, setIDStr, header, data)
		if err != nil {
			blog.Errorf("AsyncRunSetSyncTask failed, CreateModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
		if resp.Result == false || resp.Code != 0 {
			blog.Errorf("AsyncRunSetSyncTask failed, CreateModule failed, set: %s, moduleDiff: %s, rid: %s", set, moduleDiff, rid)
			return errors.New(resp.Code, resp.ErrMsg)
		}
	case metadata.ModuleDiffChanged:
		data := mapstr.MapStr(map[string]interface{}{
			common.BKModuleNameField: moduleDiff.ModuleName,
		})
		resp, err := bw.ClientSet.TopoServer().Instance().UpdateModule(ctx, bizIDStr, setIDStr, moduleIDStr, header, data)
		if err != nil {
			blog.Errorf("AsyncRunSetSyncTask failed, UpdateModule failed, set: %s, moduleDiff: %s, err: %s, rid: %s", set, moduleDiff, err.Error(), rid)
			return err
		}
		if resp.Result == false || resp.Code != 0 {
			blog.Errorf("AsyncRunSetSyncTask failed, UpdateModule failed, set: %s, moduleDiff: %s, rid: %s", set, moduleDiff, rid)
			return errors.New(resp.Code, resp.ErrMsg)
		}
	}
	return nil
}
