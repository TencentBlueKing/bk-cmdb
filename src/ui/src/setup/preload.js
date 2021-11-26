import { getAuthorizedBusiness } from '@/router/business-interceptor.js'
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
  const [{ is_pass: globalConfigAuth }] = await verifyAuth([{
    action: 'update',
    resource_type: 'configAdmin'
  }])
  store.commit('globalConfig/setAuth', globalConfigAuth)
}

export default async function (app) {
  if (window.Site.authscheme === 'iam') {
    verifyPlatformManagementAuth()
  }
  getAuthorizedBusiness()
  return Promise.all([
    getGlobalConfig(app),
    getClassifications(app),
    getUserCustom(app),
    getGlobalUsercustom(app)
  ])
}
