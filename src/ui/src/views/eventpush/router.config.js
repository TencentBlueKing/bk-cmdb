import { NAV_MODEL_MANAGEMENT } from '@/types/nav'
import {
    G_C_EVENT,
    G_U_EVENT,
    G_D_EVENT,
    G_R_EVENT
} from '@/types/auth'

export default {
    name: 'eventpush',
    path: '/eventpush',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'eventpush',
            i18n: 'Nav["事件推送"]',
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            view: G_R_EVENT,
            operation: [
                G_C_EVENT,
                G_U_EVENT,
                G_D_EVENT
            ]
        }
    }
}
