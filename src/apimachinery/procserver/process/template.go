package process

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateProcessTemplate TODO
func (p *process) CreateProcessTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/createmany/proc/proc_template"

	err = p.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteProcessTemplate TODO
func (p *process) DeleteProcessTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/deletemany/proc/proc_template"

	err = p.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchProcessTemplate TODO
func (p *process) SearchProcessTemplate(ctx context.Context, h http.Header,
	i *metadata.ListProcessTemplateWithServiceTemplateInput) (*metadata.MultipleProcessTemplate, errors.CCErrorCoder) {

	resp := new(metadata.MultipleProcessTemplateResult)
	subPath := "/findmany/proc/proc_template"

	err := p.client.Post().
		WithContext(ctx).
		Body(i).
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

// UpdateProcessTemplate TODO
func (p *process) UpdateProcessTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/proc/proc_template"

	err = p.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
