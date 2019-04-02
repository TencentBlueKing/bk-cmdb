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

import {
    GET_AUTH_META,
    GET_MODEL_INST_AUTH_META
} from '@/dictionary/auth'

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

const getAuthMeta = (type, to, meta) => {
    if (meta === GET_MODEL_INST_AUTH_META) {
        const models = router.app.$store.getters['objectModelClassify/models']
        return GET_MODEL_INST_AUTH_META(to.params.objId, type, models)
    }
    return GET_AUTH_META(type)
}

const getAuth = to => {
    const auth = to.meta.auth || {}
    const view = auth.view
    const operation = Array.isArray(auth.operation) ? auth.operation : []
    const operationAuthMeta = operation.map(type => getAuthMeta(type, to, auth.meta))
    if (view) {
        operationAuthMeta.push(getAuthMeta(view, to, auth.meta))
    }
    if (operationAuthMeta.length) {
        return router.app.$store.dispatch('auth/getOperationAuth', operationAuthMeta)
    }
    return Promise.resolve([])
}

const isViewAuthorized = to => {
    const auth = to.meta.auth || {}
    const view = auth.view
    if (!view) {
        return true
    }
    const authMeta = getAuthMeta(view, to, auth.meta)
    const viewAuth = router.app.$store.getters['auth/isAuthorized'](authMeta.resource_type, authMeta.action)
    return viewAuth
}

const cancelRequest = () => {
    const allRequest = $http.queue.get()
    const requestQueue = allRequest.filter(request => request.cancelWhenRouteChange)
    return $http.cancel(requestQueue.map(request => request.requestId))
}

const setLoading = loading => router.app.$store.commit('setGlobalLoading', loading)

const setMenuState = to => {
    const menu = to.meta.menu || {}
    const menuId = menu.id
    const parentId = menu.parent
    router.app.$store.commit('menu/setActiveMenu', menuId)
    if (parentId) {
        router.app.$store.commit('menu/setOpenMenu', parentId)
    }
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
                const isStatusPage = statusRouters.some(status => status.name === to.name)
                if (isStatusPage) {
                    next()
                } else {
                    const auth = await getAuth(to)
                    const viewAuth = isViewAuthorized(to)
                    if (viewAuth) {
                        next()
                    } else {
                        setLoading(false)
                        next({ name: '403' })
                    }
                }
            }
        } catch (e) {
            setLoading(false)
            next({name: 'error'})
        } finally {
            setupStatus.preload = false
        }
    })
})

router.afterEach((to, from) => {
    try {
        if (setupStatus.afterload) {
            afterload(router.app)
        }
    } catch (e) {
        // ignore
    } finally {
        setupStatus.afterload = false
        setLoading(false)
    }
})

export default router
