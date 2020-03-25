import { language } from '@/i18n'
import $http from '@/api'

const state = {
    site: window.Site,
    user: window.User,
    supplier: window.Supplier,
    language: language,
    globalLoading: true,
    nav: {
        stick: window.localStorage.getItem('navStick') !== 'false',
        fold: window.localStorage.getItem('navStick') === 'false'
    },
    header: {
        back: false
    },
    layout: {
        mainFullScreen: false
    },
    userList: [],
    headerTitle: '',
    permission: [],
    appHeight: window.innerHeight,
    isAdminView: true,
    title: null,
    businessSelectorVisible: false,
    businessSelectorPromise: null,
    businessSelectorResolver: null,
    scrollerState: {
        scrollbar: false
    }
}

const getters = {
    site: state => state.site,
    user: state => state.user,
    userName: state => state.user.name,
    admin: state => state.user.admin === '1',
    isAdminView: state => state.isAdminView,
    isBusinessSelected: (state, getters, rootState, rootGetters) => {
        return rootGetters['objectBiz/bizId'] !== null
    },
    language: state => state.language,
    supplier: state => state.supplier,
    supplierAccount: state => state.supplier.account,
    globalLoading: state => state.globalLoading,
    navStick: state => state.nav.stick,
    navFold: state => state.nav.fold,
    mainFullScreen: state => state.layout.mainFullScreen,
    showBack: state => state.header.back,
    userList: state => state.userList,
    headerTitle: state => state.headerTitle,
    permission: state => state.permission,
    title: state => state.title,
    businessSelectorVisible: state => state.businessSelectorVisible,
    scrollerState: state => state.scrollerState
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
    },
    getBlueKingEditStatus ({ commit }, { config }) {
        return $http.post('system/config/user_config/blueking_modify', {}, config)
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
    setLayoutStatus (state, status) {
        Object.assign(state.layout, status)
    },
    setUserList (state, list) {
        state.userList = list
    },
    setAdminView (state, isAdminView) {
        state.isAdminView = isAdminView
    },
    setPermission (state, permission) {
        state.permission = permission
    },
    setAppHeight (state, height) {
        state.appHeight = height
    },
    setTitle (state, title) {
        state.title = title
    },
    setBusinessSelectorVisible (state, visible) {
        state.businessSelectorVisible = visible
    },
    createBusinessSelectorPromise (state) {
        state.businessSelectorPromise = new Promise(resolve => {
            state.businessSelectorResolver = resolve
        })
    },
    resolveBusinessSelectorPromise (state, val) {
        state.businessSelectorResolver && state.businessSelectorResolver(val)
    },
    setScrollerState (state, scrollerState) {
        Object.assign(state.scrollerState, scrollerState)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
