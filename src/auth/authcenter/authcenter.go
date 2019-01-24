package authcenter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/auth"
)

const (
	authHeaderKey string = "X-BK-APP-CODE和X-BK-APP-SECRET"
	cmdbUser      string = "user"
	cmdbUserID    string = "system"
)

// NewAuthCenter create a instance to handle resources with blueking's AuthCenter.
func NewAuthCenter(tls *util.TLSClientConfig, cfg *AuthConfig) (auth.ResourceHandler, error) {
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
	header.Set(authHeaderKey, cfg.AppSecret)

	return &authCenter{
		client: rest.NewRESTClient(c, ""),
		header: header,
	}, nil
}

// authCenter means BlueKing's authorize center,
// which is also a open source product.
type authCenter struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	header http.Header
}

func (ac *authCenter) Register(ctx context.Context, r *auth.ResourceAttribute) (requestID string, err error) {
	if len(r.Object) == 0 {
		return "", errors.New("invalid resource attribute with empty object")
	}
	scope, err := ac.getScopeInfo(r)
	if err != nil {
		return "", err
	}

	rscID, err := ac.getResourceID(r.Layers)
	if err != nil {
		return "", err
	}
	info := &RegisterInfo{
		CreatorType: cmdbUser,
		CreatorID:   cmdbUserID,
		ScopeInfo:   *scope,
		ResourceInfo: ResourceInfo{
			ResourceType: r.Object,
			ResourceName: r.ObjectName,
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
		return "", err
	}

	if resp.Code != 0 {
		return resp.RequestID, fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)
	}

	if !resp.Data.IsCreated {
		return resp.RequestID, fmt.Errorf("register resource failed, error code: %d", resp.Code)
	}

	return resp.RequestID, nil
}

func (ac *authCenter) Deregister(ctx context.Context, r *auth.ResourceAttribute) (requestID string, err error) {
	if len(r.Object) == 0 {
		return "", errors.New("invalid resource attribute with empty object")
	}

	scope, err := ac.getScopeInfo(r)
	if err != nil {
		return "", err
	}

	rscID, err := ac.getResourceID(r.Layers)
	if err != nil {
		return "", err
	}

	info := &DeregisterInfo{
		ScopeInfo: *scope,
		ResourceInfo: ResourceInfo{
			ResourceType: r.Object,
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
		return "", err
	}

	if resp.Code != 0 {
		return resp.RequestID, fmt.Errorf("deregister resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)
	}

	if !resp.Data.IsDeleted {
		return resp.RequestID, fmt.Errorf("deregister resource failed, error code: %d", resp.Code)
	}

	return resp.RequestID, nil
}

func (ac *authCenter) Update(ctx context.Context, r *auth.ResourceAttribute) (requestID string, err error) {
	if len(r.Object) == 0 || len(r.ObjectName) == 0 {
		return "", errors.New("invalid resource attribute with empty object or object name")
	}

	scope, err := ac.getScopeInfo(r)
	if err != nil {
		return "", err
	}

	rscID, err := ac.getResourceID(r.Layers)
	if err != nil {
		return "", err
	}
	info := &UpdateInfo{
		ScopeInfo: *scope,
		ResourceInfo: ResourceInfo{
			ResourceType: r.Object,
			ResourceName: r.ObjectName,
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
		return "", err
	}

	if resp.Code != 0 {
		return resp.RequestID, fmt.Errorf("update resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)
	}

	if !resp.Data.IsUpdated {
		return resp.RequestID, fmt.Errorf("update resource failed, error code: %d", resp.Code)
	}

	return resp.RequestID, nil
}

func (ac *authCenter) Get(ctx context.Context) error {
	panic("implement me")
}

func (ac *authCenter) getScopeInfo(r *auth.ResourceAttribute) (*ScopeInfo, error) {
	s := new(ScopeInfo)
	switch r.Object {
	case "set", "module":
		s.ScopeType = "biz"
	// TODO: add filter rules for scope info.
	default:
		return nil, fmt.Errorf("unsupported scope type or info for %s", r.Object)
	}
	return s, nil
}

func (ac *authCenter) getResourceID(layers []auth.Item) (string, error) {
	var id string
	for _, item := range layers {
		if len(item.Object) == 0 || len(item.Object) == 0 {
			return "", fmt.Errorf("invalid resoutece item %s/%d", item.Object, item.InstanceID)
		}
		id = fmt.Sprintf("%s/%s:%d", id, item.Object, item.InstanceID)
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
