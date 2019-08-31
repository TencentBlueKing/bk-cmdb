import { language } from '@/i18n'
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
    headerTitle: '',
    featureTipsParams: {
        process: true,
        customQuery: true,
        model: true,
        modelBusiness: true,
        association: true,
        eventpush: true,
        adminTips: true,
        serviceTemplate: true,
        category: true,
        hostServiceInstanceCheckView: true
    },
    permission: [],
    appHeight: window.innerHeight
}

const getters = {
    site: state => state.site,
    user: state => state.user,
    userName: state => state.user.name,
    admin: state => state.user.admin === '1',
    isAdminView: (state, getters, rootState, rootGetters) => {
        const adminEntranceAuth = rootState.auth.adminEntranceAuth
        if (!adminEntranceAuth.is_pass) {
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
    headerTitle: state => state.headerTitle,
    featureTipsParams: state => state.featureTipsParams,
    permission: state => state.permission
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
    setAdminView (state, isAdminView) {
        window.sessionStorage.setItem('isAdminView', isAdminView)
        window.location = '/'
    },
    setFeatureTipsParams (state, tab) {
        const local = window.localStorage.getItem('featureTipsParams')
        if (tab) {
            state.featureTipsParams[tab] = false
            window.localStorage.setItem('featureTipsParams', JSON.stringify(state.featureTipsParams))
        } else if (local) {
            state.featureTipsParams = {
                ...state.featureTipsParams,
                ...JSON.parse(window.localStorage.getItem('featureTipsParams'))
            }
        } else {
            window.localStorage.setItem('featureTipsParams', JSON.stringify(state.featureTipsParams))
        }
    },
    setPermission (state, permission) {
        state.permission = permission
    },
    setAppHeight (state, height) {
        state.appHeight = height
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
