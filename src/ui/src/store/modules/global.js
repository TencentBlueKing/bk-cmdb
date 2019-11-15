import { language } from '@/i18n'
import $http from '@/api'

let businessSelectorResolver
const businessSelectorPromise = new Promise(resolve => {
    businessSelectorResolver = resolve
})

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
        hostServiceInstanceCheckView: true,
        customFields: true,
        hostApply: true,
        hostApplyConfirm: true,
        hostApplyConflict: true
    },
    permission: [],
    appHeight: window.innerHeight,
    isAdminView: true,
    breadcrumbs: [],
    title: null,
    businessSelectorVisible: false,
    businessSelectorPromise,
    businessSelectorResolver,
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
    featureTipsParams: state => state.featureTipsParams,
    permission: state => state.permission,
    breadcrumbs: state => state.breadcrumbs,
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
    },
    setBreadcrumbs (state, breadcrumbs) {
        state.breadcrumbs = breadcrumbs
    },
    setTitle (state, title) {
        state.title = title
    },
    setBusinessSelectorVisible (state, visible) {
        state.businessSelectorVisible = visible
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
