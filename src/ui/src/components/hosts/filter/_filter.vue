<template>
    <div class="filter-container">
        <div class="filter-layout" ref="filterLayout">
            <slot name="business"></slot>
            <div class="filter-group">
                <label for="filterIp" class="filter-label">IP</label>
                <textarea id="filterIp" class="filter-field filter-field-ip" v-model.trim="ip.text"></textarea>
                <cmdb-form-bool class="filter-field-bool" v-model="ip.bk_host_innerip" :disabled="!ip.bk_host_outerip">
                    <span class="filter-field-bool-label">{{$t('HostResourcePool[\'内网\']')}}</span>
                </cmdb-form-bool>
                <cmdb-form-bool class="filter-field-bool" v-model="ip.bk_host_outerip" :disabled="!ip.bk_host_innerip">
                    <span class="filter-field-bool-label">{{$t('HostResourcePool[\'外网\']')}}</span>
                </cmdb-form-bool>
                <cmdb-form-bool class="filter-field-bool" v-model="ip.exact" :true-value="1" :false-value="0">
                    <span class="filter-field-bool-label">{{$t('HostResourcePool[\'精确\']')}}</span>
                </cmdb-form-bool>
            </div>
            <slot name="scope"></slot>
            <div class="filter-group">
                <strong class="filter-setting-label">{{$t('Common["其他筛选项"]')}}</strong>
                <i class="filter-setting icon-cc-setting"
                    v-tooltip="$t('HostResourcePool[\'设置筛选项\']')"
                    @click="filterConfig.show = true">
                </i>
            </div>
            <div class="filter-group"
                v-for="(property, index) in customFieldProperties"
                :key="index">
                <label class="filter-label">{{getFilterLabel(property)}}</label>
                <div class="filter-field clearfix">
                    <filter-field-operator class="filter-field-operator fl"
                        v-if="!['date', 'time'].includes(property['bk_property_type'])"
                        :type="getOperatorType(property)"
                        v-model="condition[property['bk_obj_id']][property['bk_property_id']]['operator']">
                    </filter-field-operator>
                    <cmdb-form-enum class="filter-field-value fr"
                        v-if="property['bk_property_type'] === 'enum'"
                        :allow-clear="true"
                        :options="property.option || []"
                        v-model="condition[property['bk_obj_id']][property['bk_property_id']]['value']">
                    </cmdb-form-enum>
                    <cmdb-form-bool-input class="filter-field-value filter-field-bool-input fr"
                        v-else-if="property['bk_property_type'] === 'bool'"
                        v-model="condition[property['bk_obj_id']][property['bk_property_id']]['value']">
                    </cmdb-form-bool-input>
                    <cmdb-form-associate-input class="filter-field-value filter-field-associate fr"
                        v-else-if="['singleasst', 'multiasst'].includes(property['bk_property_type'])"
                        v-model="condition[property['bk_obj_id']][property['bk_property_id']]['value']">
                    </cmdb-form-associate-input>
                    <component class="filter-field-value fr" :class="`filter-field-${property['bk_property_type']}`"
                        v-else
                        :is="`cmdb-form-${property['bk_property_type']}`"
                        v-model="condition[property['bk_obj_id']][property['bk_property_id']]['value']">
                    </component>
                </div>
            </div>
            <div class="filter-button clearfix" :class="{sticky: layout.scroll}">
                <bk-button type="primary" @click="refresh" :disabled="$loading()">{{$t('Common["查询"]')}}</bk-button>
                <bk-button type="default" @click="reset">{{$t('Common["清空"]')}}</bk-button>
                <bk-button class="collection-button fr" type="default" v-if="activeSetting.includes('collection')"
                    :class="{collecting: collection.show}"
                    @click="collection.show = true">
                    <i class="icon-cc-collection"></i>
                </bk-button>
                <div class="collection-form" v-click-outside="handleCloseCollection" v-if="collection.show">
                    <div class="form-title">{{$t('Hosts[\'收藏此查询\']')}}</div>
                    <div class="form-group">
                        <input type="text" class="form-name cmdb-form-input"
                            v-validate="'required'"
                            v-model.trim="collection.name"
                            data-vv-name="collectionName"
                            :placeholder="$t('Hosts[\'请填写名称\']')">
                        <span v-show="errors.has('collectionName')" class="form-error">{{errors.first('collectionName')}}</span>
                    </div>
                    <div class="form-group">
                        <div class="form-content">{{collection.content}}</div>
                    </div>
                    <div class="form-group form-group-button">
                        <bk-button type="primary"
                            :loading="$loading('create_collection')"
                            :disabled="$loading('create_collection') || !collection.name"
                            @click="handleSaveCollection">
                            {{$t('Hosts[\'确认\']')}}
                        </bk-button>
                        <bk-button type="default" @click="handleCloseCollection">
                            {{$t('Common[\'取消\']')}}
                        </bk-button>
                    </div>
                </div>
            </div>
        </div>
        <cmdb-slider :is-show.sync="filterConfig.show" :title="$t('HostResourcePool[\'主机筛选项设置\']')" :width="600">
            <cmdb-filter-config slot="content"
                :properties="filterConfigProperties"
                :selected="customFields"
                @on-cancel="filterConfig.show = false"
                @on-apply="handleApplyFilterConfig">
            </cmdb-filter-config>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import filterFieldOperator from './_filter-field-operator.vue'
    import cmdbFilterConfig from './_filter-config.vue'
    import RESIZE_EVENTS from '@/utils/resize-events'
    export default {
        components: {
            filterFieldOperator,
            cmdbFilterConfig
        },
        props: {
            filterConfigKey: {
                type: String,
                required: true
            },
            collectionContent: {
                type: Object,
                default () {
                    return {}
                }
            },
            activeSetting: Array
        },
        data () {
            return {
                properties: {
                    biz: [],
                    host: [],
                    set: [],
                    module: []
                },
                ip: {
                    text: '',
                    'bk_host_innerip': true,
                    'bk_host_outerip': true,
                    exact: 0
                },
                condition: {},
                associateFieldMap: {
                    'biz': 'bk_biz_name',
                    'plat': 'bk_cloud_name',
                    'module': 'bk_module_name',
                    'set': 'bk_set_name'
                },
                collection: {
                    show: false,
                    name: '',
                    content: ''
                },
                filterConfig: {
                    show: false
                },
                layout: {
                    scroll: false
                }
            }
        },
        computed: {
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('hostFavorites', [
                'applying',
                'applyingInfo',
                'applyingProperties',
                'applyingConditions'
            ]),
            filterConfigProperties () {
                const properties = {}
                Object.keys(this.properties).forEach(objId => {
                    if (!['biz'].includes(objId)) {
                        properties[objId] = this.properties[objId]
                    }
                })
                return properties
            },
            customFields () {
                return this.applyingProperties.length ? this.applyingProperties : (this.usercustom[this.filterConfigKey] || [])
            },
            customFieldProperties () {
                const customFieldProperties = []
                this.customFields.forEach(field => {
                    const objId = field['bk_obj_id']
                    const propertyId = field['bk_property_id']
                    const property = this.$tools.getProperty(this.properties[objId], propertyId)
                    if (property) {
                        customFieldProperties.push(property)
                    }
                })
                return customFieldProperties
            },
            ipArray () {
                const ipArray = []
                this.ip.text.split(/\n|;|；|,|，/).forEach(text => {
                    if (text.trim().length) {
                        ipArray.push(text.trim())
                    }
                })
                return ipArray
            }
        },
        watch: {
            applyingInfo (info) {
                if (info) {
                    this.ip.text = info['ip_list'].join('\n')
                    this.ip['bk_host_innerip'] = info['bk_host_innerip']
                    this.ip['bk_host_outerip'] = info['bk_host_outerip']
                    this.ip.exact = info['exact_search']
                } else {
                    this.ip.text = ''
                    this.ip['bk_host_innerip'] = true
                    this.ip['bk_host_outerip'] = true
                    this.ip.exact = 0
                }
            },
            applyingProperties (properties) {
                let hasUnloadObj = false
                properties.forEach(property => {
                    if (!this.properties.hasOwnProperty(property['bk_obj_id'])) {
                        hasUnloadObj = true
                        this.$set(this.properties, property['bk_obj_id'], [])
                    }
                })
                if (hasUnloadObj) {
                    this.$http.cancel('hostsAttribute')
                    this.getProperties()
                }
            },
            customFieldProperties () {
                this.setCondition()
                this.updateFilterButtonStyles()
            },
            'collection.show' (show) {
                if (show) {
                    this.setCollectionContent()
                }
            }
        },
        async created () {
            this.setQueryParams()
            await this.getProperties()
            this.refresh()
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$refs.filterLayout, this.updateFilterButtonStyles)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$refs.filterLayout, this.updateFilterButtonStyles)
        },
        methods: {
            ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
            ...mapActions('hostFavorites', ['createFavorites']),
            setQueryParams () {
                const query = this.$route.query
                Object.keys(query).forEach(key => {
                    if (key === 'ip') {
                        this.ip.text = query.ip
                    } else if (key === 'inner') {
                        this.ip['bk_host_innerip'] = ['true', 'false'].includes(query.inner) ? query.inner === 'true' : !!query.inner
                    } else if (key === 'outer') {
                        this.ip['bk_host_outerip'] = ['true', 'false'].includes(query.outer) ? query.outer === 'true' : !!query.outer
                    } else if (key === 'exact') {
                        this.ip.exact = parseInt(query.exact)
                    }
                })
            },
            getProperties () {
                return this.batchSearchObjectAttribute({
                    params: {
                        bk_obj_id: {'$in': Object.keys(this.properties)},
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: `post_batchSearchObjectAttribute_${Object.keys(this.properties).join('_')}`,
                        requestGroup: Object.keys(this.properties).map(id => `post_searchObjectAttribute_${id}`),
                        fromCache: true
                    }
                }).then(result => {
                    Object.keys(this.properties).forEach(objId => {
                        this.properties[objId] = result[objId]
                    })
                    return result
                })
            },
            getFilterLabel (property) {
                const objId = property['bk_obj_id']
                const propertyModel = this.$allModels.find(model => model['bk_obj_id'] === objId)
                return `${propertyModel['bk_obj_name']} - ${property['bk_property_name']}`
            },
            getOperatorType (property) {
                const propertyType = property['bk_property_type']
                const propertyId = property['bk_property_id']
                if (['bk_set_name', 'bk_module_name'].includes(propertyId)) {
                    return 'name'
                } else if (['singlechar', 'longchar'].includes(propertyType)) {
                    return 'char'
                }
                return 'common'
            },
            getConditionOperator (property) {
                const fields = this.condition[property['bk_obj_id']]
                const target = fields.find(field => field.field === property['bk_property_id'])
                return target.operator
            },
            getParams () {
                const params = {
                    ip: {
                        flag: ['bk_host_innerip', 'bk_host_outerip'].filter(flag => this.ip[flag]).join('|'),
                        exact: this.ip.exact,
                        data: this.ipArray
                    },
                    condition: []
                }
                // 填充必要模型查询参数
                const requiredObj = ['biz', 'host', 'set', 'module']
                const normalProperties = this.customFieldProperties.filter(property => !['singleasst', 'multiasst'].includes(property['bk_property_type']))
                requiredObj.forEach(objId => {
                    const objParams = {
                        'bk_obj_id': objId,
                        condition: [],
                        fields: []
                    }
                    const objProperties = normalProperties.filter(property => property['bk_obj_id'] === objId)
                    objProperties.forEach(property => {
                        const propertyCondition = this.condition[objId][property['bk_property_id']]
                        // 必要模型参数合法时，填充对应模型的condition
                        let value = propertyCondition.value
                        if (!['', null].includes(value)) {
                            if (propertyCondition.operator === '$in') {
                                let splitValue = [...(new Set(value.split(',').map(val => val.trim())))]
                                value = splitValue.length > 1 ? [...splitValue, value] : splitValue
                            }
                            objParams.condition.push({
                                ...propertyCondition,
                                value
                            })
                        }
                    })
                    params.condition.push(objParams)
                })
                // 关联属性额外填充模型查询参数
                const associateProperties = this.customFieldProperties.filter(property => ['singleasst', 'multiasst'].includes(property['bk_property_type']))
                associateProperties.forEach(property => {
                    const associateObjId = property['bk_asst_obj_id']
                    const propertyCondition = this.condition[property['bk_obj_id']][property['bk_property_id']]
                    // 关联模型存在且查询参数合法时，填充对应关联模型的condition
                    if (associateObjId && !['', null].includes(propertyCondition.value)) {
                        let objParams = params.condition.find(condition => condition['bk_obj_id'] === associateObjId)
                        if (!objParams) {
                            objParams = {
                                'bk_obj_id': associateObjId,
                                condition: [],
                                fields: []
                            }
                            params.condition.push(objParams)
                        }
                        objParams.condition.push({
                            field: this.associateFieldMap[associateObjId] || 'bk_inst_name',
                            operator: propertyCondition.operator,
                            value: propertyCondition.value
                        })
                    }
                })
                return params
            },
            setCondition () {
                const condition = {
                    host: {},
                    set: {},
                    module: {}
                }
                this.customFieldProperties.forEach(property => {
                    condition[property['bk_obj_id']][property['bk_property_id']] = this.getPropertyCondition(property)
                })
                this.condition = condition
            },
            getPropertyCondition (property) {
                const objId = property['bk_obj_id']
                const propertyId = property['bk_property_id']
                const condition = {
                    field: propertyId,
                    operator: '',
                    value: ''
                }
                const collectionConditon = (this.applyingConditions[objId] || []).find(condition => condition.field === propertyId)
                if (collectionConditon) {
                    condition.operator = collectionConditon.operator
                    condition.value = collectionConditon.value
                }
                return condition
            },
            reset () {
                this.ip = {
                    text: '',
                    'bk_host_innerip': true,
                    'bk_host_outerip': true,
                    exact: 0
                }
                for (let objId in this.condition) {
                    for (let propertyId in this.condition[objId]) {
                        this.condition[objId][propertyId].operator = ''
                        this.condition[objId][propertyId].value = ''
                    }
                }
                this.refresh()
            },
            refresh () {
                this.$emit('on-refresh', this.getParams())
            },
            setCollectionContent () {
                const content = []
                const params = this.getParams()
                if (this.collectionContent.hasOwnProperty('business')) {
                    content.push(`bk_biz_id:${this.collectionContent.business}`)
                }
                if (params.ip.data.length) {
                    content.push(`ip:${params.ip.data.join(',')}`)
                }
                const operatorMap = {
                    '$ne': '!=',
                    '$eq': '=',
                    '$regex': '~',
                    '$in': '~'
                }
                params.condition.forEach(({condition, bk_obj_id: objId}) => {
                    if (!['biz'].includes(objId) && condition.length) {
                        const objContent = []
                        condition.forEach(({field, operator, value}) => {
                            objContent.push(`${field}${operatorMap[operator]}${Array.isArray(value) ? value.join(',') : value}`)
                        })
                        content.push(`${objId}: ${objContent.join(' | ')}`)
                    }
                })
                this.collection.content = content.join(' | ')
            },
            handleSaveCollection () {
                this.$validator.validate('collectionName').then(result => {
                    if (result) {
                        this.createFavorites({
                            params: this.getCollectionParams(),
                            config: {
                                requestId: 'create_collection'
                            }
                        }).then(() => {
                            this.handleCloseCollection()
                            this.$success(this.$t('Common["收藏成功"]'))
                        })
                    }
                })
            },
            getCollectionParams () {
                const params = this.getParams()
                const info = {
                    'bk_biz_id': this.collectionContent.business || -1,
                    'exact_search': this.ip.exact,
                    'bk_host_innerip': this.ip['bk_host_innerip'],
                    'bk_host_outerip': this.ip['bk_host_outerip'],
                    'ip_list': params.ip.data
                }
                const queryParams = []
                params.condition.forEach(({condition, bk_obj_id: objId}) => {
                    condition.forEach(({field, operator, value}) => {
                        queryParams.push({
                            'bk_obj_id': objId,
                            field,
                            operator,
                            value: Array.isArray(value) ? value.join(',') : value
                        })
                    })
                })
                return {
                    name: this.collection.name,
                    info: JSON.stringify(info),
                    'query_params': JSON.stringify(queryParams),
                    'is_default': 2
                }
            },
            handleCloseCollection () {
                this.collection.show = false
                this.collection.name = ''
                this.collection.content = ''
            },
            handleApplyFilterConfig (properties) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.filterConfigKey]: properties.map(property => {
                        return {
                            'bk_property_id': property['bk_property_id'],
                            'bk_obj_id': property['bk_obj_id']
                        }
                    })
                }).then(() => {
                    this.$store.commit('hostFavorites/setApplying', null)
                })
                this.filterConfig.show = false
            },
            updateFilterButtonStyles () {
                this.$nextTick(() => {
                    const $filterLayout = this.$refs.filterLayout
                    this.layout.scroll = $filterLayout.offsetHeight !== $filterLayout.scrollHeight
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .filter-container {
        height: 100%;
        position: relative;
    }
    .filter-layout{
        position: relative;
        height: 100%;
        @include scrollbar-y;
    }
    .filter-group {
        margin: 20px 0 0 0;
        padding: 0 20px;
        position: relative;
        background-color: #fff;
        .filter-label{
            display: block;
            font-size: 14px;
        }
        .filter-field{
            margin-top: 10px;
            &-ip{
                width: 100%;
                height: 70px;
                padding: 10px;
                border-radius: 2px;
                border: 1px solid $cmdbBorderColor;
                font-size: 14px;
                outline: none;
                resize: none;
                &:focus{
                    border-color: $cmdbBorderFocusColor;
                }
            }
            &-bool{
                margin-right: 15px;
            }
            &-bool-label{
                display: inline-block;
                vertical-align: middle;
                margin: 0 0 0 5px;
                font-size: 14px;
            }
            &-operator{
                width: 77px;
            }
            &-value{
                width: 224px;
            }
            &-date,
            &-time{
                width: 100%;
            }
        }
    }
    .filter-setting-label {
        display: inline-block;
        vertical-align: middle;
        font-size: 14px;
    }
    .filter-setting {
        display: inline-block;
        vertical-align: middle;
        cursor: pointer;
        font-size: 16px;
    }
    .filter-button{
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 20px 20px 0;
        background-color: #fff;
        &.sticky {
            border-top: 1px solid $cmdbBorderColor;
            padding: 10px 20px 0;
        }
        .collection-button.collecting {
            color: #ffb400;
        }
    }
    .collection-form {
        position: absolute;
        bottom: 100%;
        right: 0;
        width: 100%;
        margin: 0 0 10px 0;
        padding: 20px;
        border-radius: 2px;
        box-shadow: 0 2px 10px 4px rgba(12,34,59,.13);
        color: #3c96ff;
        z-index: 9999;
        background-color: #fff;
        &:before,
        &:after {
            position: absolute;
            right: 20px;
            top: 100%;
            width: 0;
            height: 0;
            content: "";
            border-left: 6px solid transparent;
            border-right: 6px solid transparent;
        }
        &:before {
            border-top: 10px solid #e7e9ef;
            margin-top: 2px;
        }
        &:after {
            border-top: 10px solid #fff;
        }
        .form-group {
            margin: 15px 0 0 0;
            position: relative;
            &-button {
                text-align: right;
            }
        }
        .form-title {
            border-left: 2px solid #6b7baa;
            padding-left: 5px;
            font-size: 12px;
        }
        .form-name {
            color: $cmdbTextColor;    
        }
        .form-error {
            position: absolute;
            top: 100%;
            left: 0;
            color: $cmdbDangerColor;
            font-size: 12px;
        }
        .form-content {
            min-height: 70px;
            max-height: 300px;
            font-size: 14px;
            padding: 10px;
            word-break: break-all;
            background-color: #f9f9f9;
            @include scrollbar-y;
        }
    }
</style>