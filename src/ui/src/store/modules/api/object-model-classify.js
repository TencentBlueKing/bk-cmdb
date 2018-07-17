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
     * 添加模型分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createClassification ({ commit, state, dispatch }, { params }) {
        return $axios.post(`object/classification`, params)
    },

    /**
     * 删除模型分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 分类数据记录id
     * @return {promises} promises 对象
     */
    deleteClassification ({ commit, state, dispatch }, { id }) {
        return $axios.delete(`object/classification/${id}`)
    },

    /**
     * 更新模型分类数据
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 分类数据记录id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateClassification ({ commit, state, dispatch }, { id, params }) {
        return $axios.delete(`object/classification/${id}`, params)
    },

    /**
     * 查询模型分类列表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchClassifications ({ commit, state, dispatch }) {
        return $axios.post(`object/classifications`)
    },

    /**
     * 查询模型分类及附属模型信息
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchClassificationsObjects ({ commit, state, dispatch }, { bkSupplierAccount, params }) {
        return $axios.post(`object/classification/${bkSupplierAccount}/objects`, params)
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
