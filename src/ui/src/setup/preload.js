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

import { getAuthorizedBusiness, getAuthorizedBusinessSet } from '@/router/business-interceptor.js'
import { verifyAuth } from '@/services/auth.js'
import store from '@/store'

const preloadConfig = {
  fromCache: false,
  cancelWhenRouteChange: false
}

export function getClassifications(app) {
  return app.$store.dispatch('objectModelClassify/searchClassificationsObjects', {
    params: {},
    config: {
      ...preloadConfig,
      requestId: 'post_searchClassificationsObjects'
    }
  })
}

export function getUserCustom(app) {
  return app.$store.dispatch('userCustom/searchUsercustom', {
    config: {
      ...preloadConfig,
      fromCache: false,
      requestId: 'post_searchUsercustom'
    }
  })
}

export function getGlobalUsercustom(app) {
  return app.$store.dispatch('userCustom/getGlobalUsercustom', {
    config: {
      ...preloadConfig,
      fromCache: false,
      globalError: false
    }
  }).catch(() => ({}))
}

/**
 * 初始化全局配置
 * @param {Object} app Vue 应用实例
 * @returns
 */
export async function getGlobalConfig(app) {
  return app.$store.dispatch('globalConfig/fetchConfig', {
    config: {
      ...preloadConfig,
      fromCache: false,
      globalError: false
    }
  })
}

/**
 * 验证平台管理模块的权限
 */
export const verifyPlatformManagementAuth = async () => {
  const [{ is_pass: isPass }] = await verifyAuth([{
    action: 'update',
    resource_type: 'configAdmin'
  }])

  if (isPass) {
    store.commit('globalConfig/setAuth', isPass)
  }
}

export default async function (app) {
  if (window.Site.authscheme === 'iam') {
    verifyPlatformManagementAuth()
  } else {
    // 开源版的可能没有 IAM，不需要鉴权
    store.commit('globalConfig/setAuth', true)
  }

  // 获取有访问权限的业务
  getAuthorizedBusiness()

  // 获取有访问权限的业务集
  getAuthorizedBusinessSet()

  return Promise.all([
    getGlobalConfig(app),
    getClassifications(app),
    getUserCustom(app),
    getGlobalUsercustom(app)
  ])
}
