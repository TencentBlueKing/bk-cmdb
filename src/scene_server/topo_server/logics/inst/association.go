package inst

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	// SearchInstanceAssociations searches object instance associations.
	SearchInstanceAssociations(kit *rest.Kit, objID string, input *metadata.CommonSearchFilter) (
		[]metadata.InstAsst, error)
	// CountInstanceAssociations counts object instance associations num.
	CountInstanceAssociations(kit *rest.Kit, objID string, input *metadata.CommonCountFilter) (
		*metadata.CommonCountResult, error)
	// SearchInstAssociationUIList instance association data related to instances, return by pagination
	SearchInstAssociationUIList(kit *rest.Kit, objID string, query *metadata.QueryCondition) (
		*metadata.SearchInstAssociationListResult, uint64, error)
	// SearchInstAssociationSingleObjectInstInfo 与实例有关系的实例关系数据,以分页的方式返回
	SearchInstAssociationSingleObjectInstInfo(kit *rest.Kit, objID string, query *metadata.QueryCondition,
		isTargetObject bool) ([]metadata.InstBaseInfo, uint64, error)
	// CreateInstanceAssociation create an association between instances
	CreateInstanceAssociation(kit *rest.Kit, request *metadata.CreateAssociationInstRequest) (
		*metadata.RspID, error)
	// CreateManyInstAssociation create many associations between instances
	CreateManyInstAssociation(kit *rest.Kit, request *metadata.CreateManyInstAsstRequest) (
		*metadata.CreateManyInstAsstResultDetail, error)
	// DeleteInstAssociation delete association between instances
	DeleteInstAssociation(kit *rest.Kit, objID string, asstIDList []int64) (uint64, error)
	// CheckAssociations returns error if the instances has associations with exist instances, clear dirty associations
	CheckAssociations(*rest.Kit, string, []int64) error

	// SearchMainlineAssociationInstTopo search mainline association topo by objID and instID
	SearchMainlineAssociationInstTopo(kit *rest.Kit, objID string, instID int64,
		withStatistics bool, withDefault bool) ([]*metadata.TopoInstRst, errors.CCError)
	// ResetMainlineInstAssociation reset mainline instance association
	ResetMainlineInstAssociation(kit *rest.Kit, currentObjID, childObjID string) error
	// SetMainlineInstAssociation set mainline instance association by parent object and current object
	SetMainlineInstAssociation(kit *rest.Kit, parentObjID, childObjID, currObjID, currObjName string) ([]int64, error)
	// TopoNodeHostAndSerInstCount get topo node host and service instance count
	TopoNodeHostAndSerInstCount(kit *rest.Kit, input *metadata.HostAndSerInstCountOption) (
		[]*metadata.TopoNodeHostAndSerInstCount, errors.CCError)

	// SetProxy proxy the interface
	SetProxy(inst InstOperationInterface)
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
	inst        InstOperationInterface
}

// SetProxy proxy the interface
func (assoc *association) SetProxy(inst InstOperationInterface) {
	assoc.inst = inst
}

// SearchInstanceAssociations searches object instance associations.
func (assoc *association) SearchInstanceAssociations(kit *rest.Kit, objID string, input *metadata.CommonSearchFilter) (
	[]metadata.InstAsst, error) {

	// search conditions.
	cond, err := input.GetConditions()
	if err != nil {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, err)
	}

	conditions := &metadata.InstAsstQueryCondition{
		ObjID: objID,
		Cond: metadata.QueryCondition{
			Fields:         input.Fields,
			Condition:      cond,
			Page:           input.Page,
			DisableCounter: true,
		},
	}

	// search object instance associations.
	resp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, conditions)
	if err != nil {
		blog.Errorf("search instance associations failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return resp.Info, nil
}

// CountInstanceAssociations counts object instance associations num.
func (assoc *association) CountInstanceAssociations(kit *rest.Kit, objID string, input *metadata.CommonCountFilter) (
	*metadata.CommonCountResult, error) {

	// count conditions.
	cond, err := input.GetConditions()
	if err != nil {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, err)
	}
	cond[common.BKObjIDField] = objID

	conditions := &metadata.Condition{
		Condition: cond,
	}

	// count object instance associations num.
	resp, err := assoc.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header, objID,
		conditions)
	if err != nil {
		blog.Errorf("count instance associations failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return &metadata.CommonCountResult{Count: resp.Count}, nil
}

// SearchInstAssociationUIList 与实例有关系的实例关系数据,以分页的方式返回
func (assoc *association) SearchInstAssociationUIList(kit *rest.Kit, objID string, query *metadata.QueryCondition) (
	*metadata.SearchInstAssociationListResult, uint64, error) {

	queryCond := &metadata.InstAsstQueryCondition{ObjID: objID}
	if query != nil {
		queryCond.Cond = *query
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("search instance association failed, query: %#v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, 0, err
	}

	objIDInstIDMap := make(map[string][]int64, 0)
	objSrcAsstArr := make([]metadata.InstAsst, 0)
	objDstAsstArr := make([]metadata.InstAsst, 0)
	for _, instAsst := range rsp.Info {
		objIDInstIDMap[instAsst.ObjectID] = append(objIDInstIDMap[instAsst.ObjectID], instAsst.InstID)
		objIDInstIDMap[instAsst.AsstObjectID] = append(objIDInstIDMap[instAsst.AsstObjectID], instAsst.AsstInstID)
		if instAsst.ObjectID == objID {
			objSrcAsstArr = append(objSrcAsstArr, instAsst)
		} else {
			objDstAsstArr = append(objDstAsstArr, instAsst)

		}
	}

	instInfo := make(map[string][]mapstr.MapStr, 0)
	for instObjID, instIDArr := range objIDInstIDMap {

		idField := metadata.GetInstIDFieldByObjID(instObjID)
		input := &metadata.QueryCondition{
			Condition: mapstr.MapStr{
				idField: mapstr.MapStr{common.BKDBIN: instIDArr},
			},
			Page: metadata.BasePage{
				Start: 0,
				Limit: common.BKNoLimit,
			},
			Fields: []string{metadata.GetInstNameFieldName(instObjID), idField},
		}
		instResp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, instObjID, input)
		if err != nil {
			blog.Errorf("search instance failed, query: %#v, err: %v, rid: %s", query, err, kit.Rid)
			return nil, 0, err
		}

		instInfo[instObjID] = instResp.Info

	}

	result := &metadata.SearchInstAssociationListResult{}
	result.Inst = instInfo
	result.Association.Src = objSrcAsstArr
	result.Association.Dst = objDstAsstArr

	return result, rsp.Count, nil
}

// checkInstAsstMapping use to check if instance association mapping correct, used by CreateInstanceAssociation
func (assoc *association) checkInstAsstMapping(kit *rest.Kit, objID string, mapping metadata.AssociationMapping,
	input *metadata.CreateAssociationInstRequest) error {

	tableName := common.GetObjectInstAsstTableName(objID, kit.SupplierAccount)
	switch mapping {
	case metadata.OneToOneMapping:
		// search instances belongs to this association.
		queryFilter := []map[string]interface{}{
			{
				common.AssociationObjAsstIDField: input.ObjectAsstID,
				common.BKInstIDField:             input.InstID,
			},
			{
				common.AssociationObjAsstIDField: input.ObjectAsstID,
				common.BKAsstInstIDField:         input.AsstInstID,
			},
		}
		instCnt, err := assoc.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, tableName,
			queryFilter)
		if err != nil {
			blog.Errorf("check instance with cond[%#v] failed, err: %v, rid: %s", queryFilter, err, kit.Rid)
			return err
		}

		for _, cnt := range instCnt {
			if cnt >= 1 {
				return kit.CCError.Error(common.CCErrorTopoCreateMultipleInstancesForOneToOneAssociation)
			}
		}
	case metadata.OneToManyMapping:
		queryFilter := []map[string]interface{}{
			{
				common.AssociationObjAsstIDField: input.ObjectAsstID,
				common.BKAsstInstIDField:         input.AsstInstID,
			},
		}
		instCnt, err := assoc.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, tableName,
			queryFilter)
		if err != nil {
			blog.Errorf("check instance with cond[%#v] failed, err: %v, rid: %s", queryFilter, err, kit.Rid)
			return err
		}

		if instCnt[0] >= 1 {
			return kit.CCError.Error(common.CCErrorTopoCreateMultipleInstancesForOneToOneAssociation)
		}

	default:
		// after all the check, new association instance can be created.
	}

	return nil
}

// CreateInstanceAssociation create an association between instances
func (assoc *association) CreateInstanceAssociation(kit *rest.Kit, request *metadata.CreateAssociationInstRequest) (
	*metadata.RspID, error) {

	cond := &metadata.QueryCondition{Condition: mapstr.MapStr{common.AssociationObjAsstIDField: request.ObjectAsstID}}
	result, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("search object association with cond[%#v] failed, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	if len(result.Info) == 0 {
		blog.Errorf("can not find object association[%s]. rid: %s", request.ObjectAsstID, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrorTopoObjectAssociationNotExist)
	}

	if err := assoc.checkInstAsstMapping(kit, result.Info[0].ObjectID, result.Info[0].Mapping, request); err != nil {
		blog.Errorf("check mapping failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	input := metadata.CreateOneInstanceAssociation{
		Data: metadata.InstAsst{
			ObjectAsstID:      request.ObjectAsstID,
			InstID:            request.InstID,
			AsstInstID:        request.AsstInstID,
			ObjectID:          result.Info[0].ObjectID,
			AsstObjectID:      result.Info[0].AsstObjID,
			AssociationKindID: result.Info[0].AsstKindID,
		},
	}
	createResult, err := assoc.clientSet.CoreService().Association().CreateInstAssociation(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("create instance association failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, err
	}

	instAssociationID := int64(createResult.Created.ID)
	input.Data.ID = int64(createResult.Created.ID)

	// generate audit log.
	audit := auditlog.NewInstanceAssociationAudit(assoc.clientSet.CoreService())
	generateAuditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParam, instAssociationID, result.Info[0].ObjectID, &input.Data)
	if err != nil {
		blog.Errorf(" delete inst asst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, *auditLog)
	if err != nil {
		blog.Errorf("delete inst asst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return &metadata.RspID{ID: instAssociationID}, err
}

// CreateManyInstAssociation create many associations between instances
func (assoc *association) CreateManyInstAssociation(kit *rest.Kit, request *metadata.CreateManyInstAsstRequest) (
	*metadata.CreateManyInstAsstResultDetail, error) {

	rawErr := request.Validate()
	if rawErr.ErrCode != 0 {
		blog.Errorf("validate parameter failed, err: %v, rid: %s", rawErr.ToCCError(kit.CCError), kit.Rid)
		return nil, rawErr.ToCCError(kit.CCError)
	}

	param := &metadata.CreateManyInstanceAssociation{}
	for _, item := range request.Details {
		param.Datas = append(param.Datas, metadata.InstAsst{
			InstID:       item.InstID,
			ObjectID:     request.ObjectID,
			AsstInstID:   item.AsstInstID,
			AsstObjectID: request.AsstObjectID,
			ObjectAsstID: request.ObjectAsstID,
		})
	}

	res, err := assoc.clientSet.CoreService().Association().CreateManyInstAssociation(kit.Ctx, kit.Header, param)
	if err != nil {
		blog.Errorf("create many instance association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	resp := metadata.NewManyInstAsstResultDetail()
	for _, item := range res.Created {
		resp.SuccessCreated[item.OriginIndex] = int64(item.ID)
	}

	for _, item := range res.Repeated {
		itemObjID, _ := item.Data.Get(common.BKObjIDField)
		itemAsstObjID, _ := item.Data.Get(common.BKAsstObjIDField)
		resp.Error[item.OriginIndex] = kit.CCError.CCErrorf(common.CCErrTopoAssociationAlreadyExist,
			itemObjID, itemAsstObjID).Error()
	}

	for _, item := range res.Exceptions {
		resp.Error[item.OriginIndex] = item.Message
	}

	if len(resp.SuccessCreated) == 0 {
		return resp, nil
	}

	// generate audit log.
	audit := auditlog.NewInstanceAssociationAudit(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	var auditList []metadata.AuditLog
	for _, asstID := range resp.SuccessCreated {
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, asstID, request.ObjectID, nil)
		if err != nil {
			blog.Errorf("generate audit log failed, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, kit.CCError.Error(common.CCErrAuditGenerateLogFailed)
		}
		auditList = append(auditList, *auditLog)
	}
	// save audit log.
	err = audit.SaveAuditLog(kit, auditList...)
	if err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return resp, nil
}

// DeleteInstAssociation method will remove docs from both source-asst-collection and target-asst-collection
// which is atomicity.
func (assoc *association) DeleteInstAssociation(kit *rest.Kit, objID string, asstIDList []int64) (uint64, error) {

	if len(asstIDList) == 0 {
		blog.Errorf("association ID list can not be empty, rid: %s", kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "id")
	}

	// search association Instances
	cond := mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: asstIDList}}
	searchCondition := &metadata.InstAsstQueryCondition{
		Cond:  metadata.QueryCondition{Condition: cond, DisableCounter: true},
		ObjID: objID,
	}
	data, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, searchCondition)
	if err != nil {
		blog.Errorf("get instance association failed, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	input := metadata.InstAsstDeleteOption{Opt: metadata.DeleteOption{Condition: cond}, ObjID: objID}
	rsp, err := assoc.clientSet.CoreService().Association().DeleteInstAssociation(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("delete instance association failed, err: %s, input: %#v, rid: %s", err, input, kit.Rid)
		return 0, err
	}

	if rsp.Count == 0 {
		return rsp.Count, nil
	}

	// generate audit log.
	audit := auditlog.NewInstanceAssociationAudit(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditList := make([]metadata.AuditLog, 0)
	for i, asstID := range asstIDList {
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, asstID, objID, &data.Info[i])
		if err != nil {
			blog.Errorf("delete instance association failed, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
			return 0, err
		}
		auditList = append(auditList, *auditLog)
	}
	// save audit log.
	err = audit.SaveAuditLog(kit, auditList...)
	if err != nil {
		blog.Errorf("delete instance association failed, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.CCError(common.CCErrAuditSaveLogFailed)
	}

	return rsp.Count, nil
}

// CheckAssociations returns error if the instances has associations with exist instances, clear dirty associations
func (assoc *association) CheckAssociations(kit *rest.Kit, objectID string, instIDs []int64) error {
	if len(instIDs) == 0 {
		return nil
	}

	// get all associations for the instances
	instAsstCond := &metadata.InstAsstQueryCondition{
		Cond: metadata.QueryCondition{Condition: mapstr.MapStr{
			common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: objectID, common.BKInstIDField: mapstr.MapStr{common.BKDBIN: instIDs}},
				{common.BKAsstObjIDField: objectID, common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: instIDs}},
			},
		}},
		ObjID: objectID,
	}
	associations, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header,
		instAsstCond)
	if err != nil {
		blog.Errorf("search instance associations failed, condition: %#v, err: %v, rid: %s", instAsstCond, err, kit.Rid)
		return err
	}

	if len(associations.Info) == 0 {
		return nil
	}

	instIDExistsMap := make(map[int64]bool)
	for _, instID := range instIDs {
		instIDExistsMap[instID] = true
	}

	// get all associated inst IDs grouped by object ID, then check if any inst exists, clear not exist one's assts
	asstObjInstIDsMap := make(map[string][]int64)
	for _, asst := range associations.Info {
		if asst.ObjectID == objectID && instIDExistsMap[asst.InstID] {
			asstObjInstIDsMap[asst.AsstObjectID] = append(asstObjInstIDsMap[asst.AsstObjectID], asst.AsstInstID)
		} else if asst.AsstObjectID == objectID && instIDExistsMap[asst.AsstInstID] {
			asstObjInstIDsMap[asst.ObjectID] = append(asstObjInstIDsMap[asst.ObjectID], asst.InstID)
		}
	}

	for asstObjID, asstInstIDs := range asstObjInstIDsMap {
		query := &metadata.Condition{
			Condition: mapstr.MapStr{
				common.GetInstIDField(asstObjID): mapstr.MapStr{common.BKDBIN: asstInstIDs},
			},
		}
		asstInstCnt, err := assoc.clientSet.CoreService().Instance().CountInstances(kit.Ctx, kit.Header, asstObjID,
			query)
		if err != nil {
			blog.ErrorJSON("check instance existence failed, err: %s, query: %s, rid: %s", err, query, kit.Rid)
			return err
		}

		if asstInstCnt.Count > 0 {
			return kit.CCError.CCError(common.CCErrorInstHasAsst)
		}

		delOpt := &metadata.InstAsstDeleteOption{
			Opt: metadata.DeleteOption{Condition: mapstr.MapStr{
				common.BKDBOR: []mapstr.MapStr{
					{common.BKObjIDField: asstObjID, common.BKInstIDField: mapstr.MapStr{common.BKDBIN: asstInstIDs}},
					{common.BKAsstObjIDField: asstObjID, common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: asstInstIDs}},
				},
			}},
			ObjID: asstObjID,
		}
		_, err = assoc.clientSet.CoreService().Association().DeleteInstAssociation(kit.Ctx, kit.Header, delOpt)
		if err != nil {
			blog.ErrorJSON("delete dirty assts failed, err: %s, cond: %s, rid: %s", err, delOpt, kit.Rid)
			return err
		}
	}
	return nil
}

// SearchInstAssociationSingleObjectInstInfo 与实例有关系的实例关系数据,以分页的方式返回
// objID 根据条件查询出来关联关系，需要返回实例信息（实例名，实例ID）的模型ID
func (assoc *association) SearchInstAssociationSingleObjectInstInfo(kit *rest.Kit, objID string,
	query *metadata.QueryCondition, isTargetObject bool) ([]metadata.InstBaseInfo, uint64, error) {

	queryCond := &metadata.InstAsstQueryCondition{ObjID: objID}
	if query != nil {
		queryCond.Cond = *query
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, queryCond)
	if nil != err {
		blog.Errorf("ReadInstAssociation http do error, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, 0, kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if len(rsp.Info) == 0 {
		return nil, 0, nil
	}

	objIDInstIDArr := make([]int64, 0)
	for _, instAsst := range rsp.Info {
		if isTargetObject {
			objIDInstIDArr = append(objIDInstIDArr, instAsst.InstID)
		} else {
			objIDInstIDArr = append(objIDInstIDArr, instAsst.AsstInstID)
		}
	}

	idField := metadata.GetInstIDFieldByObjID(objID)
	nameField := metadata.GetInstNameFieldName(objID)
	input := &metadata.QueryCondition{
		Condition: mapstr.MapStr{idField: mapstr.MapStr{common.BKDBIN: objIDInstIDArr}},
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
		Fields:    []string{nameField, idField},
	}
	instResp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
	if err != nil {
		blog.Errorf("search object(%s) inst failed, input:%s, err: %v, rid: %s", objID, err, input, kit.Rid)
		return nil, 0, err
	}

	result := make([]metadata.InstBaseInfo, 0)
	for _, row := range instResp.Info {
		id, err := row.Int64(idField)
		if err != nil {
			blog.Errorf("get inst id field(%s) failed. err: %v, inst: %s, rid: %s", idField, err, row, kit.Rid)
			// CCErrCommInstFieldConvertFail  convert %s  field %s to %s error %s
			return nil, 0, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, objID, idField, "int", err.Error())
		}
		name, err := row.String(nameField)
		if err != nil {
			blog.Errorf("get inst name field(%s) failed. err: %v, inst: %s, rid: %s", nameField, err, row, kit.Rid)
			// CCErrCommInstFieldConvertFail  convert %s  field %s to %s error %s
			return nil, 0, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, objID, nameField, "string",
				err.Error())
		}

		result = append(result, metadata.InstBaseInfo{ID: id, Name: name})
	}

	return result, rsp.Count, nil
}
