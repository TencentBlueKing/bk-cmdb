import Meta from '@/router/meta'
import { MENU_RESOURCE_CLOUD_AREA } from '@/dictionary/menu-symbol'
export default {
    name: MENU_RESOURCE_CLOUD_AREA,
    path: 'cloud-area',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '云区域'
        },
        auth: {
            view: null
        }
    })
}
