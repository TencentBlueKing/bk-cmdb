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

// Package audit is the data reporting service for audit center
package audit

import (
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/cacheservice/audit/config"
	"configcenter/src/storage/stream"
)

// Audit is the data reporting service for audit center
type Audit struct {
	conf   *config.Config
	loopW  stream.LoopInterface
	client *auditClient
}

// RunAuditDataReporting run audit data reporting
func RunAuditDataReporting(conf *config.Config, loopW stream.LoopInterface) error {
	if conf == nil || !conf.Enabled {
		blog.Info("audit data reporting is disabled")
		return nil
	}

	if err := conf.Validate(); err != nil {
		blog.Errorf("audit data reporting config(%+v) is invalid, err: %v", *conf, err)
		return err
	}

	client, err := newAuditClient(conf)
	if err != nil {
		blog.Errorf("init audit client failed, err: %v", err)
		return err
	}

	audit := &Audit{
		conf:   conf,
		loopW:  loopW,
		client: client,
	}

	if err = audit.watchAudit(); err != nil {
		blog.Errorf("watch audit event and push to audit center failed, err: %v", err)
		return err
	}

	return nil
}
