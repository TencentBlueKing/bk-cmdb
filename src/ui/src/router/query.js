import router from './index'
import throttle from 'lodash.throttle'
import { redirect } from './actions'

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

function createCallback (keys, handler, options = {}) {
    const callback = (values, oldValues = {}) => {
        if (Array.isArray(keys)) {
            const watchValues = {}
            const oldWatchValues = {}
            keys.forEach(key => {
                watchValues[key] = values[key]
                oldWatchValues[key] = oldValues[key]
            })
            const hasChange = keys.some(key => String(watchValues[key]) !== String(oldWatchValues[key]))
            hasChange && handler(watchValues, oldWatchValues)
        } else if (keys === '*') {
            const cloneValues = { ...values }
            const cloneOldValues = { ...oldValues }
            if (options.hasOwnProperty('ignore')) {
                const ignoreKeys = Array.isArray(options.ignore) ? options.ignore : [options.ignore]
                ignoreKeys.forEach(key => {
                    delete cloneValues[key]
                    delete cloneOldValues[key]
                })
            }
            const hasChange = Object.keys(cloneValues).some(key => String(cloneValues[key]) !== String(cloneOldValues[key]))
            hasChange && handler(cloneValues, cloneOldValues)
        } else {
            const value = String(values[keys])
            const oldValue = String(oldValues[keys])
            value !== oldValue && handler(value, oldValue)
        }
    }

    if (options.hasOwnProperty('throttle')) {
        const interval = typeof options.throttle === 'number' ? options.throttle : 100
        return throttle(callback, interval, { leading: false, trailing: true })
    }

    return callback
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
