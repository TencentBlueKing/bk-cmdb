<template>
    <div class="conflict-layout">
        <cmdb-tips class="resolve-tips"
            :tips-style="{
                height: '32px',
                'line-height': '32px',
                'font-size': '12px'
            }"
            :icon-style="{ 'font-size': '16px' }">
            {{$t('解决冲突：同一主机的属性在不同的模块下需保持一致')}}
        </cmdb-tips>
        <div class="conflict-table" ref="conflictTable">
            <bk-table :data="list">
                <bk-table-column :label="$t('字段名称')" width="160">
                    <template slot-scope="{ row }">
                        <span class="conflict-name">{{row.name}}</span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('所属模块')" :render-header="hanldeColumnRender">
                    <template slot-scope="{ row }">
                        <div class="conflict-modules">
                            <div v-for="(module, index) in row.conflictList"
                                :class="['module-item', { 'selected': module.selected }]"
                                :key="index">
                                <span>{{module.path || '--'}}</span>
                                <span>{{module.value | formatter(row.property.bk_property_type, row.property.option)}}</span>
                                <i class="check-model-value" @click="handleSelectDefaultValue(row, index, module.value)">选定</i>
                            </div>
                        </div>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('修改后')" width="230" class-name="conflict-custom-column">
                    <template slot-scope="{ row: { property = {} } = {} }">
                        <component class="property-component"
                            :is="`cmdb-form-${property.bk_property_type}`"
                            :class="[property.bk_property_type, { error: errors.has(property.bk_property_id) }]"
                            :options="property.option || []"
                            :data-vv-name="property.bk_property_id"
                            :data-vv-as="property.bk_property_name"
                            :placeholder="$t('请输入xx', { name: property.bk_property_name })"
                            v-validate="$tools.getValidateRules(property)"
                            v-model.trim="result[property.bk_property_id]">
                        </component>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div :class="['footer-btns', { 'sticky': scrollbar }]">
            <bk-button theme="primary" class="mr10">{{$t('确定')}}</bk-button>
            <bk-button theme="default">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import RESIZE_EVENTS from '@/utils/resize-events'
    export default {
        data () {
            return {
                list: [],
                properties: [],
                result: {},
                scrollbar: false
            }
        },
        async created () {
            const [list] = await Promise.all([
                this.getData(),
                this.getProperties()
            ])
            this.getPropertyInfo(list)
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$refs.conflictTable, this.checkScrollbar)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$refs.conflictTable, this.checkScrollbar)
        },
        methods: {
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            getPropertyInfo (list) {
                list.forEach((item, index) => {
                    const property = this.properties.find(property => property.bk_property_id === item.bk_property_id)
                    item.property = property
                    this.result[item.bk_property_id] = item[item.bk_property_id]
                })
                this.list = list
            },
            getData () {
                return [{
                    name: '主机名称',
                    bk_host_name: 'Selina',
                    bk_property_id: 'bk_host_name',
                    conflictList: [{
                        path: '当前主机值',
                        value: 'Quelly',
                        selected: true
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: 'Selina,Quelly',
                        selected: false
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块dddddddd',
                        value: 'Selina,Quelly',
                        selected: false
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块ccccc',
                        value: 'Selina,Quelly',
                        selected: false
                    }]
                }, {
                    name: '所属运营商',
                    bk_isp_name: '0',
                    bk_property_id: 'bk_isp_name',
                    conflictList: [{
                        path: '当前主机值',
                        value: '0',
                        selected: true
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: '1',
                        selected: false
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块dddddddd',
                        value: '2',
                        selected: false
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块ccccc',
                        value: '3',
                        selected: false
                    }]
                }, {
                    name: 'CPU逻辑核心数',
                    bk_cpu: 1,
                    bk_property_id: 'bk_cpu',
                    conflictList: [{
                        path: '当前主机值',
                        value: 1,
                        selected: true
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: 2,
                        selected: false
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块dddddddd',
                        value: 3,
                        selected: false
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块ccccc',
                        value: 4,
                        selected: false
                    }]
                }, {
                    name: '备份维护人',
                    operator: 'admin',
                    bk_property_id: 'operator',
                    conflictList: [{
                        path: '当前主机值',
                        value: 'admin',
                        selected: true
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: 'peibiaoyang',
                        selected: false
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块dddddddd',
                        value: 'yunyaoyang',
                        selected: false
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块ccccc',
                        value: 'v_dgzheng',
                        selected: false
                    }]
                }]
            },
            async getProperties () {
                try {
                    const properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: this.$injectMetadata({
                            bk_obj_id: 'host'
                        })
                    })
                    this.properties = properties
                } catch (e) {
                    console.error(e)
                }
            },
            hanldeColumnRender (h, { column, $index }) {
                const style = {
                    width: '50%',
                    display: 'inline-block'
                }
                return h(
                    'div',
                    { class: 'conflict-custom-label' },
                    [
                        h('span', { style }, this.$t('所属模块')),
                        h('span', { style }, this.$t('冲突值'))
                    ]
                )
            },
            handleSelectDefaultValue (row, index, value) {
                row.conflictList.forEach((module, _index) => {
                    module.selected = _index === index
                })
                this.result[row.bk_property_id] = value
            }
        }
    }
</script>

<style lang="scss" scoped>
    .conflict-layout {
        height: 100%;
        @include scrollbar-y;
    }
    .resolve-tips {
        margin: 12px 20px;
    }
    .conflict-table {
        padding: 0 20px;
        .conflict-name {
            font-weight: bold;
        }
        .conflict-modules {
            height: 100%;
            .module-item {
                display: flex;
                line-height: 22px;
                span {
                    max-width: 50%;
                    padding-right: 10px;
                    @include ellipsis;
                    &:first-child {
                        width: calc(50% - 30px);
                    }
                }
                .check-model-value {
                    display: none;
                    font-style: normal;
                    cursor: pointer;
                    color: #3A84FF;
                }
                &.selected {
                    color: #3A84FF;
                    background-color: #E1ECFF;
                }
                &:not(.selected):hover {
                    .check-model-value {
                        display: inline-block;
                    }
                }
            }
        }
    }
    .footer-btns {
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 20px 20px 0;
        background-color: #FFFFFF;
        z-index: 10;
        &.sticky {
            height: 54px;
            line-height: 54px;
            padding: 0 20px;
            margin-top: 20px;
            border-top: 1px solid #dcdee5;
            font-size: 0;
        }
    }
</style>

<style lang="scss">
    .conflict-custom-column {
        .cell {
            overflow: unset;
        }
    }
</style>
