import Vue from 'vue'
class HostStore {
    constructor () {
        this.hosts = []
    }

    get isSelected () {
        return !!this.hosts.length
    }
    get isSameBiz () {
        const bizSet = new Set()
        this.hosts.forEach(host => {
            const [biz] = host.biz
            bizSet.add(biz.bk_biz_id)
        })
        return bizSet.size === 1
    }
    get isAllIdleModule () {
        return this.hosts.every(host => {
            const [module] = host.module
            return module.default === 1
        })
    }
    setSelected (hosts = []) {
        this.hosts = hosts
    }

    getSelected () {
        return this.hosts
    }
}

export default Vue.observable(new HostStore())
