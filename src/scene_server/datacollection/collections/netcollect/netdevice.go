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

package netcollect

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"configcenter/src/ac/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
)

// NetCollect collect the net information
type NetCollect struct {
	ctx         context.Context
	db          dal.RDB
	authManager *extensions.AuthManager
}

// NewNetCollect returns a new netcollector
func NewNetCollect(ctx context.Context, db dal.RDB, authManager *extensions.AuthManager) *NetCollect {
	h := &NetCollect{
		ctx:         ctx,
		db:          db,
		authManager: authManager,
	}
	return h
}

// Hash returns hash value base on message.
func (h *NetCollect) Hash(cloudid, ip string) (string, error) {
	if len(cloudid) == 0 {
		return "", fmt.Errorf("can't make hash from invalid message format, cloudid empty")
	}
	if len(ip) == 0 {
		return "", fmt.Errorf("can't make hash from invalid message format, ip empty")
	}

	hash := fmt.Sprintf("%s:%s", cloudid, ip)

	return hash, nil
}

// Mock returns local mock message for testing.
func (h *NetCollect) Mock() string {
	return MockMessage
}

// Analyze implements the Analyzer interface
func (h *NetCollect) Analyze(msg *string) error {
	if msg == nil {
		return fmt.Errorf("message nil")
	}

	data := ReportMessage{}

	if err := json.Unmarshal([]byte(*msg), &data); err != nil {
		return fmt.Errorf("unmarshal message error: %+v", err)
	}

	for _, report := range data.Data {
		if err := h.handleReport(&report); err != nil {
			blog.Errorf("[data-collection][netcollect] handleData failed: %+v", err)
		}
	}

	return nil
}

func (h *NetCollect) handleReport(report *metadata.NetcollectReport) (err error) {
	// TODO compare 若有变化才插入
	if err = h.upsertReport(report); err != nil {
		blog.Errorf("[data-collection][netcollect] upsert association error: %v", err)
		return err
	}

	return nil
}

func (h *NetCollect) upsertReport(report *metadata.NetcollectReport) error {
	existFilter := map[string]interface{}{
		common.BKCloudIDField: report.CloudID,
		common.BKObjIDField:   report.ObjectID,
		common.BKInstKeyField: report.InstKey,
	}

	count, err := h.db.Table(common.BKTableNameNetcollectReport).Find(existFilter).Count(h.ctx)
	if err != nil {
		return err
	}
	if count <= 0 {
		err = h.db.Table(common.BKTableNameNetcollectReport).Insert(h.ctx, report)
		return err
	}

	return h.db.Table(common.BKTableNameNetcollectReport).Update(h.ctx, existFilter, report)
}

// ReportMessage define a netcollect message
type ReportMessage struct {
	Timestamp time.Time                   `json:"timestamp"`
	Dataid    int                         `json:"dataid"`
	Type      string                      `json:"type"`
	Counter   int                         `json:"counter"`
	Build     CollectorBuild              `json:"build"`
	Data      []metadata.NetcollectReport `json:"data"`
}

// CollectorBuild define a netcollector build information
type CollectorBuild struct {
	Version     string `json:"version"`
	BuildCommit string `json:"build_commit"`
	BuildTime   string `json:"build_time"`
	GoVersion   string `json:"go_version"`
}

const MockMessage = `{
    "dataid": 1014,
    "type": "netdevicebeat",
    "counter": 1,
    "Build": {
        "version": "1.0.0",
        "build_commit": "3fb6cb0b5a55cffae028d3df7bee71f90155a2f5",
        "buildtime": "2018-10-03 17:09:00",
        "go_version": "1.11.2"
    },
    "data": [
        {
            "bk_obj_id": "bk_switch",
            "bk_inst_key": "huawei 5789#56-79-9a-ii",
            "bk_host_innerip": "192.168.1.1",
			"bk_cloud_id": 0,
			"last_time": "2018-10-03 17:09:00",
            "attributes": [
                {
                    "bk_property_id": "bk_inst_name",
                    "value": "huawei 5789#56-79-9a-ii"
                }
            ],
            "associations": [
				{
					"bk_asst_inst_name": "192.168.1.1",
                    "bk_asst_obj_id": "bk_host",
                    "bk_asst_obj_name": "主机",
                    "bk_asst_property_id": "bk_host_id"
				}
			]
        }
    ]
}
`
