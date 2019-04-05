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
     * 创建分组基本信息
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createGroup ({ commit, state, dispatch }, { params, config }) {
        return $http.post('create/objectattgroup', params, config)
    },

    /**
     * 查询分组基本信息
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkObjId 模型id
     * @return {promises} promises 对象
     */
    searchGroup ({ commit, state, dispatch, rootGetters }, { objId, params, config }) {
        return $http.post(`find/objectattgroup/object/${objId}`, params, config)
    },

    /**
     * 修改分组基本信息
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateGroup ({ commit, state, dispatch }, { params, config }) {
        return $http.put('update/objectattgroup', params, config)
    },

    /**
     * 删除分组基本信息
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 分组记录标识
     * @return {promises} promises 对象
     */
    deleteGroup ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`delete/objectattgroup/${id}`, config)
    },

    /**
     * 更新模型属性分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updatePropertyGroup ({ commit, state, dispatch }, { params, config }) {
        return $http.put(`objectatt/group/property`, params, config)
    },

    /**
     * 删除模型属性分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkObjId 模型id
     * @param {String} bkPropertyId 属性id
     * @param {String} bkGroupId 分组id
     * @return {promises} promises 对象
     */
    deleteObjectPropertyGroup ({ commit, state, dispatch }, { objId, propertyId, groupId }) {
        return $http.delete(`delete/objectattgroupasst/object/${objId}/property/${propertyId}/group/${groupId}`)
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
