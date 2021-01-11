<template>
    <div class="category-wrapper" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <cmdb-tips class="mb10" tips-key="categoryTips">{{$t('服务分类功能提示')}}</cmdb-tips>
        <div class="category-filter">
            <bk-input class="filter-input"
                :clearable="true"
                :right-icon="'bk-icon icon-search'"
                :placeholder="$t('请输入关键字')"
                v-model.trim="keyword">
            </bk-input>
        </div>
        <div class="category-list">
            <div class="category-item bgc-white" v-for="(mainCategory, index) in displayList" :key="index">
                <div class="category-title" :style="{ 'background-color': editMainStatus === mainCategory['id'] ? '#f0f1f5' : '' }">
                    <div class="main-edit"
                        :style="{ width: editMainStatus === mainCategory['id'] ? '100%' : 'auto' }"
                        v-if="editMainStatus === mainCategory['id']">
                        <category-input
                            ref="editInput"
                            :input-ref="'categoryInput'"
                            :placeholder="$t('请输入一级分类')"
                            name="categoryName"
                            v-validate="'required|namedCharacter|length:128'"
                            v-model="mainCategoryName"
                            @on-confirm="handleEditCategory(mainCategory, 'main', index)"
                            @on-cancel="handleCloseEditMain">
                        </category-input>
                    </div>
                    <template v-else>
                        <div class="category-name">
                            <template v-if="mainCategory['is_built_in']">
                                <div class="category-name-text is-built-in">
                                    <div class="text-inner">
                                        <span class="main-name" :title="mainCategory.name">{{mainCategory.name}}</span>
                                        <span class="main-id">{{mainCategory.id}}</span>
                                    </div>
                                </div>
                                <span class="built-in-sign">{{$t('内置')}}</span>
                            </template>
                            <cmdb-auth v-else
                                :auth="{ type: $OPERATION.U_SERVICE_CATEGORY, relation: [bizId] }">
                                <div slot-scope="{ disabled }" :class="['category-name-text', { disabled }]">
                                    <div class="text-inner" @click.stop="handleEditMain(mainCategory['id'], mainCategory['name'])">
                                        <span class="main-name" :title="mainCategory.name">{{mainCategory.name}}</span>
                                        <span class="main-id">{{mainCategory.id}}</span>
                                    </div>
                                </div>
                            </cmdb-auth>
                        </div>
                        <div class="menu-operational" v-if="!mainCategory['is_built_in']">
                            <cmdb-auth :auth="{ type: $OPERATION.C_SERVICE_CATEGORY, relation: [bizId] }">
                                <bk-button slot-scope="{ disabled }"
                                    class="menu-btn"
                                    :disabled="disabled"
                                    :text="true"
                                    @click="handleShowAddChild(mainCategory['id'])">
                                    <i class="bk-cmdb-icon icon-cc-plus"></i>
                                </bk-button>
                            </cmdb-auth>
                            <cmdb-auth :auth="{ type: $OPERATION.D_SERVICE_CATEGORY, relation: [bizId] }">
                                <template slot-scope="{ disabled }">
                                    <bk-button v-if="disabled || !mainCategory['child_category_list'].length"
                                        class="menu-btn"
                                        :text="true"
                                        :disabled="disabled"
                                        @click="handleDeleteCategory(mainCategory['id'], 'main', index)">
                                        <i class="bk-cmdb-icon icon-cc-del"></i>
                                    </bk-button>
                                    <span class="menu-btn no-allow-btn" v-else v-bk-tooltips="deleteBtnTips">
                                        <i class="bk-cmdb-icon icon-cc-del"></i>
                                    </span>
                                </template>
                            </cmdb-auth>
                        </div>
                    </template>
                </div>
                <div class="child-category">
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
                            v-validate="'required|namedCharacter|length:128'"
                            v-model="childCategoryName"
                            @on-confirm="handleEditCategory(childCategory, 'child', index)"
                            @on-cancel="handleCloseEditChild">
                        </category-input>
                        <template v-else>
                            <div class="child-title">
                                <span :title="childCategory['name']">{{childCategory['name']}}</span>
                                <span class="child-id" :title="childCategory['id']">{{childCategory['id']}}</span>
                                <div class="child-edit" v-if="!childCategory['is_built_in']">
                                    <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_SERVICE_CATEGORY, relation: [bizId] }">
                                        <bk-button slot-scope="{ disabled }"
                                            class="child-edit-btn"
                                            theme="primary"
                                            :text="true"
                                            :disabled="disabled"
                                            @click.stop="handleEditChild(childCategory['id'], childCategory['name'])">
                                            <i class="icon-cc-edit-shape"></i>
                                        </bk-button>
                                    </cmdb-auth>
                                    <cmdb-auth :auth="{ type: $OPERATION.D_SERVICE_CATEGORY, relation: [bizId] }"
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
                                        v-bk-tooltips.top="tooltips">
                                    </i>
                                </div>
                            </div>
                        </template>
                    </div>
                    <div class="child-item is-add" v-if="!mainCategory['is_built_in'] && !addChildStatus">
                        <div class="child-title">
                            <cmdb-auth @update-auth="isAuthcompleted = true" :auth="{ type: $OPERATION.C_SERVICE_CATEGORY, relation: [bizId] }">
                                <bk-button slot-scope="{ disabled }"
                                    class="add-btn"
                                    :disabled="disabled"
                                    :text="true"
                                    v-show="isAuthcompleted"
                                    @click="handleShowAddChild(mainCategory['id'])">
                                    <i class="bk-cmdb-icon icon-cc-plus"></i>{{$t('添加')}}
                                </bk-button>
                            </cmdb-auth>
                        </div>
                    </div>
                    <div class="child-item child-edit" v-if="addChildStatus === mainCategory['id']">
                        <category-input
                            class="child-input"
                            ref="editInput"
                            :input-ref="'categoryInput'"
                            :placeholder="$t('请输入二级分类')"
                            :edit-id="mainCategory['bk_root_id']"
                            name="categoryName"
                            v-validate="'required|namedCharacter|length:128'"
                            v-model="categoryName"
                            @on-confirm="handleAddCategory"
                            @on-cancel="handleCloseAddChild">
                        </category-input>
                    </div>
                </div>
            </div>
            <div class="category-item add-item"
                :style="{ 'border-style': showAddMianCategory ? 'solid' : 'dashed' }"
                v-show="!keyword">
                <div class="category-title" :style="{ 'border-bottom-style': showAddMianCategory ? 'solid' : 'dashed' }">
                    <div class="main-edit" style="width: 100%;" v-if="showAddMianCategory">
                        <category-input
                            ref="addCategoryInput"
                            :input-ref="'categoryInput'"
                            :placeholder="$t('请输入一级分类')"
                            name="categoryName"
                            v-validate="'required|namedCharacter|length:128'"
                            v-model="categoryName"
                            @on-confirm="handleAddCategory"
                            @on-cancel="handleCloseAddBox">
                        </category-input>
                    </div>
                </div>
                <div class="child-category"></div>
                <cmdb-auth :auth="{ type: $OPERATION.C_SERVICE_CATEGORY, relation: [bizId] }"
                    v-show="!showAddMianCategory">
                    <bk-button slot-scope="{ disabled }"
                        class="add-btn"
                        :disabled="disabled"
                        @click="handleAddBox">
                    </bk-button>
                </cmdb-auth>
            </div>
            <bk-exception v-show="!displayList.length && !$loading(Object.values(request))" type="search-empty" scene="part"></bk-exception>
        </div>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import debounce from 'lodash.debounce'
    import categoryInput from './children/category-input'
    export default {
        components: {
            categoryInput
        },
        data () {
            return {
                tooltips: {
                    content: this.$t('二级分类删除提示'),
                    onShow: this.handleCategoryTipsToggle,
                    onHide: this.handleCategoryTipsToggle
                },
                deleteBtnTips: {
                    content: this.$t('请先清空二级分类'),
                    placements: ['right']
                },
                showAddMianCategory: false,
                showAddChildCategory: false,
                editMainStatus: null,
                editChildStatus: null,
                addChildStatus: null,
                categoryName: '',
                mainCategoryName: '',
                childCategoryName: '',
                list: [],
                displayList: [],
                keyword: '',
                isAuthcompleted: false,
                request: {
                    category: Symbol('category')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId'])
        },
        watch: {
            list (list) {
                this.displayList = list
            },
            keyword () {
                this.handleFilter()
            }
        },
        created () {
            this.handleFilter = debounce(this.searchList, 300)
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
                    params: { bk_biz_id: this.bizId },
                    config: { requestId: this.request.category }
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
            searchList () {
                if (this.keyword) {
                    const reg = new RegExp(this.keyword, 'i')
                    this.displayList = this.list.filter(mainCategory => {
                        if (reg.test(mainCategory.name) || reg.test(mainCategory.id)) {
                            return true
                        }
                        return mainCategory.child_category_list.findIndex(subCategory => {
                            return reg.test(subCategory.name) || reg.test(subCategory.id)
                        }) !== -1
                    })
                } else {
                    this.displayList = this.list
                }
            },
            createdCategory (name, rootId) {
                this.createServiceCategory({
                    params: {
                        bk_biz_id: this.bizId,
                        bk_root_id: rootId,
                        bk_parent_id: rootId,
                        name
                    }
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
                        childList.push(res)
                        this.$set(this.list[markIndex], 'child_category_list', childList)
                    } else {
                        this.getCategoryList()
                    }
                })
            },
            handleCategoryTipsToggle (tipsInstance) {
                const willShow = !tipsInstance.state.isVisible
                tipsInstance.reference.parentElement.classList[willShow ? 'add' : 'remove']('tips-active')
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
                        params: {
                            bk_biz_id: this.bizId,
                            id: data.id,
                            name: type === 'main' ? this.mainCategoryName : this.childCategoryName
                        }
                    }).then(res => {
                        this.$success(this.$t('保存成功'))
                        this.handleCloseEditChild()
                        // this.handleCloseEditMain()
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
                                data: { id, bk_biz_id: this.bizId }
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
                this.isAuthcompleted = false
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
        padding: 15px 20px 0;
        .category-list {
            display: flex;
            flex-flow: row wrap;
        }
        .category-item {
            position: relative;
            flex: 0 0 calc(25% - 15px);
            border: 1px solid #dcdee5;
            border-radius: 0px 0px 2px 2px;
            margin-left: 20px;
            margin-bottom: 20px;
            overflow: hidden;
            &:hover:not(.add-item) {
                box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.1);
                .menu-operational {
                    display: flex;
                }
            }
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
            padding: 0 12px 0 12px;
            height: 52px;
            font-size: 14px;
            color: #63656e;
            font-weight: bold;
            border-bottom: 1px solid #dcdee5;
            .main-edit {
                display: flex;
                align-items: center;
                /deep/ .cagetory-input .bk-form-input {
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
                display: flex;
                flex: 1;
                width: 100%;
                overflow: hidden;
                .auth-box {
                    width: 100%;
                }
                .category-name-text {
                    max-width: 100%;
                    &.disabled {
                        color: #dcdee5;
                    }
                    .text-inner {
                        max-width: 100%;
                        display: inline-flex;
                        flex-direction: column;
                        padding: 2px 6px;
                        line-height: normal;
                        cursor: pointer;
                        &:hover {
                            background: #f0f1f5;
                        }
                        .main-id {
                            font-size: 12px;
                            font-weight: 400;
                            color: #C4C6CC;
                            &::before {
                                content: "#";
                            }
                        }
                        .main-id,
                        .main-name {
                            @include ellipsis;
                        }
                        .main-name {
                            height: 20px;
                            line-height: 20px;
                        }
                    }
                    &.is-built-in {
                        max-width: calc(100% - 40px);
                        .text-inner {
                            cursor: initial;
                            &:hover {
                                background: transparent;
                            }
                        }
                    }
                }
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
                    height: 20px;
                    line-height: 20px;
                    margin: 2px 0 0 4px;
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
        }
        .child-category {
            height: 280px;
            padding: 0 10px 10px 38px;
            @include scrollbar-y;
            .child-item {
                @include space-between;
                position: relative;
                z-index: 10;
                line-height: 32px;
                &.child-edit {
                    &:first-child::after {
                        height: 32px;
                    }

                    .child-input {
                        margin-left: 10px;
                        padding-left: 8px;
                    }
                }
                &:hover:not(.is-built-in):not(.is-add) {
                    .child-title {
                        background-color: #fafbfd;
                        color: #3a84ff;
                    }
                    > span {
                        display: none;
                    }
                    .child-edit {
                        display: block;
                    }
                    .child-id {
                        display: none;
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
                    padding-right: 8px;
                    padding-left: 8px;
                    margin-left: 10px;
                    > span {
                        @include ellipsis;
                        padding-right: 10px;
                    }

                    .child-id {
                        min-width: 42px;
                        font-size: 12px;
                        color: #C4C6CC;
                        padding-right: 6px;
                        text-align: right;
                        &::before {
                            content: "#";
                        }
                    }
                }
                > span {
                    color: #c4c6cc;
                    padding-right: 18px;
                }
                .child-edit {
                    display: none;
                    margin-left: auto;
                    &.tips-active {
                        opacity: 1;
                    }
                    .child-edit-btn {
                        .icon-cc-tips-close {
                            font-size: 12px;
                        }
                        .icon-cc-edit-shape {
                            font-size: 16px;
                        }
                    }
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
                &.is-add {
                    .add-btn {
                        color: #979BA5;
                        &:hover {
                            color: #3a84ff;
                        }
                    }
                }
            }
        }
        .category-filter {
            margin-bottom: 12px;
            .filter-input {
                width: 260px;
            }
        }
    }
    .menu-operational {
        display: none;
        padding: 6px 0;
        line-height: 30px;
        .auth-box {
            display: block;
        }
        .menu-btn {
            display: block;
            width: 100%;
            height: 30px;
            line-height: 30px;
            padding: 0 7px;
            text-align: left;
            color: #979BA5;
            outline: none;
            &:hover {
                color: #3a84ff;
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
            .bk-cmdb-icon {
                font-size: 16px;
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
