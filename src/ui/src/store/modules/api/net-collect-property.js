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
     * 新增属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    createNetcollectProperty ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/property/action/create`, params, config)
    },
    /**
     * 更新属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} propertyId 属性id
     * @return {promises} promises 对象
     */
    updateNetcollectProperty ({ commit, state, dispatch, rootGetters }, { propertyId, params, config }) {
        return $http.post(`collector/netcollect/property/${propertyId}/action/update`, params, config)
    },
    /**
     * 查询属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    searchNetcollectProperty ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/property/action/search`, params, config)
    },
    /**
     * 批量删除属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 属性id
     * @return {promises} promises 对象
     */
    deleteNetcollectProperty ({ commit, state, dispatch, rootGetters }, { config }) {
        return $http.delete(`collector/netcollect/property/action/delete`, config)
    },
    /**
     * 导入属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} deviceId 属性id
     * @return {promises} promises 对象
     */
    importNetNetcollectProperty ({ commit, state, dispatch, rootGetters }, { config }) {
        return $http.post(`${window.API_HOST}collector/netproperty/import`, config)
    },
    /**
     * 导出属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    exportNetcollectProperty ({ commit, state, dispatch, rootGetters }, { config }) {
        return $http.post(`${window.API_HOST}collector/netproperty/export`, config)
    },
    /**
     * 获取导入属性模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    getNetDeviceImportTemplate ({ commit, state, dispatch, rootGetters }, { config }) {
        return $http.get(`${window.API_HOST}collector/netcollect/importtemplate/netproperty`, config)
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
