const state = {
    breadcrumbs: []
}

const getters = {
    breadcrumbs: state => state.breadcrumbs
}

const actions = {}

const mutations = {
    updateBreadcrumbs (state, breadcrumbs) {
        state.breadcrumbs = breadcrumbs
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
