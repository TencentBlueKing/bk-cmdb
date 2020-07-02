import Meta from '@/router/meta'

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
            view: null
        },
        layout: {
            breadcrumbs: false
        }
    })
}]
