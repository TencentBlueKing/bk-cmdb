<template>
    <div class="directory-layout">
        <div class="directory-options">
            <bk-input class="dir-search"
                v-model.trim="dirSearch"
                clearable
                right-icon="icon-search"
                :placeholder="$t('分组目录')">
            </bk-input>
            <cmdb-auth class="icon-cc-plus"
                :auth="{ type: $OPERATION.C_RESOURCE_DIRECTORY }"
                v-bk-tooltips.top="$t('新建目录')"
                @click="handleShowCreate">
            </cmdb-auth>
        </div>
        <ul class="dir-list" ref="dirList">
            <li
                :class="{
                    'dir-item': true,
                    'dir-item-resource': true,
                    'selected': acitveDirId === null
                }"
                @click="handleResourceClick">
                <span class="dir-name" :title="$t('主机池')">{{$t('主机池')}}</span>
                <span class="host-count">{{totalCount}}</span>
            </li>
            <li class="dir-item edit-status" v-if="createInfo.active" ref="createDirItem">
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
            <li v-for="(dir, index) in filterDirList"
                :key="index"
                :class="{
                    'dir-item': true,
                    'edit-status': editDir.id === dir.bk_module_id,
                    'selected': acitveDirId === dir.bk_module_id,
                    'is-sticky': isSticky(dir)
                }"
                :data-default="dir.default"
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
                    <i class="icon-cc-folder"></i>
                    <span class="dir-name" :title="dir.bk_module_name">{{dir.bk_module_name}}</span>
                    <template v-if="dir.default !== 1">
                        <i :class="['dir-sticky-icon', isSticky(dir) ? 'icon-cc-cancel-sticky' : 'icon-cc-sticky']"
                            v-bk-tooltips.top="isSticky(dir) ? $t('取消置顶') : $t('置顶')"
                            @click="handleToggleSticky(dir)">
                        </i>
                        <cmdb-dot-menu class="dir-operation" color="#3A84FF" @click.native.stop="handleCloseInput">
                            <div class="dot-content">
                                <cmdb-auth :auth="{ type: $OPERATION.U_RESOURCE_DIRECTORY, relation: [dir.bk_module_id] }">
                                    <bk-button slot-scope="{ disabled }"
                                        class="menu-btn"
                                        :text="true"
                                        :disabled="disabled"
                                        @click="handleResetName(dir)">
                                        {{$t('重命名')}}
                                    </bk-button>
                                </cmdb-auth>
                                <cmdb-auth :auth="{ type: $OPERATION.D_RESOURCE_DIRECTORY, relation: [dir.bk_module_id] }">
                                    <div slot-scope="{ disabled }"
                                        v-bk-tooltips.right="{
                                            content: $t('主机不为空，不能删除'),
                                            disabled: !dir.host_count
                                        }">
                                        <bk-button
                                            class="menu-btn"
                                            :text="true"
                                            :disabled="!!dir.host_count || disabled"
                                            @click="handleDelete(dir, index)">
                                            {{$t('删除')}}
                                        </bk-button>
                                    </div>
                                </cmdb-auth>
                            </div>
                        </cmdb-dot-menu>
                    </template>
                    <span class="host-count">{{dir.host_count}}</span>
                </template>
            </li>
        </ul>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import Bus from '@/utils/bus.js'
    import RouterQuery from '@/router/query'
    const CUSTOM_STICKY_KEY = 'sticky-directory'
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
                acitveDirId: null
            }
        },
        computed: {
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('resourceHost', [
                'directoryList'
            ]),
            stickyDirectory () {
                return this.usercustom[CUSTOM_STICKY_KEY] || []
            },
            filterDirList () {
                let list = [...this.directoryList]
                if (this.dirSearch) {
                    const lowerCaseSearch = this.dirSearch.toLowerCase()
                    list = this.directoryList.filter(module => module.bk_module_name.toLowerCase().indexOf(lowerCaseSearch) > -1)
                }
                const count = this.stickyDirectory.length
                list.sort((dirA, dirB) => {
                    const stickyIndexA = this.stickyDirectory.indexOf(dirA.bk_module_id) + 1
                    const stickyIndexB = this.stickyDirectory.indexOf(dirB.bk_module_id) + 1

                    return (stickyIndexA || (count + 1)) - (stickyIndexB || (count + 1))
                })
                return list
            },
            totalCount () {
                return this.directoryList.reduce((accumulator, directory) => {
                    return accumulator + directory.host_count
                }, 0)
            }
        },
        watch: {
            acitveDirId (id) {
                RouterQuery.set({
                    directory: id,
                    page: 1,
                    _t: Date.now()
                })
            }
        },
        async created () {
            Bus.$on('refresh-dir-count', this.getDirectoryList)
            this.getDirectoryList()
        },
        beforeDestroy () {
            Bus.$off('refresh-dir-count', this.getDirectoryList)
        },
        methods: {
            async getDirectoryList () {
                try {
                    const { info } = await this.$store.dispatch('resourceDirectory/getDirectoryList', {
                        params: {
                            page: {
                                sort: 'bk_module_name'
                            }
                        },
                        config: {
                            requestId: 'getDirectoryList'
                        }
                    })
                    this.$store.commit('resourceHost/setDirectoryList', info)
                    let directoryId = RouterQuery.get('directory')
                    if (directoryId) {
                        directoryId = Number(directoryId)
                        const directory = info.find(directory => directory.bk_module_id === directoryId)
                        directory && this.handleSearchHost(directory, false)
                    }
                } catch (error) {
                    console.error(error)
                }
            },
            isSticky (dir) {
                return this.stickyDirectory.includes(dir.bk_module_id)
            },
            async handleToggleSticky (dir) {
                try {
                    const previous = this.stickyDirectory
                    const isSticky = this.isSticky(dir)
                    const current = isSticky ? previous.filter(id => id !== dir.bk_module_id) : [...previous, dir.bk_module_id]
                    await this.$store.dispatch('userCustom/saveUsercustom', {
                        [CUSTOM_STICKY_KEY]: current
                    })
                    this.$success(isSticky ? this.$t('已取消置顶') : this.$t('已置顶'))
                } catch (error) {
                    console.error(error)
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
                    this.$store.commit('resourceHost/addDirectory', newDir)
                    this.$success(this.$t('新建成功'))
                    this.handleCancelCreate()
                    this.handleSearchHost(newDir)
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
                    const target = this.directoryList.find(dir => dir.bk_module_id === this.editDir.id)
                    this.$store.commit('resourceHost/updateDirectory', Object.assign({}, target, {
                        bk_module_id: this.editDir.id,
                        bk_module_name: this.editDir.name
                    }))
                    this.$success(this.$t('修改成功'))
                    this.handleCancelEdit()
                } catch (e) {
                    console.error(e)
                }
            },
            handleSearchHost (active = {}, dispatchEvent = true) {
                this.$store.commit('resourceHost/setActiveDirectory', active)
                this.acitveDirId = active.bk_module_id
                dispatchEvent && Bus.$emit('refresh-resource-list')
            },
            handleResourceClick () {
                this.$store.commit('resourceHost/setActiveDirectory', null)
                this.acitveDirId = null
                Bus.$emit('refresh-resource-list')
            },
            handleCancelCreate () {
                this.createInfo.active = false
                this.createInfo.name = ''
            },
            handleShowCreate () {
                this.createInfo.active = true
                this.$nextTick(() => {
                    const createDirItem = this.$refs.createDirItem
                    const idleNextItem = this.$refs.dirList.querySelector('[data-default="1"]').nextElementSibling
                    this.$refs.dirList.insertBefore(createDirItem, idleNextItem)
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
            async handleDelete (dir, index) {
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
                            if (dir.bk_module_id === this.acitveDirId) {
                                this.handleSearchHost(this.filterDirList[index - 1])
                            }
                            this.$store.commit('resourceHost/deleteDirectory', dir.bk_module_id)
                            this.$success(this.$t('删除成功'))
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .directory-layout {
        height: 100%;
        overflow: hidden;
        .directory-options {
            display: flex;
            align-items: center;
            padding: 18px 20px 14px;
            .dir-search {
                flex: 1;
                display: block;
                width: auto;
            }
            .icon-cc-plus {
                flex: 20px 0 0;
                font-size: 20px;
                margin-left: 10px;
                cursor: pointer;
                &:hover {
                    color: $primaryColor;
                }
                &.disabled {
                    color: $textDisabledColor;
                }
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
            margin: 1px 0;
            cursor: pointer;
            &.dir-item-resource {
                background-color: #F0F1F5;
            }
            &.is-sticky {
                background-color: #F0F1F5;
            }
            &:first-child {
                margin-top: 0;
            }
            &:not(.edit-status):not(.disabled):hover,
            &:not(.edit-status).selected {
                background-color: #E1ECFF;
                .icon-cc-folder {
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
                .dir-sticky-icon {
                    display: inline-block;
                }
            }
            &.disabled {
                .icon-cc-folder {
                    color: #DCDEE5 !important;
                }
                .dir-name {
                    color: #C4C6CC;
                }
            }
            .icon-cc-folder {
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
                display: flex;
                width: 28px;
                height: 28px;
                line-height: 28px;
                align-items: center;
                justify-content: center;
                margin-right: 8px;
                opacity: 0;
                border-radius: 50%;
                &:hover {
                    background-color: #fff;
                }
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
            .dir-sticky-icon {
                display: none;
                width: 28px;
                height: 28px;
                margin: 0 0 0 5px;
                line-height: 28px;
                text-align: center;
                color: $primaryColor;
                border-radius: 50%;
                &:hover {
                    background-color: #fff;
                }
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
