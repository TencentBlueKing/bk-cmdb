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

import has from 'has'

export default class CachedPromise {
  constructor() {
    this.cache = {}
  }

  get(id) {
    if (typeof id === 'undefined') {
      return Object.keys(this.cache).map(requestId => this.cache[requestId].promise)
    }
    return has(this.cache, id) ? this.cache[id].promise : null
  }

  set(id, promise, config) {
    Object.assign(this.cache, { [id]: { promise, config } })
  }

  getGroupedIds(id) {
    const groupedIds = []
    Object.keys(this.cache).forEach((requestId) => {
      const isInclude = groupedIds.includes(requestId)
      const isMatch = this.cache[requestId].config.requestGroup.includes(id)
      if (!isInclude && isMatch) {
        groupedIds.push(requestId)
      }
    })
    return groupedIds
  }

  getDeleteIds(id) {
    const deleteIds = this.getGroupedIds(id)
    if (has(this.cache, id)) {
      deleteIds.push(id)
    }
    return deleteIds
  }

  delete(deleteIds) {
    let requestIds = []
    if (typeof deleteIds === 'undefined') {
      requestIds = Object.keys(this.cache)
    } else if (deleteIds instanceof Array) {
      deleteIds.forEach((id) => {
        requestIds = [...requestIds, ...this.getDeleteIds(id)]
      })
    } else {
      requestIds = this.getDeleteIds(deleteIds)
    }
    requestIds = [...new Set(requestIds)]
    requestIds.forEach((requestId) => {
      delete this.cache[requestId]
    })
    return Promise.resolve(deleteIds)
  }
}
