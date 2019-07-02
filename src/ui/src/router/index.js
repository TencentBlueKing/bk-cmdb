import Vue from 'vue'
import Router from 'vue-router'

import StatusError from './StatusError.js'

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
import hostDetails from '@/views/host-details/router.config'
import model from '@/views/model-manage/router.config'
import modelAssociation from '@/views/model-association/router.config'
import modelTopology from '@/views/model-topology/router.config'
import process from '@/views/process/router.config'
import resource from '@/views/resource/router.config'
import topology from '@/views/topology/router.config'
import generalModel from '@/views/general-model/router.config'
import permission from '@/views/permission/router.config'

Vue.use(Router)

export const viewRouters = [
    index,
    audit,
    businessModel,
    customQuery,
    eventpush,
    history,
    hosts,
    ...hostDetails,
    modelAssociation,
    modelTopology,
    process,
    resource,
    topology,
    ...generalModel,
    ...business,
    ...model,
    ...permission
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

const setAuthScope = (to, from) => {
    const auth = to.meta.auth || {}
    if (typeof auth.setAuthScope === 'function') {
        auth.setAuthScope(to, from, router.app)
    }
}

const checkAuthDynamicMeta = (to, from) => {
    router.app.$store.commit('auth/clearDynamicMeta')
    const auth = to.meta.auth || {}
    const setDynamicMeta = auth.setDynamicMeta
    if (typeof setDynamicMeta === 'function') {
        setDynamicMeta(to, from, router.app)
    }
}

const checkAvailable = (to, from) => {
    if (typeof to.meta.checkAvailable === 'function') {
        return to.meta.checkAvailable(to, from, router.app)
    }
    return true
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
    return menu
        ? isAdminView
            ? menu.adminView
            : menu.businessView
        : true
}

const setupStatus = {
    preload: true,
    afterload: true
}

router.beforeEach((to, from, next) => {
    Vue.nextTick(async () => {
        try {
            setLoading(true)
            await cancelRequest()
            if (setupStatus.preload) {
                await preload(router.app)
            }
            if (!isShouldShow(to)) {
                next({ name: index.name })
            } else {
                setMenuState(to)
                setAuthScope(to, from)
                checkAuthDynamicMeta(to, from)

                const isAvailable = checkAvailable(to, from)
                if (!isAvailable) {
                    throw new StatusError({ name: '404' })
                }
                await getAuth(to)
                const viewAuth = isViewAuthorized(to)
                if (!viewAuth) {
                    throw new StatusError({ name: '403' })
                }

                const isBusinessCheckPass = checkBusiness(to)
                if (!isBusinessCheckPass) {
                    throw new StatusError({ name: 'requireBusiness' })
                }

                next()
            }
        } catch (e) {
            if (e.__CANCEL__) {
                next()
            } else if (e instanceof StatusError) {
                next({ name: e.name })
            } else {
                console.error(e)
                next({ name: 'error' })
            }
            setLoading(false)
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
        setLoading(false)
    }
})

export default router
