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
    cloudName: null
}

const getters = {
    cloudName: state => state.cloudName
}

const actions = {
    /**
     * 网络发现报告简要查询
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchNetcollect ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/summary/action/search`, params, config)
    },
    /**
     * 网络发现报告详情列表查询
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchNetcollectList ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/report/action/search`, params, config)
    },
    /**
     * 网络发现报告变更确认
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    confirmNetcollectChange ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/report/action/confirm`, params, config)
    },
    /**
     * 网络发现完成历史
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchNetcollectHistory ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/history/action/search`, params, config)
    }
}

const mutations = {
    setCloudName (state, cloudName) {
        state.cloudName = cloudName
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
