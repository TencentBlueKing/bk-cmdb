<template>
    <bk-sideslider v-transfer-dom
        :is-show.sync="isShow"
        :title="title"
        :width="900"
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
    import TaskForm from './task-form.vue'
    import TaskDetails from './task-details.vue'
    export default {
        components: {
            [TaskForm.name]: TaskForm,
            [TaskDetails.name]: TaskDetails
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
            show ({ mode, props, title }) {
                const componentMap = {
                    create: TaskForm.name,
                    update: TaskForm.name,
                    details: TaskDetails.name
                }
                this.component = componentMap[mode]
                this.componentProps = props
                this.title = title
                this.isShow = true
            },
            hide (eventType) {
                this.isShow = false
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
