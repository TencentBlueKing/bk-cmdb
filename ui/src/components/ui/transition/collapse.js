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

/* eslint-disable no-restricted-syntax */
import has from 'has'
const Transition = {
  beforeEnter(el) {
    el.classList.add('collapse-transition')
    if (!el.dataset) {
      el.dataset = {}
    }

    el.dataset.oldPaddingTop = el.style.paddingTop
    el.dataset.oldPaddingBottom = el.style.paddingBottom

    el.style.height = '0'
    el.style.paddingTop = 0
    el.style.paddingBottom = 0
  },

  enter(el) {
    el.dataset.oldOverflow = el.style.overflow
    el.style.overflow = 'hidden'
    if (el.scrollHeight !== 0) {
      el.style.height = `${el.scrollHeight}px`
      el.style.paddingTop = el.dataset.oldPaddingTop
      el.style.paddingBottom = el.dataset.oldPaddingBottom
    } else {
      el.style.height = ''
      el.style.paddingTop = el.dataset.oldPaddingTop
      el.style.paddingBottom = el.dataset.oldPaddingBottom
    }
  },

  afterEnter(el) {
    el.classList.remove('collapse-transition')
    el.style.height = ''
    el.style.overflow = el.dataset.oldOverflow
  },

  beforeLeave(el) {
    if (!el.dataset) el.dataset = {}
    el.dataset.oldPaddingTop = el.style.paddingTop
    el.dataset.oldPaddingBottom = el.style.paddingBottom
    el.dataset.oldOverflow = el.style.overflow

    el.style.height = `${el.scrollHeight}px`
    el.style.overflow = 'hidden'
  },

  leave(el) {
    if (el.scrollHeight !== 0) {
      el.classList.add('collapse-transition')
      el.style.height = 0
      el.style.paddingTop = 0
      el.style.paddingBottom = 0
    }
  },

  afterLeave(el) {
    el.classList.remove('collapse-transition')
    el.style.height = ''
    el.style.overflow = el.dataset.oldOverflow
    el.style.paddingTop = el.dataset.oldPaddingTop
    el.style.paddingBottom = el.dataset.oldPaddingBottom
  }
}

const toCamelCase = function (str) {
  return str.replace(/-([a-z])/g, g => g[1].toUpperCase())
}

export default {
  name: 'cmdb-collapse-transition',
  functional: true,
  render(h, context) {
    const events = context.data.on || {}
    const camelCaseEvents = {}
    const transitionEvents = {}
    for (const event in events) {
      camelCaseEvents[toCamelCase(event)] = events[event]
    }
    for (const event in Transition) {
      if (has(camelCaseEvents, event)) {
        transitionEvents[event] = (el) => {
          Transition[event](el)
          camelCaseEvents[event]()
        }
      } else {
        transitionEvents[event] = Transition[event]
      }
    }
    const data = {
      on: transitionEvents
    }
    return h('transition', data, context.children)
  }
}
