import Meta from '@/router/meta'
import { MENU_BUSINESS, MENU_BUSINESS_HOST } from '@/dictionary/menu-symbol'
import {
    C_CUSTOM_QUERY,
    U_CUSTOM_QUERY,
    D_CUSTOM_QUERY,
    R_CUSTOM_QUERY
} from '@/dictionary/auth'

export const OPERATION = {
    R_CUSTOM_QUERY,
    C_CUSTOM_QUERY,
    U_CUSTOM_QUERY,
    D_CUSTOM_QUERY
}

export default {
    name: 'customQuery',
    path: 'custom-query',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '动态分组',
            parent: MENU_BUSINESS_HOST
        },
        auth: {
            view: {
                R_CUSTOM_QUERY
            },
            operation: OPERATION,
            authScope: 'business'
        },
        requireBusiness: true
    })
}
