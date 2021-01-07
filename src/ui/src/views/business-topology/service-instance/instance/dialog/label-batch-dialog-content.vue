<template>
    <div class="batch-edit-label">
        <div class="exisiting-label">
            <div class="title">
                {{$t('已有标签')}}
                <span>{{$t('已有标签提示')}}</span>
            </div>
            <div class="label-set">
                <span class="label-item" :key="index" v-for="(label, index) in localLabels">
                    {{`${label.key}：${label.value}`}}
                    <i class="bk-icon icon-close-circle-shape" @click="handleRemove(index)"></i>
                </span>
            </div>
        </div>
        <div class="batch-add">
            <div class="title">{{$t('批量添加标签')}}</div>
            <slot name="batch-add-label"></slot>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            labels: {
                type: Array,
                default: () => []
            }
        },
        data () {
            return {
                localLabels: [...this.labels],
                removeLabels: []
            }
        },
        methods: {
            handleRemove (index) {
                this.removeLabels.push(this.localLabels[index])
                this.localLabels.splice(index, 1)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .batch-edit-label {
        height: 321px;
        padding: 6px 16px 20px 26px;
        .title {
            color: #63656e;
            font-size: 14px;
            padding-bottom: 8px;
            span {
                color: #979ba5;
                font-size: 12px;
            }
        }
        .exisiting-label {
            padding-bottom: 12px;
            .label-set {
                &::-webkit-scrollbar, &::-webkit-scrollbar-thumb {
                    display: none;
                }
                &:hover {
                    &::-webkit-scrollbar, &::-webkit-scrollbar-thumb {
                        display: block;
                    }
                }
                @include scrollbar-y;
                height: 84px;
                padding: 12px 0 2px 12px;
                border-radius: 2px;
                border: 1px solid #c4c6cc;
                .label-item {
                    display: inline-block;
                    height: 22px;
                    line-height: 20px;
                    font-size: 12px;
                    padding: 0 7px;
                    margin: 0 10px 10px 0;
                    color: #63656e;
                    background-color: #f0f1f5;
                    border-radius: 12px;
                }
                .bk-icon {
                    color: #c4c6cc;
                    margin-left: 4px;
                    cursor: pointer;
                    &:hover {
                        color: #979ba5;
                    }
                }
            }
        }
        .edit-label {
            padding: 0;
            /deep/ .scrollbar-box {
                height: 135px !important;
            }
        }
    }
</style>

<style lang="scss">
</style>
