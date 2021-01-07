<template>
    <div class="bk-button-group">
        <bk-button
            :class="{ 'is-selected': active === 'instance' }"
            @click="handleSwitch('instance')">
            {{$t('实例')}}
        </bk-button>
        <bk-button ref="tipsReference"
            :class="{ 'is-selected': active === 'process' }"
            v-bk-tooltips="{
                content: '#tipsContent',
                allowHtml: true,
                placement: 'bottom-end',
                disabled: tipsDisabled,
                showOnInit: !tipsDisabled,
                hideOnClick: false,
                trigger: 'manual',
                theme: 'view-switer-tips',
                zIndex: getZIndex()
            }"
            @click="handleSwitch('process')">
            {{$t('进程')}}
        </bk-button>
        <span class="tips-content" id="tipsContent">
            {{$t('切换进程视角提示语')}}
            <i class="bk-icon icon-close" @click="handleCloseTips"></i>
        </span>
    </div>
</template>

<script>
    import RouterQuery from '@/router/query'
    export default {
        data () {
            return {
                tipsDisabled: !!window.localStorage.getItem('service_instance_view_switcher'),
                active: RouterQuery.get('view', 'instance')
            }
        },
        watch: {
            active (active) {
                if (active === 'process') {
                    this.hideTips()
                }
            }
        },
        mounted () {
            if (this.active === 'process') {
                this.hideTips()
            }
        },
        methods: {
            getZIndex () {
                return window.__bk_zIndex_manager.nextZIndex()
            },
            handleSwitch (active) {
                RouterQuery.set({ 'view': active })
            },
            hideTips () {
                this.tipsDisabled = true
                const tippyInstance = this.$refs.tipsReference.$el.tippyInstance
                tippyInstance && tippyInstance.hide()
            },
            handleCloseTips () {
                window.localStorage.setItem('service_instance_view_switcher', 'closed')
                this.hideTips()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .bk-button-group {
        .tips-content {
            display: none;
        }
    }
    .tips-content {
        display: flex;
        align-items: center;
        position: relative;
        .bk-icon {
            margin-left: 20px;
            cursor: pointer;
            font-size: 16px;
        }
    }
</style>

<style lang="scss">
    .tippy-tooltip.view-switer-tips-theme {
        background-color: #699df4;
        color: #fff;
        .tippy-arrow {
            border-bottom-color: #699df4;
        }
    }
</style>
