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
	"net/http"
	"sync"
	"time"

	"configcenter/src/common/metadata"
)

type chanItem struct {
	ctx       context.Context
	eventData *metadata.EventInst
	opFunc    func(ctx context.Context, eventData *metadata.EventInst) error
	retry     int
}

type refreshHostInstModuleID struct {
	Header   http.Header `json:"header"`
	AppID    int64       `json:"bk_biz_id"`
	ModuleID int64       `json:"bk_module_id"`
}

type opProcTask struct {
	GseTaskIDArr []string    `json:"gse_task_id"`
	TaskID       string      `json:"task_id"`
	OpTime       time.Time   `json:"op_time"`
	Header       http.Header `json:"header"`
}

var (
	handEventDataChan           chan chanItem
	chnOpLock                   *sync.Once = new(sync.Once)
	initDataLock                *sync.Once = new(sync.Once)
	refreshHostInstModuleIDChan chan *refreshHostInstModuleID
	gseOPProcTaskChan           chan *opProcTask
	maxRefreshModuleData        int           = 100
	maxEventDataChan            int           = 10000
	retry                       int           = 3
	SPOPINTERVAL                time.Duration = time.Second * 30
	timedTriggerTime            time.Duration = time.Hour * 1
	timedTriggerLockExpire      time.Duration = time.Hour * 23
	GETTASKIDSPOPINTERVAL       time.Duration = time.Second * 5
	timedTriggerTaskTime        time.Duration = time.Minute * 20
)
