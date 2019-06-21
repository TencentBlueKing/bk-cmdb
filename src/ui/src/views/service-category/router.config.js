import Meta from '@/router/meta'
import { NAV_SERVICE_MANAGEMENT } from '@/dictionary/menu'

const path = '/service/cagetory'

export default {
    name: 'serviceCagetory',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'serviceCagetory',
            i18n: 'Nav["服务分类"]',
            path: path,
            order: 1,
            parent: NAV_SERVICE_MANAGEMENT,
            adminView: false
        },
        i18nTitle: 'Nav["服务分类"]'
    })
}
