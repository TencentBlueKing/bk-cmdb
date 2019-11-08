import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_TRANSFER_HOST,
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_HOST_DETAILS
} from '@/dictionary/menu-symbol'

import {
    U_HOST,
    HOST_TO_RESOURCE,

    C_TOPO,
    U_TOPO,
    D_TOPO,
    R_TOPO,

    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE,

    U_SERVICE_TEMPLATE
} from '@/dictionary/auth'

const OPERATION = {
    U_HOST,
    HOST_TO_RESOURCE,

    C_TOPO,
    U_TOPO,
    D_TOPO,
    R_TOPO,

    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE,

    U_SERVICE_TEMPLATE
}

export default [{
    name: MENU_BUSINESS_HOST_AND_SERVICE,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        auth: {
            operation: OPERATION,
            authScope: 'business'
        },
        menu: {
            i18n: '业务拓扑'
        },
        customInstanceColumn: 'business_topology_table_column_config',
        customFilterProperty: 'business_topology_filter_property_config'
    })
}, {
    name: MENU_BUSINESS_TRANSFER_HOST,
    path: 'host/transfer/:type',
    component: () => import('@/views/host-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS
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
        auth: {
            view: null,
            operation: { U_HOST, D_SERVICE_INSTANCE },
            authScope: 'business'
        },
        checkAvailable: (to, from, app) => {
            return parseInt(to.params.business) === app.$store.getters['objectBiz/bizId']
        }
    })
}]
