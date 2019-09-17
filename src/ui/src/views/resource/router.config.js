import Meta from '@/router/meta'
import { MENU_RESOURCE_HOST } from '@/dictionary/menu-symbol'
import {
    C_RESOURCE_HOST,
    U_RESOURCE_HOST,
    D_RESOURCE_HOST
} from '@/dictionary/auth'

export const OPERATION = {
    C_RESOURCE_HOST,
    U_RESOURCE_HOST,
    D_RESOURCE_HOST
}

export default {
    name: MENU_RESOURCE_HOST,
    path: 'host',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '主机'
        },
        auth: {
            operation: OPERATION,
            authScope: 'global'
        },
        filterPropertyKey: 'resource_host_filter_properties'
    })
}
