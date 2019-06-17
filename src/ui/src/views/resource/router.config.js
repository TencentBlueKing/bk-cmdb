import Meta from '@/router/meta'
import { NAV_BASIC_RESOURCE } from '@/dictionary/menu'
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

const path = '/resource'

export default {
    name: 'resource',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'resource',
            i18n: 'Nav["主机"]',
            path: path,
            order: 2,
            parent: NAV_BASIC_RESOURCE
        },
        auth: {
            operation: Object.values(OPERATION)
        }
    })
}
