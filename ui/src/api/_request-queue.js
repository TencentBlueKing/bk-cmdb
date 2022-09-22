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

export default class RequestQueue {
  constructor() {
    this.queue = []
  }

  get(id) {
    if (typeof id === 'undefined') return this.queue
    return this.queue.find(request => request.requestId === id || request.requestGroup.includes(id))
  }

  set(newRequest) {
    if (!this.queue.some(request => request.requestId === newRequest.requestId)) {
      this.queue.push(newRequest)
    }
  }

  delete(id, symbol) {
    let target
    if (symbol) {
      target = this.queue.find(request => request.requestSymbol === symbol)
    } else {
      target = this.queue.find(request => request.requestId === id)
    }
    if (target) {
      const index = this.queue.indexOf(target)
      this.queue.splice(index, 1)
    }
  }

  cancel(requestIds) {
    let cancelQueue = []
    if (typeof requestIds === 'undefined') {
      cancelQueue = [...this.queue]
    } else if (requestIds instanceof Array) {
      requestIds.forEach((requestId) => {
        const cancelRequest = this.get(requestId)
        if (cancelRequest) {
          cancelQueue.push(cancelRequest)
        }
      })
    } else {
      const cancelRequest = this.get(requestIds)
      if (cancelRequest) {
        cancelQueue.push(cancelRequest)
      }
    }
    try {
      cancelQueue.forEach((request) => {
        request.cancelExcutor(request)
        this.delete(request.requestId)
      })
      return Promise.resolve(requestIds)
    } catch (error) {
      return Promise.reject(error)
    }
  }
}
