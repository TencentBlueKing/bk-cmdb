<template>
    <div class="graphics-toolbar clearfix">
        <div class="toolbar-left fl">
            <template v-if="isEditMode">
                <bk-button theme="primary"
                    @click="handleToggleMode">
                    {{$t('返回')}}
                </bk-button>
                <span class="edit-tips">{{$t('所有更改已自动保存')}}</span>
            </template>
            <template v-else>
                <bk-button theme="primary"
                    @click="handleToggleMode">
                    {{$t('编辑拓扑')}}
                </bk-button>
            </template>
        </div>
        <div class="toolbar-right">
            <i class="toolbar-icon bk-icon icon-full-screen"
                v-bk-tooltips="$t('还原')"
                @click="handleResize">
            </i>
            <i class="toolbar-icon bk-icon icon-plus"
                v-bk-tooltips="$t('放大')"
                @click="handleZoom('in')">
            </i>
            <i class="toolbar-icon bk-icon icon-minus"
                v-bk-tooltips="$t('缩小')"
                @click="handleZoom('out')">
            </i>
            <i class="toolbar-icon icon-cc-setting"
                v-bk-tooltips="$t('拓扑显示设置')"
                @click="handleSetConfig">
            </i>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        name: 'graphics-toolabr',
        inject: ['parentRefs'],
        data () {
            return {
                editMode: false
            }
        },
        computed: {
            ...mapGetters('globalModels', ['isEditMode'])
        },
        methods: {
            getGraphics () {
                const { graphics } = this.parentRefs
                return graphics.instance
            },
            handleToggleMode () {
                this.$store.commit('globalModels/changeEditMode')
            },
            handleResize () {
                const graphics = this.getGraphics()
                graphics.resize()
            },
            handleZoom (type) {
                const graphics = this.getGraphics()
                graphics.zoom(type)
            },
            handleSetConfig () {
                const { config } = this.parentRefs
                config.toggleSlider()
            }
        }
    }
</script>

<style lang="scss">
    .graphics-toolbar {
        padding: 0 20px;
        line-height: 50px;
        font-size: 0;
        .toolbar-left {
            .edit-tips {
                display: inline-block;
                margin: 0 0 0 10px;
                vertical-align: middle;
                font-size: 14px;
                color: #a4aab3;
            }
        }
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
