<template>
    <div class="graphics-toolbar clearfix">
        <div class="toolbar-left fl">
            <bk-button type="primary">{{$t('ModelManagement["编辑拓扑"]')}}</bk-button>
        </div>
        <div class="toolbar-right">
            <i class="toolbar-icon bk-icon icon-full-screen"
                v-tooltip="$t('ModelManagement[\'还原\']')"
                @click="handleResize">
            </i>
            <i class="toolbar-icon bk-icon icon-plus"
                v-tooltip="$t('ModelManagement[\'放大\']')"
                @click="handleZoom('in')">
            </i>
            <i class="toolbar-icon bk-icon icon-minus"
                v-tooltip="$t('ModelManagement[\'缩小\']')"
                @click="handleZoom('out')">
            </i>
            <i class="toolbar-icon icon-cc-setting"
                v-tooltip="$t('ModelManagement[\'拓扑显示设置\']')"
                @click="handleSetConfig">
            </i>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'graphics-toolabr',
        inject: ['parentRefs'],
        data () {
            return {}
        },
        methods: {
            getNetwork () {
                const { graphics } = this.parentRefs
                return graphics.instance.network
            },
            handleResize () {
                const network = this.getNetwork()
                network.moveTo({scale: 1})
                network.fit()
            },
            handleZoom (type) {
                const network = this.getNetwork()
                const ratio = type === 'in' ? 1.05 : 0.95
                network.moveTo({
                    scale: network.getScale() * ratio
                })
            },
            handleSetConfig () {}
        }
    }
</script>

<style lang="scss">
    .graphics-toolbar {
        padding: 0 20px;
        border-bottom: 1px solid $cmdbBorderColor;
        line-height: 49px;
        font-size: 0;
        .toolbar-right {
            text-align: right;
            overflow: hidden;
            .toolbar-icon {
                display: inline-block;
                margin: 0 0 0 32px;
                vertical-align: middle;
                font-size: 14px;
                cursor: pointer;
            }
        }
    }
</style>