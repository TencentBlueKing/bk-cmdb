import { R_AUDIT } from '@/dictionary/auth'
import Meta from '@/router/meta'
import { MENU_ANALYSIS_AUDIT } from '@/dictionary/menu-symbol'

export default {
    name: MENU_ANALYSIS_AUDIT,
    path: 'audit',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '操作审计'
        },
        auth: {
            view: { R_AUDIT },
            authScope: 'global'
        }
    })
}
