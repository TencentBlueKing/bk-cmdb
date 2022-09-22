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

import { camelize } from '@/utils/util.js'
import testIds from '@/dictionary/test-id.js'

const isString = val => Object.prototype.toString.call(val) === '[object String]'

const testIdKeys = Object.keys(testIds)
const tagName = name => camelize(name.toLowerCase())
const routeName = name => camelize(name.replace(/^menu_/, ''), '_')

const tagNames = ['nav', 'header', 'button', 'form', 'ul', 'li', 'div', 'section']
const availableTagName = el => el.tagName && tagNames.includes(el.tagName.toLowerCase())

function setDataTestId(el, binding, vnode) {
  const { value, modifiers } = binding
  const { context, componentInstance, componentOptions } = vnode
  const moduleId = testIdKeys.find(key => modifiers[key]) || routeName(context.$route.name)
  const moduleSetting = testIds[moduleId]

  if (value && !isString(value)) {
    console.warn('only accept string values!')
    return
  }

  // 定义文件中的优先级最高，其次是在常用标签上写了value的，最后是自定义组件的形式
  let blockId = 'unknown'
  if (moduleSetting?.[value]) {
    blockId = moduleSetting[value]
  } else if (availableTagName(el) && value) {
    blockId = `${tagName(el.tagName)}_${camelize(value)}`
  } else if (componentInstance) {
    blockId = `comp_${tagName(componentOptions.tag)}`
  }

  el.setAttribute('data-test-id', `${moduleId}_${blockId}`)
}

const directive = {
  bind(el, binding, vnode) {
    setDataTestId(el, binding, vnode)
  },
  componentUpdated(el, binding, vnode) {
    if (binding.modifiers.dynamic) {
      setDataTestId(el, binding, vnode)
    }
  }
}

export default {
  install: Vue => Vue.directive('test-id', directive)
}
