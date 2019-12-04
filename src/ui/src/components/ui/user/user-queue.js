import Vue from 'vue'
import debounce from 'lodash.debounce'
import $http from '@/api'

export default new Vue({
    data () {
        return {
            userMap: {},
            updateList: [],
            requestList: [],
            requestUserInfo: null
        }
    },
    watch: {
        requestList (value) {
            if (!value.length) return
            this.requestUserInfo()
        }
    },
    created () {
        this.requestUserInfo = debounce(this.searchUserInfo, 20)
    },
    methods: {
        getUserInfo (user) {
            const userList = user.split(',')
            const userInfo = []
            let exist = true
            for (const name of userList) {
                if (!this.userMap[name]) {
                    exist = false
                    return
                }
                userInfo.push(this.userMap[name])
            }
            if (!exist) return false
            return userInfo.join(',')
        },
        addUser (data) {
            this.updateList.push(data)
            const userList = data.user.split(',')
            userList.forEach(user => {
                const exist = this.userMap[user] || this.requestList.includes(user)
                if (!exist) {
                    this.requestList.push(user)
                }
            })
        },
        async searchUserInfo () {
            const requestList = [...this.requestList]
            const updateList = [...this.updateList]
            this.requestList = []
            this.updateList = []
            try {
                const data = await $http.get(`${window.API_HOST}user/detail?exact_lookups=${requestList.join(',')}`, {
                    globalError: false
                })
                data.users.forEach(user => {
                    this.userMap[user.english_name] = `${user.english_name}(${user.chinese_name})`
                })
                for (const instance of updateList) {
                    const nameList = instance.user.split(',')
                    const userInfo = nameList.map(name => this.userMap[name] || name)
                    const user = userInfo.join(',')
                    if (instance.node instanceof Vue) {
                        instance.node.updateUserText(user)
                    } else {
                        instance.node.textContent = user || instance.user || '--'
                        instance.options.title && (instance.node.title = user || instance.user || '--')
                    }
                }
            } catch (_) {
                for (const instance of updateList) {
                    const user = instance.user
                    if (instance.node instanceof Vue) {
                        instance.node.updateUserText(user)
                    } else {
                        instance.node.textContent = user || instance.user || '--'
                        instance.options.title && (instance.node.title = user || instance.user || '--')
                    }
                }
            }
        }
    }
})
