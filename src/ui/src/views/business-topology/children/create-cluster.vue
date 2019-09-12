<template>
    <div class="create-node-layout">
        <h2 class="node-create-title">{{$t('新增集群')}}</h2>
        <div class="node-create-path" :title="topoPath">{{$t('添加节点已选择')}}：{{topoPath}}</div>
        <div class="node-create-form">
            <bk-radio-group class="form-item mb20" v-model="withTemplate">
                <bk-radio :value="true">{{$t('从集群模版创建')}}</bk-radio>
                <bk-radio :value="false">{{$t('直接创建')}}</bk-radio>
            </bk-radio-group>
            <div class="form-item" v-if="withTemplate">
                <label>{{$t('集群模版')}}</label>
                <bk-select style="width: 100%;"
                    :clearable="false"
                    :searchable="clusterTemplateList.length > 7"
                    v-model="clusterTemplate"
                    v-validate.disabled="'required'"
                    data-vv-name="clusterTemplate">
                    <bk-option v-for="option in clusterTemplateList"
                        :key="option.id"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                    <div class="add-template" slot="extension" v-if="!clusterTemplateList.length">
                        <i class="bk-icon icon-plus-circle"></i>
                        <span>{{$t('创建集群模版')}}</span>
                    </div>
                </bk-select>
            </div>
            <div class="form-item">
                <label>
                    {{$t('模块名称')}}
                    <span>（{{$t('使用模版需要重命名模版的集群名称')}}）</span>
                </label>
                <bk-input class="form-textarea"
                    type="textarea"
                    data-vv-name="clusterName"
                    v-validate="'required|singlechar|length:256'"
                    v-model="clusterName"
                    :rows="rows"
                    :disabled="withTemplate"
                    :placeholder="$t('请输入集群名称，同时创建多个集群，换行分隔')"
                    @keydown="handleKeydown">
                </bk-input>
                <span class="form-error" v-if="errors.has('clusterName')">{{errors.first('clusterName')}}</span>
            </div>
        </div>
        <div class="node-create-options">
            <bk-button theme="primary" class="mr10"
                :disabled="$loading() || errors.any()">
                {{$t('确定')}}
            </bk-button>
            <bk-button theme="default">{{$t('取消')}}</bk-button>
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
                withTemplate: true,
                clusterTemplate: '',
                clusterName: '',
                rows: 1,
                clusterTemplateList: []
            }
        },
        computed: {
            topoPath () {
                const nodePath = [...this.parentNode.parents, this.parentNode]
                return nodePath.map(node => node.data.bk_inst_name).join('/')
            }
        },
        methods: {
            setRows () {
                const rows = this.clusterName.split('\n').length
                this.rows = Math.min(3, Math.max(rows, 1))
            },
            handleKeydown (value, keyEvent) {
                if (['Enter', 'NumpadEnter'].includes(keyEvent.code)) {
                    this.rows = Math.min(this.rows + 1, 3)
                } else if (keyEvent.code === 'Backspace') {
                    this.$nextTick(() => {
                        this.setRows()
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .node-create-layout {
        position: relative;
    }
    .node-create-title {
        margin-top: -15px;
        padding: 0 26px;
        line-height: 30px;
        font-size: 24px;
        color: #444444;
        font-weight: normal;
    }
    .node-create-path {
        padding: 14px 26px 0;
        margin: 0 0 -5px 0;
        font-size: 12px;
        color: #63656E;
        @include ellipsis;
    }
    .node-create-form {
        padding: 20px 26px 32px;
    }
    .form-item {
        margin: 15px 0 0 0;
        position: relative;
        .bk-form-radio {
            display: inline-block;
            margin-right: 70px;
            /deep/ input[type=radio] {
                margin-top: 2px;
            }
        }
        label {
            display: block;
            padding: 0 0 10px;
            line-height: 19px;
            font-size: 14px;
            color: #63656E;
            > span {
                color: #979BA5;
                font-size: 12px;
            }
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
        .form-textarea {
            /deep/ textarea {
                min-height: auto !important;
                line-height: 22px;
                @include scrollbar-y;
            }
        }
    }
    .add-template {
        width: 20%;
        line-height: 38px;
        cursor: pointer;
        color: #63656E;
        font-size: 12px;
        .icon-plus-circle {
            margin-top: -2px;
            font-size: 14px;
            color: #979BA5;
        }
    }
    .node-create-options {
        padding: 9px 20px;
        border-top: 1px solid $cmdbBorderColor;
        text-align: right;
        background-color: #FAFBFD;
        font-size: 0;
    }
</style>
