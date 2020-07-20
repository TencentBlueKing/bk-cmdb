import { MENU_ANALYSIS_OPERATION } from '@/dictionary/menu-symbol'
import { OPERATION, TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'
import Meta from '@/router/meta'
const [viewAuth] = TRANSFORM_TO_INTERNAL({ type: OPERATION.R_STATISTICAL_REPORT })
export default {
    name: MENU_ANALYSIS_OPERATION,
    path: 'operation',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '运营统计'
        },
        auth: {
            view: viewAuth
        }
    })
}
