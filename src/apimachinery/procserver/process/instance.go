package process

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateProcessInstance TODO
func (p *process) CreateProcessInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/proc/process_instance"

	err = p.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteProcessInstance delete process instances by biz id and process ids
func (p *process) DeleteProcessInstance(ctx context.Context, h http.Header,
	data *metadata.DeleteProcessInstanceInServiceInstanceInput) error {
	resp := new(metadata.Response)
	subPath := "/delete/proc/process_instance"

	err := p.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.CCError()
	}
	return nil
}

// SearchProcessInstance search process instances by biz id and service instance id
func (p *process) SearchProcessInstance(ctx context.Context, h http.Header, data *metadata.ListProcessInstancesOption) (
	[]metadata.ProcessInstance, error) {

	resp := new(metadata.ListProcessInstancesRsp)
	subPath := "/findmany/proc/process_instance"

	err := p.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if resp.CCError() != nil {
		return nil, resp.CCError()
	}
	return resp.Data, nil
}

// UpdateProcessInstance TODO
func (p *process) UpdateProcessInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/proc/process_instance"

	err = p.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// ListProcessRelatedInfo TODO
func (p *process) ListProcessRelatedInfo(ctx context.Context, h http.Header, bizID int64, data metadata.ListProcessRelatedInfoOption) (resp *metadata.ListProcessRelatedInfoResponse, err error) {
	resp = new(metadata.ListProcessRelatedInfoResponse)
	subPath := "/findmany/proc/process_related_info/biz/%d"

	err = p.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// ListProcessInstancesNameIDsInModule TODO
func (p *process) ListProcessInstancesNameIDsInModule(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/findmany/proc/process_instance/name_ids"

	err = p.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// ListProcessInstancesDetailsByIDs TODO
func (p *process) ListProcessInstancesDetailsByIDs(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/findmany/proc/process_instance/detail/by_ids"

	err = p.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// ListProcessInstancesDetails TODO
func (p *process) ListProcessInstancesDetails(ctx context.Context, h http.Header, bizID int64, data metadata.ListProcessInstancesDetailsOption) (resp *metadata.MapArrayResponse, err error) {
	resp = new(metadata.MapArrayResponse)
	subPath := "/findmany/proc/process_instance/detail/biz/%d"

	err = p.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateProcessInstancesByIDs TODO
func (p *process) UpdateProcessInstancesByIDs(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/proc/process_instance/by_ids"

	err = p.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
