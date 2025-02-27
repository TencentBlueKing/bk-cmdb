package service

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateServiceTemplate TODO
func (s *service) CreateServiceTemplate(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/create/proc/service_template"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteServiceTemplate TODO
func (s *service) DeleteServiceTemplate(ctx context.Context, h http.Header,
	input *metadata.DeleteServiceTemplatesInput) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/delete/proc/service_template"

	err := s.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// SearchServiceTemplate TODO
func (s *service) SearchServiceTemplate(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/findmany/proc/service_template"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// FindServiceTemplateCountInfo TODO
func (s *service) FindServiceTemplateCountInfo(ctx context.Context, h http.Header, bizID int64,
	data map[string]interface{}) (resp *metadata.ArrayResponse, err error) {
	resp = new(metadata.ArrayResponse)
	subPath := "/findmany/proc/service_template/count_info/biz/%d"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateServiceTemplate TODO
func (s *service) UpdateServiceTemplate(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/update/proc/service_template"

	err = s.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// RemoveTemplateBindingOnModule TODO
func (s *service) RemoveTemplateBindingOnModule(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.ResponseDataMapStr, err error) {
	resp = new(metadata.ResponseDataMapStr)
	subPath := "/delete/proc/template_binding_on_module"

	err = s.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// CreateServiceTemplateAllInfo TODO
func (s *service) CreateServiceTemplateAllInfo(ctx context.Context, h http.Header,
	opt *metadata.CreateSvcTempAllInfoOption) (int64, errors.CCErrorCoder) {

	resp := new(metadata.CreateResult)
	subPath := "/create/proc/service_template/all_info"

	err := s.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return 0, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return 0, err
	}

	return resp.Data.ID, nil
}

// UpdateServiceTemplateAllInfo TODO
func (s *service) UpdateServiceTemplateAllInfo(ctx context.Context, h http.Header,
	opt *metadata.UpdateSvcTempAllInfoOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/update/proc/service_template/all_info"

	err := s.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// GetServiceTemplateAllInfo TODO
func (s *service) GetServiceTemplateAllInfo(ctx context.Context, h http.Header, opt *metadata.GetSvcTempAllInfoOption) (
	*metadata.SvcTempAllInfo, errors.CCErrorCoder) {

	resp := new(metadata.GetSvcTempAllInfoResult)
	subPath := "/find/proc/service_template/all_info"

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

	return resp.Data, nil
}

// UpdateServiceTemplateAttribute TODO
func (s *service) UpdateServiceTemplateAttribute(ctx context.Context, h http.Header,
	opt *metadata.UpdateServTempAttrOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/update/proc/service_template/attribute"

	err := s.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// DeleteServiceTemplateAttribute TODO
func (s *service) DeleteServiceTemplateAttribute(ctx context.Context, h http.Header,
	opt *metadata.DeleteServTempAttrOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/delete/proc/service_template/attribute"

	err := s.client.Delete().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// ListServiceTemplateAttribute TODO
func (s *service) ListServiceTemplateAttribute(ctx context.Context, h http.Header,
	opt *metadata.ListServTempAttrOption) (*metadata.ServTempAttrData, errors.CCErrorCoder) {

	resp := new(metadata.ServiceTemplateAttributeResult)
	subPath := "/findmany/proc/service_template/attribute"

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

	return resp.Data, nil
}

// UpdateSrvTmplHostApplyStatus 更新服务模板主机自动应用状态
func (s *service) UpdateSrvTmplHostApplyStatus(ctx context.Context, h http.Header, bizID int64,
	opt *metadata.UpdateHostApplyEnableStatusOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/updatemany/proc/service_template/host_apply_enable_status/biz/%d"

	err := s.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// DeleteSrvTmplHostApplyRule 删除服务模板场景下主机自动应用规则
func (s *service) DeleteSrvTmplHostApplyRule(ctx context.Context, h http.Header, bizID int64,
	opt *metadata.DeleteHostApplyRuleOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/deletemany/proc/service_template/host_apply_rule/biz/%d"

	err := s.client.Delete().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// UpdateSrvTmplHostApplyRuleRun 编辑服务模板场景下主机自动应用规则
func (s *service) UpdateSrvTmplHostApplyRuleRun(ctx context.Context, h http.Header,
	opt *metadata.HostApplyServiceTemplateOption) (*metadata.HostApplyTaskResult, errors.CCErrorCoder) {

	resp := new(metadata.HostApplyTaskResp)
	subPath := "/updatemany/proc/service_template/host_apply_plan/run"

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

// FindSrvTmpl 获取服务模板
func (s *service) FindSrvTmpl(ctx context.Context, h http.Header, tmplID int64) (
	*metadata.ServiceTemplate, errors.CCErrorCoder) {

	resp := new(metadata.ServiceTemplateResp)
	subPath := "/find/proc/service_template/%d"

	err := s.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, tmplID).
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

// FindSrvTmplHostApplyRuleRelation 根据配置字段搜索服务模板
func (s *service) FindSrvTmplHostApplyRuleRelation(ctx context.Context, h http.Header,
	opt *metadata.RuleRelatedServiceTemplateOption) ([]metadata.SrvTemplate, errors.CCErrorCoder) {

	resp := new(metadata.ServiceTemplatesResponse)
	subPath := "/find/proc/service_template/host_apply_rule_related"

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

	return resp.Data, nil
}

// ListSrvTmplHostApplyRuleTaskStatus list host apply rule task status
func (s *service) ListSrvTmplHostApplyRuleTaskStatus(ctx context.Context, header http.Header,
	option *metadata.HostApplyTaskStatusOption) (*metadata.HostApplyTaskStatusRsp, errors.CCErrorCoder) {

	ret := new(metadata.HostApplyRulesTaskResp)
	subPath := "/findmany/proc/service_template/host_apply_plan/status"

	err := s.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return &ret.Data, nil
}
