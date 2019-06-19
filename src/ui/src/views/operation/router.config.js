import Meta from '@/router/meta'
import { NAV_AUDIT_ANALYSE } from '@/dictionary/menu'
import {
    R_STATISTICAL_REPORT
} from '@/dictionary/auth'

export const OPERATION = {
    R_STATISTICAL_REPORT
}

const path = '/statistics_server'

export default {
    name: 'statisticalReport',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'statisticalReport',
            i18n: 'Nav["统计报表"]',
            path: path,
            parent: NAV_AUDIT_ANALYSE,
            businessView: false
        },
        auth: {
            operation: Object.values(OPERATION)
        }
    })
}
