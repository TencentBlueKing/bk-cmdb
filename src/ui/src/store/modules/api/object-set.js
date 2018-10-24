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
     * 创建集群
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createSet ({ commit, state, dispatch }, { bizId, params, config }) {
        return $http.post(`set/${bizId}`, params, config)
    },

    /**
     * 删除集群
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务id
     * @param {Number} bkSetId 集群id
     * @return {promises} promises 对象
     */
    deleteSet ({ commit, state, dispatch }, { bizId, setId, config }) {
        return $http.delete(`set/${bizId}/${setId}`, config)
    },

    /**
     * 更新集群
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务id
     * @param {Number} bkSetId 集群id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateSet ({ commit, state, dispatch }, { bizId, setId, params, config }) {
        return $http.put(`set/${bizId}/${setId}`, params, config)
    },

    /**
     * 查询集群
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Number} bkBizId 业务id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchSet ({ commit, state, dispatch, rootGetters }, { bizId, params, config }) {
        return $http.post(`set/search/${rootGetters.supplierAccount}/${bizId}`, params, config)
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
