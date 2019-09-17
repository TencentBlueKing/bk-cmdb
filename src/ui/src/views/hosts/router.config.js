import Meta from '@/router/meta'
import { MENU_BUSINESS_HOST, MENU_BUSINESS, MENU_BUSINESS_HOST_MANAGEMENT } from '@/dictionary/menu-symbol'
import {
    U_HOST,
    HOST_TO_RESOURCE
} from '@/dictionary/auth'

export const OPERATION = {
    U_HOST,
    HOST_TO_RESOURCE
}

export default {
    name: MENU_BUSINESS_HOST_MANAGEMENT,
    path: 'host',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '业务主机',
            parent: MENU_BUSINESS_HOST
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        },
        showBreadcrumbs: true,
        filterPropertyKey: 'business_host_filter_properties'
    })
}
