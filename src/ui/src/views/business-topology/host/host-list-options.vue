<template>
    <div class="options-layout clearfix">
        <div class="options fl">
            <cmdb-auth class="option" :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }">
                <bk-button theme="primary" slot-scope="{ disabled }"
                    :disabled="disabled || !isNormalModuleNode"
                    :title="isNormalModuleNode ? '' : $t('仅能在业务模块下新增')"
                    @click="handleAddHost">
                    {{$t('新增')}}
                </bk-button>
            </cmdb-auth>
            <bk-button class="ml10"
                :disabled="!hasSelection"
                @click="handleMultipleEdit">
                {{$t('编辑')}}
            </bk-button>
            <bk-dropdown-menu class="option ml10" trigger="click"
                font-size="medium"
                :disabled="!hasSelection"
                @show="isTransferMenuOpen = true"
                @hide="isTransferMenuOpen = false">
                <bk-button slot="dropdown-trigger"
                    :disabled="!hasSelection">
                    <span>{{$t('转移到')}}</span>
                    <i :class="['dropdown-icon bk-icon icon-angle-down',{ 'open': isTransferMenuOpen }]"></i>
                </bk-button>
                <ul class="bk-dropdown-list" slot="dropdown-content">
                    <li class="bk-dropdown-item"
                        @click="handleTransfer($event, 'idle', false)">
                        {{$t('空闲模块')}}
                    </li>
                    <cmdb-auth tag="li" class="bk-dropdown-item"
                        :auth="[
                            { type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] },
                            { type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] },
                            { type: $OPERATION.D_SERVICE_INSTANCE, relation: [bizId] }
                        ]"
                        @click="handleTransfer($event, 'business', false)">
                        {{$t('业务模块')}}
                    </cmdb-auth>
                    <li :class="['bk-dropdown-item', { disabled: !isIdleSetModules }]"
                        v-bk-tooltips="{
                            disabled: isIdleSetModules,
                            content: $t('仅空闲模块主机才能转移到其他业务')
                        }"
                        @click="handleTransfer($event, 'acrossBusiness', !isIdleSetModules)">
                        {{$t('其他业务')}}
                    </li>
                    <!-- 暂忽略鉴权，交互待调整，需要选择目录 -->
                    <cmdb-auth tag="li" class="bk-dropdown-item with-auth"
                        ignore
                        :auth="{ type: $OPERATION.HOST_TO_RESOURCE, relation: [bizId] }">
                        <span href="javascript:void(0)" slot-scope="{ disabled }"
                            v-bk-tooltips="isIdleModule ? '' : $t('仅空闲机模块才能转移到主机池')"
                            :class="{ disabled: !isIdleModule || disabled }"
                            @click="handleTransfer($event, 'resource', !isIdleModule)">
                            {{$t('主机池')}}
                        </span>
                    </cmdb-auth>
                </ul>
            </bk-dropdown-menu>
            <cmdb-clipboard-selector class="options-button ml10"
                label-key="bk_property_name"
                :list="clipboardList"
                :disabled="!hasSelection"
                @on-copy="handleCopy">
            </cmdb-clipboard-selector>
            <bk-dropdown-menu class="option ml10" trigger="click"
                font-size="medium"
                @show="isMoreMenuOpen = true"
                @hide="isMoreMenuOpen = false">
                <bk-button slot="dropdown-trigger">
                    <span>{{$t('更多')}}</span>
                    <i :class="['dropdown-icon bk-icon icon-angle-down',{ 'open': isMoreMenuOpen }]"></i>
                </bk-button>
                <ul class="bk-dropdown-list" slot="dropdown-content">
                    <cmdb-auth tag="li" class="bk-dropdown-item with-auth"
                        v-if="showRemove"
                        :auth="{ type: $OPERATION.D_SERVICE_INSTANCE, relation: [bizId] }">
                        <span href="javascript:void(0)"
                            slot-scope="{ disabled }"
                            :class="{ disabled: !hasSelection || disabled }"
                            @click="handleRemove($event)">
                            {{$t('移除')}}
                        </span>
                    </cmdb-auth>
                    <li :class="['bk-dropdown-item', { disabled: !hasSelection }]" @click="handleExport($event)">{{$t('导出')}}</li>
                    <cmdb-auth tag="li" class="bk-dropdown-item with-auth"
                        :auth="{ type: $OPERATION.U_HOST, relation: [bizId] }">
                        <span href="javascript:void(0)"
                            slot-scope="{ disabled }"
                            :class="{ disabled: disabled }"
                            @click="handleExcelUpdate($event)">
                            {{$t('导入excel更新')}}
                        </span>
                    </cmdb-auth>
                </ul>
            </bk-dropdown-menu>
        </div>
        <div class="options fr">
            <bk-select class="option option-collection bgc-white"
                ref="collectionSelector"
                v-model="selectedCollection"
                font-size="medium"
                searchable
                :loading="$loading(request.collection)"
                :placeholder="$t('请选择收藏条件')"
                @selected="handleCollectionSelect"
                @clear="handleCollectionClear"
                @toggle="handleCollectionToggle">
                <bk-option v-for="collection in collectionList"
                    :key="collection.id"
                    :id="collection.id"
                    :name="collection.name">
                    <div class="collection-item">
                        <span class="collection-name" :title="collection.name">{{collection.name}}</span>
                        <i class="bk-icon icon-close" @click.stop="handleDeleteCollection(collection)"></i>
                    </div>
                </bk-option>
                <div slot="extension">
                    <a href="javascript:void(0)" class="collection-create" @click="handleCreateCollection">
                        <i class="bk-icon icon-plus-circle"></i>
                        {{$t('新增条件')}}
                    </a>
                </div>
            </bk-select>
            <host-filter class="ml10" ref="hostFilter" :properties="filterProperties" :section-height="$APP.height - 250"></host-filter>
            <icon-button class="option ml10" icon="icon-cc-setting" @click="handleSetColumn"></icon-button>
        </div>
        <edit-multiple-host ref="editMultipleHost"
            :properties="hostProperties"
            :selection="$parent.table.selection">
        </edit-multiple-host>
        <cmdb-dialog v-model="dialog.show" v-bind="dialog.props" :height="650">
            <component
                :is="dialog.component"
                v-bind="dialog.componentProps"
                @confirm="handleDialogConfirm"
                @cancel="handleDialogCancel">
            </component>
        </cmdb-dialog>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="sideslider.show"
            :width="600"
            :title="sideslider.title">
            <component slot="content"
                :is="sideslider.component"
                v-bind="sideslider.componentProps"
                @on-cancel="sideslider.show = false"
                @on-apply="handleApplyColumnsConfig"
                @on-reset="handleResetColumnsConfig">
            </component>
        </bk-sideslider>
    </div>
</template>

<script>
    import HostFilter from '@/components/hosts/filter'
    import CmdbImport from '@/components/import/import'
    import EditMultipleHost from './edit-multiple-host.vue'
    import HostSelector from './host-selector-new'
    import CmdbColumnsConfig from '@/components/columns-config/columns-config'
    import { mapGetters, mapState } from 'vuex'
    import {
        MENU_BUSINESS,
        MENU_BUSINESS_TRANSFER_HOST
    } from '@/dictionary/menu-symbol'
    import Formatter from '@/filters/formatter.js'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            HostFilter,
            EditMultipleHost,
            [CmdbColumnsConfig.name]: CmdbColumnsConfig,
            [HostSelector.name]: HostSelector,
            [CmdbImport.name]: CmdbImport
        },
        data () {
            return {
                isTransferMenuOpen: false,
                isMoreMenuOpen: false,
                selectedCollection: '',
                dialog: {
                    show: false,
                    props: {
                        width: 1100,
                        showCloseIcon: false
                    },
                    component: null,
                    componentProps: {}
                },
                sideslider: {
                    show: false,
                    title: '',
                    component: null,
                    componentProps: {}
                },
                request: {
                    collection: Symbol('collection')
                }
            }
        },
        computed: {
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapState('hosts', ['collectionList']),
            ...mapGetters('businessHost', [
                'getProperties',
                'selectedNode'
            ]),
            hostProperties () {
                return this.getProperties('host')
            },
            columnsConfigProperties () {
                const setProperties = this.getProperties('set').filter(property => ['bk_set_name'].includes(property['bk_property_id']))
                const moduleProperties = this.getProperties('module').filter(property => ['bk_module_name'].includes(property['bk_property_id']))
                const hostProperties = this.getProperties('host')
                return [...setProperties, ...moduleProperties, ...hostProperties]
            },
            columnDisplayProperties () {
                return this.$parent.table.header.map(property => property.bk_property_id)
            },
            filterProperties () {
                const setProperties = this.getProperties('set')
                const moduleProperties = this.getProperties('module')
                const removeProperties = ['bk_host_id', 'bk_host_innerip', 'bk_host_outerip']
                // 模块支持服务分类筛选
                const insertProperties = [
                    {
                        bk_obj_id: 'module',
                        bk_property_id: 'service_category_id',
                        bk_property_name: this.$t('服务分类'),
                        bk_property_type: 'category'
                    }
                ]
                const hostProperties = this.hostProperties.filter(property => !removeProperties.includes(property.bk_property_id))
                return {
                    host: hostProperties,
                    set: setProperties,
                    module: [...moduleProperties, ...insertProperties]
                }
            },
            hasSelection () {
                return !!this.$parent.table.selection.length
            },
            isNormalModuleNode () {
                return this.selectedNode
                    && this.selectedNode.data.bk_obj_id === 'module'
                    && this.selectedNode.data.default === 0
            },
            isIdleModule () {
                return this.$parent.table.selection.every(data => {
                    const modules = data.module
                    return modules.every(module => module.default === 1)
                })
            },
            isIdleSetModules () {
                return this.$parent.table.selection.every(data => {
                    return data.module.every(module => module.default >= 1)
                })
            },
            clipboardList () {
                return this.$parent.table.header
            },
            showRemove () {
                return this.selectedNode
                    && !this.selectedNode.data.is_idle_set
                    && this.selectedNode.data.bk_obj_id === 'module'
                    && this.selectedNode.data.default !== 1
            }
        },
        created () {
            this.getCollectionList()
        },
        methods: {
            async getCollectionList () {
                try {
                    const result = await this.$store.dispatch('hostFavorites/searchFavorites', {
                        params: {
                            condition: {
                                bk_biz_id: this.bizId
                            }
                        },
                        config: {
                            requestId: this.request.condition
                        }
                    })
                    this.$store.commit('hosts/setCollectionList', result.info)
                } catch (e) {
                    console.error(e)
                }
            },
            handleCollectionSelect (value) {
                const collection = this.collectionList.find(collection => collection.id === value)
                try {
                    const filterList = JSON.parse(collection.query_params).map(condition => {
                        return {
                            bk_obj_id: condition.bk_obj_id,
                            bk_property_id: condition.field,
                            operator: condition.operator,
                            value: condition.value
                        }
                    })
                    const condition = this.getCollectionCondition(filterList)
                    this.$store.commit('hosts/setFilterList', filterList)
                    this.$store.commit('hosts/setCollection', collection)
                    this.$store.commit('hosts/setCondition', condition)
                    const info = JSON.parse(collection.info)
                    RouterQuery.set({
                        ip: info.ip_list.join(','),
                        exact: info.exact_search ? 1 : 0,
                        inner: info.bk_host_innerip ? 1 : 0,
                        outer: info.bk_host_outerip ? 1 : 0,
                        page: 1,
                        _t: Date.now()
                    })
                } catch (e) {
                    this.$error(this.$t('应用收藏条件失败，转换数据错误'))
                    console.error(e.message)
                }
            },
            getCollectionCondition (filterList = []) {
                const conditionMap = {}
                filterList.forEach(filter => {
                    const objectConditon = conditionMap[filter.bk_obj_id] || []
                    objectConditon.push({
                        field: filter.bk_property_id,
                        operator: filter.operator,
                        value: filter.value
                    })
                    conditionMap[filter.bk_obj_id] = objectConditon
                })
                return ['biz', 'object', 'set', 'module', 'host'].map(objId => {
                    return {
                        bk_obj_id: objId,
                        condition: conditionMap[objId] || [],
                        fields: []
                    }
                })
            },
            handleCollectionClear () {
                this.$store.commit('hosts/clearFilter')
                this.$refs.hostFilter.handleReset()
                this.$refs.hostFilter.$refs.filterPopper.instance.hide()
                const key = this.$route.meta.customFilterProperty
                const customData = this.$store.getters['userCustom/getCustomData'](key, [])
                this.$store.commit('hosts/setFilterList', customData)
            },
            handleCollectionToggle (isOpen) {
                if (isOpen) {
                    this.$refs.hostFilter.$refs.filterPopper.instance.hide()
                }
            },
            handleDeleteCollection (collection) {
                this.$bkInfo({
                    title: this.$t('确定删除收藏条件', { name: collection.name }),
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('hostFavorites/deleteFavorites', {
                                id: collection.id,
                                config: {
                                    requestId: 'deleteFavorites'
                                }
                            })
                            this.$success(this.$t('删除成功'))
                            this.selectedCollection = ''
                            this.$store.commit('hosts/deleteCollection', collection.id)
                            this.handleCollectionClear()
                        } catch (error) {
                            console.error(error)
                        }
                    }
                })
            },
            handleCreateCollection () {
                this.$store.commit('hosts/clearFilter')
                this.$refs.collectionSelector.close()
                this.$refs.hostFilter.handleToggleFilter()
            },
            handleTransfer (event, type, disabled) {
                if (disabled) {
                    event.stopPropagation()
                    return false
                }
                this.$emit('transfer', type)
            },
            handleMultipleEdit () {
                this.$refs.editMultipleHost.handleMultipleEdit()
            },
            handleAddHost () {
                this.dialog.component = HostSelector.name
                this.dialog.show = true
            },
            handleRemove (event) {
                if (!this.hasSelection) {
                    event.stopPropagation()
                    return false
                }
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        type: 'remove'
                    },
                    query: {
                        sourceModel: this.selectedNode.data.bk_obj_id,
                        sourceId: this.selectedNode.data.bk_inst_id,
                        resources: this.$parent.table.selection.map(item => item.host.bk_host_id).join(','),
                        node: this.selectedNode.id
                    },
                    history: true
                })
            },
            async handleExport (event) {
                if (!this.hasSelection) {
                    event.stopPropagation()
                    return false
                }
                try {
                    this.$store.commit('setGlobalLoading', true)
                    const data = new FormData()
                    data.append('bk_biz_id', -1)
                    data.append('bk_host_id', this.$parent.table.selection.map(item => item.host.bk_host_id).join(','))
                    const customFields = this.usercustom[this.$route.meta.customInstanceColumn]
                    if (customFields) {
                        data.append('export_custom_fields', customFields)
                    }
                    if (this.$route.meta.owner === MENU_BUSINESS) {
                        data.append('bk_biz_id', this.$route.params.bizId)
                    }
                    await this.$http.download({
                        url: `${window.API_HOST}hosts/export`,
                        method: 'post',
                        data
                    })
                    this.$store.commit('setGlobalLoading', false)
                } catch (e) {
                    console.error(e)
                    this.$store.commit('setGlobalLoading', false)
                }
            },
            handleExcelUpdate (event) {
                this.sideslider.component = CmdbImport.name
                this.sideslider.componentProps = {
                    templateUrl: `${window.API_HOST}importtemplate/host`,
                    importUrl: `${window.API_HOST}hosts/update`,
                    templdateAvailable: false
                }
                this.sideslider.title = this.$t('更新主机属性')
                this.sideslider.show = true
            },
            handleCopy (target) {
                const copyList = this.$parent.table.selection
                const copyText = []
                copyList.forEach(item => {
                    const modelData = Array.isArray(item[target.bk_obj_id]) ? item[target.bk_obj_id] : [item[target.bk_obj_id]]
                    const curCopyText = []
                    modelData.forEach(data => {
                        const value = data[target.bk_property_id]
                        const formattedValue = Formatter(value, target)
                        curCopyText.push(formattedValue)
                    })
                    if (curCopyText.length) {
                        copyText.push(curCopyText.join(','))
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
            handleDialogConfirm () {
                if (this.dialog.component === HostSelector.name) {
                    this.gotoTransferPage(...arguments)
                }
            },
            gotoTransferPage (selected) {
                this.$routerActions.redirect({
                    name: 'createServiceInstance',
                    params: {
                        setId: this.selectedNode.parent.data.bk_inst_id,
                        moduleId: this.selectedNode.data.bk_inst_id
                    },
                    query: {
                        resources: selected.map(item => item.host.bk_host_id).join(','),
                        title: this.selectedNode.data.bk_inst_name,
                        node: this.selectedNode.id
                    },
                    history: true
                })
            },
            handleDialogCancel () {
                this.dialog.show = false
            },
            async handleApplyColumnsConfig (properties) {
                this.sideslider.show = false
                await this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$route.meta.customInstanceColumn]: properties.map(property => property['bk_property_id'])
                })
                RouterQuery.set({ _t: Date.now() }) // 触发刷新表格
            },
            handleResetColumnsConfig () {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$route.meta.customInstanceColumn]: []
                })
                this.sideslider.show = false
            },
            handleSetColumn () {
                this.$refs.hostFilter.$refs.filterPopper.instance.hide()
                this.sideslider.component = CmdbColumnsConfig.name
                this.sideslider.componentProps = {
                    properties: this.columnsConfigProperties,
                    selected: this.columnDisplayProperties,
                    disabledColumns: ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id']
                }
                this.sideslider.title = this.$t('列表显示属性配置')
                this.sideslider.show = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .options-layout {
        margin-top: 12px;
    }
    .options {
        font-size: 0;
        .option {
            display: inline-block;
            vertical-align: middle;
        }
        .option-collection {
            width: 200px;
        }
        .dropdown-icon {
            margin: 0 -4px;
            display: inline-block;
            vertical-align: middle;
            height: auto;
            top: 0px;
            font-size: 20px;
            &.open {
                top: -1px;
                transform: rotate(180deg);
            }
        }
    }
    .bk-dropdown-list {
        font-size: 14px;
        color: $textColor;
        .bk-dropdown-item {
            position: relative;
            display: block;
            padding: 0 20px;
            margin: 0;
            line-height: 32px;
            cursor: pointer;
            @include ellipsis;
            &:not(.disabled):not(.with-auth):hover {
                background-color: #EAF3FF;
                color: $primaryColor;
            }
            &.disabled {
                color: $textDisabledColor;
                cursor: not-allowed;
            }
            &.with-auth {
                padding: 0;
                span {
                    display: block;
                    padding: 0 20px;
                    &:not(.disabled):hover {
                        background-color: #EAF3FF;
                        color: $primaryColor;
                    }
                    &.disabled {
                        color: $textDisabledColor;
                        cursor: not-allowed;
                    }
                }
            }
        }
    }
    /deep/ {
        .collection-item {
            width: 100%;
            display: flex;
            justify-content: space-between;
            align-items: center;
            &:hover {
                .icon-close {
                    display: block;
                }
            }
            .collection-name {
                @include ellipsis;
            }
            .icon-close {
                display: none;
                color: #979BA5;
                font-size: 20px;
                margin-right: -4px;
            }
        }
    }
</style>
