import { NAV_BUSINESS_RESOURCE } from '@/dictionary/nav'
import {
    B_C_CUSTOM_QUERY,
    B_U_CUSTOM_QUERY,
    B_D_CUSTOM_QUERY,
    B_R_CUSTOM_QUERY
} from '@/dictionary/auth'

export const OPERATION = {
    B_C_CUSTOM_QUERY,
    B_U_CUSTOM_QUERY,
    B_D_CUSTOM_QUERY
}

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
            operation: Object.values(OPERATION)
        }
    }
}
