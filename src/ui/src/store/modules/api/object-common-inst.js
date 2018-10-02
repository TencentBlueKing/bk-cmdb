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
     * 添加对象实例
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} objId 模型id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createInst ({ commit, state, dispatch, rootGetters }, { objId, params, config }) {
        return $http.post(`inst/${rootGetters.supplierAccount}/${objId}`, params, config)
    },

    /**
     * 更新对象实例
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} objId 模型id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateInst ({ commit, state, dispatch, rootGetters }, { objId, instId, params, config }) {
        return $http.put(`inst/${rootGetters.supplierAccount}/${objId}/${instId}`, params, config)
    },

    /**
     * 批量更新对象实例
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} objId 模型id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    batchUpdateInst ({ commit, state, dispatch, rootGetters }, { objId, params, config }) {
        return $http.put(`inst/${rootGetters.supplierAccount}/${objId}/batch`, params, config)
    },
    
    /**
     * 查询实例
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} objId 模型id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchInst ({ commit, state, dispatch, rootGetters }, { params, config, objId }) {
        return $http.post(`inst/association/search/owner/${rootGetters.supplierAccount}/object/${objId}`, params, config)
    },

    searchInstById ({ rootGetters }, { config, objId, instId, idKey = 'bk_inst_id' }) {
        return $http.post(`inst/association/search/owner/${rootGetters.supplierAccount}/object/${objId}`, {
            condition: {
                [objId]: [{
                    field: idKey,
                    operator: '$eq',
                    value: instId
                }]
            },
            fields: {},
            page: {
                start: 0,
                limit: 1
            }
        }, config).then(data => {
            return data.info[0] || {}
        })
    },

    /**
     * 删除对象实例
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} objId 模型id
     * @param {String} inst 实例id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteInst ({ commit, state, dispatch, rootGetters }, { objId, instId, config }) {
        return $http.delete(`inst/${rootGetters.supplierAccount}/${objId}/${instId}`, config)
    },

    /**
     * 批量删除对象实例
     * @param {String} objId 模型id
     * @param {Object} config 参数
     * @return {promises} promises 对象
     */
    batchDeleteInst ({ commit, state, dispatch, rootGetters }, {objId, config}) {
        return $http.delete(`inst/${rootGetters.supplierAccount}/${objId}/batch`, config)
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
