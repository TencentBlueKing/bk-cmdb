import Meta from '@/router/meta'
import { NAV_SERVICE_MANAGEMENT } from '@/dictionary/menu'

const path = '/service/template'

export default [{
    name: 'serviceTemplate',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'serviceTemplate',
            i18n: 'Nav["服务模板"]',
            path: path,
            order: 2,
            parent: NAV_SERVICE_MANAGEMENT
        }
    })
}, {
    name: 'operationalTemplate',
    path: '/service/operational',
    component: () => import('./children/operational.vue')
}]
