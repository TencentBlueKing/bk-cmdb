/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

/* eslint-disable no-unused-vars */
import $http from '@/api'
import CombineRequest from '@/api/combine-request'

const state = {
  propertyConfig: {},
  ruleDraft: {},
  propertyList: []
}

const getters = {
  configPropertyList: (state) => {
    state.propertyList.forEach((property) => {
      // 兼容通用方法
      property.options = property.option
      // 自定义字段空间
      // eslint-disable-next-line no-underscore-dangle
      property.__extra__ = {
        visible: true
      }
    })
    const enabledList = state.propertyList.filter(item => item.host_apply_enabled)
    const disabledList = state.propertyList.filter(item => !item.host_apply_enabled)
    return [...enabledList, ...disabledList]
  }
}

const actions = {
  async getRules({ commit, state, dispatch }, { bizId, params, config }) {
    const reqs = await CombineRequest.setup(Symbol(), (ids) => {
      const subParams = { ...params, bk_module_ids: ids }
      return $http.post(`findmany/host_apply_rule/bk_biz_id/${bizId}`, subParams, config)
    }, { segment: 500 }).add(params.bk_module_ids)

    const results = []
    const res = await Promise.all(reqs)
    res.forEach(({ info: list = [] }) => {
      results.push(...list)
    })

    return {
      count: results.length,
      info: results
    }
  },
  async getTemplateRules({ commit, state, dispatch }, { bizId, params, config }) {
    const reqs = await CombineRequest.setup(Symbol(), (ids) => {
      const subParams = { ...params, service_template_ids: ids }
      return $http.post('host/findmany/service_template/host_apply_rule', subParams, config)
    }, { segment: 500 }).add(params.service_template_ids)

    const results = []
    const res = await Promise.all(reqs)
    res.forEach(({ info: list = [] }) => {
      results.push(...list)
    })

    return {
      count: results.length,
      info: results
    }
  },
  getApplyPreview({ commit, state, dispatch }, { params, config }) {
    return $http.post('host/createmany/module/host_apply_plan/preview', params, config)
  },
  getTemplateApplyPreview({ commit, state, dispatch }, { params, config }) {
    return $http.post('host/createmany/service_template/host_apply_plan/preview', params, config)
  },
  runApply({ commit, state, dispatch }, { params, config }) {
    return $http.post('host/updatemany/module/host_apply_plan/run', params, config)
  },
  runTemplateApply({ commit, state, dispatch }, { params, config }) {
    return $http.post('updatemany/proc/service_template/host_apply_plan/run', params, config)
  },
  getApplyTaskStatus({ commit, state, dispatch }, { params, config }) {
    return $http.post('host/findmany/module/host_apply_plan/status', params, config)
  },
  getTemplateApplyTaskStatus({ commit, state, dispatch }, { params, config }) {
    return $http.post('findmany/proc/service_template/host_apply_plan/status', params, config)
  },
  getTopopath({ commit, state, dispatch }, { bizId, params, config }) {
    return $http.post(`find/topopath/biz/${bizId}`, params, config)
  },
  setEnableStatus({ commit, state, dispatch }, { bizId, params, config }) {
    return $http.put(`module/host_apply_enable_status/bk_biz_id/${bizId}`, params, config)
  },
  setTemplateEnableStatus({ commit, state, dispatch }, { bizId, params, config }) {
    return $http.put(`updatemany/proc/service_template/host_apply_enable_status/biz/${bizId}`, params, config)
  },
  getConflictCount({ commit, state, dispatch }, { params, config }) {
    return $http.post('host/findmany/module/host_apply_plan/invalid_host_count', params, config)
  },
  getTemplateConflictCount({ commit, state, dispatch }, { params, config }) {
    return $http.post('host/findmany/service_template/host_apply_plan/invalid_host_count', params, config)
  },
  deleteRules({ commit, state, dispatch }, { bizId, params }) {
    return $http.delete(`host/deletemany/module/host_apply_rule/bk_biz_id/${bizId}`, params)
  },
  deleteTemplateRules({ commit, state, dispatch }, { bizId, params }) {
    return $http.delete(`deletemany/proc/service_template/host_apply_rule/biz/${bizId}`, params)
  },
  getProperties(context, { params, config }) {
    return $http.post('find/objectattr/host', params, config)
  },
  getHostRelatedRules(context, { bizId, params, config }) {
    return $http.post(`findmany/host_apply_rule/bk_biz_id/${bizId}/host_related_rules`, params, config)
  },
  searchNode(context, { bizId, params, config }) {
    return $http.post(`find/topoinst/bk_biz_id/${bizId}/host_apply_rule_related`, params, config)
  },
  getModuleApplyStatusByTemplate(context, { params, config }) {
    return $http.post('host/find/service_template/host_apply_status', params, config)
  },
  searchTemplateNode(context, { params, config }) {
    return $http.post('find/proc/service_template/host_apply_rule_related', params, config)
  },
  getTemplateRuleCount(context, { params, config }) {
    return $http.post('host/findmany/service_template/host_apply_rule_count', params, config)
  },
  getModuleFinalRules({ commit, state, dispatch }, { params, config }) {
    return $http.post('host/findmany/module/get_module_final_rules', params, config)
  }
}

const mutations = {
  setPropertyConfig(state, config) {
    state.propertyConfig = config
  },
  setRuleDraft(state, draft) {
    state.ruleDraft = draft
  },
  clearRuleDraft(state) {
    state.ruleDraft = {}
  },
  setPropertyList(state, list) {
    state.propertyList = list
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
