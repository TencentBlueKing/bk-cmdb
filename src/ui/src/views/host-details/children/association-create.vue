<template>
    <div class="new-association">
        <div class="association-filter clearfix">
            <label class="filter-label fl">{{$t('关联列表')}}</label>
            <cmdb-selector class="fl" style="width: 280px;"
                :list="options"
                setting-key="bk_obj_asst_id"
                display-key="_label"
                @on-selected="handleSelectObj">
            </cmdb-selector>
        </div>
        <div class="association-filter clearfix">
            <label class="filter-label fl">{{$t('条件筛选')}}</label>
            <div class="filter-group filter-group-property fl">
                <cmdb-association-property-filter
                    :obj-id="currentAsstObj"
                    :exclude-type="['foreignkey']"
                    @on-property-selected="handlePropertySelected"
                    @on-operator-selected="handleOperatorSelected"
                    @on-value-change="handleValueChange">
                </cmdb-association-property-filter>
            </div>
            <bk-button theme="primary" class="btn-search fr" @click="search">{{$t('搜索')}}</bk-button>
        </div>
        <bk-table class="new-association-table"
            v-bkloading="{ isLoading: $loading() }"
            :pagination="table.pagination"
            :data="table.list"
            :col-border="true"
            :max-height="$APP.height - 210"
            @page-change="setCurrentPage"
            @page-limit-change="setCurrentLimit"
            @sort-change="setCurrentSort">
            <bk-table-column v-for="column in table.header"
                sortable="custom"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
            </bk-table-column>
            <bk-table-column :label="$t('操作')">
                <template slot-scope="{ row }">
                    <a href="javascript:void(0)" class="option-link"
                        v-if="tempData.includes(row[instanceIdKey])"
                        @click="updateAssociation(row[instanceIdKey], 'remove')">
                        {{$t('取消关联')}}
                    </a>
                    <a href="javascript:void(0)" class="option-link is-associated"
                        v-else-if="isAssociated(row)">
                        {{$t('已关联')}}
                    </a>
                    <a href="javascript:void(0)" class="option-link" v-else
                        v-click-outside="handleCloseConfirm"
                        @click.stop="beforeUpdate($event, row[instanceIdKey], 'new')">
                        {{$t('添加关联')}}
                    </a>
                </template>
            </bk-table-column>
            <cmdb-table-empty
                slot="empty"
                :stuff="table.stuff"
                :auth="$authResources({ type: tableDataPermission })">
            </cmdb-table-empty>
        </bk-table>
        <div class="confirm-tips" ref="confirmTips" v-click-outside="cancelUpdate" v-show="confirm.id">
            <p class="tips-content">{{$t('更新确认')}}</p>
            <div class="tips-option">
                <bk-button class="tips-button" theme="primary" @click="confirmUpdate">{{$t('确认')}}</bk-button>
                <bk-button class="tips-button" theme="default" @click="cancelUpdate">{{$t('取消')}}</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import cmdbAssociationPropertyFilter from './association-property-filter.vue'
    import bus from '@/utils/bus.js'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        name: 'cmdb-host-association-create',
        components: {
            cmdbAssociationPropertyFilter
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
                        limit: 10
                    },
                    sort: '',
                    stuff: {
                        type: 'search',
                        payload: {}
                    }
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
                existInstAssociation: [],
                tempData: [],
                hasChange: false
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectModelClassify', ['models']),
            objId () {
                return 'host'
            },
            instId () {
                return parseInt(this.$route.params.id)
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
                    'bk_host_innerip': this.$t('内网IP'),
                    'bk_biz_name': this.$t('业务名'),
                    'bk_cloud_name': this.$t('云区域'),
                    'bk_module_name': this.$t('模块名'),
                    'bk_set_name': this.$t('集群名'),
                    'bk_inst_name': this.$t('实例名')
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
                    start: (pagination.current - 1) * pagination.limit,
                    limit: pagination.limit,
                    sort: this.table.sort
                }
            },
            multiple () {
                return this.currentOption.mapping !== '1:1'
            },
            isSource () {
                return this.currentOption['bk_obj_id'] === this.objId
            },
            tableDataPermission () {
                const map = {
                    host: this.$OPERATION.R_HOST,
                    biz: this.$OPERATION.R_BUSINESS
                }
                return map[this.currentAsstObj] || this.$OPERATION.R_INST
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
        beforeDestroy () {
            if (this.hasChange) {
                this.hasChange = false
                bus.$emit('association-change')
            }
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
                        requestId: `post_searchObjectAttribute_${this.currentAsstObj}`
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
            setCurrentLimit (limit) {
                this.table.pagination.limit = limit
                this.search()
            },
            setCurrentSort (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.search()
            },
            setTableHeader (propertyId) {
                const header = [{
                    id: this.instanceIdKey,
                    name: 'ID'
                }, {
                    id: this.instanceNameKey,
                    name: this.instanceName
                }]
                if (propertyId && propertyId !== this.instanceNameKey) {
                    header.push({
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
                return Promise.all([
                    this.searchObjectAssociation({
                        params: {
                            condition: {
                                'bk_obj_id': this.objId
                            }
                        },
                        config: {
                            requestId: 'getSourceAssocaition'
                        }
                    }),
                    this.searchObjectAssociation({
                        params: {
                            condition: {
                                'bk_asst_obj_id': this.objId
                            }
                        },
                        config: {
                            requestId: 'getTargetAssocaition'
                        }
                    }),
                    this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
                        config: {
                            requestId: 'getMainLineModels'
                        }
                    })
                ]).then(([dataAsSource, dataAsTarget, mainLineModels]) => {
                    dataAsSource = dataAsSource || []
                    dataAsTarget = dataAsTarget || []
                    mainLineModels = mainLineModels.filter(model => !['biz', 'host'].includes(model['bk_obj_id']))
                    dataAsSource = this.getAvailableAssociation(dataAsSource, mainLineModels)
                    dataAsTarget = this.getAvailableAssociation(dataAsTarget, mainLineModels)
                    this.associationObject = [...dataAsSource, ...dataAsTarget]
                })
            },
            getAvailableAssociation (data, mainLine) {
                return data.filter(relation => {
                    return !mainLine.some(model => [relation['bk_obj_id'], relation['bk_asst_obj_id']].includes(model['bk_obj_id']))
                })
            },
            setAssociationOptions () {
                const options = this.associationObject.map(option => {
                    const isSource = option['bk_obj_id'] === this.objId
                    const type = this.associationType.find(type => type['bk_asst_id'] === option['bk_asst_id'])
                    const model = this.models.find(model => {
                        if (isSource) {
                            return model['bk_obj_id'] === option['bk_asst_obj_id']
                        } else {
                            return model['bk_obj_id'] === option['bk_obj_id']
                        }
                    })
                    return {
                        ...option,
                        '_label': `${isSource ? type['src_des'] : type['dest_des']}-${model['bk_obj_name']}`
                    }
                })
                const allLabel = options.map(option => option._label)
                const uniqueLabel = [...new Set(allLabel)]
                this.options = uniqueLabel.map(label => options.find(option => option._label === label))
            },
            async handleSelectObj (asstId, option) {
                this.tempData = []
                this.currentOption = option
                this.currentAsstObj = option['bk_obj_id'] === this.objId ? option['bk_asst_obj_id'] : option['bk_obj_id']
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
                const isSource = this.isSource
                return this.searchInstAssociation({
                    params: {
                        condition: {
                            'bk_asst_id': option['bk_asst_id'],
                            'bk_obj_asst_id': option['bk_obj_asst_id'],
                            'bk_obj_id': isSource ? this.objId : option['bk_obj_id'],
                            'bk_asst_obj_id': isSource ? option['bk_asst_obj_id'] : this.objId,
                            [`${isSource ? 'bk_inst_id' : 'bk_asst_inst_id'}`]: this.instId
                        }
                    }
                }).then(data => {
                    this.existInstAssociation = data || []
                })
            },
            isAssociated (inst) {
                return this.existInstAssociation.some(exist => {
                    if (this.isSource) {
                        return exist['bk_asst_inst_id'] === inst[this.instanceIdKey]
                    }
                    return exist['bk_inst_id'] === inst[this.instanceIdKey]
                })
            },
            async updateAssociation (instId, updateType = 'new') {
                try {
                    if (updateType === 'new') {
                        await this.createAssociation(instId)
                        this.tempData.push(instId)
                        this.$success(this.$t('添加关联成功'))
                    } else if (updateType === 'remove') {
                        await this.deleteAssociation(instId)
                        this.tempData = this.tempData.filter(tempId => tempId !== instId)
                        this.$success(this.$t('取消关联成功'))
                    } else if (updateType === 'update') {
                        await this.deleteAssociation(this.isSource ? this.existInstAssociation[0]['bk_asst_inst_id'] : this.existInstAssociation[0]['bk_inst_id'])
                        this.hasChange = true
                        this.tempData = []
                        await this.createAssociation(instId)
                        this.tempData = [instId]
                        this.$success(this.$t('添加关联成功'))
                    }
                    this.hasChange = true
                } catch (e) {
                    console.log(e)
                } finally {
                    this.getExistInstAssociation()
                }
            },
            createAssociation (instId) {
                return this.createInstAssociation({
                    params: {
                        'bk_obj_asst_id': this.currentOption['bk_obj_asst_id'],
                        'bk_inst_id': this.isSource ? this.instId : instId,
                        'bk_asst_inst_id': this.isSource ? instId : this.instId
                    }
                })
            },
            deleteAssociation (instId) {
                const instAssociation = this.existInstAssociation.find(exist => {
                    if (this.isSource) {
                        return exist['bk_asst_inst_id'] === instId
                    }
                    return exist['bk_inst_id'] === instId
                })
                return this.deleteInstAssociation({
                    id: (instAssociation || {}).id,
                    config: {
                        data: {
                            inject: !!this.$tools.getMetadataBiz(instAssociation)
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
                    this.confirm.instance = this.$bkPopover(event.target, {
                        content: this.$refs.confirmTips,
                        theme: 'light',
                        zIndex: 9999,
                        width: 230,
                        trigger: 'manual',
                        boundary: 'window',
                        arrow: true,
                        interactive: true
                    })
                    this.$nextTick(() => {
                        this.confirm.instance.show()
                    })
                }
            },
            confirmUpdate () {
                this.updateAssociation(this.confirm.id, 'update')
                this.cancelUpdate()
            },
            cancelUpdate () {
                this.confirm.instance && this.confirm.instance.hide()
            },
            async getInstance () {
                const objId = this.currentAsstObj
                const config = {
                    requestId: 'get_relation_inst',
                    globalPermission: false
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
                    this.table.stuff.type = 'search'
                    this.setTableList(data, objId)
                }).catch(e => {
                    console.error(e)
                    if (e.permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission: e.permission }
                        }
                    }
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
                const condition = [{ 'bk_obj_id': 'host', 'condition': [], fields: [] }]
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
                        'bk_data_status': { '$ne': 'disabled' }
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
                    params: this.getObjParams(),
                    config
                })
            },
            getObjParams () {
                const params = {
                    page: this.page,
                    fields: {},
                    condition: {}
                }
                const property = this.getProperty(this.filter.id)
                if (this.filter.value !== '' && property) {
                    const objId = this.currentAsstObj
                    params.condition[objId] = [{
                        'field': this.filter.id,
                        'operator': this.filter.operator,
                        'value': this.filter.value
                    }]
                }
                return params
            },
            setTableList (data, asstObjId) {
                // const properties = this.properties
                this.table.pagination.count = data.count
                if (asstObjId === 'host') {
                    data.info = data.info.map(item => item['host'])
                }
                if (asstObjId === this.objId) {
                    data.info = data.info.filter(item => item[this.instanceIdKey] !== this.instId)
                }
                this.table.list = data.info.map(item => this.$tools.flattenItem(this.properties, item))
            },
            getProperty (propertyId) {
                return this.properties.find(({ bk_property_id: bkPropertyId }) => bkPropertyId === propertyId)
            },
            handleCloseConfirm (event) {
                // this.confirm.id = null
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
        padding: 10px 20px;
        background-color: #fff;
        font-size: 14px;
        position: relative;
    }
    .association-filter{
        margin: 10px 0 0;
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
        &.is-associated {
            color: #979BA5;
            cursor: not-allowed;
        }
    }
    .new-association-table{
        margin: 20px 0 0;
    }
    .confirm-tips {
        padding: 9px;
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
