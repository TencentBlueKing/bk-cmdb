import Meta from '@/router/meta'

const path = '/synchronous/set/:setId/module/:moduleId'

export default {
    name: 'synchronous',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({})
}
