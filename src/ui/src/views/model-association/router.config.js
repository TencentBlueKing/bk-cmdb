import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'
import {
    C_RELATION,
    U_RELATION,
    D_RELATION
} from '@/dictionary/auth'

export const OPERATION = {
    C_RELATION,
    U_RELATION,
    D_RELATION
}

const path = '/model/association'

export default {
    name: 'modelAssociation',
    path: path,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'modelAssociation',
            i18n: 'Nav["关联类型"]',
            path: path,
            order: 4,
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}
