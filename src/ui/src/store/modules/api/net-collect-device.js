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
     * 新增设备
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    createDevice ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/device/action/create`, params, config)
    },
    /**
     * 更新设备
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} deviceId 设备id
     * @return {promises} promises 对象
     */
    updateDevice ({ commit, state, dispatch, rootGetters }, { deviceId, params, config }) {
        return $http.post(`collector/netcollect/device/${deviceId}/action/update`, params, config)
    },
    /**
     * 查询设备
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    searchDevice ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`collector/netcollect/device/action/search`, params, config)
    },
    /**
     * 批量删除设备
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 设备id
     * @return {promises} promises 对象
     */
    deleteDevice ({ commit, state, dispatch, rootGetters }, { config }) {
        return $http.delete(`collector/netcollect/device/action/delete`, config)
    },
    /**
     * 导入设备
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} deviceId 设备id
     * @return {promises} promises 对象
     */
    importNetDevice ({ commit, state, dispatch, rootGetters }, { config }) {
        return $http.post(`${window.API_HOST}collector/netdevice/import`, config)
    },
    /**
     * 导出设备
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    exportNetDevice ({ commit, state, dispatch, rootGetters }, { config }) {
        return $http.post(`${window.API_HOST}collector/netdevice/export`, config)
    },
    /**
     * 获取导入设备模板
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    getNetDeviceImportTemplate ({ commit, state, dispatch, rootGetters }, { config }) {
        return $http.get(`${window.API_HOST}collector/netcollect/importtemplate/netdevice`, config)
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
