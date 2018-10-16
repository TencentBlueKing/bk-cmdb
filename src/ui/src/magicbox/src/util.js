export function isVNode (node) {
    return typeof node === 'object' && node.hasOwnProperty('componentOptions')
}

export function isInArray (ele, array) {
    // for (let item of array) {
    //     if (item === ele) {
    //         return true
    //     }
    // }
    const len = array.length
    for (let i = 0; i < len; i++) {
        if (array[i] === ele) {
            return true
        }
    }

    return false
}

export function isInlineElment (node) {
    let inlineElements = ['a', 'abbr', 'acronym', 'b', 'bdo', 'big', 'br', 'cite', 'code', 'dfn', 'em', 'font', 'i', 'img', 'input', 'kbd', 'label', 'q', 's', 'samp', 'select', 'small', 'span', 'strike', 'strong', 'sub', 'sup', 'textarea', 'tt', 'u', 'var']
    let tag = (node.tagName).toLowerCase()
    let display = getComputedStyle(node).display

    if ((isInArray(tag, inlineElements) && display === 'index') || display === 'inline') {
        console.warn('Binding node is displayed as inline element. To avoid some unexpected rendering error, please set binding node displayed as block element.')

        return true
    }

    return false
}

/**
 * 获取元素相对于页面的高度
 *
 * @param node {Object} 指定的 DOM 元素
 *
 * @return {number} 高度值
 */
export function getActualTop (node) {
    let actualTop = node.offsetTop
    let current = node.offsetParent

    while (current !== null) {
        actualTop += current.offsetTop
        current = current.offsetParent
    }

    return actualTop
}

/**
 * 获取元素相对于页面左侧的宽度
 *
 * @param node {Object} 指定的 DOM 元素
 *
 * @return {number} 宽度值
 */
export function getActualLeft (node) {
    let actualLeft = node.offsetLeft
    let current = node.offsetParent

    while (current !== null) {
        actualLeft += current.offsetLeft
        current = current.offsetParent
    }

    return actualLeft
}

/**
 *  对元素添加样式类
 *  @param node {NodeElement} 指定的DOM元素
 *  @param className {String} 类名
 */
export function addClass (node, className) {
    let classNames = className.split(' ')
    if (node.nodeType === 1) {
        if (!node.className && classNames.length === 1) {
            node.className = className
        } else {
            let setClass = ' ' + node.className + ' '
            classNames.forEach((cl) => {
                if (setClass.indexOf(' ' + cl + ' ') < 0) {
                    setClass += cl + ' '
                }
            })
            let rtrim = /^\s+|\s+$/
            node.className = setClass.replace(rtrim, '')
        }
    }
}

/**
 *  对元素删除样式类
 *  @param node {NodeElement} 指定的DOM元素
 *  @param className {String} 类名
 */
export function removeClass (node, className) {
    let classNames = className.split(' ')
    if (node.nodeType === 1) {
        let setClass = ' ' + node.className + ' '
        classNames.forEach((cl) => {
            setClass = setClass.replace(' ' + cl + ' ', ' ')
        })
        let rtrim = /^\s+|\s+$/
        node.className = setClass.replace(rtrim, '')
    }
}

/**
 * 字符串转换为驼峰写法
 *
 * @param {string} str 待转换字符串
 *
 * @return {string} 转换后字符串
 */
export function camelize (str) {
    return str.replace(/-(\w)/g, (strMatch, p1) => p1.toUpperCase())
}

/**
 * 获取元素的样式
 *
 * @param {Object} elem dom 元素
 * @param {string} prop 样式属性
 *
 * @return {string} 样式值
 */
export function getStyle (elem, prop) {
    if (!elem || !prop) {
        return false
    }

    // 先获取是否有内联样式
    let value = elem.style[camelize(prop)]

    if (!value) {
        // 获取的所有计算样式
        let css = ''
        if (document.defaultView && document.defaultView.getComputedStyle) {
            css = document.defaultView.getComputedStyle(elem, null)
            value = css ? css.getPropertyValue(prop) : null
        }
    }

    return String(value)
}

const monthLong = {
    '01': 'January',
    '02': 'February',
    '03': 'March',
    '04': 'April',
    '05': 'May',
    '06': 'June',
    '07': 'July',
    '08': 'August',
    '09': 'September',
    '10': 'October',
    '11': 'November',
    '12': 'December'
}

const monthShort = {
    '01': 'Jan',
    '02': 'Feb',
    '03': 'Mar',
    '04': 'Apr',
    '05': 'May',
    '06': 'Jun',
    '07': 'Jul',
    '08': 'Aug',
    '09': 'Sep',
    '10': 'Oct',
    '11': 'Nov',
    '12': 'Dec'
}

/**
 * 格式化月份
 *
 * @param {string} month 月份值
 * @param {string} locale 语言
 * @param {boolean} isShort 是否简写月份
 *
 * @return {string} 格式化后的月份
 */
export function formatMonth (month, locale = 'en-US', isShort = false) {
    if (locale === 'en-US') {
        return isShort ? monthShort[month] : monthLong[month]
    }
    return month
}

/**
 * 函数防抖
 *
 * @param {Function} func 要执行的函数
 * @param {number} wait 等待时间
 * @param {boolean} immediate 是否立即执行
 *
 * @return {Function} 防抖后的方法
 */
export function debounce (func, wait, immediate) {
    let timeout
    let result
    const debounced = function () {
        const context = this
        const args = arguments

        if (timeout) {
            clearTimeout(timeout)
        }
        if (immediate) {
            // 如果已经执行过，不再执行
            const callNow = !timeout
            timeout = setTimeout(() => {
                timeout = null
            }, wait)
            if (callNow) {
                result = func.apply(context, args)
            }
        }
        else {
            timeout = setTimeout(() => {
                func.apply(context, args)
            }, wait)
        }
        return result
    }

    debounced.cancel = () => {
        clearTimeout(timeout)
        timeout = null
    }

    return debounced
}

const OVERFLOW_PROPERTYS = ['overflow', 'overflow-x', 'overflow-y']

const SCROLL_TYPES = ['scroll', 'auto']

// 最大值
const MAX = 4

// 根元素
const ROOT = document.body

// 竖直方向
const VERTICAL = ['top', 'bottom']
// 水平方向
const HORIZONTAL = ['left', 'right']

// 默认限制显示方向如下，显示优先级按顺序以此递减
const DEFAULT_PLACEMENT_QUEUE = ['top', 'right', 'bottom', 'left']

// 是否是个可滚动的元素
export function checkScrollable (element) {
    const css = window.getComputedStyle(element, null)
    return OVERFLOW_PROPERTYS.some(property => {
        return ~SCROLL_TYPES.indexOf(css[property])
    })
}

// 获取参考元素第一个可滚动的元素
export function getScrollContainer (el) {
    if (!el) {
        return ROOT
    }

    let parent = el.parentNode
    while (parent && parent !== ROOT) {
        if (checkScrollable(parent)) {
            return parent
        }
        parent = parent.parentNode
    }
    return ROOT
}

// 获取最优展示方向，weight 越大对应方向的优先级越高
function getBestPlacement (queue) {
    return queue.sort((a, b) => b.weight - a.weight)[0]
}

// 获取目标元素相对于参考容器的位置信息
function getBoxMargin (el, parent) {
    if (!el) {
        return
    }
    const eBox = el.getBoundingClientRect()
    const pBox = parent.getBoundingClientRect()

    const {width: vw, height: vh} = pBox
    const {width, height} = eBox

    const top = eBox.top - pBox.top
    const left = eBox.left - pBox.left
    const right = eBox.right - pBox.left
    const bottom = eBox.bottom - pBox.top

    const midX = left + width / 2
    const midY = top + height / 2

    // 目标的顶点坐标 [top-left, top-right, bottom-right, botton-left]
    const vertex = {
        tl: {x: left, y: top},
        tr: {x: right, y: top},
        br: {x: right, y: bottom},
        bl: {x: left, y: bottom}
    }

    return {
        width,
        height,
        margin: {
            top: {
                placement: 'top',
                size: top,
                start: vertex.tl,
                mid: {x: midX, y: top},
                end: vertex.tr
            },
            bottom: {
                placement: 'bottom',
                size: vh - bottom,
                start: vertex.bl,
                mid: {x: midX, y: bottom},
                end: vertex.br
            },
            left: {
                placement: 'left',
                size: left,
                start: vertex.tl,
                mid: {x: left, y: midY},
                end: vertex.bl
            },
            right: {
                placement: 'right',
                size: vw - right,
                start: vertex.tr,
                mid: {x: right, y: midY},
                end: vertex.br
            }
        }
    }
}

// ref 参考元素，container 容器， target 需要动态计算坐标的元素，limitQueue 限制展示方向
export function computePlacementInfo (ref, container, target, limitQueue, offset) {
    if (!ref || !target) {
        return
    }
    const placementQueue = limitQueue && limitQueue.length ? limitQueue : DEFAULT_PLACEMENT_QUEUE
    const {width: ew, height: eh, margin} = getBoxMargin(ref, container)
    const {width: tw, height: th} = target.getBoundingClientRect()

    const dw = (tw - ew) / 2
    const dh = (th - eh) / 2

    const queueLen = placementQueue.length
    const processedQueue = Object.keys(margin).map(key => {
        const placementItem = margin[key]
        // 这里 index 可以用来标记显示方向的优先级 index 越大，优先级越高
        const index = placementQueue.indexOf(placementItem.placement)
        placementItem.weight = index > -1 ? MAX - index : MAX - queueLen

        // 上下方向上可容纳目标元素
        const verSingleBiasCheck = (~VERTICAL.indexOf(placementItem.placement) && placementItem.size > th + offset)
        // 上下方向上可容纳目标元素 && 目标元素上下显示时左右也可完整显示目标元素
        const verFullBiasCheck = (verSingleBiasCheck && margin.left.size > dw && margin.right.size > dw)
        // 左右方向上可容纳目标元素
        const horSingleBiasCheck = (HORIZONTAL.indexOf(placementItem.placement) > -1
            && placementItem.size > tw + offset)
        // 左右方向上可容纳目标元素 && 显示时上下也可完整显示目标元素
        const horFullBiasCheck = (horSingleBiasCheck && margin.top.size > dh && margin.bottom.size > dh)
        // 竖直方向上的 top 与 bottom 的间距差值
        placementItem.dVer = margin.top.size - margin.bottom.size
        // 水平方向上的 left 与 right 的间距差值
        placementItem.dHor = margin.left.size - margin.right.size
        placementItem.mod = 'edge'

        if (verFullBiasCheck || horFullBiasCheck) {
            placementItem.mod = 'mid'
            placementItem.weight += 3 + placementItem.weight / MAX
            return placementItem
        }
        if (verSingleBiasCheck || horSingleBiasCheck) {
            placementItem.weight += 2 + placementItem.weight / MAX
        }
        return placementItem
    })
    return Object.assign({ ew, eh, tw, th, dw, dh }, getBestPlacement(processedQueue))
}

// 基于参考元素的某一侧的中点来计算目标元素的坐标
export function computeCoordinateBaseMid (placementInfo, offset) {
    const {placement, mid, tw, th} = placementInfo
    switch (placement) {
        case 'top':
            return {
                placement: 'top-mid',
                x: mid.x - tw / 2,
                y: mid.y - th - offset
            }
        case 'bottom':
            return {
                placement: 'bottom-mid',
                x: mid.x - tw / 2,
                y: mid.y + offset
            }
        case 'left':
            return {
                placement: 'left-mid',
                x: mid.x - tw - offset,
                y: mid.y - th / 2
            }
        case 'right':
            return {
                placement: 'right-mid',
                x: mid.x + offset,
                y: mid.y - th / 2
            }
        default:
    }
}

// 用于计算小三角形在 tip 窗口中的位置
export function computeArrowPos (placement, offset, size) {
    const start = offset + 'px'
    const end = offset - size * 2 + 'px'
    const posMap = {
        'top-start': {top: '100%', left: start},
        'top-mid': {top: '100%', left: '50%'},
        'top-end': {top: '100%', right: end},

        'bottom-start': {top: '0', left: start},
        'bottom-mid': {top: '0', left: '50%'},
        'bottom-end': {top: '0', right: end},

        'left-start': {top: start, left: '100%'},
        'left-mid': {top: '50%', left: '100%'},
        'left-end': {bottom: end, left: '100%'},

        'right-start': {top: start, left: '0'},
        'right-mid': {top: '50%', left: '0'},
        'right-end': {bottom: end, left: '0'}
    }
    return posMap[placement]
}

// 基于参考元素某一侧的边界来计算目标元素位置
export function computeCoordinateBaseEdge (placementInfo, offset) {
    const {placement, start, end, dHor, dVer, tw, th, ew, eh} = placementInfo
    const nearRight = dHor > 0
    const nearBottom = dVer > 0
    switch (placement) {
        case 'top':
            return {
                placement: nearRight ? 'top-end' : 'top-start',
                x: nearRight ? end.x - tw : start.x,
                y: start.y - th - offset,
                arrowsOffset: ew / 2
            }
        case 'bottom':
            return {
                placement: nearRight ? 'bottom-end' : 'bottom-start',
                x: nearRight ? end.x - tw : start.x,
                y: end.y + offset,
                arrowsOffset: ew / 2
            }
        case 'left':
            return {
                placement: nearBottom ? 'left-end' : 'left-start',
                x: start.x - tw - offset,
                y: nearBottom ? end.y - th : start.y,
                arrowsOffset: eh / 2
            }
        case 'right':
            return {
                placement: nearBottom ? 'right-end' : 'right-start',
                x: end.x + offset,
                y: nearBottom ? end.y - th : start.y,
                arrowsOffset: eh / 2
            }
        default:
    }
}

export const requestAnimationFrame = window.requestAnimationFrame
    || window.webkitRequestAnimationFrame
    || window.mozRequestAnimationFrame
    || window.oRequestAnimationFrame
    || window.msRequestAnimationFrame
    || function (callback) {
        window.setTimeout(callback, 1000 / 60)
    }
export const cancelAnimationFrame = window.cancelAnimationFrame
    || window.webkitCancelAnimationFrame
    || window.mozCancelAnimationFrame
    || window.oCancelAnimationFrame
    || window.msCancelAnimationFrame
    || function (id) {
        window.clearTimeout(id)
    }

// 获取唯一随机数
export function uuid () {
    let id = ''
    let randomNum = Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1)

    for (let i = 0; i < 7; i++) {
        id += randomNum
    }
    return id
}
