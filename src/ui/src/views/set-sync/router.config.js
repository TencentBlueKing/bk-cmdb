import Meta from '@/router/meta'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
import {
    C_SET_TEMPLATE,
    U_SET_TEMPLATE,
    D_SET_TEMPLATE,
    R_SET_TEMPLATE
} from '@/dictionary/auth'

export const OPERATION = {
    C_SET_TEMPLATE,
    U_SET_TEMPLATE,
    D_SET_TEMPLATE,
    R_SET_TEMPLATE
}

export default [{
    name: 'setSync',
    path: 'set/sync/:setTemplateId',
    component: () => import('./sync-index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步集群模板'
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        }
    })
}]
