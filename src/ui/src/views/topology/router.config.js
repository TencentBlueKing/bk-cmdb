import { NAV_BUSINESS_RESOURCE } from '@/dictionary/nav'
import {
    B_C_TOPO,
    B_U_TOPO,
    B_D_TOPO,
    B_TOPO_TRANSFER_HOST,
    B_U_HOST,
    B_HOST_TO_RESOURCE
} from '@/dictionary/auth'

export const OPERATION = {
    B_C_TOPO,
    B_U_TOPO,
    B_D_TOPO,
    B_TOPO_TRANSFER_HOST,
    B_U_HOST,
    B_HOST_TO_RESOURCE
}

export default {
    name: 'topology',
    path: '/topology',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'topology',
            i18n: 'Nav["业务拓扑"]',
            parent: NAV_BUSINESS_RESOURCE
        },
        auth: {
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}
