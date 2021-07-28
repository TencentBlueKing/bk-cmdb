package model

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	// CreateMainlineAssociation create mainline object association
	CreateMainlineAssociation(kit *rest.Kit, data *metadata.Association, maxTopoLevel int) (*metadata.Object, error)
	// DeleteMainlineAssociation delete mainline association by objID
	DeleteMainlineAssociation(kit *rest.Kit, objID string) error
	// SearchMainlineAssociationTopo get mainline topo of special model
	SearchMainlineAssociationTopo(kit *rest.Kit, targetObj *metadata.Object) ([]*metadata.MainlineObjectTopo, error)
	// IsMainlineObject check whether objID is mainline object or not
	IsMainlineObject(kit *rest.Kit, objID string) (bool, error)
}

// NewAssociationOperation create a new association operation instance
func NewAssociationOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) AssociationOperationInterface {
	return &association{
		clientSet:   client,
		authManager: authManager,
	}
}

type association struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// CreateMainlineAssociation create mainline object association
func (assoc *association) CreateMainlineAssociation(kit *rest.Kit, data *metadata.Association,
	maxTopoLevel int) (*metadata.Object, error) {

	if data.AsstObjID == "" {
		blog.ErrorJSON("bk_asst_obj_id empty, input: %s, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAsstObjIDField)
	}

	if data.ClassificationID == "" {
		blog.ErrorJSON("bk_classification_id empty, input: %s, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKClassificationIDField)
	}

	// find the mainline module's head, which is biz.
	// TODO after merge change this logic to FindSingleObject
	rsp, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{
			Condition: mapstr.MapStr{
				common.BKObjIDField: mapstr.MapStr{
					common.BKDBIN: []string{common.BKInnerObjIDApp, data.AsstObjID},
				},
			},
		})
	if err != nil {
		blog.Errorf("find object failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search the objects , err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(rsp.Data.Info) == 0 {
		blog.Errorf("search object failed, return is empty, rid: %s", kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
			kit.CCError.CCError(common.CCErrCommNotFound).Error())
	}

	objIDMap := make(map[string]*metadata.Object, 0)
	for _, obj := range rsp.Data.Info {
		objIDMap[obj.ObjectID] = &obj
	}

	items, err := assoc.SearchMainlineAssociationTopo(kit, objIDMap[common.BKInnerObjIDApp])
	if err != nil {
		blog.Errorf("failed to check the mainline topo level, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(items) >= maxTopoLevel {
		blog.Errorf("the mainline topo level is %d, the max limit is %d, rid: %s", len(items), maxTopoLevel, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoBizTopoLevelOverLimit)
	}

	// find the mainline parent object
	parentObj := objIDMap[data.AsstObjID]
	if parentObj == nil {
		blog.Errorf("object(%s) is non-exist, rid: %s", data.AsstObjID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoObjectSelectFailed)
	}

	// find the mainline child object for the parent
	childObj, err := assoc.getMainlineObject(kit, parentObj.ObjectID, true)
	if err != nil {
		blog.Errorf("failed to find the child object of the object(%s), err: %v, rid: %s",
			parentObj.ObjectID, err, kit.Rid)
		return nil, err
	}

	// check and create the association mainline object
	if err = assoc.IsValidObject(kit, data.ObjectID); err == nil {
		blog.Errorf("the object(%s) is duplicate, rid: %s", data.ObjectID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, data.ObjectID)
	}

	objData := mapstr.MapStr{
		common.BKObjIDField:            data.ObjectID,
		common.BKObjNameField:          data.ObjectName,
		common.BKObjIconField:          data.ObjectIcon,
		common.BKClassificationIDField: data.ClassificationID,
	}
	currentObj, err := assoc.CreateObject(kit, true, objData)
	if err != nil {
		return nil, err
	}

	// update the mainline topo inst association
	createdInstIDs, err := assoc.SetMainlineInstAssociation(kit, parentObj, currentObj, childObj)
	if err != nil {
		blog.Errorf("failed set the mainline inst association, err: %s, rid: %s", err, kit.Rid)
		return nil, err
	}

	if err = assoc.createMainlineObjectAssociation(kit, currentObj, parentObj.ObjectID); err != nil {
		blog.Errorf("create mainline object[%s] association related to object[%s] failed, err: %v, rid: %s",
			currentObj.ObjectID, parentObj.ObjectID, err, kit.Rid)
		return nil, err
	}

	if err = assoc.setMainlineParentObject(kit, childObj, currentObj.ObjectID); err != nil {
		blog.Errorf("update mainline current object's[%s] child object[%s] association to current failed, "+
			"err: %v, rid: %s", currentObj.ObjectID, childObj.ObjectID, err, kit.Rid)
		return nil, err
	}

	// create audit log for the created instances.
	audit := auditlog.NewInstanceAudit(assoc.clientSet.CoreService())

	cond := map[string]interface{}{
		currentObj.GetInstIDFieldName(): map[string]interface{}{
			common.BKDBIN: createdInstIDs,
		},
	}

	// generate audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, currentObj.GetObjectID(), cond)
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("creat inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrAuditSaveLogFailed)
	}

	return currentObj, nil
}

// DeleteMainlineAssociation delete mainline association by objID
func (assoc *association) DeleteMainlineAssociation(kit *rest.Kit, objID string) error {

	rsp, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKObjIDField: objID}})
	if err != nil {
		blog.Errorf("find object failed, objID: %s, err: %v, rid: %s", objID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search the object , err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if len(rsp.Data.Info) != 1 {
		blog.Errorf("failed to search the single object(%s) , err: %v, rid: %s", objID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoObjectSelectFailed)
	}

	targetObj := rsp.Data.Info[0]

	parentObj, err := assoc.getMainlineObject(kit, targetObj.ObjectID, false)
	if err != nil {
		blog.Errorf("failed to find the object(%s)'s parent, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	// update associations
	childObj, err := assoc.getMainlineObject(kit, targetObj.ObjectID, true)
	if err != nil {
		blog.Errorf("failed to find the object(%s)'s child, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	if err = assoc.ResetMainlineInstAssociation(kit, &targetObj); err != nil {
		blog.Errorf("failed to delete the object(%s)'s instance, error info %s, rid: %s", objID, err, kit.Rid)
		return err
	}

	if nil != childObj {
		// FIX: 正常情况下 childObj 不可以能为 nil，只有在拓扑异常的时候才会出现
		if err = assoc.setMainlineParentObject(kit, childObj, parentObj.ObjectID); err != nil {
			blog.Errorf("failed to update the association, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}

	// delete this object related association.
	cond := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{metadata.AssociationFieldObjectID: objID},
			{metadata.AssociationFieldAssociationObjectID: objID},
		},
	}

	if err = assoc.DeleteAssociation(kit, cond); err != nil {
		return err
	}

	// delete objects
	if err = assoc.DeleteObject(kit, targetObj.ID, false); err != nil {
		blog.Errorf("failed to delete the object(%s), err: %v, rid: %s", targetObj.ID, err, kit.Rid)
		return err
	}

	return nil
}

// SearchMainlineAssociationTopo get mainline topo of special model
// result is a list with targetObj as head, so if you want a full topo, target must be biz model.
func (assoc *association) SearchMainlineAssociationTopo(kit *rest.Kit,
	targetObj *metadata.Object) ([]*metadata.MainlineObjectTopo, error) {

	// foundObjIDMap used as a map to detect whether found model is already in,
	// so that we can detect infinite loop.
	foundObjIDMap := make(map[string]bool)
	results := make([]*metadata.MainlineObjectTopo, 0)
	for {
		resultsLen := len(results)
		tmpRst := &metadata.MainlineObjectTopo{}
		tmpRst.ObjID = targetObj.ObjectID
		tmpRst.ObjName = targetObj.ObjectName
		tmpRst.OwnerID = kit.SupplierAccount

		parentObj, err := assoc.getMainlineObject(kit, targetObj.ObjectID, false)
		if err != nil {
			blog.Errorf("get object(%s)'s parent object failed, err: %v, rid: %s", targetObj.ObjectID, err, kit.Rid)
			return nil, err
		}
		tmpRst.PreObjID = parentObj.ObjectID
		tmpRst.PreObjName = parentObj.ObjectName

		childObj, err := assoc.getMainlineObject(kit, targetObj.ObjectID, true)
		if nil == err {
			tmpRst.NextObj = childObj.ObjectID
			tmpRst.NextName = childObj.ObjectName
		} else if err != nil {
			if io.EOF != err {
				return nil, err
			}
			if _, ok := foundObjIDMap[tmpRst.ObjID]; !ok {
				results = append(results, tmpRst)
				foundObjIDMap[tmpRst.ObjID] = true
			}
			return results, nil
		}

		if _, ok := foundObjIDMap[tmpRst.ObjID]; !ok {
			results = append(results, tmpRst)
			foundObjIDMap[tmpRst.ObjID] = true
		}

		targetObj = childObj

		// detect infinite loop by checking whether there are new added objects in current loop.
		if resultsLen == len(results) {
			// merely return found objects here to avoid infinite loop.
			// returned results here maybe parts of all mainline objects.
			// better to prevent loop from taking shape seriously, at adding or editing association.
			return results, nil
		}
	}

}

// IsMainlineObject check whether objID is mainline object or not
func (assoc *association) IsMainlineObject(kit *rest.Kit, objID string) (bool, error) {
	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	asst, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return false, err
	}

	if !asst.Result {
		return false, errors.New(asst.Code, asst.ErrMsg)
	}

	if len(asst.Data.Info) <= 0 {
		return false, fmt.Errorf("model association [%+v] not found", cond)
	}

	for _, mainline := range asst.Data.Info {
		if mainline.ObjectID == objID || mainline.AsstObjID == objID {
			return true, nil
		}
	}

	return false, nil
}

// TODO after merge , delete this func and use SetMainlineInstAssociation in inst/mainline_association
func (assoc *association) SetMainlineInstAssociation(kit *rest.Kit, parent, current,
	child *metadata.Object) ([]int64, error) {

	defaultCond := &metadata.QueryInput{}
	cond := mapstr.New()
	if parent.IsCommon() {
		cond.Set(common.BKObjIDField, parent.ObjectID)
	}
	defaultCond.Condition = cond
	// fetch all parent instances.
	// TODO replace to FindInst in inst/inst.go after merge
	parentInsts, err := assoc.FindInst(kit, parent.ObjectID, defaultCond)
	if err != nil {
		blog.Errorf("failed to find parent object(%s) inst, err: %v, rid: %s", parent.ObjectID, err, kit.Rid)
		return nil, err
	}

	createdInstIDs := make([]int64, len(parentInsts.Info))
	expectParent2Children := map[int64][]mapstr.MapStr{}
	// filters out special character for mainline instances
	re, _ := regexp.Compile(`[#/,><|]`)
	instanceName := re.ReplaceAllString(current.ObjectName, "")
	// create current object instance for each parent instance and insert the current instance to
	for _, parentInst := range parentInsts.Info {
		id, err := parentInst.Int64(metadata.GetInstIDFieldByObjID(parent.ObjectID))
		if err != nil {
			blog.Errorf("failed to find the inst id, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		// we create the current object's instance for each parent instance belongs to the parent object.
		currentInst := mapstr.MapStr{common.BKObjIDField: current.ObjectID}
		currentInst.Set(current.GetInstNameFieldName(), instanceName)
		// set current instance's parent id to parent instance's id, so that they can be chained.
		currentInst.Set(common.BKInstParentStr, id)

		if parent.GetObjectID() == common.BKInnerObjIDApp {
			currentInst.Set(common.BKAppIDField, id)
		} else {
			if bizID, ok := parentInst.Get(common.BKAppIDField); ok {
				currentInst.Set(common.BKAppIDField, bizID)
			}
		}

		// create the instance now.
		instID, err := assoc.createInst(kit, current.ObjectID, currentInst)
		if err != nil {
			blog.Errorf("failed to create object(%s) default inst, err: %v, rid: %s", current.ObjectID, err, kit.Rid)
			return nil, err
		}

		createdInstIDs = append(createdInstIDs, int64(instID))

		// reset the child's parent instance's parent id to current instance's id.
		children, err := assoc.getMainlineInst(kit, parent, parentInst, true)
		if err != nil {
			blog.Errorf("failed to get the object(%s) mainline child inst, err: %v, rid: %s",
				parent.ObjectID, err, kit.Rid)
			return nil, err
		}

		expectParent2Children[int64(instID)] = children
	}

	for parentID, children := range expectParent2Children {
		for _, child := range children {
			// set the child's parent
			if err = assoc.SetMainlineParentInst(kit, child, parentID); err != nil {
				blog.Errorf("failed to set the object mainline child inst, err: %v, rid: %s", err, kit.Rid)
				return nil, err
			}
		}
	}
	return createdInstIDs, nil
}

// getMainlineObject get mainline relationship model
// the parent not exactly mean parent in a tree case
// TODO after merge this function should be moved to module/object
func (assoc *association) getMainlineObject(kit *rest.Kit, objID string, isChild bool) (*metadata.Object, error) {

	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}

	if isChild {
		cond.Set(common.BKAsstObjIDField, objID)
	} else {
		cond.Set(common.BKObjIDField, objID)
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	for _, asst := range rsp.Data.Info {
		if isChild {
			cond = mapstr.MapStr{common.BKObjIDField: asst.ObjectID}
		} else {
			cond = mapstr.MapStr{common.BKObjIDField: asst.AsstObjID}
		}

		// TODO after merge this logic can be replace by FindObject
		rspRst, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
			&metadata.QueryCondition{Condition: cond})
		if err != nil {
			blog.Errorf("request to search object failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rspRst.CCError(); err != nil {
			blog.Errorf("request to search object failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		if len(rspRst.Data.Info) > 1 {
			blog.Errorf("get multiple(%d) children/parent for object(%s), rid: %s",
				len(rspRst.Data.Info), asst.ObjectID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, asst.AsstObjID)
		}

		for _, item := range rspRst.Data.Info {
			// only one parent in the main-line
			return &item, nil
		}

	}

	return &metadata.Object{}, nil
}

func (assoc *association) getMainlineInst(kit *rest.Kit, obj *metadata.Object, inst mapstr.MapStr,
	needChild bool) ([]mapstr.MapStr, error) {

	mainlineObj, err := assoc.getMainlineObject(kit, obj.ObjectID, needChild)
	if err != nil {
		blog.Errorf("failed to get object, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
		return nil, err
	}

	instID, err := inst.Int64(obj.GetInstIDFieldName())
	parentID, err := inst.Int64(common.BKInstParentStr)

	cond := mapstr.New()
	if mainlineObj.IsCommon() {
		cond.Set(metadata.ModelFieldObjectID, mainlineObj.ObjectID)
	} else if mainlineObj.ObjectID == common.BKInnerObjIDSet {
		cond.Set(common.BKDefaultField, mapstr.MapStr{common.BKDBNE: common.DefaultResSetFlag})
	}

	if needChild {
		cond.Set(common.BKInstParentStr, instID)
	} else {
		cond.Set(mainlineObj.GetInstIDFieldName(), parentID)
	}

	instRsp, err := assoc.FindInst(kit, mainlineObj.ObjectID, &metadata.QueryInput{Condition: cond})
	if err != nil {
		blog.Errorf("search inst by object(%s) failed, err: %v, rid: %s", mainlineObj.ObjectID, err, kit.Rid)
		return nil, err
	}

	return instRsp.Info, nil
}

func (assoc *association) SetMainlineParentInst(kit *rest.Kit, childInst mapstr.MapStr, instID int64) error {
	if err := assoc.updateMainlineAssociation(kit, childInst, instID); err != nil {
		blog.Errorf("failed to update the mainline association, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// IsValidObject check whether objID is a real model's bk_obj_id field in backend
// TODO this function should be delete.
// TODO every function which use this logic need to replace to IsValidObject in model/object.go.
func (assoc *association) IsValidObject(kit *rest.Kit, objID string) error {

	checkObjCond := mapstr.MapStr{
		common.BKObjIDField: objID,
	}

	objItems, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: checkObjCond},
	)
	if err != nil {
		blog.Errorf("find object failed, cond: %+v, err: %v, rid: %s", checkObjCond, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = objItems.CCError(); err != nil {
		blog.Errorf("failed to search the objects by the condition(%#v) , err: %v, rid: %s", checkObjCond, err, kit.Rid)
		return err
	}

	if len(objItems.Data.Info) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	return nil
}

// TODO this function should be delete.
// TODO every function which use this logic need to replace to IsValidObject in model/object.go.
func (assoc *association) CreateObject(kit *rest.Kit, isMainline bool, data mapstr.MapStr) (*metadata.Object, error) {

	obj, err := IsValid(kit, false, data)
	if err != nil {
		blog.Errorf("valid data(%#v) failed, err: %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	objCls, err := assoc.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKClassificationIDField: obj.ObjCls}})
	if err != nil {
		blog.Errorf("get object classification failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if err = objCls.CCError(); err != nil {
		blog.Errorf("get object classification failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(objCls.Data.Info) == 0 {
		blog.Errorf("can't find classification by params, classification: %s is not exist, rid: %s",
			obj.ObjCls, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKClassificationIDField)
	}

	filter := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{common.BKObjIDField: obj.ObjectID},
			{common.BKObjNameField: obj.ObjectName},
		}}
	cnt, err := assoc.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameObjDes,
		[]map[string]interface{}{filter})
	if err != nil {
		blog.Errorf("get object number by filter failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	for index := range cnt {
		if cnt[index] != 0 {
			return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "bk_obj_id/bk_obj_name")
		}
	}

	if len(obj.ObjIcon) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIconField)
	}

	obj.OwnerID = kit.SupplierAccount

	objRsp, err := assoc.clientSet.CoreService().Model().CreateModel(kit.Ctx, kit.Header, &metadata.CreateModel{Spec: *obj})
	if err != nil {
		blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = objRsp.CCError(); err != nil {
		blog.Errorf("failed to search the object(%s), err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
		return nil, err
	}

	obj.ID = int64(objRsp.Data.Created.ID)

	// create the default group
	groupData := metadata.Group{
		IsDefault:  true,
		GroupIndex: -1,
		GroupName:  "Default",
		GroupID:    "default",
		ObjectID:   obj.ObjectID,
		OwnerID:    obj.OwnerID,
	}

	rspGrp, err := assoc.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header,
		obj.ObjectID, metadata.CreateModelAttributeGroup{Data: groupData})
	if err != nil {
		blog.Errorf("create attribute group failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rspGrp.CCError(); err != nil {
		blog.Errorf("create attribute group[%s] failed, err: %v, rid: %s", groupData.GroupID, err, kit.Rid)
		return nil, err
	}

	keys := make([]metadata.UniqueKey, 0)
	// create the default inst attribute
	attr := metadata.Attribute{
		ObjectID:          obj.ObjectID,
		IsOnly:            true,
		IsPre:             true,
		Creator:           "user",
		IsEditable:        true,
		PropertyIndex:     -1,
		PropertyGroup:     groupData.GroupID,
		PropertyGroupName: groupData.GroupName,
		IsRequired:        true,
		PropertyType:      common.FieldTypeSingleChar,
		PropertyID:        obj.GetInstNameFieldName(),
		PropertyName:      obj.GetDefaultInstPropertyName(),
		OwnerID:           kit.SupplierAccount,
	}

	// create a new record
	rspAttr, err := assoc.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header,
		attr.ObjectID, &metadata.CreateModelAttributes{Attributes: []metadata.Attribute{attr}})
	if err != nil {
		blog.Errorf("failed to request coreService to create model attrs, err: %v, ObjectID: %s, input: %#v, rid: %s",
			err, attr.ObjectID, attr, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rspAttr.CCError(); err != nil {
		blog.Errorf("create model attrs failed, ObjectID: %s, input: %#v, err: %v, rid: %s",
			attr.ObjectID, attr, err, kit.Rid)
		return nil, rspAttr.CCError()
	}

	for _, exception := range rspAttr.Data.Exceptions {
		return nil, kit.CCError.New(int(exception.Code), exception.Message)
	}

	if len(rspAttr.Data.Repeated) > 0 {
		blog.Errorf("create model attrs failed, the attr is duplicated, ObjectID: %s, input: %#v, rid: %s",
			attr.ObjectID, attr, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	if len(rspAttr.Data.Created) != 1 {
		blog.Errorf("create model attrs created amount error, ObjectID: %s, input: %#v, rid: %s",
			attr.ObjectID, attr, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
	}

	attr.ID = int64(rspAttr.Data.Created[0].ID)

	keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(attr.ID)})

	if isMainline {
		pAttr := metadata.Attribute{
			ObjectID:          obj.ObjectID,
			IsOnly:            true,
			IsPre:             true,
			Creator:           "system",
			IsEditable:        true,
			PropertyIndex:     -1,
			PropertyGroup:     groupData.GroupID,
			PropertyGroupName: groupData.GroupName,
			IsRequired:        true,
			PropertyType:      common.FieldTypeInt,
			PropertyID:        common.BKInstParentStr,
			PropertyName:      common.BKInstParentStr,
			IsSystem:          true,
			OwnerID:           kit.SupplierAccount,
		}

		rsppAttr, err := assoc.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header,
			pAttr.ObjectID, &metadata.CreateModelAttributes{Attributes: []metadata.Attribute{pAttr}})
		if err != nil {
			blog.Errorf("failed to request coreService to create model attrs, err: %v, ObjectID: %s, input: %#v, "+
				"rid: %s", err, pAttr.ObjectID, pAttr, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rsppAttr.CCError(); err != nil {
			blog.Errorf("create model attrs failed, ObjectID: %s, input: %#v, rid: %s", pAttr.ObjectID, pAttr, kit.Rid)
			return nil, err
		}

		for _, exception := range rsppAttr.Data.Exceptions {
			return nil, kit.CCError.New(int(exception.Code), exception.Message)
		}

		if len(rsppAttr.Data.Repeated) > 0 {
			blog.Errorf("create model attrs failed, the attr is duplicated, ObjectID: %s, input: %#v, rid: %s",
				pAttr.ObjectID, pAttr, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
		}

		if len(rsppAttr.Data.Created) != 1 {
			blog.Errorf("create model attrs created amount error, ObjectID: %s, input: %#v, rid: %s",
				pAttr.ObjectID, pAttr, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
		}
		pAttr.ID = int64(rsppAttr.Data.Created[0].ID)

		keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(pAttr.ID)})
	}

	uni := metadata.ObjectUnique{
		ObjID:   obj.ObjectID,
		OwnerID: kit.SupplierAccount,
		Keys:    keys,
		Ispre:   false,
	}
	// NOTICE: 2021年03月29日  唯一索引与index.MainLineInstanceUniqueIndex,index.InstanceUniqueIndex定义强依赖
	// 原因：建立模型之前要将表和表中的索引提前建立，mongodb 4.2.6(4.4之前)事务中不能建表，事务操作表中数据操作和建表，建立索引为互斥操作。
	resp, err := assoc.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header,
		uni.ObjID, metadata.CreateModelAttrUnique{Data: uni})
	if err != nil {
		blog.Errorf("create unique for %s failed, err: %v, rid: %s", uni.ObjID, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoObjectUniqueCreateFailed)
	}
	if err = resp.CCError(); err != nil {
		blog.Errorf("create unique for %s failed, err: %v, rid: %s", uni.ObjID, err, kit.Rid)
		return nil, err
	}

	// generate audit log of object attribute group.
	audit := auditlog.NewObjectAuditLog(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.ID, nil)
	if err != nil {
		blog.Errorf("create object %s success, but generate audit log failed, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object %s success, but save audit log failed, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return nil, err
	}

	return obj, nil
}

// TODO need to be deleted after merge
func IsValid(kit *rest.Kit, isUpdate bool, data mapstr.MapStr) (*metadata.Object, error) {

	obj := new(metadata.Object)
	if err := mapstruct.Decode2Struct(data, obj); err != nil {
		blog.Errorf("parse object failed, err: %v, input: %#v, rid: %s", err, data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed)
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjectID) {

		if err := util.ValidModelIDField(data[metadata.ModelFieldObjectID],
			metadata.ModelFieldObjectID, kit.CCError); err != nil {
			blog.Errorf("failed to valid the object id(%s), rid: %s", metadata.ModelFieldObjectID, kit.Rid)
			return nil, err
		}
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjectName) {
		if err := util.ValidModelNameField(data[metadata.ModelFieldObjectName],
			metadata.ModelFieldObjectName, kit.CCError); err != nil {
			blog.Errorf("failed to valid the object name(%s), rid: %s", metadata.ModelFieldObjectName, kit.Rid)
			return nil, kit.CCError.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectName+" "+err.Error())
		}
	}

	if !isUpdate && !data.Exists(metadata.ModelFieldObjCls) {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjCls)
	}

	if !isUpdate && !metadata.IsCommon(obj.ObjectID) {
		return nil, kit.CCError.New(common.CCErrCommParamsIsInvalid,
			fmt.Sprintf("'%s' the built-in object id, please use a new one", obj.ObjectID))
	}

	return obj, nil
}

// FindInst search instance by condition
// TODO need to delete after merge
func (assoc *association) FindInst(kit *rest.Kit, objID string,
	cond *metadata.QueryInput) (*metadata.InstResult, error) {

	result := new(metadata.InstResult)
	switch objID {
	case common.BKInnerObjIDHost:
		rsp, err := assoc.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, cond)
		if err != nil {
			blog.Errorf("get host failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rsp.CCError(); err != nil {
			blog.Errorf("search object(%s) inst by the condition(%#v) failed, err: %v, rid: %s",
				objID, cond, err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Data.Count
		result.Info = rsp.Data.Info
		return result, nil

	default:
		input := &metadata.QueryCondition{Condition: cond.Condition, TimeCondition: cond.TimeCondition}
		input.Page.Start = cond.Start
		input.Page.Limit = cond.Limit
		input.Page.Sort = cond.Sort
		input.Fields = strings.Split(cond.Fields, ",")
		rsp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
		if err != nil {
			blog.Errorf("search instance failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rsp.CCError(); err != nil {
			blog.Errorf("search object(%s) inst by the condition(%#v) failed, err: %v, rid: %s",
				objID, cond, err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Data.Count
		result.Info = rsp.Data.Info
		return result, nil
	}
}

func (assoc *association) createInst(kit *rest.Kit, objID string, data mapstr.MapStr) (uint64, error) {
	cond := &metadata.CreateModelInstance{
		Data: data,
	}
	rsp, err := assoc.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, objID, cond)
	if err != nil {
		blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to create object instance ,err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	return rsp.Data.Created.ID, nil
}

func (assoc *association) updateMainlineAssociation(kit *rest.Kit, child mapstr.MapStr, parentID int64) error {

	childObj, err := child.String(common.BKObjIDField)
	if err != nil {
		blog.Errorf("get object id in child instance failed, child: %#v, err: %v, rid: %s", child, err, kit.Rid)
		return err
	}

	childID, err := child.Int64(metadata.GetInstIDFieldByObjID(childObj))
	if err != nil {
		return err
	}

	cond := mapstr.MapStr{metadata.GetInstIDFieldByObjID(childObj): childID}
	if metadata.IsCommon(childObj) {
		cond.Set(metadata.ModelFieldObjectID, childObj)
	}

	input := metadata.UpdateOption{
		Data: mapstr.MapStr{
			common.BKInstParentStr: parentID,
		},
		Condition: cond,
	}
	rsp, err := assoc.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, childObj, &input)
	if err != nil {
		blog.Errorf("failed to request object controller, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to update the association, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (assoc *association) setMainlineParentObject(kit *rest.Kit, obj *metadata.Object, relateToObjID string) error {
	cond := mapstr.MapStr{
		common.BKObjIDField:           obj.ObjectID,
		common.AssociationKindIDField: common.AssociationKindMainline,
	}

	resp, err := assoc.clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header,
		&metadata.DeleteOption{Condition: cond})
	if err != nil {
		blog.Errorf("update mainline object[%s] association to %s, search object association failed, err: %v, rid: %s",
			obj.ObjectID, relateToObjID, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = resp.CCError(); err != nil {
		blog.Errorf("update mainline object[%s] association to %s, search object association failed, err: %v, rid: %s",
			obj.ObjectID, relateToObjID, err, kit.Rid)
		return err
	}
	return assoc.createMainlineObjectAssociation(kit, obj, relateToObjID)
}

func (assoc *association) createMainlineObjectAssociation(kit *rest.Kit, obj *metadata.Object,
	relateToObjID string) error {
	objAsstID := fmt.Sprintf("%s_%s_%s", obj.ObjectID, common.AssociationKindMainline, relateToObjID)
	defined := false
	association := metadata.Association{
		OwnerID:              kit.SupplierAccount,
		AssociationName:      objAsstID,
		AssociationAliasName: objAsstID,
		ObjectID:             obj.ObjectID,
		// related to it's parent object id
		AsstObjID:  relateToObjID,
		AsstKindID: common.AssociationKindMainline,
		Mapping:    metadata.OneToOneMapping,
		OnDelete:   metadata.NoAction,
		IsPre:      &defined,
	}

	result, err := assoc.clientSet.CoreService().Association().CreateMainlineModelAssociation(kit.Ctx, kit.Header,
		&metadata.CreateModelAssociation{Spec: association})
	if err != nil {
		blog.Errorf("create mainline object association failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err = result.CCError(); err != nil {
		blog.Errorf("create mainline object association failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// TODO this func need to replace to ResetMainlineInstAssociation in inst/mainline_association
func (assoc *association) ResetMainlineInstAssociation(kit *rest.Kit, current *metadata.Object) error {

	cond := mapstr.New()
	if current.IsCommon() {
		cond.Set(common.BKObjIDField, current.ObjectID)
	}
	defaultCond := &metadata.QueryInput{Condition: cond}

	// 获取 current 模型的所有实例
	currentInsts, err := assoc.FindInst(kit, current.ObjectID, defaultCond)
	if err != nil {
		blog.Errorf("failed to find current object(%s) inst, err: %v, rid: %s", current.ObjectID, err, kit.Rid)
		return err
	}

	// 检查实例删除后，会不会出现重名冲突
	var canReset bool
	var repeatedInstName string
	if canReset, repeatedInstName, err = assoc.checkInstNameRepeat(kit, current, currentInsts.Info); err != nil {
		blog.Errorf("can not be reset, err: %+v, rid: %s", err, kit.Rid)
		return err
	}

	if canReset == false {
		blog.Errorf("can not be reset, inst name repeated, inst: %s, rid: %s", repeatedInstName, kit.Rid)
		errMsg := kit.CCError.CCError(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat).Error() +
			" " + repeatedInstName
		return kit.CCError.New(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat, errMsg)
	}

	// NEED FIX: 下面循环中的continue ，会在处理实例异常的时候跳过当前拓扑的处理，此方式可能会导致某个业务拓扑失败，但是不会影响所有。
	// 修改 currentInsts 所有孩子结点的父节点，为 currentInsts 的父节点，并删除 currentInsts
	for _, currentInst := range currentInsts.Info {
		instID, err := currentInst.Int64(metadata.GetInstIDFieldByObjID(current.ObjectID))
		if err != nil {
			blog.Errorf("failed to get the inst id from the inst(%#v), rid: %s", currentInst, kit.Rid)
			continue
		}

		parentID, err := currentInst.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("failed to get the object(%s) mainline parent id, the current inst(%v), err: %v, rid: %s",
				current.ObjectID, currentInst, err, kit.Rid)
			continue
		}

		// reset the child's parent
		children, err := assoc.getMainlineInst(kit, current, currentInst, true)
		if err != nil {
			blog.Errorf("failed to get the object(%s) mainline child inst, err: %v, rid: %s",
				current.ObjectID, err, kit.Rid)
			continue
		}
		for _, child := range children {
			// set the child's parent
			if err = assoc.SetMainlineParentInst(kit, child, parentID); err != nil {
				blog.Errorf("failed to set the object mainline child inst, err: %v, rid: %s", err, kit.Rid)
				continue
			}
		}

		// delete the current inst
		if err := assoc.DeleteMainlineInstWithID(kit, current, instID); err != nil {
			blog.Errorf("failed to delete the current inst(%#v), err: %v, rid: %s", currentInst, err, kit.Rid)
			continue
		}
	}

	return nil
}

// checkInstNameRepeat 检查如果将 currentInsts 都删除之后，拥有共同父节点的孩子结点会不会出现名字冲突
// 如果有冲突，返回 (false, 冲突实例名, nil)
func (assoc *association) checkInstNameRepeat(kit *rest.Kit, currentObj *metadata.Object,
	currentInsts []mapstr.MapStr) (canReset bool, repeatedInstName string, err error) {

	// TODO: 返回值将bool值与出错情况分开 (bool, error)
	instNames := map[string]bool{}
	for _, currInst := range currentInsts {
		currInstParentID, err := currInst.Int64(common.BKInstParentStr)
		if err != nil {
			return false, "", err
		}

		children, err := assoc.getMainlineInst(kit, currentObj, currInst, true)
		if err != nil {
			return false, "", err
		}

		for _, child := range children {
			instName, err := child.String(common.BKInstNameField)
			if err != nil {
				return false, "", err
			}
			key := fmt.Sprintf("%d_%s", currInstParentID, instName)
			if _, ok := instNames[key]; ok {
				return false, instName, nil
			}

			instNames[key] = true
		}
	}

	return true, "", nil
}

// TODO should be deleted after merge, and which call this func use DeleteMainlineInstWithID in inst/inst.go to replace
// DeleteMainlineInstWithID delete mainline instance by it's bk_inst_id
func (assoc *association) DeleteMainlineInstWithID(kit *rest.Kit, obj *metadata.Object, instID int64) error {

	// if this instance has been bind to a instance by the association, then this instance should not be deleted.
	cnt, err := assoc.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header, obj.ObjectID,
		&metadata.Condition{
			Condition: mapstr.MapStr{common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: obj.ObjectID, common.BKInstIDField: instID},
				{common.BKAsstObjIDField: obj.ObjectID, common.BKAsstInstIDField: instID},
			}},
		})
	if err != nil {
		blog.Errorf("count association by object(%s) failed, err: %s, rid: %s", obj.ObjectID, err, kit.Rid)
		return err
	}

	if err = cnt.CCError(); err != nil {
		blog.Errorf("count association by object(%s) failed, err: %s, rid: %s", obj.ObjectID, err, kit.Rid)
		return err
	}

	if cnt.Data.Count > 0 {
		return kit.CCError.CCError(common.CCErrorInstHasAsst)
	}

	// delete this instance now.
	delCond := mapstr.MapStr{obj.GetInstIDFieldName(): instID}
	if obj.IsCommon() {
		delCond.Set(common.BKObjIDField, obj.ObjectID)
	}

	// generate audit log.
	audit := auditlog.NewInstanceAudit(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, obj.GetObjectID(), delCond)
	if err != nil {
		blog.Errorf(" delete inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// to delete.
	ops := metadata.DeleteOption{
		Condition: delCond,
	}
	rsp, err := assoc.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, obj.ObjectID, &ops)
	if err != nil {
		blog.Errorf("request to delete instance by condition failed, cond: %#v, err: %v", ops, err)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to delete the object(%s) inst by the condition(%#v), err: %v",
			obj.ObjectID, ops, err)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("delete inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

// TODO should be deleted, and use DeleteAssociationWithPreCheck in model/association to replace where used this func
func (assoc *association) DeleteAssociation(kit *rest.Kit, cond mapstr.MapStr) error {
	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("get association with cond[%v] failed, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("get association with cond[%v] failed, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(rsp.Data.Info) < 1 {
		// we assume this association has already been deleted.
		blog.Errorf("can not get association with cond[%v], rid: %s", cond, kit.Rid)
		return nil
	}

	// a pre-defined association can not be updated.
	if rsp.Data.Info[0].IsPre != nil && *rsp.Data.Info[0].IsPre {
		blog.Errorf("it's a pre-defined association, can not be deleted, cond: %#v, rid: %s", cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoDeletePredefinedAssociation)
	}

	// delete the object association
	result, err := assoc.clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header,
		&metadata.DeleteOption{Condition: cond})
	if err != nil {
		blog.Errorf("delete object association failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = result.CCError(); err != nil {
		blog.Errorf("failed to create the association (%#v) , err: %s, rid: %s", cond, err, kit.Rid)
		return err
	}

	return nil
}

// TODO should be deleted after merge, model/object has been already implement
func (assoc *association) DeleteObject(kit *rest.Kit, id int64, needCheckInst bool) error {
	if id <= 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID)
	}

	// get model by id
	cond := mapstr.MapStr{
		metadata.ModelFieldID: id,
	}

	resp, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond},
	)
	if err != nil {
		blog.Errorf("find object failed, cond: %+v, err: %s, rid: %s", cond, err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = resp.CCError(); err != nil {
		blog.Errorf("failed to search the objects by the condition(%#v) , err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	objs := resp.Data.Info

	// shouldn't return 404 error here, legacy implements just ignore not found error
	if len(objs) == 0 {
		blog.V(3).Infof("object not found, condition: %v, rid: %s", cond, kit.Rid)
		return nil
	}
	if len(objs) > 1 {
		return kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	obj := objs[0]

	// check whether it can be deleted
	if needCheckInst {
		if err = assoc.CanDelete(kit, obj); err != nil {
			return err
		}
	}

	// generate audit log of object.
	audit := auditlog.NewObjectAuditLog(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.ID, &obj)
	if err != nil {
		blog.Errorf("generate audit log failed before delete object, objName: %s, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return err
	}

	// DeleteModelCascade 将会删除模型/模型属性/属性分组/唯一校验
	rsp, err := assoc.clientSet.CoreService().Model().DeleteModelCascade(kit.Ctx, kit.Header, id)
	if err != nil {
		blog.Errorf("failed to request the object controller, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to delete the object by the id(%d), err: %v, rid: %s", id, err, kit.Rid)
		return err
	}

	// save audit log.
	if err = audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("delete object %s success, save audit log failed, err: %v, rid: %s", obj.ObjectName, err, kit.Rid)
		return err
	}

	return nil
}

// TODO should be deleted after merge
// CanDelete return nil only when:
// 1. not inner model
// 2. model has no instances
// 3. model has no association with other model
func (assoc *association) CanDelete(kit *rest.Kit, targetObj metadata.Object) error {
	// step 1. ensure not inner model
	if common.IsInnerModel(targetObj.GetObjectID()) {
		return kit.CCError.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
	}

	cond := mapstr.New()
	if targetObj.IsCommon() {
		cond.Set(common.BKObjIDField, targetObj.ObjectID)
	}

	// step 2. ensure model has no instances
	input := &metadata.Condition{Condition: cond}
	findInstResponse, err := assoc.clientSet.CoreService().Instance().CountInstances(kit.Ctx, kit.Header,
		targetObj.ObjectID, input)

	if err != nil {
		blog.Errorf("failed to check if it (%s) has some insts, err: %v, rid: %s", targetObj.ObjectID, err, kit.Rid)
		return err
	}

	if err = findInstResponse.CCError(); err != nil {
		blog.Errorf("failed to check if it (%s) has some insts, err: %v, rid: %s", targetObj.ObjectID, err, kit.Rid)
		return err
	}

	if findInstResponse.Data.Count != 0 {
		blog.Errorf("the object [%s] has been instantiated and cannot be deleted, rid: %s",
			targetObj.ObjectID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete, targetObj.ObjectID)
	}

	// step 3. ensure model has no association with other model
	or := make([]interface{}, 0)
	or = append(or, mapstr.MapStr{common.BKObjIDField: targetObj.ObjectID})
	or = append(or, mapstr.MapStr{common.AssociatedObjectIDField: targetObj.ObjectID})

	assocResult, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: map[string]interface{}{common.BKDBOR: or}})
	if err != nil {
		blog.Errorf("check object[%s] can be deleted, but get object associate info failed, err: %v, rid: %s",
			targetObj.ObjectID, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = assocResult.CCError(); err != nil {
		blog.Errorf("get object[%s] associate info failed, err: %v, rid: %s", targetObj.ObjectID, err, kit.Rid)
		return kit.CCError.Error(assocResult.Code)
	}

	if len(assocResult.Data.Info) != 0 {
		blog.Errorf("object[%s] has already associate to another one., rid: %s", targetObj.ObjectID, kit.Rid)
		return kit.CCError.Error(common.CCErrorTopoObjectHasAlreadyAssociated)
	}

	return nil
}
