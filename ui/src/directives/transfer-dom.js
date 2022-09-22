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
 * @file transfer-dom
 * https://github.com/airyland/vux/blob/v2/src/directives/transfer-dom/index.js
 * https://github.com/calebroseland/vue-dom-portal
 */

/**
 * Get target DOM Node
 * @param {(Node|string|Boolean)} [node=document.body] DOM Node, CSS selector, or Boolean
 * @return {Node} The target that the el will be appended to
 */

/* eslint-disable no-underscore-dangle */
function getTarget(node) {
  if (node === void 0) {
    node = document.body
  }
  if (node === true) {
    return document.body
  }
  return node instanceof window.Node ? node : document.querySelector(node)
}

const directive = {
  inserted(el, { value }) {
    el.className = el.className ? `${el.className} v-transfer-dom` : 'v-transfer-dom'
    const { parentNode } = el
    if (!parentNode) {
      return
    }
    const home = document.createComment('')
    let hasMovedOut = false

    if (value !== false) {
      parentNode.replaceChild(home, el) // moving out, el is no longer in the document
      getTarget(value).appendChild(el) // moving into new place
      hasMovedOut = true
    }
    if (!el.__transferDomData) {
      el.__transferDomData = {
        parentNode,
        home,
        target: getTarget(value),
        hasMovedOut
      }
    }
  },
  componentUpdated(el, { value }) {
    // need to make sure children are done updating (vs. `update`)
    const ref$1 = el.__transferDomData
    if (!ref$1) {
      return
    }
    // homes.get(el)
    const { parentNode } = ref$1
    const { home } = ref$1
    const { hasMovedOut } = ref$1 // recall where home is

    if (!hasMovedOut && value) {
      // remove from document and leave placeholder
      parentNode.replaceChild(home, el)
      // append to target
      getTarget(value).appendChild(el)
      el.__transferDomData = Object.assign(
        {},
        el.__transferDomData,
        {
          hasMovedOut: true,
          target: getTarget(value)
        }
      )
    } else if (hasMovedOut && value === false) {
      // previously moved, coming back home
      parentNode.replaceChild(el, home)
      el.__transferDomData = Object.assign(
        {},
        el.__transferDomData,
        {
          hasMovedOut: false,
          target: getTarget(value)
        }
      )
    } else if (value) {
      // already moved, going somewhere else
      getTarget(value).appendChild(el)
    }
  },
  unbind(el) {
    if (el.nodeType !== 1) {
      return false
    }
    el.className = el.className.replace('v-transfer-dom', '')
    const ref$1 = el.__transferDomData
    if (!ref$1) {
      return
    }
    if (el.__transferDomData.hasMovedOut === true) {
      el.__transferDomData.parentNode && el.__transferDomData.parentNode.appendChild(el)
    }
    el.__transferDomData = null
  }
}

export default directive
