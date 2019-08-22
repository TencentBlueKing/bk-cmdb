<template>
    <transition name="fade" duration="200">
        <div class="selector-layout" v-if="visible">
            <div class="host-wrapper">
                <h2 class="title">{{$t('添加主机')}}</h2>
                <div class="options">
                    <cmdb-selector class="options-selector"
                        v-model="filter.module"
                        :list="modules">
                    </cmdb-selector>
                    <cmdb-input class="options-filter" icon="bk-icon icon-search"
                        v-model.trim="filter.ip"
                        :placeholder="$t('请输入IP')"
                        @icon-click="infiniteIdentifier++"
                        @enter="infiniteIdentifier++">
                    </cmdb-input>
                    <i18n class="options-count fr" path="已选择主机">
                        <span place="count">{{checked.length}}</span>
                    </i18n>
                </div>
                <bk-table class="host-table"
                    :data="list"
                    :height="290"
                    @selection-change="handleSelectHost">
                    <bk-table-column type="selection" fixed width="60" align="center" class-name="bk-table-selection"></bk-table-column>
                    <bk-table-column v-for="column in header"
                        :key="column.id"
                        :prop="column.id"
                        :label="column.name">
                    </bk-table-column>
                    <infinite-loading slot="append" v-if="ready"
                        force-use-infinite-wrapper=".host-table .bk-table-body-wrapper"
                        :distance="42"
                        :identifier="infiniteIdentifier"
                        @infinite="infiniteHandler">
                        <span slot="no-more"></span>
                        <span slot="no-results"></span>
                        <span slot="error"></span>
                        <div slot="spinner" style="height: 42px;"
                            v-bkloading="{
                                isLoading: $loading(['getServiceProcessProperties', 'getHostSelectorList'])
                            }">
                        </div>
                    </infinite-loading>
                </bk-table>
                <div class="button-wrapper">
                    <bk-button class="button" theme="primary"
                        :disabled="!checked.length"
                        @click="handleConfirm">{{$t('确定')}}
                    </bk-button>
                    <bk-button class="button" theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
                </div>
            </div>
        </div>
    </transition>
</template>

<script>
    import infiniteLoading from 'vue-infinite-loading'
    export default {
        components: {
            infiniteLoading
        },
        props: {
            visible: Boolean,
            moduleInstance: {
                type: Object,
                default () {
                    return {}
                }
            },
            exclude: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                checked: [],
                storeChecked: [],
                list: [],
                header: [],
                pagination: {
                    current: 0,
                    size: 10,
                    count: 0
                },
                properties: [],
                internalModules: [],
                filter: {
                    module: 'biz',
                    ip: ''
                },
                infiniteIdentifier: 0,
                ready: false
            }
        },
        computed: {
            modules () {
                return [{
                    id: 'biz',
                    name: this.$t('业务主机')
                }, {
                    id: this.moduleInstance.bk_module_id,
                    name: this.moduleInstance.bk_module_name
                }].concat(this.internalModules)
            },
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            params () {
                const conditionMap = {
                    biz: [],
                    set: [],
                    module: [],
                    host: []
                }
                if (this.filter.module !== 'biz') {
                    conditionMap.module.push({
                        field: 'bk_module_id',
                        operator: '$eq',
                        value: this.filter.module
                    })
                }
                if (this.exclude.length) {
                    conditionMap.host.push({
                        field: 'bk_host_id',
                        operator: '$nin',
                        value: this.exclude
                    })
                }
                return {
                    bk_biz_id: this.business,
                    condition: ['biz', 'set', 'module', 'host'].map(id => {
                        return {
                            bk_obj_id: id,
                            fields: [],
                            condition: conditionMap[id]
                        }
                    }),
                    ip: {
                        data: this.filter.ip ? [this.filter.ip] : [],
                        exact: 0,
                        flag: 'bk_host_innerip|bk_host_outerip'
                    },
                    page: {
                        start: (this.pagination.current - 1) * this.pagination.size,
                        limit: this.pagination.size
                    }
                }
            }
        },
        watch: {
            visible (visible) {
                if (visible) {
                    (async () => {
                        await this.getProperties()
                        this.getInternalModule()
                        this.infiniteIdentifier++
                        this.ready = true
                    })()
                }
            },
            'filter.module' () {
                this.infiniteIdentifier++
            },
            infiniteIdentifier () {
                this.list = []
                this.checked = []
                this.pagination = {
                    current: 0,
                    count: 0,
                    size: 10
                }
            }
        },
        methods: {
            async getProperties () {
                try {
                    this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: {
                            bk_obj_id: 'host',
                            bk_supplier_account: this.$store.getters.supplierAccount
                        },
                        config: {
                            requestId: 'getServiceProcessProperties',
                            fromCache: true
                        }
                    })
                    this.setHeader()
                    return Promise.resolve()
                } catch (e) {
                    this.header = []
                    this.properties = []
                }
            },
            async getList () {
                try {
                    const data = await this.$store.dispatch('hostSearch/searchHost', {
                        params: this.params,
                        config: {
                            requestId: 'getHostSelectorList'
                        }
                    })
                    this.pagination.count = data.count
                    this.list.push(...this.$tools.flattenList(this.properties, data.info.map(item => item.host)))
                    return data
                } catch (e) {
                    console.error(e)
                    this.pagination.count = 0
                    this.list = []
                    this.checked = []
                    return {
                        count: 0,
                        info: []
                    }
                }
            },
            async getInternalModule () {
                try {
                    const data = await this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                        bizId: this.business,
                        config: {
                            requestId: `get_business_${this.business}_internal_module`,
                            fromCache: true
                        }
                    })
                    this.internalModules = data.module.map(module => {
                        return {
                            id: module.bk_module_id,
                            name: module.bk_module_name
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            setHeader () {
                this.header = [
                    'bk_host_innerip',
                    'bk_cloud_id',
                    'bk_host_outerip',
                    'bk_os_type',
                    'bk_host_name'
                ].map(id => {
                    const property = this.properties.find(property => property.bk_property_id === id) || {}
                    return {
                        id: id,
                        name: property.bk_property_name
                    }
                })
            },
            async infiniteHandler (infiniteState) {
                console.log(infiniteState)
                try {
                    const { current } = this.pagination
                    this.pagination.current = current + 1
                    const data = await this.getList()
                    infiniteState.loaded()
                    if (!data.info.length) {
                        infiniteState.complete()
                    }
                } catch (e) {
                    infiniteState.error()
                }
            },
            handleSelectHost (selection) {
                this.checked = selection.map(row => row.bk_host_id)
            },
            handleConfirm () {
                this.$emit('host-selected', [...this.checked], this.list.filter(item => this.checked.includes(item.bk_host_id)))
            },
            handleCancel () {
                this.$emit('update:visible', false)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .selector-layout {
        position: fixed;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        text-align: center;
        background-color: rgba(0, 0, 0, .6);
        z-index: 1500;
        &:before {
            content: "";
            width: 0;
            height: 100%;
            @include inlineBlock;
        }
    }
    .host-wrapper {
        width: 850px;
        max-height: 460px;
        padding: 15px 0 0;
        text-align: left;
        background-color: #fff;
        box-shadow:0px 4px 12px 0px rgba(0,0,0,0.2);
        border-radius:2px;
        @include inlineBlock;
        .title {
            margin: 0 24px;
            font-size: 24px;
            font-weight: normal;
            line-height: 31px;
        }
        .options {
            margin: 14px 24px 0;
            .options-selector {
                width: 200px;
            }
            .options-filter {
                width: 300px;
                margin-left: 10px;
            }
            .options-count {
                line-height: 36px;
                font-size: 12px;
                color: #63656E;
                [place="count"] {
                    font-weight: bold;
                    color: #2DCB56;
                }
            }
        }
        .host-table {
            width: 800px;
            margin: 10px 24px 0;
        }
    }
    .button-wrapper {
        height: 50px;
        margin: 18px 0 0;
        padding: 0 24px;
        line-height: 50px;
        text-align: right;
        font-size: 0;
        background-color: #FAFBFD;
        border-top: 1px solid #DCDEE5;
        .button {
            height: 32px;
            line-height: 30px;
            font-size: 14px;
            margin: 0 0 0 10px;
        }
    }
    .fade-enter-active, .fade-leave-active {
        transition: opacity .2s;
    }
    .fade-enter,
    .fade-leave-to {
        opacity: 0;
    }
</style>
