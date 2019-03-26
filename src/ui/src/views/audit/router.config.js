import { NAV_AUDIT_ANALYSE } from '@/dictionary/nav'

export default {
    name: 'audit',
    path: '/auditing',
    component: () => import('./index.vue'),
    meta: {
        menu: {
            id: 'audit',
            i18n: 'Nav["操作审计"]',
            parent: NAV_AUDIT_ANALYSE
        },
        auth: {
            view: '',
            operation: []
        }
    }
}
