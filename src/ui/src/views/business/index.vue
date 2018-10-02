<template>
    <div class="business-layout">
        <div class="business-options clearfix">
            <bk-button class="fl" type="primary" @click="handleCreate">{{$t("Inst['立即创建']")}}</bk-button>
            <div class="options-button fr">
                <bk-button class="button-history" v-tooltip.bottom="$t('Common[\'查看删除历史\']')" @click="routeToHistory">
                    <i class="icon-cc-history2"></i>
                </bk-button>
                <bk-button class="button-setting" v-tooltip.bottom="$t('BusinessTopology[\'列表显示属性配置\']')" @click="columnsConfig.show = true">
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
                    v-else
                    v-model.trim="filter.value"
                    :placeholder="$t('Common[\'快速查询\']')"
                    @keydown.enter="getTableData">
                <i class="filter-search bk-icon icon-search"
                    v-show="filter.type !== 'enum'"
                    @click="getTableData"></i>
            </div>
        </div>
        <cmdb-table class="business-table" ref="table"
            :loading="$loading('post_searchBusiness_list')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort"
            :wrapperMinusHeight="157"
            @handleRowClick="handleRowClick"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange">
        </cmdb-table>
        <cmdb-slider :isShow.sync="slider.show" :title="slider.title" :beforeClose="handleSliderBeforeClose">
            <bk-tab :active-name.sync="tab.active" slot="content">
                <bk-tabpanel name="attribute" :title="$t('Common[\'属性\']')" style="width: calc(100% + 40px);margin: 0 -20px;">
                    <cmdb-details v-if="attribute.type === 'details'"
                        :properties="properties"
                        :propertyGroups="propertyGroups"
                        :inst="attribute.inst.details"
                        :deleteButtonText="$t('Inst[\'归档\']')"
                        :show-delete="attribute.inst.details['bk_biz_name'] !== '蓝鲸'"
                        @on-edit="handleEdit"
                        @on-delete="handleDelete">
                    </cmdb-details>
                    <cmdb-form v-else-if="['update', 'create'].includes(attribute.type)"
                        ref="form"
                        :properties="properties"
                        :propertyGroups="propertyGroups"
                        :inst="attribute.inst.edit"
                        :type="attribute.type"
                        @on-submit="handleSave"
                        @on-cancel="handleCancel">
                    </cmdb-form>
                </bk-tabpanel>
                <bk-tabpanel name="relevance" :title="$t('HostResourcePool[\'关联\']')" :show="attribute.type !== 'create'">
                    <cmdb-relation
                        v-if="tab.active === 'relevance'"
                        obj-id="biz"
                        :inst-id="attribute.inst.details['bk_biz_id']">
                    </cmdb-relation>
                </bk-tabpanel>
                <bk-tabpanel name="history" :title="$t('HostResourcePool[\'变更记录\']')" :show="attribute.type !== 'create'">
                    <cmdb-audit-history v-if="tab.active === 'history'"
                        target="biz"
                        :instId="attribute.inst.details['bk_biz_id']">
                    </cmdb-audit-history>
                </bk-tabpanel>
            </bk-tab>
        </cmdb-slider>
        <cmdb-slider :isShow.sync="columnsConfig.show" :width="600" :title="$t('BusinessTopology[\'列表显示属性配置\']')">
            <cmdb-columns-config slot="content"
                :properties="properties"
                :selected="columnsConfig.selected"
                :disabled-columns="columnsConfig.disabledColumns"
                @on-apply="handleApplayColumnsConfig"
                @on-cancel="columnsConfig.show = false"
                @on-reset="handleResetColumnsConfig">
            </cmdb-columns-config>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import cmdbColumnsConfig from '@/components/columns-config/columns-config'
    import cmdbAuditHistory from '@/components/audit-history/audit-history.vue'
    import cmdbRelation from '@/components/relation'
    export default {
        components: {
            cmdbColumnsConfig,
            cmdbAuditHistory,
            cmdbRelation
        },
        data () {
            return {
                properties: [],
                propertyGroups: [],
                table: {
                    header: [],
                    list: [],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: 'bk_biz_id',
                    sort: 'bk_biz_id'
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
                        edit: {},
                        details: {}
                    }
                },
                columnsConfig: {
                    show: false,
                    selected: [],
                    disabledColumns: ['bk_biz_name']
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('userCustom', ['usercustom']),
            customBusinessColumns () {
                return this.usercustom['biz_table_columns'] || []
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
            customBusinessColumns () {
                this.setTableHeader()
            }
        },
        async created () {
            try {
                this.properties = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: 'biz',
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: 'post_searchObjectAttribute_biz',
                        fromCache: true
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
        methods: {
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectBiz', [
                'searchBusiness',
                'archiveBusiness',
                'updateBusiness',
                'createBusiness',
                'searchBusinessById'
            ]),
            getPropertyGroups () {
                return this.searchGroup({
                    objId: 'biz',
                    config: {
                        fromCache: true,
                        requestId: 'post_searchGroup_biz'
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            setTableHeader () {
                return new Promise((resolve, reject) => {
                    const headerProperties = this.$tools.getHeaderProperties(this.properties, this.customBusinessColumns, this.columnsConfig.disabledColumns)
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
                    id: 'bk_biz_id',
                    name: 'ID'
                }].concat(properties.map(property => {
                    return {
                        id: property['bk_property_id'],
                        name: property['bk_property_name']
                    }
                }))
            },
            handleRowClick (item) {
                this.slider.show = true
                this.slider.title = `${this.$t("Common['编辑']")} ${item['bk_biz_name']}`
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
            getBusinessList (config = {cancelPrevious: true}) {
                return this.searchBusiness({
                    params: this.getSearchParams(),
                    config: Object.assign({requestId: 'post_searchBusiness_list'}, config)
                })
            },
            getTableData () {
                this.getBusinessList().then(data => {
                    this.table.list = this.$tools.flatternList(this.properties, data.info)
                    this.table.pagination.count = data.count
                    return data
                })
            },
            getSearchParams () {
                const params = {
                    condition: {
                        'bk_data_status': {'$ne': 'disabled'}
                    },
                    fields: [],
                    page: {
                        start: this.table.pagination.size * (this.table.pagination.current - 1),
                        limit: this.table.pagination.size,
                        sort: this.table.sort
                    }
                }
                if (this.filter.id && this.filter.value) {
                    const filterType = this.filter.type
                    let filterValue = this.filter.value
                    if (filterType === 'bool') {
                        const convertValue = [true, false].find(bool => bool.toString() === filterValue)
                        filterValue = convertValue === undefined ? filterValue : convertValue
                    } else if (filterType === 'int') {
                        filterValue = isNaN(parseInt(filterValue)) ? filterValue : parseInt(filterValue)
                    }
                    params.condition[this.filter.id] = filterValue
                }
                return params
            },
            async handleEdit (flatternItem) {
                const list = await this.getBusinessList({fromCache: true})
                const inst = list.info.find(item => item['bk_biz_id'] === flatternItem['bk_biz_id'])
                const bizNameProperty = this.$tools.getProperty(this.properties, 'bk_biz_name')
                bizNameProperty.isreadonly = inst['bk_biz_name'] === '蓝鲸'
                this.attribute.inst.edit = inst
                this.attribute.type = 'update'
            },
            handleCreate () {
                this.attribute.type = 'create'
                this.attribute.inst.edit = {}
                this.slider.show = true
                this.slider.title = `${this.$t("Common['创建']")} ${this.$model['bk_obj_name']}`
            },
            handleDelete (inst) {
                this.$bkInfo({
                    title: this.$t("Common['确认要归档']", {name: inst['bk_biz_name']}),
                    confirmFn: () => {
                        this.archiveBusiness(inst['bk_biz_id']).then(() => {
                            this.slider.show = false
                            this.$success(this.$t('Common["归档成功"]'))
                            this.handlePageChange(1)
                            this.$http.cancel('post_searchBusiness_$ne_disabled')
                        })
                    }
                })
            },
            handleSave (values, changedValues, originalValues, type) {
                if (type === 'update') {
                    this.updateBusiness({
                        bizId: originalValues['bk_biz_id'],
                        params: values
                    }).then(() => {
                        this.getTableData()
                        this.searchBusinessById({bizId: originalValues['bk_biz_id']}).then(item => {
                            this.attribute.inst.details = this.$tools.flatternItem(this.properties, item)
                        })
                        this.handleCancel()
                        this.$success(this.$t("Common['修改成功']"))
                        this.$http.cancel('post_searchBusiness_$ne_disabled')
                    })
                } else {
                    this.createBusiness({
                        params: values
                    }).then(() => {
                        this.handlePageChange(1)
                        this.handleCancel()
                        this.$success(this.$t("Inst['创建成功']"))
                        this.$http.cancel('post_searchBusiness_$ne_disabled')
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
            handleApplayColumnsConfig (properties) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    'biz_table_columns': properties.map(property => property['bk_property_id'])
                })
                this.columnsConfig.show = false
            },
            handleResetColumnsConfig () {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    'biz_table_columns': []
                })
            },
            routeToHistory () {
                this.$router.push('/history/biz?relative=/business')
            },
            handleSliderBeforeClose () {
                if (this.tab.active === 'attribute' && this.attribute.type !== 'details') {
                    const $form = this.$refs.form
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
.options-button{
    font-size: 0;
    .button-history{
        border-radius: 2px 0 0 2px;
    }
    .button-setting{
        border-radius: 0 2px 2px 0;
        margin-left: -1px;
    }
}
.business-table{
    margin-top: 20px;
}
</style>