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
	"configcenter/src/auth/meta"
)

const (
	authAppCodeHeaderKey   string = "X-BK-APP-CODE"
	authAppSecretHeaderKey string = "X-BK-APP-SECRET"
	cmdbUser               string = "user"
	cmdbUserID             string = "system"
)

// NewAuthCenter create a instance to handle resources with blueking's AuthCenter.
func NewAuthCenter(tls *util.TLSClientConfig, authCfg map[string]string) (*AuthCenter, error) {
	client, err := util.NewClient(tls)
	if err != nil {
		return nil, err
	}

	var cfg AuthConfig
	address, exist := authCfg["auth.address"]
	if !exist {
		return nil, errors.New(`missing "address" configuration for auth center`)
	}

	cfg.Address = strings.Split(strings.Replace(address, " ", "", -1), ",")
	if len(cfg.Address) == 0 {
		return nil, errors.New(`invalid "address" configuration for auth center`)
	}

	cfg.AppSecret, exist = authCfg["auth.appSecret"]
	if !exist {
		return nil, errors.New(`missing "appSecret" configuration for auth center`)
	}

	if len(cfg.AppSecret) == 0 {
		return nil, errors.New(`invalid "appSecret" configuration for auth center`)
	}

	cfg.AppCode, exist = authCfg["auth.appCode"]
	if !exist {
		return nil, errors.New(`missing "appCode" configuration for auth center`)
	}

	if len(cfg.AppCode) == 0 {
		return nil, errors.New(`invalid "appCode" configuration for auth center`)
	}

	cfg.SystemID, exist = authCfg["auth.systemID"]
	if !exist {
		return nil, errors.New(`missing "systemID" configuration for auth center`)
	}

	if len(cfg.SystemID) == 0 {
		return nil, errors.New(`invalid "systemID" configuration for auth center`)
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
		client: rest.NewRESTClient(c, ""),
		Config: cfg,
		header: header,
	}, nil
}

// authCenter means BlueKing's authorize center,
// which is also a open source product.
type AuthCenter struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	header http.Header
}

func (ac *AuthCenter) Authorize(ctx context.Context, a *meta.Attribute) (decision meta.Decision, err error) {
	// TODO: fill this struct.
	info := &AuthBatch{
		Principal: Principal{
			Type: "user",
			ID:   a.User.UserName,
		},
	}

	// TODO: this operation may be wrong, because some api filters does not
	// fill the business id field, so these api should be normalized.
	if a.BusinessID != 0 {
		info.ScopeType = "biz"
		info.ScopeID = strconv.FormatInt(a.BusinessID, 10)
	} else {
		info.ScopeType = "system"
		info.ScopeID = "bk-cmdb"
	}

	info.ResourceActions = make([]ResourceAction, 0)
	for _, rsc := range a.Resources {
		info.ResourceActions = append(info.ResourceActions, ResourceAction{
			ActionID: rsc.Action.String(),
			ResourceInfo: ResourceInfo{
				ResourceType: rsc.Type.String(),
				// TODO: add resource id field.
				ResourceID: "",
			},
		})
	}

	resp := new(BatchResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources-perms/verify", ac.Config.SystemID)
	err = ac.client.Post().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(ac.header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return meta.Decision{}, err
	}

	if resp.Code != 0 {
		return meta.Decision{}, &AuthError{
			RequestID: resp.RequestID,
			Reason:    fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg),
		}
	}

	noAuth := make([]string, 0)
	for _, item := range resp.Data {
		if !item.IsPass {
			noAuth = append(noAuth, item.ResourceType)
		}
	}

	if len(noAuth) != 0 {
		return meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("resource [%s] do not have permission", strings.Join(noAuth, ",")),
		}, nil
	}

	return meta.Decision{Authorized: true}, nil
}

func (ac *AuthCenter) Register(ctx context.Context, r *meta.ResourceAttribute) error {
	if len(r.Basic.Type) == 0 {
		return errors.New("invalid resource attribute with empty object")
	}
	scope, err := ac.getScopeInfo(r)
	if err != nil {
		return err
	}

	rscID, err := ac.getResourceID(r.Layers)
	if err != nil {
		return err
	}
	info := &RegisterInfo{
		CreatorType: cmdbUser,
		CreatorID:   cmdbUserID,
		ScopeInfo:   *scope,
		ResourceInfo: ResourceInfo{
			ResourceType: string(r.Basic.Type),
			ResourceName: r.Basic.Name,
			ResourceID:   rscID,
		},
	}

	resp := new(ResourceResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources", ac.Config.SystemID)
	err = ac.client.Post().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(ac.header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{RequestID: resp.RequestID, Reason: fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)}
	}

	if !resp.Data.IsCreated {
		return &AuthError{resp.RequestID, fmt.Errorf("register resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (ac *AuthCenter) Deregister(ctx context.Context, r *meta.ResourceAttribute) error {
	if len(r.Basic.Type) == 0 {
		return errors.New("invalid resource attribute with empty object")
	}

	scope, err := ac.getScopeInfo(r)
	if err != nil {
		return err
	}

	rscID, err := ac.getResourceID(r.Layers)
	if err != nil {
		return err
	}

	info := &DeregisterInfo{
		ScopeInfo: *scope,
		ResourceInfo: ResourceInfo{
			ResourceType: r.Basic.Type.String(),
			ResourceID:   rscID,
		},
	}

	resp := new(ResourceResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources", ac.Config.SystemID)
	err = ac.client.Delete().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(ac.header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{resp.RequestID, fmt.Errorf("deregister resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)}
	}

	if !resp.Data.IsDeleted {
		return &AuthError{resp.RequestID, fmt.Errorf("deregister resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (ac *AuthCenter) Update(ctx context.Context, r *meta.ResourceAttribute) error {
	if len(r.Basic.Type) == 0 || len(r.Basic.Name) == 0 {
		return errors.New("invalid resource attribute with empty object or object name")
	}

	scope, err := ac.getScopeInfo(r)
	if err != nil {
		return err
	}

	rscID, err := ac.getResourceID(r.Layers)
	if err != nil {
		return err
	}
	info := &UpdateInfo{
		ScopeInfo: *scope,
		ResourceInfo: ResourceInfo{
			ResourceType: r.Basic.Type.String(),
			ResourceName: r.Basic.Name,
			ResourceID:   rscID,
		},
	}

	resp := new(ResourceResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources", ac.Config.SystemID)
	err = ac.client.Put().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(ac.header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{resp.RequestID, fmt.Errorf("update resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)}
	}

	if !resp.Data.IsUpdated {
		return &AuthError{resp.RequestID, fmt.Errorf("update resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (ac *AuthCenter) Get(ctx context.Context) error {
	panic("implement me")
}

func (ac *AuthCenter) getScopeInfo(r *meta.ResourceAttribute) (*ScopeInfo, error) {
	s := new(ScopeInfo)
	switch r.Basic.Name {
	case "set", "module":
		s.ScopeType = "biz"
	// TODO: add filter rules for scope info.
	default:
		return nil, fmt.Errorf("unsupported scope type or info for %s", r.Basic.Name)
	}
	return s, nil
}

func (ac *AuthCenter) getResourceID(layers []meta.Item) (string, error) {
	var id string
	for _, item := range layers {
		if len(item.Name) == 0 || len(item.Type) == 0 {
			return "", fmt.Errorf("invalid resoutece item %s/%d", item.Name, item.InstanceID)
		}
		id = fmt.Sprintf("%s/%s:%d", item.Type, item.Name, item.InstanceID)
	}
	id = strings.TrimLeft(id, "/")

	return id, nil
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
