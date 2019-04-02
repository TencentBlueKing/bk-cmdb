import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
import {
    C_TOPO,
    U_TOPO,
    D_TOPO,
    TOPO_TRANSFER_HOST,
    U_HOST,
    HOST_TO_RESOURCE
} from '@/dictionary/auth'

export const OPERATION = {
    C_TOPO,
    U_TOPO,
    D_TOPO,
    TOPO_TRANSFER_HOST,
    U_HOST,
    HOST_TO_RESOURCE
}

const path = '/topology'

export default {
    name: 'topology',
    path: path,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'topology',
            i18n: 'Nav["业务拓扑"]',
            path: path,
            order: 2,
            parent: NAV_BUSINESS_RESOURCE,
            adminView: false
        },
        auth: {
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}
