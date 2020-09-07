import Meta from '@/router/meta'
import {
    MENU_RESOURCE_INSTANCE,
    MENU_RESOURCE_MANAGEMENT,
    MENU_RESOURCE_INSTANCE_DETAILS
} from '@/dictionary/menu-symbol'

export default [{
    name: MENU_RESOURCE_INSTANCE,
    path: 'instance/:objId',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            relative: MENU_RESOURCE_MANAGEMENT
        },
        checkAvailable: (to, from, app) => {
            const modelId = to.params.objId
            const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
            return model && !model.bk_ispaused
        },
        layout: {}
    }),
    children: [{
        name: MENU_RESOURCE_INSTANCE_DETAILS,
        path: ':instId',
        component: () => import('./details.vue'),
        meta: new Meta({
            menu: {
                relative: MENU_RESOURCE_MANAGEMENT
            },
            checkAvailable: (to, from, app) => {
                const modelId = to.params.objId
                const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
                return model && !model.bk_ispaused
            },
            layout: {}
        })
    }]
}, {
    name: 'instanceHistory',
    path: 'history/instance/:objId',
    component: () => import('@/views/history/index.vue'),
    meta: new Meta({
        menu: {
            relative: MENU_RESOURCE_MANAGEMENT
        },
        layout: {},
        checkAvailable: (to, from, app) => {
            const modelId = to.params.objId
            const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
            return model && !model.bk_ispaused
        }
    })
}]
