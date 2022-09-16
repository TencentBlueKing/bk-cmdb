package service

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateServiceCategory TODO
func (s *service) CreateServiceCategory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/proc/service_category"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteServiceCategory TODO
func (s *service) DeleteServiceCategory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/delete/proc/service_category"

	err = s.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchServiceCategory TODO
func (s *service) SearchServiceCategory(ctx context.Context, h http.Header, opt *metadata.ListServiceCategoryOption) (
	*metadata.MultipleServiceCategory, errors.CCErrorCoder) {

	resp := new(metadata.MultipleServiceCategoryResult)
	subPath := "/findmany/proc/service_category"

	err := s.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// UpdateServiceCategory TODO
func (s *service) UpdateServiceCategory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/proc/service_category"

	err = s.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
