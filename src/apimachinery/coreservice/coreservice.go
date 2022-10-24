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

// Package coreservice TODO
package coreservice

import (
	"fmt"

	"configcenter/src/apimachinery/coreservice/association"
	"configcenter/src/apimachinery/coreservice/auditlog"
	"configcenter/src/apimachinery/coreservice/auth"
	"configcenter/src/apimachinery/coreservice/cloud"
	"configcenter/src/apimachinery/coreservice/common"
	"configcenter/src/apimachinery/coreservice/count"
	"configcenter/src/apimachinery/coreservice/host"
	"configcenter/src/apimachinery/coreservice/hostapplyrule"
	"configcenter/src/apimachinery/coreservice/instance"
	"configcenter/src/apimachinery/coreservice/label"
	"configcenter/src/apimachinery/coreservice/mainline"
	"configcenter/src/apimachinery/coreservice/model"
	"configcenter/src/apimachinery/coreservice/operation"
	"configcenter/src/apimachinery/coreservice/process"
	"configcenter/src/apimachinery/coreservice/settemplate"
	"configcenter/src/apimachinery/coreservice/synchronize"
	ccSystem "configcenter/src/apimachinery/coreservice/system"
	"configcenter/src/apimachinery/coreservice/topographics"
	"configcenter/src/apimachinery/coreservice/transaction"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

// CoreServiceClientInterface TODO
type CoreServiceClientInterface interface {
	Instance() instance.InstanceClientInterface
	Model() model.ModelClientInterface
	Association() association.AssociationClientInterface
	Synchronize() synchronize.SynchronizeClientInterface
	Mainline() mainline.MainlineClientInterface
	Host() host.HostClientInterface
	Audit() auditlog.AuditClientInterface
	Process() process.ProcessInterface
	Operation() operation.OperationClientInterface
	Label() label.LabelInterface
	TopoGraphics() topographics.TopoGraphicsInterface
	SetTemplate() settemplate.SetTemplateInterface
	HostApplyRule() hostapplyrule.HostApplyRuleInterface
	System() ccSystem.SystemClientInterface
	Txn() transaction.Interface
	Count() count.CountClientInterface
	Cloud() cloud.CloudInterface
	Auth() auth.AuthClientInterface
	Common() common.CommonInterface
}

// NewCoreServiceClient TODO
func NewCoreServiceClient(c *util.Capability, version string) CoreServiceClientInterface {
	base := fmt.Sprintf("/api/%s", version)
	return &coreService{
		restCli: rest.NewRESTClient(c, base),
	}
}

type coreService struct {
	restCli rest.ClientInterface
}

// Instance TODO
func (c *coreService) Instance() instance.InstanceClientInterface {
	return instance.NewInstanceClientInterface(c.restCli)
}

// Model TODO
func (c *coreService) Model() model.ModelClientInterface {
	return model.NewModelClientInterface(c.restCli)
}

// Association TODO
func (c *coreService) Association() association.AssociationClientInterface {
	return association.NewAssociationClientInterface(c.restCli)
}

// Mainline TODO
func (c *coreService) Mainline() mainline.MainlineClientInterface {
	return mainline.NewMainlineClientInterface(c.restCli)
}

// Synchronize TODO
func (c *coreService) Synchronize() synchronize.SynchronizeClientInterface {
	return synchronize.NewSynchronizeClientInterface(c.restCli)
}

// Host TODO
func (c *coreService) Host() host.HostClientInterface {
	return host.NewHostClientInterface(c.restCli)
}

// Audit TODO
func (c *coreService) Audit() auditlog.AuditClientInterface {
	return auditlog.NewAuditClientInterface(c.restCli)
}

// Process TODO
func (c *coreService) Process() process.ProcessInterface {
	return process.NewProcessInterfaceClient(c.restCli)

}

// Operation TODO
func (c *coreService) Operation() operation.OperationClientInterface {
	return operation.NewOperationClientInterface(c.restCli)
}

// Label TODO
func (c *coreService) Label() label.LabelInterface {
	return label.NewLabelInterfaceClient(c.restCli)
}

// TopoGraphics TODO
func (c *coreService) TopoGraphics() topographics.TopoGraphicsInterface {
	return topographics.NewTopoGraphicsInterface(c.restCli)
}

// System TODO
func (c *coreService) System() ccSystem.SystemClientInterface {
	return ccSystem.NewSystemClientInterface(c.restCli)
}

// SetTemplate TODO
func (c *coreService) SetTemplate() settemplate.SetTemplateInterface {
	return settemplate.NewSetTemplateInterfaceClient(c.restCli)
}

// HostApplyRule TODO
func (c *coreService) HostApplyRule() hostapplyrule.HostApplyRuleInterface {
	return hostapplyrule.NewHostApplyRuleClient(c.restCli)
}

// Txn TODO
func (c *coreService) Txn() transaction.Interface {
	return transaction.NewTxn(c.restCli)
}

// Count TODO
func (c *coreService) Count() count.CountClientInterface {
	return count.NewCountClientInterface(c.restCli)
}

// Cloud TODO
func (c *coreService) Cloud() cloud.CloudInterface {
	return cloud.NewCloudInterfaceClient(c.restCli)
}

// Auth TODO
func (c *coreService) Auth() auth.AuthClientInterface {
	return auth.NewAuthClientInterface(c.restCli)
}

// Common TODO
func (c *coreService) Common() common.CommonInterface {
	return common.NewCommonInterfaceClient(c.restCli)
}
