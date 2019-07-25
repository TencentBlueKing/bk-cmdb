import Meta from '@/router/meta'
import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
import {
    R_PROCESS,
    C_PROCESS,
    U_PROCESS,
    D_PROCESS,
    PROCESS_BIND_MODULE,
    PROCESS_UNBIND_MODULE
} from '@/dictionary/auth'

export const OPERATION = {
    R_PROCESS,
    C_PROCESS,
    U_PROCESS,
    D_PROCESS,
    PROCESS_BIND_MODULE,
    PROCESS_UNBIND_MODULE
}

const path = '/process'

export default {
    name: 'process',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'process',
            i18n: 'Nav["进程管理"]',
            path: path,
            order: 3,
            parent: NAV_BUSINESS_RESOURCE,
            adminView: false
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        requireBusiness: true,
        i18nTitle: 'Nav["进程管理"]'
    })
}
