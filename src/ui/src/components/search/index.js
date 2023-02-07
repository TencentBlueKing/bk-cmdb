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

import bool from './bool'
import date from './date'
import enumComponent from './enum'
import float from './float'
import foreignkey from './foreignkey'
import int from './int'
import list from './list'
import longchar from './longchar'
import objuser from './objuser'
import organization from './organization'
import singlechar from './singlechar'
import table from './table'
import time from './time'
import timezone from './timezone'
import serviceTemplate from './service-template'
import module from './module'
import set from './set'
import biz from './biz'
import array from './array.vue'
import object from './object.vue'
import map from './map.vue'

export default {
  install(Vue) {
    const components = [
      bool,
      date,
      enumComponent,
      float,
      foreignkey,
      int,
      list,
      longchar,
      objuser,
      organization,
      singlechar,
      table,
      time,
      timezone,
      serviceTemplate,
      module,
      set,
      biz,
      array,
      object,
      map
    ]
    components.forEach((component) => {
      Vue.component(component.name, component)
    })
  }
}
