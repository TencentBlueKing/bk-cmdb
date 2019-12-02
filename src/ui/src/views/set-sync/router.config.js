import Meta from '@/router/meta'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
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
        }
    })
}]
