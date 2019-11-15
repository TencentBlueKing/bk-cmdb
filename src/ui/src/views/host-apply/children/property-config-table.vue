<template>
    <bk-table
        header-cell-class-name="header-cell"
        :data="modulePropertyList"
        @selection-change="handleSelectionChange"
    >
        <bk-table-column type="selection" width="60" align="center" v-if="deletable"></bk-table-column>
        <bk-table-column
            :width="readonly ? 300 : 200"
            :label="$t('字段名称')"
            prop="bk_property_name"
        >
        </bk-table-column>
        <bk-table-column
            v-if="multiple"
            :label="$t('已配置的模块')"
            :render-header="(h, data) => renderTableHeader(h, data, '主机属于多个模块，但主机的属性值仅能有唯一值')"
            class-name="table-cell-module-path"
        >
            <template slot-scope="{ row }">
                <template v-if="row.module_list">
                    <div v-for="(item, index) in row.module_list" :key="index">
                        {{item.path}}
                    </div>
                </template>
                <span v-else>--</span>
            </template>
        </bk-table-column>
        <bk-table-column
            v-if="multiple || readonly"
            :width="readonly ? null : 200"
            :label="$t('当前值')"
        >
            <template slot-scope="{ row }">
                <template v-if="multiple">
                    <div v-if="row.module_list">
                        <div v-for="item in row.module_list" :key="item.bk_inst_id">
                            {{getDisplayValue(item.value, row)}}
                        </div>
                    </div>
                    <span v-else>--</span>
                </template>
                <template v-else>
                    {{getDisplayValue(row.__extra__.value, row)}}
                </template>
            </template>
        </bk-table-column>
        <bk-table-column v-if="!readonly || !deletable" :label="$t(multiple ? '修改后' : '值')" class-name="table-cell-form-element">
            <template slot-scope="{ row }">
                <div class="form-element-content">
                    <property-form-element :property="row" @value-change="handlePropertyValueChange"></property-form-element>
                </div>
            </template>
        </bk-table-column>
        <bk-table-column
            v-if="!readonly || !deletable"
            width="200"
            :label="$t('操作')"
            :render-header="multiple ? (h, data) => renderTableHeader(h, data, '删除操作不影响原有配置') : null"
        >
            <template slot-scope="{ row }">
                <bk-button theme="primary" text @click="handlePropertyRowDel(row.bk_property_id)">删除</bk-button>
            </template>
        </bk-table-column>
    </bk-table>
</template>
<script>
    import { mapGetters } from 'vuex'
    import propertyFormElement from './property-form-element'
    export default {
        components: {
            propertyFormElement
        },
        props: {
            checkedPropertyIdList: {
                type: Array,
                default: () => ([])
            },
            initData: {
                type: Array,
                default: () => ([])
            },
            multiple: {
                type: Boolean,
                default: false
            },
            readonly: {
                type: Boolean,
                default: false
            },
            deletable: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                modulePropertyList: []
            }
        },
        computed: {
            ...mapGetters('hosts', ['configPropertyList'])
        },
        watch: {
            checkedPropertyIdList () {
                this.setModulePropertyList()
            }
        },
        created () {
            this.setModulePropertyList()
        },
        methods: {
            setModulePropertyList () {
                // 当前模块属性列表中不存在，则添加
                this.checkedPropertyIdList.forEach(propertyId => {
                    if (this.modulePropertyList.findIndex(property => propertyId === property.bk_property_id) === -1) {
                        // 使用原始主机属性对象
                        const findProperty = this.configPropertyList.find(item => propertyId === item.bk_property_id)
                        if (findProperty) {
                            // 初始化值
                            if (this.multiple) {
                                findProperty.module_list = (this.initData.find(item => item.bk_property_id === findProperty.bk_property_id) || {}).module_list
                            } else {
                                findProperty.__extra__.value = (this.initData.find(item => item.bk_property_id === findProperty.bk_property_id) || {}).value
                            }
                            this.modulePropertyList.push(this.$tools.clone(findProperty))
                        }
                    }
                })

                // 删除或取消选择的，则去除
                this.modulePropertyList = this.modulePropertyList.filter(property => this.checkedPropertyIdList.includes(property.bk_property_id))
            },
            getDisplayValue (value, property) {
                let displayValue = value
                switch (property.bk_property_type) {
                    case 'enum':
                        displayValue = (property.option.find(item => item.id === value) || {}).name
                        break
                    case 'bool':
                        displayValue = ['否', '是'][+value]
                        break
                }

                return displayValue
            },
            renderTableHeader (h, data, tips) {
                const directive = {
                    content: tips,
                    placement: 'bottom-start'
                }
                return <span>{ data.column.label } <i class="bk-cc-icon icon-cc-tips" v-bk-tooltips={ directive }></i></span>
            },
            handleSelectionChange (value) {
                this.$emit('selection-change', value)
            },
            handlePropertyRowDel (id) {
                const checkedIndex = this.checkedPropertyIdList.findIndex(propertyId => propertyId === id)
                this.checkedPropertyIdList.splice(checkedIndex, 1)
            },
            handlePropertyValueChange (value) {
                this.$emit('property-value-change', value)
            }
        }
    }
</script>

<style lang="scss">
    .table-cell-module-path:not(.header-cell) {
        .cell {
            padding-top: 8px;
            padding-bottom: 8px;
        }
    }
    .table-cell-form-element {
        .cell {
            overflow: unset;
            display: block;
        }
        .search-input-wrapper {
            position: relative;
        }
        .form-objuser .suggestion-list {
            z-index: 1000;
        }
    }
</style>
