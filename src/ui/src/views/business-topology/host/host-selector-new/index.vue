<template>
    <div class="host-selector-layout">
        <div class="layout-header">
            <h2 class="title">{{title || $t('选择主机')}}</h2>
        </div>
        <div class="layout-content">
            <div class="topo-table">
                <bk-tab :active.sync="tab.active" type="border-card">
                    <bk-tab-panel
                        v-for="(panel, index) in tab.panels"
                        v-bind="panel"
                        :key="index">
                        <div class="tab-content">
                            <keep-alive>
                                <component
                                    :is="components[panel.name]"
                                    class="selector-component"
                                    :selected="selected"
                                    @select-change="handleSelectChange">
                                </component>
                            </keep-alive>
                        </div>
                    </bk-tab-panel>
                </bk-tab>
            </div>
            <div class="result-preview">
                <div class="preview-title">{{$t('结果预览')}}</div>
                <div class="preview-content" v-show="selected.length > 0">
                    <bk-collapse v-model="collapse.active" @item-click="handleCollapse">
                        <bk-collapse-item name="host" hide-arrow>
                            <div class="collapse-title">
                                <div class="text">
                                    <i :class="['bk-icon icon-angle-right', { expand: collapse.expanded.host }]"></i>
                                    <i18n path="已选择N台主机">
                                        <span class="count" place="count">{{selected.length}}</span>
                                    </i18n>
                                </div>
                                <div class="more" @click.stop="handleClickMore">
                                    <i class="bk-icon icon-more"></i>
                                </div>
                            </div>
                            <ul slot="content" class="host-list">
                                <li class="host-item" v-for="(row, index) in selected" :key="index">
                                    <div class="ip">
                                        {{row.host.bk_host_innerip}}
                                        <span class="repeat-tag" v-if="repeatSelected.includes(row)" v-bk-tooltips="{ content: `${$t('云区域')}：${foreignkey(row.host.bk_cloud_id)}` }">{{$t('IP重复')}}</span>
                                    </div>
                                    <i class="bk-icon icon-close-line" @click="handleRemove(row)"></i>
                                </li>
                            </ul>
                        </bk-collapse-item>
                    </bk-collapse>
                </div>
            </div>
        </div>
        <div class="layout-footer">
            <bk-button class="mr10" theme="primary" :disabled="!selected.length" @click="handleNextStep">{{confirmText || $t('下一步')}}</bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
        <ul class="more-menu" ref="moreMenu" v-show="more.show">
            <li class="menu-item" @click="handleCopyIp">{{$t('复制IP')}}</li>
            <li class="menu-item" @click="handleRemoveAll">{{$t('移除所有')}}</li>
        </ul>
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
            confirmText: {
                type: String,
                default: ''
            },
            title: String
        },
        data () {
            return {
                tab: {
                    panels: [
                        { name: 'topology', label: this.$t('静态拓扑') },
                        { name: 'custom', label: this.$t('自定义输入') }
                    ],
                    active: 'topology'
                },
                collapse: {
                    expanded: { host: true },
                    active: 'host'
                },
                components: {
                    topology: HostSelectorTopology,
                    custom: HostSelectorCustom
                },
                more: {
                    instance: null,
                    show: false
                },
                repeatSelected: [],
                uniqueSelected: []
            }
        },
        computed: {
            selected () {
                return [...this.repeatSelected, ...this.uniqueSelected]
            }
        },
        created () {
            this.setSelected(this.exist)
        },
        methods: {
            handleSelectChange ({ removed, selected }) {
                this.handleRemove(removed)
                this.handleSelect(selected)
            },
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
            handleClickMore (event) {
                this.more.instance && this.more.instance.destroy()
                this.more.instance = this.$bkPopover(event.target, {
                    content: this.$refs.moreMenu,
                    allowHTML: true,
                    delay: 300,
                    trigger: 'manual',
                    boundary: 'window',
                    placement: 'bottom-end',
                    theme: 'light host-selector-popover',
                    distance: 6,
                    interactive: true
                })
                this.more.show = true
                this.$nextTick(() => {
                    this.more.instance.show()
                })
            },
            handleCopyIp () {
                const ipList = this.selected.map(item => item.host.bk_host_innerip)
                this.$copyText(ipList.join('\n')).then(() => {
                    this.$success(this.$t('复制成功'))
                }, () => {
                    this.$error(this.$t('复制失败'))
                }).finally(() => {
                    this.more.instance.hide()
                })
            },
            handleRemoveAll () {
                this.repeatSelected = []
                this.uniqueSelected = []
                this.more.instance.hide()
            },
            handleCancel () {
                this.$emit('cancel')
            },
            handleNextStep () {
                this.$emit('confirm', this.selected)
            },
            handleCollapse (names) {
                Object.keys(this.collapse.expanded).forEach(key => (this.collapse.expanded[key] = false))
                names.forEach(name => {
                    this.$set(this.collapse.expanded, name, true)
                })
            },
            foreignkey (value) {
                return foreignkey(value)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-selector-layout {
        display: flex;
        flex-direction: column;
        height: 650px;
        padding: 20px 20px 0 20px;

        .layout-header {
            margin-bottom: 20px;
            .title {
                font-size: 20px;
                font-weight: normal;
                color: $textColor;
            }
        }

        .layout-content {
            display: flex;
            flex: auto;
            overflow: hidden;

            .tab-content {
                padding: 0 24px 0 24px;
                height: 100%;
            }
            .topo-table {
                flex: auto;
            }
            .result-preview {
                flex: none;
                width: 280px;
                border: 1px solid #dcdee5;
                border-left: none;
                background: #f5f6fa;

                .preview-title {
                    color: #313238;
                    font-size: 14px;
                    line-height: 22px;
                    padding: 10px 24px;
                }
                .preview-content {
                    height: calc(100% - 42px);
                    @include scrollbar;

                    /deep/ .bk-collapse-item .bk-collapse-item-header {
                        padding: 0;
                        height: 24px;
                        line-height: 24px;
                        &:hover {
                            color: #63656e;
                        }
                    }
                }
            }
        }

        .layout-footer {
            font-size: 0;
            text-align: right;
            padding: 12px 0;
        }

        .result-preview {
            .collapse-title {
                display: flex;
                align-items: center;
                justify-content: space-between;
                padding: 0 24px 0 18px;

                .icon-angle-right {
                    font-size: 24px;
                    transition: transform .2s ease-in-out;

                    &.expand {
                        transform: rotate(90deg);
                    }
                }
                .count {
                    font-weight: 700;
                    color: #3a84ff;
                    padding: 0 .1em;
                }
                .text {
                    display: flex;
                    align-items: center;
                }
                .more {
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    width: 24px;
                    height: 24px;
                    border-radius: 2px;
                    .icon-more {
                        font-size: 18px;
                        outline: 0;
                    }
                    &:hover {
                        background: #e1ecff;
                        color: #3a84ff;
                    }
                }
            }
            .host-list {
                padding: 0 14px;
                margin: 6px 0;
            }
            .host-item {
                display: flex;
                justify-content: space-between;
                align-items: center;
                height: 32px;
                line-height: 32px;
                background: #fff;
                padding: 0 12px;
                border-radius: 2px;
                box-shadow: 0 1px 2px 0 rgba(0,0,0,.06);
                margin-bottom: 2px;
                font-size: 12px;

                .ip {
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    word-break: break-all;

                    .repeat-tag {
                        height: 16px;
                        line-height: 16px;
                        background: #FFE8C3;
                        color: #FE9C00;
                        font-size: 12px;
                        padding: 0 4px;
                        border-radius: 2px;
                    }
                }
                .icon-close-line {
                    display: none;
                    cursor: pointer;
                    color: #3a84ff;
                    font-weight: 700;
                }

                &:hover {
                    .icon-close-line {
                        display: block;
                    }
                }
            }
        }

        /deep/ .bk-tab {
            height: 100%;
            .bk-tab-header {
                padding: 0;
                height: 43px;
                background-image: linear-gradient(transparent 41px,#dcdee5 0);
                .bk-tab-label-list {
                    height: 42px;
                    .bk-tab-label-item {
                        line-height: 42px;
                        min-width: auto;
                        &.active {
                            background-color: #fff;
                        }
                    }
                }
            }
            .bk-tab-header-setting {
                height: 42px;
                line-height: 42px;
            }
            .bk-tab-section {
                padding: 0;
                height: calc(100% - 43px);
                overflow: visible;
                .bk-tab-content {
                    height: 100%;
                }
            }
        }
    }
</style>
<style lang="scss">
    .host-selector-popover-theme {
        padding: 0;
        .more-menu {
            font-size: 12px;
            padding: 6px 0;
            min-width: 84px;
            background: #fff;
            .menu-item {
                height: 32px;
                line-height: 32px;
                padding: 0 10px;
                cursor: pointer;
                &:hover {
                    background: #f5f6fa;
                    color: #3a84ff;
                }
            }
        }
    }
</style>
