export default class CachedPromise {
    constructor () {
        this.cache = {}
    }

    get (id) {
        if (typeof id === 'undefined') {
            return Object.keys(this.cache).map(requestId => this.cache[requestId])
        }
        return this.cache[id]
    }

    set (id, promise) {
        Object.assign(this.cache, {[id]: promise})
    }

    delete (deleteIds) {
        let requestIds = []
        if (typeof deleteIds === 'undefined') {
            requestIds = Object.keys(this.cache)
        } else if (deleteIds instanceof Array) {
            deleteIds.forEach(id => {
                if (this.get(id)) {
                    requestIds.push(id)
                }
            })
        } else if (this.get(deleteIds)) {
            requestIds.push(deleteIds)
        }
        requestIds.forEach(requestId => {
            delete this.cache[requestId]
        })
        return Promise.resolve(deleteIds)
    }
}
