import Vue from 'vue'
import Router from 'vue-router'

import StatusError from './StatusError.js'

import preload from '@/setup/preload'
import afterload from '@/setup/afterload'
import { translateAuth } from '@/setup/permission'
import $http from '@/api'

import index from '@/views/index/router.config'
import audit from '@/views/audit/router.config'
import business from '@/views/business/router.config'
import businessModel from '@/views/business-model/router.config'
import businessTopology from '@/views/business-topology/router.config'
import customQuery from '@/views/custom-query/router.config'
import eventpush from '@/views/eventpush/router.config'
import history from '@/views/history/router.config'
import hosts from '@/views/hosts/router.config'
import hostDetails from '@/views/host-details/router.config'
import model from '@/views/model-manage/router.config'
import modelAssociation from '@/views/model-association/router.config'
import modelTopology from '@/views/model-topology/router.config'
import resource from '@/views/resource/router.config'
import generalModel from '@/views/general-model/router.config'
import permission from '@/views/permission/router.config'
import template from '@/views/service-template/router.config'
import category from '@/views/service-category/router.config'

import serviceInstance from '@/views/service-instance/router.config'
import synchronous from '@/views/business-synchronous/router.config'

Vue.use(Router)

export const viewRouters = [
    ...index,
    audit,
    businessModel,
    businessTopology,
    customQuery,
    eventpush,
    hosts,
    ...hostDetails,
    modelAssociation,
    modelTopology,
    resource,
    ...template,
    ...generalModel,
    ...business,
    ...model,
    ...permission,
    category,
    synchronous,
    ...serviceInstance,
    history
]

const indexName = index[0].name

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
        name: indexName
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
    if (!to.meta.resetMenu) {
        return false
    }
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

const setTitle = to => {
    const { i18nTitle, title } = to.meta
    let headerTitle
    if (!i18nTitle && !title) {
        return false
    } else if (i18nTitle) {
        headerTitle = router.app.$t(i18nTitle)
    } else if (title) {
        headerTitle = title
    }
    router.app.$store.commit('setHeaderTitle', headerTitle)
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

const setPermission = async to => {
    const permission = []
    const authMeta = to.meta.auth
    if (authMeta) {
        const { view, operation } = authMeta
        const auth = [...operation]
        if (view) {
            auth.push(view)
        }
        const translated = await translateAuth(auth)
        permission.push(...translated)
    }
    router.app.$store.commit('setPermission', permission)
    return permission
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
                next({ name: indexName })
            } else {
                setMenuState(to)
                setTitle(to)
                setAuthScope(to, from)
                checkAuthDynamicMeta(to, from)

                if (to.name === 'requireBusiness' && !router.app.$store.getters.permission.length) {
                    next({ name: indexName })
                }

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
                    await setPermission(to)
                    throw new StatusError({ name: 'requireBusiness', query: { _t: Date.now() } })
                }
                next()
            }
        } catch (e) {
            if (e.__CANCEL__) {
                next()
            } else if (e instanceof StatusError) {
                next({ name: e.name, query: e.query })
            } else {
                console.error(e)
                next({ name: 'error' })
            }
        } finally {
            setLoading(false)
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
        console.error(e)
    } finally {
        setLoading(false)
    }
})

export default router
