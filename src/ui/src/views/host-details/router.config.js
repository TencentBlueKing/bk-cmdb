export default {
    name: 'hostDetails',
    path: '/host/:business?/:id',
    component: () => import(/* webpackChunkName: "hostDetails" */ './index.vue'),
    meta: {
        auth: {
            view: null,
            operation: []
        }
    }
}
