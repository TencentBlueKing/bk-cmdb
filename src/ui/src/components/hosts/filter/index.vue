<template>
    <bk-popover
        ref="filterPopper"
        placement="bottom"
        theme="light"
        trigger="manual"
        :width="350"
        :on-show="handleShow"
        :tippy-options="{
            zIndex: 1001,
            interactive: true,
            hideOnClick: false
        }">
        <bk-button class="options-button"
            theme="default"
            v-bk-tooltips.top="$t('高级筛选')"
            icon="icon-cc-funnel"
            @click="handleToggleFilter">
        </bk-button>
        <section class="filter-content" slot="content"
            :style="{
                height: $APP.height - 150 + 'px'
            }">
            <h2 class="filter-title">
                {{$t('条件筛选')}}
                <bk-button class="close-trigger" text icon="close"></bk-button>
            </h2>
            <div class="filter-scroller" ref="scroller">
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
                    <bk-button class="filter-add-button" type="primary" icon="plus" text @click="handleAddFilter">{{$t('更多条件')}}</bk-button>
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
            </div>
            <div class="filter-options clearfix"
                :class="{
                    'is-sticky': isScrolling
                }">
                <div class="fl">
                    <bk-button theme="primary" @click="handleSearch">{{$t('查询')}}</bk-button>
                    <bk-button theme="default" @click="handleCollect">{{$t('收藏条件')}}</bk-button>
                </div>
                <div class="fr">
                    <bk-button theme="default" @click="handleReset">{{$t('清空')}}</bk-button>
                </div>
            </div>
        </section>
        <property-selector :properties="properties" ref="propertySelector"></property-selector>
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
            const defaultIpConfig = {
                text: '',
                inner: true,
                outer: true,
                exact: false
            }
            return {
                ip: {
                    ...defaultIpConfig
                },
                filterCondition: [],
                defaultIpConfig,
                isScrolling: false
            }
        },
        computed: {
            ...mapState('hosts', ['filterList'])
        },
        watch: {
            filterList () {
                this.setFilterCondition()
            },
            filterCondition () {
                this.checkIsScrolling()
            }
        },
        beforeRouteLeave () {
            this.$store.commit('hosts/clearFilter')
        },
        methods: {
            handleToggleFilter () {
                const [instance] = this.$refs.filterPopper.instance.instances
                const state = instance.state
                if (state.isVisible) {
                    instance.hide()
                } else {
                    instance.show()
                }
            },
            handleAddFilter () {
                this.$refs.propertySelector.isShow = true
            },
            handleSearch () {
                const params = this.getParams()
                this.$store.commit('hosts/setFilterParams', params)
                this.handleToggleFilter()
            },
            handleCollect () {},
            handleReset () {
                this.ip = { ...this.defaultIpConfig }
                this.filterCondition.forEach(filterItem => {
                    filterItem.value = ''
                })
                const params = this.getParams()
                this.$store.commit('hosts/setFilterParams', params)
            },
            getParams () {
                const params = {
                    ip: {
                        data: [],
                        exact: this.ip.exact ? 1 : 0,
                        flag: ['bk_host_innerip', 'bk_host_outerip'].filter((flag, index) => {
                            return index === 0 ? this.ip.inner : this.ip.outer
                        }).join('|')
                    },
                    host: [],
                    module: [],
                    set: []
                }
                this.ip.text.split(/\n|;|；|,|，/).forEach(text => {
                    const trimStr = text.trim()
                    if (trimStr.length) {
                        params.ip.data.push(trimStr)
                    }
                })
                this.filterCondition.forEach(filterItem => {
                    const filterValue = filterItem.value
                    if (filterValue !== null && filterValue !== undefined && String(filterValue).length) {
                        const modelId = filterItem.bk_obj_id
                        params[modelId].push({
                            field: filterItem.bk_property_id,
                            operator: filterItem.operator,
                            value: filterValue
                        })
                    }
                })
                return params
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
            checkIsScrolling () {
                const scroller = this.$refs.scroller
                this.isScrolling = scroller.scrollHeight > scroller.offsetHeight
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
    .filter-title {
        position: relative;
        padding: 10px 20px;
        font-size:14px;
        color: #63656E;
        .close-trigger {
            position: absolute;
            right: 0px;
            top: 0px;
        }
    }
    .filter-scroller {
        position: relative;
        max-height: calc(100% - 90px);
        padding: 10px 20px;
        overflow: auto;
        @include scrollbar-y;
        &.is-scrolling {
            padding-bottom: 62px;
            .filter-options {
                position: absolute;
                bottom: 0;
                left: 0;
                width: 100%;
                background-color: #FAFBFD;
                border-top: 1px solid #DCDEE5;
            }
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
    .filter-options {
        margin: 10px 0 0 0;
        padding: 10px 0;
    }
</style>
