/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package options

import (
	"errors"
	"fmt"
	"sort"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"
)

// Config is the transfer service config
type Config struct {
	Sync       *SyncConfig
	DestExConf *DestExSyncConf
	Mongo      mongo.Config
	WatchMongo mongo.Config
	Redis      redis.Config
}

// SyncConfig is the transfer service sync config
type SyncConfig struct {
	// EnableSync defines if data sync is enabled
	EnableSync bool `mapstructure:"enableSync"`
	// EnableIncrSync defines if incremental sync is enabled
	EnableIncrSync bool `mapstructure:"enableIncrSync"`
	// Name is the transfer service name
	Name string `mapstructure:"name"`
	// Role is the transfer service role
	Role SyncRole `mapstructure:"role"`
	// SyncIntervalHours is the full sync interval, unit: hour
	SyncIntervalHours int `mapstructure:"syncIntervalHours"`
	// TransMediumAddr is the transfer medium addresses
	TransMediumAddr []string `mapstructure:"transferMediumAddress"`
}

// Validate SyncConfig
func (s *SyncConfig) Validate() error {
	if !s.EnableSync {
		return nil
	}

	if len(s.Name) == 0 {
		return errors.New("sync config name is not set")
	}

	switch s.Role {
	case SyncRoleSrc:
		if s.SyncIntervalHours <= 0 {
			return fmt.Errorf("invalid sync interval hours: %d", s.SyncIntervalHours)
		}
	case SyncRoleDest:
	default:
		return fmt.Errorf("invalid sync role: %s", s.Role)
	}

	if len(s.TransMediumAddr) == 0 {
		return fmt.Errorf("transfer medium address is not set")
	}

	return nil
}

// SyncRole is the transfer service role in cmdb synchronization
type SyncRole string

const (
	// SyncRoleSrc is the role of the source cmdb
	SyncRoleSrc SyncRole = "src"
	// SyncRoleDest is the role of the destination cmdb
	SyncRoleDest SyncRole = "dest"
)

// DestExSyncConf is the destination cmdb transfer service extra sync config
type DestExSyncConf struct {
	IDRules     []IDRuleEnvConf   `mapstructure:"idRules"`
	InnerDataID []InnerDataIDConf `mapstructure:"innerDataID"`
}

// Validate DestExSyncConf
func (s *DestExSyncConf) Validate() error {
	if len(s.IDRules) == 0 {
		return errors.New("id rules are not set")
	}

	resIDRuleMap := make(map[types.ResType][]IDRuleInfo)
	for i, rule := range s.IDRules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("validate id rule(index: %d) failed, err: %v", i, err)
		}

		for _, info := range rule.Rules {
			resIDRuleMap[info.Resource] = append(resIDRuleMap[info.Resource], info.Rules...)
		}
	}

	// check if id rules with the same interval have the same step and different remainders
	for res, infos := range resIDRuleMap {
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].StartID < infos[j].StartID
		})

		var step, end int64
		remainderMap := make(map[int64]struct{})
		for i, rule := range infos {
			if i == 0 || (rule.StartID > end && end != types.InfiniteEndID) {
				step = rule.Step
				end = rule.EndID
				remainderMap = map[int64]struct{}{rule.StartID % rule.Step: {}}
				continue
			}

			if rule.Step != step {
				return fmt.Errorf("%s id rule(%+v) has different step with this interval's step %d", res, rule, step)
			}

			remainder := rule.StartID % step
			_, exists := remainderMap[remainder]
			if exists {
				return fmt.Errorf("%s id rule(%+v) has duplicate remainder %d", res, rule, remainder)
			}

			remainderMap[remainder] = struct{}{}
			if end != types.InfiniteEndID && (rule.EndID > end || rule.EndID == types.InfiniteEndID) {
				end = rule.EndID
			}
		}
	}

	// validate inner data id info
	for i, innerID := range s.InnerDataID {
		if err := innerID.Validate(); err != nil {
			return fmt.Errorf("validate inner id info(index: %d) failed, err: %v", i, err)
		}
	}

	return nil
}

// IDRuleEnvConf is the id rule config for one environment
type IDRuleEnvConf struct {
	// Name is the cmdb transfer service name
	Name string `mapstructure:"name"`
	// Rules are all the id rule config for this environment
	Rules []IDRuleResInfo `mapstructure:"rules"`
}

// Validate IDRuleEnvConf
func (s *IDRuleEnvConf) Validate() error {
	if s.Name == "" {
		return errors.New("id rule service name is not set")
	}

	if len(s.Rules) == 0 {
		return fmt.Errorf("%s resource id rule infos are not set", s.Name)
	}

	for i, rule := range s.Rules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("validate %s resource id rule info(index: %d) failed, err: %v", s.Name, i, err)
		}
	}

	return nil
}

// IDRuleResInfo is the id rule config for one resource
type IDRuleResInfo struct {
	// Resource is the resource type of the id generator
	Resource types.ResType `mapstructure:"resource"`
	// Rules are all the id rule config for this resource
	Rules []IDRuleInfo `mapstructure:"rules"`
}

// Validate IDRuleResInfo
func (s *IDRuleResInfo) Validate() error {
	if s.Resource == "" {
		return errors.New("id rule resource is not set")
	}

	if len(s.Rules) == 0 {
		return fmt.Errorf("%s id rule infos are not set", s.Resource)
	}

	sort.Slice(s.Rules, func(i, j int) bool {
		return s.Rules[i].StartID < s.Rules[j].StartID
	})

	var prevRule IDRuleInfo
	for i, rule := range s.Rules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("validate %s id rule info(index: %d) failed, err: %v", s.Resource, i, err)
		}

		if i != 0 && (rule.StartID < prevRule.EndID || prevRule.EndID == types.InfiniteEndID) {
			return fmt.Errorf("%s id rule info(index: %d, start: %d) overlaps with previous id rule(end: %d)",
				s.Resource, i, rule.StartID, prevRule.EndID)
		}

		prevRule = rule
	}

	return nil
}

// IDRuleInfo is the id rule info
type IDRuleInfo struct {
	// StartID is the start id of the interval in which the id generator takes effect
	StartID int64 `mapstructure:"startID"`
	// EndID is the end id of the interval in which the id generator takes effect
	EndID int64 `mapstructure:"endID"`
	// Step is the step of the id generator
	Step int64 `mapstructure:"step"`
}

// Validate IDRuleInfo
func (s *IDRuleInfo) Validate() error {
	if s.StartID < 0 {
		return fmt.Errorf("start id %d is invalid", s.StartID)
	}

	if s.EndID != types.InfiniteEndID && s.EndID <= s.StartID {
		return fmt.Errorf("end id %d is invalid", s.EndID)
	}

	if s.Step <= 0 {
		return fmt.Errorf("step %d is invalid", s.Step)
	}

	return nil
}

// InnerDataIDConf is the source cmdb inner data id config
type InnerDataIDConf struct {
	// Name is the source cmdb transfer service name
	Name string `mapstructure:"name"`
	// HostPool is the source cmdb host pool id info
	HostPool *HostPoolInfo `mapstructure:"hostPool"`
}

// Validate InnerDataIDConf
func (s *InnerDataIDConf) Validate() error {
	if s.Name == "" {
		return errors.New("inner data id info service name is not set")
	}

	return nil
}

// HostPoolInfo is the host pool id info
type HostPoolInfo struct {
	// Biz is the host pool biz id
	Biz int64 `mapstructure:"biz"`
	// Set is the host pool idle set id
	Set int64 `mapstructure:"set"`
	// Module is the host pool idle module id
	Module int64 `mapstructure:"module"`
}

// Validate HostPoolInfo
func (s *HostPoolInfo) Validate() error {
	if s.Biz <= 0 {
		return fmt.Errorf("host pool biz id %d is invalid", s.Biz)
	}

	if s.Set <= 0 {
		return fmt.Errorf("host pool set id %d is invalid", s.Set)
	}

	if s.Module <= 0 {
		return fmt.Errorf("host pool module id %d is invalid", s.Module)
	}

	return nil
}
