import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
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

const path = '/custom-query'

export default {
    name: 'customQuery',
    path: path,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'customQuery',
            i18n: 'Nav["动态分组"]',
            path: path,
            order: 4,
            parent: NAV_BUSINESS_RESOURCE
        },
        auth: {
            view: B_R_CUSTOM_QUERY,
            operation: Object.values(OPERATION)
        }
    }
}
