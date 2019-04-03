import $http from '@/api'

const state = {
    operation: [],
    view: []
}

const getters = {
    operation: state => state.operation,
    isAuthorized: state => (type, action, option = { isView: false }) => {
        const authList = option.isView ? state.view : state.operation
        const auth = authList.find(auth => {
            return auth.resource_type === type && auth.action === action
        })
        return (auth || {}).is_pass
    }
}

const actions = {
    async getOperationAuth ({commit}, list = []) {
        const authList = await $http.post('auth/verify', {
            resources: list
        })
        commit('setOperationAuth', authList)
        return authList
    },
    async getViewAuth ({ commit }, list) {
        const viewAuth = await $http.post('auth/verify', {
            resources: list
        }, {
            requestId: 'getViewAuth',
            fromCache: true,
            cancelWhenRouteChange: false
        })
        commit('setViewAuth', viewAuth)
        return viewAuth
    }
}

const mutations = {
    setOperationAuth (state, operationAuth) {
        state.operation = operationAuth
    },
    setViewAuth (state, viewAuth) {
        state.view = viewAuth
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
