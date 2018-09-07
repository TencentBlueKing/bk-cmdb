import Vue from 'vue'
import Router from 'vue-router'
import store from '@/store'
import preload from '@/setup/preload'
import $http from '@/api'

const index = () => import(/* webpackChunkName: index */ '@/views/index')
const model = () => import(/* webpackChunkName: model */ '@/views/model')
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
const error = () => import(/* webpackChunkName: error */ '@/views/status/error')

Vue.use(Router)

const router = new Router({
    linkActiveClass: 'active',
    routes: [{
        path: '/error',
        component: error
    }, {
        path: '/',
        redirect: '/index'
    }, {
        path: '/index',
        component: index
    }, {
        path: '/business',
        component: business
    }, {
        path: '/model/:classifyId',
        component: model,
        meta: {
            relative: '/model'
        }
    }, {
        path: '/model',
        component: model
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
        component: hosts
    }, {
        path: '/resource',
        component: resource
    }, {
        path: '/auditing',
        component: audit
    }, {
        path: '/topology',
        component: topology
    }, {
        path: '/process',
        component: process
    }]
})

const cancelRequest = () => {
    const allRequest = $http.queue.get()
    const requestQueue = allRequest.filter(request => request.cancelWhenRouteChange)
    return $http.cancel(requestQueue.map(request => request.requestId))
}

router.beforeEach(async (to, from, next) => {
    try {
        if (to.path !== '/error') {
            router.app.$store.commit('setGlobalLoading', true)
            await cancelRequest()
            await preload(router.app)
            next()
        } else {
            next()
        }
    } catch (e) {
        next('/error', e)
    }
})

router.afterEach(() => {
    router.app.$store.commit('setGlobalLoading', false)
})

export default router
