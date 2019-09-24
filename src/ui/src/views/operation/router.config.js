import { MENU_ANALYSIS_OPERATION } from '@/dictionary/menu-symbol'
import {
    C_STATISTICAL_REPORT,
    D_STATISTICAL_REPORT,
    R_STATISTICAL_REPORT,
    U_STATISTICAL_REPORT
} from '@/dictionary/auth'
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
                C_STATISTICAL_REPORT,
                D_STATISTICAL_REPORT,
                R_STATISTICAL_REPORT,
                U_STATISTICAL_REPORT
            },
            authScope: 'global'
        }
    })
}
