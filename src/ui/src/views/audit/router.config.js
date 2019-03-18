import { NAV_AUDIT_ANALYSE } from '@/types/nav'
import { G_R_AUDIT } from '@/types/auth'

export default {
    name: 'audit',
    path: '/audit',
    components: () => import('index.vue'),
    meta: {
        menu: {
            id: "audit",
            i18n: "Nav['操作审计']",
            parent: NAV_AUDIT_ANALYSE,
        },
        auth: {
            view: [G_R_AUDIT],
            operation: []
        }
    }
}