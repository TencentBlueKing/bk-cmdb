<template>
    <div class="cagetory-wrapper" :style="{ 'padding-top': showFeatureTips ? '10px' : '' }">
        <feature-tips
            :feature-name="'cagetory'"
            :show-tips="showFeatureTips"
            :desc="$t('ServiceCagetory[\'服务分类功能提示\']')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="cagetory-list">
            <div class="cagetory-item">
                <div class="cagetory-title">
                    <div class="cagetory-name">
                        <span>数据库</span>
                        <i class="bk-icon icon-edit2"></i>
                    </div>
                    <cmdb-dot-menu class="dot-menu">
                        <div class="menu-operational">
                            <i>添加</i>
                            <i>编辑</i>
                            <i>删除</i>
                        </div>
                    </cmdb-dot-menu>
                </div>
                <div class="child-cagetory">
                    <div class="child-item child-edit">
                        <!-- <div class="child-title">Mysql</div>
                        <span>8</span> -->
                        <div class="edit-box">
                            <input type="text" class="bk-form-input" placeholder="请输入二级分类">
                            <span class="text-primary btn-confirm">确定</span>
                            <span class="text-primary">取消</span>
                        </div>
                    </div>
                    <div class="child-item">
                        <div class="child-title">
                            <span>Mysql</span>
                            <div class="child-edit">
                                <i class="bk-icon icon-edit2 mr10"></i>
                                <i class="icon-cc-tips-close"></i>
                            </div>
                        </div>
                        <span>8</span>
                    </div>
                    <div class="child-item">
                        <div class="child-title">Mysql</div>
                        <span>8</span>
                    </div>
                    <div class="child-item">
                        <div class="child-title">Mysql</div>
                        <span>8</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    export default {
        components: {
            featureTips
        },
        data () {
            return {
                showFeatureTips: false
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Nav["服务分类"]'))
            this.showFeatureTips = this.featureTipsParams['cagetory']
        }
    }
</script>

<style lang="scss" scoped>
    .cagetory-wrapper {
        .cagetory-list {
            display: flex;
            flex-wrap: wrap;
        }
        .cagetory-item {
            width: 294px;
            border: 1px solid #dcdee5;
            margin-right: 30px;
            margin-bottom: 20px;
        }
        .cagetory-title {
            @include space-between;
            background-color: #fafbfd;
            padding: 0 20px;
            height: 42px;
            line-height: 42px;
            font-size: 14px;
            color: #63656e;
            font-weight: bold;
            border-bottom: 1px solid #dcdee5;
            &:hover {
                background-color: #f0f1f5;
                .bk-icon {
                    display: inline !important;
                }
            }
            .cagetory-name {
                @include ellipsis;
                flex: 1;
                padding-right: 20px;
                .bk-icon {
                    display: none;
                    cursor: pointer;
                    color: #3a84ff;
                }
            }
            .dot-menu {
                cursor: pointer;
            }
        }
        .child-cagetory {
            padding: 0 10px 10px 38px;
            .child-item {
                @include space-between;
                position: relative;
                z-index: 10;
                line-height: 32px;
                &.child-edit {
                    .bk-form-input {
                        font-size: 12px;
                        width: 170px;
                        height: 32px;
                        margin-right: 4px;
                    }
                    .edit-box .text-primary {
                        display: inline-block;
                        line-height: normal;
                        font-size: 12px;
                        &.btn-confirm {
                            position: relative;
                            margin-right: 6px;
                            &::after {
                                content: '';
                                position: absolute;
                                top: 2px;
                                right: -6px;
                                display: inline-block;
                                width: 1px;
                                height: 14px;
                                background-color: #dcdee5;
                            }
                        }
                    }
                    &:first-child::after {
                        height: 32px;
                    }
                }
                &:hover {
                    .child-title {
                        padding-right: 10px;
                        background-color: #fafbfd;
                        color: #3a84ff;
                    }
                    >span {
                        display: none;
                    }
                    .child-edit {
                        display: block;
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
                    @include space-between;
                    color: #63656e;
                    flex: 1;
                    padding-right: 20px;
                    padding-left: 16px;
                    margin-left: 10px;
                }
                >span {
                    color: #c4c6cc;
                    padding-right: 18px;
                }
                .child-edit {
                    display: none;
                    i {
                        color: #3a84ff;
                        cursor: pointer;
                    }
                }
            }
        }
    }
    .menu-operational {
        width: 68px;
        padding: 6px 0;
        font-size: 12px;
        text-align: center;
        line-height: 32px;
        color: #c4c6cc;
        i {
            display: block;
            font-style: normal;
            cursor: pointer;
            &:hover {
                color: #3a84ff;
                background-color: #e1ecff;
            }
        }
    }
</style>
