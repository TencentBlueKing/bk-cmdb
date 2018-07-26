/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"
	sceneUtil "configcenter/src/scene_server/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
	"configcenter/src/scene_server/validator"
)

func (lgc *Logics) AddHost(appID int64, moduleID []int64, ownerID string, pheader http.Header, hostInfos map[int64]map[string]interface{}, importType metadata.HostInputType) ([]string, []string, []string, error) {

	instance := NewImportInstance(ownerID, pheader, lgc.Engine)
	defLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))

	var err error
	instance.defaultFields, err = lgc.getHostFields(ownerID, pheader)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get host fields failed, err: %v", err)
	}

	cond := hutil.NewOperation().WithOwnerID(ownerID).WithObjID(common.BKInnerObjIDHost).Data()
	assResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAssociations(context.Background(), pheader, cond)
	if err != nil || (err == nil && !assResult.Result) {
		return nil, nil, nil, fmt.Errorf("search host assosications failed, err: %v, result err: %s", err, assResult.ErrMsg)
	}
	instance.asstDes = assResult.Data

	hostMap, err := lgc.getAddHostIDMap(pheader, hostInfos)
	if err != nil {
		blog.Errorf("get hosts failed, err:%s", err.Error())
		return nil, nil, nil, fmt.Errorf("get hosts failed, err: %v", err)
	}

	_, _, err = instance.getHostAsstHande(importType, hostInfos)
	if nil != err {
		blog.Errorf("get host assocate info  error, errror:%s", err.Error())
		return nil, nil, nil, err
	}

	var errMsg, updateErrMsg, succMsg []string
	logConents := make([]auditoplog.AuditLogExt, 0)
	auditHeaders, err := lgc.GetHostAttributes(ownerID, pheader)
	if err != nil {
		return nil, nil, nil, err
	}

	for index, host := range hostInfos {
		if nil == host {
			continue
		}

		if importType == common.InputTypeExcel {

			err = instance.assObjectInt.SetObjAsstPropertyVal(host)
			if nil != err {
				blog.Errorf("host assocate property error %v %v", index, err)
				updateErrMsg = append(updateErrMsg, defLang.Languagef("import_row_int_error_str", index, err.Error()))
				continue
			}
		}

		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk == false || "" == innerIP {
			errMsg = append(errMsg, defLang.Languagef("host_import_innerip_empty", strconv.FormatInt(index, 10)))
			continue
		}

		var iSubArea interface{}
		iSubArea, ok := host[common.BKCloudIDField]
		if false == ok {
			iSubArea = host[common.BKCloudIDField]
		}
		if nil == iSubArea {
			iSubArea = common.BKDefaultDirSubArea
		}

		var iHostID interface{}
		var isOK bool
		// host not db ,check params host info with host id
		iHostID, isOK = host[common.BKHostIDField]

		if false == isOK {
			key := fmt.Sprintf("%s-%v", innerIP, iSubArea)
			iHost, isDBOK := hostMap[key]
			if isDBOK {
				isOK = isDBOK
				iHostID, _ = iHost[common.BKHostIDField]
			}

		}

		var err error
		var intHostID int64
		preData := make(map[string]interface{}, 0)
		if isOK {
			intHostID, err = util.GetInt64ByInterface(iHostID)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("invalid host id: %v", iHostID)
			}
			// delete system fields
			delete(host, common.BKHostIDField)
			preData, _, _ = lgc.GetHostInstanceDetails(pheader, ownerID, strconv.FormatInt(intHostID, 10))
			// update host instance.
			if err := instance.updateHostInstance(index, host, intHostID); err != nil {
				updateErrMsg = append(updateErrMsg, err.Error())
				continue
			}

		} else {
			intHostID, err = instance.addHostInstance(int64(common.BKDefaultDirSubArea), index, appID, moduleID, host)
			if err != nil {
				errMsg = append(errMsg, err.Error())
				continue
			}
		}

		succMsg = append(succMsg, strconv.FormatInt(index, 10))
		curData, _, err := lgc.GetHostInstanceDetails(pheader, ownerID, strconv.FormatInt(intHostID, 10))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("generate audit log, but get host instance defail failed, err: %v", err)
		}

		logConents = append(logConents, auditoplog.AuditLogExt{
			ID: intHostID,
			Content: metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: auditHeaders,
			},
			ExtKey: innerIP,
		})
	}

	if len(logConents) > 0 {
		user := util.GetUser(pheader)
		log := map[string]interface{}{
			common.BKContentField: logConents,
			common.BKOpDescField:  "import host",
			common.BKOpTypeField:  auditoplog.AuditOpTypeAdd,
		}
		_, err := lgc.CoreAPI.AuditController().AddHostLogs(context.Background(), ownerID, strconv.FormatInt(appID, 10), user, pheader, log)
		if err != nil {
			return succMsg, updateErrMsg, errMsg, fmt.Errorf("generate audit log, but get host instance defail failed, err: %v", err)
		}
	}

	if 0 < len(errMsg) || 0 < len(updateErrMsg) {
		return succMsg, updateErrMsg, errMsg, errors.New(lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader)).Language("host_import_err"))
	}

	return succMsg, updateErrMsg, errMsg, nil
}

func (lgc *Logics) getHostFields(ownerID string, pheader http.Header) (map[string]*metadata.ObjAttDes, error) {
	page := metadata.BasePage{Start: 0, Limit: common.BKNoLimit}
	opt := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(ownerID).WithPage(page).Data()
	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, opt)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("search host attributes failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	attributesDesc := make([]metadata.ObjAttDes, 0)
	for _, att := range result.Data {
		attributesDesc = append(attributesDesc, metadata.ObjAttDes{Attribute: att})
	}
	for idx, a := range attributesDesc {
		if !util.IsAssocateProperty(a.PropertyType) {
			continue
		}

		cond := hutil.NewOperation().WithPropertyID(a.PropertyID).WithOwnerID(ownerID).WithObjID(a.ObjectID).Data()
		assResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAssociations(context.Background(), pheader, cond)
		if err != nil || (err == nil && !result.Result) {
			return nil, fmt.Errorf("search host obj associations failed, err: %v, result err: %s", err, result.ErrMsg)
		}

		if 0 < len(assResult.Data) {
			attributesDesc[idx].AssociationID = assResult.Data[0].AsstObjID // by the rules, only one id
			attributesDesc[idx].AsstForward = assResult.Data[0].AsstForward // by the rules, only one id
		}
	}

	fields := make(map[string]*metadata.ObjAttDes)
	for index, f := range attributesDesc {
		fields[f.PropertyID] = &attributesDesc[index]
	}
	return fields, nil
}

func (lgc *Logics) getAddHostIDMap(pheader http.Header, hostInfos map[int64]map[string]interface{}) (map[string]map[string]interface{}, error) {
	var ipArr []string
	for _, host := range hostInfos {
		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk && "" != innerIP {
			ipArr = append(ipArr, innerIP)
		}
	}

	if 0 == len(ipArr) {
		return nil, fmt.Errorf("not found host inner ip fields")
	}

	var conds map[string]interface{}
	if 0 < len(ipArr) {
		conds = map[string]interface{}{common.BKHostInnerIPField: common.KvMap{common.BKDBIN: ipArr}}

	}

	query := &metadata.QueryInput{
		Condition: conds,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKHostIDField,
	}
	hResult, err := lgc.CoreAPI.HostController().Host().GetHosts(context.Background(), pheader, query)
	if err != nil {
		return nil, errors.New(lgc.Language.Languagef("host_search_fail_with_errmsg", err.Error()))
	} else if err == nil && !hResult.Result {
		return nil, errors.New(lgc.Language.Languagef("host_search_fail_with_errmsg", hResult.ErrMsg))
	}

	hostMap := make(map[string]map[string]interface{})
	for _, h := range hResult.Data.Info {
		key := fmt.Sprintf("%v-%v", h[common.BKHostInnerIPField], h[common.BKCloudIDField])
		hostMap[key] = h
	}

	return hostMap, nil
}

func (lgc *Logics) getObjAsstObjectPrimaryKey(ownerID string, defaultFields map[string]*metadata.ObjAttDes, pheader http.Header) (map[string][]metadata.ObjAttDes, error) {
	asstPrimaryKey := make(map[string][]metadata.ObjAttDes)
	for _, f := range defaultFields {
		if util.IsAssocateProperty(f.PropertyType) {
			page := metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKPropertyIDField}
			query := hutil.NewOperation().WithOwnerID(ownerID).WithObjID(f.AssociationID).WithPage(page).Data()
			result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, query)
			if err != nil || (err == nil && !result.Result) {
				return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
			}

			attributesDesc := make([]metadata.ObjAttDes, 0)
			for _, att := range result.Data {
				attributesDesc = append(attributesDesc, metadata.ObjAttDes{Attribute: att})
			}
			for idx, a := range attributesDesc {
				if !util.IsAssocateProperty(a.PropertyType) {
					continue
				}

				cond := hutil.NewOperation().WithPropertyID(a.PropertyID).WithOwnerID(ownerID).WithObjID(a.ObjectID).Data()
				assResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAssociations(context.Background(), pheader, cond)
				if err != nil || (err == nil && !result.Result) {
					return nil, fmt.Errorf("search host obj associations failed, err: %v, result err: %s", err, result.ErrMsg)
				}

				if 0 < len(assResult.Data) {
					attributesDesc[idx].AssociationID = assResult.Data[0].AsstObjID // by the rules, only one id
					attributesDesc[idx].AsstForward = assResult.Data[0].AsstForward // by the rules, only one id
				}
			}

			primaryFields := make([]metadata.ObjAttDes, 0)
			for _, f := range attributesDesc {
				if f.IsOnly {
					primaryFields = append(primaryFields, f)
				}
			}
			asstPrimaryKey[f.AssociationID] = primaryFields
		}
	}

	return asstPrimaryKey, nil
}

func (lgc *Logics) getAsstObjectConds(defaultFields map[string]*metadata.ObjAttDes, asstPrimaryKey map[string][]metadata.ObjAttDes, hostInfos map[int]map[string]interface{}) (map[string][]interface{}, map[int]error) {
	errs := make(map[int]error, 0)
	asstMap := make(map[string][]interface{}) //map[AssociationID][]condition

	for rowIndex, info := range hostInfos {
		for key, val := range info {
			f, ok := defaultFields[key]
			if false == ok {
				continue
			}
			if util.IsAssocateProperty(f.PropertyType) {

				asstFields, ok := asstPrimaryKey[f.AssociationID]
				if false == ok {
					errs[rowIndex] = errors.New(lgc.Language.Languagef("import_asst_property_str_not_found", key))
					continue
				}

				strVal, ok := val.(string)
				if false == ok {
					errs[rowIndex] = errors.New(lgc.Language.Languagef("import_property_str_format_error", key))
					continue
				}

				if common.ExcelDelAsstObjectRelation == strings.TrimSpace(strVal) {
					continue
				}
				rows := strings.Split(strVal, common.ExcelAsstPrimaryKeyRowChar)

				asstConds := make([]interface{}, 0)
				for _, row := range rows {
					if "" == row {
						continue
					}
					primaryKeys := strings.Split(row, common.ExcelAsstPrimaryKeySplitChar)
					if len(primaryKeys) != len(asstFields) {
						errs[rowIndex] = errors.New(lgc.Language.Languagef("import_asst_property_str_primary_count_len", key))
						continue
					}
					conds := common.KvMap{}
					if false == util.IsInnerObject(f.AssociationID) {
						conds[common.BKObjIDField] = f.AssociationID
					}
					for i, val := range primaryKeys {

						asstf := asstFields[i]
						var err error
						conds[asstf.PropertyID], err = sceneUtil.ConvByPropertytype(&asstf, val)
						if nil != err {
							errs[rowIndex] = errors.New(lgc.Language.Languagef("import_asst_property_str_primary_count_len", key))
							continue
						}
					}
					asstConds = append(asstConds, conds)

				}

				_, ok = asstMap[f.AssociationID]
				if ok {
					asstMap[f.AssociationID] = append(asstMap[f.AssociationID], asstConds...)
				} else {
					asstMap[f.AssociationID] = asstConds
				}

			}

		}

	}

	return asstMap, errs
}

func (lgc *Logics) getAsstInstByAsstObjectConds(pheader http.Header, asstInstConds map[string][]interface{}, defaultFields map[string]*metadata.ObjAttDes,
	asstPrimaryKey map[string][]metadata.ObjAttDes) (map[string]map[string]int64, error) {

	instPrimayIDMap := make(map[string]map[string]int64)
	for objID, conds := range asstInstConds {
		isExist := false
		for _, f := range defaultFields {
			if f.AssociationID == objID {
				isExist = true
			}
		}

		if false == isExist {
			continue
		}
		searchObjID := objID

		if !util.IsInnerObject(objID) {
			searchObjID = common.BKINnerObjIDObject
		}

		query := &metadata.QueryInput{
			Condition: common.KvMap{common.BKDBOR: conds},
		}
		result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), searchObjID, pheader, query)
		if err != nil || (err == nil && !result.Result) {
			return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
		}

		if len(result.Data.Info) == 0 {
			continue
		}

		primaryKey, _ := asstPrimaryKey[objID]
		for _, item := range result.Data.Info {
			keys := []string{}
			for _, f := range primaryKey {
				key, ok := item[f.PropertyID]
				if false == ok {
					errMsg := lgc.Language.Languagef("import_str_asst_str_query_data_format_error", objID, f.PropertyID)
					return nil, errors.New(errMsg)
				}
				keys = append(keys, fmt.Sprintf("%v", key))
			}
			_, ok := instPrimayIDMap[objID]
			if false == ok {
				instPrimayIDMap[objID] = make(map[string]int64)
			}

			idField := common.GetInstIDField(objID)
			idVal, exist := item[idField]
			if !exist {
				return nil, fmt.Errorf("%s %s not found", objID, idField)
			}
			id, err := util.GetInt64ByInterface(idVal)
			if nil != err {
				return nil, fmt.Errorf("get object %s  inst id error, inst info:%v, err:%s ", objID, item, err.Error())
			}
			instPrimayIDMap[objID][strings.Join(keys, common.ExcelAsstPrimaryKeySplitChar)] = id

		}

	}
	return instPrimayIDMap, nil
}

type importInstance struct {
	*backbone.Engine
	pheader   http.Header
	inputType metadata.HostInputType
	ownerID   string
	// cloudID       int64
	// hostInfos     map[int64]map[string]interface{}
	assObjectInt  *asstObjectInst
	defaultFields map[string]*metadata.ObjAttDes
	asstDes       []metadata.Association
	rowErr        map[int64]error
}

func NewImportInstance(ownerID string, pheader http.Header, engine *backbone.Engine) *importInstance {
	return &importInstance{
		pheader: pheader,
		Engine:  engine,
		ownerID: ownerID,
	}
}

func (h *importInstance) getHostAsstHande(inputType metadata.HostInputType, hostInfos map[int64]map[string]interface{}) (*asstObjectInst, map[int64]error, error) {
	if common.InputTypeExcel == inputType {
		h.assObjectInt = NewAsstObjectInst(h.pheader, h.Engine, h.ownerID, h.defaultFields)
		err := h.assObjectInt.GetObjAsstObjectPrimaryKey()
		if nil != err {
			return nil, nil, fmt.Errorf("get host assocate object  property failure, error:%s", err.Error())
		}
		h.rowErr, err = h.assObjectInt.InitInstFromData(hostInfos)
		if nil != err {
			return nil, nil, fmt.Errorf("get host assocate object instance data failure, error:%s", err.Error())
		}

	}
	return h.assObjectInt, h.rowErr, nil
}

func (h *importInstance) parseHostInstanceAssocate(inputType metadata.HostInputType, index int64, host map[string]interface{}) {
	if common.InputTypeExcel == inputType {
		if _, ok := h.rowErr[index]; true == ok {
			return
		}
		err := h.assObjectInt.SetObjAsstPropertyVal(host)
		if nil != err {
			blog.Error("host assocate property error %d %s", index, err.Error())
			h.rowErr[index] = err
		}

	}
}

func (h *importInstance) updateHostInstance(index int64, host map[string]interface{}, hostID int64) error {
	delete(host, "import_from")
	delete(host, common.CreateTimeField)

	filterFields := []string{common.CreateTimeField}

	valid := validator.NewValidMapWithKeyFields(util.GetOwnerID(h.pheader), common.BKInnerObjIDHost, filterFields, h.pheader, h.Engine)
	err := valid.ValidMap(host, common.ValidUpdate, hostID)
	if nil != err {
		blog.Errorf("host valid error %v %v", index, err.Error())
		return fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("import_row_int_error_str", index, err.Error()))

	}

	err = h.UpdateInstAssociation(h.pheader, hostID, h.ownerID, common.BKInnerObjIDHost, host)
	if err != nil {
		blog.Errorf("update host asst attr error : %s", err.Error())
		return fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("import_row_int_error_str", index, err.Error()))
	}

	input := make(map[string]interface{}, 2) //更新主机数据
	input["condition"] = map[string]interface{}{common.BKHostIDField: hostID}
	input["data"] = host

	uResult, err := h.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDHost, h.pheader, input)
	if err != nil || (err == nil && !uResult.Result) {
		ip, _ := host[common.BKHostInnerIPField].(string)
		return fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_update_fail", index, ip, fmt.Sprintf("%v, %v", err, uResult.ErrMsg)))
	}
	return nil
}

func (h *importInstance) UpdateInstAssociation(pheader http.Header, instID int64, ownerID, objID string, input map[string]interface{}) error {
	opt := hutil.NewOperation().WithOwnerID(h.ownerID).WithObjID(objID).Data()
	result, err := h.CoreAPI.ObjectController().Meta().SelectObjectAssociations(context.Background(), h.pheader, opt)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("search host attribute failed, err: %v, result err: %s", err, result.ErrMsg)
		return fmt.Errorf("search host attribute failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	for _, asst := range result.Data {
		if _, ok := input[asst.ObjectAttID]; ok {
			err := h.deleteInstAssociation(instID, objID, asst.AsstObjID)
			if nil != err {
				blog.Errorf("failed to delete the old inst association, error info is %s", err.Error())
				return err
			}
		}
	}

	asstFieldVals := ExtractDataFromAssociationField(int64(instID), input, result.Data)
	if 0 < len(asstFieldVals) {
		for _, asstFieldVal := range asstFieldVals {
			oResult, err := h.CoreAPI.ObjectController().Instance().CreateObject(context.Background(), common.BKTableNameInstAsst, h.pheader, asstFieldVal)
			if err != nil || (err == nil && !oResult.Result) {
				blog.Errorf("create host attribute failed, err: %v, result err: %s", err, oResult.ErrMsg)
				return fmt.Errorf("create host attribute failed, err: %v, result err: %s", err, oResult.ErrMsg)
			}
		}
	}

	return nil

}

func (h *importInstance) deleteInstAssociation(instID int64, objID, asstObjID string) error {
	opt := hutil.NewOperation().WithInstID(instID).WithObjID(objID)
	if "" != asstObjID {
		opt.WithAssoObjID(asstObjID)
	}

	result, err := h.CoreAPI.ObjectController().Instance().DelObject(context.Background(), common.BKTableNameInstAsst, h.pheader, opt.Data())
	if err != nil || (err == nil && !result.Result) {
		return fmt.Errorf("delete object [%v] failed, err: %v, result err: %s", instID, err, result.ErrMsg)
	}

	return nil
}

func (h *importInstance) addHostInstance(cloudID, index, appID int64, moduleID []int64, host map[string]interface{}) (int64, error) {
	ip, _ := host[common.BKHostInnerIPField].(string)
	_, ok := host[common.BKCloudIDField]
	if false == ok {
		host[common.BKCloudIDField] = cloudID
	}
	filterFields := []string{common.CreateTimeField}
	valid := validator.NewValidMapWithKeyFields(util.GetOwnerID(h.pheader), common.BKInnerObjIDHost, filterFields, h.pheader, h.Engine)
	err := valid.ValidMap(host, common.ValidCreate, 0)

	if nil != err {
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("import_row_int_error_str", index, err.Error()))
	}
	host[common.CreateTimeField] = time.Now().UTC()
	result, err := h.CoreAPI.HostController().Host().AddHost(context.Background(), h.pheader, host)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("add host by ip:%s, err:%v, reply err:%s", ip, err, result.ErrMsg)
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, result.ErrMsg))
	}

	hID, ok := result.Data.(map[string]interface{})[common.BKHostIDField]
	if !ok {
		blog.Errorf("add host by ip:%s reply not found hostID, err:%v, reply :%v", ip, err, result.Data)
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip))
	}

	hostID, err := util.GetInt64ByInterface(hID)
	if err != nil {
		blog.Errorf("add host by ip:%s reply hostID not interger, err:%v, reply :%v", ip, err, result.Data)
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, err.Error()))
	}

	hostAsstData := ExtractDataFromAssociationField(hostID, host, h.asstDes)

	for _, item := range hostAsstData {
		cResult, err := h.CoreAPI.ObjectController().Instance().CreateObject(context.Background(), common.BKTableNameInstAsst, h.pheader, item)
		if err != nil {
			return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, err.Error()))
		} else if err == nil && !cResult.Result {
			return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, cResult.ErrMsg))
		}
	}

	opt := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		ModuleID:      moduleID,
		HostID:        hostID,
	}
	hResult, err := h.CoreAPI.HostController().Module().AddModuleHostConfig(context.Background(), h.pheader, opt)
	if err != nil {
		blog.Errorf("add host module by ip:%s  err:%s", ip, err.Error())
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, err.Error()))
	} else if err == nil && !hResult.Result {
		blog.Errorf("add host module by ip:%s  err:%s", ip, hResult.ErrMsg)
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, hResult.ErrMsg))
	}

	return hostID, nil
}
