<template>
    <div class="bk-transfer" ref="transfer">
        <div class="source-list">
            <div class="slot-header" v-if="slot['left-header']">
                <div class="slot-content">
                    <slot name="left-header"></slot>
                </div>
            </div>
            <div class="header" v-else>
                {{title[0] ? title[0] : '左侧列表'}}（共{{dataList.length}}条）
                <span v-if="dataList.length === 0" class="disabled">全部添加</span>
                <span v-else @click="allToRight">全部添加</span>
            </div>
            <template v-if="dataList.length">
                <ul class="content">
                    <li v-for="(item, index) in dataList"
                        @click.stop.prevent="leftClick(index)"
                        @mouseover.stop.prevent="leftMouseover(index)"
                        @mouseleave.stop.prevent="leftMouseleave(index)" :key="index">
                        <span>{{item[displayCode]}}</span>
                        <span class="icon-wrapper" :class="[index === leftHoverIndex ? 'hover' : '']"><i class="bk-icon icon-arrows-right"></i></span>
                    </li>
                </ul>
            </template>
            <template v-else>
                <div v-if="slot['left-empty-content']">
                    <slot name="left-empty-content"></slot>
                </div>
                <div class="empty" v-else>
                    {{emptyContent[0] ? emptyContent[0] : '无数据'}}
                </div>
            </template>
        </div>
        <div class="transfer">
        </div>
        <div class="target-list">
            <div class="slot-header" v-if="slot['right-header']">
                <div class="slot-content">
                    <slot name="right-header"></slot>
                </div>
            </div>
            <div class="header" v-else>
                {{title[1] ? title[1] : '右侧列表'}}（共{{hasSelectedList.length}}条）
                <span v-if="hasSelectedList.length === 0" class="disabled">全部移除</span>
                <span v-else @click="allToLeft">全部移除</span>
            </div>
            <template v-if="hasSelectedList.length">
                <ul class="content">
                    <li v-for="(item, index) in hasSelectedList"
                        @click.stop.prevent="rightClick(index)"
                        @mouseover.stop.prevent="rightMouseover(index)"
                        @mouseleave.stop.prevent="rightMouseleave(index)" :key="index">
                        <span>{{item[displayCode]}}</span>
                        <span class="icon-wrapper" :class="[index === rightHoverIndex ? 'hover' : '']"><i class="bk-icon icon-close"></i></span>
                    </li>
                </ul>
            </template>
            <template v-else>
                <div v-if="slot['right-empty-content']">
                    <slot name="right-empty-content"></slot>
                </div>
                <div class="empty" v-else>
                    {{emptyContent[1] ? emptyContent[1] : '未选择任何项'}}
                </div>
            </template>
        </div>
    </div>
</template>


<script>
    /**
     *  bk-transfer
     *  @module components/transfer
     *  @desc 穿梭框组件
     *  @param title {Array} - 顶部title(title[0]: 左侧title,title[1]: 右侧title,)
     *  @param emptyContent {Array} - 穿梭框无数据时提示文案
     *  @param displayKey {String} - 循环list时，显示字段的key值(当list为普通数组时可不传传了也无效)
     *  @param settingKey {String} - 唯一key值
     *  @param sortKey {String} - 排序所依据的key(当list为普通数组时可不传，默认按照index值排序)
     *  @param sortable {Boolean} - 是否开启排序功能
     *  @param sourceList {Array} - 穿梭框数据源(支持普通数组)
     *  @param targetList {Array} - 默认已选择的数据源
     *  @example
    <bk-transfer
        :source-list="list"
        :target-list="value"
        :title="title"
        :empty-content="emptyContent"
        :display-key="'value'"
        :setting-key="'label'"
        :sort-key="'label'"
        :is-sort="true"
        @change="change">
    </bk-transfer>
    */
    export default {
        name: 'bk-transfer',
        props: {
            title: {
                type: Array,
                default () {
                    return ['左侧列表', '右侧列表']
                }
            },
            emptyContent: {
                type: Array,
                default () {
                    return ['无数据', '未选择任何项']
                }
            },
            displayKey: {
                type: String,
                default: 'value'
            },
            settingKey: {
                type: String,
                default: 'id'
            },
            sortKey: {
                type: String,
                default: ''
            },
            sourceList: {
                type: Array,
                default () {
                    return []
                }
            },
            targetList: {
                type: Array,
                default () {
                    return []
                }
            },
            hasHeader: {
                type: Boolean,
                default: false
            },
            sortable: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                dataList: [],
                hasSelectedList: [],
                sortList: [],
                leftHoverIndex: -1,
                rightHoverIndex: -1,
                slot: {},
                displayCode: this.displayKey,
                settingCode: this.settingKey,
                sortCode: this.sortKey,
                isSortFlag: this.sortable
            }
        },
        computed: {
            typeFlag () {
                if (!this.sourceList || !this.sourceList.length) {
                    return 'empty'
                } else {
                    const str = this.sourceList.toString()
                    if (str.indexOf('[object Object]') !== -1) {
                        return true
                    } else {
                        return false
                    }
                }
            }
        },
        watch: {
            sourceList: {
                handler: function (value) {
                    if (!value || !value.length) {
                        return
                    }
                    this.initData()
                    this.initSort()
                },
                deep: true
            },
            targetList: {
                handler: function (value) {
                    this.initData()
                    this.initSort()
                },
                deep: true
            },
            displayKey (value) {
                this.displayCode = value
                this.initData()
            },
            settingKey (value) {
                this.settingCode = value
                this.initData()
            },
            sortKey (value) {
                this.sortCode = value
                this.initSort()
            },
            sortable (value) {
                this.isSortFlag = value
                this.initSort()
            }
        },
        created () {
            if (this.typeFlag !== 'empty') {
                if (!this.typeFlag) {
                    this.generalInit()
                } else {
                    this.init()
                }
                this.initSort()
            }
            this.slot = Object.assign({}, this.$slots)
        },
        methods: {
            /**
             * 数据改变初始化
             */
            initData () {
                if (this.typeFlag !== 'empty') {
                    if (!this.typeFlag) {
                        this.generalInit()
                    } else {
                        this.init()
                    }
                }
            },

            /**
             * 普通数组数据初始化
             */
            generalInit () {
                if (!this.targetList.length || this.targetList.length > this.sourceList.length) {
                    const list = []
                    for (let i = 0; i < this.sourceList.length; i++) {
                        list.push({
                            index: i,
                            value: this.sourceList[i]
                        })
                    }
                    this.dataList = [...list]
                    this.hasSelectedList.splice(0, this.hasSelectedList.length, ...[])
                    this.$emit('change', this.dataList, [], [])
                } else {
                    const list = []
                    const valueList = []
                    for (let i = 0; i < this.sourceList.length; i++) {
                        list.push({
                            index: i,
                            value: this.sourceList[i]
                        })
                    }
                    for (let j = 0; j < list.length; j++) {
                        for (let k = 0; k < this.targetList.length; k++) {
                            if (list[j].value === this.targetList[k]) {
                                valueList.push(list[j])
                            }
                        }
                    }
                    this.hasSelectedList = [...valueList]
                    const result = list.filter(item1 => {
                        return valueList.every(item2 => {
                            return item2['index'] !== item1['index']
                        })
                    })
                    this.dataList = [...result]
                    this.$emit('change', this.dataList, [...this.generalListHandler(this.hasSelectedList)], [])
                }
            },

            /**
             * 对象数组数据初始化
             */
            init () {
                if (!this.targetList.length || this.targetList.length > this.sourceList.length) {
                    this.dataList = [...this.sourceList]
                    this.hasSelectedList = []
                    this.$emit('change', this.dataList, [], [])
                } else {
                    const result = this.sourceList.filter(item1 => {
                        return this.targetList.every(item2 => {
                            return item2 !== item1[this.settingCode]
                        })
                    })
                    const hasTempList = []
                    this.sourceList.forEach(item1 => {
                        this.targetList.forEach(item2 => {
                            if (item1[this.settingCode] === item2) {
                                hasTempList.push(item1)
                            }
                        })
                    })
                    this.hasSelectedList = [...hasTempList]
                    this.dataList = [...result]
                    const list = [...this.sourceListHandler(this.hasSelectedList)]
                    this.$emit('change', this.dataList, this.hasSelectedList, list)
                }
            },

            /**
             * 普通数组数据处理
             */
            generalListHandler (list) {
                const templateList = []
                if (!list.length) {
                    return []
                } else {
                    const dataList = [...list]
                    dataList.forEach(item => {
                        templateList.push(item.value)
                    })
                    return templateList
                }
            },

            /**
             * 对象数组数据处理
             */
            sourceListHandler (list) {
                const templateList = []
                if (!list.length) {
                    return []
                } else {
                    const dataList = [...list]
                    dataList.forEach(item => {
                        for (let key in item) {
                            if (key === this.settingCode) {
                                templateList.push(item[key])
                            }
                        }
                    })
                    return templateList
                }
            },

            /**
             * 初始化排序
             */
            initSort () {
                let templateList = []
                if (!this.typeFlag) {
                    if (this.isSortFlag) {
                        this.sortCode = 'index'
                    } else {
                        this.sortCode = ''
                    }
                    for (let k = 0; k < this.sourceList.length; k++) {
                        templateList.push({
                            index: k,
                            value: this.sourceList[k]
                        })
                    }
                } else {
                    if (!this.isSortFlag) {
                        this.sortCode = ''
                    }
                    templateList = [...this.sourceList]
                }
                if (this.sortCode) {
                    const arr = []
                    templateList.forEach(item => {
                        arr.push(item[this.sortCode])
                    })
                    this.sortList = [...arr]
                    if (this.sortList.length === this.sourceList.length) {
                        const list = [...this.dataList]
                        this.dataList = [...this.sortDataList(list, this.sortCode, this.sortList)]
                    }
                }
            },

            /**
             * 排序处理
             */
            sortDataList (list, key, sortList) {
                const arr = sortList
                return list.sort((a, b) => arr.indexOf(a[key]) - arr.indexOf(b[key]) >= 0)
            },
            
            /**
             * 全部添加到右侧穿梭框
             */
            allToRight () {
                this.leftHoverIndex = -1
                const dataList = this.dataList
                let index = 0
                const hasSelectedList = this.hasSelectedList
                while (dataList.length) {
                    const transferItem = dataList.shift()
                    hasSelectedList.push(transferItem)
                    if (this.sortList.length === this.sourceList.length) {
                        this.hasSelectedList = [...this.sortDataList(hasSelectedList, this.sortCode, this.sortList)]
                    } else {
                        this.hasSelectedList = [...hasSelectedList]
                    }
                    index++
                }
                if (!this.typeFlag) {
                    this.$emit('change', this.dataList, [...this.generalListHandler(this.hasSelectedList)], [])
                } else {
                    const list = [...this.sourceListHandler(this.hasSelectedList)]
                    this.$emit('change', this.dataList, this.hasSelectedList, list)
                }
            },

            /**
             * 全部移除到左侧穿梭框
             */
            allToLeft () {
                this.rightHoverIndex = -1
                const hasSelectedList = this.hasSelectedList
                let index = 0
                const dataList = this.dataList
                while (hasSelectedList.length) {
                    const transferItem = hasSelectedList.shift()
                    dataList.push(transferItem)
                    if (this.sortList.length === this.sourceList.length) {
                        this.dataList = [...this.sortDataList(dataList, this.sortCode, this.sortList)]
                    } else {
                        this.dataList = [...dataList]
                    }
                    index++
                }
                if (!this.typeFlag) {
                    this.$emit('change', this.dataList, [...this.generalListHandler(this.hasSelectedList)], [])
                } else {
                    const list = [...this.sourceListHandler(this.hasSelectedList)]
                    this.$emit('change', this.dataList, this.hasSelectedList, list)
                }
            },

            /**
             * 左侧穿梭框 点击 事件
             *
             * @param {number} index item 的索引
             */
            leftClick (index) {
                this.leftHoverIndex = -1
                const transferItem = this.dataList.splice(index, 1)[0]
                const hasSelectedList = this.hasSelectedList
                hasSelectedList.push(transferItem)
                if (this.sortList.length === this.sourceList.length) {
                    this.hasSelectedList = [...this.sortDataList(hasSelectedList, this.sortCode, this.sortList)]
                } else {
                    this.hasSelectedList = [...hasSelectedList]
                }
                if (!this.typeFlag) {
                    this.$emit('change', this.dataList, [...this.generalListHandler(this.hasSelectedList)], [])
                } else {
                    const list = [...this.sourceListHandler(this.hasSelectedList)]
                    this.$emit('change', this.dataList, this.hasSelectedList, list)
                }
            },

            /**
             * 右侧穿梭框 点击 事件
             *
             * @param {number} index item 的索引
             */
            rightClick (index) {
                this.rightHoverIndex = -1
                const transferItem = this.hasSelectedList.splice(index, 1)[0]
                const dataList = this.dataList
                dataList.push(transferItem)
                if (this.sortList.length === this.sourceList.length) {
                    this.dataList = [...this.sortDataList(dataList, this.sortCode, this.sortList)]
                } else {
                    this.dataList = [...dataList]
                }
                if (!this.typeFlag) {
                    this.$emit('change', this.dataList, [...this.generalListHandler(this.hasSelectedList)], [])
                } else {
                    const list = [...this.sourceListHandler(this.hasSelectedList)]
                    this.$emit('change', this.dataList, this.hasSelectedList, list)
                }
            },

            /**
             * 左侧穿梭框 mouseover 事件
             *
             * @param {number} index hover 的索引
             */
            leftMouseover (index) {
                this.leftHoverIndex = index
            },

            /**
             * 左侧穿梭框 mouseleave 事件
             *
             * @param {number} index hover 的索引
             */
            leftMouseleave (index) {
                this.leftHoverIndex = -1
            },

            /**
             * 右侧穿梭框 mouseover 事件
             *
             * @param {number} index hover 的索引
             */
            rightMouseover (index) {
                this.rightHoverIndex = index
            },

            /**
             * 右侧穿梭框 mouseleave 事件
             *
             * @param {number} index hover 的索引
             */
            rightMouseleave (index) {
                this.rightHoverIndex = -1
            }
        }
    }
</script>

<style lang="scss">
    @import '../../bk-magic-ui/src/transfer.scss';
</style>
