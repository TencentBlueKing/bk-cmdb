const state = {
    active: null,
    open: null
}

const getters = {
    active: state => state.active,
    open: state => state.open
}

const mutations = {
    setActiveMenu (state, menuId) {
        state.active = menuId
    },
    setOpenMenu (state, menuId) {
        state.open = menuId
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
