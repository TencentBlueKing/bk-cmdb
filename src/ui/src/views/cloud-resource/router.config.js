import Meta from '@/router/meta'
import { MENU_RESOURCE_CLOUD_RESOURCE } from '@/dictionary/menu-symbol'
import { OPERATION } from '@/dictionary/iam-auth'
export default {
    name: MENU_RESOURCE_CLOUD_RESOURCE,
    path: 'cloud-resource',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '云资源发现'
        },
        auth: {
            view: { type: OPERATION.R_CLOUD_RESOURCE_TASK }
        }
    })
}
