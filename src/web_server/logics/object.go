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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/chalice-1831/zip"
	yl "github.com/ghodss/yaml"
)

// GetObjectData get object data
func (lgc *Logics) GetObjectData(objID string, header http.Header, modelBizID int64) ([]interface{}, error) {

	condition := mapstr.MapStr{
		"condition": []string{
			objID,
		},
		common.BKAppIDField: modelBizID,
	}

	result, err := lgc.Engine.CoreAPI.ApiServer().GetObjectData(context.Background(), header, condition)
	if nil != err {
		return nil, fmt.Errorf("get object data failed, err: %v", err)
	}

	if !result.Result {
		return nil, fmt.Errorf("get object data failed, but got err: %s", result.ErrMsg)
	}

	return result.Data[objID].Attr, nil

}

func GetPropertyFieldType(lang language.DefaultCCLanguageIf) map[string]string {
	var fieldType = map[string]string{
		"bk_property_id":         lang.Language("val_type_text"), // "文本",
		"bk_property_name":       lang.Language("val_type_text"), // "文本",
		"bk_property_type":       lang.Language("val_type_text"), // "文本",
		"bk_property_group_name": lang.Language("val_type_text"), //  文本
		"option":                 lang.Language("val_type_text"), // "文本",
		"unit":                   lang.Language("val_type_text"), // "文本",
		"description":            lang.Language("val_type_text"), // "文本",
		"placeholder":            lang.Language("val_type_text"), // "文本",
		"editable":               lang.Language("val_type_bool"), // "布尔",
		"isrequired":             lang.Language("val_type_bool"), // "布尔",
		"isreadonly":             lang.Language("val_type_bool"), // "布尔",
	}
	return fieldType
}

func GetPropertyFieldDesc(lang language.DefaultCCLanguageIf) map[string]string {

	var fields = map[string]string{
		"bk_property_id":         lang.Language("web_en_name_required"),       // "英文名(必填)",
		"bk_property_name":       lang.Language("web_bk_alias_name_required"), // "中文名(必填)",
		"bk_property_type":       lang.Language("web_bk_data_type_required"),  // "数据类型(必填)",
		"bk_property_group_name": lang.Language("property_group"),             //  字段分组
		"option":                 lang.Language("property_option"),            // "数据配置",
		"unit":                   lang.Language("unit"),                       // "单位",
		"description":            lang.Language("desc"),                       // "描述",
		"placeholder":            lang.Language("placeholder"),                // "提示",
		"editable":               lang.Language("is_editable"),                // "是否可编辑",
		"isrequired":             lang.Language("property_is_required"),       // "是否必填",
		"isreadonly":             lang.Language("property_is_readonly"),       // "是否只读",
	}

	return fields
}

func ConvAttrOption(attrItems map[int]map[string]interface{}) {
	for index, attr := range attrItems {

		option, ok := attr[common.BKOptionField].(string)
		if false == ok {
			continue
		}

		if "\"\"" == option {
			option = ""
			attrItems[index][common.BKOptionField] = option
			continue
		}
		fieldType, _ := attr[common.BKPropertyTypeField].(string)
		if common.FieldTypeEnum != fieldType && common.FieldTypeInt != fieldType && common.FieldTypeList != fieldType {
			continue
		}

		var iOption interface{}
		err := json.Unmarshal([]byte(option), &iOption)
		if nil == err {
			attrItems[index][common.BKOptionField] = iOption
		}
	}
}

// GetObjectCount search object count
func (lgc *Logics) GetObjectCount(ctx context.Context, header http.Header, cond *metadata.ObjectCountParams) (
	*metadata.ObjectCountResult, error) {
	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	objIDs := cond.Condition.ObjectIDs
	if len(objIDs) > 20 {
		blog.Errorf("field obj_ids number %d exceeds the max number 20, rid: %s", len(objIDs), rid)
		return nil, defErr.CCErrorf(common.CCErrCommValExceedMaxFailed, "obj_ids", 20)
	}
	objArray, nonexistentObjects, err := lgc.ProcessObjectIDArray(ctx, header, objIDs)
	if err != nil {
		blog.Errorf("process object array failed, err %s, rid: %s", err.Error(), rid)
		return nil, err
	}
	resp := &metadata.ObjectCountResult{}
	for _, nonexistentObject := range nonexistentObjects {
		objCount := metadata.ObjectCountDetails{
			ObjectID: nonexistentObject,
			Error:    defErr.CCError(common.CCErrorModelNotFound).Error(),
		}
		resp.Data = append(resp.Data, objCount)
	}

	var wg sync.WaitGroup
	var apiErr errors.CCErrorCoder
	existObjResult := make([]metadata.ObjectCountDetails, len(objArray))
	pipeline := make(chan bool, 10)
	for idx, objID := range objArray {
		objCount := metadata.ObjectCountDetails{ObjectID: objID}
		pipeline <- true
		wg.Add(1)
		go func(ctx context.Context, header http.Header, objID string, idx int, objCount metadata.ObjectCountDetails) {
			defer func() {
				wg.Done()
				<-pipeline
			}()
			if metadata.IsCommon(objID) {
				params := &metadata.Condition{Condition: map[string]interface{}{common.BKObjIDField: objID}}
				count, err := lgc.CoreAPI.CoreService().Instance().CountInstances(ctx, header, objID, params)
				if err != nil {
					blog.Errorf("get %s instance count failed, err: %s, rid: %s", objID, err.Error(), rid)
					apiErr = defErr.CCErrorf(common.CCErrCommHTTPDoRequestFailed)
					return
				}
				if count == nil {
					apiErr = defErr.CCErrorf(common.CCErrCommHTTPDoRequestFailed)
					return
				}

				objCount.InstCount = count.Count
			} else {
				tableName := common.GetInstTableName(objID, common.BKDefaultOwnerID)
				count, err := lgc.CoreAPI.CoreService().Count().GetCountByFilter(ctx, header, tableName,
					[]map[string]interface{}{{}})
				if err != nil {
					blog.Errorf("get %s instance count failed, err: %s, rid: %s", objID, err.Error(), rid)
					apiErr = defErr.CCErrorf(common.CCErrCommHTTPDoRequestFailed)
					return
				}
				if len(count) == 0 {
					apiErr = defErr.CCErrorf(common.CCErrCommHTTPDoRequestFailed)
					return
				}
				objCount.InstCount = uint64(count[0])
			}

			existObjResult[idx] = objCount

		}(ctx, header, objID, idx, objCount)
	}

	wg.Wait()
	resp.Data = append(resp.Data, existObjResult...)

	if apiErr != nil {
		return nil, apiErr
	}

	resp.Result = true
	return resp, nil
}

// ProcessObjectIDArray process objectIDs
func (lgc *Logics) ProcessObjectIDArray(ctx context.Context, header http.Header, objectArray []string) ([]string,
	[]string, error) {
	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	objArray := util.RemoveDuplicatesAndEmpty(objectArray)
	objects, err := lgc.CoreAPI.CoreService().Model().ReadModel(ctx, header, &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKObjIDField: map[string]interface{}{
				common.BKDBIN: objArray,
			},
		},
	})
	if err != nil {
		blog.Errorf("get object by object id failed, err: %s, rid: %s", err.Error(), rid)
		return nil, nil, defErr.CCError(common.CCErrTopoModuleSelectFailed)
	}

	if objects.Count == 0 {
		blog.Errorf("none object exist, object id: %v, err: %s, rid: %s", objectArray,
			defErr.CCError(common.CCErrorModelNotFound).Error(), rid)
		return nil, nil, defErr.CCError(common.CCErrorModelNotFound)
	}
	var nonexistentObjects []string
	if objects.Count < int64(len(objArray)) {
		var temp []string
		for _, object := range objects.Info {
			temp = append(temp, object.ObjectID)
		}
		nonexistentObjects = util.StrArrDiff(objArray, temp)
		objArray = temp
	}

	return objArray, nonexistentObjects, nil
}

// BuildExportObjectYaml build export yaml
func (lgc *Logics) BuildExportYaml(header http.Header, expiration int64, data interface{},
	exportType string) ([]byte, error) {

	rid := util.GetHTTPCCRequestID(header)
	user := util.GetUser(header)
	nowTime := time.Now()

	expirationTime := nowTime.Unix()
	if expiration != 0 {
		expirationTime = nowTime.AddDate(0, 0, int(expiration)).Unix()
	}

	aesEncrpytor := cryptor.NewAesEncrpytor(fmt.Sprintf("export%d", nowTime.Unix()))
	aesExpiration, err := aesEncrpytor.Encrypt(fmt.Sprint(expirationTime))
	if err != nil {
		blog.Errorf("encrypt expiration time failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	exportInfo := mapstr.MapStr{
		"creator":     user,
		"create_time": nowTime.Unix(),
		"expire_time": aesExpiration,
		exportType:    data,
	}

	yamlData, err := yl.Marshal(exportInfo)
	if err != nil {
		blog.Errorf("marshal data[%v] into yaml failed, err: %v, rid: %s", data, err, rid)
		return nil, err
	}

	return yamlData, nil
}

// BuildZipFile build zip file with password
func (lgc *Logics) BuildZipFile(header http.Header, zipw *zip.Writer, fileName string, password string,
	data []byte) error {

	rid := util.GetHTTPCCRequestID(header)
	if password != "" {
		w, err := zipw.Encrypt(fileName, password, zip.AES256Encryption)
		if err != nil {
			blog.Errorf("encrypt zip file failed, err: %v, rid: %s", err, rid)
			return err
		}

		if _, err = io.Copy(w, bytes.NewReader(data)); err != nil {
			blog.Errorf("copy data into zip file failed, err: %v, rid: %s", err, rid)
			return err
		}
	} else {
		w, err := zipw.Create(fileName)
		if err != nil {
			blog.Errorf("encrypt zip file failed, err: %v, rid: %s", err, rid)
			return err
		}

		if _, err = io.Copy(w, bytes.NewReader(data)); err != nil {
			blog.Errorf("copy data into zip file failed, err: %v, rid: %s", err, rid)
			return err
		}
	}

	return nil
}

// GetDataFromZipFile get data from zip file
func (lgc *Logics) GetDataFromZipFile(header http.Header, file *zip.File, password string,
	result *metadata.AnalysisResult) error {

	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	if file.IsEncrypted() {
		if len(password) == 0 {
			blog.Errorf("zip is encrypted but no password provided, rid: %s", rid)
			return defErr.Errorf(common.CCErrWebVerifyYamlPwdFail, "no password")
		}

		file.SetPassword(password)
	}

	zipFile, err := file.Open()
	if err != nil {
		blog.Errorf("open zip file failed, err: %+v, rid: %s", err, rid)
		return err
	}

	defer zipFile.Close()

	receiver := new(bytes.Buffer)
	io.Copy(receiver, zipFile)
	if strings.Contains(file.Name, "asst_kind") {
		fileInfo := metadata.AssociationKindYaml{}
		if err := yl.Unmarshal(receiver.Bytes(), &fileInfo); err != nil {
			blog.Errorf("unmarshal zip file failed, err: %+v, rid: %s", err, rid)
			return err
		}

		if err := fileInfo.Validate(); err.ErrCode != 0 {
			blog.Errorf("validate file failed, err: %+v, rid: %s", err, rid)
			return err.ToCCError(defErr)
		}

		result.Data.Asst = append(result.Data.Asst, fileInfo.AsstKind...)
	} else {
		fileInfo := metadata.ObjectYaml{}
		if err := yl.Unmarshal(receiver.Bytes(), &fileInfo); err != nil {
			blog.Errorf("unmarshal zip file failed, err: %+v, rid: %s", err, rid)
			return err
		}
		if err := fileInfo.Validate(); err.ErrCode != 0 {
			blog.Errorf("validate file failed, err: %+v, rid: %s", err, rid)
			return err.ToCCError(defErr)
		}
		result.Data.Object = append(result.Data.Object, fileInfo.Object)
	}

	return nil
}

// CreateOrUpdateAssociationType create association type, if association type exist, update it
func (lgc *Logics) CreateOrUpdateAssociationType(ctx context.Context, header http.Header,
	asst []metadata.AssociationKind) error {

	rid := util.GetHTTPCCRequestID(header)
	asstKindMap := make(map[string]metadata.AssociationKind)
	asstKindID := make([]string, 0)
	for _, item := range asst {
		asstKindMap[item.AssociationKindID] = metadata.AssociationKind{
			AssociationKindID:       item.AssociationKindID,
			AssociationKindName:     item.AssociationKindName,
			DestinationToSourceNote: item.DestinationToSourceNote,
			SourceToDestinationNote: item.SourceToDestinationNote,
			Direction:               item.Direction,
		}
		asstKindID = append(asstKindID, item.AssociationKindID)
	}

	asstQuery := metadata.SearchAssociationTypeRequest{
		Condition: map[string]interface{}{common.AssociationKindIDField: mapstr.MapStr{common.BKDBIN: asstKindID}}}
	rsp, err := lgc.Engine.CoreAPI.ApiServer().SearchAssociationType(ctx, header, asstQuery)
	if err != nil {
		blog.Errorf("search asstkind failed, cond: %v, err: %v, rid: %s", asstQuery, err, rid)
		return err
	}

	existAsstKind := make([]string, 0)
	kindIDMap := make(map[string]int64)
	for _, item := range rsp.Info {
		existAsstKind = append(existAsstKind, item.AssociationKindID)
		kindIDMap[item.AssociationKindID] = item.ID
	}

	needCreateAsst := util.StrArrDiff(asstKindID, existAsstKind)
	for _, item := range existAsstKind {
		updataCond := metadata.UpdateAssociationTypeRequest{
			AsstName:  asstKindMap[item].AssociationKindName,
			SrcDes:    asstKindMap[item].SourceToDestinationNote,
			DestDes:   asstKindMap[item].DestinationToSourceNote,
			Direction: string(asstKindMap[item].Direction),
		}
		if err := lgc.Engine.CoreAPI.ApiServer().UpdateAssociationType(ctx, header, kindIDMap[item],
			updataCond); err != nil {
			blog.Errorf("update asstkind failed, cond: %v,err: %s, rid: %s", updataCond, err.Error(), rid)
			return err
		}
	}

	for _, item := range needCreateAsst {
		createCond := asstKindMap[item]
		if _, err := lgc.Engine.CoreAPI.ApiServer().CreateAssociationType(ctx, header, createCond); err != nil {
			blog.Errorf("create asstkind failed, cond: %v, err: %s, rid: %s", createCond, err.Error(), rid)
			return err
		}
	}

	return nil
}
