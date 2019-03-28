import { NAV_BASIC_RESOURCE } from '@/dictionary/menu'
import {
    G_C_HOST,
    G_U_HOST,
    G_D_HOST,
    G_R_HOST,
    G_HOST_ASSIGN
} from '@/dictionary/auth'

export const OPERATION = {
    G_C_HOST,
    G_U_HOST,
    G_D_HOST,
    G_HOST_ASSIGN
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
            parent: NAV_BASIC_RESOURCE
        },
        auth: {
            view: G_R_HOST,
            operation: Object.values(OPERATION)
        }
    }
}
