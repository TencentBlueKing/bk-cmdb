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

const state = {

}

const getters = {

}

const actions = {
    /**
     * 新增主机
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    addHostToResource ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/add`, params, config)
    },

    /**
     * 主机转移到业务内模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    transferHostModule ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/modules`, params, config)
    },

    /**
     * 资源池主机分配至业务的空闲机模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    transferResourcehostToIdleModule ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/modules/resource/idle`, params, config)
    },

    /**
     * 主机上交至业务的故障机模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    transferHostToFaultModule ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/modules/fault`, params, config)
    },

    /**
     * 主机上交至业务的空闲机模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    transferHostToIdleModule ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/modules/idle`, params, config)
    },

    /**
     * 主机回收至资源池
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    transferHostToResourceModule ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/modules/resource`, params, config)
    },

    /**
     * 转移主机至模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    transferHostToMutipleBizModule ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/modules/biz/mutilple`, params, config)
    }
}

const mutations = {

}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
