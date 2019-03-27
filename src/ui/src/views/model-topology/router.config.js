import { NAV_MODEL_MANAGEMENT } from '@/dictionary/nav'
import {
    G_SYSTEM_TOPOLOGY,
    G_U_MODEL
} from '@/dictionary/auth'

export const OPERATION = {
    G_U_MODEL
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
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}
