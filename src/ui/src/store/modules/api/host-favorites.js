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
     * 新加收藏
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    createFavorites ({ commit, state, dispatch }, { params }) {
        return $axios.post(`hosts/favorites`, params)
    },

    /**
     * 编辑收藏
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} id 收藏的主键
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    udpateFavorites ({ commit, state, dispatch }, { id, params }) {
        return $axios.put(`hosts/favorites/${id}`, params)
    },

    /**
     * 删除收藏
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} id 收藏的主键
     * @return {Promise} promise 对象
     */
    deleteFavorites ({ commit, state, dispatch }, { id }) {
        return $axios.delete(`hosts/favorites/${id}`)
    },

    /**
     * 收藏使用次数加一
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} id 收藏的主键
     * @return {Promise} promise 对象
     */
    incrFavorites ({ commit, state, dispatch }, { id }) {
        return $axios.put(`hosts/favorites/${id}/incr`)
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
