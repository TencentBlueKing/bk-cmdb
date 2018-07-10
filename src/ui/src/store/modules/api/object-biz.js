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
     * 添加业务
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createBusiness ({ commit, state, dispatch }, { bkSupplierAccount, params }) {
        return $axios.post(`biz/${bkSupplierAccount}`, params)
    },

    /**
     * 删除业务
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Number} bkBizId 业务id
     * @return {promises} promises 对象
     */
    deleteBusiness ({ commit, state, dispatch }, { bkSupplierAccount, bkBizId }) {
        return $axios.delete(`biz/${bkSupplierAccount}/${bkBizId}`)
    },

    /**
     * 修改业务
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Number} bkBizId 业务id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateBusiness ({ commit, state, dispatch }, { bkSupplierAccount, bkBizId, params }) {
        return $axios.put(`biz/${bkSupplierAccount}/${bkBizId}`, params)
    },

    /**
     * 查询业务
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchBusiness ({ commit, state, dispatch }, { bkSupplierAccount, params }) {
        return $axios.post(`biz/search/${bkSupplierAccount}`, params)
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
