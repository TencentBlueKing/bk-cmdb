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
            <bk-button class="process-btn" type="primary" @click="createProcess">{{$t("ProcessManagement['新增进程']")}}</bk-button>
            <div class="filter-text fr">
                <input type="text" class="bk-form-input" :placeholder="$t('ProcessManagement[\'进程名称搜索\']')" 
                    v-model.trim="filter.text" @keyup.enter="setCurrentPage(1)">
                    <i class="bk-icon icon-search" @click="setCurrentPage(1)"></i>
            </div>
        </div>
        <cmdb-table class="process-table" ref="table"
            :loading="$loading()"
            :checked.sync="table.checked"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort"
            :wrapperMinusHeight="157"
            @handleRowClick="handleRowClick"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange"
            @handleCheckAll="handleCheckAll"></cmdb-table>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        data () {
            return {
                filter: {
                    bizId: '',
                    text: ''
                },
                table: {
                    header: [],
                    list: [],
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
        created () {
            this.reload()
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            createProcess () {

            },
            multipleUpdate () {

            },
            async reload () {
                this.properties = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: this.objId,
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: `${this.objId}Attribute`,
                        fromCache: true
                    }
                })
                // await Promise.all([
                //     this.getPropertyGroups(),
                //     this.setTableHeader(),
                //     this.setFilterOptions()
                // ])
                this.getTableData()
            },
            // getPropertyGroups () {
            //     return this.searchGroup({
            //         objId: this.objId,
            //         config: {
            //             fromCache: true,
            //             requestId: `${this.objId}AttributeGroup`
            //         }
            //     }).then(groups => {
            //         this.propertyGroups = groups
            //         return groups
            //     })
            // },
            async handleCheckAll (type) {
                if (type === 'current') {
                    this.table.checked = this.table.list.map(inst => inst['bk_inst_id'])
                } else {
                    // const allData = await this.getAllInstList()
                    // this.table.checked = allData.info.map(inst => inst['bk_inst_id'])
                }
            },
            handleRowClick (item) {
                this.slider.show = true
                // this.slider.title = `${this.$t("Common['编辑']")} ${item['bk_inst_name']}`
                // this.attribute.inst.details = item
                // this.attribute.type = 'details'
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
                width: 135px;
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
