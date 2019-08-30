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

const path = '/model/association'

export default {
    name: 'modelAssociation',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'modelAssociation',
            i18n: '关联类型',
            path: path,
            order: 4,
            parent: NAV_MODEL_MANAGEMENT,
            businessView: false
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
