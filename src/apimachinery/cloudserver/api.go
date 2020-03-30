package cloudserver

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (c *cloudserver) CreateAccount(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/cloud/account"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) SearchAccount(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error) {
	resp = new(metadata.SearchResp)
	subPath := "/findmany/cloud/account"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) UpdateAccount(ctx context.Context, h http.Header, accountID int64, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/cloud/account/%d"

	err = c.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, accountID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) DeleteAccount(ctx context.Context, h http.Header, accountID int64) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/delete/cloud/account/%d"

	err = c.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, accountID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) CreateSyncTask(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/cloud/sync/task"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) SearchSyncTask(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error) {
	resp = new(metadata.SearchResp)
	subPath := "/findmany/cloud/sync/task"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) UpdateSyncTask(ctx context.Context, h http.Header, taskID int64, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/cloud/sync/task/%d"

	err = c.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, taskID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) DeleteSyncTask(ctx context.Context, h http.Header, taskID int64) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/delete/cloud/sync/task/%d"

	err = c.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, taskID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) SearchSyncHistory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error) {
	resp = new(metadata.SearchResp)
	subPath := "/findmany/cloud/sync/history"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloudserver) SearchSyncRegion(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error) {
	resp = new(metadata.SearchResp)
	subPath := "/findmany/cloud/sync/region"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
