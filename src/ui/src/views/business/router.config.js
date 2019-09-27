import Meta from '@/router/meta'
import {
    MENU_RESOURCE_BUSINESS,
    MENU_RESOURCE_BUSINESS_HISTORY,
    MENU_RESOURCE_MANAGEMENT
} from '@/dictionary/menu-symbol'
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

export default [{
    name: MENU_RESOURCE_BUSINESS,
    path: 'business',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '业务',
            relative: MENU_RESOURCE_MANAGEMENT
        },
        auth: {
            operation: OPERATION,
            authScope: 'global'
        }
    })
}, {
    name: MENU_RESOURCE_BUSINESS_HISTORY,
    path: 'history/biz',
    component: () => import('./archived.vue'),
    meta: new Meta({
        menu: {
            i18n: '已归档业务'
        },
        auth: {
            view: { BUSINESS_ARCHIVE },
            operation: { BUSINESS_ARCHIVE },
            authScope: 'global'
        },
        checkAvailable: (to, from, app) => {
            return app.$store.getters.isAdminView
        }
    })
}]
