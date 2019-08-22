import Meta from '@/router/meta'
import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
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

const path = '/custom-query'

export default {
    name: 'customQuery',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'customQuery',
            i18n: '动态分组',
            path: path,
            order: 4,
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
        i18nTitle: '动态分组'
    })
}
