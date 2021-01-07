export default class CachedPromise {
    constructor () {
        this.cache = {}
    }

    get (id) {
        if (typeof id === 'undefined') {
            return Object.keys(this.cache).map(requestId => this.cache[requestId].promise)
        }
        return this.cache.hasOwnProperty(id) ? this.cache[id].promise : null
    }

    set (id, promise, config) {
        Object.assign(this.cache, {[id]: {promise, config}})
    }

    getGroupedIds (id) {
        const groupedIds = []
        for (let requestId in this.cache) {
            const isInclude = groupedIds.includes(requestId)
            const isMatch = this.cache[requestId].config.requestGroup.includes(id)
            if (!isInclude && isMatch) {
                groupedIds.push(requestId)
            }
        }
        return groupedIds
    }

    getDeleteIds (id) {
        const deleteIds = this.getGroupedIds(id)
        if (this.cache.hasOwnProperty(id)) {
            deleteIds.push(id)
        }
        return deleteIds
    }

    delete (deleteIds) {
        let requestIds = []
        if (typeof deleteIds === 'undefined') {
            requestIds = Object.keys(this.cache)
        } else if (deleteIds instanceof Array) {
            deleteIds.forEach(id => {
                requestIds = [...requestIds, ...this.getDeleteIds(id)]
            })
        } else {
            requestIds = this.getDeleteIds(deleteIds)
        }
        requestIds = [...new Set(requestIds)]
        requestIds.forEach(requestId => {
            delete this.cache[requestId]
        })
        return Promise.resolve(deleteIds)
    }
}
