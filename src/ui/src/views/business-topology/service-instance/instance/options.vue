<template>
    <div class="options">
        <div class="left">
            <cmdb-auth :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }">
                <bk-button slot-scope="{ disabled }" theme="primary"
                    :disabled="disabled"
                    @click="handleCreate">
                    {{$t('新增')}}
                </bk-button>
            </cmdb-auth>
            <bk-dropdown-menu class="ml10" trigger="click" font-size="medium">
                <bk-button slot="dropdown-trigger">
                    {{$t('实例操作')}}
                    <i class="bk-icon icon-angle-down"></i>
                </bk-button>
                <ul class="menu-list" slot="dropdown-content">
                    <cmdb-auth tag="li" class="menu-item"
                        :auth="{ type: $OPERATION.D_SERVICE_INSTANCE, relation: [bizId] }">
                        <span class="menu-option" slot-scope="{ disabled: authDisabled }"
                            :class="{ disabled: authDisabled || !selection.length }"
                            @click="handleBatchDelete(authDisabled || !selection.length)">
                            {{$t('批量删除')}}
                        </span>
                    </cmdb-auth>
                    <li class="menu-item">
                        <span :class="{ 'menu-option': true, disabled: !selection.length }" @click="handleCopy(!selection.length)">
                            {{$t('复制IP')}}
                        </span>
                    </li>
                    <cmdb-auth tag="li" class="menu-item"
                        :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
                        <span class="menu-option" slot-scope="{ disabled: authDisabled }"
                            :class="{ disabled: authDisabled || !selection.length }"
                            @click="handleBatchEditLabels(authDisabled || !selection.length)">
                            {{$t('编辑标签')}}
                        </span>
                    </cmdb-auth>
                </ul>
            </bk-dropdown-menu>
            <cmdb-auth class="options-sync" v-if="withTemplate"
                :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
                <bk-button slot-scope="{ disabled: authDisabled }"
                    :disabled="authDisabled || !hasDifference"
                    @click="handleSyncTemplate">
                    <span class="sync-wrapper">
                        <i class="bk-icon icon-refresh"></i>
                        {{$t('同步模板')}}
                    </span>
                    <span class="topo-status" v-show="hasDifference"></span>
                </bk-button>
            </cmdb-auth>
        </div>
        <div class="right">
            <bk-checkbox class="options-expand-all" v-model="allExpanded" @change="handleExpandAll">{{$t('全部展开')}}</bk-checkbox>
            <bk-search-select class="options-search ml10"
                ref="searchSelect"
                :show-condition="false"
                :placeholder="$t('请输入实例名称或选择标签')"
                :data="searchMenuList"
                v-model="searchValue"
                @change="handleSearch">
            </bk-search-select>
            <view-switcher class="ml10" active="instance"></view-switcher>
        </div>
    </div>
</template>

<script>
    import ViewSwitcher from '../common/view-switcher'
    import Bus from '../common/bus'
    import RouterQuery from '@/router/query'
    import { MENU_BUSINESS_DELETE_SERVICE } from '@/dictionary/menu-symbol'
    import { mapGetters } from 'vuex'
    import { Validator } from 'vee-validate'
    import { MULTIPLE_IP_REGEXP } from '@/dictionary/regexp.js'
    export default {
        components: {
            ViewSwitcher
        },
        data () {
            return {
                selection: [],
                hasDifference: false,
                allExpanded: false,
                historyLabels: {},
                searchValue: [],
                request: {
                    label: Symbol('label')
                }
            }
        },
        computed: {
            ...mapGetters('businessHost', ['selectedNode']),
            ...mapGetters('objectBiz', ['bizId']),
            withTemplate () {
                return this.selectedNode && this.selectedNode.data.service_template_id
            },
            nameFilterIndex () {
                return this.searchValue.findIndex(data => data.id === 'name')
            },
            searchMenuList () {
                const hasHistoryLables = Object.keys(this.historyLabels).length > 0
                const list = [{
                    id: 'name',
                    name: this.$t('服务实例名')
                }, {
                    id: 'tagValue',
                    name: this.$t('标签值'),
                    conditions: Object.keys(this.historyLabels).map(key => ({ id: key, name: key + ':' })),
                    disabled: !hasHistoryLables
                }, {
                    id: 'tagKey',
                    name: this.$t('标签键'),
                    disabled: !hasHistoryLables
                }]
                if (this.nameFilterIndex > -1) {
                    return list.slice(1)
                }
                return list.slice(0)
            }
        },
        watch: {
            withTemplate: {
                immediate: true,
                handler (withTemplate) {
                    withTemplate && this.checkDifference()
                }
            },
            searchMenuList () {
                this.$nextTick(() => {
                    const menu = this.$refs.searchSelect && this.$refs.searchSelect.menuInstance
                    menu && (menu.list = this.$refs.searchSelect.data)
                })
            },
            selectedNode () {
                this.searchValue = []
                RouterQuery.set({
                    instanceName: ''
                })
                Bus.$emit('filter-change', this.searchValue)
            }
        },
        created () {
            this.unwatch = RouterQuery.watch(['node', 'page'], () => {
                this.allExpanded = false
            })
            Bus.$on('instance-selection-change', this.handleInstanceSelectionChange)
            Bus.$on('update-labels', this.updateHistoryLabels)
            this.setFilter()
            this.updateHistoryLabels()
        },
        beforeDestroy () {
            this.unwatch()
            Bus.$off('instance-selection-change', this.handleInstanceSelectionChange)
            Bus.$off('update-labels', this.updateHistoryLabels)
        },
        methods: {
            async updateHistoryLabels () {
                try {
                    this.historyLabels = await this.$store.dispatch('instanceLabel/getHistoryLabel', {
                        params: {
                            bk_biz_id: this.bizId
                        }
                    })
                } catch (error) {
                    console.error(error)
                }
            },
            setFilter () {
                const filterName = RouterQuery.get('instanceName')
                if (!filterName) {
                    return false
                }
                this.searchValue.push({
                    'id': 'name',
                    'name': this.$t('服务实例名'),
                    'values': [{
                        'id': filterName,
                        'name': filterName
                    }]
                })
                this.$nextTick(() => {
                    // list.vue注册的监听晚于派发，因此nextTick再触发
                    Bus.$emit('filter-change', this.searchValue)
                })
            },
            handleInstanceSelectionChange (selection) {
                this.selection = selection
            },
            handleCreate () {
                this.$routerActions.redirect({
                    name: 'createServiceInstance',
                    params: {
                        moduleId: this.selectedNode.data.bk_inst_id,
                        setId: this.selectedNode.parent.data.bk_inst_id
                    },
                    query: {
                        title: this.selectedNode.data.bk_inst_name,
                        node: this.selectedNode.id,
                        tab: 'serviceInstance'
                    },
                    history: true
                })
            },
            async handleBatchDelete (disabled) {
                if (disabled) {
                    return false
                }
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_DELETE_SERVICE,
                    params: {
                        ids: this.selection.map(row => row.id).join('/'),
                        moduleId: this.selectedNode.data.bk_inst_id
                    },
                    history: true
                })
            },
            async handleCopy (disabled) {
                if (disabled) {
                    return false
                }
                try {
                    const validator = new Validator()
                    const validPromise = []
                    this.selection.forEach(row => {
                        const ip = row.name.split('_')[0]
                        validPromise.push(new Promise(async resolve => {
                            const { valid } = await validator.verify(ip, MULTIPLE_IP_REGEXP)
                            resolve({ valid, ip })
                        }))
                    })
                    const results = await Promise.all(validPromise)
                    const validResult = results.filter(result => result.valid).map(result => result.ip)
                    const unique = [...new Set(validResult)]
                    if (unique.length) {
                        await this.$copyText(unique.join('\n'))
                        this.$success(this.$t('复制成功'))
                    } else {
                        this.$warn(this.$t('暂无可复制的IP'))
                    }
                } catch (e) {
                    console.error(e)
                    this.$error(this.$t('复制失败'))
                }
            },
            handleEidtTag (disabled) {
                if (disabled) {
                    return false
                }
            },
            async checkDifference () {
                try {
                    const data = await this.$store.dispatch('businessSynchronous/searchServiceInstanceDifferences', {
                        params: {
                            bk_biz_id: this.bizId,
                            bk_module_ids: [this.selectedNode.data.bk_inst_id],
                            service_template_id: this.selectedNode.data.service_template_id
                        },
                        config: {
                            cancelPrevious: true
                        }
                    })
                    const difference = data.find(difference => difference.bk_module_id === this.selectedNode.data.bk_inst_id)
                    this.hasDifference = !!difference && difference.has_difference
                } catch (error) {
                    console.error(error)
                }
            },
            handleSyncTemplate () {
                this.$routerActions.redirect({
                    name: 'syncServiceFromModule',
                    params: {
                        modules: String(this.selectedNode.data.bk_inst_id),
                        template: this.selectedNode.data.service_template_id
                    },
                    history: true
                })
            },
            handleSearch (filters) {
                const transformedFilters = []
                filters.forEach(data => {
                    if (!data.values) {
                        const nameIndex = transformedFilters.findIndex(filter => filter.id === 'name')
                        if (nameIndex > -1) {
                            transformedFilters[nameIndex].values = [{ ...data }]
                        } else {
                            transformedFilters.push({
                                id: 'name',
                                name: this.$t('服务实例名'),
                                values: [{ ...data }]
                            })
                        }
                    } else if (data.id === 'tagValue') {
                        const [{ name }] = data.values
                        if (data.condition && name) {
                            transformedFilters.push(data)
                        } else {
                            this.$warn(this.$t('服务实例标签值搜索提示语'))
                        }
                    } else {
                        transformedFilters.push(data)
                    }
                })
                this.$nextTick(() => {
                    this.searchValue = transformedFilters
                    Bus.$emit('filter-change', this.searchValue)
                })
            },
            handleExpandAll (expanded) {
                Bus.$emit('expand-all-change', expanded)
            },
            handleBatchEditLabels (disabled) {
                if (disabled) {
                    return false
                }
                Bus.$emit('batch-edit-labels')
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
    .menu-list {
        .menu-item {
            line-height: 32px;
            .menu-option {
                display: block;
                padding: 0 10px;
                font-size: 14px;
                cursor: pointer;
                &:hover {
                    color: $primaryColor;
                }
                &.disabled {
                    color: $textDisabledColor;
                    cursor: not-allowed;;
                }
            }
        }
    }
    .options-sync {
        display: inline-block;
        position: relative;
        margin-left: 18px;
        padding: 0;
        &::before {
            content: '';
            position: absolute;
            top: 7px;
            left: -11px;;
            width: 1px;
            height: 20px;
            background-color: #dcdee5;
        }
        .sync-wrapper {
            display: flex;
            align-items: center;
        }
        .icon-refresh {
            font-size: 12px;
            margin-right: 4px;
        }
        .topo-status {
            position: absolute;
            top: -4px;
            right: -4px;
            width: 8px;
            height: 8px;
            background-color: #ea3636;
            border-radius: 50%;
        }
    }
    .options-search {
        width: 200px;
    }
</style>
