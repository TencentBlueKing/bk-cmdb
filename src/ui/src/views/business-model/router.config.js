import { MENU_MODEL_BUSINESS_TOPOLOGY } from '@/dictionary/menu-symbol'
import { SYSTEM_TOPOLOGY } from '@/dictionary/auth'
import Meta from '@/router/meta'

export const OPERATION = {
    SYSTEM_TOPOLOGY
}

export default {
    name: MENU_MODEL_BUSINESS_TOPOLOGY,
    path: 'business/topology',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '业务层级'
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'global'
            }
        }
    })
}
