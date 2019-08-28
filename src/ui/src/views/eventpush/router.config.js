import Meta from '@/router/meta'
import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'
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
    name: 'eventpush',
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '事件推送',
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            view: OPERATION.R_EVENT,
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'global'
            }
        },
        i18nTitle: '事件推送'
    })
}
