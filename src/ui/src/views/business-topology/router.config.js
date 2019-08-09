import Meta from '@/router/meta'
import { NAV_SERVICE_MANAGEMENT } from '@/dictionary/menu'
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
    R_SERVICE_INSTANCE
} from '@/dictionary/auth'

export const OPERATION = {
    U_HOST,
    HOST_TO_RESOURCE,

    C_TOPO,
    U_TOPO,
    D_TOPO,
    R_TOPO,

    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE
}

const path = '/topology'

export default {
    name: 'topology',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'topology',
            i18n: '服务拓扑',
            path: path,
            order: 1,
            parent: NAV_SERVICE_MANAGEMENT,
            adminView: false
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        i18nTitle: '服务拓扑',
        requireBusiness: true
    })
}
