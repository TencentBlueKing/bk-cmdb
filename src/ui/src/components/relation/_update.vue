<template>
    <div class="new-association">
        <div class="association-filter clearfix">
            <label class="filter-label fl">{{$t('Association["关联列表"]')}}</label>
            <cmdb-selector class="fl" style="width: 280px;"
                :list="options"
                setting-key="bk_obj_asst_id"
                display-key="_label"
                @on-selected="handleSelectObj">
            </cmdb-selector>
        </div>
        <div class="association-filter clearfix">
            <label class="filter-label fl">{{$t('Association["条件筛选"]')}}</label>
            <div class="filter-group filter-group-property fl">
                <cmdb-property-filter
                    :objId="currentAsstObj"
                    @on-property-selected="handlePropertySelected"
                    @on-operator-selected="handleOperatorSelected"
                    @on-value-change="handleValueChange">
                </cmdb-property-filter>
            </div>
            <bk-button type="primary" class="btn-search fr" @click="search">{{$t('Association["搜索"]')}}</bk-button>
        </div>
        <cmdb-table class="new-association-table"
            :loading="$loading()"
            :height="500"
            :pagination.sync="table.pagination"
            :sort="table.sort"
            :header="table.header"
            :list="table.list"
            :colBorder="true"
            @handlePageChange="setCurrentPage"
            @handleSizeChange="search"
            @handleSortChange="setCurrentSort">
            <template slot="options" slot-scope="{ item }">
                <a href="javascript:void(0)" class="option-link"
                    v-if="isAssociated(item)"
                    @click="updateAssociation(item[instanceIdKey], 'remove')">
                    {{$t('Association["取消关联"]')}}
                </a>
                <a href="javascript:void(0)" class="option-link" v-else
                    v-click-outside="handleCloseConfirm"
                    @click.stop="beforeUpdate($event, item[instanceIdKey], 'new')">
                    {{$t('Association["添加关联"]')}}
                </a>
            </template>
        </cmdb-table>
        <div class="confirm-tips" ref="confirmTips" v-click-outside="cancelUpdate" v-show="confirm.id">
            <p class="tips-content">{{$t('Association["更新确认"]')}}</p>
            <div class="tips-option">
                <bk-button class="tips-button" type="primary" @click="confirmUpdate">{{$t('Common["确认"]')}}</bk-button>
                <bk-button class="tips-button" type="default" @click="cancelUpdate">{{$t('Common["取消"]')}}</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import cmdbPropertyFilter from './_property-filter.vue'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            cmdbPropertyFilter
        },
        data () {
            return {
                properties: [],
                filter: {
                    id: '',
                    name: '',
                    operator: '',
                    value: ''
                },
                table: {
                    header: [],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        size: 10
                    },
                    sort: ''
                },
                specialObj: {
                    'host': 'bk_host_innerip',
                    'biz': 'bk_biz_name',
                    'plat': 'bk_cloud_name',
                    'module': 'bk_module_name',
                    'set': 'bk_set_name'
                },
                confirm: {
                    instance: null,
                    id: null
                },
                associationType: [],
                associationObject: [],
                options: [],
                currentOption: {},
                currentAsstObj: '',
                existInstAssociation: []
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            objId () {
                return this.$parent.objId
            },
            instId () {
                return this.$parent.formatedInst['bk_inst_id']
            },
            instanceIdKey () {
                const specialObj = {
                    'host': 'bk_host_id',
                    'biz': 'bk_biz_id',
                    'plat': 'bk_cloud_id',
                    'module': 'bk_module_id',
                    'set': 'bk_set_id'
                }
                if (specialObj.hasOwnProperty(this.currentAsstObj)) {
                    return specialObj[this.currentAsstObj]
                }
                return 'bk_inst_id'
            },
            instanceNameKey () {
                const nameKey = {
                    'bk_host_id': 'bk_host_innerip',
                    'bk_biz_id': 'bk_biz_name',
                    'bk_cloud_id': 'bk_cloud_name',
                    'bk_module_id': 'bk_module_name',
                    'bk_set_id': 'bk_set_name',
                    'bk_inst_id': 'bk_inst_name'
                }
                return nameKey[this.instanceIdKey]
            },
            instanceName () {
                const name = {
                    'bk_host_innerip': this.$t('Common["内网IP"]'),
                    'bk_biz_name': this.$t('Association["业务名"]'),
                    'bk_cloud_name': this.$t('Hosts["云区域"]'),
                    'bk_module_name': this.$t('Hosts["模块名"]'),
                    'bk_set_name': this.$t('Hosts["集群名"]'),
                    'bk_inst_name': this.$t('Association["实例名"]')
                }
                if (name.hasOwnProperty(this.filter.id)) {
                    return this.filter.name
                }
                return name[this.instanceNameKey]
            },
            dataIdKey () {
                const specialObj = {
                    'host': 'bk_host_id',
                    'biz': 'bk_biz_id',
                    'plat': 'bk_cloud_id',
                    'module': 'bk_module_id',
                    'set': 'bk_set_id'
                }
                if (specialObj.hasOwnProperty(this.objId)) {
                    return specialObj[this.objId]
                }
                return 'bk_inst_id'
            },
            page () {
                const pagination = this.table.pagination
                return {
                    start: (pagination.current - 1) * pagination.size,
                    limit: pagination.size,
                    sort: this.table.sort
                }
            },
            multiple () {
                return this.currentOption.mapping !== '1:1'
            }
        },
        watch: {
            'filter.id' (id) {
                this.setTableHeader(id)
            }
        },
        async created () {
            await Promise.all([
                this.getAssociationType(),
                this.getObjAssociation()
            ])
            this.setAssociationOptions()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType',
                'searchInstAssociation',
                'createInstAssociation',
                'deleteInstAssociation',
                'searchObjectAssociation'
            ]),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectCommonInst', ['searchInst']),
            ...mapActions('objectBiz', ['searchBusiness']),
            ...mapActions('hostSearch', ['searchHost']),
            getAsstObjProperties () {
                return this.searchObjectAttribute({
                    params: {
                        'bk_obj_id': this.currentAsstObj,
                        'bk_supplier_account': this.supplierAccount
                    },
                    config: {
                        requestId: `post_searchObjectAttribute_${this.currentAsstObj}`,
                        fromCache: true
                    }
                }).then(properties => {
                    this.properties = properties
                    return properties
                })
            },
            close () {
                this.$emit('on-new-relation-close')
            },
            search () {
                this.setCurrentPage(1)
            },
            setCurrentPage (page) {
                this.table.pagination.current = page
                this.getInstance()
            },
            setCurrentSort (sort) {
                this.table.sort = sort
                this.search()
            },
            setTableHeader (propertyId) {
                const header = [{
                    id: this.instanceIdKey,
                    name: 'ID'
                }, {
                    id: this.instanceNameKey,
                    name: this.instanceName
                }, {
                    id: 'options',
                    name: this.$t('Association["操作"]'),
                    sortable: false
                }]
                if (propertyId && propertyId !== this.instanceNameKey) {
                    header.splice(2, 0, {
                        id: propertyId,
                        name: (this.getProperty(propertyId) || {})['bk_property_name']
                    })
                }
                this.table.header = header
            },
            getAssociationType () {
                return this.searchAssociationType({}).then(data => {
                    this.associationType = data.info
                    return data
                })
            },
            getObjAssociation () {
                return this.searchObjectAssociation({
                    params: {
                        condition: {
                            'bk_obj_id': this.objId
                        }
                    }
                }).then(data => {
                    this.associationObject = data
                    return data
                })
            },
            setAssociationOptions () {
                const options = this.associationObject.map(option => {
                    const type = this.associationType.find(type => type['bk_asst_id'] === option['bk_asst_id'])
                    const model = this.$allModels.find(model => model['bk_obj_id'] === option['bk_asst_obj_id'])
                    return {
                        ...option,
                        '_label': `${type['src_des']}-${model['bk_obj_name']}`
                    }
                })
                this.options = options
            },
            async handleSelectObj (asstId, option) {
                this.currentOption = option
                this.currentAsstObj = option['bk_asst_obj_id']
                this.table.pagination.current = 1
                this.table.pagination.count = 0
                this.table.list = []
                this.setTableHeader()
                await Promise.all([
                    this.getAsstObjProperties(),
                    this.getExistInstAssociation()
                ])
                this.getInstance()
            },
            getExistInstAssociation () {
                const option = this.currentOption
                return this.searchInstAssociation({
                    params: {
                        condition: {
                            'bk_asst_id': option['bk_asst_id'],
                            'bk_obj_asst_id': option['bk_obj_asst_id'],
                            'bk_obj_id': this.objId,
                            'bk_asst_obj_id': option['bk_asst_obj_id']
                        }
                    }
                }).then(data => {
                    this.existInstAssociation = data
                })
            },
            isAssociated (inst) {
                return this.existInstAssociation.some(exist => exist['bk_asst_inst_id'] === inst[this.instanceIdKey])
            },
            async updateAssociation (instId, updateType = 'new') {
                if (updateType === 'new') {
                    await this.createAssociation(instId)
                    this.$success(this.$t('Association["添加关联成功"]'))
                } else if (updateType === 'remove') {
                    await this.deleteAssociation(instId)
                    this.$success(this.$t('Association["取消关联成功"]'))
                } else if (updateType === 'update') {
                    await this.deleteAssociation(this.existInstAssociation[0]['bk_asst_inst_id'])
                    await this.createAssociation(instId)
                    this.$success(this.$t('Association["添加关联成功"]'))
                }
                this.getExistInstAssociation()
            },
            createAssociation (instId) {
                return this.createInstAssociation({
                    params: {
                        'bk_obj_asst_id': this.currentOption['bk_obj_asst_id'],
                        'bk_inst_id': this.instId,
                        'bk_asst_inst_id': instId
                    }
                })
            },
            deleteAssociation (instId) {
                return this.deleteInstAssociation({
                    config: {
                        params: {
                            'bk_obj_asst_id': this.currentOption['bk_obj_asst_id'],
                            'bk_inst_id': this.instId,
                            'bk_asst_inst_id': instId
                        }
                    }
                })
            },
            beforeUpdate (event, instId, updateType = 'new') {
                if (this.multiple || !this.existInstAssociation.length) {
                    this.updateAssociation(instId, updateType)
                } else {
                    this.confirm.id = instId
                    this.confirm.instance && this.confirm.instance.destroy()
                    this.confirm.instance = this.$tooltips({
                        duration: -1,
                        theme: 'light',
                        zIndex: 9999,
                        width: 230,
                        container: document.body,
                        target: event.target
                    })
                    this.confirm.instance.$el.append(this.$refs.confirmTips)
                }
            },
            confirmUpdate () {
                this.updateAssociation(this.confirm.id, 'update')
                this.cancelUpdate()
            },
            cancelUpdate () {
                this.confirm.instance && this.confirm.instance.setVisible(false)
            },
            async getInstance () {
                const objId = this.currentAsstObj
                const config = {
                    requestId: 'get_relation_inst',
                    cancelPrevious: true
                }
                let promise
                switch (objId) {
                    case 'host':
                        promise = this.getHostInstance(config)
                        break
                    case 'biz':
                        promise = this.getBizInstance(config)
                        break
                    default:
                        promise = this.getObjInstance(objId, config)
                }
                promise.then(data => {
                    this.setTableList(data, objId)
                })
            },
            getHostInstance (config) {
                const ipFields = ['bk_host_innerip', 'bk_host_outerip']
                const filter = this.filter
                const hostParams = {
                    condition: this.getHostCondition(),
                    ip: {
                        flag: ipFields.includes(filter.id) ? filter.id : 'bk_host_innerip|bk_host_outerip',
                        exact: 0,
                        data: ipFields.includes(filter.id) && filter.value.length ? filter.value.split(',') : []
                    },
                    page: this.page
                }
                return this.searchHost({
                    params: hostParams,
                    config
                })
            },
            getHostCondition () {
                let condition = [{'bk_obj_id': 'host', 'condition': [], fields: []}]
                const property = this.getProperty(this.filter.id)
                if (this.filter.value !== '' && property) {
                    condition[0]['condition'].push({
                        'field': this.filter.id,
                        'operator': this.filter.operator,
                        'value': this.filter.value
                    })
                }
                return condition
            },
            getBizInstance (config) {
                const params = {
                    condition: {
                        'bk_data_status': {'$ne': 'disabled'}
                    },
                    fields: [],
                    page: this.page
                }
                if (this.filter.value !== '') {
                    params.condition[this.filter.id] = this.filter.value
                }
                return this.searchBusiness({
                    params,
                    config
                })
            },
            getObjInstance (objId, config) {
                return this.searchInst({
                    objId: objId,
                    params: this.getObjCondition(),
                    config
                })
            },
            getObjCondition () {
                let condition = {}
                const property = this.getProperty(this.filter.id)
                if (this.filter.value !== '' && property) {
                    const objId = this.currentAsstObj
                    condition[objId] = [{
                        'field': this.filter.id,
                        'operator': this.filter.operator,
                        'value': this.filter.value
                    }]
                }
                return condition
            },
            setTableList (data, asstObjId) {
                const properties = this.properties
                this.table.pagination.count = data.count
                if (asstObjId === 'host') {
                    data.info = data.info.map(item => item['host'])
                }
                if (asstObjId === this.objId) {
                    data.info = data.info.filter(item => item[this.instanceIdKey] !== this.instId)
                }
                this.table.list = data.info.map(item => this.$tools.flatternItem(this.properties, item))
            },
            getProperty (propertyId) {
                return this.properties.find(({bk_property_id: bkPropertyId}) => bkPropertyId === propertyId)
            },
            handleCloseConfirm () {
                this.confirm.id = null
            },
            handlePropertySelected (value, data) {
                this.filter.id = data['bk_property_id']
                this.filter.name = data['bk_property_name']
            },
            handleOperatorSelected (value, data) {
                this.filter.operator = value
            },
            handleValueChange (value) {
                this.filter.value = value
            }
        }
    }
</script>

<style lang="scss" scoped>
    .new-association{
        background-color: #fff;
        font-size: 14px;
        position: relative;
        border: 1px solid $cmdbBorderColor;
    }
    .association-filter{
        margin: 10px 20px 0;
    }
    .filter-label{
        text-align: right;
        width: 56px;
        height: 36px;
        line-height: 36px;
        margin: 0 10px 0 0;
    }
    .filter-group{
        &.filter-group-name{
            .filter-name{
                width: 170px;
            }
        }
    }
    .btn-search{
        margin: 0 0 0 8px;
    }
    .option-link{
        font-size: 12px;
        color: #3c96ff;
    }
    .new-association-table{
        margin: 20px 0 0;
        border: none;
    }
    .confirm-tips {
        padding: 9px 22px;
        .tips-content {
            color: $cmdbTextColor;
            line-height: 20px;
        }
        .tips-option {
            margin: 12px 0 0 0;
            text-align: right;
            .tips-button {
                height: 26px;
                line-height: 24px;
                padding: 0 16px;
                min-width: 56px;
                font-size: 12px;
            }
        }
    }
</style>
