<template>
    <div class="create-node-layout">
        <h2 class="node-create-title">{{$t('新增子节点')}}</h2>
        <div class="node-create-path" :title="topoPath">{{$t('添加节点已选择')}}：{{topoPath}}</div>
        <div class="node-create-form">
            <div class="form-item">
                <label>{{$t('节点名称')}} <font color="red">*</font></label>
                <bk-input class="form-textarea"
                    type="textarea"
                    data-vv-name="nodeName"
                    v-validate="'required|singlechar|length:256'"
                    v-model="nodeName"
                    :rows="rows"
                    :placeholder="$t('请输入节点名称，多个同级节点换行分隔')"
                    @keydown="handleKeydown">
                </bk-input>
                <span class="form-error" v-if="errors.has('nodeName')">{{errors.first('nodeName')}}</span>
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
                nodeName: '',
                rows: 1
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
                const rows = this.nodeName.split('\n').length
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
        padding: 12px 26px 32px;
    }
    .form-item {
        margin: 15px 0 0 0;
        position: relative;
        label {
            display: block;
            padding: 0 0 10px;
            line-height: 19px;
            font-size: 14px;
            color: #63656E;
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
    .node-create-options {
        padding: 9px 20px;
        border-top: 1px solid $cmdbBorderColor;
        text-align: right;
        background-color: #FAFBFD;
        font-size: 0;
    }
</style>
