<template>
    <div class="vue-waterfall" :style="style">
      <slot></slot>
    </div>
</template>

<style>
    .vue-waterfall {
        position: relative;
        /*overflow: hidden; cause clientWidth = 0 in IE if height not bigger than 0 */
    }
</style>

<script>
/* eslint-disable */

const MOVE_CLASS_PROP = '_wfMoveClass'

export default {
    props: {
        autoResize: {
            default: true
        },
        interval: {
            default: 200,
            validator: (val) => val >= 0
        },
        align: {
            default: 'left',
            validator: (val) => ~['left', 'right', 'center'].indexOf(val)
        },
        line: {
            default: 'v',
            validator: (val) => ~['v', 'h'].indexOf(val)
        },
        lineGap: {
            required: true,
            validator: (val) => val >= 0
        },
        minLineGap: {
            validator: (val) => val >= 0
        },
        maxLineGap: {
            validator: (val) => val >= 0
        },
        singleMaxWidth: {
            validator: (val) => val >= 0
        },
        fixedHeight: {
            default: false
        },
        grow: {
            validator: (val) => val instanceof Array
        },
        watch: {
            default: () => ({})
        }
    },
    data: () => ({
        style: {
            height: '',
            overflow: ''
        },
        token: null
    }),
    methods: {
        reflowHandler,
        autoResizeHandler,
        reflow
    },
    created () {
        this.virtualRects = []
        this.$on('reflow', () => {
            this.reflowHandler()
        })
        this.$watch(() => (
            this.align,
            this.line,
            this.lineGap,
            this.minLineGap,
            this.maxLineGap,
            this.singleMaxWidth,
            this.fixedHeight,
            this.watch
        ), this.reflowHandler)
        this.$watch('grow', this.reflowHandler)
    },
    mounted () {
        this.$watch('autoResize', this.autoResizeHandler)
        on(this.$el, getTransitionEndEvent(), tidyUpAnimations, true)
        this.autoResizeHandler(this.autoResize)
    },
    beforeDestroy () {
        this.autoResizeHandler(false)
        off(this.$el, getTransitionEndEvent(), tidyUpAnimations, true)
    }
}

function autoResizeHandler (autoResize) {
    if (autoResize === false || !this.autoResize) {
        off(window, 'resize', this.reflowHandler, false)
    } else {
        on(window, 'resize', this.reflowHandler, false)
    }
}

function tidyUpAnimations (event) {
    let node = event.target
    let moveClass = node[MOVE_CLASS_PROP]
    if (moveClass) {
        removeClass(node, moveClass)
    }
}

function reflowHandler () {
    clearTimeout(this.token)
    this.token = setTimeout(this.reflow, this.interval)
}

function reflow () {
    if (!this.$el) { return }
    let width = this.$el.clientWidth
    let metas = this.$children.map((slot) => slot.getMeta())
    metas.sort((a, b) => a.order - b.order)
    this.virtualRects = metas.map(() => ({}))
    calculate(this, metas, this.virtualRects)
    setTimeout(() => {
        if (isScrollBarVisibilityChange(this.$el, width)) {
            calculate(this, metas, this.virtualRects)
        }
        this.style.overflow = 'hidden'
        render(this.virtualRects, metas)
        this.$emit('reflowed', this)
    }, 0)
}

function isScrollBarVisibilityChange (el, lastClientWidth) {
    return lastClientWidth !== el.clientWidth
}

function calculate (vm, metas, styles) {
    let options = getOptions(vm)
    let processor = vm.line === 'h' ? horizontalLineProcessor : verticalLineProcessor
    processor.calculate(vm, options, metas, styles)
}

function getOptions (vm) {
    const maxLineGap = vm.maxLineGap ? +vm.maxLineGap : vm.lineGap
    return {
        align: ~['left', 'right', 'center'].indexOf(vm.align) ? vm.align : 'left',
        line: ~['v', 'h'].indexOf(vm.line) ? vm.line : 'v',
        lineGap: +vm.lineGap,
        minLineGap: vm.minLineGap ? +vm.minLineGap : vm.lineGap,
        maxLineGap: maxLineGap,
        singleMaxWidth: Math.max(vm.singleMaxWidth || 0, maxLineGap),
        fixedHeight: !!vm.fixedHeight,
        grow: vm.grow && vm.grow.map(val => +val)
    }
}

var verticalLineProcessor = (() => {

    function calculate (vm, options, metas, rects) {
        let width = vm.$el.clientWidth
        let grow = options.grow
        let strategy = grow
        ? getRowStrategyWithGrow(width, grow)
        : getRowStrategy(width, options)
        let tops = getArrayFillWith(0, strategy.count)
        metas.forEach((meta, index) => {
            let offset = tops.reduce((last, top, i) => top < tops[last] ? i : last, 0)
            let width = strategy.width[offset % strategy.count]
            let rect = rects[index]
            rect.top = tops[offset]
            rect.left = strategy.left + (offset ? sum(strategy.width.slice(0, offset)) : 0)
            rect.width = width
            rect.height = meta.height * (options.fixedHeight ? 1 : width / meta.width)
            tops[offset] = tops[offset] + rect.height
        })
        vm.style.height = Math.max.apply(Math, tops) + 'px'
    }

    function getRowStrategy (width, options) {
        let count = width / options.lineGap
        let slotWidth
        if (options.singleMaxWidth >= width) {
            count = 1
            slotWidth = Math.max(width, options.minLineGap)
        } else {
            let maxContentWidth = options.maxLineGap * ~~count
            let minGreedyContentWidth = options.minLineGap * ~~(count + 1)
            let canFit = maxContentWidth >= width
            let canFitGreedy = minGreedyContentWidth <= width
            if (canFit && canFitGreedy) {
                count = Math.round(count)
                slotWidth = width / count
            } else if (canFit) {
                count = ~~count
                slotWidth = width / count
            } else if (canFitGreedy) {
                count = ~~(count + 1)
                slotWidth = width / count
            } else {
                count = ~~count
                slotWidth = options.maxLineGap
            }
            if (count === 1) {
                slotWidth = Math.min(width, options.singleMaxWidth)
                slotWidth = Math.max(slotWidth, options.minLineGap)
            }
        }
        return {
            width: getArrayFillWith(slotWidth, count),
            count: count,
            left: getLeft(width, slotWidth * count, options.align)
        }
    }

    function getRowStrategyWithGrow (width, grow) {
        let total = sum(grow)
        return {
            width: grow.map(val => width * val / total),
            count: grow.length,
            left: 0
        }
    }

    return {
        calculate
    }
})()

var horizontalLineProcessor = (() => {

    function calculate (vm, options, metas, rects) {
        let width = vm.$el.clientWidth
        let total = metas.length
        let top = 0
        let offset = 0
        while (offset < total) {
            let strategy = getRowStrategy(width, options, metas, offset)
            for (let i = 0, left = 0, meta, rect; i < strategy.count; i++) {
                meta = metas[offset + i]
                rect = rects[offset + i]
                rect.top = top
                rect.left = strategy.left + left
                rect.width = meta.width * strategy.height / meta.height
                rect.height = strategy.height
                left += rect.width
            }
            offset += strategy.count
            top += strategy.height
        }
        vm.style.height = top + 'px'
    }

    function getRowStrategy (width, options, metas, offset) {
        let greedyCount = getGreedyCount(width, options.lineGap, metas, offset)
        let lazyCount = Math.max(greedyCount - 1, 1)
        let greedySize = getContentSize(width, options, metas, offset, greedyCount)
        let lazySize = getContentSize(width, options, metas, offset, lazyCount)
        let finalSize = chooseFinalSize(lazySize, greedySize, width)
        let height = finalSize.height
        let fitContentWidth = finalSize.width
        if (finalSize.count === 1) {
            fitContentWidth = Math.min(options.singleMaxWidth, width)
            height = metas[offset].height * fitContentWidth / metas[offset].width
        }
        return {
            left: getLeft(width, fitContentWidth, options.align),
            count: finalSize.count,
            height: height
        }
    }

    function getGreedyCount (rowWidth, rowHeight, metas, offset) {
        let count = 0
        for (let i = offset, width = 0; i < metas.length && width <= rowWidth; i++) {
            width += metas[i].width * rowHeight / metas[i].height
            count++
        }
        return count
    }

    function getContentSize (rowWidth, options, metas, offset, count) {
        let originWidth = 0
        for (let i = count - 1; i >= 0; i--) {
            let meta = metas[offset + i]
            originWidth += meta.width * options.lineGap / meta.height
        }
        let fitHeight = options.lineGap * rowWidth / originWidth
        let canFit = (fitHeight <= options.maxLineGap && fitHeight >= options.minLineGap)
        if (canFit) {
            return {
                cost: Math.abs(options.lineGap - fitHeight),
                count: count,
                width: rowWidth,
                height: fitHeight
            }
        } else {
            let height = originWidth > rowWidth ? options.minLineGap : options.maxLineGap
            return {
                cost: Infinity,
                count: count,
                width: originWidth * height / options.lineGap,
                height: height
            }
        }
    }

    function chooseFinalSize (lazySize, greedySize, rowWidth) {
        if (lazySize.cost === Infinity && greedySize.cost === Infinity) {
            return greedySize.width < rowWidth ? greedySize : lazySize
        } else {
            return greedySize.cost >= lazySize.cost ? lazySize : greedySize
        }
    }

    return {
        calculate
    }
})()

function getLeft (width, contentWidth, align) {
    switch (align) {
        case 'right':
            return width - contentWidth
        case 'center':
            return (width - contentWidth) / 2
        default:
        return 0
    }
}

function sum (arr) {
    return arr.reduce((sum, val) => sum + val)
}

function render (rects, metas) {
    let metasNeedToMoveByTransform = metas.filter((meta) => meta.moveClass)
    let firstRects = getRects(metasNeedToMoveByTransform)
    applyRects(rects, metas)
    let lastRects = getRects(metasNeedToMoveByTransform)
    metasNeedToMoveByTransform.forEach((meta, i) => {
        meta.node[MOVE_CLASS_PROP] = meta.moveClass
        setTransform(meta.node, firstRects[i], lastRects[i])
    })
    document.body.clientWidth // forced reflow
    metasNeedToMoveByTransform.forEach((meta) => {
        addClass(meta.node, meta.moveClass)
        clearTransform(meta.node)
    })
}

function getRects (metas) {
    return metas.map((meta) => meta.vm.rect)
}

function applyRects (rects, metas) {
    rects.forEach((rect, i) => {
        let style = metas[i].node.style
        metas[i].vm.rect = rect
        for (let prop in rect) {
            style[prop] = rect[prop] + 'px'
        }
    })
}

function setTransform (node, firstRect, lastRect) {
    let dx = firstRect.left - lastRect.left
    let dy = firstRect.top - lastRect.top
    let sw = firstRect.width / lastRect.width
    let sh = firstRect.height / lastRect.height
    node.style.transform =
    node.style.WebkitTransform = `translate(${dx}px,${dy}px) scale(${sw},${sh})`
    node.style.transitionDuration = '0s'
}

function clearTransform (node) {
    node.style.transform = node.style.WebkitTransform = ''
    node.style.transitionDuration = ''
}

function getTransitionEndEvent () {
    let isWebkitTrans =
        window.ontransitionend === undefined &&
        window.onwebkittransitionend !== undefined
    let transitionEndEvent = isWebkitTrans
        ? 'webkitTransitionEnd'
        : 'transitionend'
    return transitionEndEvent
}

/**
 * util
 */

function getArrayFillWith (item, count) {
    let getter = (typeof item === 'function') ? () => item() : () => item
    let arr = []
    for (let i = 0; i < count; i++) {
        arr[i] = getter()
    }
    return arr
}

function addClass (elem, name) {
    if (!hasClass(elem, name)) {
        let cur = attr(elem, 'class').trim()
        let res = (cur + ' ' + name).trim()
        attr(elem, 'class', res)
    }
}

function removeClass (elem, name) {
    let reg = new RegExp('\\s*\\b' + name + '\\b\\s*', 'g')
    let res = attr(elem, 'class').replace(reg, ' ').trim()
    attr(elem, 'class', res)
}

function hasClass (elem, name) {
    return (new RegExp('\\b' + name + '\\b')).test(attr(elem, 'class'))
}

function attr (elem, name, value) {
    if (typeof value !== 'undefined') {
        elem.setAttribute(name, value)
    } else {
        return elem.getAttribute(name) || ''
    }
}

function on (elem, type, listener, useCapture = false) {
    elem.addEventListener(type, listener, useCapture)
}

function off (elem, type, listener, useCapture = false) {
    elem.removeEventListener(type, listener, useCapture)
}

</script>
