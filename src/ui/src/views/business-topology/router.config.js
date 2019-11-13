import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_TRANSFER_HOST,
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_HOST_DETAILS
} from '@/dictionary/menu-symbol'

import {
    U_HOST,
    C_TOPO,
    R_TOPO,
    U_TOPO,
    D_TOPO,
    C_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    HOST_TO_RESOURCE
} from '@/dictionary/auth'

export default [{
    name: MENU_BUSINESS_HOST_AND_SERVICE,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '业务拓扑'
        },
        auth: {
            operation: {
                U_HOST,
                C_TOPO,
                R_TOPO,
                U_TOPO,
                D_TOPO,
                C_SERVICE_INSTANCE,
                R_SERVICE_INSTANCE,
                U_SERVICE_INSTANCE,
                D_SERVICE_INSTANCE,
                HOST_TO_RESOURCE
            }
        },
        customInstanceColumn: 'business_topology_table_column_config',
        customFilterProperty: 'business_topology_filter_property_config'
    })
}, {
    name: MENU_BUSINESS_TRANSFER_HOST,
    path: 'host/transfer/:type',
    component: () => import('@/views/host-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        }
    })
}, {
    name: MENU_BUSINESS_HOST_DETAILS,
    path: ':business/host/:id',
    component: () => import('@/views/host-details/index'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机详情',
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        },
        checkAvailable: (to, from, app) => {
            return parseInt(to.params.business) === app.$store.getters['objectBiz/bizId']
        }
    })
}]
