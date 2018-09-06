<template>
    <div class="define-wrapper">
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
                <input type="text" class="cmdb-form-input">
            </div>
            <div class="userapi-group content">
                <label class="userapi-label">
                    {{$t("CustomQuery['查询内容']")}}<span class="color-danger"> * </span>
                </label>
                <div class="userapi-content-display clearfix">
                    <textarea class="userapi-textarea" v-model="selectedName" disabled name="" id="" cols="30" rows="10"></textarea>
                    <bk-button :disabled="attribute.isShow" v-tooltip="$t('Common[\'新增\']')" type="primary" class="btn-icon icon-cc-plus" @click="toggleContentSelector(true)"></bk-button>
                </div>
                <div class="content-selector" v-show="attribute.isShow" ref="userapiContentSelector">
                    <bk-selector class="fl userapi-content-selector"
                        ref="content"
                        :list="attribute.list"
                        @visible-toggle="toggleContentSelector"
                        setting-key="bk_property_id"
                        display-key="bk_property_name"
                        :selected.sync="attribute.selected"
                        :multiSelect="true">
                    </bk-selector>
                </div>
            </div>
            <ul class="userapi-list">
                <li v-for="(property, index) in userProperties" :key="`${property.propertyId}-${property.objId}`">
                    <label class="filter-label">
                        {{property.objName}} - {{property.propertyName}}
                    </label>
                    <div class="filter-content clearfix">
                        <cmdb-form-enum class="filter-field-value fr"
                            v-if="property.propertyType === 'enum'"
                            :allow-clear="true"
                            :options="property.option || []"
                            v-model="property.value">
                        </cmdb-form-enum>
                        <cmdb-form-bool-input class="filter-field-value filter-field-bool-input fr"
                            v-else-if="property.propertyType === 'bool'"
                            v-model="property.value">
                        </cmdb-form-bool-input>
                        <cmdb-form-associate-input class="filter-field-value filter-field-associate fr"
                            v-else-if="['singleasst', 'multiasst'].includes(property.propertyType)"
                            v-model="property.value">
                        </cmdb-form-associate-input>
                        <component class="filter-field-value fr" :class="`filter-field-${property.propertyType}`"
                            v-else
                            :is="`cmdb-form-${property.propertyType}`"
                            v-model="property.value">
                        </component>
                    </div>
                </li>
            </ul>
            <div class="userapi-new">
                <button class="userapi-new-btn" @click="toggleUserAPISelector(true)">{{$t("CustomQuery['新增查询条件']")}}</button>
                <div class="userapi-pop-wrapper" ref="userapiPop">
                    <div class="userapi-new-selector-pop" v-show="objectInfo.isPropertiesShow">
                        <p class="pop-title">{{$t("CustomQuery['新增查询条件']")}}</p>
                        <bk-selector class="userapi-new-selector" 
                            :list="objectInfo.list"
                            :selected.sync="objectInfo.selected">
                        </bk-selector>
                        <div class="userapi-new-selector-wrapper">
                            <bk-selector
                                ref="propertySelector"
                                :list="object[objectInfo.selected]['properties']"
                                :selected.sync="propertySelected[objectInfo.selected]"
                                setting-key="bk_property_id"
                                display-key="bk_property_name"
                                :multiSelect="true">
                            </bk-selector>
                        </div>
                        <div class="btn-wrapper">
                            <bk-button type="primary" class="btn confirm" @click="addUserProperties">{{$t("Common['确定']")}}</bk-button>
                            <bk-button type="default" class="btn vice-btn" @click="toggleUserAPISelector(false)">{{$t("Common['取消']")}}</bk-button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
        },
        data () {
            return {
                name: '',
                attribute: {
                    list: [],
                    selected: [],
                    isShow: false,
                    default: [{
                        'bk_property_id': 'bk_host_innerip',
                        'bk_property_name': this.$t("Common['内网IP']")
                    }, {
                        'bk_property_id': 'bk_biz_name',
                        'bk_property_name': this.$t("Common['业务']")
                    }, {
                        'bk_property_id': 'bk_set_name',
                        'bk_property_name': this.$t("Hosts['集群']")
                    }, {
                        'bk_property_id': 'bk_module_name',
                        'bk_property_name': this.$t("Hosts['模块']")
                    }, {
                        'bk_property_id': 'bk_cloud_id',
                        'bk_property_name': this.$t("Hosts['云区域']")
                    }]
                },
                objectInfo: {
                    isPropertiesShow: false,
                    selected: 'host',
                    list: [{
                        id: 'host',
                        name: this.$t("Hosts['主机']")
                    }, {
                        id: 'set',
                        name: this.$t("Hosts['集群']")
                    }, {
                        id: 'module',
                        name: this.$t("Hosts['模块']")
                    }]
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
                    }
                },
                propertySelected: {
                    host: [],
                    set: [],
                    module: []
                },
                userProperties: [],
                operatorMap: {
                    'time': '$in',
                    'enum': '$eq'
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
            }
        },
        watch: {
            'object.host.properties' (properties) {
                let selected = []
                let tempList = []
                console.log(properties)
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
                this.attribute.list = tempList.concat(this.attribute.default)
                this.attribute.selected = selected
            }
        },
        created () {
            this.initObjectProperties()
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            async initObjectProperties () {
                const res = await Promise.all([
                    this.searchObjectAttribute({
                        params: {
                            bk_obj_id: 'host',
                            bk_supplier_account: this.supplierAccount
                        }
                    }),
                    this.searchObjectAttribute({
                        params: {
                            bk_obj_id: 'set',
                            bk_supplier_account: this.supplierAccount
                        }
                    }),
                    this.searchObjectAttribute({
                        params: {
                            bk_obj_id: 'module',
                            bk_supplier_account: this.supplierAccount
                        }
                    }),
                    this.searchObjectAttribute({
                        params: {
                            bk_obj_id: 'biz',
                            bk_supplier_account: this.supplierAccount
                        }
                    })
                ])
                this.object['host']['properties'] = res[0].filter(property => !property['bk_isapi'])
                this.object['set']['properties'] = res[1].filter(property => !property['bk_isapi'])
                this.object['module']['properties'] = res[2].filter(property => !property['bk_isapi'])
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
            addUserProperties () {
                let selectedList = []
                for (let key in this.propertySelected) {
                    if (this.propertySelected[key].length) {
                        this.propertySelected[key].map(propertyId => {
                            let property = this.getOriginalProperty(propertyId, key)
                            let {
                                'bk_property_name': propertyName,
                                'bk_property_type': propertyType,
                                'bk_asst_obj_id': asstObjId,
                                'bk_obj_id': objId
                            } = property
                            selectedList.push({
                                propertyId,
                                objId
                            })
                            property.disabled = true
                            let isExist = this.userProperties.findIndex(property => {
                                return propertyId === property.propertyId
                            }) > -1
                            if (!isExist) {
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
                        })
                    }
                }
                this.userProperties = this.userProperties.filter(property => {
                    return selectedList.findIndex(({propertyId, objId}) => {
                        return propertyId === property.propertyId && objId === property.objId
                    }) > -1
                })
                this.toggleUserAPISelector(false)
            },
            toggleContentSelector (isShow) {
                this.$refs.content.open = isShow
                this.attribute.isShow = isShow
            },
            toggleUserAPISelector (isPropertiesShow) {
                if (!isPropertiesShow) {
                    let properties = {
                        host: [],
                        set: [],
                        module: []
                    }
                    this.userProperties.map(property => {
                        properties[property.objId].push(property.bkPropertyId)
                    })
                    this.object.host.selected = properties.host
                    this.object.set.selected = properties.set
                    this.object.module.selected = properties.module
                }
                this.objectInfo.isPropertiesShow = isPropertiesShow
                this.$refs.userapiPop.style.zIndex = ++this.zIndex
            }
        }
    }
</script>


<style lang="scss" scoped>
    .define-wrapper {
        padding: 30px;
        .userapi-group {
            margin-bottom: 15px;
            &.content {
                margin-bottom: 40px;
                .content-selector {
                    position: relative;
                }
            }
            .userapi-label {
                display: block;
                margin-bottom: 5px;
            }
            .userapi-textarea {
                float: left;
                width: 334px;
                height: 80px;
                padding: 5px 16px;
                margin-bottom: 10px;
                font-size: 14px;
                resize: none;
                outline: none;
                vertical-align: bottom;
                color: #666;
                background: #fff;
                border-color: $cmdbBorderColor;
            }
            .btn-icon {
                margin-left: 10px;
                vertical-align: top;
                width: 26px;
                height: 26px;
                padding: 0;
                margin-top: 0;
                font-size: 20px;
                line-height: 25px;
            }
        }
        .userapi-list {
            .filter-label {
                display: block;
            }
            .filter-content {

            }
        }
        .userapi-new{
            width: 334px;
            margin-top: 20px;
            font-size: 14px;
            .userapi-new-btn{
                width: 100%;
                height: 32px;
                background-color: #ffffff;
                border-radius: 2px;
                border: 1px dashed #c3cdd7;
                outline: 0;
                color: #c7ced6;
                &:hover{
                    box-shadow: 0px 3px 6px 0px rgba(51, 60, 72, 0.1);
                }
            }
            .userapi-pop-wrapper {
                position: fixed;
                top: 150px;
                left: 0;
                width: 100%;
                z-index: 99;
            }
            .userapi-new-selector-pop {
                margin: 0 auto;
                padding: 30px;
                background: #fff;
                box-shadow: 0px 3px 6px 0.12px rgba(175, 177, 180, 0.61);
                width: 530px;
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
    }
</style>
