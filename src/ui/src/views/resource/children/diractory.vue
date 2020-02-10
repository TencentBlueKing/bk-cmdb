<template>
    <div class="diractory-layout">
        <bk-input class="dir-search" v-model="dirSearch" :placeholder="$t('分组目录')"></bk-input>
        <div class="dir-header">
            <span class="title">{{$t('资源池')}}</span>
            <i class="icon-cc-plus" v-bk-tooltips.top="$t('新建目录')" @click.stop="handleCreateDir"></i>
        </div>
        <ul class="dir-list">
            <li class="dir-item" :class="{ 'selected': curActiveDir === -1 }" @click="handleSearchHost(-1)">
                <i class="icon-cc-memory"></i>
                <span class="dir-name" :title="$t('默认')">{{$t('默认')}}</span>
                <span class="host-count">999+</span>
            </li>
            <li class="dir-item edit-status" v-if="createDir.active">
                <bk-input
                    ref="createdDir"
                    v-click-outside="hanldeCancelCreateDir"
                    class="reset-name"
                    :placeholder="$t('请输入目录名称，回车结束')"
                    v-model="createDir.name"
                    @enter="handleConfirm">
                </bk-input>
            </li>
            <cmdb-auth
                style="display: block;"
                tag="li"
                :auth="$authResources({ type: $OPERATION.C_RESOURCE_HOST })"
                v-for="(dir, index) in dirList"
                :key="index">
                <template slot-scope="{ disabled }">
                    <div
                        class="dir-item"
                        :class="{
                            'edit-status': editDir.id === dir.id,
                            'disabled': disabled,
                            'selected': curActiveDir === dir.id && !disabled
                        }"
                        @click="handleSearchHost(dir.id)">
                        <template v-if="editDir.id === dir.id">
                            <bk-input
                                class="reset-name"
                                v-click-outside="hanldeCancelEdit"
                                :placeholder="$t('请输入目录名称，回车结束')"
                                :ref="`dir-node-${dir.id}`"
                                v-model="editDir.name"
                                @enter="handleConfirm"
                                @click.native.stop>
                            </bk-input>
                        </template>
                        <template v-else>
                            <i class="icon-cc-memory"></i>
                            <span class="dir-name" :title="dir.name">{{dir.name}}</span>
                            <cmdb-dot-menu class="dir-operation" color="#3A84FF" @click.native.stop="handleCloseInput">
                                <div class="dot-content">
                                    <cmdb-auth :auth="$authResources({ type: $OPERATION.C_RESOURCE_HOST })">
                                        <bk-button slot-scope="{ disabled: btnDisabled }"
                                            class="menu-btn"
                                            :disabled="btnDisabled"
                                            :text="true"
                                            @click="handleResetName(dir)">
                                            {{$t('重命名')}}
                                        </bk-button>
                                    </cmdb-auth>
                                    <cmdb-auth :auth="$authResources({ type: $OPERATION.C_RESOURCE_HOST })">
                                        <template slot-scope="{ disabled: btnDisabled }">
                                            <bk-button v-if="true"
                                                class="menu-btn"
                                                :text="true"
                                                :disabled="btnDisabled"
                                                @click="handleRemoveDir">
                                                {{$t('删除')}}
                                            </bk-button>
                                            <span class="menu-btn no-allow-btn" v-else v-bk-tooltips.right="$t('主机不为空，不能删除')">
                                                {{$t('删除')}}
                                            </span>
                                        </template>
                                    </cmdb-auth>
                                </div>
                            </cmdb-dot-menu>
                            <i v-if="disabled" class="icon-cc-lock" v-bk-tooltips.top="$t('无权限')"></i>
                            <span class="host-count" v-else>{{dir.count}}</span>
                        </template>
                    </div>
                </template>
            </cmdb-auth>
        </ul>
    </div>
</template>

<script>
    import Bus from '@/utils/bus.js'
    export default {
        data () {
            return {
                dirSearch: '',
                resetName: false,
                createDir: {
                    active: false,
                    name: ''
                },
                editDir: {
                    id: null,
                    name: ''
                },
                curActiveDir: -1,
                dirList: []
            }
        },
        methods: {
            handleSearchHost (activeId) {
                Bus.$emit('refresh-list')
                this.curActiveDir = activeId
            },
            hanldeCancelCreateDir () {
                this.createDir.active = false
                this.createDir.name = ''
            },
            handleCreateDir () {
                this.createDir.active = true
                this.$nextTick(() => {
                    this.$refs.createdDir.$refs.input.focus()
                })
            },
            hanldeCancelEdit () {
                this.editDir.id = null
                this.editDir.name = ''
            },
            handleConfirm () {
                this.$success(this.$t('新建成功'))
                this.hanldeCancelCreateDir()
            },
            handleResetName (dir) {
                this.editDir.id = dir.id
                this.editDir.name = dir.name
                this.$nextTick(() => {
                    this.$refs[`dir-node-${dir.id}`][0].$refs.input.focus()
                })
            },
            handleCloseInput () {
                this.hanldeCancelCreateDir()
                this.hanldeCancelEdit()
            },
            handleRemoveDir () {
                this.$bkInfo({
                    title: this.$t('确认确定删除目录'),
                    subTitle: this.$t('即将删除目录', { name: 'LOL专用-新入库' }),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: () => {
                        console.log(1)
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .diractory-layout {
        height: 100%;
        overflow: hidden;
        .dir-search {
            padding: 18px 20px 14px;
        }
        .dir-header {
            @include space-between;
            padding: 0 20px;
            height: 42px;
            line-height: 42px;
            background-color: #F0F1F5;
            &:hover {
                background-color: #E1ECFF;
                .icon-cc-plus {
                    background-color: #3A84FF;
                }
            }
            .title {
                font-weight: bold;
                font-size: 14px;
            }
            .icon-cc-plus {
                width: 18px;
                height: 18px;
                line-height: 18px;
                text-align: center;
                color: #FFFFFF;
                background-color: #C4C6CC;
                border-radius: 2px;
                cursor: pointer;
            }
        }
        .dir-list {
            height: calc(100% - 106px);
            padding-bottom: 10px;
            @include scrollbar-y;
        }
        .dir-item {
            display: flex;
            align-items: center;
            height: 36px;
            padding: 0 20px;
            margin: 6px 0;
            cursor: pointer;
            &:first-child {
                margin-top: 0;
            }
            &:not(.edit-status):not(.disabled):hover,
            &:not(.edit-status).selected {
                background-color: #E1ECFF;
                .icon-cc-memory {
                    color: #3A84FF;
                }
                .dir-name {
                    color: #3A84FF;
                }
                .dir-operation {
                    display: block;
                    opacity: 1;
                }
                .host-count {
                    color: #FFFFFF;
                    background-color: #A2C5FD;
                }
            }
            &.disabled {
                .icon-cc-memory {
                    color: #DCDEE5 !important;
                }
                .dir-name {
                    color: #C4C6CC;
                }
            }
            .reset-name {
                width: 100%;
            }
            .icon-cc-memory {
                font-size: 16px;
                margin-right: 10px;
                color: #C4C6CC;
            }
            .dir-name {
                flex: 1;
                font-size: 14px;
                color: #63656E;
                @include ellipsis;
            }
            .dir-operation {
                width: 20px;
                margin-right: 8px;
                opacity: 0;
            }
            .host-count {
                height: 18px;
                line-height: 17px;
                font-size: 12px;
                padding: 0 5px;
                color: #979BA5;
                text-align: center;
                background-color: #F0F1F5;
                border-radius: 2px;
            }
            .icon-cc-lock {
                font-size: 14px;
                color: #C4C6CC;
            }
        }
    }
    .dot-content {
        width: 90px;
        padding: 6px 0;
        .auth-box {
            display: block;
        }
        .menu-btn {
            display: block;
            width: 100%;
            height: 32px;
            line-height: 32px;
            padding: 0 8px;
            text-align: left;
            color: #63656E;
            outline: none;
            &:hover {
                color: #3A84FF;
                background-color: #E1ECFF;
            }
            &:disabled {
                color: #DCDEE5;
                background-color: transparent;
            }
            &.no-allow-btn {
                cursor: not-allowed;
                color: #DCDEE5;
                background-color: transparent;
            }
        }
    }
</style>
