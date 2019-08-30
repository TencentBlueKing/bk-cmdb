import Meta from '@/router/meta'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'

export default {
    name: 'synchronous',
    path: '/business/:business/synchronous/module/:moduleId/set/:setId',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS
    })
}
