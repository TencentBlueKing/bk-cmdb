/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

/* eslint-disable no-unused-vars */
import $http from '@/api'
import { CONTAINER_OBJECTS, CONTAINER_OBJECT_NAMES } from '@/dictionary/container.js'

const state = {
  classifications: []
}

const getters = {
  classifications: state => state.classifications,
  models: (state) => {
    const models = []
    state.classifications.forEach((classification) => {
      (classification.bk_objects || []).forEach((model) => {
        models.push({
          ...model,
          bk_classification_name: classification.bk_classification_name,
          bk_classification_id: classification.bk_classification_id
        })
      })
    })
    return models
  },
  getModelById: (state, getters) => id => getters.models.find(model => model.bk_obj_id === id)
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
  createClassification({ commit, state, dispatch }, { params, config }) {
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
  deleteClassification({ commit, state, dispatch }, { id, config }) {
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
  updateClassification({ commit, state, dispatch }, { id, params }) {
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
  searchClassifications({ commit, state, dispatch }, { params, config }) {
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
  searchClassificationsObjects({ commit, state, dispatch, rootGetters }, { params = {}, config }) {
    return $http.post('find/classificationobject', params, config).then((data) => {
      const classification = data || []

      // 注入容器分组和对象
      const containerClassification = {
        id: Date.now(),
        bk_ishidden: true, // 在页面中不显示
        bk_classification_icon: '',
        bk_classification_id: 'bk_container',
        bk_classification_name: '容器',
        bk_classification_type: '',
        bk_objects: []
      }
      Object.keys(CONTAINER_OBJECTS).forEach((objKey) => {
        const objId = CONTAINER_OBJECTS[objKey]
        containerClassification.bk_objects.push({
          id: Date.now(),
          bk_classification_id: 'bk_container',
          bk_ishidden: true,
          bk_ispaused: false,
          ispre: false,
          bk_obj_icon: 'icon-cc-default',
          bk_obj_id: objId,
          bk_obj_name: CONTAINER_OBJECT_NAMES[objId].FULL,
          bk_supplier_account: '0',
          position: ''
        })
      })

      classification.push(containerClassification)

      commit('setClassificationsObjects', data)
      return data
    })
  },

  getClassificationsObjectStatistics({ state }, { config }) {
    return $http.get('object/statistics', config)
  }
}

const mutations = {
  setClassificationsObjects(state, classifications) {
    state.classifications = classifications
  },
  updateClassify(state, classification) {
    // eslint-disable-next-line max-len
    const activeClassification = state.classifications.find(({ bk_classification_id: bkClassificationId }) => bkClassificationId === classification.bk_classification_id)
    if (activeClassification) {
      activeClassification.bk_classification_icon = classification.bk_classification_icon
      activeClassification.bk_classification_name = classification.bk_classification_name
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
          id: 0,
          isNewClassify: classification.isNewClassify
        },
        ...classification
      })
    }
  },
  deleteClassify(state, classificationId) {
    // eslint-disable-next-line max-len
    const index = state.classifications.findIndex(({ bk_classification_id: bkClassificationId }) => bkClassificationId === classificationId)
    state.classifications.splice(index, 1)
  },
  updateModel(state, data) {
    const models = []
    state.classifications.forEach((classification) => {
      (classification.bk_objects || []).forEach((model) => {
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
