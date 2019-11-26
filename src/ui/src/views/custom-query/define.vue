<template>
    <div class="define-wrapper" v-bkloading="{
        isLoading: $loading([
            'post_searchObjectAttribute_host',
            'post_searchObjectAttribute_set',
            'post_searchObjectAttribute_module',
            'getUserAPIDetail'
        ])
    }">
        <div class="define-box" ref="defineBox">
            <div class="userapi-group">
                <label class="userapi-label">
                    {{$t('业务')}}<span class="color-danger"> * </span>
                </label>
                <cmdb-business-selector
                    class="business-selector"
                    :disabled="true"
                ></cmdb-business-selector>
            </div>
            <div class="userapi-group">
                <label class="userapi-label">
                    {{$t('查询名称')}}<span class="color-danger"> * </span>
                </label>
                <cmdb-auth style="display: block;" :auth="authResources">
                    <bk-input slot-scope="{ disabled }"
                        type="text"
                        class="cmdb-form-input"
                        v-model.trim="name"
                        :name="$t('查询名称')"
                        :placeholder="$t('请输入xx', { name: $t('查询名称') })"
                        :disabled="disabled"
                        v-validate="'required|length:256'">
                    </bk-input>
                </cmdb-auth>
                <span v-show="errors.has($t('查询名称'))" class="color-danger">{{ errors.first($t('查询名称')) }}</span>
            </div>
            <div class="query-conditons" ref="queryConditions">
                <div class="query-title">
                    <span>{{$t('查询条件')}}</span>
                    <i class="icon-cc-tips" v-bk-tooltips.right="$t('针对查询内容进行条件过滤')"></i>
                </div>
                <ul class="userapi-list">
                    <li v-for="(property, index) in userProperties" :key="`${property.propertyId}-${property.objId}`">
                        <label class="filter-label">
                            {{property.objName}} - {{property.propertyName}}
                        </label>
                        <cmdb-auth style="display: block;"
                            :auth="authResources">
                            <template slot-scope="{ disabled }">
                                <div class="filter-main">
                                    <div class="filter-content clearfix" :class="{ disabled: disabled }">
                                        <filter-field-operator class="filter-field-operator fl"
                                            v-if="!['date', 'time'].includes(property.propertyType)"
                                            :type="getOperatorType(property)"
                                            :disabled="disabled"
                                            v-model="property.operator">
                                        </filter-field-operator>
                                        <component class="filter-field-value filter-field-enum fl"
                                            v-if="['list', 'enum'].includes(property.propertyType)"
                                            :is="`cmdb-form-${property.propertyType}`"
                                            v-validate="'required'"
                                            :data-vv-name="property.propertyId"
                                            :allow-clear="true"
                                            :options="getEnumOptions(property)"
                                            :disabled="disabled"
                                            v-model="property.value">
                                        </component>
                                        <cmdb-form-bool-input class="filter-field-value filter-field-bool-input fl"
                                            v-else-if="property.propertyType === 'bool'"
                                            v-model="property.value"
                                            v-validate="'required'"
                                            :data-vv-name="property.propertyId"
                                            :disabled="disabled">
                                        </cmdb-form-bool-input>
                                        <cmdb-search-input class="filter-field-value filter-field-char fl" :style="{ '--index': 99 - index }"
                                            v-else-if="['singlechar', 'longchar'].includes(property.propertyType)"
                                            v-model="property.value"
                                            v-validate="'required'"
                                            :data-vv-name="property.propertyId"
                                            :disabled="disabled">
                                        </cmdb-search-input>
                                        <cmdb-form-date-range class="filter-field-value"
                                            v-validate="'required'"
                                            :data-vv-name="property.propertyId"
                                            v-else-if="['date', 'time'].includes(property.propertyType)"
                                            v-model="property.value">
                                        </cmdb-form-date-range>
                                        <cmdb-cloud-selector
                                            v-else-if="property.propertyId === 'bk_cloud_id'"
                                            class="filter-field-value fl"
                                            v-validate="'required'"
                                            :data-vv-name="property.propertyId"
                                            :allow-clear="true"
                                            v-model="property.value">
                                        </cmdb-cloud-selector>
                                        <component class="filter-field-value fl" :class="`filter-field-${property.propertyType}`"
                                            v-else
                                            v-validate="'required'"
                                            :data-vv-name="property.propertyId"
                                            :is="`cmdb-form-${property.propertyType}`"
                                            :disabled="disabled"
                                            v-model="property.value">
                                        </component>
                                        <i class="userapi-delete fr bk-icon icon-close"
                                            v-if="!disabled"
                                            @click="deleteUserProperty(property, index)">
                                        </i>
                                    </div>
                                    <span class="error-tips" v-if="errors.has(property.propertyId)"
                                        :style="{ 'margin-left': ['date', 'time'].includes(property.propertyType) ? '0' : '120px' }">
                                        {{errors.first(property.propertyId)}}
                                    </span>
                                </div>
                            </template>
                        </cmdb-auth>
                    </li>
                </ul>
                <cmdb-auth :auth="authResources">
                    <bk-button slot-scope="{ disabled }"
                        class="add-conditon-btn"
                        theme="primary"
                        :text="true"
                        :disabled="disabled"
                        icon="icon-plus-circle"
                        @click="handleAddQueryCondition">
                        {{$t('继续添加')}}
                    </bk-button>
                </cmdb-auth>
            </div>
        </div>
        <div class="userapi-btn-group" :class="{ 'sticky': hasScrollbar }">
            <cmdb-auth :auth="authResources">
                <bk-button slot-scope="{ disabled }"
                    theme="primary"
                    class="userapi-btn"
                    v-bk-tooltips="$t('保存后的查询可通过接口调用生效')"
                    :loading="$loading(['createCustomQuery', 'updateCustomQuery'])"
                    :disabled="errors.any() || disabled"
                    @click="saveUserAPI">
                    {{type === 'create' ? $t('提交') : $t('保存')}}
                </bk-button>
            </cmdb-auth>
            <bk-button class="userapi-btn" :disabled="errors.any()" @click.stop="previewUserAPI">
                {{$t('预览')}}
            </bk-button>
            <bk-button theme="default" class="userapi-btn" @click="closeSlider">
                {{$t('取消')}}
            </bk-button>
        </div>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="propertySlider.isShow"
            :width="394"
            :title="$t('添加分组条件')"
            :before-close="handleBeforeClose">
            <property-selector slot="content"
                ref="propertySeletor"
                v-if="propertySlider.isShow"
                :properties="propertySlider.properties"
                :selected-properties="selectedProperties">
            </property-selector>
            <div slot="footer" class="property-btn-group">
                <bk-button theme="primary"
                    @click="addUserProperties">
                    {{$t('确定')}}
                </bk-button>
                <bk-button @click="handleHideQueryCondition">{{$t('取消')}}</bk-button>
            </div>
        </bk-sideslider>
        <!-- eslint-disable vue/space-infix-ops -->
        <v-preview ref="preview"
            v-if="isPreviewShow"
            :api-params="apiParams"
            :attribute="object"
            :table-header="attribute.selected"
            @close="isPreviewShow = false">
        </v-preview>
        <!-- eslint-disable end -->
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import filterFieldOperator from '@/components/hosts/filter/_filter-field-operator'
    import vPreview from './preview'
    import propertySelector from './query-property-seletor'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    export default {
        components: {
            filterFieldOperator,
            vPreview,
            propertySelector
        },
        props: {
            type: {
                type: String
            },
            bizId: {
                type: Number
            },
            id: {
                type: [String, Number],
                default: ''
            },
            object: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                name: '',
                attribute: {
                    list: [],
                    selected: [],
                    isShow: false,
                    defaultName: ['内网IP', '集群', '模块', '业务', '云区域'].map(i18n => this.$t(i18n)).join(','),
                    default: [{
                        'bk_property_id': 'bk_host_innerip',
                        'bk_property_name': this.$t('内网IP')
                    }, {
                        'bk_property_id': 'bk_set_name',
                        'bk_property_name': this.$t('集群')
                    }, {
                        'bk_property_id': 'bk_module_name',
                        'bk_property_name': this.$t('模块')
                    }, {
                        'bk_property_id': 'bk_biz_name',
                        'bk_property_name': this.$t('业务')
                    }, {
                        'bk_property_id': 'bk_cloud_id',
                        'bk_property_name': this.$t('云区域')
                    }]
                },
                filter: {
                    isShow: false,
                    allList: []
                },
                userProperties: [],
                operatorMap: {
                    'time': '$in',
                    'enum': '$eq'
                },
                isPreviewShow: false,
                dataCopy: {
                    name: '',
                    userProperties: []
                },
                propertySlider: {
                    isShow: false,
                    properties: {}
                },
                hasScrollbar: false
            }
        },
        computed: {
            ...mapGetters([
                'supplierAccount'
            ]),
            filterList () {
                return this.filter.allList.filter(item => {
                    return !this.userProperties.some(property => {
                        return item['bk_obj_id'] === property.objId && item['bk_property_id'] === property.propertyId
                    })
                })
            },
            /* 生成保存自定义API的参数 */
            apiParams () {
                const paramsMap = [
                    { 'bk_obj_id': 'set', condition: [], fields: [] },
                    { 'bk_obj_id': 'module', condition: [], fields: [] },
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
                    const param = paramsMap.find(({ bk_obj_id: objId }) => {
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
                        const value = property['value']
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
                                const arr = value.replace('，', ',').split(',')
                                const isExist = arr.findIndex(val => {
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
                const params = {
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
            },
            selectedProperties () {
                return this.userProperties.map(property => `${property.objId}-${property.propertyId}`)
            },
            authResources () {
                if (this.type === 'update') {
                    return this.$authResources({ type: this.$OPERATION.U_CUSTOM_QUERY })
                }
                return {}
            }
        },
        async created () {
            await this.initObjectProperties()
            if (this.type !== 'create') {
                await this.getUserAPIDetail()
            }
            await this.initAttributeObject()
        },
        mounted () {
            addResizeListener(this.$refs.queryConditions, this.handleResize)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.queryConditions, this.handleResize)
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
            handleResize () {
                this.$nextTick(() => {
                    const scroller = this.$refs.defineBox
                    if (scroller) {
                        this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
                    }
                })
            },
            isCloseConfirmShow () {
                if (this.name !== this.dataCopy.name || this.userProperties.length !== this.dataCopy.userProperties.length) {
                    return true
                }
                return this.userProperties.some((property, index) => {
                    const propertyCopy = this.dataCopy.userProperties[index]
                    let res = false
                    for (const key in property) {
                        if (JSON.stringify(property[key]) !== JSON.stringify(propertyCopy[key])) {
                            res = true
                            break
                        }
                    }
                    return res
                })
            },
            initAttributeObject () {
                const properties = this.object.host.properties
                let selected = []
                const tempList = []
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
                this.attribute.list = tempList
                this.attribute.selected = selected
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
                const properties = []
                const info = JSON.parse(detail['info'])
                info.condition.forEach(condition => {
                    condition['condition'].forEach(property => {
                        const originalProperty = this.getOriginalProperty(property.field, condition['bk_obj_id'])
                        if (originalProperty) {
                            if (['time', 'date'].includes(originalProperty['bk_property_type']) && properties.some(({ propertyId }) => propertyId === originalProperty['bk_property_id'])) {
                                const repeatProperty = properties.find(({ propertyId }) => propertyId === originalProperty['bk_property_id'])
                                repeatProperty.value = [repeatProperty.value, property.value]
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
                })
                this.userProperties = properties
                this.name = detail['name']
                const timer = setTimeout(() => {
                    this.dataCopy = {
                        name: detail['name'],
                        userProperties: this.$tools.clone(properties)
                    }
                    clearTimeout(timer)
                })
            },
            getUserPropertyValue (property, originalProperty) {
                if (
                    property.operator === '$in'
                    && ['bk_module_name', 'bk_set_name'].includes(originalProperty['bk_property_id'])
                ) {
                    return property.value[property.value.length - 1]
                } else if (property.operator === '$multilike' && Array.isArray(property.value)) {
                    return property.value.join('\n')
                }
                return (property.value === null || property.value === undefined) ? '' : property.value
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
                const params = Object.assign({}, this.apiParams, { 'info': JSON.stringify(this.apiParams['info']) })
                // 将Info字段转为JSON字符串提交
                if (this.type === 'create') {
                    const res = await this.createCustomQuery({
                        params,
                        config: {
                            requestId: 'createCustomQuery'
                        }
                    })
                    this.$success(this.$t('保存成功'))
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
                    this.$success(this.$t('修改成功'))
                    this.$emit('update', res)
                }
                this.dataCopy = {
                    name: this.name,
                    userProperties: this.$tools.clone(this.userProperties)
                }
            },
            closeSlider () {
                this.$emit('cancel')
            },
            deleteUserProperty (userProperty, index) {
                this.userProperties.splice(index, 1)
            },
            getEnumOptions (userProperty) {
                const property = this.getOriginalProperty(userProperty.propertyId, userProperty.objId)
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
                        params: this.$injectMetadata({
                            bk_obj_id: 'host',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'post_searchObjectAttribute_host',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'set',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'post_searchObjectAttribute_set',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'module',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'post_searchObjectAttribute_module',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'biz',
                            bk_supplier_account: this.supplierAccount
                        }),
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
                            filter_name: `${this.$t('主机')}-${property['bk_property_name']}`
                        }
                    }
                })
                setList = setList.map(property => {
                    return {
                        ...property,
                        ...{
                            filter_id: `${property['bk_obj_id']}-${property['bk_property_id']}`,
                            filter_name: `${this.$t('集群')}-${property['bk_property_name']}`
                        }
                    }
                })
                moduleList = moduleList.map(property => {
                    return {
                        ...property,
                        ...{
                            filter_id: `${property['bk_obj_id']}-${property['bk_property_id']}`,
                            filter_name: `${this.$t('模块')}-${property['bk_property_name']}`
                        }
                    }
                })
                this.filter.allList = [...hostList, ...setList, ...moduleList]
                const propertyMap = {}
                this.filter.allList.forEach(item => {
                    if (propertyMap.hasOwnProperty(item['bk_obj_id'])) {
                        propertyMap[item['bk_obj_id']].push({
                            ...item,
                            __selected__: false
                        })
                    } else {
                        propertyMap[item['bk_obj_id']] = [{
                            ...item,
                            __selected__: false
                        }]
                    }
                })
                this.propertySlider.properties = propertyMap
            },
            /* 通过选择的propertyId, 查找其对应的对象，以获得更多信息 */
            getOriginalProperty (bkPropertyId, bkObjId) {
                let property = null
                for (const objId in this.object) {
                    for (let i = 0; i < this.object[objId]['properties'].length; i++) {
                        const loopProperty = this.object[objId]['properties'][i]
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
            addUserProperties () {
                const propertySeletorElm = this.$refs.propertySeletor
                const hasChanged = propertySeletorElm.hasChanged
                if (!hasChanged) {
                    this.handleHideQueryCondition()
                    return
                }
                const addPropertyList = propertySeletorElm.addPropertyList
                const removePropertyList = propertySeletorElm.removePropertyList
                this.userProperties = this.userProperties.filter(property => !removePropertyList.includes(`${property.objId}-${property.propertyId}`))
                for (let i = 0; i < addPropertyList.length; i++) {
                    const {
                        'bk_property_id': propertyId,
                        'bk_property_name': propertyName,
                        'bk_property_type': propertyType,
                        'bk_asst_obj_id': asstObjId,
                        'bk_obj_id': objId
                    } = this.filterList.find(property => property.filter_id === addPropertyList[i].filter_id)
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
                }
                this.handleHideQueryCondition()
            },
            handleAddQueryCondition () {
                this.propertySlider.isShow = true
            },
            handleHideQueryCondition () {
                this.propertySlider.isShow = false
            },
            handleBeforeClose () {
                const hasChanged = this.$refs.propertySeletor.hasChanged
                if (hasChanged) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.handleHideQueryCondition()
                                resolve(true)
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.handleHideQueryCondition()
                return true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .define-wrapper {
        height: 100%;
        .define-box {
            max-height: calc(100% - 55px);
            padding: 18px 20px;
            overflow: auto;
        }
        .userapi-group {
            margin-bottom: 20px;
            width: 100%;
            font-size: 14px;
            .userapi-label {
                display: block;
                margin-bottom: 5px;
            }
            .business-selector {
                width: 100%;
            }
        }
        .query-conditons {
            .query-title {
                font-size: 14px;
                span {
                    color: #63656e;
                    font-weight: bold;
                }
                .icon-cc-tips {
                    color: #979ba5;
                }
            }
            .add-conditon-btn {
                margin: 10px 0 0 0;
                /deep/ .icon-plus-circle {
                    font-size: 16px;
                    margin: -2px 2px 0 0;
                }
            }
        }
        .userapi-list {
            font-size: 14px;
            .filter-label {
                display: block;
                margin-top: 20px;
                line-height: 1;
            }
            .filter-main {
                position: relative;
                .error-tips {
                    color: #ff5656;
                }
            }
            .filter-content {
                display: flex;
                margin-top: 10px;
                &:hover {
                    .userapi-delete {
                        opacity: 1;
                    }
                }
                .content-right {
                    margin-left: 97px;
                }
                .filter-field-operator {
                    flex: 110px 0 0;
                    margin-right: 10px;
                }
                .filter-field-value {
                    flex: 1;
                    width: 0;
                    &.cmdb-search-input {
                        /deep/ .search-input-wrapper {
                            z-index: var(--index);
                        }
                    }
                    &.filter-field-objuser {
                        /deep/ .suggestion-list {
                            z-index: 100 !important;
                        }
                    }
                }
                .userapi-delete {
                    width: 32px;
                    height: 32px;
                    line-height: 32px;
                    text-align: center;
                    font-size: 16px;
                    color: #C4C6CC;
                    cursor: pointer;
                    opacity: 0;
                    &:hover {
                        color: #7d8088;
                    }
                }
            }
        }
        .userapi-btn-group {
            position: sticky;
            bottom: 0;
            display: flex;
            align-items: center;
            background-color: #fff;
            font-size: 0;
            padding-left: 10px;
            &.sticky {
                border-top: 1px solid #dcdee5;
                width: 100%;
                height: 54px;
                line-height: 54px;
            }
            .bk-button {
                margin-left: 10px;
            }
            .button-delete {
                background-color: #fff;
                color: #ff5656;
                &:disabled {
                    color: #dcdee5;
                }
            }
        }
    }
    .property-btn-group {
        font-size: 0;
        padding: 0 20px;
        /deep/ .bk-button {
            margin: 0 10px 0 0;
        }
    }
</style>
