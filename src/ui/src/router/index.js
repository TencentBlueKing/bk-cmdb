import Vue from 'vue'
import Router from 'vue-router'

import preload from '@/setup/preload'
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

Vue.use(Router)

const statusRouter = [
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

const router = new Router({
    mode: 'history',
    routes: [
        {
            path: '*',
            redirect: {
                name: '404'
            }
        }, {
            path: '/',
            redirect: {
                name: index.name
            }
        },
        ...statusRouter,
        index,
        audit,
        business,
        businessModel,
        customQuery,
        history,
        hosts,
        model,
        modelAssociation,
        modelTopology,
        process,
        resource
    ]
})

const getAuthMeta = (auth, to) => {
    if (typeof auth === 'function') {
        const dynamicParams = to.meta.dynamicParams || []
        auth = auth.apply(null, dynamicParams.map(param => to.params[param]))
    }
    const [type, action] = auth.split('.')
    return { type, action }
}

const checkAuth = async to => {
    const auth = to.meta.auth || {}
    const view = auth.view
    if (!view) {
        return Promise.resolve(true)
    }
    const operation = Array.isArray(auth.operation) ? auth.operation : []
    const authParams = [view, ...operation].map(auth => {
        const authMeta = getAuthMeta(auth, to)
        return {
            resource_type: authMeta.type,
            action: authMeta.action
        }
    })
    const $store = router.app.$store
    const authList = await $store.dispatch('auth/getAuthList', authParams)
    const authMeta = getAuthMeta(view, to)
    const viewAuth = $store.getters['auth/checkAuth'](authMeta.type, authMeta.action)
    return viewAuth
}

const cancelRequest = () => {
    const allRequest = $http.queue.get()
    const requestQueue = allRequest.filter(request => request.cancelWhenRouteChange)
    return $http.cancel(requestQueue.map(request => request.requestId))
}

const setLoading = loading => router.app.$store.commit('setGlobalLoading', loading)

router.beforeEach((to, from, next) => {
    const isStatusPage = statusRouter.some(status => status.name === to.name)
    if (isStatusPage) {
        next()
    } else {
        router.app.$nextTick(async () => {
            try {
                setLoading(true)
                await cancelRequest()
                await preload(router.app)
                const authorized = await checkAuth(to)
                if (authorized) {
                    next()
                } else {
                    next({ name: '403' })
                }
            } catch (e) {
                console.log(e)
                next({name: 'error'})
            }
        })
    }
})

router.afterEach((to, from) => {
    router.app.$nextTick(() => {
        setLoading(false)
    })
})

export default router
