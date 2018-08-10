<template>
    <div class="process-wrapper">
        <div class="process-filter clearfix">
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
            <v-application-selector class="filter-selector fl"
                :filterable="true" 
                :selected.sync="filter.bkBizId">
            </v-application-selector>
            <button class="mr10 fl bk-button bk-default"
                :disabled="!table.chooseId.length" 
                @click="multipleUpdate">
                <i class="icon-cc-edit"></i>
                <span>{{$t("BusinessTopology['修改']")}}</span>
            </button>
            <bk-button class="fl" type="primary" @click="createProcess">{{$t("ProcessManagement['新增进程']")}}</bk-button>
            <div class="filter-text fr">
                <input type="text" class="bk-form-input" :placeholder="$t('ProcessManagement[\'进程名称搜索\']')" 
                    v-model.trim="filter.text" @keyup.enter="setCurrentPage(1)">
                    <i class="bk-icon icon-search" @click="setCurrentPage(1)"></i>
            </div>
        </div>
        <v-table class="process-table" 
            :header="table.header"
            :list="table.list"
            :defaultSort="table.defaultSort"
            :pagination.sync="table.pagination"
            :loading="table.isLoading"
            :wrapperMinusHeight="150"
            :checked="table.chooseId"
            :isCheckboxShow="false"
            @handlePageChange="setCurrentPage"
            @handleSizeChange="setCurrentSize"
            @handleSortChange="setCurrentSort"
            @handleCheckAll="getAllProcessID"
            @handleRowClick="showProcessAttribute">
            <template v-for="({property, id, name}) in table.header" :slot="id" slot-scope="{ item }">
                <template v-if="id === 'bk_process_id'">
                    <label style="width:100%;text-align:center;" class="bk-form-checkbox bk-checkbox-small" @click.stop>
                        <input type="checkbox"
                            :value="item['bk_process_id']" 
                            v-model="table.chooseId">
                    </label>
                </template>
                <template v-else-if="property['bk_property_type'] === 'enum'">
                    {{getEnumCell(item[id], property)}}
                </template>
                <template v-else-if="!!property['bk_asst_obj_id']">
                    {{getAssociateCell(item[id])}}
                </template>
                <template v-else>{{item[id]}}</template>
            </template>
        </v-table>
        <v-side-slider
            :isShow.sync="slider.isShow"
            :isCloseConfirmShow="slider.isCloseConfirmShow"
            :hasCloseConfirm="true"
            :title="slider.title"
            @closeSlider="closeSliderConfirm">
            <bk-tab :active-name="slider.tab.active" @tab-changed="tabChanged" slot="content" :class="['tab-wrapper', slider.tab.active]" style="border:none;">
                <bk-tabpanel name="attribute" :title="$t('ProcessManagement[\'属性\']')">
                    <v-attribute :active="slider.isShow && slider.tab.active === 'attribute'"
                        ref="attribute"
                        :type="slider.tab.attribute.type"
                        :formFields="formFields"
                        :formValues="slider.tab.attribute.formValues"
                        :isMultipleUpdate="slider.tab.attribute.isMultipleUpdate"
                        :objId="'process'"
                        @submit="saveProcess"
                        @delete="deleteProcess">
                    </v-attribute>
                </bk-tabpanel>
                <bk-tabpanel name="bind" :title="$t('ProcessManagement[\'模块绑定\']')" :show="slider.tab.type === 'update' && !slider.tab.attribute.isMultipleUpdate">
                    <v-module :bkProcessId="slider.tab.module.bkProcessId" :bkBizId="filter.bkBizId"></v-module>
                </bk-tabpanel>
                <bk-tabpanel name="history" :title="$t('HostResourcePool[\'变更记录\']')" :show="slider.tab.type === 'update' && !slider.tab.attribute.isMultipleUpdate">
                    <v-history :active="slider.tab.active === 'history'" :type="'process'" :instId="slider.tab.history.bkProcessId"></v-history>
                </bk-tabpanel>
            </bk-tab>
        </v-side-slider>
    </div>
</template>
<script>
    import vApplicationSelector from '@/components/common/selector/application'
    import vTable from '@/components/table/table'
    import vSideSlider from '@/components/slider/sideslider'
    import vAttribute from '@/components/object/attribute'
    import vHistory from '@/components/history/history'
    import vModule from './children/module'
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                filter: {
                    bkBizId: '',
                    text: ''
                },
                attribute: [],
                table: {
                    header: [],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    chooseId: [],
                    defaultSort: '-bk_process_id',
                    sort: '-bk_process_id',
                    isLoading: false
                },
                slider: {
                    isShow: false,
                    isCloseConfirmShow: false,
                    title: {
                        text: this.$t("ProcessManagement['新增进程']"),
                        icon: 'icon-cc-cube-entity'
                    },
                    tab: {
                        active: 'attribute',
                        type: 'create',
                        attribute: {
                            type: 'create',
                            isMultipleUpdate: false,
                            formValues: {}
                        },
                        module: {
                            bkProcessId: 0
                        },
                        history: {
                            bkProcessId: 0
                        }
                    }
                }
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            searchParams () {
                let params = {
                    condition: {
                        'bk_biz_id': this.filter.bkBizId
                    },
                    fields: [],
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.size,
                        limit: this.table.pagination.size,
                        sort: this.table.sort
                    }
                }
                if (this.filter.text) {
                    params['condition']['bk_process_name'] = this.filter.text
                }
                return params
            },
            formFields () {
                let list = this.attribute
                if (this.slider.tab.attribute.isMultipleUpdate) {
                    list = this.attribute.filter(({bk_property_id: bkPropertyId}) => {
                        return bkPropertyId !== 'bk_process_name' && bkPropertyId !== 'bk_func_id'
                    })
                }
                return list
            }
        },
        watch: {
            attribute (attribute) {
                let headerLead = []
                let headerMiddle = []
                let headerTail = []
                attribute.map(property => {
                    let {
                        'bk_property_id': bkPropertyId,
                        'bk_property_name': bkPropertyName
                    } = property
                    let headerItem = {
                        id: bkPropertyId,
                        name: bkPropertyName,
                        property: property
                    }
                    switch (bkPropertyId) {
                        case 'bk_process_name':
                            headerMiddle[0] = headerItem
                            break
                        case 'bk_func_id':
                            headerMiddle[1] = headerItem
                            break
                        case 'bind_ip':
                            headerMiddle[2] = headerItem
                            break
                        case 'port':
                            headerMiddle[3] = headerItem
                            break
                        case 'protocol':
                            headerMiddle[4] = headerItem
                            break
                        case 'bk_func_name':
                            headerMiddle[5] = headerItem
                            break
                    }
                })
                headerMiddle = headerMiddle.filter(header => {
                    return header !== undefined
                })
                headerLead.unshift({
                    id: 'bk_process_id',
                    name: 'bk_process_id',
                    type: 'checkbox',
                    width: 50
                })
                this.table.header = headerLead.concat(headerMiddle, headerTail).slice(0, 7)
            },
            'filter.bkBizId' (bkBizId) {
                this.$nextTick(() => {
                    this.setCurrentPage(1)
                    this.table.chooseId = []
                })
            }
        },
        methods: {
            tableChecked (type) {
                if (type === 'all') {
                    this.getAllProcessID(true)
                } else if (type === 'cancel') {
                    this.table.chooseId = []
                } else {
                    let chooseId = []
                    this.table.list.map(({bk_process_id: bkProcessId}) => chooseId.push(bkProcessId))
                    this.table.chooseId = chooseId
                }
            },
            multipleUpdate () {
                this.slider.title.text = this.$t("ProcessManagement['批量编辑']")
                this.slider.tab.attribute.isMultipleUpdate = true
                this.slider.tab.attribute.type = 'create'
                this.slider.tab.attribute.formValues = {bk_process_id: this.table.chooseId.join(',')}
                this.slider.isShow = true
            },
            async getAllProcessID (isChecked) {
                if (isChecked) {
                    let allProcessId = []
                    let params = this.$deepClone(this.searchParams)
                    params.page = {}
                    this.table.isLoading = true
                    try {
                        const res = await this.$axios.post(`proc/search/${this.bkSupplierAccount}/${this.filter.bkBizId}`, params)
                        res.data.info.map((item, index) => {
                            allProcessId.push(item['bk_process_id'])
                        })
                        this.table.chooseId = allProcessId
                    } catch (e) {
                        this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                    } finally {
                        this.table.isLoading = false
                    }
                } else {
                    this.table.chooseId = []
                }
            },
            closeSliderConfirm () {
                this.slider.tab.module.bkProcessId = null
                this.slider.isCloseConfirmShow = this.$refs.attribute.isCloseConfirmShow()
            },
            getEnumCell (data, property) {
                let obj = property.option.find(({id}) => {
                    return id === data
                })
                if (obj) {
                    return obj.name
                }
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
            createProcess () {
                this.slider.title.text = this.$t("ProcessManagement['新增进程']")
                this.slider.tab.active = 'attribute'
                this.slider.tab.type = 'create'
                this.slider.tab.attribute.isMultipleUpdate = false
                this.slider.tab.attribute.type = 'create'
                this.slider.tab.attribute.formValues = {}
                this.slider.isShow = true
            },
            saveProcess (formData, originalData) {
                if (this.slider.tab.attribute.type === 'create' && !this.slider.tab.attribute.isMultipleUpdate) {
                    this.$axios.post(`proc/${this.bkSupplierAccount}/${this.filter.bkBizId}`, formData, {id: 'editAttr'}).then(res => {
                        if (res.result) {
                            this.$alertMsg(this.$t("ProcessManagement['新增进程成功']"), 'success')
                            this.setCurrentPage(1)
                            this.closeSlider()
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                } else if (!this.slider.tab.attribute.isMultipleUpdate) {
                    this.$axios.put(`proc/${this.bkSupplierAccount}/${this.filter.bkBizId}/${originalData['bk_process_id']}`, formData, {id: 'editAttr'}).then(res => {
                        if (res.result) {
                            this.$alertMsg(this.$t("ProcessManagement['修改进程成功']"), 'success')
                            this.setCurrentPage(1)
                            this.closeSlider()
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                } else {
                    let {
                        bk_process_id: bkProcessId
                    } = originalData
                    this.$axios.put(`proc/${this.bkSupplierAccount}/${this.filter.bkBizId}`, {...formData, ...{bk_process_id: bkProcessId.toString()}}, {id: 'editAttr'}).then(res => {
                        if (res.result) {
                            this.$alertMsg(this.$t("ProcessManagement['修改进程成功']"), 'success')
                            this.setCurrentPage(1)
                            this.closeSlider()
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            deleteProcess (data) {
                this.$bkInfo({
                    title: this.$t("ProcessManagement['确认要删除进程']", {processName: data['bk_process_name']}),
                    confirmFn: () => {
                        this.$axios.delete(`proc/${this.bkSupplierAccount}/${this.filter.bkBizId}/${data['bk_process_id']}`, {id: 'instDelete'}).then((res) => {
                            if (res.result) {
                                this.closeSlider()
                                this.setCurrentPage(1)
                                this.$alertMsg(this.$t("ProcessManagement['删除进程成功']"), 'success')
                                this.table.chooseId = this.table.chooseId.filter(id => id !== data['bk_process_id'])
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                })
            },
            showProcessAttribute (item) {
                this.slider.title.text = `${this.$t("ProcessManagement['编辑进程']")}`
                this.slider.tab.active = 'attribute'
                this.slider.tab.type = 'update'
                this.slider.tab.attribute.isMultipleUpdate = false
                this.slider.tab.attribute.type = 'update'
                this.slider.tab.attribute.formValues = {}
                this.slider.tab.module.bkProcessId = item['bk_process_id']
                this.slider.tab.history.bkProcessId = item['bk_process_id']
                this.slider.isShow = true
                this.$axios.get(`proc/${this.bkSupplierAccount}/${this.filter.bkBizId}/${item['bk_process_id']}`).then(res => {
                    if (res.result) {
                        let values = {
                            'bk_process_id': item['bk_process_id']
                        }
                        res.data.map(column => {
                            values[column['bk_property_id']] = column['bk_property_value']
                        })
                        this.slider.tab.attribute.formValues = values
                        this.slider.title.text = `${this.$t("ProcessManagement['编辑进程']")} ${this.slider.tab.attribute.formValues['bk_process_name']}`
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            tabChanged (active) {
                this.slider.tab.active = active
            },
            getProcessAttribute () {
                let params = {
                    'bk_obj_id': 'process',
                    'bk_supplier_account': this.bkSupplierAccount
                }
                this.$axios.post('object/attr/search', params).then(res => {
                    if (res.result) {
                        this.attribute = res.data
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            getTableList () {
                this.table.isLoading = true
                this.$axios.post(`proc/search/${this.bkSupplierAccount}/${this.filter.bkBizId}`, this.searchParams).then(res => {
                    if (res.result) {
                        this.table.list = res.data.info
                        this.table.pagination.count = res.data.count
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.table.isLoading = false
                }).catch((e) => {
                    this.table.isLoading = false
                    this.table.list = []
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t("Common['您没有当前业务的权限']"))
                    }
                })
            },
            setCurrentSize (size) {
                this.table.pagination.size = size
                this.setCurrentPage(1)
            },
            setCurrentSort (sort) {
                this.table.sort = sort
                this.setCurrentPage(1)
            },
            setCurrentPage (page) {
                this.table.pagination.current = page
                this.getTableList()
            },
            closeSlider () {
                this.slider.isShow = false
            },
            showFiling () {
                this.filing.isShow = true
            }
        },
        created () {
            this.getProcessAttribute()
        },
        components: {
            vApplicationSelector,
            vTable,
            vSideSlider,
            vAttribute,
            vModule,
            vHistory
        }
    }
</script>
<style lang="scss" scoped>
    .process-wrapper{
        height: 100%;
        padding: 20px;
    }
    .process-filter{
        .dropdown-btn {
            width: 75px;
        }
        .filter-selector{
            width: 170px;
            margin-right: 10px;
        }
        .filter-text{
            position: relative;
            .bk-form-input{
                width: 320px;
            }
            .icon-search{
                position: absolute;
                right: 10px;
                top: 11px;
                cursor: pointer;
            }
        }
    }
    .process-table{
        margin-top: 20px;
    }
    .tab-wrapper{
        height: calc(100% - 60px);
        padding: 0 20px;
        &.attribute{
            height: calc(100% - 121px);
        }
    }
</style>