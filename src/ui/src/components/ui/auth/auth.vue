<template>
    <component class="auth-box"
        :is="tag"
        v-cursor="{
            active: !isAuthorized,
            auth: resources
        }"
        :class="{ disabled }"
        @click="handleClick">
        <slot :disabled="disabled"></slot>
    </component>
</template>

<script>
    import resourceOperation, { deepEqual } from './auth-queue'
    export default {
        name: 'cmdb-auth',
        props: {
            auth: {
                type: Object,
                default: () => ({})
            },
            tag: {
                type: String,
                default: 'span'
            },
            trigger: {
                type: String,
                default: 'initial',
                validator (trigger) {
                    return ['initial', 'click'].includes(trigger)
                }
            }
        },
        data () {
            return {
                pending: false,
                isAuthorized: true,
                disabled: true,
                turnOnVerify: window.CMDB_CONFIG.site.authscheme === 'iam'
            }
        },
        computed: {
            useIAM () {
                return this.turnOnVerify && Object.keys(this.auth).length
            },
            resources () {
                if (!this.auth.type) return []
                const types = Array.isArray(this.auth.type) ? this.auth.type : [this.auth.type]
                return types.map(type => {
                    return {
                        ...this.auth,
                        type: type
                    }
                })
            }
        },
        watch: {
            auth: {
                immediate: true,
                deep: true,
                handler (value, oldValue) {
                    if (!this.useIAM) {
                        this.disabled = false
                        this.$emit('update-auth', true)
                    } else if (this.trigger === 'initial' && !deepEqual(value, oldValue)) {
                        this.setup()
                    }
                }
            }
        },
        methods: {
            setup () {
                this.pending = true
                resourceOperation.pushQueue({
                    component: this,
                    data: this.auth
                })
            },
            updateAuth (auths) {
                const isPass = auths.every(auth => auth.is_pass)
                this.isAuthorized = isPass
                this.disabled = !isPass
                this.pending = false
                this.$emit('update-auth', isPass)
            },
            async handleClick () {
                if (this.disabled) {
                    return
                }
                if (this.useIAM && this.trigger === 'click') {
                    const dynamicResult = await new Promise(resolve => {
                        this.setup()
                        const unwatch = this.$watch(this.pending, () => {
                            if (!this.isAuthorized) {
                                this.$el.__cursor__.globalCallback({ auth: this.auth })
                            }
                            unwatch()
                            resolve(this.isAuthorized)
                        })
                    })
                    if (!dynamicResult) {
                        return
                    }
                }
                this.$emit('click')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .auth-box {
        display: inline-block;
    }
</style>
