<template>
    <bk-popover
        ref="filterPopper"
        placement="bottom"
        theme="light hosts-filter-shadow"
        trigger="manual"
        :width="350"
        :on-show="handleShow"
        :on-hide="handleHide"
        :z-index="1001"
        :tippy-options="{
            interactive: true,
            hideOnClick: false,
            onShown: checkIsScrolling
        }">
        <icon-button class="filter-trigger"
            icon="icon-cc-funnel"
            v-bk-tooltips.top="$t('高级筛选')"
            :class="{
                'is-active': isFilterActive
            }"
            @click="handleToggleFilter">
        </icon-button>
        <section class="filter-content" slot="content"
            :style="{
                height: (sectionHeight ? sectionHeight : ($APP.height - 200)) + 'px'
            }">
            <h2 class="filter-title">
                {{$t('高级筛选')}}
                <bk-button class="close-trigger" text icon="close" @click="handleToggleFilter"></bk-button>
            </h2>
            <div class="filter-scroller" ref="scroller">
                <div class="filter-group" style="padding: 0;">
                    <label class="filter-label">IP</label>
                    <bk-input type="textarea" v-model="ip.text" :rows="4" :placeholder="$t('请输入IP，多个IP请使用换行分隔')"></bk-input>
                </div>
                <div class="filter-group checkbox-group">
                    <bk-checkbox class="filter-checkbox"
                        v-model="ip.inner"
                        :disabled="!ip.outer">
                        {{$t('内网IP')}}
                    </bk-checkbox>
                    <bk-checkbox class="filter-checkbox"
                        v-model="ip.outer"
                        :disabled="!ip.inner">
                        {{$t('外网IP')}}
                    </bk-checkbox>
                    <bk-checkbox class="filter-checkbox" v-model="ip.exact">{{$t('精确')}}</bk-checkbox>
                </div>
                <div class="filter-group"
                    v-for="(filterItem, index) in filterCondition"
                    :key="filterItem.bk_property_id">
                    <label class="filter-label">{{getFilterLabel(filterItem)}}</label>
                    <div class="filter-condition">
                        <filter-operator class="filter-operator"
                            v-if="!['date', 'time', 'category'].includes(filterItem.bk_property_type)"
                            v-model="filterItem.operator"
                            :type="getOperatorType(filterItem)">
                        </filter-operator>
                        <component class="filter-value"
                            v-if="['enum', 'list'].includes(filterItem.bk_property_type)"
                            :is="`cmdb-form-${filterItem.bk_property_type}`"
                            :options="filterItem.option || []"
                            v-model="filterItem.value">
                        </component>
                        <cmdb-form-bool-input class="filter-value"
                            v-else-if="filterItem.bk_property_type === 'bool'"
                            v-model="filterItem.value">
                        </cmdb-form-bool-input>
                        <cmdb-search-input class="filter-value" :style="{ '--index': 99 - index }"
                            v-else-if="['singlechar', 'longchar'].includes(filterItem.bk_property_type)"
                            v-model="filterItem.value">
                        </cmdb-search-input>
                        <cmdb-form-date-range class="filter-value"
                            v-else-if="['date', 'time'].includes(filterItem.bk_property_type)"
                            v-model="filterItem.value">
                        </cmdb-form-date-range>
                        <cmdb-cloud-selector
                            v-else-if="filterItem.bk_property_id === 'bk_cloud_id'"
                            class="filter-value"
                            v-model="filterItem.value">
                        </cmdb-cloud-selector>
                        <cmdb-service-category
                            v-else-if="filterItem.bk_property_id === 'service_category_id'"
                            class="filter-value"
                            v-model="filterItem.value">
                        </cmdb-service-category>
                        <component class="filter-value"
                            v-else
                            :is="`cmdb-form-${filterItem.bk_property_type}`"
                            :unit="filterItem.unit"
                            v-model="filterItem.value">
                        </component>
                        <i class="bk-icon icon-close" @click.stop="handleDeleteFilter(filterItem)"></i>
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
                <template v-if="isBusinessHost">
                    <div class="fl">
                        <bk-button theme="primary" style="margin-right: 6px;" @click="handleSearch">{{$t('查询')}}</bk-button>
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
                            :z-index="1002"
                            :tippy-options="{
                                interactive: true,
                                hideOnClick: false,
                                onShown: focusCollectionName,
                                onHidden: clearCollectionName
                            }">
                            <bk-button theme="default" @click="handleCreateCollection">{{$t('收藏此条件')}}</bk-button>
                            <section class="collection" slot="content">
                                <label class="collection-title">{{$t('收藏此条件')}}</label>
                                <bk-input class="collection-name"
                                    ref="collectionName"
                                    :placeholder="$t('请填写名称')"
                                    :data-vv-as="$t('名称')"
                                    data-vv-name="collectionName"
                                    v-model="collectionName"
                                    v-validate="'required|length:256'"
                                    @enter="handleSaveCollection">
                                </bk-input>
                                <p class="collection-error" v-if="errors.has('collectionName')">{{errors.first('collectionName')}}</p>
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
                </template>
                <template v-else>
                    <bk-button theme="primary" class="mr5" @click="handleSearch">{{$t('查询')}}</bk-button>
                    <bk-button theme="default" @click="handleReset">{{$t('清空')}}</bk-button>
                </template>
            </div>
        </section>
        <property-selector :properties="properties" ref="propertySelector"></property-selector>
    </bk-popover>
</template>

<script>
    import filterOperator from './_filter-field-operator.vue'
    import propertySelector from './filter-property-selector.vue'
    import { mapState, mapGetters } from 'vuex'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    import Bus from '@/utils/bus'
    import RouterQuery from '@/router/query'
    import { getIPPayload } from '@/utils/host'
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
            sectionHeight: {
                type: Number,
                default: null
            }
        },
        data () {
            const defaultIpConfig = {
                text: '',
                inner: true,
                outer: true,
                exact: true
            }
            return {
                ip: {
                    ...defaultIpConfig
                },
                filterCondition: [],
                defaultIpConfig,
                isScrolling: false,
                collectionName: '',
                propertyPromise: null,
                propertyResolver: null,
                isShow: false
            }
        },
        computed: {
            ...mapState('hosts', ['filterList', 'collection', 'shouldInjectAsset']),
            ...mapGetters('hosts', ['isCollection']),
            isFilterActive () {
                const hasIP = !!this.ip.text.replace(/\n|;|；|,|，/g, '').length
                const hasField = this.filterCondition.some(condition => {
                    const filterValue = condition.value
                    return filterValue !== null
                        && filterValue !== undefined
                        && !!String(filterValue).length
                })
                return hasIP || hasField || this.isShow
            },
            isBusinessHost () {
                return this.$route.name === MENU_BUSINESS_HOST_AND_SERVICE
            }
        },
        watch: {
            filterList (newList, oldList) {
                this.setFilterCondition()
            },
            filterCondition () {
                this.checkIsScrolling()
            },
            properties (properties) {
                this.propertyResolver && this.propertyResolver()
            },
            '$route.name' () {
                this.handleToggleFilter(false)
            }
        },
        async created () {
            Bus.$on('toggle-host-filter', this.handleToggleFilter)
            Bus.$on('reset-host-filter', this.handleReset)
            this.propertyPromise = new Promise((resolve, reject) => {
                this.propertyResolver = () => {
                    this.propertyResolver = null
                    resolve()
                }
            })
            this.unwatch = RouterQuery.watch(['ip', 'exact', 'inner', 'outer'], ({
                ip = '',
                exact = this.ip.exact ? '1' : '0',
                inner = this.ip.inner ? '1' : '0',
                outer = this.ip.outer ? '1' : '0'
            }) => {
                this.ip.text = ip.split(',').join('\n')
                this.ip.exact = parseInt(exact) === 1
                this.ip.inner = parseInt(inner) === 1
                this.ip.outer = parseInt(outer) === 1
            }, { immediate: true })
            await this.initCustomFilterList()
        },
        beforeDestroy () {
            Bus.$off('toggle-host-filter', this.handleToggleFilter)
            Bus.$off('reset-host-filter', this.handleReset)
            this.unwatch()
            this.$store.commit('hosts/clearFilter')
        },
        methods: {
            initCustomFilterList () {
                const key = this.$route.meta.filterPropertyKey
                const customData = this.$store.getters['userCustom/getCustomData'](key, [])
                if (!customData.length && !this.isCollection) {
                    customData.push({
                        bk_obj_id: 'host',
                        bk_property_id: 'operator',
                        operator: '$eq',
                        value: ''
                    }, {
                        bk_obj_id: 'host',
                        bk_property_id: 'bk_cloud_id',
                        operator: '$eq',
                        value: ''
                    })
                }
                this.$store.commit('hosts/setFilterList', customData)
            },
            handleToggleFilter (visible) {
                const instance = this.$refs.filterPopper.instance
                if (typeof visible === 'boolean') {
                    visible ? instance.show() : instance.hide(0)
                } else {
                    const state = instance.state
                    if (state.isVisible) {
                        instance.hide()
                    } else {
                        instance.show()
                    }
                }
            },
            handleAddFilter () {
                // fix遮盖层问题
                window.__bk_zIndex_manager.nextZIndex()
                this.$refs.propertySelector.isShow = true
            },
            async handleDeleteFilter (filterItem) {
                const conditionList = this.filterCondition.filter(item => !(item.bk_obj_id === filterItem.bk_obj_id && item.bk_property_id === filterItem.bk_property_id))
                const list = conditionList.map(condition => {
                    return {
                        bk_obj_id: condition.bk_obj_id,
                        bk_property_id: condition.bk_property_id,
                        operator: condition.operator,
                        value: condition.value
                    }
                })
                if (!this.isCollection) {
                    const userCustomList = list.map(item => {
                        return {
                            ...item,
                            operator: '',
                            value: ''
                        }
                    })
                    const key = this.$route.meta.filterPropertyKey
                    await this.$store.dispatch('userCustom/saveUsercustom', {
                        [key]: userCustomList
                    })
                }
                this.$store.commit('hosts/setShouldInjectAsset', false)
                this.$store.commit('hosts/setFilterList', list)
            },
            handleSearch (toggle = true) {
                const params = this.getParams()
                this.updateQuery(params)
                toggle && this.handleToggleFilter()
            },
            handleCreateCollection () {
                const instance = this.$refs.collectionPopover.instance
                this.errors.clear()
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
                    const isValid = await this.$validator.validate('collectionName')
                    if (!isValid) {
                        return false
                    }
                    const data = this.getCollectionParams()
                    const result = await this.$store.dispatch('hostFavorites/createFavorites', {
                        params: data,
                        config: {
                            requestId: 'createCollection'
                        }
                    })
                    this.$success(this.$t('收藏成功'))
                    this.$store.commit('hosts/addCollection', Object.assign({}, data, result))
                    this.handleCancelCollection()
                } catch (e) {
                    console.error(e)
                }
            },
            getCollectionParams () {
                return {
                    bk_biz_id: this.$store.getters['objectBiz/bizId'],
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
            focusCollectionName () {
                this.$refs.collectionName.$refs.input.focus()
            },
            clearCollectionName () {
                this.collectionName = ''
            },
            handleCancelCollection () {
                const instance = this.$refs.collectionPopover.instance
                instance.hide()
            },
            handleReset () {
                this.ip = { ...this.defaultIpConfig }
                this.filterCondition.forEach(filterItem => {
                    filterItem.value = ''
                })
                const params = this.getParams()
                this.updateQuery(params)
            },
            getParams () {
                const inputIPPlayload = {
                    data: this.getIPList(),
                    inner: this.ip.inner,
                    outer: this.ip.outer,
                    exact: this.ip.exact
                }
                const params = {
                    ip: getIPPayload(inputIPPlayload),
                    host: [],
                    module: [],
                    set: [],
                    biz: [],
                    object: []
                }
                this.filterCondition.forEach(filterItem => {
                    const filterValue = filterItem.value
                    if (filterValue !== null && filterValue !== undefined && String(filterValue).length) {
                        const modelId = filterItem.bk_obj_id
                        if (['date', 'time'].includes(filterItem.bk_property_type)) {
                            params[modelId].push(...[{
                                field: filterItem.bk_property_id,
                                operator: '$gte',
                                value: filterItem.value[0]
                            }, {
                                field: filterItem.bk_property_id,
                                operator: '$lte',
                                value: filterItem.value[1]
                            }])
                        } else {
                            let operator = filterItem.operator
                            let value = filterValue

                            if (['category', 'organization'].includes(filterItem.bk_property_type)) {
                                operator = '$in'
                            }

                            if (filterItem.operator === '$multilike' || ['bk_set_name', 'bk_module_name'].includes(filterItem.bk_property_id)) {
                                value = filterValue.split(/\n|;|；|,|，/).filter(str => str.trim().length).map(str => str.trim())
                            }

                            params[modelId].push({
                                field: filterItem.bk_property_id,
                                operator,
                                value
                            })
                        }
                    }
                })
                return params
            },
            updateQuery (params) {
                this.$store.commit('hosts/setCondition', ['biz', 'set', 'module', 'host', 'object'].map(modelId => {
                    return {
                        bk_obj_id: modelId,
                        fields: [],
                        condition: params[modelId]
                    }
                }))
                const flags = params.ip.flag.split('|')
                const query = {
                    ip: params.ip.data.join(','),
                    exact: params.ip.exact,
                    inner: flags.includes('bk_host_innerip') ? 1 : 0,
                    outer: flags.includes('bk_host_outerip') ? 1 : 0,
                    page: 1,
                    _t: Date.now()
                }
                const assetCondition = params.host.find(condition => {
                    return condition.field === 'bk_asset_id' && condition.operator === '$in'
                })
                query.bk_asset_id = assetCondition ? assetCondition.value.toString() : ''
                RouterQuery.set(query)
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
            async setFilterCondition () {
                try {
                    await this.propertyPromise
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
                                value: filter.value,
                                unit: property.unit
                            }
                            const existCondition = oldCondition.find(old => {
                                return old.bk_obj_id === property.bk_obj_id
                                    && old.bk_property_id === property.bk_property_id
                            })
                            if (existCondition) {
                                if (this.isCollection) {
                                    filterCondition.push(Object.assign(existCondition, newCondition))
                                } else {
                                    filterCondition.push(Object.assign(newCondition, existCondition))
                                }
                            } else {
                                filterCondition.push(newCondition)
                            }
                        }
                    })
                    this.injectAssetCondition(filterCondition)
                    this.filterCondition = filterCondition
                } catch (e) {
                    console.error(e)
                }
            },
            injectAssetCondition (filterCondition) {
                const assetIds = RouterQuery.get('bk_asset_id', '').split(',').filter(id => !!id.length)
                if (!this.shouldInjectAsset || !assetIds.length) {
                    this.$store.commit('hosts/setShouldInjectAsset', true)
                    return
                }
                const injectCondition = {
                    bk_obj_id: 'host',
                    bk_property_id: 'bk_asset_id',
                    bk_property_type: 'singlechar',
                    operator: '$in',
                    option: '',
                    value: assetIds.join('\n')
                }
                const assetCondition = filterCondition.find(condition => {
                    return condition.bk_obj_id === injectCondition.bk_obj_id
                        && condition.bk_property_id === injectCondition.bk_property_id
                        && condition.bk_property_type === injectCondition.bk_property_type
                })
                if (assetCondition) {
                    Object.assign(assetCondition, injectCondition)
                } else {
                    filterCondition.unshift(injectCondition)
                }
            },
            checkIsScrolling () {
                this.$nextTick(() => {
                    const scroller = this.$refs.scroller
                    this.isScrolling = scroller.scrollHeight > scroller.offsetHeight
                })
            },
            handleShow (popper) {
                this.isShow = true
                popper.popperChildren.tooltip.style.padding = 0
            },
            handleHide () {
                this.isShow = false
                const collectionPopover = this.$refs.collectionPopover
                if (collectionPopover && collectionPopover.instance.state.isShown) {
                    collectionPopover.instance.hide()
                }
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
            },
            hide () {
                try {
                    const instance = this.$refs.filterPopper.instance
                    instance.hide()
                } catch (e) {}
            }
        }
    }
</script>

<style lang="scss" scoped="true">
    .filter-content {
        border: 1px solid #DCDEE5;
    }
    .filter-trigger {
        width: 32px;
        padding: 0;
        line-height: 30px;
    }
    .filter-trigger {
        &.is-active {
            color: #3A84FF;
            &:hover {
                color: #3A84FF;
            }
        }
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
            font-size: 20px;
            color: #979BA5;
            &:hover {
                color: #63656E;
            }

            /deep/ .bk-icon {
                font-size: 22px;
            }
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
                > div {
                    display: flex;
                    justify-content: center;
                    align-items: center;
                }
                .bk-icon {
                    position: initial;
                    font-size: 20px;
                }
            }
        }
    }
    .filter-condition {
        display: flex;
        &:hover .icon-close{
            opacity: 1;
        }
        .filter-operator {
            flex: 75px 0 0;
            margin-right: 8px;
        }
        .filter-value {
            width: 0;
            flex: 1;
            &.cmdb-search-input {
                /deep/ .search-input-wrapper {
                    z-index: var(--index);
                }
            }
        }
        .icon-close {
            color: #d8d8d8;
            font-size: 20px;
            font-weight: bold;
            line-height: 32px;
            height: 32px;
            margin: 0 0 0 3px;
            cursor: pointer;
            opacity: 0;
            &:hover {
                color: #63656e;
            }
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
        .collection-error {
            color: $dangerColor;
        }
        .collection-options {
            padding: 20px 0 10px;
            text-align: right;
        }
    }
</style>

<style lang="scss">
    .hosts-filter-shadow-theme {
        box-shadow: 0px 1px 6px 0px rgba(220,222,229,1) !important;
    }
</style>
