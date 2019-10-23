<template>
    <span class="auth-box"
        v-cursor="{
            active: isAuthorized,
            auth: [resource]
        }">
        <slot :disabled="disabled"></slot>
    </span>
</template>

<script>
    import temp, { deepEqual } from './auth-queue'
    export default {
        name: 'cmdb-auth',
        props: {
            authResource: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                isAuthorized: false,
                disabled: true
            }
        },
        computed: {
            resource () {
                return this.authResource.type || ''
            }
        },
        watch: {
            authResource: {
                immediate: true,
                handler (value, oldValue) {
                    if (!deepEqual(value, oldValue)) {
                        temp.pushQueue({
                            component: this,
                            data: this.authResource
                        })
                    }
                }
            }
        },
        methods: {
            updateAuth (auth) {
                const isPass = auth.is_pass
                this.isAuthorized = !isPass
                this.disabled = !isPass
            }
        }
    }
</script>

<style lang="scss" scoped>
    .auth-box {
        display: inline-block;
    }
</style>
