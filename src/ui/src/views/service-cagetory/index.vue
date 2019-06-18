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
                        :style="{ width: editMainStatus === mainCagetory['id'] ? '100%' : 'auto' }"
                        v-if="editMainStatus === mainCagetory['id']">
                        <cagetory-input
                            ref="editInput"
                            :input-ref="'cagetoryInput'"
                            :set-style="{ border: 'none', outline: 'none', padding: 0, 'background-color': 'transparent !important' }"
                            :placeholder="$t('ServiceCagetory[\'请输入一级分类\']')"
                            name="cagetoryName"
                            v-validate="'required|namedCharacter'"
                            v-model="mainCagetoryName"
                            @on-confirm="handleEditCagetory(mainCagetory['id'], mainCagetory['name'], 'main')"
                            @on-cancel="handleCloseEditMain">
                        </cagetory-input>
                    </div>
                    <template v-else>
                        <div class="cagetory-name">
                            <span>{{mainCagetory['name']}}</span>
                            <i class="property-edit icon-cc-edit-shape" @click.stop="handleEditMain(mainCagetory['id'], mainCagetory['name'])"></i>
                        </div>
                        <cmdb-dot-menu class="dot-menu">
                            <div class="menu-operational">
                                <i @click="handleShowAddChild(mainCagetory['id'])">{{$t("ServiceCagetory['添加二级分类']")}}</i>
                                <i class="not-allowed" v-if="mainCagetory['child_cagetory_list'].length || mainCagetory['is_built_in']">{{$t("Common['删除']")}}</i>
                                <i v-else @click="handleDeleteCagetory(mainCagetory['id'])">{{$t("Common['删除']")}}</i>
                            </div>
                        </cmdb-dot-menu>
                    </template>
                </div>
                <div class="child-cagetory">
                    <div class="child-item child-edit" v-if="addChildStatus === mainCagetory['id']">
                        <cagetory-input
                            ref="editInput"
                            :input-ref="'cagetoryInput'"
                            :placeholder="$t('ServiceCagetory[\'请输入二级分类\']')"
                            :edit-id="mainCagetory['bk_root_id']"
                            name="cagetoryName"
                            v-validate="'required|namedCharacter'"
                            v-model="cagetoryName"
                            @on-confirm="handleAddCagetory"
                            @on-cancel="handleCloseAddChild">
                        </cagetory-input>
                    </div>
                    <div :class="['child-item', editChildStatus === childCagetory['id'] ? 'child-edit' : '']"
                        v-for="(childCagetory, childIndex) in mainCagetory['child_cagetory_list']"
                        :key="childIndex">
                        <cagetory-input
                            v-if="editChildStatus === childCagetory['id']"
                            ref="editInput"
                            :input-ref="'cagetoryInput'"
                            :placeholder="$t('ServiceCagetory[\'请输入二级分类\']')"
                            name="cagetoryName"
                            v-validate="'required|namedCharacter'"
                            v-model="childCagetoryName"
                            @on-confirm="handleEditCagetory(childCagetory['id'], childCagetory['name'], 'child')"
                            @on-cancel="handleCloseEditChild">
                        </cagetory-input>
                        <template v-else>
                            <div class="child-title">
                                <span>{{childCagetory['name']}}</span>
                                <div class="child-edit" v-if="!childCagetory['is_built_in']">
                                    <i class="property-edit icon-cc-edit-shape mr10"
                                        @click.stop="handleEditChild(childCagetory['id'], childCagetory['name'])">
                                    </i>
                                    <i class="icon-cc-tips-close"
                                        v-if="!childCagetory['usage_amount']"
                                        @click.stop="handleDeleteCagetory(childCagetory['id'])">
                                    </i>
                                    <i class="icon-cc-tips-close" v-else
                                        style="color: #c4c6cc; cursor: not-allowed;"
                                        v-bktooltips="tooltips">
                                    </i>
                                </div>
                            </div>
                            <!-- <span>{{childCagetory['usage_amount']}}</span> -->
                        </template>
                    </div>
                </div>
            </div>
            <div class="cagetory-item add-item" :style="{ 'border-style': showAddMianCagetory ? 'solid' : 'dashed' }">
                <div class="cagetory-title" :style="{ 'border-bottom-style': showAddMianCagetory ? 'solid' : 'dashed' }">
                    <div class="main-edit" style="width: 100%;" v-if="showAddMianCagetory">
                        <cagetory-input
                            ref="addCagetoryInput"
                            :input-ref="'cagetoryInput'"
                            :set-style="{ border: 'none', outline: 'none', padding: 0, 'background-color': 'transparent !important' }"
                            :placeholder="$t('ServiceCagetory[\'请输入一级分类\']')"
                            name="cagetoryName"
                            v-validate="'required|namedCharacter'"
                            v-model="cagetoryName"
                            @on-confirm="handleAddCagetory"
                            @on-cancel="handleCloseAddBox">
                        </cagetory-input>
                    </div>
                </div>
                <div class="child-cagetory"></div>
                <span class="add-box" v-if="!showAddMianCagetory" @click="handleAddBox"></span>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import cagetoryInput from './children/cagetory-input'
    export default {
        components: {
            featureTips,
            cagetoryInput
        },
        data () {
            return {
                tooltips: {
                    content: this.$t("ServiceCagetory['二级分类删除提示']"),
                    arrowsSize: 5
                },
                showFeatureTips: false,
                showAddMianCagetory: false,
                showAddChildCagetory: false,
                editMainStatus: null,
                editChildStatus: null,
                addChildStatus: null,
                cagetoryName: '',
                mainCagetoryName: '',
                childCagetoryName: '',
                list: []
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['cagetory']
            this.getCagetoryList()
        },
        methods: {
            ...mapActions('serviceClassification', [
                'searchServiceCategory',
                'createServiceCategory',
                'updateServiceCategory',
                'deleteServiceCategory'
            ]),
            getCagetoryList () {
                this.searchServiceCategory({
                    params: this.$injectMetadata({})
                }).then((data) => {
                    const cagetoryList = data.info.map(item => {
                        return {
                            usage_amount: item['usage_amount'],
                            ...item['category']
                        }
                    })
                    const list = cagetoryList.filter(cagetory => !cagetory.hasOwnProperty('bk_parent_id'))
                    this.list = list.map(mainCagetory => {
                        return {
                            ...mainCagetory,
                            child_cagetory_list: cagetoryList.filter(cagetory => cagetory['bk_parent_id'] === mainCagetory['id'])
                        }
                    })
                })
            },
            createdCagetory (name, rootId) {
                this.createServiceCategory({
                    params: this.$injectMetadata({
                        bk_root_id: rootId,
                        bk_parent_id: rootId,
                        name
                    })
                }).then(() => {
                    this.$bkMessage({
                        message: this.$t("Common['保存成功']"),
                        theme: 'success'
                    })
                    this.showAddMianCagetory = false
                    this.handleCloseAddChild()
                    this.getCagetoryList()
                })
            },
            async handleAddCagetory (name, bk_root_id = 0) {
                if (!await this.$validator.validateAll()) {
                    this.$bkMessage({
                        message: this.errors.first('cagetoryName') || this.$t("ServiceCagetory['请输入分类名称']"),
                        theme: 'error'
                    })
                } else {
                    this.createdCagetory(name, bk_root_id)
                }
            },
            async handleEditCagetory (id, name, type) {
                if (!await this.$validator.validateAll()) {
                    this.$bkMessage({
                        message: this.errors.first('cagetoryName') || this.$t("ServiceCagetory['请输入分类名称']"),
                        theme: 'error'
                    })
                } else if (name === this.mainCagetoryName || name === this.childCagetoryName) {
                    this.handleCloseEditChild()
                    this.handleCloseEditMain()
                } else {
                    this.updateServiceCategory({
                        params: this.$injectMetadata({
                            id,
                            name: type === 'main' ? this.mainCagetoryName : this.childCagetoryName
                        })
                    }).then(() => {
                        this.$bkMessage({
                            message: this.$t("Common['保存成功']"),
                            theme: 'success'
                        })
                        this.handleCloseEditChild()
                        this.handleCloseEditMain()
                        this.getCagetoryList()
                    })
                }
            },
            handleDeleteCagetory (id) {
                this.$bkInfo({
                    title: this.$t("ServiceCagetory['确认删除分类']"),
                    confirmFn: async () => {
                        await this.deleteServiceCategory({
                            params: {
                                data: this.$injectMetadata({ id })
                            },
                            config: {
                                requestId: 'delete_proc_services_category'
                            }
                        }).then(() => {
                            this.$bkMessage({
                                message: this.$t("Common['删除成功']"),
                                theme: 'success'
                            })
                            this.getCagetoryList()
                        })
                    }
                })
            },
            handleEditMain (id, name) {
                this.editMainStatus = id
                this.mainCagetoryName = name
                this.handleCloseEditChild()
                this.handleCloseAddChild()
                this.handleCloseAddBox()
                this.$nextTick(() => {
                    this.$refs.editInput[0].$refs.cagetoryInput.focus()
                })
            },
            handleCloseEditMain () {
                this.editMainStatus = null
            },
            handleEditChild (id, name) {
                this.editChildStatus = id
                this.childCagetoryName = name
                this.handleCloseAddChild()
                this.handleCloseEditMain()
                this.handleCloseAddBox()
                this.$nextTick(() => {
                    this.$refs.editInput[0].$refs.cagetoryInput.focus()
                })
            },
            handleCloseEditChild () {
                this.editChildStatus = null
            },
            handleAddBox () {
                this.showAddMianCagetory = true
                this.$nextTick(() => {
                    this.$refs.addCagetoryInput.$refs.cagetoryInput.focus()
                })
            },
            handleCloseAddBox () {
                this.showAddMianCagetory = false
                this.cagetoryName = ''
            },
            handleShowAddChild (id) {
                this.addChildStatus = id
                this.$nextTick(() => {
                    this.$refs.editInput[0].$refs.cagetoryInput.focus()
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
        min-width: 1442px;
        .cagetory-list {
            display: flex;
            flex-flow: row wrap;
        }
        .cagetory-item {
            position: relative;
            min-width: 320px;
            flex: 0 0 22%;
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
            }
            .cagetory-name {
                @include ellipsis;
                flex: 1;
                padding-right: 20px;
                .icon-cc-edit-shape {
                    display: none;
                    cursor: pointer;
                    color: #3a84ff;
                }
                &:hover .icon-cc-edit-shape {
                    display: inline !important;
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
                    font-size: 14px;
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
                .edit-box {
                    width: 100%;
                }
            }
        }
    }
    .menu-operational {
        padding: 6px 0;
        line-height: 32px;
        color: #63656e;
        i {
            display: block;
            font-style: normal;
            padding: 0 8px;
            cursor: pointer;
            &:hover {
                color: #3a84ff;
                background-color: #e1ecff;
            }
            &.not-allowed {
                color: #c4c6cc;
                background-color: transparent;
                cursor: not-allowed;
            }
        }
    }
    @media screen and (min-width: 1920px){
        .cagetory-wrapper {
            min-width: 1650px;
            .cagetory-item {
                min-width: auto;
                flex: 0 0 18%;
            }
        }
    }
</style>
