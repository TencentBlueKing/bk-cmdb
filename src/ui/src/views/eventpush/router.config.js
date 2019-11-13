import Meta from '@/router/meta'
import { MENU_RESOURCE_EVENTPUSH } from '@/dictionary/menu-symbol'
import {
    C_EVENT,
    U_EVENT,
    D_EVENT,
    R_EVENT,
    GET_AUTH_META
} from '@/dictionary/auth'
export default {
    name: MENU_RESOURCE_EVENTPUSH,
    path: 'eventpush',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '事件订阅'
        },
        auth: {
            view: { ...GET_AUTH_META(R_EVENT) },
            operation: {
                C_EVENT,
                U_EVENT,
                D_EVENT,
                R_EVENT
            }
        }
    })
}
