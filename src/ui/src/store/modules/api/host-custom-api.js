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
     * 新加自定义API
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    createCustomQuery ({ commit, state, dispatch }, { params }) {
        return $axios.post(`userapi`, params)
    },
    
    /**
     * 更新自定义API
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务Id
     * @param {String} id 主键id
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    updateCustomQuery ({ commit, state, dispatch }, { bkBizId, id, params }) {
        return $axios.post(`userapi/${bkBizId}/${id}`, params)
    },

    /**
     * 删除自定义API
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务Id
     * @param {String} id 主键id
     * @return {Promise} promise 对象
     */
    deleteCustomQuery ({ commit, state, dispatch }, { bkBizId, id }) {
        return $axios.delete(`userapi/${bkBizId}/${id}`)
    },

    /**
     * 查询自定义API
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务Id
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    searchCustomQuery ({ commit, state, dispatch }, { bkBizId, params }) {
        return $axios.post(`userapi/search/${bkBizId}`, params)
    },

    /**
     * 获取自定义API详情
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务Id
     * @param {String} id 主键id
     * @return {Promise} promise 对象
     */
    getCustomQueryDetail ({ commit, state, dispatch }, { bkBizId, id }) {
        return $axios.get(`userapi/detail/${bkBizId}/${id}`)
    },

    /**
     * 根据自定义API获取数据
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkBizId 业务Id
     * @param {String} id 主键id
     * @param {Number} start 记录开始位置
     * @param {Number} limit 每夜限制条数，最大200
     * @return {Promise} promise 对象
     */
    getCustomQueryData ({ commit, state, dispatch }, { bkBizId, id, start, limit, params }) {
        return $axios.get(`userapi/detail/${bkBizId}/${id}/${start}/${limit}`)
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
