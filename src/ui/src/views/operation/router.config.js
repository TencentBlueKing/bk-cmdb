import { MENU_ANALYSIS_OPERATION } from '@/dictionary/menu-symbol'
import { R_STATISTICAL_REPORT, GET_AUTH_META } from '@/dictionary/auth'
import Meta from '@/router/meta'

export default {
    name: MENU_ANALYSIS_OPERATION,
    path: 'operation',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '运营统计'
        },
        auth: {
            view: {
                ...GET_AUTH_META(R_STATISTICAL_REPORT)
            }
        }
    })
}
