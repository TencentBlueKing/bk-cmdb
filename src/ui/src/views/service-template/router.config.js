import Meta from '@/router/meta'
import { MENU_BUSINESS_SERVICE, MENU_BUSINESS } from '@/dictionary/menu-symbol'
import {
    C_SERVICE_TEMPLATE,
    U_SERVICE_TEMPLATE,
    D_SERVICE_TEMPLATE,
    R_SERVICE_TEMPLATE
} from '@/dictionary/auth'

export const OPERATION = {
    C_SERVICE_TEMPLATE,
    U_SERVICE_TEMPLATE,
    D_SERVICE_TEMPLATE,
    R_SERVICE_TEMPLATE
}

export default [{
    name: 'serviceTemplate',
    path: '/business/:business/service/template',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '服务模板',
            parent: MENU_BUSINESS_SERVICE
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        i18nTitle: '服务模板',
        requireBusiness: true
    })
}, {
    name: 'operationalTemplate',
    path: '/business/:business/service/operational/template/:templateId?',
    component: () => import('./children/operational.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        requireBusiness: true
    })
}]
