import { getAuthorizedBusiness } from '@/router/business-interceptor.js'
const preloadConfig = {
    fromCache: true,
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

export default async function (app) {
    getAuthorizedBusiness(app)
    return Promise.all([
        getClassifications(app),
        getUserCustom(app),
        getGlobalUsercustom(app)
    ])
}
