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
            i18n: '事件推送'
        },
        auth: {
            view: OPERATION.R_EVENT,
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'global'
            }
        }
    })
}
