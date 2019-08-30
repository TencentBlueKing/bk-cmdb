import Meta from '@/router/meta'
import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'

import { SYSTEM_MODEL_GRAPHICS } from '@/dictionary/auth'

export const OPERATION = {
    SYSTEM_MODEL_GRAPHICS
}

const path = '/model/topology'

export default {
    name: 'modelTopology',
    path: path,
    component: () => import('./index.old.vue'),
    meta: new Meta({
        menu: {
            id: 'modelTopology',
            i18n: '模型拓扑',
            path: path,
            order: 2,
            parent: NAV_MODEL_MANAGEMENT,
            businessView: false
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'global'
            }
        },
        i18nTitle: '模型拓扑'
            
    })
}
