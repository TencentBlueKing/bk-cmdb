import { NAV_AUDIT_ANALYSE } from '@/dictionary/menu'
import Meta from '@/router/meta'

const path = '/auditing'

export default {
    name: 'audit',
    path: path,
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'audit',
            i18n: 'Nav["操作审计"]',
            path: path,
            order: 1,
            parent: NAV_AUDIT_ANALYSE,
            adminView: true
        },
        auth: {
            view: '',
            operation: []
        }
    }
}
