import { NAV_BASIC_RESOURCE } from '@/dictionary/menu'
import {
    G_C_BUSINESS,
    G_U_BUSINESS,
    G_R_BUSINESS,
    G_BUSINESS_ARCHIVE
} from '@/dictionary/auth'

export const OPERATION = {
    G_C_BUSINESS,
    G_U_BUSINESS,
    G_BUSINESS_ARCHIVE
}

const businessPath = '/business'
const historyPath = '/history/biz'

export default [{
    name: 'business',
    path: businessPath,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'business',
            i18n: 'Nav["业务"]',
            path: businessPath,
            order: 1,
            parent: NAV_BASIC_RESOURCE
        },
        auth: {
            view: G_R_BUSINESS,
            operation: Object.values(OPERATION)
        }
    }
}, {
    name: 'businessHistory',
    path: historyPath,
    component: () => import('./archived.vue'),
    meta: {
        auth: {
            view: '',
            operation: [G_BUSINESS_ARCHIVE]
        }
    }
}]
