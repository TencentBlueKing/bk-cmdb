import { NAV_COLLECT } from '@/dictionary/nav'
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

export default {
    name: 'resource',
    path: '/resource',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'resource',
            i18n: 'Nav["主机"]',
            parent: NAV_COLLECT
        },
        auth: {
            view: G_R_HOST,
            operation: Object.values(OPERATION)
        }
    }
}
