import Meta from '@/router/meta'
import {
    R_CONFIG_ADMIN,
    GET_AUTH_META
} from '@/dictionary/auth'

export default [{
    name: 'admin_index',
    path: 'index',
    redirect: 'config'
}, {
    name: 'admin_config',
    path: 'config',
    component: () => import('@/views/admin/config'),
    meta: new Meta({
        menu: {
            i18n: '配置'
        },
        auth: {
            view: GET_AUTH_META(R_CONFIG_ADMIN)
        },
        layout: {
            breadcrumbs: false
        }
    })
}]
