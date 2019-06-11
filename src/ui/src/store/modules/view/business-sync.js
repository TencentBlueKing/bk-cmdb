import Vue from 'vue'
const state = {
    instanceMap: {}
}

const getters = {
}

const actions = {
}

const mutations = {
    setInstance (state, data) {
        Vue.set(state.instanceMap, data.id, data.instanceProperty)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
