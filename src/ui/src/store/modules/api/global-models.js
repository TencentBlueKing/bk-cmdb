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
    graphicsData: []
}

const getters = {
    modelConfig: (state, getters, rootState, rootGetters) => {
        const modelConfig = {}
        rootGetters['objectModelClassify/models'].forEach(model => {
            modelConfig[model.bk_obj_id] = true
        })
        return Object.assign(modelConfig, rootGetters['userCustom/usercustom'].topoModelConfig)
    },
    displayModelGroups: (state, getters, rootState, rootGetters) => {
        const modelGroups = []
        const graphicsData = state.graphicsData
        rootGetters['objectModelClassify/classifications'].forEach(group => {
            const models = group.bk_objects || []
            if (models.length) {
                const availableModels = models.filter(model => {
                    const isSpecialModel = ['process', 'plat'].includes(model.bk_obj_id)
                    const asstData = graphicsData.find(data => data.bk_obj_id === model.bk_obj_id) || {}
                    const position = asstData.position || {}
                    const hasPosition = position.x !== null && position.y !== null
                    return !isSpecialModel && hasPosition
                })
                if (availableModels.length) {
                    modelGroups.push({
                        ...group,
                        bk_objects: availableModels
                    })
                }
            }
        })
        return modelGroups
    }
}

const actions = {
    /**
     * 订阅事件
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    searchModelAction ({ commit, state, dispatch, rootGetters }) {
        return $http.post(`objects/topographics/scope_type/global/scope_id/0/action/search`).then(data => {
            data = data.filter(node => {
                const model = rootGetters['objectModelClassify/getModelById'](node.bk_obj_id)
                return model && !model.bk_ispaused
            })
            commit('setGraphicsData', data)
            return data
        })
    },

    /**
     * 批量更新节点位置信息
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    updateModelAction ({ commit, state, dispatch }, { params }) {
        return $http.post(`objects/topographics/scope_type/global/scope_id/0/action/update`, params).then(() => {
            commit('updateGraphicsData', params)
        })
    }
}

const mutations = {
    setGraphicsData (state, data) {
        state.graphicsData = data
    },
    updateGraphicsData (state, data) {
        data.forEach(node => {
            const exist = state.graphicsData.find(exist => exist.bk_obj_id === node.bk_obj_id)
            if (exist) {
                Object.assign(exist, node)
            }
        })
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
