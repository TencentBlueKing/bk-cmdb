import { NAV_MODEL_MANAGEMENT } from '@/types/nav'
import { G_SYSTEM_TOPOLOGY } from '@/types/auth'

export default {
    name: 'businessModel',
    path: '/model/business',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'businessModel',
            i18n: 'Nav["业务模型"]',
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            view: G_SYSTEM_TOPOLOGY,
            operation: []
        }
    }
}
