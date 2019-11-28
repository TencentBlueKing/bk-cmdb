import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_SERVICE,
    MENU_BUSINESS_SERVICE_TEMPLATE
} from '@/dictionary/menu-symbol'
import {
    C_SERVICE_TEMPLATE,
    U_SERVICE_TEMPLATE,
    D_SERVICE_TEMPLATE,
    R_SERVICE_TEMPLATE
} from '@/dictionary/auth'

export default [{
    name: MENU_BUSINESS_SERVICE_TEMPLATE,
    path: 'service/template',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '服务模板',
            parent: MENU_BUSINESS_SERVICE
        },
        auth: {
            operation: {
                C_SERVICE_TEMPLATE,
                U_SERVICE_TEMPLATE,
                D_SERVICE_TEMPLATE,
                R_SERVICE_TEMPLATE
            }
        }
    })
}, {
    name: 'operationalTemplate',
    path: 'service/operational/template/:templateId?',
    component: () => import('./template.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            relative: MENU_BUSINESS_SERVICE_TEMPLATE
        },
        auth: {
            operation: {
                C_SERVICE_TEMPLATE,
                U_SERVICE_TEMPLATE,
                D_SERVICE_TEMPLATE,
                R_SERVICE_TEMPLATE
            }
        }
    })
}]
