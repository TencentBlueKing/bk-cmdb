<template>
    <div class="resource-layout clearfix">
        <cmdb-hosts-table class="resource-main" ref="resourceTable"
            :columns-config-key="columnsConfigKey"
            :columns-config-properties="columnsConfigProperties"
            :columns-config-disabled-columns="['bk_host_innerip', 'bk_cloud_id', 'bk_biz_name', 'bk_module_name']"
            :edit-auth="$OPERATION.U_RESOURCE_HOST"
            :delete-auth="$OPERATION.D_RESOURCE_HOST"
            :save-auth="$OPERATION.U_RESOURCE_HOST"
            :show-scope="true"
            :show-history="true"
            @on-checked="handleChecked"
            @on-set-header="handleSetHeader">
            <template slot="options-left">
                <cmdb-auth v-if="isAdminView" :auth="$authResources({ type: $OPERATION.C_RESOURCE_HOST })">
                    <bk-button slot-scope="{ disabled }"
                        class="options-button"
                        theme="primary"
                        style="margin-left: 0"
                        :disabled="disabled"
                        @click="importInst.show = true">
                        {{$t('导入主机')}}
                    </bk-button>
                </cmdb-auth>
                <bk-select class="options-business-selector"
                    v-if="isAdminView"
                    font-size="medium"
                    :popover-width="180"
                    :searchable="businessList.length > 7"
                    :disabled="!table.checked.length"
                    :clearable="false"
                    :placeholder="$t('分配到')"
                    v-model="assignBusiness"
                    @selected="handleAssignHosts">
                    <bk-option id="empty" :name="$t('分配到')" hidden></bk-option>
                    <bk-option v-for="option in businessList"
                        :key="option.bk_biz_id"
                        :id="option.bk_biz_id"
                        :name="option.bk_biz_name">
                    </bk-option>
                </bk-select>
                <cmdb-clipboard-selector class="options-clipboard"
                    :list="clipboardList"
                    :disabled="!table.checked.length"
                    @on-copy="handleCopy">
                </cmdb-clipboard-selector>
                <cmdb-button-group
                    :buttons="buttons"
                    :expand="!isAdminView">
                </cmdb-button-group>
            </template>
        </cmdb-hosts-table>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="importInst.show"
            :width="800"
            :title="$t('批量导入')">
            <bk-tab :active.sync="importInst.active" type="unborder-card" slot="content" v-if="importInst.show">
                <bk-tab-panel name="import" :label="$t('批量导入')">
                    <cmdb-import v-if="importInst.show && importInst.active === 'import'"
                        :template-url="importInst.templateUrl"
                        :import-url="importInst.importUrl"
                        @success="getHostList(true)"
                        @partialSuccess="getHostList(true)">
                        <span slot="download-desc" style="display: inline-block;vertical-align: top;">
                            {{$t('说明：内网IP为必填列')}}
                        </span>
                    </cmdb-import>
                </bk-tab-panel>
                <bk-tab-panel name="agent" :label="$t('自动导入')">
                    <div class="automatic-import">
                        <p>{{$t("agent安装说明")}}</p>
                        <div class="back-contain">
                            <i class="icon-cc-skip"></i>
                            <a href="javascript:void(0)" @click="openAgentApp">{{$t('点此进入节点管理')}}</a>
                        </div>
                    </div>
                </bk-tab-panel>
            </bk-tab>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapGetters, mapActions, mapState } from 'vuex'
    import cmdbHostsTable from '@/components/hosts/table'
    import cmdbImport from '@/components/import/import'
    import cmdbButtonGroup from '@/components/ui/other/button-group'
    import { MENU_RESOURCE_MANAGEMENT } from '@/dictionary/menu-symbol'
    export default {
        components: {
            cmdbHostsTable,
            cmdbImport,
            cmdbButtonGroup
        },
        data () {
            return {
                properties: {
                    biz: [],
                    host: [],
                    set: [],
                    module: []
                },
                table: {
                    checked: [],
                    header: [],
                    columnsConfigKey: 'resource_table_columns',
                    exportUrl: `${window.API_HOST}hosts/export`
                },
                assignBusiness: 'empty',
                importInst: {
                    show: false,
                    active: 'import',
                    templateUrl: `${window.API_HOST}importtemplate/host`,
                    importUrl: `${window.API_HOST}hosts/import`
                },
                isDropdownShow: false,
                ready: false,
                businessList: []
            }
        },
        computed: {
            ...mapGetters(['userName', 'isAdminView']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapState('hosts', ['filterParams']),
            buttons () {
                return [{
                    id: 'edit',
                    text: this.$t('修改'),
                    handler: this.handleMultipleEdit,
                    disabled: !this.table.checked.length,
                    auth: this.$authResources({ type: this.$OPERATION.U_RESOURCE_HOST })
                }, {
                    id: 'delete',
                    text: this.$t('删除'),
                    handler: this.handleMultipleDelete,
                    disabled: !this.table.checked.length,
                    auth: this.$authResources({ type: this.$OPERATION.D_RESOURCE_HOST }),
                    available: this.isAdminView
                }, {
                    id: 'export',
                    text: this.$t('导出'),
                    handler: this.exportField,
                    disabled: !this.table.checked.length
                }]
            },
            columnsConfigKey () {
                return `${this.userName}_$resource_${this.isAdminView ? 'adminView' : this.bizId}_table_columns`
            },
            customColumns () {
                return this.usercustom[this.columnsConfigKey]
            },
            clipboardList () {
                return this.table.header.filter(header => header.type !== 'checkbox')
            },
            columnsConfigProperties () {
                const setProperties = this.properties.set.filter(property => ['bk_set_name'].includes(property['bk_property_id']))
                const moduleProperties = this.properties.module.filter(property => ['bk_module_name'].includes(property['bk_property_id']))
                const businessProperties = this.properties.biz.filter(property => ['bk_biz_name'].includes(property['bk_property_id']))
                const hostProperties = this.properties.host
                return [...setProperties, ...moduleProperties, ...businessProperties, ...hostProperties]
            }
        },
        watch: {
            'importInst.show' (show) {
                if (!show) {
                    this.importInst.active = 'import'
                }
            },
            filterParams () {
                if (this.ready) {
                    this.getHostList(false, true)
                }
            }
        },
        async created () {
            try {
                this.setDynamicBreadcrumbs()
                await this.getFullAmountBusiness()
                await this.getProperties()
                this.getHostList()
                this.ready = true
            } catch (e) {
                console.error(e)
            }
        },
        beforeDestroy () {
            this.ready = false
        },
        methods: {
            ...mapActions('hostBatch', ['exportHost']),
            ...mapActions('hostSearch', ['searchHost']),
            ...mapActions('hostDelete', ['deleteHost']),
            ...mapActions('hostRelation', ['transferResourcehostToIdleModule']),
            ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
            setDynamicBreadcrumbs () {
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('资源目录'),
                    route: {
                        name: MENU_RESOURCE_MANAGEMENT
                    }
                }, {
                    label: this.$t('主机')
                }])
            },
            async getFullAmountBusiness () {
                try {
                    const data = await this.$store.dispatch('objectBiz/getFullAmountBusiness')
                    this.businessList = data.info || []
                } catch (e) {
                    console.error(e)
                    this.businessList = []
                }
            },
            getProperties () {
                return this.batchSearchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: { '$in': Object.keys(this.properties) },
                        bk_supplier_account: this.supplierAccount
                    }, { inject: false }),
                    config: {
                        requestId: `post_batchSearchObjectAttribute_${Object.keys(this.properties).join('_')}`,
                        requestGroup: Object.keys(this.properties).map(id => `post_searchObjectAttribute_${id}`)
                    }
                }).then(result => {
                    Object.keys(this.properties).forEach(objId => {
                        this.properties[objId] = result[objId]
                    })
                    return result
                })
            },
            getHostList (resetPage = false, event = false) {
                const params = this.getParams()
                this.$refs.resourceTable.search(-1, params, resetPage, event)
            },
            getParams () {
                const defaultModel = ['biz', 'set', 'module', 'host']
                const params = {
                    bk_biz_id: -1,
                    ip: this.filterParams.ip,
                    condition: defaultModel.map(model => {
                        return {
                            bk_obj_id: model,
                            condition: this.filterParams[model] || [],
                            fields: []
                        }
                    })
                }
                return params
            },
            handleAssignHosts (businessId, option) {
                const business = {
                    bk_biz_id: businessId,
                    bk_biz_name: option.name
                }
                if (!businessId) return
                if (this.hasSelectAssignedHost()) {
                    this.$error(this.$t('请勿选择已分配主机'))
                    this.$nextTick(() => {
                        this.assignBusiness = 'empty'
                    })
                } else {
                    this.$bkInfo({
                        title: this.$t('确认分配到业务'),
                        subHeader: this.getConfirmContent(business),
                        confirmFn: () => {
                            this.assignHosts(business)
                        },
                        cancelFn: () => {
                            this.assignBusiness = 'empty'
                        }
                    })
                }
            },
            hasSelectAssignedHost () {
                const allList = this.$refs.resourceTable.table.list
                const list = allList.filter(item => this.table.checked.includes(item['host']['bk_host_id']))
                const existAssigned = list.some(item => item['biz'].some(biz => biz.default !== 1))
                return existAssigned
            },
            assignHosts (business) {
                this.transferResourcehostToIdleModule({
                    params: {
                        'bk_biz_id': business['bk_biz_id'],
                        'bk_host_id': this.table.checked
                    }
                }).then(() => {
                    this.$success(this.$t('分配成功'))
                    this.assignBusiness = 'empty'
                    this.$refs.resourceTable.table.checked = []
                    this.$refs.resourceTable.handlePageChange(1)
                }).catch(e => {
                    this.assignBusiness = 'empty'
                })
            },
            getConfirmContent (business) {
                const render = this.$createElement
                const content = render('i18n', {
                    props: {
                        path: '确认转移主机信息'
                    }
                }, [
                    render('span', {
                        attrs: { place: 'num' },
                        style: { color: '#3c96ff' }
                    }, this.table.checked.length),
                    render('span', {
                        attrs: { place: 'name' },
                        style: { color: '#3c96ff' }
                    }, business['bk_biz_name'])
                ])
                return content
            },
            handleChecked (checked) {
                this.table.checked = checked
            },
            handleSetHeader (header) {
                this.table.header = header
            },
            handleCopy (target) {
                this.$refs.resourceTable.handleCopy(target)
            },
            handleMultipleEdit () {
                if (this.hasSelectAssignedHost()) {
                    this.$error(this.$t('请勿选择已分配主机'))
                    return false
                }
                this.$refs.resourceTable.handleMultipleEdit()
            },
            handleMultipleDelete () {
                if (this.hasSelectAssignedHost()) {
                    this.$error(this.$t('请勿选择已分配主机'))
                    return false
                }
                this.$bkInfo({
                    title: `${this.$t('确定删除选中的主机')}？`,
                    confirmFn: () => {
                        this.deleteHost({
                            params: {
                                data: {
                                    'bk_host_id': this.table.checked.join(','),
                                    'bk_supplier_account': this.supplierAccount
                                }
                            }
                        }).then(() => {
                            this.$success(this.$t('成功删除选中的主机'))
                            this.$refs.resourceTable.table.checked = []
                            this.$refs.resourceTable.handlePageChange(1)
                        })
                    }
                })
            },
            openAgentApp () {
                const agent = window.Site.agent
                if (agent) {
                    const topWindow = window.top
                    const isPaasConsole = topWindow !== window
                    if (isPaasConsole) {
                        topWindow.postMessage(JSON.stringify({
                            action: 'open_other_app',
                            app_code: 'bk_nodeman'
                        }), '*')
                    } else {
                        window.open(agent)
                    }
                } else {
                    this.$warn(this.$t('未配置Agent安装APP地址'))
                }
            },
            exportExcel (response) {
                const contentDisposition = response.headers['content-disposition']
                const fileName = contentDisposition.substring(contentDisposition.indexOf('filename') + 9)
                const url = window.URL.createObjectURL(new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }))
                const link = document.createElement('a')
                link.style.display = 'none'
                link.href = url
                link.setAttribute('download', fileName)
                document.body.appendChild(link)
                link.click()
                document.body.removeChild(link)
            },
            async exportField () {
                const formData = new FormData()
                formData.append('bk_host_id', this.table.checked)
                if (this.customColumns) {
                    formData.append('export_custom_fields', this.customColumns)
                }
                formData.append('bk_biz_id', '-1')
                const res = await this.exportHost({
                    params: formData,
                    config: {
                        globalError: false,
                        originalResponse: true,
                        responseType: 'blob'
                    }
                })
                this.exportExcel(res)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .resource-layout{
        padding: 0;
        overflow: hidden;
        .resource-main{
            height: 100%;
            padding: 0 20px;
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
    .options-button{
        position: relative;
        display: inline-block;
        vertical-align: middle;
        font-size: 14px;
        margin: 0 10px 0 0;
        padding: 0 10px;
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
        &.options-icon {
            border-radius: 0;
            margin: 0 -1px 0 0;
        }
    }
    .options-clipboard {
        margin: 0 10px 0 0px;
    }
    .options-business-selector{
        display: inline-block;
        vertical-align: middle;
        margin: 0 10px 0 0;
    }
    .automatic-import{
        padding:40px 30px 0 30px;
        .back-contain{
            cursor:pointer;
            color: #3c96ff;
            img{
                margin-right: 5px;
            }
            a{
                color:#3c96ff;
            }
        }
    }
</style>
