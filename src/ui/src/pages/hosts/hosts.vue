/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <div class="hosts-wrapper clearfix">
        <slot name="filter">
            <div class="filter-container fr">
                <v-filter
                    :queryColumns="filter.queryColumns"
                    :queryColumnData="filter.queryColumnData"
                    :attribute="attribute"
                    :isShowBiz="isShowBiz"
                    :isShowCollect="isShowCollect"
                    :isShowHistory="isShowHistory"
                    :isShowScope="isShowScope"
                    @refresh="setTableCurrentPage(1, true)"
                    @bkBizSelected="bkBizSelected"
                    @showField="setFilterField"
                    @applyCollect="setQueryColumnData"
                    @applyHistory="setQueryColumnData"
                    @filterChange="setFilterParams"
                    @emptyField="emptyFilterField">
                </v-filter>
            </div>
        </slot>
        <div class="table-container">
            <bk-dropdown-menu ref="dropdown" class="fl mr10" :trigger="'click'">
                <bk-button class="dropdown-btn checkbox" type="default" slot="dropdown-trigger">
                    <template v-if="table.chooseId.length">
                        <i class="checkbox-btn" :class="{'checked': table.chooseId.length!==table.pagination.count, 'checked-all': table.chooseId.length===table.pagination.count}" @click.stop="tableChecked('cancel')"></i>
                    </template>
                    <template v-else>
                        <i class="checkbox-btn" @click.stop="tableChecked('current')"></i>
                    </template>
                    <i class="bk-icon icon-angle-down"></i>
                </bk-button>
                <ul class="bk-dropdown-list" slot="dropdown-content">
                    <li>
                        <a href="javascript:;" @click="tableChecked('current')">{{$t("Common['全选本页']")}}</a>
                    </li>
                    <li>
                        <a href="javascript:;" @click="tableChecked('all')">{{$t("Common['跨页全选']")}}</a>
                    </li>
                </ul>
            </bk-dropdown-menu>
            <div class="btn-wrapper clearfix" :class="{'disabled': !table.chooseId.length}">
                <bk-dropdown-menu ref="dropdown" class="mr10" :trigger="'click'">
                    <bk-button class="dropdown-btn" type="default" slot="dropdown-trigger" style="width:100px" :disabled="!table.chooseId.length">
                        <span>{{$t('Common["复制"]')}}</span>
                        <i :class="['bk-icon icon-angle-down',{'icon-flip': isDropdownShow}]"></i>
                    </bk-button>
                    <ul class="bk-dropdown-list" slot="dropdown-content">
                        <template v-for="(item, index) in table.tableHeader">
                            <li v-if="index">
                                <a href="javascript:;" class="copy" :data-clipboard-text="getClipText(item)">{{item.name}}</a>
                            </li>
                        </template>
                    </ul>
                </bk-dropdown-menu>
                <slot name="btnGroup">
                    <div class="bk-group btn-group bk-button-group clearfix">
                        <div class="btn-tooltip-wrapper" v-tooltip="$t('BusinessTopology[\'修改\']')">
                            <button class="bk-button bk-default"
                                :disabled="!table.chooseId.length" 
                                @click="multipleUpdate">
                                <i class="icon-cc-edit"></i>
                            </button>
                        </div>
                        <div class="btn-tooltip-wrapper" v-tooltip="$t('BusinessTopology[\'转移\']')">
                            <bk-button type="default"
                                :disabled="!table.chooseId.length"
                                @click="transferHost">
                                <i class="icon-cc-shift"></i>
                            </bk-button>
                        </div>
                        <div class="btn-tooltip-wrapper" v-tooltip="$t('HostResourcePool[\'导出选中\']')">
                            <form ref="exportForm" :action="exportUrl" method="POST" style="display: inline-block;">
                                <input type="hidden" name="bk_host_id" :value="table.chooseId">
                                <input type="hidden" name="bk_biz_id" value="-1">
                                <bk-button type="default"
                                    class="center"
                                    btnType="submit"
                                    :disabled="!table.chooseId.length"
                                    @click.prevent="exportChoose">
                                    <i class="icon-cc-derivation"></i>
                                </bk-button>
                            </form>
                        </div>
                        <bk-button type="default" v-if="isShowCrossImport" @click="handleCrossImport">{{$t("Common['跨业务导入']")}}</bk-button>
                        <bk-button type="default" class="button-setting last" @click="setTableField" v-tooltip="$t('BusinessTopology[\'列表显示属性配置\']')">
                            <i class="icon-cc-setting"></i>
                        </bk-button>
                        <bk-button type="primary" :loading="$loading('hostSearch')" v-show="isShowRefresh" @click="setTableCurrentPage(1, true)" class="fr mr0">
                            {{$t("HostResourcePool['刷新查询']")}}
                        </bk-button>
                    </div>
                </slot>
            </div>
            <v-table class="index-table"
                ref="indexTable"
                :header="table.tableHeader"
                :list="table.tableList"
                :defaultSort="table.defaultSort"
                :pagination.sync="table.pagination"
                :loading="table.isLoading || outerLoading"
                :checked="table.chooseId"
                :wrapperMinusHeight="wrapperMinusHeight"
                :visible="tableVisible"
                :isCheckboxShow="false"
                @handlePageChange="setTableCurrentPage"
                @handleSizeChange="setTablePageSize"
                @handleSortChange="setTableSort"
                @handleCheckAll="getAllHostID"
                @handleRowClick="showHostAttribute">
                    <template v-for="({id,name, property}, index) in table.tableHeader" :slot="id" slot-scope="{ item }">
                        <label v-if="id === 'bk_host_id'" style="width:100%;text-align:center;" class="bk-form-checkbox bk-checkbox-small" @click.stop>
                            <input type="checkbox"
                                :value="item['host']['bk_host_id']" 
                                v-model="table.chooseId">
                        </label>
                        <template v-else>{{getCellValue(property, item)}}</template>
                    </template>
                </v-table>
        </div>
        <v-sideslider 
            :isShow.sync="sideslider.isShow" 
            :title="sideslider.title"
            :hasCloseConfirm="true"
            :isCloseConfirmShow="sideslider.isCloseConfirmShow"
            :width="sideslider.width"
            @closeSlider="closeSliderConfirm">
            <div slot="content" class="sideslider-content" :class="`sideslider-content-${sideslider.type}`">
                <bk-tab class="attribute-tab" style="border:none;"
                    v-show="sideslider.type === 'attribute'"
                    :active-name="sideslider.attribute.active" 
                    @tab-changed="attributeTabChanged">
                    <bk-tabpanel name="attribute" :title="$t('HostResourcePool[\'主机属性\']')">
                        <v-attribute ref="hostAttribute"
                            :objId="'host'"
                            :showDelete="false"
                            :formValues="sideslider.attribute.form.formValues"
                            :formFields="sideslider.attribute.form.formFields"
                            :type="sideslider.attribute.form.type"
                            :isMultipleUpdate="sideslider.attribute.form.isMultipleUpdate"
                            :active="sideslider.isShow && sideslider.attribute.active === 'attribute'"
                            @submit="saveHostAttribute">
                            <div slot="list" class="attribute-group relation-list">
                                <h3 class="title">{{$t("BusinessTopology['业务拓扑']")}}</h3>
                                <ul class="attribute-list clearfix">
                                    <li class="attribute-item" v-for="item in sideslider.hostRelation">
                                        <span class="attribute-item-value">{{item}}</span>
                                    </li>
                                </ul>
                            </div>
                        </v-attribute>
                    </bk-tabpanel>
                    <bk-tabpanel name="relevance" :title="$t('HostResourcePool[\'关联\']')" :show="!sideslider.attribute.form.isMultipleUpdate">
                        <v-relevance
                            :isShow="sideslider.attribute.active==='relevance'"
                            :objId="'host'"
                            :ObjectID="sideslider.attribute.form.formValues['bk_host_id']"
                            :instance="sideslider.attribute.form.formValues"
                            @handleUpdate="getTableList"
                        ></v-relevance>
                    </bk-tabpanel>
                    <bk-tabpanel name="status" :title="$t('HostResourcePool[\'实时状态\']')"
                        :show="!sideslider.attribute.form.isMultipleUpdate">
                        <v-status :isShow="sideslider.attribute.active==='status'" 
                            :isSidesliderShow="sideslider.isShow"
                            :isWindowsOSType="sideslider.attribute.isWindowsOSType"
                            :isLoaded.sync="sideslider.attribute.status.isLoaded">
                        </v-status>
                    </bk-tabpanel>
                    <bk-tabpanel name="host" title="Host" 
                        :show="!sideslider.attribute.form.isMultipleUpdate && !sideslider.attribute.isWindowsOSType && hostSnapshot !== ''">
                        <v-host></v-host>
                    </bk-tabpanel>
                    <bk-tabpanel name="router" title="Router" 
                        :show="!sideslider.attribute.form.isMultipleUpdate && !sideslider.attribute.isWindowsOSType && hostSnapshot !== ''">
                        <v-router></v-router>
                    </bk-tabpanel>
                    <bk-tabpanel name="history" :title="$t('HostResourcePool[\'变更记录\']')"
                        :show="!sideslider.attribute.form.isMultipleUpdate && sideslider.attribute.form.type === 'update'">
                        <v-history 
                            :type="'host'" 
                            :active="sideslider.attribute.active === 'history'" 
                            :innerIP="sideslider.attribute.form.formValues.bk_host_innerip">
                        </v-history>
                    </bk-tabpanel>
                </bk-tab>
                <v-field v-show="sideslider.type === 'field' && sideslider.isShow"
                    :isShow="sideslider.type === 'field' && sideslider.isShow"
                    :shownFields="sideslider.fields.shownFields"
                    :fieldOptions="sideslider.fields.fieldOptions"
                    :isShowExclude="sideslider.fields.isShowExclude"
                    :minField="sideslider.fields.minField"
                    @apply="applyField"
                    @cancel="cancelSetField">
                </v-field>
            </div>
        </v-sideslider>
        <v-host-transfer-pop v-if="isShowTransfer"
            :isShow.sync="transfer.isShow"
            :chooseId="table.chooseId"
            :hosts="selectedList"
            @success="transferSuccess">
        </v-host-transfer-pop>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import vTable from '@/components/table/table'
    import vFilter from '@/components/filter/filter'
    import vSideslider from '@/components/slider/sideslider'
    import vAttribute from '@/components/object/attribute'
    import vRelevance from '@/components/relevance/relevance'
    import vHostTransferPop from '@/components/hostTransferPop/hostTransferPop'
    import vHistory from '@/components/history/history'
    import vField from '@/components/field/field'
    import vStatus from './children/status.vue'
    import vHost from './children/host'
    import vRouter from './children/router'
    import bus from '@/eventbus/bus'
    import Clipboard from 'clipboard'
    import { getHostRelation } from '@/utils/util'
    export default {
        props: {
            outerParams: {
                type: Object,
                default () {
                    return {}
                }
            },
            isShowCrossImport: {
                type: Boolean,
                default: false
            },
            isShowBiz: {
                type: Boolean,
                default: true
            },
            isShowCollect: {
                type: Boolean,
                default: true
            },
            isShowHistory: {
                type: Boolean,
                default: true
            },
            isShowTransfer: {
                type: Boolean,
                default: true
            },
            isShowRefresh: {
                type: Boolean,
                default: false
            },
            isShowScope: {
                type: Boolean,
                default: false
            },
            outerLoading: {
                type: Boolean,
                default: false
            },
            tableVisible: {
                type: Boolean,
                default: true
            },
            wrapperMinusHeight: {
                type: Number,
                default: 150
            }
        },
        data () {
            return {
                isDropdownShow: false,
                selectedList: [],
                forSelectedList: [],
                bkBizId: '',
                table: {
                    tableHeader: [],
                    useDefaultHeader: true,
                    tableList: [],
                    defaultSort: 'bk_host_id',
                    sort: 'bk_host_id',
                    pagination: {
                        size: 10,
                        current: 1,
                        count: 0
                    },
                    chooseId: [],
                    isLoading: false
                },
                filter: {
                    queryColumns: [],
                    queryColumnData: {}
                },
                attribute: [],
                historyParams: {
                    'bk_content': ''
                },
                sideslider: {
                    type: 'attribute',
                    width: 800,
                    isShow: false,
                    isCloseConfirmShow: false,
                    title: {
                        text: this.$t('HostResourcePool[\'主机属性\']')
                    },
                    attribute: {
                        active: 'attribute',
                        form: {
                            formValues: {},
                            formFields: [],
                            type: 'update',
                            isMultipleUpdate: false
                        },
                        status: {
                            isLoaded: false
                        },
                        isWindowsOSType: true
                    },
                    fields: {
                        shownFields: [],
                        fieldOptions: [],
                        type: 'displayColumns',
                        isShowExclude: true,
                        minField: 1
                    },
                    hostRelation: []
                },
                transfer: {
                    isShow: false
                },
                filterParams: {
                    ip: {
                        data: [],
                        exact: 0,
                        flag: 'bk_host_innerip|bk_host_outerip'
                    },
                    condition: []
                }
            }
        },
        computed: {
            ...mapGetters({
                'bkSupplierAccount': 'bkSupplierAccount',
                'hostSnapshot': 'getHostSnapshot'
            }),
            ...mapGetters('object', ['topo']),
            allProperties () {
                let allProperties = []
                this.attribute.map(({properties}) => {
                    allProperties = [...allProperties, ...properties]
                })
                return allProperties
            },
            bkObjIds () {
                return this.attribute.map(({bk_obj_id: bkObjId}) => {
                    return bkObjId
                })
            },
            attrLoaded () {
                let isLoaded = true
                this.attribute.map(({loaded}) => {
                    if (!loaded) {
                        isLoaded = false
                    }
                })
                return isLoaded
            },
            exportUrl () {
                return `${window.siteUrl}hosts/export`
            },
            searchParams () {
                let params = {
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.size,
                        limit: this.table.pagination.size,
                        sort: this.table.sort
                    },
                    pattern: ''
                }
                return Object.assign(params, this.mergeCondition(this.filterParams, this.outerParams))
            },
            defaultTableHeader () {
                let tableHeader = []
                this.attribute.map(({bk_obj_id: bkObjId, properties}) => {
                    if (bkObjId === 'host') {
                        let requiredProperties = []
                        let notRequiredProperties = []
                        properties.map(property => {
                            if (property['isrequired']) {
                                requiredProperties.push(property)
                            } else {
                                notRequiredProperties.push(property)
                            }
                        })
                        tableHeader = requiredProperties.concat(notRequiredProperties)
                    }
                })
                return tableHeader.slice(0, 7)
            }
        },
        watch: {
            historyParams (historyParams) {
                this.saveHistorySearch()
            },
            defaultTableHeader (defaultTableHeader) {
                if (this.table.useDefaultHeader) {
                    this.setTableHeader(defaultTableHeader)
                }
            },
            outerParams () {
                this.setTableCurrentPage(1)
            },
            'table.chooseId' (chooseId, oldVal) {
                this.setSelectedList(chooseId, oldVal)
                this.$emit('choose', chooseId)
            },
            attrLoaded (attrLoaded) {
                if (attrLoaded && this.bkBizId) {
                    this.setTableCurrentPage(1)
                }
                if (attrLoaded) {
                    this.$emit('attrLoaded')
                }
            },
            bkBizId (bkBizId) {
                if (this.attrLoaded) {
                    this.$nextTick(() => {
                        this.setTableCurrentPage(1)
                    })
                }
                this.$emit('bizLoaded')
            }
        },
        methods: {
            tableChecked (type) {
                if (type === 'all') {
                    this.getAllHostID(true)
                } else if (type === 'cancel') {
                    this.table.chooseId = []
                } else {
                    let chooseId = []
                    this.table.tableList.map(item => chooseId.push(item.host['bk_host_id']))
                    this.table.chooseId = chooseId
                }
            },
            getHostRelation (data) {
                return getHostRelation(data)
            },
            getClipText (item) {
                let text = []
                this.selectedList.map(selected => {
                    let value = this.getCellValue(item.property, selected)
                    if (value) {
                        text.push(value)
                    }
                })
                return text.join('\n')
            },
            setSelectedList (newId, oldId) {
                let diffIdList = newId.concat(oldId).filter(id => !newId.includes(id) || !oldId.includes(id))
                if (newId.length > oldId.length) {
                    let list = this.forSelectedList.filter(li => {
                        return diffIdList.findIndex(id => {
                            return id === li.host['bk_host_id']
                        }) !== -1
                    })
                    this.selectedList = this.selectedList.concat(list)
                } else {
                    let list = this.selectedList.filter(li => {
                        return diffIdList.findIndex(id => {
                            return id === li.host['bk_host_id']
                        }) === -1
                    })
                    this.selectedList = list
                }
            },
            closeSliderConfirm () {
                this.sideslider.isCloseConfirmShow = this.$refs.hostAttribute.isCloseConfirmShow()
            },
            clearChooseId () {
                this.table.chooseId = []
            },
            async setTopoAttribute () {
                await this.$store.dispatch('object/getTopo', true)
                this.attribute = this.topo.filter(model => ['biz', 'set', 'module', 'host'].includes(model['bk_obj_id'])).reverse().map(model => {
                    return {
                        'bk_obj_id': model['bk_obj_id'],
                        'bk_obj_name': model['bk_obj_name'],
                        'properties': [],
                        'loaded': false
                    }
                })
            },
            async getAllAttribute () {
                return this.$Axios.all(this.bkObjIds.map((bkObjId, index) => {
                    return this.getAttribute(bkObjId, index)
                }))
            },
            async getAttribute (bkObjId, index) {
                let params = {
                    'bk_supplier_account': this.bkSupplierAccount,
                    'bk_obj_id': bkObjId,
                    page: {
                        sort: 'bk_property_name'
                    }
                }
                return this.$axios.post('object/attr/search', params).then(res => {
                    if (res.result) {
                        this.attribute[index] = Object.assign(this.attribute[index], {properties: res.data, loaded: true})
                        if (bkObjId === 'host') {
                            this.sideslider.attribute.form.formFields = res.data.slice(0)
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            async getUserCustomColumn () {
                const customPrefix = this.$route.path === '/hosts' ? 'host' : 'resource'
                return this.$axios.post('usercustom/user/search', {}).then(res => {
                    if (res.result) {
                        let hostDisplayColumns = res.data['host_display_column'] || []
                        let hostQueryColumns = (res.data[`${customPrefix}_query_column`] || []).filter(({bk_obj_id: bkObjId}) => !['biz'].includes(bkObjId))
                        let availableDisplayColumn = hostDisplayColumns.filter(column => {
                            return this.getColumnProperty(column['bk_property_id'], column['bk_obj_id'])
                        })
                        if (availableDisplayColumn.length) {
                            this.table.useDefaultHeader = false
                            this.setTableHeader(availableDisplayColumn)
                        } else {
                            this.table.useDefaultHeader = true
                        }
                        let availableQueryColumn = hostQueryColumns.filter(column => {
                            return this.getColumnProperty(column['bk_property_id'], column['bk_obj_id'])
                        })
                        this.filter.queryColumns = availableQueryColumn
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    return res
                }).catch((e) => {
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t('Common[\'您没有当前业务的权限\']'))
                    }
                })
            },
            setTableHeader (columns) {
                this.table.tableHeader = [{
                    id: 'bk_host_id',
                    name: 'bk_host_id',
                    type: 'checkbox',
                    width: 50
                }].concat(columns.map(column => {
                    const property = this.getColumnProperty(column['bk_property_id'], column['bk_obj_id'])
                    return {
                        id: column['bk_property_id'],
                        name: property ? property['bk_property_name'] : column['bk_property_name'],
                        property: property
                    }
                }))
            },
            getColumnProperty (bkPropertyId, bkObjId) {
                return this.allProperties.find(property => {
                    return property['bk_property_id'] === bkPropertyId && property['bk_obj_id'] === bkObjId
                })
            },
            getCellValue (property, item) {
                if (property) {
                    let bkObjId = property['bk_obj_id']
                    let value = item[bkObjId][property['bk_property_id']]
                    if (property['bk_property_id'] === 'bk_module_name') {
                        let moduleName = []
                        item.module.map(({bk_module_name: bkModuleName}) => {
                            moduleName.push(bkModuleName)
                        })
                        return moduleName.join(',')
                    }
                    if (property['bk_property_id'] === 'bk_set_name') {
                        let setName = []
                        item.set.map(({bk_set_name: bksetName}) => {
                            setName.push(bksetName)
                        })
                        return setName.join(',')
                    }
                    if (property['bk_asst_obj_id'] && Array.isArray(value)) {
                        let tempValue = []
                        value.map(({bk_inst_name: bkInstName}) => {
                            if (bkInstName) {
                                tempValue.push(bkInstName)
                            }
                        })
                        value = tempValue.join(',')
                    } else if (property['bk_property_type'] === 'date') {
                        value = this.$formatTime(value, 'YYYY-MM-DD')
                    } else if (property['bk_property_type'] === 'time') {
                        value = this.$formatTime(value)
                    } else if (property['bk_property_type'] === 'enum') {
                        let option = property.option.find(({id}) => {
                            return id === value
                        })
                        if (option) {
                            value = option.name
                        } else {
                            value = ''
                        }
                    }
                    return value
                }
                return ''
            },
            bkBizSelected (bkBizId) {
                this.clearChooseId()
                this.bkBizId = bkBizId
            },
            setFilterParams (filter) {
                this.filterParams = filter
            },
            multipleUpdate () {
                let attribute = this.sideslider.attribute
                this.sideslider.width = 800
                this.sideslider.isShow = true
                this.sideslider.type = 'attribute'
                this.sideslider.title.text = this.$t('HostResourcePool[\'主机属性\']')
                attribute.active = 'attribute'
                attribute.form.isMultipleUpdate = true
                attribute.form.formValues = {bk_host_id: this.table.chooseId.join(',')}
                attribute.form.type = 'create'
            },
            transferHost () {
                this.transfer.isShow = true
            },
            transferSuccess () {
                this.table.chooseId = []
                this.setTableCurrentPage(1)
            },
            exportChoose () {
                this.$refs.exportForm.submit()
            },
            setTableField () {
                let extraProperty = [{
                    bk_property_name: this.$t('Hosts[\'集群名\']'),
                    bk_property_id: 'bk_set_name',
                    bk_obj_id: 'set',
                    bk_isapi: false
                }, {
                    bk_property_name: this.$t('Hosts[\'模块名\']'),
                    bk_property_id: 'bk_module_name',
                    bk_obj_id: 'module',
                    bk_isapi: false
                }]
                const hostAttribute = this.attribute.find(({bk_obj_id: bkObjId}) => bkObjId === 'host') || {}
                const hostProperties = hostAttribute.properties || []
                let properties = [...hostProperties, ...extraProperty].sort((propertyA, propertyB) => {
                    return propertyA['bk_property_name'].localeCompare(propertyB['bk_property_name'])
                })
                this.sideslider.width = 600
                this.sideslider.isShow = true
                this.sideslider.type = 'field'
                this.sideslider.title.text = this.$t('BusinessTopology[\'列表显示属性配置\']')
                this.sideslider.fields.type = 'displayColumns'
                this.sideslider.fields.isShowExclude = true
                this.sideslider.fields.fieldOptions = [{
                    'bk_obj_id': 'host',
                    'bk_obj_name': this.$t('Hosts[\'主机\']'),
                    'properties': properties,
                    'loaded': true
                }]
                this.sideslider.fields.minField = 1
                this.sideslider.fields.shownFields = this.table.tableHeader.slice(1).map(({property}) => {
                    return {
                        bk_property_id: property['bk_property_id'],
                        bk_property_name: property['bk_property_name'],
                        bk_obj_id: property['bk_obj_id']
                    }
                })
            },
            setFilterField () {
                this.sideslider.width = 600
                this.sideslider.isShow = true
                this.sideslider.type = 'field'
                this.sideslider.title.text = this.$t('HostResourcePool[\'主机筛选项设置\']')
                this.sideslider.fields.type = 'queryColumns'
                this.sideslider.fields.isShowExclude = false
                this.sideslider.fields.minField = 0
                this.sideslider.fields.shownFields = this.filter.queryColumns.slice(0)
                let fieldOptions = []
                if (this.$route.path === '/hosts') {
                    fieldOptions = this.attribute.filter(({bk_obj_id: bkObjId}) => !['biz'].includes(bkObjId))
                } else if (this.$route.path === '/resource') {
                    fieldOptions = this.attribute.filter(({bk_obj_id: bkObjId}) => bkObjId === 'host')
                }
                this.sideslider.fields.fieldOptions = fieldOptions
            },
            applyField (fields) {
                if (this.sideslider.fields.type === 'displayColumns') {
                    this.setTableHeader(fields)
                    this.updateUserCustomDisplayColumn(fields)
                } else {
                    this.updateUserCustomQueryColumn(fields)
                }
            },
            emptyFilterField () {
                this.$nextTick(() => {
                    this.setTableCurrentPage(1)
                })
            },
            setQueryColumnData (collect) {
                let queryColumnData = {
                    condition: [{
                        'bk_obj_id': 'host',
                        fields: [],
                        condition: []
                    }, {
                        'bk_obj_id': 'biz',
                        fields: [],
                        condition: []
                    }, {
                        'bk_obj_id': 'module',
                        fields: [],
                        condition: []
                    }, {
                        'bk_obj_id': 'set',
                        fields: [],
                        condition: []
                    }]
                }
                let queryColumns = this.filter.queryColumns.slice(0)
                let info = JSON.parse(collect['info'])
                let queryParams = JSON.parse(collect['query_params'])
                queryColumnData['bk_biz_id'] = info['bk_biz_id']
                queryColumnData['ip'] = {
                    data: info.ip_list,
                    exact: info.exact_search,
                    'bk_host_innerip': info['bk_host_innerip'],
                    'bk_host_outerip': info['bk_host_outerip']
                }
                queryParams.map((params) => {
                    queryColumnData.condition.map(({bk_obj_id: bkObjId, condition}) => {
                        if (params['bk_obj_id'] === bkObjId) {
                            let isInclude = queryColumns.find(({bk_property_id: bkPropertyId, bk_obj_id: columnObjId}) => {
                                return bkPropertyId === params.field && columnObjId === params['bk_obj_id']
                            })
                            if (!isInclude) {
                                let collectQueryColumnProperty = this.getColumnProperty(params.field, params['bk_obj_id'])
                                if (collectQueryColumnProperty) {
                                    let collectQueryColumn = {
                                        'bk_property_id': collectQueryColumnProperty['bk_property_id'],
                                        'bk_property_name': collectQueryColumnProperty['bk_property_name'],
                                        'bk_property_type': collectQueryColumnProperty['bk_property_type']
                                    }
                                    if (collectQueryColumnProperty['option']) {
                                        collectQueryColumn['option'] = collectQueryColumnProperty['option']
                                    }
                                    queryColumns.push(collectQueryColumn)
                                }
                            }
                            condition.push({
                                field: params.field,
                                operator: params.operator,
                                value: params.value
                            })
                        }
                    })
                })
                this.filter.queryColumns = queryColumns
                this.filter.queryColumnData = queryColumnData
                this.$nextTick(() => {
                    this.setTableCurrentPage(1)
                })
            },
            updateUserCustomDisplayColumn (fields) {
                const customPrefix = this.$route.path === '/hosts' ? 'host' : 'resource'
                let updateParams = {}
                updateParams['host_display_column'] = fields.map(({bk_property_id, bk_property_name, bk_obj_id}) => {
                    return {bk_property_id, bk_property_name, bk_obj_id}
                })
                this.$axios.post('usercustom', JSON.stringify(updateParams), {id: 'userCustom'}).then(res => {
                    this.cancelSetField()
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            updateUserCustomQueryColumn (fields) {
                let columns = fields.map(property => {
                    let {
                        bk_property_id: bkPropertyId,
                        bk_property_name: bkPropertyName,
                        bk_property_type: bkPropertyType,
                        bk_obj_id: bkObjId,
                        option
                    } = property
                    let column = {
                        bk_property_id: bkPropertyId,
                        bk_property_name: bkPropertyName,
                        bk_property_type: bkPropertyType,
                        bk_obj_id: bkObjId
                    }
                    if (option) {
                        column['option'] = option
                    }
                    return column
                })
                this.filter.queryColumns = columns
                const customPrefix = this.$route.path === '/hosts' ? 'host' : 'resource'
                let updateParams = {}
                updateParams[`${customPrefix}_query_column`] = columns
                this.$axios.post('usercustom', JSON.stringify(updateParams), {id: 'userCustom'}).then(res => {
                    this.cancelSetField()
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            cancelSetField () {
                this.sideslider.isShow = false
            },
            getTableList () {
                this.table.isLoading = true
                this.$axios.post('hosts/search', this.searchParams, {id: 'hostSearch'}).then(res => {
                    this.table.isLoading = false
                    if (res.result) {
                        this.table.pagination.count = res.data.count
                        this.table.tableList = res.data.info
                        this.forSelectedList = this.$deepClone(res.data.info)
                        this.historyParams = {
                            content: JSON.stringify(this.searchParams)
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                }).catch((e) => {
                    this.table.isLoading = false
                    this.table.tableList = []
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t('Common[\'您没有当前业务的权限\']'))
                    }
                })
            },
            saveHistorySearch () {
                this.$axios.post('hosts/history', this.historyParams).then(res => {
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            showHostAttribute (item, index) {
                this.sideslider.hostRelation = getHostRelation(item)
                let bkHostId = item['host']['bk_host_id']
                let attribute = this.sideslider.attribute
                this.sideslider.width = 800
                this.sideslider.isShow = true
                this.sideslider.type = 'attribute'
                this.sideslider.title.text = this.$t('HostResourcePool[\'主机属性\']')
                attribute.active = 'attribute'
                attribute.form.formValues = {}
                attribute.form.isMultipleUpdate = false
                attribute.form.type = 'update'
                this.getHostDetails(bkHostId)
            },
            getHostDetails (bkHostId) {
                this.$axios.get(`hosts/${this.bkSupplierAccount}/${bkHostId}`).then((res) => {
                    if (res.result) {
                        let values = {
                            bk_host_id: bkHostId
                        }
                        res.data.map(({bk_property_id: bkPropertyId, bk_property_value: bkPropertyValue}) => {
                            values[bkPropertyId] = bkPropertyValue !== null ? bkPropertyValue : ''
                            if (bkPropertyId === 'OSType') {
                                this.sideslider.attribute.isWindowsOSType = bkPropertyValue !== 'Linux'
                            }
                        })
                        this.sideslider.attribute.form.formValues = values
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
                this.$axios.get(`hosts/snapshot/${bkHostId}`).then(res => {
                    if (res.result) {
                        this.sideslider.attribute.status.isLoaded = true
                        this.$store.commit('setHostSnapshot', res.data)
                    }
                })
            },
            saveHostAttribute (formData, formValues) {
                let { bk_host_id: bkHostID } = formValues
                this.$axios.put('hosts/batch', Object.assign(formData, {bk_host_id: bkHostID.toString()}), {id: 'editAttr'}).then(res => {
                    if (res.result) {
                        this.$alertMsg(this.$t('Common[\'保存成功\']'), 'success')
                        this.setTableCurrentPage(1)
                        if (!this.sideslider.attribute.form.isMultipleUpdate) {
                            this.$refs.hostAttribute.displayType = 'list'
                            this.getHostDetails(bkHostID)
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                }).catch(e => {
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t('Common[\'权限不足\']'))
                    }
                })
            },
            getAllHostID (isChecked) {
                if (isChecked) {
                    let allHostId = []
                    let searchParams = JSON.parse(this.historyParams['content'])
                    searchParams.page = {}
                    // searchParams.condition.map(({bk_obj_id: bkObjId, fields}) => {
                    //     if (bkObjId === 'host') {
                    //         fields.push('bk_host_id')
                    //     }
                    // })
                    this.table.isLoading = true
                    this.$axios.post('hosts/search/', searchParams).then(res => {
                        if (res.result) {
                            res.data.info.forEach((item, index) => {
                                allHostId.push(item['host']['bk_host_id'])
                            })
                            this.forSelectedList = res.data.info
                            this.table.chooseId = allHostId
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                        this.table.isLoading = false
                    }).catch(() => {
                        this.table.isLoading = false
                    })
                } else {
                    this.table.chooseId = []
                }
            },
            handleCrossImport () {
                this.$emit('handleCrossImport')
            },
            attributeTabChanged (activeName) {
                this.sideslider.attribute.active = activeName
            },
            setTableCurrentPage (current, reset = false) {
                this.table.pagination.current = current
                if (reset) {
                    this.table.chooseId = []
                }
                this.getTableList()
            },
            setTablePageSize (size) {
                this.table.pagination.size = size
                this.setTableCurrentPage(1)
            },
            setTableSort (sort) {
                this.table.sort = sort
                this.setTableCurrentPage(1)
            },
            mergeCondition (targetParams, sourceParams) {
                let mergedParams = this.$deepClone(targetParams)
                if (sourceParams && sourceParams.hasOwnProperty('condition')) {
                    sourceParams.condition.forEach(sourceCondition => {
                        let targetCondition = mergedParams.condition.find(targetCondition => targetCondition['bk_obj_id'] === sourceCondition['bk_obj_id'])
                        if (targetCondition) {
                            targetCondition.fields = [...new Set([...targetCondition.fields, ...sourceCondition.fields])]
                            sourceCondition.condition.forEach(sourceMeta => {
                                let targetMeta = targetCondition.condition.find(targetMeta => targetMeta.field === sourceMeta.field)
                                if (targetMeta) {
                                    targetMeta.operator = sourceMeta.operator
                                    targetMeta.value = sourceMeta.value
                                } else {
                                    targetCondition.condition.push(this.$deepClone(sourceMeta))
                                }
                            })
                        } else {
                            mergedParams.condition.push(sourceCondition)
                        }
                    })
                }
                if (sourceParams && sourceParams.hasOwnProperty('bk_biz_id')) {
                    mergedParams['bk_biz_id'] = sourceParams['bk_biz_id']
                }
                return mergedParams
            },
            async init () {
                await this.setTopoAttribute()
                await this.getAllAttribute()
                await this.getUserCustomColumn()
            }
        },
        created () {
            this.init()
            this.$nextTick(() => {
                let clipboard = new Clipboard('.copy')
                clipboard.on('success', () => {
                    this.$alertMsg(this.$t('Common["复制成功"]'), 'success')
                })
                clipboard.on('error', () => {
                    this.$alertMsg(this.$t('Common["复制失败"]'))
                })
            })
        },
        beforeRouteLeave (to, from, next) {
            this.$store.commit('resetHostSearch')
            next()
        },
        components: {
            vTable,
            vFilter,
            vSideslider,
            vAttribute,
            vRelevance,
            vStatus,
            vHost,
            vRouter,
            vHostTransferPop,
            vHistory,
            vField
        }
    }
</script>

<style lang="scss" scoped>
.hosts-wrapper{
    height: 100%;
}
.table-container{
    padding: 20px 20px 0;
    height: 100%;
    overflow: hidden;
    .dropdown-btn{
        width: 100px;
        cursor: pointer;
        &.checkbox {
            width: 75px;
        }
    }
    .btn-group{
        display: inline-block;
        width: calc(100% - 196px);
        vertical-align: middle;
        font-size: 0;
        .bk-button{
            font-size: 14px;
            &:disabled{
                cursor: not-allowed !important;
            }
            &.button-search{
                width: 178px;
                margin: 0;
                background-color: #3c96ff;
                border-color: #3c96ff;
                color:#fff;
            }
        }
    }
    .index-table{
        margin-top: 10px;
    }
}
.filter-container{
    height: 100%;
    overflow: visible;
}
.sideslider-content{
    padding: 0 20px;
    &.sideslider-content-attribute{
        height: calc(100% - 132px);
    }
    &.sideslider-content-field{
        height: 100%;
        padding: 0;
    }
}
.attribute-tab{
    height: 100%;
}
.attribute-group.relation-list {
    .attribute-item {
        width: 100%;
        .attribute-item-value {
            max-width: 100%;
        }
    }
}
</style>