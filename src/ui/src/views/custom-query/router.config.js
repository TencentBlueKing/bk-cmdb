import Meta from '@/router/meta'
import { MENU_BUSINESS, MENU_BUSINESS_CUSTOM_QUERY } from '@/dictionary/menu-symbol'
import {
    C_CUSTOM_QUERY,
    U_CUSTOM_QUERY,
    D_CUSTOM_QUERY,
    R_CUSTOM_QUERY,
    GET_AUTH_META
} from '@/dictionary/auth'

export default {
    name: MENU_BUSINESS_CUSTOM_QUERY,
    path: 'custom-query',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '动态分组'
        },
        auth: {
            view: (to, app) => {
                const bizId = app.$store.getters['objectBiz/bizId']
                return {
                    bk_biz_id: bizId,
                    ...GET_AUTH_META(R_CUSTOM_QUERY)
                }
            },
            operation: {
                C_CUSTOM_QUERY,
                U_CUSTOM_QUERY,
                D_CUSTOM_QUERY,
                R_CUSTOM_QUERY
            }
        }
    })
}
