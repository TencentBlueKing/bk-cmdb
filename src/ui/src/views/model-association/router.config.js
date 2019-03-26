import { NAV_MODEL_MANAGEMENT } from '@/dictionary/nav'
import {
    G_C_RELATION,
    G_U_RELATION,
    G_D_RELATION
} from '@/dictionary/auth'

export const OPERATION = {
    G_C_RELATION,
    G_U_RELATION,
    G_D_RELATION
}

export default {
    name: 'modelAssociation',
    path: '/model/association',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'modelAssociation',
            i18n: 'Nav["关联类型"]',
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}
