import Vue from 'vue'
const state = {
    propertyMap: {},
    propertyGroupMap: {},
    templateMap: {},
    categoryMap: {},
    selectedNode: null,
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
    setTemplates (state, data) {
        Vue.set(state.templateMap, data.id, data.templates)
    },
    setCategories (state, data) {
        Vue.set(state.categoryMap, data.id, data.categories)
    },
    setSelectedNode (state, node) {
        state.selectedNode = node
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
        state.templateMap = {}
        state.categoryMap = {}
        state.selectedNode = null
        state.hostSelectorVisible = false
        state.selectedHost = []
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
