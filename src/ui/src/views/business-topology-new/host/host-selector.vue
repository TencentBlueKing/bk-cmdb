<template>
    <div class="layout clearfix" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <div class="wrapper-left fl">
            <h2 class="title">{{$t('选择主机')}}</h2>
            <bk-select class="selector-type" v-model="type" :clearable="false">
                <bk-option id="topology" :name="$t('业务拓扑')"></bk-option>
                <bk-option id="custom" name="IP"></bk-option>
            </bk-select>
            <keep-alive>
                <component :is="activeComponent" class="selector-component"></component>
            </keep-alive>
        </div>
        <div class="wrapper-right fl">
            <div class="selected-count">
                <i18n path="已选择N台主机">
                    <span class="count">{{selected.length}}</span>
                </i18n>
            </div>
            <bk-table
                :data="selected"
                :outer-border="false"
                :header-border="false"
                :header-cell-style="{ background: '#fff' }"
                :height="367">
                <bk-table-column :label="$t('内网IP')">
                    <template slot-scope="{ row }">{{row.host.bk_host_innerip}}</template>
                </bk-table-column>
                <bk-table-column :label="$t('云区域')">
                    <template slot-scope="{ row }">{{row.host.bk_cloud_id | foreignkey}}</template>
                </bk-table-column>
                <bk-table-column :label="$t('操作')">
                    <bk-button slot-scope="{ row }" text @click="handleRemove(row)">{{$t('移除')}}</bk-button>
                </bk-table-column>
            </bk-table>
        </div>
        <div class="clearfix"></div>
        <div class="wrapper-footer">
            <bk-button class="mr10" theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
            <bk-button theme="primary">{{$t('下一步')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { foreignkey } from '@/filters/formatter.js'
    import HostSelectorTopology from './host-selector-topology.vue'
    import HostSelectorCustom from './host-selector-custom.vue'
    export default {
        components: {
            HostSelectorTopology,
            HostSelectorCustom
        },
        filters: {
            foreignkey
        },
        data () {
            return {
                components: {
                    topology: HostSelectorTopology,
                    custom: HostSelectorCustom
                },
                activeComponent: null,
                type: 'custom',
                filter: '',
                selected: [],
                request: {
                    internal: Symbol('topology'),
                    instance: Symbol('instance'),
                    host: Symbol('host')
                }
            }
        },
        watch: {
            type: {
                immediate: true,
                handler (type) {
                    this.activeComponent = this.components[type]
                }
            }
        },
        methods: {
            handleRemove (row) {
                this.selected = this.selected.filter(target => target.host.bk_host_id !== row.host.bk_host_id)
            },
            handleSelect (item) {
                const isExist = this.selected.some(target => target.host.bk_host_id === item.host.bk_host_id)
                if (!isExist) {
                    this.selected = [...this.selected, item]
                }
            },
            handleCancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .layout {
        width: 850px;
        height: 460px;
    }
    .wrapper-left {
        width: 240px;
        height: 410px;
        border-right: 1px solid $borderColor;
        .title {
            padding: 0 20px;
            margin: 15px 0 20px 0;
            line-height:26px;
            font-size:20px;
            font-weight: normal;
            color: $textColor;
        }
        .selector-type {
            display: block;
            margin: 0px 20px;
        }
        .selector-component {
            margin-top: 10px;
            height: 307px;
        }
    }
    .wrapper-right {
        width: 610px;
        height: 410px;
        .selected-count {
            height: 42px;
            line-height: 42px;
            padding: 0 12px;
            background-color: #FAFBFD;
            font-size: 12px;
            font-weight: bold;
            color: $textColor;
        }
    }
    .wrapper-footer {
        height: 50px;
        border-top: 1px solid $borderColor;
        font-size: 0;
        text-align: right;
        background-color: #FAFBFD;
        padding: 8px 20px 9px;
    }
</style>
