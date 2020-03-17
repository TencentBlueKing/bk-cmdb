<template>
    <div v-bkloading="{ isLoading: $loading(request.findOneTask) }">
        <bk-tab class="details-tab" :active.sync="active" type="unborder-card" slot="content">
            <bk-tab-panel name="details" :label="$t('任务详情')">
            </bk-tab-panel>
            <bk-tab-panel name="history" :label="$t('录入历史')">
            </bk-tab-panel>
        </bk-tab>
        <keep-alive>
            <component class="details-component" :is="component"
                :container="this"
                :task="task">
            </component>
        </keep-alive>
    </div>
</template>

<script>
    import TaskDetailsInfo from './task-details-info.vue'
    import TaskDetailsHistory from './task-details-history.vue'
    import TaskForm from './task-form.vue'
    export default {
        name: 'task-details',
        components: {
            [TaskDetailsHistory.name]: TaskDetailsHistory,
            [TaskDetailsInfo.name]: TaskDetailsInfo,
            [TaskForm.name]: TaskForm
        },
        props: {
            id: {
                type: Number,
                required: true
            }
        },
        data () {
            return {
                task: null,
                title: '',
                active: 'details',
                detailsComponent: TaskDetailsInfo.name,
                request: {
                    findOneTask: Symbol('findOneTask')
                }
            }
        },
        computed: {
            component () {
                if (!this.task) {
                    return null
                }
                if (this.active === 'details') {
                    return this.detailsComponent
                }
                return TaskDetailsHistory.name
            }
        },
        created () {
            this.getTaskDetails()
        },
        methods: {
            async getTaskDetails () {
                try {
                    this.task = await this.$store.dispatch('cloud/resource/findOneTask', {
                        id: this.id,
                        config: {
                            requestId: this.request.findOneTask
                        }
                    })
                } catch (e) {
                    console.error(e)
                    this.task = null
                }
            },
            show (options) {
                this.title = options.title || this.title
                this.detailsComponent = options.detailsComponent || TaskDetailsInfo.name
            },
            hide (eventType) {
                eventType && this.$emit(eventType)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-tab {
        height: auto;
        /deep/ {
            .bk-tab-header {
                padding: 0;
                margin: 0 30px;
            }
            .bk-tab-section {
                height: 0;
            }
        }
    }
    .details-component {
        height: calc(100% - 58px);
        @include scrollbar-y;
    }
</style>
