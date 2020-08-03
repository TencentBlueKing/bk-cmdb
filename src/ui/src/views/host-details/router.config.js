import Meta from '@/router/meta'

export default [{
    name: 'host-landing',
    path: '/host-landing/:ip/:cloudId?',
    component: () => import('./landing.vue'),
    meta: new Meta()
}]
