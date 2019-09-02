import Meta from '@/router/meta'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
import {
    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE
} from '@/dictionary/auth'

export const OPERATION = {
    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE
}

export default [{
    name: 'createServiceInstance',
    path: '/business/:business/service/instance/create/set/:setId/module/:moduleId',
    component: () => import('./create.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        resetMenu: false,
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        }
    })
}, {
    name: 'cloneServiceInstance',
    path: '/business/:business/service/instance/clone/set/:setId/module/:moduleId/instance/:instanceId/host/:hostId',
    props: true,
    component: () => import('./clone.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        resetMenu: false,
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'business'
            }
        }
    })
}]
