<template>
    <div class="category-wrapper" :style="{ 'padding-top': showFeatureTips ? '10px' : '' }">
        <feature-tips
            :feature-name="'category'"
            :show-tips="showFeatureTips"
            :desc="$t('ServiceCategory[\'服务分类功能提示\']')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="category-list">
            <div class="category-item" v-for="(mainCategory, index) in list" :key="index">
                <div class="category-title" :style="{ 'background-color': mainCategory['editStatus'] ? '#f0f1f5' : '' }">
                    <div class="main-edit"
                        :style="{ width: editMainStatus === mainCategory['id'] ? '100%' : 'auto' }"
                        v-if="editMainStatus === mainCategory['id']">
                        <category-input
                            ref="editInput"
                            :input-ref="'categoryInput'"
                            :set-style="{ border: 'none', outline: 'none', padding: 0, 'background-color': 'transparent !important' }"
                            :placeholder="$t('ServiceCategory[\'请输入一级分类\']')"
                            name="categoryName"
                            v-validate="'required|namedCharacter'"
                            v-model="mainCategoryName"
                            @on-confirm="handleEditCategory(mainCategory['id'], mainCategory['name'], 'main')"
                            @on-cancel="handleCloseEditMain">
                        </category-input>
                    </div>
                    <template v-else>
                        <div class="category-name">
                            <span>{{mainCategory['name']}}</span>
                            <i v-if="!mainCategory['is_built_in']"
                                class="property-edit icon-cc-edit-shape"
                                @click.stop="handleEditMain(mainCategory['id'], mainCategory['name'])">
                            </i>
                        </div>
                        <cmdb-dot-menu class="dot-menu" v-if="!mainCategory['is_built_in']">
                            <div class="menu-operational">
                                <i @click="handleShowAddChild(mainCategory['id'])">{{$t("ServiceCategory['添加二级分类']")}}</i>
                                <i class="not-allowed" v-if="mainCategory['child_category_list'].length">{{$t("Common['删除']")}}</i>
                                <i v-else @click="handleDeleteCategory(mainCategory['id'])">{{$t("Common['删除']")}}</i>
                            </div>
                        </cmdb-dot-menu>
                    </template>
                </div>
                <div class="child-category">
                    <div class="child-item child-edit" v-if="addChildStatus === mainCategory['id']">
                        <category-input
                            ref="editInput"
                            :input-ref="'categoryInput'"
                            :placeholder="$t('ServiceCategory[\'请输入二级分类\']')"
                            :edit-id="mainCategory['bk_root_id']"
                            name="categoryName"
                            v-validate="'required|namedCharacter'"
                            v-model="categoryName"
                            @on-confirm="handleAddCategory"
                            @on-cancel="handleCloseAddChild">
                        </category-input>
                    </div>
                    <div :class="['child-item', editChildStatus === childCategory['id'] ? 'child-edit' : '']"
                        v-for="(childCategory, childIndex) in mainCategory['child_category_list']"
                        :key="childIndex">
                        <category-input
                            v-if="editChildStatus === childCategory['id']"
                            ref="editInput"
                            :input-ref="'categoryInput'"
                            :placeholder="$t('ServiceCategory[\'请输入二级分类\']')"
                            name="categoryName"
                            v-validate="'required|namedCharacter'"
                            v-model="childCategoryName"
                            @on-confirm="handleEditCategory(childCategory['id'], childCategory['name'], 'child')"
                            @on-cancel="handleCloseEditChild">
                        </category-input>
                        <template v-else>
                            <div class="child-title">
                                <span>{{childCategory['name']}}</span>
                                <div class="child-edit" v-if="!childCategory['is_built_in']">
                                    <i class="property-edit icon-cc-edit-shape mr10"
                                        @click.stop="handleEditChild(childCategory['id'], childCategory['name'])">
                                    </i>
                                    <i class="icon-cc-tips-close"
                                        v-if="!childCategory['usage_amount']"
                                        @click.stop="handleDeleteCategory(childCategory['id'])">
                                    </i>
                                    <i class="icon-cc-tips-close" v-else
                                        style="color: #c4c6cc; cursor: not-allowed;"
                                        v-bktooltips="tooltips">
                                    </i>
                                </div>
                            </div>
                            <!-- <span>{{childCategory['usage_amount']}}</span> -->
                        </template>
                    </div>
                </div>
            </div>
            <div class="category-item add-item" :style="{ 'border-style': showAddMianCategory ? 'solid' : 'dashed' }">
                <div class="category-title" :style="{ 'border-bottom-style': showAddMianCategory ? 'solid' : 'dashed' }">
                    <div class="main-edit" style="width: 100%;" v-if="showAddMianCategory">
                        <category-input
                            ref="addCategoryInput"
                            :input-ref="'categoryInput'"
                            :set-style="{ border: 'none', outline: 'none', padding: 0, 'background-color': 'transparent !important' }"
                            :placeholder="$t('ServiceCategory[\'请输入一级分类\']')"
                            name="categoryName"
                            v-validate="'required|namedCharacter'"
                            v-model="categoryName"
                            @on-confirm="handleAddCategory"
                            @on-cancel="handleCloseAddBox">
                        </category-input>
                    </div>
                </div>
                <div class="child-category"></div>
                <span class="add-box" v-if="!showAddMianCategory" @click="handleAddBox"></span>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import categoryInput from './children/category-input'
    export default {
        components: {
            featureTips,
            categoryInput
        },
        data () {
            return {
                tooltips: {
                    content: this.$t("ServiceCategory['二级分类删除提示']"),
                    arrowsSize: 5
                },
                showFeatureTips: false,
                showAddMianCategory: false,
                showAddChildCategory: false,
                editMainStatus: null,
                editChildStatus: null,
                addChildStatus: null,
                categoryName: '',
                mainCategoryName: '',
                childCategoryName: '',
                list: []
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['category']
            this.getCategoryList()
        },
        methods: {
            ...mapActions('serviceClassification', [
                'searchServiceCategory',
                'createServiceCategory',
                'updateServiceCategory',
                'deleteServiceCategory'
            ]),
            getCategoryList () {
                this.searchServiceCategory({
                    params: this.$injectMetadata({})
                }).then((data) => {
                    const categoryList = data.info.map(item => {
                        return {
                            usage_amount: item['usage_amount'],
                            ...item['category']
                        }
                    })
                    const list = categoryList.filter(category => !category.hasOwnProperty('bk_parent_id') && category.id !== 1)
                    this.list = list.map(mainCategory => {
                        return {
                            ...mainCategory,
                            child_category_list: categoryList.filter(category => category['bk_parent_id'] === mainCategory['id'])
                        }
                    })
                })
            },
            createdCategory (name, rootId) {
                this.createServiceCategory({
                    params: this.$injectMetadata({
                        bk_root_id: rootId,
                        bk_parent_id: rootId,
                        name
                    })
                }).then(() => {
                    this.$success(this.$t('Common["保存成功"]'))
                    this.showAddMianCategory = false
                    this.handleCloseAddChild()
                    this.getCategoryList()
                })
            },
            async handleAddCategory (name, bk_root_id = 0) {
                if (!await this.$validator.validateAll()) {
                    this.$bkMessage({
                        message: this.errors.first('categoryName') || this.$t("ServiceCategory['请输入分类名称']"),
                        theme: 'error'
                    })
                } else {
                    this.createdCategory(name, bk_root_id)
                }
            },
            async handleEditCategory (id, name, type) {
                if (!await this.$validator.validateAll()) {
                    this.$bkMessage({
                        message: this.errors.first('categoryName') || this.$t("ServiceCategory['请输入分类名称']"),
                        theme: 'error'
                    })
                } else if (name === this.mainCategoryName || name === this.childCategoryName) {
                    this.handleCloseEditChild()
                    this.handleCloseEditMain()
                } else {
                    this.updateServiceCategory({
                        params: this.$injectMetadata({
                            id,
                            name: type === 'main' ? this.mainCategoryName : this.childCategoryName
                        })
                    }).then(() => {
                        this.$success(this.$t('Common["保存成功"]'))
                        this.handleCloseEditChild()
                        this.handleCloseEditMain()
                        this.getCategoryList()
                    })
                }
            },
            handleDeleteCategory (id) {
                this.$bkInfo({
                    title: this.$t("ServiceCategory['确认删除分类']"),
                    confirmFn: async () => {
                        await this.deleteServiceCategory({
                            params: {
                                data: this.$injectMetadata({ id })
                            },
                            config: {
                                requestId: 'delete_proc_services_category'
                            }
                        }).then(() => {
                            this.$success(this.$t('Common["删除成功"]'))
                            this.getCategoryList()
                        })
                    }
                })
            },
            handleEditMain (id, name) {
                this.editMainStatus = id
                this.mainCategoryName = name
                this.handleCloseEditChild()
                this.handleCloseAddChild()
                this.handleCloseAddBox()
                this.$nextTick(() => {
                    this.$refs.editInput[0].$refs.categoryInput.focus()
                })
            },
            handleCloseEditMain () {
                this.editMainStatus = null
            },
            handleEditChild (id, name) {
                this.editChildStatus = id
                this.childCategoryName = name
                this.handleCloseAddChild()
                this.handleCloseEditMain()
                this.handleCloseAddBox()
                this.$nextTick(() => {
                    this.$refs.editInput[0].$refs.categoryInput.focus()
                })
            },
            handleCloseEditChild () {
                this.editChildStatus = null
            },
            handleAddBox () {
                this.showAddMianCategory = true
                this.$nextTick(() => {
                    this.$refs.addCategoryInput.$refs.categoryInput.focus()
                })
            },
            handleCloseAddBox () {
                this.showAddMianCategory = false
                this.categoryName = ''
            },
            handleShowAddChild (id) {
                this.addChildStatus = id
                this.$nextTick(() => {
                    this.$refs.editInput[0].$refs.categoryInput.focus()
                })
            },
            handleCloseAddChild () {
                this.addChildStatus = null
                this.categoryName = ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .category-wrapper {
        min-width: 1442px;
        .category-list {
            display: flex;
            flex-flow: row wrap;
        }
        .category-item {
            position: relative;
            min-width: 320px;
            flex: 0 0 22%;
            border: 1px solid #dcdee5;
            margin-right: 30px;
            margin-bottom: 20px;
            &.add-item {
                .category-name {
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
        .category-title {
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
            .category-name {
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
        .child-category {
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
        .category-wrapper {
            min-width: 1650px;
            .category-item {
                min-width: auto;
                flex: 0 0 18%;
            }
        }
    }
</style>
