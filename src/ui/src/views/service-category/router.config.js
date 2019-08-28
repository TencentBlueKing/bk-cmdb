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

const path = '/service/cagetory'

export default {
    name: 'serviceCagetory',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            id: 'serviceCagetory',
            i18n: '服务分类',
            path: path,
            order: 3,
            parent: MENU_BUSINESS_SERVICE,
            adminView: false
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        i18nTitle: '服务分类',
        requireBusiness: true
    })
}
