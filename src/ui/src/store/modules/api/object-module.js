/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import { $Axios, $axios } from '@/api/axios'

const state = {

}

const getters = {

}

const actions = {
    /**
     * 创建模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务id
     * @param {Number} bkSetId 集群id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createModule ({ commit, state, dispatch }, { bkBizId, bkSetId, params }) {
        return $axios.post(`module/${bkBizId}/${bkSetId}`, params)
    },

    /**
     * 删除模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务id
     * @param {Number} bkSetId 集群id
     * @param {Number} bkModuleId 模块id
     * @return {promises} promises 对象
     */
    deleteModule ({ commit, state, dispatch }, { bkBizId, bkSetId, bkModuleId }) {
        return $axios.delete(`module/${bkBizId}/${bkSetId}/${bkModuleId}`)
    },

    /**
     * 更新模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务id
     * @param {Number} bkSetId 集群id
     * @param {Number} bkModuleId 模块id
     * @return {promises} promises 对象
     */
    updateModule ({ commit, state, dispatch }, { bkBizId, bkSetId, bkModuleId, params }) {
        return $axios.put(`module/${bkBizId}/${bkSetId}/${bkModuleId}`, params)
    },

    /**
     * 查询模块
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 模块id
     * @param {Number} bkBizId 业务id
     * @param {Number} bkSetId 集群id
     * @return {promises} promises 对象
     */
    searchModule ({ commit, state, dispatch }, { bkSupplierAccount, bkBizId, bkSetId, params }) {
        return $axios.post(`module/${bkSupplierAccount}/${bkBizId}/${bkSetId}`, params)
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
