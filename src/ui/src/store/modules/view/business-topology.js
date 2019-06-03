import Vue from 'vue'
const state = {
    propertyMap: {},
    propertyGroupMap: {},
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
