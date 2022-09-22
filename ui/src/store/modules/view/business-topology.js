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
import Vue from 'vue'
const state = {
  propertyMap: {},
  propertyGroupMap: {},
  serviceTemplateMap: {},
  setTemplateMap: {},
  processTemplateMap: {},
  categoryMap: {},
  instanceIpMap: {},
  selectedNode: null,
  selectedNodeInstance: null,
  hostSelectorVisible: false,
  selectedHost: []
}

const getters = {
}

const actions = {
  getModelProperty(context) {
    console.log(context)
  }
}

const mutations = {
  setProperties(state, data) {
    Vue.set(state.propertyMap, data.id, data.properties)
  },
  setPropertyGroups(state, data) {
    Vue.set(state.propertyGroupMap, data.id, data.groups)
  },
  setServiceTemplate(state, data) {
    Vue.set(state.serviceTemplateMap, data.id, data.templates)
  },
  setSetTemplate(state, data) {
    Vue.set(state.setTemplateMap, data.id, data.templates)
  },
  setProcessTemplate(state, data) {
    Vue.set(state.processTemplateMap, data.id, data.template)
  },
  setCategories(state, data) {
    Vue.set(state.categoryMap, data.id, data.categories)
  },
  setSelectedNode(state, node) {
    state.selectedNode = node
  },
  setSelectedNodeInstance(state, instance) {
    state.selectedNodeInstance = instance
  },
  setHostSelectorVisible(state, visible) {
    state.hostSelectorVisible = visible
  },
  setSelectedHost(state, selectedHost) {
    state.selectedHost = selectedHost
  },
  resetState(state) {
    state.propertyMap = {}
    state.propertyGroupMap = {}
    state.serviceTemplateMap = {}
    state.setTemplateMap = {}
    state.processTemplateMap = {}
    state.categoryMap = {}
    state.instanceIpMap = {}
    state.selectedNode = null
    state.selectedNodeInstance = null
    state.hostSelectorVisible = false
    state.selectedHost = []
  },
  setInstanceIp(state, { hostId, res }) {
    Vue.set(state.instanceIpMap, hostId, res)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
