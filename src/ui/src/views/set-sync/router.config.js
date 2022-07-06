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
import { MENU_BUSINESS, MENU_BUSINESS_SET_TEMPLATE } from '@/dictionary/menu-symbol'
export default [{
  name: 'setSync',
  path: 'set/sync/:setTemplateId',
  component: () => import('./sync-index.vue'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '同步集群模板',
      relative: MENU_BUSINESS_SET_TEMPLATE
    }
  })
}]
