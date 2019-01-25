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
     *  采集器查询
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    searchDataCollection ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/collector/action/search`, params, config)
    },
    /**
     *  采集器配置保存
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    updateDataCollection ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/collector/action/update`, params, config)
    },
    /**
     *  采集器执行发现
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    collectDataCollection ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/collector/action/discover`, params, config)
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
