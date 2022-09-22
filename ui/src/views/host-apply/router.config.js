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

import Meta from '@/router/meta'
import {
  MENU_BUSINESS,
  MENU_BUSINESS_HOST,
  MENU_BUSINESS_HOST_APPLY,
  MENU_BUSINESS_HOST_APPLY_EDIT,
  MENU_BUSINESS_HOST_APPLY_CONFIRM,
  MENU_BUSINESS_HOST_APPLY_RUN,
  MENU_BUSINESS_HOST_APPLY_CONFLICT,
  MENU_BUSINESS_HOST_APPLY_FAILED
} from '@/dictionary/menu-symbol'
import { TASK_STATUS, getTask } from './task-helper.js'
import { CONFIG_MODE } from '@/service/service-template/index.js'

// before拦截钩子，当本地存在未完成的应用任务则直接跳转到应用页面
const beforeInterceptor = async (to, from, app) => {
  const { $store, $router } = app
  const bizId = $store.getters['objectBiz/bizId']

  // 规范化配置模式路由，不合法的模式强制跳转至首页
  const allowedModes = Object.values(CONFIG_MODE)
  if (!allowedModes.includes(to?.params?.mode)) {
    $router.push({
      name: MENU_BUSINESS_HOST_APPLY,
      params: {
        bizId,
        mode: allowedModes[0]
      }
    })
    return Promise.resolve(false)
  }

  const localTaskId = getTask(`${$store.getters.userName}${bizId}`)
  if (localTaskId) {
    const result = await $store.dispatch('hostApply/getApplyTaskStatus', {
      params: {
        bk_biz_id: bizId,
        task_ids: [localTaskId]
      },
      config: {
        globalError: false
      }
    })
    const taskStatus = result?.task_info?.find(({ task_id: id }) => id === localTaskId)?.status
    if ([TASK_STATUS.NEW, TASK_STATUS.WAITING, TASK_STATUS.EXECUTING].includes(taskStatus)) {
      $router.push({
        name: MENU_BUSINESS_HOST_APPLY_RUN,
        params: { bizId }
      })
      return Promise.resolve(false)
    }
  }
  return Promise.resolve(true)
}

export default [{
  name: MENU_BUSINESS_HOST_APPLY,
  path: 'host-apply/:mode?',
  component: () => import('./index.vue'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST
    },
    before: beforeInterceptor
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_CONFIRM,
  path: 'host-apply/:mode/confirm',
  component: () => import('./property-confirm'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    },
    before: beforeInterceptor
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_EDIT,
  path: 'host-apply/:mode/edit',
  component: () => import('./edit'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    },
    before: beforeInterceptor
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_RUN,
  path: 'host-apply/:mode/run',
  component: () => import('./property-run'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    },
    layout: {
      breadcrumbs: {
        back: true
      }
    }
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_CONFLICT,
  path: 'host-apply/:mode/conflict',
  component: () => import('./conflict-list'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    }
  }),
  before: beforeInterceptor
}, {
  name: MENU_BUSINESS_HOST_APPLY_FAILED,
  path: 'host-apply/:mode/failed',
  component: () => import('./failed-list'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    }
  }),
  before: beforeInterceptor
}]
