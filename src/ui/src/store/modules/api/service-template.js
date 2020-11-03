/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import $http from '@/api'

const state = {}

const getters = {}

const actions = {
    /**
     * 查询服务模板详情列表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchServiceTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`findmany/proc/service_template/with_detail`, params, config)
    },
    /**
     * 查询服务模板详情列表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchServiceTemplateWithoutDetails ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`findmany/proc/service_template`, params, config)
    },
    /**
     * 创建服务模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createServiceTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`create/proc/service_template`, params, config)
    },
    /**
     * 编辑服务模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateServiceTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.put(`update/proc/service_template`, params, config)
    },
    /**
     * 删除服务模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteServiceTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.delete(`delete/proc/service_template`, params)
    },
    /**
     * 查看单个服务模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    findServiceTemplate ({ commit, state, dispatch, rootGetters }, { id, config }) {
        return $http.get(`find/proc/service_template/${id}/detail`, config)
    },

    getServiceTemplateModules (context, { bizId, serviceTemplateId, params, config }) {
        return $http.post(`module/bk_biz_id/${bizId}/service_template_id/${serviceTemplateId}`, params, config)
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
