<template>
    <div class="options-layout clearfix">
        <div class="options fl">
            <bk-button class="option" theme="primary">{{$t('新增')}}</bk-button>
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
                    <li :class="['bk-dropdown-item', { disabled: !isIdleSet }]"
                        @click="handleTransfer($event, 'idle', !isIdleSet)">
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
            <bk-dropdown-menu class="option ml10" trigger="click"
                font-size="large"
                @show="isMoreMenuOpen = true"
                @hide="isMoreMenuOpen = false">
                <bk-button slot="dropdown-trigger">
                    <span>{{$t('更多')}}</span>
                    <i :class="['dropdown-icon bk-icon icon-angle-down',{ 'open': isMoreMenuOpen }]"></i>
                </bk-button>
                <ul class="bk-dropdown-list" slot="dropdown-content">
                    <li :class="['bk-dropdown-item', { disabled: !hasSelection }]">{{$t('移除')}}</li>
                    <li :class="['bk-dropdown-item', { disabled: !hasSelection }]">{{$t('导出')}}</li>
                    <li :class="['bk-dropdown-item', { disabled: !hasSelection }]">{{$t('复制')}}</li>
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
            <icon-button class="option ml10" icon="icon-cc-setting"></icon-button>
        </div>
        <edit-multiple-host ref="editMultipleHost"
            :properties="hostProperties"
            :selection="$parent.table.selection">
        </edit-multiple-host>
    </div>
</template>

<script>
    import EditMultipleHost from './edit-multiple-host.vue'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            EditMultipleHost
        },
        data () {
            return {
                isTransferMenuOpen: false,
                isMoreMenuOpen: false,
                selectedCollection: '',
                collectionList: [],
                request: {
                    collection: Symbol('collection')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', [
                'getProperties',
                'currentNode'
            ]),
            hostProperties () {
                return this.getProperties('host')
            },
            hasSelection () {
                return !!this.$parent.table.selection.length
            },
            isIdleModule () {
                return this.currentNode && this.currentNode.data.default === 1
            },
            isIdleSet () {
                return this.currentNode && this.currentNode.data.default !== 0
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
</style>
