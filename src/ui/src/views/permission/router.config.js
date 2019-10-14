import Meta from '@/router/meta'
import { NAV_PERMISSION } from '@/dictionary/menu'

const path = {
    business: '/permission/business',
    system: '/permission/system'
}

export default [{
    name: 'businessPermission',
    path: path.business,
    component: () => import('./business.vue'),
    meta: new Meta({
        menu: {
            i18n: '业务权限管理'
        }
    })
}, {
    name: 'systemPermission',
    path: path.system,
    component: () => import('./role.vue'),
    meta: new Meta({
        menu: {
            i18n: '系统权限管理'
        }
    })
}]
