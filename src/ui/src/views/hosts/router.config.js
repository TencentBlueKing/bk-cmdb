import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
import {
    B_R_HOST,
    B_U_HOST,
    B_HOST_TO_RESOURCE,
    B_TOPO_TRANSFER_HOST
} from '@/dictionary/auth'

export const OPERATION = {
    B_U_HOST,
    B_HOST_TO_RESOURCE,
    B_TOPO_TRANSFER_HOST
}

const path = '/hosts'

export default {
    name: 'hosts',
    path: path,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'hosts',
            i18n: 'Nav["主机查询"]',
            path: path,
            order: 1,
            parent: NAV_BUSINESS_RESOURCE
        },
        auth: {
            view: B_R_HOST,
            operation: Object.values(OPERATION)
        }
    }
}
