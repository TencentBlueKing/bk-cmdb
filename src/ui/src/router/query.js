import router from './index'
import throttle from 'lodash.throttle'
import deepEqual from 'deep-equal'
import { redirect } from './actions'

// Vue.watch的options
function createWatchOptions (key, options) {
    const watchOptions = {
        immediate: false,
        deep: false
    }
    for (const key in watchOptions) {
        if (options.hasOwnProperty(key)) {
            watchOptions[key] = options[key]
        }
    }
    if (key === '*') {
        watchOptions.deep = true
    }
    return watchOptions
}

// 始终watch router.query, 根据watch的key再做变更比对，除immediate外，无变更时不触发注册的handler
function createCallback (keys, handler, options = {}) {
    let immediateCalled = false
    const callback = (values, oldValues = {}) => {
        let execValue, execOldValue
        if (Array.isArray(keys)) {
            execValue = {}
            execOldValue = {}
            keys.forEach(key => {
                execValue[key] = values[key]
                execOldValue[key] = oldValues[key]
            })
        } else if (keys === '*') {
            execValue = { ...values }
            execOldValue = { ...oldValues }
            if (options.hasOwnProperty('ignore')) {
                const ignoreKeys = Array.isArray(options.ignore) ? options.ignore : [options.ignore]
                ignoreKeys.forEach(key => {
                    delete execValue[key]
                    delete execOldValue[key]
                })
            }
        } else {
            execValue = values[keys]
            execOldValue = oldValues[keys]
        }
        if (options.immediate && !immediateCalled) {
            immediateCalled = true
            handler(execValue, execOldValue)
        } else {
            const hasChange = !deepEqual(execValue, execOldValue)
            hasChange && handler(execValue, execOldValue)
        }
    }

    if (options.hasOwnProperty('throttle')) {
        const interval = typeof options.throttle === 'number' ? options.throttle : 100
        return throttle(callback, interval, { leading: false, trailing: true })
    }

    return callback
}

function isEmpty (value) {
    return value === '' || value === undefined || value === null
}

class RouterQuery {
    constructor () {
        this.router = router
    }
    get app () {
        return router.app
    }

    get route () {
        return this.app.$route
    }

    get (key, defaultValue) {
        if (this.route.query.hasOwnProperty(key)) {
            return this.route.query[key]
        }
        if (arguments.length === 2) {
            return defaultValue
        }
    }

    getAll () {
        return this.route.query
    }

    set (key, value) {
        const query = { ...this.route.query }
        if (typeof key === 'object') {
            Object.assign(query, key)
        } else {
            query[key] = value
        }
        Object.keys(query).forEach(queryKey => {
            if (isEmpty(query[queryKey])) {
                delete query[queryKey]
            }
        })
        redirect({
            ...this.route,
            query: query
        })
    }

    delete (key) {
        const query = {
            ...this.route.query
        }
        delete query[key]
        redirect({
            ...this.route,
            query: query
        })
    }

    clear () {
        redirect({
            ...this.route,
            query: {}
        })
    }

    watch (key, handler, options = {}) {
        const watchOptions = createWatchOptions(key, options)
        const callback = createCallback(key, handler, options)
        const expression = () => this.route.query
        return this.app.$watch(expression, callback, watchOptions)
    }
}

export default new RouterQuery()
