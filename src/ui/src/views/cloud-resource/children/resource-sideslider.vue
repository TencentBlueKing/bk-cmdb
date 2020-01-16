<template>
    <bk-sideslider v-transfer-dom
        :is-show.sync="isShow"
        :title="title"
        :width="695"
        @hidden="handleHidden">
        <component slot="content" ref="component"
            class="slider-content"
            :is="component"
            :container="this"
            v-bind="componentProps">
        </component>
    </bk-sideslider>
</template>

<script>
    import ResourceForm from './resource-form.vue'
    import ResourceDetails from './resource-details.vue'
    export default {
        components: {
            [ResourceForm.name]: ResourceForm,
            [ResourceDetails.name]: ResourceDetails
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
                    this.component = ResourceForm.name
                } else if (options.type === 'details') {
                    this.component = ResourceDetails.name
                }
                this.title = options.title
                this.isShow = true
            },
            hide (eventType) {
                this.isShow = false
                eventType && this.$emit(eventType)
            },
            handleHidden () {
                this.component = null
                this.componentProps = {}
            }
        }
    }
</script>

<style lang="scss" scoped>
    .slider-content {
        height: 100%;
        @include scrollbar-y;
    }
</style>
