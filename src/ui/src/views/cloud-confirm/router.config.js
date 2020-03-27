import Meta from '@/router/meta'
import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'

export default [{
    name: 'resourceConfirm',
    path: 'resource-confirm',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '资源确认'
        }
    })
}, {
    name: 'confirmHistory',
    path: 'confirm-history',
    component: () => import('./history.vue'),
    meta: new Meta()
}]
