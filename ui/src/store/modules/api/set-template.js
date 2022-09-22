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

const state = {}

const getters = {}

const mutations = {}

const actions = {
  createSetTemplate({ commit }, { bizId, params, config }) {
    return $http.post(`create/topo/set_template/bk_biz_id/${bizId}`, params, config)
  },
  updateSetTemplate({ commit }, { bizId, setTemplateId, params, config }) {
    return $http.put(`update/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}`, params, config)
  },
  deleteSetTemplate({ commit }, { bizId, config }) {
    return $http.delete(`deletemany/topo/set_template/bk_biz_id/${bizId}`, config)
  },
  getSetTemplates({ commit }, { bizId, params, config }) {
    return $http.post(`findmany/topo/set_template/bk_biz_id/${bizId}/web`, params, config)
  },
  getSingleSetTemplateInfo({ commit }, { bizId, setTemplateId, config }) {
    return $http.get(`find/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}`, config)
  },
  getSetTemplateServices({ commit }, { bizId, setTemplateId }, config) {
    return $http.get(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/service_templates`)
  },
  getSetInstancesWithStatus({ commit }, { bizId, params, config }) {
    return $http.post(`findmany/topo/set_template_sync_status/bk_biz_id/${bizId}`, params, config)
  },
  getSetInstancesWithTopo({ commit }, { bizId, setTemplateId, params, config }) {
    return $http.post(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/sets/web`, params, config)
  },
  getSyncHistory({ commit }, { bizId, params, config }) {
    return $http.post(`findmany/topo/set_template_sync_history/bk_biz_id/${bizId}`, params, config)
  },
  getSetTemplateServicesStatistics({ commit }, { bizId, setTemplateId, config }) {
    return $http.get(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/service_templates/with_statistics`, config)
  },
  getSetTemplateStatus(context, { bizId, params, config }) {
    return $http.post(`findmany/topo/set_template/bk_biz_id/${bizId}/set_template_status`, params, config)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
