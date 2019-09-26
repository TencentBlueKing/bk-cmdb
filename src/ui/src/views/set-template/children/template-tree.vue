<template>
    <div class="template-tree">
        <div class="node-root clearfix">
            <i class="folder-icon bk-icon icon-down-shape fl" @click="handleCollapse"></i>
            <i class="root-icon icon icon-cc-nav-model-02 fl"></i>
            <span class="root-name" :title="templateName">{{templateName}}</span>
        </div>
        <cmdb-collapse-transition>
            <ul class="node-children" v-show="!collapse">
                <li class="node-child clearfix"
                    v-for="module in modules"
                    :key="module.id"
                    :class="{ selected: selected === module.id }"
                    @click="handleChildClick(module)">
                    <i class="child-icon icon icon-cc-cube fl"></i>
                    <span class="child-options fr">
                        <i class="options-view icon icon-cc-show" @click="handleViewService"></i>
                        <i class="options-delete icon icon-cc-tips-close" @click="handleDeleteService"></i>
                    </span>
                    <span class="child-name">{{module.name}}</span>
                </li>
                <li class="options-child node-child clearfix"
                    v-if="['create', 'edit'].includes(mode)"
                    @click="handleAddService">
                    <i class="child-icon icon icon-cc-zoom-in fl"></i>
                    <span class="child-name">{{$t('添加服务模板')}}</span>
                </li>
            </ul>
        </cmdb-collapse-transition>
    </div>
</template>

<script>
    export default {
        props: {
            mode: {
                type: String,
                default: 'create',
                validator (value) {
                    return ['create', 'edit', 'view'].includes(value)
                }
            }
        },
        data () {
            return {
                templateName: this.$t('模板集群名称'),
                modules: [{ id: 1, name: 'gameserver' }, { id: 2, name: 'gameserver' }],
                collapse: false,
                selected: null,
                unwatch: null
            }
        },
        watch: {
            mode: {
                immediate: true,
                handler (value, oldValue) {
                    if (value === 'create') {
                        this.initMonitorTemplateName()
                    }
                }
            }
        },
        beforeDestory () {
            this.unwatch && this.unwatch()
        },
        methods: {
            initMonitorTemplateName () {
                this.unwatch = this.$watch(() => {
                    return this.$parent.templateName
                }, value => {
                    if (value) {
                        this.templateName = value
                    } else {
                        this.templateName = this.$t('模板集群名称')
                    }
                })
            },
            handleCollapse () {
                this.collapse = !this.collapse
            },
            handleChildClick (module) {
                this.selected = module.id
            },
            handleAddService () {},
            handleViewService () {},
            handleDeleteService () {}
        }
    }
</script>

<style lang="scss" scoped>
    $iconColor: #C4C6CC;
    $fontColor: #63656E;
    .template-tree {
        padding: 10px 0 10px 20px;
        border: 1px dashed #C4C6CC;
        background-color: #fff;
    }
    .node-root {
        line-height: 36px;
        cursor: default;
        .folder-icon {
            width: 23px;
            height: 36px;
            line-height: 36px;
            text-align: center;
            font-size: 12px;
            color: $iconColor;
            cursor: pointer;
        }
        .root-icon {
            position: relative;
            margin: 8px 4px 8px 0px;
            font-size: 20px;
            color: $iconColor;
            background-color: #fff;
            z-index: 2;
        }
        .root-name {
            display: block;
            padding: 0 10px 0 0;
            font-size: 14px;
            color: $fontColor;
            @include ellipsis;
        }
    }
    .node-children {
        line-height: 36px;
        margin-left: 32px;
        cursor: default;
        .node-child {
            padding: 0 10px 0 32px;
            position: relative;
            &:hover {
                background-color: rgba(240,241,245, .6);
            }
            &.selected {
                background-color: #F0F1F5;
            }
            &:hover,
            &.selected {
                .child-options {
                    display: block;
                }
            }
            &:before {
                position: absolute;
                left: 0px;
                top: -18px;
                content: "";
                width: 25px;
                height: 36px;
                border-left: 1px dashed #DCDEE5;
                border-bottom: 1px dashed #DCDEE5;
                z-index: 1;
            }
            &.options-child {
                cursor: pointer;
                .child-icon {
                    font-size: 18px;
                    background-color: #fff;
                    color: #3A84FF;
                }
                .child-name {
                    color: #3A84FF;
                }
            }
            .child-name {
                display: block;
                padding: 0 10px 0 0;
                font-size: 14px;
                color: $fontColor;
                @include ellipsis;
            }
            .child-icon {
                width: 20px;
                height: 20px;
                margin: 8px 7px 0 0;
                border-radius: 50%;
                text-align: center;
                color: #fff;
                line-height: 20px;
                font-size: 12px;
                background-color: #3A84FF;
            }
            .child-options {
                display: none;
                margin-right: 9px;
                font-size: 0;
                .options-view {
                    font-size: 18px;
                    color: #3A84FF;
                    cursor: pointer;
                    &:hover {
                        opacity: .75;
                    }
                }
                .options-delete {
                    width: 24px;
                    height: 24px;
                    margin-left: 14px;
                    font-size: 12px;
                    text-align: center;
                    line-height: 24px;
                    color: $iconColor;
                    cursor: pointer;
                    &:hover {
                        opacity: .75;
                    }
                }
            }
        }
    }
</style>
