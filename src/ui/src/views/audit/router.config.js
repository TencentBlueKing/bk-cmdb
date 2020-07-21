import { MENU_ANALYSIS_AUDIT } from '@/dictionary/menu-symbol'
import { OPERATION } from '@/dictionary/iam-auth'
import Meta from '@/router/meta'
export default {
    name: MENU_ANALYSIS_AUDIT,
    path: 'audit',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '操作审计'
        },
        auth: {
            view: { type: OPERATION.R_AUDIT }
        }
    })
}
