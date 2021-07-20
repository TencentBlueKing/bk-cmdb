package inst

import (
	"fmt"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	// SearchInstanceAssociations searches object instance associations.
	SearchInstanceAssociations(kit *rest.Kit, objID string,
		input *metadata.CommonSearchFilter) (*metadata.CommonSearchResult, error)

	// CountInstanceAssociations counts object instance associations num.
	CountInstanceAssociations(kit *rest.Kit, objID string,
		input *metadata.CommonCountFilter) (*metadata.CommonCountResult, error)

	// SearchInstAssociation search instance association by metadata.SearchAssociationInstRequest
	SearchInstAssociation(kit *rest.Kit, request *metadata.SearchAssociationInstRequest) (
		resp *metadata.SearchAssociationInstResult, err error)

	// SearchInstAssociationUIList instance association data related to instances, return by pagination
	SearchInstAssociationUIList(kit *rest.Kit, objID string,
		query *metadata.QueryCondition) (result interface{}, asstCnt uint64, err error)

	// CreateInstanceAssociation create an association between instances
	CreateInstanceAssociation(kit *rest.Kit, request *metadata.CreateAssociationInstRequest) (
		resp *metadata.CreateAssociationInstResult, err error)

	// CreateManyInstAssociation create many associations between instances
	CreateManyInstAssociation(kit *rest.Kit,
		request *metadata.CreateManyInstAsstRequest) (*metadata.CreateManyInstAsstResultDetail, error)

	// DeleteInstAssociation delete association between instances
	DeleteInstAssociation(kit *rest.Kit, objID string,
		asstIDList []int64) (resp *metadata.DeleteAssociationInstResult, err error)

	// CheckAssociation check a association, clear dirty associations
	CheckAssociation(kit *rest.Kit, objectID string, instID int64) error

	// CheckAssociations returns error if the instances has associations with exist instances, clear dirty associations
	CheckAssociations(kit *rest.Kit, objectID string, instIDs []int64) error
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

// SearchInstanceAssociations searches object instance associations.
func (assoc *association) SearchInstanceAssociations(kit *rest.Kit, objID string,
	input *metadata.CommonSearchFilter) (*metadata.CommonSearchResult, error) {

	// search conditions.
	cond, err := input.GetConditions()
	if err != nil {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, err)
	}
	cond[common.BKObjIDField] = objID

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
		blog.Errorf("search instance associations failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result || resp.Code != 0 {
		return nil, kit.CCError.New(resp.Code, resp.ErrMsg)
	}

	result := &metadata.CommonSearchResult{}
	for idx := range resp.Data.Info {
		result.Info = append(result.Info, &resp.Data.Info[idx])
	}

	return result, nil
}

// CountInstanceAssociations counts object instance associations num.
func (assoc *association) CountInstanceAssociations(kit *rest.Kit, objID string,
	input *metadata.CommonCountFilter) (*metadata.CommonCountResult, error) {

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
		blog.Errorf("count instance associations failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result || resp.Code != 0 {
		return nil, kit.CCError.New(resp.Code, resp.ErrMsg)
	}

	return &metadata.CommonCountResult{Count: resp.Data.Count}, nil
}

// SearchInstAssociation search instance association by metadata.SearchAssociationInstRequest
func (assoc *association) SearchInstAssociation(kit *rest.Kit, request *metadata.SearchAssociationInstRequest) (
	resp *metadata.SearchAssociationInstResult, err error) {
	queryCond := &metadata.InstAsstQueryCondition{
		Cond:  metadata.QueryCondition{Condition: request.Condition},
		ObjID: request.ObjID,
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		return nil, err
	}

	resp = &metadata.SearchAssociationInstResult{BaseResp: rsp.BaseResp, Data: []*metadata.InstAsst{}}
	for index := range rsp.Data.Info {
		resp.Data = append(resp.Data, &rsp.Data.Info[index])
	}

	return resp, nil
}

// SearchInstAssociationUIList 与实例有关系的实例关系数据,以分页的方式返回
func (assoc *association) SearchInstAssociationUIList(kit *rest.Kit, objID string,
	query *metadata.QueryCondition) (result interface{}, asstCnt uint64, err error) {
	queryCond := &metadata.InstAsstQueryCondition{
		ObjID: objID,
	}
	if query != nil {
		queryCond.Cond = *query
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, queryCond)
	if nil != err {
		blog.Errorf("search instance association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, 0, kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.ErrorJSON("search instance association failed, query: %s, response: %s, rid: %s", query, rsp, kit.Rid)
		return nil, 0, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	objIDInstIDMap := make(map[string][]int64, 0)
	var objSrcAsstArr []metadata.InstAsst
	var objDstAsstArr []metadata.InstAsst
	for _, instAsst := range rsp.Data.Info {
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
		cond := condition.CreateCondition()
		cond.Field(idField).In(instIDArr)
		input := &metadata.QueryCondition{
			Condition: cond.ToMapStr(),
			Page: metadata.BasePage{
				Start: 0,
				Limit: common.BKNoLimit,
			},
			Fields: []string{metadata.GetInstNameFieldName(instObjID), idField},
		}
		instResp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, instObjID, input)
		if err != nil {
			blog.Errorf("search instance failed, err: %s, input:%s, rid: %s", err.Error(), input, kit.Rid)
			return nil, 0, kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
		}
		if !instResp.Result {
			blog.ErrorJSON("search instance failed, query: %s, response: %s, rid: %s", query, rsp, kit.Rid)
			return nil, 0, kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
		instInfo[instObjID] = instResp.Data.Info

	}
	instAsstMap := map[string][]metadata.InstAsst{
		"src": objSrcAsstArr,
		"dst": objDstAsstArr,
	}

	result = mapstr.MapStr{
		"association": instAsstMap,
		"instance":    instInfo,
	}

	return result, rsp.Data.Count, nil
}

func (assoc *association) CreateInstanceAssociation(kit *rest.Kit, request *metadata.CreateAssociationInstRequest) (
	resp *metadata.CreateAssociationInstResult, err error) {

	cond := mapstr.MapStr{common.AssociationObjAsstIDField: request.ObjectAsstID}
	result, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search object association with cond[%v] failed, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	if !result.Result {
		blog.Errorf("search object association with cond[%v] failed, err: %s, rid: %s", cond, result.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		blog.Errorf("can not find object association[%s]. rid: %s", request.ObjectAsstID, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrorTopoObjectAssociationNotExist)
	}

	objectAsst := result.Data.Info[0]

	objID := objectAsst.ObjectID
	asstObjID := objectAsst.AsstObjID

	switch result.Data.Info[0].Mapping {
	case metadata.OneToOneMapping:
		// search instances belongs to this association.
		cond := mapstr.MapStr{
			common.AssociationObjAsstIDField: request.ObjectAsstID,
			common.BKInstIDField:             request.InstID,
		}

		instance, err := assoc.SearchInstAssociation(kit,
			&metadata.SearchAssociationInstRequest{Condition: cond, ObjID: objID})
		if err != nil {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %v, rid: %s",
				cond, err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instance.Result {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %s, rid: %s",
				cond, instance.ErrMsg, kit.Rid)
			return nil, kit.CCError.New(instance.Code, instance.ErrMsg)
		}
		if len(instance.Data) >= 1 {
			return nil, kit.CCError.Error(common.CCErrorTopoCreateMultipleInstancesForOneToOneAssociation)
		}

		cond = mapstr.MapStr{
			common.AssociationObjAsstIDField: request.ObjectAsstID,
			common.BKAsstInstIDField:         request.AsstInstID,
		}

		instance, err = assoc.SearchInstAssociation(kit,
			&metadata.SearchAssociationInstRequest{Condition: cond, ObjID: objID})
		if err != nil {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %v, rid: %s",
				cond, err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instance.Result {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %s, rid: %s",
				cond, instance.ErrMsg, kit.Rid)
			return nil, kit.CCError.New(instance.Code, instance.ErrMsg)
		}
		if len(instance.Data) >= 1 {
			return nil, kit.CCError.Error(common.CCErrorTopoCreateMultipleInstancesForOneToOneAssociation)
		}
	case metadata.OneToManyMapping:

		cond = mapstr.MapStr{
			common.AssociationObjAsstIDField: request.ObjectAsstID,
			common.BKAsstInstIDField:         request.AsstInstID,
		}

		instance, err := assoc.SearchInstAssociation(kit,
			&metadata.SearchAssociationInstRequest{Condition: cond, ObjID: objID})
		if err != nil {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %v, rid: %s",
				cond, err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instance.Result {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %s, rid: %s",
				cond, instance.ErrMsg, kit.Rid)
			return nil, kit.CCError.New(instance.Code, instance.ErrMsg)
		}
		if len(instance.Data) >= 1 {
			return nil, kit.CCError.Error(common.CCErrorTopoCreateMultipleInstancesForOneToManyAssociation)
		}

	default:
		// after all the check, new association instance can be created.
	}

	input := metadata.CreateOneInstanceAssociation{
		Data: metadata.InstAsst{
			ObjectAsstID:      request.ObjectAsstID,
			InstID:            request.InstID,
			AsstInstID:        request.AsstInstID,
			ObjectID:          objID,
			AsstObjectID:      asstObjID,
			AssociationKindID: objectAsst.AsstKindID,
		},
	}
	createResult, err := assoc.clientSet.CoreService().Association().CreateInstAssociation(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("create instance association failed, do coreservice create failed, err: %+v, rid: %s",
			err, kit.Rid)
		return nil, err
	}
	if err := createResult.CCError(); err != nil {
		blog.Errorf("create instance association failed, do coreservice create failed, err: %+v, rid: %s",
			err, kit.Rid)
		return nil, err
	}

	resp = &metadata.CreateAssociationInstResult{BaseResp: createResult.BaseResp}
	instanceAssociationID := int64(createResult.Data.Created.ID)
	resp.Data.ID = instanceAssociationID

	curData := mapstr.NewFromStruct(input.Data, "json")
	curData.Set("name", objectAsst.AssociationAliasName)

	// generate audit log.
	audit := auditlog.NewInstanceAssociationAudit(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, instanceAssociationID, objID, nil)
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

	return resp, err
}

func (assoc *association) CreateManyInstAssociation(kit *rest.Kit,
	request *metadata.CreateManyInstAsstRequest) (*metadata.CreateManyInstAsstResultDetail, error) {
	rawErr := request.Validate()
	if rawErr.ErrCode != 0 {
		blog.Errorf("validate parameter failed, err: %s, rid: %s", rawErr.ToCCError(kit.CCError).Error(), kit.Rid)
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
		blog.Errorf("create many instance association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	if !res.Result {
		blog.Errorf("create many instance association failed, err: %s, rid: %s", res.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(res.Code, res.ErrMsg)
	}

	resp := metadata.NewManyInstAsstResultDetail()
	for _, item := range res.Data.Created {
		resp.SuccessCreated[item.OriginIndex] = int64(item.ID)
	}

	for _, item := range res.Data.Repeated {
		itemObjID, _ := item.Data.Get(common.BKObjIDField)
		itemAsstObjID, _ := item.Data.Get(common.BKAsstObjIDField)
		resp.Error[item.OriginIndex] = kit.CCError.CCErrorf(common.CCErrTopoAssociationAlreadyExist,
			itemObjID, itemAsstObjID).Error()
	}

	for _, item := range res.Data.Exceptions {
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
		blog.Errorf("save audit log failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return resp, nil
}

// DeleteInstAssociation method will remove docs from both source-asst-collection and target-asst-collection, which is atomicity.
func (assoc *association) DeleteInstAssociation(kit *rest.Kit, objID string,
	asstIDList []int64) (resp *metadata.DeleteAssociationInstResult, err error) {

	// asstIDList check duplicate
	idMap := make(map[int64]struct{})
	for _, id := range asstIDList {
		if id <= 0 {
			blog.ErrorJSON("delete instance association failed, input id list contains illegal id %d, "+
				"kit: %+v, rid: %s", id, kit, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, id)
		}
		if _, exists := idMap[id]; exists {
			blog.ErrorJSON("delete instance association failed, input id list contains duplicate id %d, "+
				"kit: %+v, rid: %s", kit, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, id)
		}
		idMap[id] = struct{}{}
	}
	// search association Instances
	cond := mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: asstIDList}}
	searchCondition := metadata.InstAsstQueryCondition{
		Cond: metadata.QueryCondition{
			Condition:      cond,
			DisableCounter: true,
		},
		ObjID: objID,
	}
	data, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, &searchCondition)
	if err != nil {
		blog.Errorf("delete instance association failed, get instance association failed, "+
			"kit: %+v, err: %+v, rid: %s", kit, err, kit.Rid)
		return nil, err
	}
	if len(data.Data.Info) != len(idMap) {
		errStr := ""
		for _, dataRst := range data.Data.Info {
			delete(idMap, dataRst.ID)
		}
		for idNotExist := range idMap {
			errStr = fmt.Sprintf("%s,%d", errStr, idNotExist)
		}
		errStr = strings.TrimPrefix(errStr, ",")
		blog.Errorf("delete instance association failed, %s in ID list does not exists, "+
			"searchCondition: %+v, rid: %s", errStr, searchCondition, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, errStr)
	}

	// get different association models, check whether they are exists.
	objAsstMap := make(map[string]struct{})
	objAsstList := []string{}
	for _, instAsst := range data.Data.Info {
		objAsstMap[instAsst.ObjectAsstID] = struct{}{}
	}

	for objAsstID := range objAsstMap {
		objAsstList = append(objAsstList, objAsstID)
	}

	searchCond := mapstr.MapStr{common.AssociationObjAsstIDField: mapstr.MapStr{common.BKDBIN: objAsstList}}
	// NOTE this interface call maybe can change into SearchObject function
	assInfoResult, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: searchCond})

	if err != nil {
		blog.Errorf("search object association with cond[%v] failed, err: %v, rid: %s", searchCond, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !assInfoResult.Result {
		blog.Errorf("search object association with cond[%v] failed, err: %s, rid: %s",
			searchCond, assInfoResult.ErrMsg, kit.Rid)
		return nil, assInfoResult.CCError()
	}
	if len(assInfoResult.Data.Info) != len(objAsstList) {
		blog.Errorf("got unexpected number of model associations %d which should be %d, searchCondition: %+v, rid: %s",
			len(assInfoResult.Data.Info), len(objAsstList), searchCond, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommNotFound)
	}

	input := metadata.InstAsstDeleteOption{
		Opt: metadata.DeleteOption{
			Condition: cond,
		},
		ObjID: objID,
	}
	rsp, err := assoc.clientSet.CoreService().Association().DeleteInstAssociation(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.ErrorJSON("delete instance association failed, err: %s, input: %s, rid: %s", err, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	resp = &metadata.DeleteAssociationInstResult{
		BaseResp: rsp.BaseResp,
	}

	// generate audit log.
	audit := auditlog.NewInstanceAssociationAudit(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditList := []metadata.AuditLog{}
	for i, asstID := range asstIDList {
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, asstID, objID, &data.Data.Info[i])
		if err != nil {
			blog.Errorf("delete instance association failed, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		auditList = append(auditList, *auditLog)
	}
	// save audit log.
	err = audit.SaveAuditLog(kit, auditList...)
	if err != nil {
		blog.Errorf("delete instance association failed, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return resp, nil
}

func (assoc *association) CheckAssociation(kit *rest.Kit, objectID string, instID int64) error {

	cond := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{common.BKObjIDField: objectID, common.BKInstIDField: instID},
			{common.BKAsstObjIDField: objectID, common.BKAsstInstIDField: instID},
		},
	}

	asst, err := assoc.SearchInstAssociation(kit,
		&metadata.SearchAssociationInstRequest{Condition: cond, ObjID: objectID})
	if nil != err {
		return err
	}
	if len(asst.Data) == 0 {
		return nil
	}
	for _, asst := range asst.Data {
		var errCheck error
		isInstExist := false
		if asst.ObjectID == objectID && asst.InstID == instID {
			isInstExist, errCheck = assoc.CheckAssociationInstExist(kit, asst.AsstObjectID, asst.AsstInstID)
		} else if asst.AsstObjectID == objectID && asst.AsstInstID == instID {
			isInstExist, errCheck = assoc.CheckAssociationInstExist(kit, asst.ObjectID, asst.InstID)
		} else {
			return kit.CCError.New(common.CCErrCommDBSelectFailed, "instance is not associated in selected association")
		}
		if errCheck != nil {
			return errCheck
		}
		if isInstExist {
			return kit.CCError.CCErrorf(common.CCErrTopoInstHasBeenAssociation, instID)
		}
	}

	return nil
}

// CheckAssociations returns error if the instances has associations with exist instances, clear dirty associations
func (assoc *association) CheckAssociations(kit *rest.Kit, objectID string, instIDs []int64) error {
	if len(instIDs) == 0 {
		return nil
	}

	// get all associations for the instances
	cond := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{common.BKObjIDField: objectID, common.BKInstIDField: mapstr.MapStr{common.BKDBIN: instIDs}},
			{common.BKAsstObjIDField: objectID, common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: instIDs}},
		},
	}

	associations, err := assoc.SearchInstAssociation(kit,
		&metadata.SearchAssociationInstRequest{Condition: cond})
	if err != nil {
		blog.ErrorJSON("search instance associations failed, err: %s, cond: %s, rid: %s", err, cond, kit.Rid)
		return err
	}

	if len(associations.Data) == 0 {
		return nil
	}

	instIDExistsMap := make(map[int64]bool)
	for _, instID := range instIDs {
		instIDExistsMap[instID] = true
	}

	// get all associated inst IDs grouped by object ID, then check if any inst exists, clear not exist one's assts
	asstObjInstIDsMap := make(map[string][]int64)
	for _, asst := range associations.Data {
		if asst.ObjectID == objectID && instIDExistsMap[asst.InstID] {
			asstObjInstIDsMap[asst.AsstObjectID] = append(asstObjInstIDsMap[asst.AsstObjectID], asst.AsstInstID)
		} else if asst.AsstObjectID == objectID && instIDExistsMap[asst.AsstInstID] {
			asstObjInstIDsMap[asst.ObjectID] = append(asstObjInstIDsMap[asst.ObjectID], asst.InstID)
		}
	}

	for asstObjID, asstInstIDs := range asstObjInstIDsMap {
		query := &metadata.QueryCondition{
			Condition: mapstr.MapStr{
				common.GetInstIDField(asstObjID): mapstr.MapStr{common.BKDBIN: asstInstIDs},
			},
			Page: metadata.BasePage{Limit: 1},
		}

		asstInstRsp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, asstObjID, query)
		if err != nil {
			blog.ErrorJSON("check instance existence failed, err: %s, query: %s, rid: %s", err, query, kit.Rid)
			return kit.CCError.Error(common.CCErrObjectSelectInstFailed)
		}
		if err := asstInstRsp.CCError(); err != nil {
			blog.ErrorJSON("check instance existence failed, err: %s, query: %s, rid: %s", err, query, kit.Rid)
			return err
		}

		if len(asstInstRsp.Data.Info) > 0 {
			return kit.CCError.CCError(common.CCErrorInstHasAsst)
		}

		deleteAsstCond := mapstr.MapStr{
			common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: asstObjID, common.BKInstIDField: mapstr.MapStr{common.BKDBIN: asstInstIDs}},
				{common.BKAsstObjIDField: asstObjID, common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: asstInstIDs}},
			},
		}

		if err := assoc.deleteAssociationDirtyData(kit, asstObjID, deleteAsstCond); err != nil {
			blog.ErrorJSON("delete dirty assts failed, err: %s, cond: %s, rid: %s", err, deleteAsstCond, kit.Rid)
			return err
		}
	}
	return nil
}

func (assoc *association) CheckAssociationInstExist(kit *rest.Kit, objectID string, instID int64) (bool, error) {
	instIDField := common.GetInstIDField(objectID)
	instRsp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objectID,
		&metadata.QueryCondition{Condition: mapstr.MapStr{instIDField: instID}})
	if err != nil {
		return false, kit.CCError.Error(common.CCErrObjectSelectInstFailed)
	}
	if !instRsp.Result {
		return false, kit.CCError.New(instRsp.Code, instRsp.ErrMsg)
	}
	if len(instRsp.Data.Info) > 0 {
		return true, nil
	}
	// 实例不存在，删除实例的关联关系
	cond := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{common.BKObjIDField: objectID, common.BKInstIDField: instID},
			{common.BKAsstObjIDField: objectID, common.BKAsstInstIDField: instID},
		},
	}
	if err := assoc.deleteAssociationDirtyData(kit, objectID, cond); err != nil {
		return false, err
	}
	return false, nil
}

func (assoc *association) deleteAssociationDirtyData(kit *rest.Kit, objID string, cond mapstr.MapStr) error {
	delOpt := &metadata.InstAsstDeleteOption{
		Opt:   metadata.DeleteOption{Condition: cond},
		ObjID: objID,
	}

	rsp, err := assoc.clientSet.CoreService().Association().DeleteInstAssociation(kit.Ctx, kit.Header, delOpt)
	if nil != err {
		blog.Errorf("request to delete instance association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("failed to delete the inst association info , err: %s, rid: %s", rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}
