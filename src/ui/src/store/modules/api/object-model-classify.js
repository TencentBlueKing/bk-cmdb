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
    classifications: [],
    invisibleClassifications: ['bk_host_manage', 'bk_biz_topo']
}

const getters = {
    classifications: state => state.classifications,
    models: state => {
        const models = []
        state.classifications.forEach(classification => {
            (classification['bk_objects'] || []).forEach(model => {
                models.push(model)
            })
        })
        return models
    },
    getModelById: (state, getters) => id => {
        return getters.models.find(model => model['bk_obj_id'] === id)
    },
    activeClassifications: state => {
        const classifications = state.classifications
        // 1.去掉停用模型
        let activeClassifications = classifications.map(classification => {
            const activeClassification = { ...classification }
            activeClassification['bk_objects'] = activeClassification['bk_objects'].filter(model => !model['bk_ispaused'])
            return activeClassification
        })
        // 2.去掉无启用模型的分类和不显示的分类
        activeClassifications = activeClassifications.filter(classification => {
            const {
                'bk_classification_id': bkClassificationId,
                'bk_objects': bkObjects
            } = classification
            return !state.invisibleClassifications.includes(bkClassificationId) && Array.isArray(bkObjects) && bkObjects.length
        })
        return activeClassifications
    }
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
    createClassification ({ commit, state, dispatch }, { params, config }) {
        return $http.post('create/objectclassification', params, config)
    },

    /**
     * 删除模型分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 分类数据记录id
     * @return {promises} promises 对象
     */
    deleteClassification ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`delete/objectclassification/${id}`, config)
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
        return $http.put(`update/objectclassification/${id}`, params)
    },

    /**
     * 查询模型分类列表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchClassifications ({ commit, state, dispatch }, { params, config }) {
        return $http.post('find/objectclassification', params || {}, config)
    },

    /**
     * 查询模型分类及附属模型信息
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchClassificationsObjects ({ commit, state, dispatch, rootGetters }, { params = {}, config }) {
        return $http.post('find/classificationobject', params, config).then(data => {
            commit('setClassificationsObjects', data)
            return data
        })
    },

    getClassificationsObjectStatistics ({ state }, { config }) {
        return $http.get('object/statistics', config)
    }
}

const mutations = {
    setClassificationsObjects (state, classifications) {
        state.classifications = classifications
    },
    updateClassify (state, classification) {
        const activeClassification = state.classifications.find(({ bk_classification_id: bkClassificationId }) => bkClassificationId === classification['bk_classification_id'])
        if (activeClassification) {
            activeClassification['bk_classification_icon'] = classification['bk_classification_icon']
            activeClassification['bk_classification_name'] = classification['bk_classification_name']
        } else {
            state.classifications.push({
                ...{
                    bk_asst_objects: {},
                    bk_classification_icon: 'icon-cc-default',
                    bk_classification_id: '',
                    bk_classification_name: '',
                    bk_classification_type: '',
                    bk_objects: [],
                    bk_supplier_account: '',
                    id: 0
                },
                ...classification
            })
        }
    },
    deleteClassify (state, classificationId) {
        const index = state.classifications.findIndex(({ bk_classification_id: bkClassificationId }) => bkClassificationId === classificationId)
        state.classifications.splice(index, 1)
    },
    updateModel (state, data) {
        const models = []
        state.classifications.forEach(classification => {
            (classification['bk_objects'] || []).forEach(model => {
                models.push(model)
            })
        })
        const model = models.find(model => model.bk_obj_id === data.bk_obj_id)
        if (model) {
            Object.assign(model, data)
        }
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
