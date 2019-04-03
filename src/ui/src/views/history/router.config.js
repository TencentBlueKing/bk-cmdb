import {
    R_INST
} from '@/dictionary/auth'

export default {
    name: 'history',
    path: '/history/:objId',
    component: () => import('./index.vue'),
    meta: {
        auth: {
            view: R_INST,
            operation: [R_INST]
        }
    }
}
