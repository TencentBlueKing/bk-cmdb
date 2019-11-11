import Meta from '@/router/meta'
import { U_MODEL } from '@/dictionary/auth'
import { MENU_BUSINESS_CUSTOM_FIELDS } from '@/dictionary/menu-symbol'

export const OPERATION = { U_MODEL }

export default {
    name: MENU_BUSINESS_CUSTOM_FIELDS,
    path: 'custom-fields',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '自定义字段'
        },
        auth: {
            operation: { U_MODEL },
            authScope: 'business'
        }
    })
}
