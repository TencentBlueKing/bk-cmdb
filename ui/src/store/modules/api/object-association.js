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

const state = {
  associationList: []
}

const getters = {
  associationList: state => state.associationList
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
  searchAssociationType({ commit, state, dispatch }, { params, config }) {
    return $http.post('find/associationtype', params, config)
  },
  /**
     * 添加关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  createAssociationType({ commit, state, dispatch }, { params, config }) {
    return $http.post('create/associationtype', params, config)
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
  updateAssociationType({ commit, state, dispatch }, { id, params, config }) {
    return $http.put(`update/associationtype/${id}`, params, config)
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
  deleteAssociationType({ commit, state, dispatch }, { id, config }) {
    return $http.delete(`delete/associationtype/${id}`, config)
  },
  /**
     * 查询模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  searchObjectAssociation({ commit, state, dispatch }, { params, config }) {
    return $http.post('find/objectassociation', params, config)
  },
  /**
     * 添加模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  createObjectAssociation({ commit, state, dispatch }, { params, config }) {
    return $http.post('create/objectassociation', params, config)
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
  updateObjectAssociation({ commit, state, dispatch }, { id, params, config }) {
    return $http.put(`update/objectassociation/${id}`, params, config)
  },
  /**
     * 删除模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  deleteObjectAssociation({ commit, state, dispatch }, { id, config }) {
    return $http.delete(`delete/objectassociation/${id}`, config)
  },
  /**
     * 根据关联类型查询使用这些关联类型的关联关系列表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  searchAssociationListWithAssociationKindList({ commit, state, dispatch }, { params, config }) {
    return $http.post('find/topoassociationtype', params, config)
  },
  /**
     * 查询实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  searchInstAssociation({ commit, state, dispatch }, { params, config }) {
    return $http.post('find/instassociation', params, config)
  },
  /**
     * 添加实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  createInstAssociation({ commit, state, dispatch }, { params, config }) {
    return $http.post('create/instassociation', params, config)
  },
  /**
     * 删除实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  deleteInstAssociation({ commit, state, dispatch }, { id, objId, config }) {
    return $http.delete(`delete/instassociation/${objId}/${id}`, config)
  }
}

const mutations = {
  setAssociationList(state, list) {
    state.associationList = list
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
