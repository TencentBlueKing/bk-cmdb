import Meta from '@/router/meta'
import { NAV_SERVICE_MANAGEMENT } from '@/dictionary/menu'
import {
    C_SERVICE_TEMPLATE,
    U_SERVICE_TEMPLATE,
    D_SERVICE_TEMPLATE,
    R_SERVICE_TEMPLATE,
    
    C_PROCESS_TEMPLATE,
    U_PROCESS_TEMPLATE,
    D_PROCESS_TEMPLATE,
    R_PROCESS_TEMPLATE
} from '@/dictionary/auth'

export const OPERATION = {
    C_SERVICE_TEMPLATE,
    U_SERVICE_TEMPLATE,
    D_SERVICE_TEMPLATE,
    R_SERVICE_TEMPLATE,
    
    C_PROCESS_TEMPLATE,
    U_PROCESS_TEMPLATE,
    D_PROCESS_TEMPLATE,
    R_PROCESS_TEMPLATE
}

const path = '/service/template'

export default [{
    name: 'serviceTemplate',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'serviceTemplate',
            i18n: 'Nav["服务模板"]',
            path: path,
            order: 2,
            parent: NAV_SERVICE_MANAGEMENT,
            adminView: false
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        i18nTitle: "Nav['服务模板']"
    })
}, {
    name: 'operationalTemplate',
    path: '/service/operational/template/:templateId?',
    component: () => import('./children/operational.vue'),
    meta: {
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        }
    }
}]
