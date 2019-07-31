import index from './index.vue'
export default {
    name: 'tree',
    path: '/tree',
    component: index,
    meta: {
        auth: {
            view: null,
            operation: []
        }
    }
}
