import Meta from '@/router/meta'
import { MENU_BUSINESS_SERVICE, MENU_BUSINESS } from '@/dictionary/menu-symbol'
import {
    C_SERVICE_CATEGORY,
    U_SERVICE_CATEGORY,
    D_SERVICE_CATEGORY,
    R_SERVICE_CATEGORY
} from '@/dictionary/auth'

export const OPERATION = {
    C_SERVICE_CATEGORY,
    U_SERVICE_CATEGORY,
    D_SERVICE_CATEGORY,
    R_SERVICE_CATEGORY
}

export default {
    name: 'serviceCagetory',
    path: 'service/cagetory',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '服务分类',
            parent: MENU_BUSINESS_SERVICE
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        },
        requireBusiness: true
    })
}
