<template>
    <div class="hosts-layout clearfix">
        <cmdb-hosts-filter class="hosts-filter fr" @on-refresh="handleRefresh" :filter-config-key="filter.filterConfigKey">
            <div class="filter-group" slot="business">
                <label class="filter-label">{{$t('Hosts[\'选择业务\']')}}</label>
                <cmdb-business-selector class="filter-field" v-model="filter.business"></cmdb-business-selector>
            </div>
        </cmdb-hosts-filter>
        <div class="hosts-main">
            <div class="hosts-options">
                <bk-button class="options-button" type="default"
                    v-tooltip="$t('BusinessTopology[\'修改\']')"
                    :disabled="!table.checked.length"
                    @click="handleMultipleEdit">
                    <i class="icon-cc-edit"></i>
                </bk-button>
                <bk-button class="options-button" type="default"
                    v-tooltip="$t('BusinessTopology[\'转移\']')"
                    :disabled="!table.checked.length">
                    <i class="icon-cc-shift"></i>
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
                <bk-button class="options-button" type="default"
                    v-tooltip="$t('BusinessTopology[\'列表显示属性配置\']')"
                    @click="columnsConfig.show = true">
                    <i class="icon-cc-setting"></i>
                </bk-button>
            </div>
            <cmdb-table class="hosts-table"
                :loading="$loading()"
                :checked.sync="table.checked"
                :header="table.header"
                :list="table.list"
                :defaultSort="table.defaultSort"
                :pagination.sync="table.pagination"
                :wrapperMinusHeight="157"
                @handleRowClick="handleRowClick"
                @handleSortChange="handleSortChange"
                @handlePageChange="handlePageChange"
                @handleSizeChange="handleSizeChange"
                @handleCheckAll="handleCheckAll">
                <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                    <label class="table-checkbox bk-form-checkbox bk-checkbox-small"
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
                :properties="columnsConfigProperties"
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
                    columnsConfigKey: 'hosts_table_columns'
                },
                filter: {
                    filterConfigKey: 'hosts_filter_fields',
                    business: null,
                    businessResolver: null
                }
            }
        },
        computed: {
            columnsConfigProperties () {
                const setProperties = this.properties.set.filter(property => ['bk_set_name'].includes(property['bk_property_id']))
                const moduleProperties = this.properties.module.filter(property => ['bk_module_name'].includes(property['bk_property_id']))
                const hostProperties = this.properties.host
                return [...setProperties, ...moduleProperties, ...hostProperties]
            }
        },
        watch: {
            'filter.business' (business) {
                if (this.filter.businessResolver) {
                    this.filter.businessResolver()
                } else {
                    this.table.checked = []
                    this.getHostList()
                }
            },
            customColumns () {
                this.setTableHeader()
            }
        },
        async created () {
            try {
                await Promise.all([
                    this.getBusiness(),
                    this.getParams()
                ])
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
            getBusiness () {
                return new Promise((resolve, reject) => {
                    this.filter.businessResolver = () => {
                        this.filter.businessResolver = null
                        resolve()
                    }
                })
            },
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
                    const headerProperties = this.$tools.getHeaderProperties(this.columnsConfigProperties, this.customColumns)
                    resolve(headerProperties)
                }).then(properties => {
                    this.table.header = [{
                        id: 'bk_host_id',
                        type: 'checkbox',
                        objId: 'host',
                        width: 200
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
                        ...this.filter.params,
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
                        ...this.filter.params,
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .hosts-layout{
        height: 100%;
        padding: 0;
        overflow: hidden;
        .hosts-main{
            height: 100%;
            padding: 20px;
            overflow: hidden;
        }
        .hosts-filter{
            height: 100%;
        }
    }
    .hosts-options{
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
        }
    }
    .hosts-table{
        margin-top: 20px;
    }
</style>