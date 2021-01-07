<template>
    <bk-dialog class="bk-dialog-no-padding"
        v-model="isShow"
        :mask-close="false"
        :width="580"
        @after-leave="handleHidden">
        <div class="reset-header" slot="header">
            {{$t('批量编辑')}}
            <span>{{$tc('已选择实例', serviceInstances.length, { num: serviceInstances.length })}}</span>
        </div>
        <batch-content ref="batchContent"
            :labels="labels">
            <single-content
                ref="singleContent"
                slot="batch-add-label"
                class="edit-label">
            </single-content>
        </batch-content>
        <div class="edit-label-footer" slot="footer">
            <bk-button theme="primary" :loading="$loading(Object.values(request))" @click.stop="handleSubmit">{{$t('确定')}}</bk-button>
            <bk-button theme="default" class="ml5" @click.stop="close">{{$t('取消')}}</bk-button>
        </div>
    </bk-dialog>
</template>

<script>
    import BatchContent from './label-batch-dialog-content'
    import SingleContent from './label-dialog-content'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            BatchContent,
            SingleContent
        },
        props: {
            serviceInstances: Array,
            updateCallback: Function
        },
        data () {
            return {
                isShow: false,
                request: {
                    delete: Symbol('delete'),
                    create: Symbol('create')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            labels () {
                const labels = []
                this.serviceInstances.forEach(instance => {
                    if (!instance.labels) return
                    Object.keys(instance.labels).forEach(key => {
                        const value = instance.labels[key]
                        const exist = labels.find(label => label.key === key && label.value === value)
                        if (exist) {
                            exist.instances.push(instance)
                        } else {
                            labels.push({
                                key,
                                value,
                                instances: [instance]
                            })
                        }
                    })
                })
                return labels
            }
        },
        methods: {
            show () {
                this.isShow = true
            },
            close () {
                this.isShow = false
            },
            async handleSubmit () {
                try {
                    const singleContent = this.$refs.singleContent
                    const validateResult = await singleContent.$validator.validateAll()
                    if (!validateResult) {
                        return false
                    }
                    const batchContent = this.$refs.batchContent
                    const removeLabels = batchContent.removeLabels

                    const removeIds = []
                    const removeKeys = []
                    if (removeLabels.length) {
                        removeLabels.forEach(data => {
                            removeIds.push(...data.instances.map(instance => instance.id))
                            removeKeys.push(data.key)
                        })
                        await this.$store.dispatch('instanceLabel/deleteInstanceLabel', {
                            config: {
                                data: {
                                    bk_biz_id: this.bizId,
                                    instance_ids: removeIds,
                                    keys: removeKeys
                                },
                                requestId: this.request.delete
                            }
                        })
                    }

                    const list = singleContent.submitList
                    const labelSet = {}
                    if (list.length) {
                        list.forEach(label => {
                            labelSet[label.key] = label.value
                        })
                        await this.$store.dispatch('instanceLabel/createInstanceLabel', {
                            params: {
                                bk_biz_id: this.bizId,
                                instance_ids: this.serviceInstances.map(instance => instance.id),
                                labels: labelSet
                            },
                            config: {
                                requestId: this.request.create
                            }
                        })
                    }
                    if (removeLabels.length || list.length) {
                        this.updateCallback && this.updateCallback(removeKeys, labelSet)
                        this.$success(this.$t('保存成功'))
                    }
                    this.close()
                } catch (error) {
                    console.error(error)
                }
            },
            handleHidden () {
                this.$emit('close')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .edit-label-dialog {
        /deep/ .bk-dialog-header {
            text-align: left !important;
            font-size: 24px;
            color: #444444;
            margin-top: -15px;
        }
        .edit-label-dialog-footer {
            .bk-button {
                min-width: 76px;
            }
        }
    }
</style>
