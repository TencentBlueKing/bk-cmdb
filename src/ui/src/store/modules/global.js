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
    userList: []
}

const getters = {
    site: state => state.site,
    user: state => state.user,
    userName: state => state.user.name,
    admin: state => state.user.admin === '1',
    language: state => state.language,
    supplier: state => state.supplier,
    supplierAccount: state => state.supplier.account,
    globalLoading: state => state.globalLoading,
    navStick: state => state.nav.stick,
    navFold: state => state.nav.fold,
    showBack: state => state.header.back,
    userList: state => state.userList
}

const actions = {
    getUserList ({commit}) {
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
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
