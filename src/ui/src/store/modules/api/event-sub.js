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
     * 订阅事件
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Object} bkBizId 业务id
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    subscribeEvent ({ commit, state, dispatch }, {bkSupplierAccount, bkBizId, params}) {
        return $axios.post(`event/subscribe/${bkSupplierAccount}/${bkBizId}`, params)
    },

    /**
     * 退订事件
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Number} bkBizId 业务id
     * @param {Number} subscriptionId 订阅id
     * @return {Promise} promise 对象
     */
    unsubcribeEvent ({commit, state, dispatch}, {bkSupplierAccount, bkBizId, subscriptionId}) {
        return $axios.delete(`event/subscribe/${bkSupplierAccount}/${bkBizId}/${subscriptionId}`)
    },

    /**
     * 修改订阅
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Number} bkBizId 业务id
     * @param {Number} subscriptionId 订阅id
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    updateEventSubscribe ({commit, state, dispatch}, {bkSupplierAccount, bkBizId, subscriptionId, params}) {
        return $axios.put(`event/subscribe/${bkSupplierAccount}/${bkBizId}/${subscriptionId}`, params)
    },

    /**
     * 查询订阅
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Number} bkBizId 业务id
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    searchSubscription ({commit, state, dispatch}, {bkSupplierAccount, bkBizId, params}) {
        return $axios.post(`event/subscribe/search/${bkSupplierAccount}/${bkBizId}`, params)
    },

    /**
     * 测试推送
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    pingSubscription ({ commit, state, dispatch }, params) {
        return $axios.post(`event/subscribe/ping`, params)
    },

    /**
     * 测试推送（只测试连通性）
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    testingConnection ({ commit, state, dispatch }, params) {
        return $axios.post(`event/subscribe/telnet`, params)
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
