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
                required: true
            },
            tag: {
                type: String,
                default: 'span'
            }
        },
        data () {
            return {
                isAuthorized: true,
                disabled: true,
                turnOnVerify: window.Site.authscheme === 'iam'
            }
        },
        computed: {
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
                    if (!this.turnOnVerify || !Object.keys(this.auth).length) {
                        this.disabled = false
                        this.$emit('update-auth', true)
                    } else if (!deepEqual(value, oldValue)) {
                        resourceOperation.pushQueue({
                            component: this,
                            data: this.auth
                        })
                    }
                }
            }
        },
        methods: {
            updateAuth (auths) {
                const isPass = auths.every(auth => auth.is_pass)
                this.isAuthorized = isPass
                this.disabled = !isPass
                this.$emit('update-auth', isPass)
            },
            handleClick () {
                if (this.disabled) {
                    return
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
