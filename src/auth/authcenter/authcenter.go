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

package authcenter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/auth/authcenter/permit"
	"configcenter/src/auth/meta"
	"configcenter/src/common/blog"
)

const (
	authAppCodeHeaderKey   string = "X-BK-APP-CODE"
	authAppSecretHeaderKey string = "X-BK-APP-SECRET"
	cmdbUser               string = "user"
	cmdbUserID             string = "system"
)

// ParseConfigFromKV returns a new config
func ParseConfigFromKV(prefix string, configmap map[string]string) (AuthConfig, error) {
	var cfg AuthConfig
	enable, exist := configmap[prefix+".enable"]
	if !exist {
		return AuthConfig{}, nil
	}

	var err error
	cfg.Enable, err = strconv.ParseBool(enable)
	if err != nil {
		return AuthConfig{}, errors.New(`invalid auth "enable" value`)
	}

	if !cfg.Enable {
		return AuthConfig{}, nil
	}

	address, exist := configmap[prefix+".address"]
	if !exist {
		return cfg, errors.New(`missing "address" configuration for auth center`)
	}

	cfg.Address = strings.Split(strings.Replace(address, " ", "", -1), ",")
	if len(cfg.Address) == 0 {
		return cfg, errors.New(`invalid "address" configuration for auth center`)
	}
	for i := range cfg.Address {
		if !strings.HasSuffix(cfg.Address[i], "/") {
			cfg.Address[i] = cfg.Address[i] + "/"
		}
	}

	cfg.AppSecret, exist = configmap[prefix+".appSecret"]
	if !exist {
		return cfg, errors.New(`missing "appSecret" configuration for auth center`)
	}

	if len(cfg.AppSecret) == 0 {
		return cfg, errors.New(`invalid "appSecret" configuration for auth center`)
	}

	cfg.AppCode, exist = configmap[prefix+".appCode"]
	if !exist {
		return cfg, errors.New(`missing "appCode" configuration for auth center`)
	}

	if len(cfg.AppCode) == 0 {
		return cfg, errors.New(`invalid "appCode" configuration for auth center`)
	}

	cfg.SystemID = SystemIDCMDB

	return cfg, nil
}

// NewAuthCenter create a instance to handle resources with blueking's AuthCenter.
func NewAuthCenter(tls *util.TLSClientConfig, cfg AuthConfig) (*AuthCenter, error) {
	blog.V(5).Infof("new auth center client with parameters tls: %+v, cfg: %+v", tls, cfg)
	if !cfg.Enable {
		return new(AuthCenter), nil
	}
	client, err := util.NewClient(tls)
	if err != nil {
		return nil, err
	}

	c := &util.Capability{
		Client: client,
		Discover: &acDiscovery{
			servers: cfg.Address,
		},
		Throttle: flowctrl.NewRateLimiter(1000, 1000),
		Mock: util.MockInfo{
			Mocked: false,
		},
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set(authAppCodeHeaderKey, cfg.AppCode)
	header.Set(authAppSecretHeaderKey, cfg.AppSecret)

	return &AuthCenter{
		Config: cfg,
		authClient: &authClient{
			client:      rest.NewRESTClient(c, ""),
			Config:      cfg,
			basicHeader: header,
		},
	}, nil
}

// AuthCenter means BlueKing's authorize center,
// which is also a open source product.
type AuthCenter struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	header     http.Header
	authClient *authClient
}

func (ac *AuthCenter) Authorize(ctx context.Context, a *meta.AuthAttribute) (decision meta.Decision, err error) {
	blog.V(5).Infof("AuthCenter Config is: %+v", ac.Config)
	if !ac.Config.Enable {
		blog.V(5).Infof("AuthCenter Config is disabled. config: %+v", ac.Config)
		return meta.Decision{Authorized: true}, nil
	}

	batchresult, err := ac.AuthorizeBatch(ctx, a.User, a.Resources...)
	noAuth := make([]string, 0)
	for i, item := range batchresult {
		if !item.Authorized {
			noAuth = append(noAuth, fmt.Sprintf("resource [%v] permission deny by reason: %s", a.Resources[i].Type, item.Reason))
		}
	}

	if len(noAuth) > 0 {
		return meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("%v", noAuth),
		}, nil
	}

	return meta.Decision{Authorized: true}, nil
}

func (ac *AuthCenter) AuthorizeBatch(ctx context.Context, user meta.UserInfo, resources ...meta.ResourceAttribute) (decisions []meta.Decision, err error) {
	if !ac.Config.Enable {
		decisions = make([]meta.Decision, len(resources), len(resources))
		for i := range decisions {
			decisions[i].Authorized = true
		}
	}

	header := http.Header{}
	header.Set(AuthSupplierAccountHeaderKey, user.SupplierAccount)

	type AuthResult struct {
		meta.ResourceAttribute
		*meta.Decision     // must use pointer
		exactResourceIndex int
		anyResourceIndex   int
	}

	biz2res := map[int64][]AuthResult{}
	decisions = make([]meta.Decision, len(resources), len(resources))
	for i := 0; i < len(resources); i++ {
		biz2res[resources[i].BusinessID] = append(biz2res[resources[i].BusinessID], AuthResult{ResourceAttribute: resources[i], Decision: &decisions[i]})
	}

	exactResourceInfo := AuthBatch{
		Principal: Principal{
			Type: "user",
			ID:   user.UserName,
		},
	}
	anyResourceInfo := AuthBatch{
		Principal: Principal{
			Type: "user",
			ID:   user.UserName,
		},
	}
	for biz, ress := range biz2res {
		if biz > 0 {
			exactResourceInfo.ScopeType = ScopeTypeIDBiz
			exactResourceInfo.ScopeID = strconv.FormatInt(biz, 10)
			anyResourceInfo.ScopeType = ScopeTypeIDBiz
			anyResourceInfo.ScopeID = strconv.FormatInt(biz, 10)
		} else {
			exactResourceInfo.ScopeType = ScopeTypeIDSystem
			exactResourceInfo.ScopeID = SystemIDCMDB
			anyResourceInfo.ScopeType = ScopeTypeIDSystem
			anyResourceInfo.ScopeID = SystemIDCMDB
		}
		exactResourceIndex := 0
		anyResourceIndex := 0
		for ressindex := range ress {
			if permit.IsPermit(&ress[ressindex].ResourceAttribute) {
				blog.Debug("permited")
				ress[ressindex].Decision.Authorized = true
			} else {
				blog.Debug("query permit %+v", ress[ressindex])
				rscInfo, err := adaptor(&ress[ressindex].ResourceAttribute)
				if err != nil {
					ress[ressindex].Decision.Authorized = false
					ress[ressindex].Decision.Reason = fmt.Sprintf("adaptor resource info failed, err: %v", err)
					continue
				}

				actionID, err := adaptorAction(&ress[ressindex].ResourceAttribute)
				if err != nil {
					ress[ressindex].Decision.Authorized = false
					ress[ressindex].Decision.Reason = fmt.Sprintf("adaptor action info failed, err: %v", err)
					continue
				}

				resourceAction := ResourceAction{
					ActionID:     actionID,
					ResourceInfo: *rscInfo,
				}
				blog.Debug("query param %+v", resourceAction)
				if len(rscInfo.ResourceID) > 0 {
					exactResourceInfo.ResourceActions = append(exactResourceInfo.ResourceActions, resourceAction)
					ress[ressindex].exactResourceIndex = exactResourceIndex
					exactResourceIndex++
				} else {
					anyResourceInfo.ResourceActions = append(exactResourceInfo.ResourceActions, resourceAction)
					ress[ressindex].anyResourceIndex = anyResourceIndex
					anyResourceIndex++
				}

			}
		}

		if len(exactResourceInfo.ResourceActions) > 0 {
			batchresult, err := ac.authClient.verifyExactResourceBatch(ctx, header, &exactResourceInfo)
			if err != nil {
				reason := fmt.Sprintf("verify failed, err: %v", err)
				for _, res := range ress {
					res.Authorized = false
					res.Reason = reason
				}
				continue
			}
			for ressindex := range ress {
				if ress[ressindex].Authorized || len(ress[ressindex].Reason) > 0 {
					continue
				}
				if ress[ressindex].exactResourceIndex >= len(batchresult) {
					ress[ressindex].Authorized = false
					ress[ressindex].Reason = fmt.Sprintf("index out of range, %d:%d", ress[ressindex].exactResourceIndex, len(batchresult))
					continue
				}
				if batchresult[ress[ressindex].exactResourceIndex].IsPass {
					ress[ressindex].Authorized = true
				} else {
					ress[ressindex].Authorized = false
					ress[ressindex].Reason = "permission deny"
				}
			}
		}
		if len(anyResourceInfo.ResourceActions) > 0 {
			batchresult, err := ac.authClient.verifyAnyResourceBatch(ctx, header, &anyResourceInfo)
			if err != nil {
				reason := fmt.Sprintf("verify failed, err: %v", err)
				for _, res := range ress {
					res.Authorized = false
					res.Reason = reason
				}
				continue
			}
			for ressindex := range ress {
				if ress[ressindex].Authorized || len(ress[ressindex].Reason) > 0 {
					continue
				}
				if ress[ressindex].anyResourceIndex >= len(batchresult) {
					ress[ressindex].Authorized = false
					ress[ressindex].Reason = fmt.Sprintf("index out of range, %d:%d", ress[ressindex].anyResourceIndex, len(batchresult))
					continue
				}
				if batchresult[ress[ressindex].anyResourceIndex].IsPass {
					ress[ressindex].Authorized = true
				} else {
					ress[ressindex].Authorized = false
					ress[ressindex].Reason = "permission deny"
				}
			}
		}
		exactResourceInfo.ResourceActions = nil
		anyResourceInfo.ResourceActions = nil
	}

	return
}

func (ac *AuthCenter) RegisterResource(ctx context.Context, rs ...meta.ResourceAttribute) error {
	if !ac.Config.Enable {
		return nil
	}

	if len(rs) <= 0 {
		// not resource should be register
		return nil
	}
	info := RegisterInfo{}
	info.CreatorType = cmdbUser
	info.CreatorID = cmdbUserID
	header := http.Header{}
	for _, r := range rs {
		if len(r.Basic.Type) == 0 {
			return errors.New("invalid resource attribute with empty object")
		}
		scope, err := ac.getScopeInfo(&r)
		if err != nil {
			return err
		}

		rscInfo, err := adaptor(&r)
		if err != nil {
			return fmt.Errorf("adaptor resource info failed, err: %v", err)
		}
		entity := ResourceEntity{}
		entity.ScopeID = scope.ScopeID
		entity.ScopeType = scope.ScopeType
		entity.ResourceType = rscInfo.ResourceType
		entity.ResourceID = rscInfo.ResourceID
		entity.ResourceName = rscInfo.ResourceName

		// TODO replace register with batch createorupdate interface, currently is register one by one.
		info.Resources = make([]ResourceEntity, 0)
		info.Resources = append(info.Resources, entity)
		header.Set(AuthSupplierAccountHeaderKey, r.SupplierAccount)
		ac.authClient.registerResource(ctx, header, &info)
	}
	return nil
}

func (ac *AuthCenter) DeregisterResource(ctx context.Context, rs ...meta.ResourceAttribute) error {
	if !ac.Config.Enable {
		return nil
	}
	if len(rs) <= 0 {
		// not resource should be deregister
		return nil
	}
	info := DeregisterInfo{}
	header := http.Header{}
	for _, r := range rs {
		if len(r.Basic.Type) == 0 {
			return errors.New("invalid resource attribute with empty object")
		}

		scope, err := ac.getScopeInfo(&r)
		if err != nil {
			return err
		}

		rscInfo, err := adaptor(&r)
		if err != nil {
			return fmt.Errorf("adaptor resource info failed, err: %v", err)
		}

		entity := ResourceEntity{}
		entity.ScopeID = scope.ScopeID
		entity.ScopeType = scope.ScopeType
		entity.ResourceType = rscInfo.ResourceType
		entity.ResourceID = rscInfo.ResourceID
		entity.ResourceName = rscInfo.ResourceName

		info.Resources = append(info.Resources, entity)

		header.Set(AuthSupplierAccountHeaderKey, r.SupplierAccount)
	}

	return ac.authClient.deregisterResource(ctx, header, &info)
}

func (ac *AuthCenter) UpdateResource(ctx context.Context, r *meta.ResourceAttribute) error {
	if !ac.Config.Enable {
		return nil
	}

	if len(r.Basic.Type) == 0 || len(r.Basic.Name) == 0 {
		return errors.New("invalid resource attribute with empty object or object name")
	}

	scope, err := ac.getScopeInfo(r)
	if err != nil {
		return err
	}

	rscInfo, err := adaptor(r)
	if err != nil {
		return fmt.Errorf("adaptor resource info failed, err: %v", err)
	}
	info := &UpdateInfo{
		ScopeInfo:    *scope,
		ResourceInfo: *rscInfo,
	}

	header := http.Header{}
	header.Set(AuthSupplierAccountHeaderKey, r.SupplierAccount)
	return ac.authClient.updateResource(ctx, header, info)
}

func (ac *AuthCenter) Get(ctx context.Context) error {
	panic("implement me")
}

func (ac *AuthCenter) getScopeInfo(r *meta.ResourceAttribute) (*ScopeInfo, error) {
	s := new(ScopeInfo)
	// TODO: this operation may be wrong, because some api filters does not
	// fill the business id field, so these api should be normalized.
	if r.BusinessID > 0 {
		s.ScopeType = ScopeTypeIDBiz
		s.ScopeID = strconv.FormatInt(r.BusinessID, 10)
	} else {
		s.ScopeType = ScopeTypeIDSystem
		s.ScopeID = SystemIDCMDB
	}
	return s, nil
}

type acDiscovery struct {
	// auth's servers address, must prefixed with http:// or https://
	servers []string
	index   int
	sync.Mutex
}

func (s *acDiscovery) GetServers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	num := len(s.servers)
	if num == 0 {
		return []string{}, errors.New("oops, there is no server can be used")
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.servers[s.index-1:], s.servers[:s.index-1]...), nil
	} else {
		s.index = 0
		return append(s.servers[num-1:], s.servers[:num-1]...), nil
	}
}
