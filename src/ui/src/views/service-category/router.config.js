import Meta from '@/router/meta'
import { NAV_SERVICE_MANAGEMENT } from '@/dictionary/menu'
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
        menu: {
            id: 'serviceCagetory',
            i18n: '服务分类',
            path: path,
            order: 3,
            parent: NAV_SERVICE_MANAGEMENT,
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
