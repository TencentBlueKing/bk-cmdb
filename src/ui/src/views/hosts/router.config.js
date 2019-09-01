import Meta from '@/router/meta'
import { MENU_BUSINESS_HOST, MENU_BUSINESS } from '@/dictionary/menu-symbol'
import {
    R_HOST,
    U_HOST,
    HOST_TO_RESOURCE
} from '@/dictionary/auth'

export const OPERATION = {
    U_HOST,
    R_HOST,
    HOST_TO_RESOURCE
}

export default {
    name: 'hosts',
    path: 'host',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '业务主机',
            parent: MENU_BUSINESS_HOST
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        filterPropertyKey: 'business_host_filter_properties'
    })
}
