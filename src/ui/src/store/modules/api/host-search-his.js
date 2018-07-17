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
     * 新加主机查询历史
     * @param {Function} comit store commit mutation hander
     * @param {Object} state store state
     * @param {string} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    createHostSearchHistory ({ commit, state, dispatch }, { params }) {
        return $axios.post(`hosts/history`, params)
    },

    /**
     * 获取主机查询历史
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store disoatch action hander
     * @param {Number} start 记录开始位置
     * @param {Number} limit 每页限制条数，最大200
     * @return {Promise} promise 对象
     */
    searchHostSearchHistory ({ commit, state, dispatch }, { start, limit }) {
        return $axios.get(`hosts/history${start}/${limit}`)
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
