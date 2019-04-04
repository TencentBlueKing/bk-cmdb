import { getMetadataBiz } from '@/utils/tools'
import {
    C_INST,
    R_INST,
    U_INST,
    D_INST
} from '@/dictionary/auth'

const prefix = '/general-model/'
const param = 'objId'

export const GET_MODEL_PATH = modelId => {
    return prefix + modelId
}

export const OPERATION = {
    C_INST,
    R_INST,
    U_INST,
    D_INST
}

export default [{
    path: GET_MODEL_PATH('host'),
    redirect: {
        name: 'resource'
    }
}, {
    path: GET_MODEL_PATH('process'),
    redirect: {
        name: 'process'
    }
}, {
    path: GET_MODEL_PATH('biz'),
    redirect: {
        name: 'business'
    }
}, {
    path: GET_MODEL_PATH('plat'),
    redirect: {
        name: '404'
    }
}, {
    name: 'generalModel',
    path: GET_MODEL_PATH(`:${param}`),
    component: () => import('./index.vue'),
    meta: {
        auth: {
            view: '',
            operation: Object.values(OPERATION),
            setDynamicMeta: (to, from, app) => {
                console.log(app.$store.getters['objectModelClassify/models'])
                const modelId = to.params[param]
                const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
                const bizId = getMetadataBiz(model)
                const dynamicMeta = {}
                if (bizId) {
                    dynamicMeta.bk_biz_id = bizId
                }
                app.$store.commit('auth/setDynamicMeta', {})
            }
        }
    }
}]
