const defaultParams = {
    ip: {
        flag: 'bk_host_innerip|bk_host_outer',
        exact: 0,
        data: []
    }
}
const state = {
    filterList: [],
    filterParams: {
        ...defaultParams
    }
}

const mutations = {
    clearFilter (state) {
        state.filterList = []
        state.filterParams = {
            ...defaultParams
        }
    },
    setFilterList (state, list) {
        state.filterList = list
    },
    setFilterParams (state, params) {
        state.filterParams = params
    }
}

export default {
    namespaced: true,
    state,
    mutations
}
