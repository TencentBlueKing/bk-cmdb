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
  MENU_RESOURCE_INSTANCE,
  MENU_RESOURCE_MANAGEMENT,
  MENU_RESOURCE_INSTANCE_DETAILS
} from '@/dictionary/menu-symbol'

export default [{
  name: MENU_RESOURCE_INSTANCE,
  path: 'instance/:objId',
  component: () => import('./index.vue'),
  meta: new Meta({
    menu: {
      relative: MENU_RESOURCE_MANAGEMENT
    },
    checkAvailable: (to, from, app) => {
      const modelId = to.params.objId
      const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
      return model && !model.bk_ispaused
    },
    layout: {}
  }),
  children: [{
    name: MENU_RESOURCE_INSTANCE_DETAILS,
    path: ':instId',
    component: () => import('./details.vue'),
    meta: new Meta({
      menu: {
        relative: MENU_RESOURCE_MANAGEMENT
      },
      checkAvailable: (to, from, app) => {
        const modelId = to.params.objId
        const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
        return model && !model.bk_ispaused
      },
      layout: {}
    })
  }]
}, {
  name: 'instanceHistory',
  path: 'history/instance/:objId',
  component: () => import('@/views/history/index.vue'),
  meta: new Meta({
    menu: {
      relative: MENU_RESOURCE_MANAGEMENT
    },
    layout: {},
    checkAvailable: (to, from, app) => {
      const modelId = to.params.objId
      const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
      return model && !model.bk_ispaused
    }
  })
}]
