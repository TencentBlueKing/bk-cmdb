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

function createIdProperty (objId) {
    const keyMap = {
        biz: 'bk_biz_id',
        host: 'bk_host_id'
    }
    return {
        id: Date.now(),
        bk_obj_id: objId,
        bk_property_id: keyMap[objId] || 'bk_inst_id',
        bk_property_name: 'ID',
        bk_property_index: -1,
        bk_property_type: 'int',
        isonly: true,
        ispre: true,
        bk_isapi: true,
        bk_issystem: true,
        isreadonly: true,
        editable: false,
        bk_property_group: null,
        _is_inject_: true
    }
}

const state = {

}

const getters = {

}

const actions = {
    /**
     * 创建对象模型属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createObjectAttribute ({ commit, state, dispatch }, { params, config }) {
        return $http.post('create/objectattr', params, config)
    },

    /**
     * 业务下创建对象模型属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createBizObjectAttribute ({ commit, state, dispatch }, { bizId, params, config }) {
        return $http.post(`/create/objectattr/biz/${bizId}`, params, config)
    },

    /**
     * 删除对象模型属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} id 被删除的数据记录的唯一标识id
     * @return {promises} promises 对象
     */
    deleteObjectAttribute ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`delete/objectattr/${id}`, config)
    },

    /**
     * 更新对象属性模型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} id 被删除的数据记录的唯一标识id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateObjectAttribute ({ commit, state, dispatch }, { id, params, config }) {
        return $http.put(`update/objectattr/${id}`, params, config)
    },

    /**
     * 业务下更新对象属性模型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} id 被删除的数据记录的唯一标识id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateBizObjectAttribute ({ commit, state, dispatch }, { bizId, id, params, config }) {
        return $http.put(`update/objectattr/biz/${bizId}/id/${id}`, params, config)
    },

    /**
     * 查询对象属性模型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchObjectAttribute ({ commit, state, dispatch }, { params, config, injectId = false }) {
        return $http.post('find/objectattr', params, config).then(data => {
            if (injectId !== params.bk_obj_id) {
                return data
            }
            const alreadyInject = data.some(property => property._is_inject_)
            if (alreadyInject) {
                return data
            }
            data.unshift(createIdProperty(injectId))
            return data
        })
    },

    /**
     * 批量查询对象属性模型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    batchSearchObjectAttribute ({ commit, state, dispatch }, { params, config, injectId = false }) {
        return $http.post(`find/objectattr`, params, config).then(properties => {
            const result = {}
            params.bk_obj_id.$in.forEach(objId => {
                result[objId] = []
            })
            properties.forEach(property => {
                result[property['bk_obj_id']].push(property)
            })
            Object.keys(result).forEach(objId => {
                if (injectId === objId) {
                    const alreadyInject = result[objId].some(property => property._is_inject_)
                    if (!alreadyInject) {
                        result[objId].unshift(createIdProperty(objId))
                    }
                }
            })
            return result
        })
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
