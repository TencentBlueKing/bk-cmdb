import { NAV_BUSINESS_RESOURCE } from '@/types/nav'
import {
    B_C_CUSTOM_QUERY,
    B_U_CUSTOM_QUERY,
    B_D_CUSTOM_QUERY,
    B_R_CUSTOM_QUERY
} from '@/types/auth'

export default {
    name: 'customQuery',
    path: '/custom-query',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'customQuery',
            i18n: 'Nav["动态分组"]',
            parent: NAV_BUSINESS_RESOURCE
        },
        auth: {
            view: B_R_CUSTOM_QUERY,
            operation: [
                B_C_CUSTOM_QUERY,
                B_U_CUSTOM_QUERY,
                B_D_CUSTOM_QUERY
            ]
        }
    }
}
