import { NAV_BASIC_RESOURCE } from '@/dictionary/menu'
import {
    C_BUSINESS,
    U_BUSINESS,
    R_BUSINESS,
    BUSINESS_ARCHIVE
} from '@/dictionary/auth'

export const OPERATION = {
    R_BUSINESS,
    C_BUSINESS,
    U_BUSINESS,
    BUSINESS_ARCHIVE
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
            parent: NAV_BASIC_RESOURCE,
            adminView: true
        },
        auth: {
            view: R_BUSINESS,
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
            operation: [BUSINESS_ARCHIVE]
        }
    }
}]
