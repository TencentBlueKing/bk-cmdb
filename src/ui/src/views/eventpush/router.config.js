import Meta from '@/router/meta'
import { MENU_RESOURCE_EVENTPUSH } from '@/dictionary/menu-symbol'
import { OPERATION } from '@/dictionary/iam-auth'

export default {
    name: MENU_RESOURCE_EVENTPUSH,
    path: 'eventpush',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '事件订阅'
        },
        auth: {
            view: { type: OPERATION.R_EVENT }
        }
    })
}
