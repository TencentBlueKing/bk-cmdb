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

package index

import (
	"configcenter/src/common"
	"configcenter/src/storage/dal/types"
	"fmt"
)

func InstanceIndex() []types.Index {
	return instanceDefaultIndex
}

func InstanceAssoicationIndex() []types.Index {
	return assoicationDefaultIndex
}

func CCFieldTypeToDBType(typ string) string {
	switch typ {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar:
		return "string"
	case common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeEnum, common.FieldTypeUser, common.FieldTypeTimeZone,
		common.FieldTypeList, common.FieldTypeOrganization:
		return "number"
	case common.FieldTypeDate, common.FieldTypeTime:
		return "date"
	case common.FieldTypeBool:
		return "bool"

	}

	// other type not support
	return ""
}

func GetUniqueIndexNameByID(id uint64) string {
	return fmt.Sprintf("%s%d", common.CCLogicUniqueIdxNamePrefix, id)
}
