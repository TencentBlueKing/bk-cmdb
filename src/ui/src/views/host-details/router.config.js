import {
    U_HOST,
    U_RESOURCE_HOST
} from '@/dictionary/auth'

const component = () => import(/* webpackChunkName: "hostDetails" */ './index.vue')

export const OPERATION = {
    U_HOST,
    U_RESOURCE_HOST
}

export const RESOURCE_HOST = 'resourceHostDetails'

export const BUSINESS_HOST = 'businessHostDetails'

export default [{
    name: RESOURCE_HOST,
    path: '/host/:id',
    component: component,
    meta: {
        auth: {
            view: null,
            operation: [U_RESOURCE_HOST]
        }
    }
}, {
    name: BUSINESS_HOST,
    path: '/business/:business/host/:id',
    component: component,
    meta: {
        auth: {
            view: null,
            operation: [U_HOST]
        }
    }
}]
