<template>
    <div class="diractory-layout">
        <bk-input class="dir-search"
            v-model.trim="dirSearch"
            :placeholder="$t('分组目录')">
        </bk-input>
        <div class="dir-header">
            <span class="title">{{$t('资源池')}}</span>
            <i class="icon-cc-plus" v-bk-tooltips.top="$t('新建目录')" @click.stop="handleShowCreate"></i>
        </div>
        <ul class="dir-list">
            <li class="dir-item edit-status" v-if="createInfo.active">
                <bk-input
                    ref="createdDir"
                    v-click-outside="handleCancelCreate"
                    style="width: 100%"
                    v-validate="'required|singlechar|length:256'"
                    :placeholder="$t('请输入目录名称，回车结束')"
                    v-model="createInfo.name"
                    @enter="handleConfirm(true)">
                </bk-input>
            </li>
            <cmdb-auth
                style="display: block;"
                tag="li"
                :auth="$authResources({ type: $OPERATION.C_RESOURCE_HOST })"
                v-for="(dir, index) in filterDirList"
                :key="index">
                <template slot-scope="{ disabled }">
                    <div
                        class="dir-item"
                        :class="{
                            'edit-status': editDir.id === dir.bk_module_id,
                            'disabled': disabled,
                            'selected': acitveDirId === dir.bk_module_id && !disabled
                        }"
                        @click="handleSearchHost(dir)">
                        <template v-if="editDir.id === dir.bk_module_id">
                            <bk-input
                                style="width: 100%"
                                v-click-outside="handleCancelEdit"
                                v-validate="'required|singlechar|length:256'"
                                :placeholder="$t('请输入目录名称，回车结束')"
                                :ref="`dir-node-${dir.bk_module_id}`"
                                v-model="editDir.name"
                                @enter="handleConfirm(false)"
                                @click.native.stop>
                            </bk-input>
                        </template>
                        <template v-else>
                            <i class="icon-cc-memory"></i>
                            <span class="dir-name" :title="dir.bk_module_name">{{dir.bk_module_name}}</span>
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
                                                @click="handleDelete(dir)">
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
                            <span class="host-count" v-else>{{dir.host_count}}</span>
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
                createInfo: {
                    active: false,
                    name: ''
                },
                editDir: {
                    id: null,
                    name: ''
                },
                acitveDirId: -1,
                dirList: [],
                defaultDir: {
                    bk_module_id: -1,
                    bk_module_name: '默认',
                    host_count: 0
                }
            }
        },
        computed: {
            filterDirList () {
                if (this.dirSearch) {
                    return this.dirList.filter(module => module.bk_module_name.indexOf(this.dirSearch) > -1)
                }
                return [...this.dirList]
            }
        },
        created () {
            Bus.$on('refresh-dir-count', this.refreshCount)
            this.getDirectoryList()
        },
        beforeDestroy () {
            Bus.$off('refresh-dir-count', this.refreshCount)
        },
        methods: {
            async getDirectoryList () {
                try {
                    const data = await this.$store.dispatch('resourceDirectory/getDirectoryList', {
                        params: {
                            page: {
                                sort: '-bk_module_id'
                            }
                        },
                        config: {
                            requestId: 'getDirectoryList'
                        }
                    })
                    this.dirList = data.info || []
                    const firstDir = this.dirList[0] || {}
                    this.acitveDirId = firstDir.bk_module_id
                    this.$store.commit('resourceHost/setActiveDirectory', firstDir)
                } catch (e) {
                    console.error(e)
                    this.dirList = []
                }
            },
            async createdDir () {
                try {
                    const data = await this.$store.dispatch('resourceDirectory/createDirectory', {
                        params: {
                            bk_module_name: this.createInfo.name
                        }
                    })
                    const newDir = {
                        bk_module_id: data.created.id,
                        bk_module_name: this.createInfo.name,
                        host_count: 0
                    }
                    this.dirList.splice(1, 0, newDir)
                    this.$success(this.$t('新建成功'))
                    this.handleCancelCreate()
                } catch (e) {
                    console.error(e)
                }
            },
            async updateDir () {
                try {
                    await this.$store.dispatch('resourceDirectory/updateDirectory', {
                        moduleId: this.editDir.id,
                        params: {
                            bk_module_name: this.editDir.name
                        },
                        config: {
                            requestId: 'updateDir'
                        }
                    })
                    const index = this.dirList.findIndex(dir => dir.bk_module_id === this.editDir.id)
                    this.$set(this.dirList, index, Object.assign(this.dirList[index], {
                        bk_module_id: this.editDir.id,
                        bk_module_name: this.editDir.name
                    }))
                    this.$success(this.$t('修改成功'))
                    this.handleCancelEdit()
                } catch (e) {
                    console.error(e)
                }
            },
            handleSearchHost (active = {}) {
                this.$store.commit('resourceHost/setActiveDirectory', active)
                Bus.$emit('refresh-resource-list')
                this.acitveDirId = active.bk_module_id
            },
            handleCancelCreate () {
                this.createInfo.active = false
                this.createInfo.name = ''
            },
            handleShowCreate () {
                this.createInfo.active = true
                this.$nextTick(() => {
                    this.$refs.createdDir.$refs.input.focus()
                })
            },
            handleCancelEdit () {
                this.editDir.id = null
                this.editDir.name = ''
            },
            async handleConfirm (isCreate) {
                if (!await this.$validator.validateAll()) {
                    this.$error(this.$t('请正确目录名称'))
                    return
                }
                if (isCreate) {
                    this.createdDir()
                } else {
                    this.updateDir()
                }
            },
            handleResetName (dir) {
                this.editDir.id = dir.bk_module_id
                this.editDir.name = dir.bk_module_name
                this.$nextTick(() => {
                    this.$refs[`dir-node-${dir.bk_module_id}`][0].$refs.input.focus()
                })
            },
            handleCloseInput () {
                this.handleCancelCreate()
                this.handleCancelEdit()
            },
            async handleDelete (dir) {
                if (dir.host_count) {
                    this.$error(this.$t('目标包含主机, 不允许删除'))
                    return
                }
                this.$bkInfo({
                    title: this.$t('确认确定删除目录'),
                    subTitle: this.$t('即将删除目录', { name: dir.bk_module_name }),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('resourceDirectory/deleteDirectory', {
                                moduleId: dir.bk_module_id,
                                config: {
                                    requestId: 'deleteDirectory'
                                }
                            })
                            const index = this.dirList.findIndex(target => target.bk_module_id === dir.bk_module_id)
                            this.dirList.splice(index, 1)
                            if (dir.bk_module_id === this.acitveDirId) {
                                this.acitveDirId = this.defaultDir.bk_module_id
                            }
                            this.$success(this.$t('删除成功'))
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            refreshCount ({ reduceId, addId, count }) {
                this.dirList = this.dirList.map((dir, index) => {
                    if (dir.bk_module_id === reduceId) {
                        dir.host_count -= count
                    } else if (dir.bk_module_id === addId) {
                        dir.host_count += count
                    }
                    return dir
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
