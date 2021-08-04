package operation

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/ac/meta"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	// ImportInstAssociation add instance association by excel
	ImportInstAssociation(kit *rest.Kit, objID string, importData map[int]metadata.ExcelAssociation,
		asstObjectUniqueIDMap map[string]int64, objectUniqueID int64) (metadata.ResponeImportAssociationData, error)

	// FindAssociationByObjectAssociationID find association by objid and asstid
	FindAssociationByObjectAssociationID(ctx context.Context, kit *rest.Kit, objID string,
		asstIDArr []string) ([]metadata.Association, errors.CCError)
}

// NewAssociationOperation create a new association operation instance
func NewAssociationOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager, lang language.DefaultCCLanguageIf) AssociationOperationInterface {
	return &association{
		clientSet:   client,
		authManager: authManager,
		lang:        lang,
	}
}

type association struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	lang        language.DefaultCCLanguageIf
}

// ImportInstAssociation add instance association by excel
func (assoc *association) ImportInstAssociation(kit *rest.Kit, objID string,
	importData map[int]metadata.ExcelAssociation, asstObjectUniqueIDMap map[string]int64,
	objectUniqueID int64) (metadata.ResponeImportAssociationData, error) {

	ia := NewImportAssociation(assoc, kit, objID, importData, asstObjectUniqueIDMap, objectUniqueID)
	err := ia.ParsePrimaryKey()
	resp := metadata.ResponeImportAssociationData{}
	if err != nil {
		return resp, err
	}

	errIdxMsgMap := ia.ImportAssociation()
	if len(errIdxMsgMap) > 0 {
		err = kit.CCError.CCError(common.CCErrorTopoImportAssociation)
	}

	for row, msg := range errIdxMsgMap {
		resp.ErrMsgMap = append(resp.ErrMsgMap, metadata.RowMsgData{
			Row: row,
			Msg: msg,
		})
	}

	return resp, err
}

// FindAssociationByObjectAssociationID find association by objid and asstid
func (assoc *association) FindAssociationByObjectAssociationID(ctx context.Context, kit *rest.Kit, objID string,
	asstIDArr []string) ([]metadata.Association, errors.CCError) {

	input := &metadata.QueryCondition{}
	input.Condition = map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{common.BKObjIDField: objID},
			{common.BKAsstObjIDField: objID},
		},
		common.AssociationObjAsstIDField: map[string]interface{}{common.BKDBIN: asstIDArr},
	}
	input.Page.Limit = common.BKNoLimit

	resp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("find object by association http do error. err: %v, input: %v, rid: %s",
			err, input, kit.Rid)
		return nil, err
	}

	if err := resp.CCError(); err != nil {
		blog.Errorf("find object by association http reply error. reply: %v, input: %v, rid: %s",
			resp, input, kit.Rid)
		return nil, err
	}

	return resp.Data.Info, nil
}

type importAssociationInst struct {
	instID int64
	//strings.Joion([]string{property name, property value}, "=")|true
	attrNameVal map[string]bool
}

type importAssociation struct {
	objID      string
	cli        *association
	ctx        context.Context
	importData map[int]metadata.ExcelAssociation
	// 模型使用的唯一校验相关的信息
	asstObjectUniqueIDMap map[string]int64
	objectUniqueID        int64
	kit                   *rest.Kit
	language              language.DefaultCCLanguageIf

	// map[AssociationName]Association alias  map[association flag]Association
	asstIDInfoMap map[string]*metadata.Association
	// asst obj info  map[objID]map[property name] attribute
	asstObjIDProperty map[string]map[string]metadata.Attribute
	// 当前操作模型使用的唯一校验，用来解决自关联使用不同的唯一校验
	objIDProperty map[string]metadata.Attribute

	parseImportDataErr map[int]string
	//map[objID][]condition.Condition， 查询与当前操作模型有关联关系的实例参数
	queryAsstInstCondArr map[string][]mapstr.MapStr
	//[]condition.Condition, 查询当前操作模型的的实例参数
	queryInstCondArr []mapstr.MapStr

	// map[objID][instance id]strings.Join([]string{property name, property value}, "=")[]importAssociationInst
	asstInstIDAttrKeyValMap map[string]map[string][]*importAssociationInst
	// map[instance id]strings.Join([]string{property name, property value}, "=")[]importAssociationInst
	instIDAttrKeyValMap map[string][]*importAssociationInst
	//http header http request id
	rid string

	authManager *extensions.AuthManager
}

type importAssociationInterface interface {
	// ParsePrimaryKey parse msg about importAssociation
	ParsePrimaryKey() error

	// ImportAssociation add association by excel import
	ImportAssociation() map[int]string
}

// NewImportAssociation build an import association object
func NewImportAssociation(cli *association, kit *rest.Kit, objID string, importData map[int]metadata.ExcelAssociation,
	asstObjectUniqueIDMap map[string]int64, objectUniqueID int64) importAssociationInterface {

	return &importAssociation{
		objID:                 objID,
		cli:                   cli,
		ctx:                   kit.Ctx,
		importData:            importData,
		asstObjectUniqueIDMap: asstObjectUniqueIDMap,
		objectUniqueID:        objectUniqueID,

		kit:      kit,
		language: cli.lang,

		asstIDInfoMap:           make(map[string]*metadata.Association, 0),
		asstObjIDProperty:       make(map[string]map[string]metadata.Attribute, 0),
		objIDProperty:           make(map[string]metadata.Attribute, 0),
		parseImportDataErr:      make(map[int]string),
		queryAsstInstCondArr:    make(map[string][]mapstr.MapStr),
		queryInstCondArr:        make([]mapstr.MapStr, 0),
		asstInstIDAttrKeyValMap: make(map[string]map[string][]*importAssociationInst),
		instIDAttrKeyValMap:     make(map[string][]*importAssociationInst),
		rid:                     kit.Rid,

		authManager: cli.authManager,
	}
}

// ImportAssociation add association by excel import
func (ia *importAssociation) ImportAssociation() map[int]string {
	ia.importAssociation()

	return ia.parseImportDataErr
}

// ParsePrimaryKey parse msg about importAssociation
func (ia *importAssociation) ParsePrimaryKey() error {
	err := ia.getAssociationInfo()
	if err != nil {
		blog.Errorf("parse primary key failed, err: %v, rid: %s", err, ia.rid)
		return err
	}

	err = ia.getObjProperty()
	if err != nil {
		blog.Errorf("parse primary key failed, err: %v, rid: %s", err, ia.rid)
		return err
	}

	err = ia.getAssociationObjProperty()
	if err != nil {
		blog.Errorf("parse primary key failed, err: %v, rid: %s", err, ia.rid)
		return err
	}

	ia.parseImportDataPrimary()
	err = ia.getInstDataByQueryCondArr()
	if err != nil {
		blog.Errorf("parse primary key failed, err: %v, rid: %s", err, ia.rid)
		return err
	}

	return nil

}

func (ia *importAssociation) importAssociation() {
	for idx, asstInfo := range ia.importData {
		_, ok := ia.parseImportDataErr[idx]
		if ok {
			continue
		}

		asst, ok := ia.asstIDInfoMap[asstInfo.ObjectAsstID]
		if !ok {
			ia.parseImportDataErr[idx] = ia.language.Languagef("import_association_id_not_found",
				asstInfo.ObjectAsstID)
			continue
		}

		srcInstID, dstInstID, err := ia.getTargetIndexSrcDstInstID(idx, asst, asstInfo)
		if err != nil {
			continue
		}

		err = ia.authManager.AuthorizeByInstanceID(ia.ctx, ia.kit.Header, meta.Update, ia.objID, srcInstID)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			continue
		}

		err = ia.authManager.AuthorizeByInstanceID(ia.ctx, ia.kit.Header, meta.Update, asst.AsstObjID, dstInstID)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			continue
		}

		if ok := ia.checkExcelAssociationOperate(idx, srcInstID, dstInstID, asst, asstInfo); !ok {
			continue
		}
	}
}

func (ia *importAssociation) getTargetIndexSrcDstInstID(idx int, asst *metadata.Association,
	asstInfo metadata.ExcelAssociation) (int64, int64, error) {
	srcInstID, dstInstID, err := int64(0), int64(0), error(nil)
	if asst.ObjectID == ia.objID {
		srcInstID, err = ia.getObjectInstIDByPrimaryKey(asstInfo.SrcPrimary)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			return 0, 0, err
		}

		dstInstID, err = ia.getAssociationObjectInstIDByPrimaryKey(asst.AsstObjID, asstInfo.DstPrimary)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			return 0, 0, err
		}
	} else {
		srcInstID, err = ia.getAssociationObjectInstIDByPrimaryKey(asst.ObjectID, asstInfo.SrcPrimary)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			return 0, 0, err
		}

		dstInstID, err = ia.getObjectInstIDByPrimaryKey(asstInfo.DstPrimary)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			return 0, 0, err
		}
	}

	return srcInstID, dstInstID, nil
}

func (ia *importAssociation) checkExcelAssociationOperate(idx int, srcInstID, dstInstID int64,
	asst *metadata.Association, asstInfo metadata.ExcelAssociation) bool {

	switch asstInfo.Operate {
	case metadata.ExcelAssociationOperateAdd:

		conds := mapstr.MapStr{
			common.AssociationObjAsstIDField: asstInfo.ObjectAsstID,
			common.BKObjIDField:              asst.ObjectID,
			common.BKInstIDField:             srcInstID,
			common.AssociatedObjectIDField:   asst.AsstObjID,
		}
		isExist, err := ia.isExistInstAsst(idx, conds, dstInstID, asst.ObjectID, asst.Mapping)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			return false
		}

		if isExist {
			return false
		}

		ia.addSrcAssociation(idx, asst.AssociationName, srcInstID, dstInstID)
		return true

	case metadata.ExcelAssociationOperateDelete:
		conds := mapstr.MapStr{
			common.AssociationObjAsstIDField: asstInfo.ObjectAsstID,
			common.BKObjIDField:              asst.ObjectID,
			common.BKInstIDField:             srcInstID,
			common.AssociatedObjectIDField:   asst.AsstObjID,
			common.BKAsstInstIDField:         dstInstID,
		}
		ia.delSrcAssociation(idx, ia.objID, conds)
		return true
	default:
		ia.parseImportDataErr[idx] = ia.language.Language("import_association_operate_not_found")
		return true
	}
}

func (ia *importAssociation) getAssociationInfo() error {

	var associationFlag []string
	for _, info := range ia.importData {
		associationFlag = append(associationFlag, info.ObjectAsstID)
	}

	cond := mapstr.MapStr{
		common.AssociationObjAsstIDField: mapstr.MapStr{common.BKDBIN: associationFlag},
		common.BKDBOR: []mapstr.MapStr{
			{common.BKObjIDField: ia.objID},
			{common.BKAsstObjIDField: ia.objID},
		},
	}

	queryInput := &metadata.QueryCondition{Condition: cond}

	rsp, err := ia.cli.clientSet.CoreService().Association().ReadModelAssociation(ia.ctx, ia.kit.Header, queryInput)
	if err != nil {
		blog.Errorf("search object association failed, err: %v, input:%#v, rid:%s", err, queryInput, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search the inst association, err: %v, input:%#v, rid:%s", err, queryInput, ia.rid)
		return err
	}

	for index := range rsp.Data.Info {
		ia.asstIDInfoMap[rsp.Data.Info[index].AssociationName] = &rsp.Data.Info[index]
	}

	return nil
}

func (ia *importAssociation) getAssociationObjProperty() error {

	var objIDArr []string
	var uniqueIDArr []int64
	for objID, uniqueID := range ia.asstObjectUniqueIDMap {
		objIDArr = append(objIDArr, objID)
		uniqueIDArr = append(uniqueIDArr, uniqueID)
	}
	objIDArr = append(objIDArr, ia.objID)

	uniqueCond := mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: uniqueIDArr}}

	uniqueQueryCond := metadata.QueryCondition{Condition: uniqueCond}
	uniqueResult, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrUnique(ia.ctx, ia.kit.Header,
		uniqueQueryCond)
	if err != nil {
		blog.Errorf("search model unique attr failed, err: %v, input:%#v, rid:%s", err, uniqueQueryCond, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	var propertyIDArr []uint64
	for _, unique := range uniqueResult.Data.Info {
		for _, property := range unique.Keys {
			propertyIDArr = append(propertyIDArr, property.ID)
		}
	}

	cond := mapstr.MapStr{
		common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDArr},
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: propertyIDArr},
	}

	attrCond := &metadata.QueryCondition{Condition: cond}
	attrCond.Fields = []string{
		common.BKFieldID, common.BKObjIDField, common.BKPropertyIDField, common.BKPropertyNameField,
	}

	rsp, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrByCondition(ia.ctx, ia.kit.Header, attrCond)
	if err != nil {
		blog.Errorf("failed to search attribute , err: %c, input:%#v, rid:%s", err, attrCond, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search attribute, err: %v, input:%#v, rid:%s", err, cond, ia.rid)
		return err
	}

	for _, attr := range rsp.Data.Info {
		_, ok := ia.asstObjIDProperty[attr.ObjectID]
		if !ok {
			ia.asstObjIDProperty[attr.ObjectID] = make(map[string]metadata.Attribute)
		}
		ia.asstObjIDProperty[attr.ObjectID][attr.PropertyName] = attr
	}

	return nil

}

func (ia *importAssociation) getObjProperty() error {

	uniqueCond := mapstr.MapStr{
		common.BKFieldID: mapstr.MapStr{common.BKDBIN: ia.objectUniqueID},
	}

	uniqueQueryCond := metadata.QueryCondition{Condition: uniqueCond}
	uniqueResult, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrUnique(ia.ctx, ia.kit.Header,
		uniqueQueryCond)
	if err != nil {
		blog.Errorf("search model unique failed, err: %v, input:%#v, rid:%s", err, uniqueQueryCond, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = uniqueResult.CCError(); err != nil {
		blog.Errorf("search model unique, err: %v, input:%#v, rid:%s", err, uniqueQueryCond, ia.rid)
		return err
	}

	var propertyIDArr []uint64
	for _, unique := range uniqueResult.Data.Info {
		for _, property := range unique.Keys {
			propertyIDArr = append(propertyIDArr, property.ID)
		}
	}

	cond := mapstr.MapStr{
		common.BKObjIDField: mapstr.MapStr{common.BKDBIN: ia.objID},
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: propertyIDArr},
	}

	attrCond := &metadata.QueryCondition{Condition: cond}
	attrCond.Fields = []string{
		common.BKFieldID, common.BKObjIDField, common.BKPropertyIDField, common.BKPropertyNameField,
	}
	rsp, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrByCondition(ia.ctx, ia.kit.Header, attrCond)
	if err != nil {
		blog.Errorf("search attribute failed, err: %v, input:%#v, rid:%s", err, attrCond, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if ccErr := rsp.CCError(); ccErr != nil {
		blog.Errorf("search attribute failed, resp: %#v, input:%#v, err:%v, rid:%s", rsp, cond, ccErr, ia.rid)
		return ccErr
	}

	for _, attr := range rsp.Data.Info {
		ia.objIDProperty[attr.PropertyName] = attr
	}

	return nil

}

func (ia *importAssociation) parseImportDataPrimary() {

	for idx, info := range ia.importData {

		associationInst, ok := ia.asstIDInfoMap[info.ObjectAsstID]
		if !ok {
			ia.parseImportDataErr[idx] = ia.language.Languagef("import_asstid_not_found", info.ObjectAsstID)
			continue
		}

		var srcPropertyArr map[string]metadata.Attribute
		var dstPropertyArr map[string]metadata.Attribute

		isSelfObject := false
		if associationInst.ObjectID == ia.objID {
			srcPropertyArr = ia.objIDProperty
			dstPropertyArr = ia.asstObjIDProperty[associationInst.AsstObjID]
			isSelfObject = true

			if _, ok = ia.queryAsstInstCondArr[associationInst.AsstObjID]; !ok {
				ia.queryAsstInstCondArr[associationInst.AsstObjID] = make([]mapstr.MapStr, 0)
			}
		} else {
			srcPropertyArr = ia.asstObjIDProperty[associationInst.ObjectID]
			dstPropertyArr = ia.objIDProperty

			if _, ok = ia.queryAsstInstCondArr[associationInst.ObjectID]; !ok {
				ia.queryAsstInstCondArr[associationInst.ObjectID] = make([]mapstr.MapStr, 0)
			}
		}

		srcCond, err := ia.parseImportDataPrimaryItem(associationInst.ObjectID, info.SrcPrimary, srcPropertyArr)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
		} else {
			if isSelfObject {
				ia.queryInstCondArr = append(ia.queryInstCondArr, srcCond)
			} else {
				ia.queryAsstInstCondArr[associationInst.ObjectID] =
					append(ia.queryAsstInstCondArr[associationInst.ObjectID], srcCond)
			}
		}

		dstCond, err := ia.parseImportDataPrimaryItem(associationInst.AsstObjID, info.DstPrimary, dstPropertyArr)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
		} else {
			if isSelfObject {
				ia.queryAsstInstCondArr[associationInst.AsstObjID] =
					append(ia.queryAsstInstCondArr[associationInst.AsstObjID], dstCond)
			} else {
				ia.queryInstCondArr = append(ia.queryInstCondArr, dstCond)
			}
		}
	}

	return

}

func (ia *importAssociation) parseImportDataPrimaryItem(objID string, item string,
	propertyMap map[string]metadata.Attribute) (mapstr.MapStr, error) {

	keyValMap := mapstr.New()
	primaryArr := strings.Split(item, common.ExcelAsstPrimaryKeySplitChar)

	for _, primary := range primaryArr {

		primary = strings.TrimSpace(primary)
		keyValArr := strings.Split(primary, common.ExcelAsstPrimaryKeyJoinChar)
		if len(keyValArr) != 2 {
			blog.Errorf("parse import data primaryItem eror. primary:%s, rid:%s", primary, ia.rid)
			return nil, fmt.Errorf(ia.language.Languagef("import_asst_obj_property_str_primary_format_error",
				objID, item))
		}

		attr, ok := propertyMap[keyValArr[0]]
		if !ok {
			return nil, fmt.Errorf(ia.language.Languagef("import_asst_obj_primary_property_str_not_found",
				objID, keyValArr[0]))
		}

		realVal, err := convStrToCCType(keyValArr[1], attr)
		if err != nil {
			return nil, fmt.Errorf(ia.language.Languagef("import_asst_obj_property_str_primary_type_error",
				objID, keyValArr[0]))
		}

		keyValMap[attr.PropertyID] = realVal
	}
	if len(keyValMap) != len(propertyMap) {
		blog.Errorf("parse import inst failed. keyVal: %v, objID: %s, objIDProperty: %s, rid: %s",
			keyValMap, objID, propertyMap[objID], ia.rid)
		return nil, fmt.Errorf(ia.language.Languagef("import_asst_obj_property_str_primary_count_len", objID, item))
	}

	return keyValMap, nil
}

func (ia *importAssociation) getInstDataByQueryCondArr() error {

	for objID, valArr := range ia.queryAsstInstCondArr {
		instArr, err := ia.getObjectInstDataByCondArr(objID, valArr, ia.asstObjIDProperty[objID])
		if err != nil {
			blog.Errorf("get instance data failed, objID: %s, err: %v, rid: %s", objID, err, ia.rid)
			return err
		}

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		for _, inst := range instArr {
			ia.parseInstToImportAssociationObjectInst(objID, instIDKey, inst)
		}
	}

	instArr, err := ia.getObjectInstDataByCondArr(ia.objID, ia.queryInstCondArr, ia.objIDProperty)
	if err != nil {
		blog.Errorf("get instance data failed, objID: %s, err: %v, rid: %s", ia.objID, err, ia.rid)
		return err
	}

	instIDKey := metadata.GetInstIDFieldByObjID(ia.objID)
	for _, inst := range instArr {
		ia.parseInstToImportObjectInst(ia.objID, instIDKey, inst)
	}

	return nil
}

// 获取模型实例数据
func (ia *importAssociation) getObjectInstDataByCondArr(objID string, valArr []mapstr.MapStr,
	attrs map[string]metadata.Attribute) ([]mapstr.MapStr, error) {

	instIDKey := metadata.GetInstIDFieldByObjID(objID)
	if objID == common.BKInnerObjIDHost && len(valArr) > 0 {
		for idx, val := range valArr {
			if ok := val.Exists(common.BKCloudIDField); !ok {
				continue
			}
			intCloudID, err := val.Int64(common.BKCloudIDField)
			if err != nil {
				blog.Errorf("get cloudID failed, err: %v, rid: %s", err, ia.rid)
				return nil, err
			}
			valArr[idx][common.BKCloudIDField] = intCloudID
		}
	}

	if len(valArr) == 0 {
		return nil, nil
	}

	conds := mapstr.MapStr{common.BKDBOR: valArr}
	instArr, err := ia.getInstDataByObjIDCondArr(objID, instIDKey, conds, attrs)
	if err != nil {
		blog.Errorf("get instance data failed, objID: %s, err: %v, rid: %s", objID, err, ia.rid)
		return nil, err
	}

	return instArr, err
}

func (ia *importAssociation) getInstDataByObjIDCondArr(objID, instIDKey string, conds mapstr.MapStr,
	attrs map[string]metadata.Attribute) ([]mapstr.MapStr, error) {

	var fields []string
	for _, attr := range attrs { //ia.asstObjIDProperty[objID] {
		fields = append(fields, attr.PropertyID)
	}

	fields = append(fields, instIDKey)
	queryInput := &metadata.QueryCondition{}
	queryInput.Condition = conds
	queryInput.Fields = fields

	instSearchResult, err := ia.cli.clientSet.CoreService().Instance().ReadInstance(ia.ctx, ia.kit.Header,
		objID, queryInput)
	if err != nil {
		blog.Errorf("failed to  search %s instance , err: %v, input: %#v, rid:%s",
			objID, err, queryInput, ia.rid)
		return nil, ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err := instSearchResult.CCError(); err != nil {
		blog.Errorf("failed to search %s instance, reply: %#v, input:%#v, err: %v, rid:%s",
			objID, instSearchResult, queryInput, err, ia.rid)
		return nil, err
	}

	return instSearchResult.Data.Info, nil
}

// 导入模型关联对象实例数据
func (ia *importAssociation) parseInstToImportAssociationObjectInst(objID, instIDKey string, inst mapstr.MapStr) {

	_, ok := ia.asstInstIDAttrKeyValMap[objID]

	if !ok {
		ia.asstInstIDAttrKeyValMap[objID] = make(map[string][]*importAssociationInst)
	}

	attrs := ia.asstObjIDProperty[objID]

	instInfoArr, err := ia.parseInstToImportAssociationInstInfo(objID, instIDKey, inst, attrs)
	if err != nil {
		blog.Errorf("parse instance to import association object instance failed, rid: %s", ia.rid)
		// 沿用已有逻辑
		return
	}

	ia.asstInstIDAttrKeyValMap[objID] = mergeInstToImportAssociationInst(ia.asstInstIDAttrKeyValMap[objID], instInfoArr)

	return
}

// 导入模型数据实例查询， 自关联的时候是src 对象
func (ia *importAssociation) parseInstToImportObjectInst(objID, instIDKey string, inst mapstr.MapStr) {

	instInfoArr, err := ia.parseInstToImportAssociationInstInfo(objID, instIDKey, inst, ia.objIDProperty)
	if err != nil {
		blog.Errorf("parse instance to import object instance failed, rid: %s", ia.rid)
		// 沿用已有逻辑
		return
	}

	ia.instIDAttrKeyValMap = mergeInstToImportAssociationInst(ia.instIDAttrKeyValMap, instInfoArr)

	return
}

func (ia *importAssociation) parseInstToImportAssociationInstInfo(objID, instIDKey string, inst mapstr.MapStr,
	attrs map[string]metadata.Attribute) (map[string][]*importAssociationInst, error) {

	instID, err := inst.Int64(instIDKey)
	//inst info can not found
	if err != nil {
		blog.Warnf("get %s field from %s model error, err: %v, rid:%s", instID, objID, err, ia.rid)
		return nil, err
	}

	attrNameValMap := importAssociationInst{
		instID:      instID,
		attrNameVal: make(map[string]bool),
	}

	for _, attr := range attrs {
		val, err := inst.String(attr.PropertyID)
		//inst info can not found
		if err != nil {
			blog.Warnf("get %s field from %s model error, err: %v, rid: %s", attr.PropertyID, objID, err, ia.rid)
			return nil, err
		}
		attrNameValMap.attrNameVal[buildPrimaryStr(attr.PropertyName, val)] = true
	}

	instIDAttrKeyValMap := make(map[string][]*importAssociationInst, 0)
	for key := range attrNameValMap.attrNameVal {
		instIDAttrKeyValMap[key] = append(instIDAttrKeyValMap[key], &attrNameValMap)
	}

	return instIDAttrKeyValMap, nil
}

func (ia *importAssociation) delSrcAssociation(idx int, objID string, cond mapstr.MapStr) {

	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}

	delOpt := &metadata.InstAsstDeleteOption{
		Opt:   metadata.DeleteOption{Condition: cond},
		ObjID: objID,
	}

	result, err := ia.cli.clientSet.CoreService().Association().DeleteInstAssociation(ia.ctx, ia.kit.Header, delOpt)
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
		return
	}

	if err = result.CCError(); err != nil {
		ia.parseImportDataErr[idx] = result.ErrMsg
		return
	}

}

func (ia *importAssociation) addSrcAssociation(idx int, asstFlag string, instID, assInstID int64) {

	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}

	asstInfo := ia.asstIDInfoMap[asstFlag]

	inst := metadata.CreateOneInstanceAssociation{}
	inst.Data.ObjectAsstID = asstFlag
	inst.Data.InstID = instID
	inst.Data.ObjectID = asstInfo.ObjectID
	inst.Data.AsstObjectID = asstInfo.AsstObjID
	inst.Data.AsstInstID = assInstID
	inst.Data.AssociationKindID = asstInfo.AsstKindID
	rsp, err := ia.cli.clientSet.CoreService().Association().CreateInstAssociation(ia.ctx, ia.kit.Header, &inst)
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
	}

	if rsp == nil {
		return
	}

	if err = rsp.CCError(); err != nil {
		ia.parseImportDataErr[idx] = rsp.ErrMsg
	}
}

func (ia *importAssociation) isExistInstAsst(idx int, cond mapstr.MapStr, dstInstID int64, objID string,
	asstMapping metadata.AssociationMapping) (isExit bool, err error) {

	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}

	if asstMapping != metadata.OneToOneMapping {
		cond.Set(common.BKAsstInstIDField, dstInstID)
	}

	queryCond := &metadata.InstAsstQueryCondition{
		Cond:  metadata.QueryCondition{Condition: cond},
		ObjID: objID,
	}
	rsp, err := ia.cli.clientSet.CoreService().Association().ReadInstAssociation(ia.ctx, ia.kit.Header, queryCond)
	if err != nil {
		blog.Errorf("search instance association failed, err: %v, rid: %s", err, ia.rid)
		return false, err
	}

	if err = rsp.CCError(); err != nil {
		ia.parseImportDataErr[idx] = rsp.ErrMsg
		return false, err
	}

	if len(rsp.Data.Info) == 0 {
		return false, nil
	}

	if rsp.Data.Info[0].AsstInstID != dstInstID && asstMapping == metadata.OneToOneMapping {
		errMsg := ia.kit.CCError.Errorf(common.CCErrCommDuplicateItem, "association")
		blog.Errorf("check whether exist instance association failed, err: %s, rid: %s", errMsg, ia.rid)
		return false, errMsg
	}

	return true, nil
}

func (ia *importAssociation) getAssociationObjectInstIDByPrimaryKey(objID, primary string) (int64, error) {

	primaryArr := strings.Split(primary, common.ExcelAsstPrimaryKeySplitChar)
	if len(primaryArr) == 0 {
		errMsg := fmt.Errorf(ia.language.Languagef("import_instance_not_found", objID, primary))
		blog.Errorf("get association object(%s) instID failed, err: %s, rid: %s", objID, errMsg, ia.rid)
		return 0, errMsg
	}

	instArr, ok := ia.asstInstIDAttrKeyValMap[objID][primaryArr[0]]
	if !ok {
		errMsg := fmt.Errorf(ia.language.Languagef("import_instance_not_found", objID, primaryArr[0]))
		blog.Errorf("get association object(%s) instID failed, err: %s, rid: %s", objID, errMsg, ia.rid)
		return 0, errMsg
	}

	if instID := findInst(instArr, primaryArr); instID != 0 {
		return instID, nil
	}

	errMsg := fmt.Errorf(ia.language.Languagef("import_instance_not_found", objID, primary))
	blog.Errorf("get association object(%s) instID failed, err: %s, rid: %s", objID, errMsg, ia.rid)
	return 0, errMsg
}

func (ia *importAssociation) getObjectInstIDByPrimaryKey(primary string) (int64, error) {

	primaryArr := strings.Split(primary, common.ExcelAsstPrimaryKeySplitChar)
	if len(primaryArr) == 0 {
		errMsg := fmt.Errorf(ia.language.Languagef("import_instance_not_found", ia.objID, primary))
		blog.Errorf("get object instID failed, err: %s, rid: %s", errMsg, ia.rid)
		return 0, errMsg
	}

	instArr, ok := ia.instIDAttrKeyValMap[primaryArr[0]]
	if !ok {
		errMsg := fmt.Errorf(ia.language.Languagef("import_instance_not_found", ia.objID, primaryArr[0]))
		blog.Errorf("get object instID failed, err: %s, rid: %s", errMsg, ia.rid)
		return 0, fmt.Errorf(ia.language.Languagef("import_instance_not_found", ia.objID, primaryArr[0]))
	}

	if instID := findInst(instArr, primaryArr); instID != 0 {
		return instID, nil
	}

	errMsg := fmt.Errorf(ia.language.Languagef("import_instance_not_found", ia.objID, primary))
	blog.Errorf("get object instID failed, err: %s, rid: %s", errMsg, ia.rid)
	return 0, errMsg

}

func findInst(instArr []*importAssociationInst, primaryArr []string) int64 {

	for _, inst := range instArr {

		isEq := true
		for _, item := range primaryArr {
			if _, ok := inst.attrNameVal[item]; !ok {
				isEq = false
				break
			}
		}
		if isEq {
			return inst.instID
		}

	}

	return 0
}

func buildPrimaryStr(name, val string) string {

	return name + common.ExcelAsstPrimaryKeyJoinChar + val
}

func convStrToCCType(val string, attr metadata.Attribute) (interface{}, error) {

	switch attr.PropertyType {
	case common.FieldTypeBool:

		return strconv.ParseBool(val)
	case common.FieldTypeEnum:
		option, optionOk := attr.Option.([]interface{})
		if !optionOk {
			return nil, fmt.Errorf("not foud")
		}
		return getEnumIDByName(val, option), nil
	case common.FieldTypeInt:
		return util.GetInt64ByInterface(val)
	case common.FieldTypeFloat:
		return util.GetFloat64ByInterface(val)

	default:
		return val, nil
	}
}

// getEnumIDByName get enum name from option
func getEnumIDByName(name string, items []interface{}) string {

	id := name
	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {
			enumName, ok := mapVal["name"].(string)
			if true == ok {
				if enumName == name {
					id = mapVal["id"].(string)
				}
			}
		}
	}

	return id
}

func mergeInstToImportAssociationInst(src,
	dst map[string][]*importAssociationInst) map[string][]*importAssociationInst {

	if dst == nil {
		return src
	}

	for key, valArr := range src {
		dst[key] = append(dst[key], valArr...)
	}

	return dst
}
