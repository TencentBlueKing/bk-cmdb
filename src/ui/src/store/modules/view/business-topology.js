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
    getModelProperty (context) {
        console.log(context)
    }
}

const mutations = {
    setProperties (state, data) {
        Vue.set(state.propertyMap, data.id, data.properties)
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
    setSelectedNode (state, node) {
        state.selectedNode = node
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
    resetState (state) {
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
    setInstanceIp (state, { hostId, res }) {
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
