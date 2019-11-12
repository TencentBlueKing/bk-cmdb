import Meta from '@/router/meta'
import {
    MENU_BUSINESS_HOST,
    MENU_BUSINESS,
    MENU_BUSINESS_HOST_MANAGEMENT,
    MENU_BUSINESS_HOST_DETAILS
} from '@/dictionary/menu-symbol'

export default [{
    name: MENU_BUSINESS_HOST_MANAGEMENT,
    path: 'host',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '业务主机',
            parent: MENU_BUSINESS_HOST
        },
        filterPropertyKey: 'business_host_filter_properties'
    })
}, {
    name: MENU_BUSINESS_HOST_DETAILS,
    path: ':business/host/:id',
    component: () => import('@/views/host-details/index'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机详情',
            relative: MENU_BUSINESS_HOST_MANAGEMENT
        },
        checkAvailable: (to, from, app) => {
            return parseInt(to.params.business) === app.$store.getters['objectBiz/bizId']
        }
    })
}]
