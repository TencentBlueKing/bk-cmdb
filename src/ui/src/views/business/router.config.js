import { NAV_COLLECT } from '@/dictionary/nav'
import {
    G_C_BUSINESS,
    G_U_BUSINESS,
    G_D_BUSINESS,
    G_R_BUSINESS
} from '@/dictionary/auth'

export const OPERATION = {
    G_C_BUSINESS,
    G_U_BUSINESS,
    G_D_BUSINESS
}

export default {
    name: 'business',
    path: '/business',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'business',
            i18n: 'Nav["业务"]',
            parent: NAV_COLLECT
        },
        auth: {
            view: G_R_BUSINESS,
            operation: Object.values(OPERATION)
        }
    }
}
