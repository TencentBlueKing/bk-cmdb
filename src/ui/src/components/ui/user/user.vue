<template>
    <component :is="tag">{{localText}}</component>
</template>

<script>
    import userQueue from './user-queue.js'
    export default {
        name: 'cmdb-user',
        props: {
            tag: {
                type: String,
                default: 'span'
            },
            text: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                localText: ''
            }
        },
        watch: {
            text: {
                immediate: true,
                handler (value, oldValue) {
                    if (value === oldValue) return
                    else if ([undefined, null, '', '--'].includes(value)) {
                        this.localText = '--'
                        return
                    }
                    this.localText = value
                    this.setUserText()
                }
            }
        },
        methods: {
            setUserText () {
                const user = this.localText
                const userInfo = userQueue.getUserInfo(user)
                if (userInfo) {
                    this.localText = userInfo
                } else {
                    userQueue.addUser({
                        user,
                        node: this
                    })
                }
            },
            updateUserText (user) {
                this.localText = user || this.text || '--'
            }
        }
    }
</script>
