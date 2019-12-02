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
// import jsCookie from 'js-cookie'

const state = {
    business: [],
    bizId: null,
    authorizedBusiness: []
}

const getters = {
    business: state => state.business,
    bizId: state => state.bizId,
    currentBusiness: state => state.authorizedBusiness.find(business => business.bk_biz_id === state.bizId),
    authorizedBusiness: state => state.authorizedBusiness
}

const actions = {
    getAuthorizedBusiness ({ commit }, config = {}) {
        return $http.get('biz/with_reduced', config).then(data => {
            commit('setAuthorizedBusiness', data.info)
            return data.info
        })
    },
    /**
     * 添加业务
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createBusiness ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`biz/${rootGetters.supplierAccount}`, params, config)
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
        return $http.delete(`biz/${bkSupplierAccount}/${bkBizId}`)
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
    updateBusiness ({ commit, state, dispatch, rootGetters }, { bizId, params, config }) {
        return $http.put(`biz/${rootGetters.supplierAccount}/${bizId}`, params, config)
    },

    /**
     * 归档业务
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    archiveBusiness ({ commit, state, dispatch, rootGetters }, bizId) {
        return $http.put(`biz/status/disabled/${rootGetters.supplierAccount}/${bizId}`)
    },

    /**
     * 恢复业务
     * @param {Function} commit store commit mutation hander
     * @param {Number} bizId 业务id
     * @return {promises} promises 对象
     */
    recoveryBusiness ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.put(`biz/status/enable/${rootGetters.supplierAccount}/${params['bk_biz_id']}`, {}, config)
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
    searchBusiness ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`${window.API_HOST}biz/search/web`, params, config)
    },

    searchBusinessById ({ rootGetters }, { bizId, config }) {
        return $http.post(`biz/search/${rootGetters.supplierAccount}`, {
            condition: {
                'bk_biz_id': {
                    '$eq': bizId
                }
            },
            fields: [],
            page: {
                start: 0,
                limit: 1
            }
        }, config).then(data => {
            return data.info[0] || {}
        })
    },
    getFullAmountBusiness ({ commit }, config = {}) {
        return $http.get('biz/simplify', config)
    }
}

const mutations = {
    setBusiness (state, business) {
        state.business = business
    },
    setBizId (state, bizId) {
        state.bizId = bizId
    },
    setAuthorizedBusiness (state, list) {
        state.authorizedBusiness = list
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
