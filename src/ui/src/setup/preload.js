import {
    SYSTEM_MANAGEMENT,
    GET_AUTH_META,
    GET_MODEL_INST_AUTH_META
} from '@/dictionary/auth'

const preloadConfig = {
    fromCache: true,
    cancelWhenRouteChange: false
}

export function getViewAuth (app) {
    const viewAuthorities = [/* GET_AUTH_META(SYSTEM_MANAGEMENT) */]
    app.$router.options.routes.forEach(route => {
        const meta = route.meta || {}
        const auth = meta.auth || {}
        const staticView = auth.view && (!auth.meta || auth.meta !== GET_MODEL_INST_AUTH_META)
        if (staticView) {
            viewAuthorities.push(GET_AUTH_META(auth.view))
        }
    })
    return app.$store.dispatch('auth/getViewAuth', viewAuthorities)
}

export function getClassifications (app) {
    return app.$store.dispatch('objectModelClassify/searchClassificationsObjects', {
        params: app.$injectMetadata(),
        config: {
            ...preloadConfig,
            requestId: 'post_searchClassificationsObjects'
        }
    })
}

export function getAuthorizedBusiness (app) {
    return app.$store.dispatch('objectBiz/getAuthorizedBusiness')
}

export function getUserCustom (app) {
    return app.$store.dispatch('userCustom/searchUsercustom', {
        config: {
            ...preloadConfig,
            fromCache: false,
            requestId: 'post_searchUsercustom'
        }
    })
}

export function getUserList (app) {
    return app.$store.dispatch('getUserList').then(list => {
        window.CMDB_USER_LIST = list
        app.$store.commit('setUserList', list)
        return list
    }).catch(e => {
        window.CMDB_USER_LIST = []
    })
}

export default async function (app) {
    return Promise.all([
        getViewAuth(app),
        getAuthorizedBusiness(app),
        getClassifications(app),
        getUserCustom(app),
        getUserList(app)
    ])
}
