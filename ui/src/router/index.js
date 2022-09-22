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

import preload from '@/setup/preload'
import afterload from '@/setup/afterload'
import { setupValidator } from '@/setup/validate'
import { $error } from '@/magicbox'
import i18n from '@/i18n'
import { changeDocumentTitle } from '@/utils/change-document-title'
import { OPERATION } from '@/dictionary/iam-auth'

import {
  before as businessBeforeInterceptor
} from './business-interceptor'

import {
  MENU_ENTRY,
  MENU_BUSINESS_SET,
  MENU_BUSINESS,
  MENU_RESOURCE,
  MENU_MODEL,
  MENU_ANALYSIS,
  MENU_PLATFORM_MANAGEMENT
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
    ...hostLandingViews,
    {
      name: MENU_ENTRY,
      component: dynamicRouterView,
      children: indexViews,
      path: '/',
      redirect: '/index'
    },
    {
      name: MENU_BUSINESS_SET,
      components: {
        default: dynamicRouterView,
        error: require('@/views/status/error').default,
        permission: require('@/views/status/non-exist-business-set').default
      },
      children: businessSetViews,
      path: '/business-set/:bizSetId',
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
      path: '/business/:bizId?',
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
      path: '/model',
      redirect: '/model/index',
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
      path: '/resource',
      redirect: '/resource/index',
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
      path: '/analysis',
      redirect: '/analysis/audit',
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
      path: '/platform-management',
      redirect: '/platform-management/global-config',
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

// eslint-disable-next-line no-unused-vars
const checkViewAuthorize = async to => Promise.resolve()

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
      await checkViewAuthorize(to)

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
export default router
