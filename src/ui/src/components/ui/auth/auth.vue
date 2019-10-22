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
    import temp from './auth-queue'
    export default {
        name: 'cmdb-auth',
        props: {
            authResource: {
                type: Object,
                default: () => ({})
            },
            bizId: {
                type: Number,
                default: null
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
            },
            resourceId () {
                return this.authResource.id || 0
            }
        },
        watch: {
            authResource () {
                console.log(1)
                // temp.addQueue({
                //     id: `${this.resource}-${this.resourceId}`,
                //     component: this,
                //     bizId: this.bizId,
                //     resource: this.authResource
                // })
            }
        },
        created () {
            temp.addQueue({
                id: `${this.resource}-${this.resourceId}`,
                component: this,
                bizId: this.bizId,
                resource: this.authResource
            })
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
