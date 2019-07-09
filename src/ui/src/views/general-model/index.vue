<template>
    <div class="models-layout">
        <div class="models-options clearfix">
            <div class="options-button clearfix fl">
                <div class="fl" v-tooltip="$t('ModelManagement[\'导入\']')"
                    v-cursor="{
                        active: !$isAuthorized([$OPERATION.C_INST, $OPERATION.U_INST]),
                        auth: [$OPERATION.C_INST, $OPERATION.U_INST]
                    }">
                    <bk-button class="models-button"
                        :disabled="!$isAuthorized([$OPERATION.C_INST, $OPERATION.U_INST])"
                        @click="importSlider.show = true">
                        <i class="icon-cc-import"></i>
                    </bk-button>
                </div>
                <div class="fl" v-tooltip="$t('ModelManagement[\'导出\']')">
                    <bk-button class="models-button" type="default"
                        :disabled="!table.checked.length"
                        @click="handleExport">
                        <i class="icon-cc-derivation"></i>
                    </bk-button>
                </div>
                <div class="fl" v-tooltip="$t('Inst[\'批量更新\']')"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.U_INST),
                        auth: [$OPERATION.U_INST]
                    }">
                    <bk-button class="models-button"
                        :disabled="!table.checked.length || !$isAuthorized($OPERATION.U_INST)"
                        @click="handleMultipleEdit">
                        <i class="icon-cc-edit"></i>
                    </bk-button>
                </div>
                <div class="fl" v-tooltip="$t('Common[\'删除\']')"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.D_INST),
                        auth: [$OPERATION.D_INST]
                    }">
                    <bk-button class="models-button button-delete"
                        :disabled="!table.checked.length || !$isAuthorized($OPERATION.D_INST)"
                        @click="handleMultipleDelete">
                        <i class="icon-cc-del"></i>
                    </bk-button>
                </div>
                <div class="fl" style="margin-left: 20px;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.C_INST),
                        auth: [$OPERATION.C_INST]
                    }">
                    <bk-button type="primary"
                        :disabled="!$isAuthorized($OPERATION.C_INST)"
                        @click="handleCreate">
                        {{$t("Common['新建']")}}
                    </bk-button>
                </div>
            </div>
            <div class="options-button fr">
                <bk-button v-tooltip="$t('Common[\'查看删除历史\']')" @click="routeToHistory">
                    <i class="icon-cc-history"></i>
                </bk-button>
                <bk-button class="button-setting" v-tooltip="$t('BusinessTopology[\'列表显示属性配置\']')" @click="columnsConfig.show = true">
                    <i class="icon-cc-setting"></i>
                </bk-button>
            </div>
            <div class="options-filter clearfix fr">
                <bk-selector class="filter-selector fl"
                    :searchable="true"
                    :list="filter.options"
                    :selected.sync="filter.id">
                </bk-selector>
                <cmdb-form-enum class="filter-value fl"
                    v-if="filter.type === 'enum'"
                    :options="$tools.getEnumOptions(properties, filter.id)"
                    :allow-clear="true"
                    v-model="filter.value"
                    @on-selected="getTableData">
                </cmdb-form-enum>
                <input class="filter-value cmdb-form-input fl" type="text" maxlength="11"
                    v-else-if="filter.type === 'int'"
                    v-model.number="filter.value"
                    :placeholder="$t('Common[\'快速查询\']')"
                    @keydown.enter="getTableData">
                <input class="filter-value cmdb-form-input fl" type="text"
                    v-else-if="filter.type === 'float'"
                    v-model.number="filter.value"
                    :placeholder="$t('Common[\'快速查询\']')"
                    @keydown.enter="getTableData">
                <input class="filter-value cmdb-form-input fl" type="text"
                    v-else
                    v-model.trim="filter.value"
                    :placeholder="$t('Common[\'快速查询\']')"
                    @keydown.enter="getTableData">
                <i class="filter-search bk-icon icon-search"
                    v-show="filter.type !== 'enum'"
                    @click="getTableData"></i>
            </div>
        </div>
        <cmdb-table class="models-table" ref="table"
            :loading="$loading()"
            :checked.sync="table.checked"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :default-sort="table.defaultSort"
            :wrapper-minus-height="157"
            @handleRowClick="handleRowClick"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange"
            @handleCheckAll="handleCheckAll">
        </cmdb-table>
        <cmdb-slider :is-show.sync="slider.show" :title="slider.title" :before-close="handleSliderBeforeClose">
            <bk-tab :active-name.sync="tab.active" slot="content">
                <bk-tabpanel name="attribute" :title="$t('Common[\'属性\']')" style="width: calc(100% + 40px);margin: 0 -20px;">
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
                        :properties="properties"
                        :property-groups="propertyGroups"
                        :object-unique="objectUnique"
                        :save-auth="$OPERATION.U_INST"
                        @on-submit="handleMultipleSave"
                        @on-cancel="handleMultipleCancel">
                    </cmdb-form-multiple>
                </bk-tabpanel>
                <bk-tabpanel name="relevance" :title="$t('HostResourcePool[\'关联\']')" :show="['update', 'details'].includes(attribute.type)">
                    <cmdb-relation
                        v-if="tab.active === 'relevance'"
                        :auth="$OPERATION.U_INST"
                        :obj-id="objId"
                        :inst="attribute.inst.details">
                    </cmdb-relation>
                </bk-tabpanel>
                <bk-tabpanel name="history" :title="$t('HostResourcePool[\'变更记录\']')" :show="['update', 'details'].includes(attribute.type)">
                    <cmdb-audit-history v-if="tab.active === 'history'"
                        :target="objId"
                        :inst-id="attribute.inst.details['bk_inst_id']">
                    </cmdb-audit-history>
                </bk-tabpanel>
            </bk-tab>
        </cmdb-slider>
        <cmdb-slider :is-show.sync="columnsConfig.show" :width="600" :title="$t('BusinessTopology[\'列表显示属性配置\']')">
            <cmdb-columns-config slot="content"
                :properties="properties"
                :selected="columnsConfig.selected"
                :disabled-columns="columnsConfig.disabledColumns"
                @on-apply="handleApplyColumnsConfig"
                @on-cancel="columnsConfig.show = false"
                @on-reset="handleResetColumnsConfig">
            </cmdb-columns-config>
        </cmdb-slider>
        <cmdb-slider
            :is-show.sync="importSlider.show"
            :title="$t('HostResourcePool[\'批量导入\']')">
            <cmdb-import v-if="importSlider.show" slot="content"
                :template-url="url.template"
                :import-url="url.import"
                :download-payload="url.downloadPayload"
                :import-payload="url.importPayload"
                @success="handlePageChange(1)"
                @partialSuccess="handlePageChange(1)">
            </cmdb-import>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import cmdbColumnsConfig from '@/components/columns-config/columns-config'
    import cmdbAuditHistory from '@/components/audit-history/audit-history.vue'
    import cmdbRelation from '@/components/relation'
    import cmdbImport from '@/components/import/import'
    export default {
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
                        size: 10,
                        current: 1
                    },
                    defaultSort: 'bk_inst_id',
                    sort: 'bk_inst_id'
                },
                filter: {
                    id: '',
                    value: '',
                    type: '',
                    options: []
                },
                slider: {
                    show: false,
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
            }
        },
        watch: {
            'filter.id' (id) {
                this.filter.value = ''
                this.filter.type = (this.$tools.getProperty(this.properties, id) || {})['bk_property_type']
            },
            'slider.show' (show) {
                if (!show) {
                    this.tab.active = 'attribute'
                }
            },
            customColumns () {
                this.setTableHeader()
            },
            objId () {
                this.$store.commit('setHeaderTitle', this.model['bk_obj_name'])
                this.reload()
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.model['bk_obj_name'])
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
                        size: 10,
                        current: 1
                    },
                    defaultSort: 'bk_inst_id',
                    sort: 'bk_inst_id'
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
                    const headerProperties = this.$tools.getHeaderProperties(this.properties, this.customColumns, this.columnsConfig.disabledColumns)
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
                this.table.header = [{
                    id: 'bk_inst_id',
                    type: 'checkbox',
                    width: 50
                }].concat(properties.map(property => {
                    return {
                        id: property['bk_property_id'],
                        name: property['bk_property_name']
                    }
                }))
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
                this.table.sort = sort
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
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
            getTableData () {
                this.getInstList().then(data => {
                    this.table.list = this.$tools.flattenList(this.properties, data.info)
                    this.table.pagination.count = data.count
                    this.setAllHostList(data.info)
                    return data
                })
            },
            getSearchParams () {
                const params = {
                    condition: {
                        [this.objId]: []
                    },
                    fields: {},
                    page: {
                        start: this.table.pagination.size * (this.table.pagination.current - 1),
                        limit: this.table.pagination.size,
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
                this.slider.title = `${this.$t("Common['创建']")} ${this.model['bk_obj_name']}`
            },
            handleDelete (inst) {
                this.$bkInfo({
                    title: this.$t("Common['确认要删除']", { name: inst['bk_inst_name'] }),
                    confirmFn: () => {
                        this.deleteInst({
                            objId: this.objId,
                            instId: inst['bk_inst_id'],
                            config: {
                                data: this.$injectMetadata({}, { inject: !this.isPublicModel })
                            }
                        }).then(() => {
                            this.slider.show = false
                            this.$success(this.$t('Common["删除成功"]'))
                            this.handlePageChange(1)
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
                        this.$success(this.$t("Common['修改成功']"))
                    })
                } else {
                    this.createInst({
                        params: this.$injectMetadata(values, { inject: !this.isPublicModel }),
                        objId: this.objId
                    }).then(() => {
                        this.handlePageChange(1)
                        this.handleCancel()
                        this.$success(this.$t("Inst['创建成功']"))
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
                this.slider.title = this.$t('Inst[\'批量更新\']')
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
                    this.$success(this.$t('Common["修改成功"]'))
                    this.handlePageChange(1)
                })
            },
            handleMultipleCancel () {
                this.slider.show = false
            },
            handleMultipleDelete () {
                this.$bkInfo({
                    title: this.$t("Common['确定删除选中的实例']"),
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
                    this.$success(this.$t('Common["删除成功"]'))
                    this.table.checked = []
                    this.handlePageChange(1)
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
                    name: 'history',
                    params: {
                        objId: this.objId
                    },
                    query: {
                        from: this.$route.fullPath
                    }
                })
            },
            handleSliderBeforeClose () {
                if (this.tab.active === 'attribute' && this.attribute.type !== 'details') {
                    const $form = this.attribute.type === 'multiple' ? this.$refs.multipleForm : this.$refs.form
                    const changedValues = $form.changedValues
                    if (Object.keys(changedValues).length) {
                        return new Promise((resolve, reject) => {
                            this.$bkInfo({
                                title: this.$t('Common["退出会导致未保存信息丢失，是否确认？"]'),
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
.options-filter{
    position: relative;
    margin-right: 10px;
    .filter-selector{
        width: 115px;
        border-radius: 2px 0 0 2px;
        margin-right: -1px;
    }
    .filter-value{
        width: 320px;
        border-radius: 0 2px 2px 0;
    }
    .filter-search{
        position: absolute;
        right: 10px;
        top: 11px;
        cursor: pointer;
    }
}
.models-button{
    display: inline-block;
    border-radius: 0;
    margin-left: -1px;
    position: relative;
    &:hover{
        z-index: 1;
        &.button-delete {
            color: $cmdbDangerColor;
            border-color: $cmdbDangerColor;
        }
    }
}
.options-button{
    font-size: 0;
    white-space: nowrap;
    .button-history{
        border-radius: 2px 0 0 2px;
    }
    .button-setting{
        border-radius: 0 2px 2px 0;
        margin-left: -1px;
    }
}
.models-table{
    margin-top: 20px;
}
</style>
