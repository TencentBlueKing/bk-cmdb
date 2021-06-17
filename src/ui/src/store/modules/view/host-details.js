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
  instances: {
    source: [],
    target: []
  },
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
  sourceInstances: state => state.instances.source,
  targetInstances: state => state.instances.target,
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
    state.instances[data.type] = data.instances
  },
  setAssociationTypes(state, types) {
    state.associationTypes = types
  },
  deleteAssociation(state, data) {
    const { type } = data
    const { model } = data
    const target = data.association
    const instances = state.instances[type === 'source' ? 'target' : 'source']
    const associations = instances.find(data => data.bk_obj_id === model)
    const index = associations.children.findIndex(association => association.asso_id === target.asso_id)
    if (index > -1) {
      associations.children.splice(index, 1)
    }
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
