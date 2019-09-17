import Meta from '@/router/meta'
import {
    C_MODEL_GROUP,
    U_MODEL_GROUP,
    D_MODEL_GROUP,
    C_MODEL,
    U_MODEL,
    D_MODEL
} from '@/dictionary/auth'
import { MENU_BUSINESS_ADVANCED } from '@/dictionary/menu-symbol'

export const OPERATION = {
    C_MODEL_GROUP,
    U_MODEL_GROUP,
    D_MODEL_GROUP,
    C_MODEL,
    U_MODEL,
    D_MODEL
}

export default {
    name: 'customFields',
    path: 'custom-fields',
    component: () => import('./index.vue'),
    meta: new Meta({
        available: false,
        menu: {
            i18n: '自定义字段',
            parent: MENU_BUSINESS_ADVANCED
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        }
    })
}
