import Meta from '@/router/meta'
import { MENU_MODEL_ASSOCIATION } from '@/dictionary/menu-symbol'

export default {
    name: MENU_MODEL_ASSOCIATION,
    path: 'association',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '关联类型'
        }
    })
}
