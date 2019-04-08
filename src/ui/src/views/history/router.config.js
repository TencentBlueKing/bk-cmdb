export default {
    name: 'history',
    path: '/history/:objId',
    component: () => import('./index.vue'),
    meta: {
        auth: {
            view: '',
            operation: []
        },
        checkAvailable: (to, from, app) => {
            const modelId = to.params.objId
            const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
            return model && !model.bk_ispaused
        }
    }
}
