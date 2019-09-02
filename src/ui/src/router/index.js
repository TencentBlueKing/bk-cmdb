/* eslint-disable */
import Vue from 'vue'
import Router from 'vue-router'

import StatusError from './StatusError.js'

import preload from '@/setup/preload'
import afterload from '@/setup/afterload'
import { translateAuth } from '@/setup/permission'
import $http from '@/api'

import index from '@/views/index/router.config'

import {
    MENU_INDEX,
    MENU_BUSINESS,
    MENU_RESOURCE,
    MENU_MODEL,
    MENU_ANALYSIS
} from '@/dictionary/menu-symbol'


import {
    businessViews,
    resourceViews,
    modelViews,
    analysisViews
} from '@/views'

import dynamicRouterView from '@/components/layout/dynamic-router-view'

Vue.use(Router)

export const viewRouters = []

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
        name: MENU_INDEX
    }
}]

const router = new Router({
    mode: 'hash',
    routes: [
        ...redirectRouters,
        ...statusRouters,
        ...index,
        {
            name: MENU_BUSINESS,
            component: dynamicRouterView,
            children: businessViews,
            path: '/business'
        }, {
            name: MENU_MODEL,
            component: dynamicRouterView,
            children: modelViews,
            path: '/model',
            redirect: '/model/index'
        },
        {
            name: MENU_RESOURCE,
            component: dynamicRouterView,
            children: resourceViews,
            path: '/resource',
            redirect: '/resource/index'
        }, {
            name: MENU_ANALYSIS,
            component: dynamicRouterView,
            children: analysisViews,
            path: '/analysis',
            redirect: '/analysis/audit'
        }
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

/* eslint-disable-next-line */
const setAuthScope = (to, from) => {
    const auth = to.meta.auth || {}
    if (typeof auth.setAuthScope === 'function') {
        auth.setAuthScope(to, from, router.app)
    }
}
/* eslint-disable-next-line */
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
    if (!to.meta.requireBusiness) {
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

const checkBusinessMenuRedirect = (to) => {
    const isBusinessMenu = to.matched.length > 1 && to.matched[0].name === MENU_BUSINESS
    if (!isBusinessMenu) {
        return false
    }
    return router.app.$store.state.objectBiz.bizId === null
}

const setAdminView = to => {
    const isAdminView = to.matched.length && to.matched[0].name !== MENU_BUSINESS
    router.app.$store.commit('setAdminView', isAdminView)
}

const setupStatus = {
    preload: true,
    afterload: true
}

router.beforeEach((to, from, next) => {
    Vue.nextTick(async () => {
        try {
            setLoading(true)
            if (setupStatus.preload) {
                await preload(router.app)
            }
            if (!isShouldShow(to)) {
                next({ name: MENU_INDEX })
            } else {
                // 防止直接进去申请业务的提示界面导致无法正确跳转权限中心
                if (to.name === 'requireBusiness' && !router.app.$store.getters.permission.length) {
                    return next({ name: MENU_INDEX })
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


                // 在业务菜单下刷新页面时，先重定向到一级路由，一级路由视图中的业务选择器设定成功后再跳转到二级视图
                const shouldRedirectToBusinessMenu = checkBusinessMenuRedirect(to)
                if (shouldRedirectToBusinessMenu) {
                    return next({ name: MENU_BUSINESS })
                }
                const isBusinessCheckPass = checkBusiness(to)
                if (!isBusinessCheckPass) {
                    await setPermission(to)
                    throw new StatusError({ name: 'requireBusiness', query: { _t: Date.now() } })
                }

                setAdminView(to)

                return next()
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
