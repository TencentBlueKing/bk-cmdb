import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_HOST,
    MENU_BUSINESS_HOST_APPLY
} from '@/dictionary/menu-symbol'
import {
    U_HOST,
    D_SERVICE_INSTANCE,
    HOST_TO_RESOURCE
} from '@/dictionary/auth'

export const OPERATION = {
    U_HOST,
    D_SERVICE_INSTANCE,
    HOST_TO_RESOURCE
}

export default [{
    name: MENU_BUSINESS_HOST_APPLY,
    path: 'host-apply',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机自动应用',
            parent: MENU_BUSINESS_HOST
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        }
    })
}, {
    name: 'hostApplyConfirm',
    path: 'host-apply/confirm',
    component: () => import('./property-confirm'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        },
        auth: {
            operation: OPERATION,
            authScope: 'business'
        }
    })
}]
