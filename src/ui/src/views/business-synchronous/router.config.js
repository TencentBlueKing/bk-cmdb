import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_SERVICE_TEMPLATE
} from '@/dictionary/menu-symbol'

export default [{
    name: 'syncServiceFromModule',
    path: 'synchronous/module/:template/:modules',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步模板',
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        },
        layout: {}
    })
}, {
    name: 'syncServiceFromTemplate',
    path: 'sync/service-template/:template/:modules',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步模板',
            relative: MENU_BUSINESS_SERVICE_TEMPLATE
        },
        layout: {}
    })
}]
