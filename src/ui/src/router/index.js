import Vue from 'vue'
import Router from 'vue-router'

import StatusError from './StatusError.js'

import preload from '@/setup/preload'
import afterload from '@/setup/afterload'

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
            path: '/business',
            redirect: '/business/host'
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
    const view = auth.view || {}
    const operation = auth.operation || {}
    const routerAuth = Object.values({ ...view, ...operation })
    if (routerAuth.length) {
        return router.app.$store.dispatch('auth/getAuth', {
            type: 'operation',
            list: routerAuth
        })
    }
    return Promise.resolve([])
}

const checkViewAuthorized = to => {
    const auth = to.meta.auth || {}
    const view = auth.view
    if (view) {
        const viewAuth = router.app.$isAuthorized(Object.values(view))
        if (!viewAuth) {
            to.meta.view = 'permission'
        }
    }
}

const setLoading = loading => router.app.$store.commit('setGlobalLoading', loading)

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
    } else if (to.meta.hasOwnProperty('available')) {
        return to.meta.available
    }
    return true
}

const setAdminView = to => {
    const isAdminView = to.matched.length && to.matched[0].name !== MENU_BUSINESS
    router.app.$store.commit('setAdminView', isAdminView)
}

// 进入业务二级导航时需要先加载业务
// 在App.vue中添加一个隐藏的业务选择器，业务选择器完成设置后resolve对应的promise
const checkOwner = async to => {
    const matched = to.matched
    if (matched.length && matched[0].name === MENU_BUSINESS) {
        router.app.$store.commit('setBusinessSelectorVisible', true)
        const result = await router.app.$store.state.businessSelectorPromise
        to.meta.view = result ? 'default' : 'permission'
    } else {
        to.meta.view = 'default'
        router.app.$store.commit('setBusinessSelectorVisible', false)
    }
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
            await checkOwner(to)
            setAdminView(to)
            setAuthScope(to, from)
            checkAuthDynamicMeta(to, from)

            const isAvailable = checkAvailable(to, from)
            if (!isAvailable) {
                throw new StatusError({ name: '404' })
            }
            await getAuth(to)
            checkViewAuthorized(to)
            return next()
        } catch (e) {
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
        router.app.$store.commit('setTitle', '')
        router.app.$store.commit('setBreadcrumbs', [])
    } catch (e) {
        console.error(e)
    } finally {
        setLoading(false)
    }
})

export default router
