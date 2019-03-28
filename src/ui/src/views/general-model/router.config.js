import {
    C_INST,
    R_INST,
    U_INST,
    D_INST,
    GET_MODEL_INST_AUTH_META
} from '@/dictionary/auth'

const prefix = '/general-model/'
const param = 'objId'

export const GET_MODEL_PATH = modelId => {
    return prefix + modelId
}

export const OPERATION = {
    C_INST,
    U_INST,
    D_INST
}

export default [{
    name: 'generalModel',
    path: `${prefix}:${param}`,
    component: () => import('./index.vue'),
    meta: {
        auth: {
            view: R_INST,
            meta: GET_MODEL_INST_AUTH_META,
            operation: Object.values(OPERATION)
        },
        dynamicParams: [param]
    }
}, {
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
}]
