import Meta from '@/router/meta'
import { MENU_RESOURCE_EVENTPUSH } from '@/dictionary/menu-symbol'
import { OPERATION, TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'

const [viewAuth] = TRANSFORM_TO_INTERNAL({ type: OPERATION.R_EVENT })
export default {
    name: MENU_RESOURCE_EVENTPUSH,
    path: 'eventpush',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '事件订阅'
        },
        auth: {
            view: viewAuth
        }
    })
}
