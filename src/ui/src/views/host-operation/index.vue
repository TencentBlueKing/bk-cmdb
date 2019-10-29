<template>
    <div class="layout">
        <div class="info mb20">
            <label class="info-label">{{$t('已选主机')}}</label>
            <i18n path="N台主机" class="info-content">
                <b class="info-count" place="count">{{hostCount}}</b>
            </i18n>
        </div>
        <div class="info mb10" ref="changeInfo">
            <label class="info-label">{{$t('变更确认')}}</label>
            <div class="info-content">
                <ul class="tab clearfix">
                    <template v-for="(item, index) in tabList">
                        <li class="tab-grep fl" v-if="index" :key="index"></li>
                        <li class="tab-item fl"
                            :class="{ active: activeTab === item.id }"
                            :key="item.id"
                            @click="handleTabClick(item)">
                            <span class="tab-label">{{item.label}}</span>
                            <span class="tab-count">{{item.count}}</span>
                        </li>
                    </template>
                </ul>
                <keep-alive>
                    <component class="tab-component" :is="activeComponent"></component>
                </keep-alive>
            </div>
        </div>
        <div class="options" :class="{ 'is-sticky': hasScrollbar }">
            <bk-button theme="primary">{{$t('确认移除')}}</bk-button>
            <bk-button class="ml10" theme="default">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import DeletedServiceInstance from './children/deleted-service-instance.vue'
    import MoveToIdleHost from './children/move-to-idle-host.vue'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    export default {
        components: {
            DeletedServiceInstance,
            MoveToIdleHost
        },
        data () {
            return {
                hasScrollbar: false,
                hostCount: 108,
                activeComponent: null,
                activeTab: null,
                tabMap: {
                    deletedServiceInstance: {
                        id: Symbol('instance'),
                        label: this.$t('删除服务实例'),
                        count: 12,
                        order: 1,
                        component: DeletedServiceInstance
                    },
                    moveToIdleHost: {
                        id: Symbol('host'),
                        label: this.$t('移动到空闲机的主机'),
                        count: 12,
                        order: 2,
                        component: MoveToIdleHost
                    }
                }
            }
        },
        computed: {
            type () {
                return this.$route.params.type || 'remove'
            },
            tabList () {
                const map = {
                    remove: ['deletedServiceInstance', 'moveToIdleHost']
                }
                return map[this.type].map(tab => this.tabMap[tab]).sort((A, B) => A.order - B.order)
            }
        },
        mounted () {
            addResizeListener(this.$refs.changeInfo, this.resizeHandler)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.changeInfo, this.resizeHandler)
        },
        methods: {
            resizeHandler () {
                this.$nextTick(() => {
                    const scroller = this.$el.parentElement
                    this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
                })
            },
            handleTabClick (tab) {
                this.activeTab = tab.id
                this.activeComponent = tab.component
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info {
        display: flex;
        .info-label {
            flex: 128px 0 0;
            font-size: 14px;
            font-weight: bold;
            color: $textColor;
            text-align: right;
        }
        .info-content {
            flex: 500px 1 1;
            margin-left: 8px;
            padding-right: 20px;
            font-size: 14px;
            .info-count {
                font-weight: bold;
            }
        }
    }
    .tab {
        .tab-grep {
            width: 2px;
            height: 19px;
            margin: 0 8px;
            background-color: #C4C6CC;
        }
        .tab-item {
            position: relative;
            color: $textColor;
            font-size: 0;
            cursor: pointer;
            &.active {
                color: $primaryColor;
            }
            &.active:after {
                content: "";
                position: absolute;
                left: 0;
                top: 30px;
                width: 100%;
                height: 2px;
                background-color: $primaryColor;
            }
            .tab-label {
                display: inline-block;
                vertical-align: middle;
                margin-right: 7px;
                font-size: 14px;
            }
            .tab-count {
                display: inline-block;
                vertical-align: middle;
                height: 16px;
                padding: 0 5px;
                border-radius: 4px;
                line-height: 16px;
                font-size: 12px;
                color: #FFF;
                background-color: #979BA5;
            }
        }
    }
    .tab-component {
        margin-top: 20px;
    }
    .options {
        position: sticky;
        padding: 10px 0 10px 136px;
        font-size: 0;
        bottom: 0;
        left: 0;
        &.is-sticky {
            background-color: #FFF;
            border-top: 1px solid $borderColor;
            z-index: 100;
        }
    }
</style>
