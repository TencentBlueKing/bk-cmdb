<template>
    <bk-dialog class="bk-dialog-no-padding edit-label-dialog"
        v-model="isShow"
        :width="580"
        :mask-close="false"
        :esc-close="false"
        @after-leave="handleHidden">
        <div slot="header">
            {{$t('编辑标签')}}
        </div>
        <label-dialog-content ref="labelComp"
            :default-labels="labels">
        </label-dialog-content>
        <div class="edit-label-dialog-footer" slot="footer">
            <bk-button theme="primary" @click.stop="handleSubmit">{{$t('确定')}}</bk-button>
            <bk-button theme="default" class="ml5" @click.stop="close">{{$t('取消')}}</bk-button>
        </div>
    </bk-dialog>
</template>

<script>
    import LabelDialogContent from './label-dialog-content'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            LabelDialogContent
        },
        props: {
            serviceInstance: Object,
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
                if (!this.serviceInstance.labels) {
                    return []
                }
                return Object.entries(this.serviceInstance.labels).map(([key, value], index) => ({ id: index, key, value }))
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
                    const labelComp = this.$refs.labelComp
                    const validateResult = await labelComp.$validator.validateAll()
                    if (!validateResult) {
                        return false
                    }
                    const list = labelComp.submitList
                    const removeKeysList = labelComp.removeKeysList
                    const originList = labelComp.originList
                    const hasChange = JSON.stringify(labelComp.list) !== JSON.stringify(originList)

                    const request = []
                    if (removeKeysList.length) {
                        request.push(this.$store.dispatch('instanceLabel/deleteInstanceLabel', {
                            config: {
                                data: {
                                    bk_biz_id: this.bizId,
                                    instance_ids: [this.serviceInstance.id],
                                    keys: removeKeysList
                                },
                                requestId: this.request.delete
                            }
                        }))
                    }

                    if (list.length && hasChange) {
                        const labelSet = {}
                        list.forEach(label => {
                            labelSet[label.key] = label.value
                        })
                        request.push(this.$store.dispatch('instanceLabel/createInstanceLabel', {
                            params: {
                                bk_biz_id: this.bizId,
                                instance_ids: [this.serviceInstance.id],
                                labels: labelSet
                            },
                            config: {
                                requestId: this.request.create
                            }
                        }))
                    }
                    if (request.length) {
                        await Promise.all(request)
                        this.updateCallback && this.updateCallback(list)
                        this.$success(this.$t('保存成功'))
                    }
                    this.isShow = false
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
