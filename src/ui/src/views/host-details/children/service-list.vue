<template>
    <div class="service-wrapper">
        <div class="options">
            <bk-checkbox class="options-checkall"
                :size="16"
                v-model="isCheckAll"
                :disabled="!instances.length"
                :title="$t('全选本页')"
                @change="handleCheckALL">
            </bk-checkbox>
            <cmdb-auth :auth="HOST_AUTH.D_SERVICE_INSTANCE">
                <bk-button slot-scope="{ disabled }"
                    class="ml10"
                    :disabled="disabled || !checked.length"
                    @click="batchDelete(!checked.length)">
                    {{$t('批量删除')}}
                </bk-button>
            </cmdb-auth>
            <div class="option-right fr">
                <bk-checkbox class="options-checkbox"
                    :size="16"
                    :disabled="!instances.length"
                    v-model="isExpandAll"
                    @change="handleExpandAll">
                    <span class="checkbox-label">{{$t('全部展开')}}</span>
                </bk-checkbox>
                <div class="options-search">
                    <bk-search-select
                        ref="searchSelect"
                        :show-condition="false"
                        :placeholder="$t('请输入实例名称或选择标签')"
                        :data="searchSelect"
                        v-model="searchSelectData"
                        @menu-child-condition-select="handleConditionSelect"
                        @change="handleSearch">
                    </bk-search-select>
                </div>
                <cmdb-switcher-group v-model="currentView" tips-key="host_service_list_view_tips" :tips="$t('标签或路径切换')">
                    <cmdb-switcher-item name="path" :tips="$t('显示拓扑')">
                        <i class="icon-cc-instance-path"></i>
                    </cmdb-switcher-item>
                    <cmdb-switcher-item name="label" :tips="$t('显示标签')">
                        <i class="icon-cc-label"></i>
                    </cmdb-switcher-item>
                </cmdb-switcher-group>
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
            <cmdb-table-empty
                slot="empty"
                :stuff="emptyStuff">
                <div class="empty-text">
                    <p>{{$t('暂无服务实例')}}，<span @click="handleGoAddInstance">{{$t('去业务拓扑添加')}}</span></p>
                </div>
            </cmdb-table-empty>
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
    import {
        MENU_BUSINESS_HOST_AND_SERVICE,
        MENU_BUSINESS_DELETE_SERVICE
    } from '@/dictionary/menu-symbol'
    import { mapState } from 'vuex'
    import serviceInstanceTable from './service-instance-table.vue'
    import authMixin from '../mixin-auth'
    import CmdbSwitcherGroup from '@/components/switcher/switcher-group'
    import CmdbSwitcherItem from '@/components/switcher/switcher-item'
    export default {
        components: {
            serviceInstanceTable,
            CmdbSwitcherGroup,
            CmdbSwitcherItem
        },
        mixins: [authMixin],
        data () {
            return {
                searchSelect: [
                    {
                        name: this.$t('服务实例名'),
                        id: 0
                    },
                    {
                        name: this.$t('标签值'),
                        id: 1,
                        children: [{
                            id: '',
                            name: ''
                        }],
                        conditions: []
                    },
                    {
                        name: this.$t('标签键'),
                        id: 2,
                        children: [{
                            id: '',
                            name: ''
                        }]
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
                currentView: 'path',
                historyLabels: {}
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            host () {
                return this.info.host || {}
            },
            emptyStuff () {
                return {
                    type: this.searchSelectData.length ? 'search' : 'default',
                    payload: {}
                }
            }
        },
        watch: {
            checked () {
                this.isCheckAll = (this.checked.length === this.instances.length) && this.checked.length !== 0
            }
        },
        created () {
            this.getHostSeriveInstances()
            this.getHistoryLabel()
        },
        methods: {
            async getHostSeriveInstances () {
                try {
                    const searchKey = this.searchSelectData.find(item => (item.id === 0 && item.hasOwnProperty('values'))
                        || (![0, 1].includes(item.id) && !item.hasOwnProperty('values')))
                    const data = await this.$store.dispatch('serviceInstance/getHostServiceInstances', {
                        params: {
                            page: {
                                start: (this.pagination.current - 1) * this.pagination.size,
                                limit: this.pagination.size
                            },
                            bk_host_id: this.host.bk_host_id,
                            bk_biz_id: this.info.biz[0].bk_biz_id,
                            search_key: searchKey
                                ? searchKey.hasOwnProperty('values') ? searchKey.values[0].name : searchKey.name
                                : '',
                            selectors: this.getSelectorParams()
                        }
                    })
                    if (data.count && !data.info.length) {
                        this.pagination.current -= 1
                        this.getHostSeriveInstances()
                    }
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
            getSelectorParams () {
                try {
                    const labels = this.searchSelectData.filter(item => item.id === 1 && item.hasOwnProperty('values'))
                    const labelsKey = this.searchSelectData.filter(item => item.id === 2 && item.hasOwnProperty('values'))
                    const submitLabel = {}
                    const submitLabelKey = {}
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
                    labelsKey.forEach(label => {
                        const id = label.values[0].id
                        if (!submitLabelKey[id]) {
                            submitLabelKey[id] = id
                        }
                    })
                    const selectors = Object.keys(submitLabel).map(key => {
                        return {
                            key: key,
                            operator: 'in',
                            values: submitLabel[key]
                        }
                    })
                    const selectorsKey = Object.keys(submitLabelKey).map(key => {
                        return {
                            key: key,
                            operator: 'exists',
                            values: []
                        }
                    })
                    return selectors.concat(selectorsKey)
                } catch (e) {
                    console.error(e)
                    return []
                }
            },
            async getHistoryLabel () {
                const historyLabels = await this.$store.dispatch('instanceLabel/getHistoryLabel', {
                    params: {
                        bk_biz_id: this.info.biz[0].bk_biz_id
                    },
                    config: {
                        requestId: 'getHistoryLabel',
                        cancelPrevious: true
                    }
                })
                this.historyLabels = historyLabels
                const keys = Object.keys(historyLabels)
                const valueOption = keys.map(key => {
                    return {
                        name: key + ' : ',
                        id: key
                    }
                })
                const keyOption = keys.map(key => {
                    return {
                        name: key,
                        id: key
                    }
                })
                if (!valueOption.length) {
                    this.$set(this.searchSelect[1], 'disabled', true)
                }
                if (!keyOption.length) {
                    this.$set(this.searchSelect[2], 'disabled', true)
                }
                this.$set(this.searchSelect[1], 'conditions', valueOption)
                this.$set(this.searchSelect[2], 'children', keyOption)
            },
            handleDeleteInstance (id) {
                this.getHostSeriveInstances()
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
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_DELETE_SERVICE,
                    params: {
                        ids: this.checked.map(instance => instance.id).join('/')
                    },
                    history: true
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
                const [biz] = this.info.biz
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    params: {
                        bizId: biz.bk_biz_id
                    },
                    query: {
                        node: `biz-${biz.bk_biz_id}`,
                        ip: this.info.host.bk_host_innerip
                    }
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
            color: #c4c6cc;
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
        background-color: #699df4 !important;
        .tippy-arrow {
            border-bottom-color: #699df4 !important;
        }
    }
</style>
