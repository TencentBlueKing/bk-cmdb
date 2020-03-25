<template>
    <div class="options-layout clearfix">
        <div class="options-left">
            <template v-if="scope === 1">
                <cmdb-auth class="mr10" :auth="$authResources({ type: $OPERATION.C_RESOURCE_HOST })">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        style="margin-left: 0"
                        :disabled="disabled"
                        @click="importInst.show = true">
                        {{$t('导入主机')}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth class="mr10" :auth="$authResources({ type: $OPERATION.U_RESOURCE_HOST })">
                    <bk-select slot-scope="{ disabled }"
                        font-size="medium"
                        :popover-width="180"
                        :disabled="!table.checked.length || disabled"
                        :clearable="false"
                        :placeholder="$t('分配到')"
                        v-model="assign.curSelected"
                        @selected="handleAssignHosts">
                        <bk-option id="empty" :name="$t('分配到')" hidden></bk-option>
                        <bk-option v-for="option in assignTarget"
                            :key="option.id"
                            :id="option.id"
                            :name="option.name">
                        </bk-option>
                    </bk-select>
                </cmdb-auth>
            </template>
            <cmdb-auth v-else class="mr10" :auth="$authResources({ type: $OPERATION.U_RESOURCE_HOST })">
                <bk-button slot-scope="{ disabled }"
                    theme="primary"
                    style="margin-left: 0"
                    :disabled="disabled || !table.checked.length"
                    @click="handleMultipleEdit">
                    {{$t('编辑')}}
                </bk-button>
            </cmdb-auth>
            <cmdb-clipboard-selector class="options-clipboard mr10"
                :list="clipboardList"
                :disabled="!table.checked.length"
                @on-copy="handleCopy">
            </cmdb-clipboard-selector>
            <cmdb-button-group v-if="scope === 1"
                class="mr10"
                :buttons="buttons"
                :expand="false">
            </cmdb-button-group>
            <bk-button v-else theme="default" :disabled="!table.checked.length" @click="exportField">{{$t('导出')}}</bk-button>
        </div>
        <div class="options-right">
            <cmdb-host-filter class="ml10"
                ref="hostFilter"
                :section-height="$APP.height - 250"
                :properties="filterProperties"
                :show-scope="true">
            </cmdb-host-filter>
            <icon-button class="ml10"
                icon="icon icon-cc-setting"
                v-bk-tooltips.top="$t('列表显示属性配置')"
                @click="handleColumnConfigClick">
            </icon-button>
            <icon-button class="ml10"
                v-bk-tooltips="$t('查看删除历史')"
                icon="icon icon-cc-history"
                @click="routeToHistory">
            </icon-button>
        </div>
        
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
                        @success="$parent.getHostList(true)"
                        @partialSuccess="$parent.getHostList(true)">
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
        
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="slider.show"
            :title="slider.title"
            :width="800"
            :before-close="handleSliderBeforeClose">
            <bk-tab :active.sync="tab.active" type="unborder-card" slot="content" v-if="slider.show">
                <bk-tab-panel name="attribute" :label="$t('属性')" style="width: calc(100% + 40px);margin: 0 -20px;">
                    <cmdb-form-multiple v-if="tab.attribute.type === 'multiple'"
                        ref="multipleForm"
                        :properties="properties.host"
                        :property-groups="propertyGroups"
                        :object-unique="objectUnique"
                        @on-submit="handleMultipleSave"
                        @on-cancel="handleSliderBeforeClose">
                    </cmdb-form-multiple>
                </bk-tab-panel>
            </bk-tab>
        </bk-sideslider>

        <bk-sideslider
            v-transfer-dom
            :is-show.sync="columnsConfig.show"
            :width="600"
            :title="$t('列表显示属性配置')">
            <cmdb-columns-config slot="content"
                v-if="columnsConfig.show"
                :properties="columnsConfigProperties"
                :selected="columnsConfigSelected"
                :disabled-columns="columnsConfig.disabledColumns"
                @on-apply="handleApplyColumnsConfig"
                @on-cancel="columnsConfig.show = false"
                @on-reset="handleResetColumnsConfig">
            </cmdb-columns-config>
        </bk-sideslider>

        <bk-dialog
            class="assign-dialog"
            v-model="assign.show"
            header-position="left"
            :width="480"
            :mask-close="false"
            :esc-close="false"
            :close-icon="false"
            :title="assign.title"
            @cancel="handleCancelAssignHosts">
            <div class="assign-content">
                <i18n class="assign-count" tag="div" path="已选择主机">
                    <span place="count">{{table.checked.length}}</span>
                </i18n>
                <div class="assign-seleted">
                    <p>{{assign.label}}</p>
                    <bk-select
                        font-size="normal"
                        :searchable="true"
                        :clearable="false"
                        :disabled="$loading(assign.requestId)"
                        :placeholder="assign.placeholder"
                        v-model="assign.id">
                        <bk-option v-for="option in assignOptions"
                            :key="option.id"
                            :id="option.id"
                            :name="option.name">
                        </bk-option>
                        <div slot="extension" v-if="assign.curSelected === 'toDirs'" @click="handleApplyPermission">
                            <a href="javascript:void(0)" class="apply-others">
                                <i class="bk-icon icon-plus-circle"></i>
                                {{$t('申请其他资源目录')}}
                            </a>
                        </div>
                    </bk-select>
                </div>
            </div>
            <div class="assign-footer" slot="footer">
                <bk-button
                    class="mr10"
                    theme="primary"
                    :disabled="assign.id === ''"
                    :loading="$loading(assign.requestId)"
                    @click="handleConfirmAssign">
                    {{$t('确定')}}
                </bk-button>
                <bk-button theme="default" :disabled="$loading(assign.requestId)" @click="handleCancelAssignHosts">{{$t('取消')}}</bk-button>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { translateAuth } from '@/setup/permission'
    import cmdbImport from '@/components/import/import'
    import cmdbButtonGroup from '@/components/ui/other/button-group'
    import cmdbHostFilter from '@/components/hosts/filter/index.vue'
    import cmdbColumnsConfig from '@/components/columns-config/columns-config'
    import Bus from '@/utils/bus.js'
    export default {
        components: {
            cmdbImport,
            cmdbButtonGroup,
            cmdbHostFilter,
            cmdbColumnsConfig
        },
        data () {
            return {
                importInst: {
                    show: false,
                    active: 'import',
                    templateUrl: `${window.API_HOST}importtemplate/host`,
                    importUrl: `${window.API_HOST}hosts/import`
                },
                businessList: [],
                objectUnique: [],
                slider: {
                    show: false,
                    title: ''
                },
                tab: {
                    active: 'attribute',
                    attribute: {
                        type: 'details',
                        inst: {
                            details: {},
                            edit: {},
                            original: {}
                        }
                    }
                },
                columnsConfig: {
                    show: false,
                    selected: [],
                    disabledColumns: ['bk_host_innerip', 'bk_cloud_id', 'bk_module_name', 'bk_set_name']
                },
                assign: {
                    show: false,
                    id: '',
                    curSelected: 'empty',
                    placeholder: this.$t('请选择xx', { name: this.$t('业务') }),
                    label: this.$t('业务列表'),
                    title: this.$t('分配到业务空闲机'),
                    requestId: Symbol('assignHosts')
                },
                assignTarget: [{
                    id: 'toBusiness',
                    name: this.$t('业务空闲机')
                }, {
                    id: 'toDirs',
                    name: this.$t('资源池其他目录')
                }],
                dirList: []
            }
        },
        computed: {
            ...mapGetters('resourceHost', ['activeDirectory']),
            table () {
                return this.$parent.table
            },
            clipboardList () {
                return this.table.header.filter(header => header.type !== 'checkbox')
            },
            scope () {
                return this.$parent.scope
            },
            properties () {
                return this.$parent.properties
            },
            propertyGroups () {
                return this.$parent.propertyGroups
            },
            columnsConfigProperties () {
                return this.$parent.columnsConfigProperties
            },
            columnsConfigSelected () {
                return this.$parent.columnsConfig.selected
            },
            buttons () {
                return [{
                    id: 'edit',
                    text: this.$t('编辑'),
                    handler: this.handleMultipleEdit,
                    disabled: !this.table.checked.length,
                    auth: this.$authResources({ type: this.$OPERATION.U_RESOURCE_HOST })
                }, {
                    id: 'delete',
                    text: this.$t('删除'),
                    handler: this.handleMultipleDelete,
                    disabled: !this.table.checked.length,
                    auth: this.$authResources({ type: this.$OPERATION.D_RESOURCE_HOST })
                }, {
                    id: 'export',
                    text: this.$t('导出'),
                    handler: this.exportField,
                    disabled: !this.table.checked.length
                }]
            },
            filterProperties () {
                const { module, set, host } = this.properties
                const filterProperty = ['bk_host_innerip', 'bk_host_outerip']
                return {
                    host: host.filter(property => !filterProperty.includes(property.bk_property_id)),
                    module,
                    set
                }
            },
            assignOptions () {
                if (this.assign.curSelected === 'toBusiness') {
                    return this.businessList.map(item => ({
                        id: item.bk_biz_id,
                        name: item.bk_biz_name
                    }))
                }
                return this.dirList.filter(item => item.bk_module_id !== this.activeDirectory.bk_module_id).map(item => ({
                    id: item.bk_module_id,
                    name: item.bk_module_name
                }))
            }
        },
        watch: {
            'importInst.show' (show) {
                if (!show) {
                    this.importInst.active = 'import'
                }
            }
        },
        async created () {
            try {
                await this.getFullAmountBusiness()
                await this.getDirectoryList()
            } catch (e) {
                console.error(e.message)
            }
        },
        methods: {
            async getFullAmountBusiness () {
                try {
                    const data = await this.$http.get('biz/simplify?sort=bk_biz_name')
                    this.businessList = data.info || []
                } catch (e) {
                    console.error(e)
                    this.businessList = []
                }
            },
            async getDirectoryList () {
                try {
                    const data = await this.$store.dispatch('resourceDirectory/getDirectoryList', {
                        params: {}
                    })
                    this.dirList = data.info || []
                } catch (e) {
                    console.error(e)
                    this.dirList = []
                }
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
            handleAssignHosts (id) {
                if (id === 'toBusiness') {
                    this.assign.placeholder = this.$t('请选择xx', { name: this.$t('业务') })
                    this.assign.label = this.$t('业务列表')
                    this.assign.title = this.$t('分配到业务空闲机')
                } else {
                    this.assign.placeholder = this.$t('请选择xx', { name: this.$t('目录') })
                    this.assign.label = this.$t('目录列表')
                    this.assign.title = this.$t('分配到资源池其他目录')
                }
                this.assign.show = true
            },
            handleCancelAssignHosts () {
                this.assign.id = ''
                this.assign.show = false
                this.assign.curSelected = 'empty'
            },
            hasSelectAssignedHost () {
                const allList = this.$parent.table.list
                const list = allList.filter(item => this.table.checked.includes(item['host']['bk_host_id']))
                const existAssigned = list.some(item => item['biz'].some(biz => biz.default !== 1))
                return existAssigned
            },
            handleConfirmAssign () {
                this.assign.curSelected === 'toBusiness' ? this.assignHostsToBusiness() : this.changeHostsDir()
            },
            async assignHostsToBusiness () {
                const moduleId = this.activeDirectory.bk_module_id
                await this.$store.dispatch('resourceDirectory/assignHostsToBusiness', {
                    params: {
                        bk_module_id: moduleId,
                        bk_biz_id: this.assign.id,
                        bk_host_id: this.table.checked
                    },
                    config: {
                        requestId: this.assign.requestId
                    }
                }).then(() => {
                    Bus.$emit('refresh-dir-count', {
                        reduceId: moduleId,
                        count: this.table.checked.length
                    })
                    this.$success(this.$t('分配成功'))
                    this.$parent.table.checked = []
                    this.$parent.handlePageChange(1)
                    this.handleCancelAssignHosts()
                }).catch(e => {
                    console.error(e)
                })
            },
            async changeHostsDir () {
                try {
                    await this.$store.dispatch('resource/host/transfer/directory', {
                        params: {
                            bk_module_id: this.assign.id,
                            bk_host_id: this.table.checked
                        }
                    })
                    Bus.$emit('refresh-dir-count', {
                        reduceId: this.activeDirectory.bk_module_id,
                        addId: this.assign.id,
                        count: this.table.checked.length
                    })
                    this.$success(this.$t('转移成功'))
                    this.$parent.table.checked = []
                    this.$parent.handlePageChange(1)
                    this.handleCancelAssignHosts()
                } catch (e) {
                    console.error(e)
                }
            },
            getHostCellText (header, item) {
                const objId = header.objId
                const propertyId = header.id
                const headerProperty = this.$tools.getProperty(this.properties[objId], propertyId)
                const originalValues = item[objId] instanceof Array ? item[objId] : [item[objId]]
                const text = []
                originalValues.forEach(value => {
                    const flattenedText = this.$tools.getPropertyText(headerProperty, value)
                    flattenedText ? text.push(flattenedText) : void (0)
                })
                return text.join(',') || '--'
            },
            handleCopy (target) {
                const copyList = this.table.list.filter(item => {
                    return this.table.checked.includes(item['host']['bk_host_id'])
                })
                const copyText = []
                this.$tools.clone(copyList).forEach(item => {
                    const cellText = this.getHostCellText(target, item)
                    if (cellText !== '--') {
                        copyText.push(cellText)
                    }
                })
                if (copyText.length) {
                    this.$copyText(copyText.join('\n')).then(() => {
                        this.$success(this.$t('复制成功'))
                    }, () => {
                        this.$error(this.$t('复制失败'))
                    })
                } else {
                    this.$info(this.$t('该字段无可复制的值'))
                }
            },
            async handleMultipleEdit () {
                if (this.hasSelectAssignedHost()) {
                    this.$error(this.$t('请勿选择已分配主机'))
                    return false
                }
                this.objectUnique = await this.$store.dispatch('objectUnique/searchObjectUniqueConstraints', {
                    objId: 'host',
                    params: this.$injectMetadata({}, {
                        inject: this.$route.name !== 'resource'
                    })
                })
                this.tab.attribute.type = 'multiple'
                this.slider.title = this.$t('主机属性')
                this.slider.show = true
            },
            async handleMultipleSave (changedValues) {
                await this.$store.dispatch('hostUpdate/updateHost', {
                    params: this.$injectMetadata({
                        ...changedValues,
                        'bk_host_id': this.table.checked.join(',')
                    }, { inject: this.$route.name !== 'resource' })
                })
                this.slider.show = false
            },
            handleMultipleDelete () {
                if (this.hasSelectAssignedHost()) {
                    this.$error(this.$t('请勿选择已分配主机'))
                    return false
                }
                this.$bkInfo({
                    title: `${this.$t('确定删除选中的主机')}？`,
                    confirmFn: () => {
                        this.$store.dispatch('hostDelete/deleteHost', {
                            params: {
                                data: {
                                    'bk_host_id': this.table.checked.join(','),
                                    'bk_supplier_account': this.supplierAccount
                                }
                            }
                        }).then(() => {
                            this.$success(this.$t('成功删除选中的主机'))
                            this.$parent.table.checked = []
                            this.$parent.handlePageChange(1)
                        })
                    }
                })
            },
            handleSliderBeforeClose () {
                if (this.tab.active === 'attribute' && this.tab.attribute.type !== 'details') {
                    const $form = this.tab.attribute.type === 'update' ? this.$refs.form : this.$refs.multipleForm
                    const changedValues = $form.changedValues
                    if (Object.keys(changedValues).length) {
                        return new Promise((resolve, reject) => {
                            this.$bkInfo({
                                title: this.$t('确认退出'),
                                subTitle: this.$t('退出会导致未保存信息丢失'),
                                extCls: 'bk-dialog-sub-header-center',
                                confirmFn: () => {
                                    this.slider.show = false
                                },
                                cancelFn: () => {
                                    resolve(false)
                                }
                            })
                        })
                    }
                    this.slider.show = false
                }
                this.slider.show = false
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
                if (this.$parent.customColumns) {
                    formData.append('export_custom_fields', this.$parent.customColumns)
                }
                formData.append('bk_biz_id', '-1')
                const res = await this.$store.dispatch('hostBatch/exportHost', {
                    params: formData,
                    config: {
                        globalError: false,
                        originalResponse: true,
                        responseType: 'blob'
                    }
                })
                this.exportExcel(res)
            },
            handleColumnConfigClick () {
                Bus.$emit('toggle-host-filter', false)
                this.columnsConfig.show = true
            },
            handleCloseFilter () {
                this.$refs.hostFilter.$refs.filterPopper.instance.hide()
            },
            routeToHistory () {
                this.$router.push({ name: 'hostHistory' })
            },
            handleApplyColumnsConfig (properties) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$parent.columnsConfigKey]: properties.map(property => property['bk_property_id'])
                })
                this.columnsConfig.show = false
            },
            handleResetColumnsConfig () {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$parent.columnsConfigKey]: []
                })
            },
            async handleApplyPermission () {
                try {
                    const permission = []
                    const operation = this.$tools.getValue(this.$route.meta, 'auth.operation', {})
                    if (Object.keys(operation).length) {
                        const translated = await translateAuth(Object.values(operation))
                        permission.push(...translated)
                    }
                    const url = await this.$store.dispatch('auth/getSkipUrl', { params: permission })
                    window.open(url)
                } catch (e) {
                    console.error(e)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .options-layout {
        margin-bottom: 10px;
    }
    .options-left {
        float: left;
        font-size: 0;
    }
    .options-right {
        float: right;
        overflow: hidden;
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
    .assign-dialog {
        color: #63656E;
        .assign-count {
            font-size: 12px;
            padding-bottom: 20px;
            span {
                font-weight: bold;
            }
        }
        .assign-seleted {
            .bk-select {
                width: 100%;
                margin-top: 10px;
            }
        }
        .assign-footer {
            font-size: 0;
        }
    }
    .apply-others {
        display: inline-block;
        width: 60%;
        font-size: 12px;
        color: #63656E;
        line-height: 32px;
        cursor: pointer;
        &:hover {
            color: #3a84ff;
        }
        .bk-icon {
            font-size: 14px;
            display: inline-block;
            vertical-align: -2px;
        }
    }
</style>
