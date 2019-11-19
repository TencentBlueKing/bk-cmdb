import Meta from '@/router/meta'
import { MENU_BUSINESS, MENU_BUSINESS_CUSTOM_FIELDS } from '@/dictionary/menu-symbol'
import { U_MODEL } from '@/dictionary/auth'

export default {
    name: MENU_BUSINESS_CUSTOM_FIELDS,
    path: 'custom-fields',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '自定义字段'
        },
        auth: {
            operation: {
                U_MODEL
            }
        }
    })
}
