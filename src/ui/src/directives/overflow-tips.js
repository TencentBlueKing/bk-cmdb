import Vue from 'vue'

const _this = Vue.prototype

const instanceMap = new Map()

const popoverOptions = {
    theme: 'v-overflow-tips',
    maxWidth: 400,
    arrow: true,
    interactive: true,
    boundary: 'window'
}

const createDom = el => {
    const css = window.getComputedStyle(el, null)
    const dom = document.createElement('div')
    const width = parseFloat(css['width']) ? Math.ceil(parseFloat(css['width'])) + 'px' : css['width']
    dom.style.cssText = `width: ${width}; line-height: ${css['line-height']}; font-size: ${css['font-size']}; word-break: ${css['word-break']}`
    dom.textContent = el.textContent
    return dom
}

const isOverflow = el => {
    const css = window.getComputedStyle(el, null)
    const lineClamp = css['-webkit-line-clamp']
    if (lineClamp !== 'none' && lineClamp > 1) {
        const targetHeight = parseFloat(css['height'])
        const dom = createDom(el)
        document.body.appendChild(dom)
        const domHeight = window.getComputedStyle(dom, null)['height']
        document.body.removeChild(dom)
        return targetHeight < parseFloat(domHeight)
    }
    return el.clientWidth < el.scrollWidth
}

const setData = el => {
    const mapItem = {
        popover: null,
        textContent: el.textContent
    }
    if (isOverflow(el)) {
        const instance = _this.$bkPopover(el, {
            content: el.textContent,
            ...popoverOptions
        })
        mapItem.popover = instance
    }
    instanceMap.set(el, mapItem)
}

const overflowTips = {
    inserted (el, binding) {
        setData(el)
    },
    update (el, binding) {
        setTimeout(() => {
            const instance = instanceMap.get(el)
            const innerText = el.innerText
            const textContent = el.textContent
            if ([undefined, null, '', '--'].includes(innerText) || (instance && textContent === instance.textContent)) return
            instance && instance.popover && instance.popover.destroy()
            setData(el)
        }, 0)
    },
    unbind (el) {
        const instance = instanceMap.get(el)
        if (instance) {
            instance.popover && instance.popover.destroy()
            instanceMap.delete(el)
        }
    }
}

export default {
    install: Vue => Vue.directive('overflow-tips', overflowTips)
}
