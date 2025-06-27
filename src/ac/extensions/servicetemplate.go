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

	"configcenter/src/ac/meta"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * service template
 */

func (a *AuthManager) collectServiceTemplateByIDs(ctx context.Context, header http.Header,
	templateIDs ...int64) ([]metadata.ServiceTemplate, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	templateIDs = util.IntArrayUnique(templateIDs)
	option := &metadata.ListServiceTemplateOption{
		ServiceTemplateIDs: templateIDs,
	}
	result, err := a.clientSet.CoreService().Process().ListServiceTemplates(ctx, header, option)
	if err != nil {
		blog.V(3).Infof("list service templates by id failed, templateIDs: %+v, err: %+v, rid: %s", templateIDs, err,
			rid)
		return nil, fmt.Errorf("list service templates by id failed, err: %+v", err)
	}

	return result.Info, nil
}

func (a *AuthManager) extractBusinessIDFromServiceTemplate(templates ...metadata.ServiceTemplate) (int64, error) {
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

// MakeResourcesByServiceTemplate TODO
func (a *AuthManager) MakeResourcesByServiceTemplate(header http.Header, action meta.Action, businessID int64,
	templates ...metadata.ServiceTemplate) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, template := range templates {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ProcessServiceTemplate,
				Name:       template.Name,
				InstanceID: template.ID,
			},
			TenantID:   httpheader.GetTenantID(header),
			BusinessID: businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

// AuthorizeByServiceTemplateID TODO
func (a *AuthManager) AuthorizeByServiceTemplateID(kit *rest.Kit, action meta.Action, ids ...int64) (
	*metadata.BaseResp, bool, error) {
	if !a.Enabled() {
		return nil, true, nil
	}

	if len(ids) == 0 {
		return nil, true, nil
	}

	templates, err := a.collectServiceTemplateByIDs(kit.Ctx, kit.Header, ids...)
	if err != nil {
		return nil, true, fmt.Errorf("get service templates by id failed, err: %+v", err)
	}
	return a.AuthorizeByServiceTemplates(kit, action, templates...)
}

// AuthorizeByServiceTemplates TODO
func (a *AuthManager) AuthorizeByServiceTemplates(kit *rest.Kit, action meta.Action,
	templates ...metadata.ServiceTemplate) (*metadata.BaseResp, bool, error) {

	if !a.Enabled() {
		return nil, true, nil
	}

	if len(templates) == 0 {
		return nil, true, nil
	}
	// extract business id
	bizID, err := a.extractBusinessIDFromServiceTemplate(templates...)
	if err != nil {
		return nil, true, fmt.Errorf("extract business id from service templates failed, err: %v", err)
	}

	// make auth resources
	resources := a.MakeResourcesByServiceTemplate(kit.Header, action, bizID, templates...)
	authResp, authorized := a.Authorize(kit, resources...)
	return authResp, authorized, nil
}
