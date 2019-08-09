package service

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (s *service) CreateServiceInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/proc/service_instance"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) DeleteServiceInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/deletemany/proc/service_instance"

	err = s.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) SearchServiceInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/findmany/proc/service_instance"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
