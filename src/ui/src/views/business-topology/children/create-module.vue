<template>
    <div class="node-create-layout">
        <h2 class="node-create-title">{{$t('新增模块')}}</h2>
        <div class="node-create-path" :title="topoPath">{{$t('添加节点已选择')}}：{{topoPath}}</div>
        <div class="node-create-form"
            :style="{
                'max-height': Math.min($APP.height - 400, 400) + 'px'
            }">
            <div class="form-item clearfix mt30">
                <div class="create-type fl">
                    <input class="type-radio"
                        type="radio"
                        id="formTemplate"
                        name="createType"
                        v-model="withTemplate"
                        :value="1">
                    <label for="formTemplate">{{$t('从模板创建')}}</label>
                </div>
                <div class="create-type fl ml50">
                    <input class="type-radio"
                        type="radio"
                        id="createDirectly"
                        name="createType"
                        v-model="withTemplate"
                        :value="0">
                    <label for="createDirectly">{{$t('直接创建')}}</label>
                </div>
            </div>
            <div class="form-item" v-if="withTemplate">
                <label>{{$t('服务模板')}}</label>
                <bk-select style="width: 100%;"
                    :clearable="false"
                    :searchable="templateList.length > 7"
                    v-model="template"
                    v-validate.disabled="'required'"
                    data-vv-name="template"
                    key="template">
                    <bk-option v-for="(option, index) in templateList"
                        :key="index"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                    <div class="add-template" slot="extension" @click="jumpServiceTemplate" v-if="!templateList.length">
                        <i class="bk-icon icon-plus-circle"></i>
                        <span>{{$t('新建服务模板')}}</span>
                    </div>
                </bk-select>
                <span class="form-error" v-if="errors.has('template')">{{errors.first('template')}}</span>
            </div>
            <div class="form-item">
                <label>
                    {{$t('模块名称')}}
                    <font color="red">*</font>
                    <i class="icon-cc-tips"
                        v-bk-tooltips.top="$t('模块名称提示')"
                        v-if="withTemplate === 1">
                    </i>
                </label>
                <cmdb-form-singlechar
                    v-model="moduleName"
                    v-validate="'required|singlechar|length:256'"
                    data-vv-name="moduleName"
                    key="moduleName"
                    :placeholder="$t('请输入xx', { name: $t('模块名称') })"
                    :disabled="!!withTemplate">
                </cmdb-form-singlechar>
                <span class="form-error" v-if="errors.has('moduleName')">{{errors.first('moduleName')}}</span>
            </div>
            <div class="form-item clearfix" v-if="!withTemplate">
                <label>{{$t('所属服务分类')}}<font color="red">*</font></label>
                <cmdb-selector class="service-class fl"
                    v-model="firstClass"
                    v-validate.disabled="'required'"
                    data-vv-name="firstClass"
                    key="firstClass"
                    :auto-select="false"
                    :list="firstClassList">
                </cmdb-selector>
                <cmdb-selector class="service-class fr"
                    v-model="secondClass"
                    v-validate.disabled="'required'"
                    data-vv-name="secondClass"
                    key="secondClass"
                    :list="secondClassList">
                </cmdb-selector>
                <span class="form-error" v-if="errors.has('firstClass')">{{errors.first('firstClass')}}</span>
                <span class="form-error second-class" v-if="errors.has('secondClass')">{{errors.first('secondClass')}}</span>
            </div>
        </div>
        <div class="node-create-options">
            <bk-button theme="primary"
                :disabled="$loading() || errors.any()"
                @click="handleSave">
                {{$t('提交')}}
            </bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { MENU_BUSINESS_SERVICE_TEMPLATE } from '@/dictionary/menu-symbol'
    export default {
        props: {
            parentNode: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                withTemplate: 1,
                createTypeList: [{
                    id: 1,
                    name: this.$t('从模板创建')
                }, {
                    id: 0,
                    name: this.$t('直接创建')
                }],
                template: '',
                templateList: [],
                moduleName: '',
                firstClass: '',
                firstClassList: [],
                secondClass: '',
                values: {}
            }
        },
        computed: {
            topoPath () {
                const nodePath = [...this.parentNode.parents, this.parentNode]
                return nodePath.map(node => node.data.bk_inst_name).join('/')
            },
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            serviceTemplateMap () {
                return this.$store.state.businessHost.serviceTemplateMap
            },
            currentTemplate () {
                return this.templateList.find(item => item.id === this.template) || {}
            },
            categoryMap () {
                return this.$store.state.businessHost.categoryMap
            },
            currentCategory () {
                return this.firstClassList.find(category => category.id === this.firstClass) || {}
            },
            secondClassList () {
                return this.currentCategory.secondCategory || []
            }
        },
        watch: {
            withTemplate (withTemplate) {
                if (withTemplate) {
                    this.firstClass = ''
                    this.secondClass = ''
                } else {
                    this.template = ''
                    this.getServiceCategories()
                }
            },
            template (template) {
                if (template) {
                    this.moduleName = this.currentTemplate.name
                } else {
                    this.moduleName = ''
                }
            }
        },
        created () {
            this.getServiceTemplates()
        },
        methods: {
            async getServiceTemplates () {
                if (this.serviceTemplateMap.hasOwnProperty(this.business)) {
                    this.templateList = this.serviceTemplateMap[this.business]
                } else {
                    try {
                        const data = await this.$store.dispatch('serviceTemplate/searchServiceTemplate', {
                            params: this.$injectMetadata({}, { injectBizId: true })
                        })
                        const templates = data.info.map(item => item.service_template)
                        this.templateList = templates
                        this.$store.commit('businessHost/setServiceTemplate', {
                            id: this.business,
                            templates: templates
                        })
                    } catch (e) {
                        console.error(e)
                        this.templateList = []
                    }
                }
                this.template = this.templateList[0] ? this.templateList[0].id : ''
            },
            async getServiceCategories () {
                if (this.categoryMap.hasOwnProperty(this.business)) {
                    this.firstClassList = this.categoryMap[this.business]
                } else {
                    try {
                        const data = await this.$store.dispatch('serviceClassification/searchServiceCategory', {
                            params: this.$injectMetadata({}, { injectBizId: true })
                        })
                        const categories = this.collectServiceCategories(data.info)
                        this.firstClassList = categories
                        this.$store.commit('businessHost/setCategories', {
                            id: this.business,
                            categories: categories
                        })
                    } catch (e) {
                        console.error(e)
                        this.firstClassList = []
                    }
                }
            },
            collectServiceCategories (data) {
                const categories = []
                data.forEach(item => {
                    if (!item.category.bk_parent_id) {
                        categories.push(item.category)
                    }
                })
                categories.forEach(category => {
                    category.secondCategory = data.filter(item => item.category.bk_parent_id === category.id).map(item => item.category)
                })
                return categories
            },
            handleSave () {
                this.$validator.validateAll().then(isValid => {
                    if (isValid) {
                        this.$emit('submit', {
                            bk_module_name: this.moduleName,
                            service_category_id: this.withTemplate ? this.currentTemplate.service_category_id : this.secondClass,
                            service_template_id: this.withTemplate ? this.template : 0
                        })
                    }
                })
            },
            handleCancel () {
                this.$emit('cancel')
            },
            jumpServiceTemplate () {
                this.$router.push({ name: MENU_BUSINESS_SERVICE_TEMPLATE })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .node-create-layout {
        position: relative;
    }
    .node-create-title {
        margin-top: -14px;
        padding: 0 26px;
        line-height: 30px;
        font-size: 22px;
        color: #333948;
    }
    .node-create-path {
        padding: 23px 26px 0;
        margin: 0 0 -5px 0;
        font-size: 12px;
        @include ellipsis;
    }
    .node-create-form {
        padding: 0 26px 27px;
        overflow: visible;
    }
    .form-item {
        margin: 15px 0 0 0;
        position: relative;
        label {
            display: block;
            padding: 7px 0;
            line-height: 19px;
            font-size: 14px;
            color: #63656E;
        }
        .service-class {
            width: 260px;
            @include inlineBlock;
        }
        .form-error {
            position: absolute;
            top: 100%;
            left: 0;
            font-size: 12px;
            color: $cmdbDangerColor;
            &.second-class {
                left: 270px;
            }
        }
        .create-type {
            display: flex;
            align-items: center;
            .type-radio {
                -webkit-appearance: none;
                width: 16px;
                height: 16px;
                padding: 3px;
                border: 1px solid #979BA5;
                border-radius: 50%;
                background-clip: content-box;
                outline: none;
                cursor: pointer;
                &:checked {
                    border-color: #3A84FF;
                    background-color: #3A84FF;
                }
            }
            label {
                padding: 0 0 0 6px;
                font-size: 14px;
                cursor: pointer;
            }
        }
    }
    .node-create-options {
        padding: 9px 20px;
        border-top: 1px solid $cmdbBorderColor;
        text-align: right;
        background-color: #fafbfd;
    }
    font {
        padding: 0 2px;
    }
    .add-template {
        width: 20%;
        cursor: pointer;
        .icon-plus-circle {
            @include inlineBlock;
            font-size: 14px;
        }
        span {
            @include inlineBlock;
        }
    }
</style>
