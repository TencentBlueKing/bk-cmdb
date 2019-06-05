import Meta from '@/router/meta'

const path = '/synchronous/module/:moduleId/template/:templateId'

export default {
    name: 'synchronous',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({})
}
