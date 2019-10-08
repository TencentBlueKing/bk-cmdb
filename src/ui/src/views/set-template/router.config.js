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
    name: 'setTemplateInfo',
    path: 'set/template/info/:templateId',
    component: () => import('./template-info.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '模板信息',
            relative: MENU_BUSINESS_SET_TEMPLATE
        }
    })
}, {
    name: 'setTemplateMode',
    path: 'set/template/:mode/:templateId?',
    component: () => import('./template.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            relative: MENU_BUSINESS_SET_TEMPLATE
        }
    })
}]
