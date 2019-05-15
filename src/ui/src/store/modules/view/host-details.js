const state = {
    info: null
}

const mutations = {
    setHostInfo (state, info) {
        state.info = info
    }
}

export default {
    namespaced: true,
    state,
    mutations
}
