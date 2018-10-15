<template>
    <div class="bk-selector"
        :class="[extCls, {'open': open}]"
        @click="openFn"
        @keyup.enter="updateSelect"
        v-clickoutside="close">
        <div class="bk-selector-wrapper">
            <input class="bk-selector-input" readonly="readonly"
                :class="{placeholder: selectedText === defaultPlaceholder, active: open}"
                :value="selectedText"
                :placeholder="defaultPlaceholder"
                :disabled="disabled"
                @mouseover="showClearFn"
                @mouseleave="showClear = false">
            <i :class="['bk-icon icon-angle-down bk-selector-icon',{'disabled': disabled}]" v-show="!isLoading && !showClear"></i>
            <i :class="['bk-icon icon-close bk-selector-icon clear-icon']"
                v-show="!isLoading && showClear"
                @mouseover="showClearFn"
                @mouseleave="showClear = false"
                @click="clearSelected($event)">
            </i>
            <div class="bk-spin-loading bk-spin-loading-mini bk-spin-loading-primary selector-loading-icon" v-show="isLoading">
                <div class="rotate rotate1"></div>
                <div class="rotate rotate2"></div>
                <div class="rotate rotate3"></div>
                <div class="rotate rotate4"></div>
                <div class="rotate rotate5"></div>
                <div class="rotate rotate6"></div>
                <div class="rotate rotate7"></div>
                <div class="rotate rotate8"></div>
            </div>
        </div>

        <transition :name="listSlideName">
            <div class="bk-selector-list" v-show="!isLoading && open" :style="panelStyle">
                <!-- 搜索栏 -->
                <div class="bk-selector-search-item"
                    @click="$event.stopPropagation()"
                    v-if="searchable">
                    <i class="bk-icon icon-search"></i>
                    <input type="text" v-model="condition" @input="inputFn" ref="searchNode" :placeholder='searchPlaceholder'>
                </div>
                <ul :style="{'max-height': `${contentMaxHeight}px`}" class='outside-ul'>
                    <li :class="['bk-selector-list-item', item.children && item.children.length ? 'bk-selector-group-list-item' : '', {'active': activeIndex === i}]"
                        v-if="localList.length !== 0"
                        :key='i'
                        v-for="(item, i) in localList">
                        <!-- 有分组 start -->
                        <template v-if="item.children && item.children.length && hasChildren">
                            <div class="bk-selector-group-name">{{item[displayKey]}}</div>
                            <ul class="bk-selector-group-list">
                                <li v-for="(child, index) in item.children" class="bk-selector-list-item" :class="[{'active': itemIndex === index && groupIndex === i && hasChildren}]">
                                    <div class="bk-selector-node bk-selector-sub-node"
                                        :class="{'bk-selector-selected': !multiSelect && child[settingKey] === selected}">
                                        <div class="text" @click.stop="selectItem(child, $event)">
                                            <label class="bk-form-checkbox bk-checkbox-small mr0 bk-selector-multi-label" v-if="multiSelect">
                                                <input type="checkbox"
                                                    :name="'multiSelect' + +new Date()"
                                                    :value="child[settingKey]"
                                                    v-model="localSelected">
                                                    {{ child[displayKey] }}
                                            </label>
                                            <template v-else>
                                                {{ child[displayKey] }}
                                            </template>
                                        </div>
                                        <div class="bk-selector-tools" v-if="tools !== false">
                                            <i class="bk-icon icon-edit2 bk-selector-list-icon"
                                                v-if="tools.edit !== false"
                                                @click.stop="editFn(index)"></i>
                                            <i class="bk-icon icon-close bk-selector-list-icon"
                                                v-if="tools.del !== false"
                                                @click.stop="delFn(index)"></i>
                                        </div>
                                    </div>
                                </li>
                            </ul>
                        </template>
                        <!-- 有分组 end -->

                        <!-- 没分组 start -->
                        <template v-else>
                            <div class="bk-selector-node" :class="{'bk-selector-selected': !multiSelect && item[settingKey] === selected, 'is-disabled': item.isDisabled}">
                                <div class="text" @click.stop="selectItem(item, $event)" :title="item[displayKey]">
                                    <label class="bk-form-checkbox bk-checkbox-small mr0 bk-selector-multi-label" v-if="multiSelect">
                                        <input type="checkbox"
                                        :name="'multiSelect' + +new Date()"
                                        :value="item[settingKey]"
                                        v-model="localSelected">
                                        {{ item[displayKey] }}
                                    </label>
                                    <template v-else>
                                        {{ item[displayKey] }}
                                    </template>
                                </div>
                                <div class="bk-selector-tools" v-if="tools !== false">
                                    <i class="bk-icon icon-edit2 bk-selector-list-icon"
                                        v-if="tools.edit !== false"
                                        @click.stop="editFn(index)"></i>
                                    <i class="bk-icon icon-close bk-selector-list-icon"
                                        v-if="tools.del !== false"
                                        @click.stop="delFn(index)"></i>
                                </div>
                            </div>
                        </template>
                        <!-- 没分组 end -->
                    </li>
                    <li class="bk-selector-list-item" v-if="!isLoading && localList.length === 0">
                        <div class="text no-search-result">
                            {{list.length ? defaultSearchEmptyText : defaultEmptyText}}
                        </div>
                    </li>
                </ul>
                <!-- 新增项 start -->
                <slot></slot>
            </div>
        </transition>
    </div>
</template>

<script>
    /**
     *  bk-dropdown
     *  @module components/dropdown
     *  @desc 下拉选框组件，模拟原生select
     *  @param extCls {String} - 自定义的样式
     *  @param hasCreateItem {Boolean} - 下拉菜单中是否有新增项，默认为true
     *  @param tools {Object, Boolean} - 待选项右侧的工具按钮，有两个可配置的key：edit和del，默认为两者都不显示。
     *  @param list {Array} - 必选，下拉菜单所需的数据列表
     *  @param filterList {Array} - 过滤列表
     *  @param selected {Number} - 必选，选中的项的index值，支持.sync修饰符
     *  @param placeholder {String, Boolean} - 是否显示占位行，默认为显示，且文字为“请选择”
     *  @param displayKey {String} - 循环list时，显示字段的key值，默认为name
     *  @param disabled {Boolean} - 是否禁用组件，默认为false
     *  @param multiSelect {Boolean} - 是否支持多选，默认为false
     *  @param searchable {Boolean} - 是否支持筛选，默认为false
     *  @param searchKey {Boolean} - 筛选时，搜索的key值，默认为'name'
     *  @param allowClear {Boolean} - 是否可以清除单选时选中的项
     *  @param settingKey {String} - 根据配置这个字段，自定义在单选时，选中某项之后的回调函数的第一个返回值的内容
     *  @example
        <bk-selector
            :list="list"
            :tools="tools"
            :selected.sync="selected"
            :placeholder="placeholder"
            :displayKey="displayKey"
            :has-create-item="hasCreateItem"
            :ext-cls="extCls"></bk-dropdown>
    */

    import clickoutside from '../../directives/clickoutside'
    import {getActualTop, getActualLeft} from '../../util'
    import locale from '../../mixins/locale'

    export default {
        name: 'bk-selector',
        mixins: [locale],
        directives: {
            clickoutside
        },
        props: {
            extCls: {
                type: String
            },
            searchPlaceholder: {
                type: String
            },
            isLoading: {
                type: Boolean,
                default: false
            },
            hasCreateItem: {
                type: Boolean,
                default: false
            },
            hasChildren: {
                type: [Boolean, String],
                default: false
            },
            tools: {
                type: [Object, Boolean],
                default: false
            },
            list: {
                type: Array,
                required: true
            },
            filterList: {
                type: Array,
                default () {
                    return []
                }
            },
            selected: {
                type: [Number, Array, String],
                required: true
            },
            placeholder: {
                type: [String, Boolean],
                default: ''
            },
            // 是否联动
            isLink: {
                type: [String, Boolean],
                default: false
            },
            displayKey: {
                type: String,
                default: 'name'
            },
            disabled: {
                type: [String, Boolean, Number],
                default: false
            },
            multiSelect: {
                type: Boolean,
                default: false
            },
            searchable: {
                type: Boolean,
                default: false
            },
            searchKey: {
                type: String,
                default: 'name'
            },
            allowClear: {
                type: Boolean,
                default: false
            },
            settingKey: {
                type: String,
                default: 'id'
            },
            initPreventTrigger: {
                type: Boolean,
                default: false
            },
            emptyText: {
                type: String,
                default: ''
            },
            searchEmptyText: {
                type: String,
                default: ''
            },
            contentMaxHeight: {
                type: Number,
                default: 300
            }
        },
        data () {
            return {
                open: false,
                selectedList: this.calcSelected(this.selected),
                condition: '',
                // localList: this.list,
                localSelected: this.selected,
                // emptyText: this.list.length ? '无匹配数据' : '暂无数据',
                showClear: false,
                activeIndex: -1,
                groupIndex: 0,
                itemIndex: -1,
                isKeydown: false,
                listInterval: null,
                panelStyle: {},
                listSlideName: 'toggle-slide',
                defaultPlaceholder: this.t('selector.pleaseselect'),
                defaultEmptyText: this.t('selector.emptyText'),
                defaultSearchEmptyText: this.t('selector.searchEmptyText')
            }
        },
        computed: {
            localList () {
                if (!this.multiSelect) {
                    this.list.forEach(item => {
                        if (this.filterList.includes(item[this.settingKey])) {
                            item.isDisabled = true
                        } else {
                            item.isDisabled = false
                        }
                    })
                }
                if (this.searchable && this.condition) {
                    const arr = []
                    const key = this.searchKey
                    const len = this.list.length
                    for (let i = 0; i < len; i++) {
                        const item = this.list[i]
                        if (item.children) {
                            const results = []
                            const childLen = item.children.length
                            for (let j = 0; j < childLen; j++) {
                                const child = item.children[j]
                                if (child[key].toLowerCase().includes(this.condition.toLowerCase())) {
                                    results.push(child)
                                }
                            }
                            if (results.length) {
                                const cloneItem = Object.assign({}, item)
                                cloneItem.children = results
                                arr.push(cloneItem)
                            }
                        } else {
                            if (item[key].toLowerCase().includes(this.condition.toLowerCase())) {
                                arr.push(item)
                            }
                        }
                    }

                    // for (let item of this.list) {
                    //     if (item.children) {
                    //         const results = []
                    //         for (let child of item.children) {
                    //             if (child[key].includes(this.condition)) {
                    //                 results.push(child)
                    //             }
                    //         }
                    //         if (results.length) {
                    //             const cloneItem = Object.assign({}, item)
                    //             cloneItem.children = results
                    //             arr.push(cloneItem)
                    //         }
                    //     } else {
                    //         if (item[key].includes(this.condition)) {
                    //             arr.push(item)
                    //         }
                    //     }
                    // }

                    return arr
                }
                return this.list
            },
            currentItem () {
                return this.list[this.localSelected]
            },
            selectedText () {
                let text = this.defaultPlaceholder
                let textArr = []
                if (Array.isArray(this.selectedList) && this.selectedList.length) {
                    this.selectedList.forEach(item => {
                        textArr.push(item[this.displayKey])
                    })
                } else if (this.selectedList) {
                    this.selectedList[this.displayKey] && textArr.push(this.selectedList[this.displayKey])
                }

                return textArr.length ? textArr.join(',') : this.defaultPlaceholder
            }
        },
        watch: {
            selected (newVal) {
                // 重新生成选择列表
                if (this.list.length) {
                    this.selectedList = this.calcSelected(this.selected, this.isLink)
                }
                this.localSelected = this.selected
            },
            list (newVal) {
                // 重新生成选择列表
                // this.localList = this.list
                if (this.selected) {
                    this.selectedList = this.calcSelected(this.selected, this.isLink)
                } else {
                    this.selectedList = []
                }
            },
            localSelected (val) {
                // 重新生成选择列表
                if (this.list.length) {
                    this.selectedList = this.calcSelected(this.localSelected, this.isLink)
                }
            },
            open (newVal) {
                const searchNode = this.$refs.searchNode
                if (searchNode) {
                    if (newVal) {
                        this.$nextTick(() => {
                            searchNode.focus()
                        })
                    }
                }
                this.$emit('visible-toggle', newVal)

                if (newVal) {
                    window.onkeydown = e => {
                        switch (e.keyCode) {
                            case 38:
                                // 上
                                this.listUp()
                                break

                            case 40:
                                // 下
                                this.listDown()
                                break

                            default:
                                break
                        }
                    }

                    document.onkeydown = function (e) {
                        e = e || event
                        if (e.keyCode === 38) {
                            return false
                        }
                        if (e.keyCode === 40) {
                            return false
                        }
                    }
                } else {
                    window.onkeydown = null
                    document.onkeydown = null
                }
            }
        },
        created () {
            if (this.placeholder) {
                this.defaultPlaceholder = this.placeholder
            }
            if (this.emptyText) {
                this.defaultEmptyText = this.emptyText
            }
            if (this.searchEmptyText) {
                this.defaultSearchEmptyText = this.searchEmptyText
            }
        },
        mounted () {
            this.popup = this.$el
            if (this.isLink) {
                if (this.list.length && this.selected) {
                    this.calcSelected(this.selected, this.isLink)
                }
            }
        },
        methods: {
            // 用于展开列表时控制鼠标上下建的事件绑定
            listUp () {
                let maxIndex = 0
                if (this.hasChildren) {
                    let arr = []
                    this.localList.forEach((list, index) => {
                        list.children.forEach((l, i) => {
                            arr.push(l)
                        })
                    })

                    if (this.groupIndex > 0) {
                        this.itemIndex--
                        if (this.itemIndex < 0) {
                            this.itemIndex = this.localList[this.groupIndex - 1].children.length - 1
                            this.groupIndex--
                        }
                    } else {
                        this.groupIndex = this.localList.length - 1
                        this.itemIndex = this.localList[this.groupIndex].children.length - 1
                    }
                } else {
                    maxIndex = this.localList.length - 1
                    if (this.activeIndex > 0) {
                        this.activeIndex--
                    } else {
                        this.activeIndex = maxIndex
                    }
                }
                // 检测更新滚动
                this.updateScrollTop()
            },
            listDown () {
                let maxIndex = 0
                if (this.hasChildren) {
                    let arr = []
                    this.localList.forEach((list, index) => {
                        list.children.forEach((l, i) => {
                            arr.push(l)
                        })
                    })

                    if (this.groupIndex < this.localList.length) {
                        this.itemIndex++
                        if (this.itemIndex > this.localList[this.groupIndex].children.length - 1) {
                            this.groupIndex++
                            this.itemIndex = 0
                        }
                        if (this.groupIndex === this.localList.length) {
                            this.groupIndex = 0
                            this.itemIndex = 0
                        }
                    } else {
                        this.groupIndex = 0
                        this.itemIndex = 0
                    }
                } else {
                    maxIndex = this.localList.length - 1
                    if (this.activeIndex < maxIndex) {
                        this.activeIndex++
                    } else {
                        this.activeIndex = 0
                    }
                }

                // 检测更新滚动
                this.updateScrollTop()
            },
            updateScrollTop (type) {
                // 获取下拉列表容器的位置信息，用于判断上下键选中的元素是否在可视区域，若不在则需滚动至可视区域
                const panelObj = this.$el.querySelector('.bk-selector-list .outside-ul')
                let panelInfo = {
                    height: panelObj.clientHeight,
                    yAxios: panelObj.getBoundingClientRect().y
                }
                this.$nextTick(() => {
                    const activeObj = this.$el.querySelector('.bk-selector-list .active')
                    let activeInfo = {
                        height: activeObj.clientHeight,
                        yAxios: activeObj.getBoundingClientRect().y
                    }

                    // 选中元素顶部坐标大于容器顶部坐标时，则该元素有部分或者全部区域不在可视区域，执行滚动
                    if (activeInfo.yAxios < panelInfo.yAxios) {
                        let currentScTop = panelObj.scrollTop
                        panelObj.scrollTop = currentScTop - (panelInfo.yAxios - activeInfo.yAxios)
                    }

                    let distanceToBottom = activeInfo.yAxios + activeInfo.height - panelInfo.yAxios

                    // 选中元素底部坐标大于容器顶部坐标，且超出容器的实际高度，则该元素有部分或者全部区域不在可视区域，执行滚动
                    if (distanceToBottom > panelInfo.height) {
                        let currentScTop = panelObj.scrollTop
                        panelObj.scrollTop = currentScTop + distanceToBottom - panelInfo.height
                    }
                })
            },
            // 单选下拉列表支持键盘上下键选择
            updateSelect () {
                let arr = []
                let activeItem
                if (!this.hasChildren) {
                    activeItem = this.localList[this.activeIndex]
                } else {
                    activeItem = this.localList[this.groupIndex].children[this.itemIndex]
                }
                if (this.multiSelect) {
                    // 若是多选，则enter键触发时判断是否已选，然后进行反选操作
                    const isAdded = this.selectedList.some(item => item[this.settingKey] === activeItem[this.settingKey])
                    if (isAdded) {
                        this.selectedList = this.selectedList.filter(item => item[this.settingKey] !== activeItem[this.settingKey])
                        this.localSelected = this.localSelected.filter(item => item !== activeItem[this.settingKey])
                        this.$emit('update:selected', this.localSelected)
                    } else {
                        this.selectedList.push(activeItem)
                        this.localSelected.push(activeItem[this.settingKey])
                        this.localSelected.sort()
                    }
                    if (Array.isArray(this.localSelected)) {
                        let data = []
                        for (let key of this.localSelected) {
                            let params = this.getItem(key)
                            if (params) {
                                data.push(params)
                            }
                        }
                        this.$emit('item-selected', this.localSelected, data)
                    }
                    return
                }
                this.$el.querySelector('.bk-selector-list .active').click()
                this.selectedList = this.localList[this.activeIndex]
                this.$emit('update:selected', this.selectedList[this.settingKey])
                let data = {}
                this.list.forEach(item => {
                    if (item[this.settingKey] === this.selectedList[this.settingKey]) {
                        data = item
                    }
                })

                this.$emit('item-selected', this.selectedList[this.settingKey], data)
            },
            getItem (key) {
                let data = null

                this.list.forEach(item => {
                    if (!item.children) {
                        if (String(item[this.settingKey]) === String(key)) {
                            data = item
                        }
                    } else {
                        let children = item.children
                        children.forEach(child => {
                            if (String(child[this.settingKey]) === String(key)) {
                                data = child
                            }
                        })
                    }
                })
                return data
            },
            calcSelected (selected, isTrigger) {
                let data = null

                if (Array.isArray(selected)) {
                    data = []
                    const len = selected.length
                    for (let i = 0; i < len; i++) {
                        let item = this.getItem(selected[i])
                        if (item) {
                            data.push(item)
                        }
                    }

                    if (data.length && isTrigger && !this.initPreventTrigger) {
                        this.$emit('item-selected', selected, data, isTrigger)
                    }
                } else if (selected !== undefined) {
                    let item = this.getItem(selected)
                    if (item) {
                        data = item
                    }
                    if (data && isTrigger && !this.initPreventTrigger) {
                        this.$emit('item-selected', selected, data, isTrigger)
                    }
                }
                return data
            },
            close () {
                this.open = false
            },
            initSelectorPosition (currentTarget) {
                if (currentTarget) {
                    let distanceLeft = getActualLeft(currentTarget)
                    let distanceTop = getActualTop(currentTarget)
                    let winWidth = document.body.clientWidth
                    let winHeight = document.body.clientHeight
                    let xSet = {}
                    let ySet = {}
                    let listHeight = this.list.length * 42
                    if (listHeight > 160) {
                        listHeight = 160
                    }
                    let scrollTop = document.documentElement.scrollTop || document.body.scrollTop

                    if ((distanceTop + listHeight + 42 - scrollTop) < winHeight) {
                        ySet = {
                            top: '40px',
                            bottom: 'auto'
                        }

                        this.listSlideName = 'toggle-slide'
                    } else {
                        ySet = {
                            top: 'auto',
                            bottom: '40px'
                        }

                        this.listSlideName = 'toggle-slide2'
                    }

                    this.panelStyle = {...ySet}
                }
            },
            openFn (event) {
                if (this.disabled) {
                    return
                }

                if (!this.disabled) {
                    if (!this.open && event) {
                        this.initSelectorPosition(event.currentTarget)
                    }
                    this.open = !this.open
                    this.$emit('visible-toggle', this.open)
                }
            },

            /**
             *  计算返回渲染的数组
             */
            calcList () {
                if (this.searchable) {
                    const arr = []
                    const key = this.searchKey

                    const len = this.list.length
                    for (let i = 0; i < len; i++) {
                        const item = this.list[i]
                        if (item.children) {
                            const results = []
                            const childLen = item.children.length
                            for (let j = 0; j < childLen; j++) {
                                const child = item.children[j]
                                if (child[key].toLowerCase().includes(this.condition.toLowerCase())) {
                                    results.push(child)
                                }
                            }
                            if (results.length) {
                                const cloneItem = Object.assign({}, item)
                                cloneItem.children = results
                                arr.push(cloneItem)
                            }
                        } else {
                            if (item[key].toLowerCase().includes(this.condition.toLowerCase())) {
                                arr.push(item)
                            }
                        }
                    }

                    this.localList = arr
                } else {
                    this.localList = this.list
                }
            },

            /**
             *  是否显示清除当前选择的icon
             */
            showClearFn () {
                if (this.allowClear && !this.multiSelect && this.localSelected !== -1 && this.localSelected !== '') {
                    this.showClear = true
                }
            },

            /**
             *  清除选择
             */
            clearSelected (e) {
                this.$emit('clear', this.localSelected)
                this.localSelected = -1
                this.showClear = false
                e.stopPropagation()
                this.$emit('update:selected', '')
            },

            /**
             *  选中列表中的项
             */
            selectItem (data, event) {
                if (data.isDisabled) return
                setTimeout(() => {
                    this.toggleSelect(data, event)
                }, 10)
            },

            toggleSelect (data, event) {
                // label嵌input，触发两次click
                let item
                let $selected
                let $index
                let $selectedList
                let settingKey = this.settingKey
                let isMultiSelect = this.multiSelect
                let list = this.localList
                let index = (data && data[settingKey] !== undefined) ? data[settingKey] : undefined

                if (isMultiSelect && event.target.tagName.toLowerCase() === 'label') {
                    return
                }
                if (index !== undefined) {
                    if (!isMultiSelect) {
                        $selected = index
                    } else {
                        $selected = this.localSelected
                    }

                    item = data
                    this.$emit('update:selected', $selected)
                    $selectedList = this.calcSelected($selected)
                } else {
                    item = null
                    this.$emit('update:selected', -1)
                }

                // 单选时，返回的两个参数是选中项的id（或通过settingKey配置的值）和选中项的数据
                // 多选时，返回的是选中项的索引数组和选中项的数据数组

                this.$emit('item-selected', $selected, $selectedList)

                if (!isMultiSelect) {
                    this.openFn()
                }

                // 点击搜索出来后的列表，不应该把搜索条件清空，清空条件会导致 calcList 方法里计算 localList 的时候计算成所有的
                // this.condition = ''
            },
            editFn (index) {
                this.$emit('edit', index)
                this.openFn()
            },
            delFn (index) {
                this.$emit('del', index)
                this.openFn()
            },
            createFn (e) {
                this.$emit('create')
                this.openFn()
                e.stopPropagation()
            },
            inputFn () {
                this.$emit('typing', this.condition)
            }
        }
    }
</script>

<style lang="scss">
    @import '../../bk-magic-ui/src/selector.scss';
    @import '../../bk-magic-ui/src/loading.scss';
</style>
