import {
    D_C_MODEL,
    D_R_MODEL,
    D_U_MODEL,
    D_D_MODEL
} from '@/dictionary/auth'

const prefix = '/general-model/'
const param = 'objId'

export const GET_MODEL_PATH = modelId => {
    return prefix + modelId
}

export const OPERATION = {
    D_C_MODEL,
    D_U_MODEL,
    D_D_MODEL
}

export default {
    name: 'generalModel',
    path: `${prefix}:${param}`,
    component: () => import('./index.vue'),
    meta: {
        auth: {
            view: D_R_MODEL,
            operation: Object.values(OPERATION)
        },
        dynamicParams: [param]
    }
}
