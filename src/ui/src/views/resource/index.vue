<template>
    <div class="resource-layout clearfix">
        <bk-tab
            :active.sync="activeTab"
            class="scope-tab"
            type="unborder-card"
            @tab-change="handleTabChange">
            <bk-tab-panel v-for="item in scopeList"
                :key="item.id"
                :name="item.id"
                :label="item.label">
            </bk-tab-panel>
        </bk-tab>
        <div class="content">
            <cmdb-resize-layout
                v-if="activeTab === 1"
                :class="['resize-layout fl', { 'is-collapse': layout.collapse }]"
                :handler-offset="3"
                :min="200"
                :max="480"
                :disabled="layout.collapse"
                direction="right">
                <resource-diractory></resource-diractory>
                <i class="diractory-collapse-icon bk-icon icon-angle-left"
                    @click="layout.collapse = !layout.collapse">
                </i>
            </cmdb-resize-layout>
            <resource-hosts class="main" ref="resourceHosts"></resource-hosts>
        </div>
        <router-subview></router-subview>
    </div>
</template>

<script>
    import resourceDiractory from './children/directory.vue'
    import resourceHosts from './children/host-list.vue'
    import Bus from '@/utils/bus.js'
    import { mapGetters, mapActions } from 'vuex'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            resourceDiractory,
            resourceHosts
        },
        data () {
            return {
                layout: {
                    collapse: false
                },
                activeTab: this.$route.query.tab || 1,
                scopeList: [{
                    id: 1,
                    label: this.$t('未分配')
                }, {
                    id: 0,
                    label: this.$t('已分配')
                }, {
                    id: 'all',
                    label: this.$t('全部')
                }],
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
                businessList: [],
                assignDialog: {
                    show: false,
                    business: {}
                }
            }
        },
        computed: {
            ...mapGetters(['userName', 'isAdminView']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
            buttons () {
                return [{
                    id: 'edit',
                    text: this.$t('修改'),
                    handler: this.handleMultipleEdit,
                    disabled: !this.table.checked.length,
                    auth: this.$authResources({ type: this.$OPERATION.U_RESOURCE_HOST })
                }]
            },
            columnsConfigKey () {
                // 资源池独占components/table，无需再判断，否则会引发资源池详情第一次点击无法跳转至业务拓扑
                // 因为key的变化引发了header的变化，从而导致触发RouterQuery中的重定向，使第一次跳转失效
                // return `${this.userName}_$resource_${this.isAdminView ? 'adminView' : this.bizId}_table_columns`
                return `${this.userName}_$resource_adminView_table_columns`
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
            activeTab (tab) {
                const resourceHosts = this.$refs.resourceHosts
                if (resourceHosts) {
                    resourceHosts.scope = tab
                    resourceHosts.handlePageChange(1)
                }
            }
        },
        async created () {
            try {
                await Promise.all([
                    this.getFullAmountBusiness(),
                    this.getProperties()
                ])
                this.ready = true
            } catch (e) {
                console.error(e)
            }
        },
        methods: {
            ...mapActions('hostBatch', ['exportHost']),
            ...mapActions('hostSearch', ['searchHost']),
            ...mapActions('hostDelete', ['deleteHost']),
            ...mapActions('hostRelation', ['transferResourcehostToIdleModule']),
            ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
            handleTabChange (tab) {
                Bus.$emit('toggle-host-filter', false)
                Bus.$emit('reset-host-filter')
            },
            async getFullAmountBusiness () {
                try {
                    const data = await this.$http.get('biz/simplify?sort=bk_biz_name')
                    this.businessList = data.info || []
                } catch (e) {
                    console.error(e)
                    this.businessList = []
                }
            },
            getProperties () {
                return this.batchSearchObjectAttribute({
                    injectId: 'host',
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
                    this.assignDialog.business = business
                    this.assignDialog.show = true
                }
            },
            hasSelectAssignedHost () {
                const allList = this.$refs.resourceTable.table.list
                const list = allList.filter(item => this.table.checked.includes(item['host']['bk_host_id']))
                const existAssigned = list.some(item => item['biz'].some(biz => biz.default !== 1))
                return existAssigned
            },
            cancelAssignHosts () {
                this.assignBusiness = 'empty'
                this.assignDialog.show = false
            },
            assignHosts () {
                this.transferResourcehostToIdleModule({
                    params: {
                        'bk_biz_id': this.assignDialog.business['bk_biz_id'],
                        'bk_host_id': this.table.checked
                    },
                    config: {
                        requestId: 'transferResourcehostToIdleModule'
                    }
                }).then(() => {
                    this.$success(this.$t('分配成功'))
                    this.$refs.resourceTable.table.checked = []
                    RouterQuery.set({
                        _t: Date.now(),
                        page: 1
                    })
                }).finally(() => {
                    this.assignBusiness = 'empty'
                    this.assignDialog.show = false
                })
            },
            refreshList () {
                RouterQuery.set({
                    _t: Date.now(),
                    ip: ''
                })
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
                            RouterQuery.set({
                                _t: Date.now(),
                                page: 1
                            })
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
        .scope-tab {
            height: auto;
            margin: 0 20px;
            /deep/ .bk-tab-header {
                padding: 0;
            }
        }
        .content {
            height: calc(100% - 58px);
            overflow: hidden;
        }
        .resize-layout {
            position: relative;
            width: 280px;
            height: 100%;
            border-right: 1px solid $cmdbLayoutBorderColor;
            &.is-collapse {
                width: 0 !important;
                border-right: none;
                .diractory-collapse-icon:before {
                    display: inline-block;
                    transform: rotate(180deg);
                }
            }
            .diractory-collapse-icon {
                position: absolute;
                left: 100%;
                top: 50%;
                width: 16px;
                height: 100px;
                line-height: 100px;
                background: $cmdbLayoutBorderColor;
                border-radius: 0px 12px 12px 0px;
                transform: translateY(-50%);
                text-align: center;
                font-size: 12px;
                color: #fff;
                cursor: pointer;
                &:hover {
                    background: #699DF4;
                }
            }
        }
        .main {
            height: 100%;
            padding: 10px 20px 0 20px;
            overflow: hidden;
        }
    }
    .assign-dialog {
        /deep/ .bk-dialog-body {
            padding: 0 50px 40px;
        }
        .assign-info span {
            color: #3c96ff;
        }
        .assign-footer {
            padding-top: 20px;
            font-size: 0;
            text-align: center;
            .bk-button-normal {
                width: 96px;
            }
        }
    }
</style>
