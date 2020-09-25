<template>
    <div class="resource-layout clearfix">
        <bk-tab
            :active.sync="activeTab"
            class="scope-tab"
            type="unborder-card"
            @tab-change="handleTabChange">
            <bk-tab-panel v-for="item in scopeList"
                :key="item.id"
                :name="item.id"
                :label="item.label">
            </bk-tab-panel>
        </bk-tab>
        <div class="content">
            <cmdb-resize-layout
                v-if="isResourcePool"
                :class="['resize-layout fl', { 'is-collapse': layout.collapse }]"
                :handler-offset="3"
                :min="200"
                :max="480"
                :disabled="layout.collapse"
                direction="right">
                <resource-directory></resource-directory>
                <i class="directory-collapse-icon bk-icon icon-angle-left"
                    @click="layout.collapse = !layout.collapse">
                </i>
            </cmdb-resize-layout>
            <resource-hosts class="main"></resource-hosts>
        </div>
        <router-subview></router-subview>
    </div>
</template>

<script>
    import resourceDirectory from './children/directory.vue'
    import resourceHosts from './children/host-list.vue'
    import Bus from '@/utils/bus.js'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            resourceDirectory,
            resourceHosts
        },
        data () {
            return {
                layout: {
                    collapse: false
                },
                activeTab: RouterQuery.get('scope', '1'),
                scopeList: [{
                    id: '1',
                    label: this.$t('未分配')
                }, {
                    id: '0',
                    label: this.$t('已分配')
                }, {
                    id: 'all',
                    label: this.$t('全部')
                }]
            }
        },
        computed: {
            isResourcePool () {
                return this.activeTab.toString() === '1'
            }
        },
        methods: {
            handleTabChange (tab) {
                Bus.$emit('toggle-host-filter', false)
                Bus.$emit('reset-host-filter')
                RouterQuery.set({
                    scope: isNaN(tab) ? tab : parseInt(tab),
                    ip: '',
                    bk_asset_id: '',
                    page: 1,
                    _t: Date.now()
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .resource-layout{
        .scope-tab {
            height: auto;
            margin: 0 20px;
            /deep/ .bk-tab-header {
                padding: 0;
            }
        }
        .content {
            height: calc(100% - 58px);
            overflow: hidden;
        }
        .resize-layout {
            position: relative;
            width: 280px;
            height: 100%;
            border-right: 1px solid $cmdbLayoutBorderColor;
            &.is-collapse {
                width: 0 !important;
                border-right: none;
                .directory-collapse-icon:before {
                    display: inline-block;
                    transform: rotate(180deg);
                }
            }
            .directory-collapse-icon {
                position: absolute;
                left: 100%;
                top: 50%;
                width: 16px;
                height: 100px;
                line-height: 100px;
                background: $cmdbLayoutBorderColor;
                border-radius: 0px 12px 12px 0px;
                transform: translateY(-50%);
                text-align: center;
                text-indent: -2px;
                font-size: 20px;
                color: #fff;
                cursor: pointer;
                &:hover {
                    background: #699DF4;
                }
            }
        }
        .main {
            height: 100%;
            padding: 10px 20px 0 20px;
            overflow: hidden;
        }
    }
    .assign-dialog {
        /deep/ .bk-dialog-body {
            padding: 0 50px 40px;
        }
        .assign-info span {
            color: #3c96ff;
        }
        .assign-footer {
            padding-top: 20px;
            font-size: 0;
            text-align: center;
            .bk-button-normal {
                width: 96px;
            }
        }
    }
</style>
