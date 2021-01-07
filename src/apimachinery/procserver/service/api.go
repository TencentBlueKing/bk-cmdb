package service

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type ServiceClientInterface interface {
	CreateServiceCategory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteServiceCategory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SearchServiceCategory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	UpdateServiceCategory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)

	CreateServiceInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	UpdateServiceInstances(ctx context.Context, h http.Header, bizID int64, data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteServiceInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SearchServiceInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SearchServiceInstanceBySetTemplate(ctx context.Context, appID string, h http.Header, data map[string]interface{}) (resp *metadata.ResponseInstData, err error)
	DiffServiceInstanceWithTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SyncServiceInstanceByTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	ServiceInstanceAddLabels(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	ServiceInstanceRemoveLabels(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	ServiceInstanceFindLabels(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)

	CreateServiceTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error)
	DeleteServiceTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error)
	SearchServiceTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error)
	UpdateServiceTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error)
	RemoveTemplateBindingOnModule(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error)
}

func NewServiceClientInterface(client rest.ClientInterface) ServiceClientInterface {
	return &service{client: client}
}

type service struct {
	client rest.ClientInterface
}
