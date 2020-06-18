<template>
    <div class="define-wrapper" v-bkloading="{
        isLoading: $loading([
            'post_searchObjectAttribute_host',
            'post_searchObjectAttribute_set',
            'post_searchObjectAttribute_module',
            'getUserAPIDetail'
            ])
        }">
        <div class="define-box">
            <div class="userapi-group">
                <label class="userapi-label">
                    {{$t("Common['业务']")}}<span class="color-danger"> * </span>
                </label>
                <cmdb-business-selector
                    class="business-selector"
                    :disabled="true"
                ></cmdb-business-selector>
            </div>
            <div class="userapi-group">
                <label class="userapi-label">
                    {{$t("CustomQuery['查询名称']")}}<span class="color-danger"> * </span>
                </label>
                <input type="text" class="cmdb-form-input"
                v-model.trim="name"
                :name="$t('CustomQuery[\'查询名称\']')"
                :disabled="!authority.includes('update')"
                v-validate="'required|max:15'">
                <span v-show="errors.has($t('CustomQuery[\'查询名称\']'))" class="color-danger">{{ errors.first($t('CustomQuery[\'查询名称\']')) }}</span>
            </div>
            <div class="userapi-group content">
                <label class="userapi-label">
                    {{$t("CustomQuery['查询内容']")}}<span class="color-danger"> * </span>
                </label>
                <div class="userapi-content-display">
                    <div class="text-content"
                        :class="{
                            open: attribute.isShow,
                            disabled: !authority.includes('update')
                        }"
                        @click="toggleContentSelector(true)">
                        <span class="default-name">{{attribute.defaultName}}</span><span v-if="selectedName.length">,{{selectedName}}</span>
                        <i class="bk-icon icon-angle-down"></i>
                    </div>
                    <div class="userapi-content-display-mask" v-if="attribute.isShow"></div>
                    <bk-selector class="fl userapi-content-selector"
                        :searchable="true"
                        search-key="bk_property_name"
                        ref="content"
                        :list="attribute.list"
                        @visible-toggle="toggleContentSelector"
                        :content-max-height="200"
                        setting-key="bk_property_id"
                        display-key="bk_property_name"
                        :selected.sync="attribute.selected"
                        :multiSelect="true"
                        :disabled="!authority.includes('update')">
                    </bk-selector>
                </div>
            </div>
            <ul class="userapi-list">
                <li v-for="(property, index) in userProperties" :key="`${property.propertyId}-${property.objId}`">
                    <label class="filter-label">
                        {{property.objName}} - {{property.propertyName}}
                    </label>
                    <div class="filter-content clearfix" :class="{disabled: !authority.includes('update')}">
                        <filter-field-operator class="filter-field-operator fl"
                            v-if="!['date', 'time'].includes(property.propertyType)"
                            :type="getOperatorType(property)"
                            :disabled="!authority.includes('update')"
                            v-model="property.operator">
                        </filter-field-operator>
                        <cmdb-form-enum class="filter-field-value filter-field-enum fl"
                            v-if="property.propertyType === 'enum'"
                            :allow-clear="true"
                            :options="getEnumOptions(property)"
                            :disabled="!authority.includes('update')"
                            v-model="property.value">
                        </cmdb-form-enum>
                        <cmdb-form-bool-input class="filter-field-value filter-field-bool-input fl"
                            v-else-if="property.propertyType === 'bool'"
                            v-model="property.value"
                            :disabled="!authority.includes('update')">
                        </cmdb-form-bool-input>
                        <cmdb-form-associate-input class="filter-field-value filter-field-associate fl"
                            v-else-if="['singleasst', 'multiasst'].includes(property.propertyType)"
                            v-model="property.value"
                            :disabled="!authority.includes('update')">
                        </cmdb-form-associate-input>
                        <component class="filter-field-value fl" :class="`filter-field-${property.propertyType}`"
                            v-else
                            :is="`cmdb-form-${property.propertyType}`"
                            :disabled="!authority.includes('update')"
                            v-model="property.value">
                        </component>
                        <i class="userapi-delete fr bk-icon icon-close"
                            v-if="authority.includes('update')"
                            @click="deleteUserProperty(property, index)">
                        </i>
                    </div>
                </li>
            </ul>
            <div class="userapi-new">
                <button class="userapi-new-btn"
                    :disabled="!authority.includes('update')"
                    @click="toggleUserAPISelector(true)">
                    {{$t("CustomQuery['新增查询条件']")}}
                </button>
                <div class="userapi-new-mask" v-if="filter.isShow"></div>
                <bk-selector class="userapi-new-selector"
                    v-if="filter.isShow"
                    :searchable="true"
                    search-key="filter_name"
                    ref="propertySelector"
                    :list="filterList"
                    @visible-toggle="toggleUserAPISelector"
                    :content-max-height="200"
                    setting-key="filter_id"
                    display-key="filter_name"
                    @item-selected="addUserProperties"
                    :selected="''">
                </bk-selector>
            </div>
            <div class="userapi-btn-group">
                <bk-button type="primary" class="userapi-btn" :disabled="errors.any()" @click.stop="previewUserAPI">
                    {{$t("CustomQuery['预览']")}}
                </bk-button>
                <bk-button type="primary" class="userapi-btn"
                    v-tooltip="$t('CustomQuery[\'保存后的查询可通过接口调用生效\']')"
                    :loading="$loading(['createCustomQuery', 'updateCustomQuery'])"
                    :disabled="errors.any() || !authority.includes('update')"
                    @click="saveUserAPI">
                    {{$t("Common['保存']")}}
                </bk-button>
                <bk-button type="default" class="userapi-btn" @click="closeSlider">
                    {{$t("Common['取消']")}}
                </bk-button>
                <bk-button type="danger" class="userapi-btn button-delete"
                    v-if="type === 'update'"
                    :loading="$loading('deleteCustomQuery')"
                    :disabled="!authority.includes('delete')"
                    @click="deleteUserAPI">
                    {{$t("Common['删除']")}}
                </bk-button>
            </div>
        </div>
        <v-preview ref="preview"
            v-if="isPreviewShow"
            :apiParams="apiParams"
            :attribute="object"
            :tableHeader="attribute.selected"
            @close="isPreviewShow = false">
        </v-preview>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import filterFieldOperator from '@/components/hosts/filter/_filter-field-operator'
    import vPreview from './preview'
    export default {
        components: {
            filterFieldOperator,
            vPreview
        },
        props: {
            type: {
                type: String
            },
            bizId: {
                type: Number
            },
            id: {
                default: ''
            },
            authority: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                name: '',
                attribute: {
                    list: [],
                    selected: [],
                    isShow: false,
                    defaultName: `${this.$t("Common['内网IP']")},${this.$t("Hosts['集群']")},${this.$t("Hosts['模块']")},${this.$t("Common['业务']")},${this.$t("Hosts['云区域']")}`,
                    default: [{
                        'bk_property_id': 'bk_host_innerip',
                        'bk_property_name': this.$t("Common['内网IP']")
                    }, {
                        'bk_property_id': 'bk_set_name',
                        'bk_property_name': this.$t("Hosts['集群']")
                    }, {
                        'bk_property_id': 'bk_module_name',
                        'bk_property_name': this.$t("Hosts['模块']")
                    }, {
                        'bk_property_id': 'bk_biz_name',
                        'bk_property_name': this.$t("Common['业务']")
                    }, {
                        'bk_property_id': 'bk_cloud_id',
                        'bk_property_name': this.$t("Hosts['云区域']")
                    }]
                },
                filter: {
                    isShow: false,
                    allList: []
                },
                object: {
                    'host': {
                        id: 'host',
                        name: this.$t("Hosts['主机']"),
                        properties: [],
                        selected: []
                    },
                    'set': {
                        id: 'set',
                        name: this.$t("Hosts['集群']"),
                        properties: [],
                        selected: []
                    },
                    'module': {
                        id: 'module',
                        name: this.$t("Hosts['模块']"),
                        properties: [],
                        selected: []
                    },
                    'biz': {
                        id: 'biz',
                        name: this.$t("Common['业务']"),
                        properties: [],
                        selected: []
                    }
                },
                userProperties: [],
                operatorMap: {
                    'time': '$in',
                    'enum': '$eq'
                },
                isPreviewShow: false,
                dataCopy: {
                    name: '',
                    userProperties: [],
                    attributeSelectd: ''
                }
            }
        },
        computed: {
            ...mapGetters([
                'supplierAccount'
            ]),
            selectedName () {
                let nameList = []
                this.attribute.selected.map(propertyId => {
                    let attr = this.attribute.list.find(({bk_property_id: bkPropertyId}) => {
                        return bkPropertyId === propertyId
                    })
                    if (attr) {
                        nameList.push(attr['bk_property_name'])
                    }
                })
                return nameList.join(',')
            },
            filterList () {
                return this.filter.allList.filter(item => {
                    if (['foreignkey'].includes(item['bk_property_type'])) {
                        return false
                    }
                    return !this.userProperties.some(property => {
                        return item['bk_obj_id'] === property.objId && item['bk_property_id'] === property.propertyId
                    })
                })
            },
            /* 生成保存自定义API的参数 */
            apiParams () {
                let paramsMap = [
                    {'bk_obj_id': 'set', condition: [], fields: []},
                    {'bk_obj_id': 'module', condition: [], fields: []},
                    {
                        'bk_obj_id': 'biz',
                        condition: [{
                            field: 'default', // 该参数表明查询非资源池下的主机
                            operator: '$ne',
                            value: 1
                        }],
                        fields: []
                    }, {
                        'bk_obj_id': 'host',
                        condition: [],
                        fields: this.attribute.selected ? this.attribute.selected : []
                    }
                ]
                const specialObj = {
                    'host': 'bk_host_innerip',
                    'biz': 'bk_biz_name',
                    'plat': 'bk_cloud_name',
                    'module': 'bk_module_name',
                    'set': 'bk_set_name'
                }
                this.userProperties.forEach((property, index) => {
                    let param = paramsMap.find(({bk_obj_id: objId}) => {
                        return objId === property.objId
                    })
                    if (property.propertyType === 'singleasst' || property.propertyType === 'multiasst') {
                        paramsMap.push({
                            'bk_obj_id': property.asstObjId,
                            fields: [],
                            condition: [{
                                field: specialObj.hasOwnProperty(property.asstObjId) ? specialObj[property.asstObjId] : 'bk_inst_name',
                                operator: property.operator,
                                value: property.value
                            }]
                        })
                    } else if (property.propertyType === 'time' || property.propertyType === 'date') {
                        let value = property['value'].split(' - ')
                        param['condition'].push({
                            field: property.propertyId,
                            operator: value[0] === value[1] ? '$eq' : '$gte',
                            value: value[0]
                        })
                        param['condition'].push({
                            field: property.propertyId,
                            operator: value[0] === value[1] ? '$eq' : '$lte',
                            value: value[1]
                        })
                    } else if (property.propertyType === 'bool' && ['true', 'false'].includes(property.value)) {
                        param['condition'].push({
                            field: property.propertyId,
                            operator: property.operator,
                            value: property.value === 'true'
                        })
                    } else if (property.operator === '$multilike') {
                        param.condition.push({
                            field: property.propertyId,
                            operator: property.operator,
                            value: property.value.split('\n').filter(str => str.trim().length).map(str => str.trim())
                        })
                    } else {
                        let operator = property.operator
                        let value = property.value
                        // 多模块与多集群查询
                        if (property.propertyId === 'bk_module_name' || property.propertyId === 'bk_set_name') {
                            operator = operator === '$regex' ? '$in' : operator
                            if (operator === '$in') {
                                let arr = value.replace('，', ',').split(',')
                                let isExist = arr.findIndex(val => {
                                    return val === value
                                }) > -1
                                value = isExist ? arr : [...arr, value]
                            }
                        }
                        param['condition'].push({
                            field: property.propertyId,
                            operator: operator,
                            value: value
                        })
                    }
                })
                let params = {
                    'bk_biz_id': this.bizId,
                    'info': {
                        condition: paramsMap
                    },
                    'name': this.name
                }
                if (this.type === 'update') {
                    params['id'] = this.id
                }
                return params
            }
        },
        watch: {
            'object.host.properties' (properties) {
                let selected = []
                let tempList = []
                properties.map(property => {
                    let isDefaultPropery = false
                    selected = this.attribute.default.map(defaultProperty => {
                        if (property['bk_property_id'] === defaultProperty['bk_property_id']) {
                            isDefaultPropery = true
                        }
                        return defaultProperty['bk_property_id']
                    })
                    if (!isDefaultPropery) {
                        tempList.push(property)
                    }
                })
                // this.attribute.list = tempList.concat(this.attribute.default)
                this.attribute.list = tempList
                this.attribute.selected = selected
                this.dataCopy.attributeSelected = this.$tools.clone(selected)
            }
        },
        async created () {
            await this.initObjectProperties()
            if (this.type !== 'create') {
                await this.getUserAPIDetail()
                this.toggleUserAPISelector(false)
            }
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            ...mapActions('hostCustomApi', [
                'getCustomQueryDetail',
                'createCustomQuery',
                'updateCustomQuery',
                'deleteCustomQuery'
            ]),
            isCloseConfirmShow () {
                if (this.name !== this.dataCopy.name || JSON.stringify(this.dataCopy.attributeSelected) !== JSON.stringify(this.attribute.selected) || this.userProperties.length !== this.dataCopy.userProperties.length) {
                    return true
                }
                return this.userProperties.some((property, index) => {
                    let propertyCopy = this.dataCopy.userProperties[index]
                    let res = false
                    for (let key in property) {
                        if (property[key] !== propertyCopy[key]) {
                            res = true
                            break
                        }
                    }
                    return res
                })
            },
            async getUserAPIDetail () {
                const res = await this.getCustomQueryDetail({
                    bizId: this.bizId,
                    id: this.id,
                    config: {
                        requestId: 'getUserAPIDetail'
                    }
                })
                this.setUserProperties(res)
            },
            setUserProperties (detail) {
                let properties = []
                let info = JSON.parse(detail['info'])
                info.condition.forEach(condition => {
                    condition['condition'].forEach(property => {
                        let originalProperty = this.getOriginalProperty(property.field, condition['bk_obj_id'])
                        if (originalProperty) {
                            if (['time', 'date'].includes(originalProperty['bk_property_type']) && properties.some(({propertyId}) => propertyId === originalProperty['bk_property_id'])) {
                                let repeatProperty = properties.find(({propertyId}) => propertyId === originalProperty['bk_property_id'])
                                repeatProperty.value = [repeatProperty.value, property.value].join(' - ')
                            } else {
                                properties.push({
                                    'objId': originalProperty['bk_obj_id'],
                                    'objName': this.object[originalProperty['bk_obj_id']].name,
                                    'propertyType': originalProperty['bk_property_type'],
                                    'propertyName': originalProperty['bk_property_name'],
                                    'propertyId': originalProperty['bk_property_id'],
                                    'asstObjId': originalProperty['bk_asst_obj_id'],
                                    'operator': property.operator,
                                    'value': this.getUserPropertyValue(property, originalProperty)
                                })
                            }
                        }
                    })
                    if (condition['bk_obj_id'] === 'host') {
                        this.attribute.selected = condition['fields']
                    }
                })
                this.userProperties = properties
                this.toggleUserAPISelector(false)
                this.name = detail['name']
                this.dataCopy = {
                    name: detail['name'],
                    userProperties: this.$tools.clone(properties),
                    attributeSelected: this.attribute.selected
                }
            },
            getUserPropertyValue (property, originalProperty) {
                if (
                    property.operator === '$in' &&
                    ['bk_module_name', 'bk_set_name'].includes(originalProperty['bk_property_id'])
                ) {
                    return property.value[property.value.length - 1]
                }
                return property.value
            },
            async previewUserAPI () {
                if (!await this.$validator.validateAll()) {
                    this.$error(this.errors.all()[0])
                    return
                }
                this.isPreviewShow = true
                this.$nextTick(() => {
                    this.$refs.preview.$el.style.zIndex = ++this.zIndex
                })
            },
            async saveUserAPI () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                let params = Object.assign({}, this.apiParams, {'info': JSON.stringify(this.apiParams['info'])})
                // 将Info字段转为JSON字符串提交
                if (this.type === 'create') {
                    const res = await this.createCustomQuery({
                        params,
                        config: {
                            requestId: 'createCustomQuery'
                        }
                    })
                    this.$success(this.$t("Common['保存成功']"))
                    this.$emit('create', res)
                } else {
                    const res = await this.updateCustomQuery({
                        bizId: this.bizId,
                        id: this.id,
                        params,
                        config: {
                            requestId: 'updateCustomQuery'
                        }
                    })
                    this.$success(this.$t("Common['修改成功']"))
                    this.$emit('update', res)
                }
                this.dataCopy = {
                    name: this.name,
                    userProperties: this.$tools.clone(this.userProperties),
                    attributeSelected: this.attribute.selected
                }
            },
            closeSlider () {
                this.$emit('cancel')
            },
            deleteUserAPI () {
                this.$bkInfo({
                    title: this.$t("CustomQuery['确认要删除']", {name: this.apiParams.name}),
                    confirmFn: async () => {
                        await this.deleteCustomQuery({
                            bizId: this.bizId,
                            id: this.id,
                            config: {
                                requestId: 'deleteCustomQuery'
                            }
                        })
                        this.$success(this.$t("Common['删除成功']"))
                        this.$emit('delete')
                        this.$emit('cancel')
                    }
                })
            },
            deleteUserProperty (userProperty, index) {
                this.userProperties.splice(index, 1)
            },
            getEnumOptions (userProperty) {
                let property = this.getOriginalProperty(userProperty.propertyId, userProperty.objId)
                if (property) {
                    return property.option || []
                }
                return []
            },
            getOperatorType (property) {
                const propertyType = property.propertyType
                const propertyId = property.propertyId
                if (['bk_set_name', 'bk_module_name'].includes(propertyId)) {
                    return 'name'
                } else if (['singlechar', 'longchar', 'singleasst', 'multiasst'].includes(propertyType)) {
                    return 'char'
                }
                return 'common'
            },
            async initObjectProperties () {
                const res = await Promise.all([
                    this.searchObjectAttribute({
                        params: {
                            bk_obj_id: 'host',
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: 'post_searchObjectAttribute_host',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: {
                            bk_obj_id: 'set',
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: 'post_searchObjectAttribute_set',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: {
                            bk_obj_id: 'module',
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: 'post_searchObjectAttribute_module',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: {
                            bk_obj_id: 'biz',
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: 'post_searchObjectAttribute_biz',
                            fromCache: true
                        }
                    })
                ])
                let hostList = res[0].filter(property => !property['bk_isapi'])
                let setList = res[1].filter(property => !property['bk_isapi'])
                let moduleList = res[2].filter(property => !property['bk_isapi'])
                hostList = hostList.map(property => {
                    return {
                        ...property,
                        ...{
                            filter_id: `${property['bk_obj_id']}-${property['bk_property_id']}`,
                            filter_name: `${this.$t("Hosts['主机']")}-${property['bk_property_name']}`
                        }
                    }
                })
                setList = setList.map(property => {
                    return {
                        ...property,
                        ...{
                            filter_id: `${property['bk_obj_id']}-${property['bk_property_id']}`,
                            filter_name: `${this.$t("Hosts['集群']")}-${property['bk_property_name']}`
                        }
                    }
                })
                moduleList = moduleList.map(property => {
                    return {
                        ...property,
                        ...{
                            filter_id: `${property['bk_obj_id']}-${property['bk_property_id']}`,
                            filter_name: `${this.$t("Hosts['模块']")}-${property['bk_property_name']}`
                        }
                    }
                })
                this.filter.allList = [...hostList, ...setList, ...moduleList]
                this.object['host']['properties'] = res[0].filter(property => !property['bk_isapi'])
                this.object['set']['properties'] = res[1].filter(property => !property['bk_isapi'])
                this.object['module']['properties'] = res[2].filter(property => !property['bk_isapi'])
                this.object['biz']['properties'] = res[3].filter(property => !property['bk_isapi'])
            },
            /* 通过选择的propertyId, 查找其对应的对象，以获得更多信息 */
            getOriginalProperty (bkPropertyId, bkObjId) {
                let property = null
                for (let objId in this.object) {
                    for (var i = 0; i < this.object[objId]['properties'].length; i++) {
                        let loopProperty = this.object[objId]['properties'][i]
                        if (loopProperty['bk_property_id'] === bkPropertyId && loopProperty['bk_obj_id'] === bkObjId) {
                            property = loopProperty
                            break
                        }
                    }
                    if (property) {
                        break
                    }
                }
                return property
            },
            addUserProperties (key, property) {
                let {
                    'bk_property_id': propertyId,
                    'bk_property_name': propertyName,
                    'bk_property_type': propertyType,
                    'bk_asst_obj_id': asstObjId,
                    'bk_obj_id': objId
                } = property
                this.userProperties.push({
                    objId,
                    propertyId,
                    propertyType,
                    propertyName,
                    objName: this.object[objId].name,
                    asstObjId,
                    operator: this.operatorMap.hasOwnProperty(propertyType) ? this.operatorMap[propertyType] : '',
                    value: ''
                })
            },
            toggleContentSelector (isShow) {
                if (this.authority.includes('update')) {
                    this.$refs.content.open = isShow
                    this.attribute.isShow = isShow
                }
            },
            toggleUserAPISelector (isPropertiesShow) {
                this.filter.isShow = isPropertiesShow
                if (isPropertiesShow) {
                    this.$nextTick(() => {
                        this.$refs.propertySelector.open = isPropertiesShow
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .define-wrapper {
        padding: 30px 15px 30px 30px;
        height: 100%;
        .define-box {
            height: 100%;
            @include scrollbar;
        }
        .userapi-group {
            margin-bottom: 15px;
            width: 370px;
            &.content {
                margin-bottom: 30px;
            }
            .userapi-label {
                display: block;
                margin-bottom: 5px;
            }
            .userapi-content-display {
                position: relative;
                .text-content {
                    border-radius: 2px;
                    border: 1px solid $cmdbBorderColor;
                    padding: 0 28px 0 16px;
                    height: 32px;
                    line-height: 30px;
                    overflow: hidden;
                    &.open {
                        padding: 5px 28px 5px 16px;
                        height: auto;
                        line-height: 20px;
                        overflow: visible;
                        border-color: $cmdbBorderFocusColor;
                        .icon-angle-down {
                            color: $cmdbBorderFocusColor;
                            transform: rotate(180deg);
                        }
                    }
                    &.disabled {
                        background-color: #fafafa;
                        color: #aaa;
                        cursor: not-allowed;
                    }
                    .default-name {
                        color: $cmdbBorderColor;
                    }
                }
                .userapi-content-display-mask {
                    position: absolute;
                    left: 0;
                    top: 0;
                    right: 0;
                    bottom: 0;
                    width: 100%;
                    height: 100%;
                }
                .icon-angle-down {
                    position: absolute;
                    right: 8px;
                    top: 8px;
                    font-size: 12px;
                    transition: transform .2s linear;
                }
            }
        }
        .userapi-list {
            width: 370px;
            .filter-label {
                display: block;
                margin-top: 20px;
                line-height: 1;
            }
            .filter-content {
                margin-top: 10px;
                width: 100%;
                &.disabled {
                    .filter-field-value {
                        width: 273px;
                    }
                }
                .content-right {
                    margin-left: 97px;
                }
                .filter-field-operator {
                    width: 87px;
                    margin-right: 10px;
                }
                .filter-field-value {
                    width: 237px;
                    &.filter-field-time,
                    &.filter-field-date {
                        width: 334px;
                    }
                }
                .userapi-delete {
                    margin: 11px 12px 0 0;
                    color: #c3cdd7;
                    cursor: pointer;
                }
            }
        }
        .userapi-new{
            width: 370px;
            margin-top: 20px;
            font-size: 14px;
            .userapi-new-btn{
                width: 100%;
                height: 32px;
                background-color: #ffffff;
                border-radius: 2px;
                border: 1px dashed #c3cdd7;
                outline: 0;
                color: $cmdbBorderFocusColor;
                &:hover{
                    box-shadow: 0px 3px 6px 0px rgba(51, 60, 72, 0.1);
                }
                &:disabled {
                    background-color: #fafafa;
                    color: #aaa;
                    cursor: not-allowed;
                }
            }
            .userapi-new-mask {
                position: absolute;
                left: 0;
                top: 0;
                right: 0;
                bottom: 0;
                width: 100%;
                height: 100%;
            }
            .userapi-pop-wrapper {
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                z-index: 99;
            }
            .userapi-new-selector-pop {
                position: absolute;
                top: calc(50% - 218px);
                right: 30px;
                padding: 30px;
                width: 370px;
                background: #fff;
                box-shadow: 0px 3px 6px 0.12px rgba(175, 177, 180, 0.61);
                border: 1px solid #fff;
                border-image: linear-gradient(#f5f5f5, #d2d4d9) 30 30;
                .pop-title {
                    margin-bottom: 20px;
                    line-height: 1;
                }
                .btn-wrapper {
                    margin-top: 20px;
                    text-align: right;
                    .bk-button {
                        min-width: 110px;
                        height: 34px;
                        line-height: 32px;
                        &:first-child {
                            margin-right: 10px;
                        }
                    }
                }
            }
        }
        .userapi-btn-group {
            position: sticky;
            margin-top: 30px;
            bottom: 0;
            left: 0;
            background: #fff;
            line-height: 36px;
            height: 37px;
            font-size: 0;
            .bk-button {
                margin-right: 10px;
            }
            .button-delete {
                background-color: #fff;
                color: #ff5656;
            }
        }
    }
</style>

<style lang="scss">
    .api-wrapper {
        .define-wrapper {
            .userapi-new {
                position: relative;
                .userapi-new-selector {
                    position: absolute;
                    left: 0;
                    bottom: 32px;
                }
                .bk-selector-wrapper {
                    display: none;
                }
                .bk-selector-list {
                    top: 36px;
                    left: 1px;
                }
            }
            .userapi-new-selector-wrapper {
                .bk-selector-wrapper {
                    display: none;
                }
                .bk-selector-list {
                    display: block !important;
                    position: static;
                    margin-top: 5px;
                    z-index: 1;
                    box-shadow: none;
                    border: solid 1px $cmdbFnMainColor;
                }
            }
            .userapi-content-display {
                .userapi-content-selector {
                    position: absolute;
                    left: 0;
                    bottom: 32px;
                }
                .bk-selector-wrapper {
                    display: none;
                }
                .bk-selector-list {
                    top: 36px;
                    left: 1px;
                }
            }
        }
    }
</style>

