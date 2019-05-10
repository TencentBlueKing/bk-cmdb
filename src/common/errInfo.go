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

package common

//CC error number defined in this file
//Errno name is composed of the following format CCErr[XXX]
const (

	// the system code

	// CCSystemBusy the system is busy
	CCSystemBusy = -1
	CCSuccess    = 0
	CCSuccessStr = "success"

	// common error code 1199XXX

	// CCErrCommJSONUnmarshalFailed JSON deserialization failed
	CCErrCommJSONUnmarshalFailed = 1199000

	// CCErrCommJSONMarshalFailed JSON serialization failed
	CCErrCommJSONMarshalFailed = 1199001

	// CCErrCommHTTPDoRequestFailed the HTTP Request failed
	CCErrCommHTTPDoRequestFailed = 1199002

	// CCErrCommHTTPInputInvalid the input parameter is invalid, and the parameter here refers to the URL or Query parameter
	CCErrCommHTTPInputInvalid = 1199003

	// CCErrCommHTTPReadBodyFailed unable to read HTTP request body data
	CCErrCommHTTPReadBodyFailed = 1199004

	// CCErrCommHTTPBodyEmpty  HTTP request body data is not set
	CCErrCommHTTPBodyEmpty = 1199005

	// CCErrCommParamsInvalid parameter validation in the body is not paased
	CCErrCommParamsInvalid = 1199006

	// CCErrCommParamsNeedString  the parameter must be of type string
	CCErrCommParamsNeedString = 1199007

	// CCErrCommParamsLostField the parameter not specified
	CCErrCommParamsLostField = 1199008

	// CCErrCommParamsNeedInt the parameter must be of type int
	CCErrCommParamsNeedInt = 1199009

	// CCErrCommParamsNeedSet the parameter unassigned
	CCErrCommParamsNeedSet = 1199010

	// CCErrCommParamsIsInvalid the parameter is invalid or nonexistent
	CCErrCommParamsIsInvalid = 1199011

	// CCErrCommParseDataFailed failed to read data from data field
	CCErrCommParseDataFailed = 1199013

	// CCErrCommDuplicateItem duplicate data
	CCErrCommDuplicateItem = 1199014

	// CCErrCommOverLimit data length exceeds limit
	CCErrCommOverLimit = 1199015

	// CCErrFieldRegValidFailed regular verification failed
	CCErrFieldRegValidFailed = 1199016

	// CCErrCommDBSelectFailed database query failed
	CCErrCommDBSelectFailed = 1199017

	// CCErrCommDBInsertFailed database cannot add data
	CCErrCommDBInsertFailed = 1199018

	//CCErrCommNotFound the goal does not exist
	CCErrCommNotFound = 1199019

	//CCErrCommDBUpdateFailed database cannot update data
	CCErrCommDBUpdateFailed = 1199020

	//CCErrCommDBDeleteFailed database cannot delete data
	CCErrCommDBDeleteFailed = 1199021

	//CCErrCommRelyOnServerAddressFailed dependent service did not start
	CCErrCommRelyOnServerAddressFailed = 1199022

	//CCErrCommExcelTemplateFailed unable to generate and download
	CCErrCommExcelTemplateFailed = 1199023

	// CCErrCommParamsNeedTimeZone the parameter must be time zone type
	CCErrCommParamsNeedTimeZone = 1199024

	// CCErrCommParamsNeedBool the parameter must be bool type
	CCErrCommParamsNeedBool = 1199025

	// CCErrCommConfMissItem  missing configuration item
	CCErrCommConfMissItem = 1199026

	// CCErrCommNotAuthItem failed to get authorization information
	CCErrCommNotAuthItem = 1199027

	// CCErrCommNotAuthItem field valide failed
	CCErrCommFieldNotValid = 1199028

	//CCErrCommReplyDataFormatError Return data format error
	CCErrCommReplyDataFormatError = 1199029

	//CCErrCommReplyDataFormatError Return data format error
	CCErrCommPostInputParseError = 1199030

	// CCErrCommResourceInitFailed %s init failed
	CCErrCommResourceInitFailed = 1199031

	// CCErrCommParams should be string
	CCErrCommParamsShouldBeString = 1199032

	// CCErrCommSearchPropertyFailed get object property fields error
	CCErrCommSearchPropertyFailed = 1199033

	// CCErrCommParamsShouldBeEnum set enum
	CCErrCommParamsShouldBeEnum = 1199034

	// CCErrCommXXExceedLimit  xx exceed limit number
	CCErrCommXXExceedLimit = 1199035

	CCErrProxyRequestFailed      = 1199036
	CCErrRewriteRequestUriFailed = 1199037

	// CCErrCommLogicDataNil   need data %s is null
	CCErrCommInstDataNil = 1199038
	// CCErrCommInstFieldNotFound  %s field does not exist in %s
	CCErrCommInstFieldNotFound = 1199039
	// CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
	CCErrCommInstFieldConvFail = 1199040
	// CCErrCommUtilFail  handle %s error %s
	CCErrCommUtilHandleFail = 1199041

	// CCErrCommParamsNeedFloat the parameter must be float type
	CCErrCommParamsNeedFloat = 1199042

	// CCErrCommInternalServerError %s Internal Server Error
	CCErrCommInternalServerError = 1199999

	// apiserver 1100XXX

	// toposerver 1101XXX
	// CCErrTopoInstCreateFailed unable to create the instance
	CCErrTopoInstCreateFailed = 1101000

	// CCErrTopoInstDeleteFailed unable to delete the instance
	CCErrTopoInstDeleteFailed = 1101001

	// CCErrTopoInstUpdateFailed unable to update the instance
	CCErrTopoInstUpdateFailed = 1101002

	// CCErrTopoInstSelectFailed unable to search the instance
	CCErrTopoInstSelectFailed = 1101003

	// CCErrTopoModuleCreateFailed unable to create a module
	CCErrTopoModuleCreateFailed = 1101004

	// CCErrTopoModuleDeleteFailed unable to delete a module
	CCErrTopoModuleDeleteFailed = 1101005

	// CCErrTopoModuleUpdateFailed unable to update a module
	CCErrTopoModuleUpdateFailed = 1101006

	// CCErrTopoModuleSelectFailed unable to select a module
	CCErrTopoModuleSelectFailed = 1101007

	// CCErrTopoSetCreateFailed unable to create a set
	CCErrTopoSetCreateFailed = 1101008

	// CCErrTopoSetDeleteFailed unable to delete a set
	CCErrTopoSetDeleteFailed = 1101009

	// CCErrTopoSetUpdateFailed unable to update a set
	CCErrTopoSetUpdateFailed = 1101010

	// CCErrTopoSetSelectFailed unable to select a set
	CCErrTopoSetSelectFailed = 1101011

	// CCErrTopoInstHasHostChild include hosts
	CCErrTopoInstHasHostChild = 1101012

	// CCErrTopoObjectCreateFailed unable to create a object
	CCErrTopoObjectCreateFailed = 1101013

	// CCErrTopoObjectDeleteFailed unable to delete a object
	CCErrTopoObjectDeleteFailed = 1101014

	// CCErrTopoObjectUpdateFailed unable to update a object
	CCErrTopoObjectUpdateFailed = 1101015

	// CCErrTopoObjectSelectFailed unable to select a object
	CCErrTopoObjectSelectFailed = 1101016

	// CCErrTopoObjectAttributeCreateFailed unable to create a object attribute
	CCErrTopoObjectAttributeCreateFailed = 1101017

	// CCErrTopoObjectAttributeDeleteFailed unable to delete a object attribute
	CCErrTopoObjectAttributeDeleteFailed = 1101018

	// CCErrTopoObjectAttributeUpdateFailed unable to update a object attribute
	CCErrTopoObjectAttributeUpdateFailed = 1101019

	// CCErrTopoObjectAttributeSelectFailed unable to select a object attribute
	CCErrTopoObjectAttributeSelectFailed = 1101020

	// CCErrTopoObjectClassificationCreateFailed unable to create a object classification
	CCErrTopoObjectClassificationCreateFailed = 1101021

	// CCErrTopoObjectClassificationDeleteFailed unbale to delete a object classification
	CCErrTopoObjectClassificationDeleteFailed = 1101022

	// CCErrTopoObjectClassificationUpdateFailed unable to update a object classification
	CCErrTopoObjectClassificationUpdateFailed = 1101023

	// CCErrTopoObjectClassificationSelectFailed unable to select a object classification
	CCErrTopoObjectClassificationSelectFailed = 1101024

	// CCErrTopoObjectGroupCreateFailed unable to create object group
	CCErrTopoObjectGroupCreateFailed = 1101025

	// CCErrTopoObjectGroupDeleteFailed unable to delete a object group
	CCErrTopoObjectGroupDeleteFailed = 1101026

	// CCErrTopoObjectGroupUpdateFailed unable to update a object group
	CCErrTopoObjectGroupUpdateFailed = 1101027

	// CCErrTopoObjectGroupSelectFailed unable to select a object group
	CCErrTopoObjectGroupSelectFailed = 1101028

	// CCErrTopoObjectClassificationHasObject the object classification can't be deleted under clssification
	CCErrTopoObjectClassificationHasObject = 1101029

	// CCErrTopoHasHostCheckFailed cannot detect if host information is included
	CCErrTopoHasHostCheckFailed = 1101030

	// CCErrTopoHasHost include host
	CCErrTopoHasHost = 1101030

	// CCErrTopoGetCloudErrStrFaild get cloud error
	CCErrTopoGetCloudErrStrFaild = 1101031
	// CCErrTopoCloudNotFound   cloud area not found
	CCErrTopoCloudNotFound = 1101032

	// CCErrTopoGetAppFaild search app err %s
	CCErrTopoGetAppFaild = 1101033
	// CCErrTopoGetModuleFailed search  module err %s
	CCErrTopoGetModuleFailed = 1101034
	// CCErrTopoBizTopoOverLevel the mainline topo level over limit
	CCErrTopoBizTopoLevelOverLimit = 1101035
	// CCErrTopoInstHasBeenAssociation the mainline topo level over limit
	CCErrTopoInstHasBeenAssociation = 1101036
	// it is forbidden to delete , that has some insts
	CCErrTopoObjectHasSomeInstsForbiddenToDelete = 1101037
	// the associations %s->%s already exist
	CCErrTopoAssociationAlreadyExist = 1101038
	// the source association object does not exist
	CCErrTopoAssociationSourceObjectNotExist = 1101039
	// the destination association object does not exist
	CCErrTopoAssociationDestinationObjectNotExist = 1101040
	// invalid object association id, should be int64
	CCErrTopoInvalidObjectAssociationID = 1101041
	// got multiple object association with one association id
	CCErrTopoGotMultipleAssociationInstance = 1101042
	// association with a object has multiple instance, can not be deleted.
	CCErrTopoAssociationHasAlreadyBeenInstantiated = 1101043
	// get association kind with id failed.
	CCErrTopoGetAssociationKindFailed = 1101044
	// create object association missing object kind id, src object id or destination object id.
	CCErrorTopoAssociationMissingParameters = 1101045
	// the given association id does not exist.
	CCErrorTopoObjectAssociationNotExist = 1101046
	// update object association, but update fields that can not be updated.
	CCErrorTopoObjectAssociationUpdateForbiddenFields = 1101047
	// mainline object association do not exist
	CCErrorTopoMainlineObjectAssociationNotExist = 1101048
	// CCErrorTopoImportAssociation  import association error
	CCErrorTopoImportAssociation = 1101049
	// got multiple association kind with a id
	CCErrorTopoGetMultipleAssoKindInstWithOneID = 1101050
	// delete a pre-defined association kind.
	CCErrorTopoDeletePredefinedAssociationKind = 1101051
	// create new instance for a new association, but association map is 1:1
	CCErrorTopoCreateMultipleInstancesForOneToOneAssociation = 1101052
	// the object has associate to another object, or has been associated by another one.
	CCErrorTopoObjectHasAlreadyAssociated = 1101053
	// update a pre-defined association, it's forbidden.
	CCErrorTopoUpdatePredefinedAssociation = 1101054
	// can not delete a pre-defined association.
	CCErrorTopoDeletePredefinedAssociation = 1101055
	// association do not exist.
	CCErrorTopoAssociationDoNotExist = 1101056
	// create model's instance patch, but instance's data missing field bk_inst_name
	CCErrorTopoObjectInstanceMissingInstanceNameField = 1101057
	// object instance's bk_inst_name filed is not string
	CCErrorTopoInvalidObjectInstanceNameFieldValue = 1101058
	// create model's instance patch, but instance's name is duplicate.
	CCErrorTopoMutipleObjectInstanceName = 1101059
	// association kind has already been instantiated
	CCErrorTopoAssociationKindHasBeenUsed = 1101060
	// create new instance for a new association, but association map is 1:n
	CCErrorTopoCreateMultipleInstancesForOneToManyAssociation = 1101061

	CCErrTopoAppDeleteFailed                       = 1001031
	CCErrTopoAppUpdateFailed                       = 1001032
	CCErrTopoAppSearchFailed                       = 1001033
	CCErrTopoAppCreateFailed                       = 1001034
	CCErrTopoForbiddenToDeleteModelFailed          = 1001035
	CCErrTopoMainlineCreatFailed                   = 1001037
	CCErrTopoMainlineDeleteFailed                  = 1001038
	CCErrTopoMainlineSelectFailed                  = 1001039
	CCErrTopoTopoSelectFailed                      = 1001040
	CCErrTopoUserGroupCreateFailed                 = 1001041
	CCErrTopoUserGroupDeleteFailed                 = 1001042
	CCErrTopoUserGroupUpdateFailed                 = 1001043
	CCErrTopoUserGroupSelectFailed                 = 1001044
	CCErrTopoUserGroupPrivilegeUpdateFailed        = 1001045
	CCErrTopoUserGroupPrivilegeSelectFailed        = 1001046
	CCErrTopoUserPrivilegeSelectFailed             = 1001047
	CCErrTopoRolePrivilegeCreateFailed             = 1001048
	CCErrTopoDeleteMainLineObjectAndInstNameRepeat = 1001049
	CCErrHostNotAllowedToMutiBiz                   = 1001050
	CCErrTopoGraphicsSearchFailed                  = 1001051
	CCErrTopoGraphicsUpdateFailed                  = 1001052
	CCErrTopoObjectUniqueCreateFailed              = 1001060
	CCErrTopoObjectUniqueUpdateFailed              = 1001061
	CCErrTopoObjectUniqueDeleteFailed              = 1001062
	CCErrTopoObjectUniqueSearchFailed              = 1001063
	CCErrTopoObjectPropertyNotFound                = 1001064
	CCErrTopoObjectPropertyUsedByUnique            = 1001065
	CCErrTopoObjectUniqueKeyKindInvalid            = 1001066
	CCErrTopoObjectUniquePresetCouldNotDelOrEdit   = 1001067
	CCErrTopoObjectUniqueCanNotHasMutiMustCheck    = 1001068
	CCErrTopoObjectUniqueShouldHaveMoreThanOne     = 1001069

	CCErrTopoMulueIDNotfoundFailed = 1101080
	CCErrTopoBkAppNotAllowedDelete = 1101081
	// CCErrorTopoAssociationKindMainlineUnavailable can't use bk_mainline in this case
	CCErrorTopoAssociationKindMainlineUnavailable = 1101082
	// CCErrorTopoAssociationKindInconsistent means AssociationKind parameter Inconsistent with caller method
	CCErrorTopoAssociationKindInconsistent = 1101083
	// CCErrorTopoModleStopped means model have been stopped to use
	CCErrorTopoModleStopped = 1101084
	// objectcontroller 1102XXX

	// CCErrObjectPropertyGroupInsertFailed failed to save the property group
	CCErrObjectPropertyGroupInsertFailed = 1102000
	// CCErrObjectPropertyGroupDeleteFailed failed to delete the property group
	CCErrObjectPropertyGroupDeleteFailed = 1102001
	// CCErrObjectPropertyGroupSelectFailed failed to select the property group
	CCErrObjectPropertyGroupSelectFailed = 1102002
	// CCErrObjectPropertyGroupUpdateFailed failed to update the filed
	CCErrObjectPropertyGroupUpdateFailed = 1102003

	CCErrObjectCreateInstFailed       = 1102004
	CCErrObjectUpdateInstFailed       = 1102005
	CCErrObjectDeleteInstFailed       = 1102006
	CCErrObjectSelectInstFailed       = 1102007
	CCErrObjectSelectIdentifierFailed = 1102008

	// CCErrObjectDBOpErrno failed to operation database
	CCErrObjectDBOpErrno = 1102004

	// event_server 1103XXX
	// CCErrEventSubscribeInsertFailed failed to save the Subscribe
	CCErrEventSubscribeInsertFailed = 1103000

	// CCErrEventSubscribeDeleteFailed failed to delete the Subscribe
	CCErrEventSubscribeDeleteFailed = 1103001

	// CCErrEventSubscribeSelectFailed failed to select the Subscribe
	CCErrEventSubscribeSelectFailed = 1103002

	// CCErrEventSubscribeUpdateFailed failed to update the filed
	CCErrEventSubscribeUpdateFailed = 1103003

	// CCErrEventSubscribePingFailed failed to ping the target
	CCErrEventSubscribePingFailed = 1103004
	// CCErrEventSubscribePingFailed failed to telnet the target
	CCErrEventSubscribeTelnetFailed = 1103005
	// CCErrEventOperateSuccessBUtSentEventFailed failed to sent event
	CCErrEventPushEventFailed = 1103006

	// host 1104XXX
	CCErrHostModuleRelationAddFailed = 1104000

	// migrate 1105XXX
	//  CCErrCommMigrateFailed failed to migrate
	CCErrCommMigrateFailed = 1105000

	// hostcontroller 1106XXX
	CCErrHostSelectInst                  = 1106000
	CCErrHostCreateInst                  = 1106002
	CCErrHostGetSnapshot                 = 1106003
	CCErrHostTransferModule              = 1106004
	CCErrDelDefaultModuleHostConfig      = 1106005
	CCErrGetModule                       = 1106006
	CCErrDelOriginHostModuelRelationship = 1106007
	CCErrGetOriginHostModuelRelationship = 1106008
	CCErrTransferHostFromPool            = 1106009
	CCErrAlreadyAssign                   = 1106010
	CCErrNotBelongToIdleModule           = 1106011
	CCErrTransfer2ResourcePool           = 1106012
	CCErrCreateUserCustom                = 1106013
	CCErrHostFavouriteQueryFail          = 1106014
	CCErrHostFavouriteCreateFail         = 1106015
	CCErrHostFavouriteUpdateFail         = 1106016
	CCErrHostFavouriteDeleteFail         = 1106017
	CCErrHostFavouriteDupFail            = 1106018
	CCErrHostGetSnapshotChannelEmpty     = 1106019
	CCErrHostGetSnapshotChannelClose     = 1106020

	// proccontroller 1107XXX
	CCErrProcDeleteProc2Module   = 1107001
	CCErrProcCreateProc2Module   = 1107002
	CCErrProcSelectProc2Module   = 1107003
	CCErrProcCreateProcConf      = 1107004
	CCErrProcDeleteProcConf      = 1107005
	CCErrProcGetProcConf         = 1107006
	CCErrProcUpdateProcConf      = 1107007
	CCErrProcCreateInstanceModel = 1107008
	CCErrProcGetInstanceModel    = 1107009
	CCErrProcDeleteInstanceModel = 1107010
	CCErrProcDeleteProc2Template = 1107011
	CCErrProcCreateProc2Template = 1107012
	CCErrProcSelectProc2Template = 1107013

	// procserver 1108XXX
	CCErrProcSearchDetailFaile          = 1108001
	CCErrProcBindToMoudleFaile          = 1108002
	CCErrProcUnBindToMoudleFaile        = 1108003
	CCErrProcSelectBindToMoudleFaile    = 1108004
	CCErrProcUpdateProcessFaile         = 1108005
	CCErrProcSearchProcessFaile         = 1108006
	CCErrProcDeleteProcessFaile         = 1108007
	CCErrProcCreateProcessFaile         = 1108008
	CCErrProcFieldValidFaile            = 1108009
	CCErrProcGetByApplicationIDFail     = 1108010
	CCErrProcGetByIP                    = 1108011
	CCErrProcOperateFaile               = 1108012
	CCErrProcBindWithModule             = 1108013
	CCErrProcDeleteTemplateFail         = 1108014
	CCErrProcUpdateTemplateFail         = 1108015
	CCErrProcSearchTemplateFail         = 1108016
	CCErrProcBindToTemplateFailed       = 1108017
	CCErrProcUnBindToTemplateFailed     = 1108018
	CCErrProcSelectBindToTemplateFailed = 1108019
	CCErrProcQueryTaskInfoFail          = 1108020
	CCErrProcQueryTaskWaitOPFail        = 1108021
	CCErrProcQueryTaskOPErrFail         = 1108022
	CCErrProcCreateTemplateFail         = 1108023

	// auditlog 1109XXX
	CCErrAuditSaveLogFaile      = 1109001
	CCErrAuditTakeSnapshotFaile = 1109001

	//hostserver
	CCErrHostGetFail              = 1110001
	CCErrHostUpdateFail           = 1110002
	CCErrHostUpdateFieldFail      = 1110003
	CCErrHostCreateFail           = 1110004
	CCErrHostModifyFail           = 1110005
	CCErrHostDeleteFail           = 1110006
	CCErrHostFiledValdFail        = 1110007
	CCErrHostNotFound             = 1110008
	CCErrHostLength               = 1110009
	CCErrHostDetailFail           = 1110010
	CCErrHostSnap                 = 1110011
	CCErrHostFeildValidFail       = 1110012
	CCErrHostFavCreateFail        = 1110013
	CCErrHostEmptyFavName         = 1110014
	CCErrHostFavUpdateFail        = 1110015
	CCErrHostFavDeleteFail        = 1110016
	CCErrHostFavGetFail           = 1110017
	CCErrHostHisCreateFail        = 1110018
	CCErrHostHisGetFail           = 1110019
	CCErrHostCustomCreateFail     = 1110020
	CCErrHostCustomGetFail        = 1110021
	CCErrHostCustomGetDefaultFail = 1110022
	CCErrHostNotINAPP             = 1110023
	CCErrHostNotINAPPFail         = 1110024
	CCErrHostDELResourcePool      = 1110025
	CCErrHostAddRelationFail      = 1110026
	CCErrHostMoveResourcePoolFail = 1110027
	CCErrHostEditRelationPoolFail = 1110028
	CCErrAddHostToModule          = 1110029
	CCErrAddHostToModuleFailStr   = 1110030

	// hostserver api machinery new error code
	CCErrAddUserCustomQueryFaild       = 1110040
	CCErrUpdateUserCustomQueryFaild    = 1110041
	CCErrDeleteUserCustomQueryFaild    = 1110042
	CCErrSearchUserCustomQueryFaild    = 1110043
	CCErrGetUserCustomQueryDetailFaild = 1110044
	CCErrHostModuleConfigFaild         = 1110045
	CCErrHostGetSetFaild               = 1110046
	CCErrHostGetAPPFail                = 1110047
	CCErrHostAPPNotFoundFail           = 1110048
	CCErrHostGetModuleFail             = 1110049
	CCErrHostAgentStatusFail           = 1110050
	CCErrHostSearchNeedObjectInstIDErr = 1110058

	//web  1111XXX
	CCErrWebFileNoFound      = 1111001
	CCErrWebFileSaveFail     = 1111002
	CCErrWebOpenFileFail     = 1111003
	CCErrWebFileContentEmpty = 1111004
	CCErrWebFileContentFail  = 1111005
	CCErrWebGetHostFail      = 1111006
	CCErrWebCreateEXCELFail  = 1111007
	CCErrWebGetObjectFail    = 1111008

	// CCErrApiServerV2AppNameLenErr app name must be 1-32 len
	CCErrAPIServerV2APPNameLenErr = 1170001

	// CCErrAPIServerV2DirectErr  disply error
	CCErrAPIServerV2DirectErr = 1170002

	// CCErrAPIServerV2SetNameLenErr  set name must be < 24 len
	CCErrAPIServerV2SetNameLenErr = 1170003

	// CCErrAPIServerV2MultiModuleIDErr  single module id  is int
	CCErrAPIServerV2MultiModuleIDErr = 1170004

	// CCErrAPIServerV2MultiSetIDErr  single set id is int
	CCErrAPIServerV2MultiSetIDErr = 1170005

	// CCErrAPIServerV2OSTypeErr osType must be linux or windows
	CCErrAPIServerV2OSTypeErr = 1170006

	// CCErrAPIServerV2HostModuleContainDefaultModuleErr  translate host to multiple module not contain default module
	CCErrAPIServerV2HostModuleContainDefaultModuleErr = 1170007

	/** TODO: 以下错误码需要改造 **/

	//http
	CC_Err_Comm_http_DO              = 2000
	CC_Err_Comm_http_DO_STR          = "do http request failed!"
	CC_Err_Comm_http_ReadReqBody     = 2002
	CC_Err_Comm_http_ReadReqBody_STR = "read http request body failed!"

	//json
	CC_ERR_Comm_JSON_DECODE     = 3001
	CC_ERR_Comm_JSON_DECODE_STR = "json decode failed!"
	CC_ERR_Comm_JSON_ENCODE     = 3002
	CC_ERR_Comm_JSON_ENCODE_STR = "json encode failed!"

	CC_Err_Comm_APP_QUERY_FAIL = 4008
	//user custom
	CC_Err_Comm_USER_CUSTOM_SAVE_FAIL = 5000

	Json_Marshal_ERR     = 9000
	Json_Marshal_ERR_STR = "json marshal error"
)
