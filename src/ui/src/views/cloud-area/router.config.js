import Meta from '@/router/meta'
export default {
    name: 'cloudArea',
    path: 'cloud-area',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '首页'
        },
        auth: {
            view: null
        }
    })
}
