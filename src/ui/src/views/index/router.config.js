import Meta from '@/router/meta'
import { MENU_INDEX } from '@/dictionary/menu-symbol'

export default [{
    name: MENU_INDEX,
    path: '/index',
    component: () => import('./index.vue'),
    meta: new Meta({
        auth: {
            view: null
        },
        i18nTitle: '首页'
    })
}, {
    name: 'search',
    path: '/index/search',
    component: () => import('./children/full-text-search.vue'),
    meta: new Meta({
        checkAvailable: (to, from, app) => {
            return window.Site.fullTextSearch === 'on'
        }
    })
}]
