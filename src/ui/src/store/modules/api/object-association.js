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
     * 查询关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchAssociationType ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`topo/association/type/action/search`, params, config)
    },
    /**
     * 添加关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createAssociationType ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`topo/association/type/action/create`, params, config)
    },
    /**
     * 编辑关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateAssociationType ({ commit, state, dispatch }, { id, params, config }) {
        return $http.post(`topo/association/type/${id}/action/update`, params, config)
    },
    /**
     * 删除关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteAssociationType ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`topo/association/type/${id}/action/delete`, config)
    },
    /**
     * 查询模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchObjectAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`object/association/action/search`, params, config)
    },
    /**
     * 添加模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createObjectAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`object/association/action/create`, params, config)
    },
    /**
     * 编辑模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateObjectAssociation ({ commit, state, dispatch }, { id, params, config }) {
        return $http.post(`object/association/${id}/action/update`, params, config)
    },
    /**
     * 删除关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteObjectAssociation ({ commit, state, dispatch }, { id, params, config }) {
        return $http.delete(`object/association/${id}/action/delete`, params, config)
    },
    /**
     * 查询实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchInstAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`inst/association/action/search`, params, config)
    },
    /**
     * 添加实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createInstAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`inst/association/action/create`, params, config)
    },
    /**
     * 删除实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteInstAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`inst/association/action/delete`, params, config)
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
