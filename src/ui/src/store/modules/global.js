import { language } from '@/i18n'

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
    }
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
    showBack: state => state.header.back
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
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
