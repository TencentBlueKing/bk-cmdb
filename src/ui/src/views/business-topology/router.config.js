import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_SERVICE,
    MENU_BUSINESS_SERVICE_TOPOLOGY
} from '@/dictionary/menu-symbol'
import {
    U_HOST,
    HOST_TO_RESOURCE,

    C_TOPO,
    U_TOPO,
    D_TOPO,

    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE,

    U_SERVICE_TEMPLATE
} from '@/dictionary/auth'

export const OPERATION = {
    U_HOST,
    HOST_TO_RESOURCE,

    C_TOPO,
    U_TOPO,
    D_TOPO,

    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE,

    U_SERVICE_TEMPLATE
}

export default {
    name: MENU_BUSINESS_SERVICE_TOPOLOGY,
    path: 'topology',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '服务拓扑',
            parent: MENU_BUSINESS_SERVICE
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        }
    })
}
