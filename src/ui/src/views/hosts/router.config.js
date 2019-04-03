import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
import {
    R_HOST,
    U_HOST,
    HOST_TO_RESOURCE
} from '@/dictionary/auth'

export const OPERATION = {
    U_HOST,
    R_HOST,
    HOST_TO_RESOURCE
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
            parent: NAV_BUSINESS_RESOURCE,
            adminView: false
        },
        auth: {
            view: R_HOST,
            operation: Object.values(OPERATION)
        },
        requireBusiness: true
    }
}
