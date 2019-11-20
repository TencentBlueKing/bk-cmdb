import Meta from '@/router/meta'
import { MENU_BUSINESS, MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
import {
    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE
} from '@/dictionary/auth'

export default [{
    name: 'createServiceInstance',
    path: 'service/instance/create/set/:setId/module/:moduleId',
    component: () => import('./create.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '添加服务实例',
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        },
        auth: {
            operation: {
                C_SERVICE_INSTANCE,
                U_SERVICE_INSTANCE,
                D_SERVICE_INSTANCE,
                R_SERVICE_INSTANCE
            }
        }
    })
}, {
    name: 'cloneServiceInstance',
    path: 'service/instance/clone/set/:setId/module/:moduleId/instance/:instanceId/host/:hostId',
    props: true,
    component: () => import('./clone.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '克隆服务实例',
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        },
        auth: {
            operation: {
                C_SERVICE_INSTANCE,
                U_SERVICE_INSTANCE,
                D_SERVICE_INSTANCE,
                R_SERVICE_INSTANCE
            }
        }
    })
}]
