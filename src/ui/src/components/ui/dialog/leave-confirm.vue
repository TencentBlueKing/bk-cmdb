<template>
    <span hidden></span>
</template>

<script>
    import ConfirmStore from './confirm-store.js'
    export default {
        name: 'cmdb-leave-confirm',
        props: {
            id: {
                type: [String, Number, Symbol],
                required: true
            },
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
                default: ''
            },
            cancelText: {
                type: String,
                default: ''
            },
            reverse: Boolean
        },
        data () {
            return {
                visible: false,
                confirmPromise: Promise.resolve(true),
                confirmResolve: null
            }
        },
        mounted () {
            ConfirmStore.install(this)
            this.addListener()
        },
        beforeDestroy () {
            ConfirmStore.uninstall(this)
            this.removeListener()
        },
        methods: {
            show () {
                if (this.active) {
                    this.confirmPromise = new Promise(resolve => {
                        this.confirmResolve = resolve
                    })
                    this.$bkInfo({
                        title: this.title,
                        subHeader: this.$createElement('div', {
                            class: 'leave-confirm-content'
                        }, this.content),
                        okText: this.okText || this.$t('确认'),
                        cancelText: this.cancelText || this.$t('取消'),
                        closeIcon: false,
                        confirmFn: () => {
                            this.confirmResolve(this.reverse)
                        },
                        cancelFn: () => {
                            this.confirmResolve(!this.reverse)
                        }
                    })
                } else {
                    this.confirmPromise = Promise.resolve(true)
                }
            },
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
                    /* eslint-disable-next-line */
                    return (e || window.event).returnValue = this.title
                }
            },
            async beforeEachHook (to, from, next) {
                const result = await ConfirmStore.popup(this.id)
                result ? next() : next(false)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .leave-confirm-content {
        text-align: center;
        font-size: 14px;
    }
</style>
