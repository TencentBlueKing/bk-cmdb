import Meta from '@/router/meta'
import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
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

const path = '/hosts'

export default {
    name: 'hosts',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'hosts',
            i18n: '业务主机',
            path: path,
            parent: NAV_BUSINESS_RESOURCE,
            adminView: false
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        requireBusiness: true,
        i18nTitle: '业务主机',
        filterPropertyKey: 'business_host_filter_properties'
    })
}
