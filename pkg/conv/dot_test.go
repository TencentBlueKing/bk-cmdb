/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package conv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDot(t *testing.T) {
	str := "test1.test2.test3.test4"
	encodedStr := EncodeDot(str)
	require.Equal(t, encodedStr, "test1\\u002etest2\\u002etest3\\u002etest4")

	decodedStr := DecodeDot(encodedStr)
	require.Equal(t, decodedStr, str)
}

func TestDotEdgeCases(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"no_dot", "nodot"},
		{"k8s_label", "app.kubernetes.io/name"},
		{"docker_image", "docker.io/library/nginx:1.21.0"},
		{"semver", "v1.2.3-beta.1"},
		{"only_dots", "..."},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			encoded := EncodeDot(c.input)
			decoded := DecodeDot(encoded)
			require.Equal(t, c.input, decoded)
		})
	}
}
