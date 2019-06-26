const requestFrame = window.requestAnimationFrame
    || window.mozRequestAnimationFrame
    || window.webkitRequestAnimationFrame
    || function (fn) {
        return window.setTimeout(fn, 20)
    }

const cancelFrame = window.cancelAnimationFrame
    || window.mozCancelAnimationFrame
    || window.webkitCancelAnimationFrame
    || window.clearTimeout

const addEventListener = (el, binding) => {
    el.addEventListener('mouseenter', mouseenter)
    el.addEventListener('mousemove', mousemove)
    el.addEventListener('mouseleave', mouseleave)
    el.addEventListener('click', click)
}
const removeEventListener = el => {
    el.removeEventListener('mouseenter', mouseenter)
    el.removeEventListener('mousemove', mousemove)
    el.removeEventListener('mouseleave', mouseleave)
    el.removeEventListener('click', click)
}

const options = {
    x: 0,
    y: 0,
    width: 16,
    height: 16,
    zIndex: 100000,
    cursor: 'pointer',
    className: 'v-cursor',
    activeClass: 'v-cursor-active'
}

const mouseenter = event => {
    const el = event.currentTarget
    const data = el.__cursor__
    if (data.active) {
        el.style.cursor = data.cursor
        proxy.style.display = 'block'
        el.classList.add(data.activeClass)
        updateProxyPosition(event)
    }
}

const mousemove = event => {
    const el = event.currentTarget
    const data = el.__cursor__
    if (data.active) {
        updateProxyPosition(event)
    }
}

const mouseleave = event => {
    const el = event.currentTarget
    const data = el.__cursor__
    el.style.cursor = ''
    proxy.style.display = 'none'
    el.classList.remove(data.activeClass)
}

const click = event => {
    const el = event.currentTarget
    const data = el.__cursor__
    if (!data.active) {
        return false
    }
    const callback = data.onclick
    if (typeof callback === 'function') {
        callback(data)
    }
    const globalCallback = data.globalCallback
    if (typeof globalCallback === 'function') {
        globalCallback(data)
    }
}

let proxy = null
let frameId = null

const createProxy = () => {
    proxy = document.createElement('span')
    proxy.style.position = 'fixed'
    proxy.style.pointerEvents = 'none'
    proxy.style.zIndex = options.zIndex
    proxy.style.width = options.width + 'px'
    proxy.style.height = options.height + 'px'
    proxy.classList.add(options.className)
    document.body.append(proxy)
}

const updateProxyPosition = event => {
    const el = event.currentTarget
    const data = el.__cursor__
    if (frameId) {
        cancelFrame(frameId)
    }
    frameId = requestFrame(() => {
        proxy.style.left = event.clientX + data.x + 'px'
        proxy.style.top = event.clientY + data.y + 'px'
    })
}

const setChildrenEvents = (target, pointerEvents) => {
    Array.prototype.forEach.call(target.children, child => {
        child.style.pointerEvents = pointerEvents
    })
}

const cursor = {
    inserted (el, binding, vNode) {
        if (!proxy) {
            createProxy()
        }
        const data = { ...options }
        if (typeof binding.value !== 'object') {
            data.active = binding.value
        } else {
            Object.assign(data, binding.value)
        }
        const target = data.selector ? el.querySelector(data.selector) : el
        if (target) {
            el.__cursor_target__ = target
            target.__cursor__ = data
            addEventListener(target)
            const pointerEvents = data.active ? 'none' : ''
            setChildrenEvents(target, pointerEvents)
        }
    },
    update (el, binding) {
        const data = { ...options }
        if (typeof binding.value !== 'object') {
            data.active = binding.value
        } else {
            Object.assign(data, binding.value)
        }
        let target = el.__cursor_target__
        if (!target) {
            target = el.querySelector(data.selector)
            if (target) {
                el.__cursor_target__ = target
                target.__cursor__ = data
                addEventListener(target)
                const pointerEvents = data.active ? 'none' : ''
                setChildrenEvents(target, pointerEvents)
            }
        } else {
            Object.assign(target.__cursor__, data)
            const pointerEvents = data.active ? 'none' : ''
            setChildrenEvents(target, pointerEvents)
        }
    },
    unbind (el) {
        const target = el.__cursor_target__
        removeEventListener(target)
    }
}

export default {
    install: Vue => Vue.directive('cursor', cursor),
    directive: cursor,
    setOptions: customOptions => {
        Object.assign(options, customOptions)
    }
}
