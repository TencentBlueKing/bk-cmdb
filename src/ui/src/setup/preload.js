import { getAuthorizedBusiness } from '@/router/business-interceptor.js'
const preloadConfig = {
    fromCache: false,
    cancelWhenRouteChange: false
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

export function getUserCustom (app) {
    return app.$store.dispatch('userCustom/searchUsercustom', {
        config: {
            ...preloadConfig,
            fromCache: false,
            requestId: 'post_searchUsercustom'
        }
    })
}

export function getGlobalUsercustom (app) {
    return app.$store.dispatch('userCustom/getGlobalUsercustom', {
        config: {
            ...preloadConfig,
            fromCache: false,
            globalError: false
        }
    }).catch(() => {
        return {}
    })
}

export async function getConfig (app) {
    return app.$store.dispatch('getConfig', {
        config: {
            ...preloadConfig,
            fromCache: false,
            globalError: false
        }
    }).then(data => {
        app.$store.commit('setConfig', data)
    }).catch(() => {
        return {}
    })
}

export default async function (app) {
    getConfig(app)
    getAuthorizedBusiness(app)
    return Promise.all([
        getClassifications(app),
        getUserCustom(app),
        getGlobalUsercustom(app)
    ])
}
