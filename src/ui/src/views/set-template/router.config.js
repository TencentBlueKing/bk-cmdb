import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_SERVICE,
    MENU_BUSINESS_SET_TEMPLATE,
    MENU_BUSINESS_HOST_AND_SERVICE
} from '@/dictionary/menu-symbol'
import {
    C_SET_TEMPLATE,
    U_SET_TEMPLATE,
    D_SET_TEMPLATE,
    U_TOPO
} from '@/dictionary/auth'

export default [{
    name: MENU_BUSINESS_SET_TEMPLATE,
    path: 'set/template',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '集群模板',
            parent: MENU_BUSINESS_SERVICE
        },
        auth: {
            operation: {
                C_SET_TEMPLATE,
                U_SET_TEMPLATE,
                D_SET_TEMPLATE,
                U_TOPO
            }
        }
    })
}, {
    name: 'setTemplateConfig',
    path: 'set/template/:mode/:templateId?',
    component: () => import('./template.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '集群模板',
            relative: MENU_BUSINESS_SET_TEMPLATE
        },
        auth: {
            operation: {
                C_SET_TEMPLATE,
                U_SET_TEMPLATE,
                D_SET_TEMPLATE,
                U_TOPO
            }
        },
        layout: {
            previous: (view) => {
                const query = view.$route.query
                let params = { name: MENU_BUSINESS_SET_TEMPLATE }
                if (query.tab) {
                    params = {
                        name: MENU_BUSINESS_HOST_AND_SERVICE,
                        query: query
                    }
                }
                return params
            }
        }
    })
}, {
    name: 'syncHistory',
    path: 'set/instance/history/:templateId?',
    component: () => import('./sync-history.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步历史',
            relative: MENU_BUSINESS_SET_TEMPLATE
        },
        layout: {
            previous () {
                return {
                    name: 'setTemplateConfig',
                    params: {
                        mode: 'view',
                        templateId: window.CMDB_APP.$route.params.templateId
                    }
                }
            }
        }
    })
}]
