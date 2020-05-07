import Vue from 'vue'
import Router from 'vue-router'

import StatusError from './StatusError.js'

import preload from '@/setup/preload'
import afterload from '@/setup/afterload'

import index from '@/views/index/router.config'

import {
    before as businessBeforeInterceptor
} from './business-interceptor'

import {
    MENU_ENTRY,
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
}]

const router = new Router({
    mode: 'hash',
    routes: [
        ...redirectRouters,
        ...statusRouters,
        {
            name: MENU_ENTRY,
            component: dynamicRouterView,
            children: [...index],
            path: '/',
            redirect: '/index'
        },
        {
            name: MENU_BUSINESS,
            components: {
                default: dynamicRouterView,
                permission: require('@/views/status/non-exist-business').default
            },
            children: businessViews,
            path: '/business/:bizId?'
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

const beforeHooks = new Set()

function runBeforeHooks () {
    return Promise.all(Array.from(beforeHooks).map(callback => callback()))
}

export const addBeforeHooks = function (hook) {
    beforeHooks.add(hook)
}

const checkViewAuthorize = async to => {
    // owener判断已经发现无业务时
    if (to.meta.view === 'permission') {
        return false
    }
    const auth = to.meta.auth || {}
    const view = auth.view
    if (view) {
        const viewAuthData = typeof view === 'function' ? view(to, router.app) : view
        const viewAuth = await router.app.$store.dispatch('auth/getViewAuth', viewAuthData)
        to.meta.view = viewAuth ? 'default' : 'permission'
    }
    return Promise.resolve()
}

const setLoading = loading => router.app.$store.commit('setGlobalLoading', loading)

const checkAvailable = (to, from) => {
    if (typeof to.meta.checkAvailable === 'function') {
        return to.meta.checkAvailable(to, from, router.app)
    } else if (to.meta.hasOwnProperty('available')) {
        return to.meta.available
    }
    return true
}

// 因产品形态调整，去掉了管理模式与业务模式，为避免修改过多逻辑，此处做兼容处理
const setAdminView = to => {
    const isAdminView = to.matched.length && to.matched[0].name !== MENU_BUSINESS
    router.app.$store.commit('setAdminView', isAdminView)
}

const setupStatus = {
    preload: true,
    afterload: true
}

router.beforeEach((to, from, next) => {
    router.app.$store.commit('setTitle', '')
    if (to.name === from.name) {
        return next()
    }
    Vue.nextTick(async () => {
        try {
            setLoading(true)
            if (setupStatus.preload) {
                await preload(router.app)
            }
            await runBeforeHooks()
            const shouldContinue = await businessBeforeInterceptor(router.app, to, from, next)
            if (!shouldContinue) {
                return false
            }
            setAdminView(to)

            const isAvailable = checkAvailable(to, from)
            if (!isAvailable) {
                throw new StatusError({ name: '404' })
            }
            await checkViewAuthorize(to)
            return next()
        } catch (e) {
            console.error(e)
            if (e.__CANCEL__) {
                next()
            } else if (e instanceof StatusError) {
                next({ name: e.name, query: e.query })
            } else if (e.status !== 401) {
                console.error(e)
                next({ name: 'error' })
            } else {
                next()
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
