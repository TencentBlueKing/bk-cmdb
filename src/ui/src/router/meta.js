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

export default class Meta {
  constructor(data = {}) {
    this.owner = ''
    this.title = ''
    this.available = true
    Object.keys(data).forEach((key) => {
      this[key] = data[key]
    })

    this.menu = Object.assign({
      i18n: '',
      parent: null,
      relative: null
    }, data.menu)

    this.auth = Object.assign({
      view: null,
      operation: null,
      permission: null
    }, data.auth)

    this.layout = Object.assign({
      breadcrumbs: true,
      previous: null
    }, data.layout)

    this.view = 'default'
  }
}
