<template>
    <div class="new-association">
        <a href="javascript:void(0)" class="association-close-handle bk-icon icon-angle-double-down" @click="close"></a>
        <div class="association-filter clearfix">
            <label class="filter-label fl">{{$t('Association["关联列表"]')}}</label>
            <bk-select class="association-list-selector fl" :selected.sync="filter.objId">
                <bk-select-option v-for="(option, index) in associationOptions"
                    :key="index"
                    :label="option.label"
                    :value="option.value">
                </bk-select-option>
            </bk-select>
        </div>
        <div class="association-filter clearfix">
            <label class="filter-label fl">{{$t('Association["条件筛选"]')}}</label>
            <div class="filter-group filter-group-property fl">
                <v-property-filter
                    :excludeType="filter.objId === 'biz' ? ['singleasst', 'multiasst'] : []"
                    :objId="filter.objId"
                    @handlePropertySelected="handlePropertySelected"
                    @handleOperatorSelected="handleOperatorSelected"
                    @handleValueChange="handleValueChange">
                </v-property-filter>
            </div>
            <bk-button type="primary" class="btn-search fr" @click="search">{{$t('Association["搜索"]')}}</bk-button>
        </div>
        <v-table class="new-association-table"
            :loading="table.loading"
            :wrapperMinusHeight="400"
            :pagination.sync="table.pagination"
            :sort="table.sort"
            :header="table.header"
            :list="table.list"
            :colBorder="true"
            @handlePageChange="setCurrentPage"
            @handleSizeChange="search"
            @handleSortChange="setCurrentSort">
            <template slot="options" slot-scope="{ item }">
                <bk-button type="default" :disabled="selectedAssociationProperty && !selectedAssociationProperty['editable']" class="btn-option btn-option-remove" v-if="selectedInstId.includes(item[instanceIdKey])" @click="updateAssociation(item[instanceIdKey], 'remove')">{{$t('Association["取消关联"]')}}</bk-button>
                <bk-button type="primary" :disabled="selectedAssociationProperty && !selectedAssociationProperty['editable']" class="btn-option btn-option-new" v-else @click="updateAssociation(item[instanceIdKey], 'new')">{{$t('Association["添加关联"]')}}</bk-button>
            </template>
        </v-table>
    </div>
</template>

<script>
    import vPropertyFilter from '@/components/common/selector/property-filter'
    import vTable from '@/components/table/table'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            vPropertyFilter,
            vTable
        },
        props: {
            objId: {
                type: String,
                required: true
            },
            instance: {
                type: Object,
                defaut () {
                    return {}
                }
            }
        },
        data () {
            return {
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
                    loading: true,
                    header: [],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        size: 10
                    },
                    sort: '',
                    cancelSource: null
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
            ...mapGetters(['bkSupplierAccount']),
            ...mapGetters('object', ['attribute']),
            ...mapGetters('navigation', ['classifications']),
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
            objAttribute () {
                return this.attribute[this.objId] || []
            },
            associationOptions () {
                const validAssociation = this.association.filter(model => !['plat', 'process'].includes(model['bk_obj_id']))
                return validAssociation.map(model => {
                    return {
                        value: model['bk_obj_id'],
                        label: this.getAssociationOptionLabel(model)
                    }
                })
            },
            selectedAssociationProperty () {
                return this.objAttribute.find(property => property['bk_asst_obj_id'] === this.filter.objId)
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
            },
            searchValue () {
                let value = this.filter.property.value
                let property = this.getProperty(this.filter.property.id, this.filter.objId)
                if (property && property['bk_property_type'] === 'bool') {
                    value = ['true', 'false'].includes(this.filter.property.value) ? this.filter.property.value === 'true' : this.filter.property.value
                }
                return value
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
                    await this.$store.dispatch('object/getAttribute', {objId: filterObjId})
                    this.getInstance()
                }
            },
            'filter.property.id' (id) {
                this.setTableHeader()
            }
        },
        async created () {
            await this.$store.dispatch('object/getAttribute', {objId: this.objId})
            this.getAssociationTopo()
        },
        methods: {
            close () {
                this.$emit('handleNewAssociationClose')
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
            getAssociationTopo () {
                if (this.instance && this.instance.hasOwnProperty(this.dataIdKey)) {
                    const topoUrl = `inst/association/topo/search/owner/${this.bkSupplierAccount}/object/${this.objId}/inst/${this.instance[this.dataIdKey]}`
                    this.$axios.post(topoUrl).then(res => {
                        if (res.result) {
                            this.association = res.data.length ? res.data[0]['next'] : []
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            async updateAssociation (instId, updateType = 'new') {
                let payload = {
                    updateType: updateType,
                    objId: this.objId,
                    associated: this.selectedInstId,
                    id: this.selectedAssociationProperty['bk_property_id'],
                    value: instId,
                    multiple: this.multiple
                }
                if (this.objId === 'host') {
                    let params = {}
                    params[this.dataIdKey] = this.instance[this.dataIdKey].toString()
                    payload['params'] = params
                } else {
                    payload[this.dataIdKey] = this.instance[this.dataIdKey]
                }
                const response = await this.$store.dispatch({
                    type: 'association/updateAssociation',
                    ...payload
                })
                if (response.result) {
                    this.getAssociationTopo()
                    const msg = updateType === 'remove' ? this.$t('Association["取消关联成功"]') : this.$t('Association["添加关联成功"]')
                    this.$alertMsg(msg, 'success')
                    this.$emit('handleUpdate')
                } else {
                    this.$alertMsg(response['bk_error_msg'])
                }
            },
            async getInstance () {
                const filterObjId = this.filter.objId
                this.cancelSource && this.cancelSource.cancel()
                this.table.loading = true
                try {
                    const cancelToken = this.setCancelToken()
                    switch (filterObjId) {
                        case 'host':
                            await this.getHostInstance(cancelToken, filterObjId).then(res => this.setTableList(res, filterObjId))
                            break
                        case 'biz':
                            await this.getBizInstance(cancelToken, filterObjId).then(res => this.setTableList(res, filterObjId))
                            break
                        default:
                            await this.getObjInstance(cancelToken, filterObjId).then(res => this.setTableList(res, filterObjId))
                    }
                } catch (e) {
                    this.$alertMsg(e.message)
                } finally {
                    this.table.loading = false
                    this.cancelSource = null
                }
            },
            getHostInstance (cancelToken, filterObjId) {
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
                return this.$axios.post('hosts/search', hostParams, {cancelToken})
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
                                'value': this.searchValue
                            }]
                        })
                    } else {
                        condition[0]['condition'].push({
                            'field': this.filter.property.id,
                            'operator': this.filter.property.operator,
                            'value': this.searchValue
                        })
                    }
                }
                return condition
            },
            getBizInstance (cancelToken, filterObjId) {
                const params = {
                    condition: {
                        'bk_data_status': {'$ne': 'disabled'}
                    },
                    fields: [],
                    page: this.page
                }
                if (this.filter.property.value !== '') {
                    params.condition[this.filter.property.id] = this.searchValue
                }
                return this.$axios.post(`biz/search/${this.bkSupplierAccount}`, params, {cancelToken})
            },
            getObjInstance (cancelToken, filterObjId) {
                return this.$axios.post(`inst/association/search/owner/${this.bkSupplierAccount}/object/${filterObjId}`, {
                    condition: this.getObjCondition(),
                    fields: {},
                    page: this.page
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
                        'value': this.searchValue
                    }]
                }
                return condition
            },
            setCancelToken () {
                this.table.cancelSource = this.$Axios.CancelToken.source()
                return this.table.cancelSource.token
            },
            setTableList (res, filterObjId) {
                if (res.result) {
                    const properties = this.attribute[filterObjId] || []
                    this.table.pagination.count = res.data.count
                    if (filterObjId === 'host') {
                        res.data.info = res.data.info.map(item => item['host'])
                    }
                    this.table.list = res.data.info.map(item => this.setItem(item, properties))
                } else {
                    this.$alertMsg(res['bk_error_msg'])
                }
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
                            item[key] = this.$formatTime(item[key], type === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
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
                return this.attribute[objId].find(({bk_property_id: bkPropertyId}) => bkPropertyId === propertyId)
            },
            handlePropertySelected (data) {
                this.filter.property.id = data.value
                this.filter.property.name = data.label
            },
            handleOperatorSelected (operator) {
                this.filter.property.operator = operator
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
        .association-close-handle{
            display: block;
            height: 35px;
            line-height: 35px;
            color: $textColor;
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
        width: 280px;
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