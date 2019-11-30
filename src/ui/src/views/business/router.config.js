import Meta from '@/router/meta'
import {
    MENU_RESOURCE_BUSINESS,
    MENU_RESOURCE_BUSINESS_HISTORY,
    MENU_RESOURCE_MANAGEMENT
} from '@/dictionary/menu-symbol'

import {
    BUSINESS_ARCHIVE,
    GET_AUTH_META
} from '@/dictionary/auth'

export default [{
    name: MENU_RESOURCE_BUSINESS,
    path: 'business',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '业务',
            relative: MENU_RESOURCE_MANAGEMENT
        }
    })
}, {
    name: MENU_RESOURCE_BUSINESS_HISTORY,
    path: 'business/history',
    component: () => import('./archived.vue'),
    meta: new Meta({
        menu: {
            i18n: '已归档业务',
            relative: MENU_RESOURCE_MANAGEMENT
        },
        auth: {
            view: { ...GET_AUTH_META(BUSINESS_ARCHIVE) },
            operation: {
                BUSINESS_ARCHIVE
            }
        }
    })
}]
