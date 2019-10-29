import Meta from '@/router/meta'

export default {
    name: 'hostOperation',
    path: 'host-operation/:type',
    component: () => import('./index.vue'),
    meta: new Meta({})
}