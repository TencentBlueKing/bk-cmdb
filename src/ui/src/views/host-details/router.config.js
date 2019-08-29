import {
    U_HOST,
    U_RESOURCE_HOST,
    GET_AUTH_META,
    D_SERVICE_INSTANCE
} from '@/dictionary/auth'

const component = () => import(/* webpackChunkName: "hostDetails" */ './index.vue')

export const OPERATION = {
    U_HOST,
    U_RESOURCE_HOST,
    D_SERVICE_INSTANCE
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

            operation: [U_RESOURCE_HOST],
            setDynamicMeta (to, from, app) {
                const meta = GET_AUTH_META(U_RESOURCE_HOST)
                app.$store.commit('auth/setResourceMeta', {
                    ...meta,
                    resource_id: parseInt(to.params.id)
                })
            },
            setAuthScope () {
                this.authScope = 'global'
            }
        },
        i18nTitle: '主机详情'
    }
}, {
    name: BUSINESS_HOST,
    path: '/business/:business/host/:id',
    component: component,
    meta: {
        auth: {
            view: null,
            operation: [U_HOST, D_SERVICE_INSTANCE],
            setDynamicMeta (to, from, app) {
                const hostMeta = GET_AUTH_META(U_HOST)
                const serviceInstanceMeta = GET_AUTH_META(D_SERVICE_INSTANCE)
                app.$store.commit('auth/setResourceMeta', [{
                    ...hostMeta,
                    resource_id: parseInt(to.params.id),
                    bk_biz_id: parseInt(to.params.business)
                }, {
                    ...serviceInstanceMeta,
                    resource_id: parseInt(to.params.id),
                    bk_biz_id: parseInt(to.params.business)
                }])
            },
            setAuthScope () {
                this.authScope = 'business'
            }
        },
        i18nTitle: '主机详情',
        requireBusiness: true
    }
}]
