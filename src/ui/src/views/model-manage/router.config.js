import Meta from '@/router/meta'
import { getMetadataBiz } from '@/utils/tools'
import { MENU_MODEL_MANAGEMENT } from '@/dictionary/menu-symbol'
import {
    C_MODEL_GROUP,
    U_MODEL_GROUP,
    D_MODEL_GROUP,
    C_MODEL,
    U_MODEL,
    D_MODEL,
    GET_AUTH_META
} from '@/dictionary/auth'

export const OPERATION = {
    C_MODEL_GROUP,
    U_MODEL_GROUP,
    D_MODEL_GROUP,
    C_MODEL,
    U_MODEL,
    D_MODEL
}

export default [{
    name: MENU_MODEL_MANAGEMENT,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '模型'
        },
        auth: {
            operation: OPERATION,
            authScope: 'global'
        }
    })
}, {
    name: 'modelDetails',
    path: 'index/details/:modelId',
    component: () => import('./children/index.vue'),
    meta: new Meta({
        menu: {
            i18n: '模型详情'
        },
        auth: {
            operation: { U_MODEL, D_MODEL },
            authScope: 'global',
            setDynamicMeta: (to, from, app) => {
                const modelId = to.params.modelId
                const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
                const bizId = getMetadataBiz(model)
                const resourceMeta = [{
                    ...GET_AUTH_META(OPERATION.U_MODEL),
                    resource_id: model.id
                }, {
                    ...GET_AUTH_META(OPERATION.D_MODEL),
                    resource_id: model.id
                }]
                if (bizId) {
                    resourceMeta.forEach(meta => {
                        meta.bk_biz_id = parseInt(bizId)
                    })
                }
                app.$store.commit('auth/setResourceMeta', resourceMeta)
            }
        },
        checkAvailable: (to, from, app) => {
            const modelId = to.params.modelId
            const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
            return !!model
        }
    })
}]
