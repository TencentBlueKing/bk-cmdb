package operation

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (s *operation) SearchInstCount(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "read/operation/inst/count"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) AggregateBizHost(ctx context.Context, h http.Header, data interface{}) (resp *metadata.AggregateIntResponse, err error) {
	resp = new(metadata.AggregateIntResponse)
	subPath := "read/operation/biz/host"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) CommonAggregate(ctx context.Context, h http.Header, data metadata.ChartOption) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "read/operation/common/aggregate"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) ModelInstAggregate(ctx context.Context, h http.Header, data interface{}) (resp *metadata.AggregateStringResponse, err error) {
	resp = new(metadata.AggregateStringResponse)
	subPath := "read/operation/model/inst"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) CreateOperationChart(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CoreUint64Response, err error) {
	resp = new(metadata.CoreUint64Response)
	subPath := "/create/operation/chart"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) DeleteOperationChart(ctx context.Context, h http.Header, id string) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("delete/operation/chart/%v", id)

	err = s.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) SearchOperationChart(ctx context.Context, h http.Header, data interface{}) (resp *metadata.SearchChartConfig, err error) {
	resp = new(metadata.SearchChartConfig)
	subPath := "/search/operation/chart"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) UpdateOperationChart(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "update/operation/chart"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) SearchOperationChartData(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/search/operation/chart/data"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) UpdateOperationChartPosition(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/operation/chart/position"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
