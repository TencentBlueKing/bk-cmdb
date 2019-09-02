import Meta from '@/router/meta'
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

export default {
    name: 'customFields',
    path: '/custom-fields',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'customQuery',
            i18n: '自定义字段',
            path: '/custom-fields',
            order: 4,
            businessView: true
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        i18nTitle: '自定义字段'
    })
}
