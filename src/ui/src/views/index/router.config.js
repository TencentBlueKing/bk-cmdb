import Meta from '@/router/meta'
import { NAV_INDEX } from '@/dictionary/menu'

const path = '/index'

export default [{
    name: 'index',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: NAV_INDEX,
            i18n: '首页',
            path: path
        },
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
