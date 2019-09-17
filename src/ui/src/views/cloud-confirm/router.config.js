import Meta from '@/router/meta'
import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
import {
    C_CLOUD_CONFIRM,
    U_CLOUD_CONFIRM,
    D_CLOUD_CONFIRM,
    R_CLOUD_CONFIRM,
    R_CONFIRM_HISTORY
} from '@/dictionary/auth'

export const OPERATION = {
    R_CLOUD_CONFIRM,
    C_CLOUD_CONFIRM,
    U_CLOUD_CONFIRM,
    D_CLOUD_CONFIRM,
    R_CONFIRM_HISTORY
}

const path = 'resource-confirm'

export default [{
    name: 'resourceConfirm',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '资源确认'
        },
        auth: {
            operation: OPERATION
        },
        requireBusiness: true
    })
}, {
    name: 'confirmHistory',
    path: 'confirm-history',
    component: () => import('./history.vue'),
    meta: new Meta({
        auth: {
            operation: { R_CONFIRM_HISTORY },
            authScope: 'global'
        }
    })
}]
