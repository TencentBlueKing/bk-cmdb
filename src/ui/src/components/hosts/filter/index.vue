<template>
    <div class="hosts-filter-layout" :class="{close}">
        <i class="filter-toggle bk-icon icon-angle-right" @click="close = !close"></i>
        <div class="filter-main">
            <bk-tab class="filter-tab" size="small" :active-name.sync="tab.active" style="padding: 0">
                <bk-tabpanel name="filter" :title="$t('HostResourcePool[\'筛选\']')" v-if="activeTab.includes('filter')">
                    <keep-alive>
                        <the-filter ref="theFilter"
                            v-if="tab.active === 'filter'"
                            :filter-config-key="filterConfigKey"
                            @on-refresh="handleRefresh">
                            <slot name="business" slot="business"></slot>
                            <slot name="scope" slot="scope"></slot>
                        </the-filter>
                    </keep-alive>
                </bk-tabpanel>
                <bk-tabpanel name="collection" :title="$t('Hosts[\'收藏\']')" v-if="activeTab.includes('collection')">
                    <the-collection v-if="tab.active === 'collection'"></the-collection>
                </bk-tabpanel>
                <bk-tabpanel name="history" :title="$t('Hosts[\'历史\']')" v-if="activeTab.includes('history')">
                    <the-history v-if="tab.active === 'history'"></the-history>
                </bk-tabpanel>
                <the-setting slot="setting"
                    :active-setting="activeSetting"
                    :filter-config-key="filterConfigKey"
                    @on-reset="handleReset">
                </the-setting>
            </bk-tab>
        </div>
    </div>
</template>

<script>
    import theFilter from './_filter'
    import theCollection from './_collection'
    import theHistory from './_history'
    import theSetting from './_setting.vue'
    export default {
        components: {
            theFilter,
            theCollection,
            theHistory,
            theSetting
        },
        props: {
            activeTab: {
                type: Array,
                default () {
                    return ['filter', 'collection', 'history']
                }
            },
            activeSetting: {
                type: Array,
                default () {
                    return ['reset', 'collection', 'filter-config']
                }
            },
            filterConfigKey: {
                type: String,
                required: true
            }
        },
        data () {
            return {
                close: false,
                tab: {
                    active: 'filter'
                }
            }
        },
        methods: {
            handleReset () {
                this.$refs.theFilter.reset()
            },
            handleRefresh (params) {
                this.$emit('on-refresh', params)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .hosts-filter-layout{
        position: relative;
        height: 100%;
        &.close{
            .filter-toggle{
                transform: rotate(180deg);
                border-top-right-radius: 12px;
                border-bottom-right-radius: 12px;
                border-top-left-radius: 0;
                border-bottom-left-radius: 0;
            }
            .filter-main{
                display: none;
            }
        }
        .filter-toggle{
            position: absolute;
            right: 100%;
            top: 50%;
            width: 14px;
            height: 100px;
            margin: -50px  0 0 0;
            line-height: 100px;
            color: #fff;
            font-size: 12px;
            text-align: center;
            border-top-left-radius: 12px;
            border-bottom-left-radius: 12px;
            background-color: #c3cdd7;
            transition: background-color .2s ease;
            cursor: pointer;
            &:hover{
                background-color: #6b7baa;
            }
        }
    }
    .filter-main{
        width: 358px;
        height: 100%;
        padding: 10px 20px;
        border-left: 1px solid $cmdbBorderColor;
    }
</style>

<style lang="scss">
    .hosts-filter-layout{
        .bk-tab2.filter-tab .bk-tab2-head .bk-tab2-nav .tab2-nav-item{
            padding: 0 15px !important;
        }
    }
</style>