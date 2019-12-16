package service

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (s *service) CreateServiceTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/create/proc/service_template"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) DeleteServiceTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/delete/proc/service_template"

	err = s.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) SearchServiceTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/findmany/proc/service_template"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) UpdateServiceTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/update/proc/service_template"

	err = s.client.Put().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) RemoveTemplateBindingOnModule(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/delete/proc/template_binding_on_module"

	err = s.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
