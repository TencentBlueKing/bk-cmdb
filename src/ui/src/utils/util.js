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
import { PROPERTY_TYPE_EXCLAMATION_TIPS } from '@/dictionary/property-constants'
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
  return String(str).replace(escapeCharRE, '\\$1')
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

export const isExclmationProperty = type => PROPERTY_TYPE_EXCLAMATION_TIPS.includes(type)
