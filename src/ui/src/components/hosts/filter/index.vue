<template>
    <bk-popover
        ref="filterPopper"
        placement="bottom"
        theme="light"
        trigger="click"
        :width="350"
        :on-show="handleShow"
        :tippy-options="{
            zIndex: 1001,
            interactive: true
        }">
        <bk-button class="options-button"
            theme="default"
            v-bk-tooltips.top="$t('高级筛选')"
            icon="icon-cc-funnel">
        </bk-button>
        <section class="filter-content" slot="content"
            :style="{
                height: $APP.height - 150 + 'px'
            }">
            <h2 class="filter-title">
                {{$t('条件筛选')}}
                <bk-button class="close-trigger" text icon="close"></bk-button>
            </h2>
            <div class="filter-group">
                <label class="filter-label">IP</label>
                <bk-input type="textarea" v-model="ip.text" :rows="4"></bk-input>
            </div>
            <div class="filter-group checkbox-group">
                <bk-checkbox class="filter-checkbox"
                    v-model="ip.inner"
                    :disabled="!ip.outer">
                    {{$t('内网')}}
                </bk-checkbox>
                <bk-checkbox class="filter-checkbox"
                    v-model="ip.outer"
                    :disabled="!ip.inner">
                    {{$t('外网')}}
                </bk-checkbox>
                <bk-checkbox class="filter-checkbox" v-model="ip.exact">{{$t('精确')}}</bk-checkbox>
            </div>
            <div class="filter-add">
                <bk-button class="filter-add-button" type="primary" icon="plus" text>{{$t('更多条件')}}</bk-button>
            </div>
            <div class="filter-group"
                v-for="(filterItem, index) in filterCondition"
                :key="index">
                <label class="filter-label">{{getFilterLabel(filterItem)}}</label>
                <div class="filter-condition">
                    <filter-operator class="filter-operator"
                        :type="getOperatorType(filterItem)"
                        v-model="filterItem.operator">
                    </filter-operator>
                    <cmdb-form-enum class="filter-value"
                        v-if="filterItem.bk_property_type === 'enum'"
                        :options="filterItem.option || []"
                        v-model="filterItem.value">
                    </cmdb-form-enum>
                    <cmdb-form-bool-input class="filter-value"
                        v-else-if="filterItem.bk_property_type === 'bool'"
                        v-model="filterItem.value">
                    </cmdb-form-bool-input>
                    <component class="filter-value"
                        v-else
                        :is="`cmdb-form-${filterItem.bk_property_type}`"
                        v-model="filterItem.value">
                    </component>
                </div>
            </div>
        </section>
        <property-selector :properties="properties"></property-selector>
    </bk-popover>
</template>

<script>
    import filterOperator from './_filter-field-operator.vue'
    import propertySelector from './filter-property-selector.vue'
    import { mapState } from 'vuex'
    export default {
        components: {
            filterOperator,
            propertySelector
        },
        props: {
            properties: {
                type: Object,
                default () {
                    return {}
                }
            }
        },
        data () {
            return {
                ip: {
                    text: '',
                    inner: true,
                    outer: true,
                    exact: false
                },
                filterCondition: []
            }
        },
        computed: {
            ...mapState('hosts', ['filterList'])
        },
        watch: {
            filterList () {
                this.setFilterCondition()
            }
        },
        methods: {
            handleToggleFilter () {
                console.log(this.$refs.filterPopper)
            },
            setFilterCondition () {
                try {
                    const condition = []
                    this.filterList.forEach(filter => {
                        const modelId = filter.bk_obj_id
                        const propertyId = filter.bk_property_id
                        const property = (this.properties[modelId] || []).find(property => property.bk_property_id === propertyId)
                        if (property) {
                            condition.push({
                                bk_obj_id: modelId,
                                bk_property_id: propertyId,
                                bk_property_type: property.bk_property_type,
                                option: property.option,
                                operator: '',
                                value: ''
                            })
                        }
                    })
                    this.filterCondition = condition
                } catch (e) {
                    console.error(e)
                }
            },
            handleShow (popper) {
                popper.popperChildren.tooltip.style.padding = 0
            },
            getFilterLabel (filterItem) {
                const model = this.$store.getters['objectModelClassify/getModelById'](filterItem.bk_obj_id) || {}
                const property = (this.properties[filterItem.bk_obj_id] || []).find(property => property.bk_property_id === filterItem.bk_property_id) || {}
                return `${model.bk_obj_name} - ${property.bk_property_name}`
            },
            getOperatorType (filterItem) {
                const propertyType = filterItem.bk_property_type
                const propertyId = filterItem.bk_property_id
                if (['bk_set_name', 'bk_module_name'].includes(propertyId)) {
                    return 'name'
                } else if (['singlechar', 'longchar'].includes(propertyType)) {
                    return 'char'
                }
                return 'common'
            }
        }
    }
</script>

<style lang="scss" scoped="true">
    .filter-content {
        position: relative;
        padding: 10px 20px;
    }
    .filter-title {
        position: relative;
        font-size:14px;
        color: #63656E;
        .close-trigger {
            position: absolute;
            right: -15px;
            top: -4px;
        }
    }
    .filter-group {
        padding: 15px 0 0 0;
        &.checkbox-group {
            padding: 10px 0 0 0;
            .filter-checkbox {
                margin: 0 15px 0 0;
            }
        }
        .filter-label {
            display: block;
            line-height: 30px;
            color: #63656E;
        }
    }
    .filter-add {
        margin: 14px 0 0 0;
        .filter-add-button {
            /deep/ {
                span {
                    display: inline-block;
                    vertical-align: middle;
                }
            }
        }
    }
    .filter-condition {
        display: flex;
        .filter-operator {
            flex: 75px 0 0;
            margin-right: 8px;
        }
        .filter-value {
            flex: 1;
        }
    }
</style>
