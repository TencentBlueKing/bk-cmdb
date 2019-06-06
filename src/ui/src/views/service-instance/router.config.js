import Meta from '@/router/meta'

export default [{
    name: 'createServiceInstance',
    path: '/service/instance/create/set/:setId/module/:moduleId',
    component: () => import('./create.vue'),
    meta: new Meta()
}]
