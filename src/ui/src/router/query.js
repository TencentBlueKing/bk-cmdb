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
import router from './index'
import throttle from 'lodash.throttle'
import deepEqual from 'deep-equal'
import { redirect } from './actions'

// Vue.watch的options
function createWatchOptions(key, options) {
  const watchOptions = {
    immediate: false,
    deep: false
  }
  // eslint-disable-next-line no-restricted-syntax
  for (const key in watchOptions) {
    if (has(options, key)) {
      watchOptions[key] = options[key]
    }
  }
  if (key === '*') {
    watchOptions.deep = true
  }
  return watchOptions
}

// 始终watch router.query, 根据watch的key再做变更比对，除immediate外，无变更时不触发注册的handler
function createCallback(keys, handler, options = {}) {
  let immediateCalled = false
  const callback = (values, oldValues = {}) => {
    let execValue; let execOldValue
    if (Array.isArray(keys)) {
      execValue = {}
      execOldValue = {}
      keys.forEach((key) => {
        execValue[key] = values[key]
        execOldValue[key] = oldValues[key]
      })
    } else if (keys === '*') {
      execValue = { ...values }
      execOldValue = { ...oldValues }
      if (has(options, 'ignore')) {
        const ignoreKeys = Array.isArray(options.ignore) ? options.ignore : [options.ignore]
        ignoreKeys.forEach((key) => {
          delete execValue[key]
          delete execOldValue[key]
        })
      }
    } else {
      execValue = values[keys]
      execOldValue = oldValues[keys]
    }
    if (options.immediate && !immediateCalled) {
      immediateCalled = true
      handler(execValue, execOldValue)
    } else {
      const hasChange = !deepEqual(execValue, execOldValue)
      hasChange && handler(execValue, execOldValue)
    }
  }

  if (has(options, 'throttle')) {
    const interval = typeof options.throttle === 'number' ? options.throttle : 100
    return throttle(callback, interval, { leading: false, trailing: true })
  }

  return callback
}

function isEmpty(value) {
  return value === '' || value === undefined || value === null
}

class RouterQuery {
  constructor() {
    this.router = router
  }
  get app() {
    return router.app
  }

  get route() {
    return this.app.$route
  }

  get(key, defaultValue) {
    if (has(this.route.query, key)) {
      return this.route.query[key]
    }
    if (arguments.length === 2) {
      return defaultValue
    }
  }

  getAll() {
    return this.route.query
  }

  set(key, value) {
    const query = { ...this.route.query }
    if (typeof key === 'object') {
      Object.assign(query, key)
    } else {
      query[key] = value
    }
    Object.keys(query).forEach((queryKey) => {
      if (isEmpty(query[queryKey])) {
        delete query[queryKey]
      }
    })
    redirect({
      ...this.route,
      query
    })
  }

  setAll(value) {
    redirect({
      ...this.route,
      query: {
        ...value
      }
    })
  }

  delete(key) {
    const query = {
      ...this.route.query
    }
    delete query[key]
    redirect({
      ...this.route,
      query
    })
  }

  refresh() {
    this.set('_t', Date.now())
  }

  clear() {
    redirect({
      ...this.route,
      query: {}
    })
  }

  watch(key, handler, options = {}) {
    const watchOptions = createWatchOptions(key, options)
    const callback = createCallback(key, handler, options)
    const expression = () => this.route.query
    return this.app.$watch(expression, callback, watchOptions)
  }
}

export default new RouterQuery()
