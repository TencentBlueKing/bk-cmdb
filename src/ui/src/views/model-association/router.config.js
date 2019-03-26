import { NAV_MODEL_MANAGEMENT } from '@/dictionary/nav'
import {
    G_C_RELATION,
    G_U_RELATION,
    G_D_RELATION
} from '@/dictionary/auth'

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
            operation: [
                G_C_RELATION,
                G_U_RELATION,
                G_D_RELATION
            ]
        }
    }
}
