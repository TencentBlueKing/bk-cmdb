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
  diffTemplateAndInstances({ commit, state, dispatch, rootGetters }, { bizId, setTemplateId, params, config }) {
    return $http.post(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/diff_with_instances`, params, config)
  },
  syncTemplateToInstances({ commit }, { bizId, setTemplateId, params, config }) {
    return $http.post(`updatemany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/sync_to_instances`, params, config)
  },
  getInstancesSyncStatus({ commit }, { bizId, setTemplateId, params, config }) {
    return $http.post(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/instances_sync_status`, params, config)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
