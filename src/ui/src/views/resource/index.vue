<template>
    <div class="resource-layout clearfix">
        <cmdb-hosts-filter class="resource-filter fr"
            :active-tab="['filter']"
            :active-setting="['reset', 'filter-config']"
            :filter-config-key="filter.filterConfigKey"
            @on-refresh="handleRefresh">
            <div class="filter-group" slot="scope">
                <label class="filter-label">{{$t("Hosts['搜索范围']")}}</label>
                <cmdb-form-bool class="filter-field" :disabled="true" :checked="true">
                    <span class="filter-field-label">{{$t("Hosts['未分配主机']")}}</span>
                </cmdb-form-bool>
                <cmdb-form-bool class="filter-field" v-model="filter.assigned">
                    <span class="filter-field-label">{{$t("Hosts['已分配主机']")}}</span>
                </cmdb-form-bool>
            </div>
        </cmdb-hosts-filter>
        <div class="resource-main">
            <div class="resource-options clearfix">
                <cmdb-selector class="options-business-selector"
                    :placeholder="$t('HostResourcePool[\'分配到业务空闲机池\']')"
                    :disabled="!table.checked.length"
                    :list="business"
                    :auto-select="false"
                    setting-key="bk_biz_id"
                    display-key="bk_biz_name"
                    v-model="assignBusiness"
                    @on-selected="handleAssignHosts">
                </cmdb-selector>
                <bk-button class="options-button" type="default"
                    v-tooltip="$t('BusinessTopology[\'修改\']')"
                    :disabled="!table.checked.length"
                    @click="handleMultipleEdit">
                    <i class="icon-cc-edit"></i>
                </bk-button>
                <bk-button class="options-button options-button-delete" type="default"
                    v-tooltip="$t('Common[\'删除\']')"
                    :disabled="!table.checked.length"
                    @click="handleMultipleDelete">
                    <i class="icon-cc-del"></i>
                </bk-button>
                <bk-button class="options-button" type="submit default"
                    form="exportForm"
                    v-tooltip="$t('HostResourcePool[\'导出选中\']')"
                    :disabled="!table.checked.length">
                    <i class="icon-cc-derivation"></i>
                </bk-button>
                <form id="exportForm" :action="table.exportUrl" method="POST" hidden>
                    <input type="hidden" name="bk_host_id" :value="table.checked">
                    <input type="hidden" name="bk_biz_id" value="-1">
                </form>
                <cmdb-clipboard-selector class="options-button"
                    :list="clipboardList"
                    :disabled="!table.checked.length"
                    @on-copy="handleCopy">
                </cmdb-clipboard-selector>
                <div class="fr">
                    <bk-button class="options-button" type="default"
                        v-tooltip="$t('BusinessTopology[\'列表显示属性配置\']')"
                        @click="columnsConfig.show = true">
                        <i class="icon-cc-setting"></i>
                    </bk-button>
                    <bk-button class="options-button" type="default"
                        v-tooltip="$t('Common[\'查看删除历史\']')"
                        @click="routeToHistory">
                        <i class="icon-cc-history"></i>
                    </bk-button>
                    <bk-button class="options-button" type="primary"
                        v-tooltip="$t('HostResourcePool[\'导入主机\']')">
                        <i class="icon-cc-import"></i>
                    </bk-button>
                </div>
            </div>
            <cmdb-table class="resource-table"
                :loading="$loading()"
                :checked.sync="table.checked"
                :header="table.header"
                :list="table.list"
                :pagination.sync="table.pagination"
                :defaultSort="table.defaultSort"
                :wrapperMinusHeight="157"
                @handleRowClick="handleRowClick"
                @handleSortChange="handleSortChange"
                @handlePageChange="handlePageChange"
                @handleSizeChange="handleSizeChange"
                @handleCheckAll="handleCheckAll">
                <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                    <label style="width:100%;text-align:center;" class="bk-form-checkbox bk-checkbox-small"
                        v-if="header.id === 'bk_host_id'" 
                        @click.stop>
                        <input type="checkbox"
                            :value="item['host']['bk_host_id']" 
                            v-model="table.checked">
                    </label>
                    <span v-else>
                        {{getHostCellText(header, item)}}
                    </span>
                </template>
            </cmdb-table>
        </div>
        <cmdb-slider :isShow.sync="slider.show" :title="slider.title">
            <bk-tab :active-name.sync="tab.active" slot="content">
                <bk-tabpanel name="attribute" :title="$t('Common[\'属性\']')">
                    <cmdb-details v-if="tab.attribute.type === 'details'"
                        :properties="properties.host"
                        :propertyGroups="propertyGroups"
                        :inst="tab.attribute.inst.details"
                        :show-delete="false"
                        @on-edit="handleEdit">
                    </cmdb-details>
                    <cmdb-form v-else-if="tab.attribute.type === 'update'"
                        :properties="properties.host"
                        :propertyGroups="propertyGroups"
                        :inst="tab.attribute.inst.edit"
                        :type="tab.attribute.type"
                        @on-submit="handleSave"
                        @on-cancel="handleCancel">
                    </cmdb-form>
                    <cmdb-form-multiple v-else-if="tab.attribute.type === 'multiple'"
                        :properties="properties.host"
                        :propertyGroups="propertyGroups"
                        @on-submit="handleMultipleSave"
                        @on-cancel="handleMultipleCancel">
                    </cmdb-form-multiple>
                </bk-tabpanel>
                <bk-tabpanel name="relevance" :title="$t('HostResourcePool[\'关联\']')"></bk-tabpanel>
                <bk-tabpanel name="history" :title="$t('HostResourcePool[\'变更记录\']')">
                    <cmdb-audit-history v-if="tab.active === 'history'"
                        target="host"
                        :ext-key="{'$in': [tab.attribute.inst.details['bk_host_innerip']]}">
                    </cmdb-audit-history>
                </bk-tabpanel>
            </bk-tab>
        </cmdb-slider>
        <cmdb-slider
            :is-show.sync="columnsConfig.show"
            :width="600"
            :title="$t('BusinessTopology[\'列表显示属性配置\']')">
            <cmdb-columns-config slot="content"
                :properties="properties.host"
                :selected="columnsConfig.selected"
                @on-apply="handleApplyColumnsConfig"
                @on-cancel="columnsConfig.show = false"
                @on-reset="handleResetColumnsConfig">
            </cmdb-columns-config>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import hostsMixin from '@/mixins/hosts'
    import cmdbAuditHistory from '@/components/audit-history/audit-history.vue'
    export default {
        mixins: [hostsMixin],
        components: {
            cmdbAuditHistory
        },
        data () {
            return {
                table: {
                    columnsConfigKey: 'resource_table_columns'
                },
                filter: {
                    filterConfigKey: 'resource_filter_fields',
                    business: -1,
                    assigned: false
                },
                assignBusiness: ''
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['business'])
        },
        watch: {
            'filter.business' (business) {
                if (this.businessResolver) {
                    this.businessResolver()
                } else {
                    this.getHostList()
                }
            },
            customColumns () {
                this.setTableHeader()
            }
        },
        async created () {
            try {
                await this.getParams()
                await Promise.all([
                    this.getProperties(),
                    this.getHostPropertyGroups()
                ])
                await this.setTableHeader()
                this.getHostList()
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('hostSearch', ['searchHost']),
            ...mapActions('hostDelete', ['deleteHost']),
            ...mapActions('hostRelation', ['transferResourcehostToIdlemodule']),
            getParams () {
                return new Promise((resolve, reject) => {
                    this.filter.paramsResolver = () => {
                        this.filter.paramsResolver = null
                        resolve()
                    }
                })
            },
            setTableHeader () {
                return new Promise((resolve, reject) => {
                    const headerProperties = this.$tools.getHeaderProperties(this.properties.host, this.customColumns)
                    resolve(headerProperties)
                }).then(properties => {
                    this.table.header = [{
                        id: 'bk_host_id',
                        type: 'checkbox',
                        objId: 'host',
                        width: 50
                    }].concat(properties.map(property => {
                        return {
                            id: property['bk_property_id'],
                            name: property['bk_property_name'],
                            objId: property['bk_obj_id']
                        }
                    }))
                    this.columnsConfig.selected = properties.map(property => property['bk_property_id'])
                })
            },
            getHostList () {
                this.searchHost({
                    params: {
                        ...this.getScopedParams(),
                        'bk_biz_id': this.filter.business,
                        page: {
                            start: (this.table.pagination.current - 1) * this.table.pagination.size,
                            limit: this.table.pagination.size,
                            sort: this.table.sort
                        }
                    },
                    config: {
                        requestId: 'hostSearch'
                    }
                }).then(data => {
                    this.table.pagination.count = data.count
                    this.table.list = data.info
                    this.setAllHostList(data.info)
                    return data
                }).catch(e => {
                    this.table.checked = []
                    this.table.list = []
                    this.table.pagination.count = 0
                })
            },
            getAllHostList () {
                return this.searchHost({
                    params: {
                        ...this.getScopedParams(),
                        'bk_biz_id': this.filter.business,
                        page: {}
                    },
                    config: {
                        requestId: 'hostSearchAll'
                    }
                }).then(data => {
                    this.table.allList = data.info
                    return data
                })
            },
            getScopedParams () {
                const params = this.$tools.clone(this.filter.params)
                if (!this.filter.assigned) {
                    const businessParams = params.condition.find(condition => condition['bk_obj_id'] === 'biz')
                    businessParams.condition.push({
                        field: 'default',
                        operator: '$eq',
                        value: 1
                    })
                }
                return params
            },
            routeToHistory () {
                this.$router.push('/history/host?relative=/resource')
            },
            handleAssignHosts (businessId, business) {
                if (!businessId) return
                if (this.hasSelectAssignedHost()) {
                    this.$error(this.$t('Hosts["请勿选择已分配主机"]'))
                    this.assignBusiness = ''
                } else {
                    this.$bkInfo({
                        title: this.$t("HostResourcePool['请确认是否转移']"),
                        content: this.getConfirmContent(business),
                        confirmFn: () => {
                            this.assignHosts(business)
                        },
                        cancelFn: () => {
                            this.assignBusiness = ''
                        }
                    })
                }
            },
            hasSelectAssignedHost () {
                const list = this.table.allList.filter(item => this.table.checked.includes(item['host']['bk_host_id']))
                const existAssigned = list.some(item => item['biz'].some(biz => biz.default !== 1))
                return existAssigned
            },
            assignHosts (business) {
                this.transferResourcehostToIdlemodule({
                    params: {
                        'bk_biz_id': business['bk_biz_id'],
                        'bk_host_id': this.table.checked
                    }
                }).then(() => {
                    this.$success(this.$t("HostResourcePool['分配成功']"))
                    this.table.checked = []
                    this.assignBusiness = ''
                    this.handlePageChange(1)
                })
            },
            getConfirmContent (business) {
                const render = this.$createElement
                let content
                if (this.$i18n.locale === 'en') {
                    content = render('p', [
                        render('span', 'Selected '),
                        render('span', {
                            style: {color: '#3c96ff'}
                        }, this.table.checked.length),
                        render('span', ' Hosts Transfer to Idle machine under '),
                        render('span', {
                            style: {color: '#3c96ff'}
                        }, business['bk_biz_name'])
                    ])
                } else {
                    content = render('p', [
                        render('span', '选中的 '),
                        render('span', {
                            style: {color: '#3c96ff'}
                        }, this.table.checked.length),
                        render('span', ' 个主机转移到 '),
                        render('span', {
                            style: {color: '#3c96ff'}
                        }, business['bk_biz_name']),
                        render('span', ' 下的空闲机模块')
                    ])
                }
                return content
            },
            handleMultipleDelete () {
                this.$bkInfo({
                    title: `${this.$t("HostResourcePool['确定删除选中的主机']")}？`,
                    confirmFn: () => {
                        this.deleteHost({
                            params: {
                                data: {
                                    'bk_host_id': this.table.checked.join(','),
                                    'bk_supplier_account': this.supplierAccount
                                }
                            }
                        }).then(() => {
                            this.$bkInfo({
                                statusOpts: {
                                    title: this.$t("HostResourcePool['成功删除选中的主机']"),
                                    subtitle: false
                                },
                                type: 'success'
                            })
                            this.table.checked = []
                            this.handlePageChange(1)
                        })
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .resource-layout{
        height: 100%;
        padding: 0;
        overflow: hidden;
        .resource-main{
            height: 100%;
            padding: 20px;
            overflow: hidden;
        }
        .resource-filter{
            height: 100%;
            .filter-field{
                margin: 10px 15px 0 0;
            }
            .filter-field-label{
                display: inline-block;
                vertical-align: middle;
                line-height: 1;
                font-size: 14px;
            }
        }
    }
    .resource-options{
        font-size: 0;
        .options-button{
            position: relative;
            display: inline-block;
            vertical-align: middle;
            border-radius: 0;
            font-size: 14px;
            margin-left: -1px;
            &:hover{
                z-index: 1;
            }
            &-delete:hover{
                border-color: #ef4c4c;
                color: #ef4c4c;
            }
            &-history{
                width: 36px;
                padding: 0;
                margin: 0 0 0 10px;
                border-radius: 2px;
            }
        }
        .options-table-selector,
        .options-business-selector{
            margin: 0 10px 0 0;
        }
        .options-business-selector{
            width: 180px;
        }
    }
    .resource-table{
        margin-top: 20px;
    }
</style>