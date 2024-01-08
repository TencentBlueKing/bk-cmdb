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
import { OPERATION } from '@/dictionary/iam-auth'
import {
  MENU_MODEL_FIELD_TEMPLATE,
  MENU_MODEL_FIELD_TEMPLATE_CREATE,
  MENU_MODEL_FIELD_TEMPLATE_CREATE_BASIC,
  MENU_MODEL_FIELD_TEMPLATE_CREATE_FIELD_SETTINGS,
  MENU_MODEL_FIELD_TEMPLATE_EDIT,
  MENU_MODEL_FIELD_TEMPLATE_EDIT_BASIC,
  MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
  MENU_MODEL_FIELD_TEMPLATE_EDIT_BINDING,
  MENU_MODEL_FIELD_TEMPLATE_BIND,
  MENU_MODEL_FIELD_TEMPLATE_SYNC_MODEL
} from '@/dictionary/menu-symbol'

export default [
  {
    name: MENU_MODEL_FIELD_TEMPLATE,
    path: 'field-template',
    component: () => import('./index.vue'),
    meta: new Meta({
      menu: {
        i18n: '字段组合模板'
      },
      layout: {
      }
    })
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_CREATE,
    path: 'field-template/create',
    redirect: { name: MENU_MODEL_FIELD_TEMPLATE_CREATE_BASIC }
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_CREATE_BASIC,
    path: 'field-template/create/basic',
    component: () => import('./create-basic.vue'),
    meta: new Meta({
      menu: {
        i18n: '新建字段组合模板',
        relative: MENU_MODEL_FIELD_TEMPLATE
      }
    })
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_CREATE_FIELD_SETTINGS,
    path: 'field-template/create/field-settings',
    component: () => import('./create-field-settings.vue'),
    meta: new Meta({
      menu: {
        i18n: '新建字段组合模板',
        relative: MENU_MODEL_FIELD_TEMPLATE
      }
    })
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_EDIT,
    path: 'field-template/edit',
    redirect: { name: MENU_MODEL_FIELD_TEMPLATE_EDIT_BASIC }
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_EDIT_BASIC,
    path: 'field-template/edit/:id/basic',
    component: () => import('./edit-basic.vue'),
    meta: new Meta({
      menu: {
        i18n: '编辑字段组合模板',
        relative: MENU_MODEL_FIELD_TEMPLATE
      }
    })
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
    path: 'field-template/edit/:id/field-settings',
    component: () => import('./edit-field-settings.vue'),
    meta: new Meta({
      menu: {
        i18n: '编辑字段组合模板',
        relative: MENU_MODEL_FIELD_TEMPLATE
      },
      auth: {
        view: (to) => {
          const { id } = to.params
          return ({ type: OPERATION.R_FIELD_TEMPLATE, relation: [Number(id)] })
        }
      }
    })
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_EDIT_BINDING,
    path: 'field-template/edit/:id/binding',
    component: () => import('./edit-binding.vue'),
    meta: new Meta({
      menu: {
        i18n: '编辑字段组合模板',
        relative: MENU_MODEL_FIELD_TEMPLATE
      },
      auth: {
        view: (to) => {
          const { id } = to.params
          return ({ type: OPERATION.R_FIELD_TEMPLATE, relation: [Number(id)] })
        }
      }
    })
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_BIND,
    path: 'field-template/bind/:id',
    component: () => import('./bind.vue'),
    meta: new Meta({
      menu: {
        i18n: '绑定模型',
        relative: MENU_MODEL_FIELD_TEMPLATE
      },
      auth: {
        view: (to) => {
          const { id } = to.params
          return ({ type: OPERATION.R_FIELD_TEMPLATE, relation: [Number(id)] })
        }
      }
    })
  },
  {
    name: MENU_MODEL_FIELD_TEMPLATE_SYNC_MODEL,
    path: 'field-template/sync/:id/model/:modelId',
    component: () => import('./sync.vue'),
    meta: new Meta({
      menu: {
        i18n: '同步模型',
        relative: MENU_MODEL_FIELD_TEMPLATE
      }
    })
  }
]
