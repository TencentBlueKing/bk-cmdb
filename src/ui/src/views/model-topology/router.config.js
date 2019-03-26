import { NAV_MODEL_MANAGEMENT } from '@/dictionary/nav'
import {
    G_SYSTEM_TOPOLOGY,
    G_C_RELATION,
    G_U_RELATION,
    G_D_RELATION
} from '@/dictionary/auth'

export const OPERATION = {
    G_C_RELATION,
    G_U_RELATION,
    G_D_RELATION
}

export default {
    name: 'modelTopology',
    path: '/model/topology',
    component: () => import('./index.old.vue'),
    meta: {
        menu: {
            id: 'modelTopology',
            i18n: 'Nav["模型拓扑"]',
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            view: G_SYSTEM_TOPOLOGY,
            operation: Object.values(OPERATION)
        }
    }
}
