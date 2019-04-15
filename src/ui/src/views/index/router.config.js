import Meta from '@/router/meta'
import { NAV_INDEX } from '@/dictionary/menu'

const path = '/index'

export default {
    name: 'index',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: NAV_INDEX,
            i18n: 'Nav["首页"]',
            path: path
        },
        auth: {
            view: null
        }
    })
}
