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
	"net/http"

	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/language"
)

func GetImportNetDevices(f *xlsx.File, url string, header http.Header, defLang language.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	return nil, nil, nil
}

//BuildNetDeviceExcelTemplate  return httpcode, error
func BuildNetDeviceExcelTemplate(url, objID, filename string, header http.Header, defLang language.DefaultCCLanguageIf) error {
	var file *xlsx.File
	file = xlsx.NewFile()
	sheet, err := file.AddSheet(common.BKNetDevice)
	if err != nil {
		blog.Errorf("get %s fields error:", objID, err.Error())
		return err
	}

	fields := getNetDevicefield(defLang)

	blog.V(5).Infof("BuildNetDeviceExcelTemplate fields count:%d", fields)

	productExcelHealer(fields, nil, sheet, defLang)
	ProductExcelCommentSheet(file, defLang)

	if err = file.Save(filename); nil != err {
		return err
	}

	return nil
}

func getNetDevicefield(lang language.DefaultCCLanguageIf) map[string]Property {

	return map[string]Property{
		common.BKDeviceNameField: Property{
			Name: lang.Language("import_device_comment_device_name"), ID: common.BKDeviceNameField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 0, IsRequire: true,
		},
		common.BKDeviceModelField: Property{
			Name: lang.Language("import_device_comment_device_model"), ID: common.BKDeviceModelField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 1, IsRequire: true,
		},
		common.BKObjNameField: Property{
			Name: lang.Language("import_device_comment_obj_name"), ID: common.BKObjNameField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 2, IsRequire: true,
		},
		common.BKVendorField: Property{
			Name: lang.Language("import_device_comment_vendor"), ID: common.BKVendorField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 3, IsRequire: true,
		},
	}
}
