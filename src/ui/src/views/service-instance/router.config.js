import Meta from '@/router/meta'

export default [{
    name: 'createServiceInstance',
    path: '/service/instance/create/set/:setId/module/:moduleId',
    component: () => import('./create.vue'),
    meta: new Meta({
        resetMenu: false
    })
}, {
    name: 'cloneServiceInstance',
    path: '/service/instance/clone/set/:setId/module/:moduleId/instance/:instanceId/host/:hostId',
    props: true,
    component: () => import('./clone.vue'),
    meta: new Meta({
        resetMenu: false
    })
}]
