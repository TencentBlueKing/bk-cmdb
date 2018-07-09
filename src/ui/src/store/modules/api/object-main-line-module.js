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
     * 添加模型主关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createMainlineObject ({ commit, state, dispatch }, { params }) {
        return $axios.post(`topo/model/mainline`, params)
    },

    /**
     * 删除模型主关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkObjId 对象的模型id
     * @return {promises} promises 对象
     */
    deleteMainlineObject ({ commit, state, dispatch }, { bkSupplierAccount, bkObjId }) {
        return $axios.delete(`topo/model/mainline/owners/${bkSupplierAccount}/objectids/${bkObjId}`)
    },

    /**
     * 查询模型拓扑
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @return {promises} promises 对象
     */
    searchMainlineObject ({ commit, state, dispatch }, { bkSupplierAccount }) {
        return $axios.get(`topo/model/${bkSupplierAccount}`)
    },

    /**
     * 获取实例拓扑
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkBizId 业务id
     * @return {promises} promises 对象
     */
    getInstTopo ({ commit, state, dispatch }, { bkSupplierAccount, bkBizId }) {
        return $axios.get(`topo/inst/${bkSupplierAccount}/${bkBizId}`)
    },

    /**
     * 获取子节点实例
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkObjId 对象的模型id
     * @param {String} bkBizId 业务id
     * @param {String} bkInstId 实例id
     * @return {promises} promises 对象
     */
    searchInstTopo ({ commit, state, dispatch }, { bkSupplierAccount, bkObjId, bkBizId, bkInstId }) {
        return $axios.get(`topo/inst/child/${bkSupplierAccount}/${bkObjId}/${bkBizId}/${bkInstId}`)
    },

    /**
     * 查询内置模块集
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkBizId 业务id
     * @return {promises} promises 对象
     */
    getInternalTopo ({ commit, state, dispatch }, { bkSupplierAccount, bkBizId }) {
        return $axios.get(`topo/internal/${bkSupplierAccount}/${bkBizId}`)
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
