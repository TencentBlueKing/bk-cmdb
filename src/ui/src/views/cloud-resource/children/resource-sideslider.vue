<template>
    <bk-sideslider v-transfer-dom
        :is-show.sync="isShow"
        :title="title"
        :width="695"
        @hidden="handleHidden">
        <component slot="content" ref="component"
            class="slider-content"
            :is="component"
            :container="this">
        </component>
    </bk-sideslider>
</template>

<script>
    import ResourceForm from './resource-form.vue'
    export default {
        components: {
            [ResourceForm.name]: ResourceForm
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
                this.component = ResourceForm.name
                this.title = options.title
                this.isShow = true
            },
            hide (eventType) {
                this.isShow = false
                eventType && this.$emit(eventType)
            },
            handleHidden () {
                this.component = null
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
