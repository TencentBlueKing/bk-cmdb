const state = {
    showClassify: false
}

const getters = {
    showClassify: state => state.showClassify
}

const mutations = {
    toggleClassify (state, showClassify) {
        state.showClassify = showClassify
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
