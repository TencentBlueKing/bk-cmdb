<template>
    <component class="auth-box"
        :is="tag"
        v-cursor="{
            active: !isAuthorized,
            auth: auth,
            onclick
        }"
        :class="{ disabled }"
        @click="handleClick">
        <slot :disabled="disabled"></slot>
    </component>
</template>

<script>
    import AuthProxy from './auth-queue'
    import deepEqual from 'deep-equal'
    export default {
        name: 'cmdb-auth',
        props: {
            ignore: Boolean,
            auth: {
                type: [Object, Array]
            },
            tag: {
                type: String,
                default: 'span'
            },
            onclick: Function
        },
        data () {
            return {
                authResults: null,
                authMetas: null,
                isAuthorized: false,
                disabled: true,
                useIAM: window.CMDB_CONFIG.site.authscheme === 'iam'
            }
        },
        watch: {
            auth: {
                deep: true,
                handler (value, oldValue) {
                    !deepEqual(value, oldValue) && this.setAuthProxy()
                }
            },
            ignore () {
                this.setAuthProxy()
            }
        },
        mounted () {
            this.setAuthProxy()
        },
        methods: {
            setAuthProxy () {
                if (this.useIAM && this.auth && !this.ignore) {
                    AuthProxy.add({
                        component: this,
                        data: this.auth
                    })
                } else {
                    this.disabled = false
                    this.isAuthorized = true
                    this.$emit('update-auth', true)
                }
            },
            updateAuth (authResults, authMetas) {
                let isPass
                if (!authResults.length && authMetas.length) { // 鉴权失败
                    isPass = false
                } else {
                    isPass = authResults.every(result => result.is_pass)
                }
                this.authResults = authResults
                this.authMetas = authMetas
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
