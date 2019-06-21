import Meta from '@/router/meta'
import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'

import {
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
    meta: new Meta({
        menu: {
            id: 'modelTopology',
            i18n: 'Nav["模型拓扑"]',
            path: path,
            order: 2,
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            operation: Object.values(OPERATION)
        },
        i18nTitle: 'Nav["模型拓扑"]'
    })
}
