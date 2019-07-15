<template>
    <div class="service-wrapper">
        <div class="options">
            <cmdb-form-bool class="options-checkall"
                :size="16"
                :checked="isCheckAll"
                :title="$t('Common[\'全选本页\']')"
                @change="handleCheckALL">
            </cmdb-form-bool>
            <bk-button class="ml10" :disabled="!checked.length" @click="batchDelete(!checked.length)">{{$t("BusinessTopology['批量删除']")}}</bk-button>
            <div class="option-right fr">
                <cmdb-form-bool class="options-checkbox"
                    :size="16"
                    :checked="isExpandAll"
                    @change="handleExpandAll">
                    <span class="checkbox-label">{{$t('Common["全部展开"]')}}</span>
                </cmdb-form-bool>
                <cmdb-form-singlechar class="options-search"
                    :placeholder="$t('BusinessTopology[\'请输入实例名称搜索\']')"
                    v-model="filter"
                    @keydown.native.enter="handleSearch">
                    <i class="bk-icon icon-close"
                        v-show="filter.length"
                        @click="handleClearFilter">
                    </i>
                    <i class="bk-icon icon-search" @click.stop="handleSearch"></i>
                </cmdb-form-singlechar>
            </div>
        </div>
        <div class="tables">
            <service-instance-table
                v-for="(instance, index) in instances"
                ref="serviceInstanceTable"
                :key="instance.id"
                :instance="instance"
                :expanded="index === 0"
                @delete-instance="handleDeleteInstance"
                @check-change="handleCheckChange">
            </service-instance-table>
            <cmdb-table
                v-if="!instances.length"
                :list="[]"
                :empty-height="42"
                :sortable="false"
                :reference-document-height="false">
            </cmdb-table>
        </div>
        <bk-paging class="pagination"
            pagination-able
            location="left"
            :cur-page="pagination.current"
            :total-page="pagination.totalPage"
            :pagination-count="pagination.size"
            @page-change="handlePageChange"
            @pagination-change="handleSizeChange">
        </bk-paging>
    </div>
</template>

<script>
    import { mapState } from 'vuex'
    import serviceInstanceTable from './service-instance-table.vue'
    export default {
        components: {
            serviceInstanceTable
        },
        data () {
            return {
                pagination: {
                    current: 1,
                    totalPage: 0,
                    size: 10
                },
                checked: [],
                isExpandAll: false,
                isCheckAll: false,
                filter: '',
                instances: []
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            host () {
                return this.info.host || {}
            }
        },
        created () {
            this.getHostSeriveInstances()
        },
        methods: {
            async getHostSeriveInstances () {
                try {
                    const data = await this.$store.dispatch('serviceInstance/getHostServiceInstances', {
                        params: this.$injectMetadata({
                            page: {
                                start: (this.pagination.current - 1) * this.pagination.size,
                                limit: this.pagination.size
                            },
                            bk_host_id: this.host.bk_host_id
                        })
                    })
                    this.checked = []
                    this.isCheckAll = false
                    this.isExpandAll = false
                    this.pagination.totalPage = Math.ceil(data.count / this.pagination.size)
                    this.instances = data.info
                } catch (e) {
                    console.error(e)
                    this.instances = []
                }
            },
            handleDeleteInstance (id) {
                this.instances = this.instances.filter(instance => instance.id !== id)
            },
            handleCheckALL (checked) {
                this.filter = ''
                this.isCheckAll = checked
                this.$refs.serviceInstanceTable.forEach(table => {
                    table.checked = checked
                })
            },
            handleExpandAll (expanded) {
                this.filter = ''
                this.isExpandAll = expanded
                this.$refs.serviceInstanceTable.forEach(table => {
                    table.localExpanded = expanded
                })
            },
            batchDelete (disabled) {
                if (disabled) {
                    return false
                }
                this.$bkInfo({
                    title: this.$t('BusinessTopology["确认删除实例"]'),
                    content: this.$t('BusinessTopology["即将删除选中的实例"]', { count: this.checked.length }),
                    confirmFn: async () => {
                        try {
                            const serviceInstanceIds = this.checked.map(instance => instance.id)
                            await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                                config: {
                                    data: this.$injectMetadata({
                                        service_instance_ids: serviceInstanceIds
                                    }),
                                    requestId: 'batchDeleteServiceInstance'
                                }
                            })
                            this.instances = this.instances.filter(instance => !serviceInstanceIds.includes(instance.id))
                            this.checked = []
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            handleCheckChange (checked, instance) {
                if (checked) {
                    this.checked.push(instance)
                } else {
                    this.checked = this.checked.filter(target => target.id !== instance.id)
                }
            },
            handleClearFilter () {
                this.filter = ''
                this.handleSearch()
            },
            handleSearch () {
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.pagination.current = page
                this.getHostSeriveInstances()
            },
            handleSizeChange (size) {
                this.pagination.current = 1
                this.pagination.size = size
                this.getHostSeriveInstances()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .service-wrapper {
        padding-top: 12px;
    }
    .options-checkall {
        width: 36px;
        height: 32px;
        line-height: 30px;
        padding: 0 9px;
        text-align: center;
        border: 1px solid #f0f1f5;
        border-radius: 2px;
    }
    .options-checkbox {
        margin: 0 15px 0 0;
        .checkbox-label {
            padding: 0 0 0 4px;
        }
    }
    .options-search {
        @include inlineBlock;
        position: relative;
        width: 240px;
        .icon-search {
            position: absolute;
            top: 9px;
            right: 9px;
            font-size: 14px;
            cursor: pointer;
        }
        .icon-close {
            position: absolute;
            top: 8px;
            right: 30px;
            width: 16px;
            height: 16px;
            line-height: 16px;
            border-radius: 50%;
            text-align: center;
            background-color: #ddd;
            color: #fff;
            font-size: 12px;
            transition: backgroundColor .2s linear;
            cursor: pointer;
            &:before {
                display: block;
                transform: scale(.7);
            }
            &:hover {
                background-color: #ccc;
            }
        }
        /deep/ {
            .cmdb-form-input {
                height: 32px;
                line-height: 30px;
                padding-right: 50px;
            }
        }
    }
    .tables {
        padding-top: 16px;
    }
</style>
