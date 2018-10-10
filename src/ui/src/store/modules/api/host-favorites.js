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
    applying: null
}

const getters = {
    applying: state => state.applying,
    applyingInfo: state => {
        if (state.applying) {
            return JSON.parse(state.applying.info)
        }
        return null
    },
    applyingProperties: state => {
        const properties = []
        if (state.applying) {
            const ignore = ['biz']
            const params = JSON.parse(state.applying['query_params'])
            params.forEach(param => {
                if (!ignore.includes(param['bk_obj_id'])) {
                    properties.push({
                        'bk_obj_id': param['bk_obj_id'],
                        'bk_property_id': param.field
                    })
                }
            })
        }
        return properties
    },
    applyingConditions: state => {
        const conditions = {}
        if (state.applying) {
            const ignore = ['biz']
            const params = JSON.parse(state.applying['query_params'])
            params.forEach(param => {
                const objId = param['bk_obj_id']
                if (!ignore.includes(objId)) {
                    conditions[objId] = conditions[objId] || []
                    conditions[objId].push({
                        field: param.field,
                        operator: param.operator,
                        value: param.value
                    })
                }
            })
        }
        return conditions
    }
}

const actions = {
    /**
     * 搜索收藏
     * @param {Object} context store上下文
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    searchFavorites (context, { params, config }) {
        return $http.post('hosts/favorites/search', params, config)
    },
    /**
     * 新加收藏
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    createFavorites ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/favorites`, params, config)
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
    udpateFavorites ({ commit, state, dispatch }, { id, params, config }) {
        return $http.put(`hosts/favorites/${id}`, params, config)
    },

    /**
     * 删除收藏
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} id 收藏的主键
     * @return {Promise} promise 对象
     */
    deleteFavorites ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`hosts/favorites/${id}`, config)
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
        return $http.put(`hosts/favorites/${id}/incr`)
    }
}

const mutations = {
    setApplying (state, collection) {
        state.applying = collection
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
