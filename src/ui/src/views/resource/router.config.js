import { NAV_BASIC_RESOURCE } from '@/dictionary/menu'
import {
    C_HOST,
    U_HOST,
    D_HOST,
    R_HOST,
    HOST_ASSIGN
} from '@/dictionary/auth'

export const OPERATION = {
    C_HOST,
    R_HOST,
    U_HOST,
    D_HOST,
    HOST_ASSIGN
}

const path = '/resource'

export default {
    name: 'resource',
    path: path,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'resource',
            i18n: 'Nav["主机"]',
            path: path,
            order: 2,
            parent: NAV_BASIC_RESOURCE,
            adminView: true
        },
        auth: {
            view: '',
            operation: Object.values(OPERATION)
        },
        requireBusiness: true
    }
}
