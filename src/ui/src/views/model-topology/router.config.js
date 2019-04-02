import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'
import {
    SYSTEM_TOPOLOGY,
    U_MODEL
} from '@/dictionary/auth'

export const OPERATION = {
    U_MODEL
}

const path = '/model/topology'

export default {
    name: 'modelTopology',
    path: path,
    component: () => import('./index.old.vue'),
    meta: {
        menu: {
            id: 'modelTopology',
            i18n: 'Nav["模型拓扑"]',
            path: path,
            order: 2,
            parent: NAV_MODEL_MANAGEMENT,
            adminView: true
        },
        auth: {
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}
