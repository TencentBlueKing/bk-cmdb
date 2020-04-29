import router, { addBeforeHooks } from './index'
import throttle from 'lodash.throttle'
import { redirect } from './actions'

class RouterQuery {
    constructor () {
        this.unwatchs = []
        addBeforeHooks(() => {
            this.unwatchs.forEach(unwatch => unwatch())
            this.unwatchs = []
        })
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
        const watchOptions = {
            immediate: false,
            deep: false
        }
        for (const optionKey in watchOptions) {
            if (options.hasOwnProperty(optionKey)) {
                watchOptions[optionKey] = options[optionKey]
            }
        }
        let callback = handler
        if (options.hasOwnProperty('throttle')) {
            const interval = typeof options.throttle === 'number' ? options.throttle : 100
            callback = throttle(handler, interval, { leading: false, trailing: true })
        }
        let expression = () => this.route.query[key]
        if (key === '*') {
            expression = () => this.route.query
            watchOptions.deep = true
        } else if (Array.isArray(key)) {
            expression = () => {
                const values = {}
                key.forEach(key => {
                    values[key] = this.route.query[key]
                })
                return values
            }
        }
        this.unwatchs.push(this.app.$watch(expression, callback, watchOptions))
    }
}

export default new RouterQuery()
