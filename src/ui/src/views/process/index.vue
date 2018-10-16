<template>
    <div class="process-wrapper">
        <div class="process-filter clearfix">
            <cmdb-business-selector class="business-selector" v-model="filter.bizId">
            </cmdb-business-selector>
            <bk-button class="process-btn"
                type="default"
                :disabled="!table.checked.length" 
                @click="handleMultipleEdit">
                <i class="icon-cc-edit"></i>
                <span>{{$t("BusinessTopology['修改']")}}</span>
            </bk-button>
            <bk-button class="process-btn" type="primary" @click="handleCreate">{{$t("ProcessManagement['新增进程']")}}</bk-button>
            <div class="filter-text fr">
                <input type="text" class="bk-form-input" :placeholder="$t('ProcessManagement[\'进程名称搜索\']')" 
                    v-model.trim="filter.text" @keyup.enter="handlePageChange(1)">
                    <i class="bk-icon icon-search" @click="handlePageChange(1)"></i>
            </div>
        </div>
        <cmdb-table class="process-table" ref="table"
            :loading="$loading('post_searchProcess_list')"
            :checked.sync="table.checked"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort"
            :wrapperMinusHeight="300"
            @handleRowClick="handleRowClick"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange"
            @handleCheckAll="handleCheckAll">
            <div class="empty-info" slot="data-empty">
                <p>{{$t("Common['暂时没有数据']")}}</p>
                <p>{{$t("ProcessManagement['当前业务并无进程，可点击下方按钮新增']")}}</p>
                <bk-button class="process-btn" type="primary" @click="handleCreate">{{$t("ProcessManagement['新增进程']")}}</bk-button>
            </div>
        </cmdb-table>
        <cmdb-slider :isShow.sync="slider.show" :title="slider.title">
            <bk-tab :active-name.sync="tab.active" slot="content">
                <bk-tabpanel name="attribute" :title="$t('Common[\'属性\']')" style="width: calc(100% + 40px);margin: 0 -20px;">
                    <cmdb-details v-if="attribute.type === 'details'"
                        :properties="properties"
                        :propertyGroups="propertyGroups"
                        :inst="attribute.inst.details"
                        @on-edit="handleEdit"
                        @on-delete="handleDelete">
                    </cmdb-details>
                    <cmdb-form v-else-if="['update', 'create'].includes(attribute.type)"
                        :properties="properties"
                        :propertyGroups="propertyGroups"
                        :inst="attribute.inst.edit"
                        :type="attribute.type"
                        @on-submit="handleSave"
                        @on-cancel="handleCancel">
                    </cmdb-form>
                    <cmdb-form-multiple v-else-if="attribute.type === 'multiple'"
                        :properties="properties"
                        :propertyGroups="propertyGroups"
                        @on-submit="handleMultipleSave"
                        @on-cancel="handleMultipleCancel">
                    </cmdb-form-multiple>
                </bk-tabpanel>
                <bk-tabpanel name="moduleBind" :title="$t('ProcessManagement[\'模块绑定\']')" :show="attribute.type === 'details'">
                    <v-module v-if="tab.active === 'moduleBind'"
                        :processId="attribute.inst.details['bk_process_id']"
                        :bizId="filter.bizId"
                    ></v-module>
                </bk-tabpanel>
                <bk-tabpanel name="history" :title="$t('HostResourcePool[\'变更记录\']')" :show="attribute.type === 'details'">
                    <cmdb-audit-history v-if="tab.active === 'history'"
                        target="process"
                        :instId="attribute.inst.details['bk_process_id']">
                    </cmdb-audit-history>
                </bk-tabpanel>
            </bk-tab>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import cmdbAuditHistory from '@/components/audit-history/audit-history'
    import vModule from './module'
    export default {
        components: {
            cmdbAuditHistory,
            vModule
        },
        data () {
            return {
                properties: [],
                slider: {
                    show: false,
                    title: ''
                },
                attribute: {
                    type: null,
                    inst: {
                        details: {},
                        edit: {}
                    }
                },
                tab: {
                    active: 'attribute'
                },
                filter: {
                    bizId: '',
                    text: '',
                    businessResolver: null
                },
                table: {
                    header: [],
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    checked: [],
                    defaultSort: '-bk_process_id',
                    sort: '-bk_process_id'
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount'])
        },
        watch: {
            'filter.bizId' () {
                if (this.filter.businessResolver) {
                    this.filter.businessResolver()
                } else {
                    this.table.checked = []
                    this.handlePageChange(1)
                }
            }
        },
        created () {
            this.reload()
        },
        methods: {
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('procConfig', [
                'searchProcess',
                'deleteProcess',
                'createProcess',
                'updateProcess',
                'batchUpdateProcess',
                'searchProcessById'
            ]),
            handleMultipleEdit () {
                this.attribute.type = 'multiple'
                this.slider.title = this.$t('Inst[\'批量更新\']')
                this.slider.show = true
            },
            handleMultipleSave (values) {
                this.batchUpdateProcess({
                    bizId: this.filter.bizId,
                    params: {
                        ...values,
                        bk_process_id: this.table.checked.join(',')
                    },
                    config: {
                        requestId: `processBatchUpdate`
                    }
                }).then(() => {
                    this.$success(this.$t('Common["修改成功"]'))
                    this.handlePageChange(1)
                })
            },
            handleMultipleCancel () {
                this.slider.show = false
            },
            getBusiness () {
                return new Promise((resolve, reject) => {
                    this.filter.businessResolver = () => {
                        this.filter.businessResolver = null
                        resolve()
                    }
                })
            },
            async reload () {
                this.properties = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: 'process',
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: `post_searchObjectAttribute_process`,
                        fromCache: true
                    }
                })
                await Promise.all([
                    this.getBusiness(),
                    this.getPropertyGroups(),
                    this.setTableHeader()
                ])
                this.handlePageChange(1)
            },
            getPropertyGroups () {
                return this.searchGroup({
                    objId: 'process',
                    config: {
                        fromCache: true,
                        requestId: 'post_searchGroup_process'
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            setTableHeader () {
                let header = []
                let headerMap = ['bk_process_name', 'bk_func_id', 'bind_ip', 'port', 'protocol', 'bk_func_name']
                this.properties.map(property => {
                    let {
                        'bk_property_id': propertyId,
                        'bk_property_name': propertyName
                    } = property
                    let index = headerMap.indexOf(propertyId)
                    if (index !== -1) {
                        header[index] = {
                            id: propertyId,
                            name: propertyName
                        }
                    }
                })
                header.unshift({
                    id: 'bk_process_id',
                    type: 'checkbox',
                    width: 50
                })
                this.table.header = header
            },
            async handleEdit (flatternItem) {
                const list = await this.getProcessList({fromCache: true})
                const inst = list.info.find(item => item['bk_process_id'] === flatternItem['bk_process_id'])
                this.attribute.inst.edit = inst
                this.attribute.type = 'update'
            },
            handleCreate () {
                this.attribute.type = 'create'
                this.attribute.inst.edit = {}
                this.slider.show = true
                this.slider.title = `${this.$t("Common['创建']")} ${this.$model['bk_obj_name']}`
            },
            handleDelete (process) {
                this.$bkInfo({
                    title: this.$t("Common['确认要删除']", {name: process['bk_process_name']}),
                    confirmFn: () => {
                        this.deleteProcess({
                            bizId: this.filter.bizId,
                            processId: process['bk_process_id']
                        }).then(() => {
                            this.slider.show = false
                            this.$success(this.$t('Common["删除成功"]'))
                            this.handlePageChange(1)
                        })
                    }
                })
            },
            handleSave (values, changedValues, originalValues, type) {
                if (type === 'update') {
                    this.updateProcess({
                        bizId: this.filter.bizId,
                        processId: originalValues['bk_process_id'],
                        params: values
                    }).then(() => {
                        this.getTableData()
                        this.searchProcessById({
                            bizId: this.filter.bizId,
                            processId: originalValues['bk_process_id']
                        }).then(process => {
                            this.attribute.inst.details = this.$tools.flatternItem(this.properties, process)
                        })
                        this.handleCancel()
                        this.$success(this.$t("Common['修改成功']"))
                    })
                } else {
                    this.createProcess({
                        params: values,
                        bizId: this.filter.bizId
                    }).then(() => {
                        this.handlePageChange(1)
                        this.handleCancel()
                        this.$success(this.$t("Inst['创建成功']"))
                    })
                }
            },
            handleCancel () {
                if (this.attribute.type === 'create') {
                    this.slider.show = false
                } else {
                    this.attribute.type = 'details'
                }
            },
            getProcessList (config = {cancelPrevious: true}) {
                return this.searchProcess({
                    bizId: this.filter.bizId,
                    params: this.getSearchParams(),
                    config: Object.assign({requestId: 'post_searchProcess_list'}, config)
                })
            },
            getAllProcessList () {
                return this.searchProcess({
                    bizId: this.filter.bizId,
                    params: {
                        ...this.getSearchParams(),
                        page: {}
                    },
                    config: {
                        requestId: `post_searchProcess_list`,
                        cancelPrevious: true
                    }
                })
            },
            getTableData () {
                this.getProcessList().then(data => {
                    this.table.list = this.$tools.flatternList(this.properties, data.info)
                    this.table.pagination.count = data.count
                    return data
                })
            },
            getSearchParams () {
                let params = {
                    condition: {
                        'bk_biz_id': this.filter.bizId
                    },
                    fields: [],
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.size,
                        limit: this.table.pagination.size,
                        sort: this.table.sort
                    }
                }
                if (this.filter.text !== '') {
                    params['condition']['bk_process_name'] = this.filter.text
                }
                return params
            },
            async handleCheckAll (type) {
                if (type === 'current') {
                    this.table.checked = this.table.list.map(inst => inst['bk_process_id'])
                } else {
                    const allData = await this.getAllProcessList()
                    this.table.checked = allData.info.map(inst => inst['bk_process_id'])
                }
            },
            handleRowClick (item) {
                this.tab.active = 'attribute'
                this.slider.show = true
                this.slider.title = `${this.$t("Common['编辑']")} ${item['bk_process_name']}`
                this.attribute.inst.details = item
                this.attribute.type = 'details'
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .process-wrapper {
        .process-filter {
            .business-selector {
                float: left;
                width: 170px;
                margin-right: 10px;
            }
            .process-btn {
                float: left;
                margin-right: 10px;
            }
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
        .process-table {
            margin-top: 20px;
        }
    }
</style>
