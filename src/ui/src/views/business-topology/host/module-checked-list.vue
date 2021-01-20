<template>
    <div class="module-checked-list">
        <div class="selected-info">
            <div class="selected-title">{{$t('结果预览')}}</div>
            <div class="selected-count">
                <i18n path="已选择N个模块">
                    <em class="count" place="count">{{checked.length}}</em>
                </i18n>
                <bk-button text theme="primary"
                    v-show="checked.length"
                    @click="handleClearModule">
                    {{$t('清空')}}
                </bk-button>
            </div>
        </div>
        <ul class="module-list">
            <li class="module-item" v-for="node in checked"
                :key="node.id">
                <div class="module-info">
                    <span class="info-icon">{{node.data.bk_obj_name[0]}}</span>
                    <span class="info-name" :title="node.data.bk_inst_name">
                        {{node.data.bk_inst_name}}
                    </span>
                </div>
                <div class="module-topology" :title="getNodePath(node)">{{getNodePath(node)}}</div>
                <i class="bk-icon icon-close" @click="handleDeleteModule(node)"></i>
            </li>
        </ul>
    </div>
</template>

<script>
    export default {
        props: {
            checked: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        methods: {
            handleDeleteModule (node) {
                this.$emit('delete', node)
            },
            handleClearModule () {
                this.$emit('clear')
            },
            getNodePath (node) {
                const parents = node.parents
                return parents.map(parent => parent.data.bk_inst_name).join(' / ')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .module-checked-list {
        height: 100%;
        .selected-info {
            .selected-title {
                color: #313238;
                font-size: 14px;
                line-height: 22px;
                margin-bottom: 10px;
            }

            .selected-count {
                display: flex;
                justify-content: space-between;
                font-size: 12px;
                line-height: 20px;
                color: $textColor;
                .count {
                    padding: 0 4px;
                    font-weight: bold;
                    font-style: normal;
                    color: #3A84FF;
                }
            }

            /deep/ .bk-button-text {
                font-size: 12px;
            }
        }
    }

    .module-list {
        height: calc(100% - 58px);
        margin-top: 4px;
        @include scrollbar-y;
        .module-item {
            position: relative;
            margin-top: 2px;
            background: #fff;
            padding: 4px 12px;
            border-radius: 2px;
            box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
            .icon-close {
                display: none;
                position: absolute;
                color: #3a84ff;
                font-size: 22px;
                top: 8px;
                right: 0px;
                width: 28px;
                height: 28px;
                line-height: 28px;
                text-align: center;
                cursor: pointer;
            }
            &:hover {
                .icon-close {
                    display: block;
                }
            }
        }
    }
    .module-info {
        display: flex;
        align-items: center;
        .info-icon {
            flex: none;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #fff;
        }
        .info-name {
            padding-left: 10px;
            font-size: 12px;
            color: $textColor;
            line-height: 20px;
            @include ellipsis;
        }
    }
    .module-topology {
        padding-left: 30px;
        font-size: 12px;
        color: #C4C6CC;
        @include ellipsis;
    }
</style>
