export function getBusiness (app) {
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
            requestId: 'post_searchBusiness_$ne_disabled'
        }
    }).then(business => {
        app.$store.commit('objectBiz/setBusiness', business.info)
        return business
    })
}

export default async function (app, to, from) {
    const functions = []
    if (to.meta.requireBusiness) {
        functions.push(getBusiness)
    }
    console.log(functions, to.meta)
    return Promise.all(functions.map(func => func(app, to, from)))
}
