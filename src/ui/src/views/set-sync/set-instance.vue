<template>
    <div class="set-instance-layout"
        :class="{
            'borderBottom': !localExpand,
            'cant-sync': !canSyncStatus
        }">
        <div class="title" @click="localExpand = !localExpand">
            <div class="left-info">
                <i class="bk-icon icon-right-shape" :class="{ 'is-expand': localExpand }"></i>
                <h2 :class="['path', { 'is-read': hasRead }]">{{topoPath}}</h2>
                <span class="sync-status-tips" v-if="!canSyncStatus">{{$t('不可同步')}}</span>
            </div>
            <i class="bk-icon icon-close"
                v-if="iconClose"
                v-bk-tooltips="$t('本次不同步')"
                @click.stop="handleClose">
            </i>
        </div>
        <cmdb-collapse-transition>
            <div class="main clearfix" v-show="localExpand">
                <div class="sync fl">
                    <h3>{{$t('同步前')}}</h3>
                    <div class="sync-main fl">
                        <div class="sync-title"
                            :class="{ 'is-expand': beforeSyncExpand }"
                            :title="setDeatails.bk_set_name"
                            @click.stop="beforeSyncExpand = !beforeSyncExpand">
                            <i class="bk-icon icon-right-shape"></i>
                            <i class="sync-icon">{{$i18n.locale === 'en' ? 's' : '集'}}</i>
                            <span class="set-name">{{setDeatails.bk_set_name}}</span>
                        </div>
                        <cmdb-collapse-transition>
                            <div v-show="beforeSyncExpand">
                                <ul class="sync-info">
                                    <li class="mt15"
                                        v-for="_module in beforeChangeList"
                                        :key="_module.bk_module_id"
                                        :title="_module.bk_module_name">
                                        <i class="sync-icon">{{$i18n.locale === 'en' ? 'm' : '模'}}</i>
                                        <span class="name">{{_module.bk_module_name}}</span>
                                    </li>
                                </ul>
                            </div>
                        </cmdb-collapse-transition>
                    </div>
                </div>
                <div class="sync fl sync-after">
                    <h3>{{$t('同步后')}}</h3>
                    <div class="sync-main fl">
                        <div class="sync-title"
                            :class="{ 'is-expand': afterSyncExpand }"
                            :title="setDeatails.bk_set_name"
                            @click.stop="afterSyncExpand = !afterSyncExpand">
                            <i class="bk-icon icon-right-shape"></i>
                            <i class="sync-icon">{{$i18n.locale === 'en' ? 's' : '集'}}</i>
                            <span class="set-name">{{setDeatails.bk_set_name}}</span>
                        </div>
                        <cmdb-collapse-transition>
                            <div v-show="afterSyncExpand">
                                <ul class="sync-info">
                                    <li v-for="_module in instance.module_diffs"
                                        :key="_module.bk_module_id + _module.bk_module_name"
                                        :class="['mt15', {
                                            'has-delete': _module.diff_type === 'remove',
                                            'has-changed': _module.diff_type === 'changed',
                                            'new-add': _module.diff_type === 'add'
                                        }]">
                                        <i class="sync-icon" :title="_module.bk_module_name">{{$i18n.locale === 'en' ? 'm' : '模'}}</i>
                                        <span class="name" :title="_module.bk_module_name">{{_module.bk_module_name}}</span>
                                        <div class="tips" v-if="_module.diff_type === 'remove' && existHost(_module.bk_module_id)">
                                            <i class="bk-icon icon-exclamation"></i>
                                            <i18n path="存在主机不可同步提示" tag="p">
                                                <span place="btn" class="view-btn" @click="handleViewModule(_module.bk_module_id)">{{$t('跳转查看')}}</span>
                                            </i18n>
                                        </div>
                                    </li>
                                </ul>
                            </div>
                        </cmdb-collapse-transition>
                    </div>
                </div>
            </div>
        </cmdb-collapse-transition>
    </div>
</template>

<script>
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    export default {
        props: {
            instance: {
                type: Object,
                required: true
            },
            expand: {
                type: Boolean,
                default: false
            },
            iconClose: {
                type: Boolean,
                default: true
            },
            moduleHostCount: {
                type: Object,
                default: () => ({})
            }
        },
        data () {
            return {
                localExpand: this.expand,
                beforeSyncExpand: true,
                afterSyncExpand: true,
                hasRead: this.expand
            }
        },
        computed: {
            beforeChangeList () {
                return this.instance.module_diffs.filter(_module => _module.diff_type !== 'add')
            },
            setDeatails () {
                return this.instance.set_detail
            },
            topoPath () {
                const path = this.instance.topo_path
                if (path.length) {
                    const topoPath = this.$tools.clone(path)
                    return topoPath.reverse().map(path => path.bk_inst_name).join(' / ')
                }
                return '--'
            },
            canSyncStatus () {
                for (const _module of this.instance.module_diffs) {
                    if (_module.diff_type === 'remove' && this.moduleHostCount[_module.bk_module_id] > 0) {
                        return false
                    }
                }
                return true
            }
        },
        watch: {
            localExpand (value) {
                if (value && !this.hasRead) {
                    this.hasRead = true
                }
            }
        },
        methods: {
            existHost (moduleId) {
                return this.moduleHostCount[moduleId] > 0
            },
            handleClose () {
                this.$emit('close', this.instance.bk_set_id)
            },
            handleViewModule (moduleId) {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: {
                        node: `module-${moduleId}`
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .set-instance-layout {
        border: 1px solid #DCDEE5;
        background-color: #ffffff;
        &.borderBottom {
            border-bottom: none;
        }
        &.cant-sync {
            border-color: #F0F1F5;
            .title {
                background-color: #F0F1F5;
                .left-info {
                    color: #979BA5;
                }
            }
            .icon-right-shape {
                color: #C4C6CC;
            }
        }
    }
    .title {
        @include space-between;
        height: 42px;
        padding: 0 10px;
        background-color: #DCDEE5;
        cursor: pointer;
        .left-info {
            display: flex;
            align-items: center;
            color: #313238;
            font-size: 14px;
        }
        .path {
            position: relative;
            font-size: 14px;
            padding: 0 10px;
            font-weight: normal;
            &:not(.is-read)::before {
                content: '';
                position: absolute;
                top: 0;
                right: 2px;
                width: 6px;
                height: 6px;
                background-color: #FF5656;
                border-radius: 50%;
            }
        }
        .sync-status-tips {
            @include inlineBlock;
            height: 20px;
            line-height: 20px;
            font-size: 12px;
            padding: 0 5px;
            color: #FFFFFF;
            background-color: #FE9C9C;
            margin: 2px 0 0 12px;
        }
        .icon-right-shape {
            color: #63656E;
            transition: all .5s;
            &.is-expand {
                transform: rotateZ(90deg);
            }
        }
        .icon-close {
            color: #979BA5;
            font-size: 16px;
            font-weight: bold;
            margin-right: 2px;
            cursor: pointer;
        }
    }
    .main {
        padding: 14px 30px 26px;
    }
    .sync {
        min-width: 320px;
        max-width: 380px;
        color: #63656E;
        font-size: 14px;
        &.sync-after {
            max-width: calc(100% - 420px);
            margin-left: 40px;
            .sync-main {
                max-width: 490px;
            }
        }
        > h3 {
            font-size: 14px;
            float: left;
            margin-right: 20px;
        }
        .sync-icon {
            @include inlineBlock;
            width: 20px;
            height: 20px;
            line-height: 19px;
            font-size: 12px;
            font-style: normal;
            text-align: center;
            color: #FFFFFF;
            border-radius: 50%;
            background-color: #97AED6;
            margin-right: 7px;
        }
        .sync-title {
            display: flex;
            align-items: center;
            cursor: pointer;
            .set-name {
                max-width: 200px;
                @include ellipsis;
            }
            &.is-expand {
                .icon-right-shape {
                    transform: rotateZ(90deg);
                }
            }
            .bk-icon {
                color: #C4C6CC;
                &.icon-right-shape {
                    transition: all .5s;
                    margin-right: 10px;
                }
            }
        }
        .sync-info {
            padding-left: 48px;
            li {
                display: flex;
                align-items: center;
                .name {
                    @include inlineBlock;
                    @include ellipsis;
                    flex: 1;
                    max-width: 200px;
                    line-height: 20px;
                    font-size: 14px;
                }
                .tips {
                    display: flex;
                    align-items: center;
                    font-size: 12px;
                    color: #FF5656;
                    margin: 4px 0 0 6px;
                    .bk-icon {
                        width: 16px;
                        height: 16px;
                        line-height: 16px;
                        text-align: center;
                        color: #FFFFFF;
                        background-color: #FF5656;
                        border-radius: 50%;
                        margin-right: 4px;
                    }
                    .view-btn {
                        color: #3A84FF;
                        cursor: pointer;
                    }
                }
                &.has-delete {
                    color: #C4C6CC;
                    .sync-icon {
                        background-color: #DCDEE5;
                    }
                    .name {
                        text-decoration: line-through;
                    }
                }
                &.has-changed {
                    color: #3A84FF;
                    .sync-icon {
                        background-color: #3A84FF;
                    }
                }
                &.new-add {
                    color: #2DCB56;
                    .sync-icon {
                        background-color: #2DCB56;
                    }
                }
            }
        }
    }
</style>
