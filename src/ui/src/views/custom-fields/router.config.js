import Meta from '@/router/meta'

export default {
    name: 'customFields',
    path: '/custom-fields',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'customQuery',
            i18n: '自定义字段',
            path: '/custom-fields',
            order: 4,
            adminView: true
        },
        auth: {
            operation: [],
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        i18nTitle: '自定义字段'
    })
}
