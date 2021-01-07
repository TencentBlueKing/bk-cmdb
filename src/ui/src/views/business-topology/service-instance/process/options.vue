<template>
    <div class="options">
        <div class="left">
            <cmdb-auth :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
                <bk-button slot-scope="{ disabled }" theme="primary"
                    :disabled="disabled || !selection.value.length"
                    @click="handleBatchEdit">
                    {{$t('编辑')}}
                </bk-button>
            </cmdb-auth>
            <cmdb-auth class="ml10" v-if="!serviceTemplateId"
                :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
                <bk-button slot-scope="{ disabled }" theme="default"
                    :disabled="disabled || !selection.value.length"
                    :loading="$loading(selection.requestId)"
                    @click="handleBatchtDelete">
                    {{$t('删除')}}
                </bk-button>
            </cmdb-auth>
        </div>
        <div class="right">
            <bk-checkbox class="options-expand-all" v-model="expandAll" @change="handleExpandAllChange">{{$t('全部展开')}}</bk-checkbox>
            <bk-input class="options-search ml10"
                ref="searchSelect"
                v-model.trim="searchValue"
                right-icon="bk-icon icon-search"
                clearable
                :max-width="200"
                :placeholder="$t('进程别名')"
                @enter="handleSearch"
                @clear="handleSearch">
            </bk-input>
            <view-switcher class="ml10" active="process"></view-switcher>
        </div>
    </div>
</template>

<script>
    import ViewSwitcher from '../common/view-switcher'
    import Bus from '../common/bus'
    import FormMultiple from '@/components/service/form/form-multiple.js'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            ViewSwitcher
        },
        data () {
            return {
                withTemplate: true,
                searchData: [],
                searchValue: '',
                expandAll: false,
                selection: {
                    process: null,
                    value: [],
                    requestId: null
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['selectedNode']),
            serviceTemplateId () {
                return this.selectedNode && this.selectedNode.data.service_template_id
            }
        },
        created () {
            Bus.$on('process-selection-change', this.handleProcessSelectionChange)
            Bus.$on('process-list-change', this.handleProcessListChange)
        },
        beforeDestroy () {
            Bus.$off('process-selection-change', this.handleProcessSelectionChange)
            Bus.$off('process-list-change', this.handleProcessListChange)
        },
        methods: {
            handleBatchEdit () {
                const props = {
                    submitHandler: this.batchEditSubmitHandler
                }
                if (this.serviceTemplateId) {
                    props.serviceTemplateId = this.serviceTemplateId
                    // 按服务模板创建的服务实例聚合出来的进程，其对应的服务模板都是同一个，因此取第一个的id即可
                    props.processTemplateId = this.selection.value[0].relation.process_template_id
                }
                FormMultiple.show(props)
            },
            async batchEditSubmitHandler (values) {
                try {
                    await this.$store.dispatch('serviceInstance/batchUpdateProcess', {
                        params: {
                            bk_biz_id: this.bizId,
                            process_ids: this.selection.value.map(process => process.process_id),
                            update_data: values
                        }
                    })
                    Bus.$emit('refresh-expand-list', this.selection.process)
                    this.$success('修改成功')
                } catch (error) {
                    console.error(error)
                }
            },
            handleBatchtDelete () {
                this.$bkInfo({
                    title: this.$t('确定删除N个进程', { count: this.selection.value.length }),
                    confirmFn: () => {
                        Bus.$emit('batch-delete', this.selection.process.bk_process_name)
                    }
                })
            },
            handleProcessSelectionChange (process, selection, requestId) {
                if (selection.length) {
                    this.selection.process = process
                    this.selection.value = selection
                    this.selection.requestId = requestId
                } else if (process === this.selection.process) {
                    this.selection.process = null
                    this.selection.value = []
                    this.selection.requestId = null
                }
            },
            handleSearch () {
                Bus.$emit('filter-list', this.searchValue)
            },
            handleExpandAllChange (expand) {
                Bus.$emit('expand-all-change', expand)
            },
            handleProcessListChange () {
                this.expandAll = false
                this.selection = {
                    process: null,
                    value: [],
                    requestId: null
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .options {
        display: flex;
        justify-content: space-between;
        flex-wrap: wrap;
        .left,
        .right {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }
    }
    .options-search {
        width: 300px;
    }
</style>
