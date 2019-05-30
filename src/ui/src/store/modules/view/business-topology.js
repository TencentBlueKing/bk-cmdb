import Vue from 'vue'
const state = {
    propertyMap: {}
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
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
