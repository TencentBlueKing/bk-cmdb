import {
    D_R_MODEL
} from '@/dictionary/auth'

export default {
    name: 'history',
    path: '/history/:objId',
    component: () => import('./index.vue'),
    meta: {
        auth: {
            view: D_R_MODEL,
            operation: []
        },
        dynamicParams: ['objId']
    }
}
