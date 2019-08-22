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

const path = '/eventpush'

export default {
    name: 'eventpush',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'eventpush',
            i18n: '事件推送',
            path: path,
            order: 5,
            parent: NAV_MODEL_MANAGEMENT,
            businessView: false
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
