import Meta from '@/router/meta'
import { MENU_RESOURCE_INSTANCE, MENU_RESOURCE_MANAGEMENT } from '@/dictionary/menu-symbol'

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
        layout: {
            previous: {
                name: MENU_RESOURCE_MANAGEMENT
            }
        }
    })
}, {
    name: 'instanceHistory',
    path: 'instance/:objId/history',
    component: () => import('@/views/history/index.vue'),
    meta: new Meta({
        menu: {
            relative: MENU_RESOURCE_MANAGEMENT
        },
        layout: {
            previous () {
                return {
                    name: MENU_RESOURCE_INSTANCE,
                    params: {
                        objId: window.CMDB_APP.$route.params.objId
                    }
                }
            }
        },
        checkAvailable: (to, from, app) => {
            const modelId = to.params.objId
            const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
            return model && !model.bk_ispaused
        }
    })
}]
