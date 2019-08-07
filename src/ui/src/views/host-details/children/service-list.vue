<template>
    <div class="service-wrapper">
        <div class="options">
            <cmdb-form-bool class="options-checkall"
                :size="16"
                :checked="isCheckAll"
                :title="$t('Common[\'全选本页\']')"
                @change="handleCheckALL">
            </cmdb-form-bool>
            <span style="display: inline-block;"
                v-cursor="{
                    active: !$isAuthorized($OPERATION.D_SERVICE_INSTANCE),
                    auth: [$OPERATION.D_SERVICE_INSTANCE]
                }">
                <bk-button class="ml10"
                    :disabled="!$isAuthorized($OPERATION.D_SERVICE_INSTANCE) || !checked.length"
                    @click="batchDelete(!checked.length)">
                    {{$t("BusinessTopology['批量删除']")}}
                </bk-button>
            </span>
            <div class="option-right fr">
                <cmdb-form-bool class="options-checkbox"
                    :size="16"
                    :checked="isExpandAll"
                    @change="handleExpandAll">
                    <span class="checkbox-label">{{$t('Common["全部展开"]')}}</span>
                </cmdb-form-bool>
                <div class="options-search">
                    <bk-search-select
                        :show-condition="false"
                        :placeholder="$t('BusinessTopology[\'实例名称/标签\']')"
                        :data="searchSelect"
                        v-model="searchSelectData"
                        @change="handleSearch">
                    </bk-search-select>
                </div>
                <bk-popover
                    ref="popoverCheckView"
                    :always="true"
                    :width="210"
                    theme="check-view-color"
                    placement="bottom-end">
                    <div slot="content" class="popover-main">
                        <span>这里可以切换显示标签或路径</span>
                        <i class="bk-icon icon-close" @click="handleCheckViewStatus"></i>
                    </div>
                    <div class="options-check-view">
                        <i :class="['icon-cc-label', 'view-btn', 'pr10', { 'active': currentView === 'label' }]"
                            @click="checkView('label')"></i>
                        <span class="dividing-line"></span>
                        <i :class="['icon-cc-instance-path', 'view-btn', 'pl10', { 'active': currentView === 'path' }]"
                            @click="checkView('path')"></i>
                    </div>
                </bk-popover>
            </div>
        </div>
        <div class="tables">
            <service-instance-table
                v-for="(instance, index) in instances"
                ref="serviceInstanceTable"
                :key="instance.id"
                :instance="instance"
                :expanded="index === 0"
                :current-view="currentView"
                @delete-instance="handleDeleteInstance"
                @check-change="handleCheckChange">
            </service-instance-table>
        </div>
        <bk-table v-if="!instances.length" :data="[]"></bk-table>
        <bk-pagination class="pagination"
            align="right"
            size="small"
            :current="pagination.current"
            :count="pagination.totalPage"
            :limit="pagination.size"
            @change="handlePageChange"
            @limit-change="handleSizeChange">
        </bk-pagination>
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
                searchSelect: [
                    {
                        name: this.$t('BusinessTopology["服务实例名"]'),
                        id: 0
                    },
                    {
                        name: this.$t('BusinessTopology["标签"]'),
                        id: 1,
                        multiable: true,
                        children: []
                    }
                ],
                searchSelectData: [],
                pagination: {
                    current: 1,
                    totalPage: 0,
                    size: 10
                },
                checked: [],
                isExpandAll: false,
                isCheckAll: false,
                filter: [],
                instances: [],
                currentView: 'label',
                checkViewTipsStatus: this.$store.getters['featureTipsParams'].hostServiceInstanceCheckView
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
            this.getHistoryLabel()
        },
        mounted () {
            if (this.checkViewTipsStatus) {
                this.$refs.popoverCheckView.instance.show()
            } else {
                this.$refs.popoverCheckView.instance.destroy()
            }
        },
        methods: {
            async getHostSeriveInstances () {
                try {
                    const searchKey = this.searchSelectData.find(item => (item.id === 0 && item.hasOwnProperty('values'))
                        || (![0, 1].includes(item.id) && !item.hasOwnProperty('values')))
                    const labels = this.searchSelectData.filter(item => item.id === 1 && item.hasOwnProperty('values'))
                    const labelKeys = []
                    labels.forEach(label => {
                        label.values.forEach(value => {
                            if (!labelKeys.includes(value.name)) {
                                labelKeys.push(value.name)
                            }
                        })
                    })
                    const selectors = labelKeys.map(key => {
                        return {
                            key: key,
                            operator: 'exists',
                            values: []
                        }
                    })
                    const data = await this.$store.dispatch('serviceInstance/getHostServiceInstances', {
                        params: this.$injectMetadata({
                            page: {
                                start: (this.pagination.current - 1) * this.pagination.size,
                                limit: this.pagination.size
                            },
                            bk_host_id: this.host.bk_host_id,
                            search_key: searchKey
                                ? searchKey.hasOwnProperty('values') ? searchKey.values[0].name : searchKey.name
                                : '',
                            selectors: selectors
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
            async getHistoryLabel () {
                const historyLabels = await this.$store.dispatch('instanceLabel/getHistoryLabel', {
                    params: this.$injectMetadata({}),
                    config: {
                        requestId: 'getHistoryLabel',
                        cancelPrevious: true
                    }
                })
                const keys = Object.keys(historyLabels).map(key => {
                    return {
                        name: key,
                        id: key
                    }
                })
                this.$set(this.searchSelect[1], 'children', keys)
            },
            handleDeleteInstance (id) {
                this.instances = this.instances.filter(instance => instance.id !== id)
            },
            handleCheckALL (checked) {
                this.searchSelectData = []
                this.isCheckAll = checked
                this.$refs.serviceInstanceTable.forEach(table => {
                    table.checked = checked
                })
            },
            handleExpandAll (expanded) {
                this.searchSelectData = []
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
                    subTitle: this.$t('BusinessTopology["即将删除选中的实例"]', { count: this.checked.length }),
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
                this.handleSearch()
            },
            handleSearch (value) {
                const instanceName = this.searchSelectData.filter(item => (item.id === 0 && item.hasOwnProperty('values'))
                    || (![0, 1].includes(item.id) && !item.hasOwnProperty('values')))
                if (instanceName.length) {
                    this.searchSelect[0].id === 0 && this.searchSelect.shift()
                } else {
                    this.searchSelect[0].id !== 0 && this.searchSelect.unshift({
                        name: this.$t('BusinessTopology["服务实例名"]'),
                        id: 0
                    })
                }
                if (instanceName.length >= 2) {
                    this.searchSelectData.pop()
                    this.$bkMessage({
                        message: this.$t('BusinessTopology["服务实例名重复"]'),
                        theme: 'warning'
                    })
                    return
                }
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
            },
            checkView (value) {
                this.currentView = value
            },
            handleCheckViewStatus () {
                this.$store.commit('setFeatureTipsParams', 'hostServiceInstanceCheckView')
                this.$refs.popoverCheckView.instance.hide()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .service-wrapper {
        padding: 14px 0 0 0;
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
        min-width: 240px;
        max-width: 500px;
        z-index: 99;
    }
    .popover-main {
        @include inlineBlock;
        .bk-icon {
            margin: -2px 0 0 10px;
            cursor: pointer;
        }
    }
    .options-check-view {
        @include inlineBlock;
        height: 32px;
        line-height: 30px;
        font-size: 0;
        border: 1px solid #c4c6cc;
        padding: 0 10px;
        border-radius: 2px;
        .view-btn {
            color: #dcdee5;
            font-size: 14px;
            height: 100%;
            line-height: 30px;
            cursor: pointer;
            &:hover, &.active {
                color: #3a84ff;
            }
        }
        .dividing-line {
            @include inlineBlock;
            width: 1px;
            height: 14px;
            background-color: #dcdee5;
        }
    }
    .tables {
        padding-top: 16px;
    }
</style>

<style lang="scss">
    .check-view-color-theme {
        padding: 10px !important;
        background-color: #699df4;
        .tippy-arrow {
            border-bottom-color: #699df4 !important;
        }
    }
</style>
