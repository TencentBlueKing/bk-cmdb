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
    createObject ({ commit, state, dispatch }, { params }) {
        return $axios.post(`object`, params)
    },

    /**
     * 删除对象模型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 被删除的数据记录的id
     * @return {promises} promises 对象
     */
    deleteObject ({ commit, state, dispatch }, { id }) {
        return $axios.delete(`object/${id}`)
    },

    /**
     * 更新对象模型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 目标数据的记录id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateObject ({ commit, state, dispatch }, { id, params }) {
        return $axios.put(`object/${id}`, params)
    },

    /**
     * 查询对象模型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchObjects ({ commit, state, dispatch }, { params }) {
        return $axios.post(`objects`, params)
    },

    /**
     * 查询普通对象模型的拓扑结构
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchObjectTopo ({ commit, state, dispatch }, { params }) {
        return $axios.post(`objects/topo`, params)
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
