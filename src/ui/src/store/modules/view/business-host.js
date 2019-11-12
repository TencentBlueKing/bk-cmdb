import Vue from 'vue'
let commonRequestResolver
const commonRequest = new Promise((resolve, reject) => {
    commonRequestResolver = resolve
})
const state = {
    propertyMap: {},
    topologyModels: [],
    commonRequest,
    commonRequestResolver,
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
    propertyMap: state => state.propertyMap,
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
    selectedNode: state => state.selectedNode,
    getDefaultSearchCondition: state => () => {
        return ['biz', 'set', 'module', 'host', 'object'].map(modelId => ({
            bk_obj_id: modelId,
            condition: [],
            fields: []
        }))
    },
    commonRequest: state => state.commonRequest
}

const mutations = {
    setPropertyMap (state, propertyMap = {}) {
        state.propertyMap = propertyMap
    },
    setProperties (state, data) {
        Vue.set(state.propertyMap, data.id, data.properties)
    },
    setTopologyModels (state, topologyModels) {
        state.topologyModels = topologyModels
    },
    setSelectedNode (state, node) {
        state.selectedNode = node
    },
    resolveCommonRequest (state) {
        state.commonRequestResolver()
    },
    setPropertyGroups (state, data) {
        Vue.set(state.propertyGroupMap, data.id, data.groups)
    },
    setServiceTemplate (state, data) {
        Vue.set(state.serviceTemplateMap, data.id, data.templates)
    },
    setSetTemplate (state, data) {
        Vue.set(state.setTemplateMap, data.id, data.templates)
    },
    setProcessTemplate (state, data) {
        Vue.set(state.processTemplateMap, data.id, data.template)
    },
    setCategories (state, data) {
        Vue.set(state.categoryMap, data.id, data.categories)
    },
    setSelectedNodeInstance (state, instance) {
        state.selectedNodeInstance = instance
    },
    setHostSelectorVisible (state, visible) {
        state.hostSelectorVisible = visible
    },
    setSelectedHost (state, selectedHost) {
        state.selectedHost = selectedHost
    },
    setInstanceIp (state, { hostId, res }) {
        Vue.set(state.instanceIpMap, hostId, res)
    },
    clear (state) {
        state.propertyMap = {}
        state.topologyModels = []
        state.selectedNode = null
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
