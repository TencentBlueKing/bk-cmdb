/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template lang="html">
   <div class="host-resource-wrapper clearfix">
        <div class="bottom-contain clearfix">
            <div class="btn-group fl">
                <template v-if="objId!=='biz'">
                    <bk-dropdown-menu ref="dropdown" class="mr10" :trigger="'click'">
                        <bk-button class="dropdown-btn checkbox" type="default" slot="dropdown-trigger">
                            <template v-if="table.chooseId.length">
                                <i class="checkbox-btn" :class="{'checked': table.chooseId.length!==table.pagination.count, 'checked-all': table.chooseId.length===table.pagination.count}" @click.stop="tableChecked('cancel')"></i>
                            </template>
                            <template v-else>
                                <i class="checkbox-btn" @click.stop="tableChecked('current')"></i>
                            </template>
                            <i :class="['bk-icon icon-angle-down',{'icon-flip': isDropdownShow}]"></i>
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
                    <div class="bk-group bk-button-group mr10">
                        <bk-button v-tooltip="$t('ModelManagement[\'导入\']')" type="default" class="bk-button vice-btn" @click="importSlider.isShow = true" :disabled="unauthorized.update">
                            <i class="icon-cc-import"></i>
                        </bk-button>
                        <div class="btn-tooltip-wrapper" v-tooltip="$t('ModelManagement[\'导出\']')">
                            <form :action="exportUrl" ref="export" style="display: inline-block;" method="POST">
                                <input type="hidden" :value="table.chooseId.join(',')" name="bk_inst_id">
                                <bk-button btnType="submit" class="bk-button vice-btn first" :disabled="!table.chooseId.length" title="test">
                                    <i class="icon-cc-derivation"></i>
                                </bk-button>
                            </form>
                        </div>
                        <div class="btn-tooltip-wrapper" v-tooltip="$t('BusinessTopology[\'修改\']')">
                            <bk-button type="default" class="vice-btn"
                                :disabled="!table.chooseId.length" 
                                @click="multipleUpdate">
                                <i class="icon-cc-edit"></i>
                            </bk-button>
                        </div>
                        <div class="btn-tooltip-wrapper" v-tooltip="$t('Common[\'删除\']')">
                            <bk-button type="default" v-if="objId !== 'biz'" class="bk-button delete-button mr10" :disabled="!table.chooseId.length" @click="confirmBatchDel">
                                <i class="icon-cc-del"></i>
                            </bk-button>
                        </div>
                    </div>
                </template>
                <button class="bk-button bk-primary bk-button-componey create-btn mr10" @click="openObjectSlider('create')" :disabled="unauthorized.update">{{$t("Inst['立即创建']")}}</button>
            </div>
            <div class="fr bk-group bk-button-group">
                <bk-button v-if="objId !== 'biz'" class="bk-button ml10" @click="filing.isShow = true" v-tooltip="$t('Common[\'查看删除历史\']')">
                    <i class="icon-cc-history"></i>
                </bk-button>
                <bk-button v-else class="bk-button ml10" @click="filing.isShow = true" v-tooltip="$t('Common[\'查看归档历史\']')">
                    <i class="icon-cc-history2"></i>
                </bk-button>
                <bk-button class="bk-button " @click="settingSlider.isShow = true" v-tooltip="$t('BusinessTopology[\'列表显示属性配置\']')">
                    <i class="icon-cc-setting"></i>
                </bk-button>
            </div>
            <div class="quick-search fr">
                <div class="fl left-select">
                    <bk-select :selected.sync="filter.selected" ref="filterSelector" @on-selected="setFilterType">
                        <bk-select-option
                            v-for="(option, index) of filterList"
                            :key="option.id"
                            :value="option.id"
                            :label="option.name">
                        </bk-select-option>
                    </bk-select>
                </div>
                <template v-if="filter.type === 'enum'">
                    <bk-select class="search-options fl" :selected.sync="filter.value" @on-selected="doFilter" :showClear="true">
                        <bk-select-option v-for="option in getEnumOptions()"
                            :key="option.id"
                            :value="option.id"
                            :label="option.name">
                        </bk-select-option>
                    </bk-select>
                </template>
                <template v-else>
                    <input v-if="filter.type === 'int'" type="text" maxlength="11" class="bk-form-input search-text int" 
                    :placeholder="$t('Common[\'快速查询\']')" v-model.number="filter.value" @keyup.enter="doFilter">
                    <input v-else type="text" class="bk-form-input search-text" :placeholder="$t('Common[\'快速查询\']')" v-model.trim="filter.value" @keyup.enter="doFilter">
                    <i class="bk-icon icon-search" @click="doFilter"></i>
                </template>
            </div>
        </div>
        <div class="table-contain">
            <v-object-table
                :header="table.header" 
                :list="table.list" 
                :pagination.sync="table.pagination"
                :defaultSort="table.defaultSort"
                :checked.sync="table.chooseId"
                :wrapperMinusHeight="150"
                :loading="$loading('instSearch')"
                :isCheckboxShow="false"
                @handleRowClick="editObject"
                @handleSortChange="setTableSort"
                @handlePageChange="setTablePage"
                @handleSizeChange="setTableSize"
                @handleCheckAll="getAllObjectId">
                <template v-for="({property, id}, index) in table.header.filter(head => head.type !== 'checkbox')" :slot="id" slot-scope="{item}">
                    <template v-if="!!property['bk_asst_obj_id']">
                        {{getAssociateCell(item[id])}}
                    </template>
                    <template v-else-if="property['bk_property_type'] === 'enum'">
                        {{getEnumCell(item[id], property)}}
                    </template>
                    <template v-else>{{item[id]}}</template>
                </template>
            </v-object-table>
            <v-sideslider
                :isShow.sync="slider.isShow"
                :hasCloseConfirm="true"
                :isCloseConfirmShow="slider.isCloseConfirmShow"
                :title="slider.title"
                @closeSlider="closeObjectSliderConfirm">
                <div class="slide-content" slot="content">
                    <bk-tab :active-name="tab.activeName" style="border: none;" @tab-changed="tabChanged">
                        <bk-tabpanel name="attr" :title="$t('Common[\'属性\']')">
                            <v-object-attr 
                                ref="attribute"
                                :formFields="attr.formFields" 
                                :formValues="attr.formValues"
                                :showDelete="isShowDelete"
                                :type="attr.type"
                                :active="slider.isShow && tab.activeName === 'attr'"
                                :objId="objId"
                                :isBatchUpdate="false"
                                :isMultipleUpdate="attr.isMultipleUpdate"
                                @closeSlider="closeObjectSlider"
                                @submit="saveObjectAttr"
                                @delete="confirmDelete">
                            </v-object-attr>
                        </bk-tabpanel>
                        <bk-tabpanel name="relevance" :title="$t('HostResourcePool[\'关联\']')" :show="attr.type==='update'">
                            <v-relevance :isShow="tab.activeName==='relevance'"
                                :objId="objId"
                                :ObjectID="objId !== 'biz' ? attr.formValues['bk_inst_id'] : attr.formValues['bk_biz_id']"
                                :instance="attr.formValues"
                                @handleUpdate="getTableList">
                            </v-relevance>
                        </bk-tabpanel>
                        <bk-tabpanel name="history" :title="$t('HostResourcePool[\'变更记录\']')" :show="attr.type==='update'">
                            <v-history v-if="attr.type !== 'create'" :active="tab.activeName === 'history'" :type="objId" :instId="objId === 'biz' ? attr.formValues['bk_biz_id'] : attr.formValues['bk_inst_id']"></v-history>
                        </bk-tabpanel>
                    </bk-tab>
                </div>
            </v-sideslider>
        </div>
        <v-sideslider :isShow.sync="importSlider.isShow" :hasQuickClose="true" :title="{icon: 'icon-cc-derivation',text: `${$t('ModelManagement[\'导入\']')} ${objName}`}" @closeSlider="closeImportSlider">
            <v-import v-if="importSlider.isShow" slot="content" 
                :templateUrl="templateUrl" 
                :importUrl="importUrl" 
                @success="setTablePage(1)"
                @partialSuccess="setTablePage(1)"></v-import>
        </v-sideslider>
        <v-sideslider :isShow.sync="settingSlider.isShow" :hasQuickClose="true" :width="600" :title="settingSlider.title">
            <v-config-field 
                ref="configField"
                slot="content"
                :isShow="settingSlider.isShow"
                :attrList="attr.formFields"
                @apply="settingApply"
                @cancel="settingSlider.isShow = false"
                @resetFields="resetFields"
                :objId="objId">
            </v-config-field>
        </v-sideslider>
        <v-delete-history
            :isShow.sync="filing.isShow"
            :objId="objId"
            :objTableHeader="table.header"
        ></v-delete-history>
   </div>
</template>

<script>
    import vAssociationList from '@/components/relevance/list'
    import vNewAssociation from '@/components/relevance/new-association'
    import { mapGetters } from 'vuex'
    import vObjectTable from '@/components/table/table'
    import vObjectAttr from '@/components/object/attribute'
    import Authority from '@/mixins/authority'
    import vRelevance from '@/components/relevance/relevance'
    import vHistory from '@/components/history/history'
    import vImport from '@/components/import/import'
    import vSideslider from '@/components/slider/sideslider'
    import vConfigField from './children/configField'
    import vDeleteHistory from '@/components/history/delete'
    export default {
        mixins: [Authority],
        data () {
            return {
                isDropdownShow: false,
                isSelectShow: false,
                // 查询条件
                filter: {
                    list: [],
                    selected: '',
                    value: '',
                    type: ''
                },
                filterList: [],
                // 表格数据
                table: {
                    loading: false,
                    header: [],
                    list: [],
                    chooseId: [],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-bk_inst_id',
                    sort: ''
                },
                filing: {
                    isShow: false
                },
                // 侧滑状态
                slider: {
                    isShow: false,
                    isCloseConfirmShow: false,
                    title: {
                        icon: '',
                        text: ''
                    }
                },
                importSlider: {
                    isShow: false
                },
                settingSlider: {
                    isShow: false,
                    title: {
                        text: this.$t("Common['编辑']")
                    }
                },
                // 属性展示界面状态
                attr: {
                    type: 'update',
                    formFields: [],
                    formValues: {},
                    isMultipleUpdate: false
                },
                // 选项卡
                tab: {
                    activeName: 'attr'
                },
                prevHistoryCount: 0
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount',
                'usercustom'
            ]),
            
            ...mapGetters('navigation', ['activeClassifications']),
            // 当前路由对应的模型ID
            objId () {
                return this.$route.params.objId
            },
            // 当前路由对应的模型名称
            objName () {
                let objName = ''
                this.activeClassifications.map(classify => {
                    classify['bk_objects'].map(model => {
                        if (model['bk_obj_id'] === this.objId) {
                            objName = model['bk_obj_name']
                        }
                    })
                })
                return objName
            },
            isShowDelete () {
                return !(this.objId === 'biz' && this.attr.formValues['bk_biz_name'] === '蓝鲸')
            },
            // 表格查询需要的参数
            axiosConfig () {
                let config = {
                    url: undefined,
                    params: {
                        page: {
                            start: (this.table.pagination.current - 1) * this.table.pagination.size,
                            limit: this.table.pagination.size,
                            sort: this.table.sort ? this.table.sort : this.table.defaultSort
                        },
                        fields: this.objId === 'biz' ? [] : {},
                        condition: this.objId === 'biz' ? {bk_data_status: {'$ne': 'disabled'}} : {}
                    }
                }
                if (this.objId === 'biz') {
                    config.url = `biz/search/${this.bkSupplierAccount}`
                } else {
                    config.url = `inst/association/search/owner/${this.bkSupplierAccount}/object/${this.objId}`
                }
                if (this.filter.selected && this.filter.value !== '') {
                    if (this.objId === 'biz') {
                        if (this.filter.type === 'bool' && ['true', 'false'].includes(this.filter.value)) {
                            config.params.condition[this.filter.selected] = this.filter.value === 'true'
                        } else {
                            config.params.condition[this.filter.selected] = this.filter.value
                        }
                    } else {
                        const specialObj = {
                            'host': 'bk_host_innerip',
                            'biz': 'bk_biz_name',
                            'plat': 'bk_cloud_name',
                            'module': 'bk_module_name',
                            'set': 'bk_set_name'
                        }
                        if (this.filter.type === 'singleasst' || this.filter.type === 'multiasst') {
                            let bkAsstObjId = this.getProperty(this.filter.selected)['bk_asst_obj_id']
                            config.params.condition[bkAsstObjId] = [{
                                field: specialObj.hasOwnProperty(bkAsstObjId) ? specialObj[bkAsstObjId] : 'bk_inst_name',
                                operator: '$regex',
                                value: this.filter.value
                            }]
                        } else if (this.filter.type === 'bool') {
                            config.params.condition[this.objId] = [{
                                field: this.filter.selected,
                                operator: '$eq',
                                value: ['true', 'false'].includes(this.filter.value) ? this.filter.value === 'true' : this.filter.value
                            }]
                        } else {
                            config.params.condition[this.objId] = [{
                                field: this.filter.selected,
                                operator: ['enum', 'int'].includes(this.filter.type) ? '$eq' : '$regex',
                                value: this.filter.value
                            }]
                        }
                    }
                }
                return config
            },
            templateUrl () {
                return `${window.siteUrl}importtemplate/${this.objId}`
            },
            importUrl () {
                return `${window.siteUrl}insts/owner/${this.bkSupplierAccount}/object/${this.objId}/import`
            },
            exportUrl () {
                return `${window.siteUrl}insts/owner/${this.bkSupplierAccount}/object/${this.objId}/export`
            }
        },
        watch: {
            'filing.isShow' (isShow) {
                if (!isShow && this.objId === 'biz') {
                    this.getTableList()
                }
                this.updateHistoryCount(isShow)
            },
            // 切换模型时，重新初始化表格
            objId () {
                // 页码调整到第一页
                this.table.pagination.current = 1
                this.filter.value = ''
                this.table.chooseId = []
                // 初始化表格
                this.initTable()
            },
            filterList (filterList) {
                if (filterList.length) {
                    this.filter.type = this.getProperty(filterList[0]['id'])['bk_property_type']
                    this.filter.selected = filterList.length ? filterList[0]['id'] : ''
                } else {
                    this.filter.type = ''
                    this.filter.selected = ''
                }
            },
            'slider.isShow' (isShow) {
                if (!isShow) {
                    this.closeObjectSlider()
                }
            },
            'filter.selected' () {
                this.filter.value = ''
            }
        },
        methods: {
            tableChecked (type) {
                if (type === 'all') {
                    this.getAllObjectId(true)
                } else if (type === 'cancel') {
                    this.table.chooseId = []
                } else {
                    let chooseId = []
                    this.table.list.map(({bk_inst_id: bkInstId}) => chooseId.push(bkInstId))
                    this.table.chooseId = chooseId
                }
            },
            async resetFields () {
                await this.getTableHeader()
                this.$refs.configField.getUserAttr()
            },
            multipleUpdate () {
                this.tab.activeName = 'attr'
                this.slider.isShow = true
                this.slider.title.text = this.$t('Inst[\'批量更新\']')
                this.attr.isMultipleUpdate = true
                this.attr.formValues = {bk_inst_id: this.table.chooseId.join(',')}
                this.attr.type = 'create'
            },
            closeObjectSliderConfirm () {
                this.slider.isCloseConfirmShow = this.$refs.attribute.isCloseConfirmShow()
            },
            getProperty (id) {
                return this.attr.formFields.find(({bk_property_id: bkPropertyId}) => bkPropertyId === id)
            },
            setFilterType (option, index) {
                this.filter.type = this.getProperty(option.value)['bk_property_type']
            },
            settingApply (shownList) {
                let list = []
                shownList.map(val => {
                    for (let i = 0, formFields = this.attr.formFields; i < formFields.length; i++) {
                        if (val['bk_property_id'] === formFields[i]['bk_property_id']) {
                            list.push({
                                id: val['bk_property_id'],
                                name: val['bk_property_name'],
                                property: JSON.parse(JSON.stringify(formFields[i]))
                            })
                            break
                        }
                    }
                })
                if (this.objId === 'biz') {
                    list.unshift({
                        id: 'bk_biz_id',
                        name: 'ID',
                        property: {}
                    })
                }
                if (this.table.header[0].hasOwnProperty('type') && this.table.header[0].type === 'checkbox') {
                    this.table.header.splice(1)
                    this.table.header = this.table.header.concat(list)
                } else {
                    this.table.header = JSON.parse(JSON.stringify(list))
                }
                this.settingSlider.isShow = false
            },
            tabChanged (name) {
                this.tab.activeName = name
            },
            async getUsercustom () {
                await this.$axios.post('usercustom/user/search').then(res => {
                    if (res.result) {
                        this.$store.commit('setUsercustom', res.data)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            // 初始化表格，先获取表格字段再获取表格列表
            initTable () {
                this.getTableHeader().then(properties => {
                    this.attr.formFields = [...properties]
                    this.filterList = properties.map(property => {
                        return {
                            id: property['bk_property_id'],
                            name: property['bk_property_name'],
                            type: property['bk_property_type']
                        }
                    })
                    this.getTableList()
                })
            },
            /*
                设置字段
            */
            setUsercustom (data) {
                let list = []
                data.map(val => {
                    list.push({
                        bk_property_id: val.id,
                        bk_property_name: val.name,
                        bk_property_type: val.property['bk_property_type']
                    })
                })
                let usercustom = this.$deepClone(this.usercustom)
                usercustom[`${this.objId}DisplayColumn`] = list
                this.$store.commit('setUsercustom', usercustom)
            },
            // 获取表格属性字段
            getTableHeader () {
                return this.$axios.post('object/attr/search', {bk_obj_id: this.objId, bk_supplier_account: this.bkSupplierAccount}).then(async res => {
                    let header = []
                    let headerLead = []
                    let headerTail = []
                    if (res.result) {
                        await this.getUsercustom()
                        if (this.usercustom.hasOwnProperty(`${this.objId}DisplayColumn`) && this.usercustom[`${this.objId}DisplayColumn`].length) {
                            this.usercustom[`${this.objId}DisplayColumn`].map(val => {
                                for (let i = 0, attr = res.data; i < res.data.length; i++) {
                                    if (val['bk_property_id'] === attr[i]['bk_property_id']) {
                                        header.push({
                                            id: attr[i]['bk_property_id'],
                                            name: attr[i]['bk_property_name'],
                                            property: attr[i]
                                        })
                                        break
                                    }
                                }
                                if (val['bk_property_type'] !== 'singleasst' && val['bk_property_type'] !== 'multiasst') {
                                    let property = res.data.find(({bk_property_id: bkPropertyId}) => {
                                        return bkPropertyId === val['bk_property_id']
                                    })
                                }
                            })
                        } else { // 没有时则显示前六
                            res.data.map(attr => {
                                let headerObj = {
                                    id: attr['bk_property_id'],
                                    name: attr['bk_property_name'],
                                    property: attr
                                }
                                if (attr['bk_isonly'] && attr['bk_isrequired']) {
                                    header.unshift(headerObj)
                                } else if (attr['bk_isonly'] || attr['bk_isrequired']) {
                                    header.push(headerObj)
                                } else {
                                    headerTail.push(headerObj)
                                }
                            })
                        }
                        if (this.objId === 'biz') {
                            headerLead.unshift({
                                id: 'bk_biz_id',
                                name: 'ID',
                                property: {}
                            })
                        } else {
                            headerLead.unshift({
                                id: 'bk_inst_id',
                                name: 'ID',
                                width: 50,
                                type: 'checkbox',
                                property: {}
                            })
                        }
                        header = headerLead.concat(header).concat(headerTail)
                        if (this.usercustom.hasOwnProperty(`${this.objId}DisplayColumn`) && this.usercustom[`${this.objId}DisplayColumn`].length) {
                            this.table.header = header
                        } else {
                            this.setUsercustom(header.slice(1, 7))
                            this.table.header = header.slice(0, 7)
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    return res.data
                })
            },
            // 获取表格列表
            getTableList () {
                this.table.loading = true
                return this.$axios.post(this.axiosConfig.url, this.axiosConfig.params, {id: 'instSearch'}).then(res => {
                    let data = {
                        count: 0,
                        info: []
                    }
                    if (res.result) {
                        data.count = res.data.count || 0
                        data.info = res.data.info || []
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.table.list = data.info
                    this.table.pagination.count = data.count
                    return data
                }).catch(e => {
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t("Inst['您没有当前模型的权限']"))
                    }
                }).finally(() => {
                    this.table.loading = false
                })
            },
            getAllObjectId (isAllCheck) {
                if (isAllCheck) {
                    this.table.loading = true
                    let idKey = this.objId === 'biz' ? 'bk_biz_id' : 'bk_inst_id'
                    let params = this.$deepClone(this.axiosConfig.params)
                    params.page = {}
                    this.objId === 'biz' ? void 0 : params.fields[this.objId] = [idKey]
                    this.$axios.post(this.axiosConfig.url, params).then(res => {
                        if (res.result) {
                            let chooseId = []
                            res.data.info.map(attr => {
                                chooseId.push(attr[idKey])
                            })
                            this.table.chooseId = chooseId
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    }).finally(() => {
                        this.table.loading = false
                    })
                } else {
                    this.table.chooseId = []
                }
            },
            // 表格排序
            setTableSort (sort) {
                this.table.sort = sort
                this.getTableList()
            },
            // 表格展示条数
            setTableSize (size) {
                this.table.pagination.current = 1
                this.table.pagination.size = size
                this.getTableList()
            },
            // 表格当前页
            setTablePage (page) {
                this.table.pagination.current = page
                this.getTableList()
            },
            doFilter () {
                this.table.chooseId = []
                this.setTablePage(1)
            },
            // 保存新增/修改的属性
            saveObjectAttr (formData, {bk_biz_id: bizId, bk_inst_id: instId}) {
                if (this.attr.type === 'create' && this.attr.isMultipleUpdate) {
                    let params = {
                        update: []
                    }
                    this.table.chooseId.map(id => {
                        params.update.push({
                            datas: formData,
                            inst_id: id
                        })
                    })
                    this.$axios.put(`inst/${this.bkSupplierAccount}/${this.objId}/batch`, params).then(res => {
                        if (res.result) {
                            this.setTablePage(1)
                            this.closeObjectSlider()
                            this.$alertMsg(this.$t("Common['修改成功']"), 'success')
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                } else if (this.attr.type === 'update') {
                    let updateUrl = this.objId === 'biz'
                        ? `biz/${this.bkSupplierAccount}/${bizId}`
                        : `inst/${this.bkSupplierAccount}/${this.objId}/${instId}`
                    this.$axios.put(updateUrl, formData, {id: 'editAttr'}).then(res => {
                        if (res.result) {
                            this.setTablePage(1)
                            this.closeObjectSlider()
                            this.$alertMsg(this.$t("Common['修改成功']"), 'success')
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                } else {
                    let createUrl = this.objId === 'biz' ? `biz/${this.bkSupplierAccount}` : `inst/${this.bkSupplierAccount}/${this.objId}`
                    this.$axios.post(createUrl, formData, {id: 'editAttr'}).then(res => {
                        if (res.result) {
                            this.setTablePage(1)
                            this.closeObjectSlider()
                            this.$alertMsg(this.$t("Inst['创建成功']"), 'success')
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            confirmBatchDel () {
                this.$bkInfo({
                    title: this.$tc("Inst['确定删除选中的实例']", this.table.chooseId.length, {number: this.table.chooseId.length}),
                    confirmFn: () => {
                        this.batchDeleteInst()
                    }
                })
            },
            async batchDeleteInst () {
                try {
                    const params = {
                        delete: {
                            inst_ids: this.table.chooseId
                        }
                    }
                    const res = await this.$axios.delete(`inst/${this.bkSupplierAccount}/${this.objId}/batch`, {data: params})
                    if (res.result) {
                        this.table.chooseId = []
                        this.setTablePage(1)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            // 删除模型实例前进行确认
            confirmDelete (data) {
                let {
                    bk_inst_name: instName,
                    bk_biz_name: bizName
                } = data
                this.$bkInfo({
                    title: this.objId === 'biz' ? this.$t("Common['确认要归档']", {name: this.objId === 'biz' ? bizName : instName}) : this.$t("Common['确认要删除']", {name: this.objId === 'biz' ? bizName : instName}),
                    content: this.objId === 'biz' ? this.$t("Inst['归档提示']") : '',
                    confirmFn: () => {
                        this.deleteObject(data)
                    }
                })
            },
            // 删除模型实例
            async deleteObject ({bk_biz_id: bizId, bk_inst_id: instId}) {
                let url = ''
                let method = ''
                if (this.objId === 'biz') {
                    url = `biz/status/disabled/${this.bkSupplierAccount}/${bizId}`
                    method = 'put'
                } else {
                    url = `inst/${this.bkSupplierAccount}/${this.objId}/${instId}`
                    method = 'delete'
                }
                try {
                    const res = await this.$axios({
                        method: method,
                        url: url,
                        id: 'instDelete'
                    })
                    if (res.result) {
                        this.setTablePage(1)
                        this.closeObjectSlider()
                        this.table.chooseId = this.table.chooseId.filter(id => id !== (this.objId === 'biz' ? bizId : instId))
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            // 展示模型实例属性明细
            editObject (item) {
                this.attr.formValues = Object.assign({}, item)
                this.attr.type = 'update'
                this.attr.isMultipleUpdate = false
                this.openObjectSlider('update', item['bk_inst_name'])
            },
            // 打开侧滑界面
            openObjectSlider (type, instName) {
                this.tab.activeName = 'attr'
                this.attr.type = type
                if (type === 'create') {
                    this.attr.isMultipleUpdate = false
                    this.slider.title.icon = 'icon-cc-create-business'
                    this.slider.title.text = `${this.$t("Common['创建']")} ${this.objName}`
                } else {
                    this.slider.title.icon = 'icon-cc-edit'
                    let name = this.objId === 'biz' ? this.attr.formValues['bk_biz_name'] : this.attr.formValues['bk_inst_name']
                    this.slider.title.text = `${this.$t("Common['编辑']")} ${this.objName} ${name}`
                }
                this.slider.isShow = true
            },
            // 关闭属性侧滑界面
            closeObjectSlider () {
                this.slider.isShow = false
                this.attr.formValues = {}
            },
            // 关闭导入侧滑
            closeImportSlider () {
                this.importSlider.isShow = false
            },
            // 计算关联属性单元格显示的值
            getAssociateCell (data) {
                let label = []
                if (Array.isArray(data)) {
                    data.map(({bk_inst_name: bkInstName}) => {
                        if (bkInstName) {
                            label.push(bkInstName)
                        }
                    })
                }
                return label.join(',')
            },
            getEnumCell (data, property) {
                let obj = property.option.find(({id}) => {
                    return id === data
                })
                if (obj) {
                    return obj.name
                }
            },
            getEnumOptions () {
                let selectedPropertyId = this.filter.selected
                let property = this.attr.formFields.find(({bk_property_id: bkPropertyId}) => bkPropertyId === selectedPropertyId)
                if (property) {
                    return property.option || []
                }
                return []
            },
            updateHistoryCount (isShow) {
                if (this.prevHistoryCount) {
                    this.$store.commit('navigation/updateHistoryCount', isShow ? -1 : 1)
                }
            }
        },
        created () {
            this.prevHistoryCount = this.$store.state.navigation.historyCount
        },
        mounted () {
            this.initTable()
        },
        components: {
            vObjectTable,
            vObjectAttr,
            vRelevance,
            vHistory,
            vImport,
            vSideslider,
            vConfigField,
            vDeleteHistory,
            vAssociationList,
            vNewAssociation
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #f9f9f9; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    $primaryHoverColor: #6b7baa; //鼠标移上 主要颜色
    .main-btn{  //主要按钮
        background: $primaryHoverColor;
        &:hover{
            background: #4d597d;
        }
    }
    .no-border-btn{    //无边框按钮
        background: #fff;
        color: $primaryHoverColor;
        cursor: pointer;
        &:hover{
            background: $primaryHoverColor;
            color: #fff;
        }
    }
    .btn-group {
        font-size: 0;
        .dropdown-btn {
            width: 75px;
        }
    }
    .bk-button{
        &.import,
        &.export,
        &.setting{
            width: 111px;
            &:disabled{
                cursor: not-allowed !important;
                background-color: #fafafa;
                border: 1px solid #e7e9ef;
                i{
                    color: #bec6de;
                }
                &:hover{
                    color: #bec6de;
                    background-color: #fafafa;
                    border: 1px solid #e7e9ef;
                    i{
                        color: #bec6de;
                    }
                }
            }
        }
        &.setting{
            margin: 0 0 0 10px;
            width: 36px;
            padding: 0;
            text-align: center;
        }
        .icon-cc-derivation,
        .icon-cc-import,
        .icon-cc-history,
        .icon-cc-setting{
            font-size: 16px;
            position: relative;
            top: -2px;
        }
    }
    .host-resource-wrapper{
        .host-resource-title{
            display: inline-block;
            vertical-align: middle;
            height: 36px;
            line-height: 36px;
            font-size: 14px;
            margin: 0;
        }
        .bottom-contain{
            padding:20px 20px 0;
        }
    }
    .bk-button-componey{
        height: 36px;
    }
    .table-contain{
        padding:20px;
        width:100%;
        background:$defaultColor;
        float:left;
        .slide-content{
            height: calc(100% - 120px);
            padding: 0 20px;
            overflow: auto;
            &::-webkit-scrollbar{
                width: 6px;
                height: 5px;
            }
            &::-webkit-scrollbar-thumb{
                border-radius: 20px;
                background: #a5a5a5;
            }
        }
        .title{
            position: relative;
            .close{
                position: absolute;
                top: 0;
                left: -30px;
                width: 30px;
                height: 60px;
                padding: 10px 7px 0;
                background-color: #ef4c4c;
                box-shadow: -2px 0 2px 0 rgba(0, 0, 0, 0.2);
                cursor: pointer;
                color: #fff;
                font-size: 14px;
                line-height: 20px;
                font-weight: normal;
                &:hover{
                    background: #e13d3d;
                }
            }
        }
    }
    .bk-form-title{
        font-size: 14px;
        color: #6b7baa;
        font-weight: bold;
        padding-left: 8px;
        position: relative;
        line-height: 20px;
        margin: 6px 0;
        &:before{
            content: "";
            width: 4px;
            height: 12px;
            position: absolute;
            left: 0;
            top: 4px;
            background: #6b7baa;
        }
    }
    .menber-tip{
        position: absolute;
        left:0;
        color: #ff5656;
        top: 6px;
        left: 390px;
    }
    .bk-form{
        height: calc(100% - 60px);
    }
    .bk-form .bk-label{
        width: 113px;
        padding-right: 17px;
        color: #6b7baa;
        .bk-label-span{
            display: inline-block;
            text-overflow: ellipsis;
            width: 100%;
            white-space: nowrap;
            overflow: hidden;
        }
    }
    .bk-form-content{
        margin-left: 113px;
    }
    .bk-form-input{
        color: #6b7baa;
    }
    .content{
        margin: 40px 0 0 40px;
        padding: 0 40px 30px;
        // padding: 40px 40px 30px 40px;
    }
    .busi-content-button{
        background: #fff;
        margin-left: 114px;
    }
    .table-page-contain{
        width: 100%;
        padding: 14px 20px;
        background: $primaryColor;
        .bk-page-compact{
            float: right;
        }
        .page-info{
            font-size: 12px;
            color: #838fb6;
            padding-top: 10px;
        }
    }
    .quick-search{
        position: relative;
        .left-select{
            position: relative;
            border-radius: 2px 0 0 2px;
            margin-right: -1px;
            z-index: 2;
            &:hover,
            &:active,
            &:focus{
                z-index: 2;
            }
        }
        .datatimepicker{
            float: right;
            width: 320px;
        }
        .search-text{
            position: relative;
            width: 320px;
            height: 36px;
            padding: 0 28px 0 10px;
            line-height: 36px;
            border-radius: 0 2px 2px 0;
            z-index: 1;
            &:focus{
                z-index: 2;
            }
        }
        .search-options{
            width: 320px;
        }
        .icon-search{
            position: absolute;
            right: 10px;
            top: 11px;
            cursor: pointer;
            z-index: 3;
        }
    }
</style>

<style media="screen" lang="scss">
    .host-resource-wrapper{
        position: relative;
        .errorInfo-wrapper{
            input[name="date-select"]{
                border: 1px solid #ff3737;
            }
            input{
                border: 1px solid #ff3737;
            }
            .errorInfo{
                font-size: 12px;
                color: #ff3737;
            }
        }
        .quick-search{
            .left-select{
                width: 115px;
                .bk-select-input{
                    }
                
                .bk-form-input{
                    position: relative;
                    &:focus{
                        z-index: 2;
                    }
                }
            }
        }
    }
    .bk-form-item.is-required .bk-label:after {
        position: absolute !important;
    }
    .host-resource-wrapper{
        .slide-content{
            .bk-tab2{
                height: 100%;
                .bk-tab2-content{
                    height: calc(100% - 60px);
                    overflow: auto;
                    &::-webkit-scrollbar {
                        width: 6px;
                        height: 5px;
                    }
                    &::-webkit-scrollbar-thumb {
                        border-radius: 20px;
                        background: #a5a5a5;
                        -webkit-box-shadow: inset 0 0 6px hsla(0,0%,80%,.3);
                    }
                }
            }
        }
    }
    .business-wrapper{
        .bk-tab2{
            .bk-tab2-nav{
                .tab2-nav-item{
                    padding: 0 32px;
                }
            }
        }
    }
</style>
