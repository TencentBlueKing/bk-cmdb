package privilege

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type PrivilegeInterface interface {
	CreateUserGroup(ctx context.Context, supplierAcct string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	DeleteUserGroup(ctx context.Context, supplierAcct string, groupID string, h http.Header) (resp *metadata.Response, err error)
	UpdateUserGroup(ctx context.Context, supplierAcct string, groupID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	SearchUserGroup(ctx context.Context, supplierAcct string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	UpdateUserGroupPrivi(ctx context.Context, supplierAcct string, groupID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	GetUserGroupPrivi(ctx context.Context, supplierAcct string, groupID string, h http.Header) (resp *metadata.Response, err error)
	GetUserPrivi(ctx context.Context, supplierAcct string, userName string, h http.Header) (resp *metadata.Response, err error)
	CreatePrivilege(ctx context.Context, supplierAcct string, objID string, propertyID string, h http.Header) (resp *metadata.Response, err error)
	GetPrivilege(ctx context.Context, supplierAcct string, objID string, propertyID string, h http.Header) (resp *metadata.Response, err error)
}

func NewPrivilegeInterface(client rest.ClientInterface) PrivilegeInterface {
	return &privilege{client: client}
}

type privilege struct {
	client rest.ClientInterface
}
