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

const scrollListeners = []
const resizeListeners = []
export function addMainScrollListener(fn) {
  if (scrollListeners.includes(fn)) return
  if (typeof fn === 'function') {
    scrollListeners.push(fn)
  }
}

export function removeMainScrollListener(fn) {
  const index = scrollListeners.indexOf(fn)
  if (index !== -1) {
    scrollListeners.splice(index, 1)
  }
}

export function execMainScrollListener(event) {
  scrollListeners.forEach((fn) => {
    fn(event)
  })
}

export function addMainResizeListener(fn) {
  if (resizeListeners.includes(fn)) return
  if (typeof fn === 'function') {
    resizeListeners.push(fn)
  }
}

export function removeMainResizeListener(fn) {
  const index = resizeListeners.indexOf(fn)
  if (index !== -1) {
    resizeListeners.splice(index, 1)
  }
}

export function execMainResizeListener(event) {
  resizeListeners.forEach((fn) => {
    fn(event)
  })
}
