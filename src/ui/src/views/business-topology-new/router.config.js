import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_TRANSFER_HOST
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
    name: 'businessTopologyNew',
    path: 'topologyNew',
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
        customInstanceColumn: 'business_topology_table_column_config'
    })
}, {
    name: MENU_BUSINESS_TRANSFER_HOST,
    path: 'host/transfer/:type',
    component: () => import('@/views/host-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS
    })
}]
