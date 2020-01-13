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
            for (const name of userList) {
                if (!this.userMap[name]) return
                userInfo.push(this.userMap[name])
            }
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
                const users = await $http.get(`${window.API_HOST}user/list`, {
                    params: {
                        exact_lookups: requestList.join(',')
                    },
                    globalError: false
                })
                users.forEach(user => {
                    if (user.chinese_name) {
                        this.userMap[user.english_name] = `${user.english_name}(${user.chinese_name})`
                    } else {
                        this.userMap[user.english_name] = user.english_name
                    }
                })
            } catch (_) {}
            for (const instance of updateList) {
                const nameList = instance.user.split(',')
                const userInfo = nameList.map(name => this.userMap[name] || name)
                const user = userInfo.join(',')
                if (instance.node instanceof Vue) {
                    instance.node.updateUserText(user)
                } else {
                    instance.node.textContent = user || '--'
                    instance.options.title && (instance.node.title = user || '--')
                }
            }
        }
    }
})
