package settemplate

import (
	"net/http"
	"reflect"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
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

	// 进行集群属性的同步
	if err := bw.syncSetAttributes(kit, bizID, set.SetTemplateID, setID); err != nil {
		blog.Errorf("set attribute sync task failed, set template id: %d, setID: %d, err: %v, rid: %s",
			set.SetTemplateID, setID, err, rid)
		return err
	}

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
		blog.Errorf("module sync task diff type(%s) is invalid, rid: %s", moduleDiff.DiffType, rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "diff_type")
	}
	return nil
}

// getSetTemplateAttrIdAndPropertyValue 获取集群模板的属性id以及对应的属性值
func (bw BackendWorker) getSetTemplateAttrIdAndPropertyValue(kit *rest.Kit, bizID, setTemplateID int64) ([]int64,
	map[int64]interface{}, errors.CCErrorCoder) {

	option := &metadata.ListSetTempAttrOption{
		BizID:  bizID,
		ID:     setTemplateID,
		Fields: []string{common.BKAttributeIDField, common.BKPropertyValueField},
	}

	data, err := bw.ClientSet.CoreService().SetTemplate().ListSetTemplateAttribute(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("list set template attributes failed, bizID: %d, set template id: %d, err: %v, rid: %s", bizID,
			setTemplateID, err, kit.Rid)
		return nil, nil, err
	}

	attrIDs := make([]int64, 0)
	setTemplateAttrValueMap := make(map[int64]interface{})

	for _, attr := range data.Attributes {
		attrIDs = append(attrIDs, attr.AttributeID)
		setTemplateAttrValueMap[attr.AttributeID] = attr.PropertyValue
	}

	return attrIDs, setTemplateAttrValueMap, nil
}

func (bw *BackendWorker) getSetMapStr(kit *rest.Kit, bizID, setTemplateId int64, setID int64,
	fields []string) ([]mapstr.MapStr, errors.CCErrorCoder) {

	option := &metadata.QueryCondition{
		Fields: fields,
		Condition: map[string]interface{}{
			common.BKSetIDField:         setID,
			common.BKSetTemplateIDField: setTemplateId,
			common.BKAppIDField:         bizID,
		},
		DisableCounter: true,
	}

	set, err := bw.ClientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, option)
	if err != nil {
		blog.Errorf("get set failed, option: %+v, err: %v, rid: %s", *option, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed, err.Error())
	}
	if len(set.Info) == 0 {
		blog.Errorf("no set founded, option: %+v, err: %v, rid: %s", *option, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound)
	}
	return set.Info, nil
}

// getSetAttrIDAndPropertyID 根据集群属性ID获取对应的propertyID列表以及属性ID与propertyID的对应关系
func (bw BackendWorker) getSetAttrIDAndPropertyID(kit *rest.Kit, attrIDs []int64) ([]string, map[int64]string,
	errors.CCErrorCoder) {

	attrIdPropertyMap := make(map[int64]string)
	if len(attrIDs) == 0 {
		return []string{}, attrIdPropertyMap, nil
	}

	option := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID, common.BKPropertyIDField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: attrIDs,
			},
		},
		DisableCounter: true,
	}

	res, err := bw.ClientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDSet, option)
	if err != nil {
		blog.Errorf("read set attribute failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeSelectFailed)
	}

	propertyIDs := make([]string, 0)
	for _, attrs := range res.Info {
		propertyIDs = append(propertyIDs, attrs.PropertyID)
		attrIdPropertyMap[attrs.ID] = attrs.PropertyID
	}

	return propertyIDs, attrIdPropertyMap, nil
}

// updateSetAttributesWithSetTemplate 这里更新采用的是先对比不同再更新，因为采用模板方式进行集群属性的更新不会是一个高频操作，大概率此处
// 是不需要更新的，采用先查找属性值不同的然后针对不同值进行更新
func (bw BackendWorker) updateSetAttributesWithSetTemplate(kit *rest.Kit, setID int64, setMap mapstr.MapStr,
	setTemplateAttrValueMap map[int64]interface{}, attrIdPropertyMap map[int64]string) errors.CCErrorCoder {

	if len(setTemplateAttrValueMap) == 0 || len(attrIdPropertyMap) == 0 {
		return nil
	}

	data := make(map[string]interface{})

	// 这里模板可能没有配置属性
	for setTemplateAttrID, value := range setTemplateAttrValueMap {
		if !reflect.DeepEqual(value, setMap[attrIdPropertyMap[setTemplateAttrID]]) {
			data[attrIdPropertyMap[setTemplateAttrID]] = value
		}
	}

	// 如果没有属性可以更新直接返回
	if len(data) == 0 {
		return nil
	}

	option := &metadata.UpdateOption{
		Data: data,
		Condition: map[string]interface{}{
			common.BKSetIDField: setID,
		},
	}

	_, err := bw.ClientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, option)
	if err != nil {
		blog.Errorf("update set failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrUpdateModuleAttributesFail)
	}
	return nil
}

// syncSetAttributes TODO
func (bw BackendWorker) syncSetAttributes(kit *rest.Kit, bizID, setTemplateID, setID int64) error {
	// 1、获取集群模板中的属性id以及对应的属性值 property_value
	attrIDs, setTemplateAttrValueMap, cErr := bw.getSetTemplateAttrIdAndPropertyValue(kit, bizID, setTemplateID)
	if cErr != nil {
		return cErr
	}
	// 2、从cc_ObjAttDes 中通过上面的属性id获取对应的 bk_property_id
	propertyIDs, attrIdPropertyMap, cErr := bw.getSetAttrIDAndPropertyID(kit, attrIDs)
	if cErr != nil {
		return cErr
	}

	setMap, cErr := bw.getSetMapStr(kit, bizID, setTemplateID, setID, propertyIDs)
	if cErr != nil {
		return cErr
	}
	// 3、直接更新set对应的这些 bk_property_id 的值为 property_value
	if err := bw.updateSetAttributesWithSetTemplate(kit, setID, setMap[0], setTemplateAttrValueMap,
		attrIdPropertyMap); err != nil {
		return err
	}
	return nil
}
