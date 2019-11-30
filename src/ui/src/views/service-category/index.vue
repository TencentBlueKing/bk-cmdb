<template>
    <div class="category-wrapper">
        <feature-tips
            :feature-name="'category'"
            :show-tips="showFeatureTips"
            :desc="$t('服务分类功能提示')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="category-list">
            <div class="category-item bgc-white" v-for="(mainCategory, index) in list" :key="index">
                <div class="category-title" :style="{ 'background-color': mainCategory['editStatus'] ? '#f0f1f5' : '' }">
                    <div class="main-edit"
                        :style="{ width: editMainStatus === mainCategory['id'] ? '100%' : 'auto' }"
                        v-if="editMainStatus === mainCategory['id']">
                        <category-input
                            ref="editInput"
                            :input-ref="'categoryInput'"
                            :placeholder="$t('请输入一级分类')"
                            name="categoryName"
                            v-validate="'required|namedCharacter|length:256'"
                            v-model="mainCategoryName"
                            @on-confirm="handleEditCategory(mainCategory, 'main', index)"
                            @on-cancel="handleCloseEditMain">
                        </category-input>
                    </div>
                    <template v-else>
                        <div class="category-name">
                            <span>{{mainCategory['name']}}</span>
                            <cmdb-auth v-if="!mainCategory['is_built_in']"
                                :auth="$authResources({ type: $OPERATION.U_SERVICE_CATEGORY })">
                                <bk-button slot-scope="{ disabled }"
                                    class="main-edit-btn"
                                    theme="primary"
                                    :disabled="disabled"
                                    :text="true"
                                    @click.stop="handleEditMain(mainCategory['id'], mainCategory['name'])">
                                    <i class="icon-cc-edit-shape"></i>
                                </bk-button>
                            </cmdb-auth>
                            <span v-else class="built-in-sign">{{$t('内置')}}</span>
                        </div>
                        <cmdb-dot-menu class="dot-menu-operation" v-if="!mainCategory['is_built_in']">
                            <div class="menu-operational">
                                <cmdb-auth :auth="$authResources({ type: $OPERATION.C_SERVICE_CATEGORY })">
                                    <bk-button slot-scope="{ disabled }"
                                        class="menu-btn"
                                        :disabled="disabled"
                                        :text="true"
                                        @click="handleShowAddChild(mainCategory['id'])">
                                        {{$t('添加二级分类')}}
                                    </bk-button>
                                </cmdb-auth>
                                <cmdb-auth :auth="$authResources({ type: $OPERATION.D_SERVICE_CATEGORY })">
                                    <template slot-scope="{ disabled }">
                                        <bk-button v-if="disabled || !mainCategory['child_category_list'].length"
                                            class="menu-btn"
                                            :text="true"
                                            :disabled="disabled"
                                            @click="handleDeleteCategory(mainCategory['id'], 'main', index)">
                                            {{$t('删除')}}
                                        </bk-button>
                                        <span class="menu-btn no-allow-btn" v-else v-bk-tooltips="deleteBtnTips">
                                            {{$t('删除')}}
                                        </span>
                                    </template>
                                </cmdb-auth>
                            </div>
                        </cmdb-dot-menu>
                    </template>
                </div>
                <div class="child-category">
                    <div class="child-item child-edit" v-if="addChildStatus === mainCategory['id']">
                        <category-input
                            class="child-input"
                            ref="editInput"
                            :input-ref="'categoryInput'"
                            :placeholder="$t('请输入二级分类')"
                            :edit-id="mainCategory['bk_root_id']"
                            name="categoryName"
                            v-validate="'required|namedCharacter|length:256'"
                            v-model="categoryName"
                            @on-confirm="handleAddCategory"
                            @on-cancel="handleCloseAddChild">
                        </category-input>
                    </div>
                    <div v-for="(childCategory, childIndex) in mainCategory['child_category_list']"
                        :key="childIndex"
                        :class="['child-item', {
                            'child-edit': editChildStatus === childCategory['id'],
                            'is-built-in': childCategory['is_built_in']
                        }]">
                        <category-input
                            v-if="editChildStatus === childCategory['id']"
                            class="child-input"
                            ref="editInput"
                            :input-ref="'categoryInput'"
                            :placeholder="$t('请输入二级分类')"
                            name="categoryName"
                            v-validate="'required|namedCharacter|length:256'"
                            v-model="childCategoryName"
                            @on-confirm="handleEditCategory(childCategory, 'child', index)"
                            @on-cancel="handleCloseEditChild">
                        </category-input>
                        <template v-else>
                            <div class="child-title">
                                <span>{{childCategory['name']}}</span>
                                <div class="child-edit" v-if="!childCategory['is_built_in']">
                                    <cmdb-auth class="mr10" :auth="$authResources({ type: $OPERATION.U_SERVICE_CATEGORY })">
                                        <bk-button slot-scope="{ disabled }"
                                            class="child-edit-btn"
                                            theme="primary"
                                            :text="true"
                                            :disabled="disabled"
                                            @click.stop="handleEditChild(childCategory['id'], childCategory['name'])">
                                            <i class="icon-cc-edit-shape"></i>
                                        </bk-button>
                                    </cmdb-auth>
                                    <cmdb-auth :auth="$authResources({ type: $OPERATION.D_SERVICE_CATEGORY })"
                                        v-if="!childCategory['usage_amount']">
                                        <bk-button slot-scope="{ disabled }"
                                            class="child-edit-btn"
                                            theme="primary"
                                            :text="true"
                                            :disabled="disabled"
                                            @click.stop="handleDeleteCategory(childCategory['id'], 'child', index)">
                                            <i class="icon-cc-tips-close"></i>
                                        </bk-button>
                                    </cmdb-auth>
                                    <i class="icon-cc-tips-close" v-else
                                        style="color: #dcdee5; cursor: not-allowed; outline: none;"
                                        v-bk-tooltips.bottom="tooltips">
                                    </i>
                                </div>
                            </div>
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
                            :placeholder="$t('请输入一级分类')"
                            name="categoryName"
                            v-validate="'required|namedCharacter|length:256'"
                            v-model="categoryName"
                            @on-confirm="handleAddCategory"
                            @on-cancel="handleCloseAddBox">
                        </category-input>
                    </div>
                </div>
                <div class="child-category"></div>
                <cmdb-auth :auth="$authResources({ type: $OPERATION.C_SERVICE_CATEGORY })"
                    v-if="!showAddMianCategory">
                    <bk-button slot-scope="{ disabled }"
                        class="add-btn"
                        :disabled="disabled"
                        @click="handleAddBox">
                    </bk-button>
                </cmdb-auth>
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
                    content: this.$t('二级分类删除提示')
                },
                deleteBtnTips: {
                    content: this.$t('请先清空二级分类'),
                    placements: ['right']
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
                    params: this.$injectMetadata({}, { injectBizId: true })
                }).then((data) => {
                    const categoryList = data.info.map(item => {
                        return {
                            usage_amount: item['usage_amount'],
                            ...item['category']
                        }
                    })
                    const list = categoryList.filter(category => !category['bk_parent_id'] && !(category['name'] === 'Default' && category['is_built_in']))
                    this.list = list.map(mainCategory => {
                        return {
                            ...mainCategory,
                            child_category_list: categoryList.filter(category => category['bk_parent_id'] === mainCategory['id'])
                        }
                    }).sort((prev, next) => prev.id - next.id)
                })
            },
            createdCategory (name, rootId) {
                this.createServiceCategory({
                    params: this.$injectMetadata({
                        bk_root_id: rootId,
                        bk_parent_id: rootId,
                        name
                    }, { injectBizId: true })
                }).then(res => {
                    this.$success(this.$t('保存成功'))
                    this.showAddMianCategory = false
                    this.handleCloseAddChild()
                    if (rootId) {
                        let markIndex = null
                        const currentObj = this.list.find((category, index) => {
                            markIndex = index
                            return category.hasOwnProperty('bk_root_id') && category['bk_root_id'] === rootId
                        })
                        const childList = currentObj ? currentObj['child_category_list'] : []
                        childList.unshift(res)
                        this.$set(this.list[markIndex], 'child_category_list', childList)
                    } else {
                        this.getCategoryList()
                    }
                })
            },
            async handleAddCategory (name, bk_root_id = 0) {
                if (!await this.$validator.validateAll()) {
                    this.$bkMessage({
                        message: this.errors.first('categoryName') || this.$t('请输入分类名称'),
                        theme: 'error'
                    })
                } else {
                    this.createdCategory(name, bk_root_id)
                }
            },
            async handleEditCategory (data, type, mainIndex) {
                if (!await this.$validator.validateAll()) {
                    this.$bkMessage({
                        message: this.errors.first('categoryName') || this.$t('请输入分类名称'),
                        theme: 'error'
                    })
                } else if (data.name === this.mainCategoryName || data.name === this.childCategoryName) {
                    this.handleCloseEditChild()
                    this.handleCloseEditMain()
                } else {
                    this.updateServiceCategory({
                        params: this.$injectMetadata({
                            id: data.id,
                            name: type === 'main' ? this.mainCategoryName : this.childCategoryName
                        }, { injectBizId: true })
                    }).then(res => {
                        this.$success(this.$t('保存成功'))
                        this.handleCloseEditChild()
                        this.handleCloseEditMain()
                        if (mainIndex !== undefined && type === 'child') {
                            const childList = this.list[mainIndex].child_category_list.map(child => {
                                if (child.id === res.id) {
                                    return res
                                }
                                return child
                            })
                            this.$set(this.list[mainIndex], 'child_category_list', childList)
                        } else {
                            this.$set(this.list[mainIndex], 'name', res.name)
                        }
                    })
                }
            },
            handleDeleteCategory (id, type, index) {
                this.$bkInfo({
                    title: this.$t('确认删除分类'),
                    zIndex: 999,
                    confirmFn: async () => {
                        await this.deleteServiceCategory({
                            params: {
                                data: this.$injectMetadata({ id }, { injectBizId: true })
                            },
                            config: {
                                requestId: 'delete_proc_services_category'
                            }
                        }).then(() => {
                            this.$success(this.$t('删除成功'))
                            if (type === 'main') {
                                this.list.splice(index, 1)
                            } else {
                                let childIndex = -1
                                this.list[index]['child_category_list'].find((category, findIndex) => {
                                    childIndex = findIndex
                                    return category.id === id
                                })
                                this.list[index]['child_category_list'].splice(childIndex, 1)
                            }
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
        padding: 0 20px;
        .category-list {
            display: flex;
            flex-flow: row wrap;
        }
        .category-item {
            position: relative;
            flex: 0 0 calc(25% - 15px);
            border: 1px solid #dcdee5;
            margin-left: 20px;
            margin-bottom: 20px;
            overflow: hidden;
            &:nth-child(4n+1) {
                margin-left: 0;
            }
            &.add-item {
                .category-name {
                    color: #dcdee5 !important;
                }
                .child-title {
                    color: #dcdee5 !important;
                    background-color: transparent !important;
                }
                .auth-box {
                    position: absolute;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                }
                .add-btn {
                    position: absolute;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    background-color: transparent;
                    border: none;
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
                /deep/ .cagetory-input .bk-form-input{
                    font-size: 14px;
                    color: #63656e;
                    background-color: transparent !important;
                    border: none;
                    outline: none;
                    font-weight: normal;
                    &:focus {
                        background-color: transparent !important;
                    }
                }
            }
            .category-name {
                @include ellipsis;
                flex: 1;
                padding-right: 20px;
                .icon-cc-edit-shape {
                    font-size: 14px;
                    display: none;
                    cursor: pointer;
                }
                &:hover .icon-cc-edit-shape {
                    display: inline !important;
                }
                .built-in-sign {
                    display: inline-block;
                    height: 19px;
                    line-height: 18px;
                    margin: 0 0 0 4px;
                    padding: 0 6px;
                    font-size: 12px;
                    color: #ffffff;
                    text-align: center;
                    background-color: #d3d5dd;
                    border-radius: 2px;
                }
                .auth-box {
                    vertical-align: middle;
                }
                .main-edit-btn {
                    display: block;
                    font-size: 0;
                    line-height: 1;
                }
            }
            .dot-menu-operation {
                cursor: pointer;
                /deep/ .bk-tooltip-ref {
                    width: 100%;
                }
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
                &:hover:not(.is-built-in) {
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
                    >span {
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
                }
                .edit-box {
                    width: 100%;
                }
                .child-input {
                    /deep/ .bk-form-input {
                        vertical-align: top;
                        border: 1px solid #c4c6cc;
                        background-color: #ffffff !important;
                    }
                }
            }
        }
    }
    .menu-operational {
        padding: 6px 0;
        line-height: 32px;
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
            color: #63656e;
            outline: none;
            &:hover {
                color: #3a84ff;
                background-color: #e1ecff;
            }
            &:disabled {
                color: #dcdee5;
                background-color: transparent;
            }
            &.no-allow-btn {
                cursor: not-allowed;
                color: #dcdee5;
                background-color: transparent;
            }
        }
    }
    @media screen and (min-width: 1920px){
        .category-item {
            flex: 0 0 calc(20% - 16px) !important;
            &:nth-child(4n+1) {
                margin-left: 20px !important;
            }
            &:nth-child(5n+1) {
                margin-left: 0 !important;
            }
        }
    }
</style>
