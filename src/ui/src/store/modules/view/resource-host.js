const state = {
    activeDirectory: null
}

const getters = {
    activeDirectory: state => state.activeDirectory || {}
}

const mutations = {
    setActiveDirectory (state, active) {
        state.activeDirectory = active
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
