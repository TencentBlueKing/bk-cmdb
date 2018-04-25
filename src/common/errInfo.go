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

	// CCErrCommParamsLostField the prameter not specified
	CCErrCommParamsLostField = 1199008

	// CCErrCommParamsNeedInt the parameter must be of tyep int
	CCErrCommParamsNeedInt = 1199009

	// CCErrCommParamsNeedSet the parameter unassigned
	CCErrCommParamsNeedSet = 1199010

	// CCErrCommParamsIsInvalid the parameter is invalid or nonexistent
	CCErrCommParamsIsInvalid = 1199011

	// CCErrCommUniqueCheckFailed the uniqueness validation fails
	CCErrCommUniqueCheckFailed = 1199012

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

	CCErrTopoAppDeleteFailed                = 1101031
	CCErrTopoAppUpdateFailed                = 1101032
	CCErrTopoAppSearchFailed                = 1101033
	CCErrTopoAppCreateFailed                = 1101034
	CCErrTopoForbiddenToDeleteModelFailed   = 1101035
	CCErrTopoMainlineCreatFailed            = 1101037
	CCErrTopoMainlineDeleteFailed           = 1101038
	CCErrTopoMainlineSelectFailed           = 1101039
	CCErrTopoTopoSelectFailed               = 1101040
	CCErrTopoUserGroupCreateFailed          = 1101041
	CCErrTopoUserGroupDeleteFailed          = 1101042
	CCErrTopoUserGroupUpdateFailed          = 1101043
	CCErrTopoUserGroupSelectFailed          = 1101044
	CCErrTopoUserGroupPrivilegeUpdateFailed = 1101045
	CCErrTopoUserGroupPrivilegeSelectFailed = 1101046
	CCErrTopoUserPrivilegeSelectFailed      = 1101047
	CCErrTopoRolePrivilegeCreateFailed      = 1101048

	CCErrTopoPlatQueryFailed  = 1101049
	CCErrTopoPlatDeleteFailed = 1101050
	CCErrTopoPlatCreateFailed = 1101051
	CCErrTopoHostInPlatFailed = 1101052

	CCErrTopoMulueIDNotfoundFailed = 1101080
	CCErrTopoBkAppNotAllowedDelete = 1101081

	// objectcontroller 1102XXX

	// CCErrObjectPropertyGroupInsertFailed failed to save the property group
	CCErrObjectPropertyGroupInsertFailed = 1102000
	// CCErrObjectPropertyGroupDeleteFailed failed to delete the property group
	CCErrObjectPropertyGroupDeleteFailed = 1102001
	// CCErrObjectPropertyGroupSelectFailed failed to select the property group
	CCErrObjectPropertyGroupSelectFailed = 1102002
	// CCErrObjectPropertyGroupUpdateFailed failed to update the filed
	CCErrObjectPropertyGroupUpdateFailed = 1102003

	CCErrObjectCreateInstFailed = 1102004
	CCErrObjectUpdateInstFailed = 1102005
	CCErrObjectDeleteInstFailed = 1102006
	CCErrObjectSelectInstFailed = 1102007

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

	// CCErrEventSubscribePingFailed failed to ping the filed
	CCErrEventSubscribePingFailed = 1103004
	// CCErrEventSubscribePingFailed failed to telnet the filed
	CCErrEventSubscribeTelnetFailed = 1103005

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

	// proccontroller 1107XXX
	CCErrProcDeleteProc2Module = 1107001
	CCErrProcCreateProc2Module = 1107002
	CCErrProcSelectProc2Module = 1107003

	// procserver 1108XXX
	CCErrProcSearchDetailFaile       = 1108001
	CCErrProcBindToMoudleFaile       = 1108002
	CCErrProcUnBindToMoudleFaile     = 1108003
	CCErrProcSelectBindToMoudleFaile = 1108004
	CCErrProcUpdateProcessFaile      = 1108005
	CCErrProcSearchProcessFaile      = 1108006
	CCErrProcDeleteProcessFaile      = 1108007
	CCErrProcCreateProcessFaile      = 1108008
	CCErrProcGetFail                 = 1108009

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

	// api server v2 error 1170xxx, follow-up will be deleted

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
	CC_Err_Comm_http_DO               = 2000
	CC_Err_Comm_http_DO_STR           = "do http request failed!"
	CC_Err_Comm_http_Input_Params     = 2001
	CC_Err_Comm_http_Input_Params_STR = "input params error!"
	CC_Err_Comm_http_ReadReqBody      = 2002
	CC_Err_Comm_http_ReadReqBody_STR  = "read http request body failed!"

	//json
	CC_ERR_Comm_JSON_DECODE     = 3001
	CC_ERR_Comm_JSON_DECODE_STR = "json decode failed!"
	CC_ERR_Comm_JSON_ENCODE     = 3002
	CC_ERR_Comm_JSON_ENCODE_STR = "json encode failed!"
	CC_ERR_Comm_JSON_GET        = 3003
	cc_ERR_Comm_JSON_GET_STR    = "get data from json failed!"

	Json_Marshal_ERR     = 9000
	Json_Marshal_ERR_STR = "json marshal error"
)
