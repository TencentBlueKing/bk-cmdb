/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extensions

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * service template
 */

func (am *AuthManager) CollectServiceTemplatesByBusinessIDs(ctx context.Context, header http.Header, businessID int64) ([]metadata.ServiceTemplate, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	option := &metadata.ListServiceTemplateOption{
		BusinessID: businessID,
	}
	templates, err := am.clientSet.CoreService().Process().ListServiceTemplates(ctx, header, option)
	if err != nil {
		blog.Errorf("list service templates by businessID:%d failed, err: %+v, rid: %s", businessID, err, rid)
		return nil, fmt.Errorf("list service templates by businessID:%d failed, err: %+v", businessID, err)
	}

	return templates.Info, nil
}

func (am *AuthManager) collectServiceTemplateByIDs(ctx context.Context, header http.Header, templateIDs ...int64) ([]metadata.ServiceTemplate, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	templateIDs = util.IntArrayUnique(templateIDs)
	option := &metadata.ListServiceTemplateOption{
		ServiceTemplateIDs: templateIDs,
	}
	result, err := am.clientSet.CoreService().Process().ListServiceTemplates(ctx, header, option)
	if err != nil {
		blog.V(3).Infof("list service templates by id failed, templateIDs: %+v, err: %+v, rid: %s", templateIDs, err, rid)
		return nil, fmt.Errorf("list service templates by id failed, err: %+v", err)
	}

	return result.Info, nil
}

func (am *AuthManager) extractBusinessIDFromServiceTemplate(templates ...metadata.ServiceTemplate) (int64, error) {
	var businessID int64
	for idx, template := range templates {
		bizID := template.BizID
		// we should ignore metadata.LabelBusinessID field not found error
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("get multiple business ID from service templates")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) MakeResourcesByServiceTemplate(header http.Header, action meta.Action, businessID int64, templates ...metadata.ServiceTemplate) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, template := range templates {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ProcessServiceTemplate,
				Name:       template.Name,
				InstanceID: template.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) MakeResourcesByServiceTemplateIDs(ctx context.Context, header http.Header, action meta.Action, businessID int64, ids ...int64) ([]meta.ResourceAttribute, error) {
	templates, err := am.collectServiceTemplateByIDs(ctx, header, ids...)
	if err != nil {
		return nil, fmt.Errorf("get service templates by id failed, err: %+v", err)
	}
	resources := am.MakeResourcesByServiceTemplate(header, action, businessID, templates...)
	return resources, nil
}

func (am *AuthManager) AuthorizeByServiceTemplateID(ctx context.Context, header http.Header, action meta.Action, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	templates, err := am.collectServiceTemplateByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("get service templates by id failed, err: %+v", err)
	}
	return am.AuthorizeByServiceTemplates(ctx, header, action, templates...)
}

func (am *AuthManager) GenServiceTemplateNoPermissionResp() *metadata.BaseResp {
	var p metadata.Permission
	p.SystemID = authcenter.SystemIDCMDB
	p.SystemName = authcenter.SystemNameCMDB
	p.ScopeType = authcenter.ScopeTypeIDSystem
	p.ScopeTypeName = authcenter.ScopeTypeIDSystemName
	p.ActionID = string(authcenter.ModelTopologyOperation)
	p.ActionName = authcenter.ActionIDNameMap[authcenter.ModelTopologyOperation]
	p.Resources = [][]metadata.Resource{
		{{
			ResourceType:     string(authcenter.SysSystemBase),
			ResourceTypeName: authcenter.ResourceTypeIDMap[authcenter.SysSystemBase],
		}},
	}
	p.ResourceType = p.Resources[0][0].ResourceType
	p.ResourceTypeName = p.Resources[0][0].ResourceTypeName

	resp := metadata.NewNoPermissionResp([]metadata.Permission{p})
	return &resp
}

func (am *AuthManager) AuthorizeByServiceTemplates(ctx context.Context, header http.Header, action meta.Action, templates ...metadata.ServiceTemplate) error {
	if am.Enabled() == false {
		return nil
	}

	if len(templates) == 0 {
		return nil
	}
	// extract business id
	bizID, err := am.extractBusinessIDFromServiceTemplate(templates...)
	if err != nil {
		return fmt.Errorf("authorize service templates failed, extract business id from service templates failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByServiceTemplate(header, action, bizID, templates...)

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) UpdateRegisteredServiceTemplates(ctx context.Context, header http.Header, templates ...metadata.ServiceTemplate) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if am.Enabled() == false {
		return nil
	}

	if len(templates) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromServiceTemplate(templates...)
	if err != nil {
		return fmt.Errorf("authorize service templates failed, extract business id from service template failed, err: %+v, rid: %s", err, rid)
	}

	// make auth resources
	resources := am.MakeResourcesByServiceTemplate(header, meta.EmptyAction, bizID, templates...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredServiceTemplateByID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	templates, err := am.collectServiceTemplateByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered service templates failed, get service template by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredServiceTemplates(ctx, header, templates...)
}

func (am *AuthManager) DeregisterServiceTemplateByIDs(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	templates, err := am.collectServiceTemplateByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister service templates failed, get service templates by id failed, err: %+v", err)
	}
	return am.DeregisterServiceTemplates(ctx, header, templates...)
}

func (am *AuthManager) RegisterServiceTemplates(ctx context.Context, header http.Header, templates ...metadata.ServiceTemplate) error {
	if am.Enabled() == false {
		return nil
	}

	if len(templates) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromServiceTemplate(templates...)
	if err != nil {
		return fmt.Errorf("register service templates failed, extract business id from service templates failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByServiceTemplate(header, meta.EmptyAction, bizID, templates...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterServiceTemplateByID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	templates, err := am.collectServiceTemplateByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("register service template failed, get service template by id failed, err: %+v", err)
	}
	return am.RegisterServiceTemplates(ctx, header, templates...)
}

func (am *AuthManager) DeregisterServiceTemplates(ctx context.Context, header http.Header, templates ...metadata.ServiceTemplate) error {
	if am.Enabled() == false {
		return nil
	}

	if len(templates) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromServiceTemplate(templates...)
	if err != nil {
		return fmt.Errorf("deregister service template failed, extract business id from service template failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByServiceTemplate(header, meta.EmptyAction, bizID, templates...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) ListAuthorizedServiceTemplateIDs(ctx context.Context, header http.Header, bizID int64) ([]int64, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	listOption := &meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:   meta.ProcessServiceTemplate,
			Action: meta.FindMany,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      bizID,
	}
	resources, err := am.Authorize.ListResources(ctx, listOption)
	if err != nil {
		blog.Errorf("list authorized service template from iam failed, err: %+v, rid: %s", err, rid)
		return nil, err
	}
	ids := make([]int64, 0)
	for _, item := range resources {
		for _, resource := range item {
			id, err := strconv.ParseInt(resource.ResourceID, 10, 64)
			if err != nil {
				blog.Errorf("list authorized service template from iam failed, err: %+v, rid: %s", err, rid)
				return nil, fmt.Errorf("parse resource id into int64 failed, err: %+v", err)
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}
