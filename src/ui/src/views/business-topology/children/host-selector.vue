<template>
    <div class="selector-layout">
        <div class="host-wrapper">
            <h2 class="title">{{$t('BusinessTopology["添加主机"]')}}</h2>
            <div class="options">
                <cmdb-selector class="options-selector"></cmdb-selector>
                <cmdb-input class="options-filter" icon="bk-icon icon-search"></cmdb-input>
                <i18n class="options-count fr" path="BusinessTopology['已选择主机']">
                    <span place="count">{{checked.length}}</span>
                </i18n>
            </div>
            <cmdb-table class="host-table"
                :header="header"
                :list="list"
                :pagination="pagination"
                :height="286"
                :cross-page-check="false">
            </cmdb-table>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                checked: [],
                list: [],
                header: [],
                pagination: {
                    current: 1,
                    size: 10,
                    count: 0
                },
                properties: []
            }
        },
        computed: {
            params () {
                return {
                    bk_obj_id: this.$store.getters['objectBiz/bizId'],
                    condition: ['biz', 'set', 'module', 'host'].map(id => {
                        return {
                            bk_obj_id: id,
                            fields: [],
                            condition: []
                        }
                    }),
                    ip: {
                        data: [],
                        exact: 0,
                        flag: 'bk_host_innerip|bk_host_outerip'
                    },
                    page: {
                        start: (this.current - 1) * this.pagination.size,
                        limit: this.pagination.size
                    }
                }
            }
        },
        async created () {
            await this.getProperties()
            this.getList()
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
                            requestId: 'get_service_process_properties',
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
                    this.list = this.$tools.flatternList(this.properties, data.info.map(item => item.host))
                } catch (e) {
                    console.error(e)
                    this.pagination.count = 0
                    this.list = []
                }
            },
            setHeader () {
                this.header = [{
                    id: 'bk_host_id',
                    type: 'checkbox'
                }].concat([
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
                }))
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
        z-index: 9999;
        &:before {
            content: "";
            width: 0;
            height: 100%;
            @include inlineBlock;
        }
    }
    .host-wrapper {
        width: 850px;
        height: 460px;
        padding: 15px 24px 0;
        text-align: left;
        background-color: #fff;
        box-shadow:0px 4px 12px 0px rgba(0,0,0,0.2);
        border-radius:2px;
        @include inlineBlock;
        .title {
            font-size: 24px;
            font-weight: normal;
            line-height: 31px;
        }
        .options {
            margin: 14px 0 0 0;
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
            margin: 10px 0 0 0;
        }
    }
</style>
