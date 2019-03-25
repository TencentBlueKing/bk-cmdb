import { NAV_MODEL_MANAGEMENT } from '@/types/nav'
import {
    D_C_MODEL,
    D_R_MODEL,
    D_U_MODEL,
    D_D_MODEL
} from '@/types/auth'

export default {
    name: 'generalModel',
    path: '/general-model/:objId',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'generalModel',
            parent: NAV_COLLECT
        },
        auth: {
            view: D_R_MODEL,
            operation: [
                D_C_MODEL,
                D_U_MODEL,
                D_D_MODEL
            ]
        },
        dynamicParams: ['objId']
    }
}
