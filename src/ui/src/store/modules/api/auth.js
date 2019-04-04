import Vue from 'vue'
import $http from '@/api'
import {
    GET_AUTH_META,
    STATIC_BUSINESS_MODE,
    DYNAMIC_BUSINESS_MODE
} from '@/dictionary/auth'

const defaultAuthData = {
    bk_biz_id: 0,
    is_pass: false,
    parent_layers: null,
    reason: '',
    resource_id: 0
}

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
    async getAuth ({commit, getters, rootGetters}, params) {
        const allAuth = params.list || []
        const authType = params.type || 'operation'
        const shouldAuth = []
        const shouldNotAuth = []
        const hasBusiness = !!rootGetters['objectBiz/bizId']
        const isBusinessMode = !rootGetters.isAdminView

        allAuth.forEach(auth => {
            const isBusinessAuth = STATIC_BUSINESS_MODE.includes(auth)
            if (isBusinessMode && isBusinessAuth) {
                if (hasBusiness) {
                    shouldAuth.push(auth)
                } else {
                    shouldNotAuth.push(auth)
                }
            } else {
                shouldAuth.push(auth)
            }
        })

        let authData = []
        if (shouldAuth.length) {
            authData = await $http.post('auth/verify', {
                resources: shouldAuth.map(auth => getters.getAuthMeta(auth))
            }, params.config || {})
        }

        const allAuthData = authData.concat(shouldNotAuth.map(auth => {
            const meta = getters.getAuthMeta(auth)
            return {
                ...defaultAuthData,
                ...meta
            }
        }))

        commit('setAuth', {
            type: authType,
            auth: allAuthData
        })

        return allAuthData
    }
}

const mutations = {
    setAuth (state, data) {
        state[data.type] = data.auth
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
