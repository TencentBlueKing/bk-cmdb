import { NAV_BUSINESS_RESOURCE } from '@/dictionary/nav'
import {
    B_R_HOST,
    B_U_HOST,
    B_HOST_TO_RESOURCE
} from '@/dictionary/auth'

export default {
    name: 'hosts',
    path: '/hosts',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'hosts',
            i18n: 'Nav["主机查询"]',
            parent: NAV_BUSINESS_RESOURCE
        },
        auth: {
            view: B_R_HOST,
            operation: [
                B_U_HOST,
                B_HOST_TO_RESOURCE
            ]
        }
    }
}
