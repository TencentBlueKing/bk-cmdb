import Meta from '@/router/meta'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'

export default [{
    name: 'setSync',
    path: 'set/sync/setTemplateId/:setTemplateId',
    component: () => import('./sync-index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '同步集群模板'
        }
    })
}, {
    name: 'viewSync',
    path: 'set/sync/view',
    component: () => import('./view-sync.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '查看同步'
        }
    })
}]
