import { MENU_MODEL_BUSINESS_TOPOLOGY } from '@/dictionary/menu-symbol'
import { getMetadataBiz } from '@/utils/tools'
import {
    SYSTEM_TOPOLOGY,
    U_MODEL,
    D_MODEL,
    GET_AUTH_META
} from '@/dictionary/auth'
import Meta from '@/router/meta'

export const OPERATION = {
    SYSTEM_TOPOLOGY
}

export default [{
    name: MENU_MODEL_BUSINESS_TOPOLOGY,
    path: 'business/topology',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '业务层级'
        },
        auth: {
            operation: OPERATION,
            authScope: 'global'
        }
    })
}, {
    name: 'businessModelDetails',
    path: 'business/topology/details/:modelId',
    component: () => import('@/views/model-manage/children/index.vue'),
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
                    ...GET_AUTH_META(U_MODEL),
                    resource_id: model.id
                }, {
                    ...GET_AUTH_META(D_MODEL),
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
