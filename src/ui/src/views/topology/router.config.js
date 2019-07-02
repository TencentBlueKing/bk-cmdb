import Meta from '@/router/meta'
import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
import {
    U_HOST,
    HOST_TO_RESOURCE
} from '@/dictionary/auth'

export const OPERATION = {
    U_HOST,
    HOST_TO_RESOURCE
}

const path = '/topology'

export default {
    name: 'topology',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'topology',
            i18n: 'Nav["业务拓扑"]',
            path: path,
            order: 2,
            parent: NAV_BUSINESS_RESOURCE,
            adminView: false
        },
        auth: {
            operation: Object.values(OPERATION)
        },
        requireBusiness: true,
        i18nTitle: 'Nav["业务拓扑"]'
    })
}
