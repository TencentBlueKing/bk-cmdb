<template>
    <div class="new-association">
        <a href="javascript:void(0)" class="association-close-handle bk-icon icon-angle-double-down" @click="close"></a>
        <div class="association-filter clearfix">
            <label class="filter-label fl">{{$t('Association["关联列表"]')}}</label>
            <cmdb-selector class="association-list-selector fl"
                :list="associationOptions"
                setting-key="value"
                display-key="label"
                v-model="filter.objId">
            </cmdb-selector>
        </div>
        <div class="association-filter clearfix">
            <label class="filter-label fl">{{$t('Association["条件筛选"]')}}</label>
            <div class="filter-group filter-group-property fl">
                <cmdb-property-filter
                    :excludeType="filter.objId === 'biz' ? ['singleasst', 'multiasst'] : []"
                    :objId="filter.objId"
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
                <bk-button type="default" class="btn-option btn-option-remove"
                    v-if="selectedInstId.includes(item[instanceIdKey])"
                    :disabled="selectedAssociationProperty && !selectedAssociationProperty['editable']"
                    @click="updateAssociation(item[instanceIdKey], 'remove')">
                    {{$t('Association["取消关联"]')}}
                </bk-button>
                <bk-button type="primary" class="btn-option btn-option-new" v-else
                    :disabled="selectedAssociationProperty && !selectedAssociationProperty['editable']"
                    @click="updateAssociation(item[instanceIdKey], 'new')">
                    {{$t('Association["添加关联"]')}}
                </bk-button>
            </template>
        </cmdb-table>
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
                properties: {},
                filter: {
                    objId: '',
                    property: {
                        id: '',
                        name: '',
                        operator: '',
                        value: ''
                    }
                },
                filterObjProperties: [],
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
                association: [],
                specialObj: {
                    'host': 'bk_host_innerip',
                    'biz': 'bk_biz_name',
                    'plat': 'bk_cloud_name',
                    'module': 'bk_module_name',
                    'set': 'bk_set_name'
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectModelClassify', ['classifications']),
            objId () {
                return this.$parent.objId
            },
            instId () {
                return this.$parent.instId
            },
            instanceIdKey () {
                const specialObj = {
                    'host': 'bk_host_id',
                    'biz': 'bk_biz_id',
                    'plat': 'bk_cloud_id',
                    'module': 'bk_module_id',
                    'set': 'bk_set_id'
                }
                if (specialObj.hasOwnProperty(this.filter.objId)) {
                    return specialObj[this.filter.objId]
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
                if (name.hasOwnProperty(this.filter.property.id)) {
                    return this.filter.property.name
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
            associationOptions () {
                const validAssociation = this.association.filter(model => !['plat', 'process', 'module', 'set'].includes(model['bk_obj_id']))
                return validAssociation.map(model => {
                    return {
                        value: model['bk_obj_id'],
                        label: this.getAssociationOptionLabel(model)
                    }
                })
            },
            selectedAssociationProperty () {
                return this.properties[this.objId].find(property => property['bk_asst_obj_id'] === this.filter.objId)
            },
            multiple () {
                return this.selectedAssociationProperty && this.selectedAssociationProperty['bk_property_type'] === 'multiasst'
            },
            selectedInstId () {
                const filterObjId = this.filter.objId
                const filterObject = this.association.find(obj => obj['bk_obj_id'] === filterObjId)
                if (filterObject && filterObject.count) {
                    return filterObject.children.map(({bk_inst_id: bkInstId}) => bkInstId)
                }
                return []
            },
            page () {
                const pagination = this.table.pagination
                return {
                    start: (pagination.current - 1) * pagination.size,
                    limit: pagination.size,
                    sort: this.table.sort
                }
            }
        },
        watch: {
            associationOptions (associationOptions) {
                const option = associationOptions.find(option => option.value === this.filter.objId)
                if (!option) {
                    this.filter.objId = associationOptions.length ? associationOptions[0]['value'] : ''
                }
            },
            async 'filter.objId' (filterObjId) {
                if (filterObjId) {
                    this.table.pagination.current = 1
                    this.table.pagination.count = 0
                    this.table.list = []
                    await this.getObjProperties(filterObjId)
                    this.getInstance()
                }
            },
            'filter.property.id' (id) {
                this.setTableHeader()
            }
        },
        async created () {
            await this.getObjProperties(this.objId)
            this.getAssociationTopo()
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectRelation', ['getInstRelation', 'updateInstRelation']),
            ...mapActions('objectCommonInst', ['searchInst']),
            ...mapActions('objectBiz', ['searchBusiness']),
            ...mapActions('hostSearch', ['searchHost']),
            getObjProperties (objId) {
                return this.searchObjectAttribute({
                    params: {
                        'bk_obj_id': objId,
                        'bk_supplier_account': this.supplierAccount
                    },
                    config: {
                        requestId: `get_${objId}_attribute`,
                        fromCache: true
                    }
                }).then(properties => {
                    this.$set(this.properties, objId, properties)
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
            setTableHeader () {
                const filterObjId = this.filter.objId
                const filterPropertyId = this.filter.property.id
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
                if (filterPropertyId !== this.instanceNameKey) {
                    header.splice(2, 0, {
                        id: filterPropertyId,
                        name: this.getProperty(filterPropertyId, filterObjId)['bk_property_name']
                    })
                }
                this.table.header = header
            },
            getAssociationTopo (config = {}) {
                this.getInstRelation({
                    objId: this.objId,
                    instId: this.instId,
                    config: {
                        requestId: `get_${this.objId}_${this.instId}_relation`,
                        fromCache: true,
                        ...config
                    }
                }).then(data => {
                    this.association = data[0].next
                })
            },
            async updateAssociation (instId, updateType = 'new') {
                let payload = {
                    updateType: updateType,
                    objId: this.objId,
                    relation: this.selectedInstId,
                    id: this.selectedAssociationProperty['bk_property_id'],
                    value: instId,
                    multiple: this.multiple
                }
                if (this.objId === 'host') {
                    let params = {}
                    params[this.dataIdKey] = this.instId.toString()
                    payload['params'] = params
                } else {
                    payload[this.dataIdKey] = this.instId
                }
                const response = await this.updateInstRelation({
                    params: payload
                })
                this.getAssociationTopo({clearCache: true})
                const msg = updateType === 'remove' ? this.$t('Association["取消关联成功"]') : this.$t('Association["添加关联成功"]')
                this.$success(msg)
            },
            async getInstance () {
                const filterObjId = this.filter.objId
                const config = {
                    requestId: 'get_relation_inst',
                    cancelPrevious: true
                }
                let promise
                switch (filterObjId) {
                    case 'host':
                        promise = this.getHostInstance(filterObjId, config)
                        break
                    case 'biz':
                        promise = this.getBizInstance(filterObjId, config)
                        break
                    default:
                        promise = this.getObjInstance(filterObjId, config)
                }
                promise.then(data => {
                    this.setTableList(data, filterObjId)
                })
            },
            getHostInstance (filterObjId, config) {
                const ipFields = ['bk_host_innerip', 'bk_host_outerip']
                const filterProperty = this.filter.property
                const hostParams = {
                    condition: this.getHostCondition(),
                    ip: {
                        flag: ipFields.includes(filterProperty.id) ? filterProperty.id : 'bk_host_innerip|bk_host_outerip',
                        exact: 0,
                        data: ipFields.includes(filterProperty.id) && filterProperty.value.length ? filterProperty.value.split(',') : []
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
                const property = this.getProperty(this.filter.property.id, this.filter.objId)
                if (this.filter.property.value !== '' && property) {
                    if (['singleasst', 'multiasst'].includes(property['bk_property_type'])) {
                        condition.push({
                            'bk_obj_id': property['bk_asst_obj_id'],
                            'condition': [{
                                'field': this.specialObj.hasOwnProperty(property['bk_asst_obj_id']) ? this.specialObj[property['bk_asst_obj_id']] : 'bk_inst_name',
                                'operator': this.filter.property.operator,
                                'value': this.filter.property.value
                            }]
                        })
                    } else {
                        condition[0]['condition'].push({
                            'field': this.filter.property.id,
                            'operator': this.filter.property.operator,
                            'value': this.filter.property.value
                        })
                    }
                }
                return condition
            },
            getBizInstance (filterObjId, config) {
                const params = {
                    condition: {
                        'bk_data_status': {'$ne': 'disabled'}
                    },
                    fields: [],
                    page: this.page
                }
                if (this.filter.property.value !== '') {
                    params.condition[this.filter.property.id] = this.filter.property.value
                }
                return this.searchBusiness({
                    params,
                    config
                })
            },
            getObjInstance (filterObjId, config) {
                return this.searchInst({
                    objId: filterObjId,
                    params: this.getObjCondition(),
                    config
                })
            },
            getObjCondition () {
                let condition = {}
                const property = this.getProperty(this.filter.property.id, this.filter.objId)
                if (this.filter.property.value !== '' && property) {
                    const objId = ['singleasst', 'multiasst'].includes(property['bk_property_type']) ? property['bk_asst_obj_id'] : this.filter.objId
                    condition[objId] = [{
                        'field': this.specialObj.hasOwnProperty(property['bk_asst_obj_id']) ? this.specialObj[property['bk_asst_obj_id']] : this.filter.property.id,
                        'operator': this.filter.property.operator,
                        'value': this.filter.property.value
                    }]
                }
                return condition
            },
            setTableList (data, filterObjId) {
                const properties = this.properties[filterObjId]
                this.table.pagination.count = data.count
                if (filterObjId === 'host') {
                    data.info = data.info.map(item => item['host'])
                }
                this.table.list = data.info.map(item => this.setItem(item, properties))
            },
            setItem (item, properties) {
                for (let key in item) {
                    const property = properties.find(({bk_property_id: bkPropertyId}) => bkPropertyId === key)
                    if (property) {
                        const type = property['bk_property_type']
                        if (['singleasst', 'multiasst'].includes(type) && Array.isArray(item[key])) {
                            item[key] = item[key].map(({bk_inst_name: bkInstName}) => bkInstName).join(',')
                        } else if (['enum'].includes(type) && Array.isArray(property.option)) {
                            const option = property.option.find(({id}) => id === item[key])
                            item[key] = option ? option.name : ''
                        } else if (['date', 'time'].includes(type)) {
                            item[key] = this.$tools.formatTime(item[key], type === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
                        }
                    }
                }
                return item
            },
            getAssociationOptionLabel (model) {
                let label = ''
                for (let i = 0; i < this.classifications.length; i++) {
                    const modelInClassification = this.classifications[i]['bk_objects'].find(({bk_obj_id: bkObjId}) => bkObjId === model['bk_obj_id'])
                    if (modelInClassification) {
                        label = `${this.classifications[i]['bk_classification_name']}-${model['bk_obj_name']}`
                        break
                    }
                }
                return label
            },
            getProperty (propertyId, objId) {
                return this.properties[objId].find(({bk_property_id: bkPropertyId}) => bkPropertyId === propertyId)
            },
            handlePropertySelected (value, data) {
                this.filter.property.id = data['bk_property_id']
                this.filter.property.name = data['bk_property_name']
            },
            handleOperatorSelected (value, data) {
                this.filter.property.operator = value
            },
            handleValueChange (value) {
                this.filter.property.value = value
            }
        }
    }
</script>

<style lang="scss" scoped>
    .new-association{
        background-color: #fff;
        font-size: 14px;
        position: relative;
        top: -54px;
        .association-close-handle{
            display: block;
            height: 35px;
            line-height: 35px;
            color: $cmdbTextColor;
            margin: 0 0 25px 0;
            text-align: center;
            background-image: linear-gradient(#f9f9f9, #fff);
        }
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
    .association-list-selector{
        width: 280px !important;
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
    .btn-option{
        height: 22px;
        line-height: 20px;
        font-size: 12px;
        padding: 0 8px;
        &-new{
            background-color: #30d878;
            border-color: #30d878;
            &:hover{
                background-color: #10ed6f;
                border-color: #10ed6f;
            }
        }
        &-remove{
            background-color: #fff;
            color: #3c96ff;
            border-color: currentcolor;
            &:hover{
                color: #0082ff;
            }
        }
    }
    .new-association-table{
        margin: 20px 20px 0;
    }
</style>