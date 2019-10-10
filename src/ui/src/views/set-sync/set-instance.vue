<template>
    <div class="set-instance-layout" :class="{ 'borderBottom': !localExpand }">
        <div class="title" @click="localExpand = !localExpand">
            <div class="left-info">
                <i class="bk-icon icon-right-shape" :class="{ 'is-expand': localExpand }"></i>
                <h2 class="path">{{topoPath}}</h2>
                <span :class="['count', { 'is-read': hasRead }]">{{changeCount}}</span>
            </div>
            <i v-show="iconClose" class="bk-icon icon-close" @click.stop="handleClose"></i>
        </div>
        <cmdb-collapse-transition>
            <div class="main clearfix" v-show="localExpand">
                <div class="sync fl">
                    <h3>{{$t('同步前')}}</h3>
                    <div class="sync-main fl">
                        <div class="sync-title" :class="{ 'is-expand': beforeSyncExpand }" @click.stop="beforeSyncExpand = !beforeSyncExpand">
                            <i class="bk-icon icon-right-shape"></i>
                            <i class="bk-icon icon-cc-nav-model-02"></i>
                            <span class="set-name">{{setDeatails.bk_set_name}}</span>
                        </div>
                        <cmdb-collapse-transition>
                            <div v-show="beforeSyncExpand">
                                <ul class="sync-info"
                                    v-for="_module in beforeChangeList"
                                    :key="_module.bk_module_id">
                                    <li class="mt20">
                                        <i class="bk-icon icon-cc-nav-model-02"></i>
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
                        <div class="sync-title" :class="{ 'is-expand': afterSyncExpand }" @click.stop="afterSyncExpand = !afterSyncExpand">
                            <i class="bk-icon icon-right-shape"></i>
                            <i class="bk-icon icon-cc-nav-model-02"></i>
                            <span class="set-name">{{setDeatails.bk_set_name}}</span>
                        </div>
                        <cmdb-collapse-transition>
                            <div v-show="afterSyncExpand">
                                <ul class="sync-info"
                                    v-for="_module in instance.module_diffs"
                                    :key="_module.bk_module_id">
                                    <li :class="['mt20', {
                                        'has-delete': _module.diff_type === 'remove',
                                        'has-changed': _module.diff_type === 'changed',
                                        'new-add': _module.diff_type === 'add'
                                    }]">
                                        <i class="bk-icon icon-cc-nav-model-02"></i>
                                        <span class="name">{{_module.bk_module_name}}</span>
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
            changeCount () {
                const changeList = this.instance.module_diffs.filter(_module => _module.diff_type !== 'unchanged')
                return changeList.length
            },
            setDeatails () {
                return this.instance.set_detail
            },
            topoPath () {
                const path = this.instance.topo_path
                if (path.length) {
                    const topoPath = this.$tools.clone(path)
                    return topoPath.reverse().map(path => path.InstanceName).join(' / ')
                }
                return '--'
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
            handleClose () {
                this.$emit('close', this.instance.bk_set_id)
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
    }
    .title {
        @include space-between;
        height: 42px;
        padding: 0 10px;
        background-color: #F0F1F5;
        border-bottom: 1px solid #DCDEE5;
        cursor: pointer;
        .left-info {
            display: flex;
            align-items: center;
            color: #63656E;
            font-size: 14px;
        }
        .path {
            font-size: 14px;
            padding: 0 10px;
        }
        .icon-right-shape {
            transition: all .5s;
            &.is-expand {
                transform: rotateZ(90deg);
            }
        }
        .count {
            min-width: 18px;
            height: 18px;
            line-height: 16px;
            font-size: 12px;
            color: #ffffff;
            text-align: center;
            padding: 0 5px;
            background-color: #FF5656;
            border-radius: 100px;
            &.is-read {
                background-color: #C4C6CC;
            }
        }
        .icon-close {
            font-size: 16px;
            font-weight: bold;
            margin-right: 2px;
            cursor: pointer;
        }
    }
    .main {
        padding: 14px 42px 26px;
    }
    .sync {
        min-width: 320px;
        color: #63656E;
        font-size: 14px;
        &.sync-after {
            margin-left: 50px;
        }
        > h3 {
            font-size: 14px;
            float: left;
            margin-right: 30px;
        }
        .sync-title {
            display: flex;
            align-items: center;
            cursor: pointer;
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
                &.icon-cc-nav-model-02 {
                    font-size: 18px;
                    margin-right: 6px;
                }
            }
        }
        .sync-info {
            padding-left: 48px;
            li {
                .bk-icon {
                    width: 20px;
                    height: 20px;
                    line-height: 20px;
                    text-align: center;
                    color: #ffffff;
                    font-size: 12px;
                    background-color: #C4C6CC;
                    border-radius: 50%;
                }
                .name {
                    display: inline-block;
                    line-height: 20px;
                    padding-left: 2px;
                }
                &.has-delete {
                    color: #FF5656;
                    .bk-icon {
                        background-color: #FF5656;
                    }
                    .name {
                        text-decoration: line-through;
                    }
                }
                &.has-changed {
                    color: #3A84FF;
                    .bk-icon {
                        background-color: #3A84FF;
                    }
                }
                &.new-add {
                    color: #2DCB56;
                    .bk-icon {
                        background-color: #2DCB56;
                    }
                }
            }
        }
    }
</style>
