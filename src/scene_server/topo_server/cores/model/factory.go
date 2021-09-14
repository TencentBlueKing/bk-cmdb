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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// New create a new model factory instance
func New(clientSet apimachinery.ClientSetInterface, languageIf language.CCLanguageIf) Factory {
	return &factory{
		clientSet: clientSet,
		language:  languageIf,
	}
}

// CreateClassification create classification objects
func CreateClassification(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, clsItems []metadata.Classification) []Classification {
	results := make([]Classification, 0)
	for _, cls := range clsItems {
		results = append(results, &classification{
			cls:       cls,
			kit:       kit,
			clientSet: clientSet,
		})
	}

	return results
}

// CreateObject create  objects
func CreateObject(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, objItems []metadata.Object) []Object {
	results := make([]Object, 0)
	for _, obj := range objItems {
		results = append(results, &object{
			obj:       obj,
			kit:       kit,
			clientSet: clientSet,
		})
	}

	return results
}

// CreateGroup create group  objects
func CreateGroup(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, groupItems []metadata.Group) []GroupInterface {
	results := make([]GroupInterface, 0)
	for _, grp := range groupItems {

		results = append(results, &group{
			grp:       grp,
			kit:       kit,
			clientSet: clientSet,
		})
	}

	return results
}

// CreateAttribute create attribute  objects
func CreateAttribute(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, attrItems []metadata.Attribute) []AttributeInterface {
	results := make([]AttributeInterface, 0)
	for _, attr := range attrItems {
		results = append(results, &attribute{
			attr:      attr,
			kit:       kit,
			clientSet: clientSet,
		})

	}

	return results
}

type factory struct {
	clientSet apimachinery.ClientSetInterface
	language  language.CCLanguageIf
}

func (cli *factory) CreateObject(kit *rest.Kit) Object {
	obj := &object{
		FieldValid: FieldValid{
			lang: cli.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)),
		},
		kit:       kit,
		clientSet: cli.clientSet,
	}
	obj.SetSupplierAccount(kit.SupplierAccount)
	return obj
}

func (cli *factory) CreateClassification(kit *rest.Kit) Classification {
	cls := &classification{
		kit:       kit,
		clientSet: cli.clientSet,
	}
	cls.SetSupplierAccount(kit.SupplierAccount)
	return cls
}

func (cli *factory) CreateAttribute(kit *rest.Kit) AttributeInterface {
	attr := &attribute{
		FieldValid: FieldValid{
			lang: cli.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)),
		},
		kit:       kit,
		clientSet: cli.clientSet,
	}
	attr.SetSupplierAccount(kit.SupplierAccount)
	return attr
}

func (cli *factory) CreateGroup(kit *rest.Kit, bizID int64) GroupInterface {
	return NewGroup(kit, cli.clientSet, bizID)
}

func (cli *factory) CreateMainLineAssociation(kit *rest.Kit, obj Object, asstKey string, asstObj Object) Association {
	asst := &association{
		isMainLine: true,
		kit:        kit,
		clientSet:  cli.clientSet,
	}
	asst.SetSupplierAccount(kit.SupplierAccount)
	return asst
}
func (cli *factory) CreateCommonAssociation(kit *rest.Kit, obj Object, asstKey string, asstObj Object) Association {
	asst := &association{
		isMainLine: false,
		kit:        kit,
		clientSet:  cli.clientSet,
	}
	asst.SetSupplierAccount(kit.SupplierAccount)
	return asst
}
