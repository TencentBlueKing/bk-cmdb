import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_SERVICE_TEMPLATE
} from '@/dictionary/menu-symbol'

export default [{
    name: 'syncServiceFromModule',
    path: 'synchronous/module/:moduleId/:setId',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步模板',
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        }
    })
}, {
    name: 'syncServiceFromTemplate',
    path: 'synchronous/service-template/:moduleId/:setId',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步模板',
            relative: MENU_BUSINESS_SERVICE_TEMPLATE
        }
    })
}]
