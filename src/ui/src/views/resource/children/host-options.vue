<template>
    <div class="options-layout clearfix">
        <div class="options-left fl">
            <template v-if="scope === 1">
                <cmdb-auth class="mr10"
                    :ignore="!activeDirectory"
                    :auth="[
                        { type: $OPERATION.C_RESOURCE_HOST, relation: [directoryId] },
                        { type: $OPERATION.U_RESOURCE_HOST, relation: [directoryId] }
                    ]">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        style="margin-left: 0"
                        :disabled="disabled"
                        @click="handleNewImportInst">
                        {{$t('导入主机')}}
                    </bk-button>
                </cmdb-auth>
            </template>
            <span style="display: inline-block;"
                v-bk-tooltips="{ content: this.$t('仅主机池主机可进行此操作'), disabled: isAllResourceHost }">
                <bk-select
                    class="assign-selector mr10"
                    font-size="medium"
                    :popover-width="180"
                    :disabled="!table.checked.length || !isAllResourceHost"
                    :clearable="false"
                    :placeholder="$t('分配到')"
                    v-model="assign.curSelected"
                    @selected="handleAssignHosts">
                    <bk-option id="-1" :name="$t('分配到')" hidden></bk-option>
                    <bk-option id="toBusiness" :name="$t('业务空闲机')"></bk-option>
                    <bk-option id="toDirs" :name="$t('主机池其他目录')"></bk-option>
                </bk-select>
            </span>
            <cmdb-transfer-menu class="mr10" v-if="scope !== 1" />
            <cmdb-button-group
                class="mr10"
                trigger-text="编辑"
                :buttons="editButtonGroup"
                :expand="false">
            </cmdb-button-group>
            <cmdb-clipboard-selector class="mr10"
                label-key="bk_property_name"
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
            <filter-fast-search class="option-fast-search"></filter-fast-search>
            <icon-button class="option-filter ml10" icon="icon-cc-funnel" v-bk-tooltips.top="$t('高级筛选')" @click="handleSetFilters"></icon-button>
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
            :title="importInst.title">
            <bk-tab :active.sync="importInst.active" type="unborder-card" slot="content"
                v-if="importInst.show && importInst.type === 'new'">
                <bk-tab-panel name="import" :label="$t('新增导入')">
                    <cmdb-import v-if="importInst.show && importInst.active === 'import'"
                        :template-url="importInst.templateUrl"
                        :import-url="importInst.importUrl"
                        :import-payload="importInst.payload"
                        :global-error="false"
                        :before-upload="handleBeforeUpload"
                        @error="handleImportError"
                        @success="handleImportSuccess"
                        @partialSuccess="handleImportSuccess">
                        <bk-form class="import-prepend" slot="prepend">
                            <bk-form-item :label="$t('主机池目录')" required>
                                <bk-select v-model="importInst.directory" searchable style="display: block;">
                                    <cmdb-auth-option v-for="directory in directoryList"
                                        :key="directory.bk_module_id"
                                        :id="directory.bk_module_id"
                                        :name="directory.bk_module_name"
                                        :auth="{ type: $OPERATION.C_RESOURCE_HOST, relation: [directory.bk_module_id] }">
                                    </cmdb-auth-option>
                                </bk-select>
                                <p class="form-item-tips" v-show="importInst.showTips && !importInst.directory">{{$t('请先选择主机池目录')}}</p>
                            </bk-form-item>
                        </bk-form>
                        <span slot="download-desc" style="display: inline-block;vertical-align: top;">
                            {{$t('说明：内网IP为必填列')}}
                        </span>
                        <div slot="uploadErrorMessage" class="upload-error-message" v-if="importInstError">{{importInstError}}</div>
                    </cmdb-import>
                </bk-tab-panel>
                <bk-tab-panel name="agent" :label="$t('Agent导入')">
                    <div class="automatic-import">
                        <img src="../../../assets/images/agent-import-guide.png">
                        <p class="agent-install-tips1">{{$t("agent安装说明")}}</p>
                        <p class="agent-install-tips2">{{$t("跳转节点管理，支持远程 / 手动安装")}}</p>
                        <bk-button class="agent-install-button" theme="primary" @click="openAgentApp">{{$t('跳转安装')}}</bk-button>
                    </div>
                </bk-tab-panel>
            </bk-tab>
            <div slot="content" class="edit-import-panel" v-if="importInst.type === 'edit'">
                <bk-alert class="alert" type="warning" :title="$t('请上传导出的主机表格文件')"></bk-alert>
                <cmdb-import :template-url="importInst.templateUrl"
                    :import-url="importInst.importUrl"
                    :import-payload="importInst.payload"
                    :templdate-available="importInst.templdateAvailable"
                    :global-error="false"
                    @error="handleImportError"
                    @success="handleImportSuccess"
                    @partialSuccess="handleImportSuccess">
                    <span slot="successTips" slot-scope="{ success }">
                        {{$t('更新成功N个主机数据', { N: success.length })}}
                    </span>
                    <span slot="errorTips" slot-scope="{ error }">
                        {{$t('更新失败N个主机数据', { N: error.length })}}
                    </span>
                    <div slot="uploadErrorMessage" class="upload-error-message" v-if="importInstError">{{importInstError}}</div>
                </cmdb-import>
            </div>
        </bk-sideslider>

        <bk-sideslider
            v-transfer-dom
            :is-show.sync="slider.show"
            :title="slider.title"
            :width="800"
            :before-close="handleSliderBeforeClose">
            <div slot="content" style="height: 100%" v-bkloading="{ isLoading: $loading(Object.values(slider.request)) }">
                <component :is="slider.component"
                    ref="multipleForm"
                    :properties="properties.host"
                    :property-groups="propertyGroups"
                    :object-unique="objectUnique"
                    :save-auth="saveAuth"
                    @on-submit="handleMultipleSave"
                    @on-cancel="handleSliderBeforeClose">
                </component>
            </div>
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
                            <cmdb-auth style="display: block;" :auth="option.auth" @update-auth="handleUpdateAssignAuth(option, ...arguments)">{{option.name}}</cmdb-auth>
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
    import Bus from '@/utils/bus.js'
    import RouterQuery from '@/router/query'
    import HostStore from '../transfer/host-store'
    import cmdbTransferMenu from '../transfer/transfer-menu'
    import FilterForm from '@/components/filters/filter-form.js'
    import FilterFastSearch from '@/components/filters/filter-fast-search'
    import FilterStore from '@/components/filters/store'
    import ExportFields from '@/components/export-fields/export-fields.js'
    import FilterUtils from '@/components/filters/utils'
    import BatchExport from '@/components/batch-export/index.js'
    export default {
        components: {
            cmdbImport,
            cmdbButtonGroup,
            cmdbTransferMenu,
            FilterFastSearch
        },
        data () {
            return {
                scope: '',
                importInst: {
                    title: '',
                    show: false,
                    active: 'import',
                    templateUrl: `${window.API_HOST}importtemplate/host`,
                    importUrl: '',
                    templdateAvailable: true,
                    directory: '',
                    payload: {},
                    error: null,
                    showTips: false
                },
                businessList: [],
                objectUnique: [],
                slider: {
                    show: false,
                    title: '',
                    component: null,
                    request: {
                        objectUnique: Symbol('objectUnique')
                    }
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
                assignOptions: [],
                IPWithCloudSymbol: Symbol('IPWithCloud')
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
                return undefined
            },
            table () {
                return this.$parent.table
            },
            isAllResourceHost () {
                return this.table.selection.every(({ biz }) => {
                    const [currentBiz] = biz
                    return currentBiz.default === 1
                })
            },
            clipboardList () {
                const IPWithCloud = FilterUtils.defineProperty({
                    id: this.IPWithCloudSymbol,
                    bk_obj_id: 'host',
                    bk_property_id: this.IPWithCloudSymbol,
                    bk_property_name: `${this.$t('云区域')}ID:IP`,
                    bk_property_type: 'singlechar'
                })
                const clipboardList = FilterStore.header.slice()
                clipboardList.splice(1, 0, IPWithCloud)
                return clipboardList
            },
            properties () {
                return FilterStore.modelPropertyMap
            },
            propertyGroups () {
                return FilterStore.propertyGroups
            },
            buttons () {
                const buttonConfig = [{
                    id: 'delete',
                    text: this.$t('删除'),
                    handler: this.handleMultipleDelete,
                    disabled: !this.table.checked.length
                }, {
                    id: 'export',
                    text: this.$t('导出选中'),
                    handler: this.exportField,
                    disabled: !this.table.checked.length
                }, {
                    id: 'batchExport',
                    text: this.$t('导出全部'),
                    handler: this.batchExportField,
                    disabled: !this.table.pagination.count
                }]
                if (this.scope !== 1) {
                    buttonConfig.splice(0, 1)
                }
                return buttonConfig
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
            },
            editButtonGroup () {
                const buttonConfig = [{
                    id: 'batch-edit',
                    text: this.$t('批量编辑'),
                    handler: this.handleMultipleEdit,
                    disabled: !this.table.checked.length,
                    tooltips: { content: this.$t('请先勾选需要编辑的主机'), disabled: this.table.checked.length > 0 }
                }, {
                    id: 'import-edit',
                    text: this.$t('导入编辑'),
                    handler: this.handleEditImportInst
                }]
                if (this.scope !== 1) {
                    buttonConfig.pop()
                }
                return buttonConfig
            },
            importInstError () {
                const importInstError = this.importInst.error || {}
                if (importInstError.bk_error_msg) {
                    return importInstError.bk_error_msg
                }
                return importInstError.message || ''
            }
        },
        watch: {
            'importInst.show' (show) {
                if (!show) {
                    this.importInst.type = 'new'
                    this.importInst.active = 'import'
                    this.importInst.error = null
                    this.importInst.showTips = false
                } else {
                    this.importInst.directory = this.directoryId
                }
            },
            'importInst.directory' (directory) {
                if (this.importInst.type === 'new') {
                    this.importInst.payload = {
                        bk_module_id: directory
                    }
                }
            },
            'importInst.active' () {
                this.importInst.error = null
                this.importInst.showTips = false
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
                    const data = await this.$http.get('biz/simplify?sort=bk_biz_id')
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
                let directoryId = this.directoryId
                if (!directoryId) {
                    const hosts = HostStore.getSelected()
                    directoryId = hosts[0].module[0].bk_module_id
                    const isSameModule = hosts.every(host => {
                        const [module] = host.module
                        return module.bk_module_id === directoryId
                    })
                    if (!isSameModule) {
                        this.$error(this.$t('仅支持对相同目录下的主机进行操作'))
                        this.closeAssignDialog()
                        return false
                    }
                }

                if (id === 'toBusiness') {
                    this.assign.placeholder = this.$t('请选择xx', { name: this.$t('业务') })
                    this.assign.label = this.$t('业务列表')
                    this.assign.title = this.$t('分配到业务空闲机')
                } else {
                    this.assign.placeholder = this.$t('请选择xx', { name: this.$t('目录') })
                    this.assign.label = this.$t('目录列表')
                    this.assign.title = this.$t('分配到主机池其他目录')
                }

                this.setAssignOptions(directoryId)
                this.assign.show = true
            },
            setAssignOptions (directoryId) {
                if (this.assign.curSelected === 'toBusiness') {
                    this.assignOptions = this.businessList.map(item => ({
                        id: item.bk_biz_id,
                        name: `[${item.bk_biz_id}] ${item.bk_biz_name}`,
                        disabled: true,
                        auth: {
                            type: this.$OPERATION.TRANSFER_HOST_TO_BIZ,
                            relation: [[[directoryId], [item.bk_biz_id]]]
                        }
                    }))
                } else {
                    this.assignOptions = this.directoryList.filter(item => item.bk_module_id !== directoryId).map(item => ({
                        id: item.bk_module_id,
                        name: item.bk_module_name,
                        disabled: true,
                        auth: {
                            type: this.$OPERATION.TRANSFER_HOST_TO_DIRECTORY,
                            relation: [[[directoryId], [item.bk_module_id]]]
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
                        },
                        config: {
                            requestId: this.assign.requestId
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
            handleCopy (property) {
                const copyText = this.table.selection.map(data => {
                    const modelId = property.bk_obj_id
                    const [modelData] = Array.isArray(data[modelId]) ? data[modelId] : [data[modelId]]
                    if (property.id === this.IPWithCloudSymbol) {
                        const cloud = this.$tools.getPropertyCopyValue(modelData.bk_cloud_id, 'foreignkey')
                        const ip = this.$tools.getPropertyCopyValue(modelData.bk_host_innerip, 'singlechar')
                        return `${cloud}:${ip}`
                    }
                    if (property.bk_property_type === 'topology') {
                        return data.__bk_host_topology__.join(',').replace(/\s\/\s/g, '')
                    }
                    const value = modelData[property.bk_property_id]
                    return this.$tools.getPropertyCopyValue(value, property)
                })
                this.$copyText(copyText.join('\n')).then(() => {
                    this.$success(this.$t('复制成功'))
                }, () => {
                    this.$error(this.$t('复制失败'))
                })
            },
            async handleMultipleEdit () {
                this.slider.title = this.$t('主机属性')
                this.slider.show = true
                this.objectUnique = await this.$store.dispatch('objectUnique/searchObjectUniqueConstraints', {
                    objId: 'host',
                    params: {},
                    config: {
                        requestId: this.slider.request.objectUnique
                    }
                })
                setTimeout(() => {
                    this.slider.component = 'cmdb-form-multiple'
                }, 200)
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
                                this.slider.component = null
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.slider.show = false
                this.slider.component = null
            },
            async exportField () {
                ExportFields.show({
                    title: this.$t('导出选中'),
                    properties: FilterStore.getModelProperties('host'),
                    propertyGroups: FilterStore.propertyGroups,
                    handler: this.exportHanlder
                })
            },
            async exportHanlder (properties) {
                const formData = new FormData()
                formData.append('bk_biz_id', -1)
                formData.append('bk_host_id', this.table.selection.map(({ host }) => host.bk_host_id).join(','))
                formData.append('export_custom_fields', properties.map(property => property.bk_property_id))
                try {
                    this.$store.commit('setGlobalLoading', true)
                    await this.$http.download({
                        url: `${window.API_HOST}hosts/export`,
                        method: 'post',
                        data: formData
                    })
                } catch (error) {
                    console.error(error)
                } finally {
                    this.$store.commit('setGlobalLoading', false)
                }
            },
            async batchExportField () {
                ExportFields.show({
                    title: this.$t('导出全部'),
                    properties: FilterStore.getModelProperties('host'),
                    propertyGroups: FilterStore.propertyGroups,
                    handler: this.batchExportHandler
                })
            },
            batchExportHandler (properties) {
                BatchExport({
                    name: 'host',
                    count: this.table.pagination.count,
                    options: page => {
                        const condition = this.$parent.getParams()
                        const formData = new FormData()
                        formData.append('bk_biz_id', -1)
                        formData.append('export_custom_fields', properties.map(property => property.bk_property_id))
                        formData.append('export_condition', JSON.stringify({
                            ...condition,
                            page: {
                                ...page,
                                sort: 'bk_host_id'
                            }
                        }))
                        return {
                            url: `${window.API_HOST}hosts/export`,
                            method: 'post',
                            data: formData
                        }
                    }
                })
            },
            routeToHistory () {
                this.$routerActions.redirect({
                    name: 'hostHistory',
                    history: true
                })
            },
            handleSetFilters () {
                FilterForm.show()
            },
            handleNewImportInst () {
                this.importInst.type = 'new'
                this.importInst.show = true
                this.importInst.title = this.$t('导入主机')
                this.importInst.importUrl = `${window.API_HOST}hosts/import`
                this.importInst.templdateAvailable = true
            },
            handleEditImportInst () {
                this.importInst.type = 'edit'
                this.importInst.show = true
                this.importInst.title = this.$t('导入编辑')
                this.importInst.importUrl = `${window.API_HOST}hosts/update`
                this.importInst.templdateAvailable = false
                this.importInst.payload = {
                    // 资源池约定为0
                    bk_biz_id: 0
                }
            },
            handleImportSuccess () {
                this.$parent.getHostList(true)
                Bus.$emit('refresh-dir-count')
                this.importInst.error = null
            },
            handleImportError (error) {
                this.importInst.error = error
            },
            handleBeforeUpload () {
                this.importInst.showTips = false
                if (!this.importInst.directory) {
                    this.importInst.showTips = true
                    return false
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .options-left {
        font-size: 0;
        .assign-selector {
            min-width: 80px;
        }
    }
    .options-right {
        display: flex;
        justify-content: flex-end;
        overflow: hidden;
        .option-fast-search {
            flex: 1;
            margin-left: 10px;
            max-width: 300px;
        }
        .option-filter:hover {
            color: $primaryColor;
        }
    }
    .import-prepend {
        margin: 20px 29px -10px 33px;
        /deep/ {
            .bk-form-item {
                display: flex;
                margin-bottom: 8px;
            }
            .bk-label {
                width: auto !important;
            }
            .bk-form-content {
                flex: 1;
                margin-left: auto !important;
            }
        }
        .form-item-tips {
            position: absolute;
            color: #EA3636;
            font-size: 12px;
        }
    }
    .automatic-import{
        padding: 44px 30px 0 30px;
        display: flex;
        flex-direction: column;
        align-items: center;
        .agent-install-tips1 {
            font-size: 14px;
            margin: 10px 0;
        }
        .agent-install-tips2 {
            font-size: 12px;
            color: #979BA5;
        }
        .agent-install-button {
            margin: 14px 0;
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

    .edit-import-panel {
        .alert {
            margin: 24px 29px 0 33px
        }
        /deep/ {
            .up-file {
                margin-top: 20px;
            }
        }
    }

    .upload-error-message {
        margin: 8px 0;
        color: $dangerColor;
    }
</style>
