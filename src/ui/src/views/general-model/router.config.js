import Meta from '@/router/meta'
import { getMetadataBiz } from '@/utils/tools'
import { MENU_RESOURCE_INSTANCE } from '@/dictionary/menu-symbol'
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
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope (to, from, app) {
                const modelId = to.params.objId
                const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
                const bizId = getMetadataBiz(model)
                this.authScope = bizId ? 'business' : 'global'
            },
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
                    const bizId = getMetadataBiz(model)
                    if (bizId) {
                        app.$store.commit('auth/setBusinessMeta', {
                            bk_biz_id: parseInt(bizId)
                        })
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
}
