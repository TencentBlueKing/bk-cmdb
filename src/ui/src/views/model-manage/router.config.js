import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'
import {
    G_C_MODEL_GROUP,
    G_U_MODEL_GROUP,
    G_D_MODEL_GROUP,
    G_C_MODEL,
    G_U_MODEL,
    G_D_MODEL
} from '@/dictionary/auth'

export const OPERATION = {
    G_C_MODEL_GROUP,
    G_U_MODEL_GROUP,
    G_D_MODEL_GROUP,
    G_C_MODEL,
    G_U_MODEL,
    G_D_MODEL
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
            parent: NAV_MODEL_MANAGEMENT
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
                OPERATION.G_U_MODEL,
                OPERATION.G_D_MODEL
            ]
        }
    }
}]
