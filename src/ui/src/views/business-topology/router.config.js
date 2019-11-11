import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_SERVICE,
    MENU_BUSINESS_HOST_AND_SERVICE
} from '@/dictionary/menu-symbol'

export default {
    name: MENU_BUSINESS_HOST_AND_SERVICE,
    path: 'topology',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '服务拓扑',
            parent: MENU_BUSINESS_SERVICE
        }
    })
}
