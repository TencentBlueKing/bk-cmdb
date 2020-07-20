import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_TRANSFER_HOST,
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_HOST_DETAILS,
    MENU_BUSINESS_DELETE_SERVICE
} from '@/dictionary/menu-symbol'

export default [{
    name: MENU_BUSINESS_HOST_AND_SERVICE,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '业务拓扑'
        },
        customInstanceColumn: 'business_topology_table_column_config',
        customFilterProperty: 'business_topology_filter_property_config'
    }),
    children: [{
        name: MENU_BUSINESS_HOST_DETAILS,
        path: 'host/:id',
        component: () => import('@/views/host-details/index'),
        meta: new Meta({
            owner: MENU_BUSINESS,
            menu: {
                i18n: '主机详情',
                relative: MENU_BUSINESS_HOST_AND_SERVICE
            }
        })
    }]
}, {
    name: MENU_BUSINESS_TRANSFER_HOST,
    path: 'host/transfer/:type/:module?',
    component: () => import('@/views/host-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        }
    })
}, {
    name: MENU_BUSINESS_DELETE_SERVICE,
    path: 'service/delete/:moduleId?/:ids',
    component: () => import('@/views/service-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '删除服务实例',
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        }
    })
}]
