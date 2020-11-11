<template>
    <div class="options-layout clearfix">
        <div class="options-left">
            <template v-if="scope === 1">
                <cmdb-auth class="mr10"
                    :ignore="!!activeDirectory"
                    :auth="[
                        { type: $OPERATION.C_RESOURCE_HOST, relation: [directoryId] },
                        { type: $OPERATION.U_RESOURCE_HOST, relation: [directoryId] }
                    ]">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        style="margin-left: 0"
                        :disabled="disabled"
                        @click="importInst.show = true">
                        {{$t('导入主机')}}
                    </bk-button>
                </cmdb-auth>
                <bk-select
                    class="assign-selector mr10"
                    font-size="medium"
                    :popover-width="180"
                    :disabled="!table.checked.length"
                    :clearable="false"
                    :placeholder="$t('分配到')"
                    v-model="assign.curSelected"
                    @selected="handleAssignHosts">
                    <bk-option id="-1" :name="$t('分配到')" hidden></bk-option>
                    <bk-option id="toBusiness" :name="$t('业务空闲机')"></bk-option>
                    <bk-option id="toDirs" :name="$t('主机池其他目录')"></bk-option>
                </bk-select>
            </template>
            <cmdb-transfer-menu class="mr10"
                v-else>
            </cmdb-transfer-menu>
            <cmdb-clipboard-selector class="options-clipboard mr10"
                :list="clipboardList"
                :disabled="!table.checked.length"
                @on-copy="handleCopy">
            </cmdb-clipboard-selector>
            <cmdb-button-group
                class="mr10"
                :buttons="buttons"
                :expand="false">
            </cmdb-button-group>
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
                        <bk-form class="import-prepend" slot="prepend">
                            <bk-form-item :label="$t('主机池目录')" required>
                                <bk-select v-model="importInst.directory" style="display: block;">
                                    <bk-option v-for="directory in directoryList"
                                        :key="directory.bk_module_id"
                                        :id="directory.bk_module_id"
                                        :name="directory.bk_module_name">
                                    </bk-option>
                                </bk-select>
                            </bk-form-item>
                        </bk-form>
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
            <cmdb-form-multiple v-if="slider.show"
                slot="content"
                ref="multipleForm"
                :properties="properties.host"
                :property-groups="propertyGroups"
                :object-unique="objectUnique"
                :save-auth="saveAuth"
                @on-submit="handleMultipleSave"
                @on-cancel="handleSliderBeforeClose">
            </cmdb-form-multiple>
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
            @cancel="closeAssignDialog">
            <div class="assign-content" v-if="assign.show">
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
                            :name="option.name"
                            :disabled="option.disabled">
                            <cmdb-auth style="display: block;" ignore :auth="option.auth" @update-auth="handleUpdateAssignAuth(option, ...arguments)">{{option.name}}</cmdb-auth>
                        </bk-option>
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
                <bk-button theme="default" :disabled="$loading(assign.requestId)" @click="closeAssignDialog">{{$t('取消')}}</bk-button>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import cmdbImport from '@/components/import/import'
    import cmdbButtonGroup from '@/components/ui/other/button-group'
    import cmdbHostFilter from '@/components/hosts/filter/index.vue'
    import cmdbColumnsConfig from '@/components/columns-config/columns-config'
    import Bus from '@/utils/bus.js'
    import RouterQuery from '@/router/query'
    import HostStore from '../transfer/host-store'
    import cmdbTransferMenu from '../transfer/transfer-menu'
    export default {
        components: {
            cmdbImport,
            cmdbButtonGroup,
            cmdbHostFilter,
            cmdbColumnsConfig,
            cmdbTransferMenu
        },
        data () {
            return {
                scope: '',
                importInst: {
                    show: false,
                    active: 'import',
                    templateUrl: `${window.API_HOST}importtemplate/host`,
                    importUrl: `${window.API_HOST}hosts/import`,
                    directory: ''
                },
                businessList: [],
                objectUnique: [],
                slider: {
                    show: false,
                    title: ''
                },
                columnsConfig: {
                    show: false,
                    selected: [],
                    disabledColumns: ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id']
                },
                assign: {
                    show: false,
                    id: '',
                    curSelected: '-1',
                    placeholder: this.$t('请选择xx', { name: this.$t('业务') }),
                    label: this.$t('业务列表'),
                    title: this.$t('分配到业务空闲机'),
                    requestId: Symbol('assignHosts')
                },
                assignOptions: []
            }
        },
        computed: {
            ...mapGetters('resourceHost', [
                'activeDirectory',
                'defaultDirectory',
                'directoryList'
            ]),
            directoryId () {
                if (this.activeDirectory) {
                    return this.activeDirectory.bk_module_id
                }
                return this.defaultDirectory ? this.defaultDirectory.bk_module_id : undefined
            },
            table () {
                return this.$parent.table
            },
            clipboardList () {
                return this.table.header.map(property => {
                    return {
                        id: property.bk_property_id,
                        name: property.bk_property_name,
                        objId: property.bk_obj_id
                    }
                })
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
                const buttonConfig = [{
                    id: 'edit',
                    text: this.$t('编辑'),
                    handler: this.handleMultipleEdit,
                    disabled: !this.table.checked.length
                }, {
                    id: 'delete',
                    text: this.$t('删除'),
                    handler: this.handleMultipleDelete,
                    disabled: !this.table.checked.length
                }, {
                    id: 'export',
                    text: this.$t('导出'),
                    handler: this.exportField,
                    disabled: !this.table.checked.length
                }]
                if (this.scope !== 1) {
                    buttonConfig.splice(1, 1)
                }
                return buttonConfig
            },
            filterProperties () {
                const { module, set, host, biz } = this.properties
                const filterProperty = ['bk_host_innerip', 'bk_host_outerip']
                return {
                    host: host.filter(property => !filterProperty.includes(property.bk_property_id)),
                    module,
                    set,
                    biz
                }
            },
            saveAuth () {
                return this.table.selection.map(({ host, module, biz }) => {
                    if (biz[0].default === 0) {
                        return {
                            type: this.$OPERATION.U_HOST,
                            relation: [biz[0].bk_biz_id, host.bk_host_id]
                        }
                    }
                    return {
                        type: this.$OPERATION.U_RESOURCE_HOST,
                        relation: [module[0].bk_module_id, host.bk_host_id]
                    }
                })
            }
        },
        watch: {
            'importInst.show' (show) {
                if (!show) {
                    this.importInst.active = 'import'
                } else {
                    this.importInst.directory = this.directoryId
                }
            }
        },
        async created () {
            try {
                this.unwatchScope = RouterQuery.watch('scope', (scope = 1) => {
                    this.scope = isNaN(scope) ? 'all' : parseInt(scope)
                }, { immediate: true })
                await this.getFullAmountBusiness()
            } catch (e) {
                console.error(e.message)
            }
        },
        beforeDestroy () {
            this.unwatchScope()
        },
        methods: {
            async getFullAmountBusiness () {
                try {
                    const data = await this.$http.get('biz/simplify?sort=bk_biz_name')
                    this.businessList = data.info || []
                } catch (e) {
                    console.error(e)
                    this.businessList = []
                } finally {
                    HostStore.setBusinessList(this.businessList)
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
                    this.assign.title = this.$t('分配到主机池其他目录')
                }
                this.setAssignOptions()
                this.assign.show = true
            },
            setAssignOptions () {
                if (this.assign.curSelected === 'toBusiness') {
                    this.assignOptions = this.businessList.map(item => ({
                        id: item.bk_biz_id,
                        name: item.bk_biz_name,
                        disabled: true,
                        auth: {
                            type: this.$OPERATION.TRANSFER_HOST_TO_BIZ,
                            relation: [[[this.directoryId || '*'], [item.bk_biz_id]]]
                        }
                    }))
                } else {
                    this.assignOptions = this.directoryList.filter(item => item.bk_module_id !== (this.activeDirectory || {}).bk_module_id).map(item => ({
                        id: item.bk_module_id,
                        name: item.bk_module_name,
                        disabled: true,
                        auth: {
                            type: this.$OPERATION.TRANSFER_HOST_TO_DIRECTORY,
                            relation: [[[this.directoryId || '*'], [item.bk_module_id]]]
                        }
                    }))
                }
            },
            handleUpdateAssignAuth (option, authorized) {
                option.disabled = !authorized
            },
            closeAssignDialog () {
                this.assign.id = ''
                this.assign.show = false
                this.assign.curSelected = '-1'
            },
            handleConfirmAssign () {
                this.assign.curSelected === 'toBusiness' ? this.assignHostsToBusiness() : this.changeHostsDir()
            },
            async assignHostsToBusiness () {
                await this.$store.dispatch('resourceDirectory/assignHostsToBusiness', {
                    params: {
                        bk_biz_id: this.assign.id,
                        bk_host_id: this.table.checked
                    },
                    config: {
                        requestId: this.assign.requestId
                    }
                }).then(() => {
                    Bus.$emit('refresh-dir-count')
                    this.$success(this.$t('分配成功'))
                    this.closeAssignDialog()
                    RouterQuery.set({
                        page: 1,
                        _t: Date.now()
                    })
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
                    Bus.$emit('refresh-dir-count')
                    this.$success(this.$t('转移成功'))
                    this.closeAssignDialog()
                    RouterQuery.set({
                        page: 1,
                        _t: Date.now()
                    })
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
                copyList.forEach(item => {
                    if (target.id === '__bk_host_topology__') {
                        copyText.push((item.__bk_host_topology__ || []).join(','))
                    } else {
                        const cellText = this.getHostCellText(target, item)
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
                this.objectUnique = await this.$store.dispatch('objectUnique/searchObjectUniqueConstraints', {
                    objId: 'host',
                    params: {}
                })
                this.slider.title = this.$t('主机属性')
                this.slider.show = true
            },
            async handleMultipleSave (changedValues) {
                await this.$store.dispatch('hostUpdate/updateHost', {
                    params: {
                        ...changedValues,
                        'bk_host_id': this.table.checked.join(',')
                    }
                })
                this.slider.show = false
                RouterQuery.set({
                    _t: Date.now()
                })
            },
            handleMultipleDelete () {
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
                            RouterQuery.set({
                                page: 1,
                                _t: Date.now()
                            })
                            Bus.$emit('refresh-dir-count')
                        })
                    }
                })
            },
            handleSliderBeforeClose () {
                const $form = this.$refs.multipleForm
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
                this.$routerActions.redirect({
                    name: 'hostHistory',
                    history: true
                })
            },
            handleApplyColumnsConfig (properties) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$route.meta.customInstanceColumn]: properties.map(property => property['bk_property_id'])
                })
                this.columnsConfig.show = false
                RouterQuery.set({
                    _t: Date.now()
                })
            },
            handleResetColumnsConfig () {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$route.meta.customInstanceColumn]: []
                })
                this.columnsConfig.show = false
                RouterQuery.set({
                    _t: Date.now()
                })
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
        .assign-selector {
            min-width: 80px;
        }
    }
    .options-right {
        float: right;
        overflow: hidden;
    }
    .import-prepend {
        margin: 20px 29px -10px 33px;
        /deep/ {
            .bk-form-item {
                display: flex;
            }
            .bk-label {
                width: auto !important;
            }
            .bk-form-content {
                flex: 1;
                margin-left: auto !important;
            }
        }
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
