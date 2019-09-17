import Meta from '@/router/meta'
import { MENU_RESOURCE_EVENTPUSH } from '@/dictionary/menu-symbol'
import {
    C_EVENT,
    U_EVENT,
    D_EVENT,
    R_EVENT
} from '@/dictionary/auth'

export const OPERATION = {
    C_EVENT,
    R_EVENT,
    U_EVENT,
    D_EVENT
}

export default {
    name: MENU_RESOURCE_EVENTPUSH,
    path: 'eventpush',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '事件订阅'
        },
        auth: {
            view: { R_EVENT },
            operation: OPERATION,
            authScope: 'global'
        }
    })
}
