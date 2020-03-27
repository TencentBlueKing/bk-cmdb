import Meta from '@/router/meta'
import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'

export default {
    name: 'cloudDiscover',
    path: 'cloud-discover',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '云资源发现'
        }
    })
}
