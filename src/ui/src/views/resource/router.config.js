import Meta from '@/router/meta'
import {
    MENU_RESOURCE,
    MENU_RESOURCE_HOST,
    MENU_RESOURCE_MANAGEMENT,
    MENU_RESOURCE_HOST_DETAILS,
    MENU_RESOURCE_BUSINESS_HOST_DETAILS
} from '@/dictionary/menu-symbol'

export default [{
    name: MENU_RESOURCE_HOST,
    path: 'host',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '主机'
        },
        layout: {},
        filterPropertyKey: 'resource_host_filter_properties',
        customInstanceColumn: 'resource_host_table_column_config'
    }),
    children: [{
        name: MENU_RESOURCE_HOST_DETAILS,
        path: ':id',
        component: () => import('@/views/host-details/index'),
        meta: new Meta({
            owner: MENU_RESOURCE,
            menu: {
                i18n: '主机详情',
                relative: [MENU_RESOURCE_HOST, MENU_RESOURCE_MANAGEMENT]
            },
            layout: {}
        })
    }, {
        name: MENU_RESOURCE_BUSINESS_HOST_DETAILS,
        path: ':business/:id',
        component: () => import('@/views/host-details/index'),
        meta: new Meta({
            owner: MENU_RESOURCE,
            menu: {
                i18n: '主机详情',
                relative: [MENU_RESOURCE_HOST, MENU_RESOURCE_MANAGEMENT]
            },
            layout: {}
        })
    }]
}, {
    name: 'hostHistory',
    path: 'history/host',
    component: () => import('@/views/history/index.vue'),
    meta: new Meta({
        menu: {
            relative: MENU_RESOURCE_MANAGEMENT
        },
        layout: {}
    })
}]
