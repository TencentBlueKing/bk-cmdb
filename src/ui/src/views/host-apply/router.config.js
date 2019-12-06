import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_HOST,
    MENU_BUSINESS_HOST_APPLY
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
    name: 'hostApplyConfirm',
    path: 'host-apply/confirm',
    component: () => import('./property-confirm'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        }
    })
}, {
    name: 'hostApplyEdit',
    path: 'host-apply/edit',
    component: () => import('./edit'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        }
    })
}, {
    name: 'hostApplyConflict',
    path: 'host-apply/conflict',
    component: () => import('./conflict-list'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        }
    })
}, {
    name: 'hostApplyFailed',
    path: 'host-apply/failed',
    component: () => import('./failed-list'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        }
    })
}]
