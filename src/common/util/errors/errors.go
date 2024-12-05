/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package errors is error related util package
package errors

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/storage/driver/mongodb"
)

// ConvDBInsertError convert db insert error to
func ConvDBInsertError(kit *rest.Kit, err error) errors.CCErrorCoder {
	if err == nil {
		return nil
	}

	if mongodb.IsDuplicatedError(err) {
		return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
	}

	return kit.CCError.CCError(common.CCErrCommDBInsertFailed)
}
