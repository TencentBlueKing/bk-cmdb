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

import { MENU_BUSINESS, MENU_BUSINESS_SET } from '@/dictionary/menu-symbol'
import businessSetService from '@/service/business-set/index.js'
import {
  setBizSetIdToStorage,
  setBizSetRecentlyUsed
} from '@/utils/business-set-helper.js'
import store from '@/store'

const requestId = Symbol('getAuthorizedBusiness')

let committed = false
export async function getAuthorizedBusiness() {
  const { info } = await store.dispatch('objectBiz/getAuthorizedBusiness', {
    requestId,
    fromCache: true
  })
  if (!committed) {
    store.commit('objectBiz/setAuthorizedBusiness', Object.freeze(info))
    committed = true
  }
  return info
}

export const getAuthorizedBusinessSet = async () => businessSetService.getAuthorizedWithCache()

export const before = async function (to, from, next) {
  const [toTopRoute] = to.matched
  const [fromTopRoute] = from.matched

  // 业务集合法性检测
  if (toTopRoute?.name === MENU_BUSINESS_SET) {
    const bizSetId = Number(to.params.bizSetId)
    const authorizedList = await getAuthorizedBusinessSet()
    const found = authorizedList.some(item => item.bk_biz_set_id === bizSetId)
    store.commit('bizSet/setBizSetList', authorizedList)
    if (!found) {
      toTopRoute.meta.view = 'permission'
      next()
      return false
    }
  }

  // 记录上一次是否使用的是业务集视图并且保存id值
  const isMatchedBusinessSetView = fromTopRoute?.name === MENU_BUSINESS_SET || toTopRoute?.name === MENU_BUSINESS_SET
  const availableBusinessSetView = isMatchedBusinessSetView && fromTopRoute?.meta?.view !== 'permission'

  if (availableBusinessSetView) {
    setBizSetIdToStorage(from.params.bizSetId || to.params.bizSetId)
    store.commit('bizSet/setBizSetId', from.params.bizSetId || to.params.bizSetId)
    store.commit('bizSet/setBizId', from.query.bizId || to.query.bizId)
    setBizSetRecentlyUsed(true)
  }

  // 从业务集视图跳出的时候重置view为默认防止再次进入时仍停留在permission
  if (toTopRoute?.name !== MENU_BUSINESS_SET && fromTopRoute?.name === MENU_BUSINESS_SET) {
    fromTopRoute.meta.view = 'default'
  }

  // 拦截非业务视图的路由
  if (!toTopRoute || toTopRoute.name !== MENU_BUSINESS) {
    if (fromTopRoute && fromTopRoute.name === MENU_BUSINESS) {
      fromTopRoute.meta.view = 'default'
    }
    return true
  }

  const newBizId = parseInt(to.params.bizId, 10)
  const oldBizId = parseInt(from.params.bizId, 10)

  if (fromTopRoute && fromTopRoute.name === MENU_BUSINESS && newBizId !== oldBizId) {
    window.location.hash = to.fullPath
    window.location.reload()
    return false
  }

  if (toTopRoute.meta.view === 'permission') {
    next()
    return false
  }

  const authorizedList = await getAuthorizedBusiness()
  const id = parseInt(to.params.bizId || window.localStorage.getItem('selectedBusiness'), 10)
  const business = authorizedList.find(business => business.bk_biz_id === id)
  const hasURLId = to.params.bizId

  // URL或者缓存中的id对应的业务存在
  if (business) {
    const isSubRoute = to.matched.length > 1
    toTopRoute.meta.view = 'default'
    window.localStorage.setItem('selectedBusiness', id)
    store.commit('objectBiz/setBizId', id)
    setBizSetRecentlyUsed(false)

    if (!isSubRoute) { // 如果是一级路由，则重定向到带业务id的二级路由首页(业务拓扑)
      next({
        path: `/business/${id}/index`,
        replace: true
      })
      return false
    }
    if (!hasURLId) { // 如果是二级路由且URL中不包含业务ID，则补充业务ID到URL中
      next({
        name: to.name,
        params: {
          ...to.params,
          bizId: id
        },
        query: to.query,
        replace: true
      })
      return false
    }
    return true // 正常的有权限的业务，且URL中带了ID，则直接返回，进行后续的路由逻辑
  }

  // 未找到对应有权限的业务，且URL中有业务ID，则显示一级view的无权限视图
  if (hasURLId) {
    toTopRoute.meta.view = 'permission'
    next()
    return false
  }

  // 缓存无ID，URL无ID，则认为是首次进入业务导航，取有权限业务的第一个写入URL中
  if (authorizedList.length) {
    const [firstBusiness] = authorizedList
    toTopRoute.meta.view = 'default'
    const defaultId = firstBusiness.bk_biz_id
    window.localStorage.setItem('selectedBusiness', defaultId)
    store.commit('objectBiz/setBizId', defaultId)
    next({
      path: `/business/${defaultId}/index`,
      replace: true
    })
    return false
  }

  toTopRoute.meta.view = 'permission'
  next('/business')
  return false
}
