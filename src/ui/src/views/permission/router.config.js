import { NAV_PERMISSION } from '@/dictionary/menu'
import {
    SYSTEM_MANAGEMENT
} from '@/dictionary/auth'

export const OPERATION = {
    SYSTEM_MANAGEMENT
}

const path = {
    business: '/permission/business',
    system: '/permission/system'
}

export default [{
    name: 'businessPermission',
    path: path.business,
    component: () => import('./business.vue'),
    meta: {
        menu: {
            id: 'businessPermission',
            i18n: 'Nav["业务权限管理"]',
            path: path.business,
            order: 1,
            parent: NAV_PERMISSION
        },
        auth: {
            view: SYSTEM_MANAGEMENT,
            operation: []
        }
    }
}, {
    name: 'systemPermission',
    path: path.system,
    component: () => import('./role.vue'),
    meta: {
        menu: {
            id: 'systemPermission',
            i18n: 'Nav["系统权限管理"]',
            path: path.system,
            order: 2,
            parent: NAV_PERMISSION
        },
        auth: {
            view: SYSTEM_MANAGEMENT,
            operation: []
        }
    }
}]
