import { NAV_MODEL_MANAGEMENT } from '@/types/nav'
import {
    G_C_MODEL_GROUP,
    G_U_MODEL_GROUP,
    G_D_MODEL_GROUP,
    G_C_MODEL,
    G_R_MODEL,
    G_U_MODEL,
    G_D_MODEL
} from '@/types/auth'

export default {
    name: 'model',
    path: '/model',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'model',
            i18n: 'Nav["模型"]',
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            view: G_R_MODEL,
            operation: [
                G_C_MODEL_GROUP,
                G_U_MODEL_GROUP,
                G_D_MODEL_GROUP,
                G_C_MODEL,
                G_U_MODEL,
                G_D_MODEL
            ]
        }
    }
}
