<template>
    <div class="node-create-layout">
        <h2 class="node-create-title">{{$t('BusinessTopology["新增模块"]')}}</h2>
        <div class="node-create-path" :title="topoPath">{{$t('Common["已选择"]')}}：{{topoPath}}</div>
        <div class="node-create-form">
            <div class="form-item">
                <label>{{$t('BusinessTopology["创建方式"]')}}</label>
                <cmdb-selector
                    v-model="withTemplate"
                    :list="createTypeList">
                </cmdb-selector>
            </div>
            <div class="form-item" v-if="withTemplate">
                <label>{{$t('BusinessTopology["模板名称"]')}}</label>
                <cmdb-selector
                    v-model="template"
                    :list="templateList">
                </cmdb-selector>
            </div>
            <div class="form-item">
                <label>{{$t('BusinessTopology["模块名称"]')}}</label>
                <cmdb-form-singlechar
                    v-model="moduleName"
                    :disabled="!!withTemplate">
                </cmdb-form-singlechar>
            </div>
            <div class="form-item clearfix" v-if="!withTemplate">
                <label>{{$t('BusinessTopology["服务实例分类"]')}}</label>
                <cmdb-selector class="service-class fl"
                    v-model="firstClass"
                    :list="firstClassList">
                </cmdb-selector>
                <cmdb-selector class="service-class fr"
                    v-model="secondClass"
                    :list="secondClassList">
                </cmdb-selector>
            </div>
        </div>
        <div class="node-create-options">
            <bk-button type="primary"
                :disabled="$loading() || errors.any()"
                @click="handleSave">
                {{$t('Common["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="handleCancel">{{$t('Common["取消"]')}}</bk-button>
        </div>
    </div>
</template>

<script>
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
                    name: this.$t('BusinessTopology["从模板创建"]')
                }, {
                    id: 0,
                    name: this.$t('BusinessTopology["直接创建"]')
                }],
                template: '',
                templateList: [],
                moduleName: '',
                firstClass: '',
                firstClassList: [],
                secondClass: '',
                secondClassList: [],
                values: {}
            }
        },
        computed: {
            topoPath () {
                const nodePath = [...this.parentNode.parents, this.parentNode]
                return nodePath.map(node => node.data.bk_inst_name).join('/')
            }
        },
        methods: {
            handleSave () {
                this.$validator.validateAll().then(isValid => {
                    if (isValid) {
                        this.$emit('submit', this.values)
                    }
                })
            },
            handleCancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .node-create-layout {
        position: relative;
    }
    .node-create-title {
        position: absolute;
        top: -20px;
        left: 0;
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
        max-height: 400px;
        padding: 0 26px 27px;
        @include scrollbar-y;
    }
    .form-item {
        margin: 15px 0 0 0;
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
    }
    .node-create-options {
        padding: 9px 20px;
        border-top: 1px solid $cmdbBorderColor;
        text-align: right;
        background-color: #fafbfd;
    }
</style>
