import $http from '@/api'

const state = {
    authList: [],
    staticViewAuth: []
}

const getters = {
    authList: state => state.authList,
    checkAuth: state => (type, action) => {
        const auth = state.authList.find(auth => {
            return auth.resource_type === type && auth.action === action
        })
        return (auth || {}).is_pass
    }
}

const actions = {
    async getAuthList ({commit}, list = []) {
        const authList = await $http.post('auth/verify', {
            resources: list
        })
        commit('setAuthList', authList)
        return authList
    },
    async getStaticViewAuth ({ commit }, list) {
        const staticViewAuth = await $http.post('auth/verify', {
            resources: list
        }, {
            fromCache: true,
            cancelWhenRouteChange: false
        })
        commit('setStaticViewAuth', staticViewAuth)
    }
}

const mutations = {
    setAuthList (state, list) {
        state.authList = list
    },
    setStaticViewAuth (state, staticViewAuth) {
        state.staticViewAuth = staticViewAuth
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
