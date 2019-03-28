import {
    R_INST,
    GET_MODEL_INST_AUTH_META
} from '@/dictionary/auth'

export default {
    name: 'history',
    path: '/history/:objId',
    component: () => import('./index.vue'),
    meta: {
        auth: {
            view: R_INST,
            meta: GET_MODEL_INST_AUTH_META,
            operation: []
        }
    }
}
