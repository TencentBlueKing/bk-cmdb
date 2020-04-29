<template>
    <div class="hosts-table-layout">
        <div class="hosts-options">
            <div class="options-left">
                <slot name="options-left"></slot>
            </div>
            <div class="options-right clearfix">
                <div class="fl" v-if="showScope">
                    <i class="options-split"></i>
                    <bk-select style="width: 280px;"
                        v-model="scope"
                        font-size="medium"
                        :clearable="false">
                        <bk-option id="all" :name="$t('全部主机')"></bk-option>
                        <bk-option id="0" :name="$t('已分配主机')"></bk-option>
                        <bk-option id="1" :name="$t('未分配主机')"></bk-option>
                    </bk-select>
                </div>
                <div class="fr">
                    <cmdb-host-filter class="ml10"
                        ref="hostFilter"
                        :properties="filterProperties"
                        :show-scope="showScope">
                    </cmdb-host-filter>
                    <icon-button class="ml10"
                        icon="icon icon-cc-setting"
                        v-bk-tooltips.top="$t('列表显示属性配置')"
                        @click="handleColumnConfigClick">
                    </icon-button>
                    <icon-button class="ml10" v-if="showHistory"
                        v-bk-tooltips="$t('查看删除历史')"
                        icon="icon icon-cc-history"
                        @click="routeToHistory">
                    </icon-button>
                </div>
            </div>
        </div>
        <bk-table class="hosts-table bkc-white"
            v-bkloading="{ isLoading: $loading(['searchHosts', 'batchSearchProperties', 'post_searchGroup_host']) }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 190"
            @selection-change="handleSelectionChange"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange">
            <bk-table-column type="selection" width="60" align="center" fixed class-name="bk-table-selection"></bk-table-column>
            <bk-table-column v-for="column in table.header"
                min-width="80"
                :key="column.id"
                :label="column.name"
                :sortable="column.sortable ? 'custom' : false"
                :prop="column.id"
                show-overflow-tooltip>
                <template slot-scope="{ row }">
                    <cmdb-property-value
                        :theme="column.id === 'bk_host_id' ? 'primary' : 'default'"
                        :value="row | hostValueFilter(column.objId, column.id, column)"
                        :show-unit="false"
                        :property="column.type"
                        :options="getPropertyValue(column.objId, column.id, 'option')"
                        @click.native.stop="handleValueClick(row, column)">
                    </cmdb-property-value>
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="slider.show"
            :title="slider.title"
            :width="800"
            :before-close="handleSliderBeforeClose">
            <bk-tab :active.sync="tab.active" type="unborder-card" slot="content" v-if="slider.show">
                <bk-tab-panel name="attribute" :label="$t('属性')" style="width: calc(100% + 40px);margin: 0 -20px;">
                    <cmdb-form-multiple v-if="tab.attribute.type === 'multiple'"
                        ref="multipleForm"
                        :properties="properties.host"
                        :property-groups="propertyGroups"
                        :object-unique="objectUnique"
                        @on-submit="handleMultipleSave"
                        @on-cancel="handleSliderBeforeClose">
                    </cmdb-form-multiple>
                </bk-tab-panel>
            </bk-tab>
        </bk-sideslider>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="columnsConfig.show"
            :width="600"
            :title="$t('列表显示属性配置')">
            <cmdb-columns-config slot="content"
                v-if="columnsConfig.show"
                :properties="columnsConfigProperties"
                :selected="columnsConfig.selected"
                :disabled-columns="columnsConfig.disabledColumns"
                @on-apply="handleApplyColumnsConfig"
                @on-cancel="columnsConfig.show = false"
                @on-reset="handleResetColumnsConfig">
            </cmdb-columns-config>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapGetters, mapActions, mapState } from 'vuex'
    import cmdbColumnsConfig from '@/components/columns-config/columns-config'
    import cmdbHostFilter from '@/components/hosts/filter/index.vue'
    import hostValueFilter from '@/filters/host'
    import {
        MENU_BUSINESS,
        MENU_BUSINESS_HOST_DETAILS,
        MENU_RESOURCE_HOST_DETAILS,
        MENU_RESOURCE_BUSINESS_HOST_DETAILS
    } from '@/dictionary/menu-symbol'
    import RouterQuery from '@/router/query'
    import { getIPPayload } from '@/utils/host'
    export default {
        components: {
            cmdbColumnsConfig,
            cmdbHostFilter
        },
        filters: {
            hostValueFilter
        },
        props: {
            columnsConfigProperties: {
                type: Array,
                required: true
            },
            columnsConfigKey: {
                type: String,
                required: true
            },
            columnsConfigDisabledColumns: {
                type: Array,
                default () {
                    return ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id', 'bk_module_name']
                }
            },
            saveAuth: {
                type: [String, Array],
                default: ''
            },
            editAuth: {
                type: [String, Array],
                default: ''
            },
            deleteAuth: {
                type: [String, Array],
                default: ''
            },
            showHistory: Boolean,
            showScope: Boolean
        },
        data () {
            return {
                objectUnique: [],
                properties: {
                    biz: [],
                    host: [],
                    set: [],
                    module: []
                },
                propertyGroups: [],
                table: {
                    checked: [],
                    header: Array(8).fill({}),
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    sort: 'bk_host_id',
                    exportUrl: `${window.API_HOST}hosts/export`,
                    stuff: {
                        type: 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                slider: {
                    show: false,
                    title: ''
                },
                tab: {
                    active: 'attribute',
                    attribute: {
                        type: 'details',
                        inst: {
                            details: {},
                            edit: {},
                            original: {}
                        }
                    }
                },
                columnsConfig: {
                    show: false,
                    selected: [],
                    disabledColumns: ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id', 'bk_module_name', 'bk_set_name']
                },
                scope: RouterQuery.get('scope', '1')
            }
        },
        computed: {
            ...mapState('userCustom', ['globalUsercustom']),
            ...mapGetters(['supplierAccount']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapState('hosts', ['condition']),
            customColumns () {
                return this.usercustom[this.columnsConfigKey] || []
            },
            globalCustomColumns () {
                return this.globalUsercustom['host_global_custom_table_columns'] || []
            },
            clipboardList () {
                return this.table.header.filter(header => header.type !== 'checkbox')
            },
            selectedHosts () {
                return this.table.list.filter(host => this.table.checked.includes(host['host']['bk_host_id']))
            },
            filterProperties () {
                const { module, set, host } = this.properties
                const filterProperty = ['bk_host_id', 'bk_host_innerip', 'bk_host_outerip']
                return {
                    host: host.filter(property => !filterProperty.includes(property.bk_property_id)),
                    module,
                    set
                }
            }
        },
        watch: {
            'table.checked' (checked) {
                this.$emit('on-checked', checked)
            },
            'table.header' (header) {
                this.$emit('on-set-header', header)
            },
            'slider.show' (show) {
                if (!show) {
                    this.tab.active = 'attribute'
                }
            },
            customColumns () {
                this.setTableHeader()
            },
            columnsConfigProperties () {
                this.setTableHeader()
            },
            scope (value) {
                RouterQuery.set({
                    scope: value,
                    ip: '',
                    page: 1
                })
            },
            condition () {
                RouterQuery.set('_t', Date.now())
            }
        },
        async created () {
            try {
                RouterQuery.watch(['ip', 'scope', 'exact', 'page', 'limit', 'condition', '_t'], ({
                    scope = '1',
                    page = 1,
                    limit = this.table.pagination.limit
                }) => {
                    this.scope = scope
                    this.table.pagination.current = parseInt(page)
                    this.table.pagination.limit = parseInt(limit)
                    this.getHostList()
                }, { immediate: true, throttle: 16 })
                await Promise.all([
                    this.getProperties(),
                    this.getHostPropertyGroups()
                ])
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('hostUpdate', ['updateHost']),
            ...mapActions('hostSearch', ['searchHost']),
            getPropertyValue (modelId, propertyId, field) {
                const model = this.properties[modelId]
                if (!model) {
                    return ''
                }
                const curProperty = model.find(property => property.bk_property_id === propertyId)
                return curProperty ? curProperty[field] : ''
            },
            getProperties () {
                return this.batchSearchObjectAttribute({
                    injectId: 'host',
                    params: this.$injectMetadata({
                        bk_obj_id: { '$in': Object.keys(this.properties) },
                        bk_supplier_account: this.supplierAccount
                    }, { inject: this.$route.name !== 'resource' }),
                    config: {
                        requestId: 'batchSearchProperties'
                    }
                }).then(result => {
                    Object.keys(this.properties).forEach(objId => {
                        this.properties[objId] = result[objId]
                    })
                    return result
                })
            },
            getHostPropertyGroups () {
                return this.searchGroup({
                    objId: 'host',
                    params: this.$injectMetadata(),
                    config: {
                        requestId: 'post_searchGroup_host'
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            setTableHeader () {
                const customColumns = this.customColumns.length ? this.customColumns : this.globalCustomColumns
                const properties = this.$tools.getHeaderProperties(this.columnsConfigProperties, customColumns, this.columnsConfigDisabledColumns)
                this.table.header = properties.map(property => {
                    return {
                        id: property.bk_property_id,
                        name: this.$tools.getHeaderPropertyName(property),
                        type: property.bk_property_type,
                        objId: property.bk_obj_id,
                        sortable: property.bk_obj_id === 'host' && !['foreignkey'].includes(property.bk_property_type)
                    }
                })
                this.columnsConfig.selected = properties.map(property => property['bk_property_id'])
            },
            getHostCellText (header, item) {
                const objId = header.objId
                const propertyId = header.id
                const headerProperty = this.$tools.getProperty(this.properties[objId], propertyId)
                const originalValues = item[objId] instanceof Array ? item[objId] : [item[objId]]
                const text = []
                originalValues.forEach(value => {
                    const flattenedText = this.$tools.getPropertyText(headerProperty, value)
                    flattenedText ? text.push(flattenedText) : void (0)
                })
                return text.join(',') || '--'
            },
            getHostList (event) {
                try {
                    this.searchHost({
                        params: this.injectScope({
                            bk_biz_id: -1,
                            condition: this.condition,
                            ip: getIPPayload(),
                            page: {
                                start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                                limit: this.table.pagination.limit,
                                sort: this.table.sort
                            }
                        }),
                        config: {
                            requestId: 'searchHosts',
                            cancelPrevious: true
                        }
                    }).then(data => {
                        this.table.pagination.count = data.count
                        this.table.list = data.info
    
                        if (event) {
                            this.table.stuff.type = 'search'
                        }
    
                        return data
                    }).catch(e => {
                        this.table.checked = []
                        this.table.list = []
                        this.table.pagination.count = 0
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            injectScope (params) {
                if (!this.showScope) {
                    return params
                }
                const biz = params.condition.find(condition => condition.bk_obj_id === 'biz')
                if (this.scope === 'all') {
                    biz.condition = biz.condition.filter(condition => condition.field !== 'default')
                } else {
                    const newCondition = {
                        field: 'default',
                        operator: '$eq',
                        value: parseInt(this.scope)
                    }
                    const existCondition = biz.condition.find(condition => condition.field === 'default')
                    if (existCondition) {
                        Object.assign(existCondition, newCondition)
                    } else {
                        biz.condition.push(newCondition)
                    }
                }
                return params
            },
            handlePageChange (current, event) {
                RouterQuery.set('page', current)
            },
            handleSizeChange (limit) {
                RouterQuery.set('limit', limit)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            },
            handleCopy (target) {
                const copyList = this.table.list.filter(item => {
                    return this.table.checked.includes(item['host']['bk_host_id'])
                })
                const copyText = []
                this.$tools.clone(copyList).forEach(item => {
                    const cellText = this.getHostCellText(target, item)
                    if (cellText !== '--') {
                        copyText.push(cellText)
                    }
                })
                if (copyText.length) {
                    this.$copyText(copyText.join('\n')).then(() => {
                        this.$success(this.$t('复制成功'))
                    }, () => {
                        this.$error(this.$t('复制失败'))
                    })
                } else {
                    this.$info(this.$t('该字段无可复制的值'))
                }
            },
            handleSelectionChange (selection) {
                this.table.checked = selection.map(item => item.host.bk_host_id)
            },
            handleValueClick (item, column) {
                if (column.objId !== 'host' || column.id !== 'bk_host_id') {
                    return false
                }
                const business = item.biz[0]
                if (this.$route.meta.owner === MENU_BUSINESS) {
                    this.$routerActions.redirect({
                        name: MENU_BUSINESS_HOST_DETAILS,
                        params: {
                            bizId: business.bk_biz_id,
                            business: business.bk_biz_id,
                            id: item.host.bk_host_id
                        },
                        history: true
                    })
                } else if (business.default) {
                    this.$routerActions.redirect({
                        name: MENU_RESOURCE_HOST_DETAILS,
                        params: {
                            id: item.host.bk_host_id
                        },
                        history: true
                    })
                } else {
                    this.$routerActions.redirect({
                        name: MENU_RESOURCE_BUSINESS_HOST_DETAILS,
                        params: {
                            bizId: business.bk_biz_id,
                            business: business.bk_biz_id,
                            id: item.host.bk_host_id
                        },
                        history: true
                    })
                }
            },
            batchUpdate (params) {
                return this.updateHost({ params }).then(data => {
                    this.$success(this.$t('保存成功'))
                    this.getHostList()
                    return data
                })
            },
            async handleMultipleEdit () {
                this.objectUnique = await this.$store.dispatch('objectUnique/searchObjectUniqueConstraints', {
                    objId: 'host',
                    params: this.$injectMetadata({}, {
                        inject: this.$route.name !== 'resource'
                    })
                })
                this.tab.attribute.type = 'multiple'
                this.slider.title = this.$t('主机属性')
                this.slider.show = true
            },
            async handleMultipleSave (changedValues) {
                await this.batchUpdate(this.$injectMetadata({
                    ...changedValues,
                    'bk_host_id': this.table.checked.join(',')
                }, { inject: this.$route.name !== 'resource' }))
                this.slider.show = false
            },
            handleMultipleCancel () {
                this.slider.show = false
            },
            handleApplyColumnsConfig (properties) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.columnsConfigKey]: properties.map(property => property['bk_property_id'])
                })
                this.columnsConfig.show = false
            },
            handleResetColumnsConfig () {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.columnsConfigKey]: []
                })
            },
            handleSliderBeforeClose () {
                if (this.tab.active === 'attribute' && this.tab.attribute.type !== 'details') {
                    const $form = this.tab.attribute.type === 'update' ? this.$refs.form : this.$refs.multipleForm
                    const changedValues = $form.changedValues
                    if (Object.keys(changedValues).length) {
                        return new Promise((resolve, reject) => {
                            this.$bkInfo({
                                title: this.$t('确认退出'),
                                subTitle: this.$t('退出会导致未保存信息丢失'),
                                extCls: 'bk-dialog-sub-header-center',
                                confirmFn: () => {
                                    this.slider.show = false
                                },
                                cancelFn: () => {
                                    resolve(false)
                                }
                            })
                        })
                    }
                    this.slider.show = false
                }
                this.slider.show = false
            },
            handleQuickSearch (property, value, operator) {
                this.$emit('on-quick-search', property, value, operator)
            },
            routeToHistory () {
                this.$routerActions.redirect({ name: 'hostHistory', history: true })
            },
            handleColumnConfigClick () {
                this.$refs.hostFilter.$refs.filterPopper.instance.hide()
                this.columnsConfig.show = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .hosts-options{
        font-size: 0;
        .options-left {
            float: left;
        }
        .options-right {
            overflow: hidden;
            .options-split {
                @include inlineBlock;
                width: 2px;
                height: 20px;
                margin: 0 10px;
                background-color: #DCDEE5;
            }
        }
        .icon-btn {
            width: 32px;
            padding: 0;
            line-height: 14px;
        }
        .options-button{
            position: relative;
            display: inline-block;
            vertical-align: middle;
            font-size: 14px;
            &:first-child {
                margin-left: 0;
            }
        }
    }
    .hosts-table {
        margin-top: 14px;
    }
</style>
