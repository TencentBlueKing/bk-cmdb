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
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/alexmullins/zip"
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

// GetPropertyFieldType TODO
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

// GetPropertyFieldDesc TODO
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

// ConvAttrOption TODO
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

// BuildExportYaml build export yaml
func (lgc *Logics) BuildExportYaml(header http.Header, expiration int64, data interface{},
	exportType string) ([]byte, error) {

	rid := util.GetHTTPCCRequestID(header)
	nowTime := time.Now().Local()

	expirationTime := nowTime.UnixNano()
	if expiration != 0 {
		expirationTime = nowTime.AddDate(0, 0, int(expiration)).UnixNano()
	}

	exportInfo := mapstr.MapStr{
		"create_time": nowTime.UnixNano(),
		"expire_time": expirationTime,
		exportType:    data,
	}

	yamlData, err := yl.Marshal(exportInfo)
	if err != nil {
		blog.Errorf("marshal data[%+v] into yaml failed, err: %v, rid: %s", data, err, rid)
		return nil, err
	}

	return yamlData, nil
}

// BuildZipFile build zip file with password
func (lgc *Logics) BuildZipFile(header http.Header, zipw *zip.Writer, fileName string, password string,
	data []byte) error {

	rid := util.GetHTTPCCRequestID(header)
	fh := &zip.FileHeader{
		Name:   fileName,
		Method: zip.Deflate,
	}
	fh.SetModTime(time.Now())
	if password != "" {
		fh.SetPassword(password)
	}
	w, err := zipw.CreateHeader(fh)
	if err != nil {
		blog.Errorf("encrypt zip file failed, err: %v, rid: %s", err, rid)
		return err
	}

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		blog.Errorf("copy data into zip file failed, err: %v, rid: %s", err, rid)
		return err
	}

	return nil
}

// GetDataFromZipFile get data from zip file
func (lgc *Logics) GetDataFromZipFile(header http.Header, file *zip.File, password string,
	result *metadata.AnalysisResult) (int, error) {

	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	if file.IsEncrypted() {
		if len(password) == 0 {
			blog.Errorf("zip is encrypted but no password provided, rid: %s", rid)
			return common.CCErrWebVerifyYamlPwdFail, defErr.CCErrorf(common.CCErrWebVerifyYamlPwdFail, "no password")
		}

		file.SetPassword(password)
	}

	zipFile, err := file.Open()
	if err != nil {
		blog.Errorf("open zip file failed, err: %v, rid: %s", err, rid)
		if strings.Contains(err.Error(), "invalid password") {
			return common.CCErrWebVerifyYamlPwdFail, err
		}
		return common.CCErrWebAnalysisZipFileFail, err
	}

	defer zipFile.Close()

	receiver := new(bytes.Buffer)
	if _, err := io.Copy(receiver, zipFile); err != nil {
		blog.Errorf("copy zip file into receicer failed, err: %v, rid: %s", err, rid)
		return common.CCErrWebAnalysisZipFileFail, err
	}

	// if encryption method use ZipCrypto, decrypt it with incorrect password receiver will be empty
	if len(receiver.Bytes()) == 0 {
		blog.Errorf("get file info failed, zip use ZipCrypto, password incorrect, rid: %s", rid)
		return common.CCErrWebVerifyYamlPwdFail, defErr.CCErrorf(common.CCErrWebVerifyYamlPwdFail, "invalid password")
	}

	if strings.Contains(file.Name, "asst_kind") {
		fileInfo := metadata.AssociationKindYaml{}
		if err := yl.Unmarshal(receiver.Bytes(), &fileInfo); err != nil {
			blog.Errorf("unmarshal zip file failed, err: %v, rid: %s", err, rid)
			return common.CCErrWebAnalysisZipFileFail, err
		}

		if err := fileInfo.Validate(); err.ErrCode != 0 {
			blog.Errorf("validate file failed, fileName: %s, err: %v, rid: %s", file.Name, err, rid)
			return common.CCErrWebAnalysisZipFileFail, err.ToCCError(defErr)
		}

		result.Data.Asst = append(result.Data.Asst, fileInfo.AsstKind...)
	} else {
		fileInfo := metadata.ObjectYaml{}
		if err := yl.Unmarshal(receiver.Bytes(), &fileInfo); err != nil {
			blog.Errorf("unmarshal zip file failed, err: %v, rid: %s", err, rid)
			return common.CCErrWebAnalysisZipFileFail, err
		}
		if err := fileInfo.Validate(); err.ErrCode != 0 {
			blog.Errorf("validate file failed, fileName: %s, err: %v, rid: %s", file.Name, err, rid)
			return common.CCErrWebAnalysisZipFileFail, err.ToCCError(defErr)
		}
		result.Data.Object = append(result.Data.Object, fileInfo.Object)
	}

	return 0, nil
}
