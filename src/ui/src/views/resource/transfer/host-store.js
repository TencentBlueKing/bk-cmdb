import Vue from 'vue'
class HostStore {
    constructor () {
        this.hosts = []
        this.businessList = []
    }

    get isSelected () {
        return !!this.hosts.length
    }

    get bizSet () {
        const bizSet = new Set()
        this.hosts.forEach(host => {
            const [biz] = host.biz
            bizSet.add(biz.bk_biz_id)
        })
        return bizSet
    }

    get isSameBiz () {
        return this.bizSet.size === 1
    }
    get isAllIdleModule () {
        return this.hosts.every(host => {
            const [module] = host.module
            return module.default === 1
        })
    }

    get isAllIdleSet () {
        return this.hosts.every(host => {
            const [module] = host.module
            return module.default !== 0
        })
    }

    get uniqueBusiness () {
        if (this.isSameBiz) {
            const [bizId] = Array.from(this.bizSet)
            return this.businessList.find(business => business.bk_biz_id === bizId)
        }
        return null
    }

    clear () {
        this.hosts = []
    }

    setSelected (hosts = []) {
        this.hosts = hosts
    }

    getSelected () {
        return this.hosts
    }

    setBusinessList (businessList) {
        this.businessList = businessList
    }
}

export default Vue.observable(new HostStore())
