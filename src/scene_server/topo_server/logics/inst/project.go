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

package inst

import (
	"errors"
	"strconv"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/ac/iam"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/google/uuid"
)

// ProjectOperationInterface project operation interface
type ProjectOperationInterface interface {
	// SetProxy set proxy
	SetProxy(inst InstOperationInterface)
	// CreateProject create project
	CreateProject(kit *rest.Kit, data []mapstr.MapStr) ([]int64, error)
}

// NewProjectOperation create project instance
func NewProjectOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) ProjectOperationInterface {
	return &project{
		clientSet:   client,
		authManager: authManager,
	}
}

type project struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	inst        InstOperationInterface
}

// SetProxy set proxy
func (p *project) SetProxy(inst InstOperationInterface) {
	p.inst = inst
}

// CreateProject create project
func (p *project) CreateProject(kit *rest.Kit, data []mapstr.MapStr) ([]int64, error) {
	for idx, val := range data {
		if val[common.BKProjectIDField] == nil {
			data[idx][common.BKProjectIDField] = strings.Replace(uuid.New().String(), "-", "", -1)
		}
	}

	var err error
	input := &metadata.BatchCreateModelInstOption{Data: data}
	resp, err := p.clientSet.CoreService().Instance().BatchCreateInstance(kit.Ctx, kit.Header,
		common.BKInnerObjIDProject, input)
	if err != nil {
		blog.Errorf("create project failed, err: %v, data: %v, rid: %s", err, data, kit.Rid)
		return nil, err
	}

	if len(resp.IDs) != len(data) {
		blog.Errorf("the number of project creation is inconsistent, create success number: %d, data: %v, rid: %s",
			data, len(resp.IDs), kit.Rid)
		return nil, errors.New("the number of project creation is inconsistent")
	}

	for idx := range data {
		data[idx].Set(common.BKFieldID, resp.IDs[idx])
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewInstanceAudit(p.clientSet.CoreService())
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, common.BKInnerObjIDProject, data)
	if err != nil {
		blog.Errorf("generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	if auth.EnableAuthorize() {
		iamInstances := make([]metadata.IamInstance, len(data))
		for index := range data {
			iamInstances[index] = metadata.IamInstance{
				ID:   strconv.FormatInt(resp.IDs[index], 10),
				Name: util.GetStrByInterface(data[index][common.BKProjectNameField]),
			}
		}
		iamInstancesWithCreator := metadata.IamInstancesWithCreator{
			IamInstances: metadata.IamInstances{
				Type:      string(iam.Project),
				Instances: iamInstances,
			},
			Creator: kit.User,
		}
		_, err = p.authManager.Authorizer.BatchRegisterResourceCreatorAction(kit.Ctx, kit.Header, iamInstancesWithCreator)
		if err != nil {
			blog.Errorf("register created project to iam failed, err: %s, rid: %s", err, kit.Rid)
			return nil, err
		}
	}

	return resp.IDs, nil
}
