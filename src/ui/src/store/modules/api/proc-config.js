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
     * 新增进程
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createProcess ({ commit, state, dispatch, rootGetters }, { bizId, params }) {
        return $http.post(`proc/${rootGetters.supplierAccount}/${bizId}`, params)
    },

    /**
     * 查询进程
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchProcess ({ commit, state, dispatch, rootGetters }, { bizId, params, config }) {
        return $http.post(`proc/search/${rootGetters.supplierAccount}/${bizId}`, params, config)
    },

    searchProcessById ({ rootGetters }, { bizId, processId, config }) {
        return $http.post(`proc/search/${rootGetters.supplierAccount}/${bizId}`, {
            condition: {
                'bk_biz_id': bizId,
                'bk_process_id': {
                    '$eq': processId
                }
            }
        }, config).then(data => {
            return data.info[0] || {}
        })
    },
    /**
     * 获取进程详情
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Number} processId 进程id
     * @return {promises} promises 对象
     */
    getProcessDetail ({ commit, state, dispatch, rootGetters }, { bizId, processId }) {
        return $http.get(`proc/${rootGetters.supplierAccount}/${bizId}/${processId}`)
    },

    /**
     * 删除进程
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Number} processId 进程id
     * @return {promises} promises 对象
     */
    deleteProcess ({ commit, state, dispatch, rootGetters }, { bizId, processId }) {
        return $http.delete(`proc/${rootGetters.supplierAccount}/${bizId}/${processId}`)
    },

    /**
     * 更新进程
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Number} processId 进程id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateProcess ({ commit, state, dispatch, rootGetters }, { bizId, processId, params }) {
        return $http.put(`proc/${rootGetters.supplierAccount}/${bizId}/${processId}`, params)
    },

    /**
     * 批量更新进程
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    batchUpdateProcess ({ commit, state, dispatch, rootGetters }, { bizId, params }) {
        return $http.put(`proc/${rootGetters.supplierAccount}/${bizId}`, params)
    },

    /**
     * 获取进程绑定模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Number} processId 进程id
     * @return {promises} promises 对象
     */
    getProcessBindModule ({ commit, state, dispatch, rootGetters }, { bizId, processId, config }) {
        return $http.get(`proc/module/${rootGetters.supplierAccount}/${bizId}/${processId}`, config)
    },

    /**
     * 绑定进程到模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Number} processId 进程id
     * @param {String} moduleName 模块名称
     * @return {promises} promises 对象
     */
    bindProcessModule ({ commit, state, dispatch, rootGetters }, { bizId, processId, moduleName, config }) {
        return $http.put(`proc/module/${rootGetters.supplierAccount}/${bizId}/${processId}/${moduleName}`, {}, config)
    },

    /**
     * 解绑进程模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bizId 业务id
     * @param {Number} processId 进程id
     * @param {String} moduleName 模块名称
     * @return {promises} promises 对象
     */
    deleteProcessModuleBinding ({ commit, state, dispatch, rootGetters }, { bizId, processId, moduleName, config }) {
        return $http.delete(`proc/module/${rootGetters.supplierAccount}/${bizId}/${processId}/${moduleName}`, config)
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
