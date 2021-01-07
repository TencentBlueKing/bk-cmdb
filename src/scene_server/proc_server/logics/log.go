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

	"configcenter/src/common/auditoplog"
	meta "configcenter/src/common/metadata"
)

type TemplateLog struct {
	logic   *Logics
	header  http.Header
	ownerID string
	ip      string
	Content *meta.Content
}

func (lgc *Logics) NewTemplate(pheader http.Header, ownerID string) *TemplateLog {
	return &TemplateLog{
		logic:   lgc,
		header:  pheader,
		ownerID: ownerID,
		Content: new(meta.Content),
	}
}

func (h *TemplateLog) WithPrevious(tempID int64, headers []meta.Header) error {
	var err error
	if headers != nil || len(headers) != 0 {
		h.Content.Headers = headers
	} else {
		h.Content.Headers, err = h.logic.GetTemplateAttributes(h.ownerID, h.header)
		if err != nil {
			return err
		}
	}

	h.Content.PreData, err = h.logic.GetTemplateInstanceDetails(h.header, h.ownerID, tempID)
	if err != nil {
		return err
	}

	return nil
}

func (h *TemplateLog) WithCurrent(tempID int64) error {
	var err error
	h.Content.CurData, err = h.logic.GetTemplateInstanceDetails(h.header, h.ownerID, tempID)
	if err != nil {
		return err
	}

	return nil
}

func (h *TemplateLog) AuditLog(tempID int64) *auditoplog.AuditLogExt {
	return &auditoplog.AuditLogExt{
		ID:      tempID,
		Content: h.Content,
		ExtKey:  h.ip,
	}
}
