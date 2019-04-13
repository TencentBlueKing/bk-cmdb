import { language } from '@/i18n'
import { SYSTEM_MANAGEMENT, GET_AUTH_META } from '@/dictionary/auth'
import $http from '@/api'

const state = {
    site: window.Site,
    user: window.User,
    supplier: window.Supplier,
    language: language,
    globalLoading: false,
    nav: {
        stick: window.localStorage.getItem('navStick') !== 'false',
        fold: window.localStorage.getItem('navStick') === 'false'
    },
    header: {
        back: false
    },
    userList: [],
    headerTitle: ''
}

const getters = {
    site: state => state.site,
    user: state => state.user,
    userName: state => state.user.name,
    admin: state => state.user.admin === '1',
    isAdminView: (state, getters, rootState, rootGetters) => {
        const systemAuth = rootState.auth.system
        const managementData = systemAuth.find(data => {
            const meta = GET_AUTH_META(SYSTEM_MANAGEMENT)
            return data.resource_type === meta.resource_type && data.action === meta.action
        }) || {}
        if (!managementData.is_pass) {
            return false
        }
        if (window.sessionStorage.hasOwnProperty('isAdminView')) {
            return window.sessionStorage.getItem('isAdminView') === 'true'
        } else {
            window.sessionStorage.setItem('isAdminView', false)
            return false
        }
    },
    isBusinessSelected: (state, getters, rootState, rootGetters) => {
        return rootGetters['objectBiz/bizId'] !== null
    },
    language: state => state.language,
    supplier: state => state.supplier,
    supplierAccount: state => state.supplier.account,
    globalLoading: state => state.globalLoading,
    navStick: state => state.nav.stick,
    navFold: state => state.nav.fold,
    showBack: state => state.header.back,
    userList: state => state.userList,
    headerTitle: state => state.headerTitle
}

const actions = {
    getUserList ({ commit }) {
        return $http.get(`${window.API_HOST}user/list?_t=${(new Date()).getTime()}`, {
            requestId: 'get_user_list',
            fromCache: true,
            cancelWhenRouteChange: false
        }).then(list => {
            commit('setUserList', list)
            return list
        })
    }
}

const mutations = {
    setGlobalLoading (state, loading) {
        state.globalLoading = loading
    },
    setNavStatus (state, status) {
        Object.assign(state.nav, status)
    },
    setHeaderStatus (state, status) {
        Object.assign(state.header, status)
    },
    setUserList (state, list) {
        state.userList = list
    },
    setHeaderTitle (state, headerTitle) {
        state.headerTitle = headerTitle
    },
    setAdminView (state, params) {
        window.sessionStorage.setItem('isAdminView', params.isAdminView)
        params._this.$router.push({ name: 'index' })
        window.location.reload()
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
