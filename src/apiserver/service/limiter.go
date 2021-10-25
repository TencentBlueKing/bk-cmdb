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

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

type Limiter struct {
	rd    *registerdiscover.RegDiscv
	rules map[string]*metadata.LimiterRule
	lock  sync.RWMutex
}

// NewLimiter new a limiter struct
func NewLimiter(rd *registerdiscover.RegDiscv) *Limiter {
	return &Limiter{
		rd:    rd,
		rules: make(map[string]*metadata.LimiterRule),
	}
}

// LenOfRules get the count of limiter's rules
func (l *Limiter) LenOfRules() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.rules)
}

func (l *Limiter) addRule(rule *metadata.LimiterRule) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.rules[rule.RuleName] = rule
}

func (l *Limiter) removeRule(ruleName string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.rules, ruleName)
}

func (l *Limiter) getAllRules() []*metadata.LimiterRule {
	var limiterRules []*metadata.LimiterRule
	l.lock.RLock()
	defer l.lock.RUnlock()
	for _, rule := range l.rules {
		limiterRules = append(limiterRules, rule)
	}
	return limiterRules
}

// SyncLimiterRules sync the api limiter rules from register and discover
func (l *Limiter) SyncLimiterRules(ctx context.Context) error {
	blog.Info("begin SyncLimiterRules")
	kvs, err := l.rd.GetWithPrefix(types.CCDiscoverBaseLimiter)
	if err != nil {
		blog.Errorf("fail to get key: %s, err: %v", types.CCDiscoverBaseLimiter, err)
		return err
	}
	for _, kv := range kvs {
		rule := new(metadata.LimiterRule)
		if err := json.Unmarshal([]byte(kv.Value), rule); err != nil {
			blog.Errorf("fail to Unmarshal limiter rule: %s, err: %v", kv.Value, err)
			continue
		}

		if err := rule.Verify(); err != nil {
			blog.Errorf("fail to verify  limiter rule: %v, err:%v", rule, err)
			continue
		}
		l.addRule(rule)
	}

	limitChan, err := l.rd.Watch(ctx, types.CCDiscoverBaseLimiter)
	if err != nil {
		return err
	}

	go func() {
		for event := range limitChan {
			blog.Infof("limiter rule received one event from path: %s", types.CCDiscoverBaseLimiter)
			switch event.Type {
			case registerdiscover.EventPut:
				l.updateLimiterRule(event.Key, event.Value)
			case registerdiscover.EventDel:
				l.removeLimiterRule(event.Key, event.Value)
			default:
				blog.Errorf("limiter rule received unknown event type: %v", event.Type)
				continue
			}
		}
	}()

	return nil
}

func (l *Limiter) updateLimiterRule(key, data string) {
	if key == "" {
		blog.Errorf("limiter received invalid event, for key is empty")
		return
	}

	rule := new(metadata.LimiterRule)
	if err := json.Unmarshal([]byte(data), rule); err != nil {
		blog.Errorf("unmarshal rule info failed, key: %s, info: %s, err: %v", key, data, err)
		return
	}

	if err := rule.Verify(); err != nil {
		blog.Errorf("fail to verify limiter rule: %v, err: %v", rule, err)
		return
	}
	l.addRule(rule)
}

func (l *Limiter) removeLimiterRule(key, data string) {
	if key == "" {
		blog.Errorf("limiter received invalid event, for key is empty")
		return
	}

	splitArray := strings.Split(key, "/")
	l.removeRule(splitArray[len(splitArray)-1])

	blog.Infof("remove limiter rule: %v, key: %s", data, splitArray[len(splitArray)-1])
}

// GetMatchedRule get the matched limiter rule according request
func (l *Limiter) GetMatchedRule(req *restful.Request) *metadata.LimiterRule {
	header := req.Request.Header
	var matchedRule *metadata.LimiterRule
	var min int64 = 999999
	rules := l.getAllRules()
	for _, r := range rules {
		if r.AppCode == "" && r.User == "" && r.IP == "" && r.Url == "" && r.Method == "" {
			blog.Errorf("wrong rule format, one of appcode, user, ip, url, method must be set, r:%#v", *r)
			continue
		}
		if r.AppCode != "" {
			if r.AppCode != header.Get(common.BKHTTPRequestAppCode) {
				continue
			}
		}
		if r.User != "" {
			if r.User != header.Get(common.BKHTTPHeaderUser) {
				continue
			}
		}
		if r.IP != "" {
			hit := false
			ips := strings.Split(r.IP, ",")
			for _, ip := range ips {
				if strings.TrimSpace(ip) == strings.TrimSpace(header.Get(common.BKHTTPRequestRealIP)) {
					hit = true
					break
				}
			}
			if hit == false {
				continue
			}
		}
		if r.Method != "" {
			if util.Normalize(r.Method) != util.Normalize(req.Request.Method) {
				continue
			}
		}
		if r.Url != "" {
			match, err := regexp.MatchString(r.Url, req.Request.RequestURI)
			if err != nil {
				blog.Errorf("MatchString failed, r.Url:%s, reqURI:%s, err:%s",
					r.Url, req.Request.RequestURI, err.Error())
				continue
			}
			if !match {
				continue
			}
		}

		if r.DenyAll == true {
			matchedRule = r
			break
		}
		if r.Limit < min {
			min = r.Limit
			matchedRule = r
		}
	}
	return matchedRule
}

// AddRule add limit rule
func (l *Limiter) AddRule(rule *metadata.LimiterRule) error {
	if err := rule.Verify(); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", types.CCDiscoverBaseLimiter, rule.RuleName)
	oldRule, err := l.rd.Get(path)

	if err != nil && err != registerdiscover.ErrNoKey {
		return err
	}
	if len(oldRule) != 0 {
		return fmt.Errorf("the rule %s has already existed", rule.RuleName)
	}

	data, err := json.Marshal(rule)
	if err != nil {
		return err
	}
	err = l.rd.Put(path, string(data))
	if err != nil {
		return err
	}
	return nil
}

// GetRules get limit rules
func (l *Limiter) GetRules(ruleNames []string) ([]*metadata.LimiterRule, error) {
	if ruleNames == nil {
		return nil, fmt.Errorf("rulenames must be set")
	}

	var limiterRules []*metadata.LimiterRule
	for _, name := range ruleNames {
		path := fmt.Sprintf("%s/%s", types.CCDiscoverBaseLimiter, name)
		data, err := l.rd.Get(path)
		if err != nil {
			blog.Errorf("get rule %s err: %v", name, err)
			continue
		}
		rule := new(metadata.LimiterRule)
		if err := json.Unmarshal([]byte(data), rule); err != nil {
			blog.Errorf("unmarshal rule info failed, data: %s, err: %v", data, err)
			continue
		}
		limiterRules = append(limiterRules, rule)
	}
	return limiterRules, nil
}

// DelRules delete limit rules
func (l *Limiter) DelRules(ruleNames []string) error {
	if ruleNames == nil {
		return fmt.Errorf("rulenames must be set")
	}
	for _, name := range ruleNames {
		// delete rule in etcd
		path := fmt.Sprintf("%s/%s", types.CCDiscoverBaseLimiter, name)
		err := l.rd.Delete(path)
		if err != nil {
			blog.Warnf("delete rule error in etcd, name %s, err: %v", name, err)
		}
	}
	return nil
}

// GetAllRules get all limit rules
func (l *Limiter) GetAllRules() ([]*metadata.LimiterRule, error) {
	var limiterRules []*metadata.LimiterRule
	rules, err := l.rd.GetWithPrefix(types.CCDiscoverBaseLimiter)
	if err != nil {
		return nil, err
	}
	for _, data := range rules {
		rule := new(metadata.LimiterRule)
		if err := json.Unmarshal([]byte(data.Value), rule); err != nil {
			blog.Errorf("unmarshal rule info failed, rule: %s, err: %v", data.Value, err)
			continue
		}
		limiterRules = append(limiterRules, rule)
	}
	return limiterRules, nil
}
