<template>
    <div class="options-layout clearfix">
        <div class="options fl">
            <bk-button class="option" theme="primary"
                :disabled="!isNormalModuleNode"
                @click="handleAddHost">
                {{$t('新增')}}
            </bk-button>
            <bk-dropdown-menu class="option ml10" trigger="click"
                font-size="large"
                :disabled="!hasSelection"
                @show="isTransferMenuOpen = true"
                @hide="isTransferMenuOpen = false">
                <bk-button slot="dropdown-trigger"
                    :disabled="!hasSelection">
                    <span>{{$t('转移到')}}</span>
                    <i :class="['dropdown-icon bk-icon icon-angle-down',{ 'open': isTransferMenuOpen }]"></i>
                </bk-button>
                <ul class="bk-dropdown-list" slot="dropdown-content">
                    <li :class="['bk-dropdown-item', { disabled: isIdleSet }]"
                        @click="handleTransfer($event, 'idle', isIdleSet)">
                        {{$t('空闲模块')}}
                    </li>
                    <li class="bk-dropdown-item" @click="handleTransfer($event, 'business', false)">{{$t('业务模块')}}</li>
                    <li :class="['bk-dropdown-item', { disabled: !isIdleModule }]"
                        @click="handleTransfer($event, 'resource', !isIdleModule)">
                        {{$t('资源池')}}
                    </li>
                </ul>
            </bk-dropdown-menu>
            <bk-button class="option ml10" @click="handleMultipleEdit">{{$t('编辑')}}</bk-button>
            <cmdb-clipboard-selector class="options-button ml10"
                label-key="bk_property_name"
                :list="clipboardList"
                :disabled="!hasSelection"
                @on-copy="handleCopy">
            </cmdb-clipboard-selector>
            <bk-dropdown-menu class="option ml10" trigger="click"
                font-size="large"
                @show="isMoreMenuOpen = true"
                @hide="isMoreMenuOpen = false">
                <bk-button slot="dropdown-trigger">
                    <span>{{$t('更多')}}</span>
                    <i :class="['dropdown-icon bk-icon icon-angle-down',{ 'open': isMoreMenuOpen }]"></i>
                </bk-button>
                <ul class="bk-dropdown-list" slot="dropdown-content">
                    <li :class="['bk-dropdown-item', { disabled: !hasSelection }]" @click="handleRemove">{{$t('移除')}}</li>
                    <li :class="['bk-dropdown-item', { disabled: !hasSelection }]" @click="handleExport">{{$t('导出')}}</li>
                </ul>
            </bk-dropdown-menu>
        </div>
        <div class="options fr">
            <bk-select class="option option-collection bgc-white"
                ref="collectionSelector"
                v-model="selectedCollection"
                font-size="14"
                :loading="$loading(request.collection)"
                :placeholder="$t('请选择收藏条件')"
                @selected="handleCollectionSelect"
                @clear="handleCollectionClear"
                @toggle="handleCollectionToggle">
                <bk-option v-for="collection in collectionList"
                    :key="collection.id"
                    :id="collection.id"
                    :name="collection.name">
                    <span class="collection-name" :title="collection.name">{{collection.name}}</span>
                    <i class="bk-icon icon-close" @click.stop="handleDeleteCollection(collection)"></i>
                </bk-option>
                <div slot="extension">
                    <a href="javascript:void(0)" class="collection-create" @click="handleCreateCollection">
                        <i class="bk-icon icon-plus-circle"></i>
                        {{$t('新增条件')}}
                    </a>
                </div>
            </bk-select>
            <icon-button class="option ml10" icon="icon-cc-funnel"></icon-button>
            <icon-button class="option ml10" icon="icon-cc-setting" @click="handleSetColumn"></icon-button>
        </div>
        <edit-multiple-host ref="editMultipleHost"
            :properties="hostProperties"
            :selection="$parent.table.selection">
        </edit-multiple-host>
        <cmdb-dialog v-model="dialog.show" v-bind="dialog.props">
            <component
                :is="dialog.component"
                v-bind="dialog.componentProps"
                @cancel="handleDialogCancel">
            </component>
        </cmdb-dialog>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="sideslider.show"
            :width="600"
            :title="$t('列表显示属性配置')">
            <cmdb-columns-config slot="content"
                v-if="sideslider.render"
                :properties="columnsConfigProperties"
                :selected="displayProperties"
                :disabled-columns="['bk_host_innerip', 'bk_cloud_id', 'bk_module_name', 'bk_set_name']"
                @on-cancel="sideslider.show = false"
                @on-apply="handleApplyColumnsConfig"
                @on-reset="handleResetColumnsConfig">
            </cmdb-columns-config>
        </bk-sideslider>
    </div>
</template>

<script>
    import EditMultipleHost from './edit-multiple-host.vue'
    import HostSelector from './host-selector.vue'
    import CmdbColumnsConfig from '@/components/columns-config/columns-config'
    import { mapGetters } from 'vuex'
    import {
        MENU_BUSINESS,
        MENU_BUSINESS_TRANSFER_HOST
    } from '@/dictionary/menu-symbol'
    import Formatter from '@/filters/formatter.js'
    export default {
        components: {
            EditMultipleHost,
            CmdbColumnsConfig,
            [HostSelector.name]: HostSelector
        },
        data () {
            return {
                isTransferMenuOpen: false,
                isMoreMenuOpen: false,
                selectedCollection: '',
                collectionList: [],
                dialog: {
                    show: false,
                    props: {
                        width: 850,
                        showCloseIcon: false
                    },
                    component: null,
                    componentProps: {}
                },
                sideslider: {
                    show: false,
                    render: false
                },
                request: {
                    collection: Symbol('collection')
                }
            }
        },
        computed: {
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
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
            displayProperties () {
                return this.$parent.table.header.map(property => property.bk_property_id)
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
                return this.selectedNode && this.selectedNode.data.default === 1
            },
            isIdleSet () {
                return this.selectedNode && this.selectedNode.data.default !== 0
            },
            clipboardList () {
                return this.$parent.table.header
            }
        },
        created () {
            // this.getCollectionList()
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
                    this.collectionList = result.info
                } catch (e) {
                    this.collectionList = []
                    console.error(e)
                }
            },
            handleCollectionSelect () {

            },
            handleCollectionClear () {

            },
            handleCollectionToggle () {

            },
            handleDeleteCollection () {

            },
            handleCreateCollection () {

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
            handleRemove () {
                if (!this.hasSelection) {
                    return false
                }
                this.$router.push({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        type: 'remove'
                    },
                    query: {
                        sourceModel: this.selectedNode.data.bk_obj_id,
                        sourceId: this.selectedNode.data.bk_inst_id,
                        resources: this.$parent.table.selection.map(item => item.host.bk_host_id).join(',')
                    }
                })
            },
            async handleExport () {
                if (!this.hasSelection) {
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
                        data.append('metadata', JSON.stringify(this.$injectMetadata().metadata))
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
            handleCopy (target) {
                const copyList = this.$parent.table.selection
                const copyText = []
                copyList.forEach(item => {
                    const value = item[target.bk_obj_id][target.bk_property_id]
                    const formattedValue = Formatter(value, target)
                    if (formattedValue !== '--') {
                        copyText.push(formattedValue)
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
            handleDialogCancel () {
                this.dialog.show = false
            },
            handleApplyColumnsConfig (properties) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$route.meta.customInstanceColumn]: properties.map(property => property['bk_property_id'])
                })
                this.sideslider.show = false
            },
            handleResetColumnsConfig () {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$route.meta.customInstanceColumn]: []
                })
            },
            handleSetColumn () {
                this.sideslider.render = true
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
            display: inline-block;
            vertical-align: middle;
            line-height: 19px;
            height: auto;
            top: 0px;
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
            &:hover {
                background-color: #EAF3FF;
                color: $primaryColor;
            }
            &.disabled {
                background-color: #F4F6FA;
                color: $textColor;
                cursor: not-allowed;
            }
        }
    }
    .clipboard-list {
    }
</style>
