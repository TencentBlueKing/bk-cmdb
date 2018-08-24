const preloadConfig = {
    fromCache: true,
    cancelWhenRouteChange: false
}

function _preloadPrivilege (app) {
    return app.$store.dispatch('userPrivilege/getUserPrivilege', preloadConfig)
}

function _preloadClassifications (app) {
    return app.$store.dispatch('objectModelClassify/searchClassificationsObjects', {config: preloadConfig})
}

function _preloadBusiness (app) {
    return app.$store.dispatch('objectBiz/searchBusiness', {
        params: {
            'fields': ['bk_biz_id', 'bk_biz_name'],
            'condition': {
                'bk_data_status': {
                    '$ne': 'disabled'
                }
            }
        },
        config: preloadConfig
    }).then(business => {
        app.$store.commit('objectBiz/setBusiness', business.info)
        return business
    })
}

function _preloadUserCustom (app) {
    return app.$store.dispatch('userCustom/searchUsercustom', {
        config: preloadConfig
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
