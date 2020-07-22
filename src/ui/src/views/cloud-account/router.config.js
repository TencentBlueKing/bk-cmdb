import Meta from '@/router/meta'
import { MENU_RESOURCE_CLOUD_ACCOUNT } from '@/dictionary/menu-symbol'
import { OPERATION } from '@/dictionary/iam-auth'
export default {
    name: MENU_RESOURCE_CLOUD_ACCOUNT,
    path: 'cloud-account',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '云账户'
        },
        auth: {
            view: { type: OPERATION.R_CLOUD_ACCOUNT }
        }
    })
}
