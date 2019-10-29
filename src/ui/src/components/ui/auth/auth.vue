<template>
    <span class="auth-box"
        v-cursor="{
            active: isAuthorized,
            auth: resources
        }">
        <slot :disabled="disabled"></slot>
    </span>
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
            requestAuth: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                isAuthorized: false,
                disabled: true
            }
        },
        computed: {
            resources () {
                if (!this.auth.type) return []
                return Array.isArray(this.auth.type) ? this.auth.type : [this.auth.type]
            }
        },
        watch: {
            auth: {
                immediate: true,
                deep: true,
                handler (value, oldValue) {
                    if (this.requestAuth && !deepEqual(value, oldValue)) {
                        resourceOperation.pushQueue({
                            component: this,
                            data: this.auth
                        })
                    }
                }
            },
            requestAuth: {
                immediate: true,
                handler (value) {
                    if (!value) {
                        this.disabled = false
                    }
                }
            }
        },
        methods: {
            updateAuth (auths) {
                const passData = auths.map(auth => {
                    return auth.is_pass
                })
                const isPass = passData.every(pass => pass)
                this.isAuthorized = !isPass
                this.disabled = !isPass
                this.$emit('update-auth', isPass)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .auth-box {
        display: inline-block;
    }
</style>
