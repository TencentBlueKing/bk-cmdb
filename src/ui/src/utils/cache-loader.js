import { before } from '@/router'

const cacheMap = new Map()

before(() => cacheMap.clear())

export default {
    cacheMap: cacheMap,
    use: (id, handler) => {
        const exsit = cacheMap.has(id)
        if (exsit) {
            return cacheMap.get(id)
        }
        try {
            const result = (async () => {
                return handler()
            })()
            cacheMap.set(id, result)
            return result
        } catch (e) {
            console.error(e)
            return e
        }
    }
}
