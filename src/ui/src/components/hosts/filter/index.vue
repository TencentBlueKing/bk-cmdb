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
            hideOnClick: false,
            onShown: checkIsScrolling
        }">
        <bk-button class="filter-trigger"
            theme="default"
            v-bk-tooltips.top="$t('高级筛选')"
            icon="icon-cc-funnel"
            :class="{
                'is-active': isFilterActive
            }"
            @click="handleToggleFilter">
        </bk-button>
        <section class="filter-content" slot="content"
            :style="{
                height: $APP.height - 150 + 'px'
            }">
            <h2 class="filter-title">
                {{$t('条件筛选')}}
                <bk-button class="close-trigger" text icon="close" @click="handleToggleFilter"></bk-button>
            </h2>
            <div class="filter-scroller" ref="scroller">
                <div class="filter-group" style="padding: 0;">
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
                <div class="filter-group" v-if="showScope">
                    <label class="filter-label">{{$t('搜索范围')}}</label>
                    <bk-checkbox class="filter-checkbox mr20"
                        v-model="scope.resource"
                        :disabled="!scope.business">
                        {{$t('未分配主机')}}
                    </bk-checkbox>
                    <bk-checkbox class="filter-checkbox"
                        v-model="scope.business"
                        :disabled="!scope.resource">
                        {{$t('已分配主机')}}
                    </bk-checkbox>
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
                <div class="filter-add">
                    <bk-button class="filter-add-button" type="primary" icon="plus" text @click="handleAddFilter">{{$t('更多条件')}}</bk-button>
                </div>
            </div>
            <div class="filter-options clearfix"
                :class="{
                    'is-sticky': isScrolling
                }">
                <div class="fl">
                    <bk-button theme="primary" @click="handleSearch">{{$t('查询')}}</bk-button>
                    <bk-button theme="default"
                        v-if="isCollection"
                        :loading="$loading('updateCollection')"
                        @click="handleUpdateCollection">
                        {{$t('更新条件')}}
                    </bk-button>
                    <bk-popover v-else
                        ref="collectionPopover"
                        placement="top-end"
                        theme="light"
                        trigger="manual"
                        :width="280"
                        :tippy-options="{
                            zIndex: 1002,
                            interactive: true,
                            hideOnClick: false
                        }">
                        <bk-button theme="default" @click="handleCreateCollection">{{$t('收藏条件')}}</bk-button>
                        <section class="collection" slot="content">
                            <label class="collection-title">{{$t('收藏条件')}}</label>
                            <bk-input class="collection-name"
                                :placeholder="$t('请填写名称')"
                                v-model="collectionName">
                            </bk-input>
                            <div class="collection-options">
                                <bk-button
                                    theme="primary"
                                    size="small"
                                    :disabled="!collectionName.length"
                                    :loading="$loading('createCollection')"
                                    @click="handleSaveCollection">
                                    {{$t('确定')}}
                                </bk-button>
                                <bk-button theme="default" size="small" @click="handleCancelCollection">{{$t('取消')}}</bk-button>
                            </div>
                        </section>
                    </bk-popover>
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
    import { mapState, mapGetters } from 'vuex'
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
            },
            showScope: Boolean
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
                scope: {
                    resource: true,
                    business: true
                },
                filterCondition: [],
                defaultIpConfig,
                isScrolling: false,
                collectionName: ''
            }
        },
        computed: {
            ...mapState('hosts', ['filterList', 'filterIP', 'collection']),
            ...mapGetters('hosts', ['isCollection']),
            isFilterActive () {
                const hasIP = !!this.ip.text.replace(/\n|;|；|,|，/g, '').length
                const hasField = this.filterCondition.some(condition => {
                    const filterValue = condition.value
                    return filterValue !== null
                        && filterValue !== undefined
                        && !!String(filterValue).length
                })
                const hasScope = !this.scope.resource || !this.scope.business
                return hasIP || hasField || hasScope
            }
        },
        watch: {
            filterList (newList, oldList) {
                this.setFilterCondition()
            },
            filterIP (value) {
                this.initCustomFilterIP()
            },
            filterCondition () {
                this.checkIsScrolling()
            }
        },
        created () {
            this.initCustomFilterIP()
            this.initCustomFilterList()
        },
        beforeDestroy () {
            this.$store.commit('hosts/clearFilter')
        },
        methods: {
            initCustomFilterIP () {
                if (this.filterIP) {
                    Object.assign(this.ip, this.filterIP)
                } else {
                    this.ip = { ...this.defaultIpConfig }
                }
            },
            initCustomFilterList () {
                const key = this.$route.meta.filterPropertyKey
                const customData = this.$store.getters['userCustom/getCustomData'](key, [])
                this.$store.commit('hosts/setFilterList', customData)
            },
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
            handleSearch (toggle = true) {
                const params = this.getParams()
                this.$store.commit('hosts/setFilterParams', params)
                if (toggle) {
                    this.handleToggleFilter()
                }
            },
            handleCreateCollection () {
                const [instance] = this.$refs.collectionPopover.instance.instances
                instance.show()
            },
            async handleUpdateCollection () {
                try {
                    const params = {
                        ...this.getCollectionParams(),
                        name: this.collection.name
                    }
                    await this.$store.dispatch('hostFavorites/updateFavorites', {
                        id: this.collection.id,
                        params: params,
                        config: {
                            requestId: 'updateCollection'
                        }
                    })
                    this.$store.commit('hosts/updateCollection', params)
                    this.$success(this.$t('更新收藏成功'))
                } catch (e) {
                    console.error(e)
                }
            },
            async handleSaveCollection () {
                try {
                    await this.$store.dispatch('hostFavorites/createFavorites', {
                        params: this.getCollectionParams(),
                        config: {
                            requestId: 'createCollection'
                        }
                    })
                    this.$success(this.$t('收藏成功'))
                    this.handleCancelCollection()
                } catch (e) {
                    console.error(e)
                }
            },
            getCollectionParams () {
                return {
                    name: this.collectionName,
                    info: JSON.stringify({
                        exact_search: this.ip.exact,
                        bk_host_innerip: this.ip.inner,
                        bk_host_outerip: this.ip.outer,
                        ip_list: this.getIPList()
                    }),
                    query_params: JSON.stringify(this.filterCondition.map(condition => {
                        return {
                            bk_obj_id: condition.bk_obj_id,
                            field: condition.bk_property_id,
                            operator: condition.operator,
                            value: condition.value
                        }
                    })),
                    is_default: 2
                }
            },
            handleCancelCollection () {
                const [instance] = this.$refs.collectionPopover.instance.instances
                instance.hide()
                this.collectionName = ''
            },
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
                        data: this.getIPList(),
                        exact: this.ip.exact ? 1 : 0,
                        flag: ['bk_host_innerip', 'bk_host_outerip'].filter((flag, index) => {
                            return index === 0 ? this.ip.inner : this.ip.outer
                        }).join('|')
                    },
                    host: [],
                    module: [],
                    set: [],
                    biz: []
                }
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
                if (this.scope.resource && !this.scope.business) {
                    params.biz.push({
                        field: 'default',
                        operator: '$eq',
                        value: 1
                    })
                } else if (!this.scope.resource && this.scope.business) {
                    params.biz.push({
                        field: 'default',
                        operator: '$eq',
                        value: 0
                    })
                }
                return params
            },
            getIPList () {
                const list = []
                this.ip.text.split(/\n|;|；|,|，/).forEach(text => {
                    const trimStr = text.trim()
                    if (trimStr.length) {
                        list.push(trimStr)
                    }
                })
                return list
            },
            setFilterCondition () {
                try {
                    const filterCondition = []
                    const oldCondition = this.filterCondition
                    this.filterList.forEach(filter => {
                        const modelId = filter.bk_obj_id
                        const propertyId = filter.bk_property_id
                        const property = (this.properties[modelId] || []).find(property => property.bk_property_id === propertyId)
                        if (property) {
                            const newCondition = {
                                bk_obj_id: modelId,
                                bk_property_id: propertyId,
                                bk_property_type: property.bk_property_type,
                                option: property.option,
                                operator: filter.operator,
                                value: filter.value
                            }
                            const existCondition = oldCondition.find(old => {
                                return old.bk_obj_id === property.bk_obj_id
                                    && old.bk_property_id === property.bk_property_id
                            })
                            if (existCondition) {
                                if (this.isCollection) {
                                    filterCondition.push(Object.assign(newCondition, existCondition))
                                } else {
                                    filterCondition.push(Object.assign(existCondition, newCondition))
                                }
                            } else {
                                filterCondition.push(newCondition)
                            }
                        }
                    })
                    this.filterCondition = filterCondition
                } catch (e) {
                    console.error(e)
                }
            },
            checkIsScrolling () {
                this.$nextTick(() => {
                    const scroller = this.$refs.scroller
                    this.isScrolling = scroller.scrollHeight > scroller.offsetHeight
                })
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
    .filter-trigger.is-active {
        color: #3A84FF;
    }
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
        padding: 0px 20px 10px;
        overflow: auto;
        @include scrollbar-y;
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
        padding: 10px 20px;
        &.is-sticky {
            margin: 0;
            background-color: #FAFBFD;
            border-top: 1px solid #DCDEE5;
        }
    }
    .collection {
        .collection-title {
            display: block;
            font-size: 13px;
            color: #63656E;
            line-height:17px;
        }
        .collection-name {
            margin-top: 13px;
        }
        .collection-options {
            padding: 20px 0 10px;
            text-align: right;
        }
    }
</style>
