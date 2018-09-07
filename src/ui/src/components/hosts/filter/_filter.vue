<template>
    <div class="filter-layout">
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
        <div class="filter-button">
            <bk-button type="primary" @click="refresh" :disabled="$loading()">{{$t('HostResourcePool[\'刷新查询\']')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import filterFieldOperator from './_filter-field-operator.vue'
    export default {
        components: {
            filterFieldOperator
        },
        props: {
            filterConfigKey: {
                type: String,
                required: true
            }
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
                }
            }
        },
        computed: {
            ...mapGetters('userCustom', ['usercustom']),
            customFields () {
                return this.usercustom[this.filterConfigKey] || []
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
            customFieldProperties () {
                this.setCondition()
            }
        },
        async created () {
            this.setQueryParams()
            await this.getProperties()
            this.refresh()
        },
        methods: {
            ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
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
                        requestId: 'hostsAttribute',
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
                    condition[property['bk_obj_id']][property['bk_property_id']] = {
                        field: property['bk_property_id'],
                        operator: '',
                        value: ''
                    }
                })
                this.condition = condition
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .filter-layout{
        position: relative;
        height: 100%;
        padding: 0 2px;
        @include scrollbar-y;
    }
    .filter-group {
        margin-top: 20px;
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
    .filter-button{
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 76px;
        line-height: 76px;
        background-color: #fff;
    }
</style>