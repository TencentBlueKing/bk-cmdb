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
  MENU_BUSINESS_SERVICE,
  MENU_BUSINESS_SET_TEMPLATE,
  MENU_BUSINESS_SET_TEMPLATE_CREATE,
  MENU_BUSINESS_SET_TEMPLATE_DETAILS,
  MENU_BUSINESS_SET_TEMPLATE_EDIT,
  MENU_BUSINESS_SET_TEMPLATE_SYNC_HISTORY
} from '@/dictionary/menu-symbol'

export default [
  {
    name: MENU_BUSINESS_SET_TEMPLATE,
    path: 'set/template',
    component: () => import('./index.vue'),
    meta: new Meta({
      owner: MENU_BUSINESS,
      menu: {
        i18n: '集群模板',
        parent: MENU_BUSINESS_SERVICE
      }
    })
  },
  {
    name: MENU_BUSINESS_SET_TEMPLATE_CREATE,
    path: 'set/template/create',
    component: () => import('./create.vue'),
    meta: new Meta({
      owner: MENU_BUSINESS,
      menu: {
        i18n: '新建模板',
        relative: MENU_BUSINESS_SET_TEMPLATE
      }
    })
  },
  {
    name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
    path: 'set/template/details/:templateId',
    component: () => import('./details.vue'),
    meta: new Meta({
      owner: MENU_BUSINESS,
      menu: {
        i18n: '模板详情',
        relative: MENU_BUSINESS_SET_TEMPLATE
      }
    })
  },
  {
    name: MENU_BUSINESS_SET_TEMPLATE_EDIT,
    path: 'set/template/edit/:templateId',
    component: () => import('./edit.vue'),
    meta: new Meta({
      owner: MENU_BUSINESS,
      menu: {
        i18n: '新建模板',
        relative: MENU_BUSINESS_SET_TEMPLATE
      }
    })
  },
  {
    name: MENU_BUSINESS_SET_TEMPLATE_SYNC_HISTORY,
    path: 'set/instance/history/:templateId?',
    component: () => import('./sync-history.vue'),
    meta: new Meta({
      owner: MENU_BUSINESS,
      menu: {
        i18n: '同步历史',
        relative: MENU_BUSINESS_SET_TEMPLATE
      }
    })
  }
]
