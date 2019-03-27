import { NAV_BUSINESS_RESOURCE } from '@/dictionary/nav'
import {
    B_R_PROCESS,
    B_C_PROCESS,
    B_U_PROCESS,
    B_D_PROCESS,
    B_PROCESS_BIND_MODULE,
    B_PROCESS_UNBIND_MODULE
} from '@/dictionary/auth'

export const OPERATION = {
    B_R_PROCESS,
    B_C_PROCESS,
    B_U_PROCESS,
    B_D_PROCESS,
    B_PROCESS_BIND_MODULE,
    B_PROCESS_UNBIND_MODULE
}

export default {
    name: 'process',
    path: '/process',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'process',
            i18n: 'Nav["进程管理"]',
            parent: NAV_BUSINESS_RESOURCE
        },
        auth: {
            view: B_R_PROCESS,
            operation: Object.values(OPERATION)
        }
    }
}
