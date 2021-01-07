package operation

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (s *operation) SearchInstCount(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CoreUint64Response, err error) {
	resp = new(metadata.CoreUint64Response)
	subPath := "/find/operation/inst/count"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) SearchChartData(ctx context.Context, h http.Header, data metadata.ChartConfig) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/operation/chart/data"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
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
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) DeleteOperationChart(ctx context.Context, h http.Header, id string) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "delete/operation/chart/%v"

	err = s.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, id).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) SearchOperationCharts(ctx context.Context, h http.Header, data interface{}) (resp *metadata.SearchChartResponse, err error) {
	resp = new(metadata.SearchChartResponse)
	subPath := "/findmany/operation/chart"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
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
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) SearchTimerChartData(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/operation/timer/chart/data"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
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
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *operation) SearchChartCommon(ctx context.Context, h http.Header, data interface{}) (resp *metadata.SearchChartCommon, err error) {
	resp = new(metadata.SearchChartCommon)
	subPath := "/find/operation/chart/common"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
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
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
