/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import Vue from 'vue'
import Router from 'vue-router'
import has from 'has'

import StatusError from './StatusError.js'
import { rootPath, connectToMain } from '@blueking/sub-saas'

import preload from '@/setup/preload'
import afterload from '@/setup/afterload'
import { setupValidator } from '@/setup/validate'
import { $error } from '@/magicbox'
import i18n from '@/i18n'
import { changeDocumentTitle } from '@/utils/change-document-title'
import { OPERATION } from '@/dictionary/iam-auth'
import workerTask from '@/setup/worker-task'

import {
  before as businessBeforeInterceptor
} from './business-interceptor'

import {
  MENU_ENTRY,
  MENU_INDEX,
  MENU_BUSINESS_SET,
  MENU_BUSINESS,
  MENU_RESOURCE,
  MENU_MODEL,
  MENU_MODEL_MANAGEMENT,
  MENU_RESOURCE_MANAGEMENT,
  MENU_ANALYSIS,
  MENU_ANALYSIS_AUDIT,
  MENU_PLATFORM_MANAGEMENT,
  MENU_PLATFORM_MANAGEMENT_GLOBAL_CONFIG
} from '@/dictionary/menu-symbol'

import {
  indexViews,
  hostLandingViews,
  businessSetViews,
  businessViews,
  resourceViews,
  modelViews,
  analysisViews,
  platformManagementViews
} from '@/views'

import dynamicRouterView from '@/components/layout/dynamic-router-view'

Vue.use(Router)

export const getRoutePath = (subPath) => {
  const path = subPath.startsWith('/') ? subPath.slice(1) : subPath
  return `${rootPath}${path}`
}

const statusRouters = [
  {
    name: 'no-business',
    path: getRoutePath('/no-business'),
    components: require('@/views/status/non-exist-business')
  },
  {
    name: '404',
    path: getRoutePath('/404'),
    components: require('@/views/status/404')
  },
  {
    name: 'error',
    path: getRoutePath('/error'),
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
    ...hostLandingViews,
    {
      name: MENU_ENTRY,
      component: dynamicRouterView,
      children: indexViews,
      path: rootPath,
      redirect: { name: MENU_INDEX }
    },
    {
      name: MENU_BUSINESS_SET,
      components: {
        default: dynamicRouterView,
        error: require('@/views/status/error').default,
        permission: require('@/views/status/non-exist-business-set').default
      },
      children: businessSetViews,
      path: getRoutePath('/business-set/:bizSetId'),
      meta: {
        menu: {
          i18n: '业务集'
        }
      }
    },
    {
      name: MENU_BUSINESS,
      components: {
        default: dynamicRouterView,
        error: require('@/views/status/error').default,
        permission: require('@/views/status/non-exist-business').default
      },
      children: businessViews,
      path: getRoutePath('/business/:bizId?'),
      meta: {
        menu: {
          i18n: '业务'
        }
      }
    },
    {
      name: MENU_MODEL,
      component: dynamicRouterView,
      children: modelViews,
      path: getRoutePath('/model'),
      redirect: { name: MENU_MODEL_MANAGEMENT },
      meta: {
        menu: {
          i18n: '模型'
        }
      }
    },
    {
      name: MENU_RESOURCE,
      component: dynamicRouterView,
      children: resourceViews,
      path: getRoutePath('/resource'),
      redirect: { name: MENU_RESOURCE_MANAGEMENT },
      meta: {
        menu: {
          i18n: '资源'
        }
      }
    },
    {
      name: MENU_ANALYSIS,
      component: dynamicRouterView,
      children: analysisViews,
      path: getRoutePath('/analysis'),
      redirect: { name: MENU_ANALYSIS_AUDIT },
      meta: {
        menu: {
          i18n: '运营分析'
        }
      }
    },
    {
      name: MENU_PLATFORM_MANAGEMENT,
      component: dynamicRouterView,
      children: platformManagementViews,
      path: getRoutePath('/platform-management'),
      redirect: { name: MENU_PLATFORM_MANAGEMENT_GLOBAL_CONFIG },
      meta: {
        auth: {
          view: [{ type: OPERATION.R_CONFIG_ADMIN }, { type: OPERATION.U_CONFIG_ADMIN }]
        },
        menu: {
          i18n: '平台管理'
        }
      }
    }
  ]
})

const beforeHooks = new Set()

function runBeforeHooks() {
  return Promise.all(Array.from(beforeHooks).map(callback => callback()))
}

export const addBeforeHooks = function (hook) {
  beforeHooks.add(hook)
}

function cancelRequest(app) {
  const pendingRequest = app.$http.queue.get()
  const cancelId = pendingRequest.filter(request => request.cancelWhenRouteChange).map(request => request.requestId)
  app.$http.cancelRequest(cancelId)
}

const checkViewAuthorize = async (to) => {
  // 使用就近原则向上回溯，找到路由的auth.view配置
  const findViewAuth = (route, key) => {
    if (!route) {
      return
    }
    if (route?.meta?.auth?.[key]) {
      return route.meta.auth[key]
    }
    return findViewAuth(route?.parent, key)
  }

  const getViewAuthResult = (authView) => {
    const viewAuthData = typeof authView === 'function' ? authView(to, router.app) : authView
    return router.app.$store.dispatch('auth/getViewAuth', viewAuthData)
  }

  const { matched } = to
  const latestRoute = matched[matched.length - 1]

  const authSuperView = findViewAuth(latestRoute, 'superView')
  const authView = findViewAuth(latestRoute, 'view')

  // 存在superView和authView权限
  if (authSuperView && authView) {
    const authSuperViewResult = await getViewAuthResult(authSuperView)
    const authViewResult = await getViewAuthResult(authView)

    // 没有子权限，指定需要优先申请子权限，后期期望自动关联上父权限
    if (!authViewResult) {
      to.meta.authKey = 'view'
    } else if (!authSuperViewResult) {
      to.meta.authKey = 'superView'
    }

    // 没有父权限时才拦截入口，无权限申请时根据authKey确定需要申请哪一个权限
    // 没有子权限允许进入到页面，在页面中捕获接口无权限处理
    to.meta.view = authSuperViewResult ? 'default' : 'permission'

    // 未同时配置superView（配置在父级），但希望校验子级配置的view权限
    if (authSuperViewResult && !latestRoute.meta.auth.superView && !authViewResult) {
      to.meta.authKey = 'view'
      to.meta.view = 'permission'
    }
  } else if (authView) {
    const authViewResult = await getViewAuthResult(authView)
    to.meta.view = authViewResult ? 'default' : 'permission'
  }

  return Promise.resolve()
}

const setLoading = loading => router.app.$store.commit('setGlobalLoading', loading)

const checkAvailable = (to, from) => {
  if (typeof to.meta.checkAvailable === 'function') {
    return to.meta.checkAvailable(to, from, router.app)
  }
  if (has(to.meta, 'available')) {
    return to.meta.available
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
      /**
       * 取消上个页面中的所有请求
       */
      cancelRequest(router.app)

      /**
       * 设置当前页面的标题
       */
      to.name !== from.name && router.app.$store.commit('setTitle', '')


      /**
       * 将非 permission 的 view 设置为 default
       */
      if (to.meta.view !== 'permission') {
        Vue.set(to.meta, 'view', 'default')
      }

      /**
       * 如果路由中的业务 ID 改变了，应该继续往下执行到业务拦截器，否则在同页路由下直接执行路由跳转，不再执行往后的逻辑
       */
      const isFromBiz = from.matched[0]?.name === MENU_BUSINESS
      const bizIsChanged = isFromBiz && parseInt(to.params.bizId, 10) !== parseInt(from.params.bizId, 10)

      if (to.name === from.name && !bizIsChanged) {
        return next()
      }

      /**
       * 初始化预加载，只会加载一次
       */
      if (setupStatus.preload) {
        setLoading(true)
        setupStatus.preload = false
        await preload(router.app)
        workerTask.run()
        setupValidator(router.app)
      }

      /**
       * 执行插入的钩子
       */
      await runBeforeHooks()

      /**
       * 业务拦截器，检查是否有当前业务的权限
       */
      const shouldContinue = await businessBeforeInterceptor(to, from, next)

      if (!shouldContinue) {
        return false
      }

      /**
       * 检查页面是否被设置为 available，如果设置为 false 则会跳转到 404 页面
       */
      const isAvailable = checkAvailable(to, from)
      if (!isAvailable) {
        throw new StatusError({ name: '404' })
      }

      /**
       * 检查是否有权限访问当前页面
       */
      await checkViewAuthorize(to, router.app)

      /**
       * 执行路由配置中的before钩子
       */
      if (to.meta?.before && !await to.meta?.before?.(to, from, router.app)) {
        return false
      }

      return next()
    } catch (e) {
      console.error(e)
      setupStatus.preload = true
      // eslint-disable-next-line no-underscore-dangle
      if (e.__CANCEL__) {
        next()
      } else if (e instanceof StatusError) {
        next({ name: e.name, query: e.query })
      } else if (e.status !== 401) {
        console.error(e)
        // 保留路由，将视图切换为error
        Vue.set(to.meta, 'view', 'error')
        next()
      } else {
        next()
      }
    } finally {
      setLoading(false)
    }
  })
})

router.afterEach(async (to, from) => {
  try {
    if (setupStatus.afterload) {
      setupStatus.afterload = false
      await afterload(router.app, to, from)
    }
    changeDocumentTitle()
  } catch (e) {
    setupStatus.afterload = true
    console.error(e)
  }
})

router.onError((error) => {
  if (/Loading chunk (\d*) failed/.test(error.message)) {
    $error(i18n.t('资源请求失败提示'))
  }
})

connectToMain(router)

export const useRouter = () => router

export const useRoute = () => router.app.$route

export default router
