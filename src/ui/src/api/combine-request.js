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

/**
 * 合并请求通用工具类，适用于组件拆分但希望合并数据请求场景，同时可用于控制异步方法分片执行。
 * 理论上，此工具实现的是一个延时执行模型，即执行目标可以是异步请求或者任何普通函数，编排后的执行可由使用方自行处理。
 * usage:
 * await CombineRequest.setup(Symbol(), data => Promise.resolve(data)).add(1)
 */

export default class CombineRequest {
  static setup(id, callback, options = {}) {
    return this.getInstance(id, callback, options)
  }

  static getInstance(id, callback, options) {
    const instances = CombineRequest.instances || {}
    if (!instances[id]) {
      instances[id] = new CombineRequest(id, callback, options)
    }
    CombineRequest.instances = instances
    return instances[id]
  }

  /**
   * 构建函数
   * @param {String | Symbol} id 实例ID
   * @param {Function} callback 执行的回调函数
   * @param {Object} options 配置项
   * @param {Number} options.segment 数据分片大小
   * @param {Number} options.concurrency 回调并发数
   */
  constructor(id, callback, { segment, concurrency }) {
    this.id = id
    this.timer = null
    this.data = []
    this.callback = callback
    this.promise = null,
    this.segment = segment
    this.concurrency = concurrency
  }

  add(payload) {
    this.data.push(payload)
    if (!this.timer) {
      this.promise = new Promise((resolve, reject) => {
        this.timer = setTimeout(() => {
          try {
            const result = this.run()
            resolve(result)
          } catch (error) {
            reject(error)
          } finally {
            this.reset()
          }
        }, 0)
      })
    }
    return this.promise
  }

  run() {
    if (this.callback) {
      if (this.segment && this.segment > 0) {
        if (this.concurrency && this.concurrency > 0) {
          return this.sliceRun()
        }
        return this.splitRun()
      }
      return this.callback(this.data)
    }
  }

  /**
   * 分片参数并发起调用，仅分片参数不分片调用，调用会一次性全部触发
   * @returns 所有调用的结果列表
   */
  splitRun() {
    const res = []
    const data = this.flatten(this.data)
    this.slice(this.segment, data).forEach(params => res.push(this.callback(params)))
    return res
  }

  /**
   * 同时分片参数与调用，调用会根据并发数配置进行分片，返回分片后的迭代器
   * @returns 已分片的迭代器，可自行迭代控制调用触发
   */
  sliceRun() {
    const cb = this.callback

    const data = this.flatten(this.data)
    // 分片参数为segment个一组
    const params = this.slice(this.segment, data)

    // 基于分片后的参数，分组每一次的调用形成调用队列
    const callQueue = this.slice(this.concurrency, params)

    const genCall = function* (queue, cb) {
      for (let i = 0; i < queue.length; i++) {
        // concurrency个为一组
        yield Promise.allSettled(queue[i].map(params => cb(params)))
      }
    }

    return genCall(callQueue, cb)
  }

  reset() {
    clearTimeout(this.timer)
    this.data = []
    this.timer = null
    this.promise = null
  }

  slice(size, data) {
    const segments = []
    const max = Math.ceil(data.length / size)
    for (let index = 1; index <= max; index++) {
      const segment = data.slice((index - 1) * size, size * index)
      segments.push(segment)
    }
    return segments
  }

  flatten(data) {
    const arr = data ?? []
    const fn = arr => arr.reduce((p, c) => p.concat(Array.isArray(c) ? fn(c) : c), [])
    return fn(arr)
  }
}
