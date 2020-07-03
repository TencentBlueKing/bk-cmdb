import Meta from '@/router/meta'
import { MENU_INDEX } from '@/dictionary/menu-symbol'

export default [{
    name: MENU_INDEX,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '首页'
        },
        auth: {
            view: null
        },
        layout: {
            breadcrumbs: false
        }
    })
}, {
    name: 'search',
    path: 'index/search',
    component: () => import('./children/full-text-search.vue'),
    meta: new Meta({
        checkAvailable: (to, from, app) => {
            return window.CMDB_CONFIG.site.fullTextSearch === 'on'
        },
        layout: {
            breadcrumbs: false
        }
    })
}]
