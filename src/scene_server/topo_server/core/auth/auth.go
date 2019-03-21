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

package topoauth

import (
	"context"
	"fmt"
	"net/http"

	apiutil "configcenter/src/apimachinery/util"
	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type TopoAuth struct {
	EnableAuth bool
	Authorizer auth.Authorize
}

// initialize a topology auth instance for auth usage.
// can be used for resources authorize and resources register/deregister management.
func NewTopologyAuth(tls *apiutil.TLSClientConfig, config authcenter.AuthConfig) (*TopoAuth, error) {
	authAuthorizer, err := auth.NewAuthorize(tls, config)
	if err != nil {
		blog.Errorf("new topo server authorizer failed, err: %+v", err)
		return nil, fmt.Errorf("new topo server authorizer failed, err: %+v", err)
	}

	return &TopoAuth{
		EnableAuth: config.Enable,
		Authorizer: authAuthorizer,
	}, nil
}

func (ta *TopoAuth) UpdateRegisteredObject(ctx context.Context, header http.Header, object *metadata.Object) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.Model,
			Name:       object.ObjectName,
			InstanceID: object.ID,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if _, exist := object.Metadata.Label[metadata.LabelBusinessID]; exist {
		bizID, err := object.Metadata.Label.Int64(metadata.LabelBusinessID)
		if err != nil {
			return err
		}
		resource.BusinessID = bizID
	}

	if err := ta.Authorizer.UpdateResource(ctx, &resource); err != nil {
		return err
	}

	return nil
}

func (ta *TopoAuth) DeregisterObject(ctx context.Context, header http.Header, object *metadata.Object) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.Model,
			Name:       object.ObjectName,
			InstanceID: object.ID,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if _, exist := object.Metadata.Label[metadata.LabelBusinessID]; exist {
		bizID, err := object.Metadata.Label.Int64(metadata.LabelBusinessID)
		if err != nil {
			return err
		}
		resource.BusinessID = bizID
	}

	if err := ta.Authorizer.DeregisterResource(ctx, resource); err != nil {
		return err
	}
	return nil
}

/*
func (ta *TopoAuth) AuthorizeObject(ctx context.Context, header http.Header, businessID int64, action meta.Action, objects ...metadata.Object) error {
	resources := make([]meta.ResourceAttribute, 0)
	for _, object := range objects {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Model,
				Name:       object.ObjectName,
				InstanceID: object.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID: businessID,
		}

		resources = append(resources, resource)
	}

	return ta.authorize(ctx, header, businessID, resources...)
}
*/


func (ta *TopoAuth) RegisterClassification(ctx context.Context, header http.Header, class *metadata.Classification) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.ModelClassification,
			Name:       class.ClassificationName,
			InstanceID: class.ID,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if _, exist := class.Metadata.Label[metadata.LabelBusinessID]; exist {
		bizID, err := class.Metadata.Label.Int64(metadata.LabelBusinessID)
		if err != nil {
			return err
		}
		resource.BusinessID = bizID
	}

	if err := ta.Authorizer.RegisterResource(ctx, resource); err != nil {
		return err
	}

	return nil
}

func (ta *TopoAuth) RegisterAssociationType(ctx context.Context, header http.Header, kind *metadata.AssociationKind) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.AssociationType,
			Name:       kind.AssociationKindID,
			InstanceID: kind.ID,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if err := ta.Authorizer.RegisterResource(ctx, resource); err != nil {
		return err
	}

	return nil
}

func (ta *TopoAuth) UpdateAssociationType(ctx context.Context, header http.Header, kindID int64, kindName string) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.AssociationType,
			Name:       kindName,
			InstanceID: kindID,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if err := ta.Authorizer.UpdateResource(ctx, &resource); err != nil {
		return err
	}

	return nil
}

func (ta *TopoAuth) DeregisterAssociationType(ctx context.Context, header http.Header, assoTypeID int64) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.AssociationType,
			InstanceID: assoTypeID,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if err := ta.Authorizer.DeregisterResource(ctx, resource); err != nil {
		return err
	}
	return nil
}

func (ta *TopoAuth) RegisterBusiness(ctx context.Context, header http.Header, bizName string, id int64) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.Business,
			Name:       bizName,
			InstanceID: id,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if err := ta.Authorizer.RegisterResource(ctx, resource); err != nil {
		return err
	}
	return nil
}

func (ta *TopoAuth) UpdateBusiness(ctx context.Context, header http.Header, bizName string, id int64) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.Business,
			Name:       bizName,
			InstanceID: id,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if err := ta.Authorizer.UpdateResource(ctx, &resource); err != nil {
		return err
	}

	return nil
}

func (ta *TopoAuth) DeregisterBusiness(ctx context.Context, header http.Header, id int64) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.Business,
			InstanceID: id,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if err := ta.Authorizer.DeregisterResource(ctx, resource); err != nil {
		return err
	}

	return nil
}

func (ta *TopoAuth) RegisterInstance() error {

	return nil
}

func (ta *TopoAuth) DeregisterInstance() error {

	return nil
}
