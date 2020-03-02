import Meta from '@/router/meta'
import { MENU_RESOURCE_CLOUD_RESOURCE } from '@/dictionary/menu-symbol'
export default {
    name: MENU_RESOURCE_CLOUD_RESOURCE,
    path: 'cloud-resource',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '云资源发现'
        },
        auth: {
            view: null
        },
        available: false
    })
}
