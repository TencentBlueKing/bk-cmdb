package privilege

import (
	"context"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/core/cc/api"
)

type PrivilegeInterface interface {
	CreateUserGroup(ctx context.Context, supplierAcct string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	DeleteUserGroup(ctx context.Context, supplierAcct string, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdateUserGroup(ctx context.Context, supplierAcct string, groupID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SearchUserGroup(ctx context.Context, supplierAcct string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	UpdateUserGroupPrivi(ctx context.Context, supplierAcct string, groupID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	GetUserGroupPrivi(ctx context.Context, supplierAcct string, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	GetUserPrivi(ctx context.Context, supplierAcct string, userName string, h util.Headers) (resp *api.BKAPIRsp, err error)
	CreatePrivilege(ctx context.Context, supplierAcct string, objID string, propertyID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	GetPrivilege(ctx context.Context, supplierAcct string, objID string, propertyID string, h util.Headers) (resp *api.BKAPIRsp, err error)
}

func NewPrivilegeInterface(client rest.ClientInterface) PrivilegeInterface {
	return &privilege{client: client}
}

type privilege struct {
	client rest.ClientInterface
}
