<template>
    <div ref="bkTab2" :class="['bk-tab2', {'bk-tab2-small': size === 'small'}]">
        <template v-if="needScroll">
            <div class="bk-tab2-control left" :class="leftScrollDisabled ? 'disabled' : ''" @click="scroll('left')">
                <i class="bk-tab2-icon bk-icon icon-angle-left"></i>
            </div>
            <div class="bk-tab2-control right" :class="rightScrollDisabled ? 'disabled' : ''" @click="scroll('right')">
                <i class="bk-tab2-icon bk-icon icon-angle-right"></i>
            </div>
        </template>
        <div ref="bkTab2Header" :class="['bk-tab2-head',{'is-fill': type === 'fill'}, {'scroll': needScroll}]">
            <div ref="bkTab2Nav" class="bk-tab2-nav" :style="bkTab2NavStyleObj">
                <div class="tab2-nav-item" :title="item.title"
                    :class="{'active': calcActiveName === item.name}"
                    @click="toggleTab(index, $event)"
                    v-if="item.show"
                    v-for="(item, index) in navList" >
                    <Render v-if="item.label" :render="item.label"></Render>
                    <div class="bk-panel-label" v-if="item.tag !== undefined">
                        <span class="bk-panel-tag">{{item.tag}}</span>
                    </div>
                    <span class="bk-panel-title">{{item.title}}</span>
                </div>
            </div>
            <div class="bk-tab2-action" v-if="hasSetting && !needScroll">
                <div class="action-wrapper">
                    <slot name="setting"></slot>
                </div>
            </div>
        </div>
        <div class="bk-tab2-content">
            <slot></slot>
        </div>
    </div>
</template>
<script>
    /**
     * bk-tab
     * @module components/tab
     * @desc 选项卡
     * @param type {String} 选项卡类型
     * @param size {String} 选项卡尺寸
     * @param active-name {String} 当前显示的子面板的别名，对应`bk-tabpanel`属性中的`name`字段
     * @example
     *  <bk-tab :size= "'small'" :active-name="'project1'" @tab-changed="tabChanged">
            <bk-tabpanel name="project1" title="项目一">
                <div class="p15 f12">项目一</div>
            </bk-tabpanel>
            <bk-tabpanel name="project2" title="项目二">
                <div class="p15 f12">项目二</div>
            </bk-tabpanel>
        </bk-tab>
     */

    import {getStyle} from '../../util'
    import Render from './render.js'

    export default {
        name: 'bk-tab',
        props: {
            activeName: {
                type: String,
                default: ''
            },
            type: {
                type: String,
                default: ''
            },
            size: {
                type: String,
                default: ''
            }
        },
        components: {
            Render
        },
        data () {
            return {
                navList: [],
                calcActiveName: this.activeName,
                hasSetting: false,
                needScroll: false,
                bkTab2NavStyleObj: {},
                translateX: 0,
                bkTab2HeaderWidth: 0,
                bkTab2NavWidth: 0,
                scrollDistance: 180,
                rightScrollDisabled: false,
                leftScrollDisabled: true
            }
        },
        watch: {
            activeName (val) {
                this.calcActiveName = val
            },
            calcActiveName () {
                let index = 0
                this.navList.forEach((nav, i) => {
                    if (nav.name === this.calcActiveName) {
                        index = i
                    }
                })
                this.$emit('tab-changed', this.calcActiveName, index)
            },
            translateX () {
                const node = this.$refs.bkTab2Nav
                if (!node) {
                    return
                }
                node.style.transform = node.style.webkitTransform = node.style.MozTransform
                    = node.style.msTransform = node.style.OTransform = `translateX(${this.translateX}px)`
            }
        },
        mounted () {
            if (!this.activeName) {
                this.calcActiveName = this.navList[0].name
            }
            this.hasSetting = !!this.$slots.setting
            this.$nextTick(() => {
                const bkTab2Width = parseInt(getStyle(this.$refs.bkTab2, 'width'), 10) || 0
                const bkTab2NavWidth = parseInt(getStyle(this.$refs.bkTab2Nav, 'width'), 10) || 0
                if (bkTab2Width < bkTab2NavWidth) {
                    this.needScroll = true
                    this.bkTab2NavStyleObj = {
                        transform: `translateX(${this.translateX}px)`
                    }
                } else {
                    this.bkTab2NavStyleObj = {
                        width: '100%'
                    }
                }
                this.bkTab2HeaderWidth = parseInt(getStyle(this.$refs.bkTab2Header, 'width'), 10) || 0
                this.bkTab2NavWidth = parseInt(getStyle(this.$refs.bkTab2Nav, 'width'), 10) || 0
            })
        },
        methods: {
            /**
             * 切换 tab
             */
            toggleTab (index) {
                this.calcActiveName = this.navList[index].name
                this.$emit('update:activeName', this.calcActiveName)
            },

            /**
             * 添加导航的 item，由子组件 tab-panel 调用
             *
             * @param {Object} item nav 对象
             */
            addNavItem (item) {
                this.navList.push(item)
                return this.navList.length - 1
            },

            /**
             * 滚动
             *
             * @param {string} direction 方向
             */
            scroll (direction) {
                this.rightScrollDisabled = false
                this.leftScrollDisabled = false
                let leftDistance = 0
                if (direction === 'right') {
                    leftDistance = this.bkTab2NavWidth - (this.bkTab2HeaderWidth - 60 + Math.abs(this.translateX))
                    if (leftDistance <= this.scrollDistance) {
                        this.translateX -= leftDistance
                        this.rightScrollDisabled = true
                    } else {
                        this.translateX -= this.scrollDistance
                    }
                } else {
                    if (Math.abs(this.translateX) <= this.scrollDistance) {
                        this.leftScrollDisabled = true
                    }
                    leftDistance = Math.abs(this.translateX)
                    if (leftDistance >= this.scrollDistance) {
                        this.translateX += this.scrollDistance
                    } else {
                        this.translateX += leftDistance
                    }
                }
            },

            /**
             * 子组件更新时调用
             *
             * @param {number} index 子组件的索引值
             * @param {Object} data 更新的内容
             */
            updateList (index, data) {
                const list = []
                this.navList.forEach(nav => {
                    list.push(Object.assign({}, nav))
                })

                const item = list[index]
                Object.keys(item).forEach(key => {
                    item[key] = data[key]
                })
                this.navList = [...list]
                // this.navList.splice(0, this.navList.length, ...list)

                // 当前选中的这个就是要消失的这一个，那么就选中第一个
                if (this.calcActiveName === data.name && !data.show) {
                    this.calcActiveName = this.navList[0].name
                }
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/tag.scss';
    @import '../../bk-magic-ui/src/tab.scss';
</style>
