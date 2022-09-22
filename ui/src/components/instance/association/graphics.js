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

import memoize from 'lodash.memoize'
import cytoscape from 'cytoscape'
import popper from 'cytoscape-popper'
import dagre from 'cytoscape-dagre'
import { layout, style } from './graphics-config'
import { generateObjIcon } from '@/utils/util'
cytoscape.use(popper)
cytoscape.use(dagre)
const makeSVG = memoize((element) => {
  if (!element.isNode()) return
  return new Promise((resolve) => {
    const data = element.data()
    const image = new Image()
    image.onload = () => {
      const unselected = `data:image/svg+xml;charset=utf-8,${encodeURIComponent(generateObjIcon(image, {
        name: data.name,
        iconColor: '#798aad',
        backgroundColor: '#fff'
      }))}`
      const selected = `data:image/svg+xml;charset=utf-8,${encodeURIComponent(generateObjIcon(image, {
        name: data.name,
        iconColor: '#fff',
        backgroundColor: '#3a84ff'
      }))}`
      const hover = `data:image/svg+xml;charset=utf-8,${encodeURIComponent(generateObjIcon(image, {
        name: data.name,
        iconColor: '#3a84ff',
        backgroundColor: '#fff'
      }))}`
      resolve({ selected, unselected, hover })
    }
    image.onerror = () => resolve({ selected: 'none', unselected: 'none', hover: 'none' })
    image.src = `${window.location.origin}/static/svg/${data.icon.substr(5)}.svg`
  })
}, element => element.data().icon)
export default class Graphics {
  constructor(container) {
    this.current = null
    this.position = null
    this.cytoscape = cytoscape({
      container,
      layout,
      style,
      minZoom: 0.5,
      maxZoom: 5,
      wheelSensitivity: 0.5,
      pixelRatio: 2,
      elements: []
    })
    this.on('layoutstop', event => this.handleLayoutStop(event))
    this.on('resize', event => this.handleResize(event))
    this.on('mouseover', 'edge', event => this.handleEdgeMouseover(event))
    this.on('mouseout', 'edge', event => this.handleEdgeMouseout(event))
    this.on('mouseover', 'node', event => this.handleNodeMouseover(event))
    this.on('mouseout', 'node', event => this.handleNodeMouseout(event))
    this.on('click', 'node', event => this.setCurrent(event.target))
    this.on('click', event => this.handleCoreClick(event))
  }

  on(type, target, handler) {
    this.cytoscape.on(type, target, handler)
  }
  handleLayoutStop() {
    this.cytoscape.zoom(1)
    if (!this.position) return
    // 重新布局后，计算点击前的节点位置与当前位置的差值，平移画布，使得点击的点位置相对不变
    const currentPosition = this.current.renderedPosition()
    const delta = {
      x: this.position.x - currentPosition.x,
      y: this.position.y - currentPosition.y
    }
    this.cytoscape.panBy(delta)
  }
  handleResize() {
    this.cytoscape.fit()
    this.cytoscape.zoom(1)
    this.cytoscape.center()
  }
  /**
   * 悬浮edge时，给edge添加hover类
   * 聚焦模式下，仅高亮的edge有效
   */
  handleEdgeMouseover(event) {
    const edge = event.target
    if (this.current && !edge.hasClass('highlight')) return
    edge.addClass('hover')
  }
  /**
   * edge失焦时，移除edge的hover类
   */
  handleEdgeMouseout(event) {
    event.target.removeClass('hover')
  }
  /**
   * 给node添加hover类
   * 聚焦模式下，仅给高亮的node添加hover类
   */
  handleNodeMouseover(event) {
    const node = event.target
    if (this.current && !node.hasClass('highlight')) return
    node.addClass('hover')
  }
  /**
   * node失焦时，移除hover类
   */
  handleNodeMouseout(event) {
    event.target.removeClass('hover')
  }
  /**
   * 点击画布时，取消聚焦模式
   */
  handleCoreClick(event) {
    if (event.target !== this.cytoscape || !this.current) return
    this.cytoscape.elements().removeClass(['current', 'highlight', 'weaken'])
    this.setCurrent(null)
  }

  /**
   * 设置当前选择的node，设置了node时即进入聚焦模式
   */
  setCurrent(node) {
    // 移除上一个node的状态
    if (this.current) {
      this.current.removeClass('current')
    }
    // 设置新的当前节点
    this.current = node
    this.position = node && ({ ...node.renderedPosition() })
    if (this.current) {
      this.current.addClass(['current', 'highlight']).removeClass('weaken')
    }
  }

  /**
   * 设置聚焦模式下的高亮与淡化效果
   */
  setHightlight(id) {
    const node = this.cytoscape.getElementById(id)
    const elements = this.cytoscape.elements()
    const neighborhood = node.neighborhood()
    elements.forEach((element) => {
      if (element.same(node)) return
      if (neighborhood.has(element)) {
        element.addClass('highlight')
        element.removeClass('weaken')
        return
      }
      element.addClass('weaken')
      element.removeClass('highlight')
    })
  }

  /**
   * 添加元素，过滤已存在的元素，并加载背景
   */
  add(item) {
    const list = Array.isArray(item) ? item : [item]
    const newList = list.filter(({ data }) => {
      const existElement = this.cytoscape.getElementById(data.id)
      return !existElement.length
    })
    const elements = this.cytoscape.add(newList)
    elements.forEach(async (element) => {
      const background = await makeSVG(element)
      element.data('background', background)
    })
    this.relayout()
    return elements
  }

  /**
   * 重新布局
   */
  relayout() {
    this.cytoscape.layout(layout).run()
  }
}
