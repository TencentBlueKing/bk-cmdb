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
     * 获取所有正在统计的图表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    getCountedCharts ({ commit, state, dispatch }, { params, config }) {
        return $http.get(`findmany/operation/chart`, config)
    },

    /**
     * 获取所有正在统计的图表的对应数据
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    getCountedChartsData ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`find/operation/chart/data`, params, config)
    },

    /**
     * 获取统计维度
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    getStaticDimeObj ({ commit, state, dispatch }, { params, config }) {
        return $http.post('find/objectattr', params, config)
    },

    /**
     * 获取统计对象
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    getStaticObj ({ commit, state, dispatch }, { params, config }) {
        return $http.post('find/object', params, config)
    },

    /**
     * 新增统计图表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    newStatisticalCharts ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`create/operation/chart`, params, config)
    },

    /**
     * 编辑统计图表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateStatisticalCharts ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`update/operation/chart`, params, config)
    },

    /**
     * 删除统计图表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteOperationChart ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`delete/operation/chart/${id}`, config)
    },

    /**
     * 更新图表位置
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateChartPosition ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`/update/operation/chart/position`, params, config)
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
