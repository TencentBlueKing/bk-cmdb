<template>
    <div class="layout clearfix"
        v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <div class="wrapper clearfix">
            <div class="wrapper-column wrapper-left fl">
                <h2 class="title">{{$t('选择主机')}}</h2>
                <bk-select class="selector-type" v-model="type" :clearable="false">
                    <bk-option id="topology" :name="$t('业务拓扑')"></bk-option>
                    <bk-option id="custom" name="IP"></bk-option>
                </bk-select>
                <keep-alive>
                    <component :is="activeComponent" class="selector-component" ref="dynamic"></component>
                </keep-alive>
            </div>
            <div class="wrapper-column wrapper-right fl">
                <div class="selected-count">
                    <i18n path="已选择N台主机">
                        <span class="count" place="count">{{selected.length}}</span>
                    </i18n>
                </div>
                <bk-table
                    :data="selected"
                    :outer-border="false"
                    :header-border="false"
                    :header-cell-style="{ background: '#fff' }"
                    :height="369">
                    <bk-table-column :label="$t('内网IP')">
                        <template slot-scope="{ row }">
                            {{row.host.bk_host_innerip}}
                            <span class="repeat-row" v-if="repeatSelected.includes(row)">{{$t('IP重复')}}</span>
                        </template>
                    </bk-table-column>
                    <bk-table-column :label="$t('云区域')">
                        <template slot-scope="{ row }">{{row.host.bk_cloud_id | foreignkey}}</template>
                    </bk-table-column>
                    <bk-table-column :label="$t('操作')">
                        <bk-button slot-scope="{ row }" text @click="handleRemove(row)">{{$t('移除')}}</bk-button>
                    </bk-table-column>
                </bk-table>
            </div>
        </div>
        <div class="layout-footer">
            <bk-button class="mr10" theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
            <bk-button theme="primary" :disabled="!selected.length" @click="handleNextStep">{{confirmText || $t('下一步')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { foreignkey } from '@/filters/formatter.js'
    import HostSelectorTopology from './host-selector-topology.vue'
    import HostSelectorCustom from './host-selector-custom.vue'
    export default {
        name: 'cmdb-host-selector',
        components: {
            HostSelectorTopology,
            HostSelectorCustom
        },
        filters: {
            foreignkey
        },
        props: {
            exist: {
                type: Array,
                default: () => ([])
            },
            exclude: {
                type: Array,
                default: () => ([])
            },
            confirmText: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                components: {
                    topology: HostSelectorTopology,
                    custom: HostSelectorCustom
                },
                activeComponent: null,
                type: 'topology',
                filter: '',
                repeatSelected: [],
                uniqueSelected: [],
                request: {
                    internal: Symbol('topology'),
                    instance: Symbol('instance'),
                    host: Symbol('host')
                }
            }
        },
        computed: {
            selected () {
                return [...this.repeatSelected, ...this.uniqueSelected]
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
        created () {
            this.setSelected(this.exist)
        },
        methods: {
            handleRemove (hosts) {
                const removeData = Array.isArray(hosts) ? hosts : [hosts]
                const ids = [...new Set(removeData.map(data => data.host.bk_host_id))]
                const selected = this.selected.filter(target => !ids.includes(target.host.bk_host_id))
                this.setSelected(selected)
            },
            handleSelect (hosts) {
                const selectData = Array.isArray(hosts) ? hosts : [hosts]
                const ids = [...new Set(selectData.map(data => data.host.bk_host_id))]
                const uniqueData = ids.map(id => selectData.find(data => data.host.bk_host_id === id))
                const newSelectData = []
                uniqueData.forEach(data => {
                    const isExist = this.selected.some(target => target.host.bk_host_id === data.host.bk_host_id)
                    if (!isExist) {
                        newSelectData.push(data)
                    }
                })
                if (newSelectData.length) {
                    this.setSelected([...this.selected, ...newSelectData])
                }
            },
            setSelected (selected) {
                const ipMap = {}
                const repeat = []
                const unique = []
                selected.forEach(data => {
                    const ip = data.host.bk_host_innerip
                    if (ipMap.hasOwnProperty(ip)) {
                        ipMap[ip].push(data)
                    } else {
                        ipMap[ip] = [data]
                    }
                })
                Object.values(ipMap).forEach(value => {
                    if (value.length > 1) {
                        repeat.push(...value)
                    } else {
                        unique.push(...value)
                    }
                })
                this.repeatSelected = repeat
                this.uniqueSelected = unique
            },
            handleCancel () {
                this.$emit('cancel')
            },
            handleNextStep () {
                this.$emit('confirm', this.selected)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .layout {
        position: relative;
        height: 460px;
        min-height: 300px;
        padding: 0 0 50px;
        .layout-footer {
            position: sticky;
            bottom: 0;
            left: 0;
            height: 50px;
            border-top: 1px solid $borderColor;
            font-size: 0;
            text-align: right;
            background-color: #FAFBFD;
            padding: 8px 20px 9px;
            z-index: 100;
        }
    }
    .wrapper {
        height: 100%;
        .wrapper-left {
            width: 340px;
            height: 100%;
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
                height: calc(100% - 105px);
            }
        }
        .wrapper-right {
            width: 510px;
            height: 100%;
            .selected-count {
                height: 42px;
                line-height: 42px;
                padding: 0 12px;
                background-color: #FAFBFD;
                font-size: 12px;
                font-weight: bold;
                color: $textColor;
            }
            .repeat-row {
                padding: 0 2px;
                line-height: 18px;
                background-color: #FE9C00;
                border-radius: 2px;
                color: #FFF;
                font-size: 12px;
            }
        }
    }
</style>
