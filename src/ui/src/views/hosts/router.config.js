import Meta from '@/router/meta'
import {
    MENU_BUSINESS_HOST,
    MENU_BUSINESS,
    MENU_BUSINESS_HOST_MANAGEMENT,
    MENU_BUSINESS_HOST_DETAILS
} from '@/dictionary/menu-symbol'
import {
    U_HOST,
    D_SERVICE_INSTANCE,
    HOST_TO_RESOURCE,
    GET_AUTH_META
} from '@/dictionary/auth'

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
        auth: {
            operation: {
                U_HOST,
                HOST_TO_RESOURCE
            },
            authScope: 'business'
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
        auth: {
            view: null,
            operation: { U_HOST, D_SERVICE_INSTANCE },
            setDynamicMeta (to, from, app) {
                const hostMeta = GET_AUTH_META(U_HOST)
                const serviceInstanceMeta = GET_AUTH_META(D_SERVICE_INSTANCE)
                app.$store.commit('auth/setResourceMeta', [{
                    ...hostMeta,
                    resource_id: parseInt(to.params.id),
                    bk_biz_id: parseInt(to.params.business)
                }, {
                    ...serviceInstanceMeta,
                    resource_id: parseInt(to.params.id),
                    bk_biz_id: parseInt(to.params.business)
                }])
            },
            authScope: 'business'
        }
    })
}]
