package operation

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (s *operation) SearchInstCount(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CoreUint64Response, err error) {
	resp = new(metadata.CoreUint64Response)
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

func (s *operation) SearchChartDataCommon(ctx context.Context, h http.Header, data metadata.ChartConfig) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/read/operation/chart/data/common"

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

func (s *operation) SearchOperationChart(ctx context.Context, h http.Header, data interface{}) (resp *metadata.SearchChartResponse, err error) {
	resp = new(metadata.SearchChartResponse)
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

func (s *operation) SearchTimerChartData(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
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

func (s *operation) UpdateChartPosition(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
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

func (s *operation) SearchChartCommon(ctx context.Context, h http.Header, data interface{}) (resp *metadata.SearchChartCommon, err error) {
	resp = new(metadata.SearchChartCommon)
	subPath := "/search/operation/chart/common"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) TimerFreshData(ctx context.Context, h http.Header, data interface{}) (resp *metadata.BoolResponse, err error) {
	resp = new(metadata.BoolResponse)
	subPath := "/start/operation/chart/timer"

	err = s.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
