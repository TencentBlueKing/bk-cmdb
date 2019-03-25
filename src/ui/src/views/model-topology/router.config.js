import { NAV_MODEL_MANAGEMENT } from '@/types/nav'
import {
    G_SYSTEM_TOPOLOGY,
    G_C_RELATION,
    G_U_RELATION,
    G_D_RELATION
} from '@/types/auth'

export default {
    name: 'modelTopology',
    path: '/model/topology',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'modelTopology',
            i18n: 'Nav["模型拓扑"]',
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            view: G_SYSTEM_TOPOLOGY,
            operation: [
                G_C_RELATION,
                G_U_RELATION,
                G_D_RELATION
            ]
        }
    }
}
