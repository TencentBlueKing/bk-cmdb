/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/framework/core/errors"
	"github.com/tidwall/gjson"
)

func (ps *parseStream) processRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.process().
		ServiceInstance().
		ServiceTemplate().
		ServiceCategory().
		processTemplate()
		// remove process template and process template bound related api
		// processTemplate()
		// processTemplateBound()

	return ps
}

var (
	createProcessRegexp            = regexp.MustCompile(`^/api/v3/proc/[^\s/]+/[0-9]+/?$`)
	findProcessesInBusinessRegexp  = regexp.MustCompile(`^/api/v3/proc/search/[^\s/]+/[0-9]+/?$`)
	findProcessDetailsRegexp       = regexp.MustCompile(`^/api/v3/proc/[^\s/]+/[0-9]+/[0-9]+/?$`)
	deleteProcessRegexp            = regexp.MustCompile(`^/api/v3/proc/[^\s/]+/[0-9]+/[0-9]+/?$`)
	updateProcessRegexp            = regexp.MustCompile(`^/api/v3/proc/[^\s/]+/[0-9]+/[0-9]+/?$`)
	updateProcessBatchRegexp       = regexp.MustCompile(`^/api/v3/proc/[^\s/]+/[0-9]+/?$`)
	findModulesBindByProcessRegexp = regexp.MustCompile(`^/api/v3/proc/[^\s/]+/[0-9]+/[0-9]+/?$`)
	boundModuleToProcessRegexp     = regexp.MustCompile(`^/api/v3/proc/module/[^\s/]+/[0-9]+/[0-9]+/[^\s/]+/?$`)
	unboundModuleToProcessRegexp   = regexp.MustCompile(`^/api/v3/proc/module/[^\s/]+/[0-9]+/[0-9]+/[^\s/]+/?$`)
	findboundModuleToProcessRegexp = regexp.MustCompile(`^/api/v3/proc/module/[^\s/]+/[0-9]+/[0-9]+/?$`)
	findProcessInstanceRegexp      = regexp.MustCompile(`^/api/v3/proc/inst/[^\s/]+/[0-9]+/?$`)
	freshProcHostInstPattern       = "/api/v3/proc/process/refresh/hostinstnum"
)

func (ps *parseStream) process() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create a business operation.
	if ps.hitRegexp(createProcessRegexp, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create process, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.Create,
					Name:   string(meta.Process),
				},
			},
		}

		return ps
	}

	// find processes in a business
	if ps.hitRegexp(findProcessesInBusinessRegexp, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create process, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.FindMany,
					Name:   string(meta.Process),
				},
			},
		}

		return ps
	}

	// find a process's details
	if ps.hitRegexp(findProcessDetailsRegexp, http.MethodGet) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find process detail, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}

		procID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find process detail, but got invalid process id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.FindMany,
					Name:       string(meta.Process),
					InstanceID: procID,
				},
			},
		}

		return ps
	}

	// delete a process in a business.
	if ps.hitRegexp(deleteProcessRegexp, http.MethodDelete) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete process, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}

		procID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete process, but got invalid process id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.DeleteMany,
					Name:       string(meta.Process),
					InstanceID: procID,
				},
			},
		}

		return ps
	}

	// update a process
	if ps.hitRegexp(updateProcessRegexp, http.MethodPut) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update process, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}

		procID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update process, but got invalid process id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.Update,
					Name:       string(meta.Process),
					InstanceID: procID,
				},
			},
		}

		return ps
	}

	// update process batch.
	if ps.hitRegexp(updateProcessBatchRegexp, http.MethodPut) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update process batch, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.UpdateMany,
					Name:   string(meta.Process),
				},
			},
		}

		return ps
	}

	// find modules bounded by a process.
	if ps.hitRegexp(findModulesBindByProcessRegexp, http.MethodGet) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find modules bounded by process, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}

		procID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find modules bounded by process, but got invalid process id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.FindMany,
					Name:       string(meta.Process),
					InstanceID: procID,
				},
			},
		}

		return ps
	}

	// bounded a module to a process
	if ps.hitRegexp(boundModuleToProcessRegexp, http.MethodPut) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("bound module to process, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		procID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("bound module to process, but got invalid process id: %s", ps.RequestCtx.Elements[6])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.BoundModuleToProcess,
					InstanceID: procID,
				},
			},
		}

		return ps
	}

	// unbound a module with a process.
	if ps.hitRegexp(unboundModuleToProcessRegexp, http.MethodDelete) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("unbound module to process, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		procID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("unbound module to process, but got invalid process id: %s", ps.RequestCtx.Elements[6])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.UnboundModuleToProcess,
					InstanceID: procID,
				},
			},
		}
		return ps
	}

	// find bound a module with a process.
	if ps.hitRegexp(findboundModuleToProcessRegexp, http.MethodGet) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("unbound module to process, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		procID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("unbound module to process, but got invalid process id: %s", ps.RequestCtx.Elements[6])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.Find,
					InstanceID: procID,
				},
			},
		}
		return ps
	}

	// find a process instance details
	// TODO: config this api filter.
	if ps.hitRegexp(findProcessInstanceRegexp, http.MethodGet) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find process instance details, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.FindMany,
					Name:   string(meta.Process),
				},
			},
		}

		return ps
	}

	if ps.hitPattern(freshProcHostInstPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	return ps
}

var (
	createProcConfigTemplateRegexp     = regexp.MustCompile(`^/api/v3/template/[^\s/]+/[0-9]+/?$`)
	updateProcConfigTemplateRegexp     = regexp.MustCompile(`^/api/v3/template/[^\s/]+/[0-9]+/[0-9]+/?$`)
	findProcConfigTemplatesRegexp      = regexp.MustCompile(`^/api/v3/template/search/[^\s/]+/[0-9]+/?$`)
	deleteProcConfigTemplateRegexp     = regexp.MustCompile(`^/api/v3/template/[^\s/]+/[0-9]+/[0-9]+/?$`)
	findProcessTemplateVersionRegexp   = regexp.MustCompile(`^/api/v3/template/version/search/[^\s/]+/[0-9]+/[0-9]+/?$`)
	createProcessTemplateVersionRegexp = regexp.MustCompile(`^/api/v3/template/version/[^\s/]+/[0-9]+/[0-9]+/?$`)
	updateProcessTemplateVersionRegexp = regexp.MustCompile(`^/api/v3/template/version/[^\s/]+/[0-9]+/[0-9]+/[0-9]+/?$`)
	previewProcessConfigRegexp         = regexp.MustCompile(`^/api/v3/proc/template/[^\s/]+/[0-9]+/[0-9]+/?$`)
)

// Deprecated: unused apis

func (ps *parseStream) processTemplate() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create a process config template.
	if ps.hitRegexp(createProcConfigTemplateRegexp, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create process config template, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.Create,
					Name:   meta.ProcessConfigTemplate,
				},
			},
		}

		return ps
	}

	// update a process config template.
	if ps.hitRegexp(updateProcConfigTemplateRegexp, http.MethodPut) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update process config template, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}

		templateID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update process config template, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.Update,
					Name:       meta.ProcessConfigTemplate,
					InstanceID: templateID,
				},
			},
		}

		return ps
	}

	// find processes's config template with condition.
	if ps.hitRegexp(findProcConfigTemplatesRegexp, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find process config templates, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.FindMany,
					Name:   meta.ProcessConfigTemplate,
				},
			},
		}

		return ps
	}

	// delete process config template
	if ps.hitRegexp(deleteProcConfigTemplateRegexp, http.MethodDelete) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete process config templates, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}

		templateID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete process config template, but got invalid template id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.Delete,
					Name:       meta.ProcessConfigTemplate,
					InstanceID: templateID,
				},
			},
		}
		return ps
	}

	// get process config template version
	if ps.hitRegexp(findProcessTemplateVersionRegexp, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get process config templates version, but got invalid business id: %s", ps.RequestCtx.Elements[6])
			return ps
		}

		templateID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get process config template version, but got invalid template id: %s", ps.RequestCtx.Elements[7])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.FindMany,
					Name:       meta.ProcessConfigTemplateVersion,
					InstanceID: templateID,
				},
			},
		}

		return ps
	}

	// create process template version
	if ps.hitRegexp(createProcessTemplateVersionRegexp, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create process config templates version, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.Create,
					Name:   meta.ProcessConfigTemplateVersion,
				},
			},
		}

		return ps
	}

	// update process template version
	if ps.hitRegexp(updateProcessTemplateVersionRegexp, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update process config templates version, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		versionID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update process config template version, but got invalid version id: %s", ps.RequestCtx.Elements[7])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.Create,
					Name:       meta.ProcessConfigTemplateVersion,
					InstanceID: versionID,
				},
			},
		}

		return ps
	}

	// preview process config
	if ps.hitRegexp(previewProcessConfigRegexp, http.MethodGet) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("preview process config template, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		templateID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("preview process config template, but got invalid template id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.Find,
					Name:       meta.ProcessConfigTemplate,
					InstanceID: templateID,
				},
			},
		}

		return ps
	}
	return ps
}

var (
	findProcBoundConfigRegexp             = regexp.MustCompile(`^/api/v3/template/proc/[^\s/]+/[0-9]+/[0-9]+/?$`)
	boundTemplateToProcessRegexp          = regexp.MustCompile(`^/api/v3/template/proc/[^\s/]+/[0-9]+/[0-9]+/[0-9]+/?$`)
	unboundTemplateWithProcessRegexp      = regexp.MustCompile(`^/api/v3/template/proc/[^\s/]+/[0-9]+/[0-9]+/[0-9]+/?$`)
	unboundTemplateWithProcessBatchRegexp = regexp.MustCompile(`^/api/v3/template/proc/[^\s/]+/[0-9]+/?$`)
)

func (ps *parseStream) processTemplateBound() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find the bounded process template config content.
	if ps.hitRegexp(findProcBoundConfigRegexp, http.MethodGet) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find bound process config templates, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		procID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find bound process config template, but got invalid process id: %s", ps.RequestCtx.Elements[6])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Process,
					Action:     meta.Find,
					Name:       meta.ProcessBoundConfig,
					InstanceID: procID,
				},
			},
		}

		return ps
	}

	// bound a template to a process
	if ps.hitRegexp(boundTemplateToProcessRegexp, http.MethodPut) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("bound process config templates, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.Create,
					Name:   meta.ProcessBoundConfig,
				},
			},
		}

		return ps
	}

	// unbound a template to a process
	if ps.hitRegexp(unboundTemplateWithProcessRegexp, http.MethodDelete) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("unbound process config templates, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.Delete,
					Name:   meta.ProcessBoundConfig,
				},
			},
		}

		return ps
	}

	// unbound template with a process batch.
	if ps.hitRegexp(unboundTemplateWithProcessBatchRegexp, http.MethodDelete) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("unbound process config templates batch, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Process,
					Action: meta.DeleteMany,
					Name:   meta.ProcessBoundConfig,
				},
			},
		}

		return ps
	}

	return ps
}

const (
	createServiceInstanceWithTempPattern          = "/process/v3/create/proc/service_instance/with_template"
	createServiceInstanceWithRawPattern           = "/process/v3/create/proc/service_instance/with_raw"
	findServiceInstancePattern                    = "/process/v3/find/proc/service_instance"
	deleteServiceInstancePattern                  = "/process/v3/delete/proc/service_instance"
	findServiceInstanceDifferencePattern          = "/process/v3/find/proc/service_instance/difference"
	syncServiceInstanceAccordingToServiceTemplate = "/process/v3/update/proc/service_instance/with_template"
	listServiceInstanceWithHostPattern            = "/process/v3/findmany/proc/service_instance/with_host"
)

var deleteProcessInstanceInServiceInstanceRegexp = regexp.MustCompile(`/process/v3/delete/proc/service_instance/[0-9]+/process/?$`)

func (ps *parseStream) ServiceInstance() *parseStream {
	if ps.hitPattern(createServiceInstanceWithTempPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Update,
				},
				BusinessID: bizID,
			},
		}

		return ps
	}

	if ps.hitPattern(syncServiceInstanceAccordingToServiceTemplate, http.MethodPut) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ids := gjson.GetBytes(ps.RequestCtx.Body, "service_instances").Array()
		for _, id := range ids {
			serviceInstanceID := id.Int()
			if serviceInstanceID <= 0 {
				ps.err = errors.New("invalid service instance id")
				return ps
			}
			ps.Attribute.Resources = append(ps.Attribute.Resources, meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Update,
				},
				BusinessID: bizID,
			})

		}

		return ps
	}

	if ps.hitPattern(listServiceInstanceWithHostPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Find,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	if ps.hitPattern(createServiceInstanceWithRawPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Create,
				},
				BusinessID: bizID,
			},
		}

		return ps
	}

	if ps.hitPattern(findServiceInstancePattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Find,
				},
				BusinessID: bizID,
			},
		}

		return ps
	}

	if ps.hitPattern(findServiceInstanceDifferencePattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Find,
				},
				BusinessID: bizID,
			},
		}

		return ps
	}

	if ps.hitPattern(deleteServiceInstancePattern, http.MethodDelete) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		instanceID := gjson.GetBytes(ps.RequestCtx.Body, "id").Int()
		if instanceID <= 0 {
			ps.err = errors.New("invalid service instance id")
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       meta.ProcessServiceInstance,
					Action:     meta.Delete,
					InstanceID: instanceID,
				},
				BusinessID: bizID,
			},
		}

		return ps
	}

	if ps.hitRegexp(deleteProcessInstanceInServiceInstanceRegexp, http.MethodDelete) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Delete,
				},
				BusinessID: bizID,
			},
		}

		return ps
	}

	return ps
}

const (
	findmanyServiceCategoryPattern = "/process/v3/findmany/proc/service_category"
	createServiceCategoryPattern   = "/process/v3/create/proc/service_category"
	deleteServiceCategoryPattern   = "/process/v3/delete/proc/service_category"
)

func (ps *parseStream) ServiceCategory() *parseStream {
	if ps.hitPattern(findmanyServiceCategoryPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceCategory,
					Action: meta.FindMany,
				},
				BusinessID: bizID,
			},
		}

		return ps
	}

	if ps.hitPattern(createServiceCategoryPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceCategory,
					Action: meta.Create,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	if ps.hitPattern(deleteServiceCategoryPattern, http.MethodDelete) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}
		categoryID := gjson.GetBytes(ps.RequestCtx.Body, "id").Int()
		if categoryID <= 0 {
			ps.err = errors.New("invalid category id")
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       meta.ProcessServiceCategory,
					Action:     meta.Delete,
					InstanceID: categoryID,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	return ps
}

const (
	createServiceTemplatePattern     = "/process/v3/create/proc/service_template"
	listServiceTemplatePattern       = "/process/v3/findmany/proc/service_template"
	listServiceTemplateDetailPattern = "/process/v3/findmany/proc/service_template/with_detail"
	deleteServiceTemplatePattern     = "/process/v3/delete/proc/service_template"
)

func (ps *parseStream) ServiceTemplate() *parseStream {
	if ps.hitPattern(createServiceTemplatePattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceTemplate,
					Action: meta.Create,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	if ps.hitPattern(listServiceTemplatePattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceTemplate,
					Action: meta.FindMany,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	if ps.hitPattern(listServiceTemplateDetailPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessServiceTemplate,
					Action: meta.FindMany,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	if ps.hitPattern(deleteServiceTemplatePattern, http.MethodDelete) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		templateID := gjson.GetBytes(ps.RequestCtx.Body, common.BKServiceTemplateIDField).Int()
		if templateID <= 0 {
			ps.err = errors.New("invalid service template ")
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       meta.ProcessServiceTemplate,
					Action:     meta.Delete,
					InstanceID: templateID,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	return ps
}

const (
	createProcessTemplateBatchPattern = "/process/v3/createmany/proc/proc_template/for_service_template"
	updateProcessTemplatePattern      = "/process/v3/update/proc/proc_template/for_service_template"
	deleteProcessTemplateBatchPattern = "/process/v3/deletemany/proc/proc_template/for_service_template"
	findProcessTemplateBatchPattern   = "/process/v3/findmany/proc/proc_template"
)

var findProcessTemplateRegexp = regexp.MustCompile(`/process/v3/find/proc/proc_template/id/[0-9]+/?$`)

func (ps *parseStream) ProcessTemplate() *parseStream {
	if ps.hitPattern(createProcessTemplateBatchPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessTemplate,
					Action: meta.Create,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	if ps.hitPattern(updateProcessTemplatePattern, http.MethodPut) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}
		procTemplateID := gjson.GetBytes(ps.RequestCtx.Body, common.BKProcessTemplateIDField).Int()
		if procTemplateID <= 0 {
			ps.err = errors.New("invalid process template id")
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       meta.ProcessTemplate,
					Action:     meta.Update,
					InstanceID: procTemplateID,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	if ps.hitPattern(deleteProcessTemplateBatchPattern, http.MethodDelete) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ProcessTemplate,
					Action: meta.Delete,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	if ps.hitPattern(findProcessTemplateBatchPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		procTemplateIDs := gjson.GetBytes(ps.RequestCtx.Body, "process_template_ids").Array()
		for _, id := range procTemplateIDs {
			procTemplateID := id.Int()
			if procTemplateID <= 0 {
				ps.err = errors.New("invalid process template id")
				return ps
			}

			ps.Attribute.Resources = append(ps.Attribute.Resources, meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       meta.ProcessTemplate,
					Action:     meta.Find,
					InstanceID: procTemplateID,
				},
				BusinessID: bizID,
			})
		}

		return ps
	}

	if ps.hitRegexp(findProcessTemplateRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
			ps.err = err
			return ps
		}

		procTemplateID := gjson.GetBytes(ps.RequestCtx.Body, "id").Int()
		if procTemplateID <= 0 {
			ps.err = errors.New("invalid process template id")
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       meta.ProcessTemplate,
					Action:     meta.Find,
					InstanceID: procTemplateID,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	return ps
}
