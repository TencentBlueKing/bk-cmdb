<template>
    <div class="models-layout">
        <div class="models-options clearfix">
            <div class="options-button clearfix fl">
                <cmdb-auth class="fl mr10" :auth="$authResources({ type: $OPERATION.C_INST, parent_layers: parentLayers })">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        :disabled="disabled"
                        @click="handleCreate">
                        {{$t('新建')}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth class="fl mr10" :auth="$authResources({ type: [$OPERATION.C_INST, $OPERATION.U_INST], parent_layers: parentLayers })">
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
                <cmdb-auth class="fl mr10" :auth="$authResources({ type: $OPERATION.D_INST, parent_layers: parentLayers })">
                    <bk-button slot-scope="{ disabled }"
                        class="models-button button-delete"
                        :disabled="!table.checked.length || disabled"
                        @click="handleMultipleDelete">
                        {{$t('删除')}}
                    </bk-button>
                </cmdb-auth>
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
                <bk-select class="filter-selector fl"
                    v-model="filter.id"
                    searchable
                    font-size="medium"
                    :clearable="false">
                    <bk-option v-for="(option, index) in filter.options"
                        :key="index"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                </bk-select>
                <cmdb-form-enum class="filter-value fl"
                    v-if="filter.type === 'enum'"
                    :options="$tools.getEnumOptions(properties, filter.id)"
                    :allow-clear="true"
                    :auto-select="false"
                    font-size="medium"
                    v-model="filter.value"
                    @on-selected="getTableData(true)">
                </cmdb-form-enum>
                <bk-input class="filter-value cmdb-form-input fl" type="text" maxlength="11"
                    v-else-if="filter.type === 'int'"
                    v-model.number="filter.value"
                    clearable
                    right-icon="icon-search"
                    font-size="medium"
                    :placeholder="$t('快速查询')"
                    @enter="getTableData(true)">
                </bk-input>
                <bk-input class="filter-value cmdb-form-input fl" type="text"
                    v-else-if="filter.type === 'float'"
                    v-model.number="filter.value"
                    clearable
                    right-icon="icon-search"
                    font-size="medium"
                    :placeholder="$t('快速查询')"
                    @enter="getTableData(true)">
                </bk-input>
                <bk-input class="filter-value cmdb-form-input fl" type="text"
                    v-else
                    v-model.trim="filter.value"
                    clearable
                    right-icon="icon-search"
                    font-size="medium"
                    :placeholder="$t('快速查询')"
                    @enter="getTableData(true)">
                </bk-input>
            </div>
        </div>
        <bk-table class="models-table" ref="table"
            v-bkloading="{ isLoading: $loading() }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 190"
            :row-style="{ cursor: 'pointer' }"
            @row-click="handleRowClick"
            @sort-change="handleSortChange"
            @page-limit-change="handleSizeChange"
            @page-change="handlePageChange"
            @selection-change="handleSelectChange">
            <bk-table-column type="selection" width="60" align="center" fixed class-name="bk-table-selection"></bk-table-column>
            <bk-table-column v-for="column in table.header"
                sortable="custom"
                :key="column.id"
                :prop="column.id"
                :label="column.name"
                :class-name="column.id === 'bk_inst_name' ? 'is-highlight' : ''"
                :fixed="column.id === 'bk_inst_name'">
                <template slot-scope="{ row }">
                    <span>{{row[column.id] | addUnit(getPropertyUnit(column.id))}}</span>
                </template>
            </bk-table-column>
            <cmdb-table-empty
                slot="empty"
                :auth="$authResources({
                    type: $OPERATION.C_INST,
                    parent_layers: parentLayers
                })"
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
                    <cmdb-details v-if="attribute.type === 'details'"
                        :properties="properties"
                        :property-groups="propertyGroups"
                        :inst="attribute.inst.details"
                        :edit-auth="$OPERATION.U_INST"
                        :delete-auth="$OPERATION.D_INST"
                        @on-edit="handleEdit"
                        @on-delete="handleDelete">
                    </cmdb-details>
                    <cmdb-form v-else-if="['update', 'create'].includes(attribute.type)"
                        ref="form"
                        :properties="properties"
                        :property-groups="propertyGroups"
                        :inst="attribute.inst.edit"
                        :type="attribute.type"
                        :save-auth="attribute.type === 'update' ? $OPERATION.U_INST : $OPERATION.C_INST"
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
                <bk-tab-panel name="relevance" :label="$t('关联Relation')" :visible="['update', 'details'].includes(attribute.type)">
                    <cmdb-relation
                        v-if="tab.active === 'relevance'"
                        :auth="$OPERATION.U_INST"
                        :obj-id="objId"
                        :inst="attribute.inst.details">
                    </cmdb-relation>
                </bk-tab-panel>
                <bk-tab-panel name="history" :label="$t('变更记录')" :visible="['update', 'details'].includes(attribute.type)">
                    <cmdb-audit-history v-if="tab.active === 'history'"
                        :target="objId"
                        :inst-id="attribute.inst.details['bk_inst_id']">
                    </cmdb-audit-history>
                </bk-tab-panel>
            </bk-tab>
        </bk-sideslider>
        <bk-sideslider v-transfer-dom :is-show.sync="columnsConfig.show" :width="600" :title="$t('列表显示属性配置')">
            <cmdb-columns-config slot="content"
                v-if="columnsConfig.show"
                :properties="columnProperties"
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
                :download-payload="url.downloadPayload"
                :import-payload="url.importPayload"
                @success="handlePageChange(1)"
                @partialSuccess="handlePageChange(1)">
            </cmdb-import>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import cmdbColumnsConfig from '@/components/columns-config/columns-config'
    import cmdbAuditHistory from '@/components/audit-history/audit-history.vue'
    import cmdbRelation from '@/components/relation'
    import cmdbImport from '@/components/import/import'
    import { MENU_RESOURCE_MANAGEMENT } from '@/dictionary/menu-symbol'
    export default {
        filters: {
            addUnit (value, unit) {
                if (value === '--' || !unit) {
                    return value
                }
                return value + unit
            }
        },
        components: {
            cmdbColumnsConfig,
            cmdbAuditHistory,
            cmdbRelation,
            cmdbImport
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
                    allList: [],
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
                    id: '',
                    value: '',
                    type: '',
                    options: []
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
                    disabledColumns: ['bk_inst_name']
                },
                importSlider: {
                    show: false
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'isAdminView']),
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
                return `${this.userName}_${this.objId}_${this.isAdminView ? 'adminView' : this.bizId}_table_columns`
            },
            customColumns () {
                return this.usercustom[this.customConfigKey]
            },
            url () {
                const prefix = `${window.API_HOST}insts/owner/${this.supplierAccount}/object/${this.objId}/`
                return {
                    import: prefix + 'import',
                    export: prefix + 'export',
                    template: `${window.API_HOST}importtemplate/${this.objId}`,
                    downloadPayload: this.$injectMetadata({}, { inject: !this.isPublicModel }),
                    importPayload: this.$injectMetadata({}, { inject: !this.isPublicModel })
                }
            },
            isPublicModel () {
                const model = this.models.find(model => model['bk_obj_id'] === this.objId) || {}
                return !this.$tools.getMetadataBiz(model)
            },
            parentLayers () {
                return [{
                    resource_id: this.model.id,
                    resource_type: 'model'
                }]
            },
            columnProperties () {
                const instId = {
                    bk_property_id: 'bk_inst_id',
                    bk_property_name: 'ID'
                }
                const properties = this.properties
                properties.push(instId)
                return properties
            }
        },
        watch: {
            'filter.id' (id) {
                this.filter.value = ''
                this.filter.type = (this.$tools.getProperty(this.properties, id) || {})['bk_property_type']
            },
            'filter.value' () {
                this.$route.query.instId = null
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
            this.setDynamicBreadcrumbs()
            this.reload()
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
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('资源目录'),
                    route: {
                        name: MENU_RESOURCE_MANAGEMENT
                    }
                }, {
                    label: this.model.bk_obj_name
                }])
            },
            getPropertyUnit (propertyId) {
                const property = this.properties.find(property => property.bk_property_id === propertyId)
                if (!property) {
                    return ''
                }
                return property.unit || ''
            },
            async reload () {
                try {
                    this.setRencentlyData()
                    this.resetData()
                    this.properties = await this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: this.objId,
                            bk_supplier_account: this.supplierAccount
                        }, { inject: !this.isPublicModel }),
                        config: {
                            requestId: `post_searchObjectAttribute_${this.objId}`,
                            fromCache: false
                        }
                    })
                    await Promise.all([
                        this.getPropertyGroups(),
                        this.setTableHeader(),
                        this.setFilterOptions()
                    ])
                    this.getTableData()
                } catch (e) {
                    // ignore
                }
            },
            resetData () {
                this.table = {
                    checked: [],
                    header: [],
                    list: [],
                    allList: [],
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
                    params: this.$injectMetadata({}, { inject: !this.isPublicModel }),
                    config: {
                        fromCache: false,
                        requestId: `post_searchGroup_${this.objId}`
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            setTableHeader () {
                return new Promise((resolve, reject) => {
                    const headerProperties = this.$tools.getHeaderProperties(this.columnProperties, this.customColumns, this.columnsConfig.disabledColumns)
                    resolve(headerProperties)
                }).then(properties => {
                    this.updateTableHeader(properties)
                    this.columnsConfig.selected = properties.map(property => property['bk_property_id'])
                })
            },
            setFilterOptions () {
                this.filter.options = this.properties.map(property => {
                    return {
                        id: property['bk_property_id'],
                        name: property['bk_property_name']
                    }
                })
                this.filter.id = this.filter.options.length ? this.filter.options[0]['id'] : ''
            },
            updateTableHeader (properties) {
                this.table.header = properties.map(property => {
                    return {
                        id: property['bk_property_id'],
                        name: property['bk_property_name']
                    }
                })
            },
            async handleCheckAll (type) {
                if (type === 'current') {
                    this.table.checked = this.table.list.map(inst => inst['bk_inst_id'])
                } else {
                    const allData = await this.getAllInstList()
                    this.table.checked = allData.info.map(inst => inst['bk_inst_id'])
                }
            },
            handleRowClick (item) {
                this.slider.show = true
                this.slider.title = item['bk_inst_name']
                this.attribute.inst.details = item
                this.attribute.type = 'details'
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            },
            handleSelectChange (selection) {
                this.table.checked = selection.map(row => row.bk_inst_id)
            },
            getInstList (config = { cancelPrevious: true }) {
                return this.searchInst({
                    objId: this.objId,
                    params: this.$injectMetadata(this.getSearchParams(), { inject: !this.isPublicModel }),
                    config: Object.assign({ requestId: `post_searchInst_${this.objId}` }, config)
                })
            },
            getAllInstList () {
                return this.searchInst({
                    objId: this.objId,
                    params: this.$injectMetadata({
                        ...this.getSearchParams(),
                        page: {}
                    }, { inject: !this.isPublicModel }),
                    config: {
                        requestId: `${this.objId}AllList`,
                        cancelPrevious: true
                    }
                }).then(data => {
                    this.table.allList = data.info
                    return data
                })
            },
            setAllHostList (list) {
                const newList = []
                list.forEach(item => {
                    const existItem = this.table.allList.some(existItem => existItem['bk_inst_id'] === item['bk_inst_id'])
                    if (existItem) {
                        Object.assign(existItem, item)
                    } else {
                        newList.push(item)
                    }
                })
                this.table.allList = [...this.table.allList, ...newList]
            },
            getTableData (event) {
                this.getInstList({ cancelPrevious: true, globalPermission: false }).then(data => {
                    if (data.count && !data.info.length) {
                        this.table.pagination.current -= 1
                        this.getTableData()
                    }
                    this.table.list = this.$tools.flattenList(this.properties, data.info)
                    this.table.pagination.count = data.count
                    this.setAllHostList(data.info)

                    if (event) {
                        this.table.stuff.type = 'search'
                    }

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
                    const filterType = this.filter.type
                    let filterValue = this.filter.value
                    if (filterType === 'bool') {
                        const convertValue = [true, false].find(bool => bool.toString() === filterValue)
                        filterValue = convertValue === undefined ? filterValue : convertValue
                    } else if (filterType === 'int') {
                        filterValue = isNaN(parseInt(filterValue)) ? filterValue : parseInt(filterValue)
                    } else if (filterType === 'float') {
                        filterValue = isNaN(parseFloat(filterValue)) ? filterValue : parseFloat(filterValue)
                    }
                    if (['bool', 'int', 'enum', 'float'].includes(filterType)) {
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
                } else if (this.$route.params.instId) {
                    params.condition[this.objId].push({
                        field: 'bk_inst_id',
                        operator: '$in',
                        value: [Number(this.$route.params.instId)]
                    })
                    this.$route.params.instId = null
                }
                return params
            },
            async handleEdit (flattenItem) {
                const list = await this.getInstList({ fromCache: true })
                const inst = list.info.find(item => item['bk_inst_id'] === flattenItem['bk_inst_id'])
                this.attribute.inst.edit = inst
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
                            instId: inst['bk_inst_id'],
                            config: {
                                data: this.$injectMetadata({}, { inject: !this.isPublicModel })
                            }
                        }).then(() => {
                            this.slider.show = false
                            this.$success(this.$t('删除成功'))
                            this.getTableData()
                        })
                    }
                })
            },
            handleSave (values, changedValues, originalValues, type) {
                if (type === 'update') {
                    this.updateInst({
                        objId: this.objId,
                        instId: originalValues['bk_inst_id'],
                        params: this.$injectMetadata(values, { inject: !this.isPublicModel })
                    }).then(() => {
                        this.getTableData()
                        this.searchInstById({
                            objId: this.objId,
                            instId: originalValues['bk_inst_id'],
                            params: this.$injectMetadata({}, { inject: !this.isPublicModel })
                        }).then(item => {
                            this.attribute.inst.details = this.$tools.flattenItem(this.properties, item)
                        })
                        this.handleCancel()
                        this.$success(this.$t('修改成功'))
                    })
                } else {
                    this.createInst({
                        params: this.$injectMetadata(values, { inject: !this.isPublicModel }),
                        objId: this.objId
                    }).then(() => {
                        this.handlePageChange(1)
                        this.handleCancel()
                        this.$success(this.$t('创建成功'))
                    })
                }
            },
            handleCancel () {
                if (this.attribute.type === 'create') {
                    this.slider.show = false
                } else {
                    this.attribute.type = 'details'
                }
            },
            async handleMultipleEdit () {
                this.objectUnique = await this.$store.dispatch('objectUnique/searchObjectUniqueConstraints', {
                    objId: this.objId,
                    params: this.$injectMetadata({}, {
                        inject: !this.isPublicModel
                    })
                })
                this.attribute.type = 'multiple'
                this.slider.title = this.$t('批量更新')
                this.slider.show = true
            },
            handleMultipleSave (values) {
                this.batchUpdateInst({
                    objId: this.objId,
                    params: this.$injectMetadata({
                        update: this.table.checked.map(instId => {
                            return {
                                'datas': values,
                                'inst_id': instId
                            }
                        })
                    }, { inject: !this.isPublicModel }),
                    config: {
                        requestId: `${this.objId}BatchUpdate`
                    }
                }).then(() => {
                    this.$success(this.$t('修改成功'))
                    this.slider.show = false
                    this.handlePageChange(1)
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
                        data: this.$injectMetadata({
                            'delete': {
                                'inst_ids': this.table.checked
                            }
                        }, { inject: !this.isPublicModel })
                    }
                }).then(() => {
                    this.$success(this.$t('删除成功'))
                    this.table.checked = []
                    this.getTableData()
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
            },
            routeToHistory () {
                this.$router.push({
                    name: 'instanceHistory',
                    params: {
                        objId: this.objId
                    }
                })
            },
            handleSliderBeforeClose () {
                if (this.tab.active === 'attribute' && this.attribute.type !== 'details') {
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
                if (!this.isPublicModel) {
                    data.append('metadata', JSON.stringify(this.$injectMetadata().metadata))
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
        padding: 0 20px;
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
            &.button-delete {
                color: $cmdbDangerColor;
                border-color: $cmdbDangerColor;
            }
            /deep/ &.bk-button.bk-default[disabled] {
                border-color: #dcdee5 !important;
                color: #c4c6cc !important;
            }
        }
    }
    .models-table{
        margin-top: 14px;
    }
</style>
