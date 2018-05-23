<template>
    <div class="new-association">
        <a href="javascript:void(0)" class="association-close-handle bk-icon icon-angle-double-down"></a>
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
                <bk-button type="default" class="btn-option btn-option-remove" v-if="selectedInstId.includes(item[instanceIdKey])" @click="updateAssociation(item[instanceIdKey], 'remove')">{{$t('Association["取消关联"]')}}</bk-button>
                <bk-button type="primary" class="btn-option btn-option-new" v-else @click="updateAssociation(item[instanceIdKey], 'new')">{{$t('Association["添加关联"]')}}</bk-button>
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
            data: {
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
                    instName: '',
                    property: {
                        id: '',
                        operator: '',
                        value: ''
                    }
                },
                filterObjProperties: [],
                table: {
                    loading: true,
                    header: [{
                        id: 'bk_inst_id',
                        name: 'ID'
                    }, {
                        id: 'bk_inst_name',
                        name: this.$t('Association["实例名"]')
                    }, {
                        id: 'options',
                        name: this.$t('Association["操作"]'),
                        sortable: false
                    }],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        size: 10
                    },
                    sort: '',
                    cancelSource: null
                },
                association: []
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            ...mapGetters('object', ['attribute']),
            ...mapGetters('navigation', ['classifications']),
            instanceIdKey () {
                if (this.filter.objId === 'host') {
                    return 'bk_host_id'
                } else if (this.filter.objId === 'biz') {
                    return 'bk_biz_id'
                }
                return 'bk_inst_id'
            },
            dataIdKey () {
                if (this.objId === 'host') {
                    return 'bk_host_id'
                } else if (this.objId === 'biz') {
                    return 'bk_biz_id'
                }
                return 'bk_inst_id'
            },
            objAttribute () {
                return this.attribute[this.objId] || []
            },
            associationOptions () {
                return this.association.map(model => {
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
            }
        },
        watch: {
            associationOptions (associationOptions) {
                const option = associationOptions.find(option => option.value === this.filter.objId)
                if (!option) {
                    this.filter.objId = associationOptions.length ? associationOptions[0]['value'] : ''
                }
            },
            'filter.objId' (filterObjId) {
                if (filterObjId) {
                    this.table.pagination.current = 1
                    this.table.pagination.count = 0
                    this.table.list = []
                    this.table.header[0]['id'] = this.instanceIdKey
                    if (filterObjId === 'host') {
                        this.table.header[1]['id'] = 'bk_host_innerip'
                        this.table.header[1]['name'] = this.$t('Common["内网IP"]')
                    } else if (filterObjId === 'biz') {
                        this.table.header[1]['id'] = 'bk_biz_name'
                        this.table.header[1]['name'] = this.$t('Association["业务名"]')
                    } else {
                        this.table.header[1]['id'] = 'bk_inst_name'
                        this.table.header[1]['name'] = this.$t('Association["实例名"]')
                    }
                    this.$store.dispatch('object/getAttribute', filterObjId)
                    this.getInstance()
                }
            },
            'filter.property.id' (id) {
                console.log(id)
                if (id !== this.instanceIdKey) {
                    const property = this.getFilterProperty()
                    const spliceLength = this.table.header.length === 3 ? 0 : 1
                    this.table.header.splice(2, spliceLength, {
                        id: id,
                        name: property['bk_property_name']
                    })
                }
            },
            data (data) {
                this.getAssociationTopo()
            }
        },
        created () {
            this.getAssociationTopo()
        },
        methods: {
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
            getAssociationTopo () {
                if (this.data && this.data.hasOwnProperty(this.dataIdKey)) {
                    const topoUrl = `inst/association/topo/search/owner/${this.bkSupplierAccount}/object/${this.objId}/inst/${this.data[this.dataIdKey]}`
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
                    params[this.instanceIdKey] = this.data[this.instanceIdKey]
                    payload['params'] = params
                } else {
                    payload[this.dataIdKey] = this.data[this.dataIdKey]
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
                    this.$alertMsg(e)
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
                const property = this.getFilterProperty()
                if (this.filter.property.value !== '' && property) {
                    if (['singleasst', 'multiasst'].includes(property['bk_property_type'])) {
                        condition.push({
                            'bk_obj_id': property['bk_asst_obj_id'],
                            'condition': [{
                                'field': 'bk_inst_name',
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
            getBizInstance (cancelToken, filterObjId) {
                return this.$axios.post(`biz/search/${this.bkSupplierAccount}`, {
                    condition: this.filter.property.value === '' ? {} : {
                        'biz': [{
                            'field': this.filter.property.id,
                            'operator': this.filter.property.operator,
                            'value': this.filter.property.value
                        }]
                    },
                    fields: [],
                    page: this.page
                }, {cancelToken})
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
                const property = this.getFilterProperty()
                if (this.filter.property.value !== '' && property) {
                    const objId = ['singleasst', 'multiasst'].includes(property['bk_property_type']) ? property['bk_asst_obj_id'] : this.filter.objId
                    condition[objId] = [{
                        'field': this.filter.property.id,
                        'operator': this.filter.property.operator,
                        'value': this.filter.property.value
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
                            item[key] = this.$formatTime(item['key'], type === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
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
            getFilterProperty () {
                const objId = this.filter.objId
                const propertyId = this.filter.property.id
                return (this.attribute[objId] || []).find(({bk_property_id: bkPropertyId}) => bkPropertyId === propertyId)
            },
            handlePropertySelected (data) {
                this.filter.property.id = data.value
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
        &-remove{
            background-color: #fff;
            color: #3c96ff;
            border-color: currentcolor;
        }
    }
    .new-association-table{
        margin: 20px 20px 0;
    }
</style>