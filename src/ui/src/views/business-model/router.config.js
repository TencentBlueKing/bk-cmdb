import { NAV_MODEL_MANAGEMENT } from '@/dictionary/nav'
import { G_SYSTEM_TOPOLOGY } from '@/dictionary/auth'

export const OPERATION = {
    G_SYSTEM_TOPOLOGY
}

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
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}
