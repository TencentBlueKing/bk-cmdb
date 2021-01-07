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
     * 为服务模板新增进程
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createProcessTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`createmany/proc/proc_template`, params, config)
    },
    /**
     * 批量查询进程模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    getBatchProcessTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`findmany/proc/proc_template`, params, config)
    },
    /**
     * 查询进程模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    getProcessTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`find/proc/proc_template/id/${params.processTemplateId}`, {}, config)
    },
    /**
     * 更新服务模板中的进程模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateProcessTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.put(`update/proc/proc_template`, params, config)
    },
    /**
     * 删除服务模板中的进程模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteProcessTemplate ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.delete(`deletemany/proc/proc_template`, params, config)
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
