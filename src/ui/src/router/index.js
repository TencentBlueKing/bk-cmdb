import Vue from 'vue'
import Router from 'vue-router'
import store from '@/store'
import preload from '@/setup/preload'
import $http from '@/api'

const index = () => import(/* webpackChunkName: index */ '@/views/index')
const modelManage = () => import(/* webpackChunkName: model */ '@/views/model-manage')
const modelDetail = () => import(/* webpackChunkName: model */ '@/views/model-manage/children')
const business = () => import(/* webpackChunkName: business */ '@/views/business')
const businessArchived = () => import(/* webpackChunkName: businessArchived */ '@/views/business/archived')
const generalModel = () => import(/* webpackChunkName: generalModel */ '@/views/general-model')
const deleteHistory = () => import(/* webpackChunkName: deleteHistory */ '@/views/history')
const hosts = () => import(/* webpackChunkName: hosts */ '@/views/hosts')
const eventpush = () => import(/* webpackChunkName: hosts */ '@/views/eventpush')
const resource = () => import(/* webpackChunkName: resource */ '@/views/resource')
const audit = () => import(/* webpackChunkName: hosts */ '@/views/audit')
const topology = () => import(/* webpackChunkName: topology */ '@/views/topology')
const process = () => import(/* webpackChunkName: process */ '@/views/process')
const customQuery = () => import(/* webpackChunkName: process */ '@/views/custom-query')
const error = () => import(/* webpackChunkName: error */ '@/views/status/error')
const systemAuthority = () => import(/* webpackChunkName: systemAuthority */ '@/views/permission/role')
const businessAuthority = () => import(/* webpackChunkName: businessAuthority */ '@/views/permission/business')
const modelTopology = () => import(/* webpackChunkName: modelTopology */ '@/views/model-topology')
const businessModel = () => import(/* webpackChunkName: businessModel */ '@/views/business-model')
const modelAssociation = () => import(/* webpackChunkName: modelAssociation */ '@/views/model-association')
Vue.use(Router)

const router = new Router({
    linkActiveClass: 'active',
    routes: [{
        path: '/',
        redirect: '/index'
    }, {
        path: '/index',
        component: index,
        meta: {
            ignoreAuthorize: true,
            isModel: false
        }
    }, {
        path: '/business',
        component: business,
        meta: {
            isModel: true,
            objId: 'biz'
        }
    }, {
        path: '/model',
        component: modelManage,
        meta: {
            isModel: false
        }
    }, {
        path: '/model/details/process',
        redirect: '/status-404'
    }, {
        path: '/model/details/plat',
        redirect: '/status-404'
    }, {
        path: '/model/details/:modelId',
        component: modelDetail,
        meta: {
            returnPath: '/model',
            relative: '/model',
            ignoreAuthorize: true,
            isModel: false
        }
    }, {
        path: '/model/topology',
        component: modelTopology,
        meta: {
            isModel: false
        }
    }, {
        path: '/model/business',
        component: businessModel,
        meta: {
            isModel: false
        }
    }, {
        path: '/model/association',
        component: modelAssociation,
        meta: {
            isModel: false
        }
    }, {
        path: '/eventpush',
        component: eventpush,
        meta: {
            isModel: false,
            authority: {
                type: 'backConfig',
                id: 'event'
            }
        }
    }, {
        path: '/authority/business',
        component: businessAuthority,
        meta: {
            isModel: false,
            isAdminOnly: true
        }
    }, {
        path: '/authority/system',
        component: systemAuthority,
        meta: {
            isModel: false,
            isAdminOnly: true
        }
    }, {
        path: '/history/biz',
        component: businessArchived,
        meta: {
            relative: '/business'
        }
    }, {
        path: '/general-model/:objId',
        component: generalModel,
        meta: {
            isModel: true
        }
    }, {
        path: '/history/:objId',
        component: deleteHistory,
        meta: {
            isModel: false
        }
    }, {
        path: '/hosts',
        component: hosts,
        meta: {
            requireBusiness: true,
            isModel: false
        }
    }, {
        path: '/resource',
        component: resource,
        meta: {
            isModel: false,
            authority: {
                type: 'globalBusi',
                id: 'resource'
            }
        }
    }, {
        path: '/auditing',
        component: audit,
        meta: {
            isModel: false,
            authority: {
                type: 'backConfig',
                id: 'audit'
            }
        }
    }, {
        path: '/topology',
        component: topology,
        meta: {
            requireBusiness: true,
            isModel: false
        }
    }, {
        path: '/process',
        component: process,
        meta: {
            requireBusiness: true,
            isModel: false
        }
    }, {
        path: '/custom-query',
        component: customQuery,
        meta: {
            requireBusiness: true,
            isModel: false
        }
    }, {
        path: '/status-require-business',
        components: require('@/views/status/require-business'),
        meta: {
            ignoreAuthorize: true
        }
    }, {
        path: '/status-403',
        components: require('@/views/status/403'),
        meta: {
            ignoreAuthorize: true
        }
    }, {
        path: '/status-404',
        components: require('@/views/status/404'),
        meta: {
            ignoreAuthorize: true
        }
    }, {
        path: '/status-error',
        component: error,
        meta: {
            ignoreAuthorize: true
        }
    }, {
        path: '*',
        redirect: '/status-404'
    }]
})

const cancelRequest = () => {
    const allRequest = $http.queue.get()
    const requestQueue = allRequest.filter(request => request.cancelWhenRouteChange)
    return $http.cancel(requestQueue.map(request => request.requestId))
}

const hasPrivilegeBusiness = () => {
    const privilegeBusiness = router.app.$store.getters['objectBiz/privilegeBusiness']
    return !!privilegeBusiness.length
}

const hasAuthority = to => {
    const privilege = router.app.$store.getters['objectBiz/privilegeBusiness']
    const {type, id} = to.meta.authority
    let authority = []
    if (type === 'globalBusi') {
        authority = router.app.$store.getters['userPrivilege/globalBusiAuthority'](id)
    } else if (type === 'backConfig') {
        authority = router.app.$store.getters['userPrivilege/backConfigAuthority'](id)
    }
    return authority.includes('search')
}

router.beforeEach(async (to, from, next) => {
    try {
        if (to.path !== '/status-error') {
            router.app.$store.commit('setGlobalLoading', true)
            await cancelRequest()
            await preload(router.app)
            if (to.meta.ignoreAuthorize) {
                next()
            } else if (to.meta.hasOwnProperty('authority')) {
                if (hasAuthority(to)) {
                    next()
                } else {
                    next({
                        path: '/status-403'
                    })
                }
            } else if (to.meta.requireBusiness && !hasPrivilegeBusiness()) {
                next({
                    path: '/status-require-business',
                    query: {
                        relative: to.path
                    }
                })
            } else {
                next()
            }
        } else {
            next()
        }
    } catch (e) {
        next({
            path: '/status-error',
            query: {
                relative: to.path
            }
        })
    }
})

router.afterEach((to, from) => {
    if (to.path === '/status-error') {
        $http.cancel()
    }
    router.app.$store.commit('setGlobalLoading', false)
})

export default router
