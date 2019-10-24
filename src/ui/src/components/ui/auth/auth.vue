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
    import resourceOperation, { deepEqual } from './auth-queue'
    export default {
        name: 'cmdb-auth',
        props: {
            auth: {
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
                return this.auth.type || ''
            }
        },
        watch: {
            auth: {
                immediate: true,
                deep: true,
                handler (value, oldValue) {
                    if (!deepEqual(value, oldValue)) {
                        resourceOperation.pushQueue({
                            component: this,
                            data: this.auth
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
