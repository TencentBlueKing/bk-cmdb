/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package types

import (
	"testing"

	"configcenter/src/common/json"

	"go.mongodb.org/mongo-driver/bson"
)

func TestResourceListJSONRoundTrip(t *testing.T) {
	input := []byte(`{"cpu":"500m","memory":"128Mi","nvidia.com/gpu":"1","storage":"1.5Gi"}`)

	var resources ResourceList
	if err := json.Unmarshal(input, &resources); err != nil {
		t.Fatalf("unmarshal resource list failed: %v", err)
	}

	if got := resources[ResourceName("cpu")].String(); got != "500m" {
		t.Fatalf("unexpected cpu quantity: %s", got)
	}

	output, err := json.Marshal(resources)
	if err != nil {
		t.Fatalf("marshal resource list failed: %v", err)
	}

	const expected = `{"cpu":"500m","memory":"128Mi","nvidia.com/gpu":"1","storage":"1.5Gi"}`
	if string(output) != expected {
		t.Fatalf("unexpected resource list JSON: %s", output)
	}
}

func TestQuantityRejectsNonStringJSON(t *testing.T) {
	for _, input := range []string{`128`, `{}`, `[]`, `true`} {
		var quantity Quantity
		if err := json.Unmarshal([]byte(input), &quantity); err == nil {
			t.Fatalf("expected non-string quantity %s to be rejected", input)
		}
	}
}

func TestResourceListBSONRoundTrip(t *testing.T) {
	var resources ResourceList
	if err := json.Unmarshal([]byte(`{"cpu":"250m","memory":"64Mi","nvidia.com/gpu":"2"}`), &resources); err != nil {
		t.Fatalf("unmarshal resource list failed: %v", err)
	}

	input := struct {
		Limits ResourceList `bson:"limits"`
	}{
		Limits: resources,
	}

	data, err := bson.Marshal(input)
	if err != nil {
		t.Fatalf("marshal resource list to BSON failed: %v", err)
	}

	var stored bson.M
	if err := bson.Unmarshal(data, &stored); err != nil {
		t.Fatalf("inspect stored BSON failed: %v", err)
	}
	limits, ok := stored["limits"].(bson.M)
	if !ok {
		t.Fatalf("unexpected limits BSON type: %T", stored["limits"])
	}
	if got := limits["cpu"]; got != "250m" {
		t.Fatalf("unexpected stored cpu quantity: %#v", got)
	}
	if got := limits["memory"]; got != "64Mi" {
		t.Fatalf("unexpected stored memory quantity: %#v", got)
	}
	if got := limits["nvidia.com/gpu"]; got != "2" {
		t.Fatalf("unexpected stored gpu quantity: %#v", got)
	}

	var output struct {
		Limits ResourceList `bson:"limits"`
	}
	if err := bson.Unmarshal(data, &output); err != nil {
		t.Fatalf("unmarshal resource list from BSON failed: %v", err)
	}
	if got := output.Limits[ResourceName("cpu")].String(); got != "250m" {
		t.Fatalf("unexpected decoded cpu quantity: %s", got)
	}
}

func TestContainerResourceQuantityJSON(t *testing.T) {
	input := []byte(`{
		"name":"nginx",
		"limits":{"cpu":"500m","memory":"128Mi","nvidia.com/gpu":"1"},
		"requests":{"cpu":"250m","memory":"64Mi","nvidia.com/gpu":"1"}
	}`)

	var container Container
	if err := json.Unmarshal(input, &container); err != nil {
		t.Fatalf("unmarshal container failed: %v", err)
	}

	if container.Limits == nil || (*container.Limits)[ResourceName("memory")].String() != "128Mi" {
		t.Fatalf("unexpected container limits: %#v", container.Limits)
	}
	if container.ReqSysSpecuests == nil ||
		(*container.ReqSysSpecuests)[ResourceName("cpu")].String() != "250m" {
		t.Fatalf("unexpected container requests: %#v", container.ReqSysSpecuests)
	}
	if (*container.Limits)[ResourceName("nvidia.com/gpu")].String() != "1" {
		t.Fatalf("unexpected container gpu limit: %#v", container.Limits)
	}
}
