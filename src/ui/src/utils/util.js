/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { t } from '@/i18n'
import { BUILTIN_PASTE_SPLIT_FIELDS } from '@/dictionary/model-constants.js'

const hex2grb = (hex) => {
  const rgb = []
  hex = hex.substr(1)
  if (hex.length === 3) {
    hex = hex.replace(/(.)/g, '$1$1')
  }
  hex.replace(/../g, (color) => {
    rgb.push(parseInt(color, 0x10))
  })
  return rgb
}
const getFileExtension = fileName => fileName.substr((~-fileName.lastIndexOf('.') >>> 0) + 2)

const canvas = document.createElement('canvas')

const getBase64Image = (image, color) => {
  const ctx = canvas.getContext('2d')
  canvas.width = image.width
  canvas.height = image.height
  ctx.clearRect(0, 0, canvas.width, canvas.height)
  ctx.drawImage(image, 0, 0, image.width, image.height)
  const imageData = ctx.getImageData(0, 0, image.width, image.height)
  const rgbColor = hex2grb(color)
  for (let i = 0; i < imageData.data.length; i += 4) {
    const [r, g, b] = rgbColor
    imageData.data[i] = r
    imageData.data[i + 1] = g
    imageData.data[i + 2] = b
  }
  ctx.putImageData(imageData, 0, 0)
  return canvas.toDataURL(`image/${getFileExtension(image.src)}`)
}

const HIDDEN_STYLE = `
  height:0 !important;
  visibility:hidden !important;

  position:absolute !important;
  z-index:-1000 !important;
  top:0 !important;
  right:0 !important;
`

const CONTEXT_STYLE = [
  'letter-spacing',
  'line-height',
  'padding-top',
  'padding-bottom',
  'font-family',
  'font-weight',
  'font-size',
  'text-rendering',
  'text-transform',
  'width',
  'text-indent',
  'padding-left',
  'padding-right',
  'border-width',
  'box-sizing',
]

export function svgToImageUrl(image, options) {
  const base64Image = getBase64Image(image, options.iconColor)
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 32 32" x="28" y="25" fill="${options.iconColor}" width="100">
                    <image width="15" xlink:href="${base64Image}"></image>
                </svg>
            </svg>`)}`
}

export function generateObjIcon(image, options) {
  if (image instanceof Image) {
    const base64Image = getBase64Image(image, options.iconColor)
    return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                    <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 18 18" x="22" y="5" fill="${options.iconColor}" width="65" >
                        <image width="15" xlink:href="${base64Image}"></image>
                    </svg>
                </svg>`
  }
  options = image
  return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                </svg>`
}

export function cached(fn) {
  const cache = Object.create(null)
  return (function cachedFn(str, ...args) {
    const hit = cache[str]
    return hit || (cache[str] = fn.apply(null, [str, ...args]))
  })
}

export const camelize = cached((str, separator = '-') => {
  const camelizeRE = new RegExp(`${separator}(\\w)`, 'g')
  return str.replace(camelizeRE, (_, c) => (c ? c.toUpperCase() : ''))
})

export const swapItem = (arr, fromIndex, toIndex) => {
  // eslint-disable-next-line prefer-destructuring
  arr[toIndex] = arr.splice(fromIndex, 1, arr[toIndex])[0]
  return arr
}

export const escapeRegexChar = (str) => {
  // eslint-disable-next-line no-useless-escape
  const escapeCharRE = /([\*\.\?\+\$\^\[\]\(\)\{\}\|\\\/])/g
  return str.replace(escapeCharRE, '\\$1')
}

/**
 * @param {*} event 事件对象
 * @param {*} cb 回调
 * @param {*} keyCode 调用回调的键值数组 默认为回车键
 */
export const keyupCallMethod = (event, cb, keyCode = [13]) => {
  if (!event || typeof cb !== 'function' || !keyCode instanceof Array) return
  const { keyCode: nowKey } = event
  if (keyCode.includes(nowKey)) {
    cb?.()
  }
}

/**
 * 将内容下载为文件
 * @param {string} content 内容
 * @param {string} filename 文件名
 */
export const downloadFile = (content, filename) => {
  const blob = new Blob([content])
  const url = URL.createObjectURL(blob)

  const a = document.createElement('a')
  a.style.display = 'none'
  a.href = url
  a.download = filename

  document.body.appendChild(a)
  a.click()

  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

/**
 * 获取当前页面光标位置
 * @param {Element} element 被聚焦的元素
 * @returns 光标位置
 */
export const getCursorPosition = (element) => {
  const selection = window.getSelection()
  let caretOffset = 0
  // 选中的区域
  if (selection.rangeCount > 0) {
    // false表示进行了范围选择
    const { isCollapsed } = selection
    const range = selection.getRangeAt(0)
    // 克隆一个选中区域
    const preCaretRange = range.cloneRange()
    // 设置选中区域的节点内容为当前节点
    preCaretRange.selectNodeContents(element)
    // 重置选中区域的结束位置
    preCaretRange.setEnd(range.endContainer, range.endOffset)
    const { length } = preCaretRange.toString()
    caretOffset = isCollapsed ? length : length - selection.toString().length
  }
  return caretOffset
}

/**
 * 设置当前页面光标位置
 * @param {Element} element 被聚焦的元素
 * @param {number} cursor 要设置的位置
 */
export const setCursorPosition = (element, cursor) => {
  const selection = window.getSelection()
  // 创建一个选中区域
  const range = document.createRange()
  // 选中节点的内容
  range.selectNodeContents(element)
  // 通过计算文本节点的偏移量来设置光标
  const parentAllNodes = element.childNodes
  for (let i = 0; i < parentAllNodes.length; i++) {
    const nowNode = parentAllNodes[i]
    const nodeLength = nowNode?.length ?? nowNode?.innerText?.length
    if (cursor <= nodeLength) {
      range.setStart(nowNode?.firstChild || nowNode, cursor)
      break
    }
    cursor -= nodeLength
  }
  // 设置选中区域为一个点
  range.collapse(true)
  // 移除所有的选中范围
  selection.removeAllRanges()
  // 添加新建的范围
  selection.addRange(range)
}

function calculateNodeStyling(targetElement) {
  const style = window.getComputedStyle(targetElement)
  const boxSizing = style.getPropertyValue('box-sizing')
  const paddingSize = Number.parseFloat(style.getPropertyValue('padding-bottom'))
    + Number.parseFloat(style.getPropertyValue('padding-top'))
  const borderSize = Number.parseFloat(style.getPropertyValue('border-bottom-width'))
    + Number.parseFloat(style.getPropertyValue('border-top-width'))
  const contextStyle = CONTEXT_STYLE.map(name => `${name}:${style.getPropertyValue(name)}`).join(';')

  return { contextStyle, paddingSize, borderSize, boxSizing }
}

/**
 * textarea框悬浮计算当前高度
 */
export function calcTextareaHeight(targetElement, rows = 1) {
  if (!targetElement) return
  let hiddenTextarea = document.createElement('textarea')
  document.body.appendChild(hiddenTextarea)
  const { paddingSize, borderSize, boxSizing, contextStyle } = calculateNodeStyling(targetElement)

  hiddenTextarea.setAttribute('style', `${contextStyle};${HIDDEN_STYLE}`)
  hiddenTextarea.value = targetElement.value || ''

  let height = hiddenTextarea.scrollHeight
  const result = {}

  if (boxSizing === 'border-box') {
    height = height + borderSize
  } else if (boxSizing === 'content-box') {
    height = height - paddingSize
  }

  hiddenTextarea.value = ''
  const singleRowHeight = hiddenTextarea.scrollHeight - paddingSize

  if (Number.isInteger(rows)) {
    let minHeight = singleRowHeight * rows
    if (boxSizing === 'border-box') {
      minHeight = minHeight + paddingSize + borderSize
    }
    height = Math.max(minHeight, height)
    result.minHeight = minHeight
  }

  hiddenTextarea.parentNode?.removeChild(hiddenTextarea)
  hiddenTextarea = undefined
  result.height = height
  return result
}

/**
 * 获取 动态分组/筛选 条件选择框操作后得到的 添加/删除 数据
 */
export const getConditionSelect = (val, oldVal) => {
  const addSelect = []
  const deleteSelect = []
  const selectedSet = new Set()

  val.forEach(property => selectedSet.add(`${property.bk_property_id}-${property.bk_obj_id}`))
  oldVal.forEach((property) => {
    const { bk_property_id: propertyId, bk_obj_id: modelId } = property
    const key = `${propertyId}-${modelId}`
    if (selectedSet.has(key)) {
      selectedSet.delete(key)
    } else {
      deleteSelect.push(property)
    }
  })
  val.forEach((property) => {
    const { bk_property_id: propertyId, bk_obj_id: modelId } = property
    const key = `${propertyId}-${modelId}`
    if (selectedSet.has(key)) {
      addSelect.push(property)
    }
  })

  return {
    addSelect,
    deleteSelect
  }
}

/**
 * 更新 筛选/动态分组 添加/删除 条件UI
 */
export const updatePropertySelect = (selected, remove, addSelect, deleteSelect, type = 'push', filterCondition = []) => {
  if (filterCondition.length) {
    addSelect = addSelect.filter(item => !filterCondition.includes(item.bk_property_id))
  }
  deleteSelect.forEach(property => remove(property))
  let start = 0
  const limit = 10
  while (start < addSelect.length) {
    setTimeout(() =>  selected[type](...addSelect.splice(0, limit)))
    start += limit
  }
}

export function getIPInfo(ip, ipv6, cloudId) {
  const hostList = ip ? ip.split(',') : ipv6.split(',')
  const { length } = hostList
  let info
  if (ip) {
    info = length > 1 ? `${hostList[0]}...` : hostList[0]
  } else {
    // eslint-disable-next-line no-useless-escape
    info = length > 1 ? `\[${hostList[0]}\]...` : `\[${hostList[0]}\]`
  }
  return `${cloudId}:${info}`
}

export function getHostInfoTitle(ip, ipv6, cloudId, hostId) {
  if (!(ip || ipv6)) {
    return `(${t('主机ID')})${hostId}`
  }
  return getIPInfo(ip, ipv6, cloudId)
}

export function isNumeric(str) {
  return !isNaN(str) && !isNaN(parseFloat(str))
}

export function* paginateIterator(list, pageSize) {
  let index = 0
  while (index < list.length) {
    yield list.slice(index, index + pageSize)
    index += pageSize
  }
}

/**
 * 设置当前字段是否分割字符串
 * @param {string} id 字段ID
 * @param {funtion} fn 自定义方法
 * @return true | false
 **/
export function isPasteSplit(id, fn = () => false) {
  if (BUILTIN_PASTE_SPLIT_FIELDS.includes(id)) return true
  return fn?.()
}
