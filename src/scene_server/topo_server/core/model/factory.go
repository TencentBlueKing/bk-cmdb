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

package model

import (
	"configcenter/src/apimachinery"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// New create a new model factory instance
func New(clientSet apimachinery.ClientSetInterface) Factory {
	return &factory{
		clientSet: clientSet,
	}
}

// CreateClassification create classification objects
func CreateClassification(params types.ContextParams, clientSet apimachinery.ClientSetInterface, clsItems []metadata.Classification) []Classification {
	results := make([]Classification, 0)
	for _, cls := range clsItems {

		results = append(results, &classification{
			cls:       cls,
			params:    params,
			clientSet: clientSet,
		})
	}

	return results
}

// CreateObject create  objects
func CreateObject(params types.ContextParams, clientSet apimachinery.ClientSetInterface, objItems []metadata.Object) []Object {
	results := make([]Object, 0)
	for _, obj := range objItems {

		results = append(results, &object{
			obj:       obj,
			params:    params,
			clientSet: clientSet,
		})
	}

	return results
}

// CreateGroup create group  objects
func CreateGroup(params types.ContextParams, clientSet apimachinery.ClientSetInterface, groupItems []metadata.Group) []Group {
	results := make([]Group, 0)
	for _, grp := range groupItems {

		results = append(results, &group{
			grp:       grp,
			params:    params,
			clientSet: clientSet,
		})
	}

	return results
}

// CreateAttribute create attribute  objects
func CreateAttribute(params types.ContextParams, clientSet apimachinery.ClientSetInterface, attrItems []metadata.Attribute) []Attribute {
	results := make([]Attribute, 0)
	for _, attr := range attrItems {

		results = append(results, &attribute{
			attr:      attr,
			params:    params,
			clientSet: clientSet,
		})
	}

	return results
}

type factory struct {
	clientSet apimachinery.ClientSetInterface
}

func (cli *factory) CreaetObject(params types.ContextParams) Object {
	obj := &object{
		params:    params,
		clientSet: cli.clientSet,
	}
	obj.SetSupplierAccount(params.SupplierAccount)
	return obj
}

func (cli *factory) CreaetClassification(params types.ContextParams) Classification {
	cls := &classification{
		params:    params,
		clientSet: cli.clientSet,
	}
	cls.SetSupplierAccount(params.SupplierAccount)
	return cls
}

func (cli *factory) CreateAttribute(params types.ContextParams) Attribute {
	attr := &attribute{
		params:    params,
		clientSet: cli.clientSet,
	}
	attr.SetSupplierAccount(params.SupplierAccount)
	return attr
}

func (cli *factory) CreateGroup(params types.ContextParams) Group {
	return &group{
		params:    params,
		clientSet: cli.clientSet,
	}
}

func (cli *factory) CreateMainLineAssociatin(params types.ContextParams, obj Object, asstKey string, asstObj Object) Association {
	return &association{
		isMainLine: true,
		params:     params,
		clientSet:  cli.clientSet,
	}
}
func (cli *factory) CreateCommonAssociation(params types.ContextParams, obj Object, asstKey string, asstObj Object) Association {

	return &association{
		isMainLine: false,
		params:     params,
		clientSet:  cli.clientSet,
	}
}
