import router, { addBeforeHooks } from './index'

const unwatchs = []

class RouterQuery {
    constructor () {
        addBeforeHooks(() => {
            unwatchs.forEach(unwatch => unwatch())
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
        this.app.$router.replace({
            ...this.route,
            query: {
                ...this.route.query,
                [key]: value
            }
        })
    }

    delete (key) {
        const query = {
            ...this.route.query
        }
        delete query[key]
        this.app.$router.replace({
            ...this.route,
            query: query
        })
    }

    clear () {
        this.app.$router.replace({
            ...this.route,
            query: {}
        })
    }

    watch (key, callback) {
        unwatchs.push(this.app.$watch(() => this.route.query[key], callback))
    }
}

export default new RouterQuery()
