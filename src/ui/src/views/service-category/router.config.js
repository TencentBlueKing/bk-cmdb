import Meta from '@/router/meta'
import { MENU_BUSINESS, MENU_BUSINESS_SERVICE_CATEGORY } from '@/dictionary/menu-symbol'
import {
    C_SERVICE_CATEGORY,
    U_SERVICE_CATEGORY,
    D_SERVICE_CATEGORY,
    R_SERVICE_CATEGORY
} from '@/dictionary/auth'

export default {
    name: MENU_BUSINESS_SERVICE_CATEGORY,
    path: 'service/cagetory',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '服务分类'
        },
        auth: {
            operation: {
                C_SERVICE_CATEGORY,
                U_SERVICE_CATEGORY,
                D_SERVICE_CATEGORY,
                R_SERVICE_CATEGORY
            }
        }
    })
}
