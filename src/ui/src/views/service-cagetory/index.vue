<template>
    <div class="cagetory-wrapper" :style="{ 'padding-top': showFeatureTips ? '10px' : '' }">
        <feature-tips
            :feature-name="'cagetory'"
            :show-tips="showFeatureTips"
            :desc="$t('ServiceCagetory[\'服务分类功能提示\']')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="cagetory-list">
            <div class="cagetory-item" v-for="(mainCagetory, index) in list" :key="index">
                <div class="cagetory-title" :style="{ 'background-color': mainCagetory['editStatus'] ? '#f0f1f5' : '' }">
                    <div class="main-edit"
                        v-if="editMainStatus === mainCagetory['id']"
                        v-click-outside="handleCloseEditMain">
                        <input type="text" ref="editInput"
                            :placeholder="$t('ServiceCagetory[\'请输入一级分类\']')"
                            v-model="mainCagetoryName"
                            @keypress.enter="handleEditCagetory(cagetoryName)">
                    </div>
                    <template v-else>
                        <div class="cagetory-name">
                            <span>{{mainCagetory['name']}}</span>
                            <i class="property-edit icon-cc-edit" @click.stop="handleEditMain(mainCagetory['id'], mainCagetory['name'])"></i>
                        </div>
                        <cmdb-dot-menu class="dot-menu">
                            <div class="menu-operational">
                                <i @click="handleShowAddChild(mainCagetory['id'])">{{$t("Common['添加']")}}</i>
                                <i style="cursor: not-allowed;" v-if="mainCagetory['child_cagetory_list'].length || mainCagetory['is_built_in']">{{$t("Common['删除']")}}</i>
                                <i v-else @click.stop="handleDeleteCagetory(mainCagetory['id'])">{{$t("Common['删除']")}}</i>
                            </div>
                        </cmdb-dot-menu>
                    </template>
                </div>
                <div class="child-cagetory">
                    <div class="child-item child-edit" v-if="addChildStatus === mainCagetory['id']">
                        <div class="edit-box clearfix" v-click-outside="handleCloseAddChild">
                            <input type="text"
                                ref="editInput"
                                class="bk-form-input"
                                :placeholder="$t('ServiceCagetory[\'请输入二级分类\']')"
                                v-model="cagetoryName">
                            <span class="text-primary btn-confirm"
                                @click.stop="handleAddCagetory(cagetoryName, mainCagetory['root_id'])">{{$t("Common['确定']")}}
                            </span>
                            <span class="text-primary" @click="handleCloseAddChild">{{$t("Common['取消']")}}</span>
                        </div>
                    </div>
                    <div class="child-item" v-if="!mainCagetory['child_cagetory_list'].length && editChildStatus !== mainCagetory['id']">
                        <div class="child-title" style="color: #dcdee5 !important; background-color: transparent !important;">
                            <span>{{$t("ServiceCagetory['二级分类']")}}</span>
                        </div>
                    </div>
                    <div :class="['child-item', editChildStatus === childCagetory['id'] ? 'child-edit' : '']" v-else
                        v-for="(childCagetory, childIndex) in mainCagetory['child_cagetory_list']"
                        :key="childIndex">
                        <div class="edit-box clearfix"
                            v-if="editChildStatus === childCagetory['id']"
                            v-click-outside="handleCloseEditChild">
                            <input type="text"
                                ref="editInput"
                                class="bk-form-input"
                                :placeholder="$t('ServiceCagetory[\'请输入二级分类\']')"
                                v-model="childCagetoryName">
                            <span class="text-primary btn-confirm"
                                @click.stop="handleEditCagetory(childCagetory['name'], childCagetory['parent_id'])">{{$t("Common['确定']")}}
                            </span>
                            <span class="text-primary" @click="handleCloseEditChild">{{$t("Common['取消']")}}</span>
                        </div>
                        <template v-else>
                            <div class="child-title">
                                <span>{{childCagetory['name']}}</span>
                                <div class="child-edit" v-if="!childCagetory['is_built_in']">
                                    <i class="property-edit icon-cc-edit mr10"
                                        @click.stop="handleEditChild(childCagetory['id'], childCagetory['name'])">
                                    </i>
                                    <i class="icon-cc-tips-close" @click.stop="handleDeleteCagetory(mainCagetory['id'])"></i>
                                </div>
                            </div>
                            <span>8</span>
                        </template>
                    </div>
                </div>
            </div>
            <div class="cagetory-item add-item">
                <div class="cagetory-title" :style="{ 'background-color': showAddMianCagetory ? '#f0f1f5' : '' }">
                    <div class="main-edit" v-if="showAddMianCagetory">
                        <input type="text"
                            ref="addCagetoryInput"
                            :placeholder="$t('ServiceCagetory[\'请输入一级分类\']')"
                            v-model="cagetoryName"
                            v-click-outside="handleCloseAddBox"
                            @keypress.enter="handleAddCagetory(cagetoryName)">
                    </div>
                    <template v-else>
                        <div class="cagetory-name">
                            <span>{{$t("ServiceCagetory['一级分类']")}}</span>
                        </div>
                    </template>
                </div>
                <div class="child-cagetory">
                    <div class="child-item">
                        <div class="child-title">
                            <span>{{$t("ServiceCagetory['二级分类']")}}</span>
                        </div>
                    </div>
                </div>
                <span class="add-box" v-if="!showAddMianCagetory" @click="handleAddBox"></span>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    export default {
        components: {
            featureTips
        },
        data () {
            return {
                showFeatureTips: false,
                showAddMianCagetory: false,
                showAddChildCagetory: false,
                editMainStatus: null,
                editChildStatus: null,
                addChildStatus: null,
                cagetoryName: '',
                mainCagetoryName: '',
                childCagetoryName: '',
                list: [],
                originList: []
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Nav["服务分类"]'))
            this.showFeatureTips = this.featureTipsParams['cagetory']
            this.getCagetoryList()
        },
        methods: {
            ...mapActions('serviceClassification', [
                'searchServiceCategory',
                'createServiceCategory',
                'deleteServiceCategory'
            ]),
            getCagetoryList () {
                this.searchServiceCategory({
                    params: this.$injectMetadata({})
                }).then((data) => {
                    this.originList = data.info
                    const list = data.info.filter(cagetory => !cagetory.hasOwnProperty('parent_id'))
                    this.list = list.map(mainCagetory => {
                        return {
                            ...mainCagetory,
                            child_cagetory_list: data.info.filter(cagetory => cagetory['parent_id'] === mainCagetory['id'])
                        }
                    })
                })
            },
            createdCagetory (name, rootId) {
                this.createServiceCategory({
                    params: this.$injectMetadata({
                        root_id: rootId,
                        parent_id: rootId,
                        name
                    })
                }).then(() => {
                    this.showAddMianCagetory = false
                    this.handleCloseAddChild()
                    this.getCagetoryList()
                })
            },
            handleAddCagetory (name, root_id = 0) {
                if (!name) {
                    this.$bkMessage({
                        message: '请输入分类名称',
                        theme: 'error'
                    })
                } else {
                    this.createdCagetory(name, root_id)
                }
            },
            handleEditCagetory (name, root_id = 0) {

            },
            handleDeleteCagetory (id) {
                this.$bkInfo({
                    title: '确认删除分类?',
                    confirmFn: async () => {
                        await this.deleteServiceCategory({
                            params: {
                                data: this.$injectMetadata({ id })
                            },
                            config: {
                                requestId: 'delete_proc_services_category'
                            }
                        })
                        this.getCagetoryList()
                    }
                })
            },
            handleEditMain (id, name) {
                this.editMainStatus = id
                this.mainCagetoryName = name
                this.$nextTick(() => {
                    this.$refs.editInput[0].focus()
                })
            },
            handleCloseEditMain () {
                this.editMainStatus = null
            },
            handleEditChild (id, name) {
                this.editChildStatus = id
                this.childCagetoryName = name
                this.$nextTick(() => {
                    this.$refs.editInput[0].focus()
                })
            },
            handleCloseEditChild () {
                this.editChildStatus = null
            },
            handleAddBox () {
                this.showAddMianCagetory = true
                this.$nextTick(() => {
                    this.$refs.addCagetoryInput.focus()
                })
            },
            handleCloseAddBox () {
                this.showAddMianCagetory = false
                this.cagetoryName = ''
            },
            handleShowAddChild (id) {
                this.addChildStatus = id
                this.$nextTick(() => {
                    this.$refs.editInput[0].focus()
                })
            },
            handleCloseAddChild () {
                this.addChildStatus = null
                this.cagetoryName = ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cagetory-wrapper {
        .cagetory-list {
            display: flex;
            flex-wrap: wrap;
        }
        .cagetory-item {
            position: relative;
            width: 320px;
            border: 1px solid #dcdee5;
            margin-right: 30px;
            margin-bottom: 20px;
            &.add-item {
                .cagetory-name {
                    color: #dcdee5 !important;
                }
                .child-title {
                    color: #dcdee5 !important;
                    background-color: transparent !important;
                }
                .add-box {
                    position: absolute;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    cursor: pointer;
                    &::after, &::before {
                        content: '';
                        position: absolute;
                        top: 50%;
                        left: 50%;
                        width: 20px;
                        height: 3px;
                        background-color: #3a84ff;
                        transform: translate(-50%, -50%);
                    }
                    &::before {
                        width: 3px;
                        height: 20px;
                    }
                }
            }
        }
        .cagetory-title {
            @include space-between;
            background-color: #fafbfd;
            padding: 0 20px 0 16px;
            height: 42px;
            line-height: 42px;
            font-size: 14px;
            color: #63656e;
            font-weight: bold;
            border-bottom: 1px solid #dcdee5;
            &:hover {
                background-color: #f0f1f5;
                .icon-cc-edit {
                    display: inline !important;
                }
            }
            .main-edit {
                display: flex;
                align-items: center;
                input {
                    width: 240px;
                    height: 42px;
                    line-height: 42px;
                    color: #63656e;
                    background-color: transparent;
                    border: none;
                    padding-left: 10px;
                    outline: none;
                    font-weight: normal;
                }
                &::before {
                    content: '';
                    display: block;
                    width: 2px;
                    height: 20px;
                    background-color: #63656e;
                }
            }
            .cagetory-name {
                @include ellipsis;
                flex: 1;
                padding-right: 20px;
                .icon-cc-edit {
                    display: none;
                    cursor: pointer;
                    color: #3a84ff;
                }
            }
            .dot-menu {
                cursor: pointer;
            }
        }
        .child-cagetory {
            height: 280px;
            padding: 0 10px 10px 38px;
            overflow: hidden;
            &:hover {
                @include scrollbar-y;
            }
            .child-item {
                @include space-between;
                position: relative;
                z-index: 10;
                line-height: 32px;
                &.child-edit {
                    .bk-form-input {
                        float: left;
                        font-size: 12px;
                        width: 170px;
                        height: 32px;
                        margin-right: 4px;
                    }
                    .edit-box .text-primary {
                        display: inline-block;
                        line-height: normal;
                        font-size: 12px;
                        &.btn-confirm {
                            position: relative;
                            margin-right: 6px;
                            &::after {
                                content: '';
                                position: absolute;
                                top: 2px;
                                right: -6px;
                                display: inline-block;
                                width: 1px;
                                height: 14px;
                                background-color: #dcdee5;
                            }
                        }
                    }
                    &:first-child::after {
                        height: 32px;
                    }
                }
                &:hover {
                    .child-title {
                        padding-right: 10px;
                        background-color: #fafbfd;
                        color: #3a84ff;
                    }
                    >span {
                        display: none;
                    }
                    .child-edit {
                        display: block;
                    }
                }
                &:first-child {
                    padding-top: 14px;
                    &::after {
                        height: 30px;
                        top: 0px;
                    }
                }
                &::after {
                    content: '';
                    position: absolute;
                    top: -15px;
                    left: -20px;
                    display: block;
                    width: 30px;
                    height: 32px;
                    border-bottom: 1px solid #dcdee5;
                    border-left: 1px solid #dcdee5;
                    z-index: -1;
                }
                .child-title {
                    @include ellipsis;
                    @include space-between;
                    color: #63656e;
                    flex: 1;
                    padding-right: 20px;
                    padding-left: 16px;
                    margin-left: 10px;
                    span {
                        @include ellipsis;
                        padding-right: 10px;
                    }
                }
                >span {
                    color: #c4c6cc;
                    padding-right: 18px;
                }
                .child-edit {
                    display: none;
                    i {
                        font-size: 14px;
                        color: #3a84ff;
                        cursor: pointer;
                    }
                }
            }
        }
    }
    .menu-operational {
        width: 68px;
        padding: 6px 0;
        font-size: 12px;
        text-align: center;
        line-height: 32px;
        color: #c4c6cc;
        i {
            display: block;
            font-style: normal;
            cursor: pointer;
            &:hover {
                color: #3a84ff;
                background-color: #e1ecff;
            }
        }
    }
</style>
