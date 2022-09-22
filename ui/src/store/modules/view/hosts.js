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
function getDefaultCondition() {
  return ['biz', 'set', 'module', 'host', 'object'].map(modelId => ({
    bk_obj_id: modelId,
    fields: [],
    condition: []
  }))
}
const state = {
  filterList: [],
  collection: null,
  collectionList: [],
  propertyList: [],
  condition: getDefaultCondition(),
  shouldInjectAsset: true // 控制是否注入固资编号
}

const getters = {
  isCollection: state => !!state.collection,
  configPropertyList: (state) => {
    const disableList = ['bk_cpu']
    state.propertyList.forEach((property) => {
      property.options = property.option
      // eslint-disable-next-line no-underscore-dangle
      property.__extra__ = {
        disabled: disableList.includes(property.bk_property_id)
      }
    })

    return state.propertyList
  },
  condition: state => state.condition,
  shouldInjectAsset: state => state.shouldInjectAsset
}

const mutations = {
  setFilterList(state, list) {
    state.filterList = list
  },
  setCollectionList(state, list) {
    state.collectionList = list
  },
  setCollection(state, collection) {
    state.collection = collection
  },
  addCollection(state, collection) {
    state.collectionList.push(collection)
  },
  updateCollection(state, updatedData) {
    Object.assign(state.collection, updatedData)
  },
  deleteCollection(state, id) {
    state.collectionList = state.collectionList.filter(collection => collection.id !== id)
  },
  clearFilter(state) {
    state.filterList = []
    state.collection = null
    state.condition = getDefaultCondition()
  },
  setPropertyList(state, list) {
    state.propertyList = list
  },
  setCondition(state, condition) {
    state.condition = condition
  },
  setShouldInjectAsset(state, shouldInjectAsset) {
    state.shouldInjectAsset = !!shouldInjectAsset
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations
}
