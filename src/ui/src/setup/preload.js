const preloadConfig = {
    fromCache: true,
    cancelWhenRouteChange: false
}

export function _preloadPrivilege (app) {
    return app.$store.dispatch('userPrivilege/getUserPrivilege', {
        ...preloadConfig,
        requestId: 'get_getUserPrivilege'
    })
}

export function _preloadClassifications (app) {
    return app.$store.dispatch('objectModelClassify/searchClassificationsObjects', {
        config: {
            ...preloadConfig,
            requestId: 'post_searchClassificationsObjects'
        }
    })
}

export function _preloadBusiness (app) {
    return app.$store.dispatch('objectBiz/searchBusiness', {
        params: {
            'fields': ['bk_biz_id', 'bk_biz_name'],
            'condition': {
                'bk_data_status': {
                    '$ne': 'disabled'
                }
            }
        },
        config: {
            ...preloadConfig,
            requestId: 'post_searchBusiness_$ne_disabled'
        }
    }).then(business => {
        app.$store.commit('objectBiz/setBusiness', business.info)
        return business
    })
}

export function _preloadUserCustom (app) {
    return app.$store.dispatch('userCustom/searchUsercustom', {
        config: {
            ...preloadConfig,
            fromCache: false,
            requestId: 'post_searchUsercustom'
        }
    })
}

export default function (app) {
    return Promise.all([
        _preloadPrivilege(app),
        _preloadClassifications(app),
        _preloadBusiness(app),
        _preloadUserCustom(app)
    ])
}
