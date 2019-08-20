<template>
    <div class="service-wrapper">
        <div class="options">
            <bk-checkbox class="options-checkall"
                :size="16"
                v-model="isCheckAll"
                :title="$t('全选本页')"
                @change="handleCheckALL">
            </bk-checkbox>
            <span style="display: inline-block;"
                v-cursor="{
                    active: !$isAuthorized($OPERATION.D_SERVICE_INSTANCE),
                    auth: [$OPERATION.D_SERVICE_INSTANCE]
                }">
                <bk-button class="ml10"
                    :disabled="!$isAuthorized($OPERATION.D_SERVICE_INSTANCE) || !checked.length"
                    @click="batchDelete(!checked.length)">
                    {{$t('批量删除')}}
                </bk-button>
            </span>
            <div class="option-right fr">
                <bk-checkbox class="options-checkbox"
                    :size="16"
                    v-model="isExpandAll"
                    @change="handleExpandAll">
                    <span class="checkbox-label">{{$t('全部展开')}}</span>
                </bk-checkbox>
                <div class="options-search">
                    <bk-search-select
                        ref="searchSelect"
                        :show-condition="false"
                        :placeholder="$t('实例名称/标签')"
                        :data="searchSelect"
                        v-model="searchSelectData"
                        @menu-child-condition-select="handleConditionSelect"
                        @change="handleSearch">
                    </bk-search-select>
                </div>
                <bk-popover
                    ref="popoverCheckView"
                    :always="true"
                    :width="224"
                    theme="check-view-color"
                    placement="bottom-end">
                    <div slot="content" class="popover-main">
                        <span>{{$t('标签或路径切换')}}</span>
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
        <bk-table v-if="!instances.length" :data="[]" class="mb10">
            <div slot="empty" class="empty-text">
                <img src="../../../assets/images/empty-content.png" alt="">
                <p>暂无服务实例，<span @click="handleGoAddInstance">跳转服务拓扑添加</span></p>
            </div>
        </bk-table>
        <bk-pagination v-if="instances.length"
            class="pagination"
            align="right"
            size="small"
            :current="pagination.current"
            :count="pagination.count"
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
                        name: this.$t('服务实例名'),
                        id: 0
                    },
                    {
                        name: this.$t('标签'),
                        id: 1,
                        multiable: true,
                        children: [{
                            id: '',
                            name: ''
                        }],
                        conditions: []
                    }
                ],
                searchSelectData: [],
                pagination: {
                    current: 1,
                    count: 0,
                    size: 10
                },
                checked: [],
                isExpandAll: false,
                isCheckAll: false,
                filter: [],
                instances: [],
                currentView: 'label',
                checkViewTipsStatus: this.$store.getters['featureTipsParams'].hostServiceInstanceCheckView,
                historyLabels: {}
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
                    const submitLabel = {}
                    labels.forEach(label => {
                        const conditionId = label.condition.id
                        if (!submitLabel[conditionId]) {
                            submitLabel[conditionId] = [label.values[0].id]
                        } else {
                            if (submitLabel[conditionId].indexOf(label.values[0].id) < 0) {
                                submitLabel[conditionId].push(label.values[0].id)
                            }
                        }
                    })
                    const selectors = Object.keys(submitLabel).map(key => {
                        return {
                            key: key,
                            operator: 'in',
                            values: submitLabel[key]
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
                    this.pagination.count = data.count
                    this.instances = data.info
                } catch (e) {
                    console.error(e)
                    this.instances = []
                    this.pagination.count = 0
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
                        name: key + ' : ',
                        id: key
                    }
                })
                this.historyLabels = historyLabels
                this.$set(this.searchSelect[1], 'conditions', keys)
            },
            handleDeleteInstance (id) {
                const filterInstances = this.instances.filter(instance => instance.id !== id)
                if (!filterInstances.length && this.pagination.current > 1) {
                    this.pagination.current -= 1
                    this.getHostSeriveInstances()
                }
                this.instances = filterInstances
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
                    title: this.$t('确认删除实例'),
                    subTitle: this.$t('即将删除选中的实例', { count: this.checked.length }),
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
                            const filterInstances = this.instances.filter(instance => !serviceInstanceIds.includes(instance.id))
                            if (!filterInstances.length && this.pagination.current > 1) {
                                this.pagination.current -= 1
                                this.getHostSeriveInstances()
                            }
                            this.instances = filterInstances
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
                        name: this.$t('服务实例名'),
                        id: 0
                    })
                }
                if (instanceName.length >= 2) {
                    this.searchSelectData.pop()
                    this.$bkMessage({
                        message: this.$t('服务实例名重复'),
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
            },
            handleConditionSelect (cur, index) {
                const values = this.historyLabels[cur.id]
                const children = values.map(item => {
                    return {
                        id: item,
                        name: item
                    }
                })
                const el = this.$refs.searchSelect
                el.curItem.children = children
                el.updateChildMenu(cur, index, false)
                el.showChildMenu(children)
            },
            handleGoAddInstance () {
                this.$router.replace({
                    name: 'topology'
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .service-wrapper {
        height: calc(100% - 30px);
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
        display: flex;
        justify-content: space-between;
        align-items: center;
        .bk-icon {
            margin: 0 0 0 10px;
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
        @include scrollbar-y;
        max-height: calc(100% - 80px);
        margin: 16px 0 10px;
    }
    .empty-text {
        text-align: center;
        p {
            font-size: 14px;
            margin: -10px 0 0;
        }
        span {
            color: #3a84ff;
            cursor: pointer;
        }
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
