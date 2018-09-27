const Transition = {
    beforeEnter (el) {
        el.classList.add('collapse-transition')
        if (!el.dataset) {
            el.dataset = {}
        }

        el.dataset.oldPaddingTop = el.style.paddingTop
        el.dataset.oldPaddingBottom = el.style.paddingBottom

        el.style.height = '0'
        el.style.paddingTop = 0
        el.style.paddingBottom = 0
    },

    enter (el) {
        el.dataset.oldOverflow = el.style.overflow
        el.style.overflow = 'hidden'
        if (el.scrollHeight !== 0) {
            el.style.height = el.scrollHeight + 'px'
            el.style.paddingTop = el.dataset.oldPaddingTop
            el.style.paddingBottom = el.dataset.oldPaddingBottom
        } else {
            el.style.height = ''
            el.style.paddingTop = el.dataset.oldPaddingTop
            el.style.paddingBottom = el.dataset.oldPaddingBottom
        }
    },

    afterEnter (el) {
        el.classList.remove('collapse-transition')
        el.style.height = ''
        el.style.overflow = el.dataset.oldOverflow
    },

    beforeLeave (el) {
        if (!el.dataset) el.dataset = {}
        el.dataset.oldPaddingTop = el.style.paddingTop
        el.dataset.oldPaddingBottom = el.style.paddingBottom
        el.dataset.oldOverflow = el.style.overflow

        el.style.height = el.scrollHeight + 'px'
        el.style.overflow = 'hidden'
    },

    leave (el) {
        if (el.scrollHeight !== 0) {
            el.classList.add('collapse-transition')
            el.style.height = 0
            el.style.paddingTop = 0
            el.style.paddingBottom = 0
        }
    },

    afterLeave (el) {
        el.classList.remove('collapse-transition')
        el.style.height = ''
        el.style.overflow = el.dataset.oldOverflow
        el.style.paddingTop = el.dataset.oldPaddingTop
        el.style.paddingBottom = el.dataset.oldPaddingBottom
    }
}

const toCamelCase = function (str) {
    return str.replace(/-([a-z])/g, function (g) { return g[1].toUpperCase() })
}

export default {
    name: 'cmdb-collapse-transition',
    functional: true,
    render (h, context) {
        const events = context.data.on || {}
        const camelCaseEvents = {}
        const transitionEvents = {}
        for (let event in events) {
            camelCaseEvents[toCamelCase(event)] = events[event]
        }
        for (let event in Transition) {
            if (camelCaseEvents.hasOwnProperty(event)) {
                transitionEvents[event] = el => {
                    Transition[event](el)
                    camelCaseEvents[event]()
                }
            } else {
                transitionEvents[event] = Transition[event]
            }
        }
        const data = {
            on: transitionEvents
        }
        return h('transition', data, context.children)
    }
}
