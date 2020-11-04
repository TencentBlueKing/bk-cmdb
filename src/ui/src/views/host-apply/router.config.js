import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_HOST,
    MENU_BUSINESS_HOST_APPLY,
    MENU_BUSINESS_HOST_APPLY_EDIT,
    MENU_BUSINESS_HOST_APPLY_CONFIRM,
    MENU_BUSINESS_HOST_APPLY_CONFLICT,
    MENU_BUSINESS_HOST_APPLY_FAILED
} from '@/dictionary/menu-symbol'

export default [{
    name: MENU_BUSINESS_HOST_APPLY,
    path: 'host-apply',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST
        }
    })
}, {
    name: MENU_BUSINESS_HOST_APPLY_CONFIRM,
    path: 'host-apply/confirm',
    component: () => import('./property-confirm'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY,
            relative: MENU_BUSINESS_HOST_APPLY
        }
    })
}, {
    name: MENU_BUSINESS_HOST_APPLY_EDIT,
    path: 'host-apply/edit',
    component: () => import('./edit'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY,
            relative: MENU_BUSINESS_HOST_APPLY
        }
    })
}, {
    name: MENU_BUSINESS_HOST_APPLY_CONFLICT,
    path: 'host-apply/conflict',
    component: () => import('./conflict-list'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY,
            relative: MENU_BUSINESS_HOST_APPLY
        }
    })
}, {
    name: MENU_BUSINESS_HOST_APPLY_FAILED,
    path: 'host-apply/failed',
    component: () => import('./failed-list'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY,
            relative: MENU_BUSINESS_HOST_APPLY
        }
    })
}]
