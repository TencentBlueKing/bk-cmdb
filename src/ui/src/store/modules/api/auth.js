import Vue from 'vue'
import $http from '@/api'
import {
    GET_AUTH_META,
    STATIC_BUSINESS_MODE,
    DYNAMIC_BUSINESS_MODE
} from '@/dictionary/auth'
const state = {
    operation: [],
    view: [],
    system: [],
    dynamicMeta: {}
}

const getters = {
    operation: state => state.operation,
    isAuthorized: (state, getters) => (auth, option = { type: 'operation' }) => {
        const authMeta = getters.getAuthMeta(auth)
        const authList = state[option.type] || []
        const authData = authList.find(auth => {
            const sameType = auth.resource_type === authMeta.resource_type
            const sameAction = auth.action === authMeta.action
            const sameBusiness = authMeta.hasOwnProperty('bk_biz_id') ? auth.bk_biz_id === authMeta.bk_biz_id : true
            return sameType && sameAction && sameBusiness
        })
        return (authData || {}).is_pass
    },
    getAuthMeta: (state, getters, rootState, rootGetters) => auth => {
        const meta = GET_AUTH_META(auth)
        const isBusinessMode = !rootGetters.isAdminView
        const bizId = rootGetters['objectBiz/bizId']
        if (
            isBusinessMode &&
            bizId &&
            STATIC_BUSINESS_MODE.includes(auth)
        ) {
            meta.bk_biz_id = bizId
        }
        if (DYNAMIC_BUSINESS_MODE.includes(auth)) {
            Object.assign(meta, getters.dynamicMeta)
        }
        return meta
    }
}

const actions = {
    async getOperationAuth ({commit, getters}, list = []) {
        const authList = await $http.post('auth/verify', {
            resources: list.map(auth => getters.getAuthMeta(auth))
        })
        commit('setOperationAuth', authList)
        return authList
    },
    async getViewAuth ({ commit, getters }, list) {
        const viewAuth = await $http.post('auth/verify', {
            resources: list.map(auth => getters.getAuthMeta(auth))
        }, {
            requestId: 'getViewAuth',
            fromCache: true,
            cancelWhenRouteChange: false
        })
        commit('setViewAuth', viewAuth)
        return viewAuth
    },
    async getSystemAuth ({ commit, getters }, list) {
        const systemAuth = await $http.post('auth/verify', {
            resources: list.map(auth => getters.getAuthMeta(auth))
        }, {
            requestId: 'getSystemAuth',
            fromCache: true,
            cancelWhenRouteChange: false
        })
        commit('setSystemAuth', systemAuth)
        return systemAuth
    }
}

const mutations = {
    setOperationAuth (state, operationAuth) {
        state.operation = operationAuth
    },
    setViewAuth (state, viewAuth) {
        state.view = viewAuth
    },
    setSystemAuth (state, systemAuth) {
        state.system = systemAuth
    },
    setDynamicMeta (state, meta) {
        state.dynamicMeta = meta
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
