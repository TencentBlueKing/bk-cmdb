export default {
    name: 'index',
    path: '/index',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'index',
            i18n: 'Nav["首页"]'
        },
        auth: {
            view: null
        }
    }
}
