import { MENU_ANALYSIS_OPERATION } from '@/dictionary/menu-symbol'
import Meta from '@/router/meta'

export default {
    name: MENU_ANALYSIS_OPERATION,
    path: 'operation',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '运营统计'
        }
    })
}
