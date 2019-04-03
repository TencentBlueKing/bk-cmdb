import { getMetadataBiz } from '@/utils/tools'
import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'
import {
    C_MODEL_GROUP,
    U_MODEL_GROUP,
    D_MODEL_GROUP,
    C_MODEL,
    U_MODEL,
    D_MODEL
} from '@/dictionary/auth'

export const OPERATION = {
    C_MODEL_GROUP,
    U_MODEL_GROUP,
    D_MODEL_GROUP,
    C_MODEL,
    U_MODEL,
    D_MODEL
}

const modelPath = '/model'

export default [{
    name: 'model',
    path: modelPath,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'model',
            i18n: 'Nav["模型"]',
            path: modelPath,
            order: 1,
            parent: NAV_MODEL_MANAGEMENT,
            adminView: true
        },
        auth: {
            view: '',
            operation: Object.values(OPERATION)
        }
    }
}, {
    name: 'modelDetails',
    path: '/model/details/:modelId',
    component: () => import('./children/index.vue'),
    meta: {
        auth: {
            view: '',
            operation: [
                OPERATION.U_MODEL,
                OPERATION.D_MODEL
            ],
            setDynamicMeta: (to, from, app) => {
                const modelId = to.params.modelId
                const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
                const bizId = getMetadataBiz(model)
                const dynamicMeta = {}
                if (bizId) {
                    dynamicMeta.bk_biz_id = bizId
                }
                app.$store.commit('auth/setDynamicMeta', dynamicMeta)
            }
        }
    }
}]
