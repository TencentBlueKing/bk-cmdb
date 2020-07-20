import Meta from '@/router/meta'
import { MENU_BUSINESS, MENU_BUSINESS_SERVICE_CATEGORY } from '@/dictionary/menu-symbol'

export default {
    name: MENU_BUSINESS_SERVICE_CATEGORY,
    path: 'service/cagetory',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '服务分类'
        }
    })
}
