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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * set template
 */

func (am *AuthManager) CollectSetTemplatesByBusinessIDs(ctx context.Context, header http.Header, bizID int64) ([]metadata.SetTemplate, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	option := metadata.ListSetTemplateOption{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	templates, err := am.clientSet.CoreService().SetTemplate().ListSetTemplate(ctx, header, bizID, option)
	if err != nil {
		blog.Errorf("list set templates by bizID:%d failed, err: %+v, rid: %s", bizID, err, rid)
		return nil, fmt.Errorf("list set templates by bizID:%d failed, err: %+v", bizID, err)
	}

	return templates.Info, nil
}

func (am *AuthManager) collectSetTemplateByIDs(ctx context.Context, header http.Header, bizID int64, templateIDs ...int64) ([]metadata.SetTemplate, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	templateIDs = util.IntArrayUnique(templateIDs)
	option := metadata.ListSetTemplateOption{
		SetTemplateIDs: templateIDs,
	}
	result, err := am.clientSet.CoreService().SetTemplate().ListSetTemplate(ctx, header, bizID, option)
	if err != nil {
		blog.V(3).Infof("list set templates by id failed, templateIDs: %+v, err: %+v, rid: %s", templateIDs, err, rid)
		return nil, fmt.Errorf("list set templates by id failed, err: %+v", err)
	}

	return result.Info, nil
}

func (am *AuthManager) extractBizIDFromSetTemplate(templates ...metadata.SetTemplate) (int64, error) {
	var bizID int64
	for idx, template := range templates {
		if idx > 0 && bizID != template.BizID {
			return 0, fmt.Errorf("get multiple bk_biz_id from set templates")
		}
		bizID = template.BizID
	}
	return bizID, nil
}

func (am *AuthManager) MakeResourcesBySetTemplate(header http.Header, action meta.Action, bizID int64, templates ...metadata.SetTemplate) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, template := range templates {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.SetTemplate,
				Name:       template.Name,
				InstanceID: template.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      bizID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) MakeResourcesBySetTemplateIDs(ctx context.Context, header http.Header, action meta.Action, bizID int64, ids ...int64) ([]meta.ResourceAttribute, error) {
	templates, err := am.collectSetTemplateByIDs(ctx, header, bizID, ids...)
	if err != nil {
		return nil, fmt.Errorf("get set templates by id failed, err: %+v", err)
	}
	resources := am.MakeResourcesBySetTemplate(header, action, bizID, templates...)
	return resources, nil
}

func (am *AuthManager) AuthorizeBySetTemplateID(ctx context.Context, header http.Header, action meta.Action, bizID int64, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	templates, err := am.collectSetTemplateByIDs(ctx, header, bizID, ids...)
	if err != nil {
		return fmt.Errorf("get set templates by id failed, err: %+v", err)
	}
	return am.AuthorizeBySetTemplates(ctx, header, action, templates...)
}

func (am *AuthManager) GenSetTemplateNoPermissionResp() *metadata.BaseResp {
	permission := metadata.Permission{
		SystemID:      authcenter.SystemIDCMDB,
		SystemName:    authcenter.SystemNameCMDB,
		ScopeType:     authcenter.ScopeTypeIDSystem,
		ScopeTypeName: authcenter.ScopeTypeIDSystemName,
		ActionID:      string(authcenter.ModelTopologyOperation),
		ActionName:    authcenter.ActionIDNameMap[authcenter.ModelTopologyOperation],
		Resources: [][]metadata.Resource{
			{{
				ResourceType:     string(authcenter.SysSystemBase),
				ResourceTypeName: authcenter.ResourceTypeIDMap[authcenter.SysSystemBase],
			}},
		},
	}
	permission.ResourceType = permission.Resources[0][0].ResourceType
	permission.ResourceTypeName = permission.Resources[0][0].ResourceTypeName

	resp := metadata.NewNoPermissionResp([]metadata.Permission{permission})
	return &resp
}

func (am *AuthManager) AuthorizeBySetTemplates(ctx context.Context, header http.Header, action meta.Action, templates ...metadata.SetTemplate) error {
	if am.Enabled() == false {
		return nil
	}

	if len(templates) == 0 {
		return nil
	}
	// extract business id
	bizID, err := am.extractBizIDFromSetTemplate(templates...)
	if err != nil {
		return fmt.Errorf("authorize set templates failed, extract business id from set templates failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesBySetTemplate(header, action, bizID, templates...)

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) UpdateRegisteredSetTemplates(ctx context.Context, header http.Header, templates ...metadata.SetTemplate) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if am.Enabled() == false {
		return nil
	}

	if len(templates) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBizIDFromSetTemplate(templates...)
	if err != nil {
		return fmt.Errorf("authorize set templates failed, extract business id from set template failed, err: %+v, rid: %s", err, rid)
	}

	// make auth resources
	resources := am.MakeResourcesBySetTemplate(header, meta.EmptyAction, bizID, templates...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredSetTemplateByID(ctx context.Context, header http.Header, bizID int64, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	templates, err := am.collectSetTemplateByIDs(ctx, header, bizID, ids...)
	if err != nil {
		return fmt.Errorf("update registered set templates failed, get set template by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredSetTemplates(ctx, header, templates...)
}

func (am *AuthManager) DeregisterSetTemplateByIDs(ctx context.Context, header http.Header, bizID int64, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	templates, err := am.collectSetTemplateByIDs(ctx, header, bizID, ids...)
	if err != nil {
		return fmt.Errorf("deregister set templates failed, get set templates by id failed, err: %+v", err)
	}
	return am.DeregisterSetTemplates(ctx, header, templates...)
}

func (am *AuthManager) RegisterSetTemplates(ctx context.Context, header http.Header, templates ...metadata.SetTemplate) error {
	if am.Enabled() == false {
		return nil
	}

	if len(templates) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBizIDFromSetTemplate(templates...)
	if err != nil {
		return fmt.Errorf("register set templates failed, extract business id from set templates failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesBySetTemplate(header, meta.EmptyAction, bizID, templates...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterSetTemplateByID(ctx context.Context, header http.Header, bizID int64, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	templates, err := am.collectSetTemplateByIDs(ctx, header, bizID, ids...)
	if err != nil {
		return fmt.Errorf("register set template failed, get set template by id failed, err: %+v", err)
	}
	return am.RegisterSetTemplates(ctx, header, templates...)
}

func (am *AuthManager) DeregisterSetTemplates(ctx context.Context, header http.Header, templates ...metadata.SetTemplate) error {
	if am.Enabled() == false {
		return nil
	}

	if len(templates) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBizIDFromSetTemplate(templates...)
	if err != nil {
		return fmt.Errorf("deregister set template failed, extract business id from set template failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesBySetTemplate(header, meta.EmptyAction, bizID, templates...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) ListAuthorizedSetTemplateIDs(ctx context.Context, header http.Header, bizID int64) ([]int64, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	listOption := &meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:   meta.SetTemplate,
			Action: meta.FindMany,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      bizID,
	}
	resources, err := am.Authorize.ListResources(ctx, listOption)
	if err != nil {
		blog.Errorf("list authorized set template from iam failed, err: %+v, rid: %s", err, rid)
		return nil, err
	}
	ids := make([]int64, 0)
	for _, item := range resources {
		for _, resource := range item {
			id, err := strconv.ParseInt(resource.ResourceID, 10, 64)
			if err != nil {
				blog.Errorf("list authorized set template from iam failed, err: %+v, rid: %s", err, rid)
				return nil, fmt.Errorf("parse resource id into int64 failed, err: %+v", err)
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}
