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

import index from '@/views/index/router.config'
import hostLanding from '@/views/host-details/router.config'

import audit from '@/views/audit/router.config'
import business from '@/views/business/router.config'
import customQuery from '@/views/dynamic-group/router.config'
import model from '@/views/model-manage/router.config'
import modelAssociation from '@/views/model-association/router.config'
import modelTopology from '@/views/model-topology/router.config'
import resource from '@/views/resource/router.config'
import generalModel from '@/views/general-model/router.config'
import serviceTemplate from '@/views/service-template/router.config'
import serviceCategory from '@/views/service-category/router.config'
import serviceInstance from '@/views/service-instance/router.config'
import serviceSynchronous from '@/views/business-synchronous/router.config'
import resourceManagement from '@/views/resource-manage/router.config'
import customFields from '@/views/custom-fields/router.config'
import operation from '@/views/operation/router.config'

import setSync from '@/views/set-sync/router.config'
import setTemplate from '@/views/set-template/router.config'

import hostApply from '@/views/host-apply/router.config'
import businessTopology from '@/views/business-topology/router.config'

import cloudArea from '@/views/cloud-area/router.config'
import cloudAccount from '@/views/cloud-account/router.config'
import cloudResource from '@/views/cloud-resource/router.config'
import project from '@/views/project/router.config'

// 字段组合模板
import fieldTemplate from '@/views/field-template/router.config'

// 业务集实例
import businessSet from '@/views/business-set/router.config'

import businessSetTopology from '@/views/business-set-topology/router.config.js'

import statusPermission from '@/views/status/permission'
import statusNonExistBusiness from '@/views/status/non-exist-business.vue'
import statusError from '@/views/status/error'

/**
 * 平台管理
 */
import globalConfig from '@/views/global-config/router.config'

import { MENU_BUSINESS } from '@/dictionary/menu-symbol'

const flatternViews = (views) => {
  const flatterned = []
  views.forEach((view) => {
    if (Array.isArray(view)) {
      flatterned.push(...view)
    } else {
      flatterned.push(view)
    }
  })
  return flatterned
}

export const injectStatusComponents = (views) => {
  views.forEach((view) => {
    view.components = {
      default: view.component,
      error: statusError
    }

    // 业务视图的permission是定制的，将无权限与业务不存在统一为一个无权限页面
    if (view.meta?.owner === MENU_BUSINESS) {
      view.components.permission = statusNonExistBusiness
    } else {
      view.components.permission = statusPermission
    }
  })
  return views
}

export const indexViews = injectStatusComponents(flatternViews([index]))

export const hostLandingViews = injectStatusComponents(flatternViews([hostLanding]))

export const businessViews = injectStatusComponents(flatternViews([
  customQuery,
  businessTopology,
  serviceTemplate,
  serviceCategory,
  serviceInstance,
  serviceSynchronous,
  customFields,
  setSync,
  setTemplate,
  hostApply
]))

// 业务集消费视图
export const businessSetViews = injectStatusComponents(flatternViews([
  businessSetTopology
]))

export const resourceViews = injectStatusComponents(flatternViews([
  business,
  businessSet,
  resource,
  generalModel,
  resourceManagement,
  cloudArea,
  cloudAccount,
  cloudResource,
  project
]))

export const modelViews = injectStatusComponents(flatternViews([
  model,
  modelAssociation,
  modelTopology,
  fieldTemplate
]))

export const analysisViews = injectStatusComponents(flatternViews([
  audit,
  operation
]))


export const platformManagementViews = injectStatusComponents(flatternViews([
  globalConfig
]))

export default {
  ...indexViews,
  ...hostLandingViews,
  ...businessSetViews,
  ...businessViews,
  ...resourceViews,
  ...modelViews,
  ...analysisViews,
  ...platformManagementViews
}
