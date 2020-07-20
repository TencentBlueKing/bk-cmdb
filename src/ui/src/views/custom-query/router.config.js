import Meta from '@/router/meta'
import { MENU_BUSINESS, MENU_BUSINESS_CUSTOM_QUERY } from '@/dictionary/menu-symbol'
import { OPERATION, TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'

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
                const bizId = to.params.bizId
                const [viewAuth] = TRANSFORM_TO_INTERNAL({ type: OPERATION.R_CUSTOM_QUERY, relation: [bizId] })
                return viewAuth
            }
        }
    })
}
