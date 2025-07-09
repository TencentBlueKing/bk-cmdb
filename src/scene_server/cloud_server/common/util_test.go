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

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCovertInstState(t *testing.T) {
	states := []string{"starting", "pending", "rebooting", "STARTING", "PENDING", "REBOOTING"}
	for _, state := range states {
		require.Equal(t, "starting", CovertInstState(state))
	}

	states = []string{"running", "RUNNING"}
	for _, state := range states {
		require.Equal(t, "running", CovertInstState(state))
	}

	states = []string{"stopping", "shutting-down", "terminating", "STOPPING", "SHUTTING-DOWN", "TERMINATING"}
	for _, state := range states {
		require.Equal(t, "stopping", CovertInstState(state))
	}

	states = []string{"stopped", "shutdown", "terminated", "STOPPED", "SHUTDOWN", "TERMINATED"}
	for _, state := range states {
		require.Equal(t, "stopped", CovertInstState(state))
	}

	states = []string{"fail", "create", "aaa"}
	for _, state := range states {
		require.Equal(t, "unknow", CovertInstState(state))
	}
}
