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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
	"configcenter/src/thirdparty/hooks"
)

// AddHost TODO
func (lgc *Logics) AddHost(kit *rest.Kit, appID int64, moduleIDs []int64, ownerID string,
	hostInfos map[int64]map[string]interface{}, importType metadata.HostInputType) ([]int64, []string, []string,
	[]string, error) {
	if len(moduleIDs) == 0 {
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
		return nil, nil, nil, nil, err
	}
	var err error
	defaultModule, err := lgc.CoreAPI.CoreService().Process().GetBusinessDefaultSetModuleInfo(kit.Ctx, kit.Header,
		appID)
	if err != nil {
		blog.Errorf("AddHost failed, get biz default module info failed, appID:%d, err:%s, rid:%s", appID, err.Error(),
			kit.Rid)
		return nil, nil, nil, nil, err
	}
	isInternalModule := make([]bool, 0)
	for _, moduleID := range moduleIDs {
		isInternalModule = append(isInternalModule, defaultModule.IsInternalModule(moduleID))
	}
	isInternalModule = util.BoolArrayUnique(isInternalModule)
	if len(isInternalModule) > 1 {
		err := kit.CCError.CCError(common.CCErrHostTransferFinalModuleConflict)
		return nil, nil, nil, nil, err
	}
	toInternalModule := isInternalModule[0]

	hostIDs := make([]int64, 0)
	instance := NewImportInstance(kit, ownerID, lgc)

	hostIDMap, existsHostMap, err := instance.ExtractAlreadyExistHosts(kit.Ctx, hostInfos)
	if err != nil {
		blog.Errorf("get hosts failed, err:%s, rid:%s", err.Error(), kit.Rid)
		return nil, nil, nil, nil, err
	}

	var errMsg, updateErrMsg, successMsg []string

	// for audit log.
	logContents := make([]metadata.AuditLog, 0)
	audit := auditlog.NewHostAudit(lgc.CoreAPI.CoreService())
	ccLang := lgc.Engine.Language.CreateDefaultCCLanguageIf(httpheader.GetLanguage(kit.Header))

	for _, index := range util.SortedMapInt64Keys(hostInfos) {
		host := hostInfos[index]
		if host == nil {
			continue
		}

		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk == false || "" == innerIP {
			errMsg = append(errMsg, ccLang.Languagef("host_import_innerip_empty", index))
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

		iSubAreaVal, err := util.GetInt64ByInterface(iSubArea)
		if err != nil || iSubAreaVal < 0 {
			errMsg = append(errMsg, ccLang.Language("import_host_cloudID_invalid"))
			continue
		}
		host[common.BKCloudIDField] = iSubAreaVal

		var intHostID int64
		var existInDB bool

		// we support update host info both base on hostID and innerIP, hostID has higher priority then innerIP
		hostIDFromInput, bHostIDInInput := host[common.BKHostIDField]
		if bHostIDInInput == true {
			intHostID, err = util.GetInt64ByInterface(hostIDFromInput)
			if err != nil {
				errMsg = append(errMsg, ccLang.Language("import_host_hostID_not_int"))
				continue
			}
			existInDB = true
		} else {
			// try to get hostID from db
			key := generateHostCloudKey(innerIP, iSubAreaVal)
			intHostID, existInDB = hostIDMap[key]
		}

		// remove unchangeable fields
		delete(host, common.BKHostIDField)

		var auditLog []metadata.AuditLog
		if existInDB {
			// remove unchangeable fields
			delete(host, common.BKImportFrom)
			delete(host, common.CreateTimeField)

			var err error

			// generate audit log before really change it.
			generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit,
				metadata.AuditUpdate).WithUpdateFields(host)
			auditLog, err = audit.GenerateAuditLog(generateAuditParameter, appID,
				[]mapstr.MapStr{existsHostMap[intHostID]})
			if err != nil {
				blog.Errorf("generate host audit log failed before update host, hostID: %d, bizID: %d, err: %v, rid: %s",
					intHostID, innerIP, err, kit.Rid)
				errMsg = append(errMsg, err.Error())
				continue
			}

			// update host instance.
			if err := instance.updateHostInstance(index, host, intHostID); err != nil {
				updateErrMsg = append(updateErrMsg, err.Error())
				continue
			}
		} else {
			intHostID, err = instance.addHostInstance(iSubAreaVal, index, appID, moduleIDs, toInternalModule, host)
			if err != nil {
				errMsg = append(errMsg,
					fmt.Errorf(ccLang.Languagef("host_import_add_fail", index, innerIP, err.Error())).Error())
				continue
			}
			host[common.BKHostIDField] = intHostID
			hostIDMap[generateHostCloudKey(innerIP, iSubAreaVal)] = intHostID

			// to generate audit log.
			generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
			auditLog, err = audit.GenerateAuditLog(generateAuditParameter, appID, []mapstr.MapStr{host})
			if err != nil {
				blog.Errorf("generate host audit log failed after create host, hostID: %d, bizID: %d, err: %v, rid: %s",
					intHostID, appID, err, kit.Rid)
				errMsg = append(errMsg, err.Error())
				continue
			}
		}

		// add current host operate result to batch add result.
		successMsg = append(successMsg, strconv.FormatInt(index, 10))

		// add audit log.
		logContents = append(logContents, auditLog...)
		hostIDs = append(hostIDs, intHostID)
	}

	// to save audit log.
	if len(logContents) > 0 {
		if err := audit.SaveAuditLog(kit, logContents...); err != nil {
			return hostIDs, successMsg, updateErrMsg, errMsg, fmt.Errorf("save audit log failed, but add host success, err: %v",
				err)
		}
	}

	if 0 < len(errMsg) || 0 < len(updateErrMsg) {
		return hostIDs, successMsg, updateErrMsg, errMsg, errors.New(ccLang.Language("host_import_err"))
	}

	return hostIDs, successMsg, updateErrMsg, errMsg, nil
}

// getIpField get ipv4 and ipv6 address, ipv4 and ipv6 address cannot be null at the same time.
func getIpField(host map[string]interface{}) (string, string, string) {

	innerIP, v4Ok := host[common.BKHostInnerIPField].(string)
	innerIPv6, v6Ok := host[common.BKHostInnerIPv6Field].(string)
	if (!v4Ok || innerIP == "") && (!v6Ok || innerIPv6 == "") {
		return "host_import_innerip_v4_v6_empty", "", ""
	}
	return "", innerIP, innerIPv6
}

// UpdateHostByExcel update host by excel
// NOCC:golint/fnsize(后续重构，和实例合在一起)
func (lgc *Logics) UpdateHostByExcel(kit *rest.Kit, hosts map[int64]map[string]interface{}, hostIDArr []int64,
	indexHostIDMap map[int64]int64) ([]int64, []string, error) {

	relRes, err := lgc.getHostRelationDestMsg(kit)
	if err != nil {
		blog.Errorf("get object relation failed, err: %v, rid: %s", err, kit.Rid)
		return nil, nil, err
	}

	hostCond := map[string]interface{}{common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDArr}}
	hostInfoArr, err := lgc.GetHostInfoByConds(kit, hostCond)
	if err != nil {
		blog.Errorf("get hosts failed, err: %v, condition: %#v, rid: %s", err, hostCond, kit.Rid)
		return nil, nil, err
	}

	hostMap := make(map[int64]mapstr.MapStr)
	for _, host := range hostInfoArr {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.Errorf("parse host id failed, err: %v, host: %#v, rid: %s", err, host, kit.Rid)
			return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}
		hostMap[hostID] = host
	}

	hostRelations, err := lgc.GetHostRelations(kit, metadata.HostModuleRelationRequest{HostIDArr: hostIDArr,
		Fields: []string{common.BKAppIDField, common.BKHostIDField}})
	if err != nil {
		blog.Errorf("get host relations failed, err: %v, hostIDs: %+v, rid: %s", err, hostIDArr, kit.Rid)
		return nil, nil, err
	}

	hostBizMap := make(map[int64]int64)
	for _, relation := range hostRelations {
		hostBizMap[relation.HostID] = relation.AppID
	}

	successMsg := make([]int64, 0)
	errMsg := make([]string, 0)
	audit := auditlog.NewHostAudit(lgc.CoreAPI.CoreService())
	ccLang := lgc.Engine.Language.CreateDefaultCCLanguageIf(httpheader.GetLanguage(kit.Header))
	for _, index := range util.SortedMapInt64Keys(hosts) {
		host := hosts[index]
		delete(host, common.BKHostIDField)
		intHostID := indexHostIDMap[index]

		genAuditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(host)
		auditLog, err := audit.GenerateAuditLog(genAuditParam, hostBizMap[intHostID],
			[]mapstr.MapStr{hostMap[intHostID]})
		if err != nil {
			blog.Errorf("generate host audit log failed, hostID: %d, err: %v, rid: %s", intHostID, err, kit.Rid)
			errMsg = append(errMsg, ccLang.Languagef("import_host_update_fail", index, err.Error()))
			continue
		}

		// use new transaction, need a new header
		kit.Header = kit.NewHeader()
		_ = lgc.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
			tableData, err := metadata.GetTableData(host, relRes)
			if err != nil {
				errMsg = append(errMsg, ccLang.Languagef("host_import_add_fail", index, err.Error()))
				return err
			}

			opt := &metadata.UpdateOption{
				Condition: mapstr.MapStr{common.BKHostIDField: intHostID},
				Data:      mapstr.NewFromMap(host),
			}
			_, err = lgc.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost,
				opt)
			if err != nil {
				blog.ErrorJSON("update host instance failed, err: %v, input: %+v, param: %+v, rid: %s", err, host, opt,
					kit.Rid)
				errMsg = append(errMsg, ccLang.Languagef("import_host_update_fail", index, err.Error()))
				return err
			}

			// update instance table field type data
			if tableData != nil {
				if err := lgc.updateTableData(kit, tableData, intHostID); err != nil {
					blog.ErrorJSON("update table data failed, data: %s, err: %s, rid: %s", host, err, kit.Rid)
					errMsg = append(errMsg, ccLang.Languagef("import_host_update_fail", index, err.Error()))
					return err
				}
			}

			if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
				blog.Errorf("success update host, but add host[%v] audit failed, err: %v, rid: %s", err, kit.Rid)
				errMsg = append(errMsg, ccLang.Languagef("import_host_update_fail", index, err.Error()))
				return err
			}

			successMsg = append(successMsg, index)
			return nil
		})
	}

	return successMsg, errMsg, nil
}

// AddHostByExcel add host by import excel
// NOCC:golint/fnsize(后续重构，和实例的合成一个函数)
func (lgc *Logics) AddHostByExcel(kit *rest.Kit, appID int64, moduleID int64, ownerID string,
	hostInfos map[int64]map[string]interface{}) (hostIDs []int64, successMsg []int64, errMsg []string, err error) {

	_, toInternalModule, err := lgc.GetModuleIDAndIsInternal(kit, appID, moduleID)
	if err != nil {
		blog.Errorf("AddHostByExcel failed, GetModuleIDAndIsInternal err:%s, appID:%d, moduleID:%d", err, appID,
			moduleID)
		return nil, nil, nil, err
	}

	instance := NewImportInstance(kit, ownerID, lgc)

	relRes, err := lgc.getHostRelationDestMsg(kit)
	if err != nil {
		blog.Errorf("get object relation failed, err: %v, rid: %s", err, kit.Rid)
		return nil, nil, nil, err
	}

	// for audit log
	audit := auditlog.NewHostAudit(lgc.CoreAPI.CoreService())
	ccLang := lgc.Engine.Language.CreateDefaultCCLanguageIf(httpheader.GetLanguage(kit.Header))

	for _, index := range util.SortedMapInt64Keys(hostInfos) {
		host := hostInfos[index]
		if host == nil {
			continue
		}

		errStr, innerIP, innerIPv6 := getIpField(host)
		if errStr != "" {
			errMsg = append(errMsg, ccLang.Languagef(errStr, index))
			continue
		}

		// the bk_cloud_id is directly connected area
		if _, exist := host[common.BKCloudIDField]; !exist {
			errMsg = append(errMsg, ccLang.Languagef("import_host_not_provide_cloudID", index))
			continue
		}

		cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if err != nil {
			errMsg = append(errMsg, ccLang.Languagef("import_host_cloudID_not_exist", index,
				innerIP+"/"+innerIPv6, util.GetStrByInterface(host[common.BKCloudIDField])))
			continue
		}

		// remove unchangeable fields
		delete(host, common.BKHostIDField)

		// use new transaction, need a new header
		kit.Header = kit.NewHeader()
		_ = lgc.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
			tableData, err := metadata.GetTableData(host, relRes)
			if err != nil {
				errMsg = append(errMsg, ccLang.Languagef("host_import_add_fail", index, innerIP+"/"+innerIPv6,
					err.Error()))
				return err
			}

			intHostID, err := instance.addHostInstance(cloudID, index, appID, []int64{moduleID}, toInternalModule, host)
			if err != nil {
				blog.Errorf("add host instance failed, err: %v, index: %d, bizID: %d, moduleID: %d, "+
					"toInternalModule: %t, host: %v, rid: %s", err, index, appID, moduleID, toInternalModule, host,
					kit.Rid)
				errMsg = append(errMsg, ccLang.Languagef("host_import_add_fail", index, innerIP+"/"+innerIPv6,
					err.Error()))
				return err
			}
			host[common.BKHostIDField] = intHostID

			// add host table field type data
			if tableData != nil {
				if err := lgc.addTableData(kit, tableData, intHostID); err != nil {
					blog.ErrorJSON("add table data failed, data: %s, err: %s, rid: %s", host, err, kit.Rid)
					errMsg = append(errMsg, ccLang.Languagef("host_import_add_fail", index, innerIP+"/"+innerIPv6,
						err.Error()))
					return err
				}
			}

			// to generate audit log.
			generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
			auditLog, err := audit.GenerateAuditLog(generateAuditParameter, appID, []mapstr.MapStr{host})
			if err != nil {
				blog.Errorf("generate host audit log failed after create host, hostID: %d, bizID: %d, err: %v, rid: %s",
					intHostID, appID, err, kit.Rid)
				errMsg = append(errMsg, err.Error())
				return err
			}

			// add current host operate result to batch add result
			successMsg = append(successMsg, index)

			// add audit log
			if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
				blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
				errMsg = append(errMsg, kit.CCError.Error(common.CCErrAuditSaveLogFailed).Error())
				return err
			}
			hostIDs = append(hostIDs, intHostID)
			return nil
		})
	}

	return hostIDs, successMsg, errMsg, nil
}

// getHostRelationDestMsg get host relation, it can only get bk_property_id and dest_model field
func (lgc *Logics) getHostRelationDestMsg(kit *rest.Kit) ([]metadata.ModelQuoteRelation, error) {
	opt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: &filter.Expression{
				RuleFactory: &filter.CombinedRule{
					Condition: filter.And,
					Rules: []filter.RuleFactory{
						&filter.AtomRule{
							Field:    common.BKSrcModelField,
							Operator: filter.OpFactory(filter.Equal),
							Value:    common.BKInnerObjIDHost,
						},
					},
				},
			},
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKPropertyIDField, common.BKDestModelField},
	}

	relRes, err := lgc.CoreAPI.CoreService().ModelQuote().ListModelQuoteRelation(kit.Ctx, kit.Header, opt)
	if err != nil {
		return nil, err
	}
	return relRes.Info, nil
}

func (lgc *Logics) addTableData(kit *rest.Kit, tableData *metadata.TableData, id int64) error {
	audit := auditlog.NewQuotedInstAuditLog(lgc.CoreAPI.CoreService())
	genAuditParams := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	var auditLogs []metadata.AuditLog
	for destModel, table := range tableData.ModelData {
		for idx := range table {
			table[idx].Set(common.BKInstIDField, id)
		}

		ids, err := lgc.CoreAPI.CoreService().ModelQuote().BatchCreateQuotedInstance(kit.Ctx, kit.Header, destModel,
			table)
		if err != nil {
			return err
		}

		// generate and save audit logs
		for i := range table {
			table[i][common.BKFieldID] = ids[i]
		}

		auditLog, ccErr := audit.GenerateAuditLog(genAuditParams, destModel, tableData.SrcModel,
			tableData.ModelPropertyRel[destModel], table)
		if ccErr != nil {
			return ccErr
		}
		auditLogs = append(auditLogs, auditLog...)
	}

	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

// updateTableData update table data
func (lgc *Logics) updateTableData(kit *rest.Kit, tableData *metadata.TableData, id int64) error {
	audit := auditlog.NewQuotedInstAuditLog(lgc.CoreAPI.CoreService())
	genAuditParams := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	filterOpt := metadata.CommonFilterOption{
		Filter: filtertools.GenAtomFilter(common.BKInstIDField, filter.Equal, id),
	}
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: filterOpt,
		Page:               metadata.BasePage{Limit: common.BKMaxPageSize},
	}

	var auditLogs []metadata.AuditLog
	for destModel := range tableData.ModelData {
		listRes, err := lgc.CoreAPI.CoreService().ModelQuote().ListQuotedInstance(kit.Ctx, kit.Header, destModel,
			listOpt)
		if err != nil {
			return err
		}

		err = lgc.CoreAPI.CoreService().ModelQuote().BatchDeleteQuotedInstance(kit.Ctx, kit.Header, destModel,
			&filterOpt)
		if err != nil {
			return err
		}

		auditLog, ccErr := audit.GenerateAuditLog(genAuditParams, destModel, tableData.SrcModel,
			tableData.ModelPropertyRel[destModel], listRes.Info)
		if ccErr != nil {
			return err
		}

		auditLogs = append(auditLogs, auditLog...)
	}

	// save audit logs
	err := audit.SaveAuditLog(kit, auditLogs...)
	if err != nil {
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	if err := lgc.addTableData(kit, tableData, id); err != nil {
		return err
	}

	return nil
}

// AddHostToResourcePool TODO
func (lgc *Logics) AddHostToResourcePool(kit *rest.Kit, hostList metadata.AddHostToResourcePoolHostList) ([]int64,
	*metadata.AddHostToResourcePoolResult, error) {
	bizID, err := lgc.GetDefaultAppIDWithSupplier(kit)
	if err != nil {
		blog.ErrorJSON("add host, but get default biz id failed, err: %s, input: %s, rid: %s", err, hostList, kit.Rid)
		return nil, nil, err
	}

	var toInternalModule bool
	hostList.Directory, toInternalModule, err = lgc.GetModuleIDAndIsInternal(kit, bizID, hostList.Directory)
	if err != nil {
		return nil, nil, err
	}

	hostIDs := make([]int64, 0)
	res := new(metadata.AddHostToResourcePoolResult)
	instance := NewImportInstance(kit, kit.SupplierAccount, lgc)
	logContents := make([]metadata.AuditLog, 0)
	audit := auditlog.NewHostAudit(lgc.CoreAPI.CoreService())

	for index, host := range hostList.HostInfo {
		if nil == host {
			continue
		}

		innerIP, exist := host[common.BKHostInnerIPField].(string)
		if !exist || "" == innerIP {
			res.Error = append(res.Error, metadata.AddOneHostToResourcePoolResult{
				Index:    index,
				ErrorMsg: kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKHostInnerIPField).Error(),
			})
			continue
		}
		cloudID, exist := host[common.BKCloudIDField]
		if !exist || cloudID == nil {
			res.Error = append(res.Error, metadata.AddOneHostToResourcePoolResult{
				Index:    index,
				ErrorMsg: kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKCloudIDField).Error(),
			})
			continue
		}

		// TODO remove this when bk_cloud_id field is upgraded to int type
		cloudIDVal, err := util.GetInt64ByInterface(cloudID)
		if err != nil || cloudIDVal < 0 {
			res.Error = append(res.Error, metadata.AddOneHostToResourcePoolResult{
				Index:    index,
				ErrorMsg: kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKCloudIDField).Error(),
			})
			continue
		}
		host[common.BKCloudIDField] = cloudIDVal

		hostID, err := instance.addHostInstance(cloudIDVal, int64(index), bizID, []int64{hostList.Directory},
			toInternalModule,
			host)
		if err != nil {
			res.Error = append(res.Error, metadata.AddOneHostToResourcePoolResult{
				Index:    index,
				ErrorMsg: err.Error(),
			})
			continue
		}
		host[common.BKHostIDField] = hostID

		hostIDs = append(hostIDs, hostID)
		res.Success = append(res.Success, metadata.AddOneHostToResourcePoolResult{
			Index:  index,
			HostID: hostID,
		})

		// generate audit log for create host.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, bizID, []mapstr.MapStr{host})
		if err != nil {
			blog.Errorf("generate host audit log failed after create host, hostID: %d, bizID: %d, err: %v, rid: %s",
				hostID, bizID, err, kit.Rid)
			res.Error = append(res.Error, metadata.AddOneHostToResourcePoolResult{
				Index:    index,
				HostID:   hostID,
				ErrorMsg: err.Error(),
			})
			continue
		}

		logContents = append(logContents, auditLog...)
	}

	// save audit log.
	if len(logContents) > 0 {
		if err := audit.SaveAuditLog(kit, logContents...); err != nil {
			blog.Errorf("save host audit log failed after create host, err: %v, rid: %s", err, kit.Rid)
			return hostIDs, res, err
		}
	}

	if 0 < len(res.Error) {
		return hostIDs, res, kit.CCError.CCErrorf(common.CCErrHostCreateFail)
	}

	return hostIDs, res, nil
}

func (lgc *Logics) getHostFields(kit *rest.Kit) (map[string]*metadata.ObjAttDes, error) {
	opt := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).MapStr()

	input := &metadata.QueryCondition{
		Condition: opt,
	}
	result, err := lgc.CoreAPI.CoreService().Model().
		ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("getHostFields http do error, err:%s, input:%+v, rid:%s", err.Error(), input, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	attributesDesc := make([]metadata.ObjAttDes, 0)
	for _, att := range result.Info {
		attributesDesc = append(attributesDesc, metadata.ObjAttDes{Attribute: att})
	}

	fields := make(map[string]*metadata.ObjAttDes)
	for index, f := range attributesDesc {
		fields[f.PropertyID] = &attributesDesc[index]
	}
	return fields, nil
}

// generateHostCloudKey generate a cloudKey for host that is unique among clouds by appending the cloudID.
func generateHostCloudKey(ip, cloudID interface{}) string {
	return fmt.Sprintf("%v-%v", ip, cloudID)
}

type importInstance struct {
	*backbone.Engine
	pheader   http.Header
	inputType metadata.HostInputType
	ownerID   string
	// cloudID       int64
	// hostInfos     map[int64]map[string]interface{}
	defaultFields map[string]*metadata.ObjAttDes
	rowErr        map[int64]error
	ctx           context.Context
	ccErr         ccErr.DefaultCCErrorIf
	ccLang        language.DefaultCCLanguageIf
	rid           string
	lgc           *Logics
	kit           *rest.Kit
}

// NewImportInstance TODO
func NewImportInstance(kit *rest.Kit, ownerID string, lgc *Logics) *importInstance {
	lang := httpheader.GetLanguage(kit.Header)
	return &importInstance{
		pheader: kit.Header,
		Engine:  lgc.Engine,
		ownerID: ownerID,
		ctx:     kit.Ctx,
		ccErr:   kit.CCError,
		ccLang:  lgc.Engine.Language.CreateDefaultCCLanguageIf(lang),
		rid:     kit.Rid,
		lgc:     lgc,
		kit:     kit,
	}
}

func (h *importInstance) updateHostInstance(index int64, host map[string]interface{}, hostID int64) error {
	delete(host, "import_from")
	delete(host, common.CreateTimeField)

	// 更新主机数据
	input := &metadata.UpdateOption{}
	input.Condition = map[string]interface{}{common.BKHostIDField: hostID}
	input.Data = host
	_, err := h.CoreAPI.CoreService().Instance().UpdateInstance(h.ctx, h.pheader, common.BKInnerObjIDHost, input)
	if err != nil {
		ip, _ := host[common.BKHostInnerIPField].(string)
		blog.Errorf("updateHostInstance http do error,  err:%s,input:%+v,rid:%s", err.Error(), input, h.rid)
		return fmt.Errorf(h.ccLang.Languagef("host_import_update_fail", index, ip, err.Error()))
	}

	return nil
}

// addHostInstance  add host
// cloud id：host belong cloud area id
// index: index number
// app id : host belong app id
// module id: host belong module id
// host : host info
func (h *importInstance) addHostInstance(cloudID, index, appID int64, moduleIDs []int64, toInternalModule bool,
	host map[string]interface{}) (int64, error) {
	ip, _ := host[common.BKHostInnerIPField].(string)
	if cloudID < 0 {
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip,
			h.ccLang.Language("import_host_cloudID_invalid")))
	}

	// determine if the cloud area exists
	// default cloud area must be exist
	if cloudID != common.BKDefaultDirSubArea {
		isExist, err := h.lgc.IsPlatAllExist(h.kit, []int64{cloudID})
		if nil != err {
			return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, err.Error()))

		}
		if !isExist {
			return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip,
				h.ccErr.Errorf(common.CCErrTopoCloudNotFound).Error()))

		}
	}
	host[common.BKCloudIDField] = cloudID

	input := &metadata.CreateModelInstance{
		Data: host,
	}

	// (h.ctx, h.pheader, host)
	var err error
	result, err := h.CoreAPI.CoreService().Instance().CreateInstance(h.ctx, h.pheader, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("addHostInstance http do error,err:%s, input:%+v,rid:%s", err.Error(), host, h.rid)
		return 0, err
	}

	hostID := int64(result.Created.ID)
	var hResult []metadata.ExceptionResult
	var option interface{}
	if toInternalModule == true {
		if len(moduleIDs) == 0 {
			err := h.ccErr.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
			return 0, err
		}
		opt := &metadata.TransferHostToInnerModule{
			ApplicationID: appID,
			ModuleID:      moduleIDs[0],
			HostID:        []int64{hostID},
		}
		option = opt
		hResult, err = h.CoreAPI.CoreService().Host().TransferToInnerModule(h.ctx, h.pheader, opt)
	} else {
		opt := &metadata.HostsModuleRelation{
			ApplicationID: appID,
			ModuleID:      moduleIDs,
			HostID:        []int64{hostID},
		}
		option = opt
		hResult, err = h.CoreAPI.CoreService().Host().TransferToNormalModule(h.ctx, h.pheader, opt)

	}
	if err != nil {
		blog.Errorf("transfer host failed, err: %v, result: %#v, input: %#v, rid: %s", err, hResult, option, h.rid)
		return 0, err
	}

	return hostID, nil
}

// ExtractAlreadyExistHosts extract hosts that is already in db(same innerIP+cloudID host and updated hosts with id)
// return: map[hostKey]hostID and exists host id to host info map
func (h *importInstance) ExtractAlreadyExistHosts(ctx context.Context, hostInfos map[int64]map[string]interface{}) (
	map[string]int64, map[int64]mapstr.MapStr, error) {

	// step1. extract all innerIP from hostInfos
	var ipArr []string
	hostIDs := make([]int64, 0)
	for _, host := range hostInfos {
		hostID, exists := host[common.BKHostIDField]
		if exists {
			intHostID, err := util.GetInt64ByInterface(hostID)
			if err != nil {
				blog.Errorf("parse hostID failed, err: %v, hostInfo: %#v, rid: %s", err, host, h.rid)
				return nil, nil, err
			}
			hostIDs = append(hostIDs, intHostID)
		}
		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk && "" != innerIP {
			ipArr = append(ipArr, innerIP)
		}
	}
	if len(ipArr) == 0 {
		return make(map[string]int64), make(map[int64]mapstr.MapStr), nil
	}

	// step2. query host info by innerIPs
	ipCond := make([]map[string]interface{}, len(ipArr))
	for index, innerIP := range ipArr {
		innerIPArr := strings.Split(innerIP, ",")
		ipCond[index] = map[string]interface{}{
			common.BKHostInnerIPField: map[string]interface{}{
				common.BKDBIN: innerIPArr,
			},
		}
	}
	if len(hostIDs) > 0 {
		ipCond = append(ipCond, mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDs}})
	}

	filter := map[string]interface{}{
		common.BKDBOR: append(ipCond),
	}
	query := &metadata.QueryCondition{
		Condition: filter,
		Page: metadata.BasePage{
			Start: 0,
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKHostInnerIPField, common.BKCloudIDField, common.BKHostIDField},
	}
	hResult, err := h.CoreAPI.CoreService().Instance().ReadInstance(ctx, h.pheader, common.BKInnerObjIDHost, query)
	if err != nil {
		blog.Errorf("get host failed, err: %v, input: %#v, rid:%s", err, query, h.rid)
		return nil, nil, err
	}

	// step3. arrange data as a map, cloudKey: hostID
	hostMap := make(map[string]int64, 0)
	hostIDMap := make(map[int64]mapstr.MapStr, 0)
	for _, host := range hResult.Info {
		key := generateHostCloudKey(host[common.BKHostInnerIPField], host[common.BKCloudIDField])
		hostID, err := host.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("get hostID failed, err: %v, hostInfo: %#v, rid: %s", err, host, h.rid)
			// message format: `convert %s  field %s to %s error %s`
			return hostMap, hostIDMap, h.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost,
				common.BKHostIDField, "int", err.Error())
		}
		hostMap[key] = hostID
		hostIDMap[hostID] = host
	}

	return hostMap, hostIDMap, nil
}

// AddHosts add host to business module
func (lgc *Logics) AddHosts(kit *rest.Kit, appID int64, moduleID int64, hostInfos []mapstr.MapStr) ([]int64, error) {
	if moduleID == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}
	if appID == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// check host attribute
	if err := lgc.checkHostAttr(kit, hostInfos); err != nil {
		return nil, err
	}

	// create host instance
	input := &metadata.CreateManyModelInstance{
		Datas: hostInfos,
	}
	result, err := lgc.CoreAPI.CoreService().Instance().CreateManyInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost,
		input)
	if err != nil {
		blog.Errorf("create host instance failed, input: %v, err: %v, rid: %s", input, err, kit.Rid)
		return nil, err
	}

	if len(result.Repeated) > 0 {
		blog.Errorf("host data repeated, input: %v, result: %v, rid: %s", hostInfos, result, kit.Rid)
		errMsg := util.GetStrByInterface(result.Repeated[0].Data["err_msg"])
		return nil, ccErr.NewCCError(common.CCErrCommDuplicateItem, errMsg)
	}

	if len(result.Exceptions) > 0 {
		blog.Errorf("create host failed, input: %v, result: %v, rid: %s", hostInfos, result, kit.Rid)
		return nil, kit.CCError.CCErrorf(int(result.Exceptions[0].Code), result.Exceptions[0].Message)
	}

	hostIDs := make([]int64, 0)
	for index, item := range result.Created {
		hostInfos[index][common.BKHostIDField] = int64(item.ID)
		hostIDs = append(hostIDs, int64(item.ID))
	}

	// create host module relation
	opt := &metadata.TransferHostToInnerModule{
		ApplicationID: appID,
		ModuleID:      moduleID,
		HostID:        hostIDs,
	}
	_, err = lgc.CoreAPI.CoreService().Host().TransferToInnerModule(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("add host relation failed, input: %v, err: %v, rid: %s", opt, err, kit.Rid)
		return nil, err
	}

	// to generate audit log
	audit := auditlog.NewHostAudit(lgc.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, logErr := audit.GenerateAuditLog(generateAuditParameter, appID, hostInfos)
	if logErr != nil {
		blog.Errorf("generate host audit log failed after create host, input: %v, bizID: %d, err: %v, rid: %s",
			hostInfos, appID, err, kit.Rid)
		return nil, logErr
	}

	// to save audit log.
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("add host success, but save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, fmt.Errorf("add host success, but save audit log failed, err: %v", err)
	}

	return hostIDs, nil
}

func (lgc *Logics) checkHostAttr(kit *rest.Kit, hostInfos []mapstr.MapStr) error {
	cloudIDs := make([]int64, 0)
	for index, host := range hostInfos {
		innerIPv4, isIPv4Ok := host[common.BKHostInnerIPField].(string)
		innerIPv6, isIPv6Ok := host[common.BKHostInnerIPv6Field].(string)
		if (!isIPv4Ok || innerIPv4 == "") && (!isIPv6Ok || innerIPv6 == "") {
			return kit.CCError.CCErrorf(common.CCErrCommAtLeastSetOneVal, common.BKHostInnerIPField,
				common.BKHostInnerIPv6Field)
		}

		cloudID, ok := host[common.BKCloudIDField]
		if !ok {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
		}
		cloudIDVal, err := util.GetInt64ByInterface(cloudID)
		if err != nil || cloudIDVal < 0 {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
		}
		if err := hooks.ValidHostCloudIDHook(kit, cloudIDVal); err != nil {
			return err
		}
		hostInfos[index][common.BKCloudIDField] = cloudIDVal
		cloudIDs = append(cloudIDs, cloudIDVal)

		address, ok := host[common.BKAddressingField].(string)
		if !ok || (address != common.BKAddressingDynamic && address != common.BKAddressingStatic) {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAddressingField)
		}
	}

	// validate cloud ids
	cloudIDs = util.IntArrayUnique(cloudIDs)
	isExist, err := lgc.IsPlatAllExist(kit, cloudIDs)
	if err != nil {
		return err
	}
	if !isExist {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
	}

	return nil
}
