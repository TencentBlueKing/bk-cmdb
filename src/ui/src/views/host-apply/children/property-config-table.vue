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
                <template v-if="row.__extra__.moduleList.length">
                    <div v-for="(item, index) in row.__extra__.moduleList" :key="index">
                        {{$parent.getModulePath(item.bk_module_id)}}
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
                    <template v-if="row.__extra__.moduleList.length">
                        <div v-for="(item, index) in row.__extra__.moduleList" :key="index">
                            {{item.bk_property_value | formatter(row) | unit(row.unit)}}
                        </div>
                    </template>
                    <span v-else>--</span>
                </template>
                <template v-else>
                    {{row.__extra__.value | formatter(row) | unit(row.unit)}}
                </template>
            </template>
        </bk-table-column>
        <bk-table-column v-if="!readonly" :label="$t(multiple ? '修改后' : '值')" class-name="table-cell-form-element">
            <template slot-scope="{ row }">
                <div class="form-element-content">
                    <property-form-element :property="row" @value-change="handlePropertyValueChange"></property-form-element>
                </div>
            </template>
        </bk-table-column>
        <bk-table-column
            v-if="!readonly"
            width="200"
            :label="$t('操作')"
            :render-header="multiple ? (h, data) => renderTableHeader(h, data, '删除操作不影响原有配置') : null"
        >
            <template slot-scope="{ row }">
                <bk-button theme="primary" text @click="handlePropertyRowDel(row)">删除</bk-button>
            </template>
        </bk-table-column>
    </bk-table>
</template>
<script>
    import { mapGetters, mapState } from 'vuex'
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
            ruleList: {
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
                modulePropertyList: [],
                removeRuleIds: [],
                ignoreRuleIds: []
            }
        },
        computed: {
            ...mapGetters('hosts', ['configPropertyList']),
            ...mapState('hostApply', ['ruleDraft'])
        },
        watch: {
            checkedPropertyIdList () {
                this.setModulePropertyList()
            },
            configPropertyList () {
                this.setModulePropertyList()
            },
            modulePropertyList () {
                this.setConfigData()
            }
        },
        created () {
            this.setModulePropertyList()
            if (Object.keys(this.ruleDraft).length) {
                this.modulePropertyList = this.$tools.clone(this.ruleDraft.rules)
            }
        },
        methods: {
            setModulePropertyList () {
                // 当前模块属性列表中不存在，则添加
                this.checkedPropertyIdList.forEach(id => {
                    // 原始主机属性对象
                    const moduleIndex = this.modulePropertyList.findIndex(property => id === property.id)
                    if (moduleIndex === -1) {
                        const findProperty = this.configPropertyList.find(item => id === item.id)
                        if (findProperty) {
                            const property = this.$tools.clone(findProperty)
                            // 初始化值
                            if (this.multiple) {
                                property.__extra__.moduleList = this.ruleList.filter(item => item.bk_attribute_id === property.id)
                            } else {
                                const rule = this.ruleList.find(item => item.bk_attribute_id === property.id) || {}
                                property.__extra__.ruleId = rule.id
                                property.__extra__.value = rule.bk_property_value
                            }
                            this.modulePropertyList.push(property)
                        }
                    }
                })

                // 删除或取消选择的，则去除
                this.modulePropertyList = this.modulePropertyList.filter(property => this.checkedPropertyIdList.includes(property.id))
            },
            setConfigData () {
                this.removeRuleIds = []
                this.ignoreRuleIds = []
                // 找出不存在于初始数据中的规则
                this.ruleList.forEach(rule => {
                    const findIndex = this.modulePropertyList.findIndex(property => property.id === rule.bk_attribute_id)
                    if (findIndex === -1) {
                        // 批量模式标记为忽略，单个标记为删除
                        if (this.multiple) {
                            this.ignoreRuleIds.push(rule.id)
                        } else {
                            this.removeRuleIds.push(rule.id)
                        }
                    }
                })
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
            handlePropertyRowDel (property) {
                const checkedIndex = this.checkedPropertyIdList.findIndex(id => id === property.id)
                this.checkedPropertyIdList.splice(checkedIndex, 1)
            },
            handlePropertyValueChange (value) {
                this.$emit('property-value-change', value)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-element-content {
        width: 80%;
    }
</style>
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
