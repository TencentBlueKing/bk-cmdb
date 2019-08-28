import Meta from '@/router/meta'
import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'
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
    name: 'modelAssociation',
    path: 'association',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '关联类型',
            parent: NAV_MODEL_MANAGEMENT
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'global'
            }
        },
        i18nTitle: '关联类型'
    })
}
