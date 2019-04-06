import Vue from 'vue'
import Router from 'vue-router'

import preload from '@/setup/preload'
import afterload from '@/setup/afterload'
import $http from '@/api'

import index from '@/views/index/router.config'
import audit from '@/views/audit/router.config'
import business from '@/views/business/router.config'
import businessModel from '@/views/business-model/router.config'
import customQuery from '@/views/custom-query/router.config'
import eventpush from '@/views/eventpush/router.config'
import history from '@/views/history/router.config'
import hosts from '@/views/hosts/router.config'
import model from '@/views/model-manage/router.config'
import modelAssociation from '@/views/model-association/router.config'
import modelTopology from '@/views/model-topology/router.config'
import process from '@/views/process/router.config'
import resource from '@/views/resource/router.config'
import topology from '@/views/topology/router.config'
import generalModel from '@/views/general-model/router.config'

Vue.use(Router)

export const viewRouters = [
    index,
    audit,
    businessModel,
    customQuery,
    eventpush,
    history,
    hosts,
    modelAssociation,
    modelTopology,
    process,
    resource,
    topology,
    ...generalModel,
    ...business,
    ...model
]

const statusRouters = [
    {
        name: '403',
        path: '/403',
        components: require('@/views/status/403')
    }, {
        name: '404',
        path: '/404',
        components: require('@/views/status/404')
    }, {
        name: 'error',
        path: '/error',
        components: require('@/views/status/error')
    }, {
        name: 'requireBusiness',
        path: '/require-business',
        components: require('@/views/status/require-business')
    }
]

const redirectRouters = [{
    path: '*',
    redirect: {
        name: '404'
    }
}, {
    path: '/',
    redirect: {
        name: index.name
    }
}]

const router = new Router({
    mode: 'hash',
    routes: [
        ...redirectRouters,
        ...statusRouters,
        ...viewRouters
    ]
})

const getAuth = to => {
    const auth = to.meta.auth || {}
    const view = auth.view
    const operation = auth.operation || []
    const routerAuth = view ? [view, ...operation] : operation
    if (routerAuth.length) {
        return router.app.$store.dispatch('auth/getAuth', {
            type: 'operation',
            list: routerAuth
        })
    }
    return Promise.resolve([])
}

const isViewAuthorized = to => {
    const auth = to.meta.auth || {}
    const view = auth.view
    if (!view) {
        return true
    }
    const viewAuth = router.app.$store.getters['auth/isAuthorized'](view)
    return viewAuth
}

const cancelRequest = () => {
    const allRequest = $http.queue.get()
    const requestQueue = allRequest.filter(request => request.cancelWhenRouteChange)
    return $http.cancel(requestQueue.map(request => request.requestId))
}

const setLoading = loading => router.app.$store.commit('setGlobalLoading', loading)

const setMenuState = to => {
    const isStatusRoute = statusRouters.some(route => route.name === to.name)
    if (isStatusRoute) {
        return false
    }
    const menu = to.meta.menu || {}
    const menuId = menu.id
    const parentId = menu.parent
    router.app.$store.commit('menu/setActiveMenu', menuId)
    if (parentId) {
        router.app.$store.commit('menu/setOpenMenu', parentId)
    }
}

const checkAuthDynamicMeta = (to, from) => {
    router.app.$store.commit('auth/setDynamicMeta', {})
    const auth = to.meta.auth || {}
    const setDynamicMeta = auth.setDynamicMeta
    if (typeof setDynamicMeta === 'function') {
        setDynamicMeta(to, from, router.app)
    }
}

const checkBusiness = to => {
    const getters = router.app.$store.getters
    const isAdminView = getters.isAdminView
    if (isAdminView || !to.meta.requireBusiness) {
        return true
    }
    const authorizedBusiness = getters['objectBiz/authorizedBusiness']
    return authorizedBusiness.length
}

const isShouldShow = to => {
    const isAdminView = router.app.$store.getters.isAdminView
    const menu = to.meta.menu
    if (isAdminView && menu) {
        return menu.adminView
    }
    return true
}

const setupStatus = {
    preload: true,
    afterload: true
}

router.beforeEach((to, from, next) => {
    Vue.nextTick(async () => {
        try {
            if (!isShouldShow(to)) {
                next({ name: index.name })
            } else {
                setLoading(true)
                setMenuState(to)
                await cancelRequest()
                if (setupStatus.preload) {
                    await preload(router.app)
                }
                checkAuthDynamicMeta(to, from)
                await getAuth(to)
                const viewAuth = isViewAuthorized(to)
                if (viewAuth) {
                    const isBusinessCheckPass = checkBusiness(to)
                    if (isBusinessCheckPass) {
                        next()
                    } else {
                        setLoading(false)
                        next({ name: 'requireBusiness' })
                    }
                } else {
                    setLoading(false)
                    next({ name: '403' })
                }
            }
        } catch (e) {
            console.error(e)
            setLoading(false)
            next({ name: 'error' })
        } finally {
            setupStatus.preload = false
        }
    })
})

router.afterEach((to, from) => {
    try {
        if (setupStatus.afterload) {
            afterload(router.app, to, from)
        }
    } catch (e) {
        // ignore
    } finally {
        setupStatus.afterload = false
        setLoading(false)
    }
})

export default router
