let commonRequestResolver
const commonRequest = new Promise((resolve, reject) => {
    commonRequestResolver = resolve
})
const state = {
    propertyMap: {},
    topologyModels: [],
    currentNode: null,
    commonRequest,
    commonRequestResolver
}

const getters = {
    getProperties: state => id => {
        return state.propertyMap[id] || []
    },
    topologyModels: state => state.topologyModels,
    columnsConfigProperties: (state, getters) => {
        const setProperties = getters.getProperties('set').filter(property => ['bk_set_name'].includes(property.bk_property_id))
        const moduleProperties = getters.getProperties('module').filter(property => ['bk_module_name'].includes(property.bk_property_id))
        const hostProperties = getters.getProperties('host')
        return [...setProperties, ...moduleProperties, ...hostProperties]
    },
    currentNode: state => state.currentNode,
    getDefaultSearchCondition: state => () => {
        let topologyModels = state.topologyModels
        if (!topologyModels.length) {
            topologyModels = ['biz', 'set', 'module', 'host'].map(modelId => ({ bk_obj_id: modelId }))
        }
        return topologyModels.map(model => ({
            bk_obj_id: model.bk_obj_id,
            condition: [],
            fields: []
        }))
    },
    commonRequest: state => state.commonRequest
}

const mutations = {
    setProperties (state, propertyMap = {}) {
        state.propertyMap = propertyMap
    },
    setTopologyModels (state, topologyModels) {
        state.topologyModels = topologyModels
    },
    setCurrentNode (state, node) {
        state.currentNode = node
    },
    resolveCommonRequest (state) {
        state.commonRequestResolver()
    },
    clear (state) {
        state.propertyMap = {}
        state.topologyModels = []
        state.currentNode = null
        state.commonRequest = new Promise((resolve, reject) => {
            state.commonRequestResolver = resolve
        })
    }
}
export default {
    namespaced: true,
    state,
    mutations,
    getters
}
