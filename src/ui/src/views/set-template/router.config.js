import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_SERVICE,
    MENU_BUSINESS_SET_TEMPLATE
} from '@/dictionary/menu-symbol'
export default [{
    name: MENU_BUSINESS_SET_TEMPLATE,
    path: 'set/template',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '集群模板',
            parent: MENU_BUSINESS_SERVICE
        }
    })
}, {
    name: 'createTemplate',
    path: 'create/template',
    component: () => import('./children/create-template.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '创建集群模板'
        }
    })
}]
