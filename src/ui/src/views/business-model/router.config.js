import { NAV_MODEL_MANAGEMENT } from '@/dictionary/menu'
import { SYSTEM_TOPOLOGY } from '@/dictionary/auth'
import Meta from '@/router/meta'

export const OPERATION = {
    SYSTEM_TOPOLOGY
}

const path = '/model/business'

export default {
    name: 'businessModel',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'businessModel',
            i18n: '业务层级',
            path: path,
            order: 3,
            parent: NAV_MODEL_MANAGEMENT,
            businessView: false
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'global'
            }
        },
        i18nTitle: '业务层级'
    })
}
