const state = {
    filterList: []
}

const mutations = {
    setFilterList (state, list) {
        state.filterList = list
    }
}

export default {
    namespaced: true,
    state,
    mutations
}
