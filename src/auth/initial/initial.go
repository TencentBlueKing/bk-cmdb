/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package initial

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
)

func Init(ctx context.Context, cli authcenter.AuthCenter) error {
	detail := authcenter.SystemDetail{}
	detail.System = expectSystem
	detail.Scopes = append(detail.Scopes, struct {
		ScopeTypeID   string                    `json:"scope_type_id"`
		ResourceTypes []authcenter.ResourceType `json:"resource_types"`
	}{
		ScopeTypeID:   "system",
		ResourceTypes: expectResourceType,
	})

	info, err := cli.QuerySystemInfo(ctx, meta.SystemIDCMDB, false)
	if err != nil {
		return err
	}

	if info.SystemID != meta.SystemIDCMDB {
		if err := cli.RegistSystem(ctx, expectSystem); err != nil {
			return err
		}
	}

	return nil
}

func resourceKey(res authcenter.ResourceType) string {
	return fmt.Sprintf("%s-%s-%s", res.ResourceTypeID, res.ResourceTypeName, res.ParentResourceTypeID)
}

func actionKey(actions []authcenter.Action) string {
	sort.Slice(actions, func(i, j int) bool {
		if actions[i].ActionID < actions[j].ActionID {
			return true
		} else if actions[i].ActionName < actions[j].ActionName {
			return true
		} else {
			return actions[i].IsRelatedResource == actions[j].IsRelatedResource
		}
	})

	keys := []string{}
	for _, action := range actions {
		keys = append(keys, fmt.Sprintf("%s-%s-%v", action.ActionID, action.ActionName, action.IsRelatedResource))
	}
	return strings.Join(keys, ":")
}
