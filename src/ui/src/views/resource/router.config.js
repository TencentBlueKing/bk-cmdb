import Meta from '@/router/meta'
import { MENU_RESOURCE_HOST, MENU_RESOURCE_MANAGEMENT } from '@/dictionary/menu-symbol'

export default [{
    name: MENU_RESOURCE_HOST,
    path: 'host',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '主机'
        },
        filterPropertyKey: 'resource_host_filter_properties'
    })
}, {
    name: 'hostHistory',
    path: 'host/history',
    component: () => import('@/views/history/index.vue'),
    meta: new Meta({
        menu: {
            relative: MENU_RESOURCE_MANAGEMENT
        }
    })
}]
