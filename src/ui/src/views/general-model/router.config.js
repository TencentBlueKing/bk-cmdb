import Meta from '@/router/meta'
import { MENU_RESOURCE_INSTANCE, MENU_RESOURCE_MANAGEMENT } from '@/dictionary/menu-symbol'
import {
    C_INST,
    R_INST,
    U_INST,
    D_INST
} from '@/dictionary/auth'

export const OPERATION = {
    C_INST,
    R_INST,
    U_INST,
    D_INST
}

export default {
    name: MENU_RESOURCE_INSTANCE,
    path: 'instance/:objId',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            relative: MENU_RESOURCE_MANAGEMENT
        },
        auth: {
            operation: OPERATION,
            authScope: 'global',
            setDynamicMeta: (to, from, app) => {
                const modelId = to.params.objId
                const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
                if (model) {
                    app.$store.commit('auth/setParentMeta', {
                        parent_layers: [{
                            resource_type: 'model',
                            resource_id: model.id,
                            resource_model: modelId
                        }]
                    })
                }
            }
        },
        checkAvailable: (to, from, app) => {
            const modelId = to.params.objId
            const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
            return model && !model.bk_ispaused
        }
    })
}
