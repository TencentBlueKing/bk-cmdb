<template>
    <transition name="slide">
        <div class="slidebar-wrapper" v-if="isShow" @click.self="quickClose">
            <div class="sideslider" :style="{width: `${width}px`}">
                <div class="sideslider-title" ref="title">
                    <slot name="title">
                        <h3 class="title">
                            <span class="vm">{{title}}</span>
                        </h3>
                    </slot>
                </div>
                <div ref="content" class="sideslider-content">
                    <slot name="content"></slot>
                </div>
            </div>
        </div>
    </transition>
</template>

<script>
    export default {
        name: 'cmdb-slider',
        props: {
            /*
                标题
            */
            title: {
                default: ''
            },
            /*
                弹窗显示状态
            */
            isShow: {
                type: Boolean,
                default: false
            },
            /*
                弹窗宽度
            */
            width: {
                default: 800
            },
            /*
                是否支持点击空白处关闭
            */
            hasQuickClose: {
                type: Boolean,
                default: true
            },
            beforeClose: {
                type: Function
            }

        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.calcContentHeight()
                    setTimeout(() => {
                        this.$emit('shown')
                    }, 200)
                } else {
                    setTimeout(() => {
                        this.$emit('close')
                    }, 200)
                }
            }
        },
        methods: {
            async closeSlider () {
                if (typeof this.beforeClose === 'function') {
                    let confirmed
                    try {
                        confirmed = await Promise.resolve(this.beforeClose())
                    } catch (e) {
                        confirmed = false
                    }
                    if (confirmed) {
                        this.$emit('update:isShow', false)
                    }
                } else {
                    this.$emit('update:isShow', false)
                }
            },
            quickClose () {
                if (this.hasQuickClose) {
                    this.closeSlider()
                }
            },
            calcContentHeight () {
                this.$nextTick(() => {
                    this.$refs.content.style.height = `calc(100% - ${this.$refs.title.getBoundingClientRect().height}px)`
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    $primaryColor: #6b7baa;
    $lineColor: #e7e9ef;
    .slide-leave,
    .slide-enter-active {
        .sideslider{
            transition: all linear .2s;
            right: 0;
        }
    }
    .slide-enter,
    .slide-leave-active{
        .sideslider{
            right: -100%;
        }
    }
    .slidebar-wrapper{
        font-size: 14px;
        transition: all .2s;
        width: 100%;
        height: 100%;
        z-index: 1300;
        position: fixed;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        background-color: rgba(0,0,0,0.6);
    }
    .sideslider{
        transition: all .2s linear;
        position: absolute;
        top: 0;
        bottom: 0;
        right: 0;
        background: #fff;
        box-shadow: -4px 0px 6px 0px rgba(0, 0, 0, 0.06);
        .close{
            position: absolute;
            top: 0;
            left: -30px;
            width: 30px;
            height: 60px;
            padding: 10px 7px 0;
            background-color: #ef4c4c;
            box-shadow: -2px 0 2px 0 rgba(0, 0, 0, 0.2);
            cursor: pointer;
            color: #fff;
            font-size: 14px;
            line-height: 20px;
            font-weight: normal;
            &:hover{
                background: #e13d3d;
            }
        }
    }
    .sideslider-title{
        position: relative;
        padding-left: 20px;
        line-height: 60px;
        color: #333948;
        font-size: 14px;
        height: 60px;
        font-weight: bold;
        margin: 0;
        background: #f9f9f9;
        .icon-mainframe{
            position: relative;
            top: 0px;
            margin-right: 5px;
            display: inline-block;
            width: 19px;
            height: 19px;
        }
        .title{
            font-size: 14px;
        }
    }
    .sideslider-content{
        height: calc(100% - 60px);
        position: relative;
    }
</style>