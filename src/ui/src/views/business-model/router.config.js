import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'
import { SYSTEM_TOPOLOGY } from '@/dictionary/auth'

export const OPERATION = {
    SYSTEM_TOPOLOGY
}

const path = '/model/business'

export default {
    name: 'businessModel',
    path: path,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'businessModel',
            i18n: 'Nav["业务模型"]',
            path: path,
            order: 3,
            parent: NAV_MODEL_MANAGEMENT,
            adminView: true
        },
        auth: {
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}
