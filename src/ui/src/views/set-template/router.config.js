import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_SERVICE,
    MENU_BUSINESS_SET_TEMPLATE
} from '@/dictionary/menu-symbol'

import {
    C_SET_TEMPLATE,
    U_SET_TEMPLATE,
    D_SET_TEMPLATE
} from '@/dictionary/auth'

export const OPERATION = {
    C_SET_TEMPLATE,
    U_SET_TEMPLATE,
    D_SET_TEMPLATE
}

export default [{
    name: MENU_BUSINESS_SET_TEMPLATE,
    path: 'set/template',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '集群模板',
            parent: MENU_BUSINESS_SERVICE
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        }
    })
}, {
    name: 'setTemplateConfig',
    path: 'set/template/:mode/:templateId?',
    component: () => import('./template.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '集群模板',
            relative: MENU_BUSINESS_SET_TEMPLATE
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        }
    })
}, {
    name: 'syncHistory',
    path: 'set/template/history/:templateId?',
    component: () => import('./sync-history.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步历史',
            relative: MENU_BUSINESS_SET_TEMPLATE
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        }
    })
}]
