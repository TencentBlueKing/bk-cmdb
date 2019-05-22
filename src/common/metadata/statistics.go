/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import "time"

type StatisticChartInfo struct {
	SticId     int64     `json:"bk_stic_id"`
	ObjId      string    `json:"bk_obj_id"`
	Title      string    `json:"title"`
	ChartType  string    `json:"bk_chart_type"`
	Field      string    `json:"bk_field"`
	Width      string    `json:"bk_width"`
	Position   string    `json:"position"`
	CreateTime time.Time `json:"create_time"`
	LastTime   time.Time `json:"last_time"`
}
