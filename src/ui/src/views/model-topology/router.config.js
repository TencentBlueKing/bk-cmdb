import Meta from '@/router/meta'
import { MENU_MODEL_TOPOLOGY, MENU_MODEL_TOPOLOGY_NEW } from '@/dictionary/menu-symbol'

import { SYSTEM_MODEL_GRAPHICS, SYSTEM_TOPOLOGY } from '@/dictionary/auth'

export const OPERATION = {
    SYSTEM_MODEL_GRAPHICS,
    SYSTEM_TOPOLOGY
}

export default [{
    name: MENU_MODEL_TOPOLOGY,
    path: 'all/topology',
    component: () => import('./index.old.vue'),
    meta: new Meta({
        available: false,
        menu: {
            i18n: '模型拓扑'
        },
        auth: {
            operation: OPERATION,
            authScope: 'global'
        }
    })
}, {
    name: MENU_MODEL_TOPOLOGY_NEW,
    path: 'all/topology/new',
    component: () => import('./index.new.vue'),
    meta: new Meta({
        menu: {
            i18n: '模型关系'
        },
        auth: {
            operation: OPERATION,
            authScope: 'global'
        },
        layout: {
            breadcrumbs: false
        }
    })
}]
