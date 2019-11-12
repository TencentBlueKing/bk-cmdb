import Meta from '@/router/meta'
import { MENU_MODEL_TOPOLOGY_NEW } from '@/dictionary/menu-symbol'

export default [{
    name: MENU_MODEL_TOPOLOGY_NEW,
    path: 'all/topology/new',
    component: () => import('./index.new.vue'),
    meta: new Meta({
        menu: {
            i18n: '模型关系'
        },
        layout: {
            breadcrumbs: false
        }
    })
}]
