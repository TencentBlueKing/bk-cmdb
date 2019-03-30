import { NAV_INDEX } from '@/dictionary/menu'

const path = '/index'

export default {
    name: 'index',
    path: path,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: NAV_INDEX,
            path: path
        },
        auth: {
            view: null
        }
    }
}
