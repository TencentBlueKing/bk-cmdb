<template>
    <bk-sideslider v-transfer-dom
        :is-show.sync="isShow"
        :title="title"
        :width="540"
        @hidden="handleHidden">
        <component slot="content" ref="component"
            :is="component"
            :container="this"
            v-bind="componentProps">
        </component>
    </bk-sideslider>
</template>

<script>
    import AccountForm from './account-form.vue'
    import AccountDetails from './account-details.vue'
    export default {
        components: {
            [AccountForm.name]: AccountForm,
            [AccountDetails.name]: AccountDetails
        },
        data () {
            return {
                isShow: false,
                title: '',
                component: null,
                componentProps: {}
            }
        },
        methods: {
            show (options) {
                this.componentProps = options.props || {}
                if (options.type === 'form') {
                    this.component = AccountForm.name
                } else if (options.type === 'details') {
                    this.component = AccountDetails.name
                }
                this.title = options.title
                this.isShow = true
            },
            hide () {
                this.isShow = false
            },
            handleHidden () {
                this.component = null
                this.componentProps = {}
            }
        }
    }
</script>
