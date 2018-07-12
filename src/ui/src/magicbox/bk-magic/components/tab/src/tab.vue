<template>
    <div :class="['bk-tab2', {'bk-tab2-small': size === 'small'}]">
        <div :class="['bk-tab2-head',{'is-fill': type === 'fill'}]" :style="headStyle">
            <ul class="bk-tab2-nav">
                <li class="tab2-nav-item"
                    v-for="(item, index) in navList"
                    v-show="item.show"
                    :class="{'active': calcActiveName === item.name}"
                    @click="toggleTab(index)">
                    {{ item.title }}
                </li>
            </ul>
            <div class="bk-tab2-action" v-if="hasSetting">
                <div class="action-wrapper">
                    <slot name="setting"></slot>
                </div>
            </div>
        </div>
        <div class="bk-tab2-content" :style="contentStyle">
            <slot></slot>
        </div>
    </div>
</template>
<script>
    /**
     * bk-tab
     * @module components/tab
     * @desc 开关
     * @param type {String} 选项卡类型
     * @param size {String} 选项卡尺寸
     * @param active-name {String} 当前显示的子面板的别名，对应`bk-tabpanel`属性中的`name`字段
     * @example
     *  <bk-tab :size= "'small'" :active-name="'exceptionDetecting'" @tab-changed="tabChanged">
            <bk-tabpanel name="curveSetting" title="项目一">
                <div class="p15 f12">项目一</div>
            </bk-tabpanel>
            <bk-tabpanel name="exceptionDetecting" title="项目二">
                <div class="p15 f12">项目二</div>
            </bk-tabpanel>
        </bk-tab>
     */
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
            },
            headStyle: {
                type: Object,
                default () {
                    return {}
                }
            },
            contentStyle: {
                type: Object,
                default () {
                    return {}
                }
            }
        },
        watch: {
            activeName (val) {
                this.calcActiveName = val
            }
        },
        mounted () {
            if (!this.activeName) {
                this.calcActiveName = this.navList[0].name
            }
        },
        data () {
            return {
                navList: [],
                calcActiveName: this.activeName,
                hasSetting: this.$slots.setting
            }
        },
        methods: {
            /**
             *  切换tab
             */
            toggleTab (index) {
                // this.calcActiveName = this.navList[index].name
                // this.$emit('tab-changed', this.calcActiveName, index)
                if (this.navList[index].name !== this.activeName) {
                    this.$emit('update:activeName', this.navList[index].name)
                    this.$emit('tab-changed', this.navList[index].name)
                }
            },
            addNavItem (item) {
                this.navList.push(item)
                return this.navList.length - 1
            },
            /**
             *  子组件更新时调用
             *  @param index - 子组件的索引值
             *  @param data - 更新的内容
             */
            updateList (index, data) {
                let list = this.navList
                let item = list[index]

                for (let key in item) {
                    item[key] = data[key]
                }
            }
        }
    }
</script>
