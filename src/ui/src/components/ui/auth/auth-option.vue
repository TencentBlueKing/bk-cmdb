<template>
    <bk-option
        :id="id"
        :name="name"
        :disabled="computedDisabled">
        <cmdb-auth class="auth-option"
            :auth="auth"
            :ignore="ignore"
            :title="name"
            @update-auth="handleUpdateAuth">
            {{name}}
        </cmdb-auth>
    </bk-option>
</template>

<script>
    export default {
        name: 'cmdb-auth-option',
        props: {
            id: [String, Number],
            name: [String, Number],
            disabled: Boolean,
            auth: [Object, Array],
            ignore: Boolean
        },
        data () {
            return {
                isAuthorized: false
            }
        },
        computed: {
            computedDisabled () {
                return this.disabled || !this.isAuthorized
            }
        },
        methods: {
            handleUpdateAuth (isAuthorized) {
                this.isAuthorized = isAuthorized
            }
        }
    }
</script>

<style lang="scss" scoped>
    .auth-option {
        display: block;
        @include ellipsis;
    }
</style>
