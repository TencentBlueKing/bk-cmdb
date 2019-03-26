package model

import (
	"configcenter/src/framework/clientset/types"
	"configcenter/src/framework/common/rest"
	"fmt"
)

type InstanceInterface interface {
}

type instClient struct {
	client rest.ClientInterface
}

func (s instClient) CreateObjectInstance(ctx *types.CreateSetCtx) (int64, error) {
	resp := new(types.CreateInstanceResult)
	subPath := fmt.Sprintf("/set/%d", ctx.SetID)
	err := s.client.Post().
		WithContext(ctx.Ctx).
		Body(ctx.Set).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return 0, &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return 0, &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return resp.Data.ID, nil
}

func (s instClient) DeleteObjectInstance(ctx *types.DeleteObjectCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/inst/%s/%s/%d", ctx.Tenancy, ctx.ObjectID, ctx.InstanceID)
	err := s.client.Delete().
		WithContext(ctx.Ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return nil
}

func (s instClient) UpdateObjectInstance(ctx *types.UpdateObjectCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/inst/%s/%s/%d", ctx.Tenancy, ctx.ObjectID, ctx.InstanceID)
	err := s.client.Put().
		WithContext(ctx.Ctx).
		Body(ctx.Object).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return nil
}

func (s instClient) ListObjectInstance(ctx *types.ListInstanceCtx) (*types.ListInfo, error) {
	resp := new(types.ListInstanceResult)
	subPath := fmt.Sprintf("/inst/search/owener/%s/object/%s", ctx.Tenancy, ctx.ObjectID)
	err := s.client.Post().
		WithContext(ctx.Ctx).
		Body(ctx.Filter).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return nil, &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return nil, &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return &resp.Data, nil
}
