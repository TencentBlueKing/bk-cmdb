import $http from '@/api'
const state = {
    auth: {}
}

const getters = {
    auth: state => state.auth
}

const actions = {
    getUserAuth ({commit}) {
        return $http.get('authUrl')
    }
}

const mutations = {
    setUserAuth (state, auth = {}) {
        state.auth = auth
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
