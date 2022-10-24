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
import businessService from '@/service/business/search.js'
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

  // 获取有权限和全部业务列表（带缓存）
  const [authorizedList, allBusinessList] = await Promise.all([
    getAuthorizedBusiness(),
    businessService.findAll()
  ])

  const id = parseInt(to.params.bizId || window.localStorage.getItem('selectedBusiness'), 10)
  const business = allBusinessList.find(business => business.bk_biz_id === id)
  const hasURLId = to.params.bizId
  const isAuthorized = authorizedList.some(item => item.bk_biz_id === business?.bk_biz_id)

  // 缓存无ID，URL无ID，则认为是首次进入业务导航，取一个默认业务id进入到二级路由
  if (!id) {
    // 优先取有权限业务的第一个写入URL中，否则取系统的第一个业务
    const firstBusiness = authorizedList?.length ? authorizedList?.[0] : allBusinessList?.[0]
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

  const isSubRoute = to.matched.length > 1
  toTopRoute.meta.view = 'default'
  window.localStorage.setItem('selectedBusiness', id)
  store.commit('objectBiz/setBizId', id)
  setBizSetRecentlyUsed(false)

  // 补齐业务id，如果是一级路由，则重定向到带业务id的二级路由首页(业务拓扑)
  if (!isSubRoute) {
    // next执行完之后，会再次进入route.beforeEach即会再次进入到此拦截器中，此时的route为next中指定的
    next({
      path: `/business/${id}/index`,
      replace: true
    })
    return false
  }

  // 补齐业务id，如果是二级路由且URL中不包含业务ID，则补充业务ID到URL中
  if (!hasURLId) {
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

  // 业务不存在或无权限
  if (!business || !isAuthorized) {
    // 优先使用二级路由（内页）展示无权限
    const targetRoute = to.matched?.[1] ?? to.matched?.[0]
    targetRoute.meta.view = 'permission'
  }

  // 总是放行，因为无论如何都需要进入到二级路由，前提是之前的逻辑已经保证了路由的正确性
  return true
}
