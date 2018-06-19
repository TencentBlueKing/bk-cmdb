<template>
    <div class="process-wrapper">
        <v-breadcrumb class="breadcrumbs"></v-breadcrumb>
        <div class="process-filter clearfix">
            <v-application-selector class="filter-selector fl"
                :filterable="true" 
                :selected.sync="filter.bkBizId">
            </v-application-selector>
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
            :pagination="table.pagination"
            :loading="table.isLoading"
            :wrapperMinusHeight="150"
            @handlePageChange="setCurrentPage"
            @handleSizeChange="setCurrentSize"
            @handleSortChange="setCurrentSort"
            @handleRowClick="showProcessAttribute">
            <template v-for="({property, id, name}) in table.header" :slot="id" slot-scope="{ item }">
                <template v-if="property['bk_property_type'] === 'enum'">
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
                        :formFields="attribute"
                        :formValues="slider.tab.attribute.formValues"
                        :objId="'process'"
                        @submit="saveProcess"
                        @delete="deleteProcess">
                    </v-attribute>
                </bk-tabpanel>
                <bk-tabpanel name="bind" :title="$t('ProcessManagement[\'模块绑定\']')" :show="slider.tab.type === 'update'">
                    <v-module :bkProcessId="slider.tab.module.bkProcessId" :bkBizId="filter.bkBizId"></v-module>
                </bk-tabpanel>
                <bk-tabpanel name="history" :title="$t('HostResourcePool[\'变更记录\']')" :show="slider.tab.type === 'update'">
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
    import vBreadcrumb from '@/components/common/breadcrumb/breadcrumb'
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
            }
        },
        watch: {
            attribute (attribute) {
                let headerLead = []
                let headerMiddle = []
                let headerTail = []
                attribute.map(property => {
                    let {
                        isonly,
                        isrequired,
                        'bk_isapi': bkIsapi,
                        'bk_property_id': bkPropertyId,
                        'bk_property_name': bkPropertyName
                    } = property
                    let headerItem = {
                        id: bkPropertyId,
                        name: bkPropertyName,
                        property: property
                    }
                    if (!bkIsapi) {
                        if (isonly && isrequired) {
                            headerLead.push(headerItem)
                        } else if (isonly || isonly) {
                            headerMiddle.push(headerItem)
                        } else {
                            headerTail.push(headerItem)
                        }
                    }
                })
                this.table.header = headerLead.concat(headerMiddle, headerTail).slice(0, 6)
            },
            'filter.bkBizId' (bkBizId) {
                this.$nextTick(() => {
                    this.setCurrentPage(1)
                })
            }
        },
        methods: {
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
                this.slider.tab.attribute.type = 'create'
                this.slider.tab.attribute.formValues = {}
                this.slider.isShow = true
            },
            saveProcess (formData, originalData) {
                if (this.slider.tab.attribute.type === 'create') {
                    this.$axios.post(`proc/${this.bkSupplierAccount}/${this.filter.bkBizId}`, formData).then(res => {
                        if (res.result) {
                            this.$alertMsg(this.$t("ProcessManagement['新增进程成功']"), 'success')
                            this.setCurrentPage(1)
                            this.closeSlider()
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                } else {
                    this.$axios.put(`proc/${this.bkSupplierAccount}/${this.filter.bkBizId}/${originalData['bk_process_id']}`, formData).then(res => {
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
                        this.$axios.delete(`proc/${this.bkSupplierAccount}/${this.filter.bkBizId}/${data['bk_process_id']}`).then((res) => {
                            if (res.result) {
                                this.closeSlider()
                                this.setCurrentPage(1)
                                this.$alertMsg(this.$t("ProcessManagement['删除进程成功']"), 'success')
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                })
            },
            showProcessAttribute (item) {
                this.slider.title.text = this.$t("ProcessManagement['编辑进程']")
                this.slider.tab.active = 'attribute'
                this.slider.tab.type = 'update'
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
            vHistory,
            vBreadcrumb
        }
    }
</script>
<style lang="scss" scoped>
    .process-wrapper{
        padding: 0 20px;
        .breadcrumbs{
            padding: 8px 0;
        }
    }
    .process-filter{
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