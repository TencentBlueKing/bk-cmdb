<template>
    <div class="models-layout">
        <div class="models-options clearfix">
            <div class="options-button clearfix fl">
                <cmdb-auth class="fl mr10" :auth="{ type: $OPERATION.C_INST, relation: [model.id] }">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        :disabled="disabled"
                        @click="handleCreate">
                        {{$t('新建')}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth class="fl mr10"
                    :auth="[
                        { type: $OPERATION.C_INST, relation: [model.id] },
                        { type: $OPERATION.U_INST, relation: [model.id] }
                    ]">
                    <bk-button slot-scope="{ disabled }"
                        class="models-button"
                        :disabled="disabled"
                        @click="importSlider.show = true">
                        {{$t('导入')}}
                    </bk-button>
                </cmdb-auth>
                <div class="fl mr10">
                    <bk-button class="models-button" theme="default"
                        :disabled="!table.checked.length"
                        @click="handleExport">
                        {{$t('导出')}}
                    </bk-button>
                </div>
                <div class="fl mr10">
                    <bk-button class="models-button"
                        :disabled="!table.checked.length"
                        @click="handleMultipleEdit">
                        {{$t('批量更新')}}
                    </bk-button>
                </div>
                <bk-button class="models-button button-delete fl mr10"
                    hover-theme="danger"
                    :disabled="!table.checked.length"
                    @click="handleMultipleDelete">
                    {{$t('删除')}}
                </bk-button>
            </div>
            <div class="options-button fr">
                <icon-button class="ml5"
                    v-bk-tooltips="$t('查看删除历史')"
                    icon="icon-cc-history"
                    @click="routeToHistory">
                </icon-button>
                <icon-button class="ml5"
                    v-bk-tooltips="$t('列表显示属性配置')"
                    icon="icon-cc-setting"
                    @click="columnsConfig.show = true">
                </icon-button>
            </div>
            <div class="options-filter clearfix fr">
                <cmdb-property-selector class="filter-selector fl"
                    v-model="filter.id"
                    :properties="properties"
                    :object-unique="objectUnique">
                </cmdb-property-selector>
                <component class="filter-value fl"
                    v-if="['enum', 'list', 'organization'].includes(filterType)"
                    :is="`cmdb-form-${filterType}`"
                    :options="$tools.getEnumOptions(properties, filter.id)"
                    :allow-clear="true"
                    :clearable="true"
                    :auto-select="false"
                    v-model="filter.value"
                    font-size="medium"
                    @on-selected="handleFilterValueChange"
                    @on-checked="handleFilterValueChange">
                </component>
                <bk-input class="filter-value cmdb-form-input fl" type="text" maxlength="11"
                    v-else-if="filterType === 'int'"
                    v-model.number="filter.value"
                    clearable
                    right-icon="icon-search"
                    font-size="medium"
                    :placeholder="$t('快速查询')"
                    @enter="handleFilterValueChange"
                    @clear="handleFilterValueChange">
                </bk-input>
                <bk-input class="filter-value cmdb-form-input fl" type="text"
                    v-else-if="filterType === 'float'"
                    v-model.number="filter.value"
                    clearable
                    right-icon="icon-search"
                    font-size="medium"
                    :placeholder="$t('快速查询')"
                    @enter="handleFilterValueChange"
                    @clear="handleFilterValueChange">
                </bk-input>
                <bk-input class="filter-value cmdb-form-input fl" type="text"
                    v-else
                    v-model.trim="filter.value"
                    clearable
                    right-icon="icon-search"
                    font-size="medium"
                    :placeholder="$t('快速查询')"
                    @enter="handleFilterValueChange"
                    @clear="handleFilterValueChange">
                </bk-input>
            </div>
        </div>
        <bk-table class="models-table" ref="table"
            v-bkloading="{ isLoading: $loading('^=post_searchInst_') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 190"
            @sort-change="handleSortChange"
            @page-limit-change="handleSizeChange"
            @page-change="handlePageChange"
            @selection-change="handleSelectChange">
            <bk-table-column type="selection" width="60" align="center" fixed class-name="bk-table-selection"></bk-table-column>
            <bk-table-column v-for="column in table.header"
                sortable="custom"
                min-width="80"
                :key="column.id"
                :prop="column.id"
                :label="column.name"
                show-overflow-tooltip>
                <template slot-scope="{ row }">
                    <cmdb-property-value
                        :theme="column.id === 'bk_inst_id' ? 'primary' : 'default'"
                        :show-unit="false"
                        :value="row[column.id]"
                        :property="column.property"
                        @click.native.stop="handleValueClick(row, column)">
                    </cmdb-property-value>
                </template>
            </bk-table-column>
            <bk-table-column fixed="right" :label="$t('操作')">
                <template slot-scope="{ row }">
                    <cmdb-auth :auth="{ type: $OPERATION.D_INST, relation: [model.id, row.bk_inst_id] }">
                        <template slot-scope="{ disabled }">
                            <bk-button theme="primary" text :disabled="disabled" @click.stop="handleDelete(row)">{{$t('删除')}}</bk-button>
                        </template>
                    </cmdb-auth>
                </template>
            </bk-table-column>
            <cmdb-table-empty
                slot="empty"
                :auth="{ type: $OPERATION.C_INST, relation: [model.id] }"
                :stuff="table.stuff"
                @create="handleCreate">
            </cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="slider.show"
            :title="slider.title"
            :width="800"
            :before-close="handleSliderBeforeClose">
            <bk-tab type="unborder-card" slot="content"
                v-if="slider.contentShow"
                :active.sync="tab.active" :show-header="attribute.type !== 'create'">
                <bk-tab-panel name="attribute" :label="$t('属性')" style="width: calc(100% + 40px);margin: 0 -20px;">
                    <cmdb-form v-if="['update', 'create'].includes(attribute.type)"
                        ref="form"
                        :properties="properties"
                        :property-groups="propertyGroups"
                        :inst="attribute.inst.edit"
                        :type="attribute.type"
                        :save-auth="{ type: attribute.type === 'update' ? $OPERATION.U_INST : $OPERATION.C_INST }"
                        :object-unique="objectUnique"
                        @on-submit="handleSave"
                        @on-cancel="handleCancel">
                    </cmdb-form>
                    <cmdb-form-multiple v-else-if="attribute.type === 'multiple'"
                        ref="multipleForm"
                        :uneditable-properties="['bk_inst_name']"
                        :properties="properties"
                        :property-groups="propertyGroups"
                        :object-unique="objectUnique"
                        @on-submit="handleMultipleSave"
                        @on-cancel="handleMultipleCancel">
                    </cmdb-form-multiple>
                </bk-tab-panel>
            </bk-tab>
        </bk-sideslider>
        <bk-sideslider v-transfer-dom :is-show.sync="columnsConfig.show" :width="600" :title="$t('列表显示属性配置')">
            <cmdb-columns-config slot="content"
                v-if="columnsConfig.show"
                :properties="properties"
                :selected="columnsConfig.selected"
                :disabled-columns="columnsConfig.disabledColumns"
                @on-apply="handleApplyColumnsConfig"
                @on-cancel="columnsConfig.show = false"
                @on-reset="handleResetColumnsConfig">
            </cmdb-columns-config>
        </bk-sideslider>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="importSlider.show"
            :width="800"
            :title="$t('批量导入')">
            <cmdb-import v-if="importSlider.show" slot="content"
                :template-url="url.template"
                :import-url="url.import"
                @success="handlePageChange(1)"
                @partialSuccess="handlePageChange(1)">
            </cmdb-import>
        </bk-sideslider>
        <router-subview></router-subview>
    </div>
</template>

<script>
    import { mapState, mapGetters, mapActions } from 'vuex'
    import cmdbColumnsConfig from '@/components/columns-config/columns-config'
    import cmdbImport from '@/components/import/import'
    import { MENU_RESOURCE_INSTANCE_DETAILS } from '@/dictionary/menu-symbol'
    import cmdbPropertySelector from '@/components/property-selector'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            cmdbColumnsConfig,
            cmdbImport,
            cmdbPropertySelector
        },
        data () {
            return {
                objectUnique: [],
                properties: [],
                propertyGroups: [],
                table: {
                    checked: [],
                    header: [],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    defaultSort: 'bk_inst_id',
                    sort: 'bk_inst_id',
                    stuff: {
                        type: 'default',
                        payload: {}
                    }
                },
                filter: {
                    id: 'bk_inst_name',
                    value: RouterQuery.get('filter', ''),
                    type: 'singlechar'
                },
                slider: {
                    show: false,
                    contentShow: false,
                    title: ''
                },
                tab: {
                    active: 'attribute'
                },
                attribute: {
                    type: null,
                    inst: {
                        details: {},
                        edit: {}
                    }
                },
                columnsConfig: {
                    show: false,
                    selected: [],
                    disabledColumns: ['bk_inst_id', 'bk_inst_name']
                },
                importSlider: {
                    show: false
                }
            }
        },
        computed: {
            ...mapState('userCustom', ['globalUsercustom']),
            ...mapGetters(['supplierAccount', 'userName']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('objectModelClassify', ['models', 'getModelById']),
            objId () {
                return this.$route.params.objId
            },
            model () {
                return this.getModelById(this.objId) || {}
            },
            customConfigKey () {
                return `${this.objId}_custom_table_columns`
            },
            customColumns () {
                return this.usercustom[this.customConfigKey] || []
            },
            globalCustomColumns () {
                return this.globalUsercustom[`${this.objId}_global_custom_table_columns`] || []
            },
            url () {
                const prefix = `${window.API_HOST}insts/owner/${this.supplierAccount}/object/${this.objId}/`
                return {
                    import: prefix + 'import',
                    export: prefix + 'export',
                    template: `${window.API_HOST}importtemplate/${this.objId}`
                }
            },
            parentLayers () {
                return [{
                    resource_id: this.model.id,
                    resource_type: 'model'
                }]
            },
            filterType () {
                const propertyId = this.filter.id
                const property = this.properties.find(property => property.bk_property_id === propertyId)
                if (property) {
                    return property.bk_property_type
                }
                return 'singlechar'
            }
        },
        watch: {
            'filter.id' (id) {
                this.filter.value = ''
            },
            'slider.show' (show) {
                if (!show) {
                    this.tab.active = 'attribute'
                }
                this.$nextTick(() => {
                    this.slider.contentShow = show
                })
            },
            customColumns () {
                this.setTableHeader()
            },
            objId () {
                this.setDynamicBreadcrumbs()
                this.reload()
            }
        },
        created () {
            this.unwatch = RouterQuery.watch('*', ({
                page = 1,
                limit = this.table.pagination.limit,
                filter = '',
                field = 'bk_inst_name'
            }) => {
                this.filter.id = field
                this.filter.value = filter
                this.table.pagination.current = parseInt(page)
                this.table.pagination.limit = parseInt(limit)
                this.getTableData(!!filter)
            })
            this.setDynamicBreadcrumbs()
            this.reload()
        },
        beforeDestroy () {
            this.unwatch()
        },
        beforeRouteUpdate (to, from, next) {
            this.setDynamicBreadcrumbs()
            next()
        },
        methods: {
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectCommonInst', [
                'createInst',
                'searchInst',
                'updateInst',
                'batchUpdateInst',
                'deleteInst',
                'batchDeleteInst',
                'searchInstById'
            ]),
            setDynamicBreadcrumbs () {
                this.$store.commit('setTitle', this.model.bk_obj_name)
            },
            async reload () {
                try {
                    this.setRencentlyData()
                    this.resetData()
                    this.properties = await this.searchObjectAttribute({
                        injectId: this.objId,
                        params: {
                            bk_obj_id: this.objId,
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: `post_searchObjectAttribute_${this.objId}`,
                            fromCache: false
                        }
                    })
                    await Promise.all([
                        this.getPropertyGroups(),
                        this.getObjectUnique(),
                        this.setTableHeader()
                    ])
                    RouterQuery.set({
                        _t: Date.now()
                    })
                } catch (e) {
                    // ignore
                }
            },
            handleFilterValueChange () {
                RouterQuery.set({
                    _t: Date.now(),
                    page: 1,
                    field: this.filter.id,
                    filter: this.filter.value
                })
            },
            resetData () {
                this.table = {
                    checked: [],
                    header: [],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    defaultSort: 'bk_inst_id',
                    sort: 'bk_inst_id',
                    stuff: {
                        type: 'default',
                        payload: {
                            resource: this.model['bk_obj_name']
                        }
                    }
                }
            },
            getPropertyGroups () {
                return this.searchGroup({
                    objId: this.objId,
                    params: {},
                    config: {
                        fromCache: false,
                        requestId: `post_searchGroup_${this.objId}`
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            getObjectUnique () {
                return this.$store.dispatch('objectUnique/searchObjectUniqueConstraints', {
                    objId: this.objId,
                    params: {}
                }).then(data => {
                    this.objectUnique = data
                    return data
                })
            },
            setTableHeader () {
                return new Promise((resolve, reject) => {
                    const customColumns = this.customColumns.length ? this.customColumns : this.globalCustomColumns
                    const headerProperties = this.$tools.getHeaderProperties(this.properties, customColumns, this.columnsConfig.disabledColumns)
                    resolve(headerProperties)
                }).then(properties => {
                    this.updateTableHeader(properties)
                    this.columnsConfig.selected = properties.map(property => property['bk_property_id'])
                })
            },
            updateTableHeader (properties) {
                this.table.header = properties.map(property => {
                    return {
                        id: property.bk_property_id,
                        name: this.$tools.getHeaderPropertyName(property),
                        property
                    }
                })
            },
            handleValueClick (item, column) {
                if (column.id !== 'bk_inst_id') {
                    return false
                }
                this.$routerActions.redirect({
                    name: MENU_RESOURCE_INSTANCE_DETAILS,
                    params: {
                        objId: this.objId,
                        instId: item.bk_inst_id
                    },
                    history: true
                })
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                RouterQuery.set({
                    _t: Date.now()
                })
            },
            handleSizeChange (size) {
                RouterQuery.set({
                    limit: size,
                    page: 1,
                    _t: Date.now()
                })
            },
            handlePageChange (page, withFilter = false) {
                RouterQuery.set({
                    page: page,
                    _t: Date.now()
                })
            },
            handleSelectChange (selection) {
                this.table.checked = selection.map(row => row.bk_inst_id)
            },
            getInstList (config = { cancelPrevious: true }) {
                return this.searchInst({
                    objId: this.objId,
                    params: this.getSearchParams(),
                    config: Object.assign({ requestId: `post_searchInst_${this.objId}` }, config)
                })
            },
            getTableData (withFilter) {
                this.getInstList({ cancelPrevious: true, globalPermission: false }).then(data => {
                    if (data.count && !data.info.length) {
                        this.table.pagination.current -= 1
                        this.getTableData()
                    }
                    this.table.list = data.info
                    this.table.pagination.count = data.count

                    this.table.stuff.type = withFilter ? 'search' : 'default'

                    return data
                }).catch(({ permission }) => {
                    if (permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission }
                        }
                    }
                })
            },
            getSearchParams () {
                const params = {
                    condition: {
                        [this.objId]: []
                    },
                    fields: {},
                    page: {
                        start: this.table.pagination.limit * (this.table.pagination.current - 1),
                        limit: this.table.pagination.limit,
                        sort: this.table.sort
                    }
                }
                if (this.filter.id && String(this.filter.value).length) {
                    const filterType = this.filterType
                    let filterValue = this.filter.value
                    if (filterType === 'bool') {
                        const convertValue = [true, false].find(bool => bool.toString() === filterValue)
                        filterValue = convertValue === undefined ? filterValue : convertValue
                    } else if (filterType === 'int') {
                        filterValue = isNaN(parseInt(filterValue)) ? filterValue : parseInt(filterValue)
                    } else if (filterType === 'float') {
                        filterValue = isNaN(parseFloat(filterValue)) ? filterValue : parseFloat(filterValue)
                    }
                    if (['bool', 'int', 'enum', 'float', 'organization'].includes(filterType)) {
                        params.condition[this.objId].push({
                            field: this.filter.id,
                            operator: '$eq',
                            value: filterValue
                        })
                    } else if (['singleasst', 'multiasst'].includes(filterType)) {
                        const asstObjId = (this.$tools.getProperty(this.properties, this.filter.id) || {})['bk_asst_obj_id']
                        if (asstObjId) {
                            const fieldMap = {
                                'host': 'bk_host_innerip',
                                'biz': 'bk_biz_name',
                                'plat': 'bk_cloud_name',
                                'module': 'bk_module_name',
                                'set': 'bk_set_name'
                            }
                            params.condition[asstObjId] = [{
                                field: fieldMap.hasOwnProperty(asstObjId) ? fieldMap[asstObjId] : 'bk_inst_name',
                                operator: '$in',
                                value: filterValue.split(',')
                            }]
                        }
                    } else {
                        params.condition[this.objId].push({
                            field: this.filter.id,
                            operator: '$regex',
                            value: filterValue
                        })
                    }
                }
                return params
            },
            async handleEdit (item) {
                this.attribute.inst.edit = item
                this.attribute.type = 'update'
            },
            handleCreate () {
                this.attribute.type = 'create'
                this.attribute.inst.edit = {}
                this.slider.show = true
                this.slider.title = `${this.$t('创建')} ${this.model['bk_obj_name']}`
            },
            handleDelete (inst) {
                this.$bkInfo({
                    title: this.$t('确认要删除', { name: inst['bk_inst_name'] }),
                    confirmFn: () => {
                        this.deleteInst({
                            objId: this.objId,
                            instId: inst['bk_inst_id']
                        }).then(() => {
                            this.slider.show = false
                            this.$success(this.$t('删除成功'))
                            RouterQuery.set({
                                _t: Date.now()
                            })
                        })
                    }
                })
            },
            handleSave (values, changedValues, originalValues, type) {
                if (type === 'update') {
                    this.updateInst({
                        objId: this.objId,
                        instId: originalValues['bk_inst_id'],
                        params: values
                    }).then(() => {
                        this.attribute.inst.details = Object.assign({}, originalValues, values)
                        this.handleCancel()
                        this.$success(this.$t('修改成功'))
                        RouterQuery.set({
                            _t: Date.now()
                        })
                    })
                } else {
                    delete values.bk_inst_id // properties中注入了前端自定义的bk_inst_id属性
                    this.createInst({
                        params: values,
                        objId: this.objId
                    }).then(() => {
                        RouterQuery.set({
                            _t: Date.now(),
                            page: 1
                        })
                        this.handleCancel()
                        this.$success(this.$t('创建成功'))
                    })
                }
            },
            handleCancel () {
                if (this.attribute.type === 'create') {
                    this.slider.show = false
                }
            },
            handleMultipleEdit () {
                this.attribute.type = 'multiple'
                this.slider.title = this.$t('批量更新')
                this.slider.show = true
            },
            handleMultipleSave (values) {
                this.batchUpdateInst({
                    objId: this.objId,
                    params: {
                        update: this.table.checked.map(instId => {
                            return {
                                'datas': values,
                                'inst_id': instId
                            }
                        })
                    },
                    config: {
                        requestId: `${this.objId}BatchUpdate`
                    }
                }).then(() => {
                    this.$success(this.$t('修改成功'))
                    this.slider.show = false
                    RouterQuery.set({
                        _t: Date.now(),
                        page: 1
                    })
                })
            },
            handleMultipleCancel () {
                this.slider.show = false
            },
            handleMultipleDelete () {
                this.$bkInfo({
                    title: this.$t('确定删除选中的实例'),
                    confirmFn: () => {
                        this.doBatchDeleteInst()
                    }
                })
            },
            doBatchDeleteInst () {
                this.batchDeleteInst({
                    objId: this.objId,
                    config: {
                        data: {
                            'delete': {
                                'inst_ids': this.table.checked
                            }
                        }
                    }
                }).then(() => {
                    this.$success(this.$t('删除成功'))
                    this.table.checked = []
                    RouterQuery.set({
                        _t: Date.now()
                    })
                })
            },
            handleApplyColumnsConfig (properties) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.customConfigKey]: properties.map(property => property['bk_property_id'])
                })
                this.columnsConfig.show = false
            },
            handleResetColumnsConfig () {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.customConfigKey]: []
                })
                this.columnsConfig.show = false
            },
            routeToHistory () {
                this.$routerActions.redirect({
                    name: 'instanceHistory',
                    params: {
                        objId: this.objId
                    },
                    history: true
                })
            },
            handleSliderBeforeClose () {
                if (this.tab.active === 'attribute') {
                    const $form = this.attribute.type === 'multiple' ? this.$refs.multipleForm : this.$refs.form
                    if ($form.hasChange) {
                        return new Promise((resolve, reject) => {
                            this.$bkInfo({
                                title: this.$t('确认退出'),
                                subTitle: this.$t('退出会导致未保存信息丢失'),
                                extCls: 'bk-dialog-sub-header-center',
                                confirmFn: () => {
                                    resolve(true)
                                },
                                cancelFn: () => {
                                    resolve(false)
                                }
                            })
                        })
                    }
                    return true
                }
                return true
            },
            handleExport () {
                const data = new FormData()
                data.append('bk_inst_id', this.table.checked.join(','))
                const customFields = this.usercustom[this.customConfigKey]
                if (customFields) {
                    data.append('export_custom_fields', customFields)
                }
                this.$http.download({
                    url: this.url.export,
                    method: 'post',
                    data
                })
            },
            setRencentlyData () {
                const modelId = this.model.id
                this.$store.dispatch('userCustom/setRencentlyData', { id: modelId })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .models-layout {
        padding: 15px 20px 0;
    }
    .options-filter{
        position: relative;
        margin-right: 5px;
        .filter-selector{
            width: 120px;
            border-radius: 2px 0 0 2px;
            margin-right: -1px;
        }
        .filter-value{
            width: 320px;
            border-radius: 0 2px 2px 0;
            /deep/ .bk-form-input {
                line-height: 32px;
            }
        }
        .filter-search{
            position: absolute;
            right: 10px;
            top: 8px;
            cursor: pointer;
        }
    }
    .models-button{
        display: inline-block;
        position: relative;
        &:hover{
            z-index: 1;
        }
    }
    .models-table{
        margin-top: 14px;
    }
</style>
