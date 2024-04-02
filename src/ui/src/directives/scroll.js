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
 * @directive 目标元素滚动到可视化区域
 * @scrollEle 滚动的元素
 * @param {Element} targetClass 目标元素的className
 * @param {String} direction 滚动方向 vertical竖向 / horizontal横向
 * @param {Number} distance 距离
 * @param {Boolean} scrollInViewport 目标元素已在可视区域了是否还需要滚动
 */
// 默认参数
const defaultParams = {
  direction: 'vertical',
  distance: 10,
  scrollInViewport: false
}

// 当前目标元素是否在可视区域内
const isElementInViewport = (rect, scrollRect) => {
  const { top, left, bottom, right } = rect
  const { top: scrollTop, left: scrollLeft, bottom: scrollBottom, right: scrollRight } = scrollRect
  return (
    top >= scrollTop
        && left >= scrollLeft
        && bottom <= scrollBottom
        && right <= scrollRight
  )
}
const toScrollViewport = (scrollEle, rect, scrollRect, direction, distance) => {
  const { top, left } = rect
  if (direction === 'vertical') {
    scrollEle.scrollTop = top -  distance - scrollRect.top
  } else {
    scrollEle.scrollLeft = left -  distance - scrollRect.left
  }
}
const init = (scrollEle, value) => {
  setTimeout(() => {
    const params = Object.assign(defaultParams, value)
    const { targetClass, direction, distance, scrollInViewport } = params
    const target = scrollEle.querySelector(`.${targetClass}`)
    if (!target) return

    const rect = target.getBoundingClientRect()
    const scrollRect = scrollEle.getBoundingClientRect()
    const isInViewport = isElementInViewport(rect, scrollRect)
    if (!isInViewport || (isInViewport && scrollInViewport)) {
      toScrollViewport(scrollEle, rect, scrollRect, direction, distance)
    }
  }, 0)
}
export const scroll = {
  update: (scrollEle, { value }) => {
    init(scrollEle, value)
  },
  inserted: (scrollEle, { value }) => {
    init(scrollEle, value)
  }
}
