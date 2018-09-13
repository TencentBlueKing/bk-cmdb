import Vue from 'vue'
import Router from 'vue-router'
import store from '@/store'
import preload from '@/setup/preload'
import $http from '@/api'

const index = () => import(/* webpackChunkName: index */ '@/views/index')
const model = () => import(/* webpackChunkName: model */ '@/views/model')
const modelTopo = () => import(/* webpackChunkName: model */ '@/views/model/model-topo')
const business = () => import(/* webpackChunkName: business */ '@/views/business')
const businessArchived = () => import(/* webpackChunkName: businessArchived */ '@/views/business/archived')
const generalModel = () => import(/* webpackChunkName: generalModel */ '@/views/general-model')
const deleteHistory = () => import(/* webpackChunkName: deleteHistory */ '@/views/history')
const hosts = () => import(/* webpackChunkName: hosts */ '@/views/hosts')
const eventpush = () => import(/* webpackChunkName: hosts */ '@/views/eventpush')
const permission = () => import(/* webpackChunkName: hosts */ '@/views/permission')
const resource = () => import(/* webpackChunkName: resource */ '@/views/resource')
const audit = () => import(/* webpackChunkName: hosts */ '@/views/audit')
const topology = () => import(/* webpackChunkName: topology */ '@/views/topology')
const process = () => import(/* webpackChunkName: process */ '@/views/process')
const customQuery = () => import(/* webpackChunkName: process */ '@/views/custom-query')
const error = () => import(/* webpackChunkName: error */ '@/views/status/error')

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
            ignoreAuthorize: true
        }
    }, {
        path: '/business',
        component: business
    }, {
        path: '/model',
        component: model,
        children: [{
            path: ':classifyId',
            component: modelTopo,
            meta: {
                relative: '/model'
            }
        }, {
            path: '',
            component: modelTopo,
            meta: {
                relative: '/model'
            }
        }]
    }, {
        path: '/eventpush',
        component: eventpush
    }, {
        path: '/permission',
        component: permission
    }, {
        path: '/history/biz',
        component: businessArchived,
        meta: {
            relative: '/business'
        }
    }, {
        path: '/general-model/:objId',
        component: generalModel
    }, {
        path: '/history/:objId',
        component: deleteHistory
    }, {
        path: '/hosts',
        component: hosts,
        meta: {
            requireBusiness: true
        }
    }, {
        path: '/resource',
        component: resource
    }, {
        path: '/auditing',
        component: audit
    }, {
        path: '/topology',
        component: topology,
        meta: {
            requireBusiness: true
        }
    }, {
        path: '/process',
        component: process,
        meta: {
            requireBusiness: true
        }
    }, {
        path: '/custom-query',
        component: customQuery,
        meta: {
            requireBusiness: true
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

const hasAuthority = (to) => {
    if (to.meta.ignoreAuthorize) {
        return true
    }
    const path = to.meta.relative || to.query.relative || to.path
    const authorizedNavigation = router.app.$store.getters['objectModelClassify/authorizedNavigation']
    return authorizedNavigation.some(navigation => {
        if (navigation.hasOwnProperty('path')) {
            return navigation.path === path
        }
        return navigation.children.some(child => child.path === path || child.relative === path)
    })
}

const hasPrivilegeBusiness = () => {
    const privilegeBusiness = router.app.$store.getters['objectBiz/privilegeBusiness']
    return !!privilegeBusiness.length
}

router.beforeEach(async (to, from, next) => {
    try {
        if (to.path !== '/status-error') {
            router.app.$store.commit('setGlobalLoading', true)
            await cancelRequest()
            await preload(router.app)
            if (to.meta.ignoreAuthorize) {
                next()
            } else if (hasAuthority(to)) {
                if (to.meta.requireBusiness && !hasPrivilegeBusiness()) {
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
                next({
                    path: '/status-403',
                    query: {
                        relative: to.path
                    }
                })
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
