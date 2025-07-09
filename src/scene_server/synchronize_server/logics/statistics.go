/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package logics

import (
	"configcenter/src/common/metadata"
)

type statisticsItem struct {
	//  handle error,coreservice error message
	error []metadata.ExceptionResult
	// handle error, current server  message
	otherError []metadata.ExceptionResult
	// accept data count
	update int64
	// delete data count
	delete int64
}

type instanceStatistics struct {
	host   statisticsItem
	set    statisticsItem
	module statisticsItem
	plat   statisticsItem
	proc   statisticsItem
	// map[objID] statisticsItem
	object map[string]statisticsItem
	// synchronize with a generic model
	instance statisticsItem
}

type modelStatistics struct {
	model statisticsItem
}

type associationStatistics struct {
	moduleHostConfig statisticsItem
	association      map[string]statisticsItem
}

type statistics struct {
	instance instanceStatistics
	model    modelStatistics
}
