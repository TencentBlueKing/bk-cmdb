/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */
import Vue from 'vue'
import $http from '@/api'

const state = {
    isEditMode: false,
    topologyData: [],
    topologyMap: {},
    options: {},
    edgeOptions: [],
    association: {
        show: false,
        edge: null
    },
    addEdgePromise: {
        resolve: null,
        reject: null
    }
}

const getters = {
    isEditMode: state => state.isEditMode,
    topologyData: state => state.topologyData,
    topologyMap: state => state.topologyMap,
    options: state => state.options,
    edgeOptions: state => state.edgeOptions,
    association: state => state.association,
    addEdgePromise: state => state.addEdgePromise
}

const actions = {
    /**
     * 查询模型拓扑
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    searchModelAction ({ commit, state, dispatch, rootGetters }, params) {
        return $http.post(`find/objecttopo/scope_type/global/scope_id/0`, params).then(data => {
            return data.filter(node => {
                const model = rootGetters['objectModelClassify/getModelById'](node.bk_obj_id)
                return model && !model.bk_ispaused && !model.bk_ishidden
            })
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
        return $http.post(`update/objecttopo/scope_type/global/scope_id/0`, params)
    }
}

const mutations = {
    setTopologyData (state, topologyData) {
        const topologyMap = {}
        topologyData.forEach(data => {
            topologyMap[data['bk_obj_id']] = data
        })
        state.topologyData = topologyData
        state.topologyMap = topologyMap
    },
    updateTopologyData (state, queue) {
        const updateQueue = Array.isArray(queue) ? queue : [queue]
        const topologyMap = state.topologyMap
        updateQueue.forEach(data => {
            const modelId = data['bk_obj_id']
            Object.assign(topologyMap[modelId], data)
        })
    },
    addAssociation (state, { id, association }) {
        const data = state.topologyMap[id]
        const associations = data.assts
        if (Array.isArray(associations)) {
            associations.push(association)
        } else {
            Vue.set(data, 'assts', [association])
        }
    },
    deleteAssociation (state, associationId) {
        const topologyData = state.topologyData
        for (let i = 0; i < topologyData.length; i++) {
            const associations = topologyData[i]['assts'] || []
            const index = associations.findIndex(association => association['bk_inst_id'] === associationId)
            if (index > -1) {
                associations.splice(index, 1)
                break
            }
        }
    },
    changeEditMode (state) {
        state.isEditMode = !state.isEditMode
    },
    setOptions (state, options) {
        state.options = options
    },
    setEdgeOptions (state, edgeOptions) {
        state.edgeOptions = edgeOptions
    },
    setAssociation (state, data) {
        Object.assign(state.association, data)
    },
    setAddEdgePromise (state, promise) {
        state.addEdgePromise = promise
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
