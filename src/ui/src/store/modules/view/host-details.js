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

import has from 'has'
const state = {
  info: {},
  properties: [],
  propertyGroups: [],
  association: {
    source: [],
    target: []
  },
  mainLine: [],
  instances: [],
  associationTypes: [],
  expandAll: false
}

function isBizCustomData(data) {
  return has(data, 'bk_biz_id') && data.bk_biz_id > 0
}
const getters = {
  groupedProperties: (state) => {
    const groupedProperties = []
    state.propertyGroups.forEach((group) => {
      const properties = state.properties.filter(property => property.bk_property_group === group.bk_group_id)
      if (properties.length) {
        properties.sort((prev, next) => prev.bk_property_index - next.bk_property_index)
        groupedProperties.push({
          ...group,
          properties
        })
      }
    })
    return groupedProperties.sort((prev, next) => {
      const bizCustomPrev = isBizCustomData(prev)
      const bizCustomNext = isBizCustomData(next)
      if (
        (bizCustomPrev && bizCustomNext)
              || (!bizCustomPrev && !bizCustomNext)
      ) {
        return prev.bk_group_index - next.bk_group_index
      } if (bizCustomPrev) {
        return 1
      }
      return -1
    })
  },
  mainLine: state => state.mainLine,
  associationTypes: state => state.associationTypes,
  source: state => state.association.source,
  target: state => state.association.target,
  allInstances: state => state.instances,
  isBusinessHost: state => (state.info.biz || []).some(business => business.default === 0),
  properties: state => state.properties
}

const mutations = {
  setHostInfo(state, info) {
    state.info = info || {}
  },
  setHostProperties(state, properties) {
    state.properties = properties
  },
  setHostPropertyGroups(state, propertyGroups) {
    state.propertyGroups = propertyGroups
  },
  updateInfo(state, data) {
    state.info.host = Object.assign({}, state.info.host, data)
  },
  setAssociation(state, data) {
    state.association[data.type] = data.association
  },
  setMainLine(state, mainLine) {
    state.mainLine = mainLine
  },
  setInstances(state, data) {
    state.instances = data
  },
  setAssociationTypes(state, types) {
    state.associationTypes = types
  },
  deleteAssociation(state, id) {
    const index = state.instances.findIndex(instance => instance.id === id)
    index > -1 && state.instances.splice(index, 1)
  },
  toggleExpandAll(state, expandAll) {
    state.expandAll = expandAll
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations
}
