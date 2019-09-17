import { MENU_RESOURCE_MANAGEMENT } from '@/dictionary/menu-symbol'
import Meta from '@/router/meta'
export default {
    name: MENU_RESOURCE_MANAGEMENT,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '资源目录'
        },
        custom: {
            resourceCollection: 'resourceCollection',
            hostCollected: 'hostCollected',
            businessCollected: 'businessCollected'
        }
    })
}
