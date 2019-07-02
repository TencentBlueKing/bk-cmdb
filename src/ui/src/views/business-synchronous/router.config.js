import Meta from '@/router/meta'

export default {
    name: 'synchronous',
    path: '/synchronous/module/:moduleId/set/:setId',
    component: () => import('./index.vue'),
    meta: new Meta({})
}
