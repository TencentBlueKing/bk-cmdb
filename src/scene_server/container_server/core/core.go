/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/container_server/core/pod"
)

// PodOperation pod methods
type PodOperation interface {
	CreatePod(kit *rest.Kit, inputParam metadata.CreatePod) (*metadata.CreatePodResult, error)
	CreateManyPod(kit *rest.Kit, inputParam metadata.CreateManyPod) (*metadata.CreateManyPodResult, error)
	UpdatePod(kit *rest.Kit, podID string, inputParam metadata.UpdatePod) (*metadata.UpdatePodResult, error)
	UpdateManyPod(kit *rest.Kit, inputParam metadata.UpdateManyPod) (*metadata.UpdateManyPodResult, error)
	DeletePod(kit *rest.Kit, podID string) (*metadata.DeletePodResult, error)
	DeleteManyPod(kit *rest.Kit, inputParam metadata.DeleteManyPod) (*metadata.DeleteManyPodResult, error)
	ListPod(kit *rest.Kit, inputParam metadata.ListPod) (*metadata.ListPodResult, error)
}

// Interface core interfaces methods
type Interface interface {
	PodOperation() PodOperation
}

type core struct {
	podOp PodOperation
}

// New create core
func New(
	client apimachinery.ClientSetInterface,
	languageIf language.CCLanguageIf,
) Interface {

	podOp := pod.New(client, languageIf)

	return &core{
		podOp: podOp,
	}
}

// PodOperation return pod operation interface
func (m *core) PodOperation() PodOperation {
	return m.podOp
}
