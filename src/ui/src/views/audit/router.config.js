import { NAV_AUDIT_ANALYSE } from '@/dictionary/menu'
import { R_AUDIT } from '@/dictionary/auth'
import Meta from '@/router/meta'
const path = '/auditing'

export default {
    name: 'audit',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'audit',
            i18n: '操作审计',
            path: path,
            parent: NAV_AUDIT_ANALYSE
        },
        auth: {
            view: R_AUDIT,
            setAuthScope () {
                this.authScope = 'global'
            }
        },
        i18nTitle: '操作审计'
    })
}
