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

/* eslint-disable no-unused-vars, max-len */

import $http from '@/api'

const state = {}

const getters = {}

const actions = {
  /**
     * 根据服务模板、模块查询进程实例与服务模板之间的差异
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  searchServiceInstanceDifferences({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.post('find/proc/service_instance/difference', params, config)
  },
  /**
   * 获取进程模板 diff 状态
   */
  getProcessTplDiffState({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.post(`findmany/topo/service_template_sync_status/bk_biz_id/${params.bk_biz_id}`, params, config)
  },
  /**
   * 获取服务模板diff信息
   */
  getTplDiffs({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.post('/find/proc/service_template/general_difference', params, config)
  },
  /**
   * 获取进程下涉及到变更的实例
   * @param {number} params.bk_biz_id 业务 ID
   * @param {number} params.service_template_id 服务模板 ID
   * @param {number} params.process_template_id 进程模板 ID
   * @param {number} params.bk_module_id 模块 ID
   * @param {string} [params.process_template_name] 进程模板名称。因为被删除的模板 ID 都为 0，所以需要进程模板名称配合来获取对应进程下的实例。
   * @param {boolean} [params.service_category] 是否是服务分类变更
   * @returns {Promise}
   */
  getDiffInstances({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.post('/find/proc/difference/service_instances', params, config)
  },
  /**
   * 获取单个实例的对比详情
   * @param {number} params.bk_biz_id 业务 ID
   * @param {number} params.service_template_id 服务模板 ID
   * @param {number} params.process_template_id 进程模板 ID
   * @param {number} params.bk_module_id 模块 ID
   * @param {number} params.service_instance_id 服务实例 ID
   * @returns {Promise}
   */
  getInstanceDiff({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.post('/find/proc/service_instance/difference_detail', params, config)
  },
  /**
     * 批量更新服务实例中的进程信息，保持和服务模板一致
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  syncServiceInstanceByTemplate({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.put('update/proc/service_instance/sync', params, config)
  }
}

const mutations = {}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
