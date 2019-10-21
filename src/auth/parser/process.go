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
)

func (ps *parseStream) processRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.process().
		ServiceInstance().
		ServiceTemplate().
		ServiceCategory().
		processTemplate().
		ProcessTemplate().
		ProcessInstance().
		OperationStatistic()
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
