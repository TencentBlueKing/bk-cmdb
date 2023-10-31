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

const fs = require('node:fs')

module.exports = function (source) {
  const ignoreFiles = [
    'net-discovery.js',
    'net-data-collection.js',
    'net-collect-device.js',
    'net-collect-property.js',
    'user-privilege.js',
    'object-relation.js',
    'cloud-discover.js',
    'host-search-history.js',
    'proc-config.js',
    'components/import/import.vue'
  ]

  if (ignoreFiles.some(x => this.resourcePath.endsWith(x))) {
    return source
  }

  const reg1 = new RegExp('(http|\\$http)\\.(post|get|delete|put|download)\\(([\\w\'`/${}.\\[\\]?=()\\s]+?)(?=[,)])', 'gm')
  const reg2 = /(\w+\s*[:=]\s*|return\s*)(`\$\{window.API_HOST\}[\w`/${}.]+)/gm
  const matches1 = source.matchAll(reg1)
  const matches2 = source.matchAll(reg2)
  for (const match of matches1) {
    const m = match?.[3]?.trim?.()?.split('?')?.[0]
    if (!m) {
      continue
    }

    if (['`', '\''].includes(m.substring(0, 1))) {
      fs.appendFileSync('apipaths.txt', `${m}\n`)
    } else {
      fs.appendFileSync('apivars.txt', `${this.resourcePath}: ${m}\n`)
    }
  }

  for (const match of matches2) {
    const m = match?.[2]?.trim?.()?.split('?')?.[0]
    if (!m) {
      continue
    }

    fs.appendFileSync('apipaths.txt', `${m}\n`)
  }

  return source
}
