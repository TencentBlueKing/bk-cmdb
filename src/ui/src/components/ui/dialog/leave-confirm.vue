<template>
    <bk-dialog></bk-dialog>
</template>

<script>
    export default {
        name: 'cmdb-leave-confirm',
        props: {
            active: Boolean,
            title: {
                type: String,
                default: ''
            },
            content: {
                type: String,
                default: ''
            },
            okText: {
                type: String,
                default: this.$t('确认')
            },
            cancelText: {
                type: String,
                default: this.$t('取消')
            },
            triggers: {
                type: Array,
                default () {}
            }
        },
        mounted () {
            this.addListener()
        },
        beforeDestory () {
            this.removeListener()
        },
        methods: {
            addListener () {
                window.addEventListener('beforeunload', this.unloadHandler)
                this.$router.beforeHooks.unshift(this.beforeEachHook)
            },
            removeListener () {
                window.removeEventListener('beforeunload', this.unloadHandler)
                const beforeEachHookIndex = this.$router.beforeHooks.indexOf(this.beforeEachHook)
                beforeEachHookIndex > -1 && this.$router.beforeHooks.splice(beforeEachHookIndex, 1)
            },
            unloadHandler (e) {
                if (this.active) {
                    return (e || window.event).returnValue = this.title
                }
            },
            async beforeEachHook (to, from, next) {
                if (this.active) {
                    const result = await this.promise
                    result ? next() : next(false)
                }
                next()
            }
        }
    }
</script>