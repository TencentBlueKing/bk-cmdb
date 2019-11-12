import Meta from '@/router/meta'
import {
    MENU_RESOURCE,
    MENU_RESOURCE_HOST,
    MENU_RESOURCE_MANAGEMENT,
    MENU_RESOURCE_HOST_DETAILS,
    MENU_RESOURCE_BUSINESS_HOST_DETAILS
} from '@/dictionary/menu-symbol'
const component = () => import('./index.vue')

export default [{
    name: MENU_RESOURCE_HOST_DETAILS,
    path: '/resource/host/:id',
    component: component,
    meta: new Meta({
        owner: MENU_RESOURCE,
        menu: {
            i18n: '主机详情',
            relative: [MENU_RESOURCE_HOST, MENU_RESOURCE_MANAGEMENT]
        }
    })
}, {
    name: MENU_RESOURCE_BUSINESS_HOST_DETAILS,
    path: '/resource/host/:business/:id',
    component: component,
    meta: new Meta({
        owner: MENU_RESOURCE,
        menu: {
            i18n: '主机详情',
            relative: [MENU_RESOURCE_HOST, MENU_RESOURCE_MANAGEMENT]
        }
    })
}]
