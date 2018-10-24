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
     * 根据条件查询主机
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    searchHost ({ commit, state, dispatch }, {params, config}) {
        return $http.post(`hosts/search`, params, config)
    },

    searchHostByInnerip (context, { bizId, innerip, config }) {
        return $http.post(`hosts/search`, {
            'bk_biz_id': bizId,
            condition: ['biz', 'set', 'module', 'host'].map(model => {
                return {
                    'bk_obj_id': model,
                    condition: []
                }
            }),
            ip: {
                flag: 'bk_host_innerip',
                exact: 1,
                data: [innerip]
            },
            page: {
                start: 0,
                limit: 1
            }
        }, config).then(data => {
            return data.info[0] || {}
        })
    },

    /**
     * 获取主机详情
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Number} bkHostId 主机id
     * @return {Promise} promise 对象
     */
    getHostBaseInfo ({ commit, state, dispatch, rootGetters }, { hostId, config }) {
        return $http.get(`hosts/${rootGetters.supplierAccount}/${hostId}`)
    },

    /**
     * 根据主机id获取主机快照数据
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkHostId 主机id
     * @return {Promise} promise 对象
     */
    getHostSnapshot ({ commit, state, dispatch }, { hostId }) {
        return $http.get(`hosts/snapshot/${hostId}`)
    },

    /**
     * 根据主机id获取主机快照数据
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    searchHostByCondition ({ commit, state, dispatch }, { params }) {
        return $http.post(`hosts/snapshot/asstdetail`, params)
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
