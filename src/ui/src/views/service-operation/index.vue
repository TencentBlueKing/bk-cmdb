<template>
    <div class="layout" v-bkloading="{
        isLoading: !ready,
        immediate: true
    }">
        <template v-if="ready">
            <div class="info clearfix mb20">
                <label class="info-label fl">{{$t('已选实例')}}：</label>
                <i18n tag="div" path="N个" class="info-content">
                    <b class="info-count" place="count">{{serviceInstanceIds.length}}</b>
                </i18n>
            </div>
            <div class="info clearfix mb10" ref="changeInfo">
                <label class="info-label fl">{{$t('变更确认')}}：</label>
                <div class="info-content">
                    <ul class="tab clearfix">
                        <li class="tab-item active fl">
                            <span class="tab-label">{{$t('转移到空闲机的主机')}}</span>
                            <span class="tab-count">{{moveToIdleHosts.length}}</span>
                        </li>
                    </ul>
                    <div class="tab-empty" v-if="!moveToIdleHosts.length">
                        {{$t('仅移除服务实例，主机无变更')}}
                    </div>
                    <move-to-idle-host class="tab-component" :info="moveToIdleHosts" v-else></move-to-idle-host>
                </div>
            </div>
            <div class="options" :class="{ 'is-sticky': hasScrollbar }">
                <bk-button theme="primary" :loading="$loading(request.confirm)" @click="handleConfirm">{{$t('确认')}}</bk-button>
                <bk-button class="ml10" theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
            </div>
        </template>
    </div>
</template>

<script>
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    import { mapGetters } from 'vuex'
    import MoveToIdleHost from './children/move-to-idle-host.vue'
    export default {
        components: {
            MoveToIdleHost
        },
        data () {
            return {
                ready: false,
                hasScrollbar: false,
                moveToIdleHosts: [],
                request: {
                    preview: Symbol('review'),
                    confirm: Symbol('confirm'),
                    host: Symbol('host')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', [
                'getDefaultSearchCondition'
            ]),
            moduleId () {
                return this.$route.params.moduleId && parseInt(this.$route.params.moduleId)
            },
            serviceInstanceIds () {
                return String(this.$route.params.ids).split('/').map(id => parseInt(id, 10))
            }
        },
        watch: {
            ready (ready) {
                this.$nextTick(() => {
                    addResizeListener(this.$refs.changeInfo, this.resizeHandler)
                })
            }
        },
        async created () {
            this.getPreviewData()
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.changeInfo, this.resizeHandler)
        },
        methods: {
            async getPreviewData () {
                try {
                    const data = await this.$store.dispatch('serviceInstance/previewDeleteServiceInstances', {
                        params: {
                            bk_biz_id: this.bizId,
                            service_instance_ids: this.serviceInstanceIds
                        },
                        config: {
                            requestId: this.request.preview
                        }
                    })
                    this.setMoveToIdleHosts(data.to_move_module_hosts)
                } catch (e) {
                    console.error(e)
                }
            },
            async setMoveToIdleHosts (data = []) {
                try {
                    const hostIds = []
                    data.forEach(item => {
                        if (item.move_to_idle) {
                            hostIds.push(item.bk_host_id)
                        }
                    })
                    this.moveToIdleHosts = await this.getHostInfo(hostIds)
                    setTimeout(() => {
                        this.ready = true
                    }, 300)
                } catch (e) {
                    console.error(e)
                }
            },
            getHostInfo (hostIds) {
                return this.$store.dispatch('hostSearch/searchHost', {
                    params: this.getSearchHostParams(hostIds),
                    config: {
                        requestId: this.request.host
                    }
                }).then(data => data.info)
            },
            getSearchHostParams (hostIds) {
                const params = {
                    bk_biz_id: this.bizId,
                    ip: { data: [], exact: 0, flag: 'bk_host_innerip|bk_host_outerip' },
                    page: {},
                    condition: this.getDefaultSearchCondition()
                }
                const hostCondition = params.condition.find(target => target.bk_obj_id === 'host')
                hostCondition.condition.push({
                    field: 'bk_host_id',
                    operator: '$in',
                    value: hostIds
                })
                return params
            },
            resizeHandler () {
                this.$nextTick(() => {
                    const scroller = this.$el.parentElement
                    this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
                })
            },
            async handleConfirm () {
                try {
                    await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                        config: {
                            data: this.$injectMetadata({
                                service_instance_ids: this.serviceInstanceIds
                            }, { injectBizId: true }),
                            requestId: this.request.confirm
                        }
                    })
                    this.$success(this.$t('删除成功'))
                    this.redirect()
                } catch (e) {
                    console.error(e)
                }
            },
            handleCancel () {
                this.redirect()
            },
            redirect () {
                if (this.$route.query.from) {
                    this.$router.replace({
                        path: this.$route.query.from,
                        query: this.$route.query.query
                    })
                } else {
                    this.$router.replace({
                        name: MENU_BUSINESS_HOST_AND_SERVICE,
                        query: {
                            node: `module-${this.moduleId}`,
                            tab: 'serviceInstance'
                        }
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .layout {
        padding: 15px 0 0 0;
    }
    .info {
        .info-label {
            width: 128px;
            font-size: 14px;
            font-weight: bold;
            color: $textColor;
            text-align: right;
        }
        .info-content {
            overflow: hidden;
            padding: 0 20px 0 8px;
            font-size: 14px;
            .info-count {
                font-weight: bold;
            }
            .module-grep {
                border-top: 1px solid $borderColor;
                margin-top: 10px;
            }
        }
    }
    .module-list {
        font-size: 0;
        .module-item {
            position: relative;
            display: inline-block;
            vertical-align: middle;
            height: 26px;
            max-width: 150px;
            line-height: 24px;
            padding: 0 15px;
            margin: 0 10px 8px 0;
            border: 1px solid #C4C6CC;
            border-radius: 13px;
            color: $textColor;
            font-size: 12px;
            outline: none;
            cursor: default;
            @include ellipsis;
            &.is-business-module {
                padding: 0 12px 0 25px;
            }
            &.is-trigger {
                width: 40px;
                padding: 0;
                text-align: center;
                font-size: 0;
                cursor: pointer;
                .icon-cc-edit {
                    font-size: 14px;
                }
            }
            &:hover {
                border-color: $primaryColor;
                color: $primaryColor;
                .module-mask {
                    display: block;
                }
                .module-icon {
                    background-color: $primaryColor;
                }
            }
            .module-mask {
                display: none;
                position: absolute;
                left: 0;
                top: 0;
                width: 100%;
                height: 100%;
                color: #fff;
                background-color: rgba(0, 0, 0, 0.53);
                text-align: center;
                cursor: pointer;
            }
            .module-icon {
                position: absolute;
                left: 2px;
                top: 2px;
                width: 20px;
                height: 20px;
                border-radius: 50%;
                line-height: 20px;
                text-align: center;
                color: #FFF;
                font-size: 12px;
                background-color: #C4C6CC;
            }
            .module-remove {
                position: absolute;
                right: 4px;
                top: 4px;
                width: 16px;
                height: 16px;
                border-radius: 50%;
                text-align: center;
                line-height: 16px;
                color: #FFF;
                font-size: 0px;
                background-color: #C4C6CC;
                cursor: pointer;
                &:before {
                    display: inline-block;
                    vertical-align: middle;
                    font-size: 20px;
                    transform: translateX(-2px) scale(.5);
                }
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
                .tab-count {
                    color: #FFF;
                    background-color: $primaryColor;
                }
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
                position: relative;
                display: inline-block;
                vertical-align: middle;
                height: 16px;
                padding: 0 5px;
                border-radius: 4px;
                line-height: 16px;
                font-size: 12px;
                color: #FFF;
                background-color: #979BA5;
                &.has-badge:after {
                    position: absolute;
                    top: -3px;
                    right: -3px;
                    width: 6px;
                    height: 6px;
                    border-radius: 50%;
                    border: 1px solid #FFF;
                    background-color: $dangerColor;
                    content: "";
                }
            }
        }
    }
    .tab-component {
        margin-top: 20px;
    }
    .tab-empty {
        height: 60px;
        padding: 0 28px;
        margin-top: 24px;
        line-height: 60px;
        background-color: #F0F1F5;
        color: $textColor;
        &:before {
            content: "!";
            display: inline-block;
            width: 16px;
            height: 16px;
            line-height: 16px;
            border-radius: 50%;
            text-align: center;
            color: #FFF;
            font-size: 12px;
            background-color: #C4C6CC;
        }
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
