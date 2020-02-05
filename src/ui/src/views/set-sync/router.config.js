import Meta from '@/router/meta'
import { MENU_BUSINESS, MENU_BUSINESS_SET_TEMPLATE, MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
import { U_TOPO } from '@/dictionary/auth'
export default [{
    name: 'setSync',
    path: 'set/sync/:setTemplateId',
    component: () => import('./sync-index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步集群模板'
        },
        auth: {
            operation: {
                U_TOPO
            }
        },
        layout: {
            previous: (view) => {
                const moduleId = view.$route.params['moduleId']
                let params = {
                    name: MENU_BUSINESS_SET_TEMPLATE
                }
                if (moduleId) {
                    params = {
                        name: MENU_BUSINESS_HOST_AND_SERVICE,
                        query: {
                            node: 'set-' + moduleId
                        }
                    }
                } else {
                    params = {
                        name: 'setTemplateConfig',
                        params: {
                            templateId: view.setTemplateId,
                            mode: 'view'
                        },
                        query: {
                            tab: 'instance'
                        }
                    }
                }
                return params
            }
        }
    })
}]
