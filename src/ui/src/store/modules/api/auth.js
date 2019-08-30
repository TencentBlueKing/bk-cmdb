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
    businessMeta: {},
    parentMeta: {},
    resourceMeta: [],
    adminEntranceAuth: {}
}

const getters = {
    operation: state => state.operation,
    isAuthorized: (state, getters) => (auth, option = { type: 'operation' }) => {
        const authMeta = getters.getAuthMeta(auth, option)
        const authList = state[option.type] || []
        const authData = authList.find(auth => {
            const sameType = auth.resource_type === authMeta.resource_type
            const sameAction = auth.action === authMeta.action
            const sameBusiness = authMeta.hasOwnProperty('bk_biz_id') ? auth.bk_biz_id === authMeta.bk_biz_id : true
            return sameType && sameAction && sameBusiness
        })
        return (authData || {}).is_pass
    },
    getAuthMeta: (state, getters, rootState, rootGetters) => (auth, option = {}) => {
        const meta = GET_AUTH_META(auth, option)
        const isBusinessMode = !rootGetters.isAdminView
        const bizId = rootGetters['objectBiz/bizId']
        if (
            isBusinessMode
            && bizId
            && STATIC_BUSINESS_MODE.includes(auth)
        ) {
            meta.bk_biz_id = bizId
        }
        if (DYNAMIC_BUSINESS_MODE.includes(auth)) {
            Object.assign(meta, state.parentMeta)
            Object.assign(meta, state.businessMeta)
        }
        const resourceMeta = state.resourceMeta.find(resourceMeta => {
            return resourceMeta.resource_type === meta.resource_type
                && resourceMeta.action === meta.action
        })
        if (resourceMeta) {
            Object.assign(meta, resourceMeta)
        }
        delete meta.scope
        return meta
    }
}

const actions = {
    async getAuth ({ commit, getters, rootGetters }, params) {
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
    },
    async getAdminEntranceAuth ({ commit, getters, rootGetters }, params, config) {
        const data = {
            is_pass: false
        }
        try {
            const response = await $http.get('auth/admin_entrance', config)
            data.is_pass = response.is_pass
        } catch (e) {
            console.error(e)
        }
        commit('setAdminEntranceAuth', data)
        return Promise.resolve(data)
    },
    getSkipUrl (context, { params, config }) {
        return $http.post('auth/skip_url', params, config)
    }
}

const mutations = {
    setAuth (state, data) {
        state[data.type] = data.auth
    },
    setParentMeta (state, meta = {}) {
        state.parentMeta = meta
    },
    setBusinessMeta (state, meta = {}) {
        state.businessMeta = meta
    },
    setResourceMeta (state, meta) {
        const resourceMeta = Array.isArray(meta) ? meta : [meta]
        state.resourceMeta.push(...resourceMeta)
    },
    setAdminEntranceAuth (state, data) {
        state.adminEntranceAuth = data
    },
    clearDynamicMeta (state) {
        state.parentMeta = {}
        state.businessMeta = {}
        state.resourceMeta = []
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
