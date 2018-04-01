<template>
    <div class="bk-selector"
        :class="[extCls, {'open': open}]"
        @click="openFn"
        v-clickoutside="close">
        <div class="bk-selector-wrapper">
            <input class="bk-selector-input" readonly="readonly"
                :class="{placeholder: selectedText === placeholder, active: open}"
                :value="selectedText"
                :placeholder="placeholder"
                :disabled="disabled"
                @mouseover="showClearFn"
                @mouseleave="showClear = false">
            <i class="bk-icon icon-angle-down bk-selector-icon" v-show="!isLoading && !showClear"></i>
            <i class="bk-icon icon-close bk-selector-icon clear-icon"
                v-show="!isLoading && showClear"
                @mouseover="showClearFn"
                @mouseleave="showClear = false"
                @click="clearSelected($event)">
            </i>
            <div class="bk-spin-loading bk-spin-loading-mini bk-spin-loading-primary" v-show="isLoading">
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
        <transition name="toggle-slide">
            <div class="bk-selector-list" v-show="!isLoading && open">
                <!-- 搜索栏 -->
                <div class="bk-selector-search-item"
                    @click="$event.stopPropagation()"
                    v-if="searchable">
                    <i class="bk-icon icon-search"></i>
                    <input type="text" v-model="condition" @input="inputFn">
                </div>
                <ul>
                    <li class="bk-selector-list-item"
                        v-if="localList.length !== 0"
                        v-for="(item, index) in localList">
                        <template v-if="item.children && item.children.length">
                            <div class="bk-selector-group-name">{{item[displayKey]}}</div>
                            <ul class="bk-selector-group-list">
                                <li v-for="(child, index) in item.children" class="bk-selector-list-item">
                                    <div class="bk-selector-node bk-selector-sub-node"
                                        :class="{'bk-selector-selected': !multiSelect && child[settingKey] === selected}">
                                        <div class="text" @click.stop="selectItem(child, index, $event)">
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
                            <div class="bk-selector-node" :class="{'bk-selector-selected': !multiSelect && item[settingKey] === selected}">
                                <div class="text" @click.stop="selectItem(item, index,  $event)">
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
                            {{ emptyText }}
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
    import clickoutside from './../../../utils/clickoutside'

    /**
     *  bk-dropdown
     *  @module components/dropdown
     *  @desc 下拉选框组件，模拟原生select
     *  @param extCls {String} - 自定义的样式
     *  @param hasCreateItem {Boolean} - 下拉菜单中是否有新增项，默认为true
     *  @param createText {String} - 下拉菜单中新增项的文字
     *  @param tools {Object, Boolean} - 待选项右侧的工具按钮，有两个可配置的key：edit和del，默认为两者都不显示。
     *  @param list {Array} - 必选，下拉菜单所需的数据列表
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
        <bk-dropdown
            :list="list"
            :tools="tools"
            :selected.sync="selected"
            :placeholder="placeholder"
            :displayKey="displayKey"
            :has-create-item="hasCreateItem"
            :create-text="createText"
            :ext-cls="extCls"></bk-dropdown>
    */
    export default {
        name: 'bk-selector',
        props: {
            extCls: {
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
            createText: {
                type: String,
                default: '新增数据源'
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
            selected: {
                type: [Number, Array, String],
                required: true
            },
            placeholder: {
                type: [String, Boolean],
                default: '请选择'
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
                type: [String, Boolean],
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
            }
        },
        data () {
            return {
                open: false,
                selectedList: this.calcSelected(this.selected),
                condition: '',
                localList: this.list,
                localSelected: this.selected,
                emptyText: this.list.length ? '无匹配数据' : '暂无数据',
                showClear: false
            }
        },
        watch: {
            selected (newVal) {
                // 重新生成选择列表
                if (this.list.length) {
                    this.selectedList = this.calcSelected(this.selected, this.isLink)
                }
            },
            list (newVal) {
                // 重新生成选择列表
                this.localList = this.list
                if (this.selected) {
                    this.selectedList = this.calcSelected(this.selected, this.isLink)
                }
            },
            condition () {
                this.calcList()
            },
            localSelected (val) {
                // 重新生成选择列表
                if (this.list.length) {
                    this.selectedList = this.calcSelected(this.localSelected, this.isLink)
                }
            }
        },
        directives: {
            clickoutside
        },
        computed: {
            currentItem () {
                return this.list[this.localSelected]
            },
            selectedText () {
                let text = this.placeholder
                let textArr = []
                if (Array.isArray(this.selectedList) && this.selectedList.length) {
                    this.selectedList.forEach((item) => {
                        textArr.push(item[this.displayKey])
                    })
                } else if (this.selectedList) {
                    textArr.push(this.selectedList[this.displayKey])
                }
                return textArr.length ? textArr.join(',') : this.placeholder
            }
        },
        methods: {
            getItem (key) {
                let data = {}
                let list = this.list

                list.forEach((item, index) => {
                    if (item.children) {
                        let list = item.children
                        list.forEach((item, index) => {
                            if (item[this.settingKey] === key) {
                                data.item = item
                                data.index = index
                            }
                        })
                    } else {
                        if (item[this.settingKey] === key) {
                            data.item = item
                            data.index = index
                        }
                    }
                })
                return data
            },
            calcSelected (selected, isTrigger) {
                let list = this.list
                let displayKey = this.displayKey
                let data = null
                let dataIndex = null

                if (Array.isArray(selected)) {
                    data = []
                    dataIndex = []
                    for (let key of selected) {
                        let params = this.getItem(key)
                        if (params.item) {
                            data.push(params.item)
                            dataIndex.push(params.index)
                        }
                    }
                    if (data.length && isTrigger) {
                        this.$emit('item-selected', selected, data, dataIndex)
                    }
                } else if (selected !== undefined) {
                    let params = this.getItem(selected)
                    if (params.item) {
                        data = params.item
                        dataIndex = params.index
                    }
                    if (data && isTrigger) {
                        this.$emit('item-selected', selected, data, dataIndex)
                    }
                }
                return data
            },
            close () {
                this.open = false
                this.$emit('visible-toggle', this.open)
            },
            openFn () {
                if (!this.disabled) {
                    this.open = !this.open
                    this.$emit('visible-toggle', this.open)
                }
            },
            /**
             *  计算返回渲染的数组
             */
            calcList () {
                if (this.searchable) {
                    let arr = []
                    let key = this.searchKey

                    for (let item of this.list) {
                        if (item.children) {
                            let results = []
                            for (let child of item.children) {
                                if (child[key].includes(this.condition)) {
                                    results.push(child)
                                }
                            }
                            if (results.length) {
                                let cloneItem = Object.assign({}, item)
                                cloneItem.children = results
                                arr.push(cloneItem)
                            }
                        } else {
                            if (item[key].includes(this.condition)) {
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
                if (this.allowClear && !this.multiSelect && this.localSelected !== -1) {
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
            },
            /**
             *  选中列表中的项
             */
            selectItem (data, $index, event) {
                // label嵌input，触发两次click
                if (event.target.tagName.toLowerCase() === 'label') {
                    return
                }
                let item
                let $selected
                let $selectedList
                let settingKey = this.settingKey
                let isMultiSelect = this.multiSelect
                let list = this.localList
                let index = (data && data[settingKey] !== undefined) ? data[settingKey] : undefined
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

                this.$emit('item-selected', $selected, $selectedList, $index)

                if (!isMultiSelect) {
                    this.openFn()
                }

                this.condition = ''
            },
            editFn (e, index) {
                this.$emit('edit', index)
                this.openFn()
                e.stopPropagation()
            },
            delFn (e, index) {
                this.$emit('del', index)
                this.openFn()
                e.stopPropagation()
            },
            createFn (e) {
                this.$emit('create')
                this.openFn()
                e.stopPropagation()
            },
            inputFn () {
                this.$emit('typing', this.condition)
            }
        },
        mounted () {
            this.popup = this.$el
        }
    }
</script>

<style lang="css">
    .bk-selector-loading {
        padding: 10px;
        text-align: center;
    }
    .bk-spin-loading {
        position: absolute;
        right: 10px;
        top: 10px;
    }
    .bk-selector .bk-selector-node.bk-selector-selected {
        background-color: #eef6fe;
        color: #3c96ff;
    }
</style>
