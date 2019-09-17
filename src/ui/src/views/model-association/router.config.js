import Meta from '@/router/meta'
import { MENU_MODEL_ASSOCIATION } from '@/dictionary/menu-symbol'
import {
    C_RELATION,
    U_RELATION,
    D_RELATION
} from '@/dictionary/auth'

export const OPERATION = {
    C_RELATION,
    U_RELATION,
    D_RELATION
}

export default {
    name: MENU_MODEL_ASSOCIATION,
    path: 'association',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '关联类型'
        },
        auth: {
            operation: OPERATION,
            authScope: 'global'
        }
    })
}
