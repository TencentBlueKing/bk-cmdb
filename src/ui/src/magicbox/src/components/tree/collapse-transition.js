'use strict'

var Transition = {
    'before-enter' (el) {
        // _dom.addClass(el, 'wz-animate')
        if (!el.dataset) el.dataset = {}

        el.dataset.oldPaddingTop = el.style.paddingTop
        el.dataset.oldPaddingBottom = el.style.paddingBottom

        el.style.height = '0'
        el.style.paddingTop = 0
        el.style.paddingBottom = 0
    },
    'enter' (el) {
        el.dataset.oldOverflow = el.style.overflow
        if (el.scrollHeight !== 0) {
            el.style.height = el.scrollHeight + 'px'
            el.style.paddingTop = el.dataset.oldPaddingTop
            el.style.paddingBottom = el.dataset.oldPaddingBottom
        } else {
            el.style.height = ''
            el.style.paddingTop = el.dataset.oldPaddingTop
            el.style.paddingBottom = el.dataset.oldPaddingBottom
        }

        el.style.overflow = 'hidden'
    },
    'after-enter'  (el) {
        // for safari: remove class then reset height is necessary
        // _dom.removeClass(el, 'wz-animate')
        el.style.height = ''
        el.style.overflow = el.dataset.oldOverflow
    },
    'before-leave' (el) {
        if (!el.dataset) el.dataset = {}
        el.dataset.oldPaddingTop = el.style.paddingTop
        el.dataset.oldPaddingBottom = el.style.paddingBottom
        el.dataset.oldOverflow = el.style.overflow

        el.style.height = el.scrollHeight + 'px'
        el.style.overflow = 'hidden'
    },
    'leave' (el) {
        if (el.scrollHeight !== 0) {
            // for safari: add class after set height, or it will jump to zero height suddenly, weired
            // (0, _dom.addClass)(el, 'wz-animate')
            el.style.height = 0
            el.style.paddingTop = 0
            el.style.paddingBottom = 0
        }
    },
    'after-leave' (el) {
        // _dom.removeClass(el, 'wz-animate')
        el.style.height = ''
        el.style.overflow = el.dataset.oldOverflow
        el.style.paddingTop = el.dataset.oldPaddingTop
        el.style.paddingBottom = el.dataset.oldPaddingBottom
    }
}
export default {
    name: 'CollapseTransition',
    functional: true,
    render (h, { children }) {
        const data = {
            on: Transition
        }
        return h('transition', data, children)
    }
}
